import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import { getPlayerAvatar } from "../utils/assets";
import type { Player } from "../types";

const Players = () => {
  const { data: players, loading, error } = useApi<Player[]>("/api/v1/players");

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Loading players...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="mb-8">
        <p className="text-xs uppercase tracking-widest text-[#737373] mb-2">Roster</p>
        <h1 className="font-grotesk text-3xl font-bold text-white">PLAYERS</h1>
      </div>

      <div className="border border-[#1a1a1a] overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-[#1a1a1a] bg-[#111111]">
              <th className="text-left px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium">
                Gamertag
              </th>
              <th className="text-left px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium">
                Name
              </th>
              <th className="text-left px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium hidden sm:table-cell">
                Country
              </th>
              <th className="text-left px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium hidden md:table-cell">
                Role
              </th>
              <th className="text-center px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium">
                Status
              </th>
            </tr>
          </thead>
          <tbody>
            {players?.map((player) => (
              <tr
                key={player.id}
                className="border-b border-[#1a1a1a] hover:bg-[#111111] transition-colors"
              >
                <td className="px-4 py-3">
                  <div className="flex items-center gap-3">
                    <img
                      src={getPlayerAvatar(player.gamertag)}
                      alt={player.gamertag}
                      className="w-7 h-7 object-cover rounded-full opacity-90"
                    />
                    <Link
                      to={`/players/${player.id}`}
                      className="font-medium text-white hover:text-[#a3a3a3] transition-colors text-sm"
                    >
                      {player.gamertag}
                    </Link>
                  </div>
                </td>
                <td className="px-4 py-3 text-[#a3a3a3] text-sm">
                  {player.first_name && player.last_name
                    ? `${player.first_name} ${player.last_name}`
                    : "—"}
                </td>
                <td className="px-4 py-3 text-[#737373] text-sm hidden sm:table-cell">
                  {player.country || "—"}
                </td>
                <td className="px-4 py-3 text-[#737373] text-sm hidden md:table-cell">
                  {player.role || "—"}
                </td>
                <td className="px-4 py-3 text-center">
                  <span
                    className={`text-xs font-medium px-2 py-0.5 ${
                      player.is_active
                        ? "text-green-400 bg-green-400/10"
                        : "text-[#737373] bg-white/5"
                    }`}
                  >
                    {player.is_active ? "Active" : "Inactive"}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <p className="mt-4 text-[#737373] text-xs">
        {players?.length || 0} players
      </p>
    </div>
  );
};

export default Players;
