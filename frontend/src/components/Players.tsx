import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import type { Player } from "../types";

const Players = () => {
  const { data: players, loading, error } = useApi<Player[]>("/api/v1/players");

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-gray-400">Loading players...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-red-400">Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <h1 className="text-4xl font-bold mb-8 pb-4">PLAYERS</h1>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-gray-800">
              <th className="text-left py-4 text-gray-400 text-sm uppercase tracking-wider">
                Gamertag
              </th>
              <th className="text-left py-4 text-gray-400 text-sm uppercase tracking-wider">
                Name
              </th>
              <th className="text-left py-4 text-gray-400 text-sm uppercase tracking-wider">
                Country
              </th>
              <th className="text-left py-4 text-gray-400 text-sm uppercase tracking-wider">
                Role
              </th>
              <th className="text-center py-4 text-gray-400 text-sm uppercase tracking-wider">
                Status
              </th>
            </tr>
          </thead>
          <tbody>
            {players?.map((player) => (
              <tr
                key={player.id}
                className="border-b border-gray-900 hover:bg-gray-950 transition-colors"
              >
                <td className="py-4">
                  <Link
                    to={`/players/${player.id}`}
                    className="font-semibold hover:text-gray-400"
                  >
                    {player.gamertag}
                  </Link>
                </td>
                <td className="py-4 text-gray-400">
                  {player.first_name && player.last_name
                    ? `${player.first_name} ${player.last_name}`
                    : "—"}
                </td>
                <td className="py-4 text-gray-400">{player.country || "—"}</td>
                <td className="py-4 text-gray-400">{player.role || "—"}</td>
                <td className="py-4 text-center">
                  <span
                    className={
                      player.is_active ? "text-white" : "text-gray-600"
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

      <p className="mt-8 text-gray-400 text-sm">
        Total Players: {players?.length || 0}
      </p>
    </div>
  );
};

export default Players;
