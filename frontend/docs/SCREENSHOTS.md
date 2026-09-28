# Screenshots

## `npm run screenshots`

Script: `scripts/screenshots.mjs` (Puppeteer, mesmo padrão da Suplementos Store e do PDI).

```bash
cd backend/app && docker compose up -d mongo api        # os dados vêm da API
cd frontend && npm run screenshots                      # sobe ng serve :6112, captura, encerra
SHOTS_BASE_URL=http://localhost:6012 npm run screenshots  # reaproveita o npm start
SHOTS_VIEWPORTS=mobile npm run screenshots              # mobile | tablet | desktop
SHOTS_THEMES=light npm run screenshots                  # light | dark (padrão: os dois)
SHOTS_FULLPAGE=0 npm run screenshots                    # só a primeira dobra
SHOTS_API_URL=http://localhost:8090 npm run screenshots # backend em outra porta
```

Saída: `docs/screenshots/<NN>-<tela>-<viewport>[-dark].png`: início, medicações,
nova medicação (bottom sheet), medidas, exames, calculadora (com resultado), como usar,
FAQ, sobre e menu dropdown (só mobile/tablet).

- Verifica `GET /health` e confere se o serviço é o **DoseCerta**. Outros projetos do
  portfólio também usam a porta 8080; se a porta estiver com outra API, o script aborta
  com instrução em vez de capturar telas vazias ou com erro.
- **Dados de demonstração**: sem dados, o app mostra só estados vazios. O script cria
  medicações (manhã/tarde/noite), medidas, exames, meta de peso e duas doses tomadas hoje
  via API, e **remove tudo ao final**, restaurando a meta anterior.
  `SHOTS_SEED=0` usa só os dados existentes; `SHOTS_KEEP_DATA=1` mantém os de demonstração.
- `SHOTS_API_URL` diferente de `http://localhost:8080` (o `apiUrl` do `environment.ts`):
  as chamadas do app são redirecionadas no navegador para essa API (interceptação de
  requisições), sem precisar alterar o environment.
- Tema: aplicado pelo `localStorage` (`dc-theme`) antes do carregamento, num contexto de
  navegador isolado por viewport/tema.
- Página inteira: o viewport é esticado até a altura do documento antes da captura, para
  que a barra inferior e o botão "+" (`position: fixed`) fiquem no rodapé, e não no meio
  da imagem.
- O `ng serve` temporário sobe como grupo de processos próprio e é derrubado com
  `kill(-pid)`, sem deixar o `ng` órfão na porta.
- Chrome: usa o cache do Puppeteer (`~/.cache/puppeteer`) ou do Playwright
  (`~/.cache/ms-playwright`). Se não houver:
  `npx @puppeteer/browsers install chrome@stable` ou `CHROME_PATH=/caminho/chrome`.
