import React, { useState } from 'react';
import { Link, Outlet } from 'react-router-dom';

const Layout: React.FC = () => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const toggleMobileMenu = () => {
    setMobileMenuOpen(!mobileMenuOpen);
  };

  const closeMobileMenu = () => {
    setMobileMenuOpen(false);
  };

  return (
    <div className="min-h-screen bg-black">
      {/* Header */}
      <header className="bg-black border-b border-gray-800 relative z-50 pt-safe">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16 sm:h-20">
            {/* Logo */}
            <div className="flex items-center">
              <Link to="/" className="text-2xl sm:text-3xl font-black text-white tracking-tight">
                CDLYTICS
              </Link>
            </div>

            {/* Desktop Navigation */}
            <nav className="hidden md:flex space-x-8 lg:space-x-12">
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
              <button 
                onClick={toggleMobileMenu}
                className="text-white hover:text-gray-300 p-2 transition-colors duration-200"
                aria-label="Toggle mobile menu"
              >
                <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  {mobileMenuOpen ? (
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  ) : (
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                  )}
                </svg>
              </button>
            </div>
          </div>

                      {/* Mobile Navigation Menu */}
            {mobileMenuOpen && (
              <div className="mobile-menu md:hidden">
                <div className="mobile-menu-content">
                  <Link
                    to="/"
                    className="mobile-menu-link"
                    onClick={closeMobileMenu}
                  >
                    HOME
                  </Link>
                  <Link
                    to="/teams"
                    className="mobile-menu-link"
                    onClick={closeMobileMenu}
                  >
                    TEAMS
                  </Link>
                  <Link
                    to="/players"
                    className="mobile-menu-link"
                    onClick={closeMobileMenu}
                  >
                    PLAYERS
                  </Link>
                  <Link
                    to="/kd-stats"
                    className="mobile-menu-link"
                    onClick={closeMobileMenu}
                  >
                    KD STATS
                  </Link>
                  <Link
                    to="/transfers"
                    className="mobile-menu-link"
                    onClick={closeMobileMenu}
                  >
                    TRANSFERS
                  </Link>
                </div>
              </div>
            )}
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-8 pb-safe">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-black border-t border-gray-800 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
          <div className="text-center text-gray-400 text-xs sm:text-sm uppercase tracking-wider">
            <p>&copy; 2025 CDLYTICS. ALL RIGHTS RESERVED.</p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Layout; 