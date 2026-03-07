import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { DataTable } from './data-table';

describe('DataTable screenshots', () => {
  afterEach(() => {
    document.documentElement.style.colorScheme = '';
  });

  const template = `
    <app-data-table [columns]="['Name', 'Category', 'Description']">
      <tr>
        <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">Barbell</td>
        <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">free_weights</td>
        <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">Standard Olympic barbell</td>
      </tr>
      <tr>
        <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">Dumbbell</td>
        <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">free_weights</td>
        <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">Adjustable dumbbell</td>
      </tr>
    </app-data-table>
  `;

  it('light', async () => {
    const { fixture } = await render(template, { imports: [DataTable] });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.style.colorScheme = 'dark';
    const { fixture } = await render(template, { imports: [DataTable] });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
  });
});
