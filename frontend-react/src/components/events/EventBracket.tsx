import { useState, useRef, useEffect } from 'react'
import type { BracketData } from '../../services/api'
import BracketSkeleton from '../loaders/BracketSkeleton'
import BracketControls from './BracketControls'
import BracketCanvas from './BracketCanvas'

interface Props {
  data: BracketData | null
  loading: boolean
  error: string | null
}

export default function EventBracket({ data, loading, error }: Props) {
  const [activeRound,  setActiveRound]  = useState<string | null>(null)
  const [zoom,         setZoom]         = useState(1.0)
  const [isFullscreen, setIsFullscreen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const onFsChange = () => setIsFullscreen(!!document.fullscreenElement)
    document.addEventListener('fullscreenchange', onFsChange)
    return () => document.removeEventListener('fullscreenchange', onFsChange)
  }, [])

  function toggleFullscreen() {
    if (!document.fullscreenElement) {
      containerRef.current?.requestFullscreen()
    } else {
      document.exitFullscreen()
    }
  }

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
    <div
      ref={containerRef}
      className={`space-y-6 ${isFullscreen ? 'bg-[#09090b] p-6 h-full overflow-auto' : ''}`}
    >
      {rounds.length > 1 && (
        <BracketControls
          rounds={rounds}
          active={activeRound}
          onSelect={setActiveRound}
          zoom={zoom}
          onZoom={setZoom}
          isFullscreen={isFullscreen}
          onFullscreen={toggleFullscreen}
        />
      )}
      <BracketCanvas data={data} activeRound={activeRound} zoom={zoom} />
    </div>
  )
}
