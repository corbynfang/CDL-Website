import type {
  PlayerCareerResponse,
  PlayerFranchiseEntry,
  PlayerEraStats,
} from "../../types";

interface Props {
  careerData: PlayerCareerResponse | null;
}

export default function PlayerCareer({ careerData }: Props) {
  if (!careerData?.franchises || careerData.franchises.length === 0) {
    return <p className="text-[#737373] text-sm">No career data available</p>;
  }

  return (
    <div className="space-y-6">
      {careerData.franchises.map((franchise: PlayerFranchiseEntry, fi) => (
        <div key={fi} className="bg-[#0a0a0a] border border-[#1a1a1a]">
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
              <p className="text-[10px] text-[#737373] uppercase tracking-widest mb-0.5">
                Career K/D
              </p>
              <p className="font-mono font-bold text-white text-xl">
                {franchise.career_kd?.toFixed(2) || "0.00"}
              </p>
            </div>
          </div>

          <div className="divide-y divide-[#111111]">
            {franchise.eras.map((era: PlayerEraStats, ei) => (
              <div
                key={ei}
                className="flex items-center justify-between px-5 py-3"
              >
                <div>
                  <p className="text-[#a3a3a3] text-xs font-medium">
                    {era.team_name}
                  </p>
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
                      <p className="text-[10px] text-[#737373] uppercase mb-0.5">
                        {label}
                      </p>
                      <p className="font-mono font-bold text-white text-xs">
                        {value}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
