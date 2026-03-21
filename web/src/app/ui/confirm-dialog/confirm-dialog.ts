import { Component, input, output } from '@angular/core';
import { HlmDialogImports } from '@spartan-ng/helm/dialog';
import { HlmButton } from '@spartan-ng/helm/button';

@Component({
  selector: 'app-confirm-dialog',
  imports: [HlmDialogImports, HlmButton],
  template: `
    <hlm-dialog [state]="open() ? 'open' : 'closed'" (closed)="cancelled.emit()">
      <ng-template hlmDialogPortal>
        <hlm-dialog-content [showCloseButton]="false">
          <hlm-dialog-header>
            <h3 hlmDialogTitle>{{ title() }}</h3>
            <p hlmDialogDescription>{{ message() }}</p>
          </hlm-dialog-header>
          <hlm-dialog-footer>
            <button hlmBtn variant="outline" hlmDialogClose [disabled]="isPending()">Cancel</button>
            <button
              hlmBtn
              variant="destructive"
              (click)="confirmed.emit()"
              [disabled]="isPending()"
            >
              {{ confirmLabel() }}
            </button>
          </hlm-dialog-footer>
        </hlm-dialog-content>
      </ng-template>
    </hlm-dialog>
  `,
})
export class ConfirmDialog {
  open = input(false);
  title = input('Confirm');
  message = input('Are you sure?');
  confirmLabel = input('Delete');
  isPending = input(false);

  confirmed = output();
  cancelled = output();
}
