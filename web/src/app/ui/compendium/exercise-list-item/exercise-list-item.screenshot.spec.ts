import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideRouter } from '@angular/router';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { ExerciseListItem } from './exercise-list-item';
import { DataTable, DataTableColumn } from '$ui/data-table/data-table';
import { Exercise } from '$generated/models';

const exercise: Exercise = {
  id: 1,
  createdAt: '2024-01-15T10:00:00Z',
  updatedAt: '2024-01-15T10:00:00Z',
  names: [{ id: 1, name: 'Bench Press' }],
  type: 'STRENGTH',
  force: ['PUSH'],
  primaryMuscles: ['CHEST', 'TRICEPS'],
  secondaryMuscles: ['SHOULDERS'],
  technicalDifficulty: 'intermediate',
  bodyWeightScaling: 0,
  suggestedMeasurementParadigms: ['REP_BASED'],
  description: 'A compound upper body exercise',
  instructions: [],
  images: [],
  owner: 'admin',
  public: true,
  version: 1,
  equipmentIds: [],
};

const columns: DataTableColumn[] = [
  { label: 'Name' },
  { label: 'Mastery' },
  { label: 'Type' },
  { label: 'Difficulty' },
  { label: 'Force' },
  { label: 'Primary muscles' },
  { label: 'Secondary muscles', defaultHidden: true },
  { label: 'Slug', defaultHidden: true },
  { label: 'Body weight scaling', defaultHidden: true },
  { label: 'Measurement paradigms', defaultHidden: true },
  { label: 'Description', defaultHidden: true },
  { label: 'Alternative names', defaultHidden: true },
  { label: 'Author', defaultHidden: true },
  { label: 'Version', defaultHidden: true },
  { label: 'Created by', defaultHidden: true },
  { label: 'Created at', defaultHidden: true },
  { label: 'Updated at', defaultHidden: true },
];

describe('ExerciseListItem screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const template = `
    <app-data-table [columns]="columns">
      <tr app-exercise-list-item [exercise]="exercise" [displayName]="displayName"></tr>
    </app-data-table>
  `;

  const opts = {
    imports: [DataTable, ExerciseListItem],
    providers: [provideTranslocoForTest(), provideRouter([])],
    componentProperties: { exercise, columns, displayName: exercise.names?.[0]?.name ?? '' },
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
