import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import EventHero from "./EventHero";
import {
	completedMajor,
	upcomingQualifier,
	liveMajor,
	championshipNoSeason,
} from "../../test/fixtures/events";
import { sampleTeams } from "../../test/fixtures/teams";

vi.mock("../../utils/assets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

describe("EventHero", () => {
	it("renders the event name as h1", () => {
		render(<EventHero event={completedMajor} teamCount={12} teams={[]} />);
		expect(
			screen.getByRole("heading", { level: 1, name: "CDL Major 2 2025" }),
		).toBeInTheDocument();
	});

	it("shows the game code from season", () => {
		render(<EventHero event={completedMajor} teamCount={12} teams={[]} />);
		expect(screen.getByText("BO6")).toBeInTheDocument();
	});

	it('shows "—" when no season is attached', () => {
		render(<EventHero event={championshipNoSeason} teamCount={8} teams={[]} />);
		expect(screen.getByText("—")).toBeInTheDocument();
	});

	it("shows the tournament type", () => {
		render(<EventHero event={completedMajor} teamCount={12} teams={[]} />);
		expect(screen.getByText("major tournament")).toBeInTheDocument();
	});

	it('shows "LAN" label for LAN events', () => {
		render(<EventHero event={completedMajor} teamCount={12} teams={[]} />);
		expect(screen.getByText("LAN")).toBeInTheDocument();
	});

	it("does not show LAN label for online events", () => {
		render(<EventHero event={upcomingQualifier} teamCount={0} teams={[]} />);
		expect(screen.queryByText("LAN")).not.toBeInTheDocument();
	});

	it("shows the completed status badge for a past event", () => {
		render(<EventHero event={completedMajor} teamCount={12} teams={[]} />);
		expect(screen.getByText("completed")).toBeInTheDocument();
	});

	it("shows the upcoming status badge for a future event", () => {
		render(<EventHero event={upcomingQualifier} teamCount={0} teams={[]} />);
		expect(screen.getByText("upcoming")).toBeInTheDocument();
	});

	it("shows the live status badge for a live event", () => {
		render(<EventHero event={liveMajor} teamCount={12} teams={[]} />);
		expect(screen.getByText("live")).toBeInTheDocument();
	});

	it("shows the prize pool", () => {
		render(<EventHero event={completedMajor} teamCount={12} teams={[]} />);
		expect(screen.getByText("$375K")).toBeInTheDocument();
	});

	it("does not show prize pool when prize_pool is null", () => {
		render(<EventHero event={upcomingQualifier} teamCount={0} teams={[]} />);
		expect(screen.queryByText(/\$/)).not.toBeInTheDocument();
	});

	it("shows location with country flag", () => {
		render(<EventHero event={completedMajor} teamCount={12} teams={[]} />);
		expect(screen.getByText(/Allen, Texas/)).toBeInTheDocument();
	});

	it("shows team count when teamCount > 0", () => {
		render(<EventHero event={completedMajor} teamCount={12} teams={[]} />);
		expect(screen.getByText("12 teams")).toBeInTheDocument();
	});

	it("does not show team count when teamCount is 0", () => {
		render(<EventHero event={completedMajor} teamCount={0} teams={[]} />);
		expect(screen.queryByText(/teams/)).not.toBeInTheDocument();
	});

	it("renders TeamLogoStrip when teams are provided", () => {
		render(
			<EventHero event={completedMajor} teamCount={3} teams={sampleTeams} />,
		);
		// With mocked getTeamLogo returning null, abbreviations render
		expect(screen.getByText("OT")).toBeInTheDocument();
	});

	it("does not render TeamLogoStrip when teams array is empty", () => {
		render(<EventHero event={completedMajor} teamCount={0} teams={[]} />);
		// No abbreviation divs from TeamLogoStrip
		expect(screen.queryByText("OT")).not.toBeInTheDocument();
	});
});
