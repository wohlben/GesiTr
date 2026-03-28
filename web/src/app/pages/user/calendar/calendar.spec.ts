import { TestBed } from '@angular/core/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { Calendar } from './calendar';
import { WorkoutLog } from '$generated/user-models';

function makeLog(overrides: Partial<WorkoutLog>): WorkoutLog {
  return {
    id: 1,
    createdAt: '',
    updatedAt: '',
    deletedAt: null,
    owner: 'test',
    name: 'Test',
    status: 'in_progress',
    sections: [],
    ...overrides,
  } as WorkoutLog;
}

function daysFromNowStr(days: number): string {
  const d = new Date();
  d.setDate(d.getDate() + days);
  return d.toISOString().substring(0, 10) + 'T00:00:00Z';
}

describe('Calendar', () => {
  function setup(apiOverrides: Partial<UserApiClient> = {}) {
    const defaultApi = {
      fetchWorkoutLogs: vi.fn().mockResolvedValue([]),
      fetchWorkoutSchedules: vi.fn().mockResolvedValue([]),
      fetchSchedulePeriods: vi.fn().mockResolvedValue([]),
      fetchWorkouts: vi.fn().mockResolvedValue([]),
      fetchScheduleCommitments: vi.fn().mockResolvedValue([]),
      ...apiOverrides,
    };
    TestBed.configureTestingModule({
      providers: [
        provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
        { provide: UserApiClient, useValue: defaultApi },
        provideTranslocoForTest(),
      ],
    });
    const fixture = TestBed.createComponent(Calendar);
    fixture.detectChanges();
    return fixture.componentInstance;
  }

  describe('calendarCells', () => {
    it('generates correct cells for month starting on Monday', () => {
      const c = setup();
      // Jan 1 2024 is Monday -> offset 0
      c.currentMonth.set(new Date(2024, 0, 1));
      const cells = c.calendarCells();
      expect(cells.length).toBe(31);
      expect(cells[0]).toEqual(expect.objectContaining({ day: 1 }));
      expect(cells[30]).toEqual(expect.objectContaining({ day: 31 }));
    });

    it('generates null offsets for month starting on Sunday', () => {
      const c = setup();
      // Sep 1 2024 is Sunday -> offset 6
      c.currentMonth.set(new Date(2024, 8, 1));
      const cells = c.calendarCells();
      expect(cells.length).toBe(36); // 6 nulls + 30 days
      expect(cells.slice(0, 6).every((cell) => cell === null)).toBe(true);
      expect(cells[6]).toEqual(expect.objectContaining({ day: 1 }));
    });

    it('handles February in leap year', () => {
      const c = setup();
      // Feb 2024 is leap year, Feb 1 is Thursday -> offset 3
      c.currentMonth.set(new Date(2024, 1, 1));
      const cells = c.calendarCells();
      const dayCells = cells.filter((cell) => cell !== null);
      expect(dayCells.length).toBe(29);
    });

    it('handles February in non-leap year', () => {
      const c = setup();
      // Feb 2023, Feb 1 is Wednesday -> offset 2
      c.currentMonth.set(new Date(2023, 1, 1));
      const cells = c.calendarCells();
      const dayCells = cells.filter((cell) => cell !== null);
      expect(dayCells.length).toBe(28);
    });

    it('marks today correctly', () => {
      const c = setup();
      const today = new Date();
      c.currentMonth.set(new Date(today.getFullYear(), today.getMonth(), 1));
      const cells = c.calendarCells();
      const todayCells = cells.filter((cell) => cell !== null && cell.isToday);
      expect(todayCells.length).toBe(1);
      expect(todayCells[0]!.day).toBe(today.getDate());
    });

    it('does not mark today in other months', () => {
      const c = setup();
      // Use a month far in the past
      c.currentMonth.set(new Date(2000, 0, 1));
      const cells = c.calendarCells();
      const todayCells = cells.filter((cell) => cell !== null && cell.isToday);
      expect(todayCells.length).toBe(0);
    });
  });

  describe('logsByDate - fixed_date commitments', () => {
    it('places proposed logs at their date, not dueStart', async () => {
      const c = setup({
        fetchWorkoutLogs: vi.fn().mockResolvedValue([
          makeLog({
            id: 1,
            status: 'proposed',
            date: daysFromNowStr(3),
            dueStart: daysFromNowStr(-2),
            dueEnd: daysFromNowStr(7),
            periodId: 1,
          }),
        ]),
      });
      // Wait for query to resolve
      await vi.waitFor(() => expect(c.logsQuery.data()).toBeTruthy());

      const today = new Date();
      c.currentMonth.set(new Date(today.getFullYear(), today.getMonth(), 1));
      const cells = c.calendarCells();
      const cellsWithLogs = cells.filter((cell) => cell !== null && cell.logs.length > 0);
      expect(cellsWithLogs.length).toBe(1);

      const expectedDay = new Date();
      expectedDay.setDate(expectedDay.getDate() + 3);
      expect(cellsWithLogs[0]!.day).toBe(expectedDay.getDate());
    });
  });

  describe('logsByDate - frequency commitments', () => {
    it('places frequency log on today instead of dueStart', async () => {
      const c = setup({
        fetchWorkoutLogs: vi.fn().mockResolvedValue([
          makeLog({
            id: 1,
            status: 'proposed',
            date: undefined, // frequency — no date
            dueStart: daysFromNowStr(-3),
            dueEnd: daysFromNowStr(5),
            periodId: 10,
          }),
        ]),
      });
      await vi.waitFor(() => expect(c.logsQuery.data()).toBeTruthy());

      const today = new Date();
      c.currentMonth.set(new Date(today.getFullYear(), today.getMonth(), 1));
      const cells = c.calendarCells();
      const cellsWithLogs = cells.filter((cell) => cell !== null && cell.logs.length > 0);
      expect(cellsWithLogs.length).toBe(1);
      expect(cellsWithLogs[0]!.day).toBe(today.getDate());
    });

    it('collapses multiple frequency logs per period into one', async () => {
      const c = setup({
        fetchWorkoutLogs: vi.fn().mockResolvedValue([
          makeLog({
            id: 1,
            status: 'proposed',
            date: undefined,
            dueStart: daysFromNowStr(-3),
            dueEnd: daysFromNowStr(5),
            periodId: 10,
          }),
          makeLog({
            id: 2,
            status: 'proposed',
            date: undefined,
            dueStart: daysFromNowStr(-3),
            dueEnd: daysFromNowStr(5),
            periodId: 10,
          }),
          makeLog({
            id: 3,
            status: 'committed',
            date: undefined,
            dueStart: daysFromNowStr(-3),
            dueEnd: daysFromNowStr(5),
            periodId: 10,
          }),
        ]),
      });
      await vi.waitFor(() => expect(c.logsQuery.data()).toBeTruthy());

      const today = new Date();
      c.currentMonth.set(new Date(today.getFullYear(), today.getMonth(), 1));
      const cells = c.calendarCells();
      const todayCell = cells.find((cell) => cell !== null && cell.isToday);
      // All 3 logs collapsed into 1 representative
      expect(todayCell!.logs.length).toBe(1);
    });

    it('shows frequency logs from different periods separately', async () => {
      const c = setup({
        fetchWorkoutLogs: vi.fn().mockResolvedValue([
          makeLog({
            id: 1,
            status: 'proposed',
            date: undefined,
            dueStart: daysFromNowStr(-3),
            dueEnd: daysFromNowStr(5),
            periodId: 10,
          }),
          makeLog({
            id: 2,
            status: 'proposed',
            date: undefined,
            dueStart: daysFromNowStr(-3),
            dueEnd: daysFromNowStr(5),
            periodId: 20,
          }),
        ]),
      });
      await vi.waitFor(() => expect(c.logsQuery.data()).toBeTruthy());

      const today = new Date();
      c.currentMonth.set(new Date(today.getFullYear(), today.getMonth(), 1));
      const cells = c.calendarCells();
      const todayCell = cells.find((cell) => cell !== null && cell.isToday);
      // 2 different periods → 2 dots
      expect(todayCell!.logs.length).toBe(2);
    });

    it('does not show frequency logs when period has ended', async () => {
      const c = setup({
        fetchWorkoutLogs: vi.fn().mockResolvedValue([
          makeLog({
            id: 1,
            status: 'proposed',
            date: undefined,
            dueStart: daysFromNowStr(-10),
            dueEnd: daysFromNowStr(-1), // ended yesterday
            periodId: 10,
          }),
        ]),
      });
      await vi.waitFor(() => expect(c.logsQuery.data()).toBeTruthy());

      const today = new Date();
      c.currentMonth.set(new Date(today.getFullYear(), today.getMonth(), 1));
      const cells = c.calendarCells();
      const cellsWithLogs = cells.filter((cell) => cell !== null && cell.logs.length > 0);
      expect(cellsWithLogs.length).toBe(0);
    });

    it('does not show terminal frequency logs (skipped/broken)', async () => {
      const c = setup({
        fetchWorkoutLogs: vi.fn().mockResolvedValue([
          makeLog({
            id: 1,
            status: 'skipped',
            date: undefined,
            dueStart: daysFromNowStr(-3),
            dueEnd: daysFromNowStr(5),
            periodId: 10,
          }),
          makeLog({
            id: 2,
            status: 'broken',
            date: undefined,
            dueStart: daysFromNowStr(-3),
            dueEnd: daysFromNowStr(5),
            periodId: 10,
          }),
        ]),
      });
      await vi.waitFor(() => expect(c.logsQuery.data()).toBeTruthy());

      const today = new Date();
      c.currentMonth.set(new Date(today.getFullYear(), today.getMonth(), 1));
      const cells = c.calendarCells();
      const cellsWithLogs = cells.filter((cell) => cell !== null && cell.logs.length > 0);
      expect(cellsWithLogs.length).toBe(0);
    });
  });

  describe('monthLabel', () => {
    it('formats month and year', () => {
      const c = setup();
      c.currentMonth.set(new Date(2024, 0, 1));
      expect(c.monthLabel()).toBe('January 2024');
    });

    it('formats December correctly', () => {
      const c = setup();
      c.currentMonth.set(new Date(2024, 11, 1));
      expect(c.monthLabel()).toBe('December 2024');
    });
  });

  describe('navigation', () => {
    it('prevMonth goes to previous month', () => {
      const c = setup();
      c.currentMonth.set(new Date(2024, 3, 1));
      c.prevMonth();
      expect(c.currentMonth().getMonth()).toBe(2);
    });

    it('nextMonth goes to next month', () => {
      const c = setup();
      c.currentMonth.set(new Date(2024, 3, 1));
      c.nextMonth();
      expect(c.currentMonth().getMonth()).toBe(4);
    });

    it('prevMonth wraps year boundary', () => {
      const c = setup();
      c.currentMonth.set(new Date(2024, 0, 1));
      c.prevMonth();
      expect(c.currentMonth().getMonth()).toBe(11);
      expect(c.currentMonth().getFullYear()).toBe(2023);
    });

    it('nextMonth wraps year boundary', () => {
      const c = setup();
      c.currentMonth.set(new Date(2024, 11, 1));
      c.nextMonth();
      expect(c.currentMonth().getMonth()).toBe(0);
      expect(c.currentMonth().getFullYear()).toBe(2025);
    });
  });
});
