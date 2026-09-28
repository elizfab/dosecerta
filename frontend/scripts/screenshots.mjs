/**
 * Captura screenshots de cada tela do DoseCerta (mesmo padrão da Suplementos Store / PDI).
 *
 * Pré-requisito: backend do DoseCerta rodando (medicações, medidas e exames vêm da API).
 * Em backend/app:
 *   docker compose up -d mongo api
 *
 * Uso:
 *   npm run screenshots              -> sobe ng serve na porta 6112, captura e encerra
 *   SHOTS_BASE_URL=http://localhost:6012 npm run screenshots
 *                                    -> usa um servidor já rodando (não sobe outro)
 *   SHOTS_API_URL=http://localhost:8080 -> backend usado no seed e nas capturas; se for
 *                                    diferente do apiUrl do environment.ts, as chamadas do
 *                                    app são redirecionadas para ele no navegador
 *   SHOTS_VIEWPORTS=mobile,desktop   -> restringe os viewports (mobile, tablet, desktop)
 *   SHOTS_THEMES=light,dark          -> temas capturados (padrão: light,dark)
 *   SHOTS_FULLPAGE=0                 -> captura só a primeira dobra (padrão: página inteira)
 *   SHOTS_SEED=0                     -> não cria dados de demonstração (usa o que já existe)
 *   SHOTS_KEEP_DATA=1                -> mantém os dados de demonstração ao final
 *   CHROME_PATH=/caminho/chrome      -> força um binário específico
 *
 * Saída: docs/screenshots/<NN>-<tela>-<viewport>[-dark].png
 *
 * O app é mobile-first e sem dados fica só com "estados vazios"; por isso o script
 * cria medicações, medidas, exames, meta de peso e uma dose tomada via API antes de
 * capturar — e apaga tudo (restaurando a meta anterior) ao terminar.
 */
import { spawn } from 'node:child_process';
import { existsSync, mkdirSync, readdirSync } from 'node:fs';
import { homedir } from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import puppeteer from 'puppeteer-core';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.join(__dirname, '..');
const OUT_DIR = path.join(ROOT, 'docs', 'screenshots');
const PORT = process.env.SHOTS_PORT || '6112';
const BASE = (process.env.SHOTS_BASE_URL || `http://localhost:${PORT}`).replace(/\/$/, '');
const APP_API = 'http://localhost:8080'; // apiUrl de src/app/environment/environment.ts
const API = (process.env.SHOTS_API_URL || APP_API).replace(/\/$/, '');
const FULL_PAGE = process.env.SHOTS_FULLPAGE !== '0';
const SEED = process.env.SHOTS_SEED !== '0';
const KEEP_DATA = process.env.SHOTS_KEEP_DATA === '1';
const SETTLE_MS = 1500; // espera animações de entrada e respostas da API
const THEME_KEY = 'dc-theme'; // mesma chave do ThemeService
const GOAL_SETTING = 'weightGoal'; // mesma chave da página de Medidas

const wait = (ms) => new Promise((r) => setTimeout(r, ms));

/** Data local "YYYY-MM-DD" deslocada em dias (negativo = passado). */
const isoDate = (offsetDays = 0) => {
  const d = new Date();
  d.setDate(d.getDate() + offsetDays);
  const pad = (n) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
};

/** Chama a API e desembrulha o envelope { success, data, error }. */
const api = async (method, url, body) => {
  const res = await fetch(`${API}${url}`, {
    method,
    headers: body ? { 'Content-Type': 'application/json' } : undefined,
    body: body ? JSON.stringify(body) : undefined,
  });
  const json = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(`${method} ${url} → ${res.status}: ${json?.error ?? res.statusText}`);
  return json?.data;
};

