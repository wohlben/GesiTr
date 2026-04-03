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
} from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantWorkouts: Record<string, { workoutName: string; exerciseName: string }[]> = {
  'desktop-light': [
    { workoutName: 'Push Day A', exerciseName: 'WL DL Bench Press' },
    { workoutName: 'Pull Day A', exerciseName: 'WL DL Barbell Row' },
  ],
  'desktop-dark': [
    { workoutName: 'Leg Day A', exerciseName: 'WL DD Back Squat' },
    { workoutName: 'Upper Body B', exerciseName: 'WL DD Overhead Press' },
  ],
  'mobile-light': [
    { workoutName: 'Full Body A', exerciseName: 'WL ML Deadlift' },
    { workoutName: 'Full Body B', exerciseName: 'WL ML Front Squat' },
  ],
  'mobile-dark': [
    { workoutName: 'Push Day B', exerciseName: 'WL MD Incline Press' },
    { workoutName: 'Pull Day B', exerciseName: 'WL MD Pull Up' },
  ],
};

test.describe('/compendium/workouts', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const variant = variantWorkouts[`${viewport.name}-light`];
        const items: {
          exercise: { id: number };
          scheme: { id: number };
          workout: { id: number };
          section: { id: number };
          sectionExercise: { id: number };
        }[] = [];

        for (const v of variant) {
          const exercise = await createExercise(request, { names: [v.exerciseName] });
          const scheme = await createExerciseScheme(request, {
            exerciseId: exercise.id,
          });
          const workout = await createWorkout(request, {
            name: v.workoutName,
          });
          const section = await createWorkoutSection(request, {
            workoutId: workout.id,
          });
          const sectionExercise = await createWorkoutSectionItem(request, {
            workoutSectionId: section.id,
            exerciseSchemeId: scheme.id,
          });
          items.push({ exercise, scheme, workout, section, sectionExercise });
        }

        await page.goto('/compendium/workouts', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Workouts');
        await expect(page.locator('table tbody tr')).toHaveCount(variant.length);
        await expect(page.locator('table')).toContainText(variant[0].workoutName);
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'compendium', 'workouts.png'], { fullPage: true });

        for (const item of items) {
          await deleteWorkoutSectionItem(request, item.sectionExercise.id);
          await deleteWorkoutSection(request, item.section.id);
          await deleteWorkout(request, item.workout.id);
          await deleteExerciseScheme(request, item.scheme.id);
          await deleteExercise(request, item.exercise.id);
        }
      });

      test('dark', async ({ request, page }) => {
        const variant = variantWorkouts[`${viewport.name}-dark`];
        const items: {
          exercise: { id: number };
          scheme: { id: number };
          workout: { id: number };
          section: { id: number };
          sectionExercise: { id: number };
        }[] = [];

        for (const v of variant) {
          const exercise = await createExercise(request, { names: [v.exerciseName] });
          const scheme = await createExerciseScheme(request, {
            exerciseId: exercise.id,
          });
          const workout = await createWorkout(request, {
            name: v.workoutName,
          });
          const section = await createWorkoutSection(request, {
            workoutId: workout.id,
          });
          const sectionExercise = await createWorkoutSectionItem(request, {
            workoutSectionId: section.id,
            exerciseSchemeId: scheme.id,
          });
          items.push({ exercise, scheme, workout, section, sectionExercise });
        }

        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/workouts', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Workouts');
        await expect(page.locator('table tbody tr')).toHaveCount(variant.length);
        await expect(page.locator('table')).toContainText(variant[0].workoutName);
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'compendium', 'workouts.png'], { fullPage: true });

        for (const item of items) {
          await deleteWorkoutSectionItem(request, item.sectionExercise.id);
          await deleteWorkoutSection(request, item.section.id);
          await deleteWorkout(request, item.workout.id);
          await deleteExerciseScheme(request, item.scheme.id);
          await deleteExercise(request, item.exercise.id);
        }
      });
    });
  }
});
