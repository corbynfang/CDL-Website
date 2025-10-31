import { useParams, Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import type { Player } from "../types";

interface Team {
  id: number;
  name: string;
  abbreviation: string;
  logo_url: string;
}

const TeamDetail = () => {
  const { id } = useParams<{ id: string }>();

  const {
    data: team,
    loading: teamLoading,
    error: teamError,
  } = useApi<Team>(`/api/v1/teams/${id}`);

  const { data: players, loading: playersLoading } = useApi<Player[]>(
    `/api/v1/teams/${id}/players`,
  );

  const loading = teamLoading || playersLoading;

  if (loading) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#6B7280]">Loading team data...</p>
        </div>
      </div>
    );
  }

  if (teamError || !team) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#555555]">Team not found</p>
          <Link
            to="/teams"
            className="text-[#6B7280] mt-4 inline-block"
          >
            ← Back to Teams
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        {/* Back Button */}
        <Link
          to="/teams"
          className="text-[#6B7280] mb-8 inline-block"
        >
          ← Back to Teams
        </Link>

        {/* Team Header */}
        <div className="pb-8 mb-8">
          <div className="flex items-center space-x-6">
            {team.logo_url && (
              <img
                src={team.logo_url}
                alt={team.name}
                className="w-32 h-32 object-contain"
              />
            )}
            <div>
              <h1 className="text-5xl font-bold mb-2 text-black">{team.name}</h1>
              <p className="text-xl text-[#6B7280]">{team.abbreviation}</p>
            </div>
          </div>
        </div>

        {/* Roster */}
        <div>
          <h2 className="text-2xl font-bold mb-6 text-black">ROSTER</h2>
          {players && players.length > 0 ? (
            <div className="space-y-4">
              {players.map((player) => (
                <Link
                  key={player.id}
                  to={`/players/${player.id}`}
                  className="p-6 bg-[#F4F4F5] shadow-md shadow-[rgba(0,0,0,0.1)] flex items-center justify-between"
                >
                  <div className="flex items-center space-x-4">
                    {player.avatar_url && (
                      <img
                        src={player.avatar_url}
                        alt={player.gamertag}
                        className="w-12 h-12"
                      />
                    )}
                    <div>
                      <p className="font-bold text-lg text-black">{player.gamertag}</p>
                      <p className="text-sm text-[#6B7280]">
                        {player.first_name && player.last_name
                          ? `${player.first_name} ${player.last_name}`
                          : "—"}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-[#6B7280] text-sm">{player.role || "—"}</p>
                  </div>
                </Link>
              ))}
            </div>
          ) : (
            <p className="text-[#6B7280]">No players found</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default TeamDetail;
