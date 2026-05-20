package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeDoubleElimRoundKey(t *testing.T) {
	tests := []struct {
		raw  string
		want string
	}{
		// Clear aliases fixed
		{"winners_final", "winners_finals"},
		{"grand_final", "grand_finals"},
		{"losers_round_1", "elim_r1"},
		{"losers_round_2", "elim_r2"},
		{"losers_round_3", "elim_r3"},
		{"losers_final", "elim_finals"},
		// Already-canonical values pass through unchanged
		{"winners_r1", "winners_r1"},
		{"winners_r2", "winners_r2"},
		{"winners_r3", "winners_r3"},
		{"winners_finals", "winners_finals"},
		{"elim_r1", "elim_r1"},
		{"elim_r2", "elim_r2"},
		{"elim_r3", "elim_r3"},
		{"elim_r4", "elim_r4"}, // CW 2021 extra rounds pass through
		{"elim_r5", "elim_r5"},
		{"elim_finals", "elim_finals"},
		{"grand_finals", "grand_finals"},
		// Ambiguous group-stage values pass through — not classified into any bucket
		{"round_1", "round_1"},
		{"qualification_match", "qualification_match"},
		{"losers_bracket", "losers_bracket"},
		// Anything unknown passes through
		{"some_future_round", "some_future_round"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			assert.Equal(t, tt.want, normalizeDoubleElimRoundKey(tt.raw))
		})
	}
}

func TestNormalizeEWCRoundKey(t *testing.T) {
	assert.Equal(t, "grand_finals", normalizeEWCRoundKey("grand_final"))
	assert.Equal(t, "quarterfinal", normalizeEWCRoundKey("quarterfinal"))
	assert.Equal(t, "semifinal", normalizeEWCRoundKey("semifinal"))
	assert.Equal(t, "opening_match", normalizeEWCRoundKey("opening_match"))
	assert.Equal(t, "winners_r1", normalizeEWCRoundKey("winners_r1"))
	assert.Equal(t, "grand_finals", normalizeEWCRoundKey("grand_finals")) // already canonical
	assert.Equal(t, "", normalizeEWCRoundKey(""))
}

func TestDetectBracketFormat(t *testing.T) {
	tests := []struct {
		tournamentFormat string
		tournamentType   string
		want             bracketFormat
	}{
		// tournament_format column takes precedence over tournament_type
		{"standard_cdl_double_elim", "international_major", bracketFmtStandardCDLDoubleElim},
		{"cold_war_stage_double_elim", "major_tournament", bracketFmtColdWarStageDoubleElim},
		{"cdl_major_group_stage_bracket", "major_tournament", bracketFmtCDLMajorGroupBracket},
		{"ewc_group_stage_single_elim", "major_tournament", bracketFmtEWCGroupBracket},
		// Empty tournament_format falls back to tournament_type
		{"", "major_tournament", bracketFmtStandardCDLDoubleElim},
		{"", "championship", bracketFmtStandardCDLDoubleElim},
		{"", "kickoff", bracketFmtStandardCDLDoubleElim},
		{"", "international_major", bracketFmtEWCGroupBracket},
		{"", "qualifier", bracketFmtUnknown},
		{"", "minor_tournament", bracketFmtUnknown},
		{"", "season_summary", bracketFmtUnknown},
		{"", "unknown", bracketFmtUnknown},
		{"", "", bracketFmtUnknown},
	}
	for _, tt := range tests {
		name := tt.tournamentFormat + "/" + tt.tournamentType
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, detectBracketFormat(tt.tournamentFormat, tt.tournamentType))
		})
	}
}

func TestFormatName(t *testing.T) {
	assert.Equal(t, "standard_cdl_double_elim", formatName(bracketFmtStandardCDLDoubleElim))
	assert.Equal(t, "cold_war_stage_double_elim", formatName(bracketFmtColdWarStageDoubleElim))
	assert.Equal(t, "cdl_major_group_stage_bracket", formatName(bracketFmtCDLMajorGroupBracket))
	assert.Equal(t, "ewc_group_stage_single_elim", formatName(bracketFmtEWCGroupBracket))
	assert.Equal(t, "unknown", formatName(bracketFmtUnknown))
}

