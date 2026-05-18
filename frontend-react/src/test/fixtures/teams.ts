import type { TournamentTeam } from '../../types'

function makeTeam(overrides: Partial<TournamentTeam> & { id: number; name: string }): TournamentTeam {
  return {
    abbreviation: 'TM',
    is_active: true,
    is_cdl_franchise: true,
    matches_won: 0,
    matches_lost: 0,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    ...overrides,
  }
}

export const opticTexas = makeTeam({
  id: 1,
  name: 'OpTic Texas',
  abbreviation: 'OTX',
  placement: 1,
  matches_won: 5,
  matches_lost: 1,
})

export const atlantaFaze = makeTeam({
  id: 2,
  name: 'Atlanta FaZe',
  abbreviation: 'ATL',
  placement: 2,
  matches_won: 4,
  matches_lost: 2,
})

export const bostonBreach = makeTeam({
  id: 3,
  name: 'Boston Breach',
  abbreviation: 'BOS',
  placement: 3,
  matches_won: 3,
  matches_lost: 2,
})

/** Team with no known logo — exercises the abbreviation fallback. */
export const unknownTeam = makeTeam({
  id: 99,
  name: 'Unknown Team FC',
  abbreviation: 'UNK',
  placement: null,
  matches_won: 1,
  matches_lost: 3,
})

export const sampleTeams: TournamentTeam[] = [opticTexas, atlantaFaze, bostonBreach]

/** 14 teams — exercises the "+N" overflow count logic in TeamLogoStrip. */
export const largeTeamList: TournamentTeam[] = Array.from({ length: 14 }, (_, i) =>
  makeTeam({ id: i + 1, name: `Team ${i + 1}`, abbreviation: `T${i + 1}` })
)
