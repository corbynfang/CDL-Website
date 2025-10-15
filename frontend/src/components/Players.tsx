import React from 'react';
import { Link } from 'react-router-dom';
import type { Player } from '../types';
import { useApi } from '../hooks/useApi';
import PlayerAvatar from './PlayerAvatar';
import LoadingSkeleton, { ErrorDisplay } from './LoadingSkeleton';

const Players: React.FC = () => {
  const { data: players, loading, error, refetch } = useApi<Player[]>('/api/v1/players', {
    retries: 3,
    retryDelay: 1000
  });

  if (loading) {
    return <LoadingSkeleton variant="card" count={6} />;
  }

  if (error) {
    return <ErrorDisplay message={error} onRetry={refetch} />;
  }

  if (!players || players.length === 0) {
    return (
      <div className="text-center py-8 px-4">
        <div className="text-gray-400 text-xl mb-4">No players found</div>
      </div>
    );
  }

  return (
    <div className="space-y-6 sm:space-y-8">
      {/* Hero Section */}
      <div className="text-center py-12 sm:py-16 bg-black border border-gray-800 mb-6 sm:mb-8 px-4 sm:px-6">
        <h1 className="text-3xl sm:text-4xl md:text-heading text-white mb-4 sm:mb-6">CDL PLAYERS</h1>
        <p className="text-base sm:text-lg md:text-subheading text-gray-400 max-w-3xl mx-auto leading-relaxed uppercase tracking-wider px-4">
          PROFESSIONAL CALL OF DUTY LEAGUE PLAYERS COMPETING AT THE HIGHEST LEVEL.
          EXPLORE INDIVIDUAL STATISTICS, PERFORMANCE METRICS, AND PLAYER PROFILES.
        </p>
        <div className="mt-6 sm:mt-8 text-gray-400 text-base sm:text-lg uppercase tracking-wider">
          {players.length} ACTIVE PLAYERS
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 md:gap-8">
        {players.map((player) => (
          <div key={player.id} className="card hover:border-white transition-all duration-300 transform hover:scale-105 border border-gray-800 overflow-hidden">
            <div className="flex items-center space-x-3 sm:space-x-4 mb-4 sm:mb-6 p-4 sm:p-6">
              <PlayerAvatar player={player} size="lg" />
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <h3 className="text-lg sm:text-xl md:text-2xl font-bold text-white uppercase tracking-wider">{player.gamertag}</h3>
                </div>
              </div>
            </div>

            <div className="space-y-2 mb-4 px-4 sm:px-6">
              {(player.first_name || player.last_name) && (
                <div className="flex items-center text-gray-300">
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  <span className="uppercase tracking-wider text-sm sm:text-base">{[player.first_name, player.last_name].filter(Boolean).join(' ')}</span>
                </div>
              )}

              {player.country && (
                <div className="flex items-center text-gray-300">
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span className="uppercase tracking-wider text-sm sm:text-base">{player.country}</span>
                </div>
              )}

              {player.role && (
                <div className="flex items-center text-gray-300">
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span className="uppercase tracking-wider text-sm sm:text-base">{player.role}</span>
                </div>
              )}
            </div>

            <div className="px-4 sm:px-6 pb-4 sm:pb-6">
              <Link
                to={`/players/${player.id}`}
                className="btn-secondary w-full text-center text-sm sm:text-base py-2 sm:py-3"
              >
                VIEW PROFILE
              </Link>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Players; 