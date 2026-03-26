import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { form, FormField } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CdkDragDrop, CdkDrag, CdkDropList, CdkDragHandle } from '@angular/cdk/drag-drop';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { formatBreak } from '$core/format-utils';
import { workoutKeys, workoutLogKeys } from '$core/query-keys';
import {
  Workout,
  WorkoutSectionTypeMain,
  WorkoutSectionTypeSupplementary,
  WorkoutLog,
} from '$generated/user-models';
import { ExerciseScheme } from '$generated/models';
import { PageLayout } from '../../../layout/page-layout';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmInput } from '@spartan-ng/helm/input';
import { HlmTextarea } from '@spartan-ng/helm/textarea';
import { WorkoutStartStore, SetPreview } from './workout-start.store';
import { AddExerciseDialog } from './add-exercise-dialog';

interface StartSetModel {
  id: number | null;
  targetReps: number | null;
  targetWeight: number | null;
  targetDuration: number | null;
  targetDistance: number | null;
  targetTime: number | null;
  restAfterSeconds: number | null;
}

interface StartExerciseModel {
  id: number | null;
  sourceExerciseSchemeId: number;
  breakAfterSeconds: number | null;
  sets: StartSetModel[];
}

interface PendingGroupModel {
  groupId: number;
  groupName: string;
  members: { id: number; name: string }[];
  position: number;
}

interface StartSectionModel {
  id: number | null;
  type: string;
  label: string;
  exercises: StartExerciseModel[];
  pendingGroups: PendingGroupModel[];
}

interface StartModel {
  name: string;
  notes: string;
  sections: StartSectionModel[];
}

