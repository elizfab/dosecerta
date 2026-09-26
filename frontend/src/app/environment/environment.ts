// Ambiente de desenvolvimento (usado por `ng serve` / `npm start`).
// Sobrescrito em build de produção pelo fileReplacements do angular.json,
// que troca este arquivo por environment.prod.ts.
export const environment = {
  production: false,
  // URL base da API do backend (template Go + MongoDB + Docker).
  // Deve bater com API_PORT do .env do backend. Veja INTEGRACAO-BACKEND.md.
  apiUrl: 'http://localhost:8080',
};
