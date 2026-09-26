# PADROES — Backend Template (Go + MongoDB + Docker)

Padrões técnicos obrigatórios para qualquer projeto criado a partir deste template.

---

## Tecnologias

- **Linguagem:** Go 1.22+
- **Banco de dados:** MongoDB 7.0
- **ORM/Driver:** go.mongodb.org/mongo-driver v1.15+
- **Containers:** Docker + Docker Compose
- **Interface admin:** Mongo Express

---

## Estrutura de Diretórios

```
backend/
├── cmd/
│   └── server/
│       └── main.go              → Ponto de entrada da aplicação
│
├── config/
│   └── config.go                → Carregamento de variáveis de ambiente
│
├── internal/
│   ├── domain/
│   │   ├── entity/              → Entidades do domínio
│   │   └── repository/          → Interfaces (contratos) de repositório
│   │
│   ├── usecase/                 → Casos de uso (regras de negócio)
│   │
│   ├── handler/                 → Handlers HTTP (controllers)
│   │
│   ├── repository/
│   │   └── mongodb/             → Implementação MongoDB do repositório
│   │
│   └── middleware/              → Middlewares (CORS, auth, logging)
│
├── pkg/
│   └── response/                → Helpers de resposta HTTP
│
├── docker/
│   └── mongo-init.js            → Script de inicialização do MongoDB
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
├── .env.example
├── .gitignore
├── PADROES.md                   → Este arquivo
├── CASOS_DE_USO.md              → Documentação dos casos de uso
└── README.md                    → Setup e execução
```

---

## Padrões de Código

### Nomenclatura

- Pacotes: `lowercase` sem underscores (ex: `usecase`, `handler`, `mongodb`)
- Structs: `PascalCase` (ex: `StudyItem`, `CreateStudyItemUseCase`)
- Funções exportadas: `PascalCase`
- Funções privadas: `camelCase`
- Constantes: `PascalCase` para exportadas (ex: `StatusCompleted`)
- Variáveis: `camelCase`

### Retorno de Erros

Sempre propagar erros com contexto. Nunca engolir erros silenciosamente:

```go
// ✅ Correto
if err != nil {
    return nil, fmt.Errorf("falha ao criar item: %w", err)
}

// ❌ Evitar
if err != nil {
    return nil, err  // sem contexto
}
```

### Context

Sempre passar `context.Context` como primeiro parâmetro em funções que acessam recursos externos:

```go
func (r *Repository) Create(ctx context.Context, item *entity.StudyItem) (*entity.StudyItem, error)
```

### Injeção de Dependência

Use interfaces para desacoplar. Injete dependências via construtor:

```go
type UseCase struct {
    repo repository.StudyItemRepository  // interface, não implementação
}

func NewUseCase(repo repository.StudyItemRepository) *UseCase {
    return &UseCase{repo: repo}
}
```

---

## Padrões de API

### Prefixo de rotas

```
/api/v1/
```

### Verbos HTTP

| Ação | Verbo | Rota |
|------|-------|------|
| Listar | GET | `/api/v1/study-items` |
| Buscar por ID | GET | `/api/v1/study-items/{id}` |
| Criar | POST | `/api/v1/study-items` |
| Atualizar | PUT | `/api/v1/study-items/{id}` |
| Deletar | DELETE | `/api/v1/study-items/{id}` |

### Formato de Resposta

```json
{
  "success": true,
  "data": { ... }
}
```

```json
{
  "success": false,
  "error": "mensagem de erro"
}
```

### Códigos HTTP

| Situação | Código |
|----------|--------|
| Sucesso (listagem/get) | 200 |
| Criado com sucesso | 201 |
| Requisição inválida | 400 |
| Não encontrado | 404 |
| Erro interno | 500 |

---

## Padrões de Teste

### Cobertura Mínima

**Meta:** 90% a 100%

### Regras

- Testes unitários para todos os Use Cases
- Testes unitários para todos os Handlers
- Testes unitários para todas as Entities
- Usar mocks via interface para isolar dependências
- Nenhum teste deve acessar banco de dados real
- Testes de integração (MongoDB) são opcionais e separados em `_integration_test.go`

### Nomenclatura

```go
func Test[UnitBeingTested]_[Scenario](t *testing.T)
```

Exemplos:
```go
func TestCreateStudyItem_Success(t *testing.T)
func TestCreateStudyItem_MissingSection(t *testing.T)
func TestGetStudyItem_NotFound(t *testing.T)
```

### Executar Testes

```bash
make test
make test-cover
```

---

## Variáveis de Ambiente

| Variável | Padrão | Descrição |
|---------|--------|-----------|
| `MONGO_URI` | `mongodb://localhost:27017` | URI de conexão MongoDB |
| `DB_NAME` | `backend` | Nome do banco de dados (mude por projeto — veja SETUP-NOVO-PROJETO.md) |
| `SERVER_PORT` | `8080` | Porta interna do servidor HTTP dentro do container |
| `APP_ENV` | `development` | Ambiente da aplicação |
| `PROJECT_NAME` | `backend` | Nome usado em containers/rede/volume do docker-compose |
| `API_PORT` / `MONGO_PORT` / `MONGO_EXPRESS_PORT` | `8080` / `27017` / `8081` | Portas expostas no host — mude por projeto para evitar conflito |

Nunca commitar o arquivo `.env`. Usar `.env.example` como referência.

---

## Segurança

- Credenciais apenas em variáveis de ambiente
- Nunca expor `.env` no git (já está no `.gitignore`)
- MongoDB não deve ser exposto publicamente em produção
- CORS configurado no middleware (ajustar origens em produção)
- Usar usuário não-root no container Docker

---

## Evolução

Ao adicionar um novo recurso, siga obrigatoriamente:

1. Criar entidade em `domain/entity/`
2. Criar interface em `domain/repository/`
3. Criar use cases em `usecase/`
4. Criar testes para os use cases
5. Implementar repositório em `repository/mongodb/`
6. Criar handler em `handler/`
7. Criar testes para o handler
8. Registrar rota no handler
9. Atualizar `CASOS_DE_USO.md`
