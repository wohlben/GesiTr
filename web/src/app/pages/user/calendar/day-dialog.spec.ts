import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { DayDialog } from './day-dialog';

describe('DayDialog', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideRouter([])],
    });
  });

  function createComponent() {
    return TestBed.createComponent(DayDialog).componentInstance;
  }

  describe('statusLabel', () => {
    it('replaces underscores with spaces', () => {
      expect(createComponent().statusLabel('in_progress')).toBe('in progress');
    });

    it('replaces multiple underscores', () => {
      expect(createComponent().statusLabel('some_multi_word_status')).toBe(
        'some multi word status',
      );
    });

    it('returns status unchanged when no underscores', () => {
      expect(createComponent().statusLabel('finished')).toBe('finished');
    });
  });

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
