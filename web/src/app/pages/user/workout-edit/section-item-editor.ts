import { Component, input, model, output } from '@angular/core';
import { TranslocoDirective } from '@jsverse/transloco';
import { Exercise } from '$generated/models';
import { ExerciseGroup } from '$generated/user-models';
import {
  WorkoutSectionItemTypeExercise,
  WorkoutSectionItemTypeExerciseGroup,
} from '$generated/user-models';
import { ExerciseGroupConfig } from '$ui/exercise-group-config/exercise-group-config';
import type { GroupConfigValue } from '$ui/exercise-group-config/exercise-group-config';
import type { FormValueControl } from '@angular/forms/signals';
import type { WorkoutItemModel } from './workout-item-model';
import { ExerciseItemEditor } from './exercise-item-editor';

@Component({
  selector: 'app-section-item-editor',
  imports: [ExerciseItemEditor, ExerciseGroupConfig, TranslocoDirective],
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
            (click)="setItemType(ITEM_TYPE_EXERCISE)"
            class="flex-1 px-3 py-1.5 text-sm font-medium transition-colors"
            [class]="
              value().itemType === ITEM_TYPE_EXERCISE
                ? 'bg-blue-600 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
            "
          >
            {{ t('enums.workoutSectionItemType.exercise') }}
          </button>
          <button
            type="button"
            (click)="setItemType(ITEM_TYPE_GROUP)"
            class="flex-1 border-l border-gray-300 px-3 py-1.5 text-sm font-medium transition-colors dark:border-gray-600"
            [class]="
              value().itemType === ITEM_TYPE_GROUP
                ? 'bg-blue-600 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
            "
          >
            {{ t('enums.workoutSectionItemType.exercise_group') }}
          </button>
        </div>
      </div>

      @if (value().itemType === ITEM_TYPE_EXERCISE) {
        <app-exercise-item-editor [(value)]="value" [exercises]="exercises()" />
      }

      @if (value().itemType === ITEM_TYPE_GROUP) {
        <app-exercise-group-config
          [value]="value().groupConfig"
          (valueChange)="onGroupConfigChange($event)"
          [existingGroups]="exerciseGroups()"
          [exercises]="exercises()"
        />
      }
    </div>
  `,
})
export class SectionItemEditor implements FormValueControl<WorkoutItemModel> {
  readonly value = model.required<WorkoutItemModel>();

  exercises = input.required<Exercise[]>();
  exerciseGroups = input.required<ExerciseGroup[]>();
  itemLabel = input.required<string>();

  removed = output<void>();

  readonly ITEM_TYPE_EXERCISE = WorkoutSectionItemTypeExercise;
  readonly ITEM_TYPE_GROUP = WorkoutSectionItemTypeExerciseGroup;

  setItemType(itemType: string) {
    this.value.update((v) => ({ ...v, itemType }));
  }

  onGroupConfigChange(groupConfig: GroupConfigValue) {
    this.value.update((v) => ({ ...v, groupConfig }));
  }
}
