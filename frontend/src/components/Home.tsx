import React from 'react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => {
  return (
    <div className="min-h-screen flex flex-col" style={{ backgroundColor: '#0a0a0a' }}>
      {/* Minimal Hero Section - davis7.sh inspired */}
      <div className="flex-1 flex items-center justify-center px-4 sm:px-6">
        <div className="max-w-2xl w-full">
          {/* Main Title */}
          <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold text-white mb-6" style={{ lineHeight: '1.2' }}>
            CDLytics
          </h1>
          
          {/* Subtitle */}
          <p className="text-lg sm:text-xl text-muted mb-12" style={{ color: '#a3a3a3' }}>
            professional call of duty league statistics, player analytics, and team insights.
          </p>
          
          {/* Navigation Links - davis7.sh style */}
          <div className="space-y-4">
            <Link
              to="/players"
              className="group block"
            >
              <div className="flex items-center justify-between p-4 rounded-lg transition-all duration-200 hover:bg-card">
                <div>
                  <h3 className="text-white font-semibold mb-1">Players</h3>
                  <p className="text-sm" style={{ color: '#737373' }}>browse all 79 professional cdl players</p>
                </div>
                <svg className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform" style={{ color: '#737373' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </Link>

            <Link
              to="/teams"
              className="group block"
            >
              <div className="flex items-center justify-between p-4 rounded-lg transition-all duration-200 hover:bg-card">
                <div>
                  <h3 className="text-white font-semibold mb-1">Teams</h3>
                  <p className="text-sm" style={{ color: '#737373' }}>explore 12 cdl teams and rosters</p>
                </div>
                <svg className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform" style={{ color: '#737373' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </Link>

            <Link
              to="/kd-stats"
              className="group block"
            >
              <div className="flex items-center justify-between p-4 rounded-lg transition-all duration-200 hover:bg-card">
                <div>
                  <h3 className="text-white font-semibold mb-1">K/D Statistics</h3>
                  <p className="text-sm" style={{ color: '#737373' }}>player kill/death ratios across all tournaments</p>
                </div>
                <svg className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform" style={{ color: '#737373' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </Link>

            <Link
              to="/transfers"
              className="group block"
            >
              <div className="flex items-center justify-between p-4 rounded-lg transition-all duration-200 hover:bg-card">
                <div>
                  <h3 className="text-white font-semibold mb-1">Transfers</h3>
                  <p className="text-sm" style={{ color: '#737373' }}>track player movements and roster changes</p>
                </div>
                <svg className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform" style={{ color: '#737373' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </Link>
          </div>


          {/* Footer - minimal */}
          <div className="mt-16 pb-8 text-center">
            <p className="text-sm" style={{ color: '#737373' }}>
              2025 season • black ops 6
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home; 