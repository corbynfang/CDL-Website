import React from 'react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => {
  return (
    <div className="min-h-screen bg-black">
      {/* Hero Section */}
      <div className="relative h-screen flex items-center justify-center overflow-hidden">
        {/* Video Background */}
        <div className="absolute inset-0 z-0">
          <video
            autoPlay
            loop
            muted
            playsInline
            className="hero-video"
          >
            <source src="/src/assets/video/webvideo.mp4" type="video/mp4" />
            {/* Fallback for browsers that don't support video */}
            <div className="video-fallback"></div>
          </video>
          {/* Overlay to ensure text readability */}
          <div className="hero-overlay"></div>
        </div>
        
        {/* Content */}
        <div className="hero-content text-center px-4 max-w-6xl mx-auto">
          <h1 className="text-6xl md:text-8xl font-bold text-white mb-8 tracking-tight">
            CDL <span className="text-red-500">STATS</span>
          </h1>
          <p className="text-xl md:text-2xl text-gray-300 max-w-4xl mx-auto mb-12 leading-relaxed">
            The definitive platform for Call of Duty League statistics, player analytics, and competitive insights.
          </p>
          <div className="flex flex-col sm:flex-row gap-6 justify-center">
            <Link
              to="/teams"
              className="bg-red-500 hover:bg-red-600 text-white font-bold py-4 px-8 text-lg transition-all duration-300 transform hover:scale-105"
            >
              EXPLORE TEAMS
            </Link>
            <Link
              to="/players"
              className="border-2 border-white text-white hover:bg-white hover:text-black font-bold py-4 px-8 text-lg transition-all duration-300"
            >
              VIEW PLAYERS
            </Link>
          </div>
        </div>
      </div>

      {/* Stats Section */}
      <div className="py-20 bg-gray-900">
        <div className="max-w-7xl mx-auto px-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
            <div>
              <div className="text-4xl md:text-5xl font-bold text-red-500 mb-2">12</div>
              <div className="text-gray-400 text-sm uppercase tracking-wider">CDL Teams</div>
            </div>
            <div>
              <div className="text-4xl md:text-5xl font-bold text-red-500 mb-2">48</div>
              <div className="text-gray-400 text-sm uppercase tracking-wider">Active Players</div>
            </div>
            <div>
              <div className="text-4xl md:text-5xl font-bold text-red-500 mb-2">5</div>
              <div className="text-gray-400 text-sm uppercase tracking-wider">Major Events</div>
            </div>
            <div>
              <div className="text-4xl md:text-5xl font-bold text-red-500 mb-2">24/7</div>
              <div className="text-gray-400 text-sm uppercase tracking-wider">Live Updates</div>
            </div>
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="py-20 bg-black">
        <div className="max-w-7xl mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold text-white mb-6">PLATFORM FEATURES</h2>
            <p className="text-xl text-gray-400 max-w-3xl mx-auto">
              Comprehensive analytics and insights for the Call of Duty League community
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            <div className="bg-gray-900 p-8 hover:bg-gray-800 transition-all duration-300">
              <div className="w-16 h-16 bg-red-500 rounded-lg flex items-center justify-center mb-6">
                <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <h3 className="text-2xl font-bold text-white mb-4">TEAM ANALYTICS</h3>
              <p className="text-gray-400 leading-relaxed">
                Deep dive into team performance, roster changes, and strategic insights across all CDL organizations.
              </p>
            </div>

            <div className="bg-gray-900 p-8 hover:bg-gray-800 transition-all duration-300">
              <div className="w-16 h-16 bg-red-500 rounded-lg flex items-center justify-center mb-6">
                <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <h3 className="text-2xl font-bold text-white mb-4">PLAYER PROFILES</h3>
              <p className="text-gray-400 leading-relaxed">
                Individual player statistics, K/D ratios, tournament performance, and career progression tracking.
              </p>
            </div>

            <div className="bg-gray-900 p-8 hover:bg-gray-800 transition-all duration-300">
              <div className="w-16 h-16 bg-red-500 rounded-lg flex items-center justify-center mb-6">
                <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
              <h3 className="text-2xl font-bold text-white mb-4">LIVE STATISTICS</h3>
              <p className="text-gray-400 leading-relaxed">
                Real-time match data, tournament results, and up-to-the-minute performance metrics.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* CTA Section */}
      <div className="py-20 bg-gradient-to-r from-red-500 to-red-600">
        <div className="max-w-4xl mx-auto text-center px-4">
          <h2 className="text-4xl md:text-5xl font-bold text-white mb-6">READY TO EXPLORE?</h2>
          <p className="text-xl text-red-100 mb-8">
            Dive into the world of competitive Call of Duty and discover what makes the CDL the pinnacle of esports.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              to="/teams"
              className="bg-white text-red-600 hover:bg-gray-100 font-bold py-4 px-8 text-lg transition-all duration-300"
            >
              BROWSE TEAMS
            </Link>
            <Link
              to="/players"
              className="border-2 border-white text-white hover:bg-white hover:text-red-600 font-bold py-4 px-8 text-lg transition-all duration-300"
            >
              VIEW PLAYERS
            </Link>
          </div>
        </div>
      </div>

      {/* Season Info */}
      <div className="py-16 bg-gray-900">
        <div className="max-w-6xl mx-auto px-4">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">CDL 2025 SEASON</h2>
              <div className="space-y-4">
                <div className="flex items-center">
                  <div className="w-3 h-3 bg-red-500 rounded-full mr-4"></div>
                  <span className="text-gray-300">Call of Duty: Black Ops 6</span>
                </div>
                <div className="flex items-center">
                  <div className="w-3 h-3 bg-red-500 rounded-full mr-4"></div>
                  <span className="text-gray-300">12 Professional Teams</span>
                </div>
                <div className="flex items-center">
                  <div className="w-3 h-3 bg-red-500 rounded-full mr-4"></div>
                  <span className="text-gray-300">5 Major Tournaments</span>
                </div>
                <div className="flex items-center">
                  <div className="w-3 h-3 bg-red-500 rounded-full mr-4"></div>
                  <span className="text-gray-300">Championship Weekend</span>
                </div>
              </div>
            </div>
            <div className="text-center">
              <div className="text-6xl font-bold text-red-500 mb-2">2025</div>
              <div className="text-gray-400 uppercase tracking-wider">Season Active</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home; 