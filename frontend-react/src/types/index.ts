// TypeScript interfaces matching the Go backend models

export interface AppUser {
  id: number;
  supabase_uid: string;
  username: string;
  created_at: string;
  updated_at: string;
}

export interface ThreadPost {
  id: number;
  thread_id: number;
  user_id: number;
  body: string;
  edited: boolean;
  created_at: string;
  updated_at: string;
  user: AppUser;
}

export interface ThreadResponse {
  thread_id: number;
  data: ThreadPost[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface Season {
  id: number;
  name: string;
  game_title: string;
  game_code: string; // BO6 | CW | MW2 | MW3 | VG
  start_date: string;
  end_date?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Franchise {
  id: number;
  franchise_key: string;
  name: string;
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
  franchise_id?: number;
  game_code?: string;
  is_cdl_franchise: boolean;
  team_classification?: string;
  do_not_merge?: boolean;
  valid_from?: string;
  valid_to?: string;
  needs_manual_review?: boolean;
  franchise?: Franchise;
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
  source_profile_url?: string;
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
  slug: string;
  tournament_type: string;
  start_date: string;
  end_date?: string | null;
  prize_pool?: number | null;
  location?: string;
  country?: string;
  is_lan: boolean;
  logo_url?: string;
  tournament_format?: string;
  source_event_url?: string;
  created_at: string;
  updated_at: string;
  season?: Season;
}

export interface TournamentDetail {
  tournament: Tournament;
  team_count: number;
  event_format?: string;
}

export interface TournamentTeam extends Team {
  placement?: number | null;
  matches_won: number;
  matches_lost: number;
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

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: PaginationMeta;
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

// Matches KDRow in internal/store/stats.go
export interface TopKDPlayer {
  player_id: number;
  gamertag: string;
  avatar_url: string;
  team_abbr: string;
  season_kills: number;
  season_deaths: number;
  season_assists: number;
  season_kd: number;
}

// Returned by GET /stats/all-kd-by-tournament — same as TopKDPlayer plus derived field
export interface AllKDPlayer extends TopKDPlayer {
  season_kd_plus_minus: number;
}

export interface PlayerTransfer {
  id: number;
  player_id: number;
  from_team_id?: number;
  to_team_id?: number;
  transfer_date: string;
  transfer_type: string;
  role: string;
  game_code?: string;
  season: string;
  description: string;
  raw_from_team_name?: string;
  raw_to_team_name?: string;
  created_at: string;
  player?: Player;
  from_team?: Team;
  to_team?: Team;
}

export interface MatchMap {
  id: number;
  match_id: number;
  map_number: number;
  map_name: string;
  mode: string;
  score_1: number;
  score_2: number;
  winner_id?: number;
  played: boolean;
  duration_sec: number;
  source: string;
  winner?: Team;
}

export interface PlayerTournamentStats {
  id: number;
  player_id: number;
  team_id: number;
  tournament_id: number;
  total_kills: number;
  total_deaths: number;
  total_assists: number;
  total_damage: number;
  kd_ratio: number;
  kda_ratio: number;
  rank?: number | null;
  overall_maps: number;
  overall_plus_minus: number;
  player?: Player;
  team?: Team;
}

export interface PlayerMapStats {
  id: number;
  match_id: number;
  map_number: number;
  player_id: number;
  team_id: number;
  kills: number;
  deaths: number;
  kd_ratio: number;
  damage: number;
  assists: number;
  bp_rating: number;
  hill_time: number;
  snd_rounds: number;
  plant_count: number;
  defuse_count: number;
  first_blood_count: number;
  first_death_count: number;
  source: string;
}

export interface PlayerKDTournamentEntry {
  tournament_id: number;
  tournament_name: string;
  kills: number;
  deaths: number;
  assists: number;
  kd_ratio: number;
  maps_played: number;
}

export interface PlayerKDResponse {
  player_id: number;
  gamertag: string;
  avatar_url: string;
  total_kills: number;
  total_deaths: number;
  total_assists: number;
  avg_kd: number;
  hp_kd_ratio: number;
  snd_kd_ratio: number;
  control_kd_ratio: number;
  tournament_stats: PlayerKDTournamentEntry[];
}

export interface MatchHistoryResult {
  match_id: number;
  date: string;
  opponent: string;
  opponent_abbr: string;
  result: string;
  kd: number | null;
  kills: number;
  deaths: number;
}

export interface MatchHistoryEvent {
  event: string;
  year: number;
  tournament_id: number;
  matches: MatchHistoryResult[];
}

export interface PlayerMatchHistory {
  player_id: number;
  events: MatchHistoryEvent[];
  total: number;
}

export interface PlayerEraStats {
  team_id: number;
  team_name: string;
  game_code: string;
  season_name: string;
  matches: number;
  maps: number;
  kills: number;
  deaths: number;
  kd: number;
}

export interface PlayerFranchiseEntry {
  franchise_key: string;
  franchise_name: string;
  eras: PlayerEraStats[];
  total_matches: number;
  total_maps: number;
  total_kills: number;
  total_deaths: number;
  career_kd: number;
}

export interface PlayerCareerResponse {
  player_id: number;
  gamertag: string;
  franchises: PlayerFranchiseEntry[];
}
