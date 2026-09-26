import {
  Component,
  OnInit,
  inject,
  signal,
  computed,
  effect,
  viewChild,
  ElementRef,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Chart } from 'chart.js/auto';
import type { ChartConfiguration } from 'chart.js';
import { WeightService } from '../../core/services/weight/weight-service';
import { SettingsService } from '../../core/services/settings/settings-service';
import {
  WeightRecord,
  CreateWeightInput,
  MEASUREMENT_FIELDS,
  waistHipRatio,
} from '../../shared/types/weight-record.interface';

const C_GREEN = '#1f8a68';
const C_GOAL = '#b23a68';
const GOAL_SETTING = 'weightGoal';
const PRIMARY_KEYS = ['weightKg', 'waist', 'hip', 'bodyFatPct'];

@Component({
  selector: 'app-weight',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './weight.html',
  styleUrl: './weight.scss',
})
export class WeightPage implements OnInit {
  private readonly service = inject(WeightService);
  private readonly settings = inject(SettingsService);

  readonly fields = MEASUREMENT_FIELDS;
  readonly primaryFields = MEASUREMENT_FIELDS.filter((f) => PRIMARY_KEYS.includes(f.key));
  readonly moreFields = MEASUREMENT_FIELDS.filter((f) => !PRIMARY_KEYS.includes(f.key));

  readonly records = signal<WeightRecord[]>([]);
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);
  readonly showForm = signal(false);
  readonly showMore = signal(false);
  readonly editingId = signal<string | null>(null);

  // Meta de peso (persistida no backend em /api/v1/settings/weightGoal).
  readonly goal = signal<number | null>(null);
  readonly editingGoal = signal(false);
  goalInput: number | null = null;

  readonly last = computed<WeightRecord | null>(() => this.records()[0] ?? null);
  readonly previous = computed<WeightRecord | null>(() => this.records()[1] ?? null);

  // Variação em relação ao registro anterior.
  readonly delta = computed<number | null>(() => {
    const l = this.last();
    const p = this.previous();
    return l && p ? Number((l.weightKg - p.weightKg).toFixed(1)) : null;
  });

  // Distância até a meta.
  readonly toGoal = computed<number | null>(() => {
    const l = this.last();
    const g = this.goal();
    return l && g != null ? Number((l.weightKg - g).toFixed(1)) : null;
  });

  // Formulário: data + mapa de medidas (chave → valor numérico).
  date = this.today();
  values: Record<string, number | null> = {};

  private readonly chartCanvas = viewChild<ElementRef<HTMLCanvasElement>>('weightCanvas');
  private chart?: Chart;

  constructor() {
    effect(() => this.renderChart());
  }

  ngOnInit(): void {
    this.load();
    // Carrega a meta salva (404 = sem meta ainda).
    this.settings.get(GOAL_SETTING).subscribe({
      next: (s) => {
        const v = Number(s.value);
        this.goal.set(s.value !== '' && Number.isFinite(v) ? v : null);
        this.goalInput = this.goal();
      },
      error: () => this.goal.set(null),
    });
  }

  private today(): string {
    return new Date().toISOString().slice(0, 10);
  }

  load(): void {
    this.loading.set(true);
    this.error.set(null);
    this.service.list().subscribe({
      next: (list) => {
        this.records.set(list);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Não foi possível carregar as medidas.');
        this.loading.set(false);
      },
    });
  }

  // ---- gráfico de evolução do peso ----
  private renderChart(): void {
    const el = this.chartCanvas()?.nativeElement;
    const data = [...this.records()].sort((a, b) => a.date.localeCompare(b.date));
    if (!el) return;
    this.chart?.destroy();
    if (data.length === 0) return;

    const labels = data.map((w) => w.date.slice(8) + '/' + w.date.slice(5, 7));
    const g = this.goal();
    const config: ChartConfiguration<'line'> = {
      type: 'line',
      data: {
        labels,
        datasets: [
          {
            label: 'Peso (kg)',
            data: data.map((w) => w.weightKg),
            borderColor: C_GREEN,
            backgroundColor: C_GREEN,
            tension: 0.3,
            fill: false,
          },
          ...(g != null
            ? [
                {
                  label: 'Meta',
                  data: labels.map(() => g),
                  borderColor: C_GOAL,
                  borderDash: [6, 4],
                  pointRadius: 0,
                },
              ]
            : []),
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { labels: { boxWidth: 12, font: { size: 11 } } } },
        scales: { x: { ticks: { font: { size: 10 } } }, y: { ticks: { font: { size: 10 } } } },
      },
    };
    this.chart = new Chart(el, config);
  }

  // ---- meta (backend) ----
  startEditGoal(): void {
    this.goalInput = this.goal();
    this.editingGoal.set(true);
  }

  saveGoal(): void {
    const v = this.goalInput;
    if (v == null || v <= 0) {
      this.clearGoal();
      return;
    }
    this.goal.set(v);
    this.editingGoal.set(false);
    this.settings.set(GOAL_SETTING, String(v)).subscribe({ error: () => {} });
  }

  clearGoal(): void {
    this.goal.set(null);
    this.goalInput = null;
    this.editingGoal.set(false);
    this.settings.set(GOAL_SETTING, '').subscribe({ error: () => {} });
  }

  // ---- formulário (bottom sheet) ----
  openForm(): void {
    this.editingId.set(null);
    this.date = this.today();
    this.values = {};
    this.showMore.set(false);
    this.error.set(null);
    this.showForm.set(true);
  }

  edit(rec: WeightRecord): void {
    this.editingId.set(rec.id);
    this.date = rec.date;
    this.values = {};
    for (const f of this.fields) {
      const v = rec[f.key];
      if (typeof v === 'number') this.values[f.key] = v;
    }
    this.showMore.set(true);
    this.error.set(null);
    this.showForm.set(true);
  }

  cancel(): void {
    this.showForm.set(false);
    this.editingId.set(null);
  }

  save(): void {
    const weight = this.values['weightKg'];
    if (!this.date || weight == null || weight <= 0) {
      this.error.set('Informe a data e um peso maior que zero.');
      return;
    }
    const payload: CreateWeightInput = { date: this.date, weightKg: weight };
    for (const f of this.fields) {
      const v = this.values[f.key];
      if (f.key !== 'weightKg' && v != null && v > 0) {
        (payload as Record<string, number | string>)[f.key] = v;
      }
    }

    this.saving.set(true);
    this.error.set(null);

    const id = this.editingId();
    const request$ = id
      ? this.service.update(id, payload)
      : this.service.create(payload);

    request$.subscribe({
      next: () => {
        this.saving.set(false);
        this.showForm.set(false);
        this.editingId.set(null);
        this.load();
      },
      error: () => {
        this.error.set('Não foi possível salvar a medida.');
        this.saving.set(false);
      },
    });
  }

  remove(rec: WeightRecord): void {
    if (!confirm(`Excluir a medida de ${this.formatDate(rec.date)}?`)) return;
    this.service.remove(rec.id).subscribe({
      next: () => this.load(),
      error: () => this.error.set('Não foi possível excluir o registro.'),
    });
  }

  rcq(rec: WeightRecord): string {
    const r = waistHipRatio(rec);
    return r === null ? '—' : r.toFixed(2);
  }

  formatDate(iso: string): string {
    return iso.split('-').reverse().join('/');
  }
}
