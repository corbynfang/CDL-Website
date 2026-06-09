import { useRef, useState, useCallback, useEffect } from "react";
import type { BracketData } from "../../services/api";
import {
	computeBracketLayout,
	CARD_W,
	EWC_BAND,
	EWC_X_OFFSET,
} from "../../lib/bracketLayout";
import type { Connector } from "../../lib/bracketLayout";
import BracketMatchCard from "./BracketMatchCard";

interface Props {
	data: BracketData;
	isFullscreen?: boolean;
	onFullscreen?: () => void;
}

function ConnectorPath({ c }: { c: Connector }) {
	const midX = Math.round((c.fromX + c.toX) / 2);
	const d = `M ${c.fromX} ${c.fromY} H ${midX} V ${c.toY} H ${c.toX}`;
	return (
		<path
			d={d}
			fill="none"
			stroke={c.isLoser ? "#3f3f46" : "#71717a"}
			strokeWidth={c.isLoser ? 1 : 2}
			strokeDasharray={c.isLoser ? "4 3" : undefined}
			strokeLinecap="round"
			strokeLinejoin="round"
		/>
	);
}

function ConnectorLayer({
	connectors,
	canvasWidth,
	canvasHeight,
}: {
	connectors: Connector[];
	canvasWidth: number;
	canvasHeight: number;
}) {
	return (
		<svg
			className="absolute inset-0 pointer-events-none"
			width={canvasWidth}
			height={canvasHeight}
			style={{ overflow: "visible" }}
		>
			{connectors.map((c, i) => (
				<ConnectorPath key={`${c.fromId}-${c.toId}-${i}`} c={c} />
			))}
		</svg>
	);
}

function ColLabel({ label, x }: { label: string; x: number }) {
	return (
		<div
			className="absolute top-0 flex flex-col items-center"
			style={{ left: x, width: CARD_W }}
		>
			<span className="text-[11px] uppercase tracking-[0.12em] text-zinc-400 whitespace-nowrap font-medium">
				{label}
			</span>
			<div className="mt-1.5 w-full h-px bg-[#2a2a2a]" />
		</div>
	);
}

interface Transform {
	x: number;
	y: number;
	scale: number;
}

const ZOOM_STEPS = [0.4, 0.5, 0.6, 0.75, 0.9, 1.0];
const ZOOM_DEFAULT = 0.9;

function clamp(v: number, lo: number, hi: number) {
	return Math.max(lo, Math.min(hi, v));
}

function ZoomControls({
	scale,
	onZoom,
	isFullscreen,
	onFullscreen,
}: {
	scale: number;
	onZoom: (delta: number) => void;
	isFullscreen?: boolean;
	onFullscreen?: () => void;
}) {
	return (
		<div className="absolute top-3 right-3 flex items-center gap-1 z-10">
			<button
				onClick={() => onZoom(-1)}
				className="w-7 h-7 rounded border border-[#2e2e2e] bg-[#111] text-zinc-400 hover:text-white hover:border-zinc-600 text-sm flex items-center justify-center"
			>
				−
			</button>
			<span className="text-[10px] text-zinc-600 w-10 text-center tabular-nums">
				{Math.round(scale * 100)}%
			</span>
			<button
				onClick={() => onZoom(+1)}
				className="w-7 h-7 rounded border border-[#2e2e2e] bg-[#111] text-zinc-400 hover:text-white hover:border-zinc-600 text-sm flex items-center justify-center"
			>
				+
			</button>
			{onFullscreen && (
				<button
					onClick={onFullscreen}
					className="ml-1 w-7 h-7 rounded border border-[#2e2e2e] bg-[#111] text-zinc-400 hover:text-white hover:border-zinc-600 flex items-center justify-center"
					title={isFullscreen ? "Exit fullscreen" : "Fullscreen"}
				>
					{isFullscreen ? (
						<svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
							<path d="M1 4V1h3v1H2v2H1zm7-3h3v3h-1V2H8V1zM1 8h1v2h2v1H1V8zm9 2v-2h1v3H8v-1h2z" />
						</svg>
					) : (
						<svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
							<path d="M0 0h4v1H1v3H0V0zm8 0h4v4h-1V1H8V0zM0 8h1v3h3v1H0V8zm11 3H8v1h4V8h-1v3z" />
						</svg>
					)}
				</button>
			)}
		</div>
	);
}

const EWC_GROUPS = ["A", "B", "C", "D"] as const;

