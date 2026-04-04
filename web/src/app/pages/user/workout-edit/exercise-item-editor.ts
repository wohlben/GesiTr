import { Component, computed, input, signal } from '@angular/core';
import { FormField, type FieldTree } from '@angular/forms/signals';
import { TranslocoDirective } from '@jsverse/transloco';
import { HlmComboboxImports } from '@spartan-ng/helm/combobox';
import { Exercise } from '$generated/models';
import { ExerciseScheme } from '$generated/user-exercisescheme';
import { SchemeSelector } from '$ui/scheme-selector/scheme-selector';
import { CreateSchemeDialog } from './create-scheme-dialog';
import type { WorkoutItemModel } from './workout-item-model';

@Component({
  selector: 'app-exercise-item-editor',
  imports: [FormField, HlmComboboxImports, SchemeSelector, CreateSchemeDialog, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <div>
        <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
          t('ui.exerciseConfig.exerciseLabel')
        }}</span>
        <hlm-combobox
          class="mt-1 block"
          [formField]="itemField().exerciseId"
          [filter]="exerciseIdFilter"
          [itemToString]="exerciseIdToString"
        >
          <hlm-combobox-input
            [placeholder]="t('common.search')"
            [showClear]="!!itemField().exerciseId().value()"
          />
          <ng-template hlmComboboxPortal>
            <hlm-combobox-content>
              <hlm-combobox-input [placeholder]="t('common.search')" [showClear]="false" />
              <div hlmComboboxList>
                @for (ex of exercises(); track ex.id) {
                  <hlm-combobox-item [value]="ex.id">{{ ex.names?.[0]?.name }}</hlm-combobox-item>
                }
                <hlm-combobox-empty>{{ t('common.noResults') }}</hlm-combobox-empty>
              </div>
            </hlm-combobox-content>
          </ng-template>
        </hlm-combobox>
      </div>

      <app-scheme-selector
        [exerciseId]="itemField().exerciseId().value()"
        [selectedSchemeId]="itemField().selectedSchemeId().value()"
        (schemeSelected)="onSchemeSelected($event)"
        (createRequested)="openDialog(null)"
        (editRequested)="openDialog($event)"
      />

      <app-create-scheme-dialog
        [open]="dialogOpen()"
        [preselectedExerciseId]="itemField().exerciseId().value()"
        [editingScheme]="editingScheme()"
        (schemeSaved)="onSchemeSaved($event)"
        (cancelled)="closeDialog()"
      />
    </ng-container>
  `,
})
export class ExerciseItemEditor {
  itemField = input.required<FieldTree<WorkoutItemModel>>();
  exercises = input.required<Exercise[]>();

  dialogOpen = signal(false);
  editingScheme = signal<ExerciseScheme | null>(null);

  private exerciseMap = computed(() => new Map(this.exercises().map((e) => [e.id, e])));

  exerciseIdFilter = (id: number, search: string) =>
    this.exerciseMap()
      .get(id)
      ?.names?.some((n) => n.name.toLowerCase().includes(search.toLowerCase())) ?? false;

  exerciseIdToString = (id: number) => this.exerciseMap().get(id)?.names?.[0]?.name ?? '';

  onSchemeSelected(schemeId: number | null) {
    this.itemField().selectedSchemeId().value.set(schemeId);
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
    this.itemField().selectedSchemeId().value.set(scheme.id);
    this.closeDialog();
  }
}
