import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideExternalLink } from '@ng-icons/lucide';
import { Exercise } from '$generated/models';
import { HlmIcon } from '$ui/spartan/icon/src/lib/hlm-icon';
import { SlugifyPipe } from '$ui/pipes/slugify';

@Component({
  selector: 'tr[app-exercise-list-item]',
  imports: [RouterLink, SlugifyPipe, NgIcon, HlmIcon],
  providers: [provideIcons({ lucideExternalLink })],
  template: `
    <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">
      <span class="inline-flex items-center gap-1">
        {{ exercise().name }}
        <a
          [routerLink]="['/compendium/exercises', exercise().id, exercise().name | slugify]"
          class="inline-flex items-center justify-center rounded-full p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
        >
          <ng-icon hlm name="lucideExternalLink" size="sm" />
        </a>
      </span>
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
      {{ exercise().owner }}
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
