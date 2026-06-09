import type { Tournament } from "../../types";
import {
	formatPrize,
	formatDateRange,
	countryFlag,
	formatTournamentFormat,
} from "../../utils/eventUtils";

interface Props {
	event: Tournament;
	teamCount: number;
}

function StatCard({ label, value }: { label: string; value: string }) {
	return (
		<div className="p-4 border border-[#1a1a1a] bg-[#111111] space-y-1">
			<p className="text-xs uppercase tracking-widest text-zinc-600">{label}</p>
			<p className="font-grotesk text-xl font-bold text-white">{value}</p>
		</div>
	);
}

export default function EventOverview({ event, teamCount }: Props) {
	return (
		<div className="space-y-8">
			<div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
				<StatCard
					label="Teams"
					value={teamCount > 0 ? String(teamCount) : "—"}
				/>
				<StatCard label="Prize Pool" value={formatPrize(event.prize_pool)} />
				<StatCard
					label="Format"
					value={formatTournamentFormat(event.tournament_format)}
				/>
				<StatCard label="Type" value={event.is_lan ? "LAN" : "Online"} />
			</div>

			<div className="space-y-2">
				<p className="text-xs uppercase tracking-widest text-zinc-600">
					Details
				</p>
				<div className="border border-[#1a1a1a] divide-y divide-[#1a1a1a]">
					{[
						{
							label: "Dates",
							value: formatDateRange(event.start_date, event.end_date),
						},
						event.location
							? {
									label: "Location",
									value: `${event.country ? countryFlag(event.country) + " " : ""}${event.location}`,
								}
							: null,
						event.season?.name
							? { label: "Season", value: event.season.name }
							: null,
					]
						.filter(Boolean)
						.map((row) => (
							<div key={row!.label} className="flex gap-4 px-4 py-3">
								<span className="text-xs text-zinc-600 w-24 flex-shrink-0">
									{row!.label}
								</span>
								<span className="text-sm text-zinc-300">{row!.value}</span>
							</div>
						))}
				</div>
			</div>

			{event.source_event_url && (
				<div className="space-y-2">
					<p className="text-xs uppercase tracking-widest text-zinc-600">
						External Links
					</p>
					<div className="flex flex-wrap gap-3">
						<a
							href={event.source_event_url}
							target="_blank"
							rel="noopener noreferrer"
							className="text-xs uppercase tracking-widest text-zinc-500 hover:text-white border border-[#1a1a1a] hover:border-[#2a2a2a] px-4 py-2 transition-colors"
						>
							Event Page →
						</a>
					</div>
				</div>
			)}
		</div>
	);
}
