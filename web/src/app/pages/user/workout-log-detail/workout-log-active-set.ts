import { Component, effect, input, output, signal, linkedSignal } from '@angular/core';
import { form, FormField } from '@angular/forms/signals';
import { TranslocoDirective } from '@jsverse/transloco';
import { formatTarget, formatActual, formatSetValue } from '$core/format-utils';
import { WorkoutLogItemStatusSkipped } from '$generated/user-models';
import { ViewItemSet, SetCompletionPayload } from './workout-log-view-items';

@Component({
  selector: 'app-workout-log-active-set',
  imports: [FormField, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      @if (peeked()) {
        <div class="flex flex-col rounded-lg bg-blue-100 px-4 py-5 dark:bg-blue-900/30">
          <div class="mb-1 flex items-center justify-between">
            <span class="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {{ data().exerciseName }}
            </span>
            @if (editing()) {
              <button
                type="button"
                (click)="editing.set(false); $event.stopPropagation()"
                class="rounded-md p-1 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
              >
                <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z"
                  />
                </svg>
              </button>
            } @else if (data().role === 'completed') {
              <button
                type="button"
                (click)="startEditing(); $event.stopPropagation()"
                class="rounded-md p-1 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
              >
                <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    d="M2.695 14.763l-1.262 3.154a.5.5 0 00.65.65l3.155-1.262a4 4 0 001.343-.885L17.5 5.5a2.121 2.121 0 00-3-3L3.58 13.42a4 4 0 00-.885 1.343z"
                  />
                </svg>
              </button>
            }
          </div>
          <div class="mb-4 text-sm text-gray-500 dark:text-gray-400">
            {{
              t('user.workoutLog.setOf', { current: data().set.setNumber, total: data().setCount })
            }}
          </div>
          <div class="mb-2 text-xs text-gray-400 dark:text-gray-500">
            {{ t('user.workoutLog.target') }}
            {{ formatTarget(data().set, data().exercise.targetMeasurementType) }}
          </div>

          @if (editing()) {
            <div class="mb-5 flex flex-col justify-center gap-3">
              @if (data().exercise.targetMeasurementType === 'REP_BASED') {
                <div class="flex gap-3">
                  <label class="flex flex-1 flex-col gap-1">
                    <span class="text-xs font-medium text-gray-600 dark:text-gray-400">{{
                      t('fields.reps')
                    }}</span>
                    <input
                      type="number"
                      inputmode="numeric"
                      [formField]="editForm.reps"
                      class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                    />
                  </label>
                  <label class="flex flex-1 flex-col gap-1">
                    <span class="text-xs font-medium text-gray-600 dark:text-gray-400">{{
                      t('fields.weightKg')
                    }}</span>
                    <input
                      type="number"
                      inputmode="decimal"
                      [formField]="editForm.weight"
                      class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                    />
                  </label>
                </div>
              } @else if (data().exercise.targetMeasurementType === 'TIME_BASED') {
                <label class="flex flex-col gap-1">
                  <span class="text-xs font-medium text-gray-600 dark:text-gray-400">{{
                    t('fields.durationSeconds')
                  }}</span>
                  <input
                    type="number"
                    inputmode="numeric"
                    [formField]="editForm.duration"
                    class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                  />
                </label>
              } @else if (data().exercise.targetMeasurementType === 'DISTANCE_BASED') {
                <label class="flex flex-col gap-1">
                  <span class="text-xs font-medium text-gray-600 dark:text-gray-400">{{
                    t('fields.distanceM')
                  }}</span>
                  <input
                    type="number"
                    inputmode="decimal"
                    [formField]="editForm.distance"
                    class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                  />
                </label>
              }
            </div>

            <button
              type="button"
              (click)="saveEdit()"
              class="w-full rounded-lg bg-blue-600 px-4 py-3 text-base font-semibold text-white hover:bg-blue-700 active:bg-blue-800 dark:bg-blue-500 dark:hover:bg-blue-600 dark:active:bg-blue-700"
            >
              {{ t('common.save') }}
            </button>
          } @else {
            @if (data().role === 'completed') {
              <div class="mb-2 text-xs text-gray-400 dark:text-gray-500">
                {{ t('user.workoutLog.actual') }}
                {{ formatActual(data().set, data().exercise.targetMeasurementType) }}
              </div>
            }
            <button
              type="button"
              (click)="togglePeek.emit(); $event.stopPropagation()"
              class="mt-3 flex w-full items-center justify-center py-1 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
            >
              <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path
                  fill-rule="evenodd"
                  d="M9.47 6.47a.75.75 0 011.06 0l4.25 4.25a.75.75 0 11-1.06 1.06L10 8.06l-3.72 3.72a.75.75 0 01-1.06-1.06l4.25-4.25z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>
          }
        </div>
      } @else {
        @switch (data().role) {
          @case ('completed') {
            <div
              class="flex cursor-pointer items-center gap-3 px-2 py-4"
              role="button"
              tabindex="0"
              (click)="togglePeek.emit()"
              (keydown.enter)="togglePeek.emit()"
            >
              <span class="w-6 text-center text-sm font-medium text-gray-400 dark:text-gray-500">
                {{ data().set.setNumber }}
              </span>
              <span class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                {{ formatSetValue(data().set, data().exercise.targetMeasurementType) }}
              </span>
              @if (data().set.status === skippedStatus) {
                <svg
                  class="h-5 w-5 text-yellow-500 dark:text-yellow-400"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z"
                  />
                </svg>
              } @else {
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
          }
          @case ('active') {
            <div
              class="flex flex-1 flex-col rounded-lg bg-blue-100 px-4 py-5 dark:bg-blue-900/30 md:flex-initial md:min-h-48"
            >
              <div class="mb-1 flex items-center justify-between">
                <span class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  {{ data().exerciseName }}
                </span>
                @if (data().isOverride) {
                  <button
                    type="button"
                    (click)="resetOverride.emit(); $event.stopPropagation()"
                    class="rounded-md p-1 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
                  >
                    <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                      <path
                        d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z"
                      />
                    </svg>
                  </button>
                } @else {
                  <button
                    type="button"
                    (click)="skip.emit(); $event.stopPropagation()"
                    class="rounded-md px-2 py-1 text-xs font-medium text-yellow-700 bg-yellow-100 hover:bg-yellow-200 dark:text-yellow-300 dark:bg-yellow-900/30 dark:hover:bg-yellow-900/50"
                  >
                    {{ t('common.skip') }}
                  </button>
                }
              </div>
              <div class="mb-4 text-sm text-gray-500 dark:text-gray-400">
                {{
                  t('user.workoutLog.setOf', {
                    current: data().set.setNumber,
                    total: data().setCount,
                  })
                }}
              </div>
              <div class="mb-4 text-xs text-gray-400 dark:text-gray-500">
                {{ t('user.workoutLog.target') }}
                {{ formatTarget(data().set, data().exercise.targetMeasurementType) }}
              </div>

              <div class="mb-5 flex flex-1 flex-col justify-center gap-3">
                @if (data().exercise.targetMeasurementType === 'REP_BASED') {
                  <div class="flex gap-3">
                    <label class="flex flex-1 flex-col gap-1">
                      <span class="text-xs font-medium text-gray-600 dark:text-gray-400">{{
                        t('fields.reps')
                      }}</span>
                      <input
                        type="number"
                        inputmode="numeric"
                        [formField]="activeForm.reps"
                        class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                      />
                    </label>
                    <label class="flex flex-1 flex-col gap-1">
                      <span class="text-xs font-medium text-gray-600 dark:text-gray-400">{{
                        t('fields.weightKg')
                      }}</span>
                      <input
                        type="number"
                        inputmode="decimal"
                        [formField]="activeForm.weight"
                        class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                      />
                    </label>
                  </div>
                } @else if (data().exercise.targetMeasurementType === 'TIME_BASED') {
                  <label class="flex flex-col gap-1">
                    <span class="text-xs font-medium text-gray-600 dark:text-gray-400">{{
                      t('fields.durationSeconds')
                    }}</span>
                    <input
                      type="number"
                      inputmode="numeric"
                      [formField]="activeForm.duration"
                      class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                    />
                  </label>
                } @else if (data().exercise.targetMeasurementType === 'DISTANCE_BASED') {
                  <label class="flex flex-col gap-1">
                    <span class="text-xs font-medium text-gray-600 dark:text-gray-400">{{
                      t('fields.distanceM')
                    }}</span>
                    <input
                      type="number"
                      inputmode="decimal"
                      [formField]="activeForm.distance"
                      class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                    />
                  </label>
                }
              </div>

              <button
                type="button"
                (click)="markDone()"
                class="w-full rounded-lg bg-blue-600 px-4 py-3 text-base font-semibold text-white hover:bg-blue-700 active:bg-blue-800 dark:bg-blue-500 dark:hover:bg-blue-600 dark:active:bg-blue-700"
              >
                {{ t('common.done') }}
              </button>
            </div>
          }
          @case ('upcoming') {
            <div
              [class]="
                'flex cursor-pointer items-center gap-3 px-2 py-4' +
                (data().isNaturalNext ? ' rounded-lg bg-blue-50/50 dark:bg-blue-950/20' : '')
              "
              role="button"
              tabindex="0"
              (click)="jumpTo.emit()"
              (keydown.enter)="jumpTo.emit()"
            >
              @if (data().isNaturalNext) {
                <svg
                  class="h-4 w-4 shrink-0 text-blue-500 dark:text-blue-400"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fill-rule="evenodd"
                    d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z"
                    clip-rule="evenodd"
                  />
                </svg>
              }
              <span
                [class]="
                  'w-6 text-center text-sm font-medium ' +
                  (data().isNaturalNext
                    ? 'text-blue-600 dark:text-blue-400'
                    : 'text-gray-700 dark:text-gray-300')
                "
              >
                {{ data().set.setNumber }}
              </span>
              <span
                [class]="
                  'flex-1 text-sm ' +
                  (data().isNaturalNext
                    ? 'text-blue-600 dark:text-blue-400'
                    : 'text-gray-700 dark:text-gray-300')
                "
              >
                {{ formatSetValue(data().set, data().exercise.targetMeasurementType) }}
              </span>
            </div>
          }
        }
      }
    </ng-container>
  `,
})
export class WorkoutLogActiveSet {
  data = input.required<ViewItemSet>();
  peeked = input(false);
  initialReps = input<number | null>(null);
  initialWeight = input<number | null>(null);
  initialDuration = input<number | null>(null);
  initialDistance = input<number | null>(null);
  done = output<SetCompletionPayload>();
  skip = output<void>();
  togglePeek = output<void>();
  save = output<SetCompletionPayload>();
  jumpTo = output<void>();
  resetOverride = output<void>();

  // Active form: tracks current input values, resets when parent provides new initial values
  activeModel = linkedSignal(() => ({
    reps: this.initialReps(),
    weight: this.initialWeight(),
    duration: this.initialDuration(),
    distance: this.initialDistance(),
  }));
  activeForm = form(this.activeModel);

  // Edit form for completed sets
  editing = signal(false);
  editModel = signal({
    reps: null as number | null,
    weight: null as number | null,
    duration: null as number | null,
    distance: null as number | null,
  });
  editForm = form(this.editModel);

  formatTarget = formatTarget;
  formatActual = formatActual;
  formatSetValue = formatSetValue;
  skippedStatus = WorkoutLogItemStatusSkipped;

  constructor() {
    effect(() => {
      if (!this.peeked()) {
        this.editing.set(false);
      }
    });
  }

  startEditing() {
    const log = this.data().set.exerciseLog;
    this.editModel.set({
      reps: log?.reps ?? null,
      weight: log?.weight ?? null,
      duration: log?.duration ?? null,
      distance: log?.distance ?? null,
    });
    this.editing.set(true);
  }

  saveEdit() {
    const m = this.editModel();
    this.save.emit({
      setId: this.data().set.id,
      actualReps: m.reps ?? undefined,
      actualWeight: m.weight ?? undefined,
      actualDuration: m.duration ?? undefined,
      actualDistance: m.distance ?? undefined,
    });
    this.editing.set(false);
  }

  markDone() {
    const m = this.activeModel();
    this.done.emit({
      setId: this.data().set.id,
      actualReps: m.reps ?? undefined,
      actualWeight: m.weight ?? undefined,
      actualDuration: m.duration ?? undefined,
      actualDistance: m.distance ?? undefined,
    });
  }
}
