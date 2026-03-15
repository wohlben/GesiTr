import { Component, input } from '@angular/core';
import { ViewItemHeader } from './workout-log-view-items';

@Component({
  selector: 'app-workout-log-active-header',
  template: `
    <div class="mt-3 mb-1 text-xs font-medium text-gray-500 uppercase dark:text-gray-400">
      {{ data().exerciseName }}
    </div>
  `,
})
export class WorkoutLogActiveHeader {
  data = input.required<ViewItemHeader>();
}
