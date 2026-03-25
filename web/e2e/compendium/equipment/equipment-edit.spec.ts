import { expect, test } from '../../base-test';
import { createEquipment, deleteEquipment } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

test.describe('/compendium/equipment/:id/:slug/edit', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const equipment = await createEquipment(request, {
          name: 'jump-rope',
          displayName: 'Jump Rope',
        });
        await page.goto(`/compendium/equipment/${equipment.id}/${equipment.name}/edit`, {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toHaveText('Edit Equipment');
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'compendium', 'equipment', '[id]', 'edit.png'], { fullPage: true });
        await deleteEquipment(request, equipment.id);
      });

      test('dark', async ({ request, page }) => {
        const equipment = await createEquipment(request, {
          name: 'jump-rope',
          displayName: 'Jump Rope',
        });
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(`/compendium/equipment/${equipment.id}/${equipment.name}/edit`, {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toHaveText('Edit Equipment');
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'compendium', 'equipment', '[id]', 'edit.png'], { fullPage: true });
        await deleteEquipment(request, equipment.id);
      });
    });
  }

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
