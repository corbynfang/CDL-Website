import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import BracketMatchCard from "./BracketMatchCard";
import {
	bracketMatchComplete,
	bracketMatchNoWinner,
} from "../../test/fixtures/matches";

vi.mock("../../utils/logoAssets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

function wrap(ui: React.ReactElement) {
	return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe("BracketMatchCard", () => {
	it("renders both team names", () => {
		wrap(<BracketMatchCard match={bracketMatchComplete} />);
		expect(screen.getByText("OpTic Texas")).toBeInTheDocument();
		expect(screen.getByText("Atlanta FaZe")).toBeInTheDocument();
	});

	it("renders both team scores", () => {
		wrap(<BracketMatchCard match={bracketMatchComplete} />);
		expect(screen.getByText("3")).toBeInTheDocument();
		expect(screen.getByText("1")).toBeInTheDocument();
	});

	it("links to the correct match detail page", () => {
		wrap(<BracketMatchCard match={bracketMatchComplete} />);
		const link = screen.getByRole("link");
		expect(link).toHaveAttribute("href", "/matches/10");
	});

	it("renders team abbreviation initials when no logo is available", () => {
		wrap(<BracketMatchCard match={bracketMatchComplete} />);
		// With getTeamLogo mocked to null and logo prop empty, shows abbr initials
		expect(screen.getByText("OT")).toBeInTheDocument();
		expect(screen.getByText("AT")).toBeInTheDocument();
	});

	it("uses team logo from getTeamLogo when available", async () => {
		const { getTeamLogo } = await import("../../utils/logoAssets");
		vi.mocked(getTeamLogo).mockReturnValue("/logos/optic.png");
		wrap(<BracketMatchCard match={bracketMatchComplete} />);
		const imgs = screen.getAllByRole("img");
		expect(
			imgs.some((img) => img.getAttribute("src") === "/logos/optic.png"),
		).toBe(true);
		vi.mocked(getTeamLogo).mockReturnValue("");
	});

	it("renders a completed match without crashing", () => {
		wrap(<BracketMatchCard match={bracketMatchComplete} />);
		expect(screen.getByRole("link")).toBeInTheDocument();
	});

	it("renders a match with no winner without crashing", () => {
		wrap(<BracketMatchCard match={bracketMatchNoWinner} />);
		expect(screen.getByRole("link")).toBeInTheDocument();
	});

	it("renders em-dash scores for an unplayed match", () => {
		wrap(<BracketMatchCard match={bracketMatchNoWinner} />);
		const dashes = screen.getAllByText("–");
		expect(dashes).toHaveLength(2);
	});
});
