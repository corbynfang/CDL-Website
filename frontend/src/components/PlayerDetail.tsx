import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import type { Player } from "../types";

const PlayerDetail = () => {
  const { id } = useParams<{ id: string }>();
  const [activeTab, setActiveTab] = useState("last5");

  const {
    data: player,
    loading: playerLoading,
    error: playerError,
  } = useApi<Player>(`/api/v1/players/${id}`);

  const { data: stats, loading: statsLoading } = useApi<any>(
    `/api/v1/players/${id}/kd`,
  );

  const { data: matchesData, loading: matchesLoading } = useApi<any>(
    `/api/v1/players/${id}/matches?limit=10`,
  );

  const loading = playerLoading || statsLoading || matchesLoading;

  // Helper function to calculate progress bar percentage (normalize KD to 0-100%)
  const getKDProgress = (kd: number) => {
    if (!kd || kd <= 0) return 0;
    // Cap at 2.0 KD = 100% for visualization
    return Math.min((kd / 2.0) * 100, 100);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#6B7280]">Loading player data...</p>
        </div>
      </div>
    );
  }

  if (playerError || !player) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#555555]">Player not found</p>
          <Link
            to="/players"
            className="text-[#6B7280] hover:text-black mt-4 inline-block transition-colors"
          >
            ← Back to Players
          </Link>
        </div>
      </div>
    );
  }

  const matches = matchesData?.matches || [];
  const last5Matches = matches.slice(0, 5);
  const tournamentStats = stats?.tournament_stats || [];

  // Format birthdate for display
  const formatBirthdate = (dateString?: string) => {
    if (!dateString) return "—";
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString("en-US", {
        year: "numeric",
        month: "long",
        day: "numeric",
      });
    } catch {
      return dateString;
    }
  };

  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        {/* Back Button */}
        <Link
          to="/players"
          className="text-[#6B7280] hover:text-black mb-8 inline-block transition-colors"
        >
          ← Back to Players
        </Link>

        {/* Top Section: Player Image (Left) and KD Stats Box (Right) */}
        <div className="grid md:grid-cols-2 gap-8 mb-8">
          {/* Left: Player Image */}
          <div className="flex flex-col items-center md:items-start">
            {player.avatar_url && (
              <img
                src={player.avatar_url}
                alt={player.gamertag}
                className="w-48 h-48 object-cover shadow-md shadow-[rgba(0,0,0,0.1)] mb-6"
              />
            )}

            {/* Player Info Card */}
            <div className="w-full bg-[#F4F4F5] p-6 shadow-md shadow-[rgba(0,0,0,0.1)]">
              <div className="space-y-4">
                <div>
                  <p className="text-xs uppercase tracking-wider text-[#6B7280] mb-1">
                    Name
                  </p>
                  <p className="text-lg font-bold text-black">
                    {player.first_name && player.last_name
                      ? `${player.first_name} ${player.last_name}`
                      : player.gamertag}
                  </p>
                </div>
                {player.country && (
                  <div>
                    <p className="text-xs uppercase tracking-wider text-[#6B7280] mb-1">
                      Country
                    </p>
                    <p className="text-base text-black">{player.country}</p>
                  </div>
                )}
                {player.birthdate && (
                  <div>
                    <p className="text-xs uppercase tracking-wider text-[#6B7280] mb-1">
                      Birthday
                    </p>
                    <p className="text-base text-black">
                      {formatBirthdate(player.birthdate)}
                    </p>
                  </div>
                )}
                <div>
                  <p className="text-xs uppercase tracking-wider text-[#6B7280] mb-1">
                    Team
                  </p>
                  <p className="text-base text-black">—</p>
                </div>
                {player.role && (
                  <div>
                    <p className="text-xs uppercase tracking-wider text-[#6B7280] mb-1">
                      Role
                    </p>
                    <p className="text-base text-black">{player.role}</p>
                  </div>
                )}
                <div>
                  <p className="text-xs uppercase tracking-wider text-[#6B7280] mb-1">
                    Position
                  </p>
                  <p className="text-base text-black">
                    {player.is_active ? "Active" : "Inactive"}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Right: KD Stats Box with Progress Bars */}
          <div className="bg-[#F4F4F5] p-8 shadow-md shadow-[rgba(0,0,0,0.1)]">
            <h2 className="text-2xl font-bold text-black mb-6">
              K/D STATISTICS
            </h2>

            <div className="space-y-6">
              {/* Overall KD */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <p className="text-xs text-[#6B7280] font-medium">
                    Overall K/D
                  </p>
                  <p className="text-lg font-bold text-black">
                    {stats?.avg_kd?.toFixed(2) || "0.00"}
                  </p>
                </div>
                <div className="w-full bg-gray-200 h-3 overflow-hidden">
                  <div
                    className="bg-black/80 h-full transition-all duration-500"
                    style={{
                      width: `${getKDProgress(stats?.avg_kd || 0)}%`,
                    }}
                  />
                </div>
              </div>

              {/* Hardpoint KD */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <p className="text-xs text-[#6B7280] font-medium">
                    Hardpoint K/D
                  </p>
                  <p className="text-lg font-bold text-black">
                    {stats?.hp_kd_ratio?.toFixed(2) || "0.00"}
                  </p>
                </div>
                <div className="w-full bg-gray-200 h-3 overflow-hidden">
                  <div
                    className="bg-black/80 h-full transition-all duration-500"
                    style={{
                      width: `${getKDProgress(stats?.hp_kd_ratio || 0)}%`,
                    }}
                  />
                </div>
              </div>

              {/* Search & Destroy KD */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <p className="text-xs text-[#6B7280] font-medium">
                    S&D K/D
                  </p>
                  <p className="text-lg font-bold text-black">
                    {stats?.snd_kd_ratio?.toFixed(2) || "0.00"}
                  </p>
                </div>
                <div className="w-full bg-gray-200 h-3 overflow-hidden">
                  <div
                    className="bg-black/80 h-full transition-all duration-500"
                    style={{
                      width: `${getKDProgress(stats?.snd_kd_ratio || 0)}%`,
                    }}
                  />
                </div>
              </div>

              {/* Control KD */}
              <div>
                <div className="flex justify-between items-center mb-2">
                  <p className="text-xs text-[#6B7280] font-medium">
                    Control K/D
                  </p>
                  <p className="text-lg font-bold text-black">
                    {stats?.control_kd_ratio?.toFixed(2) || "0.00"}
                  </p>
                </div>
                <div className="w-full bg-gray-200 h-3 overflow-hidden">
                  <div
                    className="bg-black/80 h-full transition-all duration-500"
                    style={{
                      width: `${getKDProgress(stats?.control_kd_ratio || 0)}%`,
                    }}
                  />
                </div>
              </div>
            </div>

            {/* Additional Stats Summary */}
            <div className="mt-8 pt-8 border-t border-gray-300 grid grid-cols-3 gap-4">
              <div className="text-center">
                <p className="text-xs text-[#6B7280] uppercase mb-1">Kills</p>
                <p className="text-xl font-bold text-black">
                  {stats?.total_kills?.toLocaleString() || "0"}
                </p>
              </div>
              <div className="text-center">
                <p className="text-xs text-[#6B7280] uppercase mb-1">Deaths</p>
                <p className="text-xl font-bold text-black">
                  {stats?.total_deaths?.toLocaleString() || "0"}
                </p>
              </div>
              <div className="text-center">
                <p className="text-xs text-[#6B7280] uppercase mb-1">Assists</p>
                <p className="text-xl font-bold text-black">
                  {stats?.total_assists?.toLocaleString() || "0"}
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Bottom Section: Tabs and Content */}
        <div className="mt-12">
          {/* Tabs */}
          <div className="flex space-x-1 border-b border-gray-300 mb-8">
            <button
              onClick={() => setActiveTab("last5")}
              className={`px-6 py-3 text-sm font-medium ${
                activeTab === "last5"
                  ? "text-black border-b-2 border-black"
                  : "text-[#6B7280]"
              }`}
            >
              Last 5 Matches
            </button>
            <button
              onClick={() => setActiveTab("matches")}
              className={`px-6 py-3 text-sm font-medium transition-all ${
                activeTab === "matches"
                  ? "text-black border-b-2 border-black"
                  : "text-[#6B7280]"
              }`}
            >
              Matches
            </button>
            <button
              onClick={() => setActiveTab("eventStats")}
              className={`px-6 py-3 text-sm font-medium transition-all ${
                activeTab === "eventStats"
                  ? "text-black border-b-2 border-black"
                  : "text-[#6B7280]"
              }`}
            >
              Event Stats
            </button>
            <button
              onClick={() => setActiveTab("events")}
              className={`px-6 py-3 text-sm font-medium transition-all ${
                activeTab === "events"
                  ? "text-black border-b-2 border-black"
                  : "text-[#6B7280]"
              }`}
            >
              Events
            </button>
          </div>

          {/* Tab Content */}
          <div className="bg-[#F4F4F5] p-8 shadow-md shadow-[rgba(0,0,0,0.1)]">
            {activeTab === "last5" && (
              <div>
                <h3 className="text-xl font-bold text-black mb-6">
                  Last 5 Matches
                </h3>
                {last5Matches.length > 0 ? (
                  <div className="space-y-4">
                    {last5Matches.map((match: any, index: number) => (
                      <div
                        key={index}
                        className="bg-white p-6"
                      >
                        <div className="flex justify-between items-center">
                          <div className="flex items-center space-x-4">
                            <span
                              className={`text-sm font-bold w-8 text-center ${
                                match.result === "W"
                                  ? "text-black"
                                  : "text-[#6B7280]"
                              }`}
                            >
                              {match.result || "—"}
                            </span>
                            <div>
                              <p className="font-semibold text-black">
                                vs {match.opponent || "Unknown"}
                              </p>
                              <p className="text-sm text-[#6B7280]">
                                {match.tournament || "Tournament"} •{" "}
                                {match.date
                                  ? new Date(match.date).toLocaleDateString()
                                  : "—"}
                              </p>
                            </div>
                          </div>
                          <div className="flex space-x-6 text-right">
                            <div>
                              <p className="text-xs text-[#6B7280] uppercase">
                                K/D
                              </p>
                              <p className="text-lg font-bold text-black">
                                {match.kd?.toFixed(2) || "0.00"}
                              </p>
                            </div>
                            <div>
                              <p className="text-xs text-[#6B7280] uppercase">
                                Kills
                              </p>
                              <p className="text-lg font-bold text-black">
                                {match.kills || "0"}
                              </p>
                            </div>
                            <div>
                              <p className="text-xs text-[#6B7280] uppercase">
                                Deaths
                              </p>
                              <p className="text-lg font-bold text-black">
                                {match.deaths || "0"}
                              </p>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-[#6B7280]">No matches available</p>
                )}
              </div>
            )}

            {activeTab === "matches" && (
              <div>
                <h3 className="text-xl font-bold text-black mb-6">
                  All Matches
                </h3>
                {matches.length > 0 ? (
                  <div className="space-y-4">
                    {matches.map((match: any, index: number) => (
                      <div
                        key={index}
                        className="bg-white p-6"
                      >
                        <div className="flex justify-between items-center">
                          <div className="flex items-center space-x-4">
                            <span
                              className={`text-sm font-bold w-8 text-center ${
                                match.result === "W"
                                  ? "text-black"
                                  : "text-[#6B7280]"
                              }`}
                            >
                              {match.result || "—"}
                            </span>
                            <div>
                              <p className="font-semibold text-black">
                                vs {match.opponent || "Unknown"}
                              </p>
                              <p className="text-sm text-[#6B7280]">
                                {match.tournament || "Tournament"} •{" "}
                                {match.date
                                  ? new Date(match.date).toLocaleDateString()
                                  : "—"}
                              </p>
                            </div>
                          </div>
                          <div className="flex space-x-6 text-right">
                            <div>
                              <p className="text-xs text-[#6B7280] uppercase">
                                K/D
                              </p>
                              <p className="text-lg font-bold text-black">
                                {match.kd?.toFixed(2) || "0.00"}
                              </p>
                            </div>
                            <div>
                              <p className="text-xs text-[#6B7280] uppercase">
                                Kills
                              </p>
                              <p className="text-lg font-bold text-black">
                                {match.kills || "0"}
                              </p>
                            </div>
                            <div>
                              <p className="text-xs text-[#6B7280] uppercase">
                                Deaths
                              </p>
                              <p className="text-lg font-bold text-black">
                                {match.deaths || "0"}
                              </p>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-[#6B7280]">No matches available</p>
                )}
              </div>
            )}

            {activeTab === "eventStats" && (
              <div>
                <h3 className="text-xl font-bold text-black mb-6">
                  Event Statistics
                </h3>
                {tournamentStats.length > 0 ? (
                  <div className="space-y-4">
                    {tournamentStats.map((tournament: any, index: number) => (
                      <div
                        key={index}
                        className="bg-white p-6"
                      >
                        <div className="flex justify-between items-start mb-4">
                          <div>
                            <p className="font-bold text-black text-lg">
                              {tournament.tournament_name || "Tournament"}
                            </p>
                            <p className="text-sm text-[#6B7280]">
                              {tournament.matches || 0} matches •{" "}
                              {tournament.maps_played || 0} maps
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-xs text-[#6B7280] uppercase mb-1">
                              K/D
                            </p>
                            <p className="text-2xl font-bold text-black">
                              {tournament.kd_ratio?.toFixed(2) || "0.00"}
                            </p>
                          </div>
                        </div>
                        <div className="grid grid-cols-3 gap-4 text-center pt-4 border-t border-gray-200">
                          <div>
                            <p className="text-xs text-[#6B7280] uppercase mb-1">
                              Kills
                            </p>
                            <p className="text-lg font-bold text-black">
                              {tournament.kills?.toLocaleString() || "0"}
                            </p>
                          </div>
                          <div>
                            <p className="text-xs text-[#6B7280] uppercase mb-1">
                              Deaths
                            </p>
                            <p className="text-lg font-bold text-black">
                              {tournament.deaths?.toLocaleString() || "0"}
                            </p>
                          </div>
                          <div>
                            <p className="text-xs text-[#6B7280] uppercase mb-1">
                              Assists
                            </p>
                            <p className="text-lg font-bold text-black">
                              {tournament.assists?.toLocaleString() || "0"}
                            </p>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-[#6B7280]">
                    No event statistics available
                  </p>
                )}
              </div>
            )}

            {activeTab === "events" && (
              <div>
                <h3 className="text-xl font-bold text-black mb-6">Events</h3>
                {tournamentStats.length > 0 ? (
                  <div className="space-y-4">
                    {tournamentStats.map((tournament: any, index: number) => (
                      <div
                        key={index}
                        className="bg-white p-6"
                      >
                        <div className="flex justify-between items-center">
                          <div>
                            <p className="font-bold text-black text-lg">
                              {tournament.tournament_name || "Tournament"}
                            </p>
                            <p className="text-sm text-[#6B7280]">
                              {tournament.matches || 0} matches played
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-xs text-[#6B7280] uppercase mb-1">
                              Maps
                            </p>
                            <p className="text-xl font-bold text-black">
                              {tournament.maps_played || 0}
                            </p>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-[#6B7280]">No events available</p>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default PlayerDetail;
