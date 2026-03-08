import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideRouter } from '@angular/router';
import { ExerciseListItem } from './exercise-list-item';
import { DataTable, DataTableColumn } from '$ui/data-table/data-table';
import { Exercise } from '$generated/models';

const exercise: Exercise = {
  id: 1,
  createdAt: '2024-01-15T10:00:00Z',
  updatedAt: '2024-01-15T10:00:00Z',
  name: 'Bench Press',
  slug: 'bench-press',
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
  alternativeNames: [],
  createdBy: 'admin',
  version: 1,
  equipmentIds: [],
};

const columns: DataTableColumn[] = [
  { label: 'Name' },
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
      <tr app-exercise-list-item [exercise]="exercise"></tr>
    </app-data-table>
  `;

  const opts = {
    imports: [DataTable, ExerciseListItem],
    providers: [provideRouter([])],
    componentProperties: { exercise, columns },
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