export default function BracketTree({
	data,
	isFullscreen,
	onFullscreen,
}: Props) {
	const layout = computeBracketLayout(data);
	const isEWC = data.event_format === "ewc_group_stage_single_elim";
	const LABEL_H = 28; // height reserved for column labels above the canvas

	const [transform, setTransform] = useState<Transform>({
		x: 0,
		y: 0,
		scale: ZOOM_DEFAULT,
	});
	const dragging = useRef(false);
	const lastPos = useRef({ x: 0, y: 0 });
	const viewportRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		setTransform({ x: 0, y: 0, scale: ZOOM_DEFAULT });
	}, [data.tournament_id]);

	const stepZoom = useCallback((delta: number) => {
		setTransform((t) => {
			const idx = ZOOM_STEPS.findIndex((s) => s >= t.scale - 0.01);
			const next = ZOOM_STEPS[clamp(idx + delta, 0, ZOOM_STEPS.length - 1)];
			return { ...t, scale: next };
		});
	}, []);

	const onWheel = useCallback(
		(e: React.WheelEvent) => {
			e.preventDefault();
			stepZoom(e.deltaY < 0 ? +1 : -1);
		},
		[stepZoom],
	);

	const onMouseDown = useCallback((e: React.MouseEvent) => {
		if (e.button !== 0) return;
		dragging.current = true;
		lastPos.current = { x: e.clientX, y: e.clientY };
	}, []);

	const onMouseMove = useCallback((e: React.MouseEvent) => {
		if (!dragging.current) return;
		const dx = e.clientX - lastPos.current.x;
		const dy = e.clientY - lastPos.current.y;
		lastPos.current = { x: e.clientX, y: e.clientY };
		setTransform((t) => ({ ...t, x: t.x + dx, y: t.y + dy }));
	}, []);

	const onMouseUp = useCallback(() => {
		dragging.current = false;
	}, []);

	if (layout.nodes.length === 0) {
		return (
			<p className="text-center text-zinc-600 py-16 text-sm">
				Bracket data not available yet.
			</p>
		);
	}

	const innerTop = LABEL_H + 8; // canvas content starts below labels

	return (
		<div
			ref={viewportRef}
			className="relative overflow-hidden rounded border border-[#1a1a1a] bg-[#09090b] select-none"
			style={{
				height: isFullscreen ? "100%" : "70vh",
				cursor: dragging.current ? "grabbing" : "grab",
			}}
			onWheel={onWheel}
			onMouseDown={onMouseDown}
			onMouseMove={onMouseMove}
			onMouseUp={onMouseUp}
			onMouseLeave={onMouseUp}
		>
			<ZoomControls
				scale={transform.scale}
				onZoom={stepZoom}
				isFullscreen={isFullscreen}
				onFullscreen={onFullscreen}
			/>

			{/* Transformed inner canvas */}
			<div
				style={{
					transform: `translate(${transform.x + 24}px, ${transform.y + 24}px) scale(${transform.scale})`,
					transformOrigin: "0 0",
					position: "absolute",
					width: layout.canvasWidth,
					height: layout.canvasHeight + innerTop,
					willChange: "transform",
				}}
			>
				{/* Column labels */}
				<div className="relative" style={{ height: LABEL_H }}>
					{layout.colLabels.map((cl) => (
						<ColLabel key={cl.col} label={cl.label} x={cl.x} />
					))}
				</div>

				{/* Cards + SVG connectors */}
				<div className="relative" style={{ height: layout.canvasHeight }}>
					{/* EWC group dividers — thin lines separating each group band */}
					{isEWC &&
						[1, 2, 3].map((i) => (
							<div
								key={i}
								className="absolute bg-[#1e1e1e]"
								style={{
									left: 0,
									top: i * EWC_BAND,
									width: layout.canvasWidth,
									height: 1,
								}}
							/>
						))}

					{/* EWC group labels — rotated text in the left margin */}
					{isEWC &&
						EWC_GROUPS.map((g, i) => (
							<div
								key={g}
								className="absolute flex items-center justify-center"
								style={{
									left: 0,
									top: i * EWC_BAND,
									width: EWC_X_OFFSET - 8,
									height: EWC_BAND,
								}}
							>
								<span
									className="text-[10px] uppercase tracking-[0.2em] text-zinc-500 select-none"
									style={{
										writingMode: "vertical-rl",
										transform: "rotate(180deg)",
									}}
								>
									Group {g}
								</span>
							</div>
						))}

					<ConnectorLayer
						connectors={layout.connectors}
						canvasWidth={layout.canvasWidth}
						canvasHeight={layout.canvasHeight}
					/>

					{layout.nodes.map((n) => (
						<div
							key={n.match.id}
							className="absolute"
							style={{ left: n.x, top: isEWC ? n.y - 20 : n.y, width: CARD_W }}
						>
							{isEWC && (
								<div className="text-[9px] uppercase tracking-[0.15em] text-zinc-500 mb-1.5 select-none whitespace-nowrap">
									{n.roundLabel}
								</div>
							)}
							<BracketMatchCard match={n.match} />
						</div>
					))}
				</div>
			</div>
		</div>
	);
}
