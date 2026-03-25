import { expect, test } from '../../base-test';
import {
  createExerciseAs,
  deleteExercise,
  deleteExerciseAs,
  toSlug,
} from '../../helpers';

test.describe('add compendium exercise to my exercises', () => {
  test('clicking "Add to Mine" copies a public exercise and navigates to user detail', async ({
    request,
    page,
  }) => {
    // Create a public exercise owned by a different user (bypasses Angular proxy)
    const exercise = await createExerciseAs('compendium-author', {
      name: 'Add To Mine Test Exercise',
      description: 'A public exercise from another user',
      primaryMuscles: ['CHEST'],
      public: true,
    });

    // Navigate to the compendium detail page (browser uses devuser via proxy)
    await page.goto(
      `/compendium/exercises/${exercise.id}/${toSlug(exercise.name)}`,
      { waitUntil: 'networkidle' },
    );
    await expect(page.locator('h1')).toHaveText('Add To Mine Test Exercise');

    // The "Add to Mine" button should be visible (not "Already Added")
    await expect(page.getByText('Add to My Exercises')).toBeVisible();

    // Click "Add to Mine" and wait for the POST and navigation
    await Promise.all([
      page.waitForResponse(
        (r) =>
          r.url().includes('/api/exercises') &&
          r.request().method() === 'POST' &&
          r.status() === 201,
      ),
      page.getByText('Add to My Exercises').click(),
    ]);

    // Should navigate to the user exercise detail page
    await page.waitForURL(/\/user\/exercises\/\d+$/);
    await expect(page.locator('h1')).toHaveText('Add To Mine Test Exercise');

    // Clean up: delete the user's copy (as devuser via proxy) and the original (as compendium-author)
    const userExerciseId = Number(
      page.url().match(/\/user\/exercises\/(\d+)/)?.[1],
    );
    await deleteExercise(request, userExerciseId);
    await deleteExerciseAs(exercise.id, 'compendium-author');
  });
});
