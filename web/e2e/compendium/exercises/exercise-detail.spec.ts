import { expect, test } from '@playwright/test';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

test.describe('/compendium/exercises/:id', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ page }) => {
        await page.goto('/compendium/exercises/1/3-4-sit-up', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Exercise');
        await expect(page).toHaveScreenshot(`${viewport.name}-light.png`);
      });

      test('dark', async ({ page }) => {
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/exercises/1/3-4-sit-up', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Exercise');
        await expect(page).toHaveScreenshot(`${viewport.name}-dark.png`);
      });
    });
  }

  test('delete dialog cancel closes the dialog', async ({ page }) => {
    await page.goto('/compendium/exercises/1/3-4-sit-up', { waitUntil: 'networkidle' });

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();

    await page.locator('[role="dialog"] button:has-text("Cancel")').click();
    await expect(page.locator('[role="dialog"]')).not.toBeVisible();

    // Still on detail page
    await expect(page.locator('h1')).not.toHaveText('Exercise');
  });

  test('delete dialog confirm deletes and navigates to list', async ({ request, page }) => {
    // Create a temporary exercise to delete via API
    const createResponse = await request.post('/api/exercises', {
      data: {
        name: 'Temp Delete Test',
        type: 'STRENGTH',
        technicalDifficulty: 'BEGINNER',
        bodyWeightScaling: 0,
        force: [],
        primaryMuscles: [],
        secondaryMuscles: [],
        suggestedMeasurementParadigms: [],
        description: '',
        instructions: [],
        images: [],
        alternativeNames: [],
        equipmentIds: [],
      },
    });
    const created = await createResponse.json();

    // Navigate to the created exercise's detail page
    await page.goto(`/compendium/exercises/${created.id}/temp-delete-test`, {
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
          r.url().includes(`/api/exercises/${created.id}`) && r.request().method() === 'DELETE',
      ),
      page.locator('[role="dialog"] button:has-text("Delete")').click(),
    ]);

    // Should navigate to list page
    await page.waitForURL(/\/compendium\/exercises$/);
    await expect(page.locator('h1')).toHaveText('Exercises');
  });
});
