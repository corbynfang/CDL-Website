import { Link } from "react-router-dom";

export interface PlayerStat {
  player_id: number;
  gamertag: string;
  kills: number;
  deaths: number;
  kd_ratio: number;
  damage: number;
  assists: number;
  bp_rating: number;
  hill_time: number;
  snd_rounds: number;
  plant_count: number;
  defuse_count: number;
  first_blood_count: number;
  first_death_count: number;
  non_traded_kills: number;
  highest_streak: number;
  data_quality_note?: string;
}

export interface MapDetail {
  map_number: number;
  map_name: string;
  mode: string;
  score_1: number;
  score_2: number;
  winner_id?: number;
  duration_sec: number;
  played: boolean;
  team1_stats: PlayerStat[];
  team2_stats: PlayerStat[];
}

const formatDuration = (sec: number) => {
  if (!sec) return "—";
  const m = Math.floor(sec / 60);
  const s = sec % 60;
  return `${m}:${String(s).padStart(2, "0")}`;
};

const kdColor = (kd: number) => {
  if (kd >= 1.5) return "text-green-400";
  if (kd >= 1.0) return "text-white";
  if (kd >= 0.8) return "text-[#a3a3a3]";
  return "text-[#737373]";
};

const ModeIcon = ({ mode }: { mode: string }) => {
  const icons: Record<string, string> = {
    Hardpoint: "HP",
    "Search and Destroy": "SND",
    Control: "CTL",
  };
  return (
    <span className="font-mono text-[10px] text-[#737373]">
      {icons[mode] ?? mode.slice(0, 3).toUpperCase()}
    </span>
  );
};

interface ScoreboardProps {
  map: MapDetail;
  team1Name: string;
  team2Name: string;
  team1ID: number;
}

export default function Scoreboard({
  map,
  team1Name,
  team2Name,
  team1ID,
}: ScoreboardProps) {
  const team1Won = map.winner_id === team1ID;
  const isHP = map.mode === "Hardpoint";
  const isSND = map.mode === "Search and Destroy";

  const StatHeader = ({ label }: { label: string }) => (
    <th className="px-2 py-2 text-[#737373] text-[10px] uppercase tracking-widest font-medium text-right">
      {label}
    </th>
  );

  const PlayerRow = ({
    p,
    highlight,
  }: {
    p: PlayerStat;
    highlight: boolean;
  }) => (
    <tr
      className={`border-b border-[#111111] transition-colors ${highlight ? "bg-[#0f0f0f]" : ""}`}
    >
      <td className="px-3 py-2">
        <Link
          to={`/players/${p.player_id}`}
          className="font-grotesk font-semibold text-white hover:text-[#a3a3a3] text-xs transition-colors"
        >
          {p.gamertag}
        </Link>
      </td>
      <td
        className={`px-2 py-2 text-xs font-bold font-mono text-right ${kdColor(p.kd_ratio)}`}
      >
        {p.kd_ratio.toFixed(2)}
      </td>
      <td className="px-2 py-2 text-xs font-mono text-right text-white">
        {p.kills}
      </td>
      <td className="px-2 py-2 text-xs font-mono text-right text-[#737373]">
        {p.deaths}
      </td>
      <td className="px-2 py-2 text-xs font-mono text-right text-[#737373]">
        {p.damage > 0 ? p.damage.toLocaleString() : "—"}
      </td>
      {isHP && (
        <td className="px-2 py-2 text-xs font-mono text-right text-[#737373]">
          {p.hill_time > 0 ? `${p.hill_time}s` : "—"}
        </td>
      )}
      {isSND && (
        <>
          <td className="px-2 py-2 text-xs font-mono text-right text-[#737373]">
            {p.plant_count > 0 ? p.plant_count : "—"}
          </td>
          <td className="px-2 py-2 text-xs font-mono text-right text-[#737373]">
            {p.defuse_count > 0 ? p.defuse_count : "—"}
          </td>
          <td className="px-2 py-2 text-xs font-mono text-right text-[#737373]">
            {p.first_blood_count > 0 ? p.first_blood_count : "—"}
          </td>
        </>
      )}
      <td className="px-2 py-2 text-xs font-mono text-right text-[#737373]">
        {p.bp_rating > 0 ? p.bp_rating.toFixed(2) : "—"}
      </td>
    </tr>
  );

  const teamTable = (
    players: PlayerStat[],
    teamName: string,
    score: number,
    won: boolean,
  ) => (
    <div className="flex-1 min-w-0">
      <div
        className={`flex items-center justify-between px-3 py-2 border-b border-[#1a1a1a] ${won ? "bg-white/[0.03]" : ""}`}
      >
        <span
          className={`font-grotesk font-bold text-xs ${won ? "text-white" : "text-[#737373]"}`}
        >
          {teamName}
          {won && (
            <span className="ml-2 text-[10px] text-green-400 uppercase tracking-widest">
              W
            </span>
          )}
        </span>
        <span
          className={`font-mono font-bold text-lg ${won ? "text-white" : "text-[#737373]"}`}
        >
          {score}
        </span>
      </div>
      <table className="w-full">
        <thead>
          <tr className="border-b border-[#1a1a1a]">
            <th className="px-3 py-1.5 text-left text-[#737373] text-[10px] uppercase tracking-widest font-medium">
              Player
            </th>
            <StatHeader label="K/D" />
            <StatHeader label="K" />
            <StatHeader label="D" />
            <StatHeader label="DMG" />
            {isHP && <StatHeader label="Hill" />}
            {isSND && (
              <>
                <StatHeader label="Plants" />
                <StatHeader label="Defuses" />
                <StatHeader label="FB" />
              </>
            )}
            <StatHeader label="Rating" />
          </tr>
        </thead>
        <tbody>
          {players.map((p, i) => (
            <PlayerRow key={p.player_id} p={p} highlight={i % 2 === 0} />
          ))}
        </tbody>
      </table>
    </div>
  );

  return (
    <div className="border border-[#1a1a1a] overflow-hidden">
      <div className="flex items-center gap-3 px-4 py-3 bg-[#0a0a0a] border-b border-[#1a1a1a]">
        <span className="text-white font-grotesk font-bold text-sm">
          Map {map.map_number}
        </span>
        <span className="text-[#737373] text-xs">·</span>
        <span className="text-[#a3a3a3] text-xs">{map.map_name}</span>
        <span className="text-[#737373] text-xs">·</span>
        <ModeIcon mode={map.mode} />
        {map.duration_sec > 0 && (
          <>
            <span className="text-[#737373] text-xs">·</span>
            <span className="text-[#737373] text-xs">
              {formatDuration(map.duration_sec)}
            </span>
          </>
        )}
      </div>

      <div className="flex flex-col lg:flex-row divide-y lg:divide-y-0 lg:divide-x divide-[#1a1a1a]">
        {teamTable(map.team1_stats, team1Name, map.score_1, team1Won)}
        {teamTable(map.team2_stats, team2Name, map.score_2, !team1Won)}
      </div>
    </div>
  );
}
