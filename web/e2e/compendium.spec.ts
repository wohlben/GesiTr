import { expect, test } from '@playwright/test';

const routes = [
  { path: '/compendium/exercises', name: 'exercises', heading: 'Exercises' },
  { path: '/compendium/equipment', name: 'equipment', heading: 'Equipment' },
  { path: '/compendium/exercise-groups', name: 'exercise-groups', heading: 'Exercise Groups' },
];

for (const route of routes) {
  test.describe(route.name, () => {
    test('light', async ({ page }) => {
      await page.goto(route.path, { waitUntil: 'networkidle' });
      await expect(page.locator('h1')).toHaveText(route.heading);
      await expect(page).toHaveScreenshot(`${route.name}-light.png`);
    });

    test('dark', async ({ page }) => {
      await page.emulateMedia({ colorScheme: 'dark' });
      await page.goto(route.path, { waitUntil: 'networkidle' });
      await expect(page.locator('h1')).toHaveText(route.heading);
      await expect(page).toHaveScreenshot(`${route.name}-dark.png`);
    });
  });
}
