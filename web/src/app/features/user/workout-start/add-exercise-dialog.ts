import { Component, inject, input, output, signal, computed, effect } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { injectQuery, injectQueryClient } from '@tanstack/angular-query-experimental';
import { injectQueries } from '@tanstack/angular-query-experimental/inject-queries-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { userExerciseKeys, exerciseKeys, exerciseSchemeKeys } from '$core/query-keys';
import { UserExerciseScheme } from '$generated/user-models';
import { HlmComboboxImports } from '@spartan-ng/helm/combobox';
import { HlmDialogImports } from '@spartan-ng/helm/dialog';
import { HlmButton } from '@spartan-ng/helm/button';

interface EnrichedExercise {
  id: number;
  name: string;
}

@Component({
  selector: 'app-add-exercise-dialog',
  imports: [FormsModule, HlmComboboxImports, HlmDialogImports, HlmButton],
  template: `
    <hlm-dialog [state]="open() ? 'open' : 'closed'" (closed)="onCancel()">
      <ng-template hlmDialogPortal>
        <hlm-dialog-content [showCloseButton]="false">
          <hlm-dialog-header>
            <h3 hlmDialogTitle>Add Exercise</h3>
          </hlm-dialog-header>

          <!-- Step 1: Exercise search via spartan combobox -->
          @if (!selectedExerciseId()) {
            <hlm-combobox
              [value]="comboboxValue()"
              (valueChange)="onComboboxSelect($event)"
              [itemToString]="exerciseToString"
            >
              <hlm-combobox-input placeholder="Search exercises..." [showClear]="false" />
              <hlm-combobox-content *hlmComboboxPortal>
                <hlm-combobox-empty>No exercises found.</hlm-combobox-empty>
                <div hlmComboboxList>
                  @for (ex of enrichedUserExercises(); track ex.id) {
                    <hlm-combobox-item [value]="ex">
                      {{ ex.name }}
                    </hlm-combobox-item>
                  }
                </div>
              </hlm-combobox-content>
            </hlm-combobox>
          }

          <!-- Step 2: Scheme selection / configuration -->
          @if (selectedExerciseId()) {
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{
                  selectedExerciseName()
                }}</span>
                <button
                  type="button"
                  (click)="clearSelection()"
                  class="text-xs text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
                >
                  Change
                </button>
              </div>

              <!-- Existing schemes -->
              @if (schemesQuery.data(); as schemes) {
                @if (schemes.length > 0 && !editingScheme()) {
                  <div class="space-y-1">
                    <span class="text-xs font-medium text-gray-500 dark:text-gray-400"
                      >Use existing scheme</span
                    >
                    @for (s of schemes; track s.id) {
                      <button
                        type="button"
                        (click)="pickScheme(s)"
                        class="block w-full rounded-md border px-3 py-2 text-left text-sm hover:bg-gray-50 dark:hover:bg-gray-700"
                        [class]="
                          pickedSchemeId() === s.id
                            ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 dark:border-blue-400'
                            : 'border-gray-200 dark:border-gray-600'
                        "
                      >
                        {{ formatScheme(s) }}
                      </button>
                    }
                    <button
                      type="button"
                      (click)="startCustomScheme()"
                      class="mt-1 text-xs text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
                    >
                      + Custom scheme
                    </button>
                  </div>
                }
              }

              <!-- Scheme editor (shown when no schemes or custom selected) -->
              @if (editingScheme() || schemesQuery.data()?.length === 0) {
                <div class="space-y-3">
                  <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                    Measurement Type
                    <select
                      [(ngModel)]="schemeMeasurementType"
                      class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-1.5 text-sm shadow-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                    >
                      <option value="REP_BASED">Rep Based</option>
                      <option value="TIME_BASED">Time Based</option>
                      <option value="DISTANCE_BASED">Distance Based</option>
                    </select>
                  </label>

                  @if (schemeMeasurementType === 'REP_BASED') {
                    <div class="grid grid-cols-2 gap-2">
                      <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                        Sets
                        <input
                          type="number"
                          [(ngModel)]="schemeSets"
                          class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                        />
                      </label>
                      <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                        Reps
                        <input
                          type="number"
                          [(ngModel)]="schemeReps"
                          class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                        />
                      </label>
                      <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                        Weight (kg)
                        <input
                          type="number"
                          [(ngModel)]="schemeWeight"
                          step="0.5"
                          class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                        />
                      </label>
                      <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                        Rest (s)
                        <input
                          type="number"
                          [(ngModel)]="schemeRest"
                          class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                        />
                      </label>
                    </div>
                  }

                  @if (schemeMeasurementType === 'TIME_BASED') {
                    <div class="grid grid-cols-2 gap-2">
                      <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                        Sets
                        <input
                          type="number"
                          [(ngModel)]="schemeSets"
                          class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                        />
                      </label>
                      <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                        Duration (s)
                        <input
                          type="number"
                          [(ngModel)]="schemeDuration"
                          class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                        />
                      </label>
                    </div>
                  }

                  @if (schemeMeasurementType === 'DISTANCE_BASED') {
                    <div class="grid grid-cols-2 gap-2">
                      <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                        Sets
                        <input
                          type="number"
                          [(ngModel)]="schemeSets"
                          class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                        />
                      </label>
                      <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
                        Distance (m)
                        <input
                          type="number"
                          [(ngModel)]="schemeDistance"
                          step="0.1"
                          class="mt-1 block w-full rounded-md border border-gray-300 px-2 py-1.5 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                        />
                      </label>
                    </div>
                  }
                </div>
              }
            </div>
          }

          <!-- Actions -->
          <hlm-dialog-footer>
            <button hlmBtn variant="outline" hlmDialogClose [disabled]="isAdding()">Cancel</button>
            <button hlmBtn (click)="onAdd()" [disabled]="!canAdd() || isAdding()">
              @if (isAdding()) {
                Adding...
              } @else {
                Add
              }
            </button>
          </hlm-dialog-footer>
        </hlm-dialog-content>
      </ng-template>
    </hlm-dialog>
  `,
})
export class AddExerciseDialog {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private queryClient = injectQueryClient();

