import { Link } from "react-router-dom";
import type { Match } from "../../types";
import { getTeamLogo } from "../../utils/assets";
import MatchCardSkeleton from "../loaders/MatchCardSkeleton";

interface Props {
	matches: Match[] | null;
	loading: boolean;
	error: string | null;
}

function MatchRow({ match }: { match: Match }) {
	const t1 = match.team1;
	const t2 = match.team2;
	const t1Won = match.winner_id === match.team1_id;
	const t2Won = match.winner_id === match.team2_id;

	function TeamCell({
		name,
		abbr,
		won,
	}: {
		name?: string;
		abbr?: string;
		won: boolean;
	}) {
		const logo = name ? getTeamLogo(name) : undefined;
		return (
			<div className={`flex items-center gap-2 ${won ? "" : "opacity-50"}`}>
				{logo ? (
					<img
						src={logo}
						alt={name}
						className="w-6 h-6 object-contain flex-shrink-0"
					/>
				) : (
					<div className="w-6 h-6 rounded-full bg-zinc-800 flex-shrink-0" />
				)}
				<span
					className={`text-sm ${won ? "text-white font-semibold" : "text-zinc-400"} hidden sm:block`}
				>
					{name ?? "—"}
				</span>
				<span
					className={`text-sm ${won ? "text-white font-semibold" : "text-zinc-400"} sm:hidden`}
				>
					{abbr ?? "?"}
				</span>
			</div>
		);
	}

	return (
		<Link
			to={`/matches/${match.id}`}
			className="group flex items-center gap-3 px-4 py-3 border border-[#1a1a1a] hover:border-[#2a2a2a] hover:bg-[#111111] transition-all"
		>
			<div className="flex-1">
				<TeamCell name={t1?.name} abbr={t1?.abbreviation} won={t1Won} />
			</div>

			<div className="flex items-center gap-2 text-sm font-bold tabular-nums flex-shrink-0">
				<span className={t1Won ? "text-white" : "text-zinc-600"}>
					{match.team1_score}
				</span>
				<span className="text-zinc-700">–</span>
				<span className={t2Won ? "text-white" : "text-zinc-600"}>
					{match.team2_score}
				</span>
			</div>

			<div className="flex-1 flex justify-end">
				<TeamCell name={t2?.name} abbr={t2?.abbreviation} won={t2Won} />
			</div>

			{match.match_type && (
				<span className="text-[10px] uppercase tracking-widest text-zinc-700 hidden lg:block w-32 text-right">
					{match.match_type.replace(/_/g, " ")}
				</span>
			)}
		</Link>
	);
}

export default function EventMatches({ matches, loading, error }: Props) {
	if (loading) {
		return (
			<div className="space-y-2">
				{Array.from({ length: 6 }).map((_, i) => (
					<MatchCardSkeleton key={i} />
				))}
			</div>
		);
	}

	if (error) {
		return (
			<p className="text-center text-zinc-600 py-16 text-sm">
				Could not load matches.
			</p>
		);
	}

	if (!matches || matches.length === 0) {
		return (
			<p className="text-center text-zinc-600 py-16 text-sm">
				No matches recorded yet.
			</p>
		);
	}

	const groups = new Map<string, Match[]>();
	for (const m of matches) {
		const key = m.match_type ?? "Other";
		if (!groups.has(key)) groups.set(key, []);
		groups.get(key)!.push(m);
	}

	return (
		<div className="space-y-6">
			{[...groups.entries()].map(([round, roundMatches]) => (
				<div key={round} className="space-y-2">
					<p className="text-[10px] uppercase tracking-widest text-zinc-600">
						{round.replace(/_/g, " ")}
					</p>
					{roundMatches.map((m) => (
						<MatchRow key={m.id} match={m} />
					))}
				</div>
			))}
		</div>
	);
}
