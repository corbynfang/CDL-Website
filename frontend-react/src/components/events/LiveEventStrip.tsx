import { Link } from 'react-router-dom'
import type { Tournament } from '../../types'
import { eventDisplayName, formatDateRange, countryFlag } from '../../utils/eventUtils'

interface Props {
  event: Tournament
}

export default function LiveEventStrip({ event }: Props) {
  return (
    <Link
      to={`/events/${event.slug}`}
      className="flex items-center justify-between gap-4 px-5 py-3 rounded-xl border border-emerald-500/30 bg-emerald-500/5 hover:bg-emerald-500/10 transition-colors group"
    >
      <div className="flex items-center gap-3">
        <span className="relative flex h-2.5 w-2.5 flex-shrink-0">
          <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
          <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500" />
        </span>
        <span className="text-xs font-bold uppercase tracking-widest text-emerald-400">Live Now</span>
        <span className="text-xs text-zinc-400 hidden sm:block">·</span>
        <span className="text-sm font-semibold text-white hidden sm:block">{eventDisplayName(event.slug, event.name)}</span>
        {event.location && (
          <span className="text-xs text-zinc-500 hidden md:block">
            {event.country ? countryFlag(event.country) : ''} {event.location}
          </span>
        )}
      </div>
      <span className="text-xs text-zinc-500 group-hover:text-zinc-400 transition-colors whitespace-nowrap">
        {formatDateRange(event.start_date, event.end_date)} →
      </span>
    </Link>
  )
}
