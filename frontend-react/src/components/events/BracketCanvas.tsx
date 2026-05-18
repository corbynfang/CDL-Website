import type { BracketData } from '../../services/api'
import { formatRound, sortedRounds, bracketSection } from '../../utils/eventUtils'
import BracketMatchCard from './BracketMatchCard'

interface Props {
  data: BracketData
  activeRound: string | null
  zoom?: number
}

export default function BracketCanvas({ data, activeRound, zoom = 1 }: Props) {
  const allRounds = sortedRounds(Object.keys(data.bracket))
  const rounds    = activeRound ? [activeRound] : allRounds

  const winnersRounds    = rounds.filter(r => bracketSection(r) === 'winners')
  const elimRounds       = rounds.filter(r => bracketSection(r) === 'elimination')
  const grandFinalsRound = rounds.filter(r => bracketSection(r) === 'grand_finals')

  function RoundColumn({ round }: { round: string }) {
    const matches = data.bracket[round] ?? []
    return (
      <div className="flex flex-col gap-2.5 flex-shrink-0">
        <p className="text-[10px] uppercase tracking-[0.14em] text-zinc-600 text-center pb-1.5 border-b border-[#1e1e1e]">
          {formatRound(round)}
        </p>
        <div className="flex flex-col gap-2.5">
          {matches.map(m => <BracketMatchCard key={m.id} match={m} />)}
        </div>
      </div>
    )
  }

  function Section({ label, rounds: sectionRounds }: { label: string; rounds: string[] }) {
    if (sectionRounds.length === 0) return null
    return (
      <div className="space-y-3">
        <p className="text-[10px] uppercase tracking-[0.14em] text-zinc-700">{label}</p>
        <div className="flex gap-5">
          {sectionRounds.map(r => <RoundColumn key={r} round={r} />)}
        </div>
      </div>
    )
  }

  return (
    <div className="overflow-auto pb-4">
      <div style={{ zoom }} className="space-y-10 min-w-max">
        <Section label="Winners Bracket" rounds={winnersRounds} />
        <Section label="Elimination Bracket" rounds={elimRounds} />
        <Section label="Grand Finals" rounds={grandFinalsRound} />
      </div>
    </div>
  )
}
