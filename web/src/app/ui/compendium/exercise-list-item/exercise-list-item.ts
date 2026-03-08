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
      {{ exercise().force?.join(', ') }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().primaryMuscles?.join(', ') }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().secondaryMuscles?.join(', ') }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().slug }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().bodyWeightScaling }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().suggestedMeasurementParadigms?.join(', ') }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().description }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().alternativeNames?.join(', ') }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().authorName }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().version }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().createdBy }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().createdAt }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().updatedAt }}
    </td>
  `,
})
export class ExerciseListItem {
  exercise = input.required<Exercise>();
}
