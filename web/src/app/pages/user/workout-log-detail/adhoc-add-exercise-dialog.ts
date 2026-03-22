import { Component, inject, input, output, signal, ViewChild } from '@angular/core';
import { injectQueryClient } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutLogKeys } from '$core/query-keys';
import { HlmDialogImports } from '@spartan-ng/helm/dialog';
import { HlmButton } from '@spartan-ng/helm/button';
import { HlmSeparator } from '@spartan-ng/helm/separator';
import { ExerciseConfig } from '$ui/exercise-config/exercise-config';
import { ExerciseRunner } from '$ui/exercise-runner/exercise-runner';

@Component({
  selector: 'app-adhoc-add-exercise-dialog',
  imports: [HlmDialogImports, HlmButton, HlmSeparator, ExerciseConfig, ExerciseRunner],
  template: `
    <hlm-dialog [state]="open() ? 'open' : 'closed'" (closed)="onClose()">
      <ng-template hlmDialogPortal>
        <hlm-dialog-content [showCloseButton]="false" class="max-h-[90dvh] overflow-y-auto">
          <hlm-dialog-header>
            <h3 hlmDialogTitle>Add Exercise</h3>
          </hlm-dialog-header>

          <!-- Phase 1: Exercise configuration -->
          <app-exercise-config #exerciseConfig />

          <hlm-separator class="my-4" />

          <!-- Phase 2: Exercise sets planning -->
          @if (exerciseConfig.userExerciseId() && exerciseConfig.sets()) {
            <app-exercise-runner
              #runner
              [exerciseName]="exerciseConfig.selectedExerciseName()"
              [measurementType]="exerciseConfig.measurementType()"
              [setCount]="exerciseConfig.sets()!"
              [defaultReps]="exerciseConfig.reps()"
              [defaultWeight]="exerciseConfig.weight()"
              [defaultDuration]="exerciseConfig.duration()"
              [defaultDistance]="exerciseConfig.distance()"
              [defaultRest]="exerciseConfig.restBetweenSets()"
            />
          } @else {
            <div class="py-4 text-center text-sm text-gray-400 dark:text-gray-500">
              Select an exercise above to plan sets
            </div>
          }

          @if (errorMessage()) {
            <div
              class="mt-2 rounded-md border border-red-300 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-700 dark:bg-red-900/20 dark:text-red-400"
            >
              {{ errorMessage() }}
            </div>
          }

          <hlm-dialog-footer>
            <button hlmBtn variant="outline" hlmDialogClose>Cancel</button>
            <button
              hlmBtn
              (click)="onAdd()"
              [disabled]="!exerciseConfig.canConfirm() || isAdding()"
            >
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
export class AdhocAddExerciseDialog {
  private userApi = inject(UserApiClient);
  private queryClient = injectQueryClient();

  open = input(false);
  sectionId = input.required<number>();
  logId = input.required<number>();
  exerciseCount = input(0);

  closed = output<void>();

  @ViewChild('exerciseConfig') exerciseConfig!: ExerciseConfig;
  @ViewChild('runner') runner?: ExerciseRunner;

  isAdding = signal(false);
  errorMessage = signal('');

  async onAdd() {
    this.isAdding.set(true);
    this.errorMessage.set('');
    try {
      // 1. Create scheme from Phase 1 config
      const scheme = await this.exerciseConfig.confirm();

      // 2. Create workout log exercise (backend auto-creates sets from scheme)
      await this.userApi.createWorkoutLogExercise({
        workoutLogSectionId: this.sectionId(),
        sourceExerciseSchemeId: scheme.id,
        position: this.exerciseCount(),
      });

      // 3. Refresh parent log query
      this.queryClient.invalidateQueries({ queryKey: workoutLogKeys.detail(this.logId()) });

      // 4. Reset and close
      this.exerciseConfig.reset();
      this.runner?.reset();
      this.closed.emit();
    } catch (err) {
      console.error('Failed to add exercise:', err);
      this.errorMessage.set(err instanceof Error ? err.message : 'Failed to add exercise');
    } finally {
      this.isAdding.set(false);
    }
  }

  onClose() {
    this.errorMessage.set('');
    this.exerciseConfig?.reset();
    this.runner?.reset();
    this.closed.emit();
  }
}
