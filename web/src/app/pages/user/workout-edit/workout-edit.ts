import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { form, required, applyEach, FormField } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { injectQueries } from '@tanstack/angular-query-experimental/inject-queries-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { userExerciseKeys, exerciseKeys, workoutKeys, exerciseSchemeKeys } from '$core/query-keys';
import { WorkoutSectionTypeMain, WorkoutSectionTypeSupplementary } from '$generated/user-models';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmInput } from '@spartan-ng/helm/input';
import { HlmTextarea } from '@spartan-ng/helm/textarea';

interface WorkoutExerciseModel {
  existingSchemeId: number | null;
  userExerciseId: number | null;
  measurementType: string;
  sets: number | null;
  reps: number | null;
  weight: number | null;
  restBetweenSets: number | null;
  timePerRep: number | null;
  duration: number | null;
  distance: number | null;
  targetTime: number | null;
}

interface WorkoutSectionModel {
  type: string;
  label: string;
  restBetweenExercises: number | null;
  exercises: WorkoutExerciseModel[];
}

interface WorkoutModel {
  name: string;
  notes: string;
  sections: WorkoutSectionModel[];
}

@Component({
  selector: 'app-workout-edit',
  imports: [
    PageLayout,
    FormField,
    RouterLink,
    ConfirmDialog,
    BrnSelectImports,
    HlmSelectImports,
    HlmInput,
    HlmTextarea,
    TranslocoDirective,
  ],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="isCreateMode() ? t('user.workouts.newTitle') : t('user.workouts.editTitle')"
        [isPending]="!isCreateMode() && workoutQuery.isPending()"
        [errorMessage]="
          !isCreateMode() && workoutQuery.isError() ? workoutQuery.error().message : undefined
        "
      >
        @if (isCreateMode() || workoutQuery.data()) {
          <form (submit)="onSubmit(); $event.preventDefault()" class="space-y-6">
            <!-- Basic Fields -->
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.name') }} *</label
              >
              <input id="name" [formField]="workoutForm.name" hlmInput class="mt-1" />
            </div>

            <div>
              <label
                for="notes"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.notes') }}</label
              >
              <textarea
                id="notes"
                [formField]="workoutForm.notes"
                rows="2"
                hlmTextarea
                class="mt-1"
              ></textarea>
            </div>

            <!-- Sections -->
            <div class="space-y-4">
              @for (section of workoutForm.sections; track $index; let si = $index) {
                <div class="rounded-lg border border-gray-200 p-4 dark:border-gray-700">
                  <div class="mb-3 flex items-center justify-between">
                    <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">
                      {{ t('user.workouts.sectionLabel', { n: si + 1 }) }}
                    </h3>
                    <button
                      type="button"
                      (click)="removeSection(si)"
                      class="text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                    >
                      {{ t('common.remove') }}
                    </button>
                  </div>

                  <div class="mb-3 grid grid-cols-1 gap-3 sm:grid-cols-3">
                    <div>
                      <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
                        t('fields.type')
                      }}</span>
                      <brn-select [formField]="section.type" class="mt-1" hlm>
                        <hlm-select-trigger class="w-full">
                          <hlm-select-value />
                        </hlm-select-trigger>
                        <hlm-select-content>
                          <hlm-option [value]="SECTION_TYPE_MAIN">{{
                            t('enums.workoutSectionType.main')
                          }}</hlm-option>
                          <hlm-option [value]="SECTION_TYPE_SUPPLEMENTARY">{{
                            t('enums.workoutSectionType.supplementary')
                          }}</hlm-option>
                        </hlm-select-content>
                      </brn-select>
                    </div>
                    <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                      {{ t('fields.label') }}
                      <input [formField]="section.label" hlmInput class="mt-1" />
                    </label>
                    <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                      {{ t('fields.restBetweenExercises') }}
                      <input
                        type="number"
                        [formField]="section.restBetweenExercises"
                        hlmInput
                        class="mt-1"
                      />
                    </label>
                  </div>

                  <!-- Exercises in section -->
                  <div class="space-y-3">
                    @for (exercise of section.exercises; track $index; let ei = $index) {
                      <div
                        class="rounded-md border border-gray-100 bg-gray-50 p-3 dark:border-gray-600 dark:bg-gray-800/50"
                      >
                        <div class="mb-2 flex items-center justify-between">
                          <span class="text-xs font-medium text-gray-600 dark:text-gray-400">
                            {{ t('user.workouts.exerciseLabel', { n: ei + 1 }) }}
                          </span>
                          <button
                            type="button"
                            (click)="removeExercise(si, ei)"
                            class="text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                          >
                            {{ t('common.remove') }}
                          </button>
                        </div>

                        <div class="mb-2 grid grid-cols-1 gap-2 sm:grid-cols-2">
                          <div>
                            <span
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                              >{{ t('ui.exerciseConfig.exerciseLabel') }}</span
                            >
                            <brn-select
                              [formField]="exercise.userExerciseId"
                              class="mt-1"
                              hlm
                              [placeholder]="t('common.select')"
                            >
                              <hlm-select-trigger class="w-full">
                                <hlm-select-value />
                              </hlm-select-trigger>
                              <hlm-select-content>
                                @for (ue of enrichedUserExercises(); track ue.id) {
                                  <hlm-option [value]="ue.id">{{ ue.name }}</hlm-option>
                                }
                              </hlm-select-content>
                            </brn-select>
                          </div>
                          <div>
                            <span
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                              >{{ t('fields.measurementType') }}</span
                            >
                            <brn-select [formField]="exercise.measurementType" class="mt-1" hlm>
                              <hlm-select-trigger class="w-full">
                                <hlm-select-value />
                              </hlm-select-trigger>
                              <hlm-select-content>
                                <hlm-option value="REP_BASED">{{
                                  t('enums.measurementType.REP_BASED')
                                }}</hlm-option>
                                <hlm-option value="TIME_BASED">{{
                                  t('enums.measurementType.TIME_BASED')
                                }}</hlm-option>
                                <hlm-option value="DISTANCE_BASED">{{
                                  t('enums.measurementType.DISTANCE_BASED')
                                }}</hlm-option>
                              </hlm-select-content>
                            </brn-select>
                          </div>
                        </div>

                        <!-- REP_BASED fields -->
                        @if (exercise.measurementType().value() === 'REP_BASED') {
                          <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
                            <label
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              {{ t('fields.sets') }}
                              <input
                                type="number"
                                [formField]="exercise.sets"
                                hlmInput
                                class="mt-1"
                              />
                            </label>
                            <label
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              {{ t('fields.reps') }}
                              <input
                                type="number"
                                [formField]="exercise.reps"
                                hlmInput
                                class="mt-1"
                              />
                            </label>
                            <label
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              {{ t('fields.weightKg') }}
                              <input
                                type="number"
                                [formField]="exercise.weight"
                                hlmInput
                                class="mt-1"
                              />
                            </label>
                            <label
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              {{ t('fields.restSeconds') }}
                              <input
                                type="number"
                                [formField]="exercise.restBetweenSets"
                                hlmInput
                                class="mt-1"
                              />
                            </label>
                          </div>
                        }

                        <!-- TIME_BASED fields -->
                        @if (exercise.measurementType().value() === 'TIME_BASED') {
                          <div class="grid grid-cols-2 gap-2">
                            <label
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              {{ t('fields.durationSeconds') }}
                              <input
                                type="number"
                                [formField]="exercise.duration"
                                hlmInput
                                class="mt-1"
                              />
                            </label>
                            <label
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              {{ t('fields.timePerRepSeconds') }}
                              <input
                                type="number"
                                [formField]="exercise.timePerRep"
                                hlmInput
                                class="mt-1"
                              />
                            </label>
                          </div>
                        }

                        <!-- DISTANCE_BASED fields -->
                        @if (exercise.measurementType().value() === 'DISTANCE_BASED') {
                          <div class="grid grid-cols-2 gap-2">
                            <label
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              {{ t('fields.distanceM') }}
                              <input
                                type="number"
                                [formField]="exercise.distance"
                                hlmInput
                                class="mt-1"
                              />
                            </label>
                            <label
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              {{ t('fields.targetTimeSeconds') }}
                              <input
                                type="number"
                                [formField]="exercise.targetTime"
                                hlmInput
                                class="mt-1"
                              />
                            </label>
                          </div>
                        }
                      </div>
                    }
                  </div>

                  <button
                    type="button"
                    (click)="addExercise(si)"
                    class="mt-2 text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
                  >
                    {{ t('user.workouts.addExercise') }}
                  </button>
                </div>
              }
            </div>

            <button
              type="button"
              (click)="addSection()"
              class="rounded-md border border-dashed border-gray-300 px-4 py-2 text-sm text-gray-600 hover:border-gray-400 hover:text-gray-800 dark:border-gray-600 dark:text-gray-400 dark:hover:border-gray-500 dark:hover:text-gray-300"
            >
              {{ t('user.workouts.addSection') }}
            </button>

            <!-- Actions -->
            <div class="flex gap-2">
              <button
                type="submit"
                [disabled]="!workoutForm().valid() || saveMutation.isPending()"
                class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {{ t('common.save') }}
              </button>
              <a
                routerLink="/user/workouts"
                class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
              >
                {{ t('common.cancel') }}
              </a>
              @if (!isCreateMode()) {
                <button
                  type="button"
                  (click)="showDeleteDialog = true"
                  class="ml-auto rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
                >
                  {{ t('common.delete') }}
                </button>
              }
            </div>
          </form>
        }

        <app-confirm-dialog
          [open]="showDeleteDialog"
          [title]="t('user.workouts.deleteTitle')"
          [message]="t('user.workouts.deleteMessage')"
          [isPending]="deleteMutation.isPending()"
          (confirmed)="onDelete()"
          (cancelled)="showDeleteDialog = false"
        />
      </app-page-layout>
    </ng-container>
  `,
})
export class WorkoutEdit {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = inject(QueryClient);
  private params = toSignal(this.route.paramMap);

  readonly SECTION_TYPE_MAIN = WorkoutSectionTypeMain;
  readonly SECTION_TYPE_SUPPLEMENTARY = WorkoutSectionTypeSupplementary;

  private id = computed(() => Number(this.params()?.get('id')));
  isCreateMode = computed(() => !this.params()?.get('id'));

  showDeleteDialog = false;

  model = signal<WorkoutModel>({ name: '', notes: '', sections: [] });
  workoutForm = form(this.model, (f) => {
    required(f.name);
    applyEach(f.sections, (section) => {
      applyEach(section.exercises, (exercise) => {
        required(exercise.userExerciseId);
      });
    });
  });

  // Existing workout data for edit mode
  workoutQuery = injectQuery(() => ({
    queryKey: workoutKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchWorkout(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  // User exercises for the picker dropdown
  private userExercisesQuery = injectQuery(() => ({
    queryKey: userExerciseKeys.list(),
    queryFn: () => this.userApi.fetchUserExercises(),
  }));

  private snapshotQueries = injectQueries(() => ({
    queries: (this.userExercisesQuery.data() ?? []).map((ue) => ({
      queryKey: exerciseKeys.version(ue.compendiumExerciseId, ue.compendiumVersion),
      queryFn: () =>
        this.compendiumApi.fetchExerciseVersion(ue.compendiumExerciseId, ue.compendiumVersion),
      staleTime: Infinity,
    })),
  }));

  enrichedUserExercises = computed(() => {
    const userExercises = this.userExercisesQuery.data();
    if (!userExercises) return [];
    const snapshots = this.snapshotQueries();
    return userExercises.map((ue, i) => {
      const exercise = snapshots[i]?.data()?.snapshot;
      return { id: ue.id, name: exercise?.name ?? `Exercise #${ue.id}` };
    });
  });

  // Track original scheme IDs for cleanup in edit mode
  private originalSchemeIds = signal<number[]>([]);

  saveMutation = injectMutation(() => ({
    mutationFn: () => this.saveWorkout(),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutKeys.all() });
      this.queryClient.invalidateQueries({ queryKey: exerciseSchemeKeys.all() });
      this.router.navigate(['/user/workouts']);
    },
  }));

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.deleteWorkout(),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutKeys.all() });
      this.queryClient.invalidateQueries({ queryKey: exerciseSchemeKeys.all() });
      this.router.navigate(['/user/workouts']);
    },
  }));

  constructor() {
    // Populate form in edit mode
    effect(() => {
      const data = this.workoutQuery.data();
      if (!data) return;

      this.model.set({
        name: data.name,
        notes: data.notes ?? '',
        sections: (data.sections ?? []).map((section) => ({
          type: section.type,
          label: section.label ?? '',
          restBetweenExercises: section.restBetweenExercises ?? null,
          exercises: (section.exercises ?? []).map(() => ({
            existingSchemeId: null,
            userExerciseId: null,
            measurementType: 'REP_BASED',
            sets: null,
            reps: null,
            weight: null,
            restBetweenSets: null,
            timePerRep: null,
            duration: null,
            distance: null,
            targetTime: null,
          })),
        })),
      });

      const schemeIds: number[] = [];
      for (const section of data.sections ?? []) {
        for (const ex of section.exercises ?? []) {
          schemeIds.push(ex.userExerciseSchemeId);
        }
      }
      this.originalSchemeIds.set(schemeIds);
      this.loadSchemes(data.sections ?? []);
    });
  }

  private async loadSchemes(
    sections: NonNullable<ReturnType<typeof this.workoutQuery.data>>['sections'],
  ) {
    for (let si = 0; si < sections.length; si++) {
      const section = sections[si];
      for (let ei = 0; ei < (section.exercises?.length ?? 0); ei++) {
        const ex = section.exercises[ei];
        try {
          const scheme = await this.userApi.fetchExerciseScheme(ex.userExerciseSchemeId);
          this.model.update((m) => ({
            ...m,
            sections: m.sections.map((s, sIdx) =>
              sIdx !== si
                ? s
                : {
                    ...s,
                    exercises: s.exercises.map((e, eIdx) =>
                      eIdx !== ei
                        ? e
                        : {
                            existingSchemeId: scheme.id,
                            userExerciseId: scheme.userExerciseId,
                            measurementType: scheme.measurementType || 'REP_BASED',
                            sets: scheme.sets ?? null,
                            reps: scheme.reps ?? null,
                            weight: scheme.weight ?? null,
                            restBetweenSets: scheme.restBetweenSets ?? null,
                            timePerRep: scheme.timePerRep ?? null,
                            duration: scheme.duration ?? null,
                            distance: scheme.distance ?? null,
                            targetTime: scheme.targetTime ?? null,
                          },
                    ),
                  },
            ),
          }));
        } catch {
          // scheme may have been deleted
        }
      }
    }
  }

  addSection() {
    this.model.update((m) => ({
      ...m,
      sections: [
        ...m.sections,
        {
          type: WorkoutSectionTypeMain,
          label: '',
          restBetweenExercises: null,
          exercises: [],
        },
      ],
    }));
  }

  removeSection(index: number) {
    this.model.update((m) => ({
      ...m,
      sections: m.sections.filter((_, i) => i !== index),
    }));
  }

  addExercise(sectionIndex: number) {
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, i) =>
        i !== sectionIndex
          ? s
          : {
              ...s,
              exercises: [
                ...s.exercises,
                {
                  existingSchemeId: null,
                  userExerciseId: null,
                  measurementType: 'REP_BASED',
                  sets: null,
                  reps: null,
                  weight: null,
                  restBetweenSets: null,
                  timePerRep: null,
                  duration: null,
                  distance: null,
                  targetTime: null,
                },
              ],
            },
      ),
    }));
  }

  removeExercise(sectionIndex: number, exerciseIndex: number) {
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
  }

  onSubmit() {
    if (this.workoutForm().valid()) {
      this.saveMutation.mutate();
    }
  }

  private async saveWorkout() {
    const val = this.model();

    if (this.isCreateMode()) {
      return this.createFlow(val);
    } else {
      return this.editFlow(val);
    }
  }

  private async createFlow(val: WorkoutModel) {
    const workout = await this.userApi.createWorkout({
      name: val.name,
      notes: val.notes || undefined,
    });

    await this.createSectionsAndExercises(workout.id, val.sections);
  }

  private async editFlow(val: WorkoutModel) {
    const workoutId = this.id();
    const existingWorkout = this.workoutQuery.data()!;

    await this.userApi.updateWorkout(workoutId, {
      name: val.name,
      notes: val.notes || undefined,
    });

    // Delete existing section exercises, then sections
    for (const section of existingWorkout.sections ?? []) {
      for (const ex of section.exercises ?? []) {
        await this.userApi.deleteWorkoutSectionExercise(ex.id);
      }
      await this.userApi.deleteWorkoutSection(section.id);
    }

    // Delete orphaned schemes
    const newSchemeIds = new Set<number>();
    for (const section of val.sections) {
      for (const ex of section.exercises) {
        if (ex.existingSchemeId) newSchemeIds.add(ex.existingSchemeId);
      }
    }
    for (const oldId of this.originalSchemeIds()) {
      if (!newSchemeIds.has(oldId)) {
        await this.userApi.deleteExerciseScheme(oldId);
      }
    }

    await this.createSectionsAndExercises(workoutId, val.sections);
  }

  private async createSectionsAndExercises(workoutId: number, sections: WorkoutSectionModel[]) {
    for (let si = 0; si < sections.length; si++) {
      const sectionVal = sections[si];

      // Create schemes for each exercise
      const schemeIds: number[] = [];
      for (const ex of sectionVal.exercises) {
        const schemeData: Record<string, unknown> = {
          userExerciseId: ex.userExerciseId,
          measurementType: ex.measurementType,
        };
        if (ex.measurementType === 'REP_BASED') {
          if (ex.sets != null) schemeData['sets'] = ex.sets;
          if (ex.reps != null) schemeData['reps'] = ex.reps;
          if (ex.weight != null) schemeData['weight'] = ex.weight;
          if (ex.restBetweenSets != null) schemeData['restBetweenSets'] = ex.restBetweenSets;
        } else if (ex.measurementType === 'TIME_BASED') {
          if (ex.duration != null) schemeData['duration'] = ex.duration;
          if (ex.timePerRep != null) schemeData['timePerRep'] = ex.timePerRep;
        } else if (ex.measurementType === 'DISTANCE_BASED') {
          if (ex.distance != null) schemeData['distance'] = ex.distance;
          if (ex.targetTime != null) schemeData['targetTime'] = ex.targetTime;
        }

        const scheme = await this.userApi.createExerciseScheme(schemeData);
        schemeIds.push(scheme.id);
      }

      // Create section
      const section = await this.userApi.createWorkoutSection({
        workoutId,
        type: sectionVal.type,
        label: sectionVal.label || undefined,
        position: si,
        restBetweenExercises: sectionVal.restBetweenExercises ?? undefined,
      });

      // Create section exercises
      for (let ei = 0; ei < schemeIds.length; ei++) {
        await this.userApi.createWorkoutSectionExercise({
          workoutSectionId: section.id,
          userExerciseSchemeId: schemeIds[ei],
          position: ei,
        });
      }
    }
  }

  onDelete() {
    this.deleteMutation.mutate();
  }

  private async deleteWorkout() {
    const workout = this.workoutQuery.data()!;

    // Delete section exercises and sections
    for (const section of workout.sections ?? []) {
      for (const ex of section.exercises ?? []) {
        await this.userApi.deleteWorkoutSectionExercise(ex.id);
      }
      await this.userApi.deleteWorkoutSection(section.id);
    }

    // Delete schemes
    for (const schemeId of this.originalSchemeIds()) {
      await this.userApi.deleteExerciseScheme(schemeId);
    }

    await this.userApi.deleteWorkout(this.id());
  }
}
