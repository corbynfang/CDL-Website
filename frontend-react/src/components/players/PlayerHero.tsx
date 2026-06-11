import { getPlayerAvatar } from "../../utils/assets";
import { getKdColorClass } from "../../utils/kdUtils";
import type { Player, PlayerKDResponse } from "../../types";

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

interface Props {
  player: Player;
  stats: PlayerKDResponse | null;
}

export default function PlayerHero({ player, stats }: Props) {
  const overallKD = stats?.avg_kd || 0;
  const hpKD = stats?.hp_kd_ratio || 0;
  const sndKD = stats?.snd_kd_ratio || 0;
  const ctlKD = stats?.control_kd_ratio || 0;

  return (
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
            <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">
              Gamertag
            </p>
            <p className="font-grotesk text-2xl font-bold text-white">
              {player.gamertag}
            </p>
          </div>

          {(player.first_name || player.last_name) && (
            <div>
              <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">
                Name
              </p>
              <p className="text-[#a3a3a3] text-sm">
                {player.first_name} {player.last_name}
              </p>
            </div>
          )}

          {player.country && (
            <div>
              <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">
                Country
              </p>
              <p className="text-[#a3a3a3] text-sm">{player.country}</p>
            </div>
          )}

          {player.birthdate && (
            <div>
              <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">
                Birthday
              </p>
              <p className="text-[#a3a3a3] text-sm">
                {formatBirthdate(player.birthdate)}
              </p>
            </div>
          )}

          {player.role && (
            <div>
              <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">
                Role
              </p>
              <p className="text-[#a3a3a3] text-sm">{player.role}</p>
            </div>
          )}

          <div>
            <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">
              Status
            </p>
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
        <h2 className="text-xs uppercase tracking-widest text-[#737373] mb-6">
          K/D Statistics
        </h2>

        <div className="space-y-5">
          {[
            { label: "Overall", value: overallKD },
            { label: "Hardpoint", value: hpKD },
            { label: "Search & Destroy", value: sndKD },
            { label: "Control", value: ctlKD },
          ].map(({ label, value }) => (
            <div key={label}>
              <div className="flex justify-between items-baseline mb-2">
                <p className="text-xs text-[#737373]">{label}</p>
                <p
                  className={`font-mono font-bold text-sm ${getKdColorClass(value)}`}
                >
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
              <p className="text-xs text-[#737373] uppercase tracking-wider mb-1">
                {label}
              </p>
              <p className="font-mono font-bold text-white text-lg">
                {value?.toLocaleString() || "0"}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

