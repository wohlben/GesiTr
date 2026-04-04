import { Component, input, output, signal, viewChild } from '@angular/core';
import { TranslocoDirective } from '@jsverse/transloco';
import { ExerciseScheme } from '$generated/user-exercisescheme';
import { HlmDialogImports } from '@spartan-ng/helm/dialog';
import { HlmButton } from '@spartan-ng/helm/button';
import { ExerciseConfig } from '$ui/exercise-config/exercise-config';

@Component({
  selector: 'app-create-scheme-dialog',
  imports: [HlmDialogImports, HlmButton, ExerciseConfig, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <hlm-dialog [state]="open() ? 'open' : 'closed'" (closed)="onCancel()">
        <ng-template hlmDialogPortal>
          <hlm-dialog-content [showCloseButton]="false">
            <hlm-dialog-header>
              <h3 hlmDialogTitle>{{ t('user.workouts.createSchemeTitle') }}</h3>
            </hlm-dialog-header>

            <app-exercise-config
              #exerciseConfig
              [preselectedExerciseId]="preselectedExerciseId()"
            />

            <hlm-dialog-footer>
              <button hlmBtn variant="outline" hlmDialogClose [disabled]="isCreating()">
                {{ t('common.cancel') }}
              </button>
              <button
                hlmBtn
                (click)="onConfirm()"
                [disabled]="!exerciseConfig.canConfirm() || isCreating()"
              >
                @if (isCreating()) {
                  {{ t('common.creating') }}
                } @else {
                  {{ t('common.create') }}
                }
              </button>
            </hlm-dialog-footer>
          </hlm-dialog-content>
        </ng-template>
      </hlm-dialog>
    </ng-container>
  `,
})
export class CreateSchemeDialog {
  open = input(false);
  preselectedExerciseId = input<number | null>(null);

  schemeCreated = output<ExerciseScheme>();
  cancelled = output();

  exerciseConfig = viewChild.required<ExerciseConfig>('exerciseConfig');

  isCreating = signal(false);

  async onConfirm() {
    this.isCreating.set(true);
    try {
      const scheme = await this.exerciseConfig().confirm();
      this.schemeCreated.emit(scheme);
      this.exerciseConfig().reset();
    } catch (err) {
      console.error('Failed to create scheme:', err);
    } finally {
      this.isCreating.set(false);
    }
  }

  onCancel() {
    this.exerciseConfig().reset();
    this.cancelled.emit();
  }
}
