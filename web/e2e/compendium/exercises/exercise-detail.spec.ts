import { expect, test } from '../../base-test';
import { createExercise, deleteExercise, toSlug } from '../../helpers';

test.describe('/compendium/exercises/:id', () => {
  test('renders detail page with exercise data', async ({ request, page }) => {
    const exercise = await createExercise(request, {
      names: ['Bicep Curl'],
      description: 'An isolation exercise for the biceps',
      primaryMuscles: ['BICEPS'],
    });
    await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Bicep Curl');
    await expect(page.getByText('An isolation exercise for the biceps')).toBeVisible();
    await deleteExercise(request, exercise.id);
  });

  test('delete dialog cancel closes the dialog', async ({ request, page }) => {
    const exercise = await createExercise(request, { names: ['Cancel Delete Test Exercise'] });
    await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}`, {
      waitUntil: 'networkidle',
    });

    await page.getByRole('button', { name: 'Delete', exact: true }).click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();

    await page.locator('[role="dialog"] button:has-text("Cancel")').click();
    await expect(page.locator('[role="dialog"]')).not.toBeVisible();

    await expect(page.locator('h1')).toHaveText('Cancel Delete Test Exercise');
    await deleteExercise(request, exercise.id);
  });

  test('delete dialog confirm deletes and navigates to list', async ({ request, page }) => {
    const exercise = await createExercise(request, { names: ['Confirm Delete Test Exercise'] });
    await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}`, {
      waitUntil: 'networkidle',
    });
    await expect(page.locator('h1')).toHaveText('Confirm Delete Test Exercise');

    await page.getByRole('button', { name: 'Delete', exact: true }).click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();
    await expect(page.locator('[role="dialog"]')).toContainText('Confirm Delete Test Exercise');

    await Promise.all([
      page.waitForResponse(
        (r) =>
          r.url().includes(`/api/exercises/${exercise.id}`) && r.request().method() === 'DELETE',
      ),
      page.locator('[role="dialog"] button:has-text("Delete")').click(),
    ]);

    await page.waitForURL(/\/compendium\/exercises$/);
    await expect(page.locator('h1')).toHaveText('Exercises');
  });
});
