import { Outlet, Link } from "react-router-dom";

const Layout = () => {
  return (
    <div className="min-h-screen bg-black text-white">
      {/* Header */}
      <header className="border-b border-gray-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <Link
              to="/"
              className="text-xl font-bold tracking-tight hover:text-white transition-colors"
            >
              CDLYTICS
            </Link>

            <nav className="flex space-x-8">
              <Link
                to="/players"
                className="text-sm uppercase tracking-wider hover:text-white transition-colors"
              >
                Players
              </Link>
              <Link
                to="/teams"
                className="text-sm uppercase tracking-wider hover:text-white transition-colors"
              >
                Teams
              </Link>
              <Link
                to="/stats"
                className="text-sm uppercase tracking-wider hover:text-white transition-colors"
              >
                Stats
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main>
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="border-t border-gray-800 mt-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <p className="text-center text-gray-400 text-sm">CDLytics © 2025</p>
        </div>
      </footer>
    </div>
  );
};

export default Layout;
