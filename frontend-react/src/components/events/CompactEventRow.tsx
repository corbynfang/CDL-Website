import { Link } from "react-router-dom";
import type { Tournament } from "../../types";
import {
	deriveStatus,
	eventDisplayName,
	formatDateRange,
	formatPrize,
	countryFlag,
} from "../../utils/eventUtils";

interface Props {
	event: Tournament;
}

const STATUS_DOT: Record<string, string> = {
	live: "bg-emerald-500",
	upcoming: "bg-blue-500",
	completed: "bg-zinc-600",
};

export default function CompactEventRow({ event }: Props) {
	const status = deriveStatus(event.start_date, event.end_date);

	return (
		<Link
			to={`/events/${event.slug}`}
			className="group flex items-center gap-3 px-4 py-3 border border-[#1a1a1a] hover:border-[#2a2a2a] hover:bg-[#111111] transition-all"
		>
			<span
				className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${STATUS_DOT[status]}`}
			/>

			<span className="font-grotesk text-sm font-medium text-[#a3a3a3] group-hover:text-white transition-colors flex-1 min-w-0 truncate">
				{eventDisplayName(event.slug, event.name)}
			</span>

			<span className="text-xs text-zinc-600 whitespace-nowrap hidden sm:block">
				{formatDateRange(event.start_date, event.end_date)}
			</span>

			{event.location && (
				<span className="text-xs text-zinc-700 whitespace-nowrap hidden md:block">
					{event.country ? countryFlag(event.country) + " " : ""}
					{event.location}
				</span>
			)}

			{event.prize_pool ? (
				<span className="text-xs text-zinc-500 whitespace-nowrap">
					{formatPrize(event.prize_pool)}
				</span>
			) : (
				<span className="text-xs text-zinc-700 whitespace-nowrap">—</span>
			)}
		</Link>
	);
}
