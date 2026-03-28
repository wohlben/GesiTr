import { expect, test } from '../../base-test';
import {
  createWorkout,
  createWorkoutLog,
  createWorkoutSchedule,
  createSchedulePeriod,
  createScheduleCommitment,
  startWorkoutLog,
  abandonWorkoutLog,
} from '../../helpers';

const viewports = [
  { name: 'desktop', width: 1280, height: 720 },
  { name: 'mobile', width: 375, height: 667 },
];

const themes = ['light', 'dark'] as const;

function daysFromNow(days: number): string {
  const d = new Date();
  d.setDate(d.getDate() + days);
  d.setHours(0, 0, 0, 0);
  return d.toISOString();
}

test.describe('/user/calendar', () => {
  test('screenshots', async ({ request, page }) => {
    // Create workout logs (dots on the calendar)
    const logA = await createWorkoutLog(request, { name: 'Morning Push' });
    await startWorkoutLog(request, logA.id);
    const logB = await createWorkoutLog(request, { name: 'Evening Recovery' });
    const startedB = await startWorkoutLog(request, logB.id);
    await abandonWorkoutLog(request, startedB.id);

    // Create a schedule with an active period and a planned period (bars on the calendar)
    const workout = await createWorkout(request, { name: 'Back Day' });
    const schedule = await createWorkoutSchedule(request, {
      workoutId: workout.id,
      startDate: daysFromNow(-14),
    });

    // Active period: started a few days ago, ends in a few days
    await createSchedulePeriod(request, {
      scheduleId: schedule.id,
      periodStart: daysFromNow(-5),
      periodEnd: daysFromNow(3),
      type: 'fixed_date',
    });

    // Planned period: starts after the active one, with commitments (hollow dots)
    const plannedPeriod = await createSchedulePeriod(request, {
      scheduleId: schedule.id,
      periodStart: daysFromNow(4),
      periodEnd: daysFromNow(11),
      type: 'fixed_date',
    });

    // Add commitments on specific days within the planned period
    await createScheduleCommitment(request, {
      periodId: plannedPeriod.id,
      date: daysFromNow(5),
    });
    await createScheduleCommitment(request, {
      periodId: plannedPeriod.id,
      date: daysFromNow(8),
    });

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
