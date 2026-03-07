import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { PageLayout } from './page-layout';
import { DataTable } from '$ui/data-table/data-table';

describe('PageLayout screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  it('with content - light', async () => {
    const template = `
      <app-page-layout header="Equipment">
        <div filters class="flex flex-wrap gap-3">
          <input type="text" placeholder="Search..." class="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
        </div>
        <app-data-table [columns]="['Name', 'Category']">
          <tr>
            <td class="px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">Barbell</td>
            <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">free_weights</td>
          </tr>
        </app-data-table>
        <p class="text-sm text-gray-500 dark:text-gray-400">1 item</p>
      </app-page-layout>
    `;
    const { fixture } = await render(template, { imports: [PageLayout, DataTable] });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('content-light');
  });

  it('with content - dark', async () => {
    document.documentElement.classList.add('dark');
    const template = `
      <app-page-layout header="Equipment">
        <div filters class="flex flex-wrap gap-3">
          <input type="text" placeholder="Search..." class="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
        </div>
        <app-data-table [columns]="['Name', 'Category']">
          <tr>
            <td class="px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">Barbell</td>
            <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">free_weights</td>
          </tr>
        </app-data-table>
        <p class="text-sm text-gray-500 dark:text-gray-400">1 item</p>
      </app-page-layout>
    `;
    const { fixture } = await render(template, { imports: [PageLayout, DataTable] });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('content-dark');
  });

  it('pending - light', async () => {
    const template = `
      <app-page-layout header="Equipment" [isPending]="true">
        <div filters class="flex flex-wrap gap-3">
          <input type="text" placeholder="Search..." class="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
        </div>
      </app-page-layout>
    `;
    const { fixture } = await render(template, { imports: [PageLayout] });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('pending-light');
  });

  it('pending - dark', async () => {
    document.documentElement.classList.add('dark');
    const template = `
      <app-page-layout header="Equipment" [isPending]="true">
        <div filters class="flex flex-wrap gap-3">
          <input type="text" placeholder="Search..." class="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
        </div>
      </app-page-layout>
    `;
    const { fixture } = await render(template, { imports: [PageLayout] });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('pending-dark');
  });

  it('error - light', async () => {
    const template = `
      <app-page-layout header="Equipment" errorMessage="Failed to load equipment.">
        <div filters class="flex flex-wrap gap-3">
          <input type="text" placeholder="Search..." class="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
        </div>
      </app-page-layout>
    `;
    const { fixture } = await render(template, { imports: [PageLayout] });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('error-light');
  });

  it('error - dark', async () => {
    document.documentElement.classList.add('dark');
    const template = `
      <app-page-layout header="Equipment" errorMessage="Failed to load equipment.">
        <div filters class="flex flex-wrap gap-3">
          <input type="text" placeholder="Search..." class="rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
        </div>
      </app-page-layout>
    `;
    const { fixture } = await render(template, { imports: [PageLayout] });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('error-dark');
  });
});
