import { useParams, Link } from "react-router-dom";
import { useApi } from "../../hooks/useApi";
import { getTeamLogo } from "../../utils/assets";
import MatchThread from "../threads/MatchThread";
import Scoreboard from "./Scoreboard";
import type { MapDetail } from "./Scoreboard";

interface MatchInfo {
  id: number;
  tournament_name: string;
  season_name: string;
  game_code: string;
  team1_id: number;
  team1_name: string;
  team1_abbr: string;
  team1_logo: string;
  team2_id: number;
  team2_name: string;
  team2_abbr: string;
  team2_logo: string;
  team1_score: number;
  team2_score: number;
  winner_id?: number;
  match_date: string;
  format: string;
  bracket_round: string;
}

interface MatchDetailResponse {
  match: MatchInfo;
  maps: MapDetail[];
}

const formatRound = (raw: string) => {
  const map: Record<string, string> = {
    major_qualifier: "Major Qualifier",
    winners_r1: "Winners Round 1",
    winners_r2: "Winners Round 2",
    winners_r3: "Winners Round 3",
    winners_finals: "Winners Finals",
    elim_r1: "Elimination Round 1",
    elim_r2: "Elimination Round 2",
    elim_r3: "Elimination Round 3",
    elim_r4: "Elimination Round 4",
    elim_finals: "Elimination Finals",
    finals: "Finals",
    grand_finals: "Grand Finals",
    "3rd_place": "3rd Place",
  };
  return (
    map[raw] ?? raw.replace(/_/g, " ").replace(/\b\w/g, (c) => c.toUpperCase())
  );
};

const MatchDetail = () => {
  const { id } = useParams<{ id: string }>();
  const { data, loading, error } = useApi<MatchDetailResponse>(
    `/api/v1/matches/${id}`,
  );

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Loading match...</p>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Match not found</p>
      </div>
    );
  }

  const { match, maps } = data;
  const team1Won = match.winner_id === match.team1_id;
  const team2Won = match.winner_id === match.team2_id;
  const team1Logo = getTeamLogo(match.team1_name);
  const team2Logo = getTeamLogo(match.team2_name);
  const playedMaps = maps.filter((m) => m.played);

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-xs text-[#737373] uppercase tracking-widest mb-10">
        <span>{match.season_name}</span>
        <span>·</span>
        <span>{match.tournament_name}</span>
        {match.bracket_round && (
          <>
            <span>·</span>
            <span>{formatRound(match.bracket_round)}</span>
          </>
        )}
      </div>

      {/* Match header */}
      <div className="bg-[#111111] border border-[#1a1a1a] p-6 mb-8">
        <div className="flex items-center justify-between gap-4">
          <Link
            to={`/teams/${match.team1_id}`}
            className="flex flex-col items-center gap-3 flex-1 group"
          >
            {team1Logo ? (
              <img
                src={team1Logo}
                alt={match.team1_name}
                className="w-20 h-20 object-contain opacity-90 group-hover:opacity-100 transition-opacity"
              />
            ) : (
              <div className="w-20 h-20 bg-[#1a1a1a] flex items-center justify-center text-sm font-bold text-[#737373] font-mono">
                {match.team1_abbr}
              </div>
            )}
            <div className="text-center">
              <p
                className={`font-grotesk font-bold text-sm ${team1Won ? "text-white" : "text-[#737373]"}`}
              >
                {match.team1_name}
              </p>
            </div>
          </Link>

          <div className="text-center flex-shrink-0 px-6">
            <div className="flex items-center gap-4">
              <span
                className={`font-mono font-black text-5xl ${team1Won ? "text-white" : "text-[#404040]"}`}
              >
                {match.team1_score}
              </span>
              <span className="text-[#404040] font-mono text-3xl">—</span>
              <span
                className={`font-mono font-black text-5xl ${team2Won ? "text-white" : "text-[#404040]"}`}
              >
                {match.team2_score}
              </span>
            </div>
            <p className="text-[#737373] text-xs uppercase tracking-widest mt-2">
              {match.format}
            </p>
            <p className="text-[#737373] text-xs mt-1">
              {new Date(match.match_date).toLocaleDateString("en-US", {
                year: "numeric",
                month: "long",
                day: "numeric",
              })}
            </p>
          </div>

          <Link
            to={`/teams/${match.team2_id}`}
            className="flex flex-col items-center gap-3 flex-1 group"
          >
            {team2Logo ? (
              <img
                src={team2Logo}
                alt={match.team2_name}
                className="w-20 h-20 object-contain opacity-90 group-hover:opacity-100 transition-opacity"
              />
            ) : (
              <div className="w-20 h-20 bg-[#1a1a1a] flex items-center justify-center text-sm font-bold text-[#737373] font-mono">
                {match.team2_abbr}
              </div>
            )}
            <div className="text-center">
              <p
                className={`font-grotesk font-bold text-sm ${team2Won ? "text-white" : "text-[#737373]"}`}
              >
                {match.team2_name}
              </p>
            </div>
          </Link>
        </div>
      </div>

      {/* Map-by-map scoreboards */}
      <div className="mb-4">
        <p className="text-xs uppercase tracking-widest text-[#737373] mb-4">
          Map Breakdown · {playedMaps.length} maps played
        </p>
        <div className="space-y-4">
          {playedMaps.length > 0 ? (
            playedMaps.map((map) => (
              <Scoreboard
                key={map.map_number}
                map={map}
                team1Name={match.team1_name}
                team2Name={match.team2_name}
                team1ID={match.team1_id}
              />
            ))
          ) : (
            <div className="border border-[#1a1a1a] p-12 text-center">
              <p className="text-[#737373] text-sm">
                Map stats not available for this match
              </p>
            </div>
          )}
        </div>
      </div>

      <MatchThread matchId={match.id} />
    </div>
  );
};

export default MatchDetail;
