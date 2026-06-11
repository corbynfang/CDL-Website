import { Link } from "react-router-dom";
import { getKdColorClass } from "../../utils/kdUtils";
import type { MatchHistoryEvent } from "../../types";

interface Props {
  events: MatchHistoryEvent[];
}

export default function PlayerMatchHistory({ events }: Props) {
  if (events.length === 0) {
    return <p className="text-[#737373] text-sm">No match data available</p>;
  }

  return (
    <div className="space-y-8">
      {events.map((event, ei) => (
        <div key={ei}>
          <h3 className="text-xs uppercase tracking-widest text-[#737373] mb-4">
            {event.event} {event.year}
          </h3>
          {event.matches?.length > 0 ? (
            <div className="overflow-x-auto border border-[#1a1a1a]">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-[#1a1a1a] bg-[#0a0a0a]">
                    {[
                      "Date",
                      "Opp",
                      "Result",
                      "KD",
                      "K",
                      "D",
                      "HP KD",
                      "SND KD",
                      "CTL KD",
                      "Slayer",
                      "Rating",
                    ].map((h) => (
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
                  {event.matches.map((match, mi) => (
                    <tr
                      key={mi}
                      className="border-b border-[#1a1a1a] hover:bg-[#0a0a0a] transition-colors cursor-pointer"
                      onClick={() =>
                        match.match_id &&
                        (window.location.href = `/matches/${match.match_id}`)
                      }
                    >
                      <td className="px-3 py-2 text-[#737373] text-xs font-mono">
                        {match.date
                          ? new Date(match.date).toLocaleDateString()
                          : "—"}
                      </td>
                      <td className="px-3 py-2 text-[#a3a3a3] text-xs font-medium">
                        {match.opponent_abbr || match.opponent || "—"}
                      </td>
                      <td
                        className={`px-3 py-2 text-xs font-bold font-mono ${match.result?.startsWith("W") ? "text-green-400" : "text-[#737373]"}`}
                      >
                        {match.result || "—"}
                      </td>
                      <td
                        className={`px-3 py-2 text-xs font-bold font-mono ${getKdColorClass(match.kd)}`}
                      >
                        {typeof match.kd === "number"
                          ? match.kd.toFixed(2)
                          : "0.00"}
                      </td>
                      <td className="px-3 py-2 text-[#737373] text-xs font-mono">
                        {match.kills || "0"}
                      </td>
                      <td className="px-3 py-2 text-[#737373] text-xs font-mono">
                        {match.deaths || "0"}
                      </td>
                      <td className="px-3 py-2 text-[#737373] text-xs font-mono">—</td>
                      <td className="px-3 py-2 text-[#737373] text-xs font-mono">—</td>
                      <td className="px-3 py-2 text-[#737373] text-xs font-mono">—</td>
                      <td className="px-3 py-2 text-[#737373] text-xs font-mono">—</td>
                      <td className="px-3 py-2 text-[#737373] text-xs font-mono text-right">
                        {match.match_id ? (
                          <Link
                            to={`/matches/${match.match_id}`}
                            className="text-[#404040] hover:text-[#737373] transition-colors"
                            onClick={(e) => e.stopPropagation()}
                          >
                            →
                          </Link>
                        ) : (
                          "—"
                        )}
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
  );
}
