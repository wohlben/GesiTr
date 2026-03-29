import { expect, test } from '../../base-test';
import { createExerciseAs, deleteExerciseAs } from '../../helpers';

test.describe('Create workout with public exercises from another user', () => {
  test('exercise dropdown shows public exercises created by another user', async ({ page }) => {
    // User "alice" creates two public exercises via the API
    const exerciseA = await createExerciseAs('alice', {
      name: 'Alice Bench Press',
      public: true,
    });
    const exerciseB = await createExerciseAs('alice', {
      name: 'Alice Deadlift',
      public: true,
    });

    try {
      // Default user (devuser) navigates to create a new workout
      await page.goto('/user/workouts/new', { waitUntil: 'networkidle' });
      await expect(page.locator('h1')).toHaveText('New Workout');

      // Fill in the workout name
      await page.locator('input#name').fill('Cross-User Workout');

      // Add a section
      await page.getByRole('button', { name: /add section/i }).click();

      // Add an exercise item to the section
      await page.getByRole('button', { name: '+ Add Exercise', exact: true }).click();

      // Open the exercise picker dropdown (the one with placeholder "-- Select --")
      const exerciseTrigger = page.getByRole('combobox', { name: '-- Select --' });
      await exerciseTrigger.click();

      // Verify alice's public exercises appear in the dropdown
      const options = page.locator('hlm-option');
      await expect(options.filter({ hasText: 'Alice Bench Press' })).toBeVisible({ timeout: 5000 });
      await expect(options.filter({ hasText: 'Alice Deadlift' })).toBeVisible();

      // Select one of alice's exercises
      await options.filter({ hasText: 'Alice Bench Press' }).click();

      // Verify the exercise was selected (trigger should show the name)
      await expect(exerciseTrigger).toContainText('Alice Bench Press');
    } finally {
      await deleteExerciseAs(exerciseA.id, 'alice');
      await deleteExerciseAs(exerciseB.id, 'alice');
    }
  });
});
