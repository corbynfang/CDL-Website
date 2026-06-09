import { useState, useEffect, useMemo, useCallback, useRef } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useApi } from "../hooks/useApi";
import { isFeatured, isHidden } from "../utils/eventUtils";
import type { Tournament, Team, PaginatedResponse, Player } from "../types";

type QuickLink =
	| { label: string; kind: "search"; query: string }
	| { label: string; kind: "link"; href: string };

const QUICK_LINKS: QuickLink[] = [
	{ label: "Simp", kind: "search", query: "Simp" },
	{ label: "Shotzzy", kind: "search", query: "Shotzzy" },
	{ label: "Scrap", kind: "search", query: "Scrap" },
	{ label: "OpTic Texas", kind: "search", query: "OpTic Texas" },
	{ label: "EWC 2025", kind: "link", href: "/events/esports-world-cup-2025" },
];

const BROWSE_CARDS = [
	{ to: "/players", title: "PLAYERS", desc: "Player stats and career data" },
	{ to: "/teams", title: "TEAMS", desc: "Rosters and franchise history" },
	{ to: "/events", title: "EVENTS", desc: "CDL tournaments and results" },
	{ to: "/stats", title: "STATS", desc: "K/D leaderboards and rankings" },
];

// Each entry is a list of fragments that must all appear in the tournament name (lowercase).
const FEATURED_FRAGS = [
	["esports world cup", "2025"],
	["major", "1", "2023"],
	["major", "3", "2025"],
];

interface SearchResult {
	type: "Player" | "Team" | "Event";
	id: number;
	label: string;
	href: string;
}

function fragMatch(name: string, frags: string[]): boolean {
	const lower = name.toLowerCase();
	return frags.every((f) => lower.includes(f));
}

