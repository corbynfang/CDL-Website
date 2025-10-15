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
    <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8">
      <div className="mb-12">
        <h1 className="text-3xl sm:text-4xl font-bold text-white mb-3">players</h1>
        <p className="text-lg" style={{ color: '#a3a3a3' }}>
          {players.length} professional cdl players
        </p>
      </div>

      <div className="space-y-2">
        {players.map((player) => (
          <Link
            key={player.id}
            to={`/players/${player.id}`}
            className="group block"
          >
            <div className="flex items-center justify-between p-4 rounded-lg transition-all duration-200 hover:bg-card">
              <div className="flex items-center space-x-4 flex-1 min-w-0">
                <PlayerAvatar player={player} size="md" />
                <div className="flex-1 min-w-0">
                  <h3 className="text-white font-semibold truncate">{player.gamertag}</h3>
                  <p className="text-sm truncate" style={{ color: '#737373' }}>
                    {[player.first_name, player.last_name].filter(Boolean).join(' ') || 'cdl player'} 
                    {player.country && ` • ${player.country}`}
                    {player.role && ` • ${player.role}`}
                  </p>
                </div>
              </div>
              <svg className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform flex-shrink-0" style={{ color: '#737373' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
};

export default Players; 