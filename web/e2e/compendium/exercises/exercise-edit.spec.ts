import { expect, test } from '../../base-test';
import { createExercise, deleteExercise, toSlug } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantNames: Record<string, string> = {
  'desktop-light': 'Tricep Extension',
  'desktop-dark': 'Tricep Dip',
  'mobile-light': 'Skull Crusher',
  'mobile-dark': 'Cable Pushdown',
};

test.describe('/compendium/exercises/:id/:slug/edit', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const name = variantNames[`${viewport.name}-light`];
        const exercise = await createExercise(request, { names: [name] });
        await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}/edit`, {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toHaveText('Edit Exercise');
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'compendium', 'exercises', '[id]', 'edit.png'], { fullPage: true });
        await deleteExercise(request, exercise.id);
      });

      test('dark', async ({ request, page }) => {
        const name = variantNames[`${viewport.name}-dark`];
        const exercise = await createExercise(request, { names: [name] });
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}/edit`, {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toHaveText('Edit Exercise');
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'compendium', 'exercises', '[id]', 'edit.png'], { fullPage: true });
        await deleteExercise(request, exercise.id);
      });
    });
  }

  test('edits name and verifies detail and list views update', async ({ request, page }) => {
    const exercise = await createExercise(request, { names: ['Edit Test Exercise'] });
    await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}/edit`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Edit Exercise');

    const nameInput = page.locator('fieldset').first().locator('input').first();
    const editedName = 'Edit Test Exercise (edited)';
    await nameInput.clear();
    await nameInput.fill(editedName);

    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/exercises/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(new RegExp(`/compendium/exercises/${exercise.id}/`));

    await expect(page.locator('h1')).toHaveText(editedName);

    await page.goto('/compendium/exercises', { waitUntil: 'networkidle' });
    await expect(page.locator('table')).toContainText(editedName);

    await deleteExercise(request, exercise.id);
  });
});
