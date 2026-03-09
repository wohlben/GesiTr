import { expect, test } from '@playwright/test';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

test.describe('/compendium/exercises/new', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ page }) => {
        await page.goto('/compendium/exercises/new', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('New Exercise');
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'compendium', 'exercises', 'new.png']);
      });

      test('dark', async ({ page }) => {
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/exercises/new', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('New Exercise');
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'compendium', 'exercises', 'new.png']);
      });
    });
  }

  test('creates a new exercise and navigates to detail page', async ({ page }) => {
    await page.goto('/compendium/exercises/new', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('New Exercise');

    const testName = `E2E Test Exercise ${Date.now()}`;
    await page.locator('#name').fill(testName);
    await page.locator('#description').fill('Created by e2e test');

    // Submit and wait for POST response
    const [response] = await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/exercises') && r.request().method() === 'POST',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    const created = await response.json();

    // Should navigate to detail page
    await page.waitForURL(/\/compendium\/exercises\/\d+\//);
    await expect(page.locator('h1')).toHaveText(testName);

    // Clean up: delete the created exercise
    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();
    await Promise.all([
      page.waitForResponse(
        (r) =>
          r.url().includes(`/api/exercises/${created.id}`) && r.request().method() === 'DELETE',
      ),
      page.locator('[role="dialog"] button:has-text("Delete")').click(),
    ]);
    await page.waitForURL(/\/compendium\/exercises$/);
  });

  test('cancel navigates to list page', async ({ page }) => {
    await page.goto('/compendium/exercises/new', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('New Exercise');

    await page.locator('a:has-text("Cancel")').click();
    await page.waitForURL(/\/compendium\/exercises$/);
    await expect(page.locator('h1')).toHaveText('Exercises');
  });
});
