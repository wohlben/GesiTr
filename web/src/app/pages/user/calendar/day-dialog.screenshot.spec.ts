import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideRouter } from '@angular/router';
import { DayDialog } from './day-dialog';
import { WorkoutLog } from '$generated/user-models';

describe('DayDialog screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const logs = [
    { id: 1, name: 'Morning Push', status: 'finished' },
    { id: 2, name: 'Evening Pull', status: 'in_progress' },
    { id: 3, name: 'Recovery', status: 'aborted' },
  ] as WorkoutLog[];

  const renderDialog = () =>
    render(DayDialog, {
      inputs: {
        open: true,
        date: new Date(2024, 0, 15),
        logs,
      },
      providers: [provideRouter([])],
    });

  it('light', async () => {
    await renderDialog();
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('light', {
      allowedMismatchedPixelRatio: 0.02,
    } as never);
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    await renderDialog();
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('dark', { allowedMismatchedPixelRatio: 0.02 } as never);
  });
});
