import { Component, input } from '@angular/core';

@Component({
  selector: 'app-data-table',
  template: `
    <div class="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-800">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-800">
        <thead class="bg-gray-50 dark:bg-gray-900">
          <tr>
            @for (col of columns(); track col) {
              <th class="px-4 py-2 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400">{{ col }}</th>
            }
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 bg-white dark:divide-gray-800 dark:bg-gray-950">
          <ng-content />
        </tbody>
      </table>
    </div>
  `,
})
export class DataTable {
  columns = input.required<string[]>();
}
