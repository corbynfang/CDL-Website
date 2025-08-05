import React, { useState, useEffect } from 'react';
import type { PlayerTransfer } from '../types';
import { transfersApi } from '../services/api';
import PlayerAvatar from './PlayerAvatar';
import TeamLogo from './TeamLogo';

const Transfers: React.FC = () => {
  const [transfers, setTransfers] = useState<PlayerTransfer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState({
    season: 'Black Ops 6',
    team_id: '',
    type: ''
  });

  useEffect(() => {
    fetchTransfers();
  }, [filters]);

  const fetchTransfers = async () => {
    try {
      setLoading(true);
      const params: any = {};
      if (filters.season) params.season = filters.season;
      if (filters.team_id) params.team_id = filters.team_id;
      if (filters.type) params.type = filters.type;
      
      const data = await transfersApi.getTransfers(params);
      setTransfers(data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch transfers');
      console.error('Error fetching transfers:', err);
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const getTransferDescription = (transfer: PlayerTransfer) => {
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
        <button 
          onClick={fetchTransfers}
          className="btn-primary"
        >
          RETRY
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="text-center py-16 bg-black border border-gray-800 mb-8">
        <h1 className="text-heading text-white mb-6">PLAYER TRANSFERS</h1>
        <p className="text-subheading text-gray-400 max-w-3xl mx-auto leading-relaxed uppercase tracking-wider">
          TRACK ALL PLAYER MOVEMENTS, ROSTER CHANGES, AND TEAM TRANSITIONS THROUGHOUT THE CDL SEASON
        </p>
        <div className="mt-8 text-gray-400 text-lg uppercase tracking-wider">
          {transfers.length} TRANSFERS
        </div>
      </div>

      {/* Filters */}
      <div className="card">
        <div className="grid md:grid-cols-3 gap-4">
          <div>
            <label className="block text-white text-sm font-medium mb-2 uppercase tracking-wider">SEASON</label>
            <select
              value={filters.season}
              onChange={(e) => setFilters({ ...filters, season: e.target.value })}
              className="input w-full"
            >
              <option value="Black Ops 6">BLACK OPS 6</option>
            </select>
          </div>
          <div>
            <label className="block text-white text-sm font-medium mb-2 uppercase tracking-wider">TEAM</label>
            <select
              value={filters.team_id}
              onChange={(e) => setFilters({ ...filters, team_id: e.target.value })}
              className="input w-full"
            >
              <option value="">ALL TEAMS</option>
              <option value="1">OPTIC TEX</option>
              <option value="2">TOR ULTRA</option>
              <option value="3">BOS BREACH</option>
              <option value="4">CAR ROYAL RAVENS</option>
              <option value="5">LA THIEVES</option>
              <option value="6">ATL FAZE</option>
              <option value="7">VAN SURGE</option>
              <option value="8">MIA HERETICS</option>
              <option value="9">LA GUERRILLAS M8</option>
              <option value="10">MIN RØKKR</option>
              <option value="11">CLOUD9 NY</option>
              <option value="12">LV FALCONS</option>
            </select>
          </div>
          <div>
            <label className="block text-white text-sm font-medium mb-2 uppercase tracking-wider">TYPE</label>
            <select
              value={filters.type}
              onChange={(e) => setFilters({ ...filters, type: e.target.value })}
              className="input w-full"
            >
              <option value="">ALL TYPES</option>
              <option value="transfer">TRANSFER</option>
              <option value="signing">SIGNING</option>
              <option value="release">RELEASE</option>
            </select>
          </div>
        </div>
      </div>

      {/* Transfers List */}
      <div className="space-y-4">
        {transfers.map((transfer) => (
          <div key={transfer.id} className="card hover:border-white transition-all duration-300">
            <div className="flex items-start space-x-4">
              <div className="flex-shrink-0">
                <PlayerAvatar player={transfer.player!} size="md" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="text-xl font-bold text-white uppercase tracking-wider">
                    {transfer.player?.gamertag}
                  </h3>
                  <span className="text-sm text-gray-400 uppercase tracking-wider">
                    {formatDate(transfer.transfer_date)}
                  </span>
                </div>
                <p className="text-gray-300 mb-4">{getTransferDescription(transfer)}</p>
                <div className="flex items-center space-x-4">
                  {transfer.from_team && (
                    <div className="flex items-center space-x-2">
                      <TeamLogo team={transfer.from_team} size="sm" />
                      <span className="text-gray-400 uppercase tracking-wider">{transfer.from_team.name}</span>
                    </div>
                  )}
                  <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
                  </svg>
                  {transfer.to_team && (
                    <div className="flex items-center space-x-2">
                      <TeamLogo team={transfer.to_team} size="sm" />
                      <span className="text-gray-400 uppercase tracking-wider">{transfer.to_team.name}</span>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {transfers.length === 0 && (
        <div className="text-center py-12">
          <div className="text-gray-400 text-xl mb-4 uppercase tracking-wider">NO TRANSFERS FOUND</div>
          <p className="text-gray-500 uppercase tracking-wider">No transfers match the current filters.</p>
        </div>
      )}
    </div>
  );
};

export default Transfers; 