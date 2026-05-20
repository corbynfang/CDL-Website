// Fixtures for PlayerDetail component tests.
// Shapes mirror what the Go backend actually returns (see players.go GetPlayerKDStats / GetPlayerMatches).

export interface PlayerKDResponse {
  player_id: number
  gamertag: string
  avatar_url: string
  total_kills: number
  total_deaths: number
  total_assists: number
  avg_kd: number
  hp_kd_ratio: number
  snd_kd_ratio: number
  control_kd_ratio: number
  tournament_stats: Array<{
    tournament_id: number
    tournament_name: string
    kills: number
    deaths: number
    assists: number
    kd_ratio: number
    maps_played: number
  }>
}

export interface PlayerMatchesResponse {
  player_id: number
  events: Array<{
    event: string
    year: number
    tournament_id: number
    matches: Array<{
      match_id: number
      date: string
      opponent: string
      opponent_abbr: string
      result: string
      kd: number | null
      kills: number
      deaths: number
    }>
  }>
  total: number
}

export const playerFixture = {
  id: 1,
  gamertag: 'Shotzzy',
  first_name: 'Anthony',
  last_name: 'Cuevas-Castro',
  country: 'US',
  role: 'slayer',
  is_active: true,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

export const playerKDFixture: PlayerKDResponse = {
  player_id: 1,
  gamertag: 'Shotzzy',
  avatar_url: '',
  total_kills: 500,
  total_deaths: 380,
  total_assists: 40,
  avg_kd: 1.32,
  hp_kd_ratio: 1.45,
  snd_kd_ratio: 1.21,
  control_kd_ratio: 0.98,
  tournament_stats: [
    {
      tournament_id: 1,
      tournament_name: 'CDL Major 1 2025',
      kills: 120,
      deaths: 90,
      assists: 10,
      kd_ratio: 1.33,
      maps_played: 10,
    },
  ],
}

// Backend returns 0 (not null) for control_kd_ratio when no Control maps played.
// This fixture documents that contract so tests can assert the current (buggy) display.
export const playerKDNoControlFixture: PlayerKDResponse = {
  ...playerKDFixture,
  control_kd_ratio: 0,
}

// 6 total matches across 2 events — last5 should show exactly 5.
export const playerMatchesFixture: PlayerMatchesResponse = {
  player_id: 1,
  events: [
    {
      event: 'CDL Major 1 2025',
      year: 2025,
      tournament_id: 1,
      matches: [
        { match_id: 10, date: '2025-03-05T18:00:00Z', opponent: 'Atlanta FaZe',     opponent_abbr: 'ATL', result: 'W 3:1', kd: 1.50, kills: 24, deaths: 16 },
        { match_id: 9,  date: '2025-03-04T18:00:00Z', opponent: 'Boston Breach',    opponent_abbr: 'BOS', result: 'L 1:3', kd: 0.80, kills: 16, deaths: 20 },
        { match_id: 8,  date: '2025-03-03T18:00:00Z', opponent: 'Miami Heretics',   opponent_abbr: 'MIA', result: 'W 3:2', kd: 1.10, kills: 22, deaths: 20 },
      ],
    },
    {
      event: 'CDL Major 2 2025',
      year: 2025,
      tournament_id: 2,
      matches: [
        { match_id: 7,  date: '2025-02-10T18:00:00Z', opponent: 'Toronto Ultra',         opponent_abbr: 'TOR', result: 'W 3:0', kd: 2.00, kills: 18, deaths: 9  },
        { match_id: 6,  date: '2025-02-09T18:00:00Z', opponent: 'Vancouver Surge',       opponent_abbr: 'VAN', result: 'L 2:3', kd: 0.90, kills: 18, deaths: 20 },
        { match_id: 5,  date: '2025-02-08T18:00:00Z', opponent: 'Los Angeles Thieves',   opponent_abbr: 'LAT', result: 'W 3:1', kd: 1.40, kills: 21, deaths: 15 },
      ],
    },
  ],
  total: 6,
}

export const playerMatchesEmptyFixture: PlayerMatchesResponse = {
  player_id: 1,
  events: [],
  total: 0,
}

// Null kd — documents the bug where the frontend renders "0.00" instead of "—"
export const playerMatchesNullKDFixture: PlayerMatchesResponse = {
  player_id: 2,
  events: [
    {
      event: 'CDL Major 1 2025',
      year: 2025,
      tournament_id: 1,
      matches: [
        { match_id: 1, date: '2025-03-05T18:00:00Z', opponent: 'Atlanta FaZe', opponent_abbr: 'ATL', result: 'W 3:1', kd: null, kills: 0, deaths: 0 },
      ],
    },
  ],
  total: 1,
}
