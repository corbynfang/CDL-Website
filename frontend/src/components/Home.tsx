import { Link } from "react-router-dom";

const Home = () => {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <div className="text-center mb-16">
        <h1 className="text-5xl font-bold mb-4">CDLYTICS</h1>
        <p className="text-xl text-gray-400">
          Call of Duty League Statistics & Analytics
        </p>
      </div>

      <div className="grid md:grid-cols-3 gap-8 max-w-4xl mx-auto">
        <Link
          to="/players"
          className="p-8 hover:bg-gray-950 hover:border-gray-700 border border-transparent transition-all"
        >
          <h2 className="text-2xl font-bold mb-2 text-white">PLAYERS</h2>
          <p className="text-gray-400">
            View player statistics and performance data
          </p>
        </Link>

        <Link
          to="/teams"
          className="p-8 hover:bg-gray-950 hover:border-gray-700 border border-transparent transition-all"
        >
          <h2 className="text-2xl font-bold mb-2 text-white">TEAMS</h2>
          <p className="text-gray-400">Browse teams and roster information</p>
        </Link>

        <Link
          to="/stats"
          className="p-8 hover:bg-gray-950 hover:border-gray-700 border border-transparent transition-all"
        >
          <h2 className="text-2xl font-bold mb-2 text-white">STATS</h2>
          <p className="text-gray-400">K/D leaderboards and rankings</p>
        </Link>
      </div>
    </div>
  );
};

export default Home;
