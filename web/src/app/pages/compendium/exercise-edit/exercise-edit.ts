import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { form, required, FormField } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { SlugifyPipe } from '$ui/pipes/slugify';
import { PageLayout } from '../../../layout/page-layout';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmInput } from '@spartan-ng/helm/input';
import { HlmTextarea } from '@spartan-ng/helm/textarea';
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

const DIFFICULTIES: TechnicalDifficulty[] = [
  DifficultyBeginner,
  DifficultyIntermediate,
  DifficultyAdvanced,
];

const FORCES: Force[] = [
  ForcePull,
  ForcePush,
  ForceStatic,
  ForceDynamic,
  ForceHinge,
  ForceRotation,
];

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
  imports: [
    PageLayout,
    FormField,
    RouterLink,
    BrnSelectImports,
    HlmSelectImports,
    HlmInput,
    HlmTextarea,
    TranslocoDirective,
  ],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="
          isCreateMode() ? t('compendium.exercises.newTitle') : t('compendium.exercises.editTitle')
        "
        [isPending]="!isCreateMode() && exerciseQuery.isPending()"
        [errorMessage]="
          !isCreateMode() && exerciseQuery.isError() ? exerciseQuery.error().message : undefined
        "
      >
        @if (isCreateMode() || exerciseQuery.data()) {
          <form (submit)="onSubmit(); $event.preventDefault()" class="space-y-4">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.name') }} *</label
              >
              <input hlmInput id="name" [formField]="exerciseForm.name" class="mt-1" />
            </div>

            <div>
              <label
                for="description"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.description') }}</label
              >
              <textarea
                hlmTextarea
                id="description"
                [formField]="exerciseForm.description"
                rows="3"
                class="mt-1"
              ></textarea>
            </div>

            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label for="type" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                  >{{ t('fields.type') }} *</label
                >
                <brn-select [formField]="exerciseForm.type" class="mt-1" hlm>
                  <hlm-select-trigger class="w-full">
                    <hlm-select-value />
                  </hlm-select-trigger>
                  <hlm-select-content>
                    @for (tp of exerciseTypes; track tp) {
                      <hlm-option [value]="tp">{{ t('enums.exerciseType.' + tp) }}</hlm-option>
                    }
                  </hlm-select-content>
                </brn-select>
              </div>

              <div>
                <label
                  for="technicalDifficulty"
                  class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                  >{{ t('fields.difficulty') }} *</label
                >
                <brn-select [formField]="exerciseForm.technicalDifficulty" class="mt-1" hlm>
                  <hlm-select-trigger class="w-full">
                    <hlm-select-value />
                  </hlm-select-trigger>
                  <hlm-select-content>
                    @for (d of difficulties; track d) {
                      <hlm-option [value]="d">{{ t('enums.difficulty.' + d) }}</hlm-option>
                    }
                  </hlm-select-content>
                </brn-select>
              </div>

              <div>
                <label
                  for="bodyWeightScaling"
                  class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                  >{{ t('fields.bodyWeightScaling') }}</label
                >
                <input
                  hlmInput
                  id="bodyWeightScaling"
                  type="number"
                  step="0.01"
                  [formField]="exerciseForm.bodyWeightScaling"
                  class="mt-1"
                />
              </div>

              <div>
                <label
                  for="authorName"
                  class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                  >{{ t('fields.authorName') }}</label
                >
                <input
                  hlmInput
                  id="authorName"
                  [formField]="exerciseForm.authorName"
                  class="mt-1"
                />
              </div>

              <div class="sm:col-span-2">
                <label
                  for="authorUrl"
                  class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                  >{{ t('fields.authorUrl') }}</label
                >
                <input hlmInput id="authorUrl" [formField]="exerciseForm.authorUrl" class="mt-1" />
              </div>
            </div>

            <!-- Force checkboxes -->
            <fieldset>
              <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('fields.force') }}
              </legend>
              <div class="mt-1 flex flex-wrap gap-3">
                @for (f of forces; track f) {
                  <label class="flex items-center gap-1 text-sm text-gray-700 dark:text-gray-300">
                    <input
                      type="checkbox"
                      [checked]="isForceSelected(f)"
                      (change)="toggleArrayValue(selectedForce, f)"
                      class="rounded"
                    />
                    {{ t('enums.force.' + f) }}
                  </label>
                }
              </div>
            </fieldset>

            <!-- Primary Muscles checkboxes -->
            <fieldset>
              <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('fields.primaryMuscles') }}
              </legend>
              <div class="mt-1 flex flex-wrap gap-3">
                @for (m of muscles; track m) {
                  <label class="flex items-center gap-1 text-sm text-gray-700 dark:text-gray-300">
                    <input
                      type="checkbox"
                      [checked]="isPrimaryMuscleSelected(m)"
                      (change)="toggleArrayValue(selectedPrimaryMuscles, m)"
                      class="rounded"
                    />
                    {{ t('enums.muscle.' + m) }}
                  </label>
                }
              </div>
            </fieldset>

            <!-- Secondary Muscles checkboxes -->
            <fieldset>
              <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('fields.secondaryMuscles') }}
              </legend>
              <div class="mt-1 flex flex-wrap gap-3">
                @for (m of muscles; track m) {
                  <label class="flex items-center gap-1 text-sm text-gray-700 dark:text-gray-300">
                    <input
                      type="checkbox"
                      [checked]="isSecondaryMuscleSelected(m)"
                      (change)="toggleArrayValue(selectedSecondaryMuscles, m)"
                      class="rounded"
                    />
                    {{ t('enums.muscle.' + m) }}
                  </label>
                }
              </div>
            </fieldset>

            <!-- Measurement Paradigms checkboxes -->
            <fieldset>
              <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('fields.suggestedMeasurementParadigms') }}
              </legend>
              <div class="mt-1 flex flex-wrap gap-3">
                @for (p of measurementParadigms; track p) {
                  <label class="flex items-center gap-1 text-sm text-gray-700 dark:text-gray-300">
                    <input
                      type="checkbox"
                      [checked]="isParadigmSelected(p)"
                      (change)="toggleArrayValue(selectedParadigms, p)"
                      class="rounded"
                    />
                    {{ t('enums.measurementType.' + p) }}
                  </label>
                }
              </div>
            </fieldset>

            <!-- Instructions -->
            <fieldset>
              <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('fields.instructions') }}
              </legend>
              <div class="mt-1 space-y-2">
                @for (item of exerciseForm.instructions; track $index; let i = $index) {
                  <div class="flex gap-2">
                    <input hlmInput [formField]="item" />
                    <button
                      type="button"
                      (click)="removeInstruction(i)"
                      class="rounded-md border border-red-300 px-2 py-1 text-sm text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20"
                    >
                      {{ t('common.remove') }}
                    </button>
                  </div>
                }
                <button
                  type="button"
                  (click)="addInstruction()"
                  class="text-sm text-blue-600 hover:underline dark:text-blue-400"
                >
                  {{ t('compendium.exercises.addInstruction') }}
                </button>
              </div>
            </fieldset>

            <!-- Images -->
            <fieldset>
              <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('fields.images') }}
              </legend>
              <div class="mt-1 space-y-2">
                @for (item of exerciseForm.images; track $index; let i = $index) {
                  <div class="flex gap-2">
                    <input hlmInput [formField]="item" />
                    <button
                      type="button"
                      (click)="removeImage(i)"
                      class="rounded-md border border-red-300 px-2 py-1 text-sm text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20"
                    >
                      {{ t('common.remove') }}
                    </button>
                  </div>
                }
                <button
                  type="button"
                  (click)="addImage()"
                  class="text-sm text-blue-600 hover:underline dark:text-blue-400"
                >
                  {{ t('compendium.exercises.addImage') }}
                </button>
              </div>
            </fieldset>

            <!-- Alternative Names -->
            <fieldset>
              <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('fields.alternativeNames') }}
              </legend>
              <div class="mt-1 space-y-2">
                @for (item of exerciseForm.alternativeNames; track $index; let i = $index) {
                  <div class="flex gap-2">
                    <input hlmInput [formField]="item" />
                    <button
                      type="button"
                      (click)="removeAlternativeName(i)"
                      class="rounded-md border border-red-300 px-2 py-1 text-sm text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20"
                    >
                      {{ t('common.remove') }}
                    </button>
                  </div>
                }
                <button
                  type="button"
                  (click)="addAlternativeName()"
                  class="text-sm text-blue-600 hover:underline dark:text-blue-400"
                >
                  {{ t('compendium.exercises.addAlternativeName') }}
                </button>
              </div>
            </fieldset>

            <!-- Equipment IDs -->
            <fieldset>
              <legend class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('fields.equipmentIds') }}
              </legend>
              <div class="mt-1 space-y-2">
                @for (item of exerciseForm.equipmentIds; track $index; let i = $index) {
                  <div class="flex gap-2">
                    <input hlmInput type="number" [formField]="item" />
                    <button
                      type="button"
                      (click)="removeEquipmentId(i)"
                      class="rounded-md border border-red-300 px-2 py-1 text-sm text-red-600 hover:bg-red-50 dark:border-red-700 dark:text-red-400 dark:hover:bg-red-900/20"
                    >
                      {{ t('common.remove') }}
                    </button>
                  </div>
                }
                <button
                  type="button"
                  (click)="addEquipmentId()"
                  class="text-sm text-blue-600 hover:underline dark:text-blue-400"
                >
                  {{ t('compendium.exercises.addEquipmentId') }}
                </button>
              </div>
            </fieldset>

            <div class="flex gap-2">
              <button
                type="submit"
                [disabled]="
                  !exerciseForm().valid() || mutation.isPending() || createMutation.isPending()
                "
                class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {{ t('common.save') }}
              </button>
              <a
                [routerLink]="isCreateMode() ? ['/compendium/exercises'] : ['..']"
                class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
              >
                {{ t('common.cancel') }}
              </a>
            </div>
          </form>
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class ExerciseEdit {
  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = inject(QueryClient);
  private slugify = new SlugifyPipe();
  private params = toSignal(this.route.paramMap);

  private id = computed(() => Number(this.params()?.get('id')));
  isCreateMode = computed(() => !this.params()?.get('id'));

  exerciseTypes = EXERCISE_TYPES;
  difficulties = DIFFICULTIES;
  forces = FORCES;
  muscles = MUSCLES;
  measurementParadigms = MEASUREMENT_PARADIGMS;

  selectedForce = new Set<string>();
  selectedPrimaryMuscles = new Set<string>();
  selectedSecondaryMuscles = new Set<string>();
  selectedParadigms = new Set<string>();

  model = signal({
    name: '',
    description: '',
    type: ExerciseTypeStrength as string,
    technicalDifficulty: DifficultyBeginner as string,
    bodyWeightScaling: 0,
    authorName: '',
    authorUrl: '',
    instructions: [] as string[],
    images: [] as string[],
    alternativeNames: [] as string[],
    equipmentIds: [] as number[],
  });

  exerciseForm = form(this.model, (f) => {
    required(f.name);
    required(f.type);
    required(f.technicalDifficulty);
  });

  exerciseQuery = injectQuery(() => ({
    queryKey: exerciseKeys.detail(this.id()),
    queryFn: () => this.api.fetchExercise(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  permissionsQuery = injectQuery(() => ({
    queryKey: exerciseKeys.permissions(this.id()),
    queryFn: () => this.api.fetchExercisePermissions(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  mutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.updateExercise>[1]) =>
      this.api.updateExercise(this.id(), data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: exerciseKeys.all() });
      this.router.navigate(['..'], { relativeTo: this.route });
    },
  }));

  createMutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.createExercise>[0]) =>
      this.api.createExercise(data),
    onSuccess: (result) => {
      this.queryClient.invalidateQueries({ queryKey: exerciseKeys.all() });
      this.router.navigate([
        '/compendium/exercises',
        result.id,
        this.slugify.transform(result.name),
      ]);
    },
  }));

  constructor() {
    effect(() => {
      const perms = this.permissionsQuery.data();
      if (perms && !perms.permissions.includes('MODIFY')) {
        this.router.navigate(['..'], { relativeTo: this.route });
      }
    });

    effect(() => {
      const data = this.exerciseQuery.data();
      if (data) {
        this.model.set({
          name: data.name,
          description: data.description,
          type: data.type,
          technicalDifficulty: data.technicalDifficulty,
          bodyWeightScaling: data.bodyWeightScaling,
          authorName: data.authorName ?? '',
          authorUrl: data.authorUrl ?? '',
          instructions: data.instructions ?? [],
          images: data.images ?? [],
          alternativeNames: data.alternativeNames ?? [],
          equipmentIds: data.equipmentIds ?? [],
        });

        this.selectedForce = new Set(data.force ?? []);
        this.selectedPrimaryMuscles = new Set(data.primaryMuscles ?? []);
        this.selectedSecondaryMuscles = new Set(data.secondaryMuscles ?? []);
        this.selectedParadigms = new Set(data.suggestedMeasurementParadigms ?? []);
      }
    });
  }

  isForceSelected(f: string) {
    return this.selectedForce.has(f);
  }
  isPrimaryMuscleSelected(m: string) {
    return this.selectedPrimaryMuscles.has(m);
  }
  isSecondaryMuscleSelected(m: string) {
    return this.selectedSecondaryMuscles.has(m);
  }
  isParadigmSelected(p: string) {
    return this.selectedParadigms.has(p);
  }

  toggleArrayValue(set: Set<string>, value: string) {
    if (set.has(value)) {
      set.delete(value);
    } else {
      set.add(value);
    }
  }

  addInstruction() {
    this.model.update((m) => ({ ...m, instructions: [...m.instructions, ''] }));
  }
  removeInstruction(i: number) {
    this.model.update((m) => ({
      ...m,
      instructions: m.instructions.filter((_, idx) => idx !== i),
    }));
  }

  addImage() {
    this.model.update((m) => ({ ...m, images: [...m.images, ''] }));
  }
  removeImage(i: number) {
    this.model.update((m) => ({ ...m, images: m.images.filter((_, idx) => idx !== i) }));
  }

  addAlternativeName() {
    this.model.update((m) => ({ ...m, alternativeNames: [...m.alternativeNames, ''] }));
  }
  removeAlternativeName(i: number) {
    this.model.update((m) => ({
      ...m,
      alternativeNames: m.alternativeNames.filter((_, idx) => idx !== i),
    }));
  }

  addEquipmentId() {
    this.model.update((m) => ({ ...m, equipmentIds: [...m.equipmentIds, 0] }));
  }
  removeEquipmentId(i: number) {
    this.model.update((m) => ({
      ...m,
      equipmentIds: m.equipmentIds.filter((_, idx) => idx !== i),
    }));
  }

  onSubmit() {
    if (this.exerciseForm().valid()) {
      const val = this.model();
      const data = this.exerciseQuery.data();
      const payload = {
        ...(this.isCreateMode()
          ? {}
          : {
              public: data!.public,
              parentExerciseId: data!.parentExerciseId,
            }),
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
      };
      if (this.isCreateMode()) {
        this.createMutation.mutate(payload);
      } else {
        this.mutation.mutate(payload);
      }
    }
  }
}
