import { Component, input, output } from '@angular/core';
import { DatePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { WorkoutLog } from '$generated/user-models';
import { HlmDialogImports } from '@spartan-ng/helm/dialog';

@Component({
  selector: 'app-day-dialog',
  imports: [DatePipe, RouterLink, HlmDialogImports],
  template: `
    <hlm-dialog [state]="open() ? 'open' : 'closed'" (closed)="closed.emit()">
      <ng-template hlmDialogPortal>
        <hlm-dialog-content>
          <hlm-dialog-header>
            <h3 hlmDialogTitle>{{ date() | date }}</h3>
          </hlm-dialog-header>

          <div class="space-y-2">
            @for (log of logs(); track log.id) {
              <a
                [routerLink]="['/user/workout-logs', log.id]"
                (click)="closed.emit()"
                class="flex items-center justify-between rounded-md border border-gray-200 px-3 py-2 transition-colors hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-700/50"
              >
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{
                  log.name
                }}</span>
                <span
                  class="rounded-full px-2 py-0.5 text-xs font-medium"
                  [class]="statusClass(log.status)"
                >
                  {{ statusLabel(log.status) }}
                </span>
              </a>
            }
          </div>
        </hlm-dialog-content>
      </ng-template>
    </hlm-dialog>
  `,
})
export class DayDialog {
  open = input(false);
  date = input<Date>();
  logs = input<WorkoutLog[]>([]);

  closed = output();

  statusClass(status: string): string {
    switch (status) {
      case 'finished':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300';
      case 'in_progress':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300';
      case 'aborted':
        return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300';
      default:
        return 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400';
    }
  }

  statusLabel(status: string): string {
    return status.replaceAll('_', ' ');
  }
}
