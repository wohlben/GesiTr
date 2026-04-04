import { Component, computed, input, model, signal } from '@angular/core';
import { TranslocoDirective } from '@jsverse/transloco';
import { HlmComboboxImports } from '@spartan-ng/helm/combobox';
import { Exercise } from '$generated/models';
import { ExerciseScheme } from '$generated/user-exercisescheme';
import { SchemeSelector } from '$ui/scheme-selector/scheme-selector';
import { CreateSchemeDialog } from './create-scheme-dialog';
import type { FormValueControl } from '@angular/forms/signals';
import type { WorkoutItemModel } from './workout-item-model';

@Component({
  selector: 'app-exercise-item-editor',
  imports: [HlmComboboxImports, SchemeSelector, CreateSchemeDialog, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <div>
        <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
          t('ui.exerciseConfig.exerciseLabel')
        }}</span>
        <hlm-combobox
          class="mt-1 block"
          [value]="selectedExercise()"
          (valueChange)="onExerciseSelected($event)"
          [filter]="exerciseFilter"
          [itemToString]="exerciseToString"
        >
          <hlm-combobox-input
            [placeholder]="t('common.search')"
            [showClear]="!!value().exerciseId"
          />
          <ng-template hlmComboboxPortal>
            <hlm-combobox-content>
              <hlm-combobox-input [placeholder]="t('common.search')" [showClear]="false" />
              <div hlmComboboxList>
                @for (ex of exercises(); track ex.id) {
                  <hlm-combobox-item [value]="ex">{{ ex.names?.[0]?.name }}</hlm-combobox-item>
                }
                <hlm-combobox-empty>{{ t('common.noResults') }}</hlm-combobox-empty>
              </div>
            </hlm-combobox-content>
          </ng-template>
        </hlm-combobox>
      </div>

      <app-scheme-selector
        [exerciseId]="value().exerciseId"
        [selectedSchemeId]="value().selectedSchemeId"
        (schemeSelected)="onSchemeSelected($event)"
        (createRequested)="openDialog(null)"
        (editRequested)="openDialog($event)"
      />

      <app-create-scheme-dialog
        [open]="dialogOpen()"
        [preselectedExerciseId]="value().exerciseId"
        [editingScheme]="editingScheme()"
        (schemeSaved)="onSchemeSaved($event)"
        (cancelled)="closeDialog()"
      />
    </ng-container>
  `,
})
export class ExerciseItemEditor implements FormValueControl<WorkoutItemModel> {
  readonly value = model.required<WorkoutItemModel>();
  exercises = input.required<Exercise[]>();

  dialogOpen = signal(false);
  editingScheme = signal<ExerciseScheme | null>(null);

  selectedExercise = computed(() => {
    const id = this.value().exerciseId;
    if (!id) return null;
    return this.exercises().find((e) => e.id === id) ?? null;
  });

  exerciseFilter = (exercise: Exercise, search: string) =>
    exercise.names?.some((n) => n.name.toLowerCase().includes(search.toLowerCase())) ?? false;

  exerciseToString = (exercise: Exercise) => exercise.names?.[0]?.name ?? '';

  onExerciseSelected(exercise: Exercise | null) {
    this.value.update((v) => ({ ...v, exerciseId: exercise?.id ?? null }));
  }

  onSchemeSelected(schemeId: number | null) {
    this.value.update((v) => ({ ...v, selectedSchemeId: schemeId }));
  }

  openDialog(scheme: ExerciseScheme | null) {
    this.editingScheme.set(scheme);
    this.dialogOpen.set(true);
  }

  closeDialog() {
    this.dialogOpen.set(false);
    this.editingScheme.set(null);
  }

  onSchemeSaved(scheme: ExerciseScheme) {
    this.value.update((v) => ({ ...v, selectedSchemeId: scheme.id }));
    this.closeDialog();
  }
}
