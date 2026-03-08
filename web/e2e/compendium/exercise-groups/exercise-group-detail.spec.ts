import { expect, test } from '@playwright/test';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

test.describe('/compendium/exercise-groups/:id', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ page }) => {
        await page.goto('/compendium/exercise-groups/1/ab-wheel', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Exercise Group');
        await expect(page).toHaveScreenshot(`${viewport.name}-light.png`);
      });

      test('dark', async ({ page }) => {
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/exercise-groups/1/ab-wheel', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Exercise Group');
        await expect(page).toHaveScreenshot(`${viewport.name}-dark.png`);
      });
    });
  }
});
