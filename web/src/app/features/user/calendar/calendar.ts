import { Component, inject, computed, signal } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutLogKeys } from '$core/query-keys';
import { WorkoutLog, WorkoutLogStatusPlanning } from '$generated/user-models';
import { PageLayout } from '../../../layout/page-layout';
import { DayDialog } from './day-dialog';

@Component({
  selector: 'app-calendar',
  imports: [PageLayout, DayDialog],
  template: `
    <app-page-layout
      header="Calendar"
      [isPending]="logsQuery.isPending()"
      [errorMessage]="logsQuery.isError() ? logsQuery.error().message : undefined"
    >
      <!-- Month navigation -->
      <div class="mb-4 flex items-center justify-between">
        <button
          type="button"
          (click)="prevMonth()"
          class="rounded-md p-2 text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 19l-7-7 7-7"
            />
          </svg>
        </button>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
          {{ monthLabel() }}
        </h2>
        <button
          type="button"
          (click)="nextMonth()"
          class="rounded-md p-2 text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 5l7 7-7 7"
            />
          </svg>
        </button>
      </div>

      <!-- Weekday headers -->
      <div
        class="grid grid-cols-7 gap-px text-center text-xs font-medium text-gray-500 dark:text-gray-400"
      >
        @for (day of weekDays; track day) {
          <div class="py-2">{{ day }}</div>
        }
      </div>

      <!-- Calendar grid -->
      <div class="grid grid-cols-7 gap-px">
        @for (cell of calendarCells(); track $index) {
          @if (cell) {
            <button
              type="button"
              [disabled]="!cell.logs.length"
              (click)="openDay(cell)"
              class="flex min-h-12 flex-col items-center rounded-md py-1.5 text-sm transition-colors"
              [class]="
                cell.isToday
                  ? cell.logs.length
                    ? 'cursor-pointer bg-blue-50 font-semibold text-blue-700 hover:bg-blue-100 dark:bg-blue-900/20 dark:text-blue-300 dark:hover:bg-blue-900/40'
                    : 'bg-blue-50 font-semibold text-blue-700 dark:bg-blue-900/20 dark:text-blue-300'
                  : cell.logs.length
                    ? 'cursor-pointer text-gray-900 hover:bg-gray-100 dark:text-gray-100 dark:hover:bg-gray-800'
                    : 'text-gray-400 dark:text-gray-600'
              "
            >
              <span>{{ cell.day }}</span>
              @if (cell.logs.length) {
                <div class="mt-0.5 flex gap-0.5">
                  @for (log of cell.logs; track log.id) {
                    <span
                      class="h-1.5 w-1.5 rounded-full"
                      [class]="
                        log.status === 'finished'
                          ? 'bg-green-500'
                          : log.status === 'in_progress'
                            ? 'bg-blue-500'
                            : 'bg-red-500'
                      "
                    ></span>
                  }
                </div>
              }
            </button>
          } @else {
            <div></div>
          }
        }
      </div>
    </app-page-layout>

    <app-day-dialog
      [open]="dialogOpen()"
      [date]="selectedDate()"
      [logs]="selectedLogs()"
      (closed)="dialogOpen.set(false)"
    />
  `,
})
export class Calendar {
  private userApi = inject(UserApiClient);

  currentMonth = signal(new Date());
  dialogOpen = signal(false);
  selectedDate = signal<Date | undefined>(undefined);
  selectedLogs = signal<WorkoutLog[]>([]);

  weekDays = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

  logsQuery = injectQuery(() => ({
    queryKey: workoutLogKeys.list(),
    queryFn: () => this.userApi.fetchWorkoutLogs(),
  }));

  private logsByDate = computed(() => {
    const logs = this.logsQuery.data();
    if (!logs) return new Map<string, WorkoutLog[]>();

    const map = new Map<string, WorkoutLog[]>();
    for (const log of logs) {
      if (log.status === WorkoutLogStatusPlanning || !log.date) continue;
      const key = log.date.substring(0, 10);
      const existing = map.get(key);
      if (existing) {
        existing.push(log);
      } else {
        map.set(key, [log]);
      }
    }
    return map;
  });

  monthLabel = computed(() => {
    const d = this.currentMonth();
    return d.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  });

  calendarCells = computed(() => {
    const d = this.currentMonth();
    const year = d.getFullYear();
    const month = d.getMonth();
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    const daysInMonth = lastDay.getDate();

    // Monday=0 ... Sunday=6
    let startOffset = firstDay.getDay() - 1;
    if (startOffset < 0) startOffset = 6;

    const today = new Date();
    const todayKey =
      today.getFullYear() === year && today.getMonth() === month ? today.getDate() : -1;

    const logsByDate = this.logsByDate();
    const cells: (null | { day: number; logs: WorkoutLog[]; isToday: boolean })[] = [];

    for (let i = 0; i < startOffset; i++) {
      cells.push(null);
    }

    for (let day = 1; day <= daysInMonth; day++) {
      const key = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
      cells.push({
        day,
        logs: logsByDate.get(key) ?? [],
        isToday: day === todayKey,
      });
    }

    return cells;
  });

  prevMonth() {
    const d = this.currentMonth();
    this.currentMonth.set(new Date(d.getFullYear(), d.getMonth() - 1, 1));
  }

  nextMonth() {
    const d = this.currentMonth();
    this.currentMonth.set(new Date(d.getFullYear(), d.getMonth() + 1, 1));
  }

  openDay(cell: { day: number; logs: WorkoutLog[] }) {
    if (!cell.logs.length) return;
    const d = this.currentMonth();
    this.selectedDate.set(new Date(d.getFullYear(), d.getMonth(), cell.day));
    this.selectedLogs.set(cell.logs);
    this.dialogOpen.set(true);
  }
}
