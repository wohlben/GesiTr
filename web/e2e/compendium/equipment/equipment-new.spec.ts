import { expect, test } from '@playwright/test';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

test.describe('/compendium/equipment/new', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ page }) => {
        await page.goto('/compendium/equipment/new', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('New Equipment');
        await expect(page).toHaveScreenshot(`${viewport.name}-light.png`);
      });

      test('dark', async ({ page }) => {
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/equipment/new', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('New Equipment');
        await expect(page).toHaveScreenshot(`${viewport.name}-dark.png`);
      });
    });
  }

  test('creates new equipment and navigates to detail page', async ({ page }) => {
    await page.goto('/compendium/equipment/new', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('New Equipment');

    const testName = `e2e-equipment-${Date.now()}`;
    const testDisplayName = `E2E Test Equipment ${Date.now()}`;
    await page.locator('#name').fill(testName);
    await page.locator('#displayName').fill(testDisplayName);
    await page.locator('#description').fill('Created by e2e test');

    // Submit and wait for POST response
    const [response] = await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/equipment') && r.request().method() === 'POST',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    const created = await response.json();

    // Should navigate to detail page
    await page.waitForURL(/\/compendium\/equipment\/\d+\//);
    await expect(page.locator('h1')).toHaveText(testDisplayName);

    // Clean up: delete the created equipment
    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();
    await Promise.all([
      page.waitForResponse(
        (r) =>
          r.url().includes(`/api/equipment/${created.id}`) && r.request().method() === 'DELETE',
      ),
      page.locator('[role="dialog"] button:has-text("Delete")').click(),
    ]);
    await page.waitForURL(/\/compendium\/equipment$/);
  });

  test('cancel navigates to list page', async ({ page }) => {
    await page.goto('/compendium/equipment/new', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('New Equipment');

    await page.locator('a:has-text("Cancel")').click();
    await page.waitForURL(/\/compendium\/equipment$/);
    await expect(page.locator('h1')).toHaveText('Equipment');
  });
});
