import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import { getPlayerAvatar } from "../utils/avatarAssets";
import type { Player, PaginatedResponse } from "../types";

const LIMIT = 25;

const Players = () => {
  const [page, setPage] = useState(1);
  const [searchInput, setSearchInput] = useState("");
  const [search, setSearch] = useState("");

  useEffect(() => {
    const timer = setTimeout(() => {
      setSearch(searchInput);
      setPage(1);
    }, 300);
    return () => clearTimeout(timer);
  }, [searchInput]);

  const params = new URLSearchParams({
    page: String(page),
    limit: String(LIMIT),
  });
  if (search) params.set("search", search);

  const { data, loading, error } = useApi<PaginatedResponse<Player>>(
    `/api/v1/players?${params}`,
  );

  const players = data?.data ?? [];
  const meta = data?.pagination;

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
      <div className="mb-8 flex items-end justify-between gap-4">
        <div>
          <p className="text-xs uppercase tracking-widest text-[#737373] mb-2">
            Roster
          </p>
          <h1 className="font-grotesk text-3xl font-bold text-white">
            PLAYERS
          </h1>
        </div>
        <input
          type="text"
          placeholder="Search gamertag..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          className="bg-transparent border border-[#1a1a1a] text-white text-sm px-3 py-2 placeholder-[#737373] focus:outline-none focus:border-[#333] w-52"
        />
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
            {players.map((player) => (
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

      {/* Pagination controls */}
      <div className="mt-4 flex items-center justify-between">
        <p className="text-[#737373] text-xs">
          {meta
            ? `${meta.total} players — page ${meta.page} of ${meta.total_pages}`
            : `${players.length} players`}
        </p>

        {meta && meta.total_pages > 1 && (
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => setPage((p) => p - 1)}
              disabled={page === 1}
              className="px-3 py-1 text-xs text-[#737373] border border-[#1a1a1a] hover:text-white hover:border-[#333] disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
            >
              Prev
            </button>
            <button
              type="button"
              onClick={() => setPage((p) => p + 1)}
              disabled={page >= meta.total_pages}
              className="px-3 py-1 text-xs text-[#737373] border border-[#1a1a1a] hover:text-white hover:border-[#333] disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
            >
              Next
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default Players;
