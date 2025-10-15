import React, { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import PlayerAvatar from './PlayerAvatar';
import LoadingSkeleton, { ErrorDisplay } from './LoadingSkeleton';

const PlayerKDStats: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [activeTab, setActiveTab] = useState('overview');

  const { data: playerStats, loading, error, refetch } = useApi<any>(
    `/api/v1/players/${id}/kd`,
    { retries: 3, retryDelay: 1000 }
  );

  const getKDColor = (kd: number) => {
    if (kd >= 1.2) return 'text-green-400';
    if (kd >= 1.0) return 'text-blue-400';
    if (kd >= 0.9) return 'text-yellow-400';
    return 'text-red-400';
  };

  const getPlusMinusColor = (value: number) => {
    if (value > 0) return 'text-green-400';
    if (value < 0) return 'text-red-400';
    return 'text-gray-300';
  };

  if (loading) {
    return <LoadingSkeleton variant="profile" />;
  }

  if (error) {
    return <ErrorDisplay message={error} onRetry={refetch} />;
  }

  if (!playerStats) {
    return (
      <div className="text-center py-8">
        <div className="text-gray-400 text-xl mb-4">No statistics found</div>
        <Link to="/kd-stats" className="btn-primary">Back to KD Stats</Link>
      </div>
    );
  }

  const tabs = [
    { id: 'overview', label: 'Overview', href: `/players/${id}` },
    { id: 'stats', label: 'Stats', href: `/players/${id}/kd-stats` },
    { id: 'matches', label: 'Matches', href: `/players/${id}/matches` },
    { id: 'events', label: 'Events', href: `/players/${id}/events` },
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
                <PlayerAvatar 
                  player={{
                    id: parseInt(id!),
                    gamertag: playerStats.gamertag || 'Unknown',
                    avatar_url: playerStats.avatar_url
                  }} 
                  size="2xl" 
                />
                <div className="hero-status-indicator"></div>
              </div>
            </div>

            {/* Player Info */}
            <div className="hero-info">
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h1 className="hero-title">
                    {playerStats.gamertag || 'Unknown Player'}
                  </h1>
                  <p className="hero-subtitle">
                    Statistics & Performance
                  </p>
                </div>
              </div>

              {/* Quick Stats */}
              <div className="hero-stats-grid">
                <div className="hero-stat-item">
                  <div className={`hero-stat-value ${getKDColor(playerStats.avg_kd || 0)}`}>
                    {(playerStats.avg_kd || 0).toFixed(3)}
                  </div>
                  <div className="hero-stat-label">Avg KD</div>
                </div>
                <div className="hero-stat-item">
                  <div className="hero-stat-value">
                    {playerStats.total_kills || 0}
                  </div>
                  <div className="hero-stat-label">Total Kills</div>
                </div>
                <div className="hero-stat-item">
                  <div className="hero-stat-value">
                    {playerStats.total_deaths || 0}
                  </div>
                  <div className="hero-stat-label">Total Deaths</div>
                </div>
                <div className="hero-stat-item">
                  <div className="hero-stat-value">
                    {(playerStats.avg_adr || 0).toFixed(0)}
                  </div>
                  <div className="hero-stat-label">Avg ADR</div>
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
          {/* Overall Stats Card */}
          <div className="card">
            <h2 className="text-xl font-semibold text-white mb-6">Overall Statistics</h2>
            <div className="grid md:grid-cols-4 gap-6 mb-6">
              <div className="card-stats">
                <div className="text-2xl font-black text-white">{playerStats.total_matches || 0}</div>
                <div className="text-gray-400 text-sm uppercase tracking-wider">Tournaments</div>
              </div>
              <div className="card-stats">
                <div className="text-2xl font-black text-white">{playerStats.total_maps || 0}</div>
                <div className="text-gray-400 text-sm uppercase tracking-wider">Maps Played</div>
              </div>
              <div className="card-stats">
                <div className="text-2xl font-black text-white">{playerStats.total_kills || 0}</div>
                <div className="text-gray-400 text-sm uppercase tracking-wider">Total Kills</div>
              </div>
              <div className="card-stats">
                <div className="text-2xl font-black text-white">{playerStats.total_deaths || 0}</div>
                <div className="text-gray-400 text-sm uppercase tracking-wider">Total Deaths</div>
              </div>
            </div>
            <div className="grid md:grid-cols-3 gap-6">
              <div className="card-stats">
                <div className={`text-2xl font-black ${getKDColor(playerStats.avg_kd || 0)}`}>
                  {(playerStats.avg_kd || 0).toFixed(3)}
                </div>
                <div className="text-gray-400 text-sm uppercase tracking-wider">Average KD</div>
              </div>
              <div className="card-stats">
                <div className={`text-2xl font-black ${getKDColor(playerStats.avg_kda || 0)}`}>
                  {(playerStats.avg_kda || 0).toFixed(3)}
                </div>
                <div className="text-gray-400 text-sm uppercase tracking-wider">Average KDA</div>
              </div>
              <div className="card-stats">
                <div className="text-2xl font-black text-white">
                  {(playerStats.avg_adr || 0).toFixed(0)}
                </div>
                <div className="text-gray-400 text-sm uppercase tracking-wider">Average ADR</div>
              </div>
            </div>
          </div>

          {/* Tournament Stats */}
          <div className="card">
            <h2 className="text-xl font-semibold text-white mb-6">Tournament Statistics</h2>
            <div className="overflow-x-auto">
              <table className="table-modern">
                <thead>
                  <tr>
                    <th>Tournament</th>
                    <th className="text-right">Kills</th>
                    <th className="text-right">Deaths</th>
                    <th className="text-right">KD Ratio</th>
                    <th className="text-right">KDA Ratio</th>
                    <th className="text-right">Maps</th>
                  </tr>
                </thead>
                <tbody>
                  {playerStats.tournament_stats?.map((tournament: any, index: number) => (
                    <tr key={index}>
                      <td className="font-medium">
                        {tournament.tournament_name}
                      </td>
                      <td className="text-right">{tournament.kills}</td>
                      <td className="text-right">{tournament.deaths}</td>
                      <td className={`text-right font-bold ${getKDColor(tournament.kd_ratio)}`}>
                        {tournament.kd_ratio.toFixed(3)}
                      </td>
                      <td className={`text-right font-bold ${getKDColor(tournament.kda_ratio)}`}>
                        {tournament.kda_ratio.toFixed(3)}
                      </td>
                      <td className="text-right">{tournament.maps_played}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* EWC2025 Detailed Stats (if available) */}
          {playerStats.tournament_stats?.some((t: any) => t.tournament_id === 7) && (
            <div className="card">
              <h2 className="text-xl font-semibold text-white mb-6">EWC 2025 Detailed Statistics</h2>
              <div className="grid md:grid-cols-3 gap-6">
                {/* Search & Destroy */}
                <div className="game-mode-card">
                  <h3 className="game-mode-header">
                    <span className="game-mode-indicator game-mode-indicator-snd"></span>
                    Search & Destroy
                  </h3>
                  <div className="space-y-3">
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">Kills</span>
                      <span className="game-mode-stat-value">{playerStats.ewc_snd_kills || 0}</span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">Deaths</span>
                      <span className="game-mode-stat-value">{playerStats.ewc_snd_deaths || 0}</span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">KD Ratio</span>
                      <span className={`game-mode-stat-value game-mode-stat-value-bold ${getKDColor(playerStats.ewc_snd_kd_ratio || 0)}`}>
                        {(playerStats.ewc_snd_kd_ratio || 0).toFixed(3)}
                      </span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">+/-</span>
                      <span className={`game-mode-stat-value game-mode-stat-value-bold ${getPlusMinusColor(playerStats.ewc_snd_plus_minus || 0)}`}>
                        {(playerStats.ewc_snd_plus_minus || 0).toFixed(0)}
                      </span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">K/M</span>
                      <span className="game-mode-stat-value">{(playerStats.ewc_snd_k_per_map || 0).toFixed(1)}</span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">First Kills</span>
                      <span className="game-mode-stat-value">{playerStats.ewc_snd_first_kills || 0}</span>
                    </div>
                  </div>
                </div>

                {/* Hardpoint */}
                <div className="game-mode-card">
                  <h3 className="game-mode-header">
                    <span className="game-mode-indicator game-mode-indicator-hp"></span>
                    Hardpoint
                  </h3>
                  <div className="space-y-3">
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">Kills</span>
                      <span className="game-mode-stat-value">{playerStats.ewc_hp_kills || 0}</span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">Deaths</span>
                      <span className="game-mode-stat-value">{playerStats.ewc_hp_deaths || 0}</span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">KD Ratio</span>
                      <span className={`game-mode-stat-value game-mode-stat-value-bold ${getKDColor(playerStats.ewc_hp_kd_ratio || 0)}`}>
                        {(playerStats.ewc_hp_kd_ratio || 0).toFixed(3)}
                      </span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">+/-</span>
                      <span className={`game-mode-stat-value game-mode-stat-value-bold ${getPlusMinusColor(playerStats.ewc_hp_plus_minus || 0)}`}>
                        {(playerStats.ewc_hp_plus_minus || 0).toFixed(0)}
                      </span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">K/M</span>
                      <span className="game-mode-stat-value">{(playerStats.ewc_hp_k_per_map || 0).toFixed(1)}</span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">Time (sec)</span>
                      <span className="game-mode-stat-value">{Math.round((playerStats.ewc_hp_time_milliseconds || 0) / 1000)}</span>
                    </div>
                  </div>
                </div>

                {/* Control */}
                <div className="game-mode-card">
                  <h3 className="game-mode-header">
                    <span className="game-mode-indicator game-mode-indicator-control"></span>
                    Control
                  </h3>
                  <div className="space-y-3">
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">Kills</span>
                      <span className="game-mode-stat-value">{playerStats.ewc_control_kills || 0}</span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">Deaths</span>
                      <span className="game-mode-stat-value">{playerStats.ewc_control_deaths || 0}</span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">KD Ratio</span>
                      <span className={`game-mode-stat-value game-mode-stat-value-bold ${getKDColor(playerStats.ewc_control_kd_ratio || 0)}`}>
                        {(playerStats.ewc_control_kd_ratio || 0).toFixed(3)}
                      </span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">+/-</span>
                      <span className={`game-mode-stat-value game-mode-stat-value-bold ${getPlusMinusColor(playerStats.ewc_control_plus_minus || 0)}`}>
                        {(playerStats.ewc_control_plus_minus || 0).toFixed(0)}
                      </span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">K/M</span>
                      <span className="game-mode-stat-value">{(playerStats.ewc_control_k_per_map || 0).toFixed(1)}</span>
                    </div>
                    <div className="game-mode-stat">
                      <span className="game-mode-stat-label">Captures</span>
                      <span className="game-mode-stat-value">{playerStats.ewc_control_captures || 0}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default PlayerKDStats; 