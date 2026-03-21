import { Component, input } from '@angular/core';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideCircleX } from '@ng-icons/lucide';
import { HlmSpinner } from '@spartan-ng/helm/spinner';
import { HlmAlertImports } from '@spartan-ng/helm/alert';
import { HlmIconImports } from '@spartan-ng/helm/icon';

@Component({
  selector: 'app-page-layout',
  imports: [HlmSpinner, HlmAlertImports, HlmIconImports, NgIcon],
  providers: [provideIcons({ lucideCircleX })],
  template: `
    <div class="space-y-4">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ header() }}</h1>
        <div
          class="sticky top-0 z-10 -mx-6 bg-gray-50 px-6 py-2 empty:hidden dark:bg-gray-950 sm:static sm:mx-0 sm:bg-transparent sm:p-0 sm:dark:bg-transparent"
        >
          <ng-content select="[actions]" />
        </div>
      </div>

      <ng-content select="[filters]" />

      @if (isPending()) {
        <hlm-spinner />
      } @else if (errorMessage()) {
        <hlm-alert variant="destructive">
          <ng-icon hlm name="lucideCircleX" />
          <p hlmAlertDescription>{{ errorMessage() }}</p>
        </hlm-alert>
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
