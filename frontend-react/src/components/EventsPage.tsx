import { useState, useMemo } from 'react'
import { useApi } from '../hooks/useApi'
import type { Tournament } from '../types'
import {
  deriveStatus,
  isHidden,
  isFeatured,
  groupByMonth,
} from '../utils/eventUtils'
import EventFilters from './events/EventFilters'
import FeaturedEventCard from './events/FeaturedEventCard'
import CompactEventRow from './events/CompactEventRow'
import LiveEventStrip from './events/LiveEventStrip'
import EventCardSkeleton from './loaders/EventCardSkeleton'

interface Filters {
  game: string
  type: string
  status: string
}

export default function EventsPage() {
  const [filters, setFilters] = useState<Filters>({ game: '', type: '', status: '' })

  const { data: events, loading } = useApi<Tournament[]>('/api/v1/tournaments')

  const visible = useMemo(() => {
    if (!events) return []
    return events.filter(e => {
      if (isHidden(e.tournament_type)) return false
      if (filters.game && e.season?.game_code !== filters.game) return false
      if (filters.type && e.tournament_type !== filters.type) return false
      if (filters.status && deriveStatus(e.start_date, e.end_date) !== filters.status) return false
      return true
    })
  }, [events, filters])

  const liveEvents = useMemo(
    () => visible.filter(e => deriveStatus(e.start_date, e.end_date) === 'live'),
    [visible]
  )

  // Group remaining events by month (newest first)
  const reversed = [...visible].reverse()
  const byMonth = groupByMonth(reversed)

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4 mb-8">
        <div>
          <p className="text-xs uppercase tracking-widest text-[#737373] mb-2">CDL</p>
          <h1 className="font-grotesk text-3xl font-bold text-white">EVENTS</h1>
          {!loading && (
            <p className="text-[#737373] text-sm mt-1">{visible.length} events</p>
          )}
        </div>
        <EventFilters filters={filters} onChange={setFilters} />
      </div>

      {liveEvents.length > 0 && (
        <div className="mb-6 space-y-2">
          {liveEvents.map(e => <LiveEventStrip key={e.id} event={e} />)}
        </div>
      )}

      {loading && (
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-3">
          {Array.from({ length: 9 }, (_, i) => <EventCardSkeleton key={`skeleton-${i}`} />)}
        </div>
      )}

      {!loading && visible.length === 0 && (
        <p className="text-center text-[#737373] py-20 text-sm">No events match your filters.</p>
      )}

      {!loading && byMonth.map(([month, monthEvents]) => (
        <div key={month} className="mb-10">
          <p className="text-xs uppercase tracking-widest text-zinc-700 mb-4">{month}</p>

          {/* Featured cards for majors/championships */}
          {monthEvents.some(e => isFeatured(e.tournament_type)) && (
            <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-3 mb-3">
              {monthEvents
                .filter(e => isFeatured(e.tournament_type))
                .map(e => <FeaturedEventCard key={e.id} event={e} />)
              }
            </div>
          )}

          {/* Compact rows for qualifiers/kickoffs */}
          {monthEvents.some(e => !isFeatured(e.tournament_type)) && (
            <div className="border border-[#1a1a1a] divide-y divide-[#1a1a1a]">
              {monthEvents
                .filter(e => !isFeatured(e.tournament_type))
                .map(e => <CompactEventRow key={e.id} event={e} />)
              }
            </div>
          )}
        </div>
      ))}
    </div>
  )
}
