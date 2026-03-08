import { Component, inject, computed, effect } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { ReactiveFormsModule, FormGroup, FormControl, FormArray, Validators } from '@angular/forms';
import { injectQuery, injectMutation, injectQueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import {
  ExerciseType,
  ExerciseTypeStrength,
  ExerciseTypeCardio,
  ExerciseTypeStretching,
  ExerciseTypeStrongman,
  TechnicalDifficulty,
  DifficultyBeginner,
  DifficultyIntermediate,
  DifficultyAdvanced,
  Force,
  ForcePull,
  ForcePush,
  ForceStatic,
  ForceDynamic,
  ForceHinge,
  ForceRotation,
  Muscle,
  MuscleAbs,
  MuscleAdductors,
  MuscleBiceps,
  MuscleCalves,
  MuscleChest,
  MuscleForearms,
  MuscleGlutes,
  MuscleHamstrings,
  MuscleHipFlexors,
  MuscleLats,
  MuscleLowerBack,
  MuscleNeck,
  MuscleObliques,
  MuscleQuads,
  MuscleTraps,
  MuscleTriceps,
  MuscleFrontDelts,
  MuscleRearDelts,
  MuscleRhomboids,
  MuscleSideDelts,
  MeasurementParadigm,
  MeasurementRepBased,
  MeasurementAMRAP,
  MeasurementTimeBased,
  MeasurementDistanceBased,
  MeasurementEMOM,
  MeasurementRoundsForTime,
  MeasurementTime,
  MeasurementDistance,
} from '$generated/models';

const EXERCISE_TYPES: ExerciseType[] = [
  ExerciseTypeStrength,
  ExerciseTypeCardio,
  ExerciseTypeStretching,
  ExerciseTypeStrongman,
];

const DIFFICULTIES: TechnicalDifficulty[] = [DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced];

const FORCES: Force[] = [ForcePull, ForcePush, ForceStatic, ForceDynamic, ForceHinge, ForceRotation];

const MUSCLES: Muscle[] = [
  MuscleAbs,
  MuscleAdductors,
  MuscleBiceps,
  MuscleCalves,
  MuscleChest,
  MuscleForearms,
  MuscleGlutes,
  MuscleHamstrings,
  MuscleHipFlexors,
  MuscleLats,
  MuscleLowerBack,
  MuscleNeck,
  MuscleObliques,
  MuscleQuads,
  MuscleTraps,
  MuscleTriceps,
  MuscleFrontDelts,
  MuscleRearDelts,
  MuscleRhomboids,
  MuscleSideDelts,
];

const MEASUREMENT_PARADIGMS: MeasurementParadigm[] = [
  MeasurementRepBased,
  MeasurementAMRAP,
  MeasurementTimeBased,
  MeasurementDistanceBased,
  MeasurementEMOM,
  MeasurementRoundsForTime,
  MeasurementTime,
  MeasurementDistance,
];

@Component({
  selector: 'app-exercise-edit',
  imports: [PageLayout, ReactiveFormsModule, RouterLink],
  template: `
    <app-page-layout
      header="Edit Exercise"
      [isPending]="exerciseQuery.isPending()"
      [errorMessage]="exerciseQuery.isError() ? exerciseQuery.error().message : undefined"
    >
      @if (exerciseQuery.data(); as exercise) {
        <form [formGroup]="form" (ngSubmit)="onSubmit()" class="space-y-4">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Name *</label>
            <input id="name" formControlName="name" class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100" />
          </div>

          <div>
            <label for="description" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Description</label>
            <textarea id="description" formControlName="description" rows="3" class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"></textarea>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label for="type" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Type *</label>
              <select id="type" formControlName="type" class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100">
                @for (t of exerciseTypes; track t) {
                  <option [value]="t">{{ t }}</option>
                }
              </select>
            </div>

            <div>
              <label for="technicalDifficulty" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Difficulty *</label>
              <select id="technicalDifficulty" formControlName="technicalDifficulty" class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100">
                @for (d of difficulties; track d) {
                  <option [value]="d">{{ d }}</option>
                }
              </select>
            </div>

            <div>
              <label for="bodyWeightScaling" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Body Weight Scaling</label>
              <input id="bodyWeightScaling" type="number" step="0.01" formControlName="bodyWeightScaling" class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100" />
            </div>

            <div>
              <label for="authorName" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Author Name</label>
              <input id="authorName" formControlName="authorName" class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100" />
            </div>

            <div class="sm:col-span-2">
              <label for="authorUrl" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Author URL</label>
              <input id="authorUrl" formControlName="authorUrl" class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100" />
            </div>
          </div>

          <!-- Force checkboxes -->
          <fieldset>
            <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">Force</legend>
            <div class="mt-1 flex flex-wrap gap-3">
              @for (f of forces; track f) {
                <label class="flex items-center gap-1 text-sm text-gray-700 dark:text-gray-300">
                  <input type="checkbox" [checked]="isForceSelected(f)" (change)="toggleArrayValue(selectedForce, f)" class="rounded" />
                  {{ f }}
                </label>
              }
            </div>
          </fieldset>

          <!-- Primary Muscles checkboxes -->
          <fieldset>
            <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">Primary Muscles</legend>
            <div class="mt-1 flex flex-wrap gap-3">
              @for (m of muscles; track m) {
                <label class="flex items-center gap-1 text-sm text-gray-700 dark:text-gray-300">
                  <input type="checkbox" [checked]="isPrimaryMuscleSelected(m)" (change)="toggleArrayValue(selectedPrimaryMuscles, m)" class="rounded" />
                  {{ m }}
                </label>
              }
            </div>
          </fieldset>

          <!-- Secondary Muscles checkboxes -->
          <fieldset>
            <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">Secondary Muscles</legend>
            <div class="mt-1 flex flex-wrap gap-3">
              @for (m of muscles; track m) {
                <label class="flex items-center gap-1 text-sm text-gray-700 dark:text-gray-300">
                  <input type="checkbox" [checked]="isSecondaryMuscleSelected(m)" (change)="toggleArrayValue(selectedSecondaryMuscles, m)" class="rounded" />
                  {{ m }}
                </label>
              }
            </div>
          </fieldset>

          <!-- Measurement Paradigms checkboxes -->
          <fieldset>
            <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">Suggested Measurement Paradigms</legend>
            <div class="mt-1 flex flex-wrap gap-3">
              @for (p of measurementParadigms; track p) {
                <label class="flex items-center gap-1 text-sm text-gray-700 dark:text-gray-300">
                  <input type="checkbox" [checked]="isParadigmSelected(p)" (change)="toggleArrayValue(selectedParadigms, p)" class="rounded" />
                  {{ p }}
                </label>
              }
            </div>
          </fieldset>

          <!-- Instructions FormArray -->
          <fieldset>
            <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">Instructions</legend>
            <div class="mt-1 space-y-2">
              @for (ctrl of instructions.controls; track $index) {
                <div class="flex gap-2">
                  <input [formControl]="ctrl" class="block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100" />
                  <button type="button" (click)="instructions.removeAt($index)" class="rounded-md border border-red-300 px-2 py-1 text-sm text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20">Remove</button>
                </div>
              }
              <button type="button" (click)="instructions.push(newStringControl())" class="text-sm text-blue-600 hover:underline dark:text-blue-400">+ Add instruction</button>
            </div>
          </fieldset>

          <!-- Images FormArray -->
          <fieldset>
            <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">Images</legend>
            <div class="mt-1 space-y-2">
              @for (ctrl of images.controls; track $index) {
                <div class="flex gap-2">
                  <input [formControl]="ctrl" class="block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100" />
                  <button type="button" (click)="images.removeAt($index)" class="rounded-md border border-red-300 px-2 py-1 text-sm text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20">Remove</button>
                </div>
              }
              <button type="button" (click)="images.push(newStringControl())" class="text-sm text-blue-600 hover:underline dark:text-blue-400">+ Add image</button>
            </div>
          </fieldset>

          <!-- Alternative Names FormArray -->
          <fieldset>
            <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">Alternative Names</legend>
            <div class="mt-1 space-y-2">
              @for (ctrl of alternativeNames.controls; track $index) {
                <div class="flex gap-2">
                  <input [formControl]="ctrl" class="block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100" />
                  <button type="button" (click)="alternativeNames.removeAt($index)" class="rounded-md border border-red-300 px-2 py-1 text-sm text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20">Remove</button>
                </div>
              }
              <button type="button" (click)="alternativeNames.push(newStringControl())" class="text-sm text-blue-600 hover:underline dark:text-blue-400">+ Add alternative name</button>
            </div>
          </fieldset>

          <!-- Equipment IDs FormArray -->
          <fieldset>
            <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">Equipment IDs</legend>
            <div class="mt-1 space-y-2">
              @for (ctrl of equipmentIds.controls; track $index) {
                <div class="flex gap-2">
                  <input [formControl]="ctrl" class="block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100" />
                  <button type="button" (click)="equipmentIds.removeAt($index)" class="rounded-md border border-red-300 px-2 py-1 text-sm text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20">Remove</button>
                </div>
              }
              <button type="button" (click)="equipmentIds.push(newStringControl())" class="text-sm text-blue-600 hover:underline dark:text-blue-400">+ Add equipment ID</button>
            </div>
          </fieldset>

          <div class="flex gap-2">
            <button
              type="submit"
              [disabled]="form.invalid || mutation.isPending()"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              Save
            </button>
            <a [routerLink]="['..']" class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800">
              Cancel
            </a>
          </div>
        </form>
      }
    </app-page-layout>
  `,
})
export class ExerciseEdit {
  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = injectQueryClient();
  private params = toSignal(this.route.paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  exerciseTypes = EXERCISE_TYPES;
  difficulties = DIFFICULTIES;
  forces = FORCES;
  muscles = MUSCLES;
  measurementParadigms = MEASUREMENT_PARADIGMS;

  selectedForce: Set<string> = new Set();
  selectedPrimaryMuscles: Set<string> = new Set();
  selectedSecondaryMuscles: Set<string> = new Set();
  selectedParadigms: Set<string> = new Set();

  instructions = new FormArray<FormControl<string>>([]);
  images = new FormArray<FormControl<string>>([]);
  alternativeNames = new FormArray<FormControl<string>>([]);
  equipmentIds = new FormArray<FormControl<string>>([]);

  form = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl('', { nonNullable: true }),
    type: new FormControl<ExerciseType>(ExerciseTypeStrength, { nonNullable: true, validators: [Validators.required] }),
    technicalDifficulty: new FormControl<TechnicalDifficulty>(DifficultyBeginner, { nonNullable: true, validators: [Validators.required] }),
    bodyWeightScaling: new FormControl(0, { nonNullable: true }),
    authorName: new FormControl('', { nonNullable: true }),
    authorUrl: new FormControl('', { nonNullable: true }),
    instructions: this.instructions,
    images: this.images,
    alternativeNames: this.alternativeNames,
    equipmentIds: this.equipmentIds,
  });

  exerciseQuery = injectQuery(() => ({
    queryKey: exerciseKeys.detail(this.id()),
    queryFn: () => this.api.fetchExercise(this.id()),
    enabled: !!this.id(),
  }));

  mutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.updateExercise>[1]) =>
      this.api.updateExercise(this.id(), data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: exerciseKeys.all() });
      this.router.navigate(['..'], { relativeTo: this.route });
    },
  }));

  constructor() {
    effect(() => {
      const data = this.exerciseQuery.data();
      if (data) {
        this.form.patchValue({
          name: data.name,
          description: data.description,
          type: data.type,
          technicalDifficulty: data.technicalDifficulty,
          bodyWeightScaling: data.bodyWeightScaling,
          authorName: data.authorName ?? '',
          authorUrl: data.authorUrl ?? '',
        });

        this.selectedForce = new Set(data.force ?? []);
        this.selectedPrimaryMuscles = new Set(data.primaryMuscles ?? []);
        this.selectedSecondaryMuscles = new Set(data.secondaryMuscles ?? []);
        this.selectedParadigms = new Set(data.suggestedMeasurementParadigms ?? []);

        this.rebuildFormArray(this.instructions, data.instructions ?? []);
        this.rebuildFormArray(this.images, data.images ?? []);
        this.rebuildFormArray(this.alternativeNames, data.alternativeNames ?? []);
        this.rebuildFormArray(this.equipmentIds, data.equipmentIds ?? []);
      }
    });
  }

  isForceSelected(f: string) { return this.selectedForce.has(f); }
  isPrimaryMuscleSelected(m: string) { return this.selectedPrimaryMuscles.has(m); }
  isSecondaryMuscleSelected(m: string) { return this.selectedSecondaryMuscles.has(m); }
  isParadigmSelected(p: string) { return this.selectedParadigms.has(p); }

  toggleArrayValue(set: Set<string>, value: string) {
    if (set.has(value)) {
      set.delete(value);
    } else {
      set.add(value);
    }
  }

  newStringControl() {
    return new FormControl('', { nonNullable: true });
  }

  private rebuildFormArray(arr: FormArray<FormControl<string>>, values: string[]) {
    arr.clear();
    for (const v of values) {
      arr.push(new FormControl(v, { nonNullable: true }));
    }
  }

  onSubmit() {
    if (this.form.valid) {
      const val = this.form.getRawValue();
      this.mutation.mutate({
        ...this.exerciseQuery.data()!,
        name: val.name,
        description: val.description,
        type: val.type,
        technicalDifficulty: val.technicalDifficulty,
        bodyWeightScaling: val.bodyWeightScaling,
        authorName: val.authorName || undefined,
        authorUrl: val.authorUrl || undefined,
        force: [...this.selectedForce],
        primaryMuscles: [...this.selectedPrimaryMuscles],
        secondaryMuscles: [...this.selectedSecondaryMuscles],
        suggestedMeasurementParadigms: [...this.selectedParadigms],
        instructions: val.instructions,
        images: val.images,
        alternativeNames: val.alternativeNames,
        equipmentIds: val.equipmentIds,
      });
    }
  }
}
