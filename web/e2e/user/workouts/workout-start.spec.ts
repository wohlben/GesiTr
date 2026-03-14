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

interface TestExercise {
  name: string;
  scheme: { sets: number; reps: number; weight: number; restBetweenSets?: number };
}

const variantData: Record<
  string,
  { workoutName: string; sectionLabel: string; exercises: TestExercise[] }
> = {
  'desktop-light': {
    workoutName: 'Start Push Day',
    sectionLabel: 'Main Lifts',
    exercises: [
      {
        name: 'WS DL Bench Press',
        scheme: { sets: 5, reps: 5, weight: 100, restBetweenSets: 180 },
      },
      {
        name: 'WS DL Overhead Press',
        scheme: { sets: 4, reps: 8, weight: 50, restBetweenSets: 120 },
      },
    ],
  },
  'desktop-dark': {
    workoutName: 'Start Pull Day',
    sectionLabel: 'Accessories',
    exercises: [
      { name: 'WS DD Barbell Row', scheme: { sets: 5, reps: 5, weight: 80, restBetweenSets: 180 } },
      { name: 'WS DD Bicep Curl', scheme: { sets: 3, reps: 12, weight: 20, restBetweenSets: 60 } },
    ],
  },
  'mobile-light': {
    workoutName: 'Start Leg Day',
    sectionLabel: 'Main Lifts',
    exercises: [
      { name: 'WS ML Back Squat', scheme: { sets: 5, reps: 5, weight: 140, restBetweenSets: 180 } },
      { name: 'WS ML Leg Press', scheme: { sets: 3, reps: 10, weight: 200, restBetweenSets: 90 } },
    ],
  },
  'mobile-dark': {
    workoutName: 'Start Full Body',
    sectionLabel: 'Compound',
    exercises: [
      { name: 'WS MD Deadlift', scheme: { sets: 3, reps: 5, weight: 160, restBetweenSets: 180 } },
      { name: 'WS MD Pull Up', scheme: { sets: 3, reps: 8, weight: 0, restBetweenSets: 120 } },
    ],
  },
};

test.describe('/user/workouts/[id]/start', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const v = variantData[`${viewport.name}-light`];
        const cleanup: (() => Promise<void>)[] = [];

        const workout = await createWorkout(request, { name: v.workoutName });
        cleanup.push(() => deleteWorkout(request, workout.id));

        const section = await createWorkoutSection(request, {
          workoutId: workout.id,
          label: v.sectionLabel,
          restBetweenExercises: 90,
        });
        cleanup.push(() => deleteWorkoutSection(request, section.id));

        for (let i = 0; i < v.exercises.length; i++) {
          const ex = v.exercises[i];
          const exercise = await createExercise(request, { name: ex.name });
          cleanup.push(() => deleteExercise(request, exercise.id));
          const userExercise = await createUserExercise(request, exercise.templateId);
          cleanup.push(() => deleteUserExercise(request, userExercise.id));
          const scheme = await createExerciseScheme(request, {
            userExerciseId: userExercise.id,
            ...ex.scheme,
          });
          cleanup.push(() => deleteExerciseScheme(request, scheme.id));
          const sectionExercise = await createWorkoutSectionExercise(request, {
            workoutSectionId: section.id,
            userExerciseSchemeId: scheme.id,
            position: i,
          });
          cleanup.push(() => deleteWorkoutSectionExercise(request, sectionExercise.id));
        }

        await page.goto(`/user/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Start Workout');
        await expect(page.locator('input#name')).toHaveValue(v.workoutName);
        await expect(page.getByText('Section 1')).toBeVisible();
        // Wait for all exercise names to load
        for (const ex of v.exercises) {
          await expect(page.getByText(ex.name)).toBeVisible({ timeout: 10000 });
        }
        await expect(page).toHaveScreenshot(
          [viewport.name, 'light', 'user', 'workouts', '[id]', 'start.png'],
          { fullPage: true },
        );

        for (const fn of cleanup.reverse()) {
          await fn();
        }
      });

      test('dark', async ({ request, page }) => {
        const v = variantData[`${viewport.name}-dark`];
        const cleanup: (() => Promise<void>)[] = [];

        const workout = await createWorkout(request, { name: v.workoutName });
        cleanup.push(() => deleteWorkout(request, workout.id));

        const section = await createWorkoutSection(request, {
          workoutId: workout.id,
          label: v.sectionLabel,
          type: 'supplementary',
          restBetweenExercises: 90,
        });
        cleanup.push(() => deleteWorkoutSection(request, section.id));

        for (let i = 0; i < v.exercises.length; i++) {
          const ex = v.exercises[i];
          const exercise = await createExercise(request, { name: ex.name });
          cleanup.push(() => deleteExercise(request, exercise.id));
          const userExercise = await createUserExercise(request, exercise.templateId);
          cleanup.push(() => deleteUserExercise(request, userExercise.id));
          const scheme = await createExerciseScheme(request, {
            userExerciseId: userExercise.id,
            ...ex.scheme,
          });
          cleanup.push(() => deleteExerciseScheme(request, scheme.id));
          const sectionExercise = await createWorkoutSectionExercise(request, {
            workoutSectionId: section.id,
            userExerciseSchemeId: scheme.id,
            position: i,
          });
          cleanup.push(() => deleteWorkoutSectionExercise(request, sectionExercise.id));
        }

        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(`/user/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Start Workout');
        await expect(page.locator('input#name')).toHaveValue(v.workoutName);
        await expect(page.getByText('Section 1')).toBeVisible();
        for (const ex of v.exercises) {
          await expect(page.getByText(ex.name)).toBeVisible({ timeout: 10000 });
        }
        await expect(page).toHaveScreenshot(
          [viewport.name, 'dark', 'user', 'workouts', '[id]', 'start.png'],
          { fullPage: true },
        );

        for (const fn of cleanup.reverse()) {
          await fn();
        }
      });
    });
  }
});
