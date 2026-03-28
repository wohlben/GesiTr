import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { DatePipe } from '@angular/common';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutKeys, workoutScheduleKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideTrash2 } from '@ng-icons/lucide';

@Component({
  selector: 'app-workout-schedule-list',
  imports: [PageLayout, RouterLink, DatePipe, NgIcon, TranslocoDirective],
  providers: [provideIcons({ lucideTrash2 })],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('user.schedules.title') + (workoutName() ? ' — ' + workoutName() : '')"
        [isPending]="schedulesQuery.isPending()"
        [errorMessage]="schedulesQuery.isError() ? schedulesQuery.error().message : undefined"
      >
        <a
          actions
          routerLink="./new"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          {{ t('user.schedules.newSchedule') }}
        </a>

        @if (schedulesQuery.data(); as schedules) {
          @if (schedules.length === 0) {
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ t('user.schedules.noResults') }}
            </p>
          } @else {
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                <thead>
                  <tr>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('user.schedules.startDate') }}
                    </th>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('user.schedules.active') }}
                    </th>
                    <th class="px-4 py-3"></th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                  @for (schedule of schedules; track schedule.id) {
                    <tr
                      class="cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800/50"
                      [routerLink]="['./', schedule.id, 'edit']"
                    >
                      <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                        {{ schedule.startDate | date }}
                      </td>
                      <td class="px-4 py-3 text-sm">
                        <span
                          class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium"
                          [class]="
                            schedule.active
                              ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300'
                              : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'
                          "
                        >
                          {{ schedule.active ? t('common.yes') : t('common.no') }}
                        </span>
                      </td>
                      <td class="px-4 py-3 text-right">
                        <button
                          type="button"
                          (click)="$event.stopPropagation(); deleteSchedule(schedule.id)"
                          class="inline-flex items-center rounded-md p-1.5 text-red-500 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/30"
                          [title]="t('common.delete')"
                        >
                          <ng-icon name="lucideTrash2" class="text-lg" />
                        </button>
                      </td>
                    </tr>
                  }
                </tbody>
              </table>
            </div>
          }
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class WorkoutScheduleList {
  private route = inject(ActivatedRoute);
  private userApi = inject(UserApiClient);
  private queryClient = inject(QueryClient);

  private workoutId = Number(this.route.snapshot.paramMap.get('id'));

  workoutQuery = injectQuery(() => ({
    queryKey: workoutKeys.detail(this.workoutId),
    queryFn: () => this.userApi.fetchWorkout(this.workoutId),
  }));

  workoutName = computed(() => this.workoutQuery.data()?.name);

  schedulesQuery = injectQuery(() => ({
    queryKey: workoutScheduleKeys.list(this.workoutId),
    queryFn: () => this.userApi.fetchWorkoutSchedules({ workoutId: this.workoutId }),
  }));

  deleteMutation = injectMutation(() => ({
    mutationFn: (id: number) => this.userApi.deleteWorkoutSchedule(id),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutScheduleKeys.all() });
    },
  }));

  deleteSchedule(id: number) {
    if (confirm('Delete this schedule? Existing workout logs will be preserved.')) {
      this.deleteMutation.mutate(id);
    }
  }
}
