import { Component, input, output, computed } from '@angular/core';
import { TranslocoDirective } from '@jsverse/transloco';
import { HlmInput } from '@spartan-ng/helm/input';
import { ExerciseGroup } from '$generated/models';

export interface ExerciseGroupConfigValue {
  exerciseGroupId: number | null;
  newGroupName: string;
  newGroupMembers: number[];
}

interface ExerciseOption {
  id: number;
  name: string;
}

@Component({
  selector: 'app-exercise-group-config',
  imports: [HlmInput, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <!-- Group selector: existing groups or "New Group" -->
      <div class="mb-2">
        <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
          t('fields.exerciseGroup')
        }}</span>
        <select
          hlmInput
          class="mt-1 w-full"
          [value]="selectedGroupId()"
          (change)="onGroupSelect($event)"
        >
          <option [value]="NEW_GROUP_SENTINEL">
            {{ t('ui.exerciseGroupConfig.newGroup') }}
          </option>
          @for (g of existingGroups(); track g.id) {
            <option [value]="g.id">
              {{ g.name || t('common.unnamedGroup', { id: g.id }) }}
            </option>
          }
        </select>
      </div>

      <!-- Group editing: always shown -->
      <div class="space-y-2">
        <!-- Group name (optional) -->
        <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
          {{ t('fields.name') }}
          <input
            hlmInput
            [value]="currentName()"
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
              <option [value]="ex.id">{{ ex.name }}</option>
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
                <span class="text-gray-900 dark:text-gray-100">{{ member.name }}</span>
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
export class ExerciseGroupConfig {
  readonly NEW_GROUP_SENTINEL = -1;

  existingGroups = input.required<ExerciseGroup[]>();
  exercises = input.required<ExerciseOption[]>();
  value = input.required<ExerciseGroupConfigValue>();
  valueChange = output<ExerciseGroupConfigValue>();

  selectedGroupId = computed(() => this.value().exerciseGroupId ?? this.NEW_GROUP_SENTINEL);
  currentName = computed(() => this.value().newGroupName);
  currentMembers = computed(() => this.value().newGroupMembers);

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

  onGroupSelect(event: Event) {
    const id = Number((event.target as HTMLSelectElement).value);
    if (id === this.NEW_GROUP_SENTINEL || isNaN(id)) {
      this.emit({ exerciseGroupId: null, newGroupName: '', newGroupMembers: [] });
    } else {
      this.emit({ exerciseGroupId: id, newGroupName: '', newGroupMembers: [] });
    }
  }

  onNameChange(event: Event) {
    const name = (event.target as HTMLInputElement).value;
    this.emit({ ...this.value(), newGroupName: name });
  }

  onAddMember(event: Event) {
    const select = event.target as HTMLSelectElement;
    const exerciseId = Number(select.value);
    if (isNaN(exerciseId) || !exerciseId || this.currentMembers().includes(exerciseId)) return;
    this.emit({
      ...this.value(),
      newGroupMembers: [...this.currentMembers(), exerciseId],
    });
    select.value = '';
  }

  onRemoveMember(exerciseId: number) {
    this.emit({
      ...this.value(),
      newGroupMembers: this.currentMembers().filter((id) => id !== exerciseId),
    });
  }

  private emit(val: ExerciseGroupConfigValue) {
    this.valueChange.emit(val);
  }
}
