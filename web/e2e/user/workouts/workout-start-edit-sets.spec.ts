import { expect, test } from '../../base-test';
import {
  createExercise,
  deleteExercise,
  createWorkout,
  deleteWorkout,
  createExerciseScheme,
  deleteExerciseScheme,
  createWorkoutSection,
  deleteWorkoutSection,
  createWorkoutSectionItem,
  deleteWorkoutSectionItem,
  fetchWorkoutLogs,
  deleteWorkoutLog,
} from '../../helpers';

test.describe('/compendium/workouts/[id]/start — set editing', () => {
  test('editing target reps and weight on a set persists after page reload', async ({
    request,
    page,
  }) => {
    // Create fixtures: 1 exercise with 2 sets
    const exercise = await createExercise(request, { name: 'Set Edit Test Ex' });
    const scheme = await createExerciseScheme(request, {
      exerciseId: exercise.id,
      sets: 2,
      reps: 10,
      weight: 60,
      restBetweenSets: 90,
    });

    const workout = await createWorkout(request, { name: 'Set Edit Test' });
    const section = await createWorkoutSection(request, {
      workoutId: workout.id,
      label: 'Main',
    });
    const sectionEx = await createWorkoutSectionItem(request, {
      workoutSectionId: section.id,
      exerciseSchemeId: scheme.id,
      position: 0,
    });

    // Navigate to workout start page, wait for exercise to load
    await page.goto(`/compendium/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
    await expect(page.getByText('Set Edit Test Ex')).toBeVisible({ timeout: 10000 });

    // Verify default values
    const repsInputs = page.locator('input[data-field="targetReps"]');
    const weightInputs = page.locator('input[data-field="targetWeight"]');
    await expect(repsInputs.first()).toHaveValue('10');
    await expect(weightInputs.first()).toHaveValue('60');

    // Change targetReps on first set to 8
    await repsInputs.first().fill('8');
    await repsInputs.first().dispatchEvent('change');

    // Wait for the PATCH call to complete
    await page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercise-sets/') &&
        resp.request().method() === 'PATCH',
    );

    // Change targetWeight on first set to 75
    await weightInputs.first().fill('75');
    await weightInputs.first().dispatchEvent('change');

    await page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercise-sets/') &&
        resp.request().method() === 'PATCH',
    );

    // Reload the page completely
    await page.reload({ waitUntil: 'networkidle' });
    await expect(page.getByText('Set Edit Test Ex')).toBeVisible({ timeout: 10000 });

    // Verify the values persisted
    const repsAfterReload = page.locator('input[data-field="targetReps"]');
    const weightAfterReload = page.locator('input[data-field="targetWeight"]');
    await expect(repsAfterReload.first()).toHaveValue('8');
    await expect(weightAfterReload.first()).toHaveValue('75');

    // Cleanup: delete planning log first, then template entities
    const logs = await fetchWorkoutLogs(request, {
      workoutId: workout.id,
      status: 'planning',
    });
    for (const log of logs) {
      await deleteWorkoutLog(request, log.id);
    }
    await deleteWorkoutSectionItem(request, sectionEx.id);
    await deleteWorkoutSection(request, section.id);
    await deleteWorkout(request, workout.id);
    await deleteExerciseScheme(request, scheme.id);
    await deleteExercise(request, exercise.id);
  });

  test('editing rest between sets persists after page reload', async ({ request, page }) => {
    // Create fixtures: 1 exercise with 2 sets (so there's 1 rest timer between them)
    const exercise = await createExercise(request, { name: 'Rest Edit Test Ex' });
    const scheme = await createExerciseScheme(request, {
      exerciseId: exercise.id,
      sets: 2,
      reps: 10,
      weight: 60,
      restBetweenSets: 90,
    });

    const workout = await createWorkout(request, { name: 'Rest Edit Test' });
    const section = await createWorkoutSection(request, {
      workoutId: workout.id,
      label: 'Main',
    });
    const sectionEx = await createWorkoutSectionItem(request, {
      workoutSectionId: section.id,
      exerciseSchemeId: scheme.id,
      position: 0,
    });

    // Navigate to workout start page, wait for exercise to load
    await page.goto(`/compendium/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
    await expect(page.getByText('Rest Edit Test Ex')).toBeVisible({ timeout: 10000 });

    // Verify default rest between sets is 90
    const restInput = page.locator('input[data-field="restAfterSeconds"]');
    await expect(restInput).toHaveValue('90');

    // Change rest to 120
    await restInput.fill('120');
    await restInput.dispatchEvent('change');

    // Wait for the PATCH call to complete
    await page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercise-sets/') &&
        resp.request().method() === 'PATCH',
    );

    // Reload the page completely
    await page.reload({ waitUntil: 'networkidle' });
    await expect(page.getByText('Rest Edit Test Ex')).toBeVisible({ timeout: 10000 });

    // Verify the rest time persisted as 120
    const restAfterReload = page.locator('input[data-field="restAfterSeconds"]');
    await expect(restAfterReload).toHaveValue('120');

    // Cleanup
    const logs = await fetchWorkoutLogs(request, {
      workoutId: workout.id,
      status: 'planning',
    });
    for (const log of logs) {
      await deleteWorkoutLog(request, log.id);
    }
    await deleteWorkoutSectionItem(request, sectionEx.id);
    await deleteWorkoutSection(request, section.id);
    await deleteWorkout(request, workout.id);
    await deleteExerciseScheme(request, scheme.id);
    await deleteExercise(request, exercise.id);
  });
});
