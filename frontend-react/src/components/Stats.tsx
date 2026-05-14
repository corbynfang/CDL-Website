import { useState } from "react";
import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import { getPlayerAvatar } from "../utils/assets";
import type { Season } from "../types";

interface PlayerStats {
  player_id: number;
  gamertag: string;
  team_abbr: string;
  season_kd: number;
  season_kills: number;
  season_deaths: number;
  season_assists: number;
}

interface StatsResponse {
  players: PlayerStats[];
  count: number;
}

const Stats = () => {
  const [selectedSeasonId, setSelectedSeasonId] = useState<string>("");

  const { data: seasons } = useApi<Season[]>("/api/v1/seasons");

  const statsUrl = selectedSeasonId
    ? `/api/v1/stats/all-kd-by-tournament?season_id=${selectedSeasonId}`
    : "/api/v1/stats/all-kd-by-tournament";

  const { data: statsData, loading, error } = useApi<StatsResponse>(statsUrl);

  const players = [...(statsData?.players ?? [])].sort(
    (a, b) => b.season_kd - a.season_kd
  );

  const selectedSeason = seasons?.find(
    (s) => String(s.id) === selectedSeasonId
  );

  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
          <div>
            <h1 className="text-4xl font-bold text-black">K/D LEADERBOARD</h1>
            {selectedSeason && (
              <p className="text-[#6B7280] mt-1">{selectedSeason.name}</p>
            )}
          </div>

          {/* Season selector */}
          <select
            value={selectedSeasonId}
            onChange={(e) => setSelectedSeasonId(e.target.value)}
            className="border border-gray-300 px-4 py-2 text-sm text-black bg-white focus:outline-none focus:border-black"
          >
            <option value="">All Seasons</option>
            {seasons?.map((s) => (
              <option key={s.id} value={String(s.id)}>
                {s.game_title}
              </option>
            ))}
          </select>
        </div>

        {loading && <p className="text-[#6B7280]">Loading statistics...</p>}

        {error && <p className="text-[#555555]">Error: {error}</p>}

        {!loading && !error && (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-300">
                  <th className="text-left py-4 text-[#6B7280] text-sm uppercase tracking-wider w-12">
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
                </tr>
              </thead>
              <tbody>
                {players.map((player, index) => (
                  <tr key={player.player_id} className="border-b border-gray-300">
                    <td className="py-4 text-[#6B7280] text-sm">{index + 1}</td>
                    <td className="py-4">
                      <div className="flex items-center gap-3">
                        <img
                          src={getPlayerAvatar(player.gamertag)}
                          alt={player.gamertag}
                          className="w-8 h-8 rounded-full object-cover"
                        />
                        <Link
                          to={`/players/${player.player_id}`}
                          className="font-semibold text-black hover:underline"
                        >
                          {player.gamertag}
                        </Link>
                      </div>
                    </td>
                    <td className="py-4 text-[#6B7280]">
                      {player.team_abbr || "—"}
                    </td>
                    <td className="py-4 text-right font-bold text-black">
                      {player.season_kd?.toFixed(2) ?? "0.00"}
                    </td>
                    <td className="py-4 text-right text-[#6B7280]">
                      {player.season_kills?.toLocaleString() ?? "0"}
                    </td>
                    <td className="py-4 text-right text-[#6B7280]">
                      {player.season_deaths?.toLocaleString() ?? "0"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {players.length === 0 && (
              <p className="text-center text-[#6B7280] py-12">
                No stats available for this season.
              </p>
            )}
          </div>
        )}

        <p className="mt-8 text-[#6B7280] text-sm">
          {players.length} players
        </p>
      </div>
    </div>
  );
};

export default Stats;
