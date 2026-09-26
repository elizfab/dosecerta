// Calculadora de dosagem — lógica pura, portada do protótipo (index.html).
// Não depende de rede nem de estado; fácil de testar.

export type DeviceType = 'ml' | 'ui' | 'gotas';

// Constantes clínicas do protótipo:
const UNITS_PER_ML = 100; // seringa de insulina U-100
const DROPS_PER_ML = 20; // 1 mL ≈ 20 gotas

export interface DoseResult {
  volumeMl: number;
  units?: number; // preenchido quando device = 'ui'
  drops?: number; // preenchido quando device = 'gotas'
}

/**
 * Converte a dose prescrita (mg) em volume (mL) a partir da concentração
 * (mg por mL), e opcionalmente em unidades (UI) ou gotas conforme o dispositivo.
 *
 * @throws Error se concentração ou dose não forem maiores que zero.
 */
export function calculateDose(
  concentrationMgPerMl: number,
  prescribedMg: number,
  device: DeviceType,
): DoseResult {
  if (
    !Number.isFinite(concentrationMgPerMl) ||
    !Number.isFinite(prescribedMg) ||
    concentrationMgPerMl <= 0 ||
    prescribedMg <= 0
  ) {
    throw new Error('Concentração e dose devem ser números maiores que zero.');
  }

  const volumeMl = prescribedMg / concentrationMgPerMl;

  if (device === 'ui') {
    return { volumeMl, units: volumeMl * UNITS_PER_ML };
  }
  if (device === 'gotas') {
    return { volumeMl, drops: volumeMl * DROPS_PER_ML };
  }
  return { volumeMl };
}

/** Texto pronto para exibição, equivalente ao do protótipo. */
export function formatDoseResult(result: DoseResult, device: DeviceType): string {
  const ml = result.volumeMl.toFixed(2);
  if (device === 'ui' && result.units !== undefined) {
    return `Equivale a ${ml} mL, ou seja, ${result.units.toFixed(0)} Unidades (UI) na seringa de insulina.`;
  }
  if (device === 'gotas' && result.drops !== undefined) {
    return `Equivale a ${ml} mL, o que representa aproximadamente ${result.drops.toFixed(1)} gotas.`;
  }
  return `Você deve aspirar exatamente ${ml} mL na seringa padrão.`;
}
