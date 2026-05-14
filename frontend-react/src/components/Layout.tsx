import { Outlet, Link } from "react-router-dom";

const Layout = () => {
  return (
    <div className="min-h-screen bg-white text-black">
      {/* Header */}
      <header className="border-b border-gray-300">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <Link
              to="/"
              className="text-xl font-bold tracking-tight hover:text-black transition-colors text-black"
            >
              CDLYTICS
            </Link>

            <nav className="flex space-x-8">
              <Link
                to="/players"
                className="text-sm uppercase tracking-wider text-[#6B7280] hover:text-black transition-colors"
              >
                Players
              </Link>
              <Link
                to="/teams"
                className="text-sm uppercase tracking-wider text-[#6B7280] hover:text-black transition-colors"
              >
                Teams
              </Link>
              <Link
                to="/stats"
                className="text-sm uppercase tracking-wider text-[#6B7280] hover:text-black transition-colors"
              >
                Stats
              </Link>
              <Link
                to="/transfers"
                className="text-sm uppercase tracking-wider text-[#6B7280] hover:text-black transition-colors"
              >
                Transfers
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
      <footer className="border-t border-gray-300 mt-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <p className="text-center text-[#6B7280] text-sm">CDLytics © 2025</p>
        </div>
      </footer>
    </div>
  );
};

export default Layout;
