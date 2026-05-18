import type { BracketData } from '../../services/api'
import { formatRound, sortedRounds, bracketSection } from '../../utils/eventUtils'
import BracketMatchCard from './BracketMatchCard'

interface Props {
  data: BracketData
  activeRound: string | null
}

export default function BracketCanvas({ data, activeRound }: Props) {
  const allRounds = sortedRounds(Object.keys(data.bracket))
  const rounds = activeRound ? [activeRound] : allRounds

  const winnersRounds   = rounds.filter(r => bracketSection(r) === 'winners')
  const elimRounds      = rounds.filter(r => bracketSection(r) === 'elimination')
  const grandFinalsRound = rounds.filter(r => bracketSection(r) === 'grand_finals')

  function RoundColumn({ round }: { round: string }) {
    const matches = data.bracket[round] ?? []
    return (
      <div className="flex flex-col gap-3 min-w-[180px]">
        <p className="text-[10px] uppercase tracking-widest text-zinc-600 text-center">
          {formatRound(round)}
        </p>
        {matches.map(m => <BracketMatchCard key={m.id} match={m} />)}
      </div>
    )
  }

  function Section({ label, rounds: sectionRounds }: { label: string; rounds: string[] }) {
    if (sectionRounds.length === 0) return null
    return (
      <div className="space-y-3">
        <p className="text-xs uppercase tracking-widest text-zinc-700">{label}</p>
        <div className="flex gap-4 overflow-x-auto pb-2">
          {sectionRounds.map(r => <RoundColumn key={r} round={r} />)}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      <Section label="Winners Bracket" rounds={winnersRounds} />
      <Section label="Elimination Bracket" rounds={elimRounds} />
      <Section label="Grand Finals" rounds={grandFinalsRound} />
    </div>
  )
}
