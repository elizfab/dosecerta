// Modelo de medicamento — espelha a entidade Medication do backend (Go).
export interface Medication {
  id: string;
  name: string;
  dosage: string; // ex: "20 mg"
  schedule: string; // "HH:mm"
  notes?: string;
  active: boolean;

  // Campos opcionais adicionais.
  period?: string; // manhã/tarde/noite/madrugada (derivado do horário)
  startDate?: string; // "YYYY-MM-DD"
  endDate?: string; // "YYYY-MM-DD"
  stock?: number; // quantidade em estoque

  createdAt: string;
  updatedAt: string;
}

// Campos aceitos ao criar/editar. O backend gera id e timestamps; no update,
// `active` ausente mantém ativo (envie false para pausar).
export interface CreateMedicationInput {
  name: string;
  dosage: string;
  schedule: string;
  notes?: string;
  period?: string; // derivado do schedule via periodFromSchedule()
  startDate?: string;
  endDate?: string;
  stock?: number;
  active?: boolean;
}
