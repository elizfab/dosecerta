# Casos de Uso — DoseCerta Backend

Recursos implementados sobre o template Go + MongoDB. Todos sob `/api/v1`, com o
envelope `{ success, data, error }` (exceto `/health`).

## Medication (`/api/v1/medications`)

| Caso de uso | Método/Rota | Regras |
|---|---|---|
| Cadastrar medicamento | `POST /api/v1/medications` | valida `name`, `dosage`, `schedule`; define `active=true` e timestamps |
| Listar medicamentos | `GET /api/v1/medications` | ordenado por `createdAt` desc |
| Buscar por id | `GET /api/v1/medications/{id}` | 404 se não existir; 400 se id inválido |
| Atualizar | `PUT /api/v1/medications/{id}` | revalida campos; atualiza `updatedAt` |
| Remover | `DELETE /api/v1/medications/{id}` | 404 se não existir |

Entidade: `Medication{ id, name, dosage, schedule, notes?, active, createdAt, updatedAt }`.

## WeightRecord (`/api/v1/weight-records`)

| Caso de uso | Método/Rota | Regras |
|---|---|---|
| Registrar peso | `POST /api/v1/weight-records` | valida `date` e `weightKg > 0`; define `createdAt` |
| Listar registros | `GET /api/v1/weight-records` | **ordenado por `date` desc** (mais recente primeiro) |
| Buscar por id | `GET /api/v1/weight-records/{id}` | 404 / 400 |
| Atualizar | `PUT /api/v1/weight-records/{id}` | revalida |
| Remover | `DELETE /api/v1/weight-records/{id}` | 404 se não existir |

Entidade: `WeightRecord{ id, date, weightKg, + circunferências opcionais (waist, hip, glute,
chest, armRight/Left, forearmRight/Left, thighRight/Left, calfRight/Left, neck, shoulders,
bodyFatPct), createdAt }`. Apresentado na UI como "Medidas corporais".

## Exam (`/api/v1/exams`)

| Caso de uso | Método/Rota | Regras |
|---|---|---|
| Cadastrar exame | `POST /api/v1/exams` | valida `name` e `date`; define `createdAt` |
| Listar exames | `GET /api/v1/exams` | ordenado por `date` desc |
| Buscar por id | `GET /api/v1/exams/{id}` | 404 / 400 |
| Atualizar | `PUT /api/v1/exams/{id}` | revalida |
| Remover | `DELETE /api/v1/exams/{id}` | 404 se não existir |

Entidade: `Exam{ id, name, value?, unit?, referenceMin?, referenceMax?, date, notes?, createdAt }`.
A situação frente à faixa de referência é derivada no frontend.

## Health (`/health`)

`GET /health` → `{ "status": "ok", "service": "<DB_NAME>" }` (sem envelope). Usado pelo
frontend para o indicador "Backend conectado".

## Arquitetura

`handler → usecase → repository (interface) → repository/mongodb`. Use cases dependem de
interfaces (`domain/repository`), permitindo testes com mocks sem tocar o banco. Erros de
domínio (`ErrNotFound`, `ErrInvalidID`, validações da entidade) são mapeados para 404/400
em `handler/medication_handler.go` (`writeDomainError`).

## Mapa de arquivos

```
internal/
├── domain/
│   ├── entity/           medication.go, weight_record.go (+ _test.go)
│   └── repository/       medication_repository.go, weight_record_repository.go, errors.go
├── usecase/              medication_usecase.go, weight_usecase.go (+ _test.go)
├── repository/mongodb/   medication_repository.go, weight_repository.go
├── handler/              medication_handler.go, weight_handler.go (+ _test.go)
└── middleware/           cors.go (env-driven via ALLOWED_ORIGINS)
```
