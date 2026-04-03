import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { form, FormField } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CdkDragDrop, CdkDrag, CdkDropList } from '@angular/cdk/drag-drop';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { formatBreak } from '$core/format-utils';
import { workoutKeys, workoutLogKeys } from '$core/query-keys';
import { Workout, WorkoutSectionTypeMain, WorkoutLog } from '$generated/user-models';
import { ExerciseScheme } from '$generated/models';
import { PageLayout } from '../../../layout/page-layout';
import { HlmInput } from '@spartan-ng/helm/input';
import { HlmTextarea } from '@spartan-ng/helm/textarea';
import { WorkoutStartStore, SetPreview } from './workout-start.store';
import { AddExerciseDialog } from './add-exercise-dialog';
import { WorkoutStartSection } from './workout-start-section';
import { StartModel, StartExerciseModel, PendingGroupModel } from './workout-start.models';

@Component({
  selector: 'app-workout-start',
  imports: [
    PageLayout,
    FormField,
    RouterLink,
    AddExerciseDialog,
    CdkDropList,
    CdkDrag,
    HlmInput,
    HlmTextarea,
    TranslocoDirective,
    WorkoutStartSection,
  ],
  providers: [WorkoutStartStore],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('user.workoutStart.title')"
        [isPending]="isPending()"
        [errorMessage]="workoutQuery.isError() ? workoutQuery.error().message : undefined"
      >
        @if (workoutQuery.data() && currentLogId()) {
          <form (submit)="onSubmit(); $event.preventDefault()" class="space-y-6">
            <!-- Basic Fields -->
            <div>
              <label
                for="name"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.name') }}</label
              >
              <input
                hlmInput
                id="name"
                [formField]="startForm.name"
                (change)="onLogChange()"
                class="mt-1"
              />
            </div>

            <div>
              <label
                for="notes"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.notes') }}</label
              >
              <textarea
                hlmTextarea
                id="notes"
                [formField]="startForm.notes"
                (change)="onLogChange()"
                rows="2"
                class="mt-1"
              ></textarea>
            </div>

            <!-- Sections -->
            <div
              cdkDropList
              [cdkDropListData]="startForm.sections"
              [cdkDropListDisabled]="isReadonly()"
              (cdkDropListDropped)="onSectionDrop($event)"
              class="space-y-4"
            >
              @for (section of startForm.sections; track $index; let si = $index) {
                <app-workout-start-section
                  cdkDrag
                  [cdkDragDisabled]="isReadonly()"
                  [section]="section"
                  [sectionIndex]="si"
                  [exerciseDisplayMap]="store.exerciseDisplay()"
                  [readonly]="isReadonly()"
                  (removed)="removeSection(si)"
                  (sectionChanged)="onSectionChange(si)"
                  (exerciseRemoved)="removeExercise(si, $event.exerciseIndex)"
                  (exerciseChanged)="onExerciseChange(si, $event.exerciseIndex)"
                  (setChanged)="onSetChange(si, $event.exerciseIndex, $event.setIndex)"
                  (exerciseDropped)="onExerciseDropFromSection(si, $event)"
                  (addExerciseRequested)="openAddExerciseDialog(si)"
                  (groupExercisePicked)="onGroupExercisePickedFromSection(si, $event)"
                />
              }
            </div>

            <!-- Add Section button -->
            @if (!isReadonly()) {
              <button
                type="button"
                (click)="addSection()"
                class="w-full rounded-md border border-dashed border-gray-300 px-3 py-1.5 text-xs text-gray-400 hover:border-gray-400 hover:text-gray-500 dark:border-gray-600 dark:text-gray-500 dark:hover:border-gray-500 dark:hover:text-gray-400"
              >
                {{ t('user.workouts.addSection') }}
              </button>
            }

            <!-- Actions -->
            <div class="flex gap-2">
              <button
                type="submit"
                [disabled]="startMutation.isPending() || hasPendingGroups()"
                class="rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
              >
                @if (startMutation.isPending()) {
                  {{ t('common.starting') }}
                } @else {
                  {{ t('user.workoutStart.startWorkout') }}
                }
              </button>
              <a
                routerLink="/compendium/workouts"
                class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
              >
                {{ t('common.cancel') }}
              </a>
            </div>
          </form>
        }

        <app-add-exercise-dialog
          [open]="addDialogOpen()"
          [sectionId]="addDialogSectionId()"
          [logId]="currentLogId() ?? 0"
          [exerciseCount]="addDialogExerciseCount()"
          [preselectedExerciseId]="addDialogPreselectedExerciseId()"
          (exerciseAdded)="onExerciseAdded($event)"
          (cancelled)="addDialogOpen.set(false)"
        />
      </app-page-layout>
    </ng-container>
  `,
  styles: `
    .cdk-drag-preview {
      box-sizing: border-box;
      border-radius: 8px;
      border: 1px solid #d1d5db;
      background: white;
      padding: 12px 16px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      font-size: 13px;
      font-weight: 600;
    }

    .cdk-drag-placeholder {
      background: #e0f2fe;
      border: 2px dashed #7dd3fc;
      border-radius: 8px;
      min-height: 3rem;
    }

    .cdk-drag-animating {
      transition: transform 200ms cubic-bezier(0, 0, 0.2, 1);
    }

    .cdk-drop-list-dragging .cdk-drag:not(.cdk-drag-placeholder) {
      transition: transform 200ms cubic-bezier(0, 0, 0.2, 1);
    }

    :host-context(.dark) .cdk-drag-preview {
      background: #1f2937;
      border-color: #4b5563;
      color: #f3f4f6;
    }

    :host-context(.dark) .cdk-drag-placeholder {
      background: rgba(14, 116, 144, 0.2);
      border-color: #0e7490;
    }
  `,
})
export class WorkoutStart {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = inject(QueryClient);
  private params = toSignal(this.route.paramMap);

  readonly store = inject(WorkoutStartStore);

  id = computed(() => Number(this.params()?.get('id')));
  private initialized = false;
  currentLogId = signal<number | null>(null);
  private creating = signal(false);

  isPending = computed(
    () =>
      this.workoutQuery.isPending() ||
      this.permissionsQuery.isPending() ||
      this.planningLogQuery.isPending() ||
      this.creating(),
  );

  isReadonly = computed(() => {
    const perms = this.permissionsQuery.data()?.permissions;
    if (!perms) return false;
    return !perms.includes('MODIFY');
  });

  hasPendingGroups = computed(() => this.model().sections.some((s) => s.pendingGroups.length > 0));

  // Add exercise dialog state
  addDialogOpen = signal(false);
  addDialogSectionIndex = signal(0);
  addDialogSectionId = signal(0);
  addDialogExerciseCount = signal(0);
  addDialogPreselectedExerciseId = signal<number | null>(null);
  addDialogPendingGroupIndex = signal<number | null>(null);

  model = signal<StartModel>({ name: '', notes: '', sections: [] });
  startForm = form(this.model);

  workoutQuery = injectQuery(() => ({
    queryKey: workoutKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchWorkout(this.id()),
    enabled: !!this.id(),
  }));

  permissionsQuery = injectQuery(() => ({
    queryKey: workoutKeys.permissions(this.id()),
    queryFn: () => this.userApi.fetchWorkoutPermissions(this.id()),
    enabled: !!this.id(),
  }));

  planningLogQuery = injectQuery(() => ({
    queryKey: workoutLogKeys.list({ workoutId: this.id(), status: 'planning' }),
    queryFn: () => this.userApi.fetchWorkoutLogs({ workoutId: this.id(), status: 'planning' }),
    enabled: !!this.id(),
  }));

  startMutation = injectMutation(() => ({
    mutationFn: () => this.startWorkout(),
  }));

  constructor() {
    // Fetch-or-create effect
    effect(() => {
      const workout = this.workoutQuery.data();
      const planningLogs = this.planningLogQuery.data();
      if (!workout || planningLogs === undefined || this.initialized) return;
      this.initialized = true;

      if (planningLogs.length > 0) {
        this.populateFromLog(planningLogs[0]);
      } else {
        this.createPlanningLog(workout);
      }
    });
  }

  // CDK Drag & Drop
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  onSectionDrop(event: CdkDragDrop<any>) {
    if (event.previousIndex === event.currentIndex) return;
    this.model.update((m) => {
      const sections = [...m.sections];
      const [moved] = sections.splice(event.previousIndex, 1);
      sections.splice(event.currentIndex, 0, moved);
      return { ...m, sections };
    });
    this.persistSectionPositions();
  }

  onExerciseDropFromSection(
    sectionIndex: number,
    event: { previousIndex: number; currentIndex: number },
  ) {
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, i) =>
        i !== sectionIndex
          ? s
          : {
              ...s,
              exercises: (() => {
                const exercises = [...s.exercises];
                const [moved] = exercises.splice(event.previousIndex, 1);
                exercises.splice(event.currentIndex, 0, moved);
                return exercises;
              })(),
            },
      ),
    }));
    this.persistExercisePositions(sectionIndex);
  }

  async removeSection(index: number) {
    const m = this.model();
    const sectionId = m.sections[index].id;
    this.model.update((m) => ({
      ...m,
      sections: m.sections.filter((_, i) => i !== index),
    }));
    if (sectionId) {
      try {
        await this.userApi.deleteWorkoutLogSection(sectionId);
      } catch (err) {
        console.error('Failed to delete section:', err);
      }
    }
  }

  async removeExercise(sectionIndex: number, exerciseIndex: number) {
    const m = this.model();
    const exerciseId = m.sections[sectionIndex].exercises[exerciseIndex].id;
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, i) =>
        i !== sectionIndex
          ? s
          : {
              ...s,
              exercises: s.exercises.filter((_, j) => j !== exerciseIndex),
            },
      ),
    }));
    if (exerciseId) {
      try {
        await this.userApi.deleteWorkoutLogExercise(exerciseId);
      } catch (err) {
        console.error('Failed to delete exercise:', err);
      }
    }
  }

  async addSection() {
    const logId = this.currentLogId();
    if (!logId) return;
    try {
      const section = await this.userApi.createWorkoutLogSection({
        workoutLogId: logId,
        type: WorkoutSectionTypeMain,
        position: this.model().sections.length,
      });
      this.model.update((m) => ({
        ...m,
        sections: [
          ...m.sections,
          {
            id: section.id,
            type: WorkoutSectionTypeMain,
            label: '',
            exercises: [],
            pendingGroups: [],
          },
        ],
      }));
    } catch (err) {
      console.error('Failed to add section:', err);
    }
  }

  onGroupExercisePickedFromSection(
    sectionIndex: number,
    event: { groupIndex: number; exerciseId: number },
  ) {
    const m = this.model();
    const section = m.sections[sectionIndex];
    const sectionId = section.id;
    if (!sectionId) return;

    this.addDialogSectionIndex.set(sectionIndex);
    this.addDialogSectionId.set(sectionId);
    this.addDialogExerciseCount.set(section.exercises.length);
    this.addDialogPreselectedExerciseId.set(event.exerciseId);
    this.addDialogPendingGroupIndex.set(event.groupIndex);
    this.addDialogOpen.set(true);
  }

  openAddExerciseDialog(sectionIndex: number) {
    const m = this.model();
    const sectionId = m.sections[sectionIndex].id;
    if (!sectionId) return;
    this.addDialogSectionIndex.set(sectionIndex);
    this.addDialogSectionId.set(sectionId);
    this.addDialogExerciseCount.set(m.sections[sectionIndex].exercises.length);
    this.addDialogPreselectedExerciseId.set(null);
    this.addDialogPendingGroupIndex.set(null);
    this.addDialogOpen.set(true);
  }

  onExerciseAdded(event: {
    exerciseLogId: number;
    exerciseName: string;
    scheme: ExerciseScheme;
    exercise: {
      id: number;
      sourceExerciseSchemeId: number;
      sets: {
        id: number;
        setNumber: number;
        targetReps?: number;
        targetWeight?: number;
        targetDuration?: number;
        targetDistance?: number;
        targetTime?: number;
        breakAfterSeconds?: number;
      }[];
    };
  }) {
    const si = this.addDialogSectionIndex();

    const newExercise: StartExerciseModel = {
      id: event.exercise.id,
      sourceExerciseSchemeId: event.exercise.sourceExerciseSchemeId,
      breakAfterSeconds: null,
      sets: (event.exercise.sets ?? []).map((set) => ({
        id: set.id,
        targetReps: set.targetReps ?? null,
        targetWeight: set.targetWeight ?? null,
        targetDuration: set.targetDuration ?? null,
        targetDistance: set.targetDistance ?? null,
        targetTime: set.targetTime ?? null,
        restAfterSeconds: set.breakAfterSeconds ?? null,
      })),
    };

    const pendingGroupIndex = this.addDialogPendingGroupIndex();
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, i) =>
        i !== si
          ? s
          : {
              ...s,
              exercises: [...s.exercises, newExercise],
              pendingGroups:
                pendingGroupIndex != null
                  ? s.pendingGroups.filter((_, j) => j !== pendingGroupIndex)
                  : s.pendingGroups,
            },
      ),
    }));

    // Update store display
    const numSets = event.scheme.sets ?? 0;
    const sets: SetPreview[] = [];
    for (let i = 1; i <= numSets; i++) {
      sets.push({
        setNumber: i,
        targetReps: event.scheme.reps,
        targetWeight: event.scheme.weight,
        targetDuration: event.scheme.duration,
        targetDistance: event.scheme.distance,
        targetTime: event.scheme.targetTime,
        restAfterSeconds: i < numSets ? (event.scheme.restBetweenSets ?? null) : null,
      });
    }
    this.store.addExerciseDisplay(event.exercise.id, event.exerciseName, event.scheme, sets);

    this.addDialogOpen.set(false);
  }

  private async persistSectionPositions() {
    const m = this.model();
    const updates = [];
    for (let i = 0; i < m.sections.length; i++) {
      const sectionId = m.sections[i].id;
      if (sectionId) {
        updates.push(this.userApi.updateWorkoutLogSection(sectionId, { position: i }));
      }
    }
    try {
      await Promise.all(updates);
    } catch (err) {
      console.error('Failed to persist section positions:', err);
    }
  }

  private async persistExercisePositions(sectionIndex: number) {
    const m = this.model();
    const exercises = m.sections[sectionIndex].exercises;
    const updates = [];
    for (let i = 0; i < exercises.length; i++) {
      const exerciseId = exercises[i].id;
      if (exerciseId) {
        updates.push(this.userApi.updateWorkoutLogExercise(exerciseId, { position: i }));
      }
    }
    try {
      await Promise.all(updates);
    } catch (err) {
      console.error('Failed to persist exercise positions:', err);
    }
  }

  formatBreak = formatBreak;

  onSubmit() {
    this.startMutation.mutate();
  }

  async onLogChange() {
    const logId = this.currentLogId();
    if (!logId) return;
    const val = this.model();
    try {
      await this.userApi.updateWorkoutLog(logId, {
        name: val.name,
        notes: val.notes || undefined,
      });
    } catch (err) {
      console.error('Failed to save log changes:', err);
    }
  }

  async onSectionChange(si: number) {
    const m = this.model();
    const section = m.sections[si];
    const sectionId = section.id;
    if (!sectionId) return;
    try {
      await this.userApi.updateWorkoutLogSection(sectionId, {
        type: section.type,
        label: section.label || undefined,
      });
    } catch (err) {
      console.error('Failed to save section changes:', err);
    }
  }

  async onExerciseChange(si: number, ei: number) {
    const m = this.model();
    const exercise = m.sections[si].exercises[ei];
    const exerciseId = exercise.id;
    if (!exerciseId) return;
    try {
      await this.userApi.updateWorkoutLogExercise(exerciseId, {
        breakAfterSeconds: exercise.breakAfterSeconds ?? undefined,
      });
    } catch (err) {
      console.error('Failed to save exercise changes:', err);
    }
  }

  async onSetChange(si: number, ei: number, setIdx: number) {
    const m = this.model();
    const set = m.sections[si].exercises[ei].sets[setIdx];
    const setId = set.id;
    if (!setId) return;

    try {
      await this.userApi.updateWorkoutLogExerciseSet(setId, {
        targetReps: set.targetReps ?? undefined,
        targetWeight: set.targetWeight ?? undefined,
        targetDuration: set.targetDuration ?? undefined,
        targetDistance: set.targetDistance ?? undefined,
        targetTime: set.targetTime ?? undefined,
        breakAfterSeconds: set.restAfterSeconds ?? undefined,
      });
    } catch (err) {
      console.error('Failed to save set changes:', err);
    }
  }

  private populateFromLog(log: WorkoutLog) {
    this.currentLogId.set(log.id);

    this.model.set({
      name: log.name,
      notes: log.notes ?? '',
      sections: (log.sections ?? []).map((section) => ({
        id: section.id,
        type: section.type,
        label: section.label ?? '',
        exercises: (section.exercises ?? []).map((ex) => ({
          id: ex.id,
          sourceExerciseSchemeId: ex.sourceExerciseSchemeId,
          breakAfterSeconds: ex.breakAfterSeconds ?? null,
          sets: (ex.sets ?? []).map((set) => ({
            id: set.id,
            targetReps: set.targetReps ?? null,
            targetWeight: set.targetWeight ?? null,
            targetDuration: set.targetDuration ?? null,
            targetDistance: set.targetDistance ?? null,
            targetTime: set.targetTime ?? null,
            restAfterSeconds: set.breakAfterSeconds ?? null,
          })),
        })),
        pendingGroups: [],
      })),
    });

    this.store.loadExerciseDisplayFromLog(log.sections ?? []);
  }

  private async createPlanningLog(workout: Workout) {
    this.creating.set(true);

    // Pre-fetch user-specific schemes for all exercise items
    const exerciseItems = (workout.sections ?? []).flatMap((s) =>
      (s.items ?? []).filter((i) => i.type === 'exercise'),
    );
    const schemeResults = await Promise.all(
      exerciseItems.map((item) =>
        this.userApi.fetchExerciseSchemes({ workoutSectionItemId: item.id }),
      ),
    );
    const userSchemeByItemId = new Map(
      schemeResults
        .flat()
        .filter((s) => s.workoutSectionItemId != null)
        .map((s) => [s.workoutSectionItemId!, s]),
    );

    const log = await this.userApi.createWorkoutLog({
      name: workout.name,
      notes: workout.notes || undefined,
      workoutId: workout.id,
    });

    for (let si = 0; si < (workout.sections ?? []).length; si++) {
      const templateSection = workout.sections[si];

      const section = await this.userApi.createWorkoutLogSection({
        workoutLogId: log.id,
        type: templateSection.type,
        label: templateSection.label || undefined,
        position: si,
        restBetweenExercises: templateSection.restBetweenExercises ?? undefined,
      });

      for (let ei = 0; ei < (templateSection.items ?? []).length; ei++) {
        const templateItem = templateSection.items[ei];
        // Only create log exercises for exercise-type items;
        // exercise_group items are resolved on the start page
        if (templateItem.type === 'exercise') {
          // Use user-specific scheme (looked up by item ID), falling back to item's embedded scheme
          const myScheme = userSchemeByItemId.get(templateItem.id);
          const schemeId = myScheme?.id ?? templateItem.exerciseSchemeId;
          if (schemeId) {
            await this.userApi.createWorkoutLogExercise({
              workoutLogSectionId: section.id,
              sourceExerciseSchemeId: schemeId,
              position: ei,
            });
          }
        }
      }
    }

    // Collect pending exercise groups and resolve their members
    const pendingGroupsBySection = new Map<number, PendingGroupModel[]>();
    const allExercises = await this.compendiumApi.fetchExercises({ limit: 200 });
    const exerciseNameMap = new Map(
      allExercises.items.map((e) => [e.id, e.names?.[0]?.name ?? '']),
    );

    for (let si = 0; si < (workout.sections ?? []).length; si++) {
      const templateSection = workout.sections[si];
      const groups: PendingGroupModel[] = [];

      for (let ei = 0; ei < (templateSection.items ?? []).length; ei++) {
        const templateItem = templateSection.items[ei];
        if (templateItem.type === 'exercise_group' && templateItem.exerciseGroupId) {
          const [group, members] = await Promise.all([
            this.compendiumApi.fetchExerciseGroup(templateItem.exerciseGroupId),
            this.compendiumApi.fetchExerciseGroupMembers({
              groupId: templateItem.exerciseGroupId,
            }),
          ]);
          groups.push({
            groupId: templateItem.exerciseGroupId,
            groupName: group.name ?? '',
            members: members.map((m) => ({
              id: m.exerciseId,
              name: exerciseNameMap.get(m.exerciseId) ?? `Exercise #${m.exerciseId}`,
            })),
            position: ei,
          });
        }
      }
      if (groups.length > 0) {
        pendingGroupsBySection.set(si, groups);
      }
    }

    // Fetch the full log with nested structure
    const fullLog = await this.userApi.fetchWorkoutLog(log.id);
    this.creating.set(false);
    this.populateFromLog(fullLog);

    // Inject pending groups into the model sections
    if (pendingGroupsBySection.size > 0) {
      this.model.update((m) => ({
        ...m,
        sections: m.sections.map((s, i) => ({
          ...s,
          pendingGroups: pendingGroupsBySection.get(i) ?? [],
        })),
      }));
    }

    // Invalidate planning log query so it reflects the new log
    this.queryClient.invalidateQueries({
      queryKey: workoutLogKeys.list({ workoutId: workout.id, status: 'planning' }),
    });
  }

  private async startWorkout() {
    await this.userApi.startWorkoutLog(this.currentLogId()!);
    this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.all() });
    this.router.navigate(['/user/workout-logs', this.currentLogId()]);
  }
}
