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

test.describe('/compendium/workouts', () => {
  test('renders workouts list with expected content', async ({ request, page }) => {
    const variant = [
      { workoutName: 'Push Day A', exerciseName: 'WL Bench Press' },
      { workoutName: 'Pull Day A', exerciseName: 'WL Barbell Row' },
    ];
    const items: {
      exercise: { id: number };
      scheme: { id: number };
      workout: { id: number };
      section: { id: number };
      sectionExercise: { id: number };
    }[] = [];

    for (const v of variant) {
      const exercise = await createExercise(request, { names: [v.exerciseName] });
      const scheme = await createExerciseScheme(request, {
        exerciseId: exercise.id,
      });
      const workout = await createWorkout(request, {
        name: v.workoutName,
      });
      const section = await createWorkoutSection(request, {
        workoutId: workout.id,
      });
      const sectionExercise = await createWorkoutSectionItem(request, {
        workoutSectionId: section.id,
        exerciseId: exercise.id,
      });
      await upsertSchemeSectionItem(request, {
        exerciseSchemeId: scheme.id,
        workoutSectionItemId: sectionExercise.id,
      });
      items.push({ exercise, scheme, workout, section, sectionExercise });
    }

    await page.goto('/compendium/workouts', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('Workouts');
    await expect(page.locator('table tbody tr')).toHaveCount(variant.length);
    await expect(page.locator('table')).toContainText('Push Day A');

    for (const item of items) {
      await deleteWorkoutSectionItem(request, item.sectionExercise.id);
      await deleteWorkoutSection(request, item.section.id);
      await deleteWorkout(request, item.workout.id);
      await deleteExerciseScheme(request, item.scheme.id);
      await deleteExercise(request, item.exercise.id);
    }
  });
});
