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
  WorkoutLogItemStatusFinished,
  WorkoutLogItemStatusSkipped,
  WorkoutLogItemStatusPartiallyFinished,
  WorkoutLogItemStatusAborted,
} from '$generated/user-models';
import { TranslocoDirective, TranslocoService } from '@jsverse/transloco';
import {
  ViewItem,
  ViewItemSet,
  SetCompletionPayload,
  AsHeaderPipe,
  AsSetPipe,
  AsBreakPipe,
} from './workout-log-view-items';
import { WorkoutLogActiveHeader } from './workout-log-active-header';
import { WorkoutLogActiveSet } from './workout-log-active-set';
import { WorkoutLogActiveBreak } from './workout-log-active-break';

function isItemTerminal(status: string): boolean {
  return (
    status === WorkoutLogItemStatusFinished ||
    status === WorkoutLogItemStatusSkipped ||
    status === WorkoutLogItemStatusPartiallyFinished ||
    status === WorkoutLogItemStatusAborted
  );
}

@Component({
  selector: 'app-workout-log-active',
  imports: [
    AsHeaderPipe,
    AsSetPipe,
    AsBreakPipe,
    WorkoutLogActiveHeader,
    WorkoutLogActiveSet,
    WorkoutLogActiveBreak,
    TranslocoDirective,
  ],
  template: `
    <ng-container *transloco="let t">
      <div class="flex min-h-[calc(100dvh-12rem)] flex-col md:min-h-0">
        @for (item of viewItems(); track item.id) {
          @if (item | asHeader; as h) {
            <app-workout-log-active-header [data]="h" />
          } @else if (item | asSet; as s) {
            <app-workout-log-active-set
              [data]="s"
              [peeked]="peekedItemId() === s.id && s.role !== 'active'"
              [(actualReps)]="actualReps"
              [(actualWeight)]="actualWeight"
              [(actualDuration)]="actualDuration"
              [(actualDistance)]="actualDistance"
              (done)="markDone()"
              (skip)="markSkipped()"
              (togglePeek)="togglePeek(s.id)"
              (save)="saveCompletedSet($event)"
              (jumpTo)="jumpToSet(s.set.id)"
              (resetOverride)="resetOverride()"
            />
          } @else if (item | asBreak; as b) {
            <app-workout-log-active-break
              [data]="b"
              [peeked]="peekedItemId() === b.id && b.role !== 'active-timer'"
              [remainingSeconds]="restSecondsRemaining()"
              (skip)="skipRest()"
              (togglePeek)="togglePeek(b.id)"
            />
          }
        }

        @if (allCompleted()) {
          <div class="py-8 text-center text-gray-500 dark:text-gray-400">
            {{ t('user.workoutLog.allCompleted') }}
          </div>
        }
      </div>
    </ng-container>
  `,
})
export class WorkoutLogActive {
  private destroyRef = inject(DestroyRef);
  private transloco = inject(TranslocoService);

