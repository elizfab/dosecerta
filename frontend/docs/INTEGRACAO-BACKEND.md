# Integração com o Backend

Este frontend consome o [template de backend Go + MongoDB + Docker](../../template-backend-go-mongo-docker)
via HTTP. Este documento explica como a conexão foi montada e como ajustá-la ao clonar este
template para um novo projeto.

## Onde fica a configuração

A URL da API vive em `src/environments/`:

- `environment.ts` — usado por `ng serve` / `npm start` (dev). Aponta para `http://localhost:8080`.
- `environment.prod.ts` — usado por `ng build` (produção). O `angular.json` troca automaticamente
  um pelo outro via `fileReplacements` na configuration `production`.

```ts
// src/environments/environment.ts
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080',
};
```

Não existe chamada HTTP no código com URL "hardcoded" — todo serviço lê `environment.apiUrl`.

## Como as portas precisam bater com o backend

| Onde                  | Porta padrão | Arquivo                                 |
| --------------------- | ------------ | --------------------------------------- |
| Frontend (dev server) | `6012`       | `package.json` (`ng serve --port 6012`) |
| Backend (API)         | `8080`       | `.env` do backend (`API_PORT`)          |

Se você mudar `API_PORT` no `.env` do backend (por exemplo, ao rodar vários projetos deste template
ao mesmo tempo — veja `SETUP-NOVO-PROJETO.md` no repositório do backend), atualize o mesmo valor em
`environment.apiUrl` aqui no frontend (dev e prod).

## CORS

O backend já libera CORS para qualquer origem (`Access-Control-Allow-Origin: *`, em
`internal/middleware/cors.go`), então não é preciso nenhuma configuração extra no frontend para
rodar local. Antes de ir para produção, restrinja as origens permitidas no backend para o domínio
real do frontend.

## Camadas criadas nesta integração

- **`src/app/shared/types/api-response.interface.ts`** — espelha o envelope `{ success, data, error }`
  que todo endpoint em `/api/v1/*` do backend retorna.
- **`src/app/core/services/api/api-service.ts`** (`ApiService`) — wrapper genérico sobre o
  `HttpClient` com `get/post/put/delete`, que já desembrulha o envelope acima. Use como base para o
  serviço de cada recurso.
- **`src/app/core/services/health/health-service.ts`** (`HealthService`) — exemplo real, consome
  `GET /health`. Esse endpoint é a única exceção no backend que **não** usa o envelope
  `{ success, data, error }` (retorna `{ status, service }` puro), por isso chama o `HttpClient`
  direto em vez de passar pelo `ApiService`.
- **`src/app/core/interceptors/api-error-interceptor.ts`** — interceptor funcional registrado em
  `app.config.ts` (`provideHttpClient(withFetch(), withInterceptors([apiErrorInterceptor]))`) que
  padroniza e loga erros de qualquer chamada à API.
- **`src/app/pages/dashboard/dashboard.ts`** — exemplo de uso ponta a ponta: injeta o
  `HealthService`, chama `check()` no `ngOnInit`, e mostra um indicador visual
  ("Backend conectado" / "Backend indisponível") no topo do dashboard.

## Como consumir um novo recurso (ex: `Product`)

Quando o backend ganhar um recurso novo em `/api/v1/products` (seguindo o `PADROES.md` do backend),
crie o espelho no frontend:

1. Crie a interface do modelo em `src/app/shared/types/product.interface.ts`.
2. Crie `src/app/core/services/product/product-service.ts`, injetando `ApiService`:

   ```ts
   @Injectable({ providedIn: 'root' })
   export class ProductService {
     private readonly api = inject(ApiService);

     list(): Observable<Product[]> {
       return this.api.get<Product[]>('/api/v1/products');
     }

     create(input: CreateProductInput): Observable<Product> {
       return this.api.post<Product>('/api/v1/products', input);
     }
   }
   ```

3. Use o serviço no componente da feature, com `inject()` + `signal`, seguindo o mesmo padrão do
   `Dashboard` (veja `dashboard.ts`).

## Rodando os dois juntos localmente

```bash
# 1. Backend (na pasta do template-backend-go-mongo-docker/app)
cp .env.example .env
make docker-up

# 2. Frontend (nesta pasta)
npm install
npm start          # http://localhost:6012, consumindo http://localhost:8080
```

Conferir a conexão: abra `http://localhost:6012`, o dashboard deve mostrar "Backend conectado". Se
mostrar "Backend indisponível", veja Troubleshooting abaixo.

## Troubleshooting

| Sintoma                              | Causa provável                         | Solução                                                                                 |
| ------------------------------------ | -------------------------------------- | --------------------------------------------------------------------------------------- |
| "Backend indisponível" no dashboard  | Backend não está rodando               | Rode `make docker-up` na pasta do backend e confira `curl http://localhost:8080/health` |
| Erro de CORS no console do navegador | Origem do frontend não liberada        | Confira `internal/middleware/cors.go` no backend (por padrão libera tudo)               |
| `ERR_CONNECTION_REFUSED`             | Porta errada entre frontend e backend  | Confira se `environment.apiUrl` bate com `API_PORT` do `.env` do backend                |
| Erro de "Mixed Content" em produção  | Frontend em HTTPS chamando API em HTTP | `environment.prod.ts` deve apontar para uma URL `https://`                              |

## Ao clonar este template para um novo projeto

- Atualize `apiUrl` em `environment.ts` (dev) e `environment.prod.ts` (produção) para a porta/URL
  do backend daquele projeto.
- Não precisa mexer em `ApiService`, `HealthService` nem no interceptor — eles são genéricos.
