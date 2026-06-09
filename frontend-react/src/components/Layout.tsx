import { Outlet, Link, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import AuthModal from "./auth/AuthModal";

const NAV_LINKS = [
  { to: "/events", label: "Events" },
  { to: "/players", label: "Players" },
  { to: "/teams", label: "Teams" },
  { to: "/stats", label: "Stats" },
  { to: "/transfers", label: "Transfers" },
];

const Layout = () => {
  const { pathname } = useLocation();
  const {
    session,
    openAuthModal,
    closeAuthModal,
    showAuthModal,
    needsProfileSetup,
    signOut,
  } = useAuth();

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

            <nav className="flex items-center space-x-6">
              {NAV_LINKS.map(({ to, label }) => {
                const active = pathname.startsWith(to);
                return (
                  <Link
                    key={to}
                    to={to}
                    className={`text-xs uppercase tracking-widest transition-colors ${
                      active
                        ? "text-white"
                        : "text-[#737373] hover:text-[#a3a3a3]"
                    }`}
                  >
                    {label}
                  </Link>
                );
              })}
              {session ? (
                <button
                  type="button"
                  onClick={() => signOut()}
                  className="text-xs uppercase tracking-widest text-[#737373] hover:text-white transition-colors"
                >
                  Sign Out
                </button>
              ) : (
                <button
                  type="button"
                  onClick={openAuthModal}
                  className="text-xs uppercase tracking-widest text-white border border-[#333] px-3 py-1 hover:border-[#555] transition-colors"
                >
                  Sign In
                </button>
              )}
            </nav>
          </div>
        </div>
      </header>

      <main>
        <Outlet />
      </main>
      {(showAuthModal || needsProfileSetup) && (
        <AuthModal onClose={closeAuthModal} />
      )}

      <footer className="border-t border-[#1a1a1a] mt-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-3">
          <p className="text-center text-[#737373] text-xs tracking-wider">
            CDLYTICS · INDEPENDENT CDL STATISTICS &amp; ANALYTICS
          </p>
          <p className="text-center text-[#4a4a4a] text-xs">
            Not affiliated with, endorsed by, or sponsored by Activision, Call
            of Duty League, Esports World Cup, or any listed team. All
            trademarks belong to their respective owners.
          </p>
          <div className="flex justify-center gap-6 text-[10px] uppercase tracking-widest text-[#4a4a4a]">
            <Link
              to="/privacy"
              className="hover:text-[#737373] transition-colors"
            >
              Privacy
            </Link>
            <Link
              to="/terms"
              className="hover:text-[#737373] transition-colors"
            >
              Terms
            </Link>
            <Link
              to="/disclaimer"
              className="hover:text-[#737373] transition-colors"
            >
              Disclaimer
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Layout;
