import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideRouter } from '@angular/router';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { EquipmentListItem } from './equipment-list-item';
import { DataTable, DataTableColumn } from '$ui/data-table/data-table';
import { Equipment } from '$generated/models';

const equipment: Equipment = {
  id: 1,
  createdAt: '2024-01-15T10:00:00Z',
  updatedAt: '2024-01-15T10:00:00Z',
  name: 'barbell',
  displayName: 'Barbell',
  description: 'Standard Olympic barbell, 20kg',
  category: 'free_weights',
  templateId: '',
  owner: 'admin',
  public: true,
  version: 1,
};

const columns: DataTableColumn[] = [
  { label: 'Name' },
  { label: 'Category' },
  { label: 'Description' },
  { label: 'Internal name', defaultHidden: true },
  { label: 'Version', defaultHidden: true },
  { label: 'Created by', defaultHidden: true },
  { label: 'Created at', defaultHidden: true },
  { label: 'Updated at', defaultHidden: true },
];

describe('EquipmentListItem screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const template = `
    <app-data-table [columns]="columns">
      <tr app-equipment-list-item [equipment]="equipment"></tr>
    </app-data-table>
  `;

  const opts = {
    imports: [DataTable, EquipmentListItem],
    providers: [provideTranslocoForTest(), provideRouter([])],
    componentProperties: { equipment, columns },
  };

  it('light', async () => {
    const { fixture } = await render(template, opts);
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    const { fixture } = await render(template, opts);
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
  });
});
