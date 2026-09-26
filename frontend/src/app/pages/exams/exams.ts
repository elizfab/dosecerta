import { Component, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ExamService } from '../../core/services/exam/exam-service';
import {
  Exam,
  CreateExamInput,
  ExamStatus,
  examStatus,
} from '../../shared/types/exam.interface';

@Component({
  selector: 'app-exams',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './exams.html',
  styleUrl: './exams.scss',
})
export class ExamsPage implements OnInit {
  private readonly service = inject(ExamService);

  readonly exams = signal<Exam[]>([]);
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);
  readonly showForm = signal(false);
  readonly editingId = signal<string | null>(null);

  form: {
    name: string;
    value: number | null;
    unit: string;
    referenceMin: number | null;
    referenceMax: number | null;
    date: string;
    notes: string;
  } = this.emptyForm();

  ngOnInit(): void {
    this.load();
  }

  private today(): string {
    return new Date().toISOString().slice(0, 10);
  }

  private emptyForm() {
    return {
      name: '',
      value: null as number | null,
      unit: '',
      referenceMin: null as number | null,
      referenceMax: null as number | null,
      date: this.today(),
      notes: '',
    };
  }

  load(): void {
    this.loading.set(true);
    this.error.set(null);
    this.service.list().subscribe({
      next: (list) => {
        this.exams.set(list);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Não foi possível carregar os exames.');
        this.loading.set(false);
      },
    });
  }

  openForm(): void {
    this.editingId.set(null);
    this.form = this.emptyForm();
    this.showForm.set(true);
  }

  edit(exam: Exam): void {
    this.editingId.set(exam.id);
    this.form = {
      name: exam.name,
      value: exam.value ?? null,
      unit: exam.unit ?? '',
      referenceMin: exam.referenceMin ?? null,
      referenceMax: exam.referenceMax ?? null,
      date: exam.date,
      notes: exam.notes ?? '',
    };
    this.showForm.set(true);
  }

  cancel(): void {
    this.showForm.set(false);
    this.editingId.set(null);
  }

  save(): void {
    if (!this.form.name.trim() || !this.form.date) {
      this.error.set('Informe ao menos o nome do exame e a data.');
      return;
    }
    const payload: CreateExamInput = { name: this.form.name.trim(), date: this.form.date };
    if (this.form.value != null) payload.value = this.form.value;
    if (this.form.unit.trim()) payload.unit = this.form.unit.trim();
    if (this.form.referenceMin != null) payload.referenceMin = this.form.referenceMin;
    if (this.form.referenceMax != null) payload.referenceMax = this.form.referenceMax;
    if (this.form.notes.trim()) payload.notes = this.form.notes.trim();

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
        this.error.set('Não foi possível salvar o exame.');
        this.saving.set(false);
      },
    });
  }

  remove(exam: Exam): void {
    if (!confirm(`Excluir o exame "${exam.name}" de ${this.formatDate(exam.date)}?`)) return;
    this.service.remove(exam.id).subscribe({
      next: () => this.load(),
      error: () => this.error.set('Não foi possível excluir o exame.'),
    });
  }

  status(exam: Exam): ExamStatus {
    return examStatus(exam);
  }

  statusLabel(exam: Exam): string {
    switch (examStatus(exam)) {
      case 'low':
        return 'Abaixo';
      case 'high':
        return 'Acima';
      case 'ok':
        return 'Normal';
      default:
        return '—';
    }
  }

  // Modificador visual do selo de status (usa .dc-status--*).
  statusVariant(exam: Exam): string {
    switch (examStatus(exam)) {
      case 'ok':
        return 'dc-status--success';
      case 'low':
        return 'dc-status--warning';
      case 'high':
        return 'dc-status--danger';
      default:
        return '';
    }
  }

  // Ícone do selo conforme a situação.
  statusIcon(exam: Exam): string {
    switch (examStatus(exam)) {
      case 'ok':
        return 'ti-circle-check';
      case 'low':
        return 'ti-arrow-down';
      case 'high':
        return 'ti-arrow-up';
      default:
        return 'ti-minus';
    }
  }

  reference(exam: Exam): string {
    if (exam.referenceMin != null && exam.referenceMax != null) {
      return `${exam.referenceMin}–${exam.referenceMax}`;
    }
    if (exam.referenceMin != null) return `≥ ${exam.referenceMin}`;
    if (exam.referenceMax != null) return `≤ ${exam.referenceMax}`;
    return '—';
  }

  formatDate(iso: string): string {
    return iso.split('-').reverse().join('/');
  }
}
