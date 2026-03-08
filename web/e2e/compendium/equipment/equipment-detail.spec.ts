import { expect, test } from '@playwright/test';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

test.describe('/compendium/equipment/:id', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ page }) => {
        await page.goto('/compendium/equipment/1/dumbbells-pair', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Equipment');
        await expect(page).toHaveScreenshot(`${viewport.name}-light.png`);
      });

      test('dark', async ({ page }) => {
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/equipment/1/dumbbells-pair', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Equipment');
        await expect(page).toHaveScreenshot(`${viewport.name}-dark.png`);
      });
    });
  }

  test('delete dialog cancel closes the dialog', async ({ page }) => {
    await page.goto('/compendium/equipment/1/dumbbells-pair', { waitUntil: 'networkidle' });

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();

    await page.locator('[role="dialog"] button:has-text("Cancel")').click();
    await expect(page.locator('[role="dialog"]')).not.toBeVisible();

    // Still on detail page
    await expect(page.locator('h1')).not.toHaveText('Equipment');
  });

  test('delete dialog confirm deletes and navigates to list', async ({ request, page }) => {
    // Create a temporary equipment item to delete via API
    const createResponse = await request.post('/api/equipment', {
      data: {
        name: 'temp-delete-test',
        displayName: 'Temp Delete Test',
        description: '',
        category: 'FREE_WEIGHTS',
      },
    });
    const created = await createResponse.json();

    // Navigate to the created equipment's detail page
    await page.goto(`/compendium/equipment/${created.id}/temp-delete-test`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Temp Delete Test');

    // Open delete dialog and confirm
    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();
    await expect(page.locator('[role="dialog"]')).toContainText('Temp Delete Test');

    await Promise.all([
      page.waitForResponse(
        (r) =>
          r.url().includes(`/api/equipment/${created.id}`) && r.request().method() === 'DELETE',
      ),
      page.locator('[role="dialog"] button:has-text("Delete")').click(),
    ]);

    // Should navigate to list page
    await page.waitForURL(/\/compendium\/equipment$/);
    await expect(page.locator('h1')).toHaveText('Equipment');
  });
});
