import { formatRound, sortedRounds } from "../../utils/eventUtils";

const ZOOM_LEVELS = [0.6, 0.75, 0.9, 1.0, 1.25, 1.5] as const;
type ZoomLevel = (typeof ZOOM_LEVELS)[number];

interface Props {
	rounds: string[];
	active: string | null;
	onSelect: (round: string | null) => void;
	zoom?: number;
	onZoom?: (z: number) => void;
	isFullscreen?: boolean;
	onFullscreen?: () => void;
}

export default function BracketControls({
	rounds,
	active,
	onSelect,
	zoom = 1.0,
	onZoom,
	isFullscreen = false,
	onFullscreen,
}: Props) {
	const sorted = sortedRounds(rounds);
	const zoomIdx = ZOOM_LEVELS.findIndex((z) => z === (zoom as ZoomLevel));
	const canZoomIn = zoomIdx < ZOOM_LEVELS.length - 1;
	const canZoomOut = zoomIdx > 0;

	const pill =
		"text-[10px] uppercase tracking-widest px-3 py-1.5 border transition-colors";
	const active_cls = "border-white text-white";
	const inactive_cls =
		"border-[#1a1a1a] text-zinc-600 hover:text-zinc-400 hover:border-[#2a2a2a]";
	const ctrl_cls =
		"border-[#1a1a1a] text-zinc-600 hover:text-zinc-300 hover:border-[#2a2a2a] disabled:opacity-30 disabled:cursor-not-allowed";

	return (
		<div className="flex flex-wrap items-center gap-2 justify-between">
			{/* Round filter */}
			<div className="flex flex-wrap gap-2">
				<button
					onClick={() => onSelect(null)}
					className={`${pill} ${active === null ? active_cls : inactive_cls}`}
				>
					All Rounds
				</button>
				{sorted.map((r) => (
					<button
						key={r}
						onClick={() => onSelect(r)}
						className={`${pill} ${active === r ? active_cls : inactive_cls}`}
					>
						{formatRound(r)}
					</button>
				))}
			</div>

			{/* Zoom + fullscreen */}
			{(onZoom || onFullscreen) && (
				<div className="flex items-center gap-1">
					<button
						onClick={() => canZoomOut && onZoom?.(ZOOM_LEVELS[zoomIdx - 1])}
						disabled={!canZoomOut || !onZoom}
						className={`${pill} ${ctrl_cls} px-2.5`}
						aria-label="Zoom out"
					>
						−
					</button>
					<button
						onClick={() => onZoom?.(1.0)}
						className={`${pill} ${ctrl_cls} tabular-nums min-w-[46px] text-center`}
						aria-label="Reset zoom"
					>
						{Math.round(zoom * 100)}%
					</button>
					<button
						onClick={() => canZoomIn && onZoom?.(ZOOM_LEVELS[zoomIdx + 1])}
						disabled={!canZoomIn || !onZoom}
						className={`${pill} ${ctrl_cls} px-2.5`}
						aria-label="Zoom in"
					>
						+
					</button>
					{onFullscreen && (
						<>
							<div className="w-px h-4 bg-[#1e1e1e] mx-1" />
							<button
								onClick={onFullscreen}
								className={`${pill} ${ctrl_cls} px-2.5`}
								aria-label={isFullscreen ? "Exit fullscreen" : "Fullscreen"}
							>
								{isFullscreen ? <CollapseIcon /> : <ExpandIcon />}
							</button>
						</>
					)}
				</div>
			)}
		</div>
	);
}

function ExpandIcon() {
	return (
		<svg
			width="11"
			height="11"
			viewBox="0 0 11 11"
			fill="none"
			stroke="currentColor"
			strokeWidth="1.5"
			strokeLinecap="round"
			aria-hidden="true"
		>
			<path d="M1 3.5V1h2.5M7.5 1H10v2.5M10 7.5V10H7.5M3.5 10H1V7.5" />
		</svg>
	);
}

function CollapseIcon() {
	return (
		<svg
			width="11"
			height="11"
			viewBox="0 0 11 11"
			fill="none"
			stroke="currentColor"
			strokeWidth="1.5"
			strokeLinecap="round"
			aria-hidden="true"
		>
			<path d="M3.5 1v2.5H1M10 3.5H7.5V1M7.5 10V7.5H10M1 7.5h2.5V10" />
		</svg>
	);
}
