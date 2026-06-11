import { getTeamLogo } from "../../utils/assets";
import type { TournamentTeam } from "../../types";

interface Props {
	teams: TournamentTeam[];
	max?: number;
	size?: "sm" | "md";
}

export default function TeamLogoStrip({ teams, max = 12, size = "sm" }: Props) {
	const visible = teams.slice(0, max);
	const dim = size === "md" ? "w-9 h-9" : "w-6 h-6";

	return (
		<div className="flex items-center gap-2 flex-wrap">
			{visible.map((team) => {
				const logo = getTeamLogo(team.name);
				return logo ? (
					<img
						key={team.id}
						src={logo}
						alt={team.name}
						title={team.name}
						className={`${dim} object-contain opacity-60`}
					/>
				) : (
					<div
						key={team.id}
						className={`${dim} rounded-full bg-zinc-800 flex items-center justify-center text-[9px] font-mono text-zinc-500`}
					>
						{(team.abbreviation ?? "?").slice(0, 2)}
					</div>
				);
			})}
			{teams.length > max && (
				<span className="text-xs text-zinc-600">+{teams.length - max}</span>
			)}
		</div>
	);
}
