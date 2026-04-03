import { expect, test } from '../../base-test';
import { createExercise, updateExercise, deleteExercise, toSlug } from '../../helpers';

test.describe('/compendium/exercises/:id/:slug/history', () => {
  test('renders history page with version entries', async ({ request, page }) => {
    const exercise = await createExercise(request, { names: ['Lat Pulldown'] });
    await updateExercise(request, exercise.id, { names: ['Lat Pulldown (v1)'] });
    await page.goto(
      `/compendium/exercises/${exercise.id}/${toSlug('Lat Pulldown')}/history`,
      { waitUntil: 'networkidle' },
    );
    await expect(page.locator('h1')).toContainText('History');
    const versionLabels = page.locator('text=/Version \\d+/');
    expect(await versionLabels.count()).toBeGreaterThanOrEqual(2);
    await deleteExercise(request, exercise.id);
  });

  test('shows history button on detail page after edits and navigates to history', async ({
    request,
    page,
  }) => {
    const exercise = await createExercise(request, { names: ['History Navigation Test'] });
    await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}/edit`, {
      waitUntil: 'networkidle',
    });
    const nameInput = page.locator('fieldset').first().locator('input').first();
    await nameInput.clear();
    await nameInput.fill('History Navigation Test (edited)');
    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/exercises/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(new RegExp(`/compendium/exercises/${exercise.id}/`));

    const historyLink = page.locator('a:has-text("History")');
    await expect(historyLink).toBeVisible();

    await historyLink.click();
    await page.waitForURL(/\/history$/);
    await expect(page.locator('h1')).toContainText('History');

    const versionLabels = page.locator('text=/Version \\d+/');
    expect(await versionLabels.count()).toBeGreaterThanOrEqual(2);

    await deleteExercise(request, exercise.id);
  });
});
