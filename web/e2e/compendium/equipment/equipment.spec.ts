import { expect, test } from '@playwright/test';
import { createEquipment, deleteEquipment } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const equipmentItems = [
  { name: 'dumbbells', displayName: 'Dumbbells (Pair)' },
  { name: 'barbell', displayName: 'Olympic Barbell' },
  { name: 'kettlebell', displayName: 'Kettlebell' },
  { name: 'resistance-band', displayName: 'Resistance Band' },
  { name: 'pull-up-bar', displayName: 'Pull-Up Bar' },
];

for (const viewport of viewports) {
  test.describe(viewport.name, () => {
    test.use({ viewport: { width: viewport.width, height: viewport.height } });

    test('light', async ({ request, page }) => {
      const items = [];
      for (const eq of equipmentItems) {
        items.push(await createEquipment(request, eq));
      }
      await page.goto('/compendium/equipment', { waitUntil: 'networkidle' });
      await expect(page.locator('h1')).toHaveText('Equipment');
      await expect(page).toHaveScreenshot(`${viewport.name}-light.png`);
      for (const item of items) await deleteEquipment(request, item.id);
    });

    test('dark', async ({ request, page }) => {
      const items = [];
      for (const eq of equipmentItems) {
        items.push(await createEquipment(request, eq));
      }
      await page.emulateMedia({ colorScheme: 'dark' });
      await page.goto('/compendium/equipment', { waitUntil: 'networkidle' });
      await expect(page.locator('h1')).toHaveText('Equipment');
      await expect(page).toHaveScreenshot(`${viewport.name}-dark.png`);
      for (const item of items) await deleteEquipment(request, item.id);
    });
  });
}
