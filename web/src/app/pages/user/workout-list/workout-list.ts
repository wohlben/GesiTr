import { Component, inject, computed } from '@angular/core';
import { RouterLink } from '@angular/router';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-workout-list',
  imports: [PageLayout, RouterLink],
  template: `
    <app-page-layout
      header="My Workouts"
      [isPending]="workoutsQuery.isPending()"
      [errorMessage]="workoutsQuery.isError() ? workoutsQuery.error().message : undefined"
    >
      <a
        actions
        routerLink="./new"
        class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
      >
        New Workout
      </a>

      @if (enrichedWorkouts(); as workouts) {
        @if (workouts.length === 0) {
          <p class="text-sm text-gray-500 dark:text-gray-400">No workouts yet.</p>
        } @else {
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead>
                <tr>
                  <th
                    class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                  >
                    Name
                  </th>
                  <th
                    class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                  >
                    Sections
                  </th>
                  <th
                    class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                  >
                    Exercises
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
                    <td class="px-4 py-3 text-right">
                      <a
                        [routerLink]="['./', workout.id, 'start']"
                        (click)="$event.stopPropagation()"
                        class="inline-flex items-center rounded-md p-1.5 text-green-600 hover:bg-green-50 dark:text-green-400 dark:hover:bg-green-900/30"
                        title="Start workout"
                      >
                        <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                          <path
                            d="M6.3 2.84A1.5 1.5 0 004 4.11v11.78a1.5 1.5 0 002.3 1.27l9.344-5.891a1.5 1.5 0 000-2.538L6.3 2.841z"
                          />
                        </svg>
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
  `,
})
export class WorkoutList {
  private userApi = inject(UserApiClient);

  workoutsQuery = injectQuery(() => ({
    queryKey: workoutKeys.list(),
    queryFn: () => this.userApi.fetchWorkouts(),
  }));

  enrichedWorkouts = computed(() => {
    const workouts = this.workoutsQuery.data();
    if (!workouts) return undefined;

    return workouts.map((w) => ({
      id: w.id,
      name: w.name,
      sectionCount: w.sections?.length ?? 0,
      exerciseCount: w.sections?.reduce((sum, s) => sum + (s.exercises?.length ?? 0), 0) ?? 0,
    }));
  });
}
