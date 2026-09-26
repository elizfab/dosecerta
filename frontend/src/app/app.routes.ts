import { Routes } from '@angular/router';
import { MainLayout } from './layouts/main-layout/main-layout';

export const routes: Routes = [
  {
    path: '',
    component: MainLayout,
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      {
        path: 'dashboard',
        loadComponent: () =>
          import('./pages/dashboard/dashboard').then((m) => m.Dashboard),
      },
      {
        path: 'medications',
        loadComponent: () =>
          import('./pages/medications/medications').then((m) => m.MedicationsPage),
      },
      {
        path: 'weight',
        loadComponent: () =>
          import('./pages/weight/weight').then((m) => m.WeightPage),
      },
      {
        path: 'exams',
        loadComponent: () =>
          import('./pages/exams/exams').then((m) => m.ExamsPage),
      },
      {
        path: 'calculator',
        loadComponent: () =>
          import('./pages/calculator/calculator').then((m) => m.CalculatorPage),
      },
      {
        path: 'how-to-use',
        loadComponent: () =>
          import('./pages/content/how-to-use/how-to-use').then((m) => m.HowToUsePage),
      },
      {
        path: 'faq',
        loadComponent: () =>
          import('./pages/content/faq/faq').then((m) => m.FaqPage),
      },
      {
        path: 'about',
        loadComponent: () =>
          import('./pages/content/about/about').then((m) => m.AboutPage),
      },
    ],
  },
  { path: '**', redirectTo: '' },
];
