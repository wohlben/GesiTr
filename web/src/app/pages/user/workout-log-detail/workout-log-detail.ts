import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { DatePipe } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutLogKeys } from '$core/query-keys';
import {
  WorkoutLogExerciseSet,
  WorkoutLogItemStatusFinished,
  WorkoutLogItemStatusSkipped,
  WorkoutLogStatusAdhoc,
} from '$generated/user-models';
import { SetCompletionPayload } from './workout-log-view-items';
import { PageLayout } from '../../../layout/page-layout';
import { WorkoutLogDetailStore } from './workout-log-detail.store';
import { WorkoutLogReview } from './workout-log-review';
import { WorkoutLogActive } from './workout-log-active';
import { AdhocAddExerciseDialog } from './adhoc-add-exercise-dialog';
import { HlmButton } from '@spartan-ng/helm/button';

@Component({
  selector: 'app-workout-log-detail',
  imports: [
    PageLayout,
    RouterLink,
    DatePipe,
    WorkoutLogReview,
    WorkoutLogActive,
    AdhocAddExerciseDialog,
    HlmButton,
  ],
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
            @if (log.status === 'adhoc') {
              <button
                hlmBtn
                variant="default"
                (click)="finishWorkout()"
                [disabled]="finishMutation.isPending()"
                class="bg-green-600 text-white hover:bg-green-700 dark:bg-green-600 dark:hover:bg-green-700"
              >
                @if (finishMutation.isPending()) {
                  Finishing...
                } @else {
                  Finish Workout
                }
              </button>
            }
            @if (log.status === 'in_progress' || log.status === 'adhoc') {
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
                  : log.status === 'partially_finished'
                    ? 'text-yellow-500 dark:text-yellow-400'
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
              } @else if (log.status === 'partially_finished') {
                <svg class="h-8 w-8" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fill-rule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-11.25a.75.75 0 00-1.5 0v4.59L7.3 9.24a.75.75 0 00-1.1 1.02l3.25 3.5a.75.75 0 001.1 0l3.25-3.5a.75.75 0 10-1.1-1.02l-1.95 2.1V6.75z"
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
              } @else if (log.status !== 'adhoc') {
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

        @if (log.status === 'in_progress' || log.status === 'adhoc') {
          <app-workout-log-active
            [log]="log"
            [exerciseNames]="store.exerciseNames()"
            (setCompleted)="completeSet($event)"
            (setSkipped)="skipSet($event)"
          />

          <!-- Add Exercise button for adhoc workouts -->
          @if (log.status === 'adhoc') {
            <div class="mt-6">
              <button
                hlmBtn
                variant="outline"
                (click)="addExerciseDialogOpen.set(true)"
                class="w-full"
              >
                + Add Exercise
              </button>
            </div>

            <app-adhoc-add-exercise-dialog
              [open]="addExerciseDialogOpen()"
              [sectionId]="adhocSectionId()"
              [logId]="id()"
              [exerciseCount]="adhocExerciseCount()"
              (closed)="onAddExerciseDialogClosed()"
            />
          }
        } @else {
          <app-workout-log-review [log]="log" [exerciseNames]="store.exerciseNames()" />
        }
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

  id = computed(() => Number(this.params()?.get('id')));

  addExerciseDialogOpen = signal(false);

  logQuery = injectQuery(() => ({
    queryKey: workoutLogKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchWorkoutLog(this.id()),
    enabled: !!this.id(),
  }));

  adhocSectionId = computed(() => {
    const log = this.logQuery.data();
    if (!log || log.status !== WorkoutLogStatusAdhoc) return 0;
    return log.sections?.[0]?.id ?? 0;
  });

  adhocExerciseCount = computed(() => {
    const log = this.logQuery.data();
    if (!log || log.status !== WorkoutLogStatusAdhoc) return 0;
    return log.sections?.[0]?.exercises?.length ?? 0;
  });

  private completeMutation = injectMutation(() => ({
    mutationFn: (payload: SetCompletionPayload) =>
      this.userApi.updateWorkoutLogExerciseSet(payload.setId, {
        status: WorkoutLogItemStatusFinished,
        actualReps: payload.actualReps,
        actualWeight: payload.actualWeight,
        actualDuration: payload.actualDuration,
        actualDistance: payload.actualDistance,
        actualTime: payload.actualTime,
      }),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.detail(this.id()) });
    },
  }));

  private skipMutation = injectMutation(() => ({
    mutationFn: (set: WorkoutLogExerciseSet) =>
      this.userApi.updateWorkoutLogExerciseSet(set.id, {
        status: WorkoutLogItemStatusSkipped,
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

  finishMutation = injectMutation(() => ({
    mutationFn: () => this.userApi.finishWorkoutLog(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.detail(this.id()) });
    },
  }));

  constructor() {
    // Load exercise names whenever log data changes
    effect(() => {
      const data = this.logQuery.data();
      if (!data) return;
      this.store.loadExerciseNames(data.sections ?? []);
    });
  }

  completeSet(payload: SetCompletionPayload) {
    this.completeMutation.mutate(payload);
  }

  skipSet(set: WorkoutLogExerciseSet) {
    this.skipMutation.mutate(set);
  }

  abandonWorkout() {
    this.abandonMutation.mutate();
  }

  finishWorkout() {
    this.finishMutation.mutate();
  }

  onAddExerciseDialogClosed() {
    this.addExerciseDialogOpen.set(false);
    // Refresh the log to pick up the newly added exercise
    this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.detail(this.id()) });
  }
}
