import React from 'react';
import { Link, Outlet } from 'react-router-dom';

const Layout: React.FC = () => {
  return (
    <div className="min-h-screen bg-black">
      {/* Header */}
      <header className="bg-black border-b border-gray-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-20">
            {/* Logo */}
            <div className="flex items-center">
              <Link to="/" className="text-3xl font-black text-white tracking-tight">
                CDL STATS
              </Link>
            </div>

            {/* Navigation */}
            <nav className="hidden md:flex space-x-12">
              <Link
                to="/"
                className="nav-link text-sm"
              >
                HOME
              </Link>
              <Link
                to="/teams"
                className="nav-link text-sm"
              >
                TEAMS
              </Link>
              <Link
                to="/players"
                className="nav-link text-sm"
              >
                PLAYERS
              </Link>
              <Link
                to="/kd-stats"
                className="nav-link text-sm"
              >
                KD STATS
              </Link>
              <Link
                to="/transfers"
                className="nav-link text-sm"
              >
                TRANSFERS
              </Link>
            </nav>

            {/* Mobile menu button */}
            <div className="md:hidden">
              <button className="text-white hover:text-gray-300 p-2 transition-colors duration-200">
                <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-black border-t border-gray-800 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="text-center text-gray-400 text-sm uppercase tracking-wider">
            <p>&copy; 2025 CDL STATS. ALL RIGHTS RESERVED.</p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Layout; 