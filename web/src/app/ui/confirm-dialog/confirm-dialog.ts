import { Component, HostListener, input, output } from '@angular/core';

@Component({
  selector: 'app-confirm-dialog',
  template: `
    @if (open()) {
      <!-- eslint-disable-next-line @angular-eslint/template/click-events-have-key-events, @angular-eslint/template/interactive-supports-focus -->
      <div
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        (click)="cancelled.emit()"
      >
        <!-- eslint-disable-next-line @angular-eslint/template/click-events-have-key-events -->
        <div
          class="mx-4 w-full max-w-sm rounded-lg bg-white p-6 shadow-xl dark:bg-gray-800"
          (click)="$event.stopPropagation()"
          role="dialog"
          aria-modal="true"
        >
          <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100">{{ title() }}</h3>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ message() }}</p>
          <div class="mt-4 flex justify-end gap-2">
            <button
              type="button"
              (click)="cancelled.emit()"
              [disabled]="isPending()"
              class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
            >
              Cancel
            </button>
            <button
              type="button"
              (click)="confirmed.emit()"
              [disabled]="isPending()"
              class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 disabled:opacity-50"
            >
              {{ confirmLabel() }}
            </button>
          </div>
        </div>
      </div>
    }
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

  @HostListener('document:keydown.escape')
  onEscape() {
    if (this.open()) {
      this.cancelled.emit();
    }
  }
}
