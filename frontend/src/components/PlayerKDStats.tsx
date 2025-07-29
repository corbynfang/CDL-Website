import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { playerApi } from '../services/api';
import PlayerAvatar from './PlayerAvatar';



const PlayerKDStats: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [playerStats, setPlayerStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPlayerStats = async () => {
      if (!id) return;
      
      try {
        setLoading(true);
        const stats = await playerApi.getPlayerKDStats(parseInt(id));
        setPlayerStats(stats);
      } catch (err) {
        setError('Failed to fetch player KD statistics');
        console.error('Error fetching player KD stats:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchPlayerStats();
  }, [id]);

  const getKDColor = (kd: number) => {
    if (kd >= 1.2) return 'text-green-400';
    if (kd >= 1.0) return 'text-blue-400';
    if (kd >= 0.9) return 'text-yellow-400';
    return 'text-red-400';
  };

  const getPlusMinusColor = (value: number) => {
    if (value > 0) return 'text-green-400';
    if (value < 0) return 'text-red-400';
    return 'text-gray-300';
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <div className="text-red-500 text-xl mb-4">{error}</div>
        <Link to="/kd-stats" className="btn-primary">Back to KD Stats</Link>
      </div>
    );
  }

  if (!playerStats) {
    return (
      <div className="text-center py-8">
        <div className="text-gray-400 text-xl mb-4">No statistics found</div>
        <Link to="/kd-stats" className="btn-primary">Back to KD Stats</Link>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-6">
          <PlayerAvatar 
            player={{
              id: parseInt(id!),
              gamertag: playerStats.gamertag || 'Unknown',
              avatar_url: playerStats.avatar_url
            }} 
            size="xl" 
          />
          <div>
            <h1 className="text-3xl font-bold text-white">{playerStats.gamertag || 'Unknown Player'}</h1>
            <p className="text-gray-400">KD Statistics</p>
          </div>
        </div>
        <div className="flex space-x-2">
          <Link to="/kd-stats" className="btn-secondary">Back to KD Stats</Link>
          <Link to={`/players/${id}`} className="btn-primary">Player Details</Link>
        </div>
      </div>

      {/* Overall Stats Card */}
      <div className="card">
        <h2 className="text-xl font-semibold text-white mb-4">Overall Statistics</h2>
        <div className="grid md:grid-cols-4 gap-4">
          <div className="text-center">
            <div className="text-2xl font-bold text-white">{playerStats.total_matches || 0}</div>
            <div className="text-gray-400 text-sm">Tournaments</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-white">{playerStats.total_maps || 0}</div>
            <div className="text-gray-400 text-sm">Maps Played</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-white">{playerStats.total_kills || 0}</div>
            <div className="text-gray-400 text-sm">Total Kills</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-white">{playerStats.total_deaths || 0}</div>
            <div className="text-gray-400 text-sm">Total Deaths</div>
          </div>
        </div>
        <div className="mt-4 grid md:grid-cols-3 gap-4">
          <div className="text-center">
            <div className={`text-2xl font-bold ${getKDColor(playerStats.avg_kd || 0)}`}>
              {(playerStats.avg_kd || 0).toFixed(3)}
            </div>
            <div className="text-gray-400 text-sm">Average KD</div>
          </div>
          <div className="text-center">
            <div className={`text-2xl font-bold ${getKDColor(playerStats.avg_kda || 0)}`}>
              {(playerStats.avg_kda || 0).toFixed(3)}
            </div>
            <div className="text-gray-400 text-sm">Average KDA</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-white">
              {(playerStats.avg_adr || 0).toFixed(0)}
            </div>
            <div className="text-gray-400 text-sm">Average ADR</div>
          </div>
        </div>
      </div>

      {/* Tournament Stats */}
      <div className="card">
        <h2 className="text-xl font-semibold text-white mb-4">Tournament Statistics</h2>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-800">
                <th className="py-2 px-4 text-left text-white font-bold">Tournament</th>
                <th className="py-2 px-4 text-right text-white font-bold">Kills</th>
                <th className="py-2 px-4 text-right text-white font-bold">Deaths</th>
                <th className="py-2 px-4 text-right text-white font-bold">KD Ratio</th>
                <th className="py-2 px-4 text-right text-white font-bold">KDA Ratio</th>
                <th className="py-2 px-4 text-right text-white font-bold">Maps</th>
              </tr>
            </thead>
            <tbody>
              {playerStats.tournament_stats?.map((tournament: any, index: number) => (
                <tr key={index} className="border-b border-gray-800 hover:bg-gray-900">
                  <td className="py-2 px-4 text-white font-medium">
                    {tournament.tournament_name}
                  </td>
                  <td className="py-2 px-4 text-right text-white">{tournament.kills}</td>
                  <td className="py-2 px-4 text-right text-white">{tournament.deaths}</td>
                  <td className={`py-2 px-4 text-right font-bold ${getKDColor(tournament.kd_ratio)}`}>
                    {tournament.kd_ratio.toFixed(3)}
                  </td>
                  <td className={`py-2 px-4 text-right font-bold ${getKDColor(tournament.kda_ratio)}`}>
                    {tournament.kda_ratio.toFixed(3)}
                  </td>
                  <td className="py-2 px-4 text-right text-white">{tournament.maps_played}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* EWC2025 Detailed Stats (if available) */}
      {playerStats.tournament_stats?.some((t: any) => t.tournament_id === 7) && (
        <div className="card">
          <h2 className="text-xl font-semibold text-white mb-4">EWC 2025 Detailed Statistics</h2>
          <div className="grid md:grid-cols-2 gap-6">
            {/* Search & Destroy */}
            <div>
              <h3 className="text-lg font-semibold text-white mb-3">Search & Destroy</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-gray-400">Kills:</span>
                  <span className="text-white">{playerStats.ewc_snd_kills || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Deaths:</span>
                  <span className="text-white">{playerStats.ewc_snd_deaths || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">KD Ratio:</span>
                  <span className={`font-bold ${getKDColor(playerStats.ewc_snd_kd_ratio || 0)}`}>
                    {(playerStats.ewc_snd_kd_ratio || 0).toFixed(3)}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">+/-:</span>
                  <span className={`font-bold ${getPlusMinusColor(playerStats.ewc_snd_plus_minus || 0)}`}>
                    {(playerStats.ewc_snd_plus_minus || 0).toFixed(0)}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">K/M:</span>
                  <span className="text-white">{(playerStats.ewc_snd_k_per_map || 0).toFixed(1)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">First Kills:</span>
                  <span className="text-white">{playerStats.ewc_snd_first_kills || 0}</span>
                </div>
              </div>
            </div>

            {/* Hardpoint */}
            <div>
              <h3 className="text-lg font-semibold text-white mb-3">Hardpoint</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-gray-400">Kills:</span>
                  <span className="text-white">{playerStats.ewc_hp_kills || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Deaths:</span>
                  <span className="text-white">{playerStats.ewc_hp_deaths || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">KD Ratio:</span>
                  <span className={`font-bold ${getKDColor(playerStats.ewc_hp_kd_ratio || 0)}`}>
                    {(playerStats.ewc_hp_kd_ratio || 0).toFixed(3)}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">+/-:</span>
                  <span className={`font-bold ${getPlusMinusColor(playerStats.ewc_hp_plus_minus || 0)}`}>
                    {(playerStats.ewc_hp_plus_minus || 0).toFixed(0)}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">K/M:</span>
                  <span className="text-white">{(playerStats.ewc_hp_k_per_map || 0).toFixed(1)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Time (sec):</span>
                  <span className="text-white">{Math.round((playerStats.ewc_hp_time_milliseconds || 0) / 1000)}</span>
                </div>
              </div>
            </div>

            {/* Control */}
            <div>
              <h3 className="text-lg font-semibold text-white mb-3">Control</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-gray-400">Kills:</span>
                  <span className="text-white">{playerStats.ewc_control_kills || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Deaths:</span>
                  <span className="text-white">{playerStats.ewc_control_deaths || 0}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">KD Ratio:</span>
                  <span className={`font-bold ${getKDColor(playerStats.ewc_control_kd_ratio || 0)}`}>
                    {(playerStats.ewc_control_kd_ratio || 0).toFixed(3)}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">+/-:</span>
                  <span className={`font-bold ${getPlusMinusColor(playerStats.ewc_control_plus_minus || 0)}`}>
                    {(playerStats.ewc_control_plus_minus || 0).toFixed(0)}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">K/M:</span>
                  <span className="text-white">{(playerStats.ewc_control_k_per_map || 0).toFixed(1)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Captures:</span>
                  <span className="text-white">{playerStats.ewc_control_captures || 0}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PlayerKDStats; 