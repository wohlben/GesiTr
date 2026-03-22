import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideRouter } from '@angular/router';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { ExerciseGroupListItem } from './exercise-group-list-item';
import { DataTable, DataTableColumn } from '$ui/data-table/data-table';
import { ExerciseGroup } from '$generated/models';

const group: ExerciseGroup = {
  id: 1,
  createdAt: '2024-01-15T10:00:00Z',
  updatedAt: '2024-01-15T10:00:00Z',
  templateId: '',
  name: 'Push Exercises',
  description: 'All pushing movements including bench press variations',
  createdBy: 'admin',
};

const columns: DataTableColumn[] = [
  { label: 'Name' },
  { label: 'Description' },
  { label: 'Created by', defaultHidden: true },
  { label: 'Created at', defaultHidden: true },
  { label: 'Updated at', defaultHidden: true },
];

describe('ExerciseGroupListItem screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const template = `
    <app-data-table [columns]="columns">
      <tr app-exercise-group-list-item [group]="group"></tr>
    </app-data-table>
  `;

  const opts = {
    imports: [DataTable, ExerciseGroupListItem],
    providers: [provideTranslocoForTest(), provideRouter([])],
    componentProperties: { group, columns },
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
