import { Component, input, model } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-filter-select',
  imports: [FormsModule],
  template: `
    <select
      [ngModel]="value()"
      (ngModelChange)="value.set($event)"
      class="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm shadow-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100"
    >
      <option value="">{{ allLabel() }}</option>
      @for (opt of options(); track opt) {
        <option [value]="opt">{{ opt }}</option>
      }
    </select>
  `,
})
export class FilterSelect {
  allLabel = input('All');
  options = input<string[]>([]);
  value = model('');
}
