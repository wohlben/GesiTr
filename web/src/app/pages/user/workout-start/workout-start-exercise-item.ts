import { Component, input, output } from '@angular/core';
import { FormField } from '@angular/forms/signals';
import { CdkDragHandle } from '@angular/cdk/drag-drop';
import { TranslocoDirective } from '@jsverse/transloco';
import { HlmInput } from '@spartan-ng/helm/input';
import { ExerciseDisplayInfo } from './workout-start.store';

@Component({
  selector: 'app-workout-start-exercise-item',
  imports: [FormField, CdkDragHandle, HlmInput, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      @let info = displayInfo();
      <div class="rounded-md border border-gray-200 dark:border-gray-600">
        <!-- Exercise header -->
        <div
          class="flex items-center justify-between border-b border-gray-100 px-3 py-2 dark:border-gray-700"
        >
          <div class="flex items-center gap-2 text-sm text-gray-900 dark:text-gray-100">
            @if (!readonly()) {
              <!-- Exercise drag handle -->
              <div
                cdkDragHandle
                class="flex cursor-grab flex-col gap-0.5 text-gray-400 active:cursor-grabbing dark:text-gray-500"
              >
                <div class="flex gap-0.5">
                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                </div>
                <div class="flex gap-0.5">
                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                </div>
                <div class="flex gap-0.5">
                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                </div>
              </div>
            }
            <div>
              <span class="font-semibold">{{ info?.name ?? t('common.loading') }}</span>
              <span class="ml-2 text-gray-500 dark:text-gray-400">{{ info?.summary }}</span>
            </div>
          </div>
          @if (!readonly()) {
            <button
              type="button"
              (click)="removed.emit()"
              class="text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
            >
              {{ t('common.remove') }}
            </button>
          }
        </div>

        <!-- Editable sets -->
        @if (exercise().sets.length) {
          <div class="px-3 py-2">
            <!-- Header -->
            <div
              class="mb-1 grid text-left text-xs text-gray-500 uppercase dark:text-gray-400"
              [class]="
                info?.measurementType === 'REP_BASED'
                  ? 'grid-cols-[2rem_5rem_6rem]'
                  : 'grid-cols-[2rem_6rem]'
              "
            >
              <span>{{ t('fields.set') }}</span>
              @if (info?.measurementType === 'REP_BASED') {
                <span>{{ t('fields.reps') }}</span>
                <span>{{ t('fields.weight') }}</span>
              }
              @if (info?.measurementType === 'TIME_BASED') {
                <span>{{ t('fields.duration') }}</span>
              }
              @if (info?.measurementType === 'DISTANCE_BASED') {
                <span>{{ t('fields.distance') }}</span>
              }
            </div>

            @for (set of exercise().sets; track $index; let setIdx = $index; let lastSet = $last) {
              <!-- Set row -->
              <div
                class="grid items-center py-1.5"
                [class]="
                  info?.measurementType === 'REP_BASED'
                    ? 'grid-cols-[2rem_5rem_6rem]'
                    : 'grid-cols-[2rem_6rem]'
                "
              >
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{
                  setIdx + 1
                }}</span>
                @if (info?.measurementType === 'REP_BASED') {
                  <div>
                    <input
                      hlmInput
                      type="number"
                      [formField]="set.targetReps"
                      data-field="targetReps"
                      (change)="setChanged.emit({ setIndex: setIdx })"
                      class="mt-1"
                    />
                  </div>
                  <div>
                    <input
                      hlmInput
                      type="number"
                      [formField]="set.targetWeight"
                      data-field="targetWeight"
                      (change)="setChanged.emit({ setIndex: setIdx })"
                      class="mt-1"
                      step="0.5"
                    />
                  </div>
                }
                @if (info?.measurementType === 'TIME_BASED') {
                  <div>
                    <input
                      hlmInput
                      type="number"
                      [formField]="set.targetDuration"
                      (change)="setChanged.emit({ setIndex: setIdx })"
                      class="mt-1"
                    />
                  </div>
                }
                @if (info?.measurementType === 'DISTANCE_BASED') {
                  <div>
                    <input
                      hlmInput
                      type="number"
                      [formField]="set.targetDistance"
                      (change)="setChanged.emit({ setIndex: setIdx })"
                      class="mt-1"
                      step="0.1"
                    />
                  </div>
                }
              </div>

              <!-- Rest between sets: line with centered badge -->
              @if (!lastSet && set.restAfterSeconds().value() !== null) {
                <div class="relative flex items-center justify-center py-0.5">
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
                    <input
                      type="number"
                      [formField]="set.restAfterSeconds"
                      data-field="restAfterSeconds"
                      (change)="setChanged.emit({ setIndex: setIdx })"
                      class="w-12 border-0 bg-transparent p-0 text-center text-xs text-gray-400 focus:ring-0 dark:text-gray-500"
                    />
                    <span>{{ t('common.unitSeconds') }}</span>
                  </div>
                </div>
              }
            }
          </div>
        }
      </div>

      <!-- Break after exercise (editable, not shown after last) -->
      @if (!isLast()) {
        <div class="relative flex items-center justify-center py-3">
          <div
            class="absolute inset-x-0 top-1/2 border-t border-gray-200 dark:border-gray-700"
          ></div>
          <div
            class="relative z-10 flex items-center gap-1.5 rounded-full bg-white px-3 py-1 text-xs text-gray-500 shadow-sm ring-1 ring-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:ring-gray-600"
          >
            <svg class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z"
                clip-rule="evenodd"
              />
            </svg>
            <input
              type="number"
              [formField]="exercise().breakAfterSeconds"
              data-field="breakAfterSeconds"
              (change)="exerciseChanged.emit()"
              class="w-12 border-0 bg-transparent p-0 text-center text-xs text-gray-500 focus:ring-0 dark:text-gray-400"
            />
            <span>{{ t('common.unitSeconds') }}</span>
          </div>
        </div>
      }
    </ng-container>
  `,
})
export class WorkoutStartExerciseItem {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  exercise = input.required<any>();
  displayInfo = input<ExerciseDisplayInfo | undefined>();
  isLast = input(false);
  readonly = input(false);

  removed = output<void>();
  exerciseChanged = output<void>();
  setChanged = output<{ setIndex: number }>();
}
