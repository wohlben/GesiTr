import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { EquipmentListItem } from './equipment-list-item';
import { DataTable } from '$ui/data-table/data-table';
import { Equipment } from '$generated/models';

const equipment: Equipment = {
  id: 1,
  createdAt: '',
  updatedAt: '',
  name: 'barbell',
  displayName: 'Barbell',
  description: 'Standard Olympic barbell, 20kg',
  category: 'free_weights',
  templateId: '',
  createdBy: '',
  version: 1,
};

describe('EquipmentListItem screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const template = `
    <app-data-table [columns]="['Name', 'Category', 'Description']">
      <tr app-equipment-list-item [equipment]="equipment"></tr>
    </app-data-table>
  `;

  it('light', async () => {
    const { fixture } = await render(template, {
      imports: [DataTable, EquipmentListItem],
      componentProperties: { equipment },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    const { fixture } = await render(template, {
      imports: [DataTable, EquipmentListItem],
      componentProperties: { equipment },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
  });
});
