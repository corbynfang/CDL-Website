import React, { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import type { Player } from '../types';
import { useApi } from '../hooks/useApi';
import PlayerAvatar from './PlayerAvatar';
import LoadingSkeleton, { ErrorDisplay } from './LoadingSkeleton';

const PlayerDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [activeTab, setActiveTab] = useState('overview');

  // Use the new useApi hook with automatic retry logic
  const { data: player, loading, error, refetch } = useApi<Player>(
    `/api/v1/players/${id}`,
    { retries: 3, retryDelay: 1000 }
  );

  if (loading) {
    return <LoadingSkeleton variant="profile" />;
  }

  if (error) {
    return <ErrorDisplay message={error} onRetry={refetch} />;
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

  const tabs = [
    { id: 'overview', label: 'Overview', href: `/players/${player.id}` },
    { id: 'stats', label: 'Stats', href: `/players/${player.id}/kd-stats` },
    { id: 'matches', label: 'Matches', href: `/players/${player.id}/matches` },
    { id: 'events', label: 'Events', href: `/players/${player.id}/events` },
  ];

  return (
    <div className="space-y-0">
      {/* Hero Section - Breaking Point Style */}
      <div className="hero-section">
        <div className="hero-content">
          <div className="flex flex-col lg:flex-row items-start lg:items-center space-y-6 lg:space-y-0 lg:space-x-8">
            {/* Player Avatar */}
            <div className="hero-avatar">
              <div className="relative">
                <PlayerAvatar player={player} size="2xl" />
                <div className="hero-status-indicator"></div>
              </div>
            </div>

            {/* Player Info */}
            <div className="hero-info">
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h1 className="hero-title">
                    {player.gamertag}
                  </h1>
                  <p className="hero-subtitle">
                    {[player.first_name, player.last_name].filter(Boolean).join(' ') || 'Unknown Player'}
                  </p>
                  {player.country && (
                    <p className="text-gray-500 text-sm mt-1">
                      {player.country}
                    </p>
                  )}
                </div>
                
                {/* Status Badge */}
                <div className="mt-4 sm:mt-0">
                  <span className={`status-badge ${
                    player.is_active ? 'status-badge-active' : 'status-badge-inactive'
                  }`}>
                    <span className={`status-indicator ${
                      player.is_active ? 'status-indicator-active' : 'status-indicator-inactive'
                    }`}></span>
                    {player.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
              </div>

              {/* Quick Stats */}
              <div className="hero-stats-grid">
                <div className="hero-stat-item">
                  <div className="hero-stat-value">-</div>
                  <div className="hero-stat-label">KD Ratio</div>
                </div>
                <div className="hero-stat-item">
                  <div className="hero-stat-value">-</div>
                  <div className="hero-stat-label">Kills</div>
                </div>
                <div className="hero-stat-item">
                  <div className="hero-stat-value">-</div>
                  <div className="hero-stat-label">Deaths</div>
                </div>
                <div className="hero-stat-item">
                  <div className="hero-stat-value">-</div>
                  <div className="hero-stat-label">ADR</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Navigation Tabs */}
      <div className="bg-gray-900 border-b border-gray-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <nav className="nav-tabs">
            {tabs.map((tab) => (
              <Link
                key={tab.id}
                to={tab.href}
                className={`nav-tab ${
                  activeTab === tab.id ? 'nav-tab-active' : 'nav-tab-inactive'
                }`}
                onClick={() => setActiveTab(tab.id)}
              >
                {tab.label}
              </Link>
            ))}
          </nav>
        </div>
      </div>

      {/* Content Area */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-8">
          {/* Player Information Card */}
          <div className="card">
            <h2 className="text-xl font-semibold text-white mb-6">Player Information</h2>
            <div className="grid md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div>
                  <span className="text-gray-400 text-sm uppercase tracking-wider">Full Name</span>
                  <p className="text-white font-medium mt-1">
                    {[player.first_name, player.last_name].filter(Boolean).join(' ') || 'N/A'}
                  </p>
                </div>
                <div>
                  <span className="text-gray-400 text-sm uppercase tracking-wider">Country</span>
                  <p className="text-white font-medium mt-1">{player.country || 'N/A'}</p>
                </div>
                <div>
                  <span className="text-gray-400 text-sm uppercase tracking-wider">Role</span>
                  <p className="text-white font-medium mt-1">{player.role || 'N/A'}</p>
                </div>
              </div>
              <div className="space-y-4">
                <div>
                  <span className="text-gray-400 text-sm uppercase tracking-wider">Status</span>
                  <p className={`font-medium mt-1 ${player.is_active ? 'text-green-400' : 'text-red-400'}`}>
                    {player.is_active ? 'Active Player' : 'Inactive Player'}
                  </p>
                </div>
                {player.twitter_handle && (
                  <div>
                    <span className="text-gray-400 text-sm uppercase tracking-wider">Social Media</span>
                    <a 
                      href={`https://twitter.com/${player.twitter_handle}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-blue-400 hover:text-blue-300 font-medium mt-1 block transition-colors duration-200"
                    >
                      @{player.twitter_handle}
                    </a>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Quick Actions Card */}
          <div className="card">
            <h2 className="text-xl font-semibold text-white mb-6">Quick Actions</h2>
            <div className="grid md:grid-cols-2 gap-4">
              <Link 
                to={`/players/${player.id}/kd-stats`}
                className="btn-primary text-center"
              >
                View Detailed Statistics
              </Link>
              <Link 
                to={`/players/${player.id}/matches`}
                className="btn-secondary text-center"
              >
                View Match History
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PlayerDetails; 