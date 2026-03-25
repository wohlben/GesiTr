import { expect, test } from '../../base-test';
import { createExercise, deleteExercise } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantExercises: Record<string, string> = {
  'desktop-light': 'User Detail Bench Press',
  'desktop-dark': 'User Detail Squat',
  'mobile-light': 'User Detail Deadlift',
  'mobile-dark': 'User Detail Overhead Press',
};

test.describe('/user/exercises/:id', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const variantKey = `${viewport.name}-light`;
        const exercise = await createExercise(request, {
          name: variantExercises[variantKey],
          description: 'A compound strength exercise',
          primaryMuscles: ['CHEST'],
        });
        await page.goto(`/user/exercises/${exercise.id}`, { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Exercise');
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'user', 'exercises', '[id].png'], { fullPage: true });
        await deleteExercise(request, exercise.id);
      });

      test('dark', async ({ request, page }) => {
        const variantKey = `${viewport.name}-dark`;
        const exercise = await createExercise(request, {
          name: variantExercises[variantKey],
          description: 'A compound strength exercise',
          primaryMuscles: ['CHEST'],
        });
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(`/user/exercises/${exercise.id}`, { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).not.toHaveText('Exercise');
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'user', 'exercises', '[id].png'], { fullPage: true });
        await deleteExercise(request, exercise.id);
      });
    });
  }

  test('delete dialog cancel closes the dialog', async ({ request, page }) => {
    const exercise = await createExercise(request, { name: 'User Cancel Delete Exercise' });
    await page.goto(`/user/exercises/${exercise.id}`, { waitUntil: 'networkidle' });

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();

    await page.locator('[role="dialog"] button:has-text("Cancel")').click();
    await expect(page.locator('[role="dialog"]')).not.toBeVisible();

    await expect(page.locator('h1')).not.toHaveText('Exercise');
    await deleteExercise(request, exercise.id);
  });

  test('delete dialog confirm deletes and navigates to list', async ({ request, page }) => {
    const exercise = await createExercise(request, { name: 'User Confirm Delete Exercise' });

    await page.goto(`/user/exercises/${exercise.id}`, { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).not.toHaveText('Exercise');

    await page.locator('button:has-text("Delete")').click();
    await expect(page.locator('[role="dialog"]')).toBeVisible();

    await Promise.all([
      page.waitForResponse(
        (r) =>
          r.url().includes(`/api/exercises/${exercise.id}`) &&
          r.request().method() === 'DELETE',
      ),
      page.locator('[role="dialog"] button:has-text("Delete")').click(),
    ]);

    await page.waitForURL(/\/user\/exercises$/);
    await expect(page.locator('h1')).toHaveText('My Exercises');
  });
});
