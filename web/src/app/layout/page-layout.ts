import { Component, input } from '@angular/core';
import { LoadingSpinner } from '$ui/loading-spinner/loading-spinner';
import { ErrorMessage } from '$ui/error-message/error-message';

@Component({
  selector: 'app-page-layout',
  imports: [LoadingSpinner, ErrorMessage],
  template: `
    <div class="space-y-4">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ header() }}</h1>

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
