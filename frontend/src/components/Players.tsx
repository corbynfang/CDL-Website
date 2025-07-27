import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import type { Player } from '../types';
import { playerApi } from '../services/api';
import PlayerAvatar from './PlayerAvatar';

const Players: React.FC = () => {
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPlayers = async () => {
      try {
        setLoading(true);
        const data = await playerApi.getPlayers();
        setPlayers(data);
      } catch (err) {
        setError('Failed to fetch players');
        console.error('Error fetching players:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchPlayers();
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
        <button
          onClick={() => window.location.reload()}
          className="btn-primary"
        >
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="text-center py-16 bg-gradient-to-r from-green-900 to-blue-900 rounded-lg mb-8">
        <h1 className="text-5xl font-bold text-white mb-6">CDL Players</h1>
        <p className="text-xl text-gray-300 max-w-3xl mx-auto leading-relaxed">
          Professional Call of Duty League players competing at the highest level.
          Explore individual statistics, performance metrics, and player profiles.
        </p>
        <div className="mt-8 text-gray-400 text-lg">
          {players.length} Active Players
        </div>
      </div>

      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        {players.map((player) => (
          <div key={player.id} className="card hover:bg-gray-750 transition-all duration-300 transform hover:scale-105 border border-gray-700 hover:border-green-500 overflow-hidden">
            <div className="flex items-center space-x-4 mb-6 p-6">
              <PlayerAvatar player={player} size="xl" />
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <h3 className="text-2xl font-bold text-white">{player.gamertag}</h3>
                  {player.is_active && (
                    <span className="text-xs text-green-400 bg-green-900 px-3 py-1 rounded-full font-bold">
                      Active
                    </span>
                  )}
                </div>
              </div>
            </div>

            <div className="space-y-2 mb-4">
              {(player.first_name || player.last_name) && (
                <div className="flex items-center text-gray-300">
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  {[player.first_name, player.last_name].filter(Boolean).join(' ')}
                </div>
              )}

              {player.country && (
                <div className="flex items-center text-gray-300">
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  {player.country}
                </div>
              )}

              {player.role && (
                <div className="flex items-center text-gray-300">
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                  </svg>
                  {player.role}
                </div>
              )}

              {player.twitter_handle && (
                <div className="flex items-center text-gray-300">
                  <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M23.953 4.57a10 10 0 01-2.825.775 4.958 4.958 0 002.163-2.723c-.951.555-2.005.959-3.127 1.184a4.92 4.92 0 00-8.384 4.482C7.69 8.095 4.067 6.13 1.64 3.162a4.822 4.822 0 00-.666 2.475c0 1.71.87 3.213 2.188 4.096a4.904 4.904 0 01-2.228-.616v.06a4.923 4.923 0 003.946 4.827 4.996 4.996 0 01-2.212.085 4.936 4.936 0 004.604 3.417 9.867 9.867 0 01-6.102 2.105c-.39 0-.779-.023-1.17-.067a13.995 13.995 0 007.557 2.209c9.053 0 13.998-7.496 13.998-13.985 0-.21 0-.42-.015-.63A9.935 9.935 0 0024 4.59z"/>
                  </svg>
                  @{player.twitter_handle}
                </div>
              )}
            </div>

            <div className="flex space-x-2">
              <Link
                to={`/players/${player.id}`}
                className="btn-primary flex-1 text-center"
              >
                View Details
              </Link>
              <Link
                to={`/players/${player.id}/kd-stats`}
                className="btn-secondary flex-1 text-center"
              >
                KD Stats
              </Link>
            </div>
          </div>
        ))}
      </div>

      {players.length === 0 && (
        <div className="text-center py-12">
          <div className="text-gray-400 text-xl mb-4">No players found</div>
          <p className="text-gray-500">There are currently no players available.</p>
        </div>
      )}
    </div>
  );
};

export default Players; 