import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import { useApi } from "../../hooks/useApi";
import PageMeta from "../PageMeta";
import PlayerHero from "./PlayerHero";
import PlayerLast5 from "./PlayerLast5";
import PlayerMatchHistory from "./PlayerMatchHistory";
import PlayerEventStats from "./PlayerEventStats";
import PlayerCareer from "./PlayerCareer";
import type {
  Player,
  PlayerKDResponse,
  PlayerMatchHistory as PlayerMatchHistoryType,
  PlayerCareerResponse,
  MatchHistoryEvent,
  MatchHistoryResult,
  PlayerKDTournamentEntry,
} from "../../types";

const TABS = [
  { id: "last5", label: "Last 5" },
  { id: "matches", label: "Matches" },
  { id: "eventStats", label: "Event Stats" },
  { id: "events", label: "Events" },
  { id: "career", label: "Career" },
];

const PlayerDetail = () => {
  const { id } = useParams<{ id: string }>();
  const [activeTab, setActiveTab] = useState("last5");

  const {
    data: player,
    loading: playerLoading,
    error: playerError,
  } = useApi<Player>(`/api/v1/players/${id}`);
  const { data: stats, loading: statsLoading } = useApi<PlayerKDResponse>(
    `/api/v1/players/${id}/kd`,
  );
  const { data: matchesData, loading: matchesLoading } =
    useApi<PlayerMatchHistoryType>(`/api/v1/players/${id}/matches`);
  const { data: careerData } = useApi<PlayerCareerResponse>(
    `/api/v1/players/${id}/franchise-career`,
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
        <Link
          to="/players"
          className="text-[#a3a3a3] hover:text-white mt-4 inline-block text-sm transition-colors"
        >
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
  const tournamentStats: PlayerKDTournamentEntry[] =
    stats?.tournament_stats || [];

  const overallKD = stats?.avg_kd || 0;
  const kdStr = overallKD > 0 ? ` — ${overallKD.toFixed(2)} K/D` : "";
  const teamName = careerData?.franchises?.[0]?.franchise_name;
  const metaDesc = `${player.gamertag}${teamName ? ` (${teamName})` : ""} CDL statistics${kdStr}. Match history, event stats, and career K/D breakdown.`;

  const personSchema = {
    "@context": "https://schema.org",
    "@type": "Person",
    name: player.gamertag,
    alternateName: [player.first_name, player.last_name].filter(Boolean).join(" ") || undefined,
    nationality: player.country || undefined,
    url: `https://cdlytics.com/players/${id}`,
    description: metaDesc,
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <PageMeta
        title={`${player.gamertag} CDL Stats`}
        description={metaDesc}
        canonical={`/players/${id}`}
        type="profile"
      />
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify(personSchema) }}
      />
      <Link
        to="/players"
        className="text-[#737373] hover:text-[#a3a3a3] text-xs uppercase tracking-widest mb-10 inline-block transition-colors"
      >
        ← Players
      </Link>

      <PlayerHero player={player} stats={stats} />

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
              <PlayerLast5 matches={last5Matches} />
            </div>
          )}

          {activeTab === "matches" && (
            <PlayerMatchHistory events={events} />
          )}

          {activeTab === "eventStats" && (
            <div>
              <h3 className="text-xs uppercase tracking-widest text-[#737373] mb-5">
                Event Statistics
              </h3>
              <PlayerEventStats tournamentStats={tournamentStats} />
            </div>
          )}

          {activeTab === "events" && (
            <div>
              <h3 className="text-xs uppercase tracking-widest text-[#737373] mb-5">
                Events
              </h3>
              {tournamentStats.length > 0 ? (
                <div className="space-y-2">
                  {tournamentStats.map((t, i) => (
                    <div
                      key={i}
                      className="flex items-center justify-between bg-[#0a0a0a] border border-[#1a1a1a] p-4"
                    >
                      <p className="font-grotesk font-semibold text-white text-sm">
                        {t.tournament_name || "Tournament"}
                      </p>
                      <div className="text-right">
                        <p className="text-xs text-[#737373] uppercase tracking-wider mb-0.5">
                          Maps
                        </p>
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
              <PlayerCareer careerData={careerData} />
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default PlayerDetail;
