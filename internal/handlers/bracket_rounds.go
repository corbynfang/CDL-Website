package handlers
type bracketFormat int

const (
	bracketFmtUnknown                bracketFormat = iota
	bracketFmtStandardCDLDoubleElim
	bracketFmtColdWarStageDoubleElim
	bracketFmtCDLMajorGroupBracket
	bracketFmtEWCGroupBracket
)

func detectBracketFormat(tournamentFormat, tournamentType string) bracketFormat {
	switch tournamentFormat {
	case "standard_cdl_double_elim":
		return bracketFmtStandardCDLDoubleElim
	case "cold_war_stage_double_elim":
		return bracketFmtColdWarStageDoubleElim
	case "cdl_major_group_stage_bracket":
		return bracketFmtCDLMajorGroupBracket
	case "ewc_group_stage_single_elim":
		return bracketFmtEWCGroupBracket
	}
	switch tournamentType {
	case "major_tournament", "championship", "kickoff":
		return bracketFmtStandardCDLDoubleElim
	case "international_major":
		return bracketFmtEWCGroupBracket
	default:
		return bracketFmtUnknown
	}
}

// formatName returns the event_format string included in the bracket API response.
func formatName(f bracketFormat) string {
	switch f {
	case bracketFmtStandardCDLDoubleElim:
		return "standard_cdl_double_elim"
	case bracketFmtColdWarStageDoubleElim:
		return "cold_war_stage_double_elim"
	case bracketFmtCDLMajorGroupBracket:
		return "cdl_major_group_stage_bracket"
	case bracketFmtEWCGroupBracket:
		return "ewc_group_stage_single_elim"
	default:
		return "unknown"
	}
}

func bracketKeysFor(f bracketFormat) map[string]struct{} {
	set := func(ss ...string) map[string]struct{} {
		m := make(map[string]struct{}, len(ss))
		for _, s := range ss {
			m[s] = struct{}{}
		}
		return m
	}
	switch f {
	case bracketFmtStandardCDLDoubleElim, bracketFmtCDLMajorGroupBracket:
		return set(
			"winners_r1", "winners_r2", "winners_finals",
			"elim_r1", "elim_r2", "elim_r3",
			"elim_finals", "grand_finals",
		)
	case bracketFmtColdWarStageDoubleElim:
		return set(
			"winners_r1", "winners_r2", "winners_finals",
			"elim_r1", "elim_r2", "elim_r3", "elim_r4", "elim_r5",
			"elim_finals", "grand_finals",
		)
	case bracketFmtEWCGroupBracket:
		return set(
			"winners_r1", "winners_r2",
			"quarterfinal", "semifinal",
			"grand_finals", "third_place_match",
		)
	default:
		return nil
	}
}

func hasGroupStage(f bracketFormat) bool {
	return f == bracketFmtCDLMajorGroupBracket || f == bracketFmtEWCGroupBracket
}

func normalizeDoubleElimRoundKey(raw string) string {
	switch raw {
	case "winners_final":
		return "winners_finals"
	case "grand_final":
		return "grand_finals"
	case "losers_round_1":
		return "elim_r1"
	case "losers_round_2":
		return "elim_r2"
	case "losers_round_3":
		return "elim_r3"
	case "losers_final":
		return "elim_finals"
	default:
		return raw
	}
}

// normalizeEWCRoundKey normalizes EWC bracket_round values.
// Only grand_final → grand_finals needs mapping; all other EWC keys are clean.
func normalizeEWCRoundKey(raw string) string {
	if raw == "grand_final" {
		return "grand_finals"
	}
	return raw
}

// ewcGroupRoundTypes are the group-stage round names used in EWC 2024 data that
// lack a group_play_X_ prefix. EWC 2025 data already has the full prefixed form.
var ewcGroupRoundTypes = map[string]bool{
	"opening_match":    true,
	"winners_match":    true,
	"decider_match":    true,
	"elimination_match": true,
}

// ewcGroupKey adds a group_play_X_ prefix to legacy EWC group-stage round keys
// using bracket_position to identify the group (1→a, 2→b, 3→c, 4→d).
// Already-prefixed keys (EWC 2025+) pass through unchanged.
// This makes EWC 2024 group data compatible with EWCGroupStageView's prefix-based grouping.
func ewcGroupKey(key string, position int) string {
	if !ewcGroupRoundTypes[key] {
		return key
	}
	if position < 1 || position > 4 {
		return key
	}
	return "group_play_" + string(rune('a'+position-1)) + "_" + key
}

// roundNormalizerFor returns the round key normalizer for the given bracket format.
// Formats without a dedicated normalizer get an identity function.
func roundNormalizerFor(f bracketFormat) func(string) string {
	switch f {
	case bracketFmtStandardCDLDoubleElim, bracketFmtColdWarStageDoubleElim,
		bracketFmtCDLMajorGroupBracket:
		return normalizeDoubleElimRoundKey
	case bracketFmtEWCGroupBracket:
		return normalizeEWCRoundKey
	default:
		return func(s string) string { return s }
	}
}
