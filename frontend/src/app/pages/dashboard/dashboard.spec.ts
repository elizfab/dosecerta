// Testa as funções puras que alimentam os cartões/gráficos do dashboard
// (RCQ das medidas e situação do exame frente à faixa de referência).
import { waistHipRatio, WeightRecord } from '../../shared/types/weight-record.interface';
import { examStatus, Exam } from '../../shared/types/exam.interface';

function wr(partial: Partial<WeightRecord>): WeightRecord {
  return { id: '1', date: '2026-07-25', weightKg: 70, createdAt: '', ...partial };
}
function ex(partial: Partial<Exam>): Exam {
  return { id: '1', name: 'Glicose', date: '2026-07-25', createdAt: '', ...partial };
}

describe('waistHipRatio', () => {
  it('calcula cintura/quadril', () => {
    expect(waistHipRatio(wr({ waist: 80, hip: 100 }))).toBeCloseTo(0.8, 5);
  });
  it('retorna null sem cintura ou quadril', () => {
    expect(waistHipRatio(wr({ waist: 80 }))).toBeNull();
  });
});

describe('examStatus', () => {
  it('ok dentro da faixa', () => {
    expect(examStatus(ex({ value: 90, referenceMin: 70, referenceMax: 99 }))).toBe('ok');
  });
  it('high acima do máximo', () => {
    expect(examStatus(ex({ value: 120, referenceMin: 70, referenceMax: 99 }))).toBe('high');
  });
  it('low abaixo do mínimo', () => {
    expect(examStatus(ex({ value: 50, referenceMin: 70, referenceMax: 99 }))).toBe('low');
  });
  it('unknown sem referência', () => {
    expect(examStatus(ex({ value: 90 }))).toBe('unknown');
  });
});
