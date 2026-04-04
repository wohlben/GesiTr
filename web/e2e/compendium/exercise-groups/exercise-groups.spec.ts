import { expect, test } from '../../base-test';
import { createExerciseGroup, deleteExerciseGroup } from '../../helpers';

test('renders exercise groups list with expected content', async ({ request, page }) => {
  const groupNames = ['Core Exercises', 'Upper Body', 'Lower Body'];
  const items = [];
  for (const name of groupNames) {
    items.push(await createExerciseGroup(request, { name }));
  }
  await page.goto('/compendium/exercise-groups', { waitUntil: 'networkidle' });
  await expect(page.locator('h1')).toHaveText('Exercise Groups');
  await expect(page.locator('table')).toBeVisible();
  await expect(page.locator('table')).toContainText('Core Exercises');
  for (const item of items) await deleteExerciseGroup(request, item.id);
});
