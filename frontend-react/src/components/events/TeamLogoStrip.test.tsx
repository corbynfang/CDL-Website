import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import TeamLogoStrip from "./TeamLogoStrip";
import {
	sampleTeams,
	largeTeamList,
	unknownTeam,
} from "../../test/fixtures/teams";

vi.mock("../../utils/logoAssets", () => ({
	getTeamLogo: vi.fn().mockReturnValue(null),
	getPlayerAvatar: vi.fn().mockReturnValue("/placeholder.png"),
}));

import { getTeamLogo } from "../../utils/logoAssets";

describe("TeamLogoStrip", () => {
	beforeEach(() => vi.mocked(getTeamLogo).mockReturnValue(""));

	it("renders without crashing with an empty array", () => {
		const { container } = render(<TeamLogoStrip teams={[]} />);
		expect(container.firstChild).toBeInTheDocument();
	});

	it("shows abbreviations when no logo is available", () => {
		render(<TeamLogoStrip teams={sampleTeams} />);
		expect(screen.getByText("OT")).toBeInTheDocument();
		expect(screen.getByText("AT")).toBeInTheDocument();
		expect(screen.getByText("BO")).toBeInTheDocument();
	});

	it("shows team logos when getTeamLogo returns a URL", () => {
		vi.mocked(getTeamLogo).mockReturnValue("/logos/optic.png");
		render(<TeamLogoStrip teams={[sampleTeams[0]]} />);
		const img = screen.getByRole("img", { name: "OpTic Texas" });
		expect(img).toBeInTheDocument();
		expect(img).toHaveAttribute("src", "/logos/optic.png");
	});

	it("shows overflow count when teams exceed max", () => {
		render(<TeamLogoStrip teams={largeTeamList} max={12} />);
		expect(screen.getByText("+2")).toBeInTheDocument();
	});

	it("does not show overflow count when teams equal max", () => {
		render(<TeamLogoStrip teams={sampleTeams} max={3} />);
		expect(screen.queryByText(/^\+/)).not.toBeInTheDocument();
	});

	it('falls back to "?" for a team with no abbreviation', () => {
		const noAbbr = {
			...unknownTeam,
			abbreviation: undefined as unknown as string,
		};
		render(<TeamLogoStrip teams={[noAbbr]} />);
		expect(screen.getByText("?")).toBeInTheDocument();
	});
});