// ---------------------------------------------------------------------------
// Dados de demonstração (criados antes, removidos depois)
// ---------------------------------------------------------------------------
const DEMO_MEDICATIONS = [
  { name: 'Levotiroxina', dosage: '50 mcg', schedule: '06:30', notes: 'Em jejum', stock: 28 },
  { name: 'Vitamina D', dosage: '2.000 UI', schedule: '09:00', stock: 60 },
  { name: 'Metformina', dosage: '500 mg', schedule: '13:00', notes: 'Após o almoço', stock: 45 },
  { name: 'Losartana', dosage: '50 mg', schedule: '20:00', stock: 30 },
  { name: 'Melatonina', dosage: '3 mg', schedule: '23:00', stock: 20 },
];

const DEMO_WEIGHTS = [
  { date: isoDate(-42), weightKg: 72.4, waist: 84, hip: 104, bodyFatPct: 31.5 },
  { date: isoDate(-28), weightKg: 71.6, waist: 83, hip: 103, bodyFatPct: 30.9 },
  { date: isoDate(-14), weightKg: 70.9, waist: 81.5, hip: 102.5, bodyFatPct: 30.2 },
  { date: isoDate(-1), weightKg: 70.1, waist: 80, hip: 102, bodyFatPct: 29.6 },
];

const DEMO_EXAMS = [
  { name: 'Glicemia em jejum', value: 92, unit: 'mg/dL', referenceMin: 70, referenceMax: 99, date: isoDate(-60) },
  { name: 'Glicemia em jejum', value: 104, unit: 'mg/dL', referenceMin: 70, referenceMax: 99, date: isoDate(-20) },
  { name: 'TSH', value: 2.1, unit: 'mUI/L', referenceMin: 0.4, referenceMax: 4.5, date: isoDate(-20) },
  { name: 'Vitamina D', value: 24, unit: 'ng/mL', referenceMin: 30, referenceMax: 100, date: isoDate(-20), notes: 'Iniciar suplementação' },
];

/** Cria os dados via API e devolve uma função que desfaz tudo. */
const seedDemoData = async () => {
  const created = { medications: [], weights: [], exams: [], doses: [] };
  let previousGoal = null;

  const cleanup = async () => {
    const ignore = () => {};
    for (const d of created.doses) {
      await api('DELETE', `/api/v1/doses?medicationId=${d.medicationId}&date=${d.date}`).catch(ignore);
    }
    for (const id of created.medications) await api('DELETE', `/api/v1/medications/${id}`).catch(ignore);
    for (const id of created.weights) await api('DELETE', `/api/v1/weight-records/${id}`).catch(ignore);
    for (const id of created.exams) await api('DELETE', `/api/v1/exams/${id}`).catch(ignore);
    await api('PUT', `/api/v1/settings/${GOAL_SETTING}`, { value: previousGoal ?? '' }).catch(ignore);
  };

  try {
    previousGoal = (await api('GET', `/api/v1/settings/${GOAL_SETTING}`).catch(() => null))?.value ?? null;

    for (const m of DEMO_MEDICATIONS) {
      created.medications.push((await api('POST', '/api/v1/medications', m)).id);
    }
    for (const w of DEMO_WEIGHTS) {
      created.weights.push((await api('POST', '/api/v1/weight-records', w)).id);
    }
    for (const e of DEMO_EXAMS) {
      created.exams.push((await api('POST', '/api/v1/exams', e)).id);
    }
    await api('PUT', `/api/v1/settings/${GOAL_SETTING}`, { value: '66' });

    // Duas doses da manhã já tomadas hoje — mostra o progresso no "Hoje".
    for (const medicationId of created.medications.slice(0, 2)) {
      const dose = { medicationId, date: isoDate() };
      await api('POST', '/api/v1/doses', dose);
      created.doses.push(dose);
    }
  } catch (err) {
    await cleanup();
    throw err;
  }

  return cleanup;
};

// ---------------------------------------------------------------------------
// Telas
// ---------------------------------------------------------------------------

