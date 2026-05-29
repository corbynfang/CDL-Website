import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import { getPlayerAvatar } from "../utils/avatarAssets";
import { getKdColorClass } from "../utils/kdUtils";
import type {
  Player,
  PlayerKDResponse,
  PlayerMatchHistory,
  PlayerCareerResponse,
  MatchHistoryEvent,
  MatchHistoryResult,
  PlayerKDTournamentEntry,
  PlayerFranchiseEntry,
  PlayerEraStats,
} from "../types";

const TABS = [
  { id: "last5", label: "Last 5" },
  { id: "matches", label: "Matches" },
  { id: "eventStats", label: "Event Stats" },
  { id: "events", label: "Events" },
  { id: "career", label: "Career" },
];

const kdBarWidth = (kd: number) => `${Math.min((kd / 2.0) * 100, 100)}%`;

const formatBirthdate = (dateString?: string) => {
  if (!dateString) return "—";
  try {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  } catch {
    return dateString;
  }
};

const PlayerDetail = () => {
  const { id } = useParams<{ id: string }>();
  const [activeTab, setActiveTab] = useState("last5");

  const { data: player, loading: playerLoading, error: playerError } = useApi<Player>(
    `/api/v1/players/${id}`
  );
  const { data: stats, loading: statsLoading } = useApi<PlayerKDResponse>(
    `/api/v1/players/${id}/kd`
  );
  const { data: matchesData, loading: matchesLoading } = useApi<PlayerMatchHistory>(
    `/api/v1/players/${id}/matches`
  );
  const { data: careerData } = useApi<PlayerCareerResponse>(
    `/api/v1/players/${id}/franchise-career`
  );

  const loading = playerLoading || statsLoading || matchesLoading;

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Loading player data...</p>
      </div>
    );
  }

  if (playerError || !player) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Player not found</p>
        <Link to="/players" className="text-[#a3a3a3] hover:text-white mt-4 inline-block text-sm transition-colors">
          ← Back to Players
        </Link>
      </div>
    );
  }

  const events: MatchHistoryEvent[] = matchesData?.events || [];
  const allMatches: MatchHistoryResult[] = events
    .flatMap((e) => e.matches || [])
    .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
  const last5Matches = allMatches.slice(0, 5);
  const tournamentStats: PlayerKDTournamentEntry[] = stats?.tournament_stats || [];

  const overallKD = stats?.avg_kd || 0;
  const hpKD = stats?.hp_kd_ratio || 0;
  const sndKD = stats?.snd_kd_ratio || 0;
  const ctlKD = stats?.control_kd_ratio || 0;

  // Some eras (e.g. BO6) have no per-mode source, so HP/SND/Control come back 0.
  // Hide those rows entirely rather than showing three empty 0.00 bars.
  const hasModeSplits = hpKD > 0 || sndKD > 0 || ctlKD > 0;
  const kdBreakdown = [
    { label: "Overall", value: overallKD },
    ...(hasModeSplits
      ? [
          { label: "Hardpoint", value: hpKD },
          { label: "Search & Destroy", value: sndKD },
          { label: "Control", value: ctlKD },
        ]
      : []),
  ];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <Link
        to="/players"
        className="text-[#737373] hover:text-[#a3a3a3] text-xs uppercase tracking-widest mb-10 inline-block transition-colors"
      >
        ← Players
      </Link>

      {/* Top Section */}
      <div className="grid md:grid-cols-2 gap-4 mb-4">
        {/* Left: Identity card */}
        <div className="flex flex-col gap-4">
          <img
            src={getPlayerAvatar(player.gamertag)}
            alt={player.gamertag}
            className="w-36 h-36 object-cover"
          />

          <div className="bg-[#111111] border border-[#1a1a1a] p-6 space-y-4 flex-1">
            <div>
              <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">Gamertag</p>
              <p className="font-grotesk text-2xl font-bold text-white">{player.gamertag}</p>
            </div>

            {(player.first_name || player.last_name) && (
              <div>
                <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">Name</p>
                <p className="text-[#a3a3a3] text-sm">
                  {player.first_name} {player.last_name}
                </p>
              </div>
            )}

            {player.country && (
              <div>
                <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">Country</p>
                <p className="text-[#a3a3a3] text-sm">{player.country}</p>
              </div>
            )}

            {player.birthdate && (
              <div>
                <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">Birthday</p>
                <p className="text-[#a3a3a3] text-sm">{formatBirthdate(player.birthdate)}</p>
              </div>
            )}

            {player.role && (
              <div>
                <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">Role</p>
                <p className="text-[#a3a3a3] text-sm">{player.role}</p>
              </div>
            )}

            <div>
              <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">Status</p>
              <span
                className={`text-xs font-medium px-2 py-0.5 ${
                  player.is_active
                    ? "text-green-400 bg-green-400/10"
                    : "text-[#737373] bg-white/5"
                }`}
              >
                {player.is_active ? "Active" : "Inactive"}
              </span>
            </div>
          </div>
        </div>

        {/* Right: KD stats */}
        <div className="bg-[#111111] border border-[#1a1a1a] p-6">
          <h2 className="text-xs uppercase tracking-widest text-[#737373] mb-6">K/D Statistics</h2>

          <div className="space-y-5">
            {kdBreakdown.map(({ label, value }) => (
              <div key={label}>
                <div className="flex justify-between items-baseline mb-2">
                  <p className="text-xs text-[#737373]">{label}</p>
                  <p className={`font-mono font-bold text-sm ${getKdColorClass(value)}`}>
                    {value.toFixed(2)}
                  </p>
                </div>
                <div className="w-full bg-[#1a1a1a] h-1.5 overflow-hidden">
                  <div
                    className="bg-white/60 h-full transition-all duration-700"
                    style={{ width: kdBarWidth(value) }}
                  />
                </div>
              </div>
            ))}
          </div>

          <div className="mt-8 pt-6 border-t border-[#1a1a1a] grid grid-cols-3 gap-4">
            {[
              { label: "Kills", value: stats?.total_kills },
              { label: "Deaths", value: stats?.total_deaths },
              { label: "Assists", value: stats?.total_assists },
            ].map(({ label, value }) => (
              <div key={label} className="text-center">
                <p className="text-xs text-[#737373] uppercase tracking-wider mb-1">{label}</p>
                <p className="font-mono font-bold text-white text-lg">
                  {value?.toLocaleString() || "0"}
                </p>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="mt-10">
        <div className="flex border-b border-[#1a1a1a]">
          {TABS.map(({ id: tabId, label }) => (
            <button
              key={tabId}
              onClick={() => setActiveTab(tabId)}
              className={`px-5 py-3 text-xs uppercase tracking-widest transition-colors border-b-2 -mb-px ${
                activeTab === tabId
                  ? "text-white border-white"
                  : "text-[#737373] border-transparent hover:text-[#a3a3a3]"
              }`}
            >
              {label}
            </button>
          ))}
        </div>

        <div className="bg-[#111111] border border-[#1a1a1a] border-t-0 p-6">
          {activeTab === "last5" && (
            <div>
              <h3 className="text-xs uppercase tracking-widest text-[#737373] mb-5">
                Last 5 Matches
              </h3>
              {last5Matches.length > 0 ? (
                <div className="space-y-2">
                  {last5Matches.map((match, i: number) => (
                    <Link
                      key={i}
                      to={match.match_id ? `/matches/${match.match_id}` : "#"}
                      className="block bg-[#0a0a0a] border border-[#1a1a1a] p-4 hover:border-[#2a2a2a] hover:bg-[#0f0f0f] transition-all"
                    >
                      <div className="flex justify-between items-center">
                        <div className="flex items-center gap-4">
                          <span
                            className={`font-mono font-bold text-sm w-6 text-center ${
                              match.result?.startsWith("W") ? "text-green-400" : "text-[#737373]"
                            }`}
                          >
                            {match.result?.charAt(0) || "—"}
                          </span>
                          <div>
                            <p className="text-white text-sm font-medium">
                              vs {match.opponent_abbr || match.opponent || "Unknown"}
                            </p>
                            <p className="text-[#737373] text-xs mt-0.5">
                              {match.date ? new Date(match.date).toLocaleDateString() : "—"}
                            </p>
                          </div>
                        </div>
                        <div className="flex gap-5 text-right items-center">
                          {[
                            { label: "K/D", value: typeof match.kd === "number" ? match.kd.toFixed(2) : "0.00", color: getKdColorClass(match.kd) },
                            { label: "K", value: match.kills || "0", color: "text-white" },
                            { label: "D", value: match.deaths || "0", color: "text-white" },
                          ].map(({ label, value, color }) => (
                            <div key={label}>
                              <p className="text-xs text-[#737373] uppercase mb-0.5">{label}</p>
                              <p className={`font-bold text-sm font-mono ${color}`}>{value}</p>
                            </div>
                          ))}
                          <span className="text-[#404040] text-xs ml-2">→</span>
                        </div>
                      </div>
                    </Link>
                  ))}
                </div>
              ) : (
                <p className="text-[#737373] text-sm">No matches available</p>
              )}
            </div>
          )}

          {activeTab === "matches" && (
            <div>
              {events.length > 0 ? (
                <div className="space-y-8">
                  {events.map((event, ei: number) => (
                    <div key={ei}>
                      <h3 className="text-xs uppercase tracking-widest text-[#737373] mb-4">
                        {event.event} {event.year}
                      </h3>
                      {event.matches?.length > 0 ? (
                        <div className="overflow-x-auto border border-[#1a1a1a]">
                          <table className="w-full">
                            <thead>
                              <tr className="border-b border-[#1a1a1a] bg-[#0a0a0a]">
                                {["Date", "Opp", "Result", "KD", "K", "D", "HP KD", "SND KD", "CTL KD", "Slayer", "Rating"].map((h) => (
                                  <th
                                    key={h}
                                    className="px-3 py-2 text-[#737373] text-xs uppercase tracking-widest font-medium text-left last:text-right"
                                  >
                                    {h}
                                  </th>
                                ))}
                              </tr>
                            </thead>
                            <tbody>
                              {event.matches.map((match, mi: number) => (
                                <tr
                                  key={mi}
                                  className="border-b border-[#1a1a1a] hover:bg-[#0a0a0a] transition-colors cursor-pointer"
                                  onClick={() => match.match_id && (window.location.href = `/matches/${match.match_id}`)}
                                >
                                  <td className="px-3 py-2 text-[#737373] text-xs font-mono">
                                    {match.date ? new Date(match.date).toLocaleDateString() : "—"}
                                  </td>
                                  <td className="px-3 py-2 text-[#a3a3a3] text-xs font-medium">
                                    {match.opponent_abbr || match.opponent || "—"}
                                  </td>
                                  <td className={`px-3 py-2 text-xs font-bold font-mono ${match.result?.startsWith("W") ? "text-green-400" : "text-[#737373]"}`}>
                                    {match.result || "—"}
                                  </td>
                                  <td className={`px-3 py-2 text-xs font-bold font-mono ${getKdColorClass(match.kd)}`}>
                                    {typeof match.kd === "number" ? match.kd.toFixed(2) : "0.00"}
                                  </td>
                                  <td className="px-3 py-2 text-[#737373] text-xs font-mono">{match.kills || "0"}</td>
                                  <td className="px-3 py-2 text-[#737373] text-xs font-mono">{match.deaths || "0"}</td>
                                  <td className="px-3 py-2 text-[#737373] text-xs font-mono">—</td>
                                  <td className="px-3 py-2 text-[#737373] text-xs font-mono">—</td>
                                  <td className="px-3 py-2 text-[#737373] text-xs font-mono">—</td>
                                  <td className="px-3 py-2 text-[#737373] text-xs font-mono">—</td>
                                  <td className="px-3 py-2 text-[#737373] text-xs font-mono text-right">
                                    {match.match_id ? (
                                      <Link to={`/matches/${match.match_id}`} className="text-[#404040] hover:text-[#737373] transition-colors" onClick={(e) => e.stopPropagation()}>→</Link>
                                    ) : "—"}
                                  </td>
                                </tr>
                              ))}
                            </tbody>
                          </table>
                        </div>
                      ) : (
                        <p className="text-[#737373] text-sm">No matches for this event</p>
                      )}
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-[#737373] text-sm">No match data available</p>
              )}
            </div>
          )}

          {activeTab === "eventStats" && (
            <div>
              <h3 className="text-xs uppercase tracking-widest text-[#737373] mb-5">
                Event Statistics
              </h3>
              {tournamentStats.length > 0 ? (
                <div className="space-y-2">
                  {tournamentStats.map((t, i: number) => (
                    <div key={i} className="bg-[#0a0a0a] border border-[#1a1a1a] p-5">
                      <div className="flex justify-between items-start mb-4">
                        <div>
                          <p className="font-grotesk font-semibold text-white text-sm">
                            {t.tournament_name || "Tournament"}
                          </p>
                          <p className="text-[#737373] text-xs mt-0.5">
                            {t.maps_played || 0} maps
                          </p>
                        </div>
                        <div className="text-right">
                          <p className="text-xs text-[#737373] uppercase tracking-wider mb-0.5">K/D</p>
                          <p className={`font-mono font-bold text-xl ${getKdColorClass(t.kd_ratio)}`}>
                            {t.kd_ratio?.toFixed(2) || "0.00"}
                          </p>
                        </div>
                      </div>
                      <div className="grid grid-cols-3 gap-4 pt-4 border-t border-[#1a1a1a] text-center">
                        {[
                          { label: "Kills", value: t.kills },
                          { label: "Deaths", value: t.deaths },
                          { label: "Assists", value: t.assists },
                        ].map(({ label, value }) => (
                          <div key={label}>
                            <p className="text-xs text-[#737373] uppercase tracking-wider mb-1">{label}</p>
                            <p className="font-mono font-bold text-white">
                              {value?.toLocaleString() || "0"}
                            </p>
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-[#737373] text-sm">No event statistics available</p>
              )}
            </div>
          )}

          {activeTab === "events" && (
            <div>
              <h3 className="text-xs uppercase tracking-widest text-[#737373] mb-5">Events</h3>
              {tournamentStats.length > 0 ? (
                <div className="space-y-2">
                  {tournamentStats.map((t, i: number) => (
                    <div key={i} className="flex items-center justify-between bg-[#0a0a0a] border border-[#1a1a1a] p-4">
                      <div>
                        <p className="font-grotesk font-semibold text-white text-sm">
                          {t.tournament_name || "Tournament"}
                        </p>
                      </div>
                      <div className="text-right">
                        <p className="text-xs text-[#737373] uppercase tracking-wider mb-0.5">Maps</p>
                        <p className="font-mono font-bold text-white text-lg">
                          {t.maps_played || 0}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-[#737373] text-sm">No events available</p>
              )}
            </div>
          )}

          {activeTab === "career" && (
            <div>
              <h3 className="text-xs uppercase tracking-widest text-[#737373] mb-5">
                Career by Franchise
              </h3>
              {careerData?.franchises && careerData.franchises.length > 0 ? (
                <div className="space-y-6">
                  {careerData.franchises.map((franchise: PlayerFranchiseEntry, fi: number) => (
                    <div key={fi} className="bg-[#0a0a0a] border border-[#1a1a1a]">
                      {/* Franchise header */}
                      <div className="flex items-center justify-between px-5 py-4 border-b border-[#1a1a1a]">
                        <div>
                          <p className="font-grotesk font-bold text-white text-sm">
                            {franchise.franchise_name}
                          </p>
                          <p className="text-[#737373] text-xs mt-0.5">
                            {franchise.total_matches} matches · {franchise.total_maps} maps
                          </p>
                        </div>
                        <div className="text-right">
                          <p className="text-[10px] text-[#737373] uppercase tracking-widest mb-0.5">Career K/D</p>
                          <p className="font-mono font-bold text-white text-xl">
                            {franchise.career_kd?.toFixed(2) || "0.00"}
                          </p>
                        </div>
                      </div>

                      {/* Era breakdown */}
                      <div className="divide-y divide-[#111111]">
                        {franchise.eras.map((era: PlayerEraStats, ei: number) => (
                          <div key={ei} className="flex items-center justify-between px-5 py-3">
                            <div>
                              <p className="text-[#a3a3a3] text-xs font-medium">{era.team_name}</p>
                              <p className="text-[#737373] text-[10px] mt-0.5 uppercase tracking-wider">
                                {era.season_name || era.game_code}
                              </p>
                            </div>
                            <div className="flex gap-5 text-right">
                              {[
                                { label: "K/D", value: era.kd?.toFixed(2) || "0.00" },
                                { label: "K", value: era.kills ?? 0 },
                                { label: "D", value: era.deaths ?? 0 },
                                { label: "Maps", value: era.maps ?? 0 },
                              ].map(({ label, value }) => (
                                <div key={label}>
                                  <p className="text-[10px] text-[#737373] uppercase mb-0.5">{label}</p>
                                  <p className="font-mono font-bold text-white text-xs">{value}</p>
                                </div>
                              ))}
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-[#737373] text-sm">No career data available</p>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default PlayerDetail;
