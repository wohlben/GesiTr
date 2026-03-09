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

const variantWorkouts: Record<string, { workoutName: string; exerciseName: string; date: string }[]> = {
  'desktop-light': [
    { workoutName: 'Push Day A', exerciseName: 'WL DL Bench Press', date: '2026-01-15T00:00:00Z' },
    { workoutName: 'Pull Day A', exerciseName: 'WL DL Barbell Row', date: '2026-01-16T00:00:00Z' },
  ],
  'desktop-dark': [
    { workoutName: 'Leg Day A', exerciseName: 'WL DD Back Squat', date: '2026-01-17T00:00:00Z' },
    { workoutName: 'Upper Body B', exerciseName: 'WL DD Overhead Press', date: '2026-01-18T00:00:00Z' },
  ],
  'mobile-light': [
    { workoutName: 'Full Body A', exerciseName: 'WL ML Deadlift', date: '2026-01-19T00:00:00Z' },
    { workoutName: 'Full Body B', exerciseName: 'WL ML Front Squat', date: '2026-01-20T00:00:00Z' },
  ],
  'mobile-dark': [
    { workoutName: 'Push Day B', exerciseName: 'WL MD Incline Press', date: '2026-01-21T00:00:00Z' },
    { workoutName: 'Pull Day B', exerciseName: 'WL MD Pull Up', date: '2026-01-22T00:00:00Z' },
  ],
};

test.describe('/user/workouts', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const variant = variantWorkouts[`${viewport.name}-light`];
        const items: {
          exercise: { id: number; templateId: string };
          userExercise: { id: number };
          scheme: { id: number };
          workout: { id: number };
          section: { id: number };
          sectionExercise: { id: number };
        }[] = [];

        for (const v of variant) {
          const exercise = await createExercise(request, { name: v.exerciseName });
          const userExercise = await createUserExercise(request, exercise.templateId);
          const scheme = await createExerciseScheme(request, {
            userExerciseId: userExercise.id,
          });
          const workout = await createWorkout(request, {
            name: v.workoutName,
            date: v.date,
          });
          const section = await createWorkoutSection(request, {
            workoutId: workout.id,
          });
          const sectionExercise = await createWorkoutSectionExercise(request, {
            workoutSectionId: section.id,
            userExerciseSchemeId: scheme.id,
          });
          items.push({ exercise, userExercise, scheme, workout, section, sectionExercise });
        }

        await page.goto('/user/workouts', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('My Workouts');
        await expect(page.locator('table tbody tr')).toHaveCount(variant.length);
        await expect(page.locator('table')).toContainText(variant[0].workoutName);
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'user', 'workouts.png'], { fullPage: true });

        for (const item of items) {
          await deleteWorkoutSectionExercise(request, item.sectionExercise.id);
          await deleteWorkoutSection(request, item.section.id);
          await deleteWorkout(request, item.workout.id);
          await deleteExerciseScheme(request, item.scheme.id);
          await deleteUserExercise(request, item.userExercise.id);
          await deleteExercise(request, item.exercise.id);
        }
      });

      test('dark', async ({ request, page }) => {
        const variant = variantWorkouts[`${viewport.name}-dark`];
        const items: {
          exercise: { id: number; templateId: string };
          userExercise: { id: number };
          scheme: { id: number };
          workout: { id: number };
          section: { id: number };
          sectionExercise: { id: number };
        }[] = [];

        for (const v of variant) {
          const exercise = await createExercise(request, { name: v.exerciseName });
          const userExercise = await createUserExercise(request, exercise.templateId);
          const scheme = await createExerciseScheme(request, {
            userExerciseId: userExercise.id,
          });
          const workout = await createWorkout(request, {
            name: v.workoutName,
            date: v.date,
          });
          const section = await createWorkoutSection(request, {
            workoutId: workout.id,
          });
          const sectionExercise = await createWorkoutSectionExercise(request, {
            workoutSectionId: section.id,
            userExerciseSchemeId: scheme.id,
          });
          items.push({ exercise, userExercise, scheme, workout, section, sectionExercise });
        }

        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/user/workouts', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('My Workouts');
        await expect(page.locator('table tbody tr')).toHaveCount(variant.length);
        await expect(page.locator('table')).toContainText(variant[0].workoutName);
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'user', 'workouts.png'], { fullPage: true });

        for (const item of items) {
          await deleteWorkoutSectionExercise(request, item.sectionExercise.id);
          await deleteWorkoutSection(request, item.section.id);
          await deleteWorkout(request, item.workout.id);
          await deleteExerciseScheme(request, item.scheme.id);
          await deleteUserExercise(request, item.userExercise.id);
          await deleteExercise(request, item.exercise.id);
        }
      });
    });
  }
});
