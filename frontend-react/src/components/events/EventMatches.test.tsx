import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import EventMatches from "./EventMatches";
import {
	sampleMatches,
	winnersR1Match,
	noTypeMatch,
} from "../../test/fixtures/matches";

vi.mock("../../utils/assets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

function wrap(ui: React.ReactElement) {
	return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe("EventMatches", () => {
	it("shows skeleton loaders while loading", () => {
		const { container } = wrap(
			<EventMatches matches={null} loading={true} error={null} />,
		);
		expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
	});

	it("shows empty state when matches is null", () => {
		wrap(<EventMatches matches={null} loading={false} error={null} />);
		expect(screen.getByText("No matches recorded yet.")).toBeInTheDocument();
	});

	it("shows empty state when matches array is empty", () => {
		wrap(<EventMatches matches={[]} loading={false} error={null} />);
		expect(screen.getByText("No matches recorded yet.")).toBeInTheDocument();
	});

	it("shows error state when error is provided", () => {
		wrap(<EventMatches matches={null} loading={false} error="Failed" />);
		expect(screen.getByText(/could not load matches/i)).toBeInTheDocument();
	});

	it("renders match team names when data is provided", () => {
		wrap(<EventMatches matches={sampleMatches} loading={false} error={null} />);
		expect(screen.getAllByText("OpTic Texas").length).toBeGreaterThan(0);
		expect(screen.getAllByText("Atlanta FaZe").length).toBeGreaterThan(0);
	});

	it("renders match scores", () => {
		wrap(
			<EventMatches matches={[winnersR1Match]} loading={false} error={null} />,
		);
		expect(screen.getByText("3")).toBeInTheDocument();
		expect(screen.getByText("0")).toBeInTheDocument();
	});

	it("groups matches by match_type", () => {
		wrap(<EventMatches matches={sampleMatches} loading={false} error={null} />);
		expect(screen.getAllByText("winners r1").length).toBeGreaterThan(0);
		expect(screen.getAllByText("grand finals").length).toBeGreaterThan(0);
	});

	it('groups matches without match_type under "Other"', () => {
		wrap(<EventMatches matches={[noTypeMatch]} loading={false} error={null} />);
		expect(screen.getByText("Other")).toBeInTheDocument();
	});

	it("links each match to /matches/:id", () => {
		wrap(
			<EventMatches matches={[winnersR1Match]} loading={false} error={null} />,
		);
		const links = screen.getAllByRole("link");
		expect(links.some((l) => l.getAttribute("href") === "/matches/1")).toBe(
			true,
		);
	});
});
