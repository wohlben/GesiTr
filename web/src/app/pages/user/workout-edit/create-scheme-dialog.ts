import { Component, effect, input, output, signal, viewChild } from '@angular/core';
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
              <h3 hlmDialogTitle>
                {{
                  editingScheme()
                    ? t('user.workouts.editSchemeTitle')
                    : t('user.workouts.createSchemeTitle')
                }}
              </h3>
            </hlm-dialog-header>

            <app-exercise-config
              #exerciseConfig
              [preselectedExerciseId]="preselectedExerciseId()"
            />

            <hlm-dialog-footer>
              <button hlmBtn variant="outline" hlmDialogClose [disabled]="isSaving()">
                {{ t('common.cancel') }}
              </button>
              <button
                hlmBtn
                (click)="onConfirm()"
                [disabled]="!exerciseConfig.canConfirm() || isSaving()"
              >
                @if (isSaving()) {
                  {{ t('common.saving') }}
                } @else {
                  {{ t('common.save') }}
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
  editingScheme = input<ExerciseScheme | null>(null);

  schemeSaved = output<ExerciseScheme>();
  cancelled = output();

  exerciseConfig = viewChild.required<ExerciseConfig>('exerciseConfig');

  isSaving = signal(false);

  constructor() {
    // Prefill form when editing an existing scheme
    effect(() => {
      const scheme = this.editingScheme();
      if (scheme && this.open()) {
        // Defer to next tick so the viewChild is available
        queueMicrotask(() => this.exerciseConfig().prefill(scheme));
      }
    });
  }

  async onConfirm() {
    this.isSaving.set(true);
    try {
      const editing = this.editingScheme();
      const scheme = editing
        ? await this.exerciseConfig().update(editing.id)
        : await this.exerciseConfig().confirm();
      this.schemeSaved.emit(scheme);
      this.exerciseConfig().reset();
    } catch (err) {
      console.error('Failed to save scheme:', err);
    } finally {
      this.isSaving.set(false);
    }
  }

  onCancel() {
    this.exerciseConfig().reset();
    this.cancelled.emit();
  }
}
