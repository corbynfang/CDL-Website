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
    `/api/v1/players/${id}/matches`,
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

  const events = matchesData?.events || [];
  // Get all matches from events for "Last 5" tab
  const allMatches = events.flatMap((event: any) => event.matches || []);
  const last5Matches = allMatches.slice(0, 5);
  const tournamentStats = stats?.tournament_stats || [];
  
  // Get overall KD stats for display
  const overallKD = stats?.avg_kd || 0;
  const hpKD = stats?.hp_kd_ratio || stats?.ewc_hp_kd_ratio || 0;
  const sndKD = stats?.snd_kd_ratio || stats?.ewc_snd_kd_ratio || 0;
  const ctlKD = stats?.control_kd_ratio || stats?.ewc_control_kd_ratio || 0;

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
                    {overallKD.toFixed(2)}
                  </p>
                </div>
                <div className="w-full bg-gray-200 h-3 overflow-hidden">
                  <div
                    className="bg-black/80 h-full transition-all duration-500"
                    style={{
                      width: `${getKDProgress(overallKD)}%`,
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
                    {hpKD.toFixed(2)}
                  </p>
                </div>
                <div className="w-full bg-gray-200 h-3 overflow-hidden">
                  <div
                    className="bg-black/80 h-full transition-all duration-500"
                    style={{
                      width: `${getKDProgress(hpKD)}%`,
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
                    {sndKD.toFixed(2)}
                  </p>
                </div>
                <div className="w-full bg-gray-200 h-3 overflow-hidden">
                  <div
                    className="bg-black/80 h-full transition-all duration-500"
                    style={{
                      width: `${getKDProgress(sndKD)}%`,
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
                    {ctlKD.toFixed(2)}
                  </p>
                </div>
                <div className="w-full bg-gray-200 h-3 overflow-hidden">
                  <div
                    className="bg-black/80 h-full transition-all duration-500"
                    style={{
                      width: `${getKDProgress(ctlKD)}%`,
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
                                vs {match.opponent_abbr || match.opponent || "Unknown"}
                              </p>
                              <p className="text-sm text-[#6B7280]">
                                {match.date
                                  ? new Date(match.date).toLocaleDateString()
                                  : "—"}
                              </p>
                            </div>
                          </div>
                          <div className="flex space-x-6 text-right">
                            <div>
                              <p className="text-xs text-[#6B7280] uppercase">
                                Result
                              </p>
                              <p className="text-lg font-bold text-black">
                                {match.result || "—"}
                              </p>
                            </div>
                            <div>
                              <p className="text-xs text-[#6B7280] uppercase">
                                K/D
                              </p>
                              <p className="text-lg font-bold text-black">
                                {typeof match.kd === 'number' ? match.kd.toFixed(2) : "0.00"}
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
                {events.length > 0 ? (
                  <div className="space-y-8">
                    {events.map((event: any, eventIndex: number) => (
                      <div key={eventIndex}>
                        <h3 className="text-xl font-bold text-black mb-4">
                          {event.event} {event.year}
                        </h3>
                        {event.matches && event.matches.length > 0 ? (
                          <div className="overflow-x-auto">
                            <table className="w-full">
                              <thead>
                                <tr className="border-b border-gray-300">
                                  <th className="text-left py-3 text-[#6B7280] text-xs uppercase tracking-wider">Date</th>
                                  <th className="text-left py-3 text-[#6B7280] text-xs uppercase tracking-wider">Opponent</th>
                                  <th className="text-left py-3 text-[#6B7280] text-xs uppercase tracking-wider">Result</th>
                                  <th className="text-right py-3 text-[#6B7280] text-xs uppercase tracking-wider">KD</th>
                                  <th className="text-right py-3 text-[#6B7280] text-xs uppercase tracking-wider">K</th>
                                  <th className="text-right py-3 text-[#6B7280] text-xs uppercase tracking-wider">D</th>
                                  <th className="text-right py-3 text-[#6B7280] text-xs uppercase tracking-wider">HP KD</th>
                                  <th className="text-right py-3 text-[#6B7280] text-xs uppercase tracking-wider">SND KD</th>
                                  <th className="text-right py-3 text-[#6B7280] text-xs uppercase tracking-wider">CTL KD</th>
                                  <th className="text-right py-3 text-[#6B7280] text-xs uppercase tracking-wider">Slayer Rating</th>
                                  <th className="text-right py-3 text-[#6B7280] text-xs uppercase tracking-wider">Rating</th>
                                </tr>
                              </thead>
                              <tbody>
                                {event.matches.map((match: any, matchIndex: number) => (
                                  <tr key={matchIndex} className="border-b border-gray-300">
                                    <td className="py-3 text-[#6B7280] text-sm">
                                      {match.date ? new Date(match.date).toLocaleDateString() : "—"}
                                    </td>
                                    <td className="py-3 text-black text-sm font-medium">
                                      {match.opponent_abbr || match.opponent || "—"}
                                    </td>
                                    <td className="py-3 text-black text-sm font-medium">
                                      {match.result || "—"}
                                    </td>
                                    <td className="py-3 text-right text-black text-sm font-bold">
                                      {typeof match.kd === 'number' ? match.kd.toFixed(2) : "0.00"}
                                    </td>
                                    <td className="py-3 text-right text-[#6B7280] text-sm">
                                      {match.kills || "0"}
                                    </td>
                                    <td className="py-3 text-right text-[#6B7280] text-sm">
                                      {match.deaths || "0"}
                                    </td>
                                    <td className="py-3 text-right text-[#6B7280] text-sm">
                                      {typeof match.hp_kd === 'number' ? match.hp_kd.toFixed(2) : "—"}
                                    </td>
                                    <td className="py-3 text-right text-[#6B7280] text-sm">
                                      {typeof match.snd_kd === 'number' ? match.snd_kd.toFixed(2) : "—"}
                                    </td>
                                    <td className="py-3 text-right text-[#6B7280] text-sm">
                                      {typeof match.ctl_kd === 'number' ? match.ctl_kd.toFixed(2) : "—"}
                                    </td>
                                    <td className="py-3 text-right text-[#6B7280] text-sm">
                                      {typeof match.slayer_rating === 'number' ? match.slayer_rating.toFixed(2) : "—"}
                                    </td>
                                    <td className="py-3 text-right text-[#6B7280] text-sm">
                                      {typeof match.rating === 'number' ? match.rating.toFixed(1) : "—"}
                                    </td>
                                  </tr>
                                ))}
                              </tbody>
                            </table>
                          </div>
                        ) : (
                          <p className="text-[#6B7280]">No matches available for this event</p>
                        )}
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-[#6B7280]">No match data available for this season</p>
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
