import { expect, test } from '../../base-test';
import { createExercise, deleteExercise } from '../../helpers';

test('renders exercises list with expected content', async ({ request, page }) => {
  const names = ['Bench Press', 'Squat', 'Deadlift'];
  const items = [];
  for (const name of names) {
    items.push(await createExercise(request, { names: [name] }));
  }
  await page.goto('/compendium/exercises', { waitUntil: 'networkidle' });
  await expect(page.locator('h1')).toHaveText('Exercises');
  await expect(page.locator('table')).toBeVisible();
  await expect(page.locator('table')).toContainText('Bench Press');
  for (const item of items) await deleteExercise(request, item.id);
});
