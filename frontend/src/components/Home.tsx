import React from 'react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => {
  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="text-center py-12">
        <h1 className="text-4xl md:text-6xl font-bold text-white mb-6">
          Welcome to <span className="text-blue-500">CDL Stats</span>
        </h1>
        <p className="text-xl text-gray-300 max-w-3xl mx-auto mb-8">
          Your comprehensive source for Call of Duty League statistics, team information, and player performance data.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link
            to="/teams"
            className="btn-primary text-lg px-8 py-3"
          >
            Meet the Teams
          </Link>
          <Link
            to="/players"
            className="btn-secondary text-lg px-8 py-3"
          >
            View Players
          </Link>
        </div>
      </div>

      {/* Featured Teams Section */}
      <div className="card">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-bold text-white mb-4">Meet the CDL Teams</h2>
          <p className="text-gray-300 max-w-2xl mx-auto">
            Discover the elite organizations competing in the Call of Duty League. 
            Explore team rosters, player profiles, and performance statistics.
          </p>
        </div>
        
        <div className="grid md:grid-cols-3 gap-6 mb-8">
          <div className="bg-gradient-to-br from-blue-600 to-purple-600 rounded-lg p-6 text-center">
            <div className="w-16 h-16 bg-white bg-opacity-20 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-2">Team Profiles</h3>
            <p className="text-blue-100 text-sm">
              Detailed team information, logos, colors, and organizational history.
            </p>
          </div>

          <div className="bg-gradient-to-br from-green-600 to-blue-600 rounded-lg p-6 text-center">
            <div className="w-16 h-16 bg-white bg-opacity-20 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-2">Player Rosters</h3>
            <p className="text-green-100 text-sm">
              Complete team rosters with player profiles, roles, and performance stats.
            </p>
          </div>

          <div className="bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg p-6 text-center">
            <div className="w-16 h-16 bg-white bg-opacity-20 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-2">Performance Stats</h3>
            <p className="text-purple-100 text-sm">
              Tournament results, match statistics, and team performance analytics.
            </p>
          </div>
        </div>

        <div className="text-center">
          <Link
            to="/teams"
            className="btn-primary text-lg px-8 py-3"
          >
            Explore All Teams
          </Link>
        </div>
      </div>

      {/* Features Grid */}
      <div className="grid md:grid-cols-3 gap-6 mt-12">
        <div className="card">
          <div className="text-center">
            <div className="w-12 h-12 bg-blue-500 rounded-lg flex items-center justify-center mx-auto mb-4">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-2">Team Statistics</h3>
            <p className="text-gray-400">
              Comprehensive team performance data, match results, and tournament statistics.
            </p>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="w-12 h-12 bg-green-500 rounded-lg flex items-center justify-center mx-auto mb-4">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-2">Player Profiles</h3>
            <p className="text-gray-400">
              Detailed player information, performance metrics, and career statistics.
            </p>
          </div>
        </div>

        <div className="card">
          <div className="text-center">
            <div className="w-12 h-12 bg-purple-500 rounded-lg flex items-center justify-center mx-auto mb-4">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-2">Live Statistics</h3>
            <p className="text-gray-400">
              Real-time match data (BETA), K/D ratios, damage statistics, and more.
            </p>
          </div>
        </div>
      </div>

      {/* Current Season Info */}
      <div className="card mt-8">
        <h2 className="text-2xl font-bold text-white mb-4">CDL 2025 Season</h2>
        <div className="grid md:grid-cols-2 gap-6">
          <div>
            <h3 className="text-lg font-semibold text-blue-400 mb-2">Game</h3>
            <p className="text-gray-300">Call of Duty: Black Ops 6</p>
          </div>
          <div>
            <h3 className="text-lg font-semibold text-blue-400 mb-2">Active Teams</h3>
            <p className="text-gray-300">12 Teams</p>
          </div>
          <div>
            <h3 className="text-lg font-semibold text-blue-400 mb-2">Season Start</h3>
            <p className="text-gray-300">January 31, 2025</p>
          </div>
          <div>
            <h3 className="text-lg font-semibold text-blue-400 mb-2">Status</h3>
            <p className="text-green-400 font-semibold">Active</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home; 