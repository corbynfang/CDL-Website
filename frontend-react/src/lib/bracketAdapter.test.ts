import { describe, it, expect } from 'vitest'
import { adaptBracketMatch, CDL_MAJOR_1_2023_FIXTURE } from './bracketAdapter'
import type { BracketMatch } from '../services/api'

function makeBracketMatch(overrides: Partial<BracketMatch> = {}): BracketMatch {
  return {
    id: 1,
    team1_id: 101,
    team2_id: 202,
    team1_name: 'OpTic Texas',
    team1_abbr: 'OTX',
    team1_logo: '',
    team2_name: 'Atlanta FaZe',
    team2_abbr: 'ATL',
    team2_logo: '',
    team1_score: 3,
    team2_score: 1,
    winner_id: 101,
    bracket_position: 1,
    match_date: '2023-02-03T18:00:00Z',
    ...overrides,
  }
}

describe('adaptBracketMatch', () => {
  it('maps id, name, and round text', () => {
    const result = adaptBracketMatch(makeBracketMatch({ id: 42 }), 'UB R1', 'ubr2m1', null)
    expect(result.id).toBe('42')
    expect(result.name).toBe('OTX vs ATL')
    expect(result.tournamentRoundText).toBe('UB R1')
  })

  it('sets nextMatchId and nextLooserMatchId', () => {
    const result = adaptBracketMatch(makeBracketMatch(), 'UB SF', 'ubf', 'lbr2m1')
    expect(result.nextMatchId).toBe('ubf')
    expect(result.nextLooserMatchId).toBe('lbr2m1')
  })

  it('sets nextLooserMatchId to undefined when not provided', () => {
    const result = adaptBracketMatch(makeBracketMatch(), 'LB R1', null)
    expect(result.nextLooserMatchId).toBeUndefined()
  })

  it('marks state SCORE_DONE for completed match', () => {
    const result = adaptBracketMatch(makeBracketMatch({ winner_id: 101 }), 'UB R1', null)
    expect(result.state).toBe('SCORE_DONE')
  })

  it('marks state NO_PARTY for match without a winner', () => {
    const result = adaptBracketMatch(makeBracketMatch({ winner_id: null }), 'UB R1', null)
    expect(result.state).toBe('NO_PARTY')
  })

  it('sets top participant fields correctly', () => {
    const result = adaptBracketMatch(makeBracketMatch(), 'UB R1', null)
    const top = result.participants[0]
    expect(top.id).toBe('101')
    expect(top.name).toBe('OpTic Texas')
    expect(top.resultText).toBe('3')
    expect(top.isWinner).toBe(true)
    expect(top.status).toBe('PLAYED')
    expect(top.abbr).toBe('OTX')
    expect(top.score).toBe(3)
  })

  it('sets bottom participant fields correctly', () => {
    const result = adaptBracketMatch(makeBracketMatch(), 'UB R1', null)
    const bot = result.participants[1]
    expect(bot.id).toBe('202')
    expect(bot.name).toBe('Atlanta FaZe')
    expect(bot.resultText).toBe('1')
    expect(bot.isWinner).toBe(false)
    expect(bot.status).toBe('PLAYED')
    expect(bot.abbr).toBe('ATL')
  })

  it('sets resultText to null for incomplete match', () => {
    const result = adaptBracketMatch(makeBracketMatch({ winner_id: null, team1_score: 0, team2_score: 0 }), 'GF', null)
    expect(result.participants[0].resultText).toBeNull()
    expect(result.participants[1].resultText).toBeNull()
  })

  it('marks neither participant as winner for incomplete match', () => {
    const result = adaptBracketMatch(makeBracketMatch({ winner_id: null }), 'GF', null)
    expect(result.participants[0].isWinner).toBe(false)
    expect(result.participants[1].isWinner).toBe(false)
  })
})


describe('CDL_MAJOR_1_2023_FIXTURE', () => {
  it('has upper and lower arrays', () => {
    expect(Array.isArray(CDL_MAJOR_1_2023_FIXTURE.upper)).toBe(true)
    expect(Array.isArray(CDL_MAJOR_1_2023_FIXTURE.lower)).toBe(true)
  })

  it('upper bracket has 8 matches (UBR1×4, UBSF×2, UBF×1, GF×1)', () => {
    expect(CDL_MAJOR_1_2023_FIXTURE.upper).toHaveLength(8)
  })

  it('lower bracket has 6 matches (LBR1×2, LBQF×2, LBSF×1, LBF×1)', () => {
    expect(CDL_MAJOR_1_2023_FIXTURE.lower).toHaveLength(6)
  })

  it('every match has exactly 2 participants', () => {
    const all = [...CDL_MAJOR_1_2023_FIXTURE.upper, ...CDL_MAJOR_1_2023_FIXTURE.lower]
    for (const match of all) {
      expect(match.participants).toHaveLength(2)
    }
  })

  it('every match has a unique id', () => {
    const all = [...CDL_MAJOR_1_2023_FIXTURE.upper, ...CDL_MAJOR_1_2023_FIXTURE.lower]
    const ids = all.map(m => m.id)
    const unique = new Set(ids)
    expect(unique.size).toBe(ids.length)
  })

  it('grand final (gf) has nextMatchId null', () => {
    const gf = CDL_MAJOR_1_2023_FIXTURE.upper.find(m => m.id === 'gf')
    expect(gf).toBeDefined()
    expect(gf?.nextMatchId).toBeNull()
  })

  it('every non-final upper match has a nextMatchId', () => {
    const nonFinal = CDL_MAJOR_1_2023_FIXTURE.upper.filter(m => m.id !== 'gf')
    for (const match of nonFinal) {
      expect(match.nextMatchId).not.toBeNull()
    }
  })

  it('every non-final lower match has a nextMatchId', () => {
    const nonFinal = CDL_MAJOR_1_2023_FIXTURE.lower.filter(m => m.id !== 'lbf')
    for (const match of nonFinal) {
      expect(match.nextMatchId).not.toBeNull()
    }
  })

  it('lb final (lbf) points nextMatchId to gf', () => {
    const lbf = CDL_MAJOR_1_2023_FIXTURE.lower.find(m => m.id === 'lbf')
    expect(lbf?.nextMatchId).toBe('gf')
  })

  it('ub r1 matches have nextLooserMatchId pointing into lower bracket', () => {
    const ubr1 = CDL_MAJOR_1_2023_FIXTURE.upper.filter(m =>
      ['ubr1m1', 'ubr1m2', 'ubr1m3', 'ubr1m4'].includes(String(m.id))
    )
    expect(ubr1).toHaveLength(4)
    for (const match of ubr1) {
      expect(match.nextLooserMatchId).not.toBeNull()
    }
  })

  it('every completed match has exactly one winner', () => {
    const all = [...CDL_MAJOR_1_2023_FIXTURE.upper, ...CDL_MAJOR_1_2023_FIXTURE.lower]
    const complete = all.filter(m => m.state === 'SCORE_DONE')
    for (const match of complete) {
      const winners = match.participants.filter((p: { isWinner?: boolean }) => p.isWinner)
      expect(winners).toHaveLength(1)
    }
  })
})
