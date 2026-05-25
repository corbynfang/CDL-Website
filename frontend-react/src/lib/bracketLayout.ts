import type { BracketData, BracketMatch } from '../services/api'
import { formatRound } from '../utils/eventUtils'

export const CARD_W = 220
export const CARD_H = 90   // two team rows + divider + padding
export const COL_STRIDE = 320  // CARD_W + 100px gap for SVG connectors
const PAD = 40

// ─── Types ────────────────────────────────────────────────────────────────────

export interface MatchNode {
  match: BracketMatch
  roundKey: string
  roundLabel: string
  col: number
  x: number
  y: number
  section: 'group' | 'winners' | 'losers' | 'final'
  feedsInto: number[]         // winner → next match IDs
  loserFeedsInto: number | null  // loser drop → next match ID
}

export interface Connector {
  fromId: number
  toId: number
  fromX: number   // right edge of source card
  fromY: number   // center-y of source card
  toX: number     // left edge of target card
  toY: number     // center-y of target card
  isLoser: boolean
}

export interface BracketLayout {
  nodes: MatchNode[]
  connectors: Connector[]
  canvasWidth: number
  canvasHeight: number
  colLabels: { col: number; x: number; label: string }[]
}

function mkNode(
  match: BracketMatch,
  roundKey: string,
  col: number,
  y: number,
  section: MatchNode['section'],
): MatchNode {
  return {
    match, roundKey,
    roundLabel: formatRound(roundKey),
    col, x: col * COL_STRIDE, y, section,
    feedsInto: [], loserFeedsInto: null,
  }
}

function sorted(matches: BracketMatch[] = []): BracketMatch[] {
  return [...matches].sort((a, b) => a.bracket_position - b.bracket_position)
}

function buildConnectors(nodes: MatchNode[], nodeById: Map<number, MatchNode>): Connector[] {
  const out: Connector[] = []
  for (const src of nodes) {
    const srcX = src.x + CARD_W
    const srcY = src.y + CARD_H / 2
    for (const toId of src.feedsInto) {
      const tgt = nodeById.get(toId)
      if (!tgt) continue
      out.push({ fromId: src.match.id, toId, fromX: srcX, fromY: srcY, toX: tgt.x, toY: tgt.y + CARD_H / 2, isLoser: false })
    }
    if (src.loserFeedsInto != null) {
      const tgt = nodeById.get(src.loserFeedsInto)
      if (tgt) out.push({ fromId: src.match.id, toId: src.loserFeedsInto, fromX: srcX, fromY: srcY, toX: tgt.x, toY: tgt.y + CARD_H / 2, isLoser: true })
    }
  }
  return out
}

function dimensions(nodes: MatchNode[]) {
  const maxX = Math.max(...nodes.map(n => n.x + CARD_W)) + PAD
  const maxY = Math.max(...nodes.map(n => n.y + CARD_H)) + PAD
  return { canvasWidth: maxX, canvasHeight: maxY }
}

function colLabels(nodes: MatchNode[]) {
  const map = new Map<number, string>()
  for (const n of nodes) {
    if (!map.has(n.col)) map.set(n.col, n.roundLabel)
  }
  return [...map.entries()].sort(([a], [b]) => a - b).map(([col, label]) => ({ col, x: col * COL_STRIDE, label }))
}

// ─── EWC Group Stage + Single Elimination ─────────────────────────────────────
//
// Column layout (0-5):
//   0 opening_match  1 winners+elim  2 decider  3 QF  4 SF  5 GF+3P
//
// Within each group band (height = BAND):
//   Opening M1 at (col 0, bandTop + OPEN_Y)
//   Opening M2 at (col 0, bandTop + OPEN_Y + OPEN_GAP)
//   Winners match at (col 1, bandTop + OPEN_Y)      ← same y as M1
//   Elim match    at (col 1, bandTop + OPEN_Y + OPEN_GAP) ← same y as M2
//   Decider match at (col 2, midpoint of winners+elim centers)
//   QF            at (col 3, midpoint of winners+decider centers)

