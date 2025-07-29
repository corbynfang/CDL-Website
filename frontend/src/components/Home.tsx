import React from 'react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => {
  return (
    <div className="min-h-screen bg-black">
      {/* Full Screen Video Background */}
      <div className="relative min-h-screen">
        {/* Video Background */}
        <div className="fixed inset-0 z-0">
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
        
        {/* Main Content Overlay */}
        <div className="relative z-10 min-h-screen flex flex-col">
          {/* Hero Section */}
          <div className="flex-1 flex items-center justify-center px-4 sm:px-6">
            <div className="text-center w-full max-w-4xl">
              <h1 className="mobile-hero-text text-white tracking-tight mb-4 sm:mb-6">
                CDL<span className="text-red-500">YTICS</span>
              </h1>
              <p className="mobile-subtitle-text text-gray-300 uppercase tracking-widest mb-6 sm:mb-8 px-4">
                CALL OF DUTY LEAGUE ANALYTICS
              </p>
              
              {/* CTA Buttons */}
              <div className="flex flex-col sm:flex-row gap-4 sm:gap-6 justify-center px-4">
                <Link
                  to="/teams"
                  className="mobile-btn btn-primary"
                >
                  EXPLORE TEAMS
                </Link>
                <Link
                  to="/players"
                  className="mobile-btn btn-secondary"
                >
                  VIEW PLAYERS
                </Link>
              </div>
            </div>
          </div>

          {/* Stats Section */}
          <div className="py-12 sm:py-16 md:py-20">
            <div className="max-w-7xl mx-auto px-4 sm:px-6">
              <div className="mobile-stats-grid text-center">
                <div className="mobile-stats-item">
                  <div className="mobile-stats-number">12</div>
                  <div className="mobile-stats-label">CDL TEAMS</div>
                </div>
                <div className="mobile-stats-item">
                  <div className="mobile-stats-number">48</div>
                  <div className="mobile-stats-label">ACTIVE PLAYERS</div>
                </div>
                <div className="mobile-stats-item">
                  <div className="mobile-stats-number">5</div>
                  <div className="mobile-stats-label">MAJOR EVENTS</div>
                </div>
                <div className="mobile-stats-item">
                  <div className="mobile-stats-number">24/7</div>
                  <div className="mobile-stats-label">LIVE UPDATES</div>
                </div>
              </div>
            </div>
          </div>

          {/* Features Section */}
          <div className="py-12 sm:py-16 md:py-20">
            <div className="max-w-7xl mx-auto px-4 sm:px-6">
              <div className="text-center mb-12 sm:mb-16">
                <h2 className="text-2xl sm:text-3xl md:text-heading text-white mb-4 sm:mb-6">PLATFORM FEATURES</h2>
                <p className="text-base sm:text-lg md:text-subheading text-gray-400 max-w-3xl mx-auto px-4">
                  COMPREHENSIVE ANALYTICS AND INSIGHTS FOR THE CALL OF DUTY LEAGUE COMMUNITY
                </p>
              </div>

              <div className="mobile-grid">
                <div className="bg-black bg-opacity-50 backdrop-blur-sm border border-gray-800 mobile-card hover:border-white transition-all duration-300">
                  <div className="w-12 h-12 sm:w-16 sm:h-16 bg-white rounded-none flex items-center justify-center mb-4 sm:mb-6">
                    <svg className="w-6 h-6 sm:w-8 sm:h-8 text-black" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                  </div>
                  <h3 className="text-lg sm:text-xl font-bold text-white mb-3 sm:mb-4 uppercase tracking-wider">TEAM ANALYTICS</h3>
                  <p className="text-gray-300 leading-relaxed text-sm sm:text-base">
                    DEEP DIVE INTO TEAM PERFORMANCE, ROSTER CHANGES, AND STRATEGIC INSIGHTS ACROSS ALL CDL ORGANIZATIONS.
                  </p>
                </div>

                <div className="bg-black bg-opacity-50 backdrop-blur-sm border border-gray-800 mobile-card hover:border-white transition-all duration-300">
                  <div className="w-12 h-12 sm:w-16 sm:h-16 bg-white rounded-none flex items-center justify-center mb-4 sm:mb-6">
                    <svg className="w-6 h-6 sm:w-8 sm:h-8 text-black" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                  </div>
                  <h3 className="text-lg sm:text-xl font-bold text-white mb-3 sm:mb-4 uppercase tracking-wider">PLAYER PROFILES</h3>
                  <p className="text-gray-300 leading-relaxed text-sm sm:text-base">
                    INDIVIDUAL PLAYER STATISTICS, K/D RATIOS, TOURNAMENT PERFORMANCE, AND CAREER PROGRESSION TRACKING.
                  </p>
                </div>

                <div className="bg-black bg-opacity-50 backdrop-blur-sm border border-gray-800 mobile-card hover:border-white transition-all duration-300">
                  <div className="w-12 h-12 sm:w-16 sm:h-16 bg-white rounded-none flex items-center justify-center mb-4 sm:mb-6">
                    <svg className="w-6 h-6 sm:w-8 sm:h-8 text-black" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                    </svg>
                  </div>
                  <h3 className="text-lg sm:text-xl font-bold text-white mb-3 sm:mb-4 uppercase tracking-wider">LIVE STATISTICS</h3>
                  <p className="text-gray-300 leading-relaxed text-sm sm:text-base">
                    REAL-TIME MATCH DATA, TOURNAMENT RESULTS, AND UP-TO-THE-MINUTE PERFORMANCE METRICS.
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* CTA Section */}
          <div className="py-12 sm:py-16 md:py-20">
            <div className="max-w-4xl mx-auto text-center px-4 sm:px-6">
              <h2 className="text-2xl sm:text-3xl md:text-heading text-white mb-4 sm:mb-6">READY TO EXPLORE?</h2>
              <p className="text-base sm:text-lg md:text-subheading text-gray-300 mb-6 sm:mb-8 px-4">
                DIVE INTO THE WORLD OF COMPETITIVE CALL OF DUTY AND DISCOVER WHAT MAKES THE CDL THE PINNACLE OF ESPORTS.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center px-4">
                <Link
                  to="/teams"
                  className="bg-white text-black hover:bg-gray-200 font-bold py-3 sm:py-4 px-6 sm:px-8 text-base sm:text-lg transition-all duration-300 uppercase tracking-wider"
                >
                  BROWSE TEAMS
                </Link>
                <Link
                  to="/players"
                  className="border-2 border-white text-white hover:bg-white hover:text-black font-bold py-3 sm:py-4 px-6 sm:px-8 text-base sm:text-lg transition-all duration-300 uppercase tracking-wider"
                >
                  VIEW PLAYERS
                </Link>
              </div>
            </div>
          </div>

          {/* Season Info */}
          <div className="py-12 sm:py-16">
            <div className="max-w-6xl mx-auto px-4 sm:px-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-8 sm:gap-12 items-center">
                <div>
                  <h2 className="text-2xl sm:text-3xl md:text-heading text-white mb-4 sm:mb-6">CDL 2025 SEASON</h2>
                  <div className="space-y-3 sm:space-y-4">
                    <div className="flex items-center">
                      <div className="w-2 h-2 sm:w-3 sm:h-3 bg-white rounded-none mr-3 sm:mr-4"></div>
                      <span className="text-gray-300 uppercase tracking-wider text-sm sm:text-base">CALL OF DUTY: BLACK OPS 6</span>
                    </div>
                    <div className="flex items-center">
                      <div className="w-2 h-2 sm:w-3 sm:h-3 bg-white rounded-none mr-3 sm:mr-4"></div>
                      <span className="text-gray-300 uppercase tracking-wider text-sm sm:text-base">12 PROFESSIONAL TEAMS</span>
                    </div>
                    <div className="flex items-center">
                      <div className="w-2 h-2 sm:w-3 sm:h-3 bg-white rounded-none mr-3 sm:mr-4"></div>
                      <span className="text-gray-300 uppercase tracking-wider text-sm sm:text-base">5 MAJOR TOURNAMENTS</span>
                    </div>
                    <div className="flex items-center">
                      <div className="w-2 h-2 sm:w-3 sm:h-3 bg-white rounded-none mr-3 sm:mr-4"></div>
                      <span className="text-gray-300 uppercase tracking-wider text-sm sm:text-base">CHAMPIONSHIP WEEKEND</span>
                    </div>
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-4xl sm:text-5xl md:text-6xl font-black text-white mb-2">2025</div>
                  <div className="text-gray-400 uppercase tracking-widest text-sm sm:text-base">SEASON ACTIVE</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home; 