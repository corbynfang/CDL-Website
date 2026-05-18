import { describe, it, expect } from 'vitest'
import { getTeamTheme } from './teamThemes'

// ── Fallback ──────────────────────────────────────────────────────────────────

describe('getTeamTheme — fallback', () => {
  it('returns a fallback theme for an unknown team name', () => {
    const theme = getTeamTheme('Some Unknown Team')
    expect(theme.primary).toBe('#4a4a5a')
    expect(theme.bg).toBe('#111111')
  })

  it('returns a fallback theme for an empty string', () => {
    const theme = getTeamTheme('')
    expect(theme.primary).toBe('#4a4a5a')
  })

  it('returns an object with primary, bg, and glow keys', () => {
    const theme = getTeamTheme('does not exist')
    expect(theme).toHaveProperty('primary')
    expect(theme).toHaveProperty('bg')
    expect(theme).toHaveProperty('glow')
  })
})

// ── Case insensitivity ────────────────────────────────────────────────────────

describe('getTeamTheme — case insensitivity', () => {
  it('matches exact lowercase', () => {
    expect(getTeamTheme('optic texas').primary).toBe('#78BE20')
  })

  it('matches uppercase input', () => {
    expect(getTeamTheme('OPTIC TEXAS').primary).toBe('#78BE20')
  })

  it('matches mixed case input', () => {
    expect(getTeamTheme('OpTic TeXas').primary).toBe('#78BE20')
  })
})

// ── CDL franchise themes ──────────────────────────────────────────────────────

describe('getTeamTheme — CDL franchise primary colors', () => {
  it.each([
    ['optic texas',          '#78BE20'],
    ['optic chicago',        '#78BE20'],
    ['dallas empire',        '#A8272A'],
    ['atlanta faze',         '#E8002D'],
    ['faze vegas',           '#E8002D'],
    ['los angeles thieves',  '#E63329'],
    ['100 thieves',          '#E63329'],
    ['boston breach',        '#3BAA35'],
    ['toronto ultra',        '#8B1FFF'],
    ['toronto koi',          '#6A2DFF'],
    ['los angeles guerrillas', '#4CAF50'],
    ['la guerrillas m8',     '#B8A7FF'],
    ['miami heretics',       '#FF6B35'],
    ['team heretics',        '#FF6B35'],
    ['new york subliners',   '#F6EB14'],
    ['cloud9 new york',      '#00AEEF'],
    ['seattle surge',        '#00FF87'],
    ['vancouver surge',      '#00A7E1'],
    ['minnesota røkkr',      '#5B2D8E'],
    ['g2 minnesota',         '#FF4C00'],
    ['carolina royal ravens', '#1B4F8E'],
    ['london royal ravens',  '#1B4F8E'],
    ['paris legion',         '#003087'],
    ['florida mutineers',    '#FFCD00'],
    ['riyadh falcons',       '#00A86B'],
    ['las vegas legion',     '#E4B062'],
  ])('%s → primary %s', (team, expectedPrimary) => {
    expect(getTeamTheme(team).primary).toBe(expectedPrimary)
  })
})

// ── Background colors ─────────────────────────────────────────────────────────

describe('getTeamTheme — background colors', () => {
  it('most teams use #111111 background', () => {
    expect(getTeamTheme('optic texas').bg).toBe('#111111')
    expect(getTeamTheme('atlanta faze').bg).toBe('#111111')
    expect(getTeamTheme('boston breach').bg).toBe('#111111')
  })

  it('los angeles thieves uses #0B0B0B background', () => {
    expect(getTeamTheme('los angeles thieves').bg).toBe('#0B0B0B')
    expect(getTeamTheme('100 thieves').bg).toBe('#0B0B0B')
  })

  it('la guerrillas m8 uses #2B2B2B background', () => {
    expect(getTeamTheme('la guerrillas m8').bg).toBe('#2B2B2B')
  })

  it('cloud9 new york uses #101820 background', () => {
    expect(getTeamTheme('cloud9 new york').bg).toBe('#101820')
  })
})

// ── Glow colors ───────────────────────────────────────────────────────────────

describe('getTeamTheme — glow format', () => {
  it('all glow values are rgba strings', () => {
    const teams = [
      'optic texas', 'atlanta faze', 'boston breach', 'toronto ultra',
      'miami heretics', 'new york subliners', 'seattle surge', 'paris legion',
    ]
    for (const team of teams) {
      expect(getTeamTheme(team).glow).toMatch(/^rgba\(/)
    }
  })

  it('fallback glow is a lower-opacity rgba', () => {
    expect(getTeamTheme('unknown team').glow).toMatch(/0\.12\)$/)
  })
})
