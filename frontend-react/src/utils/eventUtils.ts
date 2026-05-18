export type EventStatus = 'upcoming' | 'live' | 'completed'

const FEATURED_TYPES  = new Set(['major_tournament', 'championship', 'international_major'])
const BRACKET_TYPES   = new Set(['major_tournament', 'championship', 'international_major', 'kickoff'])
const HIDDEN_TYPES    = new Set(['season_summary', 'unknown'])

export const isFeatured  = (type: string) => FEATURED_TYPES.has(type)
export const hasBracket  = (type: string) => BRACKET_TYPES.has(type)
export const isHidden    = (type: string) => HIDDEN_TYPES.has(type)

export function deriveStatus(startDate: string, endDate?: string | null): EventStatus {
  const now   = new Date()
  const start = new Date(startDate)
  if (start > now) return 'upcoming'
  if (endDate) {
    const end = new Date(endDate)
    if (now <= end) return 'live'
    return 'completed'
  }
  // No end_date: treat as live if started within the last 5 days
  const fiveDays = 5 * 24 * 60 * 60 * 1000
  if (now.getTime() - start.getTime() < fiveDays) return 'live'
  return 'completed'
}

export function formatDateRange(startDate: string, endDate?: string | null): string {
  const start = new Date(startDate)
  const opts: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' }
  if (!endDate) return start.toLocaleDateString('en-US', { ...opts, year: 'numeric' })
  const end = new Date(endDate)
  const startStr = start.toLocaleDateString('en-US', opts)
  if (start.getMonth() === end.getMonth() && start.getFullYear() === end.getFullYear()) {
    return `${startStr} – ${end.getDate()}, ${end.getFullYear()}`
  }
  return `${startStr} – ${end.toLocaleDateString('en-US', { ...opts, year: 'numeric' })}`
}

export function formatPrize(amount?: number | null): string {
  if (!amount) return 'TBA'
  if (amount >= 1_000_000) return `$${(amount / 1_000_000).toFixed(amount % 1_000_000 === 0 ? 0 : 1)}M`
  if (amount >= 1_000)     return `$${(amount / 1_000).toFixed(0)}K`
  return `$${amount}`
}

export function monthLabel(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
}

const COUNTRY_FLAGS: Record<string, string> = {
  USA: '🇺🇸', CAN: '🇨🇦', ESP: '🇪🇸', SAU: '🇸🇦',
  GBR: '🇬🇧', FRA: '🇫🇷', AUS: '🇦🇺', DEU: '🇩🇪',
  KOR: '🇰🇷', JPN: '🇯🇵', NLD: '🇳🇱', BRA: '🇧🇷',
}
export const countryFlag = (code: string) => COUNTRY_FLAGS[code?.toUpperCase()] ?? '🌐'

const ROUND_LABELS: Record<string, string> = {
  // Standard CDL double-elim
  winners_r1:          'Winners Round 1',
  winners_r2:          'Winners Round 2',
  winners_r3:          'Winners Round 3',
  winners_finals:      'Winners Finals',
  elim_r1:             'Elimination Round 1',
  elim_r2:             'Elimination Round 2',
  elim_r3:             'Elimination Round 3',
  elim_r4:             'Elimination Round 4',
  elim_r5:             'Elimination Round 5',
  elim_finals:         'Elimination Finals',
  grand_finals:        'Grand Finals',
  // EWC single-elim playoff
  quarterfinal:        'Quarterfinal',
  semifinal:           'Semifinal',
  third_place_match:   'Third Place Match',
  // EWC group stage — cross-group rounds
  opening_match:       'Opening Match',
  winners_match:       'Winners Match',
  decider_match:       'Decider Match',
  elimination_match:   'Elimination Match',
  // CDL major group stage
  round_1:             'Round 1',
  qualification_match: 'Qualification Match',
  losers_bracket:      'Losers Bracket',
}

const ROUND_ORDER: Record<string, number> = {
  // Standard CDL double-elim
  winners_r1: 10, winners_r2: 20, winners_r3: 30, winners_finals: 40,
  elim_r1: 50, elim_r2: 60, elim_r3: 70, elim_r4: 73, elim_r5: 76, elim_finals: 80,
  grand_finals: 90,
  // EWC single-elim playoff (winners_r1 already at 10)
  quarterfinal: 40, semifinal: 60, third_place_match: 85,
  // EWC group stage cross-group rounds
  opening_match: 10, winners_match: 20, decider_match: 30, elimination_match: 40,
  // CDL major group stage
  round_1: 10, qualification_match: 20, losers_bracket: 30,
}

export const formatRound      = (r: string) => ROUND_LABELS[r] ?? r.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())
export const roundOrder       = (r: string) => ROUND_ORDER[r] ?? 99
export const bracketSection   = (r: string): 'winners' | 'elimination' | 'grand_finals' =>
  r === 'grand_finals' ? 'grand_finals' : r.startsWith('winners') ? 'winners' : 'elimination'

export function sortedRounds(rounds: string[]): string[] {
  return [...new Set(rounds)].sort((a, b) => roundOrder(a) - roundOrder(b))
}

export function groupByMonth<T extends { start_date: string }>(items: T[]): [string, T[]][] {
  const map = new Map<string, T[]>()
  for (const item of items) {
    const key = monthLabel(item.start_date)
    if (!map.has(key)) map.set(key, [])
    map.get(key)!.push(item)
  }
  return [...map.entries()]
}
