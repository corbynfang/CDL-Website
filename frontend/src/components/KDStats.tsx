import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { statsApi } from '../services/api';
import PlayerAvatar from './PlayerAvatar';

const MAJOR_LABELS = {
  1: 'Major 1',
  2: 'Major 2',
  3: 'Major 3',
  4: 'Major 4',
  5: 'Champs',
};

const KDStats: React.FC = () => {
  const [players, setPlayers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchAllKD = async () => {
      try {
        setLoading(true);
        const data = await statsApi.getAllPlayersKDStats();
        setPlayers(data);
      } catch (err) {
        setError('Failed to fetch KD statistics');
        console.error('Error fetching KD stats:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchAllKD();
  }, []);

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
        <button onClick={() => window.location.reload()} className="btn-primary">TRY AGAIN</button>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="text-center py-16 bg-black border border-gray-800 mb-8">
        <h1 className="text-heading text-white mb-6">KD STATISTICS</h1>
        <p className="text-subheading text-gray-400 max-w-3xl mx-auto leading-relaxed uppercase tracking-wider">
          COMPREHENSIVE KILL/DEATH RATIOS FOR ALL CDL PLAYERS ACROSS SEASON AND MAJOR TOURNAMENTS
        </p>
        <div className="mt-8 text-gray-400 text-lg uppercase tracking-wider">
          {players.length} PLAYERS
        </div>
      </div>

      <div className="card overflow-x-auto">
        <table className="table">
          <thead>
            <tr>
              <th>PLAYER</th>
              <th>TEAM</th>
              <th className="text-right">SEASON KD</th>
              <th className="text-right">KD +/-</th>
              {Object.entries(MAJOR_LABELS).map(([id, label]) => (
                <th key={id} className="text-right">{label.toUpperCase()} KD</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {players.map((player) => (
              <tr key={player.player_id} className="hover:bg-gray-900">
                <td>
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
                <td className="text-gray-300 uppercase tracking-wider">{player.team_abbr}</td>
                <td className="text-right font-bold text-white">
                  {player.season_kd ? player.season_kd.toFixed(3) : '-'}
                </td>
                <td className={`text-right font-bold ${player.season_kd_plus_minus > 0 ? 'text-green-400' : player.season_kd_plus_minus < 0 ? 'text-red-400' : 'text-gray-300'}`}>
                  {player.season_kd_plus_minus ? player.season_kd_plus_minus.toFixed(3) : '-'}
                </td>
                {Object.keys(MAJOR_LABELS).map((id) => (
                  <td key={id} className="text-right font-bold text-white">
                    {player.majors && player.majors[id] ? player.majors[id].toFixed(3) : '-'}
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