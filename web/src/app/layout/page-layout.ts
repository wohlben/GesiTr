import { Component, input } from '@angular/core';
import { LoadingSpinner } from '$ui/loading-spinner/loading-spinner';
import { ErrorMessage } from '$ui/error-message/error-message';

@Component({
  selector: 'app-page-layout',
  imports: [LoadingSpinner, ErrorMessage],
  template: `
    <div class="space-y-4">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ header() }}</h1>
        <div class="sticky top-0 z-10 -mx-6 bg-gray-50 px-6 py-2 empty:hidden dark:bg-gray-950 sm:static sm:mx-0 sm:bg-transparent sm:p-0 sm:dark:bg-transparent">
          <ng-content select="[actions]" />
        </div>
      </div>

      <ng-content select="[filters]" />

      @if (isPending()) {
        <app-loading-spinner />
      } @else if (errorMessage()) {
        <app-error-message [message]="errorMessage()!" />
      } @else {
        <ng-content />
      }
    </div>
  `,
})
export class PageLayout {
  header = input.required<string>();
  isPending = input(false);
  errorMessage = input<string>();
}
