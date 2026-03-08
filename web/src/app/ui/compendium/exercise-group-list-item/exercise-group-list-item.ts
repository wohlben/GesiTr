import { Component, input } from '@angular/core';
import { ExerciseGroup } from '$generated/models';

@Component({
  selector: 'tr[app-exercise-group-list-item]',
  template: `
    <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">
      {{ group().name }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">{{ group().description }}</td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ group().createdBy }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ group().createdAt }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ group().updatedAt }}
    </td>
  `,
})
export class ExerciseGroupListItem {
  group = input.required<ExerciseGroup>();
}
