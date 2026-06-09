import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import BracketCanvas from "./BracketCanvas";
import { sampleBracketData } from "../../test/fixtures/matches";
import type { BracketData } from "../../services/api";

vi.mock("../../utils/assets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

function wrap(ui: React.ReactElement) {
	return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe("BracketCanvas", () => {
	it("renders without crashing", () => {
		const { container } = wrap(
			<BracketCanvas data={sampleBracketData} activeRound={null} />,
		);
		expect(container.firstChild).toBeInTheDocument();
	});

	it('renders "Winners Bracket" section label when winners rounds exist', () => {
		wrap(<BracketCanvas data={sampleBracketData} activeRound={null} />);
		expect(screen.getByText("Winners Bracket")).toBeInTheDocument();
	});

	it('renders "Grand Finals" section label when grand finals round exists', () => {
		wrap(<BracketCanvas data={sampleBracketData} activeRound={null} />);
		// Both the section label and the round column header say "Grand Finals"
		expect(screen.getAllByText("Grand Finals").length).toBeGreaterThanOrEqual(
			2,
		);
	});

	it("renders round column header for winners_r1", () => {
		wrap(<BracketCanvas data={sampleBracketData} activeRound={null} />);
		expect(screen.getByText("Winners Round 1")).toBeInTheDocument();
	});

	it("renders match cards for each round", () => {
		wrap(<BracketCanvas data={sampleBracketData} activeRound={null} />);
		// Both matches in sampleBracketData should render — OpTic Texas appears in both
		const otxEls = screen.getAllByText("OpTic Texas");
		expect(otxEls.length).toBeGreaterThanOrEqual(2);
	});

	it("filters to only the active round when activeRound is set", () => {
		wrap(<BracketCanvas data={sampleBracketData} activeRound="winners_r1" />);
		expect(screen.getByText("Winners Round 1")).toBeInTheDocument();
		expect(screen.queryByText("Grand Finals")).not.toBeInTheDocument();
	});

	it("does not render Elimination Bracket section when no elim rounds exist", () => {
		wrap(<BracketCanvas data={sampleBracketData} activeRound={null} />);
		expect(screen.queryByText("Elimination Bracket")).not.toBeInTheDocument();
	});

	it("renders Elimination Bracket section when elim rounds exist", () => {
		const withElim: BracketData = {
			...sampleBracketData,
			bracket: {
				...sampleBracketData.bracket,
				elim_r1: [sampleBracketData.bracket.winners_r1[0]],
			},
		};
		wrap(<BracketCanvas data={withElim} activeRound={null} />);
		expect(screen.getByText("Elimination Bracket")).toBeInTheDocument();
	});

	it("renders an empty bracket without crashing", () => {
		const empty: BracketData = {
			tournament_id: 1,
			tournament_name: "Test",
			total_matches: 0,
			bracket: {},
		};
		const { container } = wrap(
			<BracketCanvas data={empty} activeRound={null} />,
		);
		expect(container.firstChild).toBeInTheDocument();
	});

	it("flat mode renders round columns without section labels", () => {
		const ewcData: BracketData = {
			tournament_id: 53,
			tournament_name: "EWC 2025",
			total_matches: 2,
			event_format: "ewc_group_stage_single_elim",
			bracket: {
				quarterfinal: [sampleBracketData.bracket.winners_r1[0]],
				semifinal: [sampleBracketData.bracket.grand_finals[0]],
			},
		};
		wrap(<BracketCanvas data={ewcData} activeRound={null} flat={true} />);
		expect(screen.getByText("Quarterfinal")).toBeInTheDocument();
		expect(screen.getByText("Semifinal")).toBeInTheDocument();
		expect(screen.queryByText("Winners Bracket")).not.toBeInTheDocument();
		expect(screen.queryByText("Elimination Bracket")).not.toBeInTheDocument();
	});

	it("flat mode shows elim_r4 and elim_r5 labels for Cold War format", () => {
		const cwData: BracketData = {
			tournament_id: 12,
			tournament_name: "CW Stage 1",
			total_matches: 2,
			event_format: "cold_war_stage_double_elim",
			bracket: {
				elim_r4: [sampleBracketData.bracket.winners_r1[0]],
				elim_r5: [sampleBracketData.bracket.grand_finals[0]],
			},
		};
		wrap(<BracketCanvas data={cwData} activeRound={null} flat={true} />);
		expect(screen.getByText("Elimination Round 4")).toBeInTheDocument();
		expect(screen.getByText("Elimination Round 5")).toBeInTheDocument();
	});
});
