import { Component, inject, computed } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutKeys, workoutLogKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideCalendarClock, lucideListCheck, lucideUsers } from '@ng-icons/lucide';

@Component({
  selector: 'app-workout-list',
  imports: [PageLayout, RouterLink, NgIcon, TranslocoDirective],
  providers: [provideIcons({ lucideCalendarClock, lucideListCheck, lucideUsers })],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('user.workouts.title')"
        [isPending]="workoutsQuery.isPending()"
        [errorMessage]="workoutsQuery.isError() ? workoutsQuery.error().message : undefined"
      >
        <a
          actions
          routerLink="./new"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          {{ t('user.workouts.newWorkout') }}
        </a>

        <!-- Ad-hoc Workout entry (always visible) -->
        <button
          type="button"
          (click)="startAdhoc()"
          [disabled]="adhocMutation.isPending()"
          class="mb-6 flex w-full items-center gap-3 rounded-lg border-2 border-dashed border-gray-300 p-4 text-left transition-colors hover:border-green-400 hover:bg-green-50 disabled:opacity-50 dark:border-gray-600 dark:hover:border-green-500 dark:hover:bg-green-900/20"
        >
          <span
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400"
          >
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path
                d="M6.3 2.84A1.5 1.5 0 004 4.11v11.78a1.5 1.5 0 002.3 1.27l9.344-5.891a1.5 1.5 0 000-2.538L6.3 2.841z"
              />
            </svg>
          </span>
          <div>
            <span class="text-sm font-medium text-gray-900 dark:text-gray-100">
              @if (adhocMutation.isPending()) {
                {{ t('common.starting') }}
              } @else {
                {{ t('user.workouts.adhocWorkout') }}
              }
            </span>
            <p class="text-xs text-gray-500 dark:text-gray-400">
              {{ t('user.workouts.adhocDescription') }}
            </p>
          </div>
        </button>

        @if (enrichedWorkouts(); as workouts) {
          @if (workouts.length === 0) {
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ t('user.workouts.noResults') }}
            </p>
          } @else {
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                <thead>
                  <tr>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('fields.name') }}
                    </th>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('fields.sections') }}
                    </th>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('fields.exercises') }}
                    </th>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('user.workouts.membership') }}
                    </th>
                    <th class="px-4 py-3"></th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                  @for (workout of workouts; track workout.id) {
                    <tr
                      class="cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800/50"
                      [routerLink]="['./', workout.id, 'edit']"
                    >
                      <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                        {{ workout.name }}
                      </td>
                      <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                        {{ workout.sectionCount }}
                      </td>
                      <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                        {{ workout.exerciseCount }}
                      </td>
                      <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                        @if (workout.membership) {
                          <span
                            class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium"
                            [class]="
                              workout.membership === 'admin'
                                ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300'
                                : workout.membership === 'member'
                                  ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300'
                                  : 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-300'
                            "
                          >
                            {{ t('enums.workoutGroupRole.' + workout.membership) }}
                          </span>
                        }
                      </td>
                      <td class="px-4 py-3 text-right">
                        <a
                          [routerLink]="['./', workout.id, 'schedules']"
                          (click)="$event.stopPropagation()"
                          class="inline-flex items-center rounded-md p-1.5 text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800"
                          [title]="t('user.workouts.manageSchedules')"
                        >
                          <ng-icon name="lucideCalendarClock" class="text-lg" />
                        </a>
                        <a
                          [routerLink]="['./', workout.id, 'group']"
                          (click)="$event.stopPropagation()"
                          class="inline-flex items-center rounded-md p-1.5 text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800"
                          [title]="t('user.workouts.manageGroup')"
                        >
                          <ng-icon name="lucideUsers" class="text-lg" />
                        </a>
                        <a
                          [routerLink]="['./', workout.id, 'start']"
                          (click)="$event.stopPropagation()"
                          class="inline-flex items-center rounded-md p-1.5 text-green-600 hover:bg-green-50 dark:text-green-400 dark:hover:bg-green-900/30"
                          [title]="t('user.workouts.startWorkout')"
                        >
                          <ng-icon name="lucideListCheck" class="text-lg" />
                        </a>
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
export class WorkoutList {
  private userApi = inject(UserApiClient);
  private router = inject(Router);
  private queryClient = inject(QueryClient);

  workoutsQuery = injectQuery(() => ({
    queryKey: workoutKeys.list(),
    queryFn: () => this.userApi.fetchWorkouts(),
  }));

  adhocMutation = injectMutation(() => ({
    mutationFn: () => this.userApi.startAdhocWorkoutLog(),
    onSuccess: (log) => {
      this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.all() });
      this.router.navigate(['/user/workout-logs', log.id]);
    },
  }));

  enrichedWorkouts = computed(() => {
    const workouts = this.workoutsQuery.data();
    if (!workouts) return undefined;

    return workouts.map((w) => ({
      id: w.id,
      name: w.name,
      sectionCount: w.sections?.length ?? 0,
      exerciseCount: w.sections?.reduce((sum, s) => sum + (s.items?.length ?? 0), 0) ?? 0,
      membership: w.workoutGroup?.membership,
    }));
  });

  startAdhoc() {
    this.adhocMutation.mutate();
  }
}
