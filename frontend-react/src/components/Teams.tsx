import { useState } from "react";
import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import { getTeamLogo } from "../utils/assets";
import type { Season } from "../types";
import PageMeta from "./PageMeta";

interface Team {
  id: number;
  name: string;
  abbreviation: string;
  logo_url: string;
}

const Teams = () => {
  const [selectedSeasonId, setSelectedSeasonId] = useState<string>("");

  const { data: seasons } = useApi<Season[]>("/api/v1/seasons");

  const teamsUrl = selectedSeasonId
    ? `/api/v1/teams?season_id=${selectedSeasonId}`
    : "/api/v1/teams";

  const {
    data: allTeams,
    loading,
    error,
  } = useApi<Team[]>(teamsUrl, { cacheTtl: 5 * 60 * 1000 });

  const teams = allTeams?.filter((t) => t.name !== "Unaffiliated");

  const selectedSeason = seasons?.find(
    (s) => String(s.id) === selectedSeasonId,
  );

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <PageMeta
        title="CDL Teams"
        description="Every Call of Duty League franchise and team — rosters, season history, and performance stats. Filter by CDL season."
        canonical="/teams"
      />
      <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4 mb-8">
        <div>
          <p className="text-xs uppercase tracking-widest text-[#737373] mb-2">
            League
          </p>
          <h1 className="font-grotesk text-3xl font-bold text-white">TEAMS</h1>
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

      {loading && <p className="text-[#737373] text-sm">Loading teams...</p>}
      {error && <p className="text-[#737373] text-sm">Error: {error}</p>}

      {!loading && !error && (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
            {teams?.map((team) => {
              const logo = getTeamLogo(team.name);
              return (
                <Link
                  key={team.id}
                  to={`/teams/${team.id}`}
                  className="group flex items-center gap-4 p-5 bg-[#111111] border border-[#1a1a1a] hover:border-[#2a2a2a] hover:bg-[#161616] transition-all"
                >
                  {logo ? (
                    <img
                      src={logo}
                      alt={team.name}
                      className="w-14 h-14 object-contain flex-shrink-0 opacity-90 group-hover:opacity-100 transition-opacity"
                    />
                  ) : (
                    <div className="w-14 h-14 bg-[#1a1a1a] flex items-center justify-center text-xs font-bold text-[#737373] flex-shrink-0 font-mono">
                      {team.abbreviation}
                    </div>
                  )}
                  <div>
                    <p className="font-grotesk font-semibold text-white text-sm">
                      {team.name}
                    </p>
                    <p className="text-[#737373] text-xs tracking-wider mt-0.5">
                      {team.abbreviation}
                    </p>
                  </div>
                </Link>
              );
            })}
          </div>

          {teams?.length === 0 && (
            <p className="text-center text-[#737373] py-16 text-sm">
              No teams found for this season.
            </p>
          )}

          <p className="mt-4 text-[#737373] text-xs">
            {teams?.length ?? 0} teams
          </p>
        </>
      )}
    </div>
  );
};

export default Teams;
