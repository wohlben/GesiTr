import { expect, test } from '../../base-test';
import { createExerciseGroup, deleteExerciseGroup, toSlug } from '../../helpers';

test.describe('/compendium/exercise-groups/:id', () => {
  test('renders detail page with exercise group data', async ({ request, page }) => {
    const group = await createExerciseGroup(request, {
      name: 'Balance Training',
    });
    await page.goto(`/compendium/exercise-groups/${group.id}/${toSlug(group.name)}`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Balance Training');
    await deleteExerciseGroup(request, group.id);
  });

  test('delete dialog cancel closes the dialog', async ({ request, page }) => {
    const group = await createExerciseGroup(request, { name: 'Cancel Delete Test Group' });
    await page.goto(`/compendium/exercise-groups/${group.id}/${toSlug(group.name)}`, {
      waitUntil: 'networkidle',
    });

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();

    await page.locator('[role="dialog"] button:has-text("Cancel")').click();
    await expect(page.locator('[role="dialog"]')).not.toBeVisible();

    await expect(page.locator('h1')).toHaveText('Cancel Delete Test Group');
    await deleteExerciseGroup(request, group.id);
  });

  test('delete dialog confirm deletes and navigates to list', async ({ request, page }) => {
    const group = await createExerciseGroup(request, { name: 'Confirm Delete Test Group' });
    await page.goto(`/compendium/exercise-groups/${group.id}/${toSlug(group.name)}`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Confirm Delete Test Group');

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();
    await expect(page.locator('[role="dialog"]')).toContainText('Confirm Delete Test Group');

    await Promise.all([
      page.waitForResponse(
        (r) =>
          r.url().includes(`/api/exercise-groups/${group.id}`) && r.request().method() === 'DELETE',
      ),
      page.locator('[role="dialog"] button:has-text("Delete")').click(),
    ]);

    await page.waitForURL(/\/compendium\/exercise-groups$/);
    await expect(page.locator('h1')).toHaveText('Exercise Groups');
  });
});
