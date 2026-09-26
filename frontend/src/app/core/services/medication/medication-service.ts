import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../api/api-service';
import {
  Medication,
  CreateMedicationInput,
} from '../../../shared/types/medication.interface';

/**
 * Consome os endpoints /api/v1/medications do backend, reusando o ApiService
 * genérico (que desembrulha o envelope { success, data, error }).
 */
@Injectable({ providedIn: 'root' })
export class MedicationService {
  private readonly api = inject(ApiService);
  private readonly path = '/api/v1/medications';

  list(): Observable<Medication[]> {
    return this.api.get<Medication[]>(this.path);
  }

  getById(id: string): Observable<Medication> {
    return this.api.get<Medication>(`${this.path}/${id}`);
  }

  create(input: CreateMedicationInput): Observable<Medication> {
    return this.api.post<Medication>(this.path, input);
  }

  update(
    id: string,
    input: Partial<CreateMedicationInput>,
  ): Observable<Medication> {
    return this.api.put<Medication>(`${this.path}/${id}`, input);
  }

  remove(id: string): Observable<void> {
    return this.api.delete<void>(`${this.path}/${id}`);
  }
}