/**
 * Telas capturadas. `action` roda depois do goto (ex.: abrir o formulário).
 * `only` restringe a tela a certos viewports (o menu dropdown não existe no desktop).
 */
const SCREENS = [
  { name: '01-inicio', path: '/dashboard' },
  { name: '02-medicacoes', path: '/medications' },
  {
    name: '03-medicacao-nova',
    path: '/medications',
    action: async (page) => {
      await page.click('.dc-fab');
      await wait(500);
    },
    fullPage: false,
  },
  { name: '04-medidas', path: '/weight' },
  { name: '05-exames', path: '/exams' },
  {
    name: '06-calculadora',
    path: '/calculator',
    action: async (page) => {
      await page.type('#calc-conc', '10');
      await page.type('#calc-dose', '2.5');
      await page.click('.calc-submit');
      await wait(400);
    },
  },
  { name: '07-como-usar', path: '/how-to-use' },
  { name: '08-faq', path: '/faq' },
  { name: '09-sobre', path: '/about' },
  {
    name: '10-menu',
    path: '/dashboard',
    only: ['mobile', 'tablet'],
    action: async (page) => {
      await page.click('.dc-menu-toggle');
      await wait(400);
    },
    fullPage: false,
  },
];

const ALL_VIEWPORTS = {
  mobile: { width: 390, height: 844 },
  tablet: { width: 834, height: 1194 },
  desktop: { width: 1440, height: 900 },
};
const VIEWPORTS = (process.env.SHOTS_VIEWPORTS || 'mobile,desktop')
  .split(',')
  .map((k) => k.trim())
  .filter((k) => ALL_VIEWPORTS[k])
  .map((k) => ({ key: k, ...ALL_VIEWPORTS[k] }));
const THEMES = (process.env.SHOTS_THEMES || 'light,dark')
  .split(',')
  .map((t) => t.trim())
  .filter((t) => t === 'light' || t === 'dark');

const findChrome = () => {
  if (process.env.CHROME_PATH && existsSync(process.env.CHROME_PATH)) {
    return process.env.CHROME_PATH;
  }
  const caches = [
    path.join(homedir(), '.cache/puppeteer/chrome'),
    path.join(homedir(), '.cache/ms-playwright'),
  ];
  for (const base of caches) {
    if (!existsSync(base)) continue;
    for (const dir of readdirSync(base).sort().reverse()) {
      for (const rel of [
        'chrome-linux64/chrome',
        'chrome-linux/chrome',
        'chrome-headless-shell-linux64/chrome-headless-shell',
      ]) {
        const bin = path.join(base, dir, rel);
        if (existsSync(bin)) return bin;
      }
    }
  }
  for (const bin of ['/usr/bin/google-chrome', '/usr/bin/chromium', '/usr/bin/chromium-browser']) {
    if (existsSync(bin)) return bin;
  }
  return null;
};

const waitForServer = async (url, timeoutMs = 120000) => {
  const start = Date.now();
  while (Date.now() - start < timeoutMs) {
    try {
      const res = await fetch(url, { method: 'GET' });
      if (res.ok) return true;
    } catch {
      /* ainda subindo */
    }
    await wait(1000);
  }
  return false;
};

/**
 * Confirma que quem responde em API é o backend do DoseCerta. Outros projetos do
 * portfólio usam a mesma porta 8080 — o /health devolve o nome do banco em `service`.
 */
const checkBackend = async () => {
  try {
    const res = await fetch(`${API}/health`);
    const body = await res.json();
    if (!res.ok) return `Backend respondeu ${res.status} em ${API}/health.`;
    if (!String(body?.service ?? '').includes('dose-certa')) {
      return `O backend em ${API} é "${body?.service}", não o DoseCerta. Pare o outro projeto ou use SHOTS_API_URL.`;
    }
    return null;
  } catch {
    return `Backend indisponível em ${API}/health.`;
  }
};

