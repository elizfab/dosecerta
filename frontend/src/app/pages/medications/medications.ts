import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MedicationService } from '../../core/services/medication/medication-service';
import {
  Medication,
  CreateMedicationInput,
} from '../../shared/types/medication.interface';
import { PERIODS, PERIOD_ORDER, PeriodKey, periodFromSchedule } from '../../shared/utils/period';

interface PeriodGroup {
  key: PeriodKey;
  label: string;
  icon: string;
  className: string;
  meds: Medication[];
}

@Component({
  selector: 'app-medications',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './medications.html',
  styleUrl: './medications.scss',
})
export class MedicationsPage implements OnInit {
  private readonly service = inject(MedicationService);

  readonly medications = signal<Medication[]>([]);
  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);
  readonly showForm = signal(false);
  readonly editingId = signal<string | null>(null);

  // Modelo do formulário de novo medicamento.
  form: CreateMedicationInput = this.emptyForm();

  private emptyForm(): CreateMedicationInput {
    return { name: '', dosage: '', schedule: '', notes: '', stock: undefined, startDate: '', endDate: '' };
  }

  // Medicamentos agrupados por período do dia (derivado do horário).
  readonly periodGroups = computed<PeriodGroup[]>(() => {
    const byPeriod = new Map<PeriodKey, Medication[]>();
    for (const m of this.medications()) {
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

  // Prévia do período enquanto a pessoa escolhe o horário no formulário.
  readonly formPeriod = computed(() =>
    this.form.schedule ? PERIODS[periodFromSchedule(this.form.schedule)] : null,
  );

  ngOnInit(): void {
    this.load();
  }

  load(): void {
    this.loading.set(true);
    this.error.set(null);
    this.service.list().subscribe({
      next: (list) => {
        this.medications.set(list);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Não foi possível carregar os medicamentos.');
        this.loading.set(false);
      },
    });
  }

  openForm(): void {
    this.editingId.set(null);
    this.form = this.emptyForm();
    this.error.set(null);
    this.showForm.set(true);
  }

  edit(med: Medication): void {
    this.editingId.set(med.id);
    this.form = {
      name: med.name,
      dosage: med.dosage,
      schedule: med.schedule,
      notes: med.notes ?? '',
      stock: med.stock,
      startDate: med.startDate ?? '',
      endDate: med.endDate ?? '',
      active: med.active,
    };
    this.error.set(null);
    this.showForm.set(true);
  }

  // Pausa/reativa um medicamento (envia os campos + active invertido).
  togglePause(med: Medication): void {
    this.service
      .update(med.id, {
        name: med.name,
        dosage: med.dosage,
        schedule: med.schedule,
        notes: med.notes ?? '',
        stock: med.stock,
        startDate: med.startDate ?? '',
        endDate: med.endDate ?? '',
        active: !med.active,
        period: periodFromSchedule(med.schedule),
      })
      .subscribe({
        next: () => this.load(),
        error: () => this.error.set('Não foi possível atualizar o medicamento.'),
      });
  }

  cancel(): void {
    this.showForm.set(false);
    this.editingId.set(null);
  }

  save(): void {
    if (!this.form.name.trim() || !this.form.dosage.trim() || !this.form.schedule) {
      this.error.set('Preencha nome, dose e horário.');
      return;
    }
    this.saving.set(true);
    this.error.set(null);

    // Deriva o período a partir do horário e envia ao backend para persistência.
    const payload = { ...this.form, period: periodFromSchedule(this.form.schedule) };

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
        this.error.set('Não foi possível salvar o medicamento.');
        this.saving.set(false);
      },
    });
  }

  remove(med: Medication): void {
    if (!confirm(`Excluir "${med.name}"?`)) return;
    this.service.remove(med.id).subscribe({
      next: () => this.load(),
      error: () => this.error.set('Não foi possível excluir o medicamento.'),
    });
  }
}
