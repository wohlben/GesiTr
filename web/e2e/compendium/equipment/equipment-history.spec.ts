import { expect, test, Page } from '@playwright/test';
import { createEquipment, updateEquipment, deleteEquipment } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

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

test.describe('/compendium/equipment/:id/:slug/history', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const equipment = await createEquipment(request, {
          name: 'foam-roller',
          displayName: 'Foam Roller',
        });
        await updateEquipment(request, equipment.id, { displayName: 'Foam Roller (v1)' });
        await page.goto(
          `/compendium/equipment/${equipment.id}/${equipment.name}/history`,
          { waitUntil: 'networkidle' },
        );
        await expect(page.locator('h1')).toContainText('History');
        await freezeDynamicContent(page);
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'compendium', 'equipment', '[id]', 'history.png']);
        await deleteEquipment(request, equipment.id);
      });

      test('dark', async ({ request, page }) => {
        const equipment = await createEquipment(request, {
          name: 'foam-roller',
          displayName: 'Foam Roller',
        });
        await updateEquipment(request, equipment.id, { displayName: 'Foam Roller (v1)' });
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(
          `/compendium/equipment/${equipment.id}/${equipment.name}/history`,
          { waitUntil: 'networkidle' },
        );
        await expect(page.locator('h1')).toContainText('History');
        await freezeDynamicContent(page);
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'compendium', 'equipment', '[id]', 'history.png']);
        await deleteEquipment(request, equipment.id);
      });
    });
  }

  test('shows history button on detail page after edits and navigates to history', async ({
    request,
    page,
  }) => {
    const equipment = await createEquipment(request, {
      name: 'history-nav-equipment',
      displayName: 'History Navigation Equipment',
    });
    await page.goto(`/compendium/equipment/${equipment.id}/${equipment.name}/edit`, {
      waitUntil: 'networkidle',
    });
    const displayNameInput = page.locator('#displayName');
    await displayNameInput.clear();
    await displayNameInput.fill('History Navigation Equipment (edited)');
    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/equipment/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(new RegExp(`/compendium/equipment/${equipment.id}/`));

    const historyLink = page.locator('a:has-text("History")');
    await expect(historyLink).toBeVisible();

    await historyLink.click();
    await page.waitForURL(/\/history$/);
    await expect(page.locator('h1')).toContainText('History');

    const versionLabels = page.locator('text=/Version \\d+/');
    expect(await versionLabels.count()).toBeGreaterThanOrEqual(2);

    await deleteEquipment(request, equipment.id);
  });
});
