import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
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
      <div className="min-h-screen bg-gray-900 text-white">
        <div className="container mx-auto px-4 py-8">
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-900 text-white">
        <div className="container mx-auto px-4 py-8">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-red-500 mb-4">Error</h1>
            <p className="text-gray-400">{error}</p>
            <button 
              onClick={fetchTransfers}
              className="mt-4 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <div className="container mx-auto px-4 py-8">
        {/* Hero Section */}
        <div className="text-center py-16 bg-gradient-to-r from-purple-900 to-blue-900 rounded-lg mb-8">
          <h1 className="text-5xl font-bold text-white mb-6">Player Transfers</h1>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto leading-relaxed">
            Track all player movements and roster changes throughout the Call of Duty League season.
            Stay updated with the latest transfers, signings, and team changes.
          </p>
          <div className="mt-8 text-gray-400 text-lg">
            {transfers.length} Transfers in {filters.season}
          </div>
        </div>

        {/* Filters */}
        <div className="bg-gray-800 rounded-lg p-6 mb-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Season</label>
              <select
                value={filters.season}
                onChange={(e) => setFilters({ ...filters, season: e.target.value })}
                className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
              >
                <option value="Black Ops 6">Black Ops 6</option>
                <option value="Modern Warfare 3">Modern Warfare 3</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Transfer Type</label>
              <select
                value={filters.type}
                onChange={(e) => setFilters({ ...filters, type: e.target.value })}
                className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
              >
                <option value="">All Types</option>
                <option value="CDL">CDL</option>
                <option value="Challengers">Challengers</option>
                <option value="Free Agent">Free Agent</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">Team</label>
              <select
                value={filters.team_id}
                onChange={(e) => setFilters({ ...filters, team_id: e.target.value })}
                className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:border-blue-500"
              >
                <option value="">All Teams</option>
                <option value="1">OpTic TEX</option>
                <option value="2">ATL FaZe</option>
                <option value="3">TOR Ultra</option>
                <option value="4">LV Falcons</option>
                <option value="5">CAR Royal Ravens</option>
                <option value="6">LA Guerrillas M8</option>
                <option value="7">VAN Surge</option>
                <option value="8">MIA Heretics</option>
                <option value="9">MIN RØKKR</option>
                <option value="10">BOS Breach</option>
                <option value="11">LA Thieves</option>
              </select>
            </div>
          </div>
        </div>

        {/* Transfers List */}
        <div className="space-y-6">
          {transfers.map((transfer) => (
            <div key={transfer.id} className="bg-gray-800 rounded-lg p-6 border border-gray-700 hover:border-purple-500 transition-colors">
              <div className="flex items-start space-x-4">
                {/* Player Avatar */}
                <div className="flex-shrink-0">
                  <PlayerAvatar player={transfer.player!} size="lg" />
                </div>

                {/* Transfer Details */}
                <div className="flex-1">
                  <div className="flex items-center justify-between mb-2">
                    <h3 className="text-xl font-bold text-white">
                      {transfer.player?.gamertag}
                    </h3>
                    <span className="text-sm text-gray-400">
                      {formatDate(transfer.transfer_date)}
                    </span>
                  </div>

                  <p className="text-gray-300 mb-4">
                    {getTransferDescription(transfer)}
                  </p>

                                     {/* Teams */}
                   <div className="flex items-center space-x-4">
                     {transfer.from_team && (
                       <div className="flex items-center space-x-2">
                         <span className="text-gray-400 text-sm">From:</span>
                         <TeamLogo team={transfer.from_team} size="sm" />
                         <span className="text-white text-sm">{transfer.from_team.name}</span>
                       </div>
                     )}
                     
                     <div className="text-gray-500">→</div>
                     
                     <div className="flex items-center space-x-2">
                       <span className="text-gray-400 text-sm">To:</span>
                       {transfer.to_team && (
                         <>
                           <TeamLogo team={transfer.to_team} size="sm" />
                           <span className="text-white text-sm">{transfer.to_team.name}</span>
                         </>
                       )}
                     </div>
                   </div>

                  {/* Transfer Type Badge */}
                  <div className="mt-4">
                    <span className={`inline-block px-3 py-1 rounded-full text-xs font-bold ${
                      transfer.transfer_type === 'CDL' 
                        ? 'bg-blue-600 text-white' 
                        : transfer.transfer_type === 'Challengers'
                        ? 'bg-green-600 text-white'
                        : 'bg-gray-600 text-white'
                    }`}>
                      {transfer.transfer_type}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>

        {transfers.length === 0 && !loading && (
          <div className="text-center py-12">
            <p className="text-gray-400 text-lg">No transfers found for the selected filters.</p>
          </div>
        )}

        {/* Back Button */}
        <div className="mt-8 text-center">
          <Link 
            to="/" 
            className="inline-flex items-center px-6 py-3 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition-colors"
          >
            ← Back to Home
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Transfers; 