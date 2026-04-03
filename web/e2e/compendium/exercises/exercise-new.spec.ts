import { expect, test } from '../../base-test';

test.describe('/compendium/exercises/new', () => {
  test('renders page with expected content', async ({ page }) => {
    await page.goto('/compendium/exercises/new', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('New Exercise');
    await expect(page.locator('form')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });

  test('creates a new exercise and navigates to detail page', async ({ page }) => {
    await page.goto('/compendium/exercises/new', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('New Exercise');

    const testName = `E2E Test Exercise ${Date.now()}`;
    await page.locator('fieldset').first().locator('input').first().fill(testName);
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
