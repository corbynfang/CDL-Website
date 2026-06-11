import { vi, describe, it, expect, beforeEach } from "vitest";

const mockGet = vi.hoisted(() => vi.fn());
const mockPost = vi.hoisted(() => vi.fn());

vi.mock("axios", () => ({
  default: {
    create: () => ({
      get: mockGet,
      post: mockPost,
      interceptors: {
        request: { use: vi.fn() },
        response: { use: vi.fn() },
      },
    }),
  },
}));

import { teamApi, playerApi, transfersApi } from "./api";

beforeEach(() => {
  mockGet.mockReset();
  mockPost.mockReset();
});

describe("teamApi", () => {
  it("getTeams returns the data array from the response", async () => {
    const teams = [{ id: 1, name: "Atlanta FaZe" }];
    mockGet.mockResolvedValue({ data: teams });

    const result = await teamApi.getTeams();
    expect(result).toEqual(teams);
    expect(mockGet).toHaveBeenCalledWith("/teams");
  });

  it("getTeam calls the correct endpoint with the given id", async () => {
    const team = { id: 3, name: "Boston Breach" };
    mockGet.mockResolvedValue({ data: team });

    const result = await teamApi.getTeam(3);
    expect(result).toEqual(team);
    expect(mockGet).toHaveBeenCalledWith("/teams/3");
  });

  it("getTeamPlayers calls the correct endpoint", async () => {
    mockGet.mockResolvedValue({ data: [] });
    await teamApi.getTeamPlayers(5);
    expect(mockGet).toHaveBeenCalledWith("/teams/5/players");
  });

  it("getTeamStats calls the correct endpoint", async () => {
    mockGet.mockResolvedValue({ data: [] });
    await teamApi.getTeamStats(2);
    expect(mockGet).toHaveBeenCalledWith("/teams/2/stats");
  });
});

describe("playerApi", () => {
  it("getPlayers returns the paginated response envelope", async () => {
    const envelope = {
      data: [{ id: 1, gamertag: "Scump" }],
      pagination: { page: 1, limit: 25, total: 1, total_pages: 1 },
    };
    mockGet.mockResolvedValue({ data: envelope });

    const result = await playerApi.getPlayers();
    expect(result).toEqual(envelope);
    expect(mockGet).toHaveBeenCalledWith("/players", {
      params: { page: 1, limit: 25 },
    });
  });

  it("getPlayer calls the correct endpoint", async () => {
    const player = { id: 7, gamertag: "Scump" };
    mockGet.mockResolvedValue({ data: player });

    const result = await playerApi.getPlayer(7);
    expect(result).toEqual(player);
    expect(mockGet).toHaveBeenCalledWith("/players/7");
  });

  it("getPlayerKDStats calls the correct endpoint", async () => {
    mockGet.mockResolvedValue({ data: { avg_kd: 1.5, tournament_stats: [] } });
    await playerApi.getPlayerKDStats(7);
    expect(mockGet).toHaveBeenCalledWith("/players/7/kd");
  });
});

describe("transfersApi", () => {
  it("getTransfers unwraps the transfers array from the response envelope", async () => {
    const transfers = [{ id: 1, player_id: 7 }];
    mockGet.mockResolvedValue({ data: { transfers, count: 1, timestamp: 0 } });

    const result = await transfersApi.getTransfers();
    expect(result).toEqual(transfers);
    expect(mockGet).toHaveBeenCalledWith("/transfers", { params: undefined });
  });

  it("getTransfers forwards optional filter params", async () => {
    mockGet.mockResolvedValue({
      data: { transfers: [], count: 0, timestamp: 0 },
    });
    await transfersApi.getTransfers({ season: "2024" });
    expect(mockGet).toHaveBeenCalledWith("/transfers", {
      params: { season: "2024" },
    });
  });
});
