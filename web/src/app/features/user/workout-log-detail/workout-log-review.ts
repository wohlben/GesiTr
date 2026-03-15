import { Component, input, signal } from '@angular/core';
import { formatBreak, formatSetValue } from '$core/format-utils';
import { WorkoutLog } from '$generated/user-models';

@Component({
  selector: 'app-workout-log-review',
  styles: `
    .hide-breaks .break-indicator {
      display: none;
    }
  `,
  template: `
    <div class="space-y-6" [class.hide-breaks]="!showBreaks()">
      <div class="flex justify-end">
        <button
          type="button"
          (click)="showBreaks.set(!showBreaks())"
          class="flex items-center gap-1 rounded-md px-2 py-1 text-xs"
          [class]="
            showBreaks()
              ? 'bg-gray-200 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
              : 'text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300'
          "
        >
          <svg class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
            <path
              fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z"
              clip-rule="evenodd"
            />
          </svg>
          Rest times
        </button>
      </div>
      @for (section of log().sections; track section.id) {
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
                      {{ exerciseNames()[exercise.sourceExerciseSchemeId] || 'Loading...' }}
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

                <!-- Sets -->
                <div>
                  @for (set of exercise.sets; track set.id; let lastSet = $last) {
                    <div class="flex items-center gap-3 px-2 py-1.5">
                      <span
                        class="w-6 text-center text-sm font-medium text-gray-400 dark:text-gray-500"
                      >
                        {{ set.setNumber }}
                      </span>
                      <span class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                        {{ formatSetValue(set, exercise.targetMeasurementType) }}
                      </span>
                      @if (set.status === 'finished') {
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
                    @if (
                      !lastSet &&
                      set.breakAfterSeconds !== null &&
                      set.breakAfterSeconds !== undefined
                    ) {
                      <div class="break-indicator relative flex items-center justify-center py-0.5">
                        <div
                          class="absolute inset-x-0 top-1/2 border-t border-dashed border-gray-200 dark:border-gray-700"
                        ></div>
                        <div
                          class="relative z-10 flex items-center gap-1 rounded-full bg-gray-50 px-2 py-0.5 text-xs text-gray-400 dark:bg-gray-900 dark:text-gray-500"
                        >
                          <svg class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
                            <path
                              fill-rule="evenodd"
                              d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z"
                              clip-rule="evenodd"
                            />
                          </svg>
                          {{ formatBreak(set.breakAfterSeconds) }}
                        </div>
                      </div>
                    }
                  }
                </div>

                <!-- Break after exercise -->
                @if (exercise.breakAfterSeconds && !lastEx) {
                  <div
                    class="break-indicator relative mt-2 flex items-center justify-center py-0.5"
                  >
                    <div
                      class="absolute inset-x-0 top-1/2 border-t border-dashed border-gray-200 dark:border-gray-700"
                    ></div>
                    <div
                      class="relative z-10 flex items-center gap-1 rounded-full bg-gray-50 px-2 py-0.5 text-xs text-gray-400 dark:bg-gray-900 dark:text-gray-500"
                    >
                      <svg class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
                        <path
                          fill-rule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z"
                          clip-rule="evenodd"
                        />
                      </svg>
                      {{ formatBreak(exercise.breakAfterSeconds) }}
                    </div>
                  </div>
                }
              </div>
            }
          </div>
        </div>
      }
    </div>
  `,
})
export class WorkoutLogReview {
  log = input.required<WorkoutLog>();
  exerciseNames = input.required<Record<number, string>>();

  showBreaks = signal(false);

  formatSetValue = formatSetValue;
  formatBreak = formatBreak;
}
