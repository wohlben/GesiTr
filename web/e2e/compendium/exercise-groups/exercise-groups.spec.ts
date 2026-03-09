import { expect, test } from '@playwright/test';
import { createExerciseGroup, deleteExerciseGroup } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const groupNames = ['Core Exercises', 'Upper Body', 'Lower Body', 'Cardio', 'Flexibility'];

for (const viewport of viewports) {
  test.describe(viewport.name, () => {
    test.use({ viewport: { width: viewport.width, height: viewport.height } });

    test('light', async ({ request, page }) => {
      const items = [];
      for (const name of groupNames) {
        items.push(await createExerciseGroup(request, { name }));
      }
      await page.goto('/compendium/exercise-groups', { waitUntil: 'networkidle' });
      await expect(page.locator('h1')).toHaveText('Exercise Groups');
      await expect(page).toHaveScreenshot(`${viewport.name}/light/compendium/exercise-groups.png`);
      for (const item of items) await deleteExerciseGroup(request, item.id);
    });

    test('dark', async ({ request, page }) => {
      const items = [];
      for (const name of groupNames) {
        items.push(await createExerciseGroup(request, { name }));
      }
      await page.emulateMedia({ colorScheme: 'dark' });
      await page.goto('/compendium/exercise-groups', { waitUntil: 'networkidle' });
      await expect(page.locator('h1')).toHaveText('Exercise Groups');
      await expect(page).toHaveScreenshot(`${viewport.name}/dark/compendium/exercise-groups.png`);
      for (const item of items) await deleteExerciseGroup(request, item.id);
    });
  });
}