const EWC_BAND = 320
const EWC_OPEN_Y = 50
const EWC_OPEN_GAP = 160

// winners center = OPEN_Y + CARD_H/2 = 50+45 = 95
// elim center    = OPEN_Y + OPEN_GAP + CARD_H/2 = 255
// decider top    = (95+255)/2 - 45 = 175-45 = 130
const EWC_DECIDER_Y = Math.round((EWC_OPEN_Y + CARD_H / 2 + EWC_OPEN_Y + EWC_OPEN_GAP + CARD_H / 2) / 2 - CARD_H / 2)

// QF top = midpoint of winners center (95) and decider center (175) - 45 = 135-45 = 90
const EWC_QF_Y = Math.round(
  (EWC_OPEN_Y + CARD_H / 2 + EWC_DECIDER_Y + CARD_H / 2) / 2 - CARD_H / 2,
)

function layoutEWC(data: BracketData): BracketLayout {
  const nodes: MatchNode[] = []
  const nodeById = new Map<number, MatchNode>()
  const gs = data.group_stage ?? {}
  const br = data.bracket
  const GROUPS = ['a', 'b', 'c', 'd'] as const

  // map from group index → QF node (filled below)
  const qfByGroup: (MatchNode | null)[] = [null, null, null, null]

  // ── Group stage ─────────────────────────────────────────────────────────────
  for (let gi = 0; gi < GROUPS.length; gi++) {
    const g = GROUPS[gi]
    const bandTop = gi * EWC_BAND
    const pre = `group_play_${g}_`

    const openMatches = sorted(gs[`${pre}opening_match`])
    if (openMatches.length === 0) continue  // no data for this group yet

    for (let i = 0; i < openMatches.length; i++) {
      const y = bandTop + EWC_OPEN_Y + (openMatches.length > 1 ? i * EWC_OPEN_GAP : EWC_OPEN_GAP / 2 - CARD_H / 2 + 10)
      const n = mkNode(openMatches[i], `${pre}opening_match`, 0, y, 'group')
      n.roundLabel = formatRound('opening_match')
      nodes.push(n); nodeById.set(openMatches[i].id, n)
    }

    const wMatches = sorted(gs[`${pre}winners_match`])
    let wNode: MatchNode | null = null
    if (wMatches.length > 0) {
      wNode = mkNode(wMatches[0], `${pre}winners_match`, 1, bandTop + EWC_OPEN_Y, 'group')
      wNode.roundLabel = formatRound('winners_match')
      nodes.push(wNode); nodeById.set(wMatches[0].id, wNode)
    }

    const eMatches = sorted(gs[`${pre}elimination_match`])
    let eNode: MatchNode | null = null
    if (eMatches.length > 0) {
      // If only 1 opening match, position elimination below winners
      const y = openMatches.length > 1
        ? bandTop + EWC_OPEN_Y + EWC_OPEN_GAP
        : bandTop + EWC_OPEN_Y + EWC_OPEN_GAP
      eNode = mkNode(eMatches[0], `${pre}elimination_match`, 1, y, 'group')
      eNode.roundLabel = formatRound('elimination_match')
      nodes.push(eNode); nodeById.set(eMatches[0].id, eNode)
    }

    const dMatches = sorted(gs[`${pre}decider_match`])
    let dNode: MatchNode | null = null
    if (dMatches.length > 0) {
      dNode = mkNode(dMatches[0], `${pre}decider_match`, 2, bandTop + EWC_DECIDER_Y, 'group')
      dNode.roundLabel = formatRound('decider_match')
      nodes.push(dNode); nodeById.set(dMatches[0].id, dNode)
    }

    // Wire group routing
    const openNodes = openMatches.map(m => nodeById.get(m.id)!)
    for (const on of openNodes) {
      if (wNode) on.feedsInto.push(wNode.match.id)
      if (eNode) on.loserFeedsInto = eNode.match.id
    }
    if (wNode) {
      if (dNode) wNode.loserFeedsInto = dNode.match.id
    }
    if (eNode && dNode) {
      eNode.feedsInto.push(dNode.match.id)
    }
    // wNode→QF and dNode→QF wired after QF nodes are placed
  }

  // ── Quarterfinals (col 3) ────────────────────────────────────────────────────
  const qfMatches = sorted(br['quarterfinal'] ?? br['winners_r1'] ?? [])
  for (let i = 0; i < qfMatches.length; i++) {
    const m = qfMatches[i]
    // bracket_position 1-4 maps to group band 0-3; fall back to array index
    // when bracket_position is unset (0) in the enriched-CSV seeder.
    const gi = m.bracket_position > 0 ? m.bracket_position - 1 : i
    if (gi < 0 || gi >= GROUPS.length) continue
    const y = gi * EWC_BAND + EWC_QF_Y
    const n = mkNode(m, 'quarterfinal', 3, y, 'winners')
    nodes.push(n); nodeById.set(m.id, n)
    qfByGroup[gi] = n
  }

  // Wire group → QF
  for (let gi = 0; gi < GROUPS.length; gi++) {
    const g = GROUPS[gi]
    const pre = `group_play_${g}_`
    const qfN = qfByGroup[gi]
    if (!qfN) continue

    const wNode = nodes.find(n => n.roundKey === `${pre}winners_match`) ?? null
    const dNode = nodes.find(n => n.roundKey === `${pre}decider_match`) ?? null
    if (wNode) wNode.feedsInto.push(qfN.match.id)
    if (dNode) dNode.feedsInto.push(qfN.match.id)
  }

  // ── Semifinals (col 4) ───────────────────────────────────────────────────────
  const sfMatches = sorted(br['semifinal'] ?? br['winners_r2'] ?? [])
  const sfNodes: MatchNode[] = []
  for (let i = 0; i < sfMatches.length; i++) {
    const gi1 = i * 2, gi2 = i * 2 + 1
    const qf1 = qfByGroup[gi1], qf2 = qfByGroup[gi2]
    const c1 = (qf1?.y ?? gi1 * EWC_BAND + EWC_QF_Y) + CARD_H / 2
    const c2 = (qf2?.y ?? gi2 * EWC_BAND + EWC_QF_Y) + CARD_H / 2
    const y = Math.round((c1 + c2) / 2 - CARD_H / 2)
    const n = mkNode(sfMatches[i], 'semifinal', 4, y, 'winners')
    nodes.push(n); nodeById.set(sfMatches[i].id, n)
    sfNodes.push(n)
    if (qf1) qf1.feedsInto.push(n.match.id)
    if (qf2) qf2.feedsInto.push(n.match.id)
  }

  // ── Grand Finals + Third Place (col 5) ───────────────────────────────────────
  const gfMatches = sorted(br['grand_finals'] ?? [])
  const tpMatches = sorted(br['third_place_match'] ?? [])

  let gfNode: MatchNode | null = null
  let tpNode: MatchNode | null = null

  if (sfNodes.length >= 2) {
    const sf1c = sfNodes[0].y + CARD_H / 2
    const sf2c = sfNodes[sfNodes.length - 1].y + CARD_H / 2
    const gfY = Math.round((sf1c + sf2c) / 2 - CARD_H / 2)

    if (gfMatches.length > 0) {
      gfNode = mkNode(gfMatches[0], 'grand_finals', 5, gfY, 'final')
      nodes.push(gfNode); nodeById.set(gfMatches[0].id, gfNode)
    }
    if (tpMatches.length > 0) {
      tpNode = mkNode(tpMatches[0], 'third_place_match', 5, gfY + CARD_H + 40, 'final')
      nodes.push(tpNode); nodeById.set(tpMatches[0].id, tpNode)
    }
  } else if (sfNodes.length === 1) {
    if (gfMatches.length > 0) {
      gfNode = mkNode(gfMatches[0], 'grand_finals', 5, sfNodes[0].y, 'final')
      nodes.push(gfNode); nodeById.set(gfMatches[0].id, gfNode)
    }
  }

  // Wire SF → GF / 3P
  for (const sfN of sfNodes) {
    if (gfNode) sfN.feedsInto.push(gfNode.match.id)
    if (tpNode) sfN.loserFeedsInto = tpNode.match.id
  }

  const connectors = buildConnectors(nodes, nodeById)
  const { canvasWidth, canvasHeight } = dimensions(nodes)
  return { nodes, connectors, canvasWidth, canvasHeight, colLabels: colLabels(nodes) }
}

