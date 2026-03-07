import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { ExerciseListItem } from './exercise-list-item';
import { DataTable } from '$ui/data-table/data-table';
import { Exercise } from '$generated/models';

const exercise: Exercise = {
  id: 1,
  createdAt: '',
  updatedAt: '',
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
  createdBy: '',
  version: 1,
  equipmentIds: [],
};

describe('ExerciseListItem screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const template = `
    <app-data-table [columns]="['Name', 'Type', 'Difficulty', 'Force', 'Primary muscles']">
      <tr app-exercise-list-item [exercise]="exercise"></tr>
    </app-data-table>
  `;

  it('light', async () => {
    const { fixture } = await render(template, {
      imports: [DataTable, ExerciseListItem],
      componentProperties: { exercise },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    const { fixture } = await render(template, {
      imports: [DataTable, ExerciseListItem],
      componentProperties: { exercise },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
  });
});
