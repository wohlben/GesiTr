import { expect, test } from '../../base-test';
import { createEquipment, updateEquipment, deleteEquipment } from '../../helpers';

test.describe('/compendium/equipment/:id/:slug/history', () => {
  test('renders history page with version entries', async ({ request, page }) => {
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
    const versionLabels = page.locator('text=/Version \\d+/');
    expect(await versionLabels.count()).toBeGreaterThanOrEqual(2);
    await deleteEquipment(request, equipment.id);
  });

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
