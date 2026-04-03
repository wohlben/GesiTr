import { expect, test } from '../../base-test';
import { createEquipment, deleteEquipment } from '../../helpers';

test('renders equipment list with expected content', async ({ request, page }) => {
  const equipmentItems = [
    { name: 'dumbbells', displayName: 'Dumbbells (Pair)' },
    { name: 'barbell', displayName: 'Olympic Barbell' },
    { name: 'kettlebell', displayName: 'Kettlebell' },
  ];
  const items = [];
  for (const eq of equipmentItems) {
    items.push(await createEquipment(request, eq));
  }
  await page.goto('/compendium/equipment', { waitUntil: 'networkidle' });
  await expect(page.locator('h1')).toHaveText('Equipment');
  await expect(page.locator('table')).toBeVisible();
  await expect(page.locator('table')).toContainText('Dumbbells (Pair)');
  for (const item of items) await deleteEquipment(request, item.id);
});
