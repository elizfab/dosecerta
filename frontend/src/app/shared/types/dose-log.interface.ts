// Registro de que uma dose foi tomada numa data (backend /api/v1/doses).
export interface DoseLog {
  id: string;
  medicationId: string;
  date: string; // "YYYY-MM-DD"
  createdAt: string;
}
