// Período do dia — organizador visual principal do app (ver docs/design).
// Derivado do horário "HH:mm" do medicamento; as cores vêm dos tokens --dc-period-*.

export type PeriodKey = 'morning' | 'afternoon' | 'evening' | 'night';

export interface PeriodMeta {
  key: PeriodKey;
  label: string;
  icon: string; // classe do ícone Tabler (ti-*)
  className: string; // classe do bloco: dc-period--*
  order: number; // ordem de exibição
}

export const PERIODS: Record<PeriodKey, PeriodMeta> = {
  morning: { key: 'morning', label: 'Manhã', icon: 'ti-sunrise', className: 'dc-period--morning', order: 1 },
  afternoon: { key: 'afternoon', label: 'Tarde', icon: 'ti-sun', className: 'dc-period--afternoon', order: 2 },
  evening: { key: 'evening', label: 'Noite', icon: 'ti-sunset', className: 'dc-period--evening', order: 3 },
  night: { key: 'night', label: 'Madrugada', icon: 'ti-moon', className: 'dc-period--night', order: 4 },
};

export const PERIOD_ORDER: PeriodKey[] = ['morning', 'afternoon', 'evening', 'night'];

/**
 * Deriva o período a partir de um horário "HH:mm".
 * Manhã 05–11h · Tarde 12–17h · Noite 18–23h · Madrugada 00–04h.
 * Entrada inválida cai em "morning".
 */
export function periodFromSchedule(schedule: string | undefined | null): PeriodKey {
  const h = Number((schedule ?? '').split(':')[0]);
  if (!Number.isFinite(h)) return 'morning';
  if (h >= 5 && h < 12) return 'morning';
  if (h >= 12 && h < 18) return 'afternoon';
  if (h >= 18 && h < 24) return 'evening';
  return 'night';
}
