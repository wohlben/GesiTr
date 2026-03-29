import { Component, inject, computed, signal, input, effect } from '@angular/core';
import { form, FormField, disabled } from '@angular/forms/signals';
import { injectQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseKeys, exerciseSchemeKeys, masteryKeys } from '$core/query-keys';
import { ExerciseScheme } from '$generated/models';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmInput } from '@spartan-ng/helm/input';
import { TranslocoDirective } from '@jsverse/transloco';

export interface ExerciseConfigResult {
  exerciseId: number;
  exerciseName: string;
  scheme: ExerciseScheme;
}

@Component({
  selector: 'app-exercise-config',
  imports: [FormField, BrnSelectImports, HlmSelectImports, HlmInput, TranslocoDirective],
  template: `
    <div
      *transloco="let t"
      class="rounded-md border border-gray-100 bg-gray-50 p-3 dark:border-gray-600 dark:bg-gray-800/50"
    >
      <div class="mb-2 grid grid-cols-1 gap-2 sm:grid-cols-2">
        <div>
          <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
            t('ui.exerciseConfig.exerciseLabel')
          }}</span>
          <brn-select
            [formField]="configForm.exerciseId"
            class="mt-1"
            hlm
            [placeholder]="t('common.select')"
          >
            <hlm-select-trigger class="w-full">
              <hlm-select-value />
            </hlm-select-trigger>
            <hlm-select-content>
              @for (ex of sortedExercises(); track ex.id) {
                <hlm-option [value]="ex.id">{{ ex.name }}</hlm-option>
              }
            </hlm-select-content>
          </brn-select>
        </div>
        <div>
          <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
            t('fields.measurementType')
          }}</span>
          <brn-select [formField]="configForm.measurementType" class="mt-1" hlm>
            <hlm-select-trigger class="w-full">
              <hlm-select-value />
            </hlm-select-trigger>
            <hlm-select-content>
              <hlm-option value="REP_BASED">{{ t('enums.measurementType.REP_BASED') }}</hlm-option>
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
      @if (model().measurementType === 'REP_BASED') {
        <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            {{ t('fields.sets') }}
            <input type="number" [formField]="configForm.sets" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            {{ t('fields.reps') }}
            <input type="number" [formField]="configForm.reps" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            {{ t('fields.weightKg') }}
            <input type="number" [formField]="configForm.weight" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            {{ t('fields.restSeconds') }}
            <input type="number" [formField]="configForm.restBetweenSets" hlmInput class="mt-1" />
          </label>
        </div>
      }

      <!-- TIME_BASED fields -->
      @if (model().measurementType === 'TIME_BASED') {
        <div class="grid grid-cols-2 gap-2">
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            {{ t('fields.durationSeconds') }}
            <input type="number" [formField]="configForm.duration" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            {{ t('fields.timePerRepSeconds') }}
            <input type="number" [formField]="configForm.timePerRep" hlmInput class="mt-1" />
          </label>
        </div>
      }

      <!-- DISTANCE_BASED fields -->
      @if (model().measurementType === 'DISTANCE_BASED') {
        <div class="grid grid-cols-2 gap-2">
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            {{ t('fields.distanceM') }}
            <input type="number" [formField]="configForm.distance" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            {{ t('fields.targetTimeSeconds') }}
            <input type="number" [formField]="configForm.targetTime" hlmInput class="mt-1" />
          </label>
        </div>
      }
    </div>
  `,
})
export class ExerciseConfig {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private queryClient = inject(QueryClient);

  /** When set, the exercise dropdown is locked to this user exercise. */
  preselectedExerciseId = input<number | null>(null);

  // Form state
  model = signal({
    exerciseId: null as number | null,
    measurementType: 'REP_BASED',
    sets: 3 as number | null,
    reps: 10 as number | null,
    weight: null as number | null,
    restBetweenSets: null as number | null,
    timePerRep: null as number | null,
    duration: null as number | null,
    distance: null as number | null,
    targetTime: null as number | null,
  });
  configForm = form(this.model, (f) => {
    disabled(f.exerciseId, () => !!this.preselectedExerciseId());
  });

  constructor() {
    effect(() => {
      const preselected = this.preselectedExerciseId();
      if (preselected != null) {
        this.model.update((m) => ({ ...m, exerciseId: preselected }));
      }
    });
  }

  // All exercises
  private allExercisesQuery = injectQuery(() => ({
    queryKey: exerciseKeys.list({ limit: 200 }),
    queryFn: () => this.compendiumApi.fetchExercises({ limit: 200 }),
  }));

  // Mastery list (used for sorting mastered exercises first)
  private masteryQuery = injectQuery(() => ({
    queryKey: masteryKeys.list(),
    queryFn: () => this.userApi.fetchMasteryList(),
  }));

  sortedExercises = computed(() => {
    const all = this.allExercisesQuery.data()?.items ?? [];
    const masteryIds = new Set((this.masteryQuery.data() ?? []).map((m) => m.exerciseId));
    return [...all].sort((a, b) => {
      const aHas = masteryIds.has(a.id) ? 0 : 1;
      const bHas = masteryIds.has(b.id) ? 0 : 1;
      if (aHas !== bHas) return aHas - bHas;
      return a.name.localeCompare(b.name);
    });
  });

  selectedExerciseName = computed(() => {
    const id = this.model().exerciseId;
    if (!id) return '';
    return this.sortedExercises().find((e) => e.id === id)?.name ?? '';
  });

  canConfirm = computed(() => this.model().exerciseId != null);

  /** Creates a scheme from the field values and returns it. */
  async confirm(): Promise<ExerciseScheme> {
    const m = this.model();
    const data: Record<string, unknown> = {
      exerciseId: m.exerciseId,
      measurementType: m.measurementType,
    };
    const mt = m.measurementType;
    if (mt === 'REP_BASED') {
      if (m.sets != null) data['sets'] = m.sets;
      if (m.reps != null) data['reps'] = m.reps;
      if (m.weight != null) data['weight'] = m.weight;
      if (m.restBetweenSets != null) data['restBetweenSets'] = m.restBetweenSets;
    } else if (mt === 'TIME_BASED') {
      if (m.sets != null) data['sets'] = m.sets;
      if (m.duration != null) data['duration'] = m.duration;
      if (m.timePerRep != null) data['timePerRep'] = m.timePerRep;
    } else if (mt === 'DISTANCE_BASED') {
      if (m.sets != null) data['sets'] = m.sets;
      if (m.distance != null) data['distance'] = m.distance;
      if (m.targetTime != null) data['targetTime'] = m.targetTime;
    }
    const scheme = await this.userApi.createExerciseScheme(data);
    this.queryClient.invalidateQueries({ queryKey: exerciseSchemeKeys.all() });
    return scheme;
  }

  reset() {
    this.model.set({
      exerciseId: null,
      measurementType: 'REP_BASED',
      sets: 3,
      reps: 10,
      weight: null,
      restBetweenSets: null,
      timePerRep: null,
      duration: null,
      distance: null,
      targetTime: null,
    });
  }
}
