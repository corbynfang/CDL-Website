import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import LiveEventStrip from "./LiveEventStrip";
import { liveMajor } from "../../test/fixtures/events";

vi.mock("../../utils/assets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

function wrap(ui: React.ReactElement) {
	return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe("LiveEventStrip", () => {
	it('renders the "Live Now" label', () => {
		wrap(<LiveEventStrip event={liveMajor} />);
		expect(screen.getByText("Live Now")).toBeInTheDocument();
	});

	it("renders the event name", () => {
		wrap(<LiveEventStrip event={liveMajor} />);
		expect(screen.getByText("CDL Major 3 2025")).toBeInTheDocument();
	});

	it("links to the correct event slug", () => {
		wrap(<LiveEventStrip event={liveMajor} />);
		const link = screen.getByRole("link");
		expect(link).toHaveAttribute("href", "/events/cdl-major-3-tournament-2025");
	});

	it("shows the location when provided", () => {
		wrap(<LiveEventStrip event={liveMajor} />);
		expect(screen.getByText(/Boca Raton, Florida/)).toBeInTheDocument();
	});

	it("does not crash when location is empty", () => {
		wrap(
			<LiveEventStrip event={{ ...liveMajor, location: "", country: "" }} />,
		);
		expect(screen.getByText("Live Now")).toBeInTheDocument();
	});

	it("renders a date range string in the link", () => {
		wrap(<LiveEventStrip event={liveMajor} />);
		// formatDateRange always includes "→" arrow suffix in this component
		const link = screen.getByRole("link");
		expect(link.textContent).toMatch(/→/);
	});
});
