import { Component, input, computed, signal, effect } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HlmInput } from '@spartan-ng/helm/input';

export interface ExerciseRunnerSet {
  targetReps: number | null;
  targetWeight: number | null;
  targetDuration: number | null;
  targetDistance: number | null;
  targetTime: number | null;
  restAfterSeconds: number | null;
}

@Component({
  selector: 'app-exercise-runner',
  imports: [FormsModule, HlmInput],
  template: `
    <div class="rounded-md border border-gray-200 dark:border-gray-600">
      <!-- Exercise header -->
      <div
        class="flex items-center justify-between border-b border-gray-100 px-3 py-2 dark:border-gray-700"
      >
        <span class="text-sm font-semibold text-gray-900 dark:text-gray-100">
          {{ exerciseName() }}
        </span>
        <span class="text-xs text-gray-500 dark:text-gray-400">
          {{ measurementLabel() }}
        </span>
      </div>

      @if (sets().length) {
        <div class="px-3 py-2">
          <!-- Column headers -->
          <div
            class="mb-1 grid text-left text-xs text-gray-500 uppercase dark:text-gray-400"
            [class]="
              measurementType() === 'REP_BASED'
                ? 'grid-cols-[2rem_5rem_6rem]'
                : 'grid-cols-[2rem_6rem]'
            "
          >
            <span>Set</span>
            @if (measurementType() === 'REP_BASED') {
              <span>Reps</span>
              <span>Weight</span>
            }
            @if (measurementType() === 'TIME_BASED') {
              <span>Duration</span>
            }
            @if (measurementType() === 'DISTANCE_BASED') {
              <span>Distance</span>
            }
          </div>

          @for (set of sets(); track $index; let idx = $index; let last = $last) {
            <!-- Set row -->
            <div
              class="grid items-center py-1.5"
              [class]="
                measurementType() === 'REP_BASED'
                  ? 'grid-cols-[2rem_5rem_6rem]'
                  : 'grid-cols-[2rem_6rem]'
              "
            >
              <span class="text-sm font-medium text-gray-900 dark:text-gray-100">
                {{ idx + 1 }}
              </span>
              @if (measurementType() === 'REP_BASED') {
                <div>
                  <input
                    hlmInput
                    type="number"
                    [ngModel]="set.targetReps"
                    (ngModelChange)="updateSet(idx, 'targetReps', $event)"
                  />
                </div>
                <div>
                  <input
                    hlmInput
                    type="number"
                    step="0.5"
                    [ngModel]="set.targetWeight"
                    (ngModelChange)="updateSet(idx, 'targetWeight', $event)"
                  />
                </div>
              }
              @if (measurementType() === 'TIME_BASED') {
                <div>
                  <input
                    hlmInput
                    type="number"
                    [ngModel]="set.targetDuration"
                    (ngModelChange)="updateSet(idx, 'targetDuration', $event)"
                  />
                </div>
              }
              @if (measurementType() === 'DISTANCE_BASED') {
                <div>
                  <input
                    hlmInput
                    type="number"
                    step="0.1"
                    [ngModel]="set.targetDistance"
                    (ngModelChange)="updateSet(idx, 'targetDistance', $event)"
                  />
                </div>
              }
            </div>

            <!-- Rest between sets -->
            @if (!last && set.restAfterSeconds !== null) {
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
                    [ngModel]="set.restAfterSeconds"
                    (ngModelChange)="updateSet(idx, 'restAfterSeconds', $event)"
                    class="w-12 border-0 bg-transparent p-0 text-center text-xs text-gray-400 focus:ring-0 dark:text-gray-500"
                  />
                  <span>s</span>
                </div>
              </div>
            }
          }
        </div>
      }
    </div>
  `,
})
export class ExerciseRunner {
  exerciseName = input.required<string>();
  measurementType = input.required<string>();
  setCount = input.required<number>();
  defaultReps = input<number | null>(null);
  defaultWeight = input<number | null>(null);
  defaultDuration = input<number | null>(null);
  defaultDistance = input<number | null>(null);
  defaultRest = input<number | null>(null);

  sets = signal<ExerciseRunnerSet[]>([]);

  constructor() {
    // Auto-rebuild sets whenever inputs change
    effect(
      () => {
        const count = this.setCount();
        const mt = this.measurementType();
        if (!count) return;
        const sets: ExerciseRunnerSet[] = [];
        for (let i = 0; i < count; i++) {
          sets.push({
            targetReps: mt === 'REP_BASED' ? this.defaultReps() : null,
            targetWeight: mt === 'REP_BASED' ? this.defaultWeight() : null,
            targetDuration: mt === 'TIME_BASED' ? this.defaultDuration() : null,
            targetDistance: mt === 'DISTANCE_BASED' ? this.defaultDistance() : null,
            targetTime: null,
            restAfterSeconds: i < count - 1 ? (this.defaultRest() ?? null) : null,
          });
        }
        this.sets.set(sets);
      },
      { allowSignalWrites: true },
    );
  }

  measurementLabel = computed(() => {
    const mt = this.measurementType();
    if (mt === 'REP_BASED') return 'Rep Based';
    if (mt === 'TIME_BASED') return 'Time Based';
    if (mt === 'DISTANCE_BASED') return 'Distance Based';
    return mt;
  });

  /** Called by the parent to rebuild sets from Phase 1 config. */
  rebuildSets(count: number, defaults: Partial<ExerciseRunnerSet>) {
    const sets: ExerciseRunnerSet[] = [];
    for (let i = 0; i < count; i++) {
      sets.push({
        targetReps: defaults.targetReps ?? null,
        targetWeight: defaults.targetWeight ?? null,
        targetDuration: defaults.targetDuration ?? null,
        targetDistance: defaults.targetDistance ?? null,
        targetTime: defaults.targetTime ?? null,
        restAfterSeconds: i < count - 1 ? (defaults.restAfterSeconds ?? null) : null,
      });
    }
    this.sets.set(sets);
  }

  updateSet(idx: number, field: keyof ExerciseRunnerSet, value: number | null) {
    const current = [...this.sets()];
    current[idx] = { ...current[idx], [field]: value };
    this.sets.set(current);
  }

  reset() {
    this.sets.set([]);
  }
}
