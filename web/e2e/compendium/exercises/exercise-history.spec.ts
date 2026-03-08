import { expect, test, Page } from '@playwright/test';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

// Replace dynamic timestamps and JSON snapshots with fixed values for deterministic screenshots
async function freezeDynamicContent(page: Page) {
  await page.evaluate(() => {
    document.querySelectorAll('pre').forEach((el) => {
      el.textContent = '{ "snapshot": "..." }';
    });
    for (const el of document.querySelectorAll('span')) {
      if (el.textContent?.includes(' by ')) {
        el.textContent = 'Jan 1, 2025, 12:00:00 AM by system';
      }
    }
  });
}

test.describe('/compendium/exercises/:id/:slug/history', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ page }) => {
        await page.goto('/compendium/exercises/1/3-4-sit-up/history', {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toContainText('History');
        await freezeDynamicContent(page);
        await expect(page).toHaveScreenshot(`${viewport.name}-light.png`);
      });

      test('dark', async ({ page }) => {
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/exercises/1/3-4-sit-up/history', {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toContainText('History');
        await freezeDynamicContent(page);
        await expect(page).toHaveScreenshot(`${viewport.name}-dark.png`);
      });
    });
  }

  test('shows history button on detail page after edits and navigates to history', async ({
    page,
  }) => {
    // Use exercise 10 to avoid conflicts with other tests that use exercise 1
    await page.goto('/compendium/exercises/10/alternate-hammer-curl/edit', {
      waitUntil: 'networkidle',
    });
    const nameInput = page.locator('#name');
    const originalName = await nameInput.inputValue();
    await nameInput.clear();
    await nameInput.fill(originalName + ' (history test)');
    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/exercises/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(/\/compendium\/exercises\/10\//);

    // Detail page should now show History button (>1 version entries)
    const historyLink = page.locator('a:has-text("History")');
    await expect(historyLink).toBeVisible();

    // Navigate to history page via the button
    await historyLink.click();
    await page.waitForURL(/\/history$/);
    await expect(page.locator('h1')).toContainText('History');

    // Verify at least 2 version entries are displayed
    const versionLabels = page.locator('text=/Version \\d+/');
    expect(await versionLabels.count()).toBeGreaterThanOrEqual(2);

    // Restore original name
    await page.goto('/compendium/exercises/10/alternate-hammer-curl/edit', {
      waitUntil: 'networkidle',
    });
    await page.locator('#name').clear();
    await page.locator('#name').fill(originalName);
    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/exercises/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(/\/compendium\/exercises\/10\//);
  });
});
