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
import { formatBreak } from '$core/format-utils';
import { workoutKeys, workoutLogKeys } from '$core/query-keys';
import { WorkoutSectionTypeMain, WorkoutSectionTypeSupplementary } from '$generated/user-models';
import { PageLayout } from '../../../layout/page-layout';
import { WorkoutStartStore, SetPreview } from './workout-start.store';

type SetFormGroup = FormGroup<{
  targetReps: FormControl<number | null>;
  targetWeight: FormControl<number | null>;
  targetDuration: FormControl<number | null>;
  targetDistance: FormControl<number | null>;
  targetTime: FormControl<number | null>;
  restAfterSeconds: FormControl<number | null>;
}>;

type ExerciseFormGroup = FormGroup<{
  sourceExerciseSchemeId: FormControl<number>;
  breakAfterSeconds: FormControl<number | null>;
  sets: FormArray<SetFormGroup>;
}>;

type SectionFormGroup = FormGroup<{
  type: FormControl<string>;
  label: FormControl<string>;
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

                <div class="mb-3 grid grid-cols-1 gap-3 sm:grid-cols-2">
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
                </div>

                <!-- Exercise cards -->
                <div formArrayName="exercises">
                  @for (
                    ex of getExercisesArray(si).controls;
                    track $index;
                    let ei = $index;
                    let lastEx = $last
                  ) {
                    @let info = store.exerciseDisplay()[ex.get('sourceExerciseSchemeId')!.value];
                    <div
                      [formGroupName]="ei"
                      class="rounded-md border border-gray-200 dark:border-gray-600"
                    >
                      <!-- Exercise header -->
                      <div
                        class="flex items-center justify-between border-b border-gray-100 px-3 py-2 dark:border-gray-700"
                      >
                        <div class="text-sm text-gray-900 dark:text-gray-100">
                          <span class="font-semibold">{{ info?.name ?? 'Loading...' }}</span>
                          <span class="ml-2 text-gray-500 dark:text-gray-400">{{
                            info?.summary
                          }}</span>
                        </div>
                        <button
                          type="button"
                          (click)="removeExercise(si, ei)"
                          class="text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                        >
                          Remove
                        </button>
                      </div>

                      <!-- Editable sets -->
                      @if (getSetsArray(si, ei).length) {
                        @let setsArr = getSetsArray(si, ei);
                        <div class="px-3 py-2" [formArrayName]="'sets'">
                          <!-- Header -->
                          <div
                            class="mb-1 grid text-left text-xs text-gray-500 uppercase dark:text-gray-400"
                            [class]="
                              info?.measurementType === 'REP_BASED'
                                ? 'grid-cols-[2rem_5rem_6rem]'
                                : 'grid-cols-[2rem_6rem]'
                            "
                          >
                            <span>Set</span>
                            @if (info?.measurementType === 'REP_BASED') {
                              <span>Reps</span>
                              <span>Weight</span>
                            }
                            @if (info?.measurementType === 'TIME_BASED') {
                              <span>Duration</span>
                            }
                            @if (info?.measurementType === 'DISTANCE_BASED') {
                              <span>Distance</span>
                            }
                          </div>

                          @for (
                            setCtrl of setsArr.controls;
                            track $index;
                            let setIdx = $index;
                            let lastSet = $last
                          ) {
                            <!-- Set row -->
                            <div
                              [formGroupName]="setIdx"
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
                                    type="number"
                                    formControlName="targetReps"
                                    class="w-16 rounded border border-gray-300 px-2 py-1 text-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                                  />
                                </div>
                                <div>
                                  <input
                                    type="number"
                                    formControlName="targetWeight"
                                    class="w-20 rounded border border-gray-300 px-2 py-1 text-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                                    step="0.5"
                                  />
                                </div>
                              }
                              @if (info?.measurementType === 'TIME_BASED') {
                                <div>
                                  <input
                                    type="number"
                                    formControlName="targetDuration"
                                    class="w-20 rounded border border-gray-300 px-2 py-1 text-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                                  />
                                </div>
                              }
                              @if (info?.measurementType === 'DISTANCE_BASED') {
                                <div>
                                  <input
                                    type="number"
                                    formControlName="targetDistance"
                                    class="w-20 rounded border border-gray-300 px-2 py-1 text-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                                    step="0.1"
                                  />
                                </div>
                              }
                            </div>

                            <!-- Rest between sets: line with centered badge -->
                            @if (!lastSet && setCtrl.controls.restAfterSeconds.value !== null) {
                              <div
                                [formGroupName]="setIdx"
                                class="relative flex items-center justify-center py-0.5"
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
                                  <input
                                    type="number"
                                    formControlName="restAfterSeconds"
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

                    <!-- Break after exercise (editable, not shown after last) -->
                    @if (!lastEx) {
                      <div
                        [formGroupName]="ei"
                        class="relative flex items-center justify-center py-3"
                      >
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
                            formControlName="breakAfterSeconds"
                            class="w-12 border-0 bg-transparent p-0 text-center text-xs text-gray-500 focus:ring-0 dark:text-gray-400"
                          />
                          <span>s rest</span>
                        </div>
                      </div>
                    }
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

  getSetsArray(sectionIndex: number, exerciseIndex: number) {
    return this.getExercisesArray(sectionIndex).at(exerciseIndex).controls.sets;
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
        });

        for (const ex of section.exercises ?? []) {
          sectionGroup.controls.exercises.push(
            this.createExerciseGroup(ex.userExerciseSchemeId, section.restBetweenExercises ?? null),
          );
        }

        this.sectionsArray.push(sectionGroup);
      }

      this.store.loadExerciseDisplay(data.sections ?? []);
    });

    // Populate set form arrays once store display data is loaded
    effect(() => {
      const display = this.store.exerciseDisplay();
      if (this.store.isLoadingDisplay() || Object.keys(display).length === 0) return;

      for (let si = 0; si < this.sectionsArray.length; si++) {
        const exercises = this.getExercisesArray(si);
        for (let ei = 0; ei < exercises.length; ei++) {
          const exGroup = exercises.at(ei);
          const schemeId = exGroup.controls.sourceExerciseSchemeId.value;
          const info = display[schemeId];
          if (info && exGroup.controls.sets.length === 0 && info.sets.length > 0) {
            for (const set of info.sets) {
              exGroup.controls.sets.push(this.createSetGroup(set));
            }
          }
        }
      }
    });
  }

  removeSection(index: number) {
    this.sectionsArray.removeAt(index);
  }

  removeExercise(sectionIndex: number, exerciseIndex: number) {
    this.getExercisesArray(sectionIndex).removeAt(exerciseIndex);
  }

  formatBreak = formatBreak;

  onSubmit() {
    this.startMutation.mutate();
  }

  private createSectionGroup(): SectionFormGroup {
    return new FormGroup({
      type: new FormControl(WorkoutSectionTypeMain, { nonNullable: true }),
      label: new FormControl('', { nonNullable: true }),
      exercises: new FormArray<ExerciseFormGroup>([]),
    });
  }

  private createExerciseGroup(
    schemeId: number,
    breakAfterSeconds: number | null = null,
  ): ExerciseFormGroup {
    return new FormGroup({
      sourceExerciseSchemeId: new FormControl(schemeId, { nonNullable: true }),
      breakAfterSeconds: new FormControl<number | null>(breakAfterSeconds),
      sets: new FormArray<SetFormGroup>([]),
    });
  }

  private createSetGroup(set: SetPreview): SetFormGroup {
    return new FormGroup({
      targetReps: new FormControl<number | null>(set.targetReps ?? null),
      targetWeight: new FormControl<number | null>(set.targetWeight ?? null),
      targetDuration: new FormControl<number | null>(set.targetDuration ?? null),
      targetDistance: new FormControl<number | null>(set.targetDistance ?? null),
      targetTime: new FormControl<number | null>(set.targetTime ?? null),
      restAfterSeconds: new FormControl<number | null>(set.restAfterSeconds ?? null),
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
      });

      for (let ei = 0; ei < sectionVal.exercises.length; ei++) {
        const exVal = sectionVal.exercises[ei];
        const logExercise = await this.userApi.createWorkoutLogExercise({
          workoutLogSectionId: section.id,
          sourceExerciseSchemeId: exVal.sourceExerciseSchemeId,
          position: ei,
          breakAfterSeconds: exVal.breakAfterSeconds ?? undefined,
        });

        // Apply per-set target overrides where user changed values
        for (let setIdx = 0; setIdx < exVal.sets.length; setIdx++) {
          const setVal = exVal.sets[setIdx];
          const created = logExercise.sets?.[setIdx];
          if (!created) continue;

          const overrides: Record<string, number | undefined> = {};
          if (setVal.targetReps != null && setVal.targetReps !== created.targetReps)
            overrides['targetReps'] = setVal.targetReps;
          if (setVal.targetWeight != null && setVal.targetWeight !== created.targetWeight)
            overrides['targetWeight'] = setVal.targetWeight;
          if (setVal.targetDuration != null && setVal.targetDuration !== created.targetDuration)
            overrides['targetDuration'] = setVal.targetDuration;
          if (setVal.targetDistance != null && setVal.targetDistance !== created.targetDistance)
            overrides['targetDistance'] = setVal.targetDistance;
          if (setVal.targetTime != null && setVal.targetTime !== created.targetTime)
            overrides['targetTime'] = setVal.targetTime;
          if (
            setVal.restAfterSeconds != null &&
            setVal.restAfterSeconds !== created.breakAfterSeconds
          )
            overrides['breakAfterSeconds'] = setVal.restAfterSeconds;

          if (Object.keys(overrides).length > 0) {
            await this.userApi.updateWorkoutLogExerciseSet(created.id, overrides);
          }
        }
      }
    }

    this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.all() });
    this.router.navigate(['/user/workout-logs', log.id]);
  }
}
