import axios from "axios";
import type { AxiosResponse } from "axios";
import type {
  Team,
  Player,
  Match,
  PlayerMatchStats,
  PlayerTournamentStats,
  TeamTournamentStats,
  PlayerKDStatsData,
  TopKDPlayer,
  AllKDPlayer,
  PlayerTransfer,
  Tournament,
  TournamentDetail,
  TournamentTeam,
  PaginatedResponse,
  ApiResponse,
} from "../types";

// Bracket data shape returned by GET /tournaments/:id/bracket
export interface BracketMatch {
  id: number;
  team1_id: number;
  team2_id: number;
  team1_name: string;
  team1_abbr: string;
  team1_logo: string;
  team2_name: string;
  team2_abbr: string;
  team2_logo: string;
  team1_score: number;
  team2_score: number;
  winner_id: number | null;
  bracket_position: number;
  match_date: string;
}

export interface BracketData {
  tournament_id: number;
  tournament_name: string;
  total_matches: number;
  event_format?: string;
  bracket: Record<string, BracketMatch[]>;
  group_stage?: Record<string, BracketMatch[]>;
}

// Create axios instance with base configuration
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "https://cdlytics.com/api/v1",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error("API Error:", error.response?.data || error.message);
    return Promise.reject(error);
  },
);

// Team API functions
export const teamApi = {
  // Get all teams
  getTeams: async (): Promise<Team[]> => {
    const response: AxiosResponse<Team[]> = await api.get("/teams");
    return response.data;
  },

  // Get team by ID
  getTeam: async (id: number): Promise<Team> => {
    const response: AxiosResponse<Team> = await api.get(`/teams/${id}`);
    return response.data;
  },

  // Get team players
  getTeamPlayers: async (teamId: number): Promise<Player[]> => {
    const response: AxiosResponse<Player[]> = await api.get(
      `/teams/${teamId}/players`,
    );
    return response.data;
  },

  // Get team stats
  getTeamStats: async (teamId: number): Promise<TeamTournamentStats[]> => {
    const response: AxiosResponse<TeamTournamentStats[]> = await api.get(
      `/teams/${teamId}/stats`,
    );
    return response.data;
  },

};

// Player API functions
export const playerApi = {
  // Get a page of players (?page=1&limit=25 by default)
  getPlayers: async (
    page = 1,
    limit = 25,
  ): Promise<PaginatedResponse<Player>> => {
    const response: AxiosResponse<PaginatedResponse<Player>> = await api.get(
      "/players",
      { params: { page, limit } },
    );
    return response.data;
  },

  // Get player by ID
  getPlayer: async (id: number): Promise<Player> => {
    const response: AxiosResponse<Player> = await api.get(`/players/${id}`);
    return response.data;
  },

  // Get player stats
  getPlayerStats: async (playerId: number): Promise<PlayerMatchStats[]> => {
    const response: AxiosResponse<PlayerMatchStats[]> = await api.get(
      `/players/${playerId}/stats`,
    );
    return response.data;
  },

  // Get player KD stats
  getPlayerKDStats: async (playerId: number): Promise<PlayerKDStatsData> => {
    const response: AxiosResponse<PlayerKDStatsData> = await api.get(
      `/players/${playerId}/kd`,
    );
    return response.data;
  },
};

// Stats API functions
export const statsApi = {
  // Get top KD players — backend wraps array as { timestamp, players, count }
  getTopKDPlayers: async (): Promise<TopKDPlayer[]> => {
    const response: AxiosResponse<{
      players: TopKDPlayer[];
      count: number;
      timestamp: number;
    }> = await api.get("/players/top-kd");
    return response.data.players;
  },

  // Get all player KD stats — backend wraps array as { timestamp, players, count }
  getAllPlayersKDStats: async (): Promise<AllKDPlayer[]> => {
    const response: AxiosResponse<{
      players: AllKDPlayer[];
      count: number;
      timestamp: number;
    }> = await api.get("/stats/all-kd-by-tournament");
    return response.data.players;
  },
};

// Health check
export const healthApi = {
  checkHealth: async (): Promise<
    ApiResponse<{ status: string; message: string }>
  > => {
    const response: AxiosResponse<
      ApiResponse<{ status: string; message: string }>
    > = await api.get("/health");
    return response.data;
  },
};

// Transfers API functions
export const transfersApi = {
  // Get all transfers
  getTransfers: async (params?: {
    season?: string;
    team_id?: number;
    type?: string;
    player_id?: number;
  }): Promise<PlayerTransfer[]> => {
    const response: AxiosResponse<{
      transfers: PlayerTransfer[];
      count: number;
      timestamp: number;
    }> = await api.get("/transfers", { params });
    return response.data.transfers;
  },
};

// Events (tournaments) API
export const eventsApi = {
  getAll: async (params?: { season_id?: number }): Promise<Tournament[]> => {
    const response: AxiosResponse<Tournament[]> = await api.get(
      "/tournaments",
      { params },
    );
    return response.data;
  },
  getBySlug: async (slug: string): Promise<TournamentDetail> => {
    const response: AxiosResponse<TournamentDetail> = await api.get(
      `/tournaments/slug/${slug}`,
    );
    return response.data;
  },
  getById: async (id: number): Promise<Tournament> => {
    const response: AxiosResponse<Tournament> = await api.get(
      `/tournaments/${id}`,
    );
    return response.data;
  },
  getMatches: async (id: number): Promise<Match[]> => {
    const response: AxiosResponse<Match[]> = await api.get(
      `/tournaments/${id}/matches`,
    );
    return response.data;
  },
  getTeams: async (id: number): Promise<TournamentTeam[]> => {
    const response: AxiosResponse<TournamentTeam[]> = await api.get(
      `/tournaments/${id}/teams`,
    );
    return response.data;
  },
  getBracket: async (id: number): Promise<BracketData> => {
    const response: AxiosResponse<BracketData> = await api.get(
      `/tournaments/${id}/bracket`,
    );
    return response.data;
  },
  getStats: async (id: number): Promise<PlayerTournamentStats[]> => {
    const response: AxiosResponse<PlayerTournamentStats[]> = await api.get(
      `/tournaments/${id}/stats`,
    );
    return response.data;
  },
};

export default api;
