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
  upsertSchemeSectionItem,
  fetchWorkoutLogs,
  deleteWorkoutLog,
} from '../../helpers';

test.describe('/compendium/workouts/[id]/start — break time editing', () => {
  test('editing break between exercises persists after page reload', async ({ request, page }) => {
    // Create fixtures: workout with 2 exercises in one section, default rest = 90s
    const exercise1 = await createExercise(request, { names: ['Break Test Ex A'] });
    const scheme1 = await createExerciseScheme(request, {
      exerciseId: exercise1.id,
      sets: 1,
      reps: 5,
      weight: 100,
    });

    const exercise2 = await createExercise(request, { names: ['Break Test Ex B'] });
    const scheme2 = await createExerciseScheme(request, {
      exerciseId: exercise2.id,
      sets: 1,
      reps: 8,
      weight: 50,
    });

    const workout = await createWorkout(request, { name: 'Break Edit Test' });
    const section = await createWorkoutSection(request, {
      workoutId: workout.id,
      label: 'Main',
      restBetweenExercises: 90,
    });
    const sectionEx1 = await createWorkoutSectionItem(request, {
      workoutSectionId: section.id,
      exerciseId: exercise1.id,
      position: 0,
    });
    await upsertSchemeSectionItem(request, {
      exerciseSchemeId: scheme1.id,
      workoutSectionItemId: sectionEx1.id,
    });
    const sectionEx2 = await createWorkoutSectionItem(request, {
      workoutSectionId: section.id,
      exerciseId: exercise2.id,
      position: 1,
    });
    await upsertSchemeSectionItem(request, {
      exerciseSchemeId: scheme2.id,
      workoutSectionItemId: sectionEx2.id,
    });

    // Navigate to workout start page, wait for exercises to load
    await page.goto(`/compendium/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
    await expect(page.getByText('Break Test Ex A')).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('Break Test Ex B')).toBeVisible({ timeout: 10000 });

    // Verify the default break time is 90s
    const breakInput = page.locator('input[data-field="breakAfterSeconds"]');
    await expect(breakInput).toHaveValue('90');

    // Change break time to 200s
    await breakInput.fill('200');
    await breakInput.dispatchEvent('change');

    // Wait for the API call to complete
    await page.waitForResponse(
      (resp) =>
        resp.url().includes('/api/user/workout-log-exercises/') &&
        resp.request().method() === 'PATCH',
    );

    // Reload the page completely
    await page.reload({ waitUntil: 'networkidle' });
    await expect(page.getByText('Break Test Ex A')).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('Break Test Ex B')).toBeVisible({ timeout: 10000 });

    // Verify the break time persisted as 200s
    const breakInputAfterReload = page.locator('input[data-field="breakAfterSeconds"]');
    await expect(breakInputAfterReload).toHaveValue('200');

    // Verify exercise order is preserved (A before B)
    const exA = page.getByText('Break Test Ex A');
    const exB = page.getByText('Break Test Ex B');
    const boxA = await exA.boundingBox();
    const boxB = await exB.boundingBox();
    expect(boxA).toBeTruthy();
    expect(boxB).toBeTruthy();
    expect(boxA!.y).toBeLessThan(boxB!.y);

    // Cleanup: delete planning log first, then template entities
    const logs = await fetchWorkoutLogs(request, {
      workoutId: workout.id,
      status: 'planning',
    });
    for (const log of logs) {
      await deleteWorkoutLog(request, log.id);
    }
    await deleteWorkoutSectionItem(request, sectionEx2.id);
    await deleteWorkoutSectionItem(request, sectionEx1.id);
    await deleteWorkoutSection(request, section.id);
    await deleteWorkout(request, workout.id);
    await deleteExerciseScheme(request, scheme2.id);
    await deleteExercise(request, exercise2.id);
    await deleteExerciseScheme(request, scheme1.id);
    await deleteExercise(request, exercise1.id);
  });
});
