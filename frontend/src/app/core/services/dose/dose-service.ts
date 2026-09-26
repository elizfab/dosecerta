import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../api/api-service';
import { DoseLog } from '../../../shared/types/dose-log.interface';

/** Consome os endpoints /api/v1/doses (doses tomadas por data). */
@Injectable({ providedIn: 'root' })
export class DoseService {
  private readonly api = inject(ApiService);
  private readonly path = '/api/v1/doses';

  /** Lista as doses tomadas numa data ("YYYY-MM-DD"). */
  listByDate(date: string): Observable<DoseLog[]> {
    return this.api.get<DoseLog[]>(`${this.path}?date=${encodeURIComponent(date)}`);
  }

  /** Marca a dose como tomada (idempotente). */
  mark(medicationId: string, date: string): Observable<DoseLog> {
    return this.api.post<DoseLog>(this.path, { medicationId, date });
  }

  /** Desmarca a dose (idempotente). */
  unmark(medicationId: string, date: string): Observable<void> {
    return this.api.delete<void>(
      `${this.path}?medicationId=${encodeURIComponent(medicationId)}&date=${encodeURIComponent(date)}`,
    );
  }
}
