import { expect, test } from '@playwright/test';
import { createExercise, deleteExercise } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantExercises: Record<string, string[]> = {
  'desktop-light': ['Bench Press', 'Squat', 'Deadlift', 'Overhead Press', 'Barbell Row'],
  'desktop-dark': ['Dumbbell Fly', 'Front Squat', 'Romanian Deadlift', 'Push Press', 'Pendlay Row'],
  'mobile-light': [
    'Incline Press',
    'Goblet Squat',
    'Sumo Deadlift',
    'Arnold Press',
    'Cable Row',
  ],
  'mobile-dark': [
    'Floor Press',
    'Split Squat',
    'Trap Bar Deadlift',
    'Military Press',
    'T-Bar Row',
  ],
};

for (const viewport of viewports) {
  test.describe(viewport.name, () => {
    test.use({ viewport: { width: viewport.width, height: viewport.height } });

    test('light', async ({ request, page }) => {
      const names = variantExercises[`${viewport.name}-light`];
      const items = [];
      for (const name of names) {
        items.push(await createExercise(request, { name }));
      }
      await page.goto('/compendium/exercises', { waitUntil: 'networkidle' });
      await expect(page.locator('h1')).toHaveText('Exercises');
      await expect(page).toHaveScreenshot(`${viewport.name}-light.png`);
      for (const item of items) await deleteExercise(request, item.id);
    });

    test('dark', async ({ request, page }) => {
      const names = variantExercises[`${viewport.name}-dark`];
      const items = [];
      for (const name of names) {
        items.push(await createExercise(request, { name }));
      }
      await page.emulateMedia({ colorScheme: 'dark' });
      await page.goto('/compendium/exercises', { waitUntil: 'networkidle' });
      await expect(page.locator('h1')).toHaveText('Exercises');
      await expect(page).toHaveScreenshot(`${viewport.name}-dark.png`);
      for (const item of items) await deleteExercise(request, item.id);
    });
  });
}
