import { expect, test, Page } from '@playwright/test';
import { createExercise, updateExercise, deleteExercise, toSlug } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const variantNames: Record<string, string> = {
  'desktop-light': 'Lat Pulldown',
  'desktop-dark': 'Seated Cable Row',
  'mobile-light': 'Pull Up',
  'mobile-dark': 'Chin Up',
};

async function freezeDynamicContent(page: Page) {
  await page.evaluate(() => {
    document.querySelectorAll('pre').forEach((el) => {
      el.textContent = '{ "snapshot": "..." }';
    });
    for (const el of document.querySelectorAll('span')) {
      if (el.textContent?.includes(' by ')) {
        el.textContent = 'Jan 1, 2025, 12:00:00 AM by system';
      }
    }
  });
}

test.describe('/compendium/exercises/:id/:slug/history', () => {
  for (const viewport of viewports) {
    test.describe(viewport.name, () => {
      test.use({ viewport: { width: viewport.width, height: viewport.height } });

      test('light', async ({ request, page }) => {
        const name = variantNames[`${viewport.name}-light`];
        const exercise = await createExercise(request, { name });
        await updateExercise(request, exercise.id, { name: `${name} (v1)` });
        await page.goto(
          `/compendium/exercises/${exercise.id}/${toSlug(name)}/history`,
          { waitUntil: 'networkidle' },
        );
        await expect(page.locator('h1')).toContainText('History');
        await freezeDynamicContent(page);
        await expect(page).toHaveScreenshot(`${viewport.name}/light/compendium/exercises/[id]/history.png`);
        await deleteExercise(request, exercise.id);
      });

      test('dark', async ({ request, page }) => {
        const name = variantNames[`${viewport.name}-dark`];
        const exercise = await createExercise(request, { name });
        await updateExercise(request, exercise.id, { name: `${name} (v1)` });
        await page.emulateMedia({ colorScheme: 'dark' });
        await page.goto(
          `/compendium/exercises/${exercise.id}/${toSlug(name)}/history`,
          { waitUntil: 'networkidle' },
        );
        await expect(page.locator('h1')).toContainText('History');
        await freezeDynamicContent(page);
        await expect(page).toHaveScreenshot(`${viewport.name}/dark/compendium/exercises/[id]/history.png`);
        await deleteExercise(request, exercise.id);
      });
    });
  }

  test('shows history button on detail page after edits and navigates to history', async ({
    request,
    page,
  }) => {
    const exercise = await createExercise(request, { name: 'History Navigation Test' });
    await page.goto(`/compendium/exercises/${exercise.id}/${toSlug(exercise.name)}/edit`, {
      waitUntil: 'networkidle',
    });
    const nameInput = page.locator('#name');
    await nameInput.clear();
    await nameInput.fill('History Navigation Test (edited)');
    await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/exercises/') && r.request().method() === 'PUT',
      ),
      page.locator('button[type="submit"]').click(),
    ]);
    await page.waitForURL(new RegExp(`/compendium/exercises/${exercise.id}/`));

    const historyLink = page.locator('a:has-text("History")');
    await expect(historyLink).toBeVisible();

    await historyLink.click();
    await page.waitForURL(/\/history$/);
    await expect(page.locator('h1')).toContainText('History');

    const versionLabels = page.locator('text=/Version \\d+/');
    expect(await versionLabels.count()).toBeGreaterThanOrEqual(2);

    await deleteExercise(request, exercise.id);
  });
});
