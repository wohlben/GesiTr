import { Component, input, model } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-search-input',
  imports: [FormsModule],
  template: `
    <input
      type="text"
      [placeholder]="placeholder()"
      [ngModel]="value()"
      (ngModelChange)="value.set($event)"
      class="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm shadow-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100"
    />
  `,
})
export class SearchInput {
  placeholder = input('Search...');
  value = model('');
}
