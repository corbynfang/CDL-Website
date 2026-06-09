import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import EventOverview from "./EventOverview";
import { completedMajor, upcomingQualifier } from "../../test/fixtures/events";

vi.mock("../../utils/logoAssets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

describe("EventOverview", () => {
	it("shows the team count when teamCount > 0", () => {
		render(<EventOverview event={completedMajor} teamCount={12} />);
		expect(screen.getByText("12")).toBeInTheDocument();
	});

	it('shows "—" for team count when teamCount is 0', () => {
		render(<EventOverview event={completedMajor} teamCount={0} />);
		// The Teams stat card shows "—" when count is 0
		const labels = screen.getAllByText("—");
		expect(labels.length).toBeGreaterThan(0);
	});

	it("shows formatted prize pool", () => {
		render(<EventOverview event={completedMajor} teamCount={12} />);
		expect(screen.getByText("$375K")).toBeInTheDocument();
	});

	it('shows "TBA" when prize pool is null', () => {
		render(<EventOverview event={upcomingQualifier} teamCount={0} />);
		expect(screen.getByText("TBA")).toBeInTheDocument();
	});

	it('shows "LAN" for LAN events in the Type stat card', () => {
		render(<EventOverview event={completedMajor} teamCount={12} />);
		expect(screen.getByText("LAN")).toBeInTheDocument();
	});

	it('shows "Online" for non-LAN events in the Type stat card', () => {
		render(<EventOverview event={upcomingQualifier} teamCount={0} />);
		expect(screen.getByText("Online")).toBeInTheDocument();
	});

	it("shows tournament format when available", () => {
		render(<EventOverview event={completedMajor} teamCount={12} />);
		expect(screen.getByText("Double Elimination")).toBeInTheDocument();
	});

	it('shows "—" for format when not available', () => {
		render(<EventOverview event={upcomingQualifier} teamCount={0} />);
		const dashes = screen.getAllByText("—");
		expect(dashes.length).toBeGreaterThan(0);
	});

	it("shows location in the details section", () => {
		render(<EventOverview event={completedMajor} teamCount={12} />);
		expect(screen.getByText(/Allen, Texas/)).toBeInTheDocument();
	});

	it("shows the season name in details when season is attached", () => {
		render(<EventOverview event={completedMajor} teamCount={12} />);
		expect(screen.getByText("Black Ops 6 2025")).toBeInTheDocument();
	});

	it("shows the Event Page link when source_event_url is provided", () => {
		render(<EventOverview event={completedMajor} teamCount={12} />);
		const link = screen.getByRole("link", { name: /event page/i });
		expect(link).toHaveAttribute(
			"href",
			"https://liquipedia.net/callofduty/test",
		);
	});

	it("does not render external links section when source_event_url is empty", () => {
		render(<EventOverview event={upcomingQualifier} teamCount={0} />);
		expect(screen.queryByText(/external links/i)).not.toBeInTheDocument();
		expect(screen.queryByRole("link")).not.toBeInTheDocument();
	});

	it("does not crash with no location or season", () => {
		const minimal = { ...upcomingQualifier, location: "", season: undefined };
		render(<EventOverview event={minimal} teamCount={0} />);
		expect(screen.getByText("TBA")).toBeInTheDocument();
	});
});
