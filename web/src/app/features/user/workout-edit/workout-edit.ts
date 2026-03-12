import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { ReactiveFormsModule, FormGroup, FormControl, FormArray, Validators } from '@angular/forms';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { injectQueries } from '@tanstack/angular-query-experimental/inject-queries-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { userExerciseKeys, exerciseKeys, workoutKeys, exerciseSchemeKeys } from '$core/query-keys';
import { WorkoutSectionTypeMain, WorkoutSectionTypeSupplementary } from '$generated/user-models';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';

type ExerciseFormGroup = FormGroup<{
  existingSchemeId: FormControl<number | null>;
  userExerciseId: FormControl<number | null>;
  measurementType: FormControl<string>;
  sets: FormControl<number | null>;
  reps: FormControl<number | null>;
  weight: FormControl<number | null>;
  restBetweenSets: FormControl<number | null>;
  timePerRep: FormControl<number | null>;
  duration: FormControl<number | null>;
  distance: FormControl<number | null>;
  targetTime: FormControl<number | null>;
}>;

type SectionFormGroup = FormGroup<{
  type: FormControl<string>;
  label: FormControl<string>;
  restBetweenExercises: FormControl<number | null>;
  exercises: FormArray<ExerciseFormGroup>;
}>;

@Component({
  selector: 'app-workout-edit',
  imports: [PageLayout, ReactiveFormsModule, RouterLink, ConfirmDialog],
  template: `
    <app-page-layout
      [header]="isCreateMode() ? 'New Workout' : 'Edit Workout'"
      [isPending]="!isCreateMode() && workoutQuery.isPending()"
      [errorMessage]="
        !isCreateMode() && workoutQuery.isError() ? workoutQuery.error().message : undefined
      "
    >
      @if (isCreateMode() || workoutQuery.data()) {
        <form [formGroup]="form" (ngSubmit)="onSubmit()" class="space-y-6">
          <!-- Basic Fields -->
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Name *</label
            >
            <input
              id="name"
              formControlName="name"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            />
          </div>

          <div>
            <label for="notes" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Notes</label
            >
            <textarea
              id="notes"
              formControlName="notes"
              rows="2"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            ></textarea>
          </div>

          <!-- Sections -->
          <div formArrayName="sections" class="space-y-4">
            @for (section of sectionsArray.controls; track $index; let si = $index) {
              <div
                [formGroupName]="si"
                class="rounded-lg border border-gray-200 p-4 dark:border-gray-700"
              >
                <div class="mb-3 flex items-center justify-between">
                  <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">
                    Section {{ si + 1 }}
                  </h3>
                  <button
                    type="button"
                    (click)="removeSection(si)"
                    class="text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                  >
                    Remove
                  </button>
                </div>

                <div class="mb-3 grid grid-cols-1 gap-3 sm:grid-cols-3">
                  <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                    Type
                    <select
                      formControlName="type"
                      class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                    >
                      <option [value]="SECTION_TYPE_MAIN">Main</option>
                      <option [value]="SECTION_TYPE_SUPPLEMENTARY">Supplementary</option>
                    </select>
                  </label>
                  <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                    Label
                    <input
                      formControlName="label"
                      class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                    />
                  </label>
                  <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                    Rest Between Exercises (s)
                    <input
                      type="number"
                      formControlName="restBetweenExercises"
                      class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                    />
                  </label>
                </div>

                <!-- Exercises in section -->
                <div formArrayName="exercises" class="space-y-3">
                  @for (ex of getExercisesArray(si).controls; track $index; let ei = $index) {
                    <div
                      [formGroupName]="ei"
                      class="rounded-md border border-gray-100 bg-gray-50 p-3 dark:border-gray-600 dark:bg-gray-800/50"
                    >
                      <div class="mb-2 flex items-center justify-between">
                        <span class="text-xs font-medium text-gray-600 dark:text-gray-400">
                          Exercise {{ ei + 1 }}
                        </span>
                        <button
                          type="button"
                          (click)="removeExercise(si, ei)"
                          class="text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                        >
                          Remove
                        </button>
                      </div>

                      <div class="mb-2 grid grid-cols-1 gap-2 sm:grid-cols-2">
                        <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                          Exercise *
                          <select
                            formControlName="userExerciseId"
                            class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                          >
                            <option [ngValue]="null">-- Select --</option>
                            @for (ue of enrichedUserExercises(); track ue.id) {
                              <option [ngValue]="ue.id">{{ ue.name }}</option>
                            }
                          </select>
                        </label>
                        <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                          Measurement Type
                          <select
                            formControlName="measurementType"
                            class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                          >
                            <option value="REP_BASED">Rep Based</option>
                            <option value="TIME_BASED">Time Based</option>
                            <option value="DISTANCE_BASED">Distance Based</option>
                          </select>
                        </label>
                      </div>

                      <!-- REP_BASED fields -->
                      @if (ex.get('measurementType')?.value === 'REP_BASED') {
                        <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
                          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                            Sets
                            <input
                              type="number"
                              formControlName="sets"
                              class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            />
                          </label>
                          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                            Reps
                            <input
                              type="number"
                              formControlName="reps"
                              class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            />
                          </label>
                          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                            Weight (kg)
                            <input
                              type="number"
                              formControlName="weight"
                              class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            />
                          </label>
                          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                            Rest (s)
                            <input
                              type="number"
                              formControlName="restBetweenSets"
                              class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            />
                          </label>
                        </div>
                      }

                      <!-- TIME_BASED fields -->
                      @if (ex.get('measurementType')?.value === 'TIME_BASED') {
                        <div class="grid grid-cols-2 gap-2">
                          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                            Duration (s)
                            <input
                              type="number"
                              formControlName="duration"
                              class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            />
                          </label>
                          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                            Time Per Rep (s)
                            <input
                              type="number"
                              formControlName="timePerRep"
                              class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            />
                          </label>
                        </div>
                      }

                      <!-- DISTANCE_BASED fields -->
                      @if (ex.get('measurementType')?.value === 'DISTANCE_BASED') {
                        <div class="grid grid-cols-2 gap-2">
                          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                            Distance (m)
                            <input
                              type="number"
                              formControlName="distance"
                              class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            />
                          </label>
                          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                            Target Time (s)
                            <input
                              type="number"
                              formControlName="targetTime"
                              class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
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
                  + Add Exercise
                </button>
              </div>
            }
          </div>

          <button
            type="button"
            (click)="addSection()"
            class="rounded-md border border-dashed border-gray-300 px-4 py-2 text-sm text-gray-600 hover:border-gray-400 hover:text-gray-800 dark:border-gray-600 dark:text-gray-400 dark:hover:border-gray-500 dark:hover:text-gray-300"
          >
            + Add Section
          </button>

          <!-- Actions -->
          <div class="flex gap-2">
            <button
              type="submit"
              [disabled]="form.invalid || saveMutation.isPending()"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              Save
            </button>
            <a
              routerLink="/user/workouts"
              class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
            >
              Cancel
            </a>
            @if (!isCreateMode()) {
              <button
                type="button"
                (click)="showDeleteDialog = true"
                class="ml-auto rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
              >
                Delete
              </button>
            }
          </div>
        </form>
      }

      <app-confirm-dialog
        [open]="showDeleteDialog"
        title="Delete Workout"
        message="Are you sure you want to delete this workout? This cannot be undone."
        [isPending]="deleteMutation.isPending()"
        (confirmed)="onDelete()"
        (cancelled)="showDeleteDialog = false"
      />
    </app-page-layout>
  `,
})
export class WorkoutEdit {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = injectQueryClient();
  private params = toSignal(this.route.paramMap);

