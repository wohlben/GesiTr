import {
  Component,
  input,
  output,
  computed,
  signal,
  effect,
  DestroyRef,
  inject,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { formatTarget, formatSetValue, formatBreak, formatCountdown } from '$core/format-utils';
import {
  WorkoutLog,
  WorkoutLogExercise,
  WorkoutLogExerciseSet,
  WorkoutLogSection,
  WorkoutLogStatusFinished,
  WorkoutLogStatusAborted,
} from '$generated/user-models';

type ViewItem =
  | { type: 'header'; id: string; exerciseName: string }
  | {
      type: 'set';
      id: string;
      set: WorkoutLogExerciseSet;
      exercise: WorkoutLogExercise;
      section: WorkoutLogSection;
      exerciseName: string;
      role: 'completed' | 'active' | 'upcoming';
      setCount: number;
    }
  | {
      type: 'break';
      id: string;
      seconds: number;
      label: string;
      role: 'elapsed' | 'active-timer' | 'upcoming';
    };

@Component({
  selector: 'app-workout-log-active',
  imports: [FormsModule],
  template: `
    <div class="flex min-h-[calc(100dvh-12rem)] flex-col md:min-h-0">
      @for (item of viewItems(); track item.id) {
        @switch (item.type) {
          @case ('header') {
            @let h = $any(item);
            <div class="mt-3 mb-1 text-xs font-medium text-gray-500 uppercase dark:text-gray-400">
              {{ h.exerciseName }}
            </div>
          }
          @case ('set') {
            @let s = $any(item);
            @switch (s.role) {
              @case ('completed') {
                <div class="flex items-center gap-3 px-2 py-1.5">
                  <span
                    class="w-6 text-center text-sm font-medium text-gray-400 dark:text-gray-500"
                  >
                    {{ s.set.setNumber }}
                  </span>
                  <span class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                    {{ formatSetValue(s.set, s.exercise.targetMeasurementType) }}
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
                    {{ s.exerciseName }}
                  </div>
                  <div class="mb-4 text-sm text-gray-500 dark:text-gray-400">
                    Set {{ s.set.setNumber }} of {{ s.setCount }}
                  </div>
                  <div class="mb-4 text-xs text-gray-400 dark:text-gray-500">
                    Target: {{ formatTarget(s.set, s.exercise.targetMeasurementType) }}
                  </div>

                  <div class="mb-5 flex flex-1 flex-col justify-center gap-3">
                    @if (s.exercise.targetMeasurementType === 'REP_BASED') {
                      <div class="flex gap-3">
                        <label class="flex flex-1 flex-col gap-1">
                          <span class="text-xs font-medium text-gray-600 dark:text-gray-400"
                            >Reps</span
                          >
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
                    } @else if (s.exercise.targetMeasurementType === 'TIME_BASED') {
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
                    } @else if (s.exercise.targetMeasurementType === 'DISTANCE_BASED') {
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
                    (click)="markDone()"
                    class="w-full rounded-lg bg-blue-600 px-4 py-3 text-base font-semibold text-white hover:bg-blue-700 active:bg-blue-800 dark:bg-blue-500 dark:hover:bg-blue-600 dark:active:bg-blue-700"
                  >
                    Done
                  </button>
                </div>
              }
              @case ('upcoming') {
                <div class="flex items-center gap-3 px-2 py-1.5">
                  <span
                    class="w-6 text-center text-sm font-medium text-gray-400 dark:text-gray-500"
                  >
                    {{ s.set.setNumber }}
                  </span>
                  <span class="flex-1 text-sm text-gray-400 dark:text-gray-500">
                    {{ formatSetValue(s.set, s.exercise.targetMeasurementType) }}
                  </span>
                </div>
              }
            }
          }
          @case ('break') {
            @let b = $any(item);
            @if (b.role === 'active-timer') {
              <div
                class="my-3 flex flex-col items-center justify-center rounded-xl border-2 border-amber-500 bg-amber-50/50 p-5 dark:border-amber-400 dark:bg-amber-950/20"
              >
                <div
                  class="mb-2 text-xs font-semibold tracking-wider text-amber-600 uppercase dark:text-amber-400"
                >
                  Rest
                </div>
                <div
                  class="mb-3 text-5xl font-bold tabular-nums text-amber-700 dark:text-amber-300"
                >
                  {{ formatCountdown(restSecondsRemaining()) }}
                </div>
                <div class="mb-5 text-sm text-gray-500 dark:text-gray-400">Next: {{ b.label }}</div>
                <button
                  type="button"
                  (click)="skipRest()"
                  class="w-full rounded-lg border-2 border-amber-500 bg-transparent px-4 py-3 text-base font-semibold text-amber-700 hover:bg-amber-100 active:bg-amber-200 dark:border-amber-400 dark:text-amber-300 dark:hover:bg-amber-900/30 dark:active:bg-amber-900/50"
                >
                  Skip
                </button>
              </div>
            } @else {
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
                  {{ formatBreak(b.seconds) }}
                </div>
              </div>
            }
          }
        }
      }

      @if (allCompleted()) {
        <div class="py-8 text-center text-gray-500 dark:text-gray-400">All sets completed!</div>
      }
    </div>
  `,
})
export class WorkoutLogActive {
  private destroyRef = inject(DestroyRef);

  log = input.required<WorkoutLog>();
  exerciseNames = input.required<Record<number, string>>();
  setToggled = output<WorkoutLogExerciseSet>();

  // Actual value signals for the active set inputs
  actualReps = signal<number | undefined>(undefined);
  actualWeight = signal<number | undefined>(undefined);
  actualDuration = signal<number | undefined>(undefined);
  actualDistance = signal<number | undefined>(undefined);

  // Rest timer state
  restState = signal<
    | { active: false }
    | { active: true; secondsRemaining: number; totalSeconds: number; nextLabel: string }
  >({ active: false });
  isResting = computed(() => this.restState().active);
  restSecondsRemaining = computed(() => {
    const s = this.restState();
    return s.active ? s.secondsRemaining : 0;
  });

  private timerInterval: ReturnType<typeof setInterval> | null = null;

  private flatSets = computed(() => {
    const log = this.log();
    const names = this.exerciseNames();
    const items: {
      set: WorkoutLogExerciseSet;
      exercise: WorkoutLogExercise;
      section: WorkoutLogSection;
      exerciseName: string;
    }[] = [];
    for (const section of log.sections ?? []) {
      for (const exercise of section.exercises ?? []) {
        for (const set of exercise.sets ?? []) {
          items.push({
            set,
            exercise,
            section,
            exerciseName: names[exercise.sourceExerciseSchemeId] || 'Loading...',
          });
        }
      }
    }
    return items;
  });

  private activeIdx = computed(() =>
    this.flatSets().findIndex(
      (item) =>
        item.set.status !== WorkoutLogStatusFinished && item.set.status !== WorkoutLogStatusAborted,
    ),
  );

  viewItems = computed<ViewItem[]>(() => {
    const flat = this.flatSets();
    const activeIdx = this.activeIdx();
    const resting = this.isResting();
    const items: ViewItem[] = [];

    for (let i = 0; i < flat.length; i++) {
      const curr = flat[i];
      const prev = i > 0 ? flat[i - 1] : null;

      if (!prev || prev.exercise.id !== curr.exercise.id) {
        items.push({
          type: 'header',
          id: 'header-' + curr.exercise.id,
          exerciseName: curr.exerciseName,
        });
      }

      const role: 'completed' | 'active' | 'upcoming' =
        activeIdx === -1 || i < activeIdx
          ? 'completed'
          : i === activeIdx && !resting
            ? 'active'
            : 'upcoming';

      items.push({
        type: 'set',
        id: 'set-' + curr.set.id,
        set: curr.set,
        exercise: curr.exercise,
        section: curr.section,
        exerciseName: curr.exerciseName,
        role,
        setCount: curr.exercise.sets?.length ?? 0,
      });

      if (i + 1 < flat.length) {
        const next = flat[i + 1];
        const sameExercise = curr.exercise.id === next.exercise.id;
        const breakSeconds = sameExercise
          ? (curr.set.breakAfterSeconds ?? undefined)
          : (curr.exercise.breakAfterSeconds ?? undefined);

        if (breakSeconds) {
          let breakRole: 'elapsed' | 'active-timer' | 'upcoming';
          if (activeIdx === -1 || i + 1 < activeIdx) {
            breakRole = 'elapsed';
          } else if (i === activeIdx - 1 && resting) {
            breakRole = 'active-timer';
          } else {
            breakRole = 'upcoming';
          }

          items.push({
            type: 'break',
            id: sameExercise ? 'break-set-' + curr.set.id : 'break-ex-' + curr.exercise.id,
            seconds: breakSeconds,
            label: next.exerciseName,
            role: breakRole,
          });
        }
      }
    }

    return items;
  });

  allCompleted = computed(() => this.activeIdx() === -1);

  constructor() {
    effect(() => {
      const activeItem = this.viewItems().find(
        (item): item is Extract<ViewItem, { type: 'set' }> =>
          item.type === 'set' && item.role === 'active',
      );
      if (activeItem) {
        this.actualReps.set(activeItem.set.targetReps ?? undefined);
        this.actualWeight.set(activeItem.set.targetWeight ?? undefined);
        this.actualDuration.set(activeItem.set.targetDuration ?? undefined);
        this.actualDistance.set(activeItem.set.targetDistance ?? undefined);
      }
    });

    this.destroyRef.onDestroy(() => this.clearTimer());
  }

  markDone() {
    const items = this.viewItems();
    const activeSetIdx = items.findIndex((item) => item.type === 'set' && item.role === 'active');
    if (activeSetIdx === -1) return;
    const activeItem = items[activeSetIdx] as Extract<ViewItem, { type: 'set' }>;

    this.setToggled.emit({
      ...activeItem.set,
      actualReps: this.actualReps(),
      actualWeight: this.actualWeight(),
      actualDuration: this.actualDuration(),
      actualDistance: this.actualDistance(),
    });

    const nextItem = items[activeSetIdx + 1];
    if (nextItem?.type === 'break') {
      this.startRestTimer(nextItem.seconds, nextItem.label);
    }
  }

  startRestTimer(seconds: number, nextLabel: string) {
    this.clearTimer();
    this.restState.set({
      active: true,
      secondsRemaining: seconds,
      totalSeconds: seconds,
      nextLabel,
    });
    this.timerInterval = setInterval(() => {
      const s = this.restState();
      if (!s.active) {
        this.clearTimer();
        return;
      }
      const remaining = s.secondsRemaining - 1;
      if (remaining <= 0) {
        this.clearTimer();
        this.restState.set({ active: false });
      } else {
        this.restState.set({ ...s, secondsRemaining: remaining });
      }
    }, 1000);
  }

  skipRest() {
    this.clearTimer();
    this.restState.set({ active: false });
  }

  private clearTimer() {
    if (this.timerInterval !== null) {
      clearInterval(this.timerInterval);
      this.timerInterval = null;
    }
  }

  formatTarget = formatTarget;
  formatSetValue = formatSetValue;
  formatBreak = formatBreak;
  formatCountdown = formatCountdown;
}