const main = async () => {
  const chrome = findChrome();
  if (!chrome) {
    console.error(
      'Chrome não encontrado. Rode `npx @puppeteer/browsers install chrome@stable` ou defina CHROME_PATH.',
    );
    process.exit(1);
  }

  const backendError = await checkBackend();
  if (backendError) {
    console.error(`${backendError}\nSuba a API antes: docker compose up -d mongo api (em backend/app).`);
    process.exit(1);
  }
  mkdirSync(OUT_DIR, { recursive: true });

  // Sem SHOTS_BASE_URL, sobe um ng serve temporário. `detached` cria um grupo de
  // processos próprio: o kill(-pid) derruba o npx E o ng filho.
  let server = null;
  if (!process.env.SHOTS_BASE_URL) {
    console.log(`Subindo ng serve na porta ${PORT}...`);
    server = spawn('npx', ['ng', 'serve', '--port', PORT], {
      cwd: ROOT,
      stdio: 'ignore',
      detached: true,
    });
  }

  let cleanup = null;
  try {
    if (SEED) {
      console.log('Criando dados de demonstração na API...');
      cleanup = await seedDemoData();
    }

    if (!(await waitForServer(BASE))) {
      throw new Error(`Servidor não respondeu em ${BASE}`);
    }
    console.log(
      `Servidor pronto em ${BASE}. Capturando ${SCREENS.length} telas x ${VIEWPORTS.length} viewport(s) x ${THEMES.length} tema(s)...`,
    );

    const browser = await puppeteer.launch({
      executablePath: chrome,
      headless: true,
      args: ['--no-sandbox', '--disable-dev-shm-usage'],
    });

    for (const theme of THEMES) {
      for (const viewport of VIEWPORTS) {
        // Contexto isolado por viewport/tema: localStorage (tema, histórico da calculadora) zerado
        const context = await browser.createBrowserContext();
        const page = await context.newPage();
        await page.setViewport({ width: viewport.width, height: viewport.height });
        if (API !== APP_API) {
          await page.setRequestInterception(true);
          page.on('request', (req) => {
            const url = req.url();
            if (url.startsWith(APP_API)) req.continue({ url: API + url.slice(APP_API.length) });
            else req.continue();
          });
        }
        await page.evaluateOnNewDocument(
          (key, value) => localStorage.setItem(key, value),
          THEME_KEY,
          theme,
        );

        for (const screen of SCREENS) {
          if (screen.only && !screen.only.includes(viewport.key)) continue;

          await page.goto(`${BASE}${screen.path}`, { waitUntil: 'networkidle0', timeout: 60000 });
          await wait(SETTLE_MS);
          if (screen.action) await screen.action(page);

          const suffix = theme === 'dark' ? '-dark' : '';
          const out = path.join(OUT_DIR, `${screen.name}-${viewport.key}${suffix}.png`);
          if (screen.fullPage ?? FULL_PAGE) {
            // O fullPage do Puppeteer deixa elementos `position: fixed` (barra inferior,
            // FAB) presos na altura da 1ª dobra. Esticar o viewport até a altura do
            // documento os leva para o rodapé real da página.
            const height = await page.evaluate(() => document.documentElement.scrollHeight);
            await page.setViewport({ width: viewport.width, height });
            await wait(300);
            await page.screenshot({ path: out });
            await page.setViewport({ width: viewport.width, height: viewport.height });
          } else {
            await page.screenshot({ path: out });
          }
          console.log(`  ${path.relative(ROOT, out)}`);
        }

        await context.close();
      }
    }

    await browser.close();
    console.log(`\nConcluído: screenshots em ${path.relative(ROOT, OUT_DIR)}`);
  } finally {
    if (cleanup && !KEEP_DATA) {
      console.log('Removendo dados de demonstração...');
      await cleanup();
    }
    if (server?.pid) {
      try {
        process.kill(-server.pid, 'SIGTERM');
      } catch {
        /* já encerrado */
      }
    }
  }
};

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
