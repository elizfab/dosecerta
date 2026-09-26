// Espelha o envelope de resposta do backend (pkg/response/response.go).
// Endpoints em /api/v1/* sempre respondem nesse formato.
// O endpoint /health é a exceção — responde um objeto simples, sem envelope.
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}
