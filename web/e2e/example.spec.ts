import { expect, test } from './base-test';

test('redirects to exercises', async ({ page }) => {
  await page.goto('/');
  await expect(page).toHaveURL(/\/compendium\/exercises/);
});
