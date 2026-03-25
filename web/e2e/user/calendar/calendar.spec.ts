import { expect, test } from '../../base-test';
import { createWorkoutLog, startWorkoutLog, abandonWorkoutLog } from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const themes = ['light', 'dark'] as const;

test.describe('/user/calendar', () => {
  test('screenshots', async ({ request, page }) => {
    // Create two logs: one in_progress, one aborted
    const logA = await createWorkoutLog(request, { name: 'Morning Push' });
    await startWorkoutLog(request, logA.id);
    const logB = await createWorkoutLog(request, { name: 'Evening Recovery' });
    const startedB = await startWorkoutLog(request, logB.id);
    await abandonWorkoutLog(request, startedB.id);

    for (const viewport of viewports) {
      for (const theme of themes) {
        await page.setViewportSize({ width: viewport.width, height: viewport.height });
        if (theme === 'dark') {
          await page.emulateMedia({ colorScheme: 'dark' });
        } else {
          await page.emulateMedia({ colorScheme: 'light' });
        }

        await page.goto('/user/calendar', { waitUntil: 'networkidle' });
        await expect(page.locator('h1')).toHaveText('Calendar');
        await expect(page).toHaveScreenshot(
          [viewport.name, theme, 'user', 'calendar.png'],
          { fullPage: true },
        );

        // Open day dialog
        const dayButton = page.locator('button:has(span.rounded-full)').first();
        await dayButton.click();
        await expect(page.getByRole('dialog')).toBeVisible();
        await expect(page.getByRole('dialog')).toContainText('Morning Push');
        await expect(page).toHaveScreenshot(
          [viewport.name, theme, 'user', 'calendar', 'day-dialog.png'],
          { fullPage: true },
        );

        // Close dialog before next iteration
        await page.keyboard.press('Escape');
        await expect(page.getByRole('dialog')).not.toBeVisible();
      }
    }
  });
});
