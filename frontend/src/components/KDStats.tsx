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
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <div className="text-red-500 text-xl mb-4">{error}</div>
        <button onClick={() => window.location.reload()} className="btn-primary">Try Again</button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-white">All Player KD Statistics (Season & Majors)</h1>
        <div className="text-gray-400">{players.length} Players</div>
      </div>
      <div className="card overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-gray-700">
              <th className="py-3 px-4 text-left text-gray-300 font-medium">Player</th>
              <th className="py-3 px-4 text-left text-gray-300 font-medium">Team</th>
              <th className="py-3 px-4 text-right text-gray-300 font-medium">Season KD</th>
              <th className="py-3 px-4 text-right text-gray-300 font-medium">KD +/-</th>
              {Object.entries(MAJOR_LABELS).map(([id, label]) => (
                <th key={id} className="py-3 px-4 text-right text-gray-300 font-medium">{label} KD</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {players.map((player) => (
              <tr key={player.player_id} className="border-b border-gray-800 hover:bg-gray-800">
                <td className="py-3 px-4">
                  <div className="flex items-center space-x-3">
                    <PlayerAvatar 
                      player={{
                        id: player.player_id,
                        gamertag: player.gamertag,
                        avatar_url: `/assets/avatars/${player.gamertag}.webp`
                      }} 
                      size="sm" 
                    />
                    <Link to={`/players/${player.player_id}`} className="text-white hover:text-blue-400 font-medium">
                      {player.gamertag}
                    </Link>
                  </div>
                </td>
                <td className="py-3 px-4 text-gray-300">{player.team_abbr}</td>
                <td className="py-3 px-4 text-right font-bold">
                  {player.season_kd ? player.season_kd.toFixed(3) : '-'}
                </td>
                <td className={`py-3 px-4 text-right font-bold ${player.season_kd_plus_minus > 0 ? 'text-green-400' : player.season_kd_plus_minus < 0 ? 'text-red-400' : 'text-gray-300'}`}>{player.season_kd_plus_minus ? player.season_kd_plus_minus.toFixed(3) : '-'}</td>
                {Object.keys(MAJOR_LABELS).map((id) => (
                  <td key={id} className="py-3 px-4 text-right font-bold">
                    {player.majors && player.majors[id] ? player.majors[id].toFixed(3) : '-'}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {players.length === 0 && (
        <div className="text-center py-12">
          <div className="text-gray-400 text-xl mb-4">No KD statistics available</div>
          <p className="text-gray-500">Player statistics will appear here once matches are played.</p>
        </div>
      )}
    </div>
  );
};

export default KDStats; 