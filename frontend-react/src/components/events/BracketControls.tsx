import { formatRound, sortedRounds } from '../../utils/eventUtils'

interface Props {
  rounds: string[]
  active: string | null
  onSelect: (round: string | null) => void
}

export default function BracketControls({ rounds, active, onSelect }: Props) {
  const sorted = sortedRounds(rounds)

  return (
    <div className="flex gap-2 flex-wrap">
      <button
        onClick={() => onSelect(null)}
        className={`text-[10px] uppercase tracking-widest px-3 py-1.5 border transition-colors ${
          active === null
            ? 'border-white text-white'
            : 'border-[#1a1a1a] text-zinc-600 hover:text-zinc-400 hover:border-[#2a2a2a]'
        }`}
      >
        All Rounds
      </button>
      {sorted.map(r => (
        <button
          key={r}
          onClick={() => onSelect(r)}
          className={`text-[10px] uppercase tracking-widest px-3 py-1.5 border transition-colors ${
            active === r
              ? 'border-white text-white'
              : 'border-[#1a1a1a] text-zinc-600 hover:text-zinc-400 hover:border-[#2a2a2a]'
          }`}
        >
          {formatRound(r)}
        </button>
      ))}
    </div>
  )
}
