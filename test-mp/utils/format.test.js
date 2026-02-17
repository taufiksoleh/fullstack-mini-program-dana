import { formatRupiah } from './format';

describe('formatRupiah', () => {
  // Standard cases
  test('formats 6-digit number', () => {
    expect(formatRupiah(750000)).toBe('750.000');
  });
  test('formats 8-digit number', () => {
    expect(formatRupiah(10000000)).toBe('10.000.000');
  });
  test('formats exactly 1000', () => {
    expect(formatRupiah(1000)).toBe('1.000');
  });

  // Below 1000 — no separator needed
  test('leaves numbers below 1000 unchanged', () => {
    expect(formatRupiah(500)).toBe('500');
  });
  test('handles single digit', () => {
    expect(formatRupiah(5)).toBe('5');
  });

  // Zero
  test('formats 0', () => {
    expect(formatRupiah(0)).toBe('0');
  });

  // String input
  test('accepts numeric string', () => {
    expect(formatRupiah('10000000')).toBe('10.000.000');
  });
  test('strips and reformats already-dotted string', () => {
    expect(formatRupiah('750.000')).toBe('750.000');
  });

  // Invalid / empty input
  test('returns "0" for empty string', () => {
    expect(formatRupiah('')).toBe('0');
  });
  test('returns "0" for non-numeric string', () => {
    expect(formatRupiah('abc')).toBe('0');
  });
  test('returns "0" for null', () => {
    expect(formatRupiah(null)).toBe('0');
  });
  test('returns "0" for undefined', () => {
    expect(formatRupiah(undefined)).toBe('0');
  });
});
