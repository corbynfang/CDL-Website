import React, { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { statsApi } from '../services/api';
import PlayerAvatar from './PlayerAvatar';

const MAJOR_LABELS = {
  1: 'Major 1',
  2: 'Major 2',
  3: 'Major 3',
  4: 'Major 4',
  5: 'Champs',
  7: 'EWC 2025',
};

// Players to exclude from Black Ops 6 season
const EXCLUDED_PLAYERS = ['Vikul', 'accuracy', 'Crimsix'];

const KDStats: React.FC = () => {
  const [players, setPlayers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  const fetchAllKD = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Force fresh data by adding cache-busting parameters
      const response = await statsApi.getAllPlayersKDStats();
      
      // Handle new response format with timestamp
      const data = Array.isArray(response) ? response : (response as any).players || response || [];
      
      // Filter out excluded players and ensure only players with tournament stats are shown
      const filteredPlayers = data
        .filter((player: any) => !EXCLUDED_PLAYERS.includes(player.gamertag))
        .filter((player: any) => {
          // Only include players who have at least one tournament stat
          if (!player.majors) return false;
          
          // Check if player has any valid tournament stats (including 0.0 for players with no kills/deaths)
          const hasTournamentStats = Object.values(player.majors).some((kd: any) => 
            kd !== null && kd !== undefined
          );
          
          return hasTournamentStats;
        })
        .sort((a: any, b: any) => {
          // Sort by season KD descending (highest to lowest)
          const aKD = a.season_kd || 0;
          const bKD = b.season_kd || 0;
          
          if (bKD !== aKD) {
            return bKD - aKD;
          }
          
          // If same KD, sort by gamertag alphabetically
          return a.gamertag.localeCompare(b.gamertag);
        });
      
      setPlayers(filteredPlayers);
    } catch (err) {
      setError('Failed to fetch KD statistics');
      console.error('Error fetching KD stats:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchAllKD();
  }, [fetchAllKD, refreshKey]);

  const handleRefresh = () => {
    setRefreshKey(prev => prev + 1);
  };

  const getRankColor = (index: number) => {
    if (index === 0) return 'text-yellow-400'; // Gold for #1
    if (index === 1) return 'text-gray-300'; // Silver for #2
    if (index === 2) return 'text-amber-600'; // Bronze for #3
    return 'text-gray-400'; // Default for others
  };

  const getKDColor = (kd: number) => {
    if (kd >= 1.2) return 'text-green-400'; // Excellent
    if (kd >= 1.0) return 'text-blue-400'; // Good
    if (kd >= 0.9) return 'text-yellow-400'; // Average
    return 'text-red-400'; // Below average
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-none h-12 w-12 border-b-2 border-white"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <div className="text-red-500 text-xl mb-4">{error}</div>
        <button onClick={handleRefresh} className="btn-primary">TRY AGAIN</button>
      </div>
    );
  }

  return (
    <div className="space-y-6 sm:space-y-8">
      {/* Hero Section */}
      <div className="text-center py-8 sm:py-16 bg-black border border-gray-800 mb-6 sm:mb-8">
        <h1 className="text-2xl sm:text-4xl lg:text-5xl font-black text-white mb-4 sm:mb-6 tracking-tight">KD STATISTICS</h1>
        <p className="text-sm sm:text-base lg:text-lg text-gray-400 max-w-3xl mx-auto leading-relaxed uppercase tracking-wider px-4">
          COMPREHENSIVE KILL/DEATH RATIOS FOR ALL CDL PLAYERS ACROSS SEASON AND MAJOR TOURNAMENTS
        </p>
        <div className="mt-6 sm:mt-8 text-gray-400 text-sm sm:text-lg uppercase tracking-wider">
          {players.length} PLAYERS • RANKED BY SEASON KD
        </div>
        <button 
          onClick={handleRefresh}
          className="mt-4 px-4 sm:px-6 py-2 bg-red-600 hover:bg-red-700 text-white font-bold uppercase tracking-wider transition-colors text-sm sm:text-base"
        >
          REFRESH DATA
        </button>
      </div>

      {/* Mobile Card View */}
      <div className="block lg:hidden">
        <div className="space-y-3 sm:space-y-4">
          {players.map((player, index) => (
            <div key={player.player_id} className="card p-4">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center space-x-3">
                  <div className={`font-bold text-lg ${getRankColor(index)}`}>
                    #{index + 1}
                  </div>
                  <PlayerAvatar 
                    player={{
                      id: player.player_id,
                      gamertag: player.gamertag,
                      avatar_url: player.avatar_url || `/assets/avatars/${player.gamertag}.webp`
                    }} 
                    size="sm" 
                  />
                  <Link to={`/players/${player.player_id}`} className="text-white hover:text-red-500 font-medium uppercase tracking-wider text-sm sm:text-base">
                    {player.gamertag}
                  </Link>
                </div>
                <div className="text-gray-300 uppercase tracking-wider text-xs sm:text-sm">{player.team_abbr}</div>
              </div>
              
              <div className="grid grid-cols-2 gap-3 sm:gap-4 text-sm">
                <div>
                  <div className="text-gray-400 text-xs uppercase tracking-wider">SEASON KD</div>
                  <div className={`font-bold text-base sm:text-lg ${getKDColor(player.season_kd)}`}>
                    {player.season_kd ? player.season_kd.toFixed(3) : '-'}
                  </div>
                </div>
                <div>
                  <div className="text-gray-400 text-xs uppercase tracking-wider">KD +/-</div>
                  <div className={`font-bold text-base sm:text-lg ${player.season_kd_plus_minus > 0 ? 'text-green-400' : player.season_kd_plus_minus < 0 ? 'text-red-400' : 'text-gray-300'}`}>
                    {player.season_kd_plus_minus ? player.season_kd_plus_minus.toFixed(3) : '-'}
                  </div>
                </div>
              </div>

              {/* Tournament KDs - Horizontal scroll */}
              <div className="mt-4">
                <div className="text-gray-400 text-xs uppercase tracking-wider mb-2">TOURNAMENT KDS</div>
                <div className="flex space-x-4 overflow-x-auto pb-2">
                  {Object.entries(MAJOR_LABELS).map(([id, label]) => (
                    <div key={id} className="flex-shrink-0 text-center">
                      <div className="text-gray-400 text-xs uppercase tracking-wider">{label}</div>
                      <div className="font-bold text-white text-sm">
                        {player.majors && player.majors[id] && player.majors[id] > 0 ? player.majors[id].toFixed(3) : '-'}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Desktop Table View */}
      <div className="hidden lg:block card overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-gray-800">
              <th className="py-4 px-6 text-left text-white font-bold uppercase tracking-wider">RANK</th>
              <th className="py-4 px-6 text-left text-white font-bold uppercase tracking-wider">PLAYER</th>
              <th className="py-4 px-6 text-left text-white font-bold uppercase tracking-wider">TEAM</th>
              <th className="py-4 px-6 text-right text-white font-bold uppercase tracking-wider">SEASON KD</th>
              <th className="py-4 px-6 text-right text-white font-bold uppercase tracking-wider">KD +/-</th>
              {Object.entries(MAJOR_LABELS).map(([id, label]) => (
                <th key={id} className="py-4 px-6 text-right text-white font-bold uppercase tracking-wider">{label.toUpperCase()}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {players.map((player, index) => (
              <tr key={player.player_id} className="border-b border-gray-800 hover:bg-gray-900">
                <td className="py-4 px-6">
                  <div className={`font-bold text-lg ${getRankColor(index)}`}>
                    #{index + 1}
                  </div>
                </td>
                <td className="py-4 px-6">
                  <div className="flex items-center space-x-3">
                    <div className="flex-shrink-0">
                      <PlayerAvatar 
                        player={{
                          id: player.player_id,
                          gamertag: player.gamertag,
                          avatar_url: player.avatar_url || `/assets/avatars/${player.gamertag}.webp`
                        }} 
                        size="sm" 
                      />
                    </div>
                    <Link to={`/players/${player.player_id}`} className="text-white hover:text-red-500 font-medium uppercase tracking-wider">
                      {player.gamertag}
                    </Link>
                  </div>
                </td>
                <td className="py-4 px-6 text-gray-300 uppercase tracking-wider">{player.team_abbr}</td>
                <td className={`py-4 px-6 text-right font-bold text-lg ${getKDColor(player.season_kd)}`}>
                  {player.season_kd ? player.season_kd.toFixed(3) : '-'}
                </td>
                <td className={`py-4 px-6 text-right font-bold ${player.season_kd_plus_minus > 0 ? 'text-green-400' : player.season_kd_plus_minus < 0 ? 'text-red-400' : 'text-gray-300'}`}>
                  {player.season_kd_plus_minus ? player.season_kd_plus_minus.toFixed(3) : '-'}
                </td>
                {Object.keys(MAJOR_LABELS).map((id) => (
                  <td key={id} className="py-4 px-6 text-right font-bold text-white">
                    {player.majors && player.majors[id] && player.majors[id] > 0 ? player.majors[id].toFixed(3) : '-'}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default KDStats; 