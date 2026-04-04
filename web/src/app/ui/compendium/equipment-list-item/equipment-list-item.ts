import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideExternalLink } from '@ng-icons/lucide';
import { Equipment } from '$generated/models';
import { EquipmentMastery } from '$generated/user-mastery';
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
      @if (mastery(); as m) {
        <span
          class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
          [class]="tierBadgeClass(m.tier)"
        >
          Lv.{{ m.level }}
        </span>
      }
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
      {{ equipment().createdAt }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().updatedAt }}
    </td>
  `,
})
export class EquipmentListItem {
  equipment = input.required<Equipment>();
  mastery = input<EquipmentMastery | undefined>(undefined);

  tierBadgeClass(tier: string): string {
    switch (tier) {
      case 'mastered':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300';
      case 'master':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300';
      case 'adept':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300';
      case 'journeyman':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300';
      default:
        return 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400';
    }
  }
}
