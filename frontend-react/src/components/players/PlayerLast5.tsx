import { Link } from "react-router-dom";
import { getKdColorClass } from "../../utils/kdUtils";
import type { MatchHistoryResult } from "../../types";

interface Props {
  matches: MatchHistoryResult[];
}

export default function PlayerLast5({ matches }: Props) {
  if (matches.length === 0) {
    return <p className="text-[#737373] text-sm">No matches available</p>;
  }

  return (
    <div className="space-y-2">
      {matches.map((match, i) => (
        <Link
          key={i}
          to={match.match_id ? `/matches/${match.match_id}` : "#"}
          className="block bg-[#0a0a0a] border border-[#1a1a1a] p-4 hover:border-[#2a2a2a] hover:bg-[#0f0f0f] transition-all"
        >
          <div className="flex justify-between items-center">
            <div className="flex items-center gap-4">
              <span
                className={`font-mono font-bold text-sm w-6 text-center ${
                  match.result?.startsWith("W")
                    ? "text-green-400"
                    : "text-[#737373]"
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
                {
                  label: "K/D",
                  value:
                    typeof match.kd === "number"
                      ? match.kd.toFixed(2)
                      : "0.00",
                  color: getKdColorClass(match.kd),
                },
                { label: "K", value: match.kills || "0", color: "text-white" },
                { label: "D", value: match.deaths || "0", color: "text-white" },
              ].map(({ label, value, color }) => (
                <div key={label}>
                  <p className="text-xs text-[#737373] uppercase mb-0.5">
                    {label}
                  </p>
                  <p className={`font-bold text-sm font-mono ${color}`}>
                    {value}
                  </p>
                </div>
              ))}
              <span className="text-[#404040] text-xs ml-2">→</span>
            </div>
          </div>
        </Link>
      ))}
    </div>
  );
}
