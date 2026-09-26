import { Injectable, signal } from '@angular/core';

export type ToastType = 'error' | 'success' | 'info';

export interface Toast {
  id: number;
  message: string;
  type: ToastType;
}

/**
 * Serviço global de notificações (toasts).
 * Usado pelo apiErrorInterceptor para exibir erros de API de forma padronizada,
 * sem acoplamento entre componentes.
 *
 * Uso:
 *   notificationService.show('Mensagem', 'error');
 *   notificationService.show('Salvo!', 'success');
 */
@Injectable({ providedIn: 'root' })
export class NotificationService {
  readonly toasts = signal<Toast[]>([]);

  private nextId = 0;

  show(message: string, type: ToastType = 'info', durationMs = 4000): void {
    const id = ++this.nextId;
    this.toasts.update((list) => [...list, { id, message, type }]);
    setTimeout(() => this.dismiss(id), durationMs);
  }

  dismiss(id: number): void {
    this.toasts.update((list) => list.filter((t) => t.id !== id));
  }
}
