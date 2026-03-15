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
import {
  WorkoutLog,
  WorkoutLogExercise,
  WorkoutLogExerciseSet,
  WorkoutLogSection,
  WorkoutLogStatusFinished,
  WorkoutLogStatusAborted,
} from '$generated/user-models';
import {
  ViewItem,
  ViewItemSet,
  AsHeaderPipe,
  AsSetPipe,
  AsBreakPipe,
} from './workout-log-view-items';
import { WorkoutLogActiveHeader } from './workout-log-active-header';
import { WorkoutLogActiveSet } from './workout-log-active-set';
import { WorkoutLogActiveBreak } from './workout-log-active-break';

@Component({
  selector: 'app-workout-log-active',
  imports: [
    AsHeaderPipe,
    AsSetPipe,
    AsBreakPipe,
    WorkoutLogActiveHeader,
    WorkoutLogActiveSet,
    WorkoutLogActiveBreak,
  ],
  template: `
    <div class="flex min-h-[calc(100dvh-12rem)] flex-col md:min-h-0">
      @for (item of viewItems(); track item.id) {
        @if (item | asHeader; as h) {
          <app-workout-log-active-header [data]="h" />
        } @else if (item | asSet; as s) {
          <app-workout-log-active-set
            [data]="s"
            [(actualReps)]="actualReps"
            [(actualWeight)]="actualWeight"
            [(actualDuration)]="actualDuration"
            [(actualDistance)]="actualDistance"
            (done)="markDone()"
          />
        } @else if (item | asBreak; as b) {
          <app-workout-log-active-break
            [data]="b"
            [remainingSeconds]="restSecondsRemaining()"
            (skip)="skipRest()"
          />
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
        (item): item is ViewItemSet => item.type === 'set' && item.role === 'active',
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
    const activeItem = items[activeSetIdx] as ViewItemSet;

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
}
