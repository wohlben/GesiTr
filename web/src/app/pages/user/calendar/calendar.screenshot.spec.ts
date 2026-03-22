import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutLogKeys } from '$core/query-keys';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { Calendar } from './calendar';
import { WorkoutLog } from '$generated/user-models';

describe('Calendar screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const mockLogs = [
    { id: 1, name: 'Push Day', status: 'finished', date: '2024-01-15T10:00:00Z' },
    { id: 2, name: 'Pull Day', status: 'in_progress', date: '2024-01-15T14:00:00Z' },
    { id: 3, name: 'Leg Day', status: 'aborted', date: '2024-01-20T10:00:00Z' },
  ] as WorkoutLog[];

  async function renderCalendar() {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });
    queryClient.setQueryData(workoutLogKeys.list(), mockLogs);

    const result = await render(Calendar, {
      providers: [
        provideTranslocoForTest(),
        provideTanStackQuery(queryClient),
        {
          provide: UserApiClient,
          useValue: { fetchWorkoutLogs: vi.fn().mockResolvedValue(mockLogs) },
        },
      ],
    });

    result.fixture.componentInstance.currentMonth.set(new Date(2024, 0, 1));
    result.fixture.detectChanges();
    await result.fixture.whenStable();
    return result;
  }

  it('light', async () => {
    const { fixture } = await renderCalendar();
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    const { fixture } = await renderCalendar();
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
  });
});
