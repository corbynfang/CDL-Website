// Fixtures for MatchDetail component tests.
// Shape mirrors GetMatch response in matches.go.

export interface PlayerStat {
  player_id: number
  gamertag: string
  kills: number
  deaths: number
  kd_ratio: number
  damage: number
  assists: number
  bp_rating: number
  hill_time: number
  snd_rounds: number
  plant_count: number
  defuse_count: number
  first_blood_count: number
  first_death_count: number
  non_traded_kills: number
  highest_streak: number
  data_quality_note?: string
}

export interface MapDetail {
  map_number: number
  map_name: string
  mode: string
  score_1: number
  score_2: number
  winner_id: number | null
  duration_sec: number
  played: boolean
  team1_stats: PlayerStat[]
  team2_stats: PlayerStat[]
}

export interface MatchDetailResponse {
  match: {
    id: number
    tournament_id: number
    tournament_name: string
    tournament_slug: string
    season_name: string
    game_code: string
    team1_id: number
    team1_name: string
    team1_abbr: string
    team1_logo: string
    team2_id: number
    team2_name: string
    team2_abbr: string
    team2_logo: string
    team1_score: number
    team2_score: number
    winner_id: number | null
    match_date: string
    format: string
    bracket_round: string
  }
  maps: MapDetail[]
}

function makePlayerStat(overrides: Partial<PlayerStat> & { player_id: number; gamertag: string }): PlayerStat {
  return {
    kills: 20,
    deaths: 15,
    kd_ratio: 1.33,
    damage: 3500,
    assists: 2,
    bp_rating: 1.05,
    hill_time: 0,
    snd_rounds: 0,
    plant_count: 0,
    defuse_count: 0,
    first_blood_count: 0,
    first_death_count: 0,
    non_traded_kills: 0,
    highest_streak: 5,
    ...overrides,
  }
}

export const shotzzy = makePlayerStat({ player_id: 1, gamertag: 'Shotzzy', kills: 24, deaths: 14, kd_ratio: 1.71 })
export const cellium  = makePlayerStat({ player_id: 2, gamertag: 'Cellium',  kills: 18, deaths: 16, kd_ratio: 1.13 })

const baseMatch = {
  id: 42,
  tournament_id: 1,
  tournament_name: 'CDL Major 1 2025',
  tournament_slug: 'cdl-major-1-2025',
  season_name: 'BO6 Season 2025',
  game_code: 'BO6',
  team1_id: 1,
  team1_name: 'OpTic Texas',
  team1_abbr: 'OTX',
  team1_logo: '',
  team2_id: 2,
  team2_name: 'Atlanta FaZe',
  team2_abbr: 'ATL',
  team2_logo: '',
  team1_score: 3,
  team2_score: 1,
  winner_id: 1,
  match_date: '2025-03-05T18:00:00Z',
  format: 'BO5',
  bracket_round: 'winners_r1',
}

export const matchDetailFixture: MatchDetailResponse = {
  match: baseMatch,
  maps: [
    {
      map_number: 1,
      map_name: 'Skyline',
      mode: 'Hardpoint',
      score_1: 250,
      score_2: 200,
      winner_id: 1,
      duration_sec: 480,
      played: true,
      team1_stats: [shotzzy],
      team2_stats: [cellium],
    },
    {
      map_number: 2,
      map_name: 'Rewind',
      mode: 'Search and Destroy',
      score_1: 6,
      score_2: 4,
      winner_id: 1,
      duration_sec: 600,
      played: true,
      team1_stats: [shotzzy],
      team2_stats: [cellium],
    },
  ],
}

// Match with no maps — exercises the case where maps array is empty.
export const matchDetailNoMapsFixture: MatchDetailResponse = {
  match: { ...baseMatch, id: 43 },
  maps: [],
}

// Match with a map but no player stats — exercises the empty-stats branch.
export const matchDetailEmptyStatsFixture: MatchDetailResponse = {
  match: { ...baseMatch, id: 44 },
  maps: [
    {
      map_number: 1,
      map_name: 'Skyline',
      mode: 'Hardpoint',
      score_1: 250,
      score_2: 200,
      winner_id: 1,
      duration_sec: 0,
      played: true,
      team1_stats: [],
      team2_stats: [],
    },
  ],
}

// Map with zero damage — exercises the "—" fallback in PlayerRow.
export const matchDetailZeroDamageFixture: MatchDetailResponse = {
  match: baseMatch,
  maps: [
    {
      map_number: 1,
      map_name: 'Hacienda',
      mode: 'Control',
      score_1: 3,
      score_2: 2,
      winner_id: 1,
      duration_sec: 720,
      played: true,
      team1_stats: [makePlayerStat({ player_id: 1, gamertag: 'Shotzzy', damage: 0 })],
      team2_stats: [makePlayerStat({ player_id: 2, gamertag: 'Cellium', damage: 0 })],
    },
  ],
}
