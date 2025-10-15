import React, { useState, useMemo } from 'react';
import type { PlayerTransfer } from '../types';
import { useApi } from '../hooks/useApi';
import PlayerAvatar from './PlayerAvatar';
import LoadingSkeleton, { ErrorDisplay } from './LoadingSkeleton';

const Transfers: React.FC = () => {
  const [filters] = useState({
    season: 'Black Ops 6',
    team_id: '',
    type: ''
  });

  const apiUrl = useMemo(() => {
    const params = new URLSearchParams();
    if (filters.season) params.append('season', filters.season);
    if (filters.team_id) params.append('team_id', filters.team_id);
    if (filters.type) params.append('type', filters.type);
    const query = params.toString();
    return `/api/v1/transfers${query ? `?${query}` : ''}`;
  }, [filters.season, filters.team_id, filters.type]);

  const { data: response, loading, error, refetch } = useApi<{ transfers: PlayerTransfer[]; count: number; timestamp: number }>(
    apiUrl,
    { retries: 3, retryDelay: 1000 }
  );

  const transfers = response?.transfers || [];

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const getTransferDescription = (transfer: PlayerTransfer) => {
    // Use the description field if available, otherwise fall back to the old logic
    if (transfer.description && transfer.description.trim() !== '') {
      return transfer.description;
    }

    const playerName = transfer.player?.gamertag || 'Unknown Player';
    const fromTeam = transfer.from_team?.name || 'Free Agent';
    const toTeam = transfer.to_team?.name || 'Unknown Team';
    const role = transfer.role || 'Player';

    if (fromTeam === 'Free Agent') {
      return `${playerName} joins ${toTeam} as ${role}`;
    } else if (toTeam === 'Free Agent') {
      return `${fromTeam} parts ways with ${playerName}`;
    } else {
      return `${playerName} moves from ${fromTeam} to ${toTeam}, taking on the ${role} position`;
    }
  };

  if (loading) {
    return <LoadingSkeleton variant="list" count={5} />;
  }

  if (error) {
    return <ErrorDisplay message={error} onRetry={refetch} />;
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8">
      <div className="mb-12">
        <h1 className="text-3xl sm:text-4xl font-bold text-white mb-3">transfers</h1>
        <p className="text-lg" style={{ color: '#a3a3a3' }}>
          {transfers?.length || 0} player movements • {filters.season}
        </p>
      </div>

      <div className="space-y-2">
        {transfers && transfers.length > 0 ? (
          transfers.map((transfer) => (
            <div
              key={transfer.id}
              className="block p-4 rounded-lg transition-all duration-200 hover:bg-card"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex items-start space-x-3 flex-1 min-w-0">
                  <PlayerAvatar player={transfer.player!} size="sm" />
                  <div className="flex-1 min-w-0">
                    <h3 className="text-white font-semibold mb-1">{transfer.player?.gamertag}</h3>
                    <p className="text-sm mb-2" style={{ color: '#a3a3a3' }}>
                      {getTransferDescription(transfer)}
                    </p>
                    <div className="flex items-center gap-3 text-xs" style={{ color: '#737373' }}>
                      {transfer.from_team && <span>{transfer.from_team.abbreviation}</span>}
                      <span>→</span>
                      {transfer.to_team && <span>{transfer.to_team.abbreviation}</span>}
                      <span>•</span>
                      <span>{formatDate(transfer.transfer_date)}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          ))
        ) : (
          <div className="text-center py-12">
            <p style={{ color: '#737373' }}>no transfers found</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Transfers; 