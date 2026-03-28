import { Component, inject, computed, signal } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import {
  workoutLogKeys,
  workoutKeys,
  workoutScheduleKeys,
  schedulePeriodKeys,
} from '$core/query-keys';
import {
  WorkoutLog,
  WorkoutLogStatusPlanning,
  WorkoutLogStatusProposed,
  WorkoutLogStatusCommitted,
  WorkoutLogStatusSkipped,
  WorkoutLogStatusBroken,
} from '$generated/user-models';
import {
  SchedulePeriod,
  PeriodStatusActive,
  PeriodStatusPlanned,
  PeriodStatusArchived,
} from '$generated/user-workoutschedule';
import { PageLayout } from '../../../layout/page-layout';
import { DayDialog } from './day-dialog';

interface PeriodBar {
  periodId: number;
  scheduleId: number;
  label: string;
  status: string;
  lane: number;
  isStart: boolean;
  isEnd: boolean;
}

interface CalendarCell {
  day: number;
  logs: WorkoutLog[];
  isToday: boolean;
  bars: PeriodBar[];
}

@Component({
  selector: 'app-calendar',
  imports: [PageLayout, DayDialog, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('user.calendar.title')"
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
          @for (day of weekDayKeys; track day) {
            <div class="py-2">{{ t('user.calendar.' + day) }}</div>
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
                <!-- Period bars -->
                @if (cell.bars.length) {
                  <div class="mb-0.5 flex w-full flex-col gap-px px-0.5">
                    @for (bar of cell.bars; track bar.periodId) {
                      <div
                        class="flex h-3.5 items-center overflow-hidden text-[9px] leading-none font-medium"
                        [class]="barClasses(bar)"
                      >
                        @if (bar.isStart) {
                          <span class="truncate px-1">{{ bar.label }}</span>
                        }
                      </div>
                    }
                  </div>
                }

                <span>{{ cell.day }}</span>
                @if (cell.logs.length) {
                  <div class="mt-0.5 flex gap-0.5">
                    @for (log of cell.logs; track log.id) {
                      <span
                        class="h-1.5 w-1.5 rounded-full"
                        [class]="
                          log.status === 'finished'
                            ? 'bg-green-500'
                            : log.status === 'partially_finished'
                              ? 'bg-yellow-500'
                              : log.status === 'in_progress'
                                ? 'bg-blue-500'
                                : log.status === 'proposed'
                                  ? 'bg-purple-300 dark:bg-purple-400'
                                  : log.status === 'committed'
                                    ? 'bg-purple-500'
                                    : log.status === 'broken'
                                      ? 'bg-orange-500'
                                      : log.status === 'skipped'
                                        ? 'bg-gray-400'
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
    </ng-container>
  `,
})
export class Calendar {
  private userApi = inject(UserApiClient);

  currentMonth = signal(new Date());
  dialogOpen = signal(false);
  selectedDate = signal<Date | undefined>(undefined);
  selectedLogs = signal<WorkoutLog[]>([]);

  weekDayKeys = ['mon', 'tue', 'wed', 'thu', 'fri', 'sat', 'sun'];

  logsQuery = injectQuery(() => ({
    queryKey: workoutLogKeys.list(),
    queryFn: () => this.userApi.fetchWorkoutLogs(),
  }));

  schedulesQuery = injectQuery(() => ({
    queryKey: workoutScheduleKeys.list(),
    queryFn: () => this.userApi.fetchWorkoutSchedules(),
  }));

  periodsQuery = injectQuery(() => ({
    queryKey: schedulePeriodKeys.list(),
    queryFn: () => this.userApi.fetchSchedulePeriods(),
  }));

  workoutsQuery = injectQuery(() => ({
    queryKey: workoutKeys.list(),
    queryFn: () => this.userApi.fetchWorkouts(),
  }));

  private logsByDate = computed(() => {
    const logs = this.logsQuery.data();
    if (!logs) return new Map<string, WorkoutLog[]>();

    const map = new Map<string, WorkoutLog[]>();
    for (const log of logs) {
      if (log.status === WorkoutLogStatusPlanning) continue;

      // Committed/proposed logs use dueStart as the calendar date; others use date
      const isCommitment =
        log.status === WorkoutLogStatusProposed ||
        log.status === WorkoutLogStatusCommitted ||
        log.status === WorkoutLogStatusSkipped ||
        log.status === WorkoutLogStatusBroken;
      const dateStr = isCommitment ? (log.dueStart ?? log.date) : log.date;
      if (!dateStr) continue;

      const key = dateStr.substring(0, 10);
      const existing = map.get(key);
      if (existing) {
        existing.push(log);
      } else {
        map.set(key, [log]);
      }
    }
    return map;
  });

  /** Build a lookup: scheduleId → workout name */
  private workoutNameByScheduleId = computed(() => {
    const schedules = this.schedulesQuery.data();
    const workouts = this.workoutsQuery.data();
    if (!schedules || !workouts) return new Map<number, string>();

    const workoutMap = new Map(workouts.map((w) => [w.id, w.name]));
    const result = new Map<number, string>();
    for (const s of schedules) {
      result.set(s.id, workoutMap.get(s.workoutId) ?? 'Schedule');
    }
    return result;
  });

  /**
   * Compute period bars for the current month. Each period gets a lane
   * (vertical slot) so overlapping periods from different schedules stack.
   */
  private periodLanes = computed(() => {
    const periods = this.periodsQuery.data();
    if (!periods?.length)
      return { lanes: [] as { period: SchedulePeriod; lane: number }[], maxLane: 0 };

    // Sort by start date, then longer periods first (for stable lane assignment)
    const sorted = [...periods].sort((a, b) => {
      const startCmp = a.periodStart.localeCompare(b.periodStart);
      if (startCmp !== 0) return startCmp;
      return b.periodEnd.localeCompare(a.periodEnd);
    });

    const lanes: { period: SchedulePeriod; lane: number }[] = [];
    // Track end date per lane (ISO string) for overlap detection
    const laneEnds: string[] = [];

    for (const period of sorted) {
      let lane = 0;
      while (lane < laneEnds.length && laneEnds[lane] > period.periodStart.substring(0, 10)) {
        lane++;
      }
      if (lane === laneEnds.length) laneEnds.push('');
      laneEnds[lane] = period.periodEnd.substring(0, 10);
      lanes.push({ period, lane });
    }

    return { lanes, maxLane: laneEnds.length };
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
    const nameMap = this.workoutNameByScheduleId();
    const { lanes } = this.periodLanes();

    const cells: (null | CalendarCell)[] = [];

    for (let i = 0; i < startOffset; i++) {
      cells.push(null);
    }

    for (let day = 1; day <= daysInMonth; day++) {
      const key = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;

      // Build bars for this day
      const bars: PeriodBar[] = [];
      for (const { period, lane } of lanes) {
        const pStart = period.periodStart.substring(0, 10);
        const pEnd = period.periodEnd.substring(0, 10);
        if (key >= pStart && key < pEnd) {
          bars.push({
            periodId: period.id,
            scheduleId: period.scheduleId,
            label: nameMap.get(period.scheduleId) ?? 'Schedule',
            status: period.status,
            lane,
            isStart: key === pStart,
            isEnd: this.nextDay(key) >= pEnd,
          });
        }
      }
      // Sort by lane so they render consistently
      bars.sort((a, b) => a.lane - b.lane);

      cells.push({
        day,
        logs: logsByDate.get(key) ?? [],
        isToday: day === todayKey,
        bars,
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

  openDay(cell: CalendarCell) {
    if (!cell.logs.length) return;
    const d = this.currentMonth();
    this.selectedDate.set(new Date(d.getFullYear(), d.getMonth(), cell.day));
    this.selectedLogs.set(cell.logs);
    this.dialogOpen.set(true);
  }

  barClasses(bar: PeriodBar): string {
    const rounded =
      bar.isStart && bar.isEnd
        ? 'rounded'
        : bar.isStart
          ? 'rounded-l'
          : bar.isEnd
            ? 'rounded-r'
            : '';

    switch (bar.status) {
      case PeriodStatusActive:
        return `${rounded} bg-purple-200 text-purple-900 dark:bg-purple-800/60 dark:text-purple-200`;
      case PeriodStatusPlanned:
        return `${rounded} bg-indigo-200 text-indigo-900 dark:bg-indigo-800/60 dark:text-indigo-200`;
      case PeriodStatusArchived:
        return `${rounded} bg-gray-200 text-gray-600 dark:bg-gray-700/60 dark:text-gray-400`;
      default:
        return `${rounded} bg-gray-200 text-gray-600 dark:bg-gray-700/60 dark:text-gray-400`;
    }
  }

  private nextDay(dateKey: string): string {
    const d = new Date(dateKey + 'T00:00:00');
    d.setDate(d.getDate() + 1);
    return d.toISOString().substring(0, 10);
  }
}
