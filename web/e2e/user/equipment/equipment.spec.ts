import { expect, test } from '@playwright/test';
import {
  createEquipment,
  deleteEquipment,
  createUserEquipment,
  deleteUserEquipment,
} from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantEquipment: Record<string, { name: string; displayName: string }[]> = {
  'desktop-light': [
    { name: 'my-dumbbells', displayName: 'My Dumbbells' },
    { name: 'my-barbell', displayName: 'My Barbell' },
    { name: 'my-kettlebell', displayName: 'My Kettlebell' },
  ],
  'desktop-dark': [
    { name: 'my-bench', displayName: 'My Bench' },
    { name: 'my-cable-machine', displayName: 'My Cable Machine' },
    { name: 'my-pull-up-bar', displayName: 'My Pull-Up Bar' },
  ],
  'mobile-light': [
    { name: 'my-resistance-band', displayName: 'My Resistance Band' },
    { name: 'my-foam-roller', displayName: 'My Foam Roller' },
    { name: 'my-jump-rope', displayName: 'My Jump Rope' },
  ],
  'mobile-dark': [
    { name: 'my-medicine-ball', displayName: 'My Medicine Ball' },
    { name: 'my-yoga-mat', displayName: 'My Yoga Mat' },
    { name: 'my-ab-wheel', displayName: 'My Ab Wheel' },
  ],
};

test.describe('/user/equipment', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const eqList = variantEquipment[`${viewport.name}-light`];
        const items: { equipment: { id: number; templateId: string }; userEquipment: { id: number } }[] = [];
        for (const eq of eqList) {
          const equipment = await createEquipment(request, eq);
          const userEquipment = await createUserEquipment(request, equipment.templateId);
          items.push({ equipment, userEquipment });
        }
        await page.goto('/user/equipment', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('My Equipment');
        await expect(page.locator('table tbody tr')).toHaveCount(eqList.length);
        await expect(page.locator('table')).toContainText(eqList[0].displayName);
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'user', 'equipment.png']);
        for (const item of items) {
          await deleteUserEquipment(request, item.userEquipment.id);
          await deleteEquipment(request, item.equipment.id);
        }
      });

      test('dark', async ({ request, page }) => {
        const eqList = variantEquipment[`${viewport.name}-dark`];
        const items: { equipment: { id: number; templateId: string }; userEquipment: { id: number } }[] = [];
        for (const eq of eqList) {
          const equipment = await createEquipment(request, eq);
          const userEquipment = await createUserEquipment(request, equipment.templateId);
          items.push({ equipment, userEquipment });
        }
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/user/equipment', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('My Equipment');
        await expect(page.locator('table tbody tr')).toHaveCount(eqList.length);
        await expect(page.locator('table')).toContainText(eqList[0].displayName);
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'user', 'equipment.png']);
        for (const item of items) {
          await deleteUserEquipment(request, item.userEquipment.id);
          await deleteEquipment(request, item.equipment.id);
        }
      });
    });
  }
});
