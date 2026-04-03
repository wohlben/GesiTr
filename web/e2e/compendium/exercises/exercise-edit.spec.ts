import { expect, test } from '../../base-test';
import { createExercise, deleteExercise, toSlug } from '../../helpers';

test.describe('/compendium/exercises/:id/:slug/edit', () => {
  test('renders edit form with exercise data', async ({ request, page }) => {
    const exercise = await createExercise(request, { names: ['Tricep Extension'] });
    await page.goto(
      `/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}/edit`,
      { waitUntil: 'networkidle' },
    );
    await expect(page.locator('h1')).toHaveText('Edit Exercise');
    await expect(page.locator('form')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
    await deleteExercise(request, exercise.id);
  });

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
