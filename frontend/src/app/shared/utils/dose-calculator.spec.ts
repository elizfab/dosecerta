import { calculateDose, formatDoseResult } from './dose-calculator';

describe('calculateDose', () => {
  it('calcula volume em mL (dispositivo padrão)', () => {
    // 2.5 mg / 10 (mg/mL) = 0.25 mL
    const r = calculateDose(10, 2.5, 'ml');
    expect(r.volumeMl).toBeCloseTo(0.25, 5);
    expect(r.units).toBeUndefined();
    expect(r.drops).toBeUndefined();
  });

  it('calcula unidades (UI) na seringa de insulina', () => {
    const r = calculateDose(10, 2.5, 'ui');
    expect(r.volumeMl).toBeCloseTo(0.25, 5);
    expect(r.units).toBeCloseTo(25, 5); // 0.25 * 100
  });

  it('calcula gotas', () => {
    const r = calculateDose(10, 2.5, 'gotas');
    expect(r.volumeMl).toBeCloseTo(0.25, 5);
    expect(r.drops).toBeCloseTo(5, 5); // 0.25 * 20
  });

  it('lança erro para concentração <= 0', () => {
    expect(() => calculateDose(0, 5, 'ml')).toThrow();
  });

  it('lança erro para dose <= 0', () => {
    expect(() => calculateDose(10, 0, 'ml')).toThrow();
  });

  it('lança erro para valores não numéricos', () => {
    expect(() => calculateDose(NaN, 5, 'ml')).toThrow();
  });
});

describe('formatDoseResult', () => {
  it('formata mL padrão', () => {
    expect(formatDoseResult({ volumeMl: 0.25 }, 'ml')).toContain('0.25 mL');
  });

  it('formata UI', () => {
    const txt = formatDoseResult({ volumeMl: 0.25, units: 25 }, 'ui');
    expect(txt).toContain('25 Unidades');
  });

  it('formata gotas', () => {
    const txt = formatDoseResult({ volumeMl: 0.25, drops: 5 }, 'gotas');
    expect(txt).toContain('5.0 gotas');
  });
});