func TestBracketKeysFor(t *testing.T) {
	t.Run("standard CDL has 8 keys without elim_r4/r5", func(t *testing.T) {
		keys := bracketKeysFor(bracketFmtStandardCDLDoubleElim)
		assert.Len(t, keys, 8)
		assert.Contains(t, keys, "winners_r1")
		assert.Contains(t, keys, "elim_r3")
		assert.Contains(t, keys, "grand_finals")
		assert.NotContains(t, keys, "elim_r4")
		assert.NotContains(t, keys, "elim_r5")
	})
	t.Run("cold war stage has 10 keys with elim_r4 and elim_r5", func(t *testing.T) {
		keys := bracketKeysFor(bracketFmtColdWarStageDoubleElim)
		assert.Len(t, keys, 10)
		assert.Contains(t, keys, "elim_r4")
		assert.Contains(t, keys, "elim_r5")
	})
	t.Run("CDL major group bracket has same 8 bracket keys as standard CDL", func(t *testing.T) {
		keys := bracketKeysFor(bracketFmtCDLMajorGroupBracket)
		assert.Len(t, keys, 8)
		assert.Contains(t, keys, "winners_r1")
		assert.NotContains(t, keys, "elim_r4")
		assert.NotContains(t, keys, "round_1")
		assert.NotContains(t, keys, "losers_bracket")
	})
	t.Run("EWC has single-elim playoff keys", func(t *testing.T) {
		keys := bracketKeysFor(bracketFmtEWCGroupBracket)
		assert.Contains(t, keys, "winners_r1")
		assert.Contains(t, keys, "quarterfinal")
		assert.Contains(t, keys, "semifinal")
		assert.Contains(t, keys, "grand_finals")
		assert.Contains(t, keys, "third_place_match")
		assert.NotContains(t, keys, "opening_match")
		assert.NotContains(t, keys, "decider_match")
	})
	t.Run("unknown format returns nil", func(t *testing.T) {
		assert.Nil(t, bracketKeysFor(bracketFmtUnknown))
	})
}

func TestHasGroupStage(t *testing.T) {
	assert.True(t, hasGroupStage(bracketFmtCDLMajorGroupBracket))
	assert.True(t, hasGroupStage(bracketFmtEWCGroupBracket))
	assert.False(t, hasGroupStage(bracketFmtStandardCDLDoubleElim))
	assert.False(t, hasGroupStage(bracketFmtColdWarStageDoubleElim))
	assert.False(t, hasGroupStage(bracketFmtUnknown))
}

func TestRoundNormalizerFor(t *testing.T) {
	t.Run("standard CDL applies double-elim normalizer", func(t *testing.T) {
		norm := roundNormalizerFor(bracketFmtStandardCDLDoubleElim)
		assert.Equal(t, "winners_finals", norm("winners_final"))
		assert.Equal(t, "grand_finals", norm("grand_final"))
		assert.Equal(t, "elim_r1", norm("losers_round_1"))
		assert.Equal(t, "elim_finals", norm("losers_final"))
	})
	t.Run("cold war stage applies double-elim normalizer", func(t *testing.T) {
		norm := roundNormalizerFor(bracketFmtColdWarStageDoubleElim)
		assert.Equal(t, "winners_finals", norm("winners_final"))
		assert.Equal(t, "grand_finals", norm("grand_final"))
		assert.Equal(t, "elim_r4", norm("elim_r4")) // canonical, passes through
	})
	t.Run("CDL major group bracket applies double-elim normalizer", func(t *testing.T) {
		norm := roundNormalizerFor(bracketFmtCDLMajorGroupBracket)
		assert.Equal(t, "winners_finals", norm("winners_final"))
		assert.Equal(t, "elim_r1", norm("losers_round_1"))
		// Group-stage keys pass through unchanged (routing handled by the handler)
		assert.Equal(t, "round_1", norm("round_1"))
		assert.Equal(t, "losers_bracket", norm("losers_bracket"))
	})
	t.Run("EWC applies EWC normalizer — only grand_final remapped", func(t *testing.T) {
		norm := roundNormalizerFor(bracketFmtEWCGroupBracket)
		assert.Equal(t, "grand_finals", norm("grand_final"))
		assert.Equal(t, "opening_match", norm("opening_match"))
		assert.Equal(t, "quarterfinal", norm("quarterfinal"))
		assert.Equal(t, "winners_r1", norm("winners_r1"))
	})
	t.Run("unknown format returns identity", func(t *testing.T) {
		norm := roundNormalizerFor(bracketFmtUnknown)
		assert.Equal(t, "winners_final", norm("winners_final"))
		assert.Equal(t, "losers_round_1", norm("losers_round_1"))
	})
}

