import { hasBracket } from "../../utils/eventUtils";

export type TabId = "overview" | "bracket" | "matches" | "teams" | "stats";

interface Props {
	active: TabId;
	onSelect: (tab: TabId) => void;
	tournamentType: string;
	hasStats: boolean;
}

export default function EventTabs({
	active,
	onSelect,
	tournamentType,
	hasStats,
}: Props) {
	const tabs: { id: TabId; label: string }[] = [
		{ id: "overview", label: "Overview" },
		...(hasBracket(tournamentType)
			? [{ id: "bracket" as TabId, label: "Bracket" }]
			: []),
		{ id: "matches", label: "Matches" },
		{ id: "teams", label: "Teams" },
		...(hasStats ? [{ id: "stats" as TabId, label: "Stats" }] : []),
	];

	return (
		<div className="border-b border-[#1a1a1a] bg-[#0a0a0a] sticky top-14 z-30">
			<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div className="flex gap-0 overflow-x-auto scrollbar-none">
					{tabs.map((tab) => (
						<button
							key={tab.id}
							onClick={() => onSelect(tab.id)}
							className={`px-4 py-3 text-xs uppercase tracking-widest whitespace-nowrap transition-colors border-b-2 -mb-px ${
								active === tab.id
									? "text-white border-white"
									: "text-zinc-600 border-transparent hover:text-zinc-400"
							}`}
						>
							{tab.label}
						</button>
					))}
				</div>
			</div>
		</div>
	);
}
