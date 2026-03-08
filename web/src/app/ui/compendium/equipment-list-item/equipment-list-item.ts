import { Component, input } from '@angular/core';
import { Equipment } from '$generated/models';

@Component({
  selector: 'tr[app-equipment-list-item]',
  template: `
    <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">
      {{ equipment().displayName }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().category }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ equipment().description }}
    </td>
  `,
})
export class EquipmentListItem {
  equipment = input.required<Equipment>();
}
