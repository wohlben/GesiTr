import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { form, required, FormField } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import {
  workoutKeys,
  exerciseKeys,
  exerciseSchemeKeys,
  exerciseGroupKeys,
  masteryKeys,
} from '$core/query-keys';
import {
  WorkoutSectionTypeMain,
  WorkoutSectionTypeSupplementary,
  WorkoutSectionItemTypeExercise,
  WorkoutSectionItemTypeExerciseGroup,
} from '$generated/user-models';
import { ExerciseScheme } from '$generated/user-exercisescheme';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';
import { SchemeSelector } from '$ui/scheme-selector/scheme-selector';
import { CreateSchemeDialog } from './create-scheme-dialog';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmComboboxImports } from '@spartan-ng/helm/combobox';
import { HlmInput } from '@spartan-ng/helm/input';
import { Exercise } from '$generated/models';
import { ExerciseStore } from '$core/exercise.store';
import { HlmTextarea } from '@spartan-ng/helm/textarea';
import {
  ExerciseGroupConfig,
  EMPTY_GROUP_CONFIG,
} from '$ui/exercise-group-config/exercise-group-config';
import type { GroupConfigValue } from '$ui/exercise-group-config/exercise-group-config';

interface WorkoutItemModel {
  itemType: string;
  // exercise fields (used when itemType === 'exercise')
  exerciseId: number | null;
  selectedSchemeId: number | null;
  // exercise_group fields (used when itemType === 'exercise_group')
  groupConfig: GroupConfigValue;
}

interface WorkoutSectionModel {
  type: string;
  label: string;
  restBetweenExercises: number | null;
  items: WorkoutItemModel[];
}

interface WorkoutModel {
  name: string;
  notes: string;
  sections: WorkoutSectionModel[];
}

const EMPTY_EXERCISE_ITEM: WorkoutItemModel = {
  itemType: 'exercise',
  exerciseId: null,
  selectedSchemeId: null,
  groupConfig: { ...EMPTY_GROUP_CONFIG },
};

const EMPTY_GROUP_ITEM: WorkoutItemModel = {
  itemType: 'exercise_group',
  exerciseId: null,
  selectedSchemeId: null,
  groupConfig: { ...EMPTY_GROUP_CONFIG },
};

