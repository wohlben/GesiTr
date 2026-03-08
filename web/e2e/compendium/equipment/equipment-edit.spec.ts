import { expect, test } from '@playwright/test';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

test.describe('/compendium/equipment/:id/:slug/edit', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ page }) => {
        await page.goto('/compendium/equipment/1/dumbbells-pair/edit', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Edit Equipment');
        await expect(page).toHaveScreenshot(`${viewport.name}-light.png`);
      });

      test('dark', async ({ page }) => {
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/equipment/1/dumbbells-pair/edit', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Edit Equipment');
        await expect(page).toHaveScreenshot(`${viewport.name}-dark.png`);
      });
    });
  }

  test('edits display name and verifies detail and list views update', async ({ page }) => {
    await page.goto('/compendium/equipment/1/dumbbells-pair/edit', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('Edit Equipment');

    const displayNameInput = page.locator('#displayName');
    const originalDisplayName = await displayNameInput.inputValue();
    const editedDisplayName = `${originalDisplayName} (edited)`;
    await displayNameInput.clear();
    await displayNameInput.fill(editedDisplayName);

    // Submit and wait for the PUT to complete, then navigation to detail page
    await Promise.all([
      page.waitForResponse((r) => r.url().includes('/api/equipment/') && r.request().method() === 'PUT'),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(/\/compendium\/equipment\/1\//);

    // Verify detail page shows updated display name in the header
    await expect(page.locator('h1')).toHaveText(editedDisplayName);

    // Navigate to list page and verify updated name appears
    await page.goto('/compendium/equipment', { waitUntil: 'networkidle' });
    await expect(page.locator('table')).toContainText(editedDisplayName);

    // Restore original display name
    await page.goto('/compendium/equipment/1/dumbbells-pair/edit', { waitUntil: 'networkidle' });
    await page.locator('#displayName').clear();
    await page.locator('#displayName').fill(originalDisplayName);
    await Promise.all([
      page.waitForResponse((r) => r.url().includes('/api/equipment/') && r.request().method() === 'PUT'),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(/\/compendium\/equipment\/1\//);
    await expect(page.locator('h1')).toHaveText(originalDisplayName);
  });
});
