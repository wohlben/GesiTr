import { Component, computed, forwardRef, input, signal } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

interface DayCell {
  date: Date;
  dayOffset: number;
  inPeriod: boolean;
  label: string;
}

@Component({
  selector: 'app-period-day-picker',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => PeriodDayPicker),
      multi: true,
    },
  ],
  template: `
    <!-- Weekday headers -->
    <div
      class="grid grid-cols-7 gap-1 text-center text-xs font-medium text-gray-500 dark:text-gray-400"
    >
      @for (dow of weekdays; track dow) {
        <div class="py-1">{{ dow }}</div>
      }
    </div>

    <!-- Calendar grid -->
    <div class="grid grid-cols-7 gap-1">
      @for (cell of cells(); track $index) {
        @if (cell) {
          <button
            type="button"
            [disabled]="!cell.inPeriod"
            (click)="toggleDay(cell)"
            class="relative flex flex-col items-center rounded-md py-1.5 text-sm transition-colors"
            [class]="cellClass(cell)"
          >
            <span>{{ cell.label }}</span>
            @if (cell.inPeriod && isSelected(cell)) {
              <span class="text-[10px] leading-none opacity-70">+{{ cell.dayOffset }}</span>
            }
          </button>
        } @else {
          <div></div>
        }
      }
    </div>
  `,
})
export class PeriodDayPicker implements ControlValueAccessor {
  periodStart = input.required<Date>();
  periodEnd = input.required<Date>();

  value = signal<Date[]>([]);
  weekdays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  private onChange: (value: Date[]) => void = () => {};
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  private onTouched: () => void = () => {};

  cells = computed(() => {
    const start = this.periodStart();
    const end = this.periodEnd();
    if (!start || !end) return [];

    const periodStartTime = new Date(start);
    periodStartTime.setHours(0, 0, 0, 0);
    const periodEndTime = new Date(end);
    periodEndTime.setHours(23, 59, 59, 999);

    const gridStart = new Date(periodStartTime);
    gridStart.setDate(gridStart.getDate() - gridStart.getDay());

    const gridEnd = new Date(periodEndTime);
    gridEnd.setDate(gridEnd.getDate() + (6 - gridEnd.getDay()));

    const cells: (DayCell | null)[] = [];
    const current = new Date(gridStart);

    while (current <= gridEnd) {
      const d = new Date(current);
      d.setHours(0, 0, 0, 0);
      const inPeriod = d >= periodStartTime && d <= periodEndTime;
      const dayOffset = inPeriod
        ? Math.round((d.getTime() - periodStartTime.getTime()) / (1000 * 60 * 60 * 24))
        : -1;

      cells.push({ date: new Date(d), dayOffset, inPeriod, label: String(d.getDate()) });
      current.setDate(current.getDate() + 1);
    }

    return cells;
  });

  writeValue(value: Date[] | null): void {
    this.value.set(value ?? []);
  }

  registerOnChange(fn: (value: Date[]) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }

  toggleDay(cell: DayCell) {
    if (!cell.inPeriod) return;
    this.onTouched();
    const current = this.value();
    const cellTime = cell.date.getTime();
    const exists = current.findIndex((d) => {
      const t = new Date(d);
      t.setHours(0, 0, 0, 0);
      return t.getTime() === cellTime;
    });

    const updated = exists >= 0 ? current.filter((_, i) => i !== exists) : [...current, cell.date];

    this.value.set(updated);
    this.onChange(updated);
  }

  isSelected(cell: DayCell): boolean {
    const cellTime = cell.date.getTime();
    return this.value().some((d) => {
      const t = new Date(d);
      t.setHours(0, 0, 0, 0);
      return t.getTime() === cellTime;
    });
  }

  cellClass(cell: DayCell): string {
    if (!cell.inPeriod) {
      return 'text-gray-300 dark:text-gray-600 cursor-default';
    }
    if (this.isSelected(cell)) {
      return 'bg-blue-600 text-white font-semibold cursor-pointer hover:bg-blue-700';
    }
    return 'text-gray-700 dark:text-gray-300 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800';
  }
}
