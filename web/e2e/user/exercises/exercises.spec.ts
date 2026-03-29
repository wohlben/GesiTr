import { expect, test } from '../../base-test';
import { createExercise, deleteExercise } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantExercises: Record<string, string[]> = {
  'desktop-light': ['My Bench Press', 'My Squat', 'My Deadlift'],
  'desktop-dark': ['My Dumbbell Fly', 'My Front Squat', 'My Romanian Deadlift'],
  'mobile-light': ['My Incline Press', 'My Goblet Squat', 'My Sumo Deadlift'],
  'mobile-dark': ['My Floor Press', 'My Split Squat', 'My Trap Bar Deadlift'],
};

test.describe('/compendium/exercises?mastery=me', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const names = variantExercises[`${viewport.name}-light`];
        const exercises: { id: number }[] = [];
        for (const name of names) {
          const exercise = await createExercise(request, { name });
          exercises.push(exercise);
        }
        await page.goto('/compendium/exercises?mastery=me', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Exercises');
        await expect(page.locator('table tbody tr')).toHaveCount(names.length);
        await expect(page.locator('table')).toContainText(names[0]);
        await expect(page).toHaveScreenshot([viewport.name, 'light', 'user', 'exercises.png'], { fullPage: true });
        for (const exercise of exercises) {
          await deleteExercise(request, exercise.id);
        }
      });

      test('dark', async ({ request, page }) => {
        const names = variantExercises[`${viewport.name}-dark`];
        const exercises: { id: number }[] = [];
        for (const name of names) {
          const exercise = await createExercise(request, { name });
          exercises.push(exercise);
        }
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto('/compendium/exercises?mastery=me', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Exercises');
        await expect(page.locator('table tbody tr')).toHaveCount(names.length);
        await expect(page.locator('table')).toContainText(names[0]);
        await expect(page).toHaveScreenshot([viewport.name, 'dark', 'user', 'exercises.png'], { fullPage: true });
        for (const exercise of exercises) {
          await deleteExercise(request, exercise.id);
        }
      });
    });
  }
});