  open = input(false);
  sectionId = input.required<number>();
  logId = input.required<number>();
  exerciseCount = input(0);

  exerciseAdded = output<{
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
  }>();
  cancelled = output();

  // Combobox state
  comboboxValue = signal<EnrichedExercise | null>(null);
  selectedExerciseId = signal<number | null>(null);
  selectedExerciseName = signal('');
  pickedSchemeId = signal<number | null>(null);
  pickedScheme = signal<UserExerciseScheme | null>(null);
  editingScheme = signal(false);
  isAdding = signal(false);

  // Scheme form fields
  schemeMeasurementType = 'REP_BASED';
  schemeSets: number | null = 3;
  schemeReps: number | null = 10;
  schemeWeight: number | null = null;
  schemeRest: number | null = null;
  schemeDuration: number | null = null;
  schemeDistance: number | null = null;

  // User exercises query
  private userExercisesQuery = injectQuery(() => ({
    queryKey: userExerciseKeys.list(),
    queryFn: () => this.userApi.fetchUserExercises(),
  }));

  private snapshotQueries = injectQueries(() => ({
    queries: (this.userExercisesQuery.data() ?? []).map((ue) => ({
      queryKey: exerciseKeys.version(ue.compendiumExerciseId, ue.compendiumVersion),
      queryFn: () =>
        this.compendiumApi.fetchExerciseVersion(ue.compendiumExerciseId, ue.compendiumVersion),
      staleTime: Infinity,
    })),
  }));

  enrichedUserExercises = computed(() => {
    const userExercises = this.userExercisesQuery.data();
    if (!userExercises) return [];
    const snapshots = this.snapshotQueries();
    return userExercises.map((ue, i) => {
      const exercise = snapshots[i]?.data()?.snapshot;
      return { id: ue.id, name: exercise?.name ?? `Exercise #${ue.id}` };
    });
  });

  exerciseToString = (ex: EnrichedExercise) => ex.name;

  // Schemes for selected exercise
  schemesQuery = injectQuery(() => ({
    queryKey: exerciseSchemeKeys.list({ userExerciseId: this.selectedExerciseId()! }),
    queryFn: () =>
      this.userApi.fetchExerciseSchemes({ userExerciseId: this.selectedExerciseId()! }),
    enabled: this.selectedExerciseId() != null,
  }));

  canAdd = computed(() => {
    if (!this.selectedExerciseId()) return false;
    if (this.pickedSchemeId() || this.editingScheme() || this.schemesQuery.data()?.length === 0)
      return true;
    return false;
  });

  constructor() {
    // Auto-set editing mode when no schemes exist
    effect(() => {
      const schemes = this.schemesQuery.data();
      if (schemes && schemes.length === 0 && this.selectedExerciseId()) {
        this.editingScheme.set(true);
      }
    });
  }

  onComboboxSelect(exercise: EnrichedExercise | null) {
    if (exercise) {
      this.selectedExerciseId.set(exercise.id);
      this.selectedExerciseName.set(exercise.name);
      this.comboboxValue.set(null);
      this.pickedSchemeId.set(null);
      this.pickedScheme.set(null);
      this.editingScheme.set(false);
    }
  }

