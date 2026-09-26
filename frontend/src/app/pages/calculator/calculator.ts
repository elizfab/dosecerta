import { Component, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  calculateDose,
  formatDoseResult,
  DeviceType,
} from '../../shared/utils/dose-calculator';

interface CalcHistoryItem {
  concentration: number;
  dose: number;
  device: DeviceType;
  deviceLabel: string;
  text: string;
  when: string; // ISO
}

const HISTORY_KEY = 'dc-calc-history';
const DEVICE_LABELS: Record<DeviceType, string> = {
  ml: 'Seringa em mL',
  ui: 'Insulina (UI)',
  gotas: 'Gotas',
};

@Component({
  selector: 'app-calculator',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './calculator.html',
  styleUrl: './calculator.scss',
})
export class CalculatorPage {
  concentration: number | null = null;
  dose: number | null = null;
  device: DeviceType = 'ml';

  readonly resultText = signal<string | null>(null);
  readonly error = signal<string | null>(null);
  readonly history = signal<CalcHistoryItem[]>(this.loadHistory());

  calculate(): void {
    this.error.set(null);
    this.resultText.set(null);
    try {
      const result = calculateDose(
        this.concentration ?? 0,
        this.dose ?? 0,
        this.device,
      );
      const text = formatDoseResult(result, this.device);
      this.resultText.set(text);
      this.pushHistory({
        concentration: this.concentration ?? 0,
        dose: this.dose ?? 0,
        device: this.device,
        deviceLabel: DEVICE_LABELS[this.device],
        text,
        when: new Date().toISOString(),
      });
    } catch (e) {
      this.error.set(e instanceof Error ? e.message : 'Erro ao calcular a dose.');
    }
  }

  clearHistory(): void {
    this.history.set([]);
    try {
      localStorage.removeItem(HISTORY_KEY);
    } catch {
      /* ignora */
    }
  }

  formatWhen(iso: string): string {
    const d = new Date(iso);
    return d.toLocaleString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  }

  private pushHistory(item: CalcHistoryItem): void {
    const next = [item, ...this.history()].slice(0, 20);
    this.history.set(next);
    try {
      localStorage.setItem(HISTORY_KEY, JSON.stringify(next));
    } catch {
      /* ignora */
    }
  }

  private loadHistory(): CalcHistoryItem[] {
    try {
      const raw = localStorage.getItem(HISTORY_KEY);
      if (raw) return JSON.parse(raw) as CalcHistoryItem[];
    } catch {
      /* ignora */
    }
    return [];
  }
}
