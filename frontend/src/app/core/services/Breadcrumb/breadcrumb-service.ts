import { Injectable, signal } from '@angular/core';
import { BreadcrumbItem } from '../../../shared/types/breadcrumb.interface';

@Injectable({
  providedIn: 'root',
})
export class BreadcrumbService {
  readonly extra = signal<BreadcrumbItem | null>(null);
  readonly pageTitle = signal<string | null>(null);
  readonly hidden = signal(false);

  set(item: BreadcrumbItem | null): void {
    this.extra.set(item);
  }

  setPageTitle(title: string | null): void {
    this.pageTitle.set(title);
  }

  setHidden(value: boolean): void {
    this.hidden.set(value);
  }

  clear(): void {
    this.extra.set(null);
    this.pageTitle.set(null);
    this.hidden.set(false);
  }
}
