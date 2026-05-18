import type { Tournament } from '../../types'

const BASE_SEASON = {
  id: 5,
  name: 'Black Ops 6 2025',
  game_title: 'Call of Duty: Black Ops 6',
  game_code: 'BO6',
  start_date: '2024-10-01T00:00:00Z',
  is_active: true,
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

/** A fully-completed LAN major. */
export const completedMajor: Tournament = {
  id: 1,
  season_id: 5,
  name: 'CDL Major 2 2025',
  slug: 'cdl-major-2-tournament-2025',
  tournament_type: 'major_tournament',
  start_date: '2020-03-01T00:00:00Z',
  end_date:   '2020-03-04T00:00:00Z',
  prize_pool: 375000,
  location: 'Allen, Texas',
  country: 'USA',
  is_lan: true,
  logo_url: '',
  tournament_format: 'Double Elimination',
  liquipedia_url: 'https://liquipedia.net/callofduty/test',
  breaking_point_url: 'https://breakingpoint.gg/test',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
  season: BASE_SEASON,
}

/** A future upcoming event (no prize pool, no location). */
export const upcomingQualifier: Tournament = {
  id: 2,
  season_id: 5,
  name: 'CDL Qualifier 1 2025',
  slug: 'cdl-qualifier-1-2025',
  tournament_type: 'qualifier',
  start_date: '2099-09-01T00:00:00Z',
  end_date:   '2099-09-03T00:00:00Z',
  prize_pool: null,
  location: '',
  country: '',
  is_lan: false,
  logo_url: '',
  tournament_format: '',
  liquipedia_url: '',
  breaking_point_url: '',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
  season: BASE_SEASON,
}

/** A currently-live event (started in the past, ends in the far future). */
export const liveMajor: Tournament = {
  id: 3,
  season_id: 5,
  name: 'CDL Major 3 2025',
  slug: 'cdl-major-3-tournament-2025',
  tournament_type: 'major_tournament',
  start_date: '2020-01-01T00:00:00Z',
  end_date:   '2099-12-31T00:00:00Z',
  prize_pool: 375000,
  location: 'Boca Raton, Florida',
  country: 'USA',
  is_lan: true,
  logo_url: '',
  tournament_format: 'Double Elimination',
  liquipedia_url: '',
  breaking_point_url: '',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
  season: BASE_SEASON,
}

/** A championship with no season attached (tests fallback rendering). */
export const championshipNoSeason: Tournament = {
  id: 4,
  season_id: 5,
  name: 'CDL Championship 2025',
  slug: 'cdl-league-championship-2025',
  tournament_type: 'championship',
  start_date: '2020-08-01T00:00:00Z',
  end_date:   '2020-08-05T00:00:00Z',
  prize_pool: 2000000,
  location: 'Kitchener, Ontario',
  country: 'CAN',
  is_lan: true,
  logo_url: '',
  tournament_format: 'Double Elimination',
  liquipedia_url: '',
  breaking_point_url: '',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

/** International major (hasBracket = true). */
export const internationalMajor: Tournament = {
  id: 5,
  season_id: 5,
  name: 'Esports World Cup 2025',
  slug: 'esports-world-cup-2025',
  tournament_type: 'international_major',
  start_date: '2020-07-01T00:00:00Z',
  end_date:   '2020-07-06T00:00:00Z',
  prize_pool: 1800000,
  location: 'Riyadh',
  country: 'SAU',
  is_lan: true,
  logo_url: '',
  tournament_format: '',
  liquipedia_url: '',
  breaking_point_url: '',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
  season: BASE_SEASON,
}
