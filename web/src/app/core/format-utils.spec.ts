import { formatBreak, formatCountdown } from './format-utils';

describe('formatBreak', () => {
  it('returns empty string for null', () => {
    expect(formatBreak(null)).toBe('');
  });

  it('returns empty string for undefined', () => {
    expect(formatBreak(undefined)).toBe('');
  });

  it('formats seconds under 60', () => {
    expect(formatBreak(45)).toBe('45s');
  });

  it('formats exact minutes', () => {
    expect(formatBreak(60)).toBe('1m');
    expect(formatBreak(120)).toBe('2m');
    expect(formatBreak(180)).toBe('3m');
  });

  it('formats minutes with remaining seconds', () => {
    expect(formatBreak(90)).toBe('1m 30s');
    expect(formatBreak(150)).toBe('2m 30s');
  });

  it('formats zero seconds', () => {
    expect(formatBreak(0)).toBe('0s');
  });
});

describe('formatCountdown', () => {
  it('formats seconds under 60 with s suffix', () => {
    expect(formatCountdown(10)).toBe('10s');
    expect(formatCountdown(45)).toBe('45s');
  });

  it('formats zero seconds', () => {
    expect(formatCountdown(0)).toBe('0s');
  });

  it('formats exact minutes as M:00', () => {
    expect(formatCountdown(60)).toBe('1:00');
    expect(formatCountdown(120)).toBe('2:00');
  });

  it('formats minutes with seconds as M:SS', () => {
    expect(formatCountdown(90)).toBe('1:30');
    expect(formatCountdown(65)).toBe('1:05');
  });
});
