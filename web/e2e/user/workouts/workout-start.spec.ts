import { expect, test } from '@playwright/test';
import {
  createExercise,
  deleteExercise,
  createUserExercise,
  deleteUserExercise,
  createWorkout,
  deleteWorkout,
  createExerciseScheme,
  deleteExerciseScheme,
  createWorkoutSection,
  deleteWorkoutSection,
  createWorkoutSectionExercise,
  deleteWorkoutSectionExercise,
} from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantData: Record<string, { workoutName: string; exerciseName: string }> = {
  'desktop-light': { workoutName: 'Start Push Day', exerciseName: 'WS DL Bench Press' },
  'desktop-dark': { workoutName: 'Start Pull Day', exerciseName: 'WS DD Barbell Row' },
  'mobile-light': { workoutName: 'Start Leg Day', exerciseName: 'WS ML Back Squat' },
  'mobile-dark': { workoutName: 'Start Full Body', exerciseName: 'WS MD Deadlift' },
};

test.describe('/user/workouts/[id]/start', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const v = variantData[`${viewport.name}-light`];
        const exercise = await createExercise(request, { name: v.exerciseName });
        const userExercise = await createUserExercise(request, exercise.templateId);
        const scheme = await createExerciseScheme(request, {
          userExerciseId: userExercise.id,
        });
        const workout = await createWorkout(request, {
          name: v.workoutName,
        });
        const section = await createWorkoutSection(request, {
          workoutId: workout.id,
          label: 'Main Lifts',
        });
        const sectionExercise = await createWorkoutSectionExercise(request, {
          workoutSectionId: section.id,
          userExerciseSchemeId: scheme.id,
        });

        await page.goto(`/user/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Start Workout');
        await expect(page.locator('input#name')).toHaveValue(v.workoutName);
        // Verify sections and exercises render
        await expect(page.getByText('Section 1')).toBeVisible();
        await expect(page.locator('input[formcontrolname="label"]')).toHaveValue('Main Lifts');
        // Wait for exercise name to load from async scheme fetch
        await expect(page.getByText(v.exerciseName)).toBeVisible({ timeout: 10000 });
        await expect(page).toHaveScreenshot(
          [viewport.name, 'light', 'user', 'workouts', '[id]', 'start.png'],
          { fullPage: true },
        );

        await deleteWorkoutSectionExercise(request, sectionExercise.id);
        await deleteWorkoutSection(request, section.id);
        await deleteWorkout(request, workout.id);
        await deleteExerciseScheme(request, scheme.id);
        await deleteUserExercise(request, userExercise.id);
        await deleteExercise(request, exercise.id);
      });

      test('dark', async ({ request, page }) => {
        const v = variantData[`${viewport.name}-dark`];
        const exercise = await createExercise(request, { name: v.exerciseName });
        const userExercise = await createUserExercise(request, exercise.templateId);
        const scheme = await createExerciseScheme(request, {
          userExerciseId: userExercise.id,
        });
        const workout = await createWorkout(request, {
          name: v.workoutName,
        });
        const section = await createWorkoutSection(request, {
          workoutId: workout.id,
          label: 'Accessories',
          type: 'supplementary',
        });
        const sectionExercise = await createWorkoutSectionExercise(request, {
          workoutSectionId: section.id,
          userExerciseSchemeId: scheme.id,
        });

        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(`/user/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Start Workout');
        await expect(page.locator('input#name')).toHaveValue(v.workoutName);
        await expect(page.getByText('Section 1')).toBeVisible();
        // Wait for exercise name to load from async scheme fetch
        await expect(page.getByText(v.exerciseName)).toBeVisible({ timeout: 10000 });
        await expect(page).toHaveScreenshot(
          [viewport.name, 'dark', 'user', 'workouts', '[id]', 'start.png'],
          { fullPage: true },
        );

        await deleteWorkoutSectionExercise(request, sectionExercise.id);
        await deleteWorkoutSection(request, section.id);
        await deleteWorkout(request, workout.id);
        await deleteExerciseScheme(request, scheme.id);
        await deleteUserExercise(request, userExercise.id);
        await deleteExercise(request, exercise.id);
      });
    });
  }
});
