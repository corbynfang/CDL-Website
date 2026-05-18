import { useState } from 'react'
import { useApi } from '../../hooks/useApi'
import type { BracketData } from '../../services/api'
import BracketSkeleton from '../loaders/BracketSkeleton'
import BracketControls from './BracketControls'
import BracketCanvas from './BracketCanvas'

interface Props {
  tournamentId: number
}

export default function EventBracket({ tournamentId }: Props) {
  const [activeRound, setActiveRound] = useState<string | null>(null)
  const { data, loading, error } = useApi<BracketData>(`/api/v1/tournaments/${tournamentId}/bracket`)

  if (loading) return <BracketSkeleton />

  if (error) {
    return (
      <p className="text-center text-zinc-600 py-16 text-sm">
        Bracket data not available for this event.
      </p>
    )
  }

  if (!data || data.total_matches === 0) {
    return (
      <p className="text-center text-zinc-600 py-16 text-sm">
        No bracket matches have been played yet.
      </p>
    )
  }

  const rounds = Object.keys(data.bracket)

  return (
    <div className="space-y-6">
      {rounds.length > 1 && (
        <BracketControls rounds={rounds} active={activeRound} onSelect={setActiveRound} />
      )}
      <BracketCanvas data={data} activeRound={activeRound} />
    </div>
  )
}