func TestRoundRouting_CDLMajorGroupBracket(t *testing.T) {
	keys := bracketKeysFor(bracketFmtCDLMajorGroupBracket)
	norm := roundNormalizerFor(bracketFmtCDLMajorGroupBracket)

	bracketRounds := []string{
		"winners_r1", "winners_r2", "winners_finals",
		"elim_r1", "elim_r2", "elim_r3", "elim_finals", "grand_finals",
	}
	for _, r := range bracketRounds {
		_, inBracket := keys[norm(r)]
		assert.True(t, inBracket, "%q should route to bracket", r)
	}

	// Aliased forms also land in bracket after normalization
	for raw, canonical := range map[string]string{
		"winners_final": "winners_finals",
		"grand_final":   "grand_finals",
		"losers_round_1": "elim_r1",
		"losers_final":   "elim_finals",
	} {
		normalized := norm(raw)
		assert.Equal(t, canonical, normalized)
		_, inBracket := keys[normalized]
		assert.True(t, inBracket, "%q → %q should route to bracket", raw, canonical)
	}

	// Group-stage keys do NOT land in bracket (handler sends them to group_stage)
	groupRounds := []string{"round_1", "qualification_match", "losers_bracket"}
	for _, r := range groupRounds {
		_, inBracket := keys[norm(r)]
		assert.False(t, inBracket, "%q should NOT be in bracket (goes to group_stage)", r)
	}

	assert.True(t, hasGroupStage(bracketFmtCDLMajorGroupBracket))
}

func TestRoundRouting_EWC(t *testing.T) {
	keys := bracketKeysFor(bracketFmtEWCGroupBracket)
	norm := roundNormalizerFor(bracketFmtEWCGroupBracket)

	bracketRounds := []string{"winners_r1", "quarterfinal", "semifinal", "grand_finals", "third_place_match"}
	for _, r := range bracketRounds {
		_, inBracket := keys[norm(r)]
		assert.True(t, inBracket, "%q should route to bracket", r)
	}

	// grand_final (singular) normalizes to grand_finals and lands in bracket
	normalized := norm("grand_final")
	assert.Equal(t, "grand_finals", normalized)
	_, inBracket := keys[normalized]
	assert.True(t, inBracket, "grand_final → grand_finals should be in bracket")

	// Group-stage rounds do NOT land in bracket
	groupRounds := []string{
		"opening_match", "winners_match", "decider_match", "elimination_match",
		"group_play_a_winners_round_1", "group_play_b_lower_qualifier_round",
	}
	for _, r := range groupRounds {
		_, inBracket := keys[norm(r)]
		assert.False(t, inBracket, "%q should NOT be in bracket (goes to group_stage)", r)
	}

	assert.True(t, hasGroupStage(bracketFmtEWCGroupBracket))
}

func TestRoundRouting_StandardCDL_DropsUnclassified(t *testing.T) {
	keys := bracketKeysFor(bracketFmtStandardCDLDoubleElim)
	norm := roundNormalizerFor(bracketFmtStandardCDLDoubleElim)

	// Standard CDL has no group_stage — unclassified rounds are silently dropped
	assert.False(t, hasGroupStage(bracketFmtStandardCDLDoubleElim))

	for _, r := range []string{"round_1", "qualification_match", "losers_bracket"} {
		_, inBracket := keys[norm(r)]
		assert.False(t, inBracket, "%q should be dropped (no group_stage for this format)", r)
	}

	// elim_r4/r5 are also NOT in standard CDL bracket
	for _, r := range []string{"elim_r4", "elim_r5"} {
		_, inBracket := keys[norm(r)]
		assert.False(t, inBracket, "%q should NOT be in standard CDL bracket", r)
	}
}

func TestRoundRouting_ColdWar_IncludesR4R5(t *testing.T) {
	keys := bracketKeysFor(bracketFmtColdWarStageDoubleElim)
	assert.Contains(t, keys, "elim_r4")
	assert.Contains(t, keys, "elim_r5")
	assert.False(t, hasGroupStage(bracketFmtColdWarStageDoubleElim))
}