  readonly SECTION_TYPE_MAIN = WorkoutSectionTypeMain;
  readonly SECTION_TYPE_SUPPLEMENTARY = WorkoutSectionTypeSupplementary;

  private id = computed(() => Number(this.params()?.get('id')));
  isCreateMode = computed(() => !this.params()?.get('id'));

  showDeleteDialog = false;

  form = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    notes: new FormControl('', { nonNullable: true }),
    sections: new FormArray<SectionFormGroup>([]),
  });

  get sectionsArray() {
    return this.form.controls.sections;
  }

  getExercisesArray(sectionIndex: number) {
    return this.sectionsArray.at(sectionIndex).controls.exercises;
  }

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

      this.form.patchValue({
        name: data.name,
        notes: data.notes ?? '',
      });

      // Clear existing sections and rebuild
      this.sectionsArray.clear();
      const schemeIds: number[] = [];

      for (const section of data.sections ?? []) {
        const exerciseGroups: ExerciseFormGroup[] = [];
        for (const ex of section.exercises ?? []) {
          schemeIds.push(ex.userExerciseSchemeId);
          exerciseGroups.push(this.createExerciseGroup());
        }
        const sectionGroup = this.createSectionGroup();
        sectionGroup.patchValue({
          type: section.type,
          label: section.label ?? '',
          restBetweenExercises: section.restBetweenExercises ?? null,
        });
        for (const eg of exerciseGroups) {
          sectionGroup.controls.exercises.push(eg);
        }
        this.sectionsArray.push(sectionGroup);
      }

      this.originalSchemeIds.set(schemeIds);

      // Fetch schemes and populate exercise fields
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
          const exerciseFg = this.getExercisesArray(si).at(ei);
          exerciseFg.patchValue({
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
          });
        } catch {
          // scheme may have been deleted
        }
      }
    }
  }

  addSection() {
    this.sectionsArray.push(this.createSectionGroup());
  }

  removeSection(index: number) {
    this.sectionsArray.removeAt(index);
  }

  addExercise(sectionIndex: number) {
    this.getExercisesArray(sectionIndex).push(this.createExerciseGroup());
  }

  removeExercise(sectionIndex: number, exerciseIndex: number) {
    this.getExercisesArray(sectionIndex).removeAt(exerciseIndex);
  }

  private createSectionGroup(): SectionFormGroup {
    return new FormGroup({
      type: new FormControl(WorkoutSectionTypeMain, { nonNullable: true }),
      label: new FormControl('', { nonNullable: true }),
      restBetweenExercises: new FormControl<number | null>(null),
      exercises: new FormArray<ExerciseFormGroup>([]),
    });
  }

  private createExerciseGroup(): ExerciseFormGroup {
    return new FormGroup({
      existingSchemeId: new FormControl<number | null>(null),
      userExerciseId: new FormControl<number | null>(null, [Validators.required]),
      measurementType: new FormControl('REP_BASED', { nonNullable: true }),
      sets: new FormControl<number | null>(null),
      reps: new FormControl<number | null>(null),
      weight: new FormControl<number | null>(null),
      restBetweenSets: new FormControl<number | null>(null),
      timePerRep: new FormControl<number | null>(null),
      duration: new FormControl<number | null>(null),
      distance: new FormControl<number | null>(null),
      targetTime: new FormControl<number | null>(null),
    });
  }

  onSubmit() {
    if (this.form.valid) {
      this.saveMutation.mutate();
    }
  }

  private async saveWorkout() {
    const val = this.form.getRawValue();

    if (this.isCreateMode()) {
      return this.createFlow(val);
    } else {
      return this.editFlow(val);
    }
  }

  private async createFlow(val: ReturnType<typeof this.form.getRawValue>) {
    const workout = await this.userApi.createWorkout({
      name: val.name,
      notes: val.notes || undefined,
    });

    await this.createSectionsAndExercises(workout.id, val.sections);
  }

  private async editFlow(val: ReturnType<typeof this.form.getRawValue>) {
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

  private async createSectionsAndExercises(
    workoutId: number,
    sections: ReturnType<typeof this.form.getRawValue>['sections'],
  ) {
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
