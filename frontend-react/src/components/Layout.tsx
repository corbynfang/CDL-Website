import { Outlet, Link, useLocation } from "react-router-dom";

const NAV_LINKS = [
  { to: "/players", label: "Players" },
  { to: "/teams", label: "Teams" },
  { to: "/stats", label: "Stats" },
  { to: "/transfers", label: "Transfers" },
];

const Layout = () => {
  const { pathname } = useLocation();

  return (
    <div className="min-h-screen bg-[#0a0a0a] text-[#f5f5f5]">
      <header className="sticky top-0 z-50 border-b border-[#1a1a1a] bg-[#0a0a0a]/95 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-14">
            <Link
              to="/"
              className="font-grotesk text-base font-bold tracking-tight text-white hover:text-white/80 transition-colors"
            >
              CDLYTICS
            </Link>

            <nav className="flex space-x-6">
              {NAV_LINKS.map(({ to, label }) => {
                const active = pathname.startsWith(to);
                return (
                  <Link
                    key={to}
                    to={to}
                    className={`text-xs uppercase tracking-widest transition-colors ${
                      active ? "text-white" : "text-[#737373] hover:text-[#a3a3a3]"
                    }`}
                  >
                    {label}
                  </Link>
                );
              })}
            </nav>
          </div>
        </div>
      </header>

      <main>
        <Outlet />
      </main>

      <footer className="border-t border-[#1a1a1a] mt-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <p className="text-center text-[#737373] text-xs tracking-wider">
            CDLYTICS · CDL STATISTICS & ANALYTICS
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Layout;
