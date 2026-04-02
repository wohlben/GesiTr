import { expect, test } from '../../base-test';
import { createExercise, deleteExercise } from '../../helpers';

test.describe('workout start - exercise group flow', () => {
  test('create workout with group, start, pick exercise, verify in log', async ({
    request,
    page,
  }) => {
    // Setup: create two exercises via API
    const ex1 = await createExercise(request, { name: 'E2E Start Bench Press' });
    const ex2 = await createExercise(request, { name: 'E2E Start Incline Press' });

    // === Phase 1: Create workout with exercise group via UI ===

    await page.goto('/compendium/workouts/new', { waitUntil: 'networkidle' });
    await page.locator('input#name').fill('E2E Start Group Workout');

    // Add section
    await page.getByText('+ Add Section').click();
    await expect(page.getByText('Section 1')).toBeVisible();

    // Add exercise group item
    await page.getByText('+ Add Exercise Group').click();

    // Configure the group
    const groupConfig = page.locator('app-exercise-group-config');
    await groupConfig
      .locator('input[placeholder="Group name (optional)"]')
      .fill('Push Variations');

    // Add both exercises as members
    const memberSelect = groupConfig.locator('select');
    await memberSelect.selectOption({ label: 'E2E Start Bench Press' });
    await memberSelect.selectOption({ label: 'E2E Start Incline Press' });

    // Save the workout
    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/workouts') && r.request().method() === 'POST',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(/\/compendium\/workouts$/);

    // Find the workout we just created
    const workoutLink = page.getByText('E2E Start Group Workout');
    await expect(workoutLink).toBeVisible();

    // Navigate to the workout's start page via the workouts list
    // Find the row and click the start link
    const row = page.locator('table tbody tr', { hasText: 'E2E Start Group Workout' });
    const startLink = row.locator('a[href*="/start"]');
    await startLink.click();

    // === Phase 2: Start workout page — verify pending group ===

    await expect(page.locator('h1')).toHaveText('Plan Workout', { timeout: 10000 });

    // The pending group card should be visible
    const pendingGroup = page.locator('[data-testid="pending-group"]');
    await expect(pendingGroup).toBeVisible({ timeout: 10000 });
    await expect(pendingGroup).toContainText('Push Variations');

    // The group should show the exercise picker
    await expect(pendingGroup.getByText('Pick an exercise from this group')).toBeVisible();

    // Both exercises should be available in the picker
    const groupSelect = pendingGroup.locator('select');
    const options = groupSelect.locator('option');
    // 3 options: placeholder + 2 exercises
    await expect(options).toHaveCount(3);

    // Start button should be disabled (pending groups not resolved)
    await expect(page.locator('button[type="submit"]')).toBeDisabled();

    // === Phase 3: Pick an exercise from the group ===

    await groupSelect.selectOption({ label: 'E2E Start Bench Press' });

    // The add-exercise dialog should open with the exercise pre-selected
    await expect(page.locator('hlm-dialog-content')).toBeVisible({ timeout: 5000 });

    // The exercise should be pre-selected (the dropdown should show it)
    // Configure the scheme: measurement type is already REP_BASED by default
    // Set basic params
    const dialog = page.locator('hlm-dialog-content');
    // Wait for the exercise config to be ready, then click the Add button
    const addButton = dialog.getByRole('button', { name: 'Add' });
    await expect(addButton).toBeVisible();
    await addButton.click();

    // === Phase 4: Verify the pending group is resolved ===

    // Pending group should be gone
    await expect(pendingGroup).not.toBeVisible({ timeout: 5000 });

    // The exercise should now appear as a regular exercise card
    await expect(page.getByText('E2E Start Bench Press')).toBeVisible({ timeout: 10000 });

    // Start button should now be enabled
    await expect(page.locator('button[type="submit"]')).toBeEnabled();

    // === Phase 5: Start the workout ===

    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/start') && r.request().method() === 'POST',
      ),
      page.locator('button[type="submit"]').click(),
    ]);

    // Should navigate to workout log detail
    await page.waitForURL(/\/user\/workout-logs\/\d+$/);

    // === Phase 6: Verify in workout log detail ===

    // The exercise should appear as a regular exercise in the log
    await expect(page.getByText('E2E Start Bench Press').first()).toBeVisible({ timeout: 10000 });

    // Cleanup
    await deleteExercise(request, ex1.id);
    await deleteExercise(request, ex2.id);
  });
});
