import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import type { PlayerKDStatsData, Player } from '../types';
import { playerApi } from '../services/api';

const PlayerKDStats: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [player, setPlayer] = useState<Player | null>(null);
  const [kdStats, setKdStats] = useState<PlayerKDStatsData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPlayerData = async () => {
      if (!id) return;
      
      try {
        setLoading(true);
        const [playerData, statsData] = await Promise.all([
          playerApi.getPlayer(parseInt(id)),
          playerApi.getPlayerKDStats(parseInt(id))
        ]);
        setPlayer(playerData);
        setKdStats(statsData);
      } catch (err) {
        setError('Failed to fetch player KD statistics');
        console.error('Error fetching player KD stats:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchPlayerData();
  }, [id]);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <div className="text-red-500 text-xl mb-4">{error}</div>
        <Link to="/players" className="btn-primary">
          Back to Players
        </Link>
      </div>
    );
  }

  if (!player || !kdStats) {
    return (
      <div className="text-center py-8">
        <div className="text-gray-400 text-xl mb-4">Player not found</div>
        <Link to="/players" className="btn-primary">
          Back to Players
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">{player.gamertag}</h1>
          <p className="text-gray-400">KD Statistics</p>
        </div>
        <Link to="/players" className="btn-secondary">
          Back to Players
        </Link>
      </div>

      {/* Player Info */}
      <div className="card">
        <h2 className="text-xl font-semibold text-white mb-4">Player Information</h2>
        <div className="grid md:grid-cols-2 gap-4">
          <div>
            <span className="text-gray-400">Name:</span>
            <span className="text-white ml-2">
              {[player.first_name, player.last_name].filter(Boolean).join(' ') || 'N/A'}
            </span>
          </div>
          <div>
            <span className="text-gray-400">Country:</span>
            <span className="text-white ml-2">{player.country || 'N/A'}</span>
          </div>
          <div>
            <span className="text-gray-400">Role:</span>
            <span className="text-white ml-2">{player.role || 'N/A'}</span>
          </div>
          <div>
            <span className="text-gray-400">Status:</span>
            <span className={`ml-2 ${player.is_active ? 'text-green-400' : 'text-red-400'}`}>
              {player.is_active ? 'Active' : 'Inactive'}
            </span>
          </div>
        </div>
      </div>

      {/* KD Statistics Overview */}
      <div className="grid md:grid-cols-4 gap-6">
        <div className="card text-center">
          <div className="text-3xl font-bold text-green-400">
            {kdStats.avg_kd.toFixed(2)}
          </div>
          <div className="text-gray-400 text-sm">Average KD</div>
        </div>
        <div className="card text-center">
          <div className="text-3xl font-bold text-blue-400">
            {kdStats.avg_kda.toFixed(2)}
          </div>
          <div className="text-gray-400 text-sm">Average KDA</div>
        </div>
        <div className="card text-center">
          <div className="text-3xl font-bold text-purple-400">
            {kdStats.total_kills}
          </div>
          <div className="text-gray-400 text-sm">Total Kills</div>
        </div>
        <div className="card text-center">
          <div className="text-3xl font-bold text-red-400">
            {kdStats.total_deaths}
          </div>
          <div className="text-gray-400 text-sm">Total Deaths</div>
        </div>
      </div>

      {/* Detailed Stats */}
      <div className="grid md:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="text-lg font-semibold text-white mb-4">Performance Summary</h3>
          <div className="space-y-3">
            <div className="flex justify-between">
              <span className="text-gray-400">Total Matches:</span>
              <span className="text-white font-medium">{kdStats.total_matches}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Total Maps:</span>
              <span className="text-white font-medium">{kdStats.total_maps}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Total Assists:</span>
              <span className="text-white font-medium">{kdStats.total_assists}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Average ADR:</span>
              <span className="text-white font-medium">{kdStats.avg_adr.toFixed(1)}</span>
            </div>
          </div>
        </div>

        <div className="card">
          <h3 className="text-lg font-semibold text-white mb-4">KD Performance</h3>
          <div className="space-y-3">
            <div className="flex justify-between">
              <span className="text-gray-400">KD Ratio:</span>
              <span className={`font-bold ${
                kdStats.avg_kd >= 1.5 ? 'text-green-400' :
                kdStats.avg_kd >= 1.0 ? 'text-yellow-400' :
                'text-red-400'
              }`}>
                {kdStats.avg_kd.toFixed(2)}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">KDA Ratio:</span>
              <span className={`font-bold ${
                kdStats.avg_kda >= 2.0 ? 'text-green-400' :
                kdStats.avg_kda >= 1.5 ? 'text-yellow-400' :
                'text-red-400'
              }`}>
                {kdStats.avg_kda.toFixed(2)}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Kill/Death Ratio:</span>
              <span className="text-white font-medium">
                {kdStats.total_kills}:{kdStats.total_deaths}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Tournament Statistics */}
      {kdStats.tournament_stats && kdStats.tournament_stats.length > 0 && (
        <div className="card">
          <h3 className="text-lg font-semibold text-white mb-4">Tournament Statistics</h3>
          <div className="space-y-4">
            {kdStats.tournament_stats.map((tournament) => (
              <div key={tournament.tournament_id} className="border border-gray-700 rounded-lg p-4">
                <div className="flex justify-between items-center mb-3">
                  <h4 className="text-white font-semibold">
                    {tournament.tournament_name}
                  </h4>
                  <div className="flex space-x-4 text-sm">
                    <span className="text-gray-400">
                      {tournament.matches} matches
                    </span>
                    <span className="text-gray-400">
                      {tournament.maps_played} maps
                    </span>
                  </div>
                </div>
                
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-3">
                  <div className="text-center">
                    <div className="text-green-400 font-bold text-lg">
                      {tournament.kills}
                    </div>
                    <div className="text-gray-400 text-xs">Kills</div>
                  </div>
                  <div className="text-center">
                    <div className="text-red-400 font-bold text-lg">
                      {tournament.deaths}
                    </div>
                    <div className="text-gray-400 text-xs">Deaths</div>
                  </div>
                  <div className="text-center">
                    <div className="text-blue-400 font-bold text-lg">
                      {tournament.assists}
                    </div>
                    <div className="text-gray-400 text-xs">Assists</div>
                  </div>
                  <div className="text-center">
                    <div className={`font-bold text-lg ${
                      tournament.kd_ratio >= 1.5 ? 'text-green-400' :
                      tournament.kd_ratio >= 1.0 ? 'text-yellow-400' :
                      'text-red-400'
                    }`}>
                      {tournament.kd_ratio.toFixed(2)}
                    </div>
                    <div className="text-gray-400 text-xs">KD Ratio</div>
                  </div>
                </div>
                
                <div className="flex justify-between items-center">
                  <div className="text-sm">
                    <span className="text-gray-400">KDA Ratio:</span>
                    <span className={`ml-2 font-medium ${
                      tournament.kda_ratio >= 2.0 ? 'text-green-400' :
                      tournament.kda_ratio >= 1.5 ? 'text-yellow-400' :
                      'text-red-400'
                    }`}>
                      {tournament.kda_ratio.toFixed(2)}
                    </span>
                  </div>
                  <div className="text-sm">
                    <span className="text-gray-400">Kill/Death:</span>
                    <span className="text-white ml-2">
                      {tournament.kills}:{tournament.deaths}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Match History */}
      {kdStats.match_stats.length > 0 && (
        <div className="card">
          <h3 className="text-lg font-semibold text-white mb-4">Recent Matches</h3>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-700">
                  <th className="text-left py-2 px-4 text-gray-300 text-sm">Match</th>
                  <th className="text-right py-2 px-4 text-gray-300 text-sm">Maps</th>
                  <th className="text-right py-2 px-4 text-gray-300 text-sm">Kills</th>
                  <th className="text-right py-2 px-4 text-gray-300 text-sm">Deaths</th>
                  <th className="text-right py-2 px-4 text-gray-300 text-sm">Assists</th>
                  <th className="text-right py-2 px-4 text-gray-300 text-sm">KD</th>
                  <th className="text-right py-2 px-4 text-gray-300 text-sm">KDA</th>
                </tr>
              </thead>
              <tbody>
                {kdStats.match_stats.slice(0, 10).map((match) => (
                  <tr key={match.id} className="border-b border-gray-800">
                    <td className="py-2 px-4 text-white text-sm">
                      Match #{match.match_id}
                    </td>
                    <td className="py-2 px-4 text-right text-gray-300 text-sm">
                      {match.maps_played}
                    </td>
                    <td className="py-2 px-4 text-right text-green-400 text-sm font-medium">
                      {match.total_kills}
                    </td>
                    <td className="py-2 px-4 text-right text-red-400 text-sm font-medium">
                      {match.total_deaths}
                    </td>
                    <td className="py-2 px-4 text-right text-blue-400 text-sm font-medium">
                      {match.total_assists}
                    </td>
                    <td className="py-2 px-4 text-right text-sm">
                      <span className={`font-medium ${
                        match.kd_ratio >= 1.5 ? 'text-green-400' :
                        match.kd_ratio >= 1.0 ? 'text-yellow-400' :
                        'text-red-400'
                      }`}>
                        {match.kd_ratio.toFixed(2)}
                      </span>
                    </td>
                    <td className="py-2 px-4 text-right text-sm">
                      <span className={`font-medium ${
                        match.kda_ratio >= 2.0 ? 'text-green-400' :
                        match.kda_ratio >= 1.5 ? 'text-yellow-400' :
                        'text-red-400'
                      }`}>
                        {match.kda_ratio.toFixed(2)}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {kdStats.match_stats.length === 0 && kdStats.tournament_stats.length === 0 && (
        <div className="card text-center py-8">
          <div className="text-gray-400 text-lg mb-2">No match data available</div>
          <p className="text-gray-500">This player hasn't played any matches yet.</p>
        </div>
      )}
    </div>
  );
};

export default PlayerKDStats; 