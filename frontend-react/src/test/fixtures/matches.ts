import type { Match } from "../../types";
import type { BracketMatch, BracketData } from "../../services/api";

function makeMatch(overrides: Partial<Match> & { id: number }): Match {
  return {
    tournament_id: 1,
    team1_id: 1,
    team2_id: 2,
    match_date: "2025-03-01T18:00:00Z",
    team1_score: 3,
    team2_score: 1,
    winner_id: 1,
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
    team1: {
      id: 1,
      name: "OpTic Texas",
      abbreviation: "OTX",
      is_active: true,
      is_cdl_franchise: true,
      created_at: "",
      updated_at: "",
    },
    team2: {
      id: 2,
      name: "Atlanta FaZe",
      abbreviation: "ATL",
      is_active: true,
      is_cdl_franchise: true,
      created_at: "",
      updated_at: "",
    },
    ...overrides,
  };
}

export const winnersR1Match = makeMatch({
  id: 1,
  match_type: "winners_r1",
  team1_score: 3,
  team2_score: 0,
  winner_id: 1,
});
export const grandFinalsMatch = makeMatch({
  id: 2,
  match_type: "grand_finals",
  team1_score: 3,
  team2_score: 2,
  winner_id: 1,
});
export const noTypeMatch = makeMatch({ id: 3, match_type: undefined });

export const sampleMatches: Match[] = [winnersR1Match, grandFinalsMatch];

function makeBracketMatch(
  overrides: Partial<BracketMatch> & { id: number },
): BracketMatch {
  return {
    team1_id: 1,
    team2_id: 2,
    team1_name: "OpTic Texas",
    team1_abbr: "OTX",
    team1_logo: "",
    team2_name: "Atlanta FaZe",
    team2_abbr: "ATL",
    team2_logo: "",
    team1_score: 3,
    team2_score: 1,
    winner_id: 1,
    bracket_position: 1,
    match_date: "2025-03-01T18:00:00Z",
    ...overrides,
  };
}

export const bracketMatchComplete = makeBracketMatch({ id: 10, winner_id: 1 });
export const bracketMatchNoWinner = makeBracketMatch({
  id: 11,
  winner_id: null,
  team1_score: 0,
  team2_score: 0,
});

export const sampleBracketData: BracketData = {
  tournament_id: 1,
  tournament_name: "CDL Major 2 2025",
  total_matches: 2,
  bracket: {
    winners_r1: [bracketMatchComplete],
    grand_finals: [bracketMatchNoWinner],
  },
};

export const emptyBracketData: BracketData = {
  tournament_id: 1,
  tournament_name: "CDL Major 2 2025",
  total_matches: 0,
  bracket: {},
};

export const coldWarBracketData: BracketData = {
  tournament_id: 12,
  tournament_name: "CDL Stage 1 2021",
  total_matches: 4,
  event_format: "cold_war_stage_double_elim",
  bracket: {
    elim_r4: [bracketMatchComplete],
    elim_r5: [bracketMatchComplete],
    elim_finals: [bracketMatchComplete],
    grand_finals: [bracketMatchNoWinner],
  },
};

export const cdlGroupBracketData: BracketData = {
  tournament_id: 14,
  tournament_name: "CDL Major 1 2023",
  total_matches: 6,
  event_format: "cdl_major_group_stage_bracket",
  bracket: {
    winners_r1: [bracketMatchComplete],
    grand_finals: [bracketMatchNoWinner],
  },
  group_stage: {
    round_1: [bracketMatchComplete],
    qualification_match: [bracketMatchComplete],
    losers_bracket: [bracketMatchNoWinner],
  },
};

const ewcQF1 = makeBracketMatch({ id: 20, bracket_position: 1, winner_id: 1 });
const ewcQF2 = makeBracketMatch({ id: 21, bracket_position: 2, winner_id: 1 });
const ewcSF1 = makeBracketMatch({ id: 22, bracket_position: 1, winner_id: 1 });
const ewcGF = makeBracketMatch({
  id: 23,
  bracket_position: 1,
  winner_id: null,
  team1_score: 0,
  team2_score: 0,
});
const ewcTP = makeBracketMatch({
  id: 24,
  bracket_position: 1,
  winner_id: null,
  team1_score: 0,
  team2_score: 0,
});
const ewcOpenA1 = makeBracketMatch({
  id: 30,
  bracket_position: 1,
  winner_id: 1,
});
const ewcOpenA2 = makeBracketMatch({
  id: 31,
  bracket_position: 2,
  winner_id: 1,
});

export const ewcBracketData: BracketData = {
  tournament_id: 53,
  tournament_name: "Esports World Cup 2025",
  total_matches: 10,
  event_format: "ewc_group_stage_single_elim",
  bracket: {
    quarterfinal: [ewcQF1, ewcQF2],
    semifinal: [ewcSF1],
    grand_finals: [ewcGF],
    third_place_match: [ewcTP],
  },
  group_stage: {
    group_play_a_opening_match: [ewcOpenA1, ewcOpenA2],
    group_play_a_winners_match: [bracketMatchComplete],
    group_play_a_elimination_match: [bracketMatchComplete],
    group_play_a_decider_match: [bracketMatchComplete],
    // extra keys used by GroupStageView tests (ignored by BracketTree layout)
    group_play_a_winners_round_1: [bracketMatchComplete],
    group_play_b_lower_round_1: [bracketMatchComplete],
  },
};

export const ewcNoPlayoffData: BracketData = {
  tournament_id: 52,
  tournament_name: "Esports World Cup 2024",
  total_matches: 4,
  event_format: "ewc_group_stage_single_elim",
  bracket: {
    quarterfinal: [],
    semifinal: [],
    grand_finals: [],
    third_place_match: [],
  },
  group_stage: {
    group_play_a_opening_match: [
      makeBracketMatch({ id: 40, bracket_position: 1, winner_id: 1 }),
    ],
    group_play_a_winners_match: [bracketMatchComplete],
  },
};
