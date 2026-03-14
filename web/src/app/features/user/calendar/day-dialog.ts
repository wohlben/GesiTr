import { Component, HostListener, input, output } from '@angular/core';
import { DatePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { WorkoutLog } from '$generated/user-models';

@Component({
  selector: 'app-day-dialog',
  imports: [DatePipe, RouterLink],
  template: `
    @if (open()) {
      <!-- eslint-disable-next-line @angular-eslint/template/click-events-have-key-events, @angular-eslint/template/interactive-supports-focus -->
      <div
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        (click)="closed.emit()"
      >
        <!-- eslint-disable-next-line @angular-eslint/template/click-events-have-key-events -->
        <div
          class="mx-4 w-full max-w-sm rounded-lg bg-white p-6 shadow-xl dark:bg-gray-800"
          (click)="$event.stopPropagation()"
          role="dialog"
          aria-modal="true"
        >
          <div class="mb-4 flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100">
              {{ date() | date }}
            </h3>
            <button
              type="button"
              (click)="closed.emit()"
              class="rounded-md p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>

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
        </div>
      </div>
    }
  `,
})
export class DayDialog {
  open = input(false);
  date = input<Date>();
  logs = input<WorkoutLog[]>([]);

  closed = output();

  @HostListener('document:keydown.escape')
  onEscape() {
    if (this.open()) {
      this.closed.emit();
    }
  }

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
