import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";

interface PlayerStats {
  player_id: number;
  gamertag: string;
  team_abbr: string;
  season_kd: number;
  season_kills: number;
  season_deaths: number;
  season_assists: number;
}

const Stats = () => {
  const {
    data: statsData,
    loading,
    error,
  } = useApi<any>("/api/v1/stats/all-kd-by-tournament");

  if (loading) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#6B7280]">Loading statistics...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#555555]">Error: {error}</p>
          <p className="text-[#6B7280] text-sm mt-2">
            Backend: http://localhost:8080/api/v1/stats/all-kd-by-tournament
          </p>
        </div>
      </div>
    );
  }

  // Handle the correct response structure
  const players: PlayerStats[] = statsData?.players || [];

  // Sort by K/D ratio
  const sortedPlayers = [...players].sort((a, b) => b.season_kd - a.season_kd);

  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <h1 className="text-4xl font-bold mb-8 pb-4 text-black">K/D LEADERBOARD</h1>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-300">
                <th className="text-left py-4 text-[#6B7280] text-sm uppercase tracking-wider w-16">
                  Rank
                </th>
                <th className="text-left py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Player
                </th>
                <th className="text-left py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Team
                </th>
                <th className="text-right py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  K/D
                </th>
                <th className="text-right py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Kills
                </th>
                <th className="text-right py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Deaths
                </th>
                <th className="text-right py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Assists
                </th>
              </tr>
            </thead>
            <tbody>
              {sortedPlayers.map((player, index) => (
                <tr key={player.player_id} className="border-b border-gray-300">
                  <td className="py-4 text-[#6B7280]">{index + 1}</td>
                  <td className="py-4">
                    <Link
                      to={`/players/${player.player_id}`}
                      className="font-semibold text-black"
                    >
                      {player.gamertag}
                    </Link>
                  </td>
                  <td className="py-4 text-[#6B7280]">
                    {player.team_abbr || "—"}
                  </td>
                  <td className="py-4 text-right font-bold text-black">
                    {player.season_kd?.toFixed(2) || "0.00"}
                  </td>
                  <td className="py-4 text-right text-[#6B7280]">
                    {player.season_kills?.toLocaleString() || "0"}
                  </td>
                  <td className="py-4 text-right text-[#6B7280]">
                    {player.season_deaths?.toLocaleString() || "0"}
                  </td>
                  <td className="py-4 text-right text-[#6B7280]">
                    {player.season_assists?.toLocaleString() || "0"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <p className="mt-8 text-[#6B7280] text-sm">
          Total Players: {sortedPlayers.length}
        </p>
      </div>
    </div>
  );
};

export default Stats;