  log = input.required<WorkoutLog>();
  exerciseNames = input.required<Record<number, string>>();
  setCompleted = output<SetCompletionPayload>();
  setSkipped = output<WorkoutLogExerciseSet>();

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
            exerciseName:
              names[exercise.sourceExerciseSchemeId] || this.transloco.translate('common.loading'),
          });
        }
      }
    }
    return items;
  });

  // Override: user-chosen active set (non-linear progression)
  overrideSetId = signal<number | null>(null);
  // Tracks which break is currently timing
  activeBreakId = signal<string | null>(null);

  private naturalActiveIdx = computed(() =>
    this.flatSets().findIndex((item) => !isItemTerminal(item.set.status)),
  );

  private activeIdx = computed(() => {
    const overrideId = this.overrideSetId();
    if (overrideId !== null) {
      const idx = this.flatSets().findIndex((item) => item.set.id === overrideId);
      if (idx !== -1) {
        const s = this.flatSets()[idx].set;
        if (!isItemTerminal(s.status)) {
          return idx;
        }
      }
    }
    return this.naturalActiveIdx();
  });

  viewItems = computed<ViewItem[]>(() => {
    const flat = this.flatSets();
    const activeIdx = this.activeIdx();
    const naturalIdx = this.naturalActiveIdx();
    const isOverriding = activeIdx !== naturalIdx;
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

      const isTerminal = isItemTerminal(curr.set.status);
      const hasOverride = this.overrideSetId() !== null;
      const role: 'completed' | 'active' | 'upcoming' = isTerminal
        ? 'completed'
        : i === activeIdx && (!resting || hasOverride)
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
        isNaturalNext: isOverriding && i === naturalIdx,
        isOverride: isOverriding && i === activeIdx,
      });

      if (i + 1 < flat.length) {
        const next = flat[i + 1];
        const sameExercise = curr.exercise.id === next.exercise.id;
        const breakSeconds = sameExercise
          ? (curr.set.breakAfterSeconds ?? undefined)
          : (curr.exercise.breakAfterSeconds ?? undefined);

        if (breakSeconds) {
          const breakId = sameExercise
            ? 'break-set-' + curr.set.id
            : 'break-ex-' + curr.exercise.id;
          const activeBreak = this.activeBreakId();
          const currTerminal = isItemTerminal(curr.set.status);
          const nextTerminal = isItemTerminal(next.set.status);

          let breakRole: 'elapsed' | 'active-timer' | 'upcoming';
          if (breakId === activeBreak) {
            breakRole = 'active-timer';
          } else if (currTerminal && nextTerminal) {
            breakRole = 'elapsed';
          } else {
            breakRole = 'upcoming';
          }

          items.push({
            type: 'break',
            id: breakId,
            seconds: breakSeconds,
            label: next.exerciseName,
            role: breakRole,
          });
        }
      }
    }

    return items;
  });

  allCompleted = computed(() => this.naturalActiveIdx() === -1);

  peekedItemId = signal<string | null>(null);

  constructor() {
    // Auto-reset peek when workout advances
    effect(() => {
      this.activeIdx(); // track
      this.peekedItemId.set(null);
    });

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

  togglePeek(itemId: string) {
    this.peekedItemId.set(this.peekedItemId() === itemId ? null : itemId);
  }

  markDone() {
    const items = this.viewItems();
    const activeSetIdx = items.findIndex((item) => item.type === 'set' && item.role === 'active');
    if (activeSetIdx === -1) return;
    const activeItem = items[activeSetIdx] as ViewItemSet;

    this.setCompleted.emit({
      setId: activeItem.set.id,
      actualReps: this.actualReps(),
      actualWeight: this.actualWeight(),
      actualDuration: this.actualDuration(),
      actualDistance: this.actualDistance(),
    });

    const nextItem = items[activeSetIdx + 1];
    if (nextItem?.type === 'break') {
      this.activeBreakId.set(nextItem.id);
      this.startRestTimer(nextItem.seconds, nextItem.label);
    }
    this.overrideSetId.set(null);
  }

  markSkipped() {
    const items = this.viewItems();
    const activeSetIdx = items.findIndex((item) => item.type === 'set' && item.role === 'active');
    if (activeSetIdx === -1) return;
    const activeItem = items[activeSetIdx] as ViewItemSet;

    this.setSkipped.emit(activeItem.set);
    // No rest timer for skipped sets
    this.overrideSetId.set(null);
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
        this.activeBreakId.set(null);
      } else {
        this.restState.set({ ...s, secondsRemaining: remaining });
      }
    }, 1000);
  }

  jumpToSet(setId: number) {
    this.peekedItemId.set(null);
    this.overrideSetId.set(setId);
  }

  resetOverride() {
    this.overrideSetId.set(null);
  }

  saveCompletedSet(payload: SetCompletionPayload) {
    this.setCompleted.emit(payload);
  }

  skipRest() {
    this.clearTimer();
    this.restState.set({ active: false });
    this.activeBreakId.set(null);
  }

  private clearTimer() {
    if (this.timerInterval !== null) {
      clearInterval(this.timerInterval);
      this.timerInterval = null;
    }
  }
}
