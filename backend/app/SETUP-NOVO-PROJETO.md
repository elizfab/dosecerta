# Como usar este template para um novo projeto

Este repositório é um template genérico (Go + MongoDB + Docker Compose). Para cada novo backend
(StudyPanel, Rollout Service, e os próximos projetos grandes), clone este template e siga os passos
abaixo. Você **não precisa mexer no código Go** para clonar — só no `.env`.

## 0. Primeira vez usando este template (só uma vez, aqui neste repo)

Se ainda não rodou, limpe os resquícios do StudyPanel que serviram de referência para montar o
template:

```bash
cd app
bash cleanup-template.sh
```

O script remove as entidades/handlers/use cases específicos do StudyPanel, o binário compilado
antigo e o script de seed, deixando as pastas `entity`, `repository`, `usecase` e `handler` vazias
(prontas para o primeiro recurso do próximo projeto). Ele confere no final se `go build ./...`
continua funcionando e se autodestrói ao terminar.

## 1. Clonar o template para um novo projeto

No GitHub, use o botão **"Use this template"** (ou `git clone` + trocar o remote) para criar o
repositório do novo projeto, ex: `rollout-service-backend`.

## 2. Configurar nome do projeto, banco e portas

Copie o `.env.example` para `.env`:

```bash
cp .env.example .env
```

Edite o `.env` e ajuste **apenas estas variáveis** para o novo projeto:

| Variável                 | O que é                                         | Exemplo para "rollout-service" |
| ------------------------ | ----------------------------------------------- | ------------------------------ |
| `PROJECT_NAME`           | Nome usado nos containers, rede e volume Docker | `rollout-service`              |
| `DB_NAME`                | Nome do banco MongoDB                           | `rollout_service`              |
| `API_PORT`               | Porta da API no host                            | `8090`                         |
| `MONGO_PORT`             | Porta do MongoDB no host                        | `27018`                        |
| `MONGO_EXPRESS_PORT`     | Porta do Mongo Express no host                  | `8091`                         |
| `MONGO_EXPRESS_PASSWORD` | Senha do admin do Mongo Express                 | troque o `changeme123`         |

As demais variáveis (`MONGO_URI`, `SERVER_PORT`, `APP_ENV`) podem ficar como estão — `SERVER_PORT`
é a porta _interna_ do container, ela não conflita entre projetos porque cada um roda isolado na
sua própria rede Docker.

Também atualize o nome do banco no `docker/mongo-init.js` (primeira linha,
`db.getSiblingDB("backend")`) para bater com o `DB_NAME` escolhido.

### Por que trocar as portas?

`API_PORT`, `MONGO_PORT` e `MONGO_EXPRESS_PORT` são as portas expostas na sua máquina (host). Se
você rodar dois projetos deste template ao mesmo tempo (ex: StudyPanel e Rollout Service), cada um
precisa de portas diferentes ou o `docker compose up` de um vai falhar porque a porta já está em
uso pelo outro.

Sugestão de convenção — vá somando um bloco de 10 a cada novo projeto:

| Projeto         | API_PORT | MONGO_PORT | MONGO_EXPRESS_PORT |
| --------------- | -------- | ---------- | ------------------ |
| studypanel      | 8080     | 27017      | 8081               |
| rollout-service | 8090     | 27018      | 8091               |
| próximo projeto | 8100     | 27019      | 8101               |

## 3. Subir o banco e a API

Da raiz do projeto (pasta `app/`):

```bash
make docker-up      # sobe API + MongoDB + Mongo Express
```

Ou direto com docker compose:

```bash
docker compose up -d --build
```

Conferir se subiu:

```bash
curl http://localhost:$API_PORT/health
```

Abrir o Mongo Express (interface visual do banco) em `http://localhost:$MONGO_EXPRESS_PORT`,
logando com `MONGO_EXPRESS_USER` / `MONGO_EXPRESS_PASSWORD` do seu `.env`.

Outros comandos úteis:

```bash
make docker-down     # para os containers
make docker-logs      # acompanha o log da API
```

## 4. Rodar a API localmente (fora do Docker), usando só o Mongo em container

```bash
make docker-up       # garante que o mongo está no ar (a API do compose também sobe, tudo bem)
make run              # roda a API localmente, lendo o .env
```

## 5. Adicionar os recursos do novo projeto

O template sobe vazio (sem entidades). Para cada recurso novo (ex: `Rollout`), siga o passo a
passo do **[PADROES.md](PADROES.md)**, seção "Evolução":

1. Entidade em `internal/domain/entity/`
2. Interface do repositório em `internal/domain/repository/`
3. Use cases em `internal/usecase/`
4. Testes dos use cases
5. Implementação MongoDB em `internal/repository/mongodb/`
6. Handler em `internal/handler/`
7. Testes do handler
8. Registrar a rota em `cmd/server/main.go`
9. Documentar em `CASOS_DE_USO.md` (crie esse arquivo no novo projeto — ele foi removido do
   template por ser específico de cada um)

## Checklist antes de considerar o novo projeto pronto

- [ ] `.env` criado com `PROJECT_NAME`, `DB_NAME` e portas próprias do projeto
- [ ] `docker/mongo-init.js` com o nome do banco correto
- [ ] `README.md` com o título/descrição do projeto (troque "Backend Template" pelo nome real)
- [ ] `go.mod` — não precisa trocar o `module backend`, é só um caminho de import interno
- [ ] `make docker-up` sobe sem erro e `/health` responde
- [ ] Este arquivo (`SETUP-NOVO-PROJETO.md`) pode ser apagado do novo projeto depois de configurado