@Component({
  selector: 'app-workout-edit',
  imports: [
    PageLayout,
    FormField,
    RouterLink,
    ConfirmDialog,
    BrnSelectImports,
    HlmSelectImports,
    HlmComboboxImports,
    HlmInput,
    HlmTextarea,
    TranslocoDirective,
    ExerciseGroupConfig,
    SchemeSelector,
    CreateSchemeDialog,
  ],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="isCreateMode() ? t('user.workouts.newTitle') : t('user.workouts.editTitle')"
        [isPending]="!isCreateMode() && workoutQuery.isPending()"
        [errorMessage]="
          !isCreateMode() && workoutQuery.isError() ? workoutQuery.error().message : undefined
        "
      >
        @if (workoutQuery.data()?.workoutGroup; as group) {
          <div
            class="rounded-lg border px-4 py-3 text-sm"
            [class]="
              group.membership === 'admin'
                ? 'border-purple-200 bg-purple-50 text-purple-800 dark:border-purple-800 dark:bg-purple-900/20 dark:text-purple-300'
                : group.membership === 'member'
                  ? 'border-green-200 bg-green-50 text-green-800 dark:border-green-800 dark:bg-green-900/20 dark:text-green-300'
                  : 'border-yellow-200 bg-yellow-50 text-yellow-800 dark:border-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-300'
            "
          >
            <span class="font-medium">{{ group.groupName }}</span>
            —
            {{ t('enums.workoutGroupRole.' + group.membership) }}
          </div>
        }

        @if (isCreateMode() || workoutQuery.data()) {
          <form (submit)="onSubmit(); $event.preventDefault()" class="space-y-6">
            <!-- Basic Fields -->
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.name') }} *</label
              >
              <input id="name" [formField]="workoutForm.name" hlmInput class="mt-1" />
            </div>

            <div>
              <label
                for="notes"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.notes') }}</label
              >
              <textarea
                id="notes"
                [formField]="workoutForm.notes"
                rows="2"
                hlmTextarea
                class="mt-1"
              ></textarea>
            </div>

            <!-- Sections -->
            <div class="space-y-4">
              @for (section of workoutForm.sections; track $index; let si = $index) {
                <div class="rounded-lg border border-gray-200 p-4 dark:border-gray-700">
                  <div class="mb-3 flex items-center justify-between">
                    <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">
                      {{ t('user.workouts.sectionLabel', { n: si + 1 }) }}
                    </h3>
                    <button
                      type="button"
                      (click)="removeSection(si)"
                      class="text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                    >
                      {{ t('common.remove') }}
                    </button>
                  </div>

                  <div class="mb-3 grid grid-cols-1 gap-3 sm:grid-cols-3">
                    <div>
                      <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
                        t('fields.type')
                      }}</span>
                      <brn-select [formField]="section.type" class="mt-1" hlm>
                        <hlm-select-trigger class="w-full">
                          <hlm-select-value />
                        </hlm-select-trigger>
                        <hlm-select-content>
                          <hlm-option [value]="SECTION_TYPE_MAIN">{{
                            t('enums.workoutSectionType.main')
                          }}</hlm-option>
                          <hlm-option [value]="SECTION_TYPE_SUPPLEMENTARY">{{
                            t('enums.workoutSectionType.supplementary')
                          }}</hlm-option>
                        </hlm-select-content>
                      </brn-select>
                    </div>
                    <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                      {{ t('fields.label') }}
                      <input [formField]="section.label" hlmInput class="mt-1" />
                    </label>
                    <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                      {{ t('fields.restBetweenExercises') }}
                      <input
                        type="number"
                        [formField]="section.restBetweenExercises"
                        hlmInput
                        class="mt-1"
                      />
                    </label>
                  </div>

                  <!-- Items in section -->
                  <div class="space-y-3">
                    @for (item of section.items; track $index; let ei = $index) {
                      <div
                        class="rounded-md border border-gray-100 bg-gray-50 p-3 dark:border-gray-600 dark:bg-gray-800/50"
                      >
                        <div class="mb-2 flex items-center justify-between">
                          <span class="text-xs font-medium text-gray-600 dark:text-gray-400">
                            {{ t('user.workouts.exerciseLabel', { n: ei + 1 }) }}
                          </span>
                          <button
                            type="button"
                            (click)="removeItem(si, ei)"
                            class="text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                          >
                            {{ t('common.remove') }}
                          </button>
                        </div>

                        <!-- Item type selector -->
                        <div class="mb-2">
                          <span
                            class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >{{ t('fields.itemType') }}</span
                          >
                          <brn-select [formField]="item.itemType" class="mt-1" hlm>
                            <hlm-select-trigger class="w-full">
                              <hlm-select-value />
                            </hlm-select-trigger>
                            <hlm-select-content>
                              <hlm-option [value]="ITEM_TYPE_EXERCISE">{{
                                t('enums.workoutSectionItemType.exercise')
                              }}</hlm-option>
                              <hlm-option [value]="ITEM_TYPE_GROUP">{{
                                t('enums.workoutSectionItemType.exercise_group')
                              }}</hlm-option>
                            </hlm-select-content>
                          </brn-select>
                        </div>

                        @if (item.itemType().value() === 'exercise') {
                          <!-- Exercise type: exercise picker + scheme selector -->
                          <div>
                            <span
                              class="block text-xs font-medium text-gray-700 dark:text-gray-300"
                              >{{ t('ui.exerciseConfig.exerciseLabel') }}</span
                            >
                            <hlm-combobox
                              class="mt-1 block"
                              [value]="findExerciseById(item.exerciseId().value())"
                              (valueChange)="onItemExerciseSelected(si, ei, $event)"
                              [filter]="exerciseFilter"
                              [itemToString]="exerciseToString"
                            >
                              <hlm-combobox-input
                                [placeholder]="t('common.search')"
                                [showClear]="!!item.exerciseId().value()"
                              />
                              <ng-template hlmComboboxPortal>
                                <hlm-combobox-content>
                                  <hlm-combobox-input
                                    [placeholder]="t('common.search')"
                                    [showClear]="false"
                                  />
                                  <div hlmComboboxList>
                                    @for (ue of enrichedUserExercises(); track ue.id) {
                                      <hlm-combobox-item [value]="ue">{{
                                        ue.names?.[0]?.name
                                      }}</hlm-combobox-item>
                                    }
                                    <hlm-combobox-empty>{{
                                      t('common.noResults')
                                    }}</hlm-combobox-empty>
                                  </div>
                                </hlm-combobox-content>
                              </ng-template>
                            </hlm-combobox>
                          </div>

                          <app-scheme-selector
                            [exerciseId]="item.exerciseId().value()"
                            [selectedSchemeId]="item.selectedSchemeId().value()"
                            (schemeSelected)="onSchemeSelected($event, si, ei)"
                            (createRequested)="
                              openCreateSchemeDialog(si, ei, item.exerciseId().value())
                            "
                          />
                        }

                        @if (item.itemType().value() === 'exercise_group') {
                          <app-exercise-group-config
                            [formField]="item.groupConfig"
                            [existingGroups]="exerciseGroups()"
                            [exercises]="enrichedUserExercises()"
                          />
                        }
                      </div>
                    }
                  </div>

                  <div class="mt-2 flex gap-2">
                    <button
                      type="button"
                      (click)="addItem(si, 'exercise')"
                      class="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
                    >
                      {{ t('user.workouts.addExercise') }}
                    </button>
                    <button
                      type="button"
                      (click)="addItem(si, 'exercise_group')"
                      class="text-sm text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
                    >
                      {{ t('user.workouts.addExerciseGroup') }}
                    </button>
                  </div>
                </div>
              }
            </div>

            <button
              type="button"
              (click)="addSection()"
              class="rounded-md border border-dashed border-gray-300 px-4 py-2 text-sm text-gray-600 hover:border-gray-400 hover:text-gray-800 dark:border-gray-600 dark:text-gray-400 dark:hover:border-gray-500 dark:hover:text-gray-300"
            >
              {{ t('user.workouts.addSection') }}
            </button>

            <!-- Actions -->
            <div class="flex gap-2">
              <button
                type="submit"
                [disabled]="!workoutForm().valid() || saveMutation.isPending()"
                class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {{ t('common.save') }}
              </button>
              <a
                routerLink="/compendium/workouts"
                class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
              >
                {{ t('common.cancel') }}
              </a>
              @if (!isCreateMode()) {
                <button
                  type="button"
                  (click)="showDeleteDialog = true"
                  class="ml-auto rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
                >
                  {{ t('common.delete') }}
                </button>
              }
            </div>
          </form>
        }

        <app-confirm-dialog
          [open]="showDeleteDialog"
          [title]="t('user.workouts.deleteTitle')"
          [message]="t('user.workouts.deleteMessage')"
          [isPending]="deleteMutation.isPending()"
          (confirmed)="onDelete()"
          (cancelled)="showDeleteDialog = false"
        />

        <app-create-scheme-dialog
          [open]="createSchemeDialogState().open"
          [preselectedExerciseId]="createSchemeDialogState().exerciseId"
          (schemeCreated)="onSchemeCreated($event)"
          (cancelled)="closeCreateSchemeDialog()"
        />
      </app-page-layout>
    </ng-container>
  `,
})
export class WorkoutEdit {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private exerciseStore = inject(ExerciseStore);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = inject(QueryClient);
  private params = toSignal(this.route.paramMap);

  readonly SECTION_TYPE_MAIN = WorkoutSectionTypeMain;
  readonly SECTION_TYPE_SUPPLEMENTARY = WorkoutSectionTypeSupplementary;
  readonly ITEM_TYPE_EXERCISE = WorkoutSectionItemTypeExercise;
  readonly ITEM_TYPE_GROUP = WorkoutSectionItemTypeExerciseGroup;

  private id = computed(() => Number(this.params()?.get('id')));
  isCreateMode = computed(() => !this.params()?.get('id'));

  showDeleteDialog = false;

  model = signal<WorkoutModel>({ name: '', notes: '', sections: [] });
  workoutForm = form(this.model, (f) => {
    required(f.name);
  });

  // Existing workout data for edit mode
  workoutQuery = injectQuery(() => ({
    queryKey: workoutKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchWorkout(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  // All exercises for the picker dropdown
  private allExercisesQuery = injectQuery(() => ({
    queryKey: exerciseKeys.list({ limit: 1000 }),
    queryFn: () => this.compendiumApi.fetchExercises({ limit: 1000 }),
  }));

  // Sync fetched exercises into the global store
  private syncToStore = effect(() => {
    const items = this.allExercisesQuery.data()?.items;
    if (items) this.exerciseStore.setAllFromQuery(items);
  });

  private masteryQuery = injectQuery(() => ({
    queryKey: masteryKeys.list(),
    queryFn: () => this.userApi.fetchMasteryList(),
  }));

  // Sorted with mastery exercises first
  enrichedUserExercises = computed(() => {
    const all = this.allExercisesQuery.data()?.items ?? [];
    const masteryIds = new Set((this.masteryQuery.data() ?? []).map((m) => m.exerciseId));
    return [...all].sort((a, b) => {
      const aHas = masteryIds.has(a.id) ? 0 : 1;
      const bHas = masteryIds.has(b.id) ? 0 : 1;
      if (aHas !== bHas) return aHas - bHas;
      return (a.names?.[0]?.name ?? '').localeCompare(b.names?.[0]?.name ?? '');
    });
  });

  exerciseFilter = (exercise: Exercise, search: string) =>
    exercise.names?.some((n) => n.name.toLowerCase().includes(search.toLowerCase())) ?? false;

  exerciseToString = (exercise: Exercise) => exercise.names?.[0]?.name ?? '';

  findExerciseById(id: number | null): Exercise | null {
    if (!id) return null;
    return this.enrichedUserExercises().find((e) => e.id === id) ?? null;
  }

  onItemExerciseSelected(si: number, ei: number, exercise: Exercise | null) {
    this.model.update((m) => {
      const sections = [...m.sections];
      const items = [...sections[si].items];
      items[ei] = { ...items[ei], exerciseId: exercise?.id ?? null };
      sections[si] = { ...sections[si], items };
      return { ...m, sections };
    });
  }

  // Exercise groups for the group picker
  private exerciseGroupsQuery = injectQuery(() => ({
    queryKey: exerciseGroupKeys.list({}),
    queryFn: () => this.compendiumApi.fetchExerciseGroups({}),
  }));

  exerciseGroups = computed(() => this.exerciseGroupsQuery.data()?.items ?? []);

  createSchemeDialogState = signal<{
    open: boolean;
    sectionIndex: number;
    itemIndex: number;
    exerciseId: number | null;
  }>({ open: false, sectionIndex: 0, itemIndex: 0, exerciseId: null });

  saveMutation = injectMutation(() => ({
    mutationFn: () => this.saveWorkout(),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutKeys.all() });
      this.queryClient.invalidateQueries({ queryKey: exerciseSchemeKeys.all() });
      this.router.navigate(['/compendium/workouts']);
    },
  }));

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.deleteWorkout(),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutKeys.all() });
      this.queryClient.invalidateQueries({ queryKey: exerciseSchemeKeys.all() });
      this.router.navigate(['/compendium/workouts']);
    },
  }));

  constructor() {
    // Populate form in edit mode
    effect(() => {
      const data = this.workoutQuery.data();
      if (!data) return;

      this.model.set({
        name: data.name,
        notes: data.notes ?? '',
        sections: (data.sections ?? []).map((section) => ({
          type: section.type,
          label: section.label ?? '',
          restBetweenExercises: section.restBetweenExercises ?? null,
          items: (section.items ?? []).map((item) => ({
            itemType: item.type ?? 'exercise',
            exerciseId: item.exerciseId ?? null,
            selectedSchemeId: null,
            groupConfig: {
              exerciseGroupId: item.exerciseGroupId ?? null,
              name: '',
              members: [],
            },
          })),
        })),
      });

      this.loadItemDetails(data.sections ?? []);
    });
  }

  private async loadItemDetails(
    sections: NonNullable<ReturnType<typeof this.workoutQuery.data>>['sections'],
  ) {
    // Batch-fetch scheme assignments via the join table
    const exerciseItemIds = sections
      .flatMap((s) => s.items ?? [])
      .filter((i) => i.type === 'exercise' && i.exerciseId != null)
      .map((i) => i.id);

    const schemeAssignments =
      exerciseItemIds.length > 0 ? await this.userApi.fetchSchemeSectionItems(exerciseItemIds) : [];
    const assignmentByItemId = new Map(schemeAssignments.map((a) => [a.workoutSectionItemId, a]));

    for (let si = 0; si < sections.length; si++) {
      const section = sections[si];
      for (let ei = 0; ei < (section.items?.length ?? 0); ei++) {
        const item = section.items[ei];

        if (item.exerciseId != null) {
          const assignment = assignmentByItemId.get(item.id);
          if (assignment) {
            this.model.update((m) => ({
              ...m,
              sections: m.sections.map((s, sIdx) =>
                sIdx !== si
                  ? s
                  : {
                      ...s,
                      items: s.items.map((e, eIdx) =>
                        eIdx !== ei ? e : { ...e, selectedSchemeId: assignment.exerciseSchemeId },
                      ),
                    },
              ),
            }));
          }
        }

        if (item.exerciseGroupId) {
          // Load group details (name + members) for exercise_group items
          try {
            const [group, members] = await Promise.all([
              this.compendiumApi.fetchExerciseGroup(item.exerciseGroupId),
              this.compendiumApi.fetchExerciseGroupMembers({ groupId: item.exerciseGroupId }),
            ]);
            this.model.update((m) => ({
              ...m,
              sections: m.sections.map((s, sIdx) =>
                sIdx !== si
                  ? s
                  : {
                      ...s,
                      items: s.items.map((e, eIdx) =>
                        eIdx !== ei
                          ? e
                          : {
                              ...e,
                              groupConfig: {
                                ...e.groupConfig,
                                name: group.name ?? '',
                                members: members.map((mem) => mem.exerciseId),
                              },
                            },
                      ),
                    },
              ),
            }));
          } catch {
            // group may have been deleted
          }
        }
      }
    }
  }

  addSection() {
    this.model.update((m) => ({
      ...m,
      sections: [
        ...m.sections,
        {
          type: WorkoutSectionTypeMain,
          label: '',
          restBetweenExercises: null,
          items: [],
        },
      ],
    }));
  }

  removeSection(index: number) {
    this.model.update((m) => ({
      ...m,
      sections: m.sections.filter((_, i) => i !== index),
    }));
  }

  addItem(sectionIndex: number, itemType: 'exercise' | 'exercise_group') {
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, i) =>
        i !== sectionIndex
          ? s
          : {
              ...s,
              items: [
                ...s.items,
                itemType === 'exercise' ? { ...EMPTY_EXERCISE_ITEM } : { ...EMPTY_GROUP_ITEM },
              ],
            },
      ),
    }));
  }

  removeItem(sectionIndex: number, itemIndex: number) {
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, i) =>
        i !== sectionIndex
          ? s
          : {
              ...s,
              items: s.items.filter((_, j) => j !== itemIndex),
            },
      ),
    }));
  }

  onSchemeSelected(schemeId: number | null, si: number, ei: number) {
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, sIdx) =>
        sIdx !== si
          ? s
          : {
              ...s,
              items: s.items.map((e, eIdx) =>
                eIdx !== ei ? e : { ...e, selectedSchemeId: schemeId },
              ),
            },
      ),
    }));
  }

  openCreateSchemeDialog(si: number, ei: number, exerciseId: number | null) {
    this.createSchemeDialogState.set({ open: true, sectionIndex: si, itemIndex: ei, exerciseId });
  }

  closeCreateSchemeDialog() {
    this.createSchemeDialogState.update((s) => ({ ...s, open: false }));
  }

  onSchemeCreated(scheme: ExerciseScheme) {
    const { sectionIndex, itemIndex } = this.createSchemeDialogState();
    this.onSchemeSelected(scheme.id, sectionIndex, itemIndex);
    this.closeCreateSchemeDialog();
  }

  onSubmit() {
    if (this.workoutForm().valid()) {
      this.saveMutation.mutate();
    }
  }

  private async saveWorkout() {
    const val = this.model();

    if (this.isCreateMode()) {
      return this.createFlow(val);
    } else {
      return this.editFlow(val);
    }
  }

  private async createFlow(val: WorkoutModel) {
    const workout = await this.userApi.createWorkout({
      name: val.name,
      notes: val.notes || undefined,
    });

    await this.createSectionsAndItems(workout.id, val.sections);
  }

  private async editFlow(val: WorkoutModel) {
    const workoutId = this.id();
    const existingWorkout = this.workoutQuery.data()!;

    await this.userApi.updateWorkout(workoutId, {
      name: val.name,
      notes: val.notes || undefined,
    });

    // Delete existing section items, then sections
    for (const section of existingWorkout.sections ?? []) {
      for (const item of section.items ?? []) {
        await this.userApi.deleteWorkoutSectionItem(item.id);
      }
      await this.userApi.deleteWorkoutSection(section.id);
    }

    await this.createSectionsAndItems(workoutId, val.sections);
  }

  private async createSectionsAndItems(workoutId: number, sections: WorkoutSectionModel[]) {
    for (let si = 0; si < sections.length; si++) {
      const sectionVal = sections[si];

      // Create section
      const section = await this.userApi.createWorkoutSection({
        workoutId,
        type: sectionVal.type,
        label: sectionVal.label || undefined,
        position: si,
        restBetweenExercises: sectionVal.restBetweenExercises ?? undefined,
      });

      for (let ei = 0; ei < sectionVal.items.length; ei++) {
        const item = sectionVal.items[ei];

        if (item.itemType === 'exercise') {
          const sectionItem = await this.userApi.createWorkoutSectionItem({
            workoutSectionId: section.id,
            type: 'exercise',
            exerciseId: item.exerciseId ?? undefined,
            position: ei,
          });

          // Link the selected scheme (if any) to the section item
          if (item.selectedSchemeId != null) {
            await this.userApi.upsertSchemeSectionItem({
              exerciseSchemeId: item.selectedSchemeId,
              workoutSectionItemId: sectionItem.id,
            });
          }
        } else {
          // Exercise group items — create or update group
          const gc = item.groupConfig;
          let groupId = gc.exerciseGroupId;
          if (groupId == null) {
            // Create new group
            const group = await this.compendiumApi.createExerciseGroup({
              name: gc.name || undefined,
            });
            groupId = group.id;
            for (const exerciseId of gc.members) {
              await this.compendiumApi.createExerciseGroupMember({
                groupId: group.id,
                exerciseId,
              });
            }
          } else {
            // Update existing group: sync name and members
            await this.compendiumApi.updateExerciseGroup(groupId, {
              name: gc.name || undefined,
            });
            const existingMembers = await this.compendiumApi.fetchExerciseGroupMembers({
              groupId,
            });
            const wantedSet = new Set(gc.members);
            const existingMap = new Map(existingMembers.map((m) => [m.exerciseId, m.id]));
            // Delete removed members
            for (const m of existingMembers) {
              if (!wantedSet.has(m.exerciseId)) {
                await this.compendiumApi.deleteExerciseGroupMember(m.id);
              }
            }
            // Add new members
            for (const exerciseId of gc.members) {
              if (!existingMap.has(exerciseId)) {
                await this.compendiumApi.createExerciseGroupMember({
                  groupId,
                  exerciseId,
                });
              }
            }
          }
          await this.userApi.createWorkoutSectionItem({
            workoutSectionId: section.id,
            type: 'exercise_group',
            exerciseGroupId: groupId ?? undefined,
            position: ei,
          });
        }
      }
    }
  }

  onDelete() {
    this.deleteMutation.mutate();
  }

  private async deleteWorkout() {
    const workout = this.workoutQuery.data()!;

    // Delete section items and sections
    for (const section of workout.sections ?? []) {
      for (const item of section.items ?? []) {
        await this.userApi.deleteWorkoutSectionItem(item.id);
      }
      await this.userApi.deleteWorkoutSection(section.id);
    }

    await this.userApi.deleteWorkout(this.id());
  }
}
