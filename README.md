# Frontend

Frontend do projeto Caderno Inteligente.

## Stack

- Angular 21
- SCSS
- Vitest (testes unitários)

## URLs

- Produção: `https://caderno-frontend.pages.dev`
- Local: `http://localhost:6003`

## Variáveis de ambiente

As configurações de ambiente estão em `src/environments/`:

| Arquivo | Uso |
| ------- | --- |
| `environment.ts` | Padrão |
| `environment.development.ts` | Desenvolvimento |
| `environment.production.ts` | Produção (Cloudflare Pages) |

## Comandos

```bash
npm install
npm start
npm run build
npm run test
```

## Deploy

O deploy é feito automaticamente pelo GitHub Actions para Cloudflare Pages.

### Secrets necessários

Configure no repositório em **Settings > Secrets and variables > Actions**:

- `CLOUDFLARE_ACCOUNT_ID`
- `CLOUDFLARE_API_TOKEN`
