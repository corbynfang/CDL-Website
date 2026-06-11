import { getKdColorClass } from "../../utils/kdUtils";
import type { PlayerKDTournamentEntry } from "../../types";

interface Props {
  tournamentStats: PlayerKDTournamentEntry[];
}

export default function PlayerEventStats({ tournamentStats }: Props) {
  if (tournamentStats.length === 0) {
    return (
      <p className="text-[#737373] text-sm">No event statistics available</p>
    );
  }

  return (
    <div className="space-y-2">
      {tournamentStats.map((t, i) => (
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
              <p className="text-xs text-[#737373] uppercase tracking-wider mb-0.5">
                K/D
              </p>
              <p
                className={`font-mono font-bold text-xl ${getKdColorClass(t.kd_ratio)}`}
              >
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
                <p className="text-xs text-[#737373] uppercase tracking-wider mb-1">
                  {label}
                </p>
                <p className="font-mono font-bold text-white">
                  {value?.toLocaleString() || "0"}
                </p>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
