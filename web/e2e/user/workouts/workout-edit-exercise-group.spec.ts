import { expect, test } from '../../base-test';
import { createExercise, deleteExercise } from '../../helpers';

test.describe('workout edit - exercise group', () => {
  test('create workout with exercise group, reopen and verify', async ({ request, page }) => {
    // Setup: create two exercises via API
    const ex1 = await createExercise(request, { name: 'E2E Group Bench Press' });
    const ex2 = await createExercise(request, { name: 'E2E Group Dumbbell Press' });

    // 1. Navigate to create workout page
    await page.goto('/user/workouts/new', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('New Workout');

    // 2. Fill in workout name
    await page.locator('input#name').fill('E2E Group Test Workout');

    // 3. Add a section
    await page.getByText('+ Add Section').click();
    await expect(page.getByText('Section 1')).toBeVisible();

    // 4. Add an exercise group item
    await page.getByText('+ Add Exercise Group').click();
    await expect(page.getByText('Exercise 1')).toBeVisible();

    // 5. The group config should be visible with "New Group" selected
    // Fill in group name
    const groupConfig = page.locator('app-exercise-group-config');
    const nameInput = groupConfig.locator('input[placeholder="Group name (optional)"]');
    await nameInput.fill('My Push Group');

    // 6. Add first exercise to the group (the only native <select> is the member add)
    const memberSelect = groupConfig.locator('select');
    await memberSelect.selectOption({ label: 'E2E Group Bench Press' });

    // 7. Verify exercise appears in member list
    await expect(groupConfig.getByText('E2E Group Bench Press')).toBeVisible();

    // 8. Add second exercise
    await memberSelect.selectOption({ label: 'E2E Group Dumbbell Press' });

    // 9. Verify both exercises in member list
    await expect(groupConfig.getByText('E2E Group Bench Press')).toBeVisible();
    await expect(groupConfig.getByText('E2E Group Dumbbell Press')).toBeVisible();

    // 10. Save the workout
    const [saveResponse] = await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/user/workouts') && r.request().method() === 'POST',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    const workout = await saveResponse.json();

    // Should navigate to workout list
    await page.waitForURL(/\/user\/workouts$/);

    // 11. Reopen the workout in edit mode
    await page.goto(`/user/workouts/${workout.id}/edit`, { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('Edit Workout');
    await expect(page.locator('input#name')).toHaveValue('E2E Group Test Workout');

    // 12. Verify the exercise group data is loaded
    const reopenedConfig = page.locator('app-exercise-group-config');

    // Group name should be populated
    const reopenedNameInput = reopenedConfig.locator(
      'input[placeholder="Group name (optional)"]',
    );
    await expect(reopenedNameInput).toHaveValue('My Push Group', { timeout: 10000 });

    // The group selector (brn-select) should show the existing group, NOT "New Group"
    const groupTrigger = reopenedConfig.locator('hlm-select-value');
    await expect(groupTrigger).not.toHaveText('New Group', { timeout: 5000 });
    await expect(groupTrigger).toHaveText('My Push Group');

    // Both exercises should be visible in the member list
    await expect(reopenedConfig.getByText('E2E Group Bench Press')).toBeVisible();
    await expect(reopenedConfig.getByText('E2E Group Dumbbell Press')).toBeVisible();

    // The member add select should still be present
    const reopenedMemberSelect = reopenedConfig.locator('select');
    await expect(reopenedMemberSelect).toBeVisible();

    // Cleanup
    await deleteExercise(request, ex1.id);
    await deleteExercise(request, ex2.id);
  });
});
