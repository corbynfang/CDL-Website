import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import type { Player } from "../types";

const Players = () => {
  const { data: players, loading, error } = useApi<Player[]>("/api/v1/players");

  if (loading) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#6B7280]">Loading players...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#555555]">Error: {error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <h1 className="text-4xl font-bold mb-8 pb-4 text-black">PLAYERS</h1>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-300">
                <th className="text-left py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Gamertag
                </th>
                <th className="text-left py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Name
                </th>
                <th className="text-left py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Country
                </th>
                <th className="text-left py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Role
                </th>
                <th className="text-center py-4 text-[#6B7280] text-sm uppercase tracking-wider">
                  Status
                </th>
              </tr>
            </thead>
            <tbody>
              {players?.map((player) => (
                <tr key={player.id} className="border-b border-gray-300">
                  <td className="py-4">
                    <Link to={`/players/${player.id}`} className="font-semibold text-black">
                      {player.gamertag}
                    </Link>
                  </td>
                  <td className="py-4 text-[#6B7280]">
                    {player.first_name && player.last_name
                      ? `${player.first_name} ${player.last_name}`
                      : "—"}
                  </td>
                  <td className="py-4 text-[#6B7280]">{player.country || "—"}</td>
                  <td className="py-4 text-[#6B7280]">{player.role || "—"}</td>
                  <td className="py-4 text-center">
                    <span
                      className={
                        player.is_active ? "text-black" : "text-[#6B7280]"
                      }
                    >
                      {player.is_active ? "Active" : "Inactive"}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <p className="mt-8 text-[#6B7280] text-sm">
          Total Players: {players?.length || 0}
        </p>
      </div>
    </div>
  );
};

export default Players;
