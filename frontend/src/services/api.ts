import axios from 'axios';
import type { AxiosResponse } from 'axios';
import type { 
  Team, 
  Player, 
  PlayerMatchStats, 
  TeamTournamentStats,
  PlayerKDStatsData,
  TopKDPlayer,
  PlayerTransfer,
  ApiResponse 
} from '../types';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'https://cdlytics.me/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
    'Cache-Control': 'no-cache, no-store, must-revalidate, max-age=0',
    'Pragma': 'no-cache',
    'Expires': '0',
    'X-Requested-With': 'XMLHttpRequest',
  },
});

// Request interceptor for logging and aggressive cache-busting
api.interceptors.request.use(
  (config) => {
    // console.log(`Making ${config.method?.toUpperCase()} request to ${config.url}`);
    
    // Add multiple cache-busting parameters for GET requests
    if (config.method === 'get') {
      const timestamp = Date.now();
      const random = Math.random().toString(36).substring(7);
      
      config.params = {
        ...config.params,
        _t: timestamp, // Cache-busting timestamp
        _r: random, // Random string
        _v: '1.0', // Version parameter
      };
    }
    
    // Add cache-busting headers
    if (config.headers) {
      config.headers['Cache-Control'] = 'no-cache, no-store, must-revalidate, max-age=0';
      config.headers['Pragma'] = 'no-cache';
      config.headers['Expires'] = '0';
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

// Team API functions
export const teamApi = {
  // Get all teams
  getTeams: async (): Promise<Team[]> => {
    const response: AxiosResponse<Team[]> = await api.get('/teams');
    return response.data;
  },

  // Get team by ID
  getTeam: async (id: number): Promise<Team> => {
    const response: AxiosResponse<Team> = await api.get(`/teams/${id}`);
    return response.data;
  },

  // Create new team
  createTeam: async (team: Omit<Team, 'id' | 'created_at' | 'updated_at'>): Promise<Team> => {
    const response: AxiosResponse<Team> = await api.post('/teams', team);
    return response.data;
  },

  // Get team players
  getTeamPlayers: async (teamId: number): Promise<Player[]> => {
    const response: AxiosResponse<Player[]> = await api.get(`/teams/${teamId}/players`);
    return response.data;
  },

  // Get team stats
  getTeamStats: async (teamId: number): Promise<TeamTournamentStats[]> => {
    const response: AxiosResponse<TeamTournamentStats[]> = await api.get(`/teams/${teamId}/stats`);
    return response.data;
  },

  // Create team stats
  createTeamStats: async (teamId: number, stats: Omit<TeamTournamentStats, 'id' | 'created_at'>): Promise<TeamTournamentStats> => {
    const response: AxiosResponse<TeamTournamentStats> = await api.post(`/teams/${teamId}/stats`, stats);
    return response.data;
  },
};

// Player API functions
export const playerApi = {
  // Get all players
  getPlayers: async (): Promise<Player[]> => {
    const response: AxiosResponse<Player[]> = await api.get('/players');
    return response.data;
  },

  // Get player by ID
  getPlayer: async (id: number): Promise<Player> => {
    const response: AxiosResponse<Player> = await api.get(`/players/${id}`);
    return response.data;
  },

  // Create new player
  createPlayer: async (player: Omit<Player, 'id' | 'created_at' | 'updated_at'>): Promise<Player> => {
    const response: AxiosResponse<Player> = await api.post('/players', player);
    return response.data;
  },

  // Get player stats
  getPlayerStats: async (playerId: number): Promise<PlayerMatchStats[]> => {
    const response: AxiosResponse<PlayerMatchStats[]> = await api.get(`/players/${playerId}/stats`);
    return response.data;
  },

  // Create player stats
  createPlayerStats: async (playerId: number, stats: Omit<PlayerMatchStats, 'id' | 'created_at'>): Promise<PlayerMatchStats> => {
    const response: AxiosResponse<PlayerMatchStats> = await api.post(`/players/${playerId}/stats`, stats);
    return response.data;
  },

  // Get player KD stats
  getPlayerKDStats: async (playerId: number): Promise<PlayerKDStatsData> => {
    const response: AxiosResponse<PlayerKDStatsData> = await api.get(`/players/${playerId}/kd`);
    return response.data;
  },
};

// Stats API functions
export const statsApi = {
  // Get top KD players
  getTopKDPlayers: async (): Promise<TopKDPlayer[]> => {
    const response: AxiosResponse<TopKDPlayer[]> = await api.get('/stats/top-kd-new');
    return response.data;
  },

  // Get all player KD stats for the season and all majors
  getAllPlayersKDStats: async (): Promise<any[]> => {
    const response: AxiosResponse<any[]> = await api.get('/stats/all-kd-by-tournament');
    return response.data;
  },
};

// Health check
export const healthApi = {
  checkHealth: async (): Promise<ApiResponse<{ status: string; message: string }>> => {
    const response: AxiosResponse<ApiResponse<{ status: string; message: string }>> = await api.get('/health');
    return response.data;
  },
};

// Transfers API functions
export const transfersApi = {
  // Get all transfers
  getTransfers: async (params?: { season?: string; team_id?: number; type?: string; player_id?: number }): Promise<PlayerTransfer[]> => {
    const response: AxiosResponse<{ transfers: PlayerTransfer[]; count: number; timestamp: number }> = await api.get('/transfers', { params });
    return response.data.transfers;
  },
};

export default api; 