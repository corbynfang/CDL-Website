import { Link } from "react-router-dom";

const Home = () => {
  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center mb-16">
          <h1 className="text-5xl font-bold mb-4 text-black">CDLYTICS</h1>
          <p className="text-xl text-[#6B7280]">
            Call of Duty League Statistics & Analytics
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-8 max-w-4xl mx-auto">
          <Link
            to="/players"
            className="p-8 bg-[#F4F4F5] border border-transparent shadow-md shadow-[rgba(0,0,0,0.1)]"
          >
            <h2 className="text-2xl font-bold mb-2 text-black">PLAYERS</h2>
            <p className="text-[#6B7280]">
              View player statistics and performance data
            </p>
          </Link>

          <Link
            to="/teams"
            className="p-8 bg-[#F4F4F5] border border-transparent shadow-md shadow-[rgba(0,0,0,0.1)]"
          >
            <h2 className="text-2xl font-bold mb-2 text-black">TEAMS</h2>
            <p className="text-[#6B7280]">Browse teams and roster information</p>
          </Link>

          <Link
            to="/stats"
            className="p-8 bg-[#F4F4F5] border border-transparent shadow-md shadow-[rgba(0,0,0,0.1)]"
          >
            <h2 className="text-2xl font-bold mb-2 text-black">STATS</h2>
            <p className="text-[#6B7280]">K/D leaderboards and rankings</p>
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Home;
