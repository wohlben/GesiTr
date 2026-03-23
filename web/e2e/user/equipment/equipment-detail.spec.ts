import { expect, test } from '@playwright/test';
import { createEquipment, deleteEquipment } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantEquipment: Record<string, { name: string; displayName: string }> = {
  'desktop-light': { name: 'user-detail-dumbbells', displayName: 'User Detail Dumbbells' },
  'desktop-dark': { name: 'user-detail-barbell', displayName: 'User Detail Barbell' },
  'mobile-light': { name: 'user-detail-kettlebell', displayName: 'User Detail Kettlebell' },
  'mobile-dark': { name: 'user-detail-band', displayName: 'User Detail Resistance Band' },
};

test.describe('/user/equipment/:id', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const variantKey = `${viewport.name}-light`;
        const equipment = await createEquipment(request, variantEquipment[variantKey]);
        await page.goto(`/user/equipment/${equipment.id}`, { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Equipment');
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'user', 'equipment', '[id].png'], { fullPage: true });
        await deleteEquipment(request, equipment.id);
      });

      test('dark', async ({ request, page }) => {
        const variantKey = `${viewport.name}-dark`;
        const equipment = await createEquipment(request, variantEquipment[variantKey]);
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(`/user/equipment/${equipment.id}`, { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Equipment');
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'user', 'equipment', '[id].png'], { fullPage: true });
        await deleteEquipment(request, equipment.id);
      });
    });
  }

  test('delete dialog cancel closes the dialog', async ({ request, page }) => {
    const equipment = await createEquipment(request, {
      name: 'user-cancel-delete-equip',
      displayName: 'User Cancel Delete Equipment',
    });
    await page.goto(`/user/equipment/${equipment.id}`, { waitUntil: 'networkidle' });

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();

    await page.locator('[role="dialog"] button:has-text("Cancel")').click();
    await expect(page.locator('[role="dialog"]')).not.toBeVisible();

    await expect(page.locator('h1')).not.toHaveText('Equipment');
    await deleteEquipment(request, equipment.id);
  });

  test('delete dialog confirm deletes and navigates to list', async ({ request, page }) => {
    const equipment = await createEquipment(request, {
      name: 'user-confirm-delete-equip',
      displayName: 'User Confirm Delete Equipment',
    });

    await page.goto(`/user/equipment/${equipment.id}`, { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).not.toHaveText('Equipment');

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();

    await Promise.all([
      page.waitForResponse(
        (r) =>
          r.url().includes(`/api/equipment/${equipment.id}`) &&
          r.request().method() === 'DELETE',
      ),
      page.locator('[role="dialog"] button:has-text("Delete")').click(),
    ]);

    await page.waitForURL(/\/user\/equipment$/);
    await expect(page.locator('h1')).toHaveText('My Equipment');
  });
});
