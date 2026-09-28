import { Component, DestroyRef, HostListener, inject, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { NavigationEnd, Router, RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { filter } from 'rxjs';
import { Footer } from '../../shared/components/footer/footer';
import { ToastComponent } from '../../shared/components/toast/toast';
import { ThemeService } from '../../core/services/theme/theme-service';

interface NavItem {
  path: string;
  label: string;
  icon: string; // classe do Tabler Icons
}

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive, Footer, ToastComponent],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss',
})
export class MainLayout {
  readonly themeService = inject(ThemeService);

  // Menu dropdown do header (mobile/tablet).
  readonly menuOpen = signal(false);

  // Abas principais — barra inferior no mobile / topo no desktop.
  readonly primaryNav: NavItem[] = [
    { path: '/dashboard', label: 'Início', icon: 'ti-home' },
    { path: '/medications', label: 'Medicações', icon: 'ti-pill' },
    { path: '/weight', label: 'Medidas', icon: 'ti-scale' },
    { path: '/exams', label: 'Exames', icon: 'ti-test-pipe' },
    { path: '/calculator', label: 'Calculadora', icon: 'ti-calculator' },
  ];

  // Secundárias — no dropdown (mobile) e no topo (desktop).
  readonly secondaryNav: NavItem[] = [
    { path: '/how-to-use', label: 'Como Usar', icon: 'ti-info-circle' },
    { path: '/faq', label: 'FAQ', icon: 'ti-help-circle' },
    { path: '/about', label: 'Sobre', icon: 'ti-file-text' },
  ];

  constructor() {
    // Fecha o menu sempre que a navegação terminar.
    inject(Router)
      .events.pipe(
        filter((e) => e instanceof NavigationEnd),
        takeUntilDestroyed(inject(DestroyRef)),
      )
      .subscribe(() => this.menuOpen.set(false));
  }

  toggleMenu(): void {
    this.menuOpen.update((open) => !open);
  }

  closeMenu(): void {
    this.menuOpen.set(false);
  }

  @HostListener('document:keydown.escape')
  onEscape(): void {
    this.closeMenu();
  }
}
