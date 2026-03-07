import { expect, test } from '@playwright/test';

test('redirects to exercises', async ({ page }) => {
  await page.goto('/');
  await expect(page).toHaveURL(/\/compendium\/exercises/);
});
