import { expect, test } from '../../base-test';
import { createExercise, deleteExercise } from '../../helpers';

test.describe('/compendium/exercises?mastery=me', () => {
  test('renders user exercises list with expected content', async ({ request, page }) => {
    const names = ['My Bench Press', 'My Squat', 'My Deadlift'];
    const exercises: { id: number }[] = [];
    for (const name of names) {
      const exercise = await createExercise(request, { names: [name] });
      exercises.push(exercise);
    }
    await page.goto('/compendium/exercises?mastery=me', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('Exercises');
    await expect(page.locator('table tbody tr')).toHaveCount(names.length);
    await expect(page.locator('table')).toContainText('My Bench Press');
    for (const exercise of exercises) {
      await deleteExercise(request, exercise.id);
    }
  });
});
