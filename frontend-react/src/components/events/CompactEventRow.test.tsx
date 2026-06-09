import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import CompactEventRow from "./CompactEventRow";
import {
	completedMajor,
	upcomingQualifier,
	liveMajor,
} from "../../test/fixtures/events";

vi.mock("../../utils/assets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

function wrap(ui: React.ReactElement) {
	return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe("CompactEventRow", () => {
	it("renders the event name", () => {
		wrap(<CompactEventRow event={completedMajor} />);
		expect(screen.getByText("CDL Major 2 2025")).toBeInTheDocument();
	});

	it("links to the correct event slug", () => {
		wrap(<CompactEventRow event={completedMajor} />);
		const link = screen.getByRole("link");
		expect(link).toHaveAttribute("href", "/events/cdl-major-2-tournament-2025");
	});

	it("shows the prize pool when present", () => {
		wrap(<CompactEventRow event={completedMajor} />);
		expect(screen.getByText("$375K")).toBeInTheDocument();
	});

	it('shows "—" dash when prize pool is null', () => {
		wrap(<CompactEventRow event={upcomingQualifier} />);
		expect(screen.getByText("—")).toBeInTheDocument();
	});

	it("shows location when provided", () => {
		wrap(<CompactEventRow event={completedMajor} />);
		expect(screen.getByText(/Allen, Texas/)).toBeInTheDocument();
	});

	it("does not crash when location is empty", () => {
		wrap(<CompactEventRow event={upcomingQualifier} />);
		expect(screen.getByText("CDL Qualifier 1 2025")).toBeInTheDocument();
	});

	it("renders a status dot element", () => {
		const { container } = wrap(<CompactEventRow event={completedMajor} />);
		// Status dot is a span with rounded-full
		const dot = container.querySelector("span.rounded-full");
		expect(dot).toBeInTheDocument();
	});

	it("live event renders correctly", () => {
		wrap(<CompactEventRow event={liveMajor} />);
		expect(screen.getByText("CDL Major 3 2025")).toBeInTheDocument();
	});
});
