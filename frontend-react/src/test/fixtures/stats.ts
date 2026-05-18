import type { PlayerTournamentStats } from '../../types'

function makeStat(overrides: Partial<PlayerTournamentStats> & { id: number; player_id: number }): PlayerTournamentStats {
  return {
    team_id: 1,
    tournament_id: 1,
    total_kills: 100,
    total_deaths: 80,
    total_assists: 10,
    total_damage: 50000,
    kd_ratio: 1.25,
    kda_ratio: 1.38,
    overall_maps: 8,
    overall_plus_minus: 20,
    player: { id: overrides.player_id, gamertag: `Player${overrides.player_id}`, is_active: true, created_at: '', updated_at: '' },
    team: { id: 1, name: 'OpTic Texas', abbreviation: 'OTX', is_active: true, is_cdl_franchise: true, created_at: '', updated_at: '' },
    ...overrides,
  }
}

export const highKDStat = makeStat({ id: 1, player_id: 1, kd_ratio: 1.45, total_kills: 120, overall_maps: 10,
  player: { id: 1, gamertag: 'Scump', is_active: true, created_at: '', updated_at: '' } })

export const lowKDStat = makeStat({ id: 2, player_id: 2, kd_ratio: 0.88, total_kills: 70, overall_maps: 8,
  player: { id: 2, gamertag: 'Simp', is_active: true, created_at: '', updated_at: '' },
  team: { id: 2, name: 'Atlanta FaZe', abbreviation: 'ATL', is_active: true, is_cdl_franchise: true, created_at: '', updated_at: '' } })

/** Stat with no player or team attached — exercises the fallback display. */
export const noPlayerStat = makeStat({ id: 3, player_id: 99, kd_ratio: 1.10, player: undefined, team: undefined })

export const sampleStats: PlayerTournamentStats[] = [highKDStat, lowKDStat]
