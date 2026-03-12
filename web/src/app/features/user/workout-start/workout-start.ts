import { Component, inject, computed, effect } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { ReactiveFormsModule, FormGroup, FormControl, FormArray } from '@angular/forms';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutKeys, workoutLogKeys } from '$core/query-keys';
import { WorkoutSectionTypeMain, WorkoutSectionTypeSupplementary } from '$generated/user-models';
import { PageLayout } from '../../../layout/page-layout';
import { WorkoutStartStore } from './workout-start.store';

type ExerciseFormGroup = FormGroup<{
  userExerciseSchemeId: FormControl<number>;
}>;

type SectionFormGroup = FormGroup<{
  type: FormControl<string>;
  label: FormControl<string>;
  restBetweenExercises: FormControl<number | null>;
  exercises: FormArray<ExerciseFormGroup>;
}>;

@Component({
  selector: 'app-workout-start',
  imports: [PageLayout, ReactiveFormsModule, RouterLink],
  providers: [WorkoutStartStore],
  template: `
    <app-page-layout
      header="Start Workout"
      [isPending]="workoutQuery.isPending()"
      [errorMessage]="workoutQuery.isError() ? workoutQuery.error().message : undefined"
    >
      @if (workoutQuery.data()) {
        <form [formGroup]="form" (ngSubmit)="onSubmit()" class="space-y-6">
          <!-- Basic Fields -->
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Name</label
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

          <div>
            <label for="date" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Date</label
            >
            <input
              id="date"
              type="date"
              formControlName="date"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            />
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
                <div formArrayName="exercises" class="space-y-2">
                  @for (ex of getExercisesArray(si).controls; track $index; let ei = $index) {
                    <div
                      [formGroupName]="ei"
                      class="flex items-center justify-between rounded-md border border-gray-100 bg-gray-50 px-3 py-2 dark:border-gray-600 dark:bg-gray-800/50"
                    >
                      <div class="text-sm text-gray-900 dark:text-gray-100">
                        <span class="font-medium">{{
                          store.exerciseDisplay()[ex.get('userExerciseSchemeId')!.value]?.name ??
                            'Loading...'
                        }}</span>
                        <span class="ml-2 text-gray-500 dark:text-gray-400">
                          {{
                            store.exerciseDisplay()[ex.get('userExerciseSchemeId')!.value]?.summary
                          }}
                        </span>
                      </div>
                      <button
                        type="button"
                        (click)="removeExercise(si, ei)"
                        class="text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                      >
                        Remove
                      </button>
                    </div>
                  }
                </div>
              </div>
            }
          </div>

          <!-- Actions -->
          <div class="flex gap-2">
            <button
              type="submit"
              [disabled]="startMutation.isPending()"
              class="rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
            >
              @if (startMutation.isPending()) {
                Starting...
              } @else {
                Start Workout
              }
            </button>
            <a
              routerLink="/user/workouts"
              class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
            >
              Cancel
            </a>
          </div>
        </form>
      }
    </app-page-layout>
  `,
})
export class WorkoutStart {
  private userApi = inject(UserApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = injectQueryClient();
  private params = toSignal(this.route.paramMap);

  readonly store = inject(WorkoutStartStore);
  readonly SECTION_TYPE_MAIN = WorkoutSectionTypeMain;
  readonly SECTION_TYPE_SUPPLEMENTARY = WorkoutSectionTypeSupplementary;

  private id = computed(() => Number(this.params()?.get('id')));
  private initialized = false;

  form = new FormGroup({
    name: new FormControl('', { nonNullable: true }),
    notes: new FormControl('', { nonNullable: true }),
    date: new FormControl('', { nonNullable: true }),
    sections: new FormArray<SectionFormGroup>([]),
  });

  get sectionsArray() {
    return this.form.controls.sections;
  }

  getExercisesArray(sectionIndex: number) {
    return this.sectionsArray.at(sectionIndex).controls.exercises;
  }

  workoutQuery = injectQuery(() => ({
    queryKey: workoutKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchWorkout(this.id()),
    enabled: !!this.id(),
  }));

  startMutation = injectMutation(() => ({
    mutationFn: () => this.startWorkout(),
  }));

  constructor() {
    effect(() => {
      const data = this.workoutQuery.data();
      if (!data || this.initialized) return;
      this.initialized = true;

      const today = new Date().toISOString().split('T')[0];
      this.form.patchValue({
        name: data.name,
        notes: data.notes ?? '',
        date: today,
      });

      this.sectionsArray.clear();

      for (const section of data.sections ?? []) {
        const sectionGroup = this.createSectionGroup();
        sectionGroup.patchValue({
          type: section.type,
          label: section.label ?? '',
          restBetweenExercises: section.restBetweenExercises ?? null,
        });

        for (const ex of section.exercises ?? []) {
          sectionGroup.controls.exercises.push(this.createExerciseGroup(ex.userExerciseSchemeId));
        }

        this.sectionsArray.push(sectionGroup);
      }

      this.store.loadExerciseDisplay(data.sections ?? []);
    });
  }

  removeSection(index: number) {
    this.sectionsArray.removeAt(index);
  }

  removeExercise(sectionIndex: number, exerciseIndex: number) {
    this.getExercisesArray(sectionIndex).removeAt(exerciseIndex);
  }

  onSubmit() {
    this.startMutation.mutate();
  }

  private createSectionGroup(): SectionFormGroup {
    return new FormGroup({
      type: new FormControl(WorkoutSectionTypeMain, { nonNullable: true }),
      label: new FormControl('', { nonNullable: true }),
      restBetweenExercises: new FormControl<number | null>(null),
      exercises: new FormArray<ExerciseFormGroup>([]),
    });
  }

  private createExerciseGroup(schemeId: number): ExerciseFormGroup {
    return new FormGroup({
      userExerciseSchemeId: new FormControl(schemeId, { nonNullable: true }),
    });
  }

  private async startWorkout() {
    const val = this.form.getRawValue();
    const workoutId = this.id();

    const log = await this.userApi.createWorkoutLog({
      name: val.name,
      notes: val.notes || undefined,
      date: new Date(val.date).toISOString(),
      workoutId,
    });

    for (let si = 0; si < val.sections.length; si++) {
      const sectionVal = val.sections[si];

      const section = await this.userApi.createWorkoutLogSection({
        workoutLogId: log.id,
        type: sectionVal.type,
        label: sectionVal.label || undefined,
        position: si,
        restBetweenExercises: sectionVal.restBetweenExercises ?? undefined,
      });

      for (let ei = 0; ei < sectionVal.exercises.length; ei++) {
        await this.userApi.createWorkoutLogExercise({
          workoutLogSectionId: section.id,
          userExerciseSchemeId: sectionVal.exercises[ei].userExerciseSchemeId,
          position: ei,
        });
      }
    }

    this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.all() });
    this.router.navigate(['/user/workout-logs', log.id]);
  }
}
