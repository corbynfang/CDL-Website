import { Link } from "react-router-dom";

const CARDS = [
  { to: "/players", title: "PLAYERS", desc: "Player statistics and performance data" },
  { to: "/teams", title: "TEAMS", desc: "Browse teams and roster information" },
  { to: "/stats", title: "STATS", desc: "K/D leaderboards and rankings" },
  { to: "/transfers", title: "TRANSFERS", desc: "Roster moves and signing history" },
];

const Home = () => {
  return (
    <div className="min-h-screen bg-[#0a0a0a]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-28 pb-20">
        <div className="mb-20">
          <p className="text-xs uppercase tracking-widest text-[#737373] mb-4">
            Call of Duty League
          </p>
          <h1 className="font-grotesk text-6xl font-bold tracking-tight text-white mb-4">
            CDLYTICS
          </h1>
          <p className="text-[#a3a3a3] text-lg max-w-md">
            Statistics, analytics, and roster data for the Call of Duty League.
          </p>
        </div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-3">
          {CARDS.map(({ to, title, desc }) => (
            <Link
              key={to}
              to={to}
              className="group p-6 bg-[#111111] border border-[#1a1a1a] hover:border-[#2a2a2a] hover:bg-[#161616] transition-all"
            >
              <h2 className="font-grotesk text-sm font-bold tracking-widest text-white mb-2">
                {title}
              </h2>
              <p className="text-[#737373] text-sm leading-relaxed">{desc}</p>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Home;
