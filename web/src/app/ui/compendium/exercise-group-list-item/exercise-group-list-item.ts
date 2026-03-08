import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { ExerciseGroup } from '$generated/models';
import { SlugifyPipe } from '$ui/pipes/slugify';

@Component({
  selector: 'tr[app-exercise-group-list-item]',
  imports: [RouterLink, SlugifyPipe],
  template: `
    <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">
      <span class="inline-flex items-center gap-1">
        {{ group().name }}
        <a
          [routerLink]="['/compendium/exercise-groups', group().id, group().name | slugify]"
          class="inline-flex items-center justify-center rounded-full p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
            class="h-4 w-4"
          >
            <path
              fill-rule="evenodd"
              d="M4.25 5.5a.75.75 0 0 0-.75.75v8.5c0 .414.336.75.75.75h8.5a.75.75 0 0 0 .75-.75v-4a.75.75 0 0 1 1.5 0v4A2.25 2.25 0 0 1 12.75 17h-8.5A2.25 2.25 0 0 1 2 14.75v-8.5A2.25 2.25 0 0 1 4.25 4h5a.75.75 0 0 1 0 1.5h-5Zm7.97-2.03a.75.75 0 0 1 1.06 0l4.5 4.5a.75.75 0 0 1-1.06 1.06l-3.97-3.97v7.19a.75.75 0 0 1-1.5 0V5.06L7.28 9.03a.75.75 0 0 1-1.06-1.06l4.5-4.5Z"
              clip-rule="evenodd"
            />
          </svg>
        </a>
      </span>
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
