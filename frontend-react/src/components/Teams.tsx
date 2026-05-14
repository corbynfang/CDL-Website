import { useState } from "react";
import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import { getTeamLogo } from "../utils/assets";
import type { Season } from "../types";

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

  const { data: allTeams, loading, error } = useApi<Team[]>(teamsUrl);

  const teams = allTeams?.filter((t) => t.name !== "Unaffiliated");

  const selectedSeason = seasons?.find(
    (s) => String(s.id) === selectedSeasonId
  );

  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
          <div>
            <h1 className="text-4xl font-bold text-black">TEAMS</h1>
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

        {loading && (
          <p className="text-[#6B7280]">Loading teams...</p>
        )}

        {error && (
          <p className="text-[#555555]">Error: {error}</p>
        )}

        {!loading && !error && (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {teams?.map((team) => {
                const logo = getTeamLogo(team.name);
                return (
                  <Link
                    key={team.id}
                    to={`/teams/${team.id}`}
                    className="p-8 bg-[#F4F4F5] shadow-md shadow-[rgba(0,0,0,0.1)] hover:shadow-lg transition-shadow"
                  >
                    <div className="flex items-center space-x-4">
                      {logo ? (
                        <img
                          src={logo}
                          alt={team.name}
                          className="w-16 h-16 object-contain flex-shrink-0"
                        />
                      ) : (
                        <div className="w-16 h-16 bg-gray-200 flex items-center justify-center text-xs font-bold text-gray-500 flex-shrink-0">
                          {team.abbreviation}
                        </div>
                      )}
                      <div>
                        <h2 className="text-xl font-bold text-black">
                          {team.name}
                        </h2>
                        <p className="text-[#6B7280]">{team.abbreviation}</p>
                      </div>
                    </div>
                  </Link>
                );
              })}
            </div>

            {teams?.length === 0 && (
              <p className="text-center text-[#6B7280] py-12">
                No teams found for this season.
              </p>
            )}

            <p className="mt-8 text-[#6B7280] text-sm">
              {teams?.length ?? 0} teams
            </p>
          </>
        )}
      </div>
    </div>
  );
};

export default Teams;