  clearSelection() {
    this.selectedExerciseId.set(null);
    this.selectedExerciseName.set('');
    this.comboboxValue.set(null);
    this.pickedSchemeId.set(null);
    this.pickedScheme.set(null);
    this.editingScheme.set(false);
  }

  pickScheme(scheme: UserExerciseScheme) {
    this.pickedSchemeId.set(scheme.id);
    this.pickedScheme.set(scheme);
    this.editingScheme.set(false);
  }

  startCustomScheme() {
    this.pickedSchemeId.set(null);
    this.pickedScheme.set(null);
    this.editingScheme.set(true);
  }

  formatScheme(s: UserExerciseScheme): string {
    if (s.measurementType === 'REP_BASED') {
      const parts: string[] = [];
      if (s.sets) parts.push(`${s.sets}x`);
      if (s.reps) parts.push(`${s.reps}`);
      const sr = parts.join('');
      if (s.weight) return `${sr} @ ${s.weight}kg`;
      return sr || 'Rep based';
    }
    if (s.measurementType === 'TIME_BASED') {
      return s.duration ? `${s.duration}s` : 'Time based';
    }
    if (s.measurementType === 'DISTANCE_BASED') {
      return s.distance ? `${s.distance}m` : 'Distance based';
    }
    return s.measurementType;
  }

  async onAdd() {
    this.isAdding.set(true);
    try {
      let scheme: UserExerciseScheme;

      if (this.pickedScheme()) {
        scheme = this.pickedScheme()!;
      } else {
        // Create new scheme
        const data: Record<string, unknown> = {
          userExerciseId: this.selectedExerciseId(),
          measurementType: this.schemeMeasurementType,
        };
        if (this.schemeMeasurementType === 'REP_BASED') {
          if (this.schemeSets != null) data['sets'] = this.schemeSets;
          if (this.schemeReps != null) data['reps'] = this.schemeReps;
          if (this.schemeWeight != null) data['weight'] = this.schemeWeight;
          if (this.schemeRest != null) data['restBetweenSets'] = this.schemeRest;
        } else if (this.schemeMeasurementType === 'TIME_BASED') {
          if (this.schemeSets != null) data['sets'] = this.schemeSets;
          if (this.schemeDuration != null) data['duration'] = this.schemeDuration;
        } else if (this.schemeMeasurementType === 'DISTANCE_BASED') {
          if (this.schemeSets != null) data['sets'] = this.schemeSets;
          if (this.schemeDistance != null) data['distance'] = this.schemeDistance;
        }
        scheme = await this.userApi.createExerciseScheme(data);
        this.queryClient.invalidateQueries({ queryKey: exerciseSchemeKeys.all() });
      }

      // Create log exercise
      const logExercise = await this.userApi.createWorkoutLogExercise({
        workoutLogSectionId: this.sectionId(),
        sourceExerciseSchemeId: scheme.id,
        position: this.exerciseCount(),
      });

      // Fetch the created exercise with its sets
      const section = await this.userApi.fetchWorkoutLogs({ status: 'planning' });
      let createdExercise:
        | (typeof logExercise & {
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
          })
        | undefined;
      for (const log of section) {
        for (const s of log.sections ?? []) {
          for (const ex of s.exercises ?? []) {
            if (ex.id === logExercise.id) {
              createdExercise = ex as typeof createdExercise;
            }
          }
        }
      }

      this.exerciseAdded.emit({
        exerciseLogId: logExercise.id,
        exerciseName: this.selectedExerciseName(),
        scheme,
        exercise: createdExercise ?? {
          id: logExercise.id,
          sourceExerciseSchemeId: scheme.id,
          sets: [],
        },
      });

      this.resetState();
    } catch (err) {
      console.error('Failed to add exercise:', err);
    } finally {
      this.isAdding.set(false);
    }
  }

  onCancel() {
    this.resetState();
    this.cancelled.emit();
  }

  private resetState() {
    this.comboboxValue.set(null);
    this.selectedExerciseId.set(null);
    this.selectedExerciseName.set('');
    this.pickedSchemeId.set(null);
    this.pickedScheme.set(null);
    this.editingScheme.set(false);
    this.schemeMeasurementType = 'REP_BASED';
    this.schemeSets = 3;
    this.schemeReps = 10;
    this.schemeWeight = null;
    this.schemeRest = null;
    this.schemeDuration = null;
    this.schemeDistance = null;
  }
}
