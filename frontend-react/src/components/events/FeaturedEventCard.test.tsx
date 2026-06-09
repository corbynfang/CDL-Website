import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import FeaturedEventCard from "./FeaturedEventCard";
import {
	completedMajor,
	upcomingQualifier,
	liveMajor,
} from "../../test/fixtures/events";
import { sampleTeams } from "../../test/fixtures/teams";

vi.mock("../../utils/logoAssets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

function wrap(ui: React.ReactElement) {
	return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe("FeaturedEventCard", () => {
	it("renders the event name", () => {
		wrap(<FeaturedEventCard event={completedMajor} />);
		expect(screen.getByText("CDL Major 2 2025")).toBeInTheDocument();
	});

	it("links to the correct event slug", () => {
		wrap(<FeaturedEventCard event={completedMajor} />);
		const link = screen.getByRole("link");
		expect(link).toHaveAttribute("href", "/events/cdl-major-2-tournament-2025");
	});

	it("shows the LAN badge for a LAN event", () => {
		wrap(<FeaturedEventCard event={completedMajor} />);
		expect(screen.getByText("LAN")).toBeInTheDocument();
	});

	it("does not show LAN badge for an online event", () => {
		wrap(<FeaturedEventCard event={upcomingQualifier} />);
		expect(screen.queryByText("LAN")).not.toBeInTheDocument();
	});

	it("shows the prize pool for events with a prize pool", () => {
		wrap(<FeaturedEventCard event={completedMajor} />);
		expect(screen.getByText("$375K")).toBeInTheDocument();
	});

	it("does not show prize pool when prize_pool is null", () => {
		wrap(<FeaturedEventCard event={upcomingQualifier} />);
		expect(screen.queryByText(/\$/)).not.toBeInTheDocument();
	});

	it("shows the location with country flag for a completed major", () => {
		wrap(<FeaturedEventCard event={completedMajor} />);
		expect(screen.getByText(/Allen, Texas/)).toBeInTheDocument();
	});

	it("shows completed status badge for a past event", () => {
		wrap(<FeaturedEventCard event={completedMajor} />);
		expect(screen.getByText("completed")).toBeInTheDocument();
	});

	it("shows upcoming status badge for a future event", () => {
		wrap(<FeaturedEventCard event={upcomingQualifier} />);
		expect(screen.getByText("upcoming")).toBeInTheDocument();
	});

	it("shows live status badge for a currently live event", () => {
		wrap(<FeaturedEventCard event={liveMajor} />);
		expect(screen.getByText("live")).toBeInTheDocument();
	});

	it("renders with no teams without crashing", () => {
		wrap(<FeaturedEventCard event={completedMajor} teams={[]} />);
		expect(screen.getByText("CDL Major 2 2025")).toBeInTheDocument();
	});

	it("renders TeamLogoStrip when teams are provided", () => {
		wrap(<FeaturedEventCard event={completedMajor} teams={sampleTeams} />);
		// With getTeamLogo mocked to null, abbreviations render
		expect(screen.getByText("OT")).toBeInTheDocument();
	});

	it("shows game code from season when season is attached", () => {
		wrap(<FeaturedEventCard event={completedMajor} />);
		expect(screen.getByText("BO6")).toBeInTheDocument();
	});

	it("shows tournament type when no season is attached", () => {
		wrap(
			<FeaturedEventCard event={{ ...completedMajor, season: undefined }} />,
		);
		expect(screen.getByText("major tournament")).toBeInTheDocument();
	});
});
