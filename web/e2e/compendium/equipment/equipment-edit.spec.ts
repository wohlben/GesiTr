import { expect, test } from '../../base-test';
import { createEquipment, deleteEquipment } from '../../helpers';

test.describe('/compendium/equipment/:id/:slug/edit', () => {
  test('renders edit form with equipment data', async ({ request, page }) => {
    const equipment = await createEquipment(request, {
      name: 'jump-rope',
      displayName: 'Jump Rope',
    });
    await page.goto(`/compendium/equipment/${equipment.id}/${equipment.name}/edit`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Edit Equipment');
    await expect(page.locator('form')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
    await deleteEquipment(request, equipment.id);
  });

  test('edits display name and verifies detail and list views update', async ({
    request,
    page,
  }) => {
    const equipment = await createEquipment(request, {
      name: 'edit-test-equipment',
      displayName: 'Edit Test Equipment',
    });
    await page.goto(`/compendium/equipment/${equipment.id}/${equipment.name}/edit`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Edit Equipment');

    const displayNameInput = page.locator('#displayName');
    const editedDisplayName = 'Edit Test Equipment (edited)';
    await displayNameInput.clear();
    await displayNameInput.fill(editedDisplayName);

    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/equipment/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(new RegExp(`/compendium/equipment/${equipment.id}/`));

    await expect(page.locator('h1')).toHaveText(editedDisplayName);

    await page.goto('/compendium/equipment', { waitUntil: 'networkidle' });
    await expect(page.locator('table')).toContainText(editedDisplayName);

    await deleteEquipment(request, equipment.id);
  });
});
