import { Injector, runInInjectionContext } from '@angular/core';
import { of } from 'rxjs';
import { MedicationService } from './medication-service';
import { ApiService } from '../api/api-service';
import { Medication } from '../../../shared/types/medication.interface';

// Testa o serviço com um Injector leve (Injector.create), sem TestBed —
// suficiente para resolver o inject(ApiService) do serviço e mais rápido.
describe('MedicationService', () => {
  let service: MedicationService;
  let api: jest.Mocked<Pick<ApiService, 'get' | 'post' | 'put' | 'delete'>>;

  const sample: Medication = {
    id: '1',
    name: 'Omeprazol',
    dosage: '20 mg',
    schedule: '08:00',
    active: true,
    createdAt: '2026-07-27T10:00:00Z',
    updatedAt: '2026-07-27T10:00:00Z',
  };

  beforeEach(() => {
    api = { get: jest.fn(), post: jest.fn(), put: jest.fn(), delete: jest.fn() };
    const injector = Injector.create({
      providers: [{ provide: ApiService, useValue: api }],
    });
    service = runInInjectionContext(injector, () => new MedicationService());
  });

  it('list() chama GET no path correto', (done) => {
    api.get.mockReturnValue(of([sample]));
    service.list().subscribe((res) => {
      expect(api.get).toHaveBeenCalledWith('/api/v1/medications');
      expect(res).toEqual([sample]);
      done();
    });
  });

  it('create() chama POST com o input', (done) => {
    api.post.mockReturnValue(of(sample));
    const input = { name: 'Omeprazol', dosage: '20 mg', schedule: '08:00' };
    service.create(input).subscribe(() => {
      expect(api.post).toHaveBeenCalledWith('/api/v1/medications', input);
      done();
    });
  });

  it('remove() chama DELETE no id', (done) => {
    api.delete.mockReturnValue(of(undefined));
    service.remove('1').subscribe(() => {
      expect(api.delete).toHaveBeenCalledWith('/api/v1/medications/1');
      done();
    });
  });
});
