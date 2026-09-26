import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../api/api-service';
import {
  WeightRecord,
  CreateWeightInput,
} from '../../../shared/types/weight-record.interface';

/**
 * Consome os endpoints /api/v1/weight-records do backend.
 * A listagem já vem ordenada por data desc (regra no backend).
 */
@Injectable({ providedIn: 'root' })
export class WeightService {
  private readonly api = inject(ApiService);
  private readonly path = '/api/v1/weight-records';

  list(): Observable<WeightRecord[]> {
    return this.api.get<WeightRecord[]>(this.path);
  }

  getById(id: string): Observable<WeightRecord> {
    return this.api.get<WeightRecord>(`${this.path}/${id}`);
  }

  create(input: CreateWeightInput): Observable<WeightRecord> {
    return this.api.post<WeightRecord>(this.path, input);
  }

  update(
    id: string,
    input: Partial<CreateWeightInput>,
  ): Observable<WeightRecord> {
    return this.api.put<WeightRecord>(`${this.path}/${id}`, input);
  }

  remove(id: string): Observable<void> {
    return this.api.delete<void>(`${this.path}/${id}`);
  }
}
