// Resultado de um exame laboratorial numa data.
export interface Exam {
  id: string;
  name: string;
  value?: number;
  unit?: string;
  referenceMin?: number;
  referenceMax?: number;
  date: string; // "YYYY-MM-DD"
  notes?: string;
  createdAt: string;
  updatedAt?: string;
}

export type CreateExamInput = Omit<Exam, 'id' | 'createdAt' | 'updatedAt'>;

// Situação do resultado frente à faixa de referência.
export type ExamStatus = 'ok' | 'low' | 'high' | 'unknown';

export function examStatus(e: Exam): ExamStatus {
  if (e.value == null || (e.referenceMin == null && e.referenceMax == null)) {
    return 'unknown';
  }
  if (e.referenceMin != null && e.value < e.referenceMin) return 'low';
  if (e.referenceMax != null && e.value > e.referenceMax) return 'high';
  return 'ok';
}
