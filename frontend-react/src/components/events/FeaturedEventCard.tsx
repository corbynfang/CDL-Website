import { Link } from 'react-router-dom'
import type { Tournament } from '../../types'
import type { TournamentTeam } from '../../types'
import { deriveStatus, formatDateRange, formatPrize, countryFlag } from '../../utils/eventUtils'
import TeamLogoStrip from './TeamLogoStrip'

interface Props {
  event: Tournament
  teams?: TournamentTeam[]
}

const STATUS_STYLES = {
  live:      'bg-emerald-500/10 text-emerald-400 border-emerald-500/30',
  upcoming:  'bg-blue-500/10 text-blue-400 border-blue-500/30',
  completed: 'bg-zinc-800/50 text-zinc-500 border-zinc-700/40',
}

const LEFT_BORDER = {
  live:      'border-l-emerald-500',
  upcoming:  'border-l-blue-500',
  completed: 'border-l-zinc-700',
}

export default function FeaturedEventCard({ event, teams = [] }: Props) {
  const status = deriveStatus(event.start_date, event.end_date)

  return (
    <Link
      to={`/events/${event.slug}`}
      className={`group block rounded-2xl border border-[#1a1a1a] border-l-4 ${LEFT_BORDER[status]} bg-[#111111] hover:bg-[#161616] hover:border-[#2a2a2a] transition-all overflow-hidden`}
    >
      <div className="p-5 space-y-3">
        <div className="flex items-center justify-between gap-2">
          <span
            className={`inline-flex items-center gap-1.5 text-[10px] font-bold uppercase tracking-widest px-2 py-0.5 rounded-full border ${STATUS_STYLES[status]}`}
          >
            {status === 'live' && (
              <span className="relative flex h-1.5 w-1.5">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
                <span className="relative inline-flex rounded-full h-1.5 w-1.5 bg-emerald-500" />
              </span>
            )}
            {status}
          </span>
          {event.is_lan && (
            <span className="text-[10px] uppercase tracking-widest text-zinc-600">LAN</span>
          )}
        </div>

        <div>
          <p className="text-xs uppercase tracking-widest text-zinc-600 mb-1">
            {event.season?.game_code ?? event.tournament_type.replace(/_/g, ' ')}
          </p>
          <h3 className="font-grotesk text-lg font-bold text-white leading-tight group-hover:text-zinc-100 transition-colors">
            {event.name}
          </h3>
        </div>

        <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-zinc-500">
          <span>{formatDateRange(event.start_date, event.end_date)}</span>
          {event.location && (
            <span>
              {event.country ? countryFlag(event.country) + ' ' : ''}{event.location}
            </span>
          )}
          {event.prize_pool && (
            <span className="text-zinc-400">{formatPrize(event.prize_pool)}</span>
          )}
        </div>

        {teams.length > 0 && (
          <TeamLogoStrip teams={teams} max={12} />
        )}
      </div>
    </Link>
  )
}
