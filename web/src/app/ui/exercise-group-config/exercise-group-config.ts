import { Component, model, input, computed } from '@angular/core';
import { TranslocoDirective } from '@jsverse/transloco';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmInput } from '@spartan-ng/helm/input';
import { ExerciseGroup } from '$generated/user-models';
import type { FormValueControl } from '@angular/forms/signals';

export interface GroupConfigValue {
  exerciseGroupId: number | null;
  name: string;
  members: number[];
}

export const EMPTY_GROUP_CONFIG: GroupConfigValue = {
  exerciseGroupId: null,
  name: '',
  members: [],
};

interface ExerciseOption {
  id: number;
  names?: { name: string }[];
}

@Component({
  selector: 'app-exercise-group-config',
  imports: [BrnSelectImports, HlmSelectImports, HlmInput, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <!-- Group selector -->
      <div class="mb-2">
        <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
          t('fields.exerciseGroup')
        }}</span>
        <brn-select
          [value]="value().exerciseGroupId"
          (valueChange)="onGroupSelect($event)"
          class="mt-1"
          hlm
        >
          <hlm-select-trigger class="w-full">
            <hlm-select-value />
          </hlm-select-trigger>
          <hlm-select-content>
            <hlm-option [value]="null">{{ t('ui.exerciseGroupConfig.newGroup') }}</hlm-option>
            @for (g of existingGroups(); track g.id) {
              <hlm-option [value]="g.id">
                {{ g.name || t('common.unnamedGroup', { id: g.id }) }}
              </hlm-option>
            }
          </hlm-select-content>
        </brn-select>
      </div>

      <!-- Group editing -->
      <div class="space-y-2">
        <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
          {{ t('fields.name') }}
          <input
            hlmInput
            [value]="value().name"
            (input)="onNameChange($event)"
            class="mt-1"
            [placeholder]="t('ui.exerciseGroupConfig.namePlaceholder')"
          />
        </label>

        <!-- Add member exercise -->
        <div>
          <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
            t('ui.exerciseGroupConfig.addMember')
          }}</span>
          <select hlmInput class="mt-1 w-full" (change)="onAddMember($event)">
            <option value="">{{ t('common.select') }}</option>
            @for (ex of availableExercises(); track ex.id) {
              <option [value]="ex.id">{{ ex.names?.[0]?.name }}</option>
            }
          </select>
        </div>

        <!-- Member list -->
        @if (memberDetails().length > 0) {
          <div class="space-y-1">
            <span class="block text-xs font-medium text-gray-500 dark:text-gray-400">{{
              t('ui.exerciseGroupConfig.members')
            }}</span>
            @for (member of memberDetails(); track member.id) {
              <div
                class="flex items-center justify-between rounded-md border border-gray-200 px-3 py-1.5 text-sm dark:border-gray-600"
              >
                <span class="text-gray-900 dark:text-gray-100">{{ member.names?.[0]?.name }}</span>
                <button
                  type="button"
                  (click)="onRemoveMember(member.id)"
                  class="text-red-500 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
                >
                  &times;
                </button>
              </div>
            }
          </div>
        }
      </div>
    </ng-container>
  `,
})
export class ExerciseGroupConfig implements FormValueControl<GroupConfigValue> {
  /** Synced with the parent form field via [formField]. */
  readonly value = model.required<GroupConfigValue>();

  existingGroups = input.required<ExerciseGroup[]>();
  exercises = input.required<ExerciseOption[]>();

  private currentMembers = computed(() => this.value().members);

  availableExercises = computed(() => {
    const memberSet = new Set(this.currentMembers());
    return this.exercises().filter((ex) => !memberSet.has(ex.id));
  });

  memberDetails = computed(() => {
    const exerciseMap = new Map(this.exercises().map((ex) => [ex.id, ex]));
    return this.currentMembers()
      .map((id) => exerciseMap.get(id))
      .filter((ex): ex is ExerciseOption => ex != null);
  });

  onGroupSelect(id: number | null | (number | null)[] | undefined) {
    const groupId = Array.isArray(id) ? (id[0] ?? null) : (id ?? null);
    const group = groupId != null ? this.existingGroups().find((g) => g.id === groupId) : null;
    this.value.set({ exerciseGroupId: groupId, name: group?.name ?? '', members: [] });
  }

  onNameChange(event: Event) {
    const name = (event.target as HTMLInputElement).value;
    this.value.update((v) => ({ ...v, name }));
  }

  onAddMember(event: Event) {
    const select = event.target as HTMLSelectElement;
    const exerciseId = Number(select.value);
    if (isNaN(exerciseId) || !exerciseId || this.currentMembers().includes(exerciseId)) return;
    this.value.update((v) => ({ ...v, members: [...v.members, exerciseId] }));
    select.value = '';
  }

  onRemoveMember(exerciseId: number) {
    this.value.update((v) => ({ ...v, members: v.members.filter((id) => id !== exerciseId) }));
  }
}
