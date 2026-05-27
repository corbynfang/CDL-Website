import type { BracketMatch } from '../services/api'

export interface ParticipantType {
  id: string
  name?: string | null
  resultText?: string | null
  isWinner?: boolean
  status?: string | null
  abbr?: string
  logo?: string
  score?: number
}

export interface MatchType {
  id: string
  name: string
  nextMatchId: string | null
  nextLooserMatchId?: string
  tournamentRoundText: string
  startTime: string
  state: string
  participants: ParticipantType[]
}

export interface PrototypeBracketData {
  upper: MatchType[]
  lower: MatchType[]
}

export function adaptBracketMatch(
  match: BracketMatch,
  roundText: string,
  nextMatchId: string | null,
  nextLooserMatchId?: string | null,
): MatchType {
  const team1Won = match.winner_id === match.team1_id
  const team2Won = match.winner_id === match.team2_id
  const complete = match.winner_id != null

  const top: ParticipantType = {
    id: String(match.team1_id),
    name: match.team1_name,
    resultText: complete ? String(match.team1_score) : null,
    isWinner: team1Won,
    status: complete ? 'PLAYED' : null,
    // Extra fields (ignored by the library but available to custom match cards)
    abbr: match.team1_abbr,
    logo: match.team1_logo,
    score: match.team1_score,
  }

  const bottom: ParticipantType = {
    id: String(match.team2_id),
    name: match.team2_name,
    resultText: complete ? String(match.team2_score) : null,
    isWinner: team2Won,
    status: complete ? 'PLAYED' : null,
    abbr: match.team2_abbr,
    logo: match.team2_logo,
    score: match.team2_score,
  }

  return {
    id: String(match.id),
    name: `${match.team1_abbr} vs ${match.team2_abbr}`,
    nextMatchId,
    nextLooserMatchId: nextLooserMatchId ?? undefined,
    tournamentRoundText: roundText,
    startTime: match.match_date,
    state: complete ? 'SCORE_DONE' : 'NO_PARTY',
    participants: [top, bottom],
  }
}

function makeParticipant(
  id: string,
  name: string,
  abbr: string,
  score: number,
  isWinner: boolean,
  complete: boolean,
): ParticipantType {
  return {
    id,
    name,
    resultText: complete ? String(score) : null,
    isWinner,
    status: complete ? 'PLAYED' : null,
    abbr,
    score,
  }
}

function makeMatch(
  id: string,
  roundText: string,
  nextMatchId: string | null,
  nextLooserMatchId: string | null,
  topId: string, topName: string, topAbbr: string, topScore: number, topWon: boolean,
  botId: string, botName: string, botAbbr: string, botScore: number, botWon: boolean,
  complete = true,
): MatchType {
  return {
    id,
    name: `${topAbbr} vs ${botAbbr}`,
    nextMatchId,
    nextLooserMatchId: nextLooserMatchId ?? undefined,
    tournamentRoundText: roundText,
    startTime: '2023-02-03',
    state: complete ? 'SCORE_DONE' : 'NO_PARTY',
    participants: [
      makeParticipant(topId, topName, topAbbr, topScore, topWon, complete),
      makeParticipant(botId, botName, botAbbr, botScore, botWon, complete),
    ],
  }
}

export const CDL_MAJOR_1_2023_FIXTURE: PrototypeBracketData = {
  upper: [
    // UB Round 1
    makeMatch('ubr1m1', 'UB R1', 'ubr2m1', 'lbr1m1',
      'otx', 'OpTic Texas',       'OTX', 3, true,
      'mnr', 'Minnesota RØKKR',   'MNR', 1, false),
    makeMatch('ubr1m2', 'UB R1', 'ubr2m1', 'lbr1m2',
      'atl', 'Atlanta FaZe',      'ATL', 3, true,
      'sea', 'Seattle Surge',     'SEA', 2, false),
    makeMatch('ubr1m3', 'UB R1', 'ubr2m2', 'lbr1m1',
      'lat', 'Los Angeles Thieves', 'LAT', 3, true,
      'bos', 'Boston Breach',     'BOS', 1, false),
    makeMatch('ubr1m4', 'UB R1', 'ubr2m2', 'lbr1m2',
      'nys', 'New York Subliners', 'NYS', 3, true,
      'mia', 'Miami Heretics',    'MIA', 0, false),

    // UB Semifinals
    makeMatch('ubr2m1', 'UB SF', 'ubf', 'lbr2m1',
      'otx', 'OpTic Texas',       'OTX', 3, true,
      'atl', 'Atlanta FaZe',      'ATL', 2, false),
    makeMatch('ubr2m2', 'UB SF', 'ubf', 'lbr2m2',
      'lat', 'Los Angeles Thieves', 'LAT', 3, true,
      'nys', 'New York Subliners', 'NYS', 1, false),

    // UB Final
    makeMatch('ubf', 'UB Final', 'gf', 'lbf',
      'otx', 'OpTic Texas',       'OTX', 3, true,
      'lat', 'Los Angeles Thieves', 'LAT', 2, false),

    // Grand Finals
    makeMatch('gf', 'Grand Final', null, null,
      'otx', 'OpTic Texas',       'OTX', 3, true,
      'lat', 'Los Angeles Thieves', 'LAT', 1, false),
  ],
  lower: [
    // LB Round 1
    makeMatch('lbr1m1', 'LB R1', 'lbr2m1', null,
      'mnr', 'Minnesota RØKKR',   'MNR', 3, true,
      'bos', 'Boston Breach',     'BOS', 2, false),
    makeMatch('lbr1m2', 'LB R1', 'lbr2m2', null,
      'sea', 'Seattle Surge',     'SEA', 3, true,
      'mia', 'Miami Heretics',    'MIA', 1, false),

    // LB Quarterfinals (UB SF losers enter)
    makeMatch('lbr2m1', 'LB QF', 'lbsf', null,
      'atl', 'Atlanta FaZe',      'ATL', 3, true,
      'mnr', 'Minnesota RØKKR',   'MNR', 0, false),
    makeMatch('lbr2m2', 'LB QF', 'lbsf', null,
      'nys', 'New York Subliners', 'NYS', 3, true,
      'sea', 'Seattle Surge',     'SEA', 2, false),

    // LB Semifinal
    makeMatch('lbsf', 'LB SF', 'lbf', null,
      'atl', 'Atlanta FaZe',      'ATL', 3, true,
      'nys', 'New York Subliners', 'NYS', 2, false),

    // LB Final (UB Final loser enters)
    makeMatch('lbf', 'LB Final', 'gf', null,
      'lat', 'Los Angeles Thieves', 'LAT', 3, true,
      'atl', 'Atlanta FaZe',      'ATL', 1, false),
  ],
}
