import { expect, test } from '../../base-test';
import { createEquipment, deleteEquipment } from '../../helpers';

test.describe('/compendium/equipment/:id', () => {
  test('renders detail page with equipment data', async ({ request, page }) => {
    const equipment = await createEquipment(request, {
      name: 'medicine-ball',
      displayName: 'Medicine Ball',
      description: 'A weighted ball for core and strength training',
      category: 'free_weights',
    });
    await page.goto(`/compendium/equipment/${equipment.id}/${equipment.name}`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Medicine Ball');
    await expect(page.getByText('A weighted ball for core and strength training')).toBeVisible();
    await deleteEquipment(request, equipment.id);
  });

  test('delete dialog cancel closes the dialog', async ({ request, page }) => {
    const equipment = await createEquipment(request, {
      name: 'cancel-delete-equipment',
      displayName: 'Cancel Delete Test',
    });
    await page.goto(`/compendium/equipment/${equipment.id}/${equipment.name}`, {
      waitUntil: 'networkidle',
    });

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();

    await page.locator('[role="dialog"] button:has-text("Cancel")').click();
    await expect(page.locator('[role="dialog"]')).not.toBeVisible();

    await expect(page.locator('h1')).toHaveText('Cancel Delete Test');
    await deleteEquipment(request, equipment.id);
  });

  test('delete dialog confirm deletes and navigates to list', async ({ request, page }) => {
    const equipment = await createEquipment(request, {
      name: 'confirm-delete-equipment',
      displayName: 'Confirm Delete Test',
    });
    await page.goto(`/compendium/equipment/${equipment.id}/${equipment.name}`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Confirm Delete Test');

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();
    await expect(page.locator('[role="dialog"]')).toContainText('Confirm Delete Test');

    await Promise.all([
      page.waitForResponse(
        (r) =>
          r.url().includes(`/api/equipment/${equipment.id}`) && r.request().method() === 'DELETE',
      ),
      page.locator('[role="dialog"] button:has-text("Delete")').click(),
    ]);

    await page.waitForURL(/\/compendium\/equipment$/);
    await expect(page.locator('h1')).toHaveText('Equipment');
  });
});
