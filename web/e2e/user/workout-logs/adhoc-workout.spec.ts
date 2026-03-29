import { expect, test } from '../../base-test';
import { createExercise, deleteExercise } from '../../helpers';

test.describe('/user/workout-logs — adhoc workflow', () => {
  test('full adhoc workout flow: start, add exercises, complete/skip sets, finish', async ({
    request,
    page,
  }) => {
    const cleanup: (() => Promise<void>)[] = [];

    // --- Fixtures: create exercises via API ---
    const exerciseA = await createExercise(request, { name: 'AH Bench Press' });
    cleanup.push(() => deleteExercise(request, exerciseA.id));

    const exerciseB = await createExercise(request, { name: 'AH Bicep Curl' });
    cleanup.push(() => deleteExercise(request, exerciseB.id));

    // ============================================================
    // Step 1: Start adhoc workout from workout list
    // ============================================================
    await page.goto('/user/workouts', { waitUntil: 'networkidle' });

    const adhocButton = page.getByRole('button', { name: 'Ad-hoc Workout' });
    await expect(adhocButton).toBeVisible();

    const adhocResponse = page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-logs/adhoc') &&
        resp.request().method() === 'POST',
    );
    await adhocButton.click();
    const response = await adhocResponse;
    const logData = await response.json();
    const logId = logData.id;
    expect(logId).toBeTruthy();

    await page.waitForURL(/\/user\/workout-logs\/\d+/, { timeout: 10000 });

    // ============================================================
    // Step 2: Verify empty adhoc state
    // ============================================================
    await expect(page.locator('h1')).toHaveText('Ad-hoc Workout', { timeout: 10000 });

    const finishBtn = page.getByRole('button', { name: 'Finish Workout' });
    await expect(finishBtn).toBeVisible();

    const abandonBtn = page.getByRole('button', { name: 'Abandon' });
    await expect(abandonBtn).toBeVisible();

    const addExerciseBtn = page.getByRole('button', { name: '+ Add Exercise' });
    await expect(addExerciseBtn).toBeVisible();

    // No exercises yet
    await expect(page.locator('app-workout-log-active-header')).toHaveCount(0);

    // ============================================================
    // Step 3: Open Add Exercise dialog
    // ============================================================
    await addExerciseBtn.click();

    const dialogTitle = page.locator('h3').filter({ hasText: 'Add Exercise' });
    await expect(dialogTitle).toBeVisible({ timeout: 5000 });

    // Add button should be disabled (no exercise selected)
    const addBtn = page.getByRole('button', { name: 'Add', exact: true });
    await expect(addBtn).toBeDisabled();

    // Placeholder for Phase 2
    await expect(page.getByText('Select an exercise above to plan sets')).toBeVisible();

    // ============================================================
    // Step 4: Select exercise — Phase 2 (runner) appears
    // ============================================================
    const comboboxInput = page.locator('app-exercise-config hlm-combobox-input input').first();
    await comboboxInput.fill('AH Bench');
    await page.locator('hlm-combobox-item').filter({ hasText: 'AH Bench Press' }).click();

    // Runner should appear with exercise name
    const runner = page.locator('app-exercise-runner');
    await expect(runner).toBeVisible({ timeout: 5000 });
    await expect(runner.getByText('AH Bench Press')).toBeVisible();

    // Add button should now be enabled
    await expect(addBtn).toBeEnabled();

    // Runner auto-rebuilds sets via internal effect on inputs.
    // Verify default 3 sets: REP_BASED = 2 inputs (reps + weight) per set = 6
    const runnerInputs = runner.locator('input[type="number"]');
    await expect(runnerInputs).toHaveCount(6, { timeout: 5000 });

    // ============================================================
    // Step 5: Change set count to 5
    // ============================================================
    const setsInput = page
      .locator('app-exercise-config label')
      .filter({ hasText: 'Sets' })
      .locator('input');
    await setsInput.fill('5');
    // 5 sets × 2 inputs = 10
    await expect(runnerInputs).toHaveCount(10, { timeout: 5000 });

    // ============================================================
    // Step 6: Change set count to 2
    // ============================================================
    await setsInput.fill('2');
    // 2 sets × 2 inputs = 4
    await expect(runnerInputs).toHaveCount(4, { timeout: 5000 });

    // ============================================================
    // Step 7: Fill in reps, weight, rest — verify rest pill
    // ============================================================
    const repsInput = page
      .locator('app-exercise-config label')
      .filter({ hasText: 'Reps' })
      .locator('input');
    await repsInput.fill('8');

    const weightInput = page
      .locator('app-exercise-config label')
      .filter({ hasText: 'Weight (kg)' })
      .locator('input');
    await weightInput.fill('60');

    const restInput = page
      .locator('app-exercise-config label')
      .filter({ hasText: 'Rest' })
      .locator('input');
    await restInput.fill('90');

    // Rest pill should appear between sets (1 pill for 2 sets)
    // Rest pill has a w-12 class input inside the runner
    const restPillInputs = runner.locator('input.w-12');
    await expect(restPillInputs).toHaveCount(1, { timeout: 5000 });
    await expect(restPillInputs).toHaveValue('90');

    // Now total inputs = 4 (reps+weight) + 1 (rest pill) = 5
    await expect(runnerInputs).toHaveCount(5);

    // ============================================================
    // Step 8: Add the exercise
    // ============================================================
    const schemeResponse = page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/exercise-schemes') &&
        resp.request().method() === 'POST',
    );
    const exerciseResponse = page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercises') &&
        resp.request().method() === 'POST',
    );
    await addBtn.click();
    await schemeResponse;
    await exerciseResponse;

    // Dialog should close
    await expect(dialogTitle).not.toBeVisible({ timeout: 5000 });

    // Exercise should appear in active view
    await expect(
      page.locator('app-workout-log-active-header').getByText('AH Bench Press'),
    ).toBeVisible({ timeout: 10000 });

    // Active set should be "Set 1 of 2"
    await expect(page.getByText('Set 1 of 2')).toBeVisible({ timeout: 10000 });

    // ============================================================
    // Step 9: Complete a set (click "Done")
    // ============================================================
    const doneBtn = page.getByRole('button', { name: 'Done' });
    await expect(doneBtn).toBeVisible();

    const doneResponse = page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercise-sets/') &&
        resp.request().method() === 'PATCH',
    );
    await doneBtn.click();
    await doneResponse;

    // Rest timer should appear — skip it
    const restTimerSkip = page
      .locator('app-workout-log-active-break')
      .getByRole('button', { name: 'Skip' });
    await expect(restTimerSkip).toBeVisible({ timeout: 3000 });
    await restTimerSkip.click();

    // Set 2 should now be active
    await expect(page.getByText('Set 2 of 2')).toBeVisible({ timeout: 10000 });

    // ============================================================
    // Step 10: Skip a set
    // ============================================================
    const skipSetBtn = page
      .locator('app-workout-log-active-set')
      .getByRole('button', { name: 'Skip' });
    await expect(skipSetBtn).toBeVisible();

    const skipResponse = page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercise-sets/') &&
        resp.request().method() === 'PATCH',
    );
    await skipSetBtn.click();
    await skipResponse;

    // All sets completed
    await expect(page.getByText('All sets completed!')).toBeVisible({ timeout: 10000 });

    // ============================================================
    // Step 11: Add second exercise
    // ============================================================
    await addExerciseBtn.click();
    await expect(dialogTitle).toBeVisible({ timeout: 5000 });

    // Select second exercise
    const comboboxInput2 = page.locator('app-exercise-config hlm-combobox-input input').first();
    await comboboxInput2.fill('AH Bicep');
    await page.locator('hlm-combobox-item').filter({ hasText: 'AH Bicep Curl' }).click();

    // Keep defaults (3 sets, 10 reps) — just add
    const addBtn2 = page.getByRole('button', { name: 'Add', exact: true });
    const schemeResponse2 = page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/exercise-schemes') &&
        resp.request().method() === 'POST',
    );
    const exerciseResponse2 = page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercises') &&
        resp.request().method() === 'POST',
    );
    await addBtn2.click();
    await schemeResponse2;
    await exerciseResponse2;

    // Dialog should close
    await expect(dialogTitle).not.toBeVisible({ timeout: 5000 });

    // ============================================================
    // Step 12: Verify both exercises visible
    // ============================================================
    await expect(page.getByText('AH Bench Press').first()).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('AH Bicep Curl').first()).toBeVisible({ timeout: 10000 });

    // "All sets completed!" should be gone (new exercise has in_progress sets)
    await expect(page.getByText('All sets completed!')).not.toBeVisible();

    // New exercise should have its first set active
    await expect(page.getByText('Set 1 of 3')).toBeVisible();

    // ============================================================
    // Step 13: Finish workout
    // ============================================================
    const finishResponse = page.waitForResponse(
      (resp) =>
        resp.url().includes('/finish') && resp.request().method() === 'POST',
    );
    await finishBtn.click();
    await finishResponse;

    // Finish and Abandon buttons should disappear
    await expect(finishBtn).not.toBeVisible({ timeout: 10000 });
    await expect(abandonBtn).not.toBeVisible();

    // Review mode should be shown
    await expect(page.locator('app-workout-log-review')).toBeVisible({ timeout: 10000 });

    // ============================================================
    // Cleanup
    // ============================================================
    for (const fn of cleanup.reverse()) {
      await fn();
    }
  });
});
