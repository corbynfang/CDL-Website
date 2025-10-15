import React, { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import PlayerAvatar from './PlayerAvatar';
import LoadingSkeleton, { ErrorDisplay } from './LoadingSkeleton';

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
  const [refreshKey, setRefreshKey] = useState(0);

  // Use the new API hook with cache-busting via refreshKey
  const { data: response, loading, error, refetch } = useApi<any>(
    `/api/v1/stats/all-kd-by-tournament?_refresh=${refreshKey}`,
    { retries: 3, retryDelay: 1000 }
  );

  const handleRefresh = () => {
    setRefreshKey(prev => prev + 1);
    refetch();
  };

  // Process and filter players
  const players = useMemo(() => {
    if (!response) return [];
    
    const data = Array.isArray(response) ? response : (response as any).players || response || [];
    
    const filteredPlayers = data
      .filter((player: any) => !EXCLUDED_PLAYERS.includes(player.gamertag))
      .filter((player: any) => player.majors && Object.keys(player.majors).length > 0)
      .sort((a: any, b: any) => {
        const aKD = a.season_kd || 0;
        const bKD = b.season_kd || 0;
        return bKD !== aKD ? bKD - aKD : a.gamertag.localeCompare(b.gamertag);
      });
    
    return filteredPlayers;
  }, [response]);

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
    return <LoadingSkeleton variant="table" count={10} />;
  }

  if (error) {
    return <ErrorDisplay message={error} onRetry={handleRefresh} />;
  }

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 py-8">
      <div className="mb-12">
        <div className="flex items-center justify-between mb-3">
          <h1 className="text-3xl sm:text-4xl font-bold text-white">k/d statistics</h1>
          <button 
            onClick={handleRefresh}
            className="text-sm px-4 py-2 rounded-lg transition-all"
            style={{ backgroundColor: '#1a1a1a', color: '#ffffff' }}
          >
            refresh
          </button>
        </div>
        <p className="text-lg" style={{ color: '#a3a3a3' }}>
          {players.length} players ranked by season k/d
        </p>
      </div>

      <div className="block lg:hidden">
        <div className="space-y-3 sm:space-y-4">
          {players.map((player: any, index: number) => (
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
                        {player.majors && player.majors[id] && typeof player.majors[id] === 'object' && player.majors[id].kd_ratio ? player.majors[id].kd_ratio.toFixed(3) : '-'}
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
            {players.map((player: any, index: number) => (
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
                    {player.majors && player.majors[id] && typeof player.majors[id] === 'object' && player.majors[id].kd_ratio ? player.majors[id].kd_ratio.toFixed(3) : '-'}
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