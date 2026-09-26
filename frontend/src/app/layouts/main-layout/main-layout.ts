import { Component, OnInit, inject, signal } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { Footer } from '../../shared/components/footer/footer';
import { ToastComponent } from '../../shared/components/toast/toast';
import { HealthService } from '../../core/services/health/health-service';
import { ThemeService } from '../../core/services/theme/theme-service';

type BackendStatus = 'checking' | 'online' | 'offline';

interface NavItem {
  path: string;
  label: string;
  icon: string; // emoji amigável
}

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive, Footer, ToastComponent],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss',
})
export class MainLayout implements OnInit {
  private readonly healthService = inject(HealthService);
  readonly themeService = inject(ThemeService);

  readonly backendStatus = signal<BackendStatus>('checking');

  // Abas principais — barra inferior no mobile / topo no desktop.
  readonly primaryNav: NavItem[] = [
    { path: '/dashboard', label: 'Início', icon: 'ti-home' },
    { path: '/medications', label: 'Medicações', icon: 'ti-pill' },
    { path: '/weight', label: 'Medidas', icon: 'ti-scale' },
    { path: '/exams', label: 'Exames', icon: 'ti-test-pipe' },
    { path: '/calculator', label: 'Calculadora', icon: 'ti-calculator' },
  ];

  // Secundárias — só no topo (desktop) e no rodapé.
  readonly secondaryNav: NavItem[] = [
    { path: '/how-to-use', label: 'Como Usar', icon: 'ti-info-circle' },
    { path: '/faq', label: 'FAQ', icon: 'ti-help-circle' },
    { path: '/about', label: 'Sobre', icon: 'ti-file-text' },
  ];

  ngOnInit(): void {
    this.healthService.check().subscribe({
      next: () => this.backendStatus.set('online'),
      error: () => this.backendStatus.set('offline'),
    });
  }
}
