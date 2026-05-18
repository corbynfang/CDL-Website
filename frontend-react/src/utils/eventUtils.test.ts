import { describe, it, expect } from 'vitest'
import {
  isFeatured,
  hasBracket,
  isHidden,
  deriveStatus,
  formatDateRange,
  formatPrize,
  monthLabel,
  countryFlag,
  formatRound,
  roundOrder,
  bracketSection,
  sortedRounds,
  groupByMonth,
} from './eventUtils'

// ── isFeatured ────────────────────────────────────────────────────────────────

describe('isFeatured', () => {
  it.each([
    ['major_tournament',   true],
    ['championship',       true],
    ['international_major', true],
    ['kickoff',            false],
    ['qualifier',          false],
    ['season_summary',     false],
    ['unknown',            false],
  ])('isFeatured(%s) → %s', (type, expected) => {
    expect(isFeatured(type)).toBe(expected)
  })
})

// ── hasBracket ────────────────────────────────────────────────────────────────

describe('hasBracket', () => {
  it.each([
    ['major_tournament',   true],
    ['championship',       true],
    ['international_major', true],
    ['kickoff',            true],
    ['qualifier',          false],
    ['season_summary',     false],
    ['unknown',            false],
  ])('hasBracket(%s) → %s', (type, expected) => {
    expect(hasBracket(type)).toBe(expected)
  })
})

// ── isHidden ──────────────────────────────────────────────────────────────────

describe('isHidden', () => {
  it.each([
    ['season_summary', true],
    ['unknown',        true],
    ['qualifier',      false],
    ['major_tournament', false],
    ['championship',   false],
  ])('isHidden(%s) → %s', (type, expected) => {
    expect(isHidden(type)).toBe(expected)
  })
})

// ── deriveStatus ──────────────────────────────────────────────────────────────

describe('deriveStatus', () => {
  it('returns upcoming when start date is in the future', () => {
    expect(deriveStatus('2099-01-01T00:00:00Z')).toBe('upcoming')
  })

  it('returns upcoming when start is future even with a future end date', () => {
    expect(deriveStatus('2099-01-01T00:00:00Z', '2099-01-05T00:00:00Z')).toBe('upcoming')
  })

  it('returns live when start is past and end date is in the future', () => {
    expect(deriveStatus('2020-01-01T00:00:00Z', '2099-12-31T00:00:00Z')).toBe('live')
  })

  it('returns completed when both start and end are in the past', () => {
    expect(deriveStatus('2020-01-01T00:00:00Z', '2020-01-05T00:00:00Z')).toBe('completed')
  })

  it('returns live with no end date when started within the last 5 days', () => {
    const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000).toISOString()
    expect(deriveStatus(oneHourAgo)).toBe('live')
  })

  it('returns live with no end date when started just under 5 days ago', () => {
    const fourDaysAgo = new Date(Date.now() - 4 * 24 * 60 * 60 * 1000).toISOString()
    expect(deriveStatus(fourDaysAgo)).toBe('live')
  })

  it('returns completed with no end date when started more than 5 days ago', () => {
    const sixDaysAgo = new Date(Date.now() - 6 * 24 * 60 * 60 * 1000).toISOString()
    expect(deriveStatus(sixDaysAgo)).toBe('completed')
  })

  it('returns completed when end date is null (treated as missing)', () => {
    // null end_date → falls through to the 5-day check; 2020 is way more than 5 days ago
    expect(deriveStatus('2020-01-01T00:00:00Z', null)).toBe('completed')
  })
})

// ── formatDateRange ───────────────────────────────────────────────────────────

describe('formatDateRange', () => {
  it('includes the year when no end date is provided', () => {
    const result = formatDateRange('2025-06-15T12:00:00Z')
    expect(result).toMatch(/2025/)
  })

  it('includes an en-dash separator when an end date is provided', () => {
    const result = formatDateRange('2025-06-10T12:00:00Z', '2025-06-15T12:00:00Z')
    expect(result).toContain('–')
  })

  it('includes the year in output when same-month range is given', () => {
    const result = formatDateRange('2025-06-10T12:00:00Z', '2025-06-15T12:00:00Z')
    expect(result).toMatch(/2025/)
  })

  it('includes both month abbreviations for a cross-month range', () => {
    const result = formatDateRange('2025-03-15T12:00:00Z', '2025-04-10T12:00:00Z')
    // Both month abbreviations must appear
    expect(result).toMatch(/Mar/)
    expect(result).toMatch(/Apr/)
  })

  it('uses a single month abbreviation for a same-month range', () => {
    const result = formatDateRange('2025-06-10T12:00:00Z', '2025-06-20T12:00:00Z')
    // "Jun" appears once, not twice
    const junCount = (result.match(/Jun/g) ?? []).length
    expect(junCount).toBe(1)
  })

  it('handles null end date the same as missing end date', () => {
    const withNull    = formatDateRange('2025-06-15T12:00:00Z', null)
    const withMissing = formatDateRange('2025-06-15T12:00:00Z')
    expect(withNull).toBe(withMissing)
  })
})

