import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, throwError } from 'rxjs';
import { NotificationService } from '../services/notification/notification-service';

/**
 * Interceptor funcional que padroniza erros de chamadas à API.
 * Extrai a mensagem do envelope `{ success, error }` do backend quando
 * disponível, exibe um toast global via NotificationService e loga no console.
 *
 * Erros de health-check (/health) são silenciados para não poluir a UI —
 * o MainLayout trata o status separadamente.
 */
export const apiErrorInterceptor: HttpInterceptorFn = (req, next) => {
  const notifications = inject(NotificationService);

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      const message: string =
        error.error?.error ?? error.message ?? 'Erro inesperado ao comunicar com a API.';

      console.error(`[API] ${req.method} ${req.url} → ${error.status}: ${message}`);

      // Não exibe toast para o health-check (chamada técnica, sem ação da usuária).
      const isHealthCheck = req.url.endsWith('/health');
      if (!isHealthCheck) {
        notifications.show(message, 'error');
      }

      return throwError(() => new Error(message));
    }),
  );
};
