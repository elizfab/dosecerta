import { Component, inject } from '@angular/core';
import { BreadcrumbService } from '../../../core/services/Breadcrumb/breadcrumb-service';

@Component({
  selector: 'app-header',
  imports: [],
  templateUrl: './header.html',
  standalone: true,
  styleUrl: './header.scss',
})
export class Header {
  readonly breadcrumbService = inject(BreadcrumbService);
}
