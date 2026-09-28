import { Component, inject } from '@angular/core';
import { NotificationService } from '../../../core/services/notification/notification-service';

/**
 * Exibe os toasts emitidos pelo NotificationService.
 * Adicionar <app-toast> no MainLayout para cobertura global.
 */
@Component({
  selector: 'app-toast',
  standalone: true,
  template: `
    <div class="dc-toast-container" aria-live="polite" aria-atomic="false">
      @for (toast of notifications.toasts(); track toast.id) {
        <div
          class="dc-toast"
          [class.dc-toast--error]="toast.type === 'error'"
          [class.dc-toast--success]="toast.type === 'success'"
          [class.dc-toast--info]="toast.type === 'info'"
          role="alert"
        >
          <span class="dc-toast__msg">{{ toast.message }}</span>
          <button
            type="button"
            class="dc-toast__close"
            (click)="notifications.dismiss(toast.id)"
            aria-label="Fechar notificação"
          >
            <i class="ti ti-x" aria-hidden="true"></i>
          </button>
        </div>
      }
    </div>
  `,
  styles: [`
    .dc-toast-container {
      position: fixed;
      bottom: calc(var(--dc-bottomnav-h, 4.25rem) + 1rem);
      left: 50%;
      transform: translateX(-50%);
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      z-index: 9000;
      width: min(92%, 25rem);
      pointer-events: none;
    }

    .dc-toast {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 0.75rem;
      padding: 0.75rem 1rem;
      border-radius: var(--dc-radius, 0.2rem);
      font-size: 0.875rem;
      font-weight: 500;
      pointer-events: all;
      box-shadow: var(--dc-shadow-md);
      background: var(--dc-card-bg);
      color: var(--dc-text-main);
      border-left: 0.25rem solid var(--dc-primary);
      animation: toast-in 0.2s ease;
    }

    .dc-toast--error { border-left-color: var(--dc-status-danger); }
    .dc-toast--success { border-left-color: var(--dc-status-success); }
    .dc-toast--info { border-left-color: var(--dc-status-info); }

    .dc-toast__msg { flex: 1; }

    .dc-toast__close {
      background: none;
      border: none;
      cursor: pointer;
      color: inherit;
      opacity: 0.6;
      padding: 0.125rem;
      line-height: 1;
    }
    .dc-toast__close:hover { opacity: 1; }

    @keyframes toast-in {
      from { opacity: 0; transform: translateY(0.5rem); }
      to   { opacity: 1; transform: translateY(0); }
    }
  `],
})
export class ToastComponent {
  readonly notifications = inject(NotificationService);
}
