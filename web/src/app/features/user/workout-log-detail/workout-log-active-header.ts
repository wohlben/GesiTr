import { Component, input, output } from '@angular/core';
import { ViewItemHeader } from './workout-log-view-items';

@Component({
  selector: 'app-workout-log-active-header',
  template: `
    @if (data().hasOverride) {
      <div
        class="mt-3 mb-1 flex cursor-pointer items-center gap-1 text-xs font-medium text-blue-600 uppercase dark:text-blue-400"
        role="button"
        tabindex="0"
        (click)="resetOverride.emit()"
        (keydown.enter)="resetOverride.emit()"
      >
        <svg class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
          <path
            fill-rule="evenodd"
            d="M7.793 2.232a.75.75 0 01-.025 1.06L3.622 7.25h10.003a5.375 5.375 0 010 10.75H10.75a.75.75 0 010-1.5h2.875a3.875 3.875 0 000-7.75H3.622l4.146 3.957a.75.75 0 01-1.036 1.085l-5.5-5.25a.75.75 0 010-1.085l5.5-5.25a.75.75 0 011.06.025z"
            clip-rule="evenodd"
          />
        </svg>
        {{ data().exerciseName }}
      </div>
    } @else {
      <div class="mt-3 mb-1 text-xs font-medium text-gray-500 uppercase dark:text-gray-400">
        {{ data().exerciseName }}
      </div>
    }
  `,
})
export class WorkoutLogActiveHeader {
  data = input.required<ViewItemHeader>();
  resetOverride = output<void>();
}
