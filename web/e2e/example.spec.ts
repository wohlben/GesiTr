import { expect, test } from '@playwright/test';

test('navigates to exercises list', async ({ page }) => {
  await page.goto('/');

  // Should redirect to /compendium/exercises
  await expect(page).toHaveURL(/\/compendium\/exercises/);
  await expect(page.locator('h1')).toHaveText('Exercises');
});
