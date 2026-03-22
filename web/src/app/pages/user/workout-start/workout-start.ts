import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { ReactiveFormsModule, FormGroup, FormControl, FormArray } from '@angular/forms';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { CdkDragDrop, CdkDrag, CdkDropList, CdkDragHandle } from '@angular/cdk/drag-drop';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { formatBreak } from '$core/format-utils';
import { workoutKeys, workoutLogKeys } from '$core/query-keys';
import {
  Workout,
  WorkoutSectionTypeMain,
  WorkoutSectionTypeSupplementary,
  WorkoutLog,
  UserExerciseScheme,
} from '$generated/user-models';
import { PageLayout } from '../../../layout/page-layout';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmInput } from '@spartan-ng/helm/input';
import { HlmTextarea } from '@spartan-ng/helm/textarea';
import { WorkoutStartStore, SetPreview } from './workout-start.store';
import { AddExerciseDialog } from './add-exercise-dialog';

type SetFormGroup = FormGroup<{
  id: FormControl<number | null>;
  targetReps: FormControl<number | null>;
  targetWeight: FormControl<number | null>;
  targetDuration: FormControl<number | null>;
  targetDistance: FormControl<number | null>;
  targetTime: FormControl<number | null>;
  restAfterSeconds: FormControl<number | null>;
}>;

type ExerciseFormGroup = FormGroup<{
  id: FormControl<number | null>;
  sourceExerciseSchemeId: FormControl<number>;
  breakAfterSeconds: FormControl<number | null>;
  sets: FormArray<SetFormGroup>;
}>;

type SectionFormGroup = FormGroup<{
  id: FormControl<number | null>;
  type: FormControl<string>;
  label: FormControl<string>;
  exercises: FormArray<ExerciseFormGroup>;
}>;

