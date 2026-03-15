import { Component, input, model, output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { formatTarget, formatSetValue } from '$core/format-utils';
import { ViewItemSet } from './workout-log-view-items';

@Component({
  selector: 'app-workout-log-active-set',
  imports: [FormsModule],
  template: `
    @switch (data().role) {
      @case ('completed') {
        <div class="flex items-center gap-3 px-2 py-1.5">
          <span class="w-6 text-center text-sm font-medium text-gray-400 dark:text-gray-500">
            {{ data().set.setNumber }}
          </span>
          <span class="flex-1 text-sm text-gray-500 dark:text-gray-400">
            {{ formatSetValue(data().set, data().exercise.targetMeasurementType) }}
          </span>
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
        </div>
      }
      @case ('active') {
        <div
          class="my-3 flex flex-1 flex-col rounded-xl border-2 border-blue-500 bg-blue-50/50 p-5 dark:border-blue-400 dark:bg-blue-950/20 md:flex-initial md:min-h-48"
        >
          <div class="mb-1 text-lg font-semibold text-gray-900 dark:text-gray-100">
            {{ data().exerciseName }}
          </div>
          <div class="mb-4 text-sm text-gray-500 dark:text-gray-400">
            Set {{ data().set.setNumber }} of {{ data().setCount }}
          </div>
          <div class="mb-4 text-xs text-gray-400 dark:text-gray-500">
            Target: {{ formatTarget(data().set, data().exercise.targetMeasurementType) }}
          </div>

          <div class="mb-5 flex flex-1 flex-col justify-center gap-3">
            @if (data().exercise.targetMeasurementType === 'REP_BASED') {
              <div class="flex gap-3">
                <label class="flex flex-1 flex-col gap-1">
                  <span class="text-xs font-medium text-gray-600 dark:text-gray-400">Reps</span>
                  <input
                    type="number"
                    inputmode="numeric"
                    [ngModel]="actualReps()"
                    (ngModelChange)="actualReps.set($event)"
                    class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                  />
                </label>
                <label class="flex flex-1 flex-col gap-1">
                  <span class="text-xs font-medium text-gray-600 dark:text-gray-400"
                    >Weight (kg)</span
                  >
                  <input
                    type="number"
                    inputmode="decimal"
                    [ngModel]="actualWeight()"
                    (ngModelChange)="actualWeight.set($event)"
                    class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                  />
                </label>
              </div>
            } @else if (data().exercise.targetMeasurementType === 'TIME_BASED') {
              <label class="flex flex-col gap-1">
                <span class="text-xs font-medium text-gray-600 dark:text-gray-400"
                  >Duration (s)</span
                >
                <input
                  type="number"
                  inputmode="numeric"
                  [ngModel]="actualDuration()"
                  (ngModelChange)="actualDuration.set($event)"
                  class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                />
              </label>
            } @else if (data().exercise.targetMeasurementType === 'DISTANCE_BASED') {
              <label class="flex flex-col gap-1">
                <span class="text-xs font-medium text-gray-600 dark:text-gray-400"
                  >Distance (m)</span
                >
                <input
                  type="number"
                  inputmode="decimal"
                  [ngModel]="actualDistance()"
                  (ngModelChange)="actualDistance.set($event)"
                  class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-center text-lg font-semibold text-gray-900 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                />
              </label>
            }
          </div>

          <button
            type="button"
            (click)="done.emit()"
            class="w-full rounded-lg bg-blue-600 px-4 py-3 text-base font-semibold text-white hover:bg-blue-700 active:bg-blue-800 dark:bg-blue-500 dark:hover:bg-blue-600 dark:active:bg-blue-700"
          >
            Done
          </button>
        </div>
      }
      @case ('upcoming') {
        <div class="flex items-center gap-3 px-2 py-1.5">
          <span class="w-6 text-center text-sm font-medium text-gray-400 dark:text-gray-500">
            {{ data().set.setNumber }}
          </span>
          <span class="flex-1 text-sm text-gray-400 dark:text-gray-500">
            {{ formatSetValue(data().set, data().exercise.targetMeasurementType) }}
          </span>
        </div>
      }
    }
  `,
})
export class WorkoutLogActiveSet {
  data = input.required<ViewItemSet>();
  actualReps = model<number | undefined>(undefined);
  actualWeight = model<number | undefined>(undefined);
  actualDuration = model<number | undefined>(undefined);
  actualDistance = model<number | undefined>(undefined);
  done = output<void>();

  formatTarget = formatTarget;
  formatSetValue = formatSetValue;
}
