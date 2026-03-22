import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideExternalLink } from '@ng-icons/lucide';
import { Equipment } from '$generated/models';
import { HlmIcon } from '$ui/spartan/icon/src/lib/hlm-icon';
import { SlugifyPipe } from '$ui/pipes/slugify';

@Component({
  selector: 'tr[app-equipment-list-item]',
  imports: [RouterLink, SlugifyPipe, NgIcon, HlmIcon],
  providers: [provideIcons({ lucideExternalLink })],
  template: `
    <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">
      <span class="inline-flex items-center gap-1">
        {{ equipment().displayName }}
        <a
          [routerLink]="[
            '/compendium/equipment',
            equipment().id,
            equipment().displayName | slugify,
          ]"
          class="inline-flex items-center justify-center rounded-full p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
        >
          <ng-icon hlm name="lucideExternalLink" size="sm" />
        </a>
      </span>
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().category }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().description }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().name }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().version }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().createdBy }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().createdAt }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().updatedAt }}
    </td>
  `,
})
export class EquipmentListItem {
  equipment = input.required<Equipment>();
}
