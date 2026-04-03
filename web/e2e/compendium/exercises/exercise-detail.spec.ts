import { expect, test } from '../../base-test';
import { createExercise, deleteExercise, toSlug } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantNames: Record<string, string> = {
  'desktop-light': 'Bicep Curl',
  'desktop-dark': 'Hammer Curl',
  'mobile-light': 'Concentration Curl',
  'mobile-dark': 'Preacher Curl',
};

test.describe('/compendium/exercises/:id', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const name = variantNames[`${viewport.name}-light`];
        const exercise = await createExercise(request, {
          names: [name],
          description: 'An isolation exercise for the biceps',
          primaryMuscles: ['BICEPS'],
        });
        await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}`, {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toHaveText(name);
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'compendium', 'exercises', '[id].png'], { fullPage: true });
        await deleteExercise(request, exercise.id);
      });

      test('dark', async ({ request, page }) => {
        const name = variantNames[`${viewport.name}-dark`];
        const exercise = await createExercise(request, {
          names: [name],
          description: 'An isolation exercise for the biceps',
          primaryMuscles: ['BICEPS'],
        });
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.names[0].name)}`, {
          waitUntil: 'networkidle',
        });
        await expect(page.locator('h1')).toHaveText(name);
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'compendium', 'exercises', '[id].png'], { fullPage: true });
        await deleteExercise(request, exercise.id);
      });
    });
  }

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
