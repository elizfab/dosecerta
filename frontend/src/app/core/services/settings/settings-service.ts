import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../api/api-service';

export interface Setting {
  key: string;
  value: string;
}

/** Consome os endpoints /api/v1/settings (configurações key/value). */
@Injectable({ providedIn: 'root' })
export class SettingsService {
  private readonly api = inject(ApiService);
  private readonly path = '/api/v1/settings';

  /** Lê uma configuração. Retorna erro 404 se a chave não existir. */
  get(key: string): Observable<Setting> {
    return this.api.get<Setting>(`${this.path}/${encodeURIComponent(key)}`);
  }

  /** Cria/atualiza uma configuração. */
  set(key: string, value: string): Observable<Setting> {
    return this.api.put<Setting>(`${this.path}/${encodeURIComponent(key)}`, { value });
  }
}
