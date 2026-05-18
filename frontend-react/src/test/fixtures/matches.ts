import type { Match } from '../../types'
import type { BracketMatch, BracketData } from '../../services/api'

function makeMatch(overrides: Partial<Match> & { id: number }): Match {
  return {
    tournament_id: 1,
    team1_id: 1,
    team2_id: 2,
    match_date: '2025-03-01T18:00:00Z',
    team1_score: 3,
    team2_score: 1,
    winner_id: 1,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    team1: { id: 1, name: 'OpTic Texas',   abbreviation: 'OTX', is_active: true, is_cdl_franchise: true, created_at: '', updated_at: '' },
    team2: { id: 2, name: 'Atlanta FaZe',  abbreviation: 'ATL', is_active: true, is_cdl_franchise: true, created_at: '', updated_at: '' },
    ...overrides,
  }
}

export const winnersR1Match = makeMatch({ id: 1, match_type: 'winners_r1', team1_score: 3, team2_score: 0, winner_id: 1 })
export const grandFinalsMatch = makeMatch({ id: 2, match_type: 'grand_finals', team1_score: 3, team2_score: 2, winner_id: 1 })
export const noTypeMatch = makeMatch({ id: 3, match_type: undefined })

export const sampleMatches: Match[] = [winnersR1Match, grandFinalsMatch]

function makeBracketMatch(overrides: Partial<BracketMatch> & { id: number }): BracketMatch {
  return {
    team1_id: 1,
    team2_id: 2,
    team1_name: 'OpTic Texas',
    team1_abbr: 'OTX',
    team1_logo: '',
    team2_name: 'Atlanta FaZe',
    team2_abbr: 'ATL',
    team2_logo: '',
    team1_score: 3,
    team2_score: 1,
    winner_id: 1,
    bracket_position: 1,
    match_date: '2025-03-01T18:00:00Z',
    ...overrides,
  }
}

export const bracketMatchComplete  = makeBracketMatch({ id: 10, winner_id: 1 })
export const bracketMatchNoWinner  = makeBracketMatch({ id: 11, winner_id: null, team1_score: 0, team2_score: 0 })

export const sampleBracketData: BracketData = {
  tournament_id: 1,
  tournament_name: 'CDL Major 2 2025',
  total_matches: 2,
  bracket: {
    winners_r1:   [bracketMatchComplete],
    grand_finals: [bracketMatchNoWinner],
  },
}

export const emptyBracketData: BracketData = {
  tournament_id: 1,
  tournament_name: 'CDL Major 2 2025',
  total_matches: 0,
  bracket: {},
}
