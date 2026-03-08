import { Component, input } from '@angular/core';
import { Exercise } from '$generated/models';

@Component({
  selector: 'tr[app-exercise-list-item]',
  template: `
    <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">
      {{ exercise().name }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().type }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().technicalDifficulty }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().force.join(', ') }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().primaryMuscles.join(', ') }}
    </td>
  `,
})
export class ExerciseListItem {
  exercise = input.required<Exercise>();
}