// ── formatPrize ───────────────────────────────────────────────────────────────

describe('formatPrize', () => {
  it('returns "TBA" for undefined', () => {
    expect(formatPrize(undefined)).toBe('TBA')
  })

  it('returns "TBA" for null', () => {
    expect(formatPrize(null)).toBe('TBA')
  })

  it('returns "TBA" for 0', () => {
    expect(formatPrize(0)).toBe('TBA')
  })

  it('formats thousands as K', () => {
    expect(formatPrize(375000)).toBe('$375K')
    expect(formatPrize(500000)).toBe('$500K')
  })

  it('formats exact millions without decimal', () => {
    expect(formatPrize(2000000)).toBe('$2M')
    expect(formatPrize(1000000)).toBe('$1M')
  })

  it('formats non-exact millions with one decimal', () => {
    expect(formatPrize(1500000)).toBe('$1.5M')
    expect(formatPrize(2380000)).toBe('$2.4M')
  })

  it('formats amounts under 1000 as raw dollars', () => {
    expect(formatPrize(500)).toBe('$500')
    expect(formatPrize(1)).toBe('$1')
  })
})

// ── monthLabel ────────────────────────────────────────────────────────────────

describe('monthLabel', () => {
  it('includes the year', () => {
    expect(monthLabel('2025-03-15T12:00:00Z')).toMatch(/2025/)
  })

  it('includes a long month name', () => {
    // March in any locale variant
    expect(monthLabel('2025-03-15T12:00:00Z')).toMatch(/March|Mar/)
  })

  it('returns different labels for different months', () => {
    const march = monthLabel('2025-03-15T12:00:00Z')
    const june  = monthLabel('2025-06-15T12:00:00Z')
    expect(march).not.toBe(june)
  })
})

// ── countryFlag ───────────────────────────────────────────────────────────────

describe('countryFlag', () => {
  it('returns 🇺🇸 for USA', () => {
    expect(countryFlag('USA')).toBe('🇺🇸')
  })

  it('returns 🇨🇦 for CAN', () => {
    expect(countryFlag('CAN')).toBe('🇨🇦')
  })

  it('returns 🇸🇦 for SAU', () => {
    expect(countryFlag('SAU')).toBe('🇸🇦')
  })

  it('is case-insensitive', () => {
    expect(countryFlag('usa')).toBe('🇺🇸')
    expect(countryFlag('Can')).toBe('🇨🇦')
  })

  it('returns 🌐 for an unknown country code', () => {
    expect(countryFlag('ZZZ')).toBe('🌐')
    expect(countryFlag('')).toBe('🌐')
  })
})

// ── formatRound ───────────────────────────────────────────────────────────────

describe('formatRound', () => {
  it.each([
    ['winners_r1',     'Winners Round 1'],
    ['winners_r2',     'Winners Round 2'],
    ['winners_r3',     'Winners Round 3'],
    ['winners_finals', 'Winners Finals'],
    ['elim_r1',        'Elimination Round 1'],
    ['elim_r2',        'Elimination Round 2'],
    ['elim_r3',        'Elimination Round 3'],
    ['elim_finals',    'Elimination Finals'],
    ['grand_finals',   'Grand Finals'],
  ])('formats %s → "%s"', (input, expected) => {
    expect(formatRound(input)).toBe(expected)
  })

  it('title-cases unknown round keys by replacing underscores with spaces', () => {
    expect(formatRound('custom_round')).toBe('Custom Round')
    expect(formatRound('semi_final')).toBe('Semi Final')
  })
})

// ── roundOrder ────────────────────────────────────────────────────────────────

