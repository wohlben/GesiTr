import { expect, test } from '@playwright/test';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

test.describe('/compendium/exercise-groups/:id/:slug/edit', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ page }) => {
        await page.goto('/compendium/exercise-groups/1/ab-wheel/edit', {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toHaveText('Edit Exercise Group');
        await expect(page).toHaveScreenshot(`${viewport.name}-light.png`);
      });

      test('dark', async ({ page }) => {
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/exercise-groups/1/ab-wheel/edit', {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toHaveText('Edit Exercise Group');
        await expect(page).toHaveScreenshot(`${viewport.name}-dark.png`);
      });
    });
  }

  test('edits name and verifies detail and list views update', async ({ page }) => {
    await page.goto('/compendium/exercise-groups/1/ab-wheel/edit', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('Edit Exercise Group');

    const nameInput = page.locator('#name');
    const originalName = await nameInput.inputValue();
    const editedName = `${originalName} (edited)`;
    await nameInput.clear();
    await nameInput.fill(editedName);

    // Submit and wait for the PUT to complete, then navigation to detail page
    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/exercise-groups/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(/\/compendium\/exercise-groups\/1\//);

    // Verify detail page shows updated name
    await expect(page.locator('h1')).toHaveText(editedName);

    // Navigate to list page and verify updated name appears
    await page.goto('/compendium/exercise-groups', { waitUntil: 'networkidle' });
    await expect(page.locator('table')).toContainText(editedName);

    // Restore original name
    await page.goto('/compendium/exercise-groups/1/ab-wheel/edit', { waitUntil: 'networkidle' });
    await page.locator('#name').clear();
    await page.locator('#name').fill(originalName);
    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/exercise-groups/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(/\/compendium\/exercise-groups\/1\//);
    await expect(page.locator('h1')).toHaveText(originalName);
  });
});
