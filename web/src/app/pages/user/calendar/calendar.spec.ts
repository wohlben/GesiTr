import { TestBed } from '@angular/core/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { Calendar } from './calendar';

describe('Calendar', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
        { provide: UserApiClient, useValue: { fetchWorkoutLogs: vi.fn().mockResolvedValue([]) } },
      ],
    });
  });

  function createComponent() {
    const fixture = TestBed.createComponent(Calendar);
    fixture.detectChanges();
    return fixture.componentInstance;
  }

  describe('calendarCells', () => {
    it('generates correct cells for month starting on Monday', () => {
      const c = createComponent();
      // Jan 1 2024 is Monday -> offset 0
      c.currentMonth.set(new Date(2024, 0, 1));
      const cells = c.calendarCells();
      expect(cells.length).toBe(31);
      expect(cells[0]).toEqual(expect.objectContaining({ day: 1 }));
      expect(cells[30]).toEqual(expect.objectContaining({ day: 31 }));
    });

    it('generates null offsets for month starting on Sunday', () => {
      const c = createComponent();
      // Sep 1 2024 is Sunday -> offset 6
      c.currentMonth.set(new Date(2024, 8, 1));
      const cells = c.calendarCells();
      expect(cells.length).toBe(36); // 6 nulls + 30 days
      expect(cells.slice(0, 6).every((cell) => cell === null)).toBe(true);
      expect(cells[6]).toEqual(expect.objectContaining({ day: 1 }));
    });

    it('handles February in leap year', () => {
      const c = createComponent();
      // Feb 2024 is leap year, Feb 1 is Thursday -> offset 3
      c.currentMonth.set(new Date(2024, 1, 1));
      const cells = c.calendarCells();
      const dayCells = cells.filter((cell) => cell !== null);
      expect(dayCells.length).toBe(29);
    });

    it('handles February in non-leap year', () => {
      const c = createComponent();
      // Feb 2023, Feb 1 is Wednesday -> offset 2
      c.currentMonth.set(new Date(2023, 1, 1));
      const cells = c.calendarCells();
      const dayCells = cells.filter((cell) => cell !== null);
      expect(dayCells.length).toBe(28);
    });

    it('marks today correctly', () => {
      const c = createComponent();
      const today = new Date();
      c.currentMonth.set(new Date(today.getFullYear(), today.getMonth(), 1));
      const cells = c.calendarCells();
      const todayCells = cells.filter((cell) => cell !== null && cell.isToday);
      expect(todayCells.length).toBe(1);
      expect(todayCells[0]!.day).toBe(today.getDate());
    });

    it('does not mark today in other months', () => {
      const c = createComponent();
      // Use a month far in the past
      c.currentMonth.set(new Date(2000, 0, 1));
      const cells = c.calendarCells();
      const todayCells = cells.filter((cell) => cell !== null && cell.isToday);
      expect(todayCells.length).toBe(0);
    });
  });

  describe('monthLabel', () => {
    it('formats month and year', () => {
      const c = createComponent();
      c.currentMonth.set(new Date(2024, 0, 1));
      expect(c.monthLabel()).toBe('January 2024');
    });

    it('formats December correctly', () => {
      const c = createComponent();
      c.currentMonth.set(new Date(2024, 11, 1));
      expect(c.monthLabel()).toBe('December 2024');
    });
  });

  describe('navigation', () => {
    it('prevMonth goes to previous month', () => {
      const c = createComponent();
      c.currentMonth.set(new Date(2024, 3, 1));
      c.prevMonth();
      expect(c.currentMonth().getMonth()).toBe(2);
    });

    it('nextMonth goes to next month', () => {
      const c = createComponent();
      c.currentMonth.set(new Date(2024, 3, 1));
      c.nextMonth();
      expect(c.currentMonth().getMonth()).toBe(4);
    });

    it('prevMonth wraps year boundary', () => {
      const c = createComponent();
      c.currentMonth.set(new Date(2024, 0, 1));
      c.prevMonth();
      expect(c.currentMonth().getMonth()).toBe(11);
      expect(c.currentMonth().getFullYear()).toBe(2023);
    });

    it('nextMonth wraps year boundary', () => {
      const c = createComponent();
      c.currentMonth.set(new Date(2024, 11, 1));
      c.nextMonth();
      expect(c.currentMonth().getMonth()).toBe(0);
      expect(c.currentMonth().getFullYear()).toBe(2025);
    });
  });
});
