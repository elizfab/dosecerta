// Medida corporal por data (peso + circunferências opcionais em cm).
// Mantém o nome WeightRecord por compatibilidade com o backend (/weight-records),
// mas representa "Medidas corporais" na interface.
export interface WeightRecord {
  id: string;
  date: string; // "YYYY-MM-DD"
  weightKg: number;

  // Circunferências (cm) — opcionais.
  waist?: number; // cintura/barriga
  hip?: number; // quadril
  glute?: number; // glúteo
  chest?: number; // peito
  armRight?: number;
  armLeft?: number;
  forearmRight?: number;
  forearmLeft?: number;
  thighRight?: number;
  thighLeft?: number;
  calfRight?: number;
  calfLeft?: number;
  neck?: number;
  shoulders?: number;
  bodyFatPct?: number; // % de gordura

  createdAt: string;
}

export type CreateWeightInput = Omit<WeightRecord, 'id' | 'createdAt'>;

// Metadado dos campos de medida para gerar o formulário e a tabela dinamicamente.
export interface MeasurementField {
  key: keyof WeightRecord;
  label: string;
  unit: string;
}

export const MEASUREMENT_FIELDS: MeasurementField[] = [
  { key: 'weightKg', label: 'Peso', unit: 'kg' },
  { key: 'waist', label: 'Cintura / barriga', unit: 'cm' },
  { key: 'hip', label: 'Quadril', unit: 'cm' },
  { key: 'glute', label: 'Glúteo', unit: 'cm' },
  { key: 'chest', label: 'Peito', unit: 'cm' },
  { key: 'armRight', label: 'Braço direito', unit: 'cm' },
  { key: 'armLeft', label: 'Braço esquerdo', unit: 'cm' },
  { key: 'forearmRight', label: 'Antebraço direito', unit: 'cm' },
  { key: 'forearmLeft', label: 'Antebraço esquerdo', unit: 'cm' },
  { key: 'thighRight', label: 'Coxa direita', unit: 'cm' },
  { key: 'thighLeft', label: 'Coxa esquerda', unit: 'cm' },
  { key: 'calfRight', label: 'Panturrilha direita', unit: 'cm' },
  { key: 'calfLeft', label: 'Panturrilha esquerda', unit: 'cm' },
  { key: 'neck', label: 'Pescoço', unit: 'cm' },
  { key: 'shoulders', label: 'Ombros', unit: 'cm' },
  { key: 'bodyFatPct', label: '% de gordura', unit: '%' },
];

// Relação cintura/quadril (RCQ). Referência de risco: mulheres > 0.85, homens > 0.90.
export function waistHipRatio(r: WeightRecord): number | null {
  if (r.waist && r.hip && r.hip > 0) {
    return r.waist / r.hip;
  }
  return null;
}
