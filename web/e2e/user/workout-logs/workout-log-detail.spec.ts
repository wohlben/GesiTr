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
  deleteWorkoutLog,
} from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

interface TestExercise {
  name: string;
  scheme: { sets: number; reps: number; weight: number; restBetweenSets?: number };
}

interface Variant {
  workoutName: string;
  sectionLabel: string;
  exercises: TestExercise[];
}

const variantData: Record<string, Variant> = {
  'desktop-light': {
    workoutName: 'Log Push Day',
    sectionLabel: 'Main Lifts',
    exercises: [
      { name: 'WLD DL Bench Press', scheme: { sets: 4, reps: 6, weight: 90, restBetweenSets: 150 } },
      { name: 'WLD DL Incline DB Press', scheme: { sets: 3, reps: 10, weight: 30, restBetweenSets: 90 } },
    ],
  },
  'desktop-dark': {
    workoutName: 'Log Pull Day',
    sectionLabel: 'Main Lifts',
    exercises: [
      { name: 'WLD DD Barbell Row', scheme: { sets: 4, reps: 6, weight: 80, restBetweenSets: 150 } },
      { name: 'WLD DD Lat Pulldown', scheme: { sets: 3, reps: 12, weight: 50, restBetweenSets: 90 } },
    ],
  },
  'mobile-light': {
    workoutName: 'Log Leg Day',
    sectionLabel: 'Compounds',
    exercises: [
      { name: 'WLD ML Back Squat', scheme: { sets: 5, reps: 5, weight: 120, restBetweenSets: 180 } },
    ],
  },
  'mobile-dark': {
    workoutName: 'Log Upper Day',
    sectionLabel: 'Accessories',
    exercises: [
      { name: 'WLD MD Overhead Press', scheme: { sets: 4, reps: 8, weight: 45, restBetweenSets: 120 } },
    ],
  },
};

async function createFixturesAndStartLog(
  request: Parameters<Parameters<typeof test>[2]>[0]['request'],
  page: Parameters<Parameters<typeof test>[2]>[0]['page'],
  v: Variant,
) {
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

  // Navigate to start page, wait for planning log creation, click Start
  await page.goto(`/user/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
  await expect(page.locator('h1')).toHaveText('Plan Workout');
  for (const ex of v.exercises) {
    await expect(page.getByText(ex.name)).toBeVisible({ timeout: 10000 });
  }
  await page.getByRole('button', { name: 'Start Workout' }).click();
  await page.waitForURL(/\/user\/workout-logs\/\d+/, { timeout: 10000 });

  const logId = Number(page.url().match(/\/workout-logs\/(\d+)/)![1]);
  cleanup.push(() => deleteWorkoutLog(request, logId));

  return { logId, cleanup };
}

test.describe('/user/workout-logs/[id]', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const v = variantData[`${viewport.name}-light`];
        const { cleanup } = await createFixturesAndStartLog(request, page, v);

        // Verify the log detail page rendered
        await expect(page.locator('h1')).toHaveText(v.workoutName, { timeout: 10000 });
        // Wait for exercise names to resolve
        for (const ex of v.exercises) {
          await expect(page.getByText(ex.name)).toBeVisible({ timeout: 10000 });
        }
        await expect(page).toHaveScreenshot(
          [viewport.name, 'light', 'user', 'workout-logs', '[id].png'],
          { fullPage: true },
        );

        for (const fn of cleanup.reverse()) {
          await fn();
        }
      });

      test('dark', async ({ request, page }) => {
        const v = variantData[`${viewport.name}-dark`];
        await page.emulateMedia({ colorScheme: 'dark' });
        const { cleanup } = await createFixturesAndStartLog(request, page, v);

        await expect(page.locator('h1')).toHaveText(v.workoutName, { timeout: 10000 });
        for (const ex of v.exercises) {
          await expect(page.getByText(ex.name)).toBeVisible({ timeout: 10000 });
        }
        await expect(page).toHaveScreenshot(
          [viewport.name, 'dark', 'user', 'workout-logs', '[id].png'],
          { fullPage: true },
        );

        for (const fn of cleanup.reverse()) {
          await fn();
        }
      });
    });
  }
});
