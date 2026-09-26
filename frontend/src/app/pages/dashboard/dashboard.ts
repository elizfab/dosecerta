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
import { RouterLink } from '@angular/router';
import { Observable } from 'rxjs';
import { Chart } from 'chart.js/auto';
import type { ChartConfiguration, ChartDataset, ChartOptions } from 'chart.js';
import { MedicationService } from '../../core/services/medication/medication-service';
import { WeightService } from '../../core/services/weight/weight-service';
import { ExamService } from '../../core/services/exam/exam-service';
import { DoseService } from '../../core/services/dose/dose-service';
import { Medication } from '../../shared/types/medication.interface';
import { WeightRecord } from '../../shared/types/weight-record.interface';
import { Exam } from '../../shared/types/exam.interface';
import { PERIODS, PERIOD_ORDER, PeriodKey, periodFromSchedule } from '../../shared/utils/period';

// Cores (tons fortes da paleta pastel — legíveis em gráfico).
const C_BLUE = '#2e79ad';
const C_GREEN = '#1f8a68';
const C_LAV = '#6b4fb0';
const C_PINK = '#b5487d';
const C_REF = '#c23a3a';

interface PeriodGroup {
  key: PeriodKey;
  label: string;
  icon: string;
  className: string;
  meds: Medication[];
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.scss',
})
export class Dashboard implements OnInit {
  private readonly medicationService = inject(MedicationService);
  private readonly weightService = inject(WeightService);
  private readonly examService = inject(ExamService);
  private readonly doseService = inject(DoseService);

  readonly medications = signal<Medication[]>([]);
  readonly weights = signal<WeightRecord[]>([]);
  readonly exams = signal<Exam[]>([]);
  readonly selectedExam = signal<string>('');

  // Medicamentos já marcados como tomados HOJE (persistido no backend por data).
  private readonly today = this.todayIso();
  readonly takenIds = signal<Set<string>>(new Set());

  // ---- doses de hoje agrupadas por período do dia ----
  readonly periodGroups = computed<PeriodGroup[]>(() => {
    const active = this.medications().filter((m) => m.active);
    const byPeriod = new Map<PeriodKey, Medication[]>();
    for (const m of active) {
      const p = periodFromSchedule(m.schedule);
      const list = byPeriod.get(p) ?? [];
      list.push(m);
      byPeriod.set(p, list);
    }
    return PERIOD_ORDER.map((key) => {
      const meta = PERIODS[key];
      const meds = (byPeriod.get(key) ?? []).sort((a, b) =>
        (a.schedule ?? '').localeCompare(b.schedule ?? ''),
      );
      return { key, label: meta.label, icon: meta.icon, className: meta.className, meds };
    }).filter((g) => g.meds.length > 0);
  });

  readonly activeMedsCount = computed(() => this.medications().filter((m) => m.active).length);
  readonly takenCount = computed(() => {
    const active = new Set(this.medications().filter((m) => m.active).map((m) => m.id));
    return [...this.takenIds()].filter((id) => active.has(id)).length;
  });
  readonly pendingCount = computed(() => this.activeMedsCount() - this.takenCount());
  readonly allDone = computed(() => this.activeMedsCount() > 0 && this.pendingCount() === 0);

  readonly lastWeight = computed<WeightRecord | null>(() => this.weights()[0] ?? null);
  readonly examNames = computed(() => [...new Set(this.exams().map((e) => e.name))]);

  // Referências aos <canvas> (signal queries do Angular).
  private readonly weightCanvas = viewChild<ElementRef<HTMLCanvasElement>>('weightCanvas');
  private readonly measuresCanvas = viewChild<ElementRef<HTMLCanvasElement>>('measuresCanvas');
  private readonly examCanvas = viewChild<ElementRef<HTMLCanvasElement>>('examCanvas');

  private charts: Record<string, Chart> = {};

  constructor() {
    effect(() => this.renderWeightChart());
    effect(() => this.renderMeasuresChart());
    effect(() => this.renderExamChart());
  }

  ngOnInit(): void {
    this.medicationService.list().subscribe({
      next: (l) => this.medications.set(l),
      error: () => this.medications.set([]),
    });
    this.weightService.list().subscribe({
      next: (l) => this.weights.set(l),
      error: () => this.weights.set([]),
    });
    this.examService.list().subscribe({
      next: (l) => {
        this.exams.set(l);
        if (!this.selectedExam() && l.length) this.selectedExam.set(l[0].name);
      },
      error: () => this.exams.set([]),
    });
    // Doses já tomadas hoje (backend).
    this.doseService.listByDate(this.today).subscribe({
      next: (logs) => this.takenIds.set(new Set(logs.map((d) => d.medicationId))),
      error: () => this.takenIds.set(new Set()),
    });
  }

