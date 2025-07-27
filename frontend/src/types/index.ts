// TypeScript interfaces matching the Go backend models

export interface Season {
  id: number;
  name: string;
  game_title: string;
  start_date: string;
  end_date?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Team {
  id: number;
  name: string;
  abbreviation: string;
  city?: string;
  logo_url?: string;
  primary_color?: string;
  secondary_color?: string;
  founded_date?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Player {
  id: number;
  gamertag: string;
  first_name?: string;
  last_name?: string;
  country?: string;
  birthdate?: string;
  role?: string;
  is_active: boolean;
  liquipedia_url?: string;
  twitter_handle?: string;
  avatar_url?: string;
  created_at: string;
  updated_at: string;
}

export interface TeamRoster {
  id: number;
  team_id: number;
  player_id: number;
  season_id: number;
  role?: string;
  start_date: string;
  end_date?: string;
  is_starter: boolean;
  created_at: string;
  team?: Team;
  player?: Player;
  season?: Season;
}

export interface PlayerMatchStats {
  id: number;
  match_id: number;
  player_id: number;
  team_id: number;
  maps_played: number;
  total_kills: number;
  total_deaths: number;
  total_assists: number;
  total_damage: number;
  kd_ratio: number;
  kda_ratio: number;
  adr: number;
  created_at: string;
  match?: Match;
  player?: Player;
  team?: Team;
}

export interface TeamTournamentStats {
  id: number;
  tournament_id: number;
  team_id: number;
  placement?: number;
  matches_played: number;
  matches_won: number;
  matches_lost: number;
  maps_played: number;
  maps_won: number;
  maps_lost: number;
  prize_money: number;
  created_at: string;
  tournament?: Tournament;
  team?: Team;
}

export interface Tournament {
  id: number;
  season_id: number;
  name: string;
  tournament_type?: string;
  start_date: string;
  end_date?: string;
  prize_pool?: number;
  location?: string;
  tournament_format?: string;
  liquipedia_url?: string;
  created_at: string;
  updated_at: string;
  season?: Season;
}

export interface Match {
  id: number;
  tournament_id: number;
  team1_id: number;
  team2_id: number;
  match_date: string;
  match_type?: string;
  format?: string;
  team1_score: number;
  team2_score: number;
  winner_id?: number;
  duration_minutes?: number;
  vod_url?: string;
  liquipedia_url?: string;
  created_at: string;
  updated_at: string;
  tournament?: Tournament;
  team1?: Team;
  team2?: Team;
  winner?: Team;
}

// API Response types
export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

// KD Stats interfaces
export interface TournamentKDStats {
  tournament_id: number;
  tournament_name: string;
  kills: number;
  deaths: number;
  assists: number;
  maps_played: number;
  matches: number;
  kd_ratio: number;
  kda_ratio: number;
}

export interface PlayerKDStatsData {
  player_id: number;
  total_matches: number;
  total_maps: number;
  total_kills: number;
  total_deaths: number;
  total_assists: number;
  avg_kd: number;
  avg_kda: number;
  avg_adr: number;
  tournament_stats: TournamentKDStats[];
  match_stats: PlayerMatchStats[];
}

export interface TopKDPlayer {
  player_id: number;
  gamertag: string;
  team_name: string;
  team_abbreviation: string;
  avg_kd: number;
  avg_kda: number;
  matches_played: number;
} 