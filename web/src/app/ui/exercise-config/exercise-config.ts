import { Component, inject, computed, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { injectQuery, injectQueryClient } from '@tanstack/angular-query-experimental';
import { injectQueries } from '@tanstack/angular-query-experimental/inject-queries-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { userExerciseKeys, exerciseKeys, exerciseSchemeKeys } from '$core/query-keys';
import { UserExerciseScheme } from '$generated/user-models';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmInput } from '@spartan-ng/helm/input';

export interface ExerciseConfigResult {
  userExerciseId: number;
  exerciseName: string;
  scheme: UserExerciseScheme;
}

@Component({
  selector: 'app-exercise-config',
  imports: [FormsModule, BrnSelectImports, HlmSelectImports, HlmInput],
  template: `
    <div
      class="rounded-md border border-gray-100 bg-gray-50 p-3 dark:border-gray-600 dark:bg-gray-800/50"
    >
      <div class="mb-2 grid grid-cols-1 gap-2 sm:grid-cols-2">
        <div>
          <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">Exercise *</span>
          <brn-select [(ngModel)]="userExerciseId" class="mt-1" hlm placeholder="-- Select --">
            <hlm-select-trigger class="w-full">
              <hlm-select-value />
            </hlm-select-trigger>
            <hlm-select-content>
              @for (ue of enrichedUserExercises(); track ue.id) {
                <hlm-option [value]="ue.id">{{ ue.name }}</hlm-option>
              }
            </hlm-select-content>
          </brn-select>
        </div>
        <div>
          <span class="block text-xs font-medium text-gray-700 dark:text-gray-300"
            >Measurement Type</span
          >
          <brn-select [(ngModel)]="measurementType" class="mt-1" hlm>
            <hlm-select-trigger class="w-full">
              <hlm-select-value />
            </hlm-select-trigger>
            <hlm-select-content>
              <hlm-option value="REP_BASED">Rep Based</hlm-option>
              <hlm-option value="TIME_BASED">Time Based</hlm-option>
              <hlm-option value="DISTANCE_BASED">Distance Based</hlm-option>
            </hlm-select-content>
          </brn-select>
        </div>
      </div>

      <!-- REP_BASED fields -->
      @if (measurementType() === 'REP_BASED') {
        <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            Sets
            <input type="number" [(ngModel)]="sets" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            Reps
            <input type="number" [(ngModel)]="reps" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            Weight (kg)
            <input type="number" [(ngModel)]="weight" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            Rest (s)
            <input type="number" [(ngModel)]="restBetweenSets" hlmInput class="mt-1" />
          </label>
        </div>
      }

      <!-- TIME_BASED fields -->
      @if (measurementType() === 'TIME_BASED') {
        <div class="grid grid-cols-2 gap-2">
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            Duration (s)
            <input type="number" [(ngModel)]="duration" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            Time Per Rep (s)
            <input type="number" [(ngModel)]="timePerRep" hlmInput class="mt-1" />
          </label>
        </div>
      }

      <!-- DISTANCE_BASED fields -->
      @if (measurementType() === 'DISTANCE_BASED') {
        <div class="grid grid-cols-2 gap-2">
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            Distance (m)
            <input type="number" [(ngModel)]="distance" hlmInput class="mt-1" />
          </label>
          <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
            Target Time (s)
            <input type="number" [(ngModel)]="targetTime" hlmInput class="mt-1" />
          </label>
        </div>
      }
    </div>
  `,
})
export class ExerciseConfig {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private queryClient = injectQueryClient();

  // Form fields (signal-backed for [(ngModel)])
  userExerciseId = signal<number | null>(null);
  measurementType = signal<string>('REP_BASED');
  sets = signal<number | null>(3);
  reps = signal<number | null>(10);
  weight = signal<number | null>(null);
  restBetweenSets = signal<number | null>(null);
  timePerRep = signal<number | null>(null);
  duration = signal<number | null>(null);
  distance = signal<number | null>(null);
  targetTime = signal<number | null>(null);

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

  selectedExerciseName = computed(() => {
    const id = this.userExerciseId();
    if (!id) return '';
    return this.enrichedUserExercises().find((ue) => ue.id === id)?.name ?? '';
  });

  canConfirm = computed(() => this.userExerciseId() != null);

  /** Creates a scheme from the field values and returns it. */
  async confirm(): Promise<UserExerciseScheme> {
    const data: Record<string, unknown> = {
      userExerciseId: this.userExerciseId(),
      measurementType: this.measurementType(),
    };
    const mt = this.measurementType();
    if (mt === 'REP_BASED') {
      if (this.sets() != null) data['sets'] = this.sets();
      if (this.reps() != null) data['reps'] = this.reps();
      if (this.weight() != null) data['weight'] = this.weight();
      if (this.restBetweenSets() != null) data['restBetweenSets'] = this.restBetweenSets();
    } else if (mt === 'TIME_BASED') {
      if (this.sets() != null) data['sets'] = this.sets();
      if (this.duration() != null) data['duration'] = this.duration();
      if (this.timePerRep() != null) data['timePerRep'] = this.timePerRep();
    } else if (mt === 'DISTANCE_BASED') {
      if (this.sets() != null) data['sets'] = this.sets();
      if (this.distance() != null) data['distance'] = this.distance();
      if (this.targetTime() != null) data['targetTime'] = this.targetTime();
    }
    const scheme = await this.userApi.createExerciseScheme(data);
    this.queryClient.invalidateQueries({ queryKey: exerciseSchemeKeys.all() });
    return scheme;
  }

  reset() {
    this.userExerciseId.set(null);
    this.measurementType.set('REP_BASED');
    this.sets.set(3);
    this.reps.set(10);
    this.weight.set(null);
    this.restBetweenSets.set(null);
    this.timePerRep.set(null);
    this.duration.set(null);
    this.distance.set(null);
    this.targetTime.set(null);
  }
}
