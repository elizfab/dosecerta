import { Component, computed, OnDestroy, OnInit, signal } from '@angular/core';
import { NavigationEnd, Router, RouterLink } from '@angular/router';
import { BreadcrumbItem } from '../../types/breadcrumb.interface';
import { filter, Subscription } from 'rxjs';
import { BreadcrumbService } from '../../../core/services/Breadcrumb/breadcrumb-service';

@Component({
  selector: 'app-breadcrumbs',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './breadcrumbs.html',
  styleUrl: './breadcrumbs.scss',
})
export class Breadcrumbs implements OnInit, OnDestroy {
  private readonly baseItems = signal<BreadcrumbItem[]>([]);

  readonly items = computed<BreadcrumbItem[]>(() => {
    const base = this.baseItems();
    const extra = this.breadcrumbService.extra();
    return extra ? [...base, extra] : base;
  });

  readonly isHidden = computed(() => this.breadcrumbService.hidden());

  private sub = new Subscription();

  constructor(
    private router: Router,
    private breadcrumbService: BreadcrumbService,
  ) {}

  ngOnInit(): void {
    this.build(this.router.url);

    this.sub.add(
      this.router.events
        .pipe(filter((e): e is NavigationEnd => e instanceof NavigationEnd))
        .subscribe((e) => this.build(e.urlAfterRedirects)),
    );
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

    private build(rawUrl: string): void {
    const url = rawUrl.split('?')[0];
    const segments = url.split('/').filter(Boolean);

    if (segments.length === 0 || segments[0] === 'dashboard') {
      this.baseItems.set([]);
      return;
    }

    const sectionSlug = segments[0];

    // /estudos-labs
    if (sectionSlug === 'estudos-labs') {
      this.baseItems.set([
        { label: 'Dashboard', href: '/dashboard' },
        { label: 'Estudos e Labs' },
      ]);
      return;
    }
  }
}
