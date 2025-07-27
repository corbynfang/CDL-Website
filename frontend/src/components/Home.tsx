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
            onError={(e) => {
              console.error('Video error:', e);
              // Hide video and show fallback if video fails
              const videoElement = e.target as HTMLVideoElement;
              if (videoElement) {
                videoElement.style.display = 'none';
              }
            }}
            onLoadStart={() => console.log('Video loading started')}
            onCanPlay={() => console.log('Video can play')}
          >
            <source src="/webvideo.mp4" type="video/mp4" />
          </video>
          {/* Fallback background - shows if video fails to load */}
          <div className="video-fallback"></div>
          {/* Overlay to ensure text readability */}
          <div className="hero-overlay"></div>
        </div>
        
        {/* Content */}
        <div className="hero-content text-center px-4 max-w-6xl mx-auto">
          <h1 className="text-hero text-white mb-8 tracking-tight">
            CDL <span className="text-red-500">STATS</span>
          </h1>
          <p className="text-subheading text-gray-300 max-w-4xl mx-auto mb-12 leading-relaxed">
            THE DEFINITIVE PLATFORM FOR CALL OF DUTY LEAGUE STATISTICS, PLAYER ANALYTICS, AND COMPETITIVE INSIGHTS.
          </p>
          <div className="flex flex-col sm:flex-row gap-6 justify-center">
            <Link
              to="/teams"
              className="btn-primary text-lg px-8 py-4"
            >
              EXPLORE TEAMS
            </Link>
            <Link
              to="/players"
              className="btn-secondary text-lg px-8 py-4"
            >
              VIEW PLAYERS
            </Link>
          </div>
        </div>
      </div>

      {/* Stats Section */}
      <div className="py-20 bg-black">
        <div className="max-w-7xl mx-auto px-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
            <div>
              <div className="text-4xl md:text-5xl font-black text-white mb-2">12</div>
              <div className="text-gray-400 text-xs uppercase tracking-widest">CDL TEAMS</div>
            </div>
            <div>
              <div className="text-4xl md:text-5xl font-black text-white mb-2">48</div>
              <div className="text-gray-400 text-xs uppercase tracking-widest">ACTIVE PLAYERS</div>
            </div>
            <div>
              <div className="text-4xl md:text-5xl font-black text-white mb-2">5</div>
              <div className="text-gray-400 text-xs uppercase tracking-widest">MAJOR EVENTS</div>
            </div>
            <div>
              <div className="text-4xl md:text-5xl font-black text-white mb-2">24/7</div>
              <div className="text-gray-400 text-xs uppercase tracking-widest">LIVE UPDATES</div>
            </div>
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="py-20 bg-black">
        <div className="max-w-7xl mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-heading text-white mb-6">PLATFORM FEATURES</h2>
            <p className="text-subheading text-gray-400 max-w-3xl mx-auto">
              COMPREHENSIVE ANALYTICS AND INSIGHTS FOR THE CALL OF DUTY LEAGUE COMMUNITY
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            <div className="bg-black border border-gray-800 p-8 hover:border-white transition-all duration-300">
              <div className="w-16 h-16 bg-white rounded-none flex items-center justify-center mb-6">
                <svg className="w-8 h-8 text-black" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-4 uppercase tracking-wider">TEAM ANALYTICS</h3>
              <p className="text-gray-400 leading-relaxed text-sm">
                DEEP DIVE INTO TEAM PERFORMANCE, ROSTER CHANGES, AND STRATEGIC INSIGHTS ACROSS ALL CDL ORGANIZATIONS.
              </p>
            </div>

            <div className="bg-black border border-gray-800 p-8 hover:border-white transition-all duration-300">
              <div className="w-16 h-16 bg-white rounded-none flex items-center justify-center mb-6">
                <svg className="w-8 h-8 text-black" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-4 uppercase tracking-wider">PLAYER PROFILES</h3>
              <p className="text-gray-400 leading-relaxed text-sm">
                INDIVIDUAL PLAYER STATISTICS, K/D RATIOS, TOURNAMENT PERFORMANCE, AND CAREER PROGRESSION TRACKING.
              </p>
            </div>

            <div className="bg-black border border-gray-800 p-8 hover:border-white transition-all duration-300">
              <div className="w-16 h-16 bg-white rounded-none flex items-center justify-center mb-6">
                <svg className="w-8 h-8 text-black" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-4 uppercase tracking-wider">LIVE STATISTICS</h3>
              <p className="text-gray-400 leading-relaxed text-sm">
                REAL-TIME MATCH DATA, TOURNAMENT RESULTS, AND UP-TO-THE-MINUTE PERFORMANCE METRICS.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* CTA Section */}
      <div className="py-20 bg-white">
        <div className="max-w-4xl mx-auto text-center px-4">
          <h2 className="text-heading text-black mb-6">READY TO EXPLORE?</h2>
          <p className="text-subheading text-gray-600 mb-8">
            DIVE INTO THE WORLD OF COMPETITIVE CALL OF DUTY AND DISCOVER WHAT MAKES THE CDL THE PINNACLE OF ESPORTS.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              to="/teams"
              className="bg-black text-white hover:bg-gray-800 font-bold py-4 px-8 text-lg transition-all duration-300 uppercase tracking-wider"
            >
              BROWSE TEAMS
            </Link>
            <Link
              to="/players"
              className="border-2 border-black text-black hover:bg-black hover:text-white font-bold py-4 px-8 text-lg transition-all duration-300 uppercase tracking-wider"
            >
              VIEW PLAYERS
            </Link>
          </div>
        </div>
      </div>

      {/* Season Info */}
      <div className="py-16 bg-black">
        <div className="max-w-6xl mx-auto px-4">
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <h2 className="text-heading text-white mb-6">CDL 2025 SEASON</h2>
              <div className="space-y-4">
                <div className="flex items-center">
                  <div className="w-3 h-3 bg-white rounded-none mr-4"></div>
                  <span className="text-gray-300 uppercase tracking-wider">CALL OF DUTY: BLACK OPS 6</span>
                </div>
                <div className="flex items-center">
                  <div className="w-3 h-3 bg-white rounded-none mr-4"></div>
                  <span className="text-gray-300 uppercase tracking-wider">12 PROFESSIONAL TEAMS</span>
                </div>
                <div className="flex items-center">
                  <div className="w-3 h-3 bg-white rounded-none mr-4"></div>
                  <span className="text-gray-300 uppercase tracking-wider">5 MAJOR TOURNAMENTS</span>
                </div>
                <div className="flex items-center">
                  <div className="w-3 h-3 bg-white rounded-none mr-4"></div>
                  <span className="text-gray-300 uppercase tracking-wider">CHAMPIONSHIP WEEKEND</span>
                </div>
              </div>
            </div>
            <div className="text-center">
              <div className="text-6xl font-black text-white mb-2">2025</div>
              <div className="text-gray-400 uppercase tracking-widest">SEASON ACTIVE</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home; 