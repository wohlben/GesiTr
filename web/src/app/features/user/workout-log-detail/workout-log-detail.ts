import { Component, inject, computed, effect } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { DatePipe } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { formatBreak } from '$core/format-utils';
import { workoutLogKeys } from '$core/query-keys';
import {
  WorkoutLogExerciseSet,
  WorkoutLogStatusFinished,
  WorkoutLogStatusAborted,
} from '$generated/user-models';
import { PageLayout } from '../../../layout/page-layout';
import { WorkoutLogDetailStore } from './workout-log-detail.store';

@Component({
  selector: 'app-workout-log-detail',
  imports: [PageLayout, RouterLink, DatePipe],
  providers: [WorkoutLogDetailStore],
  template: `
    <app-page-layout
      [header]="logQuery.data()?.name ?? 'Workout Log'"
      [isPending]="logQuery.isPending()"
      [errorMessage]="logQuery.isError() ? logQuery.error().message : undefined"
    >
      <a
        actions
        routerLink="/user/workouts"
        class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
      >
        Back
      </a>

      @if (logQuery.data(); as log) {
        <!-- Log header -->
        <div
          class="mb-6 flex items-center justify-between rounded-lg border border-gray-200 p-4 dark:border-gray-700"
        >
          <div>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ log.date | date }}
            </p>
            @if (log.notes) {
              <p class="mt-1 text-sm text-gray-600 dark:text-gray-300">{{ log.notes }}</p>
            }
          </div>
          <div class="flex items-center gap-2">
            @if (log.status === 'in_progress') {
              <button
                type="button"
                (click)="abandonWorkout()"
                [disabled]="abandonMutation.isPending()"
                class="rounded-md border border-red-300 px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 disabled:opacity-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20"
              >
                Abandon
              </button>
            }
            <span
              [class]="
                log.status === 'finished'
                  ? 'text-green-600 dark:text-green-400'
                  : log.status === 'aborted'
                    ? 'text-red-500 dark:text-red-400'
                    : 'text-gray-400'
              "
            >
              @if (log.status === 'finished') {
                <svg class="h-8 w-8" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fill-rule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z"
                    clip-rule="evenodd"
                  />
                </svg>
              } @else if (log.status === 'aborted') {
                <svg class="h-8 w-8" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fill-rule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z"
                    clip-rule="evenodd"
                  />
                </svg>
              } @else {
                <svg class="h-8 w-8" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fill-rule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-11.25a.75.75 0 00-1.5 0v4.59L7.3 9.24a.75.75 0 00-1.1 1.02l3.25 3.5a.75.75 0 001.1 0l3.25-3.5a.75.75 0 10-1.1-1.02l-1.95 2.1V6.75z"
                    clip-rule="evenodd"
                  />
                </svg>
              }
            </span>
          </div>
        </div>

        <!-- Sections -->
        <div class="space-y-6">
          @for (section of log.sections; track section.id) {
            <div class="rounded-lg border border-gray-200 dark:border-gray-700">
              <!-- Section header -->
              <div
                class="flex items-center justify-between border-b border-gray-200 px-4 py-3 dark:border-gray-700"
              >
                <div class="flex items-center gap-2">
                  <span
                    class="rounded-full px-2 py-0.5 text-xs font-medium"
                    [class]="
                      section.type === 'main'
                        ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300'
                        : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400'
                    "
                  >
                    {{ section.type }}
                  </span>
                  @if (section.label) {
                    <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{
                      section.label
                    }}</span>
                  }
                </div>
                @if (section.status === 'finished') {
                  <svg
                    class="h-5 w-5 text-green-500 dark:text-green-400"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z"
                      clip-rule="evenodd"
                    />
                  </svg>
                }
              </div>

              <!-- Exercises in section -->
              <div class="divide-y divide-gray-100 dark:divide-gray-800">
                @for (exercise of section.exercises; track exercise.id; let lastEx = $last) {
                  <div class="p-4">
                    <!-- Exercise header -->
                    <div class="mb-3 flex items-center justify-between">
                      <div class="flex items-center gap-2">
                        <span class="text-sm font-semibold text-gray-900 dark:text-gray-100">
                          {{
                            store.exerciseNames()[exercise.sourceExerciseSchemeId] || 'Loading...'
                          }}
                        </span>
                        @if (exercise.status === 'finished') {
                          <svg
                            class="h-4 w-4 text-green-500 dark:text-green-400"
                            viewBox="0 0 20 20"
                            fill="currentColor"
                          >
                            <path
                              fill-rule="evenodd"
                              d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z"
                              clip-rule="evenodd"
                            />
                          </svg>
                        }
                      </div>
                      <span class="text-xs text-gray-500 dark:text-gray-400">
                        {{ exercise.targetMeasurementType }}
                      </span>
                    </div>

                    <!-- Sets table -->
                    <div class="overflow-x-auto">
                      <table class="w-full text-sm">
                        <thead>
                          <tr class="text-left text-xs text-gray-500 uppercase dark:text-gray-400">
                            <th class="pb-1 pr-3">Set</th>
                            <th class="pb-1 pr-3">Target</th>
                            <th class="pb-1 pr-3">Actual</th>
                            @if (hasBreaks(exercise.sets)) {
                              <th class="pb-1 pr-3">Rest</th>
                            }
                            <th class="pb-1 w-10 text-center">Done</th>
                          </tr>
                        </thead>
                        <tbody>
                          @for (set of exercise.sets; track set.id) {
                            <tr class="border-t border-gray-100 dark:border-gray-800">
                              <td class="py-2 pr-3 font-medium text-gray-900 dark:text-gray-100">
                                {{ set.setNumber }}
                              </td>
                              <td class="py-2 pr-3 text-gray-600 dark:text-gray-400">
                                {{ formatTarget(set, exercise.targetMeasurementType) }}
                              </td>
                              <td class="py-2 pr-3 text-gray-900 dark:text-gray-100">
                                {{ formatActual(set, exercise.targetMeasurementType) }}
                              </td>
                              @if (hasBreaks(exercise.sets)) {
                                <td class="py-2 pr-3 text-gray-500 dark:text-gray-400">
                                  {{ formatBreak(set.breakAfterSeconds) }}
                                </td>
                              }
                              <td class="py-2 text-center">
                                <button
                                  type="button"
                                  (click)="toggleSet(set)"
                                  [disabled]="isSetTerminal(set)"
                                  class="inline-flex items-center justify-center rounded p-0.5 transition-colors"
                                  [class]="
                                    set.status === 'finished'
                                      ? 'text-green-600 hover:text-green-700 dark:text-green-400 dark:hover:text-green-300'
                                      : set.status === 'aborted'
                                        ? 'text-red-500 dark:text-red-400 cursor-not-allowed'
                                        : 'text-gray-300 hover:text-gray-400 dark:text-gray-600 dark:hover:text-gray-500'
                                  "
                                >
                                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                    @if (set.status === 'finished') {
                                      <path
                                        fill-rule="evenodd"
                                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z"
                                        clip-rule="evenodd"
                                      />
                                    } @else if (set.status === 'aborted') {
                                      <path
                                        fill-rule="evenodd"
                                        d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z"
                                        clip-rule="evenodd"
                                      />
                                    } @else {
                                      <path
                                        fill-rule="evenodd"
                                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm-.75-4.75a.75.75 0 001.5 0V8.66l1.95 2.1a.75.75 0 101.1-1.02l-3.25-3.5a.75.75 0 00-1.1 0L6.2 9.74a.75.75 0 001.1 1.02l1.95-2.1v4.59z"
                                        clip-rule="evenodd"
                                      />
                                    }
                                  </svg>
                                </button>
                              </td>
                            </tr>
                          }
                        </tbody>
                      </table>
                    </div>

                    <!-- Break after exercise -->
                    @if (exercise.breakAfterSeconds && !lastEx) {
                      <div
                        class="mt-2 flex items-center gap-1 text-xs text-gray-400 dark:text-gray-500"
                      >
                        <svg class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                          <path
                            fill-rule="evenodd"
                            d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z"
                            clip-rule="evenodd"
                          />
                        </svg>
                        {{ formatBreak(exercise.breakAfterSeconds) }} rest before next exercise
                      </div>
                    }
                  </div>
                }
              </div>
            </div>
          }
        </div>
      }
    </app-page-layout>
  `,
})
export class WorkoutLogDetail {
  private userApi = inject(UserApiClient);
  private route = inject(ActivatedRoute);
  private queryClient = injectQueryClient();
  private params = toSignal(this.route.paramMap);

  readonly store = inject(WorkoutLogDetailStore);

  private id = computed(() => Number(this.params()?.get('id')));
  private initialized = false;

  logQuery = injectQuery(() => ({
    queryKey: workoutLogKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchWorkoutLog(this.id()),
    enabled: !!this.id(),
  }));

  private toggleMutation = injectMutation(() => ({
    mutationFn: (set: WorkoutLogExerciseSet) =>
      this.userApi.updateWorkoutLogExerciseSet(set.id, {
        status: WorkoutLogStatusFinished,
        actualReps: set.actualReps,
        actualWeight: set.actualWeight,
        actualDuration: set.actualDuration,
        actualDistance: set.actualDistance,
        actualTime: set.actualTime,
      }),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.detail(this.id()) });
    },
  }));

  abandonMutation = injectMutation(() => ({
    mutationFn: () => this.userApi.abandonWorkoutLog(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.detail(this.id()) });
    },
  }));

  constructor() {
    effect(() => {
      const data = this.logQuery.data();
      if (!data || this.initialized) return;
      this.initialized = true;
      this.store.loadExerciseNames(data.sections ?? []);
    });
  }

  toggleSet(set: WorkoutLogExerciseSet) {
    if (this.isSetTerminal(set)) return;
    this.toggleMutation.mutate(set);
  }

  isSetTerminal(set: WorkoutLogExerciseSet): boolean {
    return set.status === WorkoutLogStatusFinished || set.status === WorkoutLogStatusAborted;
  }

  abandonWorkout() {
    this.abandonMutation.mutate();
  }

  hasBreaks(sets: WorkoutLogExerciseSet[]): boolean {
    return sets.some((s) => s.breakAfterSeconds != null);
  }

  formatTarget(set: WorkoutLogExerciseSet, measurementType: string): string {
    if (measurementType === 'REP_BASED') {
      const parts: string[] = [];
      if (set.targetReps != null) parts.push(`${set.targetReps} reps`);
      if (set.targetWeight != null) parts.push(`${set.targetWeight}kg`);
      return parts.join(' @ ') || '-';
    }
    if (measurementType === 'TIME_BASED') {
      if (set.targetDuration != null) return `${set.targetDuration}s`;
      return '-';
    }
    if (measurementType === 'DISTANCE_BASED') {
      if (set.targetDistance != null) return `${set.targetDistance}m`;
      return '-';
    }
    return '-';
  }

  formatActual(set: WorkoutLogExerciseSet, measurementType: string): string {
    if (set.status !== WorkoutLogStatusFinished) return '-';
    if (measurementType === 'REP_BASED') {
      const parts: string[] = [];
      if (set.actualReps != null) parts.push(`${set.actualReps} reps`);
      if (set.actualWeight != null) parts.push(`${set.actualWeight}kg`);
      return parts.join(' @ ') || '-';
    }
    if (measurementType === 'TIME_BASED') {
      if (set.actualDuration != null) return `${set.actualDuration}s`;
      return '-';
    }
    if (measurementType === 'DISTANCE_BASED') {
      if (set.actualDistance != null) return `${set.actualDistance}m`;
      return '-';
    }
    return '-';
  }

  formatBreak = formatBreak;

  date(value: string): string {
    return new Date(value).toLocaleDateString();
  }
}