describe('roundOrder', () => {
  it('places winners rounds before elimination rounds', () => {
    expect(roundOrder('winners_r1')).toBeLessThan(roundOrder('elim_r1'))
  })

  it('places elimination rounds before grand finals', () => {
    expect(roundOrder('elim_finals')).toBeLessThan(roundOrder('grand_finals'))
  })

  it('orders winners rounds sequentially', () => {
    expect(roundOrder('winners_r1')).toBeLessThan(roundOrder('winners_r2'))
    expect(roundOrder('winners_r2')).toBeLessThan(roundOrder('winners_r3'))
    expect(roundOrder('winners_r3')).toBeLessThan(roundOrder('winners_finals'))
  })

  it('orders elimination rounds sequentially', () => {
    expect(roundOrder('elim_r1')).toBeLessThan(roundOrder('elim_r2'))
    expect(roundOrder('elim_r2')).toBeLessThan(roundOrder('elim_r3'))
    expect(roundOrder('elim_r3')).toBeLessThan(roundOrder('elim_finals'))
  })

  it('returns 99 for unknown rounds so they sort last', () => {
    expect(roundOrder('unknown_round')).toBe(99)
  })
})

// ── bracketSection ────────────────────────────────────────────────────────────

describe('bracketSection', () => {
  it.each([
    ['winners_r1',     'winners'],
    ['winners_r2',     'winners'],
    ['winners_finals', 'winners'],
    ['elim_r1',        'elimination'],
    ['elim_r2',        'elimination'],
    ['elim_finals',    'elimination'],
    ['grand_finals',   'grand_finals'],
    ['custom_round',   'elimination'], // anything not winners/grand_finals → elimination
  ])('bracketSection(%s) → "%s"', (input, expected) => {
    expect(bracketSection(input)).toBe(expected)
  })
})

// ── sortedRounds ──────────────────────────────────────────────────────────────

describe('sortedRounds', () => {
  it('sorts rounds in bracket progression order', () => {
    const input = ['grand_finals', 'winners_r1', 'elim_r1', 'winners_finals']
    expect(sortedRounds(input)).toEqual(['winners_r1', 'winners_finals', 'elim_r1', 'grand_finals'])
  })

  it('deduplicates repeated round names', () => {
    const input = ['winners_r1', 'winners_r1', 'grand_finals']
    expect(sortedRounds(input)).toEqual(['winners_r1', 'grand_finals'])
  })

  it('returns an empty array for empty input', () => {
    expect(sortedRounds([])).toEqual([])
  })

  it('returns a single-element array unchanged', () => {
    expect(sortedRounds(['grand_finals'])).toEqual(['grand_finals'])
  })

  it('places unknown rounds at the end', () => {
    const result = sortedRounds(['custom_round', 'winners_r1'])
    expect(result[0]).toBe('winners_r1')
    expect(result[result.length - 1]).toBe('custom_round')
  })
})

// ── groupByMonth ──────────────────────────────────────────────────────────────

describe('groupByMonth', () => {
  const march1  = { id: 1, start_date: '2025-03-10T12:00:00Z' }
  const march2  = { id: 2, start_date: '2025-03-25T12:00:00Z' }
  const june1   = { id: 3, start_date: '2025-06-15T12:00:00Z' }
  const june2   = { id: 4, start_date: '2025-06-20T12:00:00Z' }

  it('returns an empty array for empty input', () => {
    expect(groupByMonth([])).toEqual([])
  })

  it('groups items with the same month together', () => {
    const result = groupByMonth([march1, march2])
    expect(result).toHaveLength(1)
    expect(result[0][1]).toHaveLength(2)
  })

  it('creates separate groups for different months', () => {
    const result = groupByMonth([march1, june1])
    expect(result).toHaveLength(2)
  })

  it('preserves item order within each group', () => {
    const result = groupByMonth([march1, march2, june1, june2])
    const marchGroup = result.find(([label]) => label.includes('March'))!
    expect(marchGroup[1][0]).toEqual(march1)
    expect(marchGroup[1][1]).toEqual(march2)
  })

  it('uses month labels that include the year', () => {
    const result = groupByMonth([march1])
    const [label] = result[0]
    expect(label).toMatch(/2025/)
  })

  it('preserves insertion order of groups (first seen month appears first)', () => {
    const result = groupByMonth([march1, june1, march2])
    // March was seen first
    expect(result[0][0]).toMatch(/March/)
  })
})
