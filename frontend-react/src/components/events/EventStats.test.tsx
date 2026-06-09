import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import EventStats from "./EventStats";
import {
	sampleStats,
	highKDStat,
	lowKDStat,
	noPlayerStat,
} from "../../test/fixtures/stats";

vi.mock("../../utils/assets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

function wrap(ui: React.ReactElement) {
	return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe("EventStats", () => {
	it("shows skeleton loaders while loading", () => {
		const { container } = wrap(<EventStats stats={null} loading={true} />);
		expect(container.querySelector(".animate-pulse")).toBeInTheDocument();
	});

	it("shows empty state when stats is null", () => {
		wrap(<EventStats stats={null} loading={false} />);
		expect(
			screen.getByText("No stats available for this event."),
		).toBeInTheDocument();
	});

	it("shows empty state when stats array is empty", () => {
		wrap(<EventStats stats={[]} loading={false} />);
		expect(
			screen.getByText("No stats available for this event."),
		).toBeInTheDocument();
	});

	it("renders player gamertags", () => {
		wrap(<EventStats stats={sampleStats} loading={false} />);
		expect(screen.getByText("Scump")).toBeInTheDocument();
		expect(screen.getByText("Simp")).toBeInTheDocument();
	});

	it("renders K/D ratios formatted to 2 decimal places", () => {
		wrap(<EventStats stats={[highKDStat]} loading={false} />);
		expect(screen.getByText("1.45")).toBeInTheDocument();
	});

	it("renders kill counts", () => {
		wrap(<EventStats stats={[highKDStat]} loading={false} />);
		expect(screen.getByText("120")).toBeInTheDocument();
	});

	it("renders team abbreviations", () => {
		wrap(<EventStats stats={[lowKDStat]} loading={false} />);
		expect(screen.getByText("ATL")).toBeInTheDocument();
	});

	it("links each player row to /players/:id", () => {
		wrap(<EventStats stats={[highKDStat]} loading={false} />);
		const links = screen.getAllByRole("link");
		expect(links.some((l) => l.getAttribute("href") === "/players/1")).toBe(
			true,
		);
	});

	it('falls back to "Player #N" when player is undefined', () => {
		wrap(<EventStats stats={[noPlayerStat]} loading={false} />);
		expect(screen.getByText("Player #99")).toBeInTheDocument();
	});

	it("renders table header columns", () => {
		wrap(<EventStats stats={sampleStats} loading={false} />);
		expect(screen.getByText("Player")).toBeInTheDocument();
		expect(screen.getByText("K")).toBeInTheDocument();
		expect(screen.getByText("K/D")).toBeInTheDocument();
	});

	it("renders row index numbers starting from 1", () => {
		wrap(<EventStats stats={sampleStats} loading={false} />);
		expect(screen.getByText("1")).toBeInTheDocument();
		expect(screen.getByText("2")).toBeInTheDocument();
	});
});