  // ---- "tomei" (persistido por dia no localStorage) ----
  isTaken(id: string): boolean {
    return this.takenIds().has(id);
  }

  toggleTaken(id: string): void {
    const wasTaken = this.takenIds().has(id);
    // Atualização otimista.
    const next = new Set(this.takenIds());
    if (wasTaken) next.delete(id);
    else next.add(id);
    this.takenIds.set(next);

    const request$: Observable<unknown> = wasTaken
      ? this.doseService.unmark(id, this.today)
      : this.doseService.mark(id, this.today);

    request$.subscribe({
      error: () => {
        // Reverte em caso de falha.
        const revert = new Set(this.takenIds());
        if (wasTaken) revert.add(id);
        else revert.delete(id);
        this.takenIds.set(revert);
      },
    });
  }

  private todayIso(): string {
    const d = new Date();
    const mm = String(d.getMonth() + 1).padStart(2, '0');
    const dd = String(d.getDate()).padStart(2, '0');
    return `${d.getFullYear()}-${mm}-${dd}`;
  }

  onExamChange(event: Event): void {
    this.selectedExam.set((event.target as HTMLSelectElement).value);
  }

  // ---- gráficos ----
  private asc<T extends { date: string }>(items: T[]): T[] {
    return [...items].sort((a, b) => a.date.localeCompare(b.date));
  }

  private draw(key: string, canvas: HTMLCanvasElement, config: ChartConfiguration<'line'>): void {
    this.charts[key]?.destroy();
    this.charts[key] = new Chart(canvas, config);
  }

  private renderWeightChart(): void {
    const el = this.weightCanvas()?.nativeElement;
    const data = this.asc(this.weights());
    if (!el || data.length === 0) return;
    this.draw('weight', el, {
      type: 'line',
      data: {
        labels: data.map((w) => this.short(w.date)),
        datasets: [
          {
            label: 'Peso (kg)',
            data: data.map((w) => w.weightKg),
            borderColor: C_GREEN,
            backgroundColor: C_GREEN,
            tension: 0.3,
          },
        ],
      },
      options: this.baseOptions(),
    });
  }

  private renderMeasuresChart(): void {
    const el = this.measuresCanvas()?.nativeElement;
    const data = this.asc(this.weights());
    if (!el || data.length === 0) return;
    const labels = data.map((w) => this.short(w.date));
    this.draw('measures', el, {
      type: 'line',
      data: {
        labels,
        datasets: [
          { label: 'Cintura', data: data.map((w) => w.waist ?? null), borderColor: C_BLUE, backgroundColor: C_BLUE, tension: 0.3, spanGaps: true },
          { label: 'Quadril', data: data.map((w) => w.hip ?? null), borderColor: C_LAV, backgroundColor: C_LAV, tension: 0.3, spanGaps: true },
          { label: 'Peito', data: data.map((w) => w.chest ?? null), borderColor: C_PINK, backgroundColor: C_PINK, tension: 0.3, spanGaps: true },
        ],
      },
      options: this.baseOptions(),
    });
  }

  private renderExamChart(): void {
    const el = this.examCanvas()?.nativeElement;
    const name = this.selectedExam();
    if (!el || !name) return;
    const data = this.asc(this.exams().filter((e) => e.name === name && e.value != null));
    if (data.length === 0) {
      this.charts['exam']?.destroy();
      return;
    }
    const labels = data.map((e) => this.short(e.date));
    const datasets: ChartDataset<'line'>[] = [
      {
        label: name,
        data: data.map((e) => e.value as number),
        borderColor: C_BLUE,
        backgroundColor: C_BLUE,
        tension: 0.3,
      },
    ];
    const min = data[data.length - 1].referenceMin;
    const max = data[data.length - 1].referenceMax;
    if (min != null) {
      datasets.push({ label: 'Ref. mín.', data: labels.map(() => min), borderColor: C_REF, borderDash: [6, 4], pointRadius: 0 });
    }
    if (max != null) {
      datasets.push({ label: 'Ref. máx.', data: labels.map(() => max), borderColor: C_REF, borderDash: [6, 4], pointRadius: 0 });
    }
    this.draw('exam', el, { type: 'line', data: { labels, datasets }, options: this.baseOptions() });
  }

  private baseOptions(): ChartOptions<'line'> {
    return {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { labels: { boxWidth: 12, font: { size: 11 } } } },
      scales: { x: { ticks: { font: { size: 10 } } }, y: { ticks: { font: { size: 10 } } } },
    };
  }

  private short(iso: string): string {
    const [, m, d] = iso.split('-');
    return `${d}/${m}`;
  }

  formatDate(iso: string): string {
    return iso.split('-').reverse().join('/');
  }
}
