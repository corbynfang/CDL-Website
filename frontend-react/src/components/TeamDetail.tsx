import { useState } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import { getTeamLogo } from "../utils/logoAssets";
import { getPlayerAvatar } from "../utils/avatarAssets";
import type { Player, Team, Franchise } from "../types";

interface FranchiseResponse {
  franchise: Franchise;
  eras: Team[];
}

const gameLabel: Record<string, string> = {
  BO6: "Black Ops 6",
  MW3: "Modern Warfare III",
  MW2: "Modern Warfare II",
  VG: "Vanguard",
  CW: "Black Ops Cold War",
};

type RosterScope = "current" | "used";

const TeamDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [rosterScope, setRosterScope] = useState<RosterScope>("current");

  const { data: team, loading: teamLoading, error: teamError } = useApi<Team>(
    `/api/v1/teams/${id}`
  );

  const rosterUrl =
    rosterScope === "used"
      ? `/api/v1/teams/${id}/players?scope=all`
      : `/api/v1/teams/${id}/players`;

  const { data: players, loading: playersLoading } = useApi<Player[]>(rosterUrl);
  const franchiseKey = team?.franchise?.franchise_key ?? "";
  const { data: franchiseData } = useApi<FranchiseResponse>(
    `/api/v1/franchises/${franchiseKey}`,
    { enabled: !!franchiseKey }
  );

  if (teamLoading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Loading team data...</p>
      </div>
    );
  }

  if (teamError || !team) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Team not found</p>
        <Link to="/teams" className="text-[#a3a3a3] hover:text-white mt-4 inline-block text-sm transition-colors">
          ← Back to Teams
        </Link>
      </div>
    );
  }

  const logo = getTeamLogo(team.name);
  // All era names for this franchise slot, newest first.
  const eras = franchiseData?.eras ? [...franchiseData.eras].reverse() : [];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <Link
        to="/teams"
        className="text-[#737373] hover:text-[#a3a3a3] text-xs uppercase tracking-widest mb-10 inline-block transition-colors"
      >
        ← Teams
      </Link>

      {/* Team header */}
      <div className="flex flex-wrap items-center gap-6 mb-12 pb-8 border-b border-[#1a1a1a]">
        {logo ? (
          <img src={logo} alt={team.name} className="w-24 h-24 object-contain opacity-90" />
        ) : (
          <div className="w-24 h-24 bg-[#1a1a1a] flex items-center justify-center text-sm font-bold text-[#737373] font-mono">
            {team.abbreviation}
          </div>
        )}
        <div>
          <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">
            {team.abbreviation} · {team.game_code ? gameLabel[team.game_code] ?? team.game_code : "CDL"}
          </p>
          <h1 className="font-grotesk text-4xl font-bold text-white">{team.name}</h1>
          {team.is_active && (
            <span className="mt-2 inline-block text-[10px] uppercase tracking-widest text-green-400 bg-green-400/10 px-2 py-0.5">
              Active
            </span>
          )}
        </div>

        {/* Franchise era selector — switches which game/era of this franchise
            is shown by navigating to that era's team page. */}
        {eras.length > 1 && (
          <div className="w-full sm:w-auto sm:ml-auto">
            <label
              htmlFor="era-select"
              className="block text-[10px] uppercase tracking-widest text-[#737373] mb-1.5"
            >
              Era
            </label>
            <select
              id="era-select"
              value={team.id}
              onChange={(e) => navigate(`/teams/${e.target.value}`)}
              className="w-full sm:w-auto bg-[#111111] border border-[#1a1a1a] px-3 py-2 text-xs text-[#a3a3a3] focus:outline-none focus:border-[#2a2a2a] uppercase tracking-wider"
            >
              {eras.map((era) => (
                <option key={era.id} value={era.id}>
                  {era.game_code ? gameLabel[era.game_code] ?? era.game_code : era.name}
                </option>
              ))}
            </select>
          </div>
        )}
      </div>

      <div className="grid md:grid-cols-3 gap-8">
        {/* Left column: roster + franchise history */}
        <div className="md:col-span-2 space-y-8">
          {/* Roster */}
          <div>
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-4">
              <h2 className="text-xs uppercase tracking-widest text-[#737373]">
                {rosterScope === "used" ? "Players Used This Season" : "Roster"}
              </h2>
              <div className="inline-flex self-start bg-[#111111] border border-[#1a1a1a] p-0.5">
                {([
                  ["current", "Current Roster"],
                  ["used", "Players Used"],
                ] as const).map(([scope, label]) => {
                  const isActive = rosterScope === scope;
                  return (
                    <button
                      key={scope}
                      type="button"
                      onClick={() => setRosterScope(scope)}
                      className={`px-3 py-1.5 text-[10px] uppercase tracking-widest transition-colors ${
                        isActive
                          ? "bg-[#1f1f1f] text-white"
                          : "text-[#737373] hover:text-[#a3a3a3]"
                      }`}
                    >
                      {label}
                    </button>
                  );
                })}
              </div>
            </div>
            {playersLoading ? (
              <p className="text-[#737373] text-sm">Loading roster...</p>
            ) : players && players.length > 0 ? (
              <div className="space-y-1">
                {players.map((player) => (
                  <Link
                    key={player.id}
                    to={`/players/${player.id}`}
                    className="group flex items-center justify-between p-4 bg-[#111111] border border-[#1a1a1a] hover:border-[#2a2a2a] hover:bg-[#161616] transition-all"
                  >
                    <div className="flex items-center gap-4">
                      <img
                        src={getPlayerAvatar(player.gamertag)}
                        alt={player.gamertag}
                        className="w-10 h-10 object-cover rounded-full opacity-90 group-hover:opacity-100 transition-opacity"
                      />
                      <div>
                        <p className="font-grotesk font-semibold text-white text-sm">
                          {player.gamertag}
                        </p>
                        <p className="text-[#737373] text-xs mt-0.5">
                          {player.first_name && player.last_name
                            ? `${player.first_name} ${player.last_name}`
                            : "—"}
                        </p>
                      </div>
                    </div>
                    <span className="text-[#737373] text-xs uppercase tracking-wider">
                      {player.role || "—"}
                    </span>
                  </Link>
                ))}
              </div>
            ) : (
              <p className="text-[#737373] text-sm">
                {rosterScope === "used"
                  ? "No players recorded for this season"
                  : "No current roster available"}
              </p>
            )}
          </div>
        </div>

        {/* Right column: franchise identity timeline */}
        {eras.length > 1 && (
          <div>
            <h2 className="text-xs uppercase tracking-widest text-[#737373] mb-4">
              Franchise History
            </h2>
            <div className="space-y-1">
              {eras.map((era) => {
                const eraLogo = getTeamLogo(era.name);
                const isCurrent = era.id === team.id;
                return (
                  <Link
                    key={era.id}
                    to={`/teams/${era.id}`}
                    className={`flex items-center gap-3 p-3 border transition-all ${
                      isCurrent
                        ? "border-[#2a2a2a] bg-[#161616]"
                        : "border-[#1a1a1a] bg-[#111111] hover:border-[#2a2a2a] hover:bg-[#161616]"
                    }`}
                  >
                    {eraLogo ? (
                      <img src={eraLogo} alt={era.name} className="w-8 h-8 object-contain opacity-80 flex-shrink-0" />
                    ) : (
                      <div className="w-8 h-8 bg-[#1a1a1a] flex items-center justify-center text-[10px] font-bold text-[#737373] flex-shrink-0 font-mono">
                        {era.abbreviation?.slice(0, 3)}
                      </div>
                    )}
                    <div className="min-w-0">
                      <p className={`text-xs font-semibold truncate ${isCurrent ? "text-white" : "text-[#a3a3a3]"}`}>
                        {era.name}
                        {isCurrent && <span className="ml-1.5 text-[9px] text-green-400">●</span>}
                      </p>
                      <p className="text-[10px] text-[#737373] mt-0.5">
                        {era.game_code ? gameLabel[era.game_code] ?? era.game_code : "—"}
                      </p>
                    </div>
                  </Link>
                );
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default TeamDetail;
