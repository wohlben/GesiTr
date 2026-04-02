import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideExternalLink } from '@ng-icons/lucide';
import { ExerciseGroup } from '$generated/user-models';
import { HlmIcon } from '$ui/spartan/icon/src/lib/hlm-icon';
import { SlugifyPipe } from '$ui/pipes/slugify';

@Component({
  selector: 'tr[app-exercise-group-list-item]',
  imports: [RouterLink, SlugifyPipe, NgIcon, HlmIcon],
  providers: [provideIcons({ lucideExternalLink })],
  template: `
    <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">
      <span class="inline-flex items-center gap-1">
        {{ group().name }}
        <a
          [routerLink]="[
            '/compendium/exercise-groups',
            group().id,
            group().name || 'group' | slugify,
          ]"
          class="inline-flex items-center justify-center rounded-full p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
        >
          <ng-icon hlm name="lucideExternalLink" size="sm" />
        </a>
      </span>
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ group().owner }}
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
