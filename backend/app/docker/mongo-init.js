// Script de inicialização do MongoDB (DoseCerta).
// Executado automaticamente na primeira subida do container
// (docker-entrypoint-initdb.d), só roda quando o volume mongo-data
// ainda não existe.
db = db.getSiblingDB("dose-certa-backend");

// Coleções e índices dos recursos do DoseCerta.
db.createCollection("medications");
db.medications.createIndex({ active: 1 });
db.medications.createIndex({ createdAt: -1 });

db.createCollection("weight_records");
db.weight_records.createIndex({ date: -1 });

db.createCollection("exams");
db.exams.createIndex({ date: -1 });
db.exams.createIndex({ name: 1, date: -1 });

// Doses tomadas: um registro por (medicamento, data). Índice único garante idempotência.
db.createCollection("dose_logs");
db.dose_logs.createIndex({ medicationId: 1, date: 1 }, { unique: true });
db.dose_logs.createIndex({ date: 1 });

// Configurações key/value (single-user), ex.: meta de peso. A própria chave é o _id.
db.createCollection("settings");

print("MongoDB (dose-certa-backend) inicializado com sucesso.");