const Home = () => {
	const [query, setQuery] = useState("");
	const [debouncedQuery, setDebouncedQuery] = useState("");

	useEffect(() => {
		const t = setTimeout(() => setDebouncedQuery(query.trim()), 300);
		return () => clearTimeout(t);
	}, [query]);

	const { data: teamsData } = useApi<Team[]>("/api/v1/teams", {
		cacheTtl: 5 * 60 * 1000,
	});
	const { data: eventsData } = useApi<Tournament[]>("/api/v1/tournaments", {
		cacheTtl: 5 * 60 * 1000,
	});

	const playerSearchUrl = debouncedQuery
		? `/api/v1/players?search=${encodeURIComponent(debouncedQuery)}&limit=5&page=1`
		: "/api/v1/players?limit=1&page=1";

	const { data: playersData, loading: playersLoading } = useApi<
		PaginatedResponse<Player>
	>(playerSearchUrl, { enabled: !!debouncedQuery });

	const results = useMemo((): SearchResult[] => {
		if (!debouncedQuery) return [];
		const q = debouncedQuery.toLowerCase();

		const playerResults: SearchResult[] = (playersData?.data ?? [])
			.slice(0, 5)
			.map((p) => ({
				type: "Player",
				id: p.id,
				label: p.gamertag,
				href: `/players/${p.id}`,
			}));

		const teamResults: SearchResult[] = (teamsData ?? [])
			.filter(
				(t) => t.name !== "Unaffiliated" && t.name.toLowerCase().includes(q),
			)
			.slice(0, 3)
			.map((t) => ({
				type: "Team",
				id: t.id,
				label: t.name,
				href: `/teams/${t.id}`,
			}));

		const eventResults: SearchResult[] = (eventsData ?? [])
			.filter(
				(e) => !isHidden(e.tournament_type) && e.name.toLowerCase().includes(q),
			)
			.slice(0, 3)
			.map((e) => ({
				type: "Event",
				id: e.id,
				label: e.name,
				href: `/events/${e.slug}`,
			}));

		return [...playerResults, ...teamResults, ...eventResults].slice(0, 10);
	}, [debouncedQuery, playersData, teamsData, eventsData]);

	const featuredEvents = useMemo((): Tournament[] => {
		if (!eventsData) return [];
		const matched = FEATURED_FRAGS.map((frags) =>
			eventsData.find((e) => fragMatch(e.name, frags)),
		).filter((e): e is Tournament => !!e);
		if (matched.length > 0) return matched;
		// Fallback: three most recent featured-type events
		return [...eventsData]
			.filter((e) => isFeatured(e.tournament_type))
			.sort((a, b) => b.start_date.localeCompare(a.start_date))
			.slice(0, 3);
	}, [eventsData]);

	const navigate = useNavigate();
	const searchRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const handler = (e: MouseEvent) => {
			if (searchRef.current && !searchRef.current.contains(e.target as Node)) {
				setQuery("");
			}
		};
		document.addEventListener("mousedown", handler);
		return () => document.removeEventListener("mousedown", handler);
	}, []);

	const isSearching = query.trim().length > 0;
	const isPending = query.trim() !== debouncedQuery;
	const isLoading = isPending || playersLoading;

	const handleKeyDown = useCallback(
		(e: React.KeyboardEvent<HTMLInputElement>) => {
			if (e.key === "Escape") {
				setQuery("");
			} else if (e.key === "Enter" && results.length > 0) {
				navigate(results[0].href);
				setQuery("");
			}
		},
		[results, navigate],
	);

	return (
		<div className="min-h-screen bg-[#0a0a0a]">
			<div className="max-w-2xl mx-auto px-4 pt-28 pb-20">
				{/* Hero */}
				<div className="text-center mb-10">
					<h1 className="font-grotesk text-5xl font-bold tracking-tight text-white mb-3">
						CDLytics
					</h1>
					<p className="text-[#a3a3a3] text-lg">
						Search Call of Duty League players, teams, events, and stats.
					</p>
				</div>

				{/* Search bar */}
				<div ref={searchRef} className="relative mb-5">
					<input
						type="text"
						value={query}
						onChange={(e) => setQuery(e.target.value)}
						placeholder="Search players, teams, events..."
						className="w-full bg-[#111111] border border-[#2a2a2a] text-white text-base px-5 py-4 placeholder-[#4a4a4a] focus:outline-none focus:border-[#404040] transition-colors"
						autoComplete="off"
						spellCheck={false}
						onKeyDown={handleKeyDown}
					/>

					{isSearching && (
						<div className="absolute top-full left-0 right-0 mt-px bg-[#111111] border border-[#2a2a2a] z-20 shadow-lg max-h-[60vh] overflow-y-auto">
							{isLoading && results.length === 0 && (
								<div className="px-5 py-3 text-[#737373] text-sm">
									Searching…
								</div>
							)}
							{!isLoading && results.length === 0 && (
								<div className="px-5 py-3 text-[#737373] text-sm">
									No results for &ldquo;{debouncedQuery}&rdquo;
								</div>
							)}
							{results.map((r, i) => (
								<Link
									key={`${r.type}-${r.id}`}
									to={r.href}
									onClick={() => setQuery("")}
									className={`flex items-center justify-between px-5 py-3 hover:bg-[#1a1a1a] transition-colors${
										i < results.length - 1 ? " border-b border-[#1a1a1a]" : ""
									}`}
								>
									<span className="text-white text-sm">{r.label}</span>
									<span className="text-[#737373] text-xs uppercase tracking-widest ml-3 shrink-0">
										{r.type}
									</span>
								</Link>
							))}
						</div>
					)}
				</div>

				{/* Quick links */}
				<div className="flex flex-wrap gap-2 justify-center mb-14">
					{QUICK_LINKS.map((link) =>
						link.kind === "link" ? (
							<Link
								key={link.label}
								to={link.href}
								className="px-3 py-1.5 text-xs text-[#737373] border border-[#1a1a1a] hover:text-white hover:border-[#2a2a2a] transition-colors"
							>
								{link.label}
							</Link>
						) : (
							<button
								key={link.label}
								type="button"
								onClick={() => setQuery(link.query)}
								className="px-3 py-1.5 text-xs text-[#737373] border border-[#1a1a1a] hover:text-white hover:border-[#2a2a2a] transition-colors"
							>
								{link.label}
							</button>
						),
					)}
				</div>

				{/* Browse cards + featured events — hidden while search is open */}
				{!isSearching && (
					<>
						<div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-12">
							{BROWSE_CARDS.map(({ to, title, desc }) => (
								<Link
									key={to}
									to={to}
									className="p-5 bg-[#111111] border border-[#1a1a1a] hover:border-[#2a2a2a] hover:bg-[#161616] transition-all"
								>
									<h2 className="font-grotesk text-xs font-bold tracking-widest text-white mb-1.5">
										{title}
									</h2>
									<p className="text-[#737373] text-xs leading-relaxed">
										{desc}
									</p>
								</Link>
							))}
						</div>

						{featuredEvents.length > 0 && (
							<div>
								<p className="text-xs uppercase tracking-widest text-[#737373] mb-4">
									Featured Events
								</p>
								<div className="flex flex-col gap-2">
									{featuredEvents.map((e) => (
										<Link
											key={e.id}
											to={`/events/${e.slug}`}
											className="flex items-center justify-between px-5 py-4 bg-[#111111] border border-[#1a1a1a] hover:border-[#2a2a2a] hover:bg-[#161616] transition-all"
										>
											<span className="text-white text-sm font-medium">
												{e.name}
											</span>
											<span className="text-[#737373] text-xs uppercase tracking-widest">
												{e.location ?? e.tournament_type.replace(/_/g, " ")}
											</span>
										</Link>
									))}
								</div>
							</div>
						)}
					</>
				)}
			</div>
		</div>
	);
};

export default Home;