// ─── CDL Double Elimination ────────────────────────────────────────────────────
//
// Column layout:
//   0: winners_r1 (4)   1: elim_r1 (2)   2: winners_r2 (2)   3: elim_r2 (2)
//   4: winners_finals   5: elim_r3        6: elim_finals       7: grand_finals
//
// Winner bracket occupies top region; loser bracket below it.
// Vertical positions are derived from bracket_position and winner/loser region split.

const CDL_ROUND_COL: Record<string, number> = {
  winners_r1: 0,
  elim_r1: 1,
  winners_r2: 2,
  elim_r2: 3,
  winners_r3: 2,  // cold war has an extra WR round
  elim_r3: 5,
  winners_finals: 4,
  elim_r4: 3,     // cold war
  elim_r5: 5,     // cold war
  elim_finals: 6,
  grand_finals: 7,
}

// For each round key, the maximum matches expected (to size the vertical grid)
const CDL_WR_ROUNDS = new Set(['winners_r1', 'winners_r2', 'winners_r3', 'winners_finals'])
const CDL_GF_ROUNDS = new Set(['grand_finals'])

// Drop routing for standard 8-team CDL double-elim.
// Key: `${fromRound}:${fromPos}` → `${toRound}:${toPos}`
const CDL_DROP_TABLE: Record<string, { round: string; pos: number }> = {
  'winners_r1:1': { round: 'elim_r1', pos: 1 },
  'winners_r1:2': { round: 'elim_r1', pos: 2 },
  'winners_r1:3': { round: 'elim_r1', pos: 1 },
  'winners_r1:4': { round: 'elim_r1', pos: 2 },
  'winners_r2:1': { round: 'elim_r2', pos: 1 },
  'winners_r2:2': { round: 'elim_r2', pos: 2 },
  'winners_finals:1': { round: 'elim_finals', pos: 1 },
}

