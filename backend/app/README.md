# Backend Template — Go + MongoDB + Docker

Template base para APIs REST em Go, seguindo Clean Architecture, com MongoDB e Docker Compose prontos para uso. Use este repositório como ponto de partida para novos projetos (backend do StudyPanel, do Rollout Service, e assim por diante).

## Stack

- **Go 1.22** — net/http com pattern matching nativo (sem framework)
- **MongoDB 7.0** — banco de dados NoSQL
- **Mongo Express** — interface web para administrar o banco
- **Docker Compose** — orquestração local

## Pré-requisitos

- Go 1.22+
- Docker e Docker Compose

## Setup rápido

```bash
cp .env.example .env
make docker-up
```

Acesse:
- **API:** http://localhost:8080/health
- **Mongo Express:** http://localhost:8081 (usuário/senha definidos no `.env`)

> Clonando este template para um novo projeto? Veja **[SETUP-NOVO-PROJETO.md](SETUP-NOVO-PROJETO.md)** — o passo a passo para renomear o projeto, trocar as portas (evitando conflito com outros projetos deste mesmo template) e subir o banco.

## Desenvolvimento local

```bash
make tidy          # instalar dependências
make docker-up      # sobe a stack completa (ou só o mongo, se preferir rodar a API localmente)
make run            # roda a API localmente
```

## Testes

```bash
make test           # todos os testes
make test-cover      # com cobertura (abre coverage.html)
```

## Estrutura

```
app/
├── cmd/server/main.go        → Entrada da aplicação
├── config/                   → Configuração via env vars
├── internal/
│   ├── domain/
│   │   ├── entity/           → Entidades do domínio (vazio — adicione as suas)
│   │   └── repository/       → Interfaces de repositório (vazio — adicione as suas)
│   ├── usecase/               → Casos de uso / regras de negócio (vazio — adicione os seus)
│   ├── handler/                → HTTP handlers (vazio — adicione os seus)
│   ├── repository/mongodb/     → Implementação MongoDB dos repositórios (vazio)
│   └── middleware/             → CORS
├── pkg/response/               → Helpers de resposta HTTP
├── docker/mongo-init.js        → Init do MongoDB
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

Este template já vem com a arquitetura pronta, mas sem recursos de negócio — as pastas `entity`, `repository`, `usecase` e `handler` estão vazias (só com um `.gitkeep`). Veja **[PADROES.md](PADROES.md)** para o padrão de código e o passo a passo de como adicionar um novo recurso.
