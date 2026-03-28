import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { DatePipe } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { schedulePeriodKeys, workoutLogKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-workout-schedule-period',
  imports: [PageLayout, RouterLink, DatePipe, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('user.schedules.periodDetail')"
        [isPending]="logsQuery.isPending()"
        [errorMessage]="logsQuery.isError() ? logsQuery.error().message : undefined"
      >
        @if (period(); as p) {
          <div class="mb-6 rounded-lg border border-gray-200 p-4 dark:border-gray-700">
            <div class="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span class="text-gray-500 dark:text-gray-400">{{
                  t('user.schedules.startDate')
                }}</span>
                <p class="font-medium text-gray-900 dark:text-gray-100">
                  {{ p.periodStart | date }}
                </p>
              </div>
              <div>
                <span class="text-gray-500 dark:text-gray-400">{{
                  t('user.schedules.endDate')
                }}</span>
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ p.periodEnd | date }}</p>
              </div>
            </div>
          </div>
        }

        <h3 class="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t('user.schedules.workoutLogs') }}
        </h3>

        @if (logsQuery.data(); as logs) {
          @if (logs.length === 0) {
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ t('common.noResults') }}
            </p>
          } @else {
            <div class="space-y-2">
              @for (log of logs; track log.id) {
                <a
                  [routerLink]="['/user/workout-logs', log.id]"
                  class="flex items-center justify-between rounded-md border border-gray-200 px-4 py-3 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-700/50"
                >
                  <div>
                    <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{
                      log.name
                    }}</span>
                    @if (log.date) {
                      <span class="ml-2 text-xs text-gray-500 dark:text-gray-400">{{
                        log.date | date
                      }}</span>
                    }
                  </div>
                  <span
                    class="rounded-full px-2 py-0.5 text-xs font-medium"
                    [class]="statusClass(log.status)"
                  >
                    {{ t('enums.workoutLogStatus.' + log.status) }}
                  </span>
                </a>
              }
            </div>
          }
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class WorkoutSchedulePeriod {
  private route = inject(ActivatedRoute);
  private userApi = inject(UserApiClient);
  private params = toSignal(this.route.paramMap);

  private scheduleId = computed(() => Number(this.params()?.get('scheduleId')));
  private periodId = computed(() => Number(this.params()?.get('periodId')));

  periodsQuery = injectQuery(() => ({
    queryKey: schedulePeriodKeys.list(this.scheduleId()),
    queryFn: () => this.userApi.fetchSchedulePeriods({ scheduleId: this.scheduleId() }),
    enabled: !!this.scheduleId(),
  }));

  period = computed(() => {
    const periods = this.periodsQuery.data();
    if (!periods) return undefined;
    return periods.find((p) => p.id === this.periodId());
  });

  logsQuery = injectQuery(() => ({
    queryKey: workoutLogKeys.list({ periodId: this.periodId() }),
    queryFn: () => this.userApi.fetchWorkoutLogs({ periodId: this.periodId() }),
    enabled: !!this.periodId(),
  }));

  statusClass(status: string): string {
    switch (status) {
      case 'finished':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300';
      case 'in_progress':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300';
      case 'proposed':
      case 'committed':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300';
      case 'broken':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-300';
      case 'skipped':
        return 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400';
      case 'aborted':
        return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300';
      default:
        return 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400';
    }
  }
}
