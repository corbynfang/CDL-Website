import { useState, useRef, useEffect } from 'react'
import type { BracketData } from '../../services/api'
import BracketSkeleton from '../loaders/BracketSkeleton'
import BracketControls from './BracketControls'
import BracketCanvas from './BracketCanvas'
import GroupStageView from './GroupStageView'

interface Props {
  data: BracketData | null
  loading: boolean
  error: string | null
}

export default function EventBracket({ data, loading, error }: Props) {
  const [activeRound,       setActiveRound]       = useState<string | null>(null)
  const [zoom,              setZoom]              = useState(1.0)
  const [isFullscreen,      setIsFullscreen]      = useState(false)
  // null = no user choice yet; derives from format on first render with data
  const [userSelectedTab,   setUserSelectedTab]   = useState<'bracket' | 'group_stage' | null>(null)
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

  const isEWC       = data.event_format === 'ewc_group_stage_single_elim'
  const hasGroupStage = !!(data.group_stage && Object.keys(data.group_stage).length > 0)
  const showTabs    = hasGroupStage

  // Derive active tab: user selection takes precedence, then format default
  const activeTab: 'bracket' | 'group_stage' =
    userSelectedTab ?? (isEWC ? 'group_stage' : 'bracket')

  const bracketRounds = Object.keys(data.bracket)
  const bracketMatchCount = bracketRounds.reduce((n, r) => n + (data.bracket[r]?.length ?? 0), 0)
  const hasPlayoffMatches = bracketMatchCount > 0

  const tabPill = 'text-[11px] uppercase tracking-widest px-4 py-2 border-b-2 transition-colors'
  const tabActive   = 'border-white text-white'
  const tabInactive = 'border-transparent text-zinc-600 hover:text-zinc-400'

  return (
    <div
      ref={containerRef}
      className={`space-y-6 ${isFullscreen ? 'bg-[#09090b] p-6 h-full overflow-auto' : ''}`}
    >
      {/* Tab switcher */}
      {showTabs && (
        <div className="flex gap-0 border-b border-[#1e1e1e]">
          <button
            onClick={() => setUserSelectedTab('bracket')}
            className={`${tabPill} ${activeTab === 'bracket' ? tabActive : tabInactive}`}
          >
            Bracket
          </button>
          <button
            onClick={() => setUserSelectedTab('group_stage')}
            className={`${tabPill} ${activeTab === 'group_stage' ? tabActive : tabInactive}`}
          >
            Group Stage
          </button>
        </div>
      )}

      {/* Bracket tab */}
      {activeTab === 'bracket' && (
        <>
          {bracketRounds.length > 1 && (
            <BracketControls
              rounds={bracketRounds}
              active={activeRound}
              onSelect={setActiveRound}
              zoom={zoom}
              onZoom={setZoom}
              isFullscreen={isFullscreen}
              onFullscreen={toggleFullscreen}
            />
          )}
          {isEWC && !hasPlayoffMatches ? (
            <p className="text-center text-zinc-600 py-16 text-sm">
              Playoff bracket data is not available yet.
            </p>
          ) : (
            <BracketCanvas
              data={data}
              activeRound={activeRound}
              zoom={zoom}
              flat={isEWC}
            />
          )}
        </>
      )}

      {/* Group stage tab */}
      {activeTab === 'group_stage' && data.group_stage && (
        <GroupStageView
          groupStage={data.group_stage}
          format={data.event_format ?? ''}
        />
      )}
    </div>
  )
}
