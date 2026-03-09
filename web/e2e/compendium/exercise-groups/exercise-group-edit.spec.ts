import { expect, test } from '@playwright/test';
import { createExerciseGroup, deleteExerciseGroup, toSlug } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

test.describe('/compendium/exercise-groups/:id/:slug/edit', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const group = await createExerciseGroup(request, { name: 'Plyometrics' });
        await page.goto(
          `/compendium/exercise-groups/${group.id}/${toSlug(group.name)}/edit`,
          { waitUntil: 'networkidle' },
        );
        await expect(page.locator('h1')).toHaveText('Edit Exercise Group');
        await expect(page).toHaveScreenshot(`${viewport.name}/light/compendium/exercise-groups/[id]/edit.png`);
        await deleteExerciseGroup(request, group.id);
      });

      test('dark', async ({ request, page }) => {
        const group = await createExerciseGroup(request, { name: 'Plyometrics' });
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(
          `/compendium/exercise-groups/${group.id}/${toSlug(group.name)}/edit`,
          { waitUntil: 'networkidle' },
        );
        await expect(page.locator('h1')).toHaveText('Edit Exercise Group');
        await expect(page).toHaveScreenshot(`${viewport.name}/dark/compendium/exercise-groups/[id]/edit.png`);
        await deleteExerciseGroup(request, group.id);
      });
    });
  }

  test('edits name and verifies detail and list views update', async ({ request, page }) => {
    const group = await createExerciseGroup(request, { name: 'Edit Test Group' });
    await page.goto(`/compendium/exercise-groups/${group.id}/${toSlug(group.name)}/edit`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Edit Exercise Group');

    const nameInput = page.locator('#name');
    const editedName = 'Edit Test Group (edited)';
    await nameInput.clear();
    await nameInput.fill(editedName);

    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/exercise-groups/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(new RegExp(`/compendium/exercise-groups/${group.id}/`));

    await expect(page.locator('h1')).toHaveText(editedName);

    await page.goto('/compendium/exercise-groups', { waitUntil: 'networkidle' });
    await expect(page.locator('table')).toContainText(editedName);

    await deleteExerciseGroup(request, group.id);
  });
});
