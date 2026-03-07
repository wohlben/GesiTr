import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { ExerciseGroupListItem } from './exercise-group-list-item';
import { DataTable } from '$ui/data-table/data-table';
import { ExerciseGroup } from '$generated/models';

const group: ExerciseGroup = {
  id: 1,
  createdAt: '',
  updatedAt: '',
  templateId: '',
  name: 'Push Exercises',
  description: 'All pushing movements including bench press variations',
  createdBy: '',
};

describe('ExerciseGroupListItem screenshots', () => {
  afterEach(() => {
    document.documentElement.style.colorScheme = '';
  });

  const template = `
    <app-data-table [columns]="['Name', 'Description']">
      <tr app-exercise-group-list-item [group]="group"></tr>
    </app-data-table>
  `;

  it('light', async () => {
    const { fixture } = await render(template, {
      imports: [DataTable, ExerciseGroupListItem],
      componentProperties: { group },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.style.colorScheme = 'dark';
    const { fixture } = await render(template, {
      imports: [DataTable, ExerciseGroupListItem],
      componentProperties: { group },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
  });
});