@Component({
  selector: 'app-workout-start',
  imports: [
    PageLayout,
    ReactiveFormsModule,
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
          <form [formGroup]="form" (ngSubmit)="onSubmit()" class="space-y-6">
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
                formControlName="name"
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
                formControlName="notes"
                (change)="onLogChange()"
                rows="2"
                class="mt-1"
              ></textarea>
            </div>

            <!-- Sections -->
            <div
              formArrayName="sections"
              cdkDropList
              [cdkDropListData]="sectionsArray.controls"
              (cdkDropListDropped)="onSectionDrop($event)"
              class="space-y-4"
            >
              @for (section of sectionsArray.controls; track $index; let si = $index) {
                <div
                  [formGroupName]="si"
                  cdkDrag
                  class="rounded-lg border border-gray-200 p-4 dark:border-gray-700"
                >
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
                        formControlName="type"
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
                        formControlName="label"
                        (change)="onSectionChange(si)"
                        class="mt-1"
                      />
                    </label>
                  </div>

                  <!-- Exercise cards -->
                  <div
                    formArrayName="exercises"
                    cdkDropList
                    [cdkDropListData]="getExercisesArray(si).controls"
                    (cdkDropListDropped)="onExerciseDrop($event, si)"
                  >
                    @for (
                      ex of getExercisesArray(si).controls;
                      track $index;
                      let ei = $index;
                      let lastEx = $last
                    ) {
                      <div cdkDrag>
                        @let info = store.exerciseDisplay()[ex.get('id')!.value!];
                        <div
                          [formGroupName]="ei"
                          class="rounded-md border border-gray-200 dark:border-gray-600"
                        >
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
                          @if (getSetsArray(si, ei).length) {
                            @let setsArr = getSetsArray(si, ei);
                            <div class="px-3 py-2" [formArrayName]="'sets'">
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
                                setCtrl of setsArr.controls;
                                track $index;
                                let setIdx = $index;
                                let lastSet = $last
                              ) {
                                <!-- Set row -->
                                <div
                                  [formGroupName]="setIdx"
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
                                        formControlName="targetReps"
                                        (change)="onSetChange(si, ei, setIdx)"
                                        class="mt-1"
                                      />
                                    </div>
                                    <div>
                                      <input
                                        hlmInput
                                        type="number"
                                        formControlName="targetWeight"
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
                                        formControlName="targetDuration"
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
                                        formControlName="targetDistance"
                                        (change)="onSetChange(si, ei, setIdx)"
                                        class="mt-1"
                                        step="0.1"
                                      />
                                    </div>
                                  }
                                </div>

                                <!-- Rest between sets: line with centered badge -->
                                @if (!lastSet && setCtrl.controls.restAfterSeconds.value !== null) {
                                  <div
                                    [formGroupName]="setIdx"
                                    class="relative flex items-center justify-center py-0.5"
                                  >
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
                                        formControlName="restAfterSeconds"
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
                          <div
                            [formGroupName]="ei"
                            class="relative flex items-center justify-center py-3"
                          >
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
                                formControlName="breakAfterSeconds"
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
                [disabled]="startMutation.isPending()"
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
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = injectQueryClient();
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

  // Add exercise dialog state
  addDialogOpen = signal(false);
  addDialogSectionIndex = signal(0);
  addDialogSectionId = signal(0);
  addDialogExerciseCount = signal(0);

  form = new FormGroup({
    name: new FormControl('', { nonNullable: true }),
    notes: new FormControl('', { nonNullable: true }),
    sections: new FormArray<SectionFormGroup>([]),
  });

  get sectionsArray() {
    return this.form.controls.sections;
  }

  getExercisesArray(sectionIndex: number) {
    return this.sectionsArray.at(sectionIndex).controls.exercises;
  }

  getSetsArray(sectionIndex: number, exerciseIndex: number) {
    return this.getExercisesArray(sectionIndex).at(exerciseIndex).controls.sets;
  }

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
  onSectionDrop(event: CdkDragDrop<SectionFormGroup[]>) {
    if (event.previousIndex === event.currentIndex) return;
    this.moveFormArrayItem(this.sectionsArray, event.previousIndex, event.currentIndex);
    this.persistSectionPositions();
  }

  onExerciseDrop(event: CdkDragDrop<ExerciseFormGroup[]>, sectionIndex: number) {
    if (event.previousIndex === event.currentIndex) return;
    this.moveFormArrayItem(
      this.getExercisesArray(sectionIndex),
      event.previousIndex,
      event.currentIndex,
    );
    this.persistExercisePositions(sectionIndex);
  }

  private moveFormArrayItem(formArray: FormArray, from: number, to: number) {
    const control = formArray.at(from);
    formArray.removeAt(from);
    formArray.insert(to, control);
  }

  async removeSection(index: number) {
    const sectionId = this.sectionsArray.at(index).controls.id.value;
    this.sectionsArray.removeAt(index);
    if (sectionId) {
      try {
        await this.userApi.deleteWorkoutLogSection(sectionId);
      } catch (err) {
        console.error('Failed to delete section:', err);
      }
    }
  }

  async removeExercise(sectionIndex: number, exerciseIndex: number) {
    const exerciseId = this.getExercisesArray(sectionIndex).at(exerciseIndex).controls.id.value;
    this.getExercisesArray(sectionIndex).removeAt(exerciseIndex);
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
        position: this.sectionsArray.length,
      });
      this.sectionsArray.push(this.createSectionGroup(section.id));
    } catch (err) {
      console.error('Failed to add section:', err);
    }
  }

  openAddExerciseDialog(sectionIndex: number) {
    const sectionId = this.sectionsArray.at(sectionIndex).controls.id.value;
    if (!sectionId) return;
    this.addDialogSectionIndex.set(sectionIndex);
    this.addDialogSectionId.set(sectionId);
    this.addDialogExerciseCount.set(this.getExercisesArray(sectionIndex).length);
    this.addDialogOpen.set(true);
  }

  onExerciseAdded(event: {
    exerciseLogId: number;
    exerciseName: string;
    scheme: UserExerciseScheme;
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

    const exGroup = this.createExerciseGroup(
      event.exercise.sourceExerciseSchemeId,
      null,
      event.exercise.id,
    );

    for (const set of event.exercise.sets ?? []) {
      exGroup.controls.sets.push(
        this.createSetGroup(
          {
            setNumber: set.setNumber,
            targetReps: set.targetReps,
            targetWeight: set.targetWeight,
            targetDuration: set.targetDuration,
            targetDistance: set.targetDistance,
            targetTime: set.targetTime,
            restAfterSeconds: set.breakAfterSeconds ?? null,
          },
          set.id,
        ),
      );
    }

    this.getExercisesArray(si).push(exGroup);

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
    const updates = [];
    for (let i = 0; i < this.sectionsArray.length; i++) {
      const sectionId = this.sectionsArray.at(i).controls.id.value;
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
    const exercises = this.getExercisesArray(sectionIndex);
    const updates = [];
    for (let i = 0; i < exercises.length; i++) {
      const exerciseId = exercises.at(i).controls.id.value;
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
    const val = this.form.getRawValue();
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
    const section = this.sectionsArray.at(si);
    const sectionId = section.controls.id.value;
    if (!sectionId) return;
    const val = section.getRawValue();
    try {
      await this.userApi.updateWorkoutLogSection(sectionId, {
        type: val.type,
        label: val.label || undefined,
      });
    } catch (err) {
      console.error('Failed to save section changes:', err);
    }
  }

  async onExerciseChange(si: number, ei: number) {
    const exercise = this.getExercisesArray(si).at(ei);
    const exerciseId = exercise.controls.id.value;
    if (!exerciseId) return;
    const val = exercise.getRawValue();
    try {
      await this.userApi.updateWorkoutLogExercise(exerciseId, {
        breakAfterSeconds: val.breakAfterSeconds ?? undefined,
      });
    } catch (err) {
      console.error('Failed to save exercise changes:', err);
    }
  }

  async onSetChange(si: number, ei: number, setIdx: number) {
    const setGroup = this.getSetsArray(si, ei).at(setIdx);
    const setId = setGroup.controls.id.value;
    if (!setId) return;

    const val = setGroup.getRawValue();
    try {
      await this.userApi.updateWorkoutLogExerciseSet(setId, {
        targetReps: val.targetReps ?? undefined,
        targetWeight: val.targetWeight ?? undefined,
        targetDuration: val.targetDuration ?? undefined,
        targetDistance: val.targetDistance ?? undefined,
        targetTime: val.targetTime ?? undefined,
        breakAfterSeconds: val.restAfterSeconds ?? undefined,
      });
    } catch (err) {
      console.error('Failed to save set changes:', err);
    }
  }

  private populateFromLog(log: WorkoutLog) {
    this.currentLogId.set(log.id);
    this.form.patchValue({
      name: log.name,
      notes: log.notes ?? '',
    });

    this.sectionsArray.clear();

    for (const section of log.sections ?? []) {
      const sectionGroup = this.createSectionGroup(section.id);
      sectionGroup.patchValue({
        type: section.type,
        label: section.label ?? '',
      });

      for (const ex of section.exercises ?? []) {
        const exGroup = this.createExerciseGroup(
          ex.sourceExerciseSchemeId,
          ex.breakAfterSeconds ?? null,
          ex.id,
        );

        for (const set of ex.sets ?? []) {
          exGroup.controls.sets.push(
            this.createSetGroup(
              {
                setNumber: set.setNumber,
                targetReps: set.targetReps,
                targetWeight: set.targetWeight,
                targetDuration: set.targetDuration,
                targetDistance: set.targetDistance,
                targetTime: set.targetTime,
                restAfterSeconds: set.breakAfterSeconds ?? null,
              },
              set.id,
            ),
          );
        }

        sectionGroup.controls.exercises.push(exGroup);
      }

      this.sectionsArray.push(sectionGroup);
    }

    this.store.loadExerciseDisplayFromLog(log.sections ?? []);
  }

  private async createPlanningLog(workout: Workout) {
    this.creating.set(true);

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

      for (let ei = 0; ei < (templateSection.exercises ?? []).length; ei++) {
        const templateEx = templateSection.exercises[ei];
        await this.userApi.createWorkoutLogExercise({
          workoutLogSectionId: section.id,
          sourceExerciseSchemeId: templateEx.userExerciseSchemeId,
          position: ei,
        });
      }
    }

    // Fetch the full log with nested structure
    const fullLog = await this.userApi.fetchWorkoutLog(log.id);
    this.creating.set(false);
    this.populateFromLog(fullLog);

    // Invalidate planning log query so it reflects the new log
    this.queryClient.invalidateQueries({
      queryKey: workoutLogKeys.list({ workoutId: workout.id, status: 'planning' }),
    });
  }

  private createSectionGroup(id: number | null = null): SectionFormGroup {
    return new FormGroup({
      id: new FormControl<number | null>(id),
      type: new FormControl(WorkoutSectionTypeMain, { nonNullable: true }),
      label: new FormControl('', { nonNullable: true }),
      exercises: new FormArray<ExerciseFormGroup>([]),
    });
  }

  private createExerciseGroup(
    schemeId: number,
    breakAfterSeconds: number | null = null,
    id: number | null = null,
  ): ExerciseFormGroup {
    return new FormGroup({
      id: new FormControl<number | null>(id),
      sourceExerciseSchemeId: new FormControl(schemeId, { nonNullable: true }),
      breakAfterSeconds: new FormControl<number | null>(breakAfterSeconds),
      sets: new FormArray<SetFormGroup>([]),
    });
  }

  private createSetGroup(set: SetPreview, id: number | null = null): SetFormGroup {
    return new FormGroup({
      id: new FormControl<number | null>(id),
      targetReps: new FormControl<number | null>(set.targetReps ?? null),
      targetWeight: new FormControl<number | null>(set.targetWeight ?? null),
      targetDuration: new FormControl<number | null>(set.targetDuration ?? null),
      targetDistance: new FormControl<number | null>(set.targetDistance ?? null),
      targetTime: new FormControl<number | null>(set.targetTime ?? null),
      restAfterSeconds: new FormControl<number | null>(set.restAfterSeconds ?? null),
    });
  }

  private async startWorkout() {
    await this.userApi.startWorkoutLog(this.currentLogId()!);
    this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.all() });
    this.router.navigate(['/user/workout-logs', this.currentLogId()]);
  }
}
