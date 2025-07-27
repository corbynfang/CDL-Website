import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import type { Player } from '../types';
import { playerApi } from '../services/api';
import PlayerAvatar from './PlayerAvatar';

const PlayerDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [player, setPlayer] = useState<Player | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPlayerData = async () => {
      if (!id) return;
      
      try {
        setLoading(true);
        const playerData = await playerApi.getPlayer(parseInt(id));
        setPlayer(playerData);
      } catch (err) {
        setError('Failed to fetch player details');
        console.error('Error fetching player details:', err);
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

  if (!player) {
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
        <div className="flex items-center space-x-6">
          <PlayerAvatar player={player} size="xl" />
          <div>
            <h1 className="text-3xl font-bold text-white">{player.gamertag}</h1>
            <p className="text-gray-400">Player Details</p>
          </div>
        </div>
        <div className="flex space-x-2">
          <Link to="/players" className="btn-secondary">
            Back to Players
          </Link>
          <Link to={`/players/${player.id}/kd-stats`} className="btn-primary">
            View KD Stats
          </Link>
        </div>
      </div>

      {/* Player Info Card */}
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
          {player.twitter_handle && (
            <div className="md:col-span-2">
              <span className="text-gray-400">Twitter:</span>
              <a 
                href={`https://twitter.com/${player.twitter_handle}`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-400 hover:text-blue-300 ml-2"
              >
                @{player.twitter_handle}
              </a>
            </div>
          )}
        </div>
      </div>

      {/* Quick Stats Card */}
      <div className="card">
        <h2 className="text-xl font-semibold text-white mb-4">Quick Actions</h2>
        <div className="grid md:grid-cols-2 gap-4">
          <Link 
            to={`/players/${player.id}/kd-stats`}
            className="btn-primary text-center"
          >
            View KD Statistics
          </Link>
          <Link 
            to={`/players/${player.id}/kd-stats`}
            className="btn-secondary text-center"
          >
            View Tournament Stats
          </Link>
        </div>
      </div>
    </div>
  );
};

export default PlayerDetails; 