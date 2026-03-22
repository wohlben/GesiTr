import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { DayDialog } from './day-dialog';

describe('DayDialog', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideRouter([]), provideTranslocoForTest()],
    });
  });

  function createComponent() {
    return TestBed.createComponent(DayDialog).componentInstance;
  }

  describe('statusClass', () => {
    it('returns green classes for finished', () => {
      const result = createComponent().statusClass('finished');
      expect(result).toContain('bg-green-100');
      expect(result).toContain('text-green-800');
    });

    it('returns blue classes for in_progress', () => {
      const result = createComponent().statusClass('in_progress');
      expect(result).toContain('bg-blue-100');
      expect(result).toContain('text-blue-800');
    });

    it('returns red classes for aborted', () => {
      const result = createComponent().statusClass('aborted');
      expect(result).toContain('bg-red-100');
      expect(result).toContain('text-red-800');
    });

    it('returns gray classes for unknown status', () => {
      const result = createComponent().statusClass('unknown');
      expect(result).toContain('bg-gray-100');
      expect(result).toContain('text-gray-600');
    });
  });
});
