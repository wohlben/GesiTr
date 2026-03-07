import { Component, input } from '@angular/core';

@Component({
  selector: 'app-error-message',
  template: `
    <div class="flex items-center gap-2 rounded-md border border-red-200 bg-red-50 px-4 py-3 dark:border-red-900 dark:bg-red-950">
      <svg class="h-5 w-5 shrink-0 text-red-500 dark:text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
      </svg>
      <span class="text-sm text-red-700 dark:text-red-300">{{ message() }}</span>
    </div>
  `,
})
export class ErrorMessage {
  message = input.required<string>();
}