@Component({
  selector: 'app-workout-start',
  imports: [
    PageLayout,
    FormField,
    RouterLink,
    AddExerciseDialog,
    CdkDropList,
    CdkDrag,
    CdkDragHandle,
    BrnSelectImports,
    HlmSelectImports,
    HlmInput,
    HlmTextarea,
    TranslocoDirective,
  ],
  providers: [WorkoutStartStore],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('user.workoutStart.title')"
        [isPending]="isPending()"
        [errorMessage]="workoutQuery.isError() ? workoutQuery.error().message : undefined"
      >
        @if (workoutQuery.data() && currentLogId()) {
          <form (submit)="onSubmit(); $event.preventDefault()" class="space-y-6">
            <!-- Basic Fields -->
            <div>
              <label
                for="name"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.name') }}</label
              >
              <input
                hlmInput
                id="name"
                [formField]="startForm.name"
                (change)="onLogChange()"
                class="mt-1"
              />
            </div>

            <div>
              <label
                for="notes"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.notes') }}</label
              >
              <textarea
                hlmTextarea
                id="notes"
                [formField]="startForm.notes"
                (change)="onLogChange()"
                rows="2"
                class="mt-1"
              ></textarea>
            </div>

            <!-- Sections -->
            <div
              cdkDropList
              [cdkDropListData]="startForm.sections"
              (cdkDropListDropped)="onSectionDrop($event)"
              class="space-y-4"
            >
              @for (section of startForm.sections; track $index; let si = $index) {
                <div cdkDrag class="rounded-lg border border-gray-200 p-4 dark:border-gray-700">
                  <div class="mb-3 flex items-center justify-between">
                    <div class="flex items-center gap-2">
                      <!-- Section drag handle -->
                      <div
                        cdkDragHandle
                        class="flex cursor-grab flex-col gap-0.5 px-1 py-1 text-gray-400 active:cursor-grabbing dark:text-gray-500"
                      >
                        <div class="flex gap-0.5">
                          <div class="h-1 w-1 rounded-full bg-current"></div>
                          <div class="h-1 w-1 rounded-full bg-current"></div>
                        </div>
                        <div class="flex gap-0.5">
                          <div class="h-1 w-1 rounded-full bg-current"></div>
                          <div class="h-1 w-1 rounded-full bg-current"></div>
                        </div>
                        <div class="flex gap-0.5">
                          <div class="h-1 w-1 rounded-full bg-current"></div>
                          <div class="h-1 w-1 rounded-full bg-current"></div>
                        </div>
                      </div>
                      <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">
                        {{ t('user.workouts.sectionLabel', { n: si + 1 }) }}
                      </h3>
                    </div>
                    <button
                      type="button"
                      (click)="removeSection(si)"
                      class="text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                    >
                      {{ t('common.remove') }}
                    </button>
                  </div>

                  <div class="mb-3 grid grid-cols-1 gap-3 sm:grid-cols-2">
                    <div>
                      <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
                        t('fields.type')
                      }}</span>
                      <brn-select
                        [formField]="section.type"
                        (valueChange)="onSectionChange(si)"
                        class="mt-1"
                        hlm
                      >
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
                      <input
                        hlmInput
                        [formField]="section.label"
                        (change)="onSectionChange(si)"
                        class="mt-1"
                      />
                    </label>
                  </div>

                  <!-- Exercise cards -->
                  <div
                    cdkDropList
                    [cdkDropListData]="section.exercises"
                    (cdkDropListDropped)="onExerciseDrop($event, si)"
                  >
                    @for (
                      exercise of section.exercises;
                      track $index;
                      let ei = $index;
                      let lastEx = $last
                    ) {
                      <div cdkDrag>
                        @let info = store.exerciseDisplay()[exercise.id().value()!];
                        <div class="rounded-md border border-gray-200 dark:border-gray-600">
                          <!-- Exercise header -->
                          <div
                            class="flex items-center justify-between border-b border-gray-100 px-3 py-2 dark:border-gray-700"
                          >
                            <div
                              class="flex items-center gap-2 text-sm text-gray-900 dark:text-gray-100"
                            >
                              <!-- Exercise drag handle -->
                              <div
                                cdkDragHandle
                                class="flex cursor-grab flex-col gap-0.5 text-gray-400 active:cursor-grabbing dark:text-gray-500"
                              >
                                <div class="flex gap-0.5">
                                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                                </div>
                                <div class="flex gap-0.5">
                                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                                </div>
                                <div class="flex gap-0.5">
                                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                                  <div class="h-0.5 w-0.5 rounded-full bg-current"></div>
                                </div>
                              </div>
                              <div>
                                <span class="font-semibold">{{
                                  info?.name ?? t('common.loading')
                                }}</span>
                                <span class="ml-2 text-gray-500 dark:text-gray-400">{{
                                  info?.summary
                                }}</span>
                              </div>
                            </div>
                            <button
                              type="button"
                              (click)="removeExercise(si, ei)"
                              class="text-xs text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                            >
                              {{ t('common.remove') }}
                            </button>
                          </div>

                          <!-- Editable sets -->
                          @if (exercise.sets.length) {
                            <div class="px-3 py-2">
                              <!-- Header -->
                              <div
                                class="mb-1 grid text-left text-xs text-gray-500 uppercase dark:text-gray-400"
                                [class]="
                                  info?.measurementType === 'REP_BASED'
                                    ? 'grid-cols-[2rem_5rem_6rem]'
                                    : 'grid-cols-[2rem_6rem]'
                                "
                              >
                                <span>{{ t('fields.set') }}</span>
                                @if (info?.measurementType === 'REP_BASED') {
                                  <span>{{ t('fields.reps') }}</span>
                                  <span>{{ t('fields.weight') }}</span>
                                }
                                @if (info?.measurementType === 'TIME_BASED') {
                                  <span>{{ t('fields.duration') }}</span>
                                }
                                @if (info?.measurementType === 'DISTANCE_BASED') {
                                  <span>{{ t('fields.distance') }}</span>
                                }
                              </div>

                              @for (
                                set of exercise.sets;
                                track $index;
                                let setIdx = $index;
                                let lastSet = $last
                              ) {
                                <!-- Set row -->
                                <div
                                  class="grid items-center py-1.5"
                                  [class]="
                                    info?.measurementType === 'REP_BASED'
                                      ? 'grid-cols-[2rem_5rem_6rem]'
                                      : 'grid-cols-[2rem_6rem]'
                                  "
                                >
                                  <span
                                    class="text-sm font-medium text-gray-900 dark:text-gray-100"
                                    >{{ setIdx + 1 }}</span
                                  >
                                  @if (info?.measurementType === 'REP_BASED') {
                                    <div>
                                      <input
                                        hlmInput
                                        type="number"
                                        [formField]="set.targetReps"
                                        data-field="targetReps"
                                        (change)="onSetChange(si, ei, setIdx)"
                                        class="mt-1"
                                      />
                                    </div>
                                    <div>
                                      <input
                                        hlmInput
                                        type="number"
                                        [formField]="set.targetWeight"
                                        data-field="targetWeight"
                                        (change)="onSetChange(si, ei, setIdx)"
                                        class="mt-1"
                                        step="0.5"
                                      />
                                    </div>
                                  }
                                  @if (info?.measurementType === 'TIME_BASED') {
                                    <div>
                                      <input
                                        hlmInput
                                        type="number"
                                        [formField]="set.targetDuration"
                                        (change)="onSetChange(si, ei, setIdx)"
                                        class="mt-1"
                                      />
                                    </div>
                                  }
                                  @if (info?.measurementType === 'DISTANCE_BASED') {
                                    <div>
                                      <input
                                        hlmInput
                                        type="number"
                                        [formField]="set.targetDistance"
                                        (change)="onSetChange(si, ei, setIdx)"
                                        class="mt-1"
                                        step="0.1"
                                      />
                                    </div>
                                  }
                                </div>

                                <!-- Rest between sets: line with centered badge -->
                                @if (!lastSet && set.restAfterSeconds().value() !== null) {
                                  <div class="relative flex items-center justify-center py-0.5">
                                    <div
                                      class="absolute inset-x-0 top-1/2 border-t border-dashed border-gray-200 dark:border-gray-700"
                                    ></div>
                                    <div
                                      class="relative z-10 flex items-center gap-1 rounded-full bg-gray-50 px-2 py-0.5 text-xs text-gray-400 dark:bg-gray-900 dark:text-gray-500"
                                    >
                                      <svg class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
                                        <path
                                          fill-rule="evenodd"
                                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z"
                                          clip-rule="evenodd"
                                        />
                                      </svg>
                                      <input
                                        type="number"
                                        [formField]="set.restAfterSeconds"
                                        data-field="restAfterSeconds"
                                        (change)="onSetChange(si, ei, setIdx)"
                                        class="w-12 border-0 bg-transparent p-0 text-center text-xs text-gray-400 focus:ring-0 dark:text-gray-500"
                                      />
                                      <span>{{ t('common.unitSeconds') }}</span>
                                    </div>
                                  </div>
                                }
                              }
                            </div>
                          }
                        </div>

                        <!-- Break after exercise (editable, not shown after last) -->
                        @if (!lastEx) {
                          <div class="relative flex items-center justify-center py-3">
                            <div
                              class="absolute inset-x-0 top-1/2 border-t border-gray-200 dark:border-gray-700"
                            ></div>
                            <div
                              class="relative z-10 flex items-center gap-1.5 rounded-full bg-white px-3 py-1 text-xs text-gray-500 shadow-sm ring-1 ring-gray-200 dark:bg-gray-800 dark:text-gray-400 dark:ring-gray-600"
                            >
                              <svg class="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                                <path
                                  fill-rule="evenodd"
                                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z"
                                  clip-rule="evenodd"
                                />
                              </svg>
                              <input
                                type="number"
                                [formField]="exercise.breakAfterSeconds"
                                data-field="breakAfterSeconds"
                                (change)="onExerciseChange(si, ei)"
                                class="w-12 border-0 bg-transparent p-0 text-center text-xs text-gray-500 focus:ring-0 dark:text-gray-400"
                              />
                              <span>{{ t('common.unitSeconds') }}</span>
                            </div>
                          </div>
                        }
                      </div>
                    }
                  </div>

                  <!-- Pending exercise groups -->
                  @for (group of section.pendingGroups; track $index; let gi = $index) {
                    <div
                      class="mt-2 rounded-md border-2 border-dashed border-amber-300 bg-amber-50 p-3 dark:border-amber-600 dark:bg-amber-950/30"
                      data-testid="pending-group"
                    >
                      <div class="mb-2 text-sm font-semibold text-amber-800 dark:text-amber-300">
                        {{
                          group.groupName().value() ||
                            t('common.unnamedGroup', { id: group.groupId().value() })
                        }}
                      </div>
                      <div class="mb-1 text-xs text-gray-600 dark:text-gray-400">
                        {{ t('user.workoutStart.pickExercise') }}
                      </div>
                      <select
                        hlmInput
                        class="w-full"
                        (change)="onGroupExercisePicked(si, gi, $event)"
                      >
                        <option value="">{{ t('common.select') }}</option>
                        @for (member of group.members().value(); track member.id) {
                          <option [value]="member.id">{{ member.name }}</option>
                        }
                      </select>
                    </div>
                  }

                  <!-- Add Exercise button -->
                  <button
                    type="button"
                    (click)="openAddExerciseDialog(si)"
                    class="mt-2 text-sm text-blue-500/70 hover:text-blue-600 dark:text-blue-400/70 dark:hover:text-blue-300"
                  >
                    {{ t('user.workouts.addExercise') }}
                  </button>
                </div>
              }
            </div>

            <!-- Add Section button -->
            <button
              type="button"
              (click)="addSection()"
              class="w-full rounded-md border border-dashed border-gray-300 px-3 py-1.5 text-xs text-gray-400 hover:border-gray-400 hover:text-gray-500 dark:border-gray-600 dark:text-gray-500 dark:hover:border-gray-500 dark:hover:text-gray-400"
            >
              {{ t('user.workouts.addSection') }}
            </button>

            <!-- Actions -->
            <div class="flex gap-2">
              <button
                type="submit"
                [disabled]="startMutation.isPending() || hasPendingGroups()"
                class="rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
              >
                @if (startMutation.isPending()) {
                  {{ t('common.starting') }}
                } @else {
                  {{ t('user.workoutStart.startWorkout') }}
                }
              </button>
              <a
                routerLink="/user/workouts"
                class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
              >
                {{ t('common.cancel') }}
              </a>
            </div>
          </form>
        }

        <app-add-exercise-dialog
          [open]="addDialogOpen()"
          [sectionId]="addDialogSectionId()"
          [logId]="currentLogId() ?? 0"
          [exerciseCount]="addDialogExerciseCount()"
          [preselectedExerciseId]="addDialogPreselectedExerciseId()"
          (exerciseAdded)="onExerciseAdded($event)"
          (cancelled)="addDialogOpen.set(false)"
        />
      </app-page-layout>
    </ng-container>
  `,
  styles: `
    .cdk-drag-preview {
      box-sizing: border-box;
      border-radius: 8px;
      border: 1px solid #d1d5db;
      background: white;
      padding: 12px 16px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      font-size: 13px;
      font-weight: 600;
    }

    .cdk-drag-placeholder {
      background: #e0f2fe;
      border: 2px dashed #7dd3fc;
      border-radius: 8px;
      min-height: 3rem;
    }

    .cdk-drag-animating {
      transition: transform 200ms cubic-bezier(0, 0, 0.2, 1);
    }

    .cdk-drop-list-dragging .cdk-drag:not(.cdk-drag-placeholder) {
      transition: transform 200ms cubic-bezier(0, 0, 0.2, 1);
    }

    :host-context(.dark) .cdk-drag-preview {
      background: #1f2937;
      border-color: #4b5563;
      color: #f3f4f6;
    }

    :host-context(.dark) .cdk-drag-placeholder {
      background: rgba(14, 116, 144, 0.2);
      border-color: #0e7490;
    }
  `,
})
export class WorkoutStart {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = inject(QueryClient);
  private params = toSignal(this.route.paramMap);

  readonly store = inject(WorkoutStartStore);
  readonly SECTION_TYPE_MAIN = WorkoutSectionTypeMain;
  readonly SECTION_TYPE_SUPPLEMENTARY = WorkoutSectionTypeSupplementary;

  id = computed(() => Number(this.params()?.get('id')));
  private initialized = false;
  currentLogId = signal<number | null>(null);
  private creating = signal(false);

  isPending = computed(
    () => this.workoutQuery.isPending() || this.planningLogQuery.isPending() || this.creating(),
  );

  hasPendingGroups = computed(() => this.model().sections.some((s) => s.pendingGroups.length > 0));

  // Add exercise dialog state
  addDialogOpen = signal(false);
  addDialogSectionIndex = signal(0);
  addDialogSectionId = signal(0);
  addDialogExerciseCount = signal(0);
  addDialogPreselectedExerciseId = signal<number | null>(null);
  addDialogPendingGroupIndex = signal<number | null>(null);

  model = signal<StartModel>({ name: '', notes: '', sections: [] });
  startForm = form(this.model);

  workoutQuery = injectQuery(() => ({
    queryKey: workoutKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchWorkout(this.id()),
    enabled: !!this.id(),
  }));

  planningLogQuery = injectQuery(() => ({
    queryKey: workoutLogKeys.list({ workoutId: this.id(), status: 'planning' }),
    queryFn: () => this.userApi.fetchWorkoutLogs({ workoutId: this.id(), status: 'planning' }),
    enabled: !!this.id(),
  }));

  startMutation = injectMutation(() => ({
    mutationFn: () => this.startWorkout(),
  }));

  constructor() {
    // Fetch-or-create effect
    effect(() => {
      const workout = this.workoutQuery.data();
      const planningLogs = this.planningLogQuery.data();
      if (!workout || planningLogs === undefined || this.initialized) return;
      this.initialized = true;

      if (planningLogs.length > 0) {
        this.populateFromLog(planningLogs[0]);
      } else {
        this.createPlanningLog(workout);
      }
    });
  }

  // CDK Drag & Drop
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  onSectionDrop(event: CdkDragDrop<any>) {
    if (event.previousIndex === event.currentIndex) return;
    this.model.update((m) => {
      const sections = [...m.sections];
      const [moved] = sections.splice(event.previousIndex, 1);
      sections.splice(event.currentIndex, 0, moved);
      return { ...m, sections };
    });
    this.persistSectionPositions();
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  onExerciseDrop(event: CdkDragDrop<any>, sectionIndex: number) {
    if (event.previousIndex === event.currentIndex) return;
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, i) =>
        i !== sectionIndex
          ? s
          : {
              ...s,
              exercises: (() => {
                const exercises = [...s.exercises];
                const [moved] = exercises.splice(event.previousIndex, 1);
                exercises.splice(event.currentIndex, 0, moved);
                return exercises;
              })(),
            },
      ),
    }));
    this.persistExercisePositions(sectionIndex);
  }

  async removeSection(index: number) {
    const m = this.model();
    const sectionId = m.sections[index].id;
    this.model.update((m) => ({
      ...m,
      sections: m.sections.filter((_, i) => i !== index),
    }));
    if (sectionId) {
      try {
        await this.userApi.deleteWorkoutLogSection(sectionId);
      } catch (err) {
        console.error('Failed to delete section:', err);
      }
    }
  }

  async removeExercise(sectionIndex: number, exerciseIndex: number) {
    const m = this.model();
    const exerciseId = m.sections[sectionIndex].exercises[exerciseIndex].id;
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, i) =>
        i !== sectionIndex
          ? s
          : {
              ...s,
              exercises: s.exercises.filter((_, j) => j !== exerciseIndex),
            },
      ),
    }));
    if (exerciseId) {
      try {
        await this.userApi.deleteWorkoutLogExercise(exerciseId);
      } catch (err) {
        console.error('Failed to delete exercise:', err);
      }
    }
  }

  async addSection() {
    const logId = this.currentLogId();
    if (!logId) return;
    try {
      const section = await this.userApi.createWorkoutLogSection({
        workoutLogId: logId,
        type: WorkoutSectionTypeMain,
        position: this.model().sections.length,
      });
      this.model.update((m) => ({
        ...m,
        sections: [
          ...m.sections,
          {
            id: section.id,
            type: WorkoutSectionTypeMain,
            label: '',
            exercises: [],
            pendingGroups: [],
          },
        ],
      }));
    } catch (err) {
      console.error('Failed to add section:', err);
    }
  }

  async onGroupExercisePicked(sectionIndex: number, groupIndex: number, event: Event) {
    const exerciseId = Number((event.target as HTMLSelectElement).value);
    if (!exerciseId || isNaN(exerciseId)) return;

    const m = this.model();
    const section = m.sections[sectionIndex];
    const sectionId = section.id;
    if (!sectionId) return;

    // Open the add exercise dialog with this exercise pre-selected
    this.addDialogSectionIndex.set(sectionIndex);
    this.addDialogSectionId.set(sectionId);
    this.addDialogExerciseCount.set(section.exercises.length);
    this.addDialogPreselectedExerciseId.set(exerciseId);
    this.addDialogPendingGroupIndex.set(groupIndex);
    this.addDialogOpen.set(true);
  }

  openAddExerciseDialog(sectionIndex: number) {
    const m = this.model();
    const sectionId = m.sections[sectionIndex].id;
    if (!sectionId) return;
    this.addDialogSectionIndex.set(sectionIndex);
    this.addDialogSectionId.set(sectionId);
    this.addDialogExerciseCount.set(m.sections[sectionIndex].exercises.length);
    this.addDialogPreselectedExerciseId.set(null);
    this.addDialogPendingGroupIndex.set(null);
    this.addDialogOpen.set(true);
  }

  onExerciseAdded(event: {
    exerciseLogId: number;
    exerciseName: string;
    scheme: ExerciseScheme;
    exercise: {
      id: number;
      sourceExerciseSchemeId: number;
      sets: {
        id: number;
        setNumber: number;
        targetReps?: number;
        targetWeight?: number;
        targetDuration?: number;
        targetDistance?: number;
        targetTime?: number;
        breakAfterSeconds?: number;
      }[];
    };
  }) {
    const si = this.addDialogSectionIndex();

    const newExercise: StartExerciseModel = {
      id: event.exercise.id,
      sourceExerciseSchemeId: event.exercise.sourceExerciseSchemeId,
      breakAfterSeconds: null,
      sets: (event.exercise.sets ?? []).map((set) => ({
        id: set.id,
        targetReps: set.targetReps ?? null,
        targetWeight: set.targetWeight ?? null,
        targetDuration: set.targetDuration ?? null,
        targetDistance: set.targetDistance ?? null,
        targetTime: set.targetTime ?? null,
        restAfterSeconds: set.breakAfterSeconds ?? null,
      })),
    };

    const pendingGroupIndex = this.addDialogPendingGroupIndex();
    this.model.update((m) => ({
      ...m,
      sections: m.sections.map((s, i) =>
        i !== si
          ? s
          : {
              ...s,
              exercises: [...s.exercises, newExercise],
              pendingGroups:
                pendingGroupIndex != null
                  ? s.pendingGroups.filter((_, j) => j !== pendingGroupIndex)
                  : s.pendingGroups,
            },
      ),
    }));

    // Update store display
    const numSets = event.scheme.sets ?? 0;
    const sets: SetPreview[] = [];
    for (let i = 1; i <= numSets; i++) {
      sets.push({
        setNumber: i,
        targetReps: event.scheme.reps,
        targetWeight: event.scheme.weight,
        targetDuration: event.scheme.duration,
        targetDistance: event.scheme.distance,
        targetTime: event.scheme.targetTime,
        restAfterSeconds: i < numSets ? (event.scheme.restBetweenSets ?? null) : null,
      });
    }
    this.store.addExerciseDisplay(event.exercise.id, event.exerciseName, event.scheme, sets);

    this.addDialogOpen.set(false);
  }

  private async persistSectionPositions() {
    const m = this.model();
    const updates = [];
    for (let i = 0; i < m.sections.length; i++) {
      const sectionId = m.sections[i].id;
      if (sectionId) {
        updates.push(this.userApi.updateWorkoutLogSection(sectionId, { position: i }));
      }
    }
    try {
      await Promise.all(updates);
    } catch (err) {
      console.error('Failed to persist section positions:', err);
    }
  }

  private async persistExercisePositions(sectionIndex: number) {
    const m = this.model();
    const exercises = m.sections[sectionIndex].exercises;
    const updates = [];
    for (let i = 0; i < exercises.length; i++) {
      const exerciseId = exercises[i].id;
      if (exerciseId) {
        updates.push(this.userApi.updateWorkoutLogExercise(exerciseId, { position: i }));
      }
    }
    try {
      await Promise.all(updates);
    } catch (err) {
      console.error('Failed to persist exercise positions:', err);
    }
  }

  formatBreak = formatBreak;

  onSubmit() {
    this.startMutation.mutate();
  }

  async onLogChange() {
    const logId = this.currentLogId();
    if (!logId) return;
    const val = this.model();
    try {
      await this.userApi.updateWorkoutLog(logId, {
        name: val.name,
        notes: val.notes || undefined,
      });
    } catch (err) {
      console.error('Failed to save log changes:', err);
    }
  }

  async onSectionChange(si: number) {
    const m = this.model();
    const section = m.sections[si];
    const sectionId = section.id;
    if (!sectionId) return;
    try {
      await this.userApi.updateWorkoutLogSection(sectionId, {
        type: section.type,
        label: section.label || undefined,
      });
    } catch (err) {
      console.error('Failed to save section changes:', err);
    }
  }

  async onExerciseChange(si: number, ei: number) {
    const m = this.model();
    const exercise = m.sections[si].exercises[ei];
    const exerciseId = exercise.id;
    if (!exerciseId) return;
    try {
      await this.userApi.updateWorkoutLogExercise(exerciseId, {
        breakAfterSeconds: exercise.breakAfterSeconds ?? undefined,
      });
    } catch (err) {
      console.error('Failed to save exercise changes:', err);
    }
  }

  async onSetChange(si: number, ei: number, setIdx: number) {
    const m = this.model();
    const set = m.sections[si].exercises[ei].sets[setIdx];
    const setId = set.id;
    if (!setId) return;

    try {
      await this.userApi.updateWorkoutLogExerciseSet(setId, {
        targetReps: set.targetReps ?? undefined,
        targetWeight: set.targetWeight ?? undefined,
        targetDuration: set.targetDuration ?? undefined,
        targetDistance: set.targetDistance ?? undefined,
        targetTime: set.targetTime ?? undefined,
        breakAfterSeconds: set.restAfterSeconds ?? undefined,
      });
    } catch (err) {
      console.error('Failed to save set changes:', err);
    }
  }

  private populateFromLog(log: WorkoutLog) {
    this.currentLogId.set(log.id);

    this.model.set({
      name: log.name,
      notes: log.notes ?? '',
      sections: (log.sections ?? []).map((section) => ({
        id: section.id,
        type: section.type,
        label: section.label ?? '',
        exercises: (section.exercises ?? []).map((ex) => ({
          id: ex.id,
          sourceExerciseSchemeId: ex.sourceExerciseSchemeId,
          breakAfterSeconds: ex.breakAfterSeconds ?? null,
          sets: (ex.sets ?? []).map((set) => ({
            id: set.id,
            targetReps: set.targetReps ?? null,
            targetWeight: set.targetWeight ?? null,
            targetDuration: set.targetDuration ?? null,
            targetDistance: set.targetDistance ?? null,
            targetTime: set.targetTime ?? null,
            restAfterSeconds: set.breakAfterSeconds ?? null,
          })),
        })),
        pendingGroups: [],
      })),
    });

    this.store.loadExerciseDisplayFromLog(log.sections ?? []);
  }

  private async createPlanningLog(workout: Workout) {
    this.creating.set(true);

    // Pre-fetch user-specific schemes for all exercise items
    const exerciseItems = (workout.sections ?? []).flatMap((s) =>
      (s.items ?? []).filter((i) => i.type === 'exercise'),
    );
    const schemeResults = await Promise.all(
      exerciseItems.map((item) =>
        this.userApi.fetchExerciseSchemes({ workoutSectionItemId: item.id }),
      ),
    );
    const userSchemeByItemId = new Map(
      schemeResults
        .flat()
        .filter((s) => s.workoutSectionItemId != null)
        .map((s) => [s.workoutSectionItemId!, s]),
    );

    const log = await this.userApi.createWorkoutLog({
      name: workout.name,
      notes: workout.notes || undefined,
      workoutId: workout.id,
    });

    for (let si = 0; si < (workout.sections ?? []).length; si++) {
      const templateSection = workout.sections[si];

      const section = await this.userApi.createWorkoutLogSection({
        workoutLogId: log.id,
        type: templateSection.type,
        label: templateSection.label || undefined,
        position: si,
        restBetweenExercises: templateSection.restBetweenExercises ?? undefined,
      });

      for (let ei = 0; ei < (templateSection.items ?? []).length; ei++) {
        const templateItem = templateSection.items[ei];
        // Only create log exercises for exercise-type items;
        // exercise_group items are resolved on the start page
        if (templateItem.type === 'exercise') {
          // Use user-specific scheme (looked up by item ID), falling back to item's embedded scheme
          const myScheme = userSchemeByItemId.get(templateItem.id);
          const schemeId = myScheme?.id ?? templateItem.exerciseSchemeId;
          if (schemeId) {
            await this.userApi.createWorkoutLogExercise({
              workoutLogSectionId: section.id,
              sourceExerciseSchemeId: schemeId,
              position: ei,
            });
          }
        }
      }
    }

    // Collect pending exercise groups and resolve their members
    const pendingGroupsBySection = new Map<number, PendingGroupModel[]>();
    const userExercises = await this.userApi.fetchUserExercises();
    const exerciseNameMap = new Map(userExercises.map((e) => [e.id, e.name]));

    for (let si = 0; si < (workout.sections ?? []).length; si++) {
      const templateSection = workout.sections[si];
      const groups: PendingGroupModel[] = [];

      for (let ei = 0; ei < (templateSection.items ?? []).length; ei++) {
        const templateItem = templateSection.items[ei];
        if (templateItem.type === 'exercise_group' && templateItem.exerciseGroupId) {
          const [group, members] = await Promise.all([
            this.compendiumApi.fetchExerciseGroup(templateItem.exerciseGroupId),
            this.compendiumApi.fetchExerciseGroupMembers({
              groupId: templateItem.exerciseGroupId,
            }),
          ]);
          groups.push({
            groupId: templateItem.exerciseGroupId,
            groupName: group.name ?? '',
            members: members.map((m) => ({
              id: m.exerciseId,
              name: exerciseNameMap.get(m.exerciseId) ?? `Exercise #${m.exerciseId}`,
            })),
            position: ei,
          });
        }
      }
      if (groups.length > 0) {
        pendingGroupsBySection.set(si, groups);
      }
    }

    // Fetch the full log with nested structure
    const fullLog = await this.userApi.fetchWorkoutLog(log.id);
    this.creating.set(false);
    this.populateFromLog(fullLog);

    // Inject pending groups into the model sections
    if (pendingGroupsBySection.size > 0) {
      this.model.update((m) => ({
        ...m,
        sections: m.sections.map((s, i) => ({
          ...s,
          pendingGroups: pendingGroupsBySection.get(i) ?? [],
        })),
      }));
    }

    // Invalidate planning log query so it reflects the new log
    this.queryClient.invalidateQueries({
      queryKey: workoutLogKeys.list({ workoutId: workout.id, status: 'planning' }),
    });
  }

  private async startWorkout() {
    await this.userApi.startWorkoutLog(this.currentLogId()!);
    this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.all() });
    this.router.navigate(['/user/workout-logs', this.currentLogId()]);
  }
}
