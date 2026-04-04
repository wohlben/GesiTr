import { expect, test } from '../../base-test';
import {
  createExercise,
  deleteExercise,
  createWorkout,
  deleteWorkout,
  createExerciseScheme,
  deleteExerciseScheme,
  createWorkoutSection,
  deleteWorkoutSection,
  createWorkoutSectionItem,
  deleteWorkoutSectionItem,
  upsertSchemeSectionItem,
} from '../../helpers';

test.describe('/compendium/workouts/[id]/edit', () => {
  test('renders edit form with workout data', async ({ request, page }) => {
    const exercise = await createExercise(request, { names: ['WE Bench Press'] });
    const scheme = await createExerciseScheme(request, {
      exerciseId: exercise.id,
    });
    const workout = await createWorkout(request, {
      name: 'Edit Push Day',
    });
    const section = await createWorkoutSection(request, {
      workoutId: workout.id,
      label: 'Main Lifts',
    });
    const sectionExercise = await createWorkoutSectionItem(request, {
      workoutSectionId: section.id,
      exerciseId: exercise.id,
    });
    await upsertSchemeSectionItem(request, {
      exerciseSchemeId: scheme.id,
      workoutSectionItemId: sectionExercise.id,
    });

    await page.goto(`/compendium/workouts/${workout.id}/edit`, { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('Edit Workout');
    await expect(page.locator('input#name')).toHaveValue('Edit Push Day');
    await expect(page.locator('form')).toBeVisible();

    await deleteWorkoutSectionItem(request, sectionExercise.id);
    await deleteWorkoutSection(request, section.id);
    await deleteWorkout(request, workout.id);
    await deleteExerciseScheme(request, scheme.id);
    await deleteExercise(request, exercise.id);
  });
});
