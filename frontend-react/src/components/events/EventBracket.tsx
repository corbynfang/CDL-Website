import { useRef, useEffect, useState } from "react";
import type { BracketData } from "../../services/api";
import BracketSkeleton from "../loaders/BracketSkeleton";
import BracketTree from "./BracketTree";
import BracketControls from "./BracketControls";
import BracketCanvas from "./BracketCanvas";
import GroupStageView from "./GroupStageView";

interface Props {
	data: BracketData | null;
	loading: boolean;
	error: string | null;
}

export default function EventBracket({ data, loading, error }: Props) {
	const [zoom, setZoom] = useState(1.0);
	const [isFullscreen, setIsFullscreen] = useState(false);
	const [activeRound, setActiveRound] = useState<string | null>(null);
	const [userTab, setUserTab] = useState<"bracket" | "group_stage" | null>(
		null,
	);
	const containerRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const onChange = () => setIsFullscreen(!!document.fullscreenElement);
		document.addEventListener("fullscreenchange", onChange);
		return () => document.removeEventListener("fullscreenchange", onChange);
	}, []);

	function toggleFullscreen() {
		if (!document.fullscreenElement) containerRef.current?.requestFullscreen();
		else document.exitFullscreen();
	}

	if (loading) return <BracketSkeleton />;

	if (error) {
		return (
			<p className="text-center text-zinc-600 py-16 text-sm">
				Bracket data not available for this event.
			</p>
		);
	}

	if (!data || data.total_matches === 0) {
		return (
			<p className="text-center text-zinc-600 py-16 text-sm">
				No bracket matches have been played yet.
			</p>
		);
	}

	const isEWC = data.event_format === "ewc_group_stage_single_elim";

	if (isEWC) {
		return (
			<div
				ref={containerRef}
				className={isFullscreen ? "bg-[#09090b] p-6 h-full overflow-auto" : ""}
			>
				<BracketTree
					data={data}
					isFullscreen={isFullscreen}
					onFullscreen={toggleFullscreen}
				/>
			</div>
		);
	}

	const hasGroupStage = !!(
		data.group_stage && Object.keys(data.group_stage).length > 0
	);
	const activeTab: "bracket" | "group_stage" = userTab ?? "bracket";
	const bracketRounds = Object.keys(data.bracket);

	const tabPill =
		"text-[11px] uppercase tracking-widest px-4 py-2 border-b-2 transition-colors";
	const tabActive = "border-white text-white";
	const tabInactive = "border-transparent text-zinc-600 hover:text-zinc-400";

	return (
		<div
			ref={containerRef}
			className={`space-y-6 ${isFullscreen ? "bg-[#09090b] p-6 h-full overflow-auto" : ""}`}
		>
			{hasGroupStage && (
				<div className="flex gap-0 border-b border-[#1e1e1e]">
					<button
						onClick={() => setUserTab("bracket")}
						className={`${tabPill} ${activeTab === "bracket" ? tabActive : tabInactive}`}
					>
						Bracket
					</button>
					<button
						onClick={() => setUserTab("group_stage")}
						className={`${tabPill} ${activeTab === "group_stage" ? tabActive : tabInactive}`}
					>
						Group Stage
					</button>
				</div>
			)}

			{activeTab === "bracket" && (
				<>
					{bracketRounds.length > 1 && (
						<BracketControls
							rounds={bracketRounds}
							active={activeRound}
							onSelect={setActiveRound}
							zoom={zoom}
							onZoom={setZoom}
							isFullscreen={isFullscreen}
							onFullscreen={toggleFullscreen}
						/>
					)}
					<BracketCanvas
						data={data}
						activeRound={activeRound}
						zoom={zoom}
						flat={false}
					/>
				</>
			)}

			{activeTab === "group_stage" && data.group_stage && (
				<GroupStageView
					groupStage={data.group_stage}
					format={data.event_format ?? ""}
				/>
			)}
		</div>
	);
}
