import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environment/environment';


export interface HealthStatus {
  status: string;
  service: string;
}

/**
 * Consome GET /health do backend — a única rota que não usa o envelope
 * `{ success, data, error }`, por isso chama o HttpClient direto (sem
 * passar pelo ApiService).
 */
@Injectable({
  providedIn: 'root',
})
export class HealthService {
  private readonly http = inject(HttpClient);

  check(): Observable<HealthStatus> {
    return this.http.get<HealthStatus>(`${environment.apiUrl}/health`);
  }
}