// Advance routing: winner of round+pos → which round+pos (winners bracket only)
const CDL_WIN_TABLE: Record<string, { round: string; pos: number }> = {
  'winners_r1:1': { round: 'winners_r2', pos: 1 },
  'winners_r1:2': { round: 'winners_r2', pos: 1 },
  'winners_r1:3': { round: 'winners_r2', pos: 2 },
  'winners_r1:4': { round: 'winners_r2', pos: 2 },
  'winners_r2:1': { round: 'winners_finals', pos: 1 },
  'winners_r2:2': { round: 'winners_finals', pos: 1 },
  'winners_finals:1': { round: 'grand_finals', pos: 1 },
  'elim_r1:1': { round: 'elim_r2', pos: 1 },
  'elim_r1:2': { round: 'elim_r2', pos: 2 },
  'elim_r2:1': { round: 'elim_r3', pos: 1 },
  'elim_r2:2': { round: 'elim_r3', pos: 1 },
  'elim_r3:1': { round: 'elim_finals', pos: 1 },
  'elim_finals:1': { round: 'grand_finals', pos: 1 },
}

function layoutCDLDoubleElim(data: BracketData): BracketLayout {
  const ROW_H = CARD_H + 30   // vertical step between sibling matches
  const WR_TOP = 0
  const LR_TOP_OFFSET = 4 * ROW_H + 60  // loser bracket starts below WR1 block

  const nodes: MatchNode[] = []
  const nodeById = new Map<number, MatchNode>()

  // Assign y position for each match based on its bracket_position
  // WR bracket positions in col 0 are 1-4 (top to bottom)
  // Derived rounds are spaced to center on their feeders

  function wrY(pos: number, count: number): number {
    // Center the count matches within the WR block height
    const blockH = (4 - 1) * ROW_H  // same height as 4 WR1 matches
    const groupH = blockH / count
    return WR_TOP + (pos - 1) * groupH + Math.round(groupH / 2 - CARD_H / 2)
  }

  function lrY(pos: number, count: number): number {
    const blockH = (2 - 1) * ROW_H  // LR has max 2 matches per round (initially)
    const groupH = blockH / Math.max(count, 1)
    return LR_TOP_OFFSET + (pos - 1) * groupH + Math.round(groupH / 2 - CARD_H / 2)
  }

  const allRounds = Object.keys(data.bracket)
  const roundsByKey = new Map<string, BracketMatch[]>()
  for (const r of allRounds) {
    roundsByKey.set(r, sorted(data.bracket[r] ?? []))
  }

  // Place WR and LR matches
  for (const [roundKey, matches] of roundsByKey) {
    if (matches.length === 0) continue
    const col = CDL_ROUND_COL[roundKey] ?? 7
    const isWR = CDL_WR_ROUNDS.has(roundKey)
    const isGF = CDL_GF_ROUNDS.has(roundKey)

    for (const m of matches) {
      let y: number
      if (isGF) {
        // GF: center between WR and LR regions
        y = Math.round((LR_TOP_OFFSET + 0 + WR_TOP + 3 * ROW_H) / 2 - CARD_H / 2)
      } else if (isWR) {
        y = wrY(m.bracket_position, matches.length)
      } else {
        y = lrY(m.bracket_position, matches.length)
      }

      const section: MatchNode['section'] = isGF ? 'final' : isWR ? 'winners' : 'losers'
      const n = mkNode(m, roundKey, col, y, section)
      nodes.push(n); nodeById.set(m.id, n)
    }
  }

  // Wire winner connections
  const matchByRoundPos = new Map<string, MatchNode>()
  for (const n of nodes) {
    matchByRoundPos.set(`${n.roundKey}:${n.match.bracket_position}`, n)
  }

  for (const n of nodes) {
    const key = `${n.roundKey}:${n.match.bracket_position}`

    const winDest = CDL_WIN_TABLE[key]
    if (winDest) {
      const tgt = matchByRoundPos.get(`${winDest.round}:${winDest.pos}`)
      if (tgt) n.feedsInto.push(tgt.match.id)
    }

    const lossDest = CDL_DROP_TABLE[key]
    if (lossDest) {
      const tgt = matchByRoundPos.get(`${lossDest.round}:${lossDest.pos}`)
      if (tgt) n.loserFeedsInto = tgt.match.id
    }
  }

  const connectors = buildConnectors(nodes, nodeById)
  const { canvasWidth, canvasHeight } = dimensions(nodes)
  return { nodes, connectors, canvasWidth, canvasHeight, colLabels: colLabels(nodes) }
}

export function computeBracketLayout(data: BracketData): BracketLayout {
  switch (data.event_format) {
    case 'ewc_group_stage_single_elim':
      return layoutEWC(data)
    case 'standard_cdl_double_elim':
    case 'cold_war_stage_double_elim':
      return layoutCDLDoubleElim(data)
    default:
      // CDL with group stage: bracket portion uses double-elim layout
      if (data.event_format === 'cdl_major_group_stage_bracket') {
        return layoutCDLDoubleElim(data)
      }
      return layoutCDLDoubleElim(data)
  }
}
