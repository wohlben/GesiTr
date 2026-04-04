import { Component, input, output } from '@angular/core';
import { FormField, type FieldTree } from '@angular/forms/signals';
import { TranslocoDirective } from '@jsverse/transloco';
import { Exercise } from '$generated/models';
import { ExerciseGroup } from '$generated/user-models';
import {
  WorkoutSectionItemTypeExercise,
  WorkoutSectionItemTypeExerciseGroup,
} from '$generated/user-models';
import { ExerciseGroupConfig } from '$ui/exercise-group-config/exercise-group-config';
import type { WorkoutItemModel } from './workout-item-model';
import { ExerciseItemEditor } from './exercise-item-editor';

@Component({
  selector: 'app-section-item-editor',
  imports: [FormField, ExerciseItemEditor, ExerciseGroupConfig, TranslocoDirective],
  template: `
    <div
      *transloco="let t"
      class="rounded-md border border-gray-100 bg-gray-50 p-3 dark:border-gray-600 dark:bg-gray-800/50"
    >
      <div class="mb-2 flex items-center justify-between">
        <span class="text-xs font-medium text-gray-600 dark:text-gray-400">
          {{ itemLabel() }}
        </span>
        <button
          type="button"
          (click)="removed.emit()"
          class="text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
        >
          {{ t('common.remove') }}
        </button>
      </div>

      <!-- Item type selector -->
      <div class="mb-2">
        <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
          t('fields.itemType')
        }}</span>
        <div
          class="mt-1 flex overflow-hidden rounded-md border border-gray-300 dark:border-gray-600"
        >
          <button
            type="button"
            (click)="itemField().itemType().value.set(ITEM_TYPE_EXERCISE)"
            class="flex-1 px-3 py-1.5 text-sm font-medium transition-colors"
            [class]="
              itemField().itemType().value() === ITEM_TYPE_EXERCISE
                ? 'bg-blue-600 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
            "
          >
            {{ t('enums.workoutSectionItemType.exercise') }}
          </button>
          <button
            type="button"
            (click)="itemField().itemType().value.set(ITEM_TYPE_GROUP)"
            class="flex-1 border-l border-gray-300 px-3 py-1.5 text-sm font-medium transition-colors dark:border-gray-600"
            [class]="
              itemField().itemType().value() === ITEM_TYPE_GROUP
                ? 'bg-blue-600 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
            "
          >
            {{ t('enums.workoutSectionItemType.exercise_group') }}
          </button>
        </div>
      </div>

      @if (itemField().itemType().value() === ITEM_TYPE_EXERCISE) {
        <app-exercise-item-editor [itemField]="itemField()" [exercises]="exercises()" />
      }

      @if (itemField().itemType().value() === ITEM_TYPE_GROUP) {
        <app-exercise-group-config
          [formField]="itemField().groupConfig"
          [existingGroups]="exerciseGroups()"
          [exercises]="exercises()"
        />
      }
    </div>
  `,
})
export class SectionItemEditor {
  itemField = input.required<FieldTree<WorkoutItemModel>>();

  exercises = input.required<Exercise[]>();
  exerciseGroups = input.required<ExerciseGroup[]>();
  itemLabel = input.required<string>();

  removed = output<void>();

  readonly ITEM_TYPE_EXERCISE = WorkoutSectionItemTypeExercise;
  readonly ITEM_TYPE_GROUP = WorkoutSectionItemTypeExerciseGroup;
}
