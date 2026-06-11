import { useState } from "react";
import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import { getPlayerAvatar } from "../utils/assets";
import type { Season } from "../types";
import PageMeta from "./PageMeta";

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

const kdColor = (kd: number) => {
  if (kd >= 1.3) return "text-green-400";
  if (kd < 1.0) return "text-[#737373]";
  return "text-white";
};

const Stats = () => {
  const [selectedSeasonId, setSelectedSeasonId] = useState<string>("");

  const { data: seasons } = useApi<Season[]>("/api/v1/seasons");

  const statsUrl = selectedSeasonId
    ? `/api/v1/stats/all-kd-by-tournament?season_id=${selectedSeasonId}`
    : "/api/v1/stats/all-kd-by-tournament";

  const { data: statsData, loading, error } = useApi<StatsResponse>(statsUrl);

  const players = [...(statsData?.players ?? [])].sort(
    (a, b) => b.season_kd - a.season_kd,
  );

  const selectedSeason = seasons?.find(
    (s) => String(s.id) === selectedSeasonId,
  );

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <PageMeta
        title="CDL K/D Rankings"
        description="Call of Duty League K/D leaderboards by season. See which players have the highest kill/death ratios across Hardpoint, Search & Destroy, and Control."
        canonical="/stats"
      />
      <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4 mb-8">
        <div>
          <p className="text-xs uppercase tracking-widest text-[#737373] mb-2">
            Leaderboard
          </p>
          <h1 className="font-grotesk text-3xl font-bold text-white">
            K/D RANKINGS
          </h1>
          {selectedSeason && (
            <p className="text-[#737373] text-sm mt-1">{selectedSeason.name}</p>
          )}
        </div>

        <select
          value={selectedSeasonId}
          onChange={(e) => setSelectedSeasonId(e.target.value)}
          className="bg-[#111111] border border-[#1a1a1a] px-3 py-2 text-xs text-[#a3a3a3] focus:outline-none focus:border-[#2a2a2a] uppercase tracking-wider"
        >
          <option value="">All Seasons</option>
          {seasons?.map((s) => (
            <option key={s.id} value={String(s.id)}>
              {s.game_title}
            </option>
          ))}
        </select>
      </div>

      {loading && <p className="text-[#737373] text-sm">Loading...</p>}
      {error && <p className="text-[#737373] text-sm">Error: {error}</p>}

      {!loading && !error && (
        <div className="border border-[#1a1a1a] overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-[#1a1a1a] bg-[#111111]">
                <th className="text-left px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium w-10">
                  #
                </th>
                <th className="text-left px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium">
                  Player
                </th>
                <th className="text-left px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium hidden sm:table-cell">
                  Team
                </th>
                <th className="text-right px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium">
                  K/D
                </th>
                <th className="text-right px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium hidden md:table-cell">
                  Kills
                </th>
                <th className="text-right px-4 py-3 text-[#737373] text-xs uppercase tracking-widest font-medium hidden md:table-cell">
                  Deaths
                </th>
              </tr>
            </thead>
            <tbody>
              {players.map((player, index) => (
                <tr
                  key={player.player_id}
                  className="border-b border-[#1a1a1a] hover:bg-[#111111] transition-colors"
                >
                  <td className="px-4 py-3 text-[#737373] font-mono text-xs w-10">
                    {index + 1}
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      <img
                        src={getPlayerAvatar(player.gamertag)}
                        alt={player.gamertag}
                        className="w-7 h-7 rounded-full object-cover opacity-90"
                      />
                      <Link
                        to={`/players/${player.player_id}`}
                        className="font-medium text-white hover:text-[#a3a3a3] transition-colors text-sm"
                      >
                        {player.gamertag}
                      </Link>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-[#737373] text-xs tracking-wider hidden sm:table-cell">
                    {player.team_abbr || "—"}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <span
                      className={`font-mono font-bold text-sm ${kdColor(player.season_kd)}`}
                    >
                      {player.season_kd?.toFixed(2) ?? "0.00"}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right text-[#737373] font-mono text-xs hidden md:table-cell">
                    {player.season_kills?.toLocaleString() ?? "0"}
                  </td>
                  <td className="px-4 py-3 text-right text-[#737373] font-mono text-xs hidden md:table-cell">
                    {player.season_deaths?.toLocaleString() ?? "0"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {players.length === 0 && (
            <p className="text-center text-[#737373] py-16 text-sm">
              No stats available for this season.
            </p>
          )}
        </div>
      )}

      <p className="mt-4 text-[#737373] text-xs">{players.length} players</p>
    </div>
  );
};

export default Stats;
