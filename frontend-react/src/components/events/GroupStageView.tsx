import type { BracketMatch } from '../../services/api'
import { formatRound, sortedRounds } from '../../utils/eventUtils'
import BracketMatchCard from './BracketMatchCard'

interface Props {
  groupStage: Record<string, BracketMatch[]>
  format: string
}

export default function GroupStageView({ groupStage, format }: Props) {
  if (format === 'ewc_group_stage_single_elim') {
    return <EWCGroupStageView groupStage={groupStage} />
  }
  return <CDLGroupStageView groupStage={groupStage} />
}

// ─── Shared round column ──────────────────────────────────────────────────────

function RoundColumn({ label, matches }: { label: string; matches: BracketMatch[] }) {
  return (
    <div className="flex flex-col gap-2.5 flex-shrink-0">
      <p className="text-[10px] uppercase tracking-[0.14em] text-zinc-600 text-center pb-1.5 border-b border-[#1e1e1e]">
        {label}
      </p>
      <div className="flex flex-col gap-2.5">
        {matches.map(m => <BracketMatchCard key={m.id} match={m} />)}
      </div>
    </div>
  )
}

// ─── CDL major group stage (flat columns) ─────────────────────────────────────
// Keys: round_1, qualification_match, losers_bracket

function CDLGroupStageView({ groupStage }: { groupStage: Record<string, BracketMatch[]> }) {
  const rounds = sortedRounds(Object.keys(groupStage))
  if (rounds.length === 0) {
    return <p className="text-center text-zinc-600 py-16 text-sm">No group stage data available.</p>
  }
  return (
    <div className="overflow-auto pb-4">
      <div className="flex gap-5 min-w-max">
        {rounds.map(r => (
          <RoundColumn key={r} label={formatRound(r)} matches={groupStage[r] ?? []} />
        ))}
      </div>
    </div>
  )
}

// ─── EWC group stage (group cards A/B/C/D + cross-group rounds) ───────────────
// Keys like group_play_a_* belong to Group A; ungrouped keys are cross-group rounds.

const EWC_GROUPS = ['a', 'b', 'c', 'd'] as const

function EWCGroupStageView({ groupStage }: { groupStage: Record<string, BracketMatch[]> }) {
  const groupKeys: Record<string, string[]> = {}
  const crossGroupKeys: string[] = []

  for (const key of Object.keys(groupStage)) {
    let placed = false
    for (const g of EWC_GROUPS) {
      if (key.startsWith(`group_play_${g}_`)) {
        if (!groupKeys[g]) groupKeys[g] = []
        groupKeys[g].push(key)
        placed = true
        break
      }
    }
    if (!placed) crossGroupKeys.push(key)
  }

  const sortedCross = sortedRounds(crossGroupKeys)
  const activeGroups = EWC_GROUPS.filter(g => groupKeys[g]?.length)

  return (
    <div className="space-y-10">
      {/* Cross-group rounds (opening_match, winners_match, etc.) */}
      {sortedCross.length > 0 && (
        <div className="space-y-3">
          <p className="text-[10px] uppercase tracking-[0.14em] text-zinc-700">All Groups</p>
          <div className="overflow-auto pb-2">
            <div className="flex gap-5 min-w-max">
              {sortedCross.map(r => (
                <RoundColumn key={r} label={formatRound(r)} matches={groupStage[r] ?? []} />
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Per-group sections */}
      {activeGroups.map(g => (
        <EWCGroupSection
          key={g}
          letter={g.toUpperCase()}
          groupPrefix={`group_play_${g}_`}
          roundKeys={groupKeys[g]}
          groupStage={groupStage}
        />
      ))}
    </div>
  )
}

interface EWCGroupSectionProps {
  letter: string
  groupPrefix: string
  roundKeys: string[]
  groupStage: Record<string, BracketMatch[]>
}

function EWCGroupSection({ letter, groupPrefix, roundKeys, groupStage }: EWCGroupSectionProps) {
  // Sort by suffix (strips group prefix, then falls back to alphabetical)
  const sorted = [...roundKeys].sort((a, b) => {
    const sa = a.slice(groupPrefix.length)
    const sb = b.slice(groupPrefix.length)
    return sa.localeCompare(sb)
  })

  return (
    <div className="space-y-3">
      <p className="text-[10px] uppercase tracking-[0.14em] text-zinc-700">Group {letter}</p>
      <div className="overflow-auto pb-2">
        <div className="flex gap-5 min-w-max">
          {sorted.map(r => {
            const suffix = r.slice(groupPrefix.length)
            return (
              <RoundColumn
                key={r}
                label={formatRound(suffix)}
                matches={groupStage[r] ?? []}
              />
            )
          })}
        </div>
      </div>
    </div>
  )
}
