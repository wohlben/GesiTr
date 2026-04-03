import { expect, test } from '../../base-test';
import {
  createExercise,
  createExerciseScheme,
  createWorkout,
  createWorkoutSection,
  createWorkoutSectionItem,
  createApiContextAs,
} from '../../helpers';

test.describe('/compendium/workouts — group member full flow', () => {
  test('member sees readonly workout-start, starts workout, tracks sets, finishes', async ({
    browser,
    request,
    page,
  }) => {
    // ============================================================
    // API Setup: devuser creates workout with exercise + group
    // ============================================================
    const exercise = await createExercise(request, {
      names: ['E2E Group Bench Press'],
      public: true,
    });
    const scheme = await createExerciseScheme(request, {
      exerciseId: exercise.id,
      sets: 2,
      reps: 10,
      weight: 60,
      restBetweenSets: 90,
    });
    const workout = await createWorkout(request, { name: 'E2E Group Push Day' });
    const section = await createWorkoutSection(request, {
      workoutId: workout.id,
      label: 'Main',
    });
    const sectionItem = await createWorkoutSectionItem(request, {
      workoutSectionId: section.id,
      exerciseSchemeId: scheme.id,
      position: 0,
    });

    const bobApi = await createApiContextAs('bob');

    // Create workout group and invite bob
    const groupRes = await request.post('/api/user/workout-groups', {
      data: { name: 'E2E Test Group', workoutId: workout.id },
    });
    expect(groupRes.ok()).toBeTruthy();
    const group = await groupRes.json();

    const memberRes = await request.post('/api/user/workout-group-memberships', {
      data: { groupId: group.id, userId: 'bob', role: 'invited' },
    });
    expect(memberRes.ok(), `Failed to invite bob: ${await memberRes.text()}`).toBeTruthy();

    // ============================================================
    // API Setup: bob creates his exercise scheme + accepts invite
    // ============================================================

    // Bob creates his own exercise scheme for the workout item
    const bobSchemeRes = await bobApi.post('/api/user/exercise-schemes', {
      data: {
        exerciseId: exercise.id,
        workoutSectionItemId: sectionItem.id,
        measurementType: 'REP_BASED',
        sets: 2,
        reps: 8,
        weight: 50,
        restBetweenSets: 60,
      },
    });
    expect(
      bobSchemeRes.ok(),
      `Failed to create bob's scheme: ${await bobSchemeRes.text()}`,
    ).toBeTruthy();

    // Bob accepts the invite
    const acceptRes = await bobApi.post(`/api/workouts/${workout.id}/group/accept`);
    expect(acceptRes.ok(), `Failed to accept invite: ${await acceptRes.text()}`).toBeTruthy();
    await bobApi.dispose();

    // ============================================================
    // UI Flow: bob opens workout list
    // ============================================================
    const baseURL = process.env['PLAYWRIGHT_TEST_BASE_URL'] ?? 'http://localhost:4200';
    const bobContext = await browser.newContext({
      baseURL,
      extraHTTPHeaders: { 'X-User-Id': 'bob' },
      serviceWorkers: 'block',
    });
    const bobPage = await bobContext.newPage();

    await bobPage.goto('/compendium/workouts', { waitUntil: 'networkidle' });
    await expect(bobPage.getByText('E2E Group Push Day')).toBeVisible({ timeout: 10000 });

    // Verify membership badge is visible
    await expect(bobPage.getByText('Member', { exact: true })).toBeVisible();

    // ============================================================
    // UI Flow: bob clicks "Start Workout"
    // ============================================================
    const startLink = bobPage.locator(`a[href*="${workout.id}/start"]`);
    await expect(startLink).toBeVisible();
    await startLink.click();

    await bobPage.waitForURL(/\/compendium\/workouts\/\d+\/start/, { timeout: 10000 });

    // Wait for the exercise to load in the planning log
    await expect(bobPage.getByText('E2E Group Bench Press')).toBeVisible({ timeout: 15000 });

    // ============================================================
    // Verify readonly mode: structural controls hidden
    // ============================================================

    // No "Add Section" button
    await expect(bobPage.locator('button:has-text("user.workouts.addSection")')).toHaveCount(0);

    // No "Add Exercise" button
    await expect(bobPage.locator('button:has-text("user.workouts.addExercise")')).toHaveCount(0);

    // No drag handles visible
    await expect(bobPage.locator('[cdkDragHandle]')).toHaveCount(0);

    // No remove buttons
    await expect(bobPage.locator('button:has-text("common.remove")')).toHaveCount(0);

    // ============================================================
    // Verify set inputs ARE editable — bob's scheme values shown
    // ============================================================
    const repsInputs = bobPage.locator('input[data-field="targetReps"]');
    const weightInputs = bobPage.locator('input[data-field="targetWeight"]');

    // Bob's scheme has 2 sets, 8 reps, 50 weight
    await expect(repsInputs).toHaveCount(2, { timeout: 5000 });
    await expect(repsInputs.first()).toHaveValue('8');
    await expect(weightInputs.first()).toHaveValue('50');

    // Bob edits reps on first set to 12
    await repsInputs.first().fill('12');
    await repsInputs.first().dispatchEvent('change');

    await bobPage.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercise-sets/') &&
        resp.request().method() === 'PATCH',
    );

    // ============================================================
    // Start the workout
    // ============================================================
    const startBtn = bobPage.getByRole('button', { name: 'Start Workout' });
    await expect(startBtn).toBeVisible();

    const startResponsePromise = bobPage.waitForResponse(
      (resp) => resp.url().includes('/start') && resp.request().method() === 'POST',
    );
    await startBtn.click();
    const startResponse = await startResponsePromise;
    const startBody = await startResponse.json();

    // Verify the log transitioned to in_progress
    expect(startBody.status).toBe('in_progress');

    // Should navigate to workout log detail
    await bobPage.waitForURL(/\/user\/workout-logs\/\d+/, { timeout: 10000 });
    await bobPage.waitForLoadState('networkidle');

    // ============================================================
    // Track exercises: complete sets
    // ============================================================

    // The Abandon button should be visible for in_progress logs
    const abandonBtn = bobPage.getByRole('button', { name: 'Abandon' });
    await expect(abandonBtn).toBeVisible({ timeout: 10000 });

    // Active set should be "Set 1 of 2"
    await expect(bobPage.getByText('Set 1 of 2')).toBeVisible({ timeout: 10000 });

    // ============================================================
    // Complete set 1
    // ============================================================
    const doneBtn = bobPage.getByRole('button', { name: 'Done' });
    await expect(doneBtn).toBeVisible();

    const doneResponse1 = bobPage.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercise-sets/') &&
        resp.request().method() === 'PATCH',
    );
    await doneBtn.click();
    await doneResponse1;

    // Rest timer should appear — skip it
    const restTimerSkip = bobPage
      .locator('app-workout-log-active-break')
      .getByRole('button', { name: 'Skip' });
    await expect(restTimerSkip).toBeVisible({ timeout: 3000 });
    await restTimerSkip.click();

    // Set 2 should now be active
    await expect(bobPage.getByText('Set 2 of 2')).toBeVisible({ timeout: 10000 });

    // ============================================================
    // Complete set 2 (last set — no rest timer after)
    // ============================================================
    const doneBtn2 = bobPage.getByRole('button', { name: 'Done' });
    await expect(doneBtn2).toBeVisible({ timeout: 5000 });

    const doneResponse2 = bobPage.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercise-sets/') &&
        resp.request().method() === 'PATCH',
    );
    await doneBtn2.click();
    await doneResponse2;

    // ============================================================
    // Verify workout auto-finished (non-adhoc logs finish when all sets complete)
    // ============================================================

    // Abandon button should disappear (workout auto-finishes)
    await expect(abandonBtn).not.toBeVisible({ timeout: 10000 });

    // The finished status checkmark icon should be visible
    await expect(bobPage.getByText('E2E Group Bench Press')).toBeVisible();

    // ============================================================
    // Cleanup
    // ============================================================
    await bobContext.close();
  });
});
