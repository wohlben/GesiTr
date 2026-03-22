import { Component, inject, input, output, signal, viewChild } from '@angular/core';
import { QueryClient } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { UserExerciseScheme } from '$generated/user-models';
import { HlmDialogImports } from '@spartan-ng/helm/dialog';
import { HlmButton } from '@spartan-ng/helm/button';
import { ExerciseConfig } from '$ui/exercise-config/exercise-config';

@Component({
  selector: 'app-add-exercise-dialog',
  imports: [HlmDialogImports, HlmButton, ExerciseConfig, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <hlm-dialog [state]="open() ? 'open' : 'closed'" (closed)="onCancel()">
        <ng-template hlmDialogPortal>
          <hlm-dialog-content [showCloseButton]="false">
            <hlm-dialog-header>
              <h3 hlmDialogTitle>{{ t('user.workoutLog.addExerciseTitle') }}</h3>
            </hlm-dialog-header>

            <app-exercise-config #exerciseConfig />

            <!-- Actions -->
            <hlm-dialog-footer>
              <button hlmBtn variant="outline" hlmDialogClose [disabled]="isAdding()">
                {{ t('common.cancel') }}
              </button>
              <button
                hlmBtn
                (click)="onAdd()"
                [disabled]="!exerciseConfig.canConfirm() || isAdding()"
              >
                @if (isAdding()) {
                  {{ t('common.adding') }}
                } @else {
                  {{ t('common.add') }}
                }
              </button>
            </hlm-dialog-footer>
          </hlm-dialog-content>
        </ng-template>
      </hlm-dialog>
    </ng-container>
  `,
})
export class AddExerciseDialog {
  private userApi = inject(UserApiClient);
  private queryClient = inject(QueryClient);

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

  exerciseConfig = viewChild.required<ExerciseConfig>('exerciseConfig');

  isAdding = signal(false);

  async onAdd() {
    this.isAdding.set(true);
    try {
      const scheme = await this.exerciseConfig().confirm();

      // Create log exercise
      const logExercise = await this.userApi.createWorkoutLogExercise({
        workoutLogSectionId: this.sectionId(),
        sourceExerciseSchemeId: scheme.id,
        position: this.exerciseCount(),
      });

      // Fetch the created exercise with its sets
      const logs = await this.userApi.fetchWorkoutLogs({ status: 'planning' });
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
      for (const log of logs) {
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
        exerciseName: this.exerciseConfig().selectedExerciseName(),
        scheme,
        exercise: createdExercise ?? {
          id: logExercise.id,
          sourceExerciseSchemeId: scheme.id,
          sets: [],
        },
      });

      this.exerciseConfig().reset();
    } catch (err) {
      console.error('Failed to add exercise:', err);
    } finally {
      this.isAdding.set(false);
    }
  }

  onCancel() {
    this.exerciseConfig().reset();
    this.cancelled.emit();
  }
}
