import type { Tournament, TournamentTeam } from '../../types'
import { deriveStatus, eventDisplayName, formatDateRange, formatPrize, countryFlag } from '../../utils/eventUtils'
import TeamLogoStrip from './TeamLogoStrip'

interface Props {
  event: Tournament
  teamCount: number
  teams: TournamentTeam[]
}

const STATUS_BADGE = {
  live:      'bg-emerald-500/10 text-emerald-400 border-emerald-500/30',
  upcoming:  'bg-blue-500/10 text-blue-400 border-blue-500/30',
  completed: 'bg-zinc-800/50 text-zinc-500 border-zinc-700/40',
}

export default function EventHero({ event, teamCount, teams }: Props) {
  const status = deriveStatus(event.start_date, event.end_date)

  return (
    <div className="border-b border-[#1a1a1a] bg-[#0a0a0a]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-14 pb-10 space-y-5">
        <div className="flex flex-wrap items-center gap-2 text-xs uppercase tracking-widest text-zinc-600">
          <span>{event.season?.game_code ?? '—'}</span>
          <span>·</span>
          <span>{event.tournament_type.replace(/_/g, ' ')}</span>
          {event.is_lan && (
            <>
              <span>·</span>
              <span>LAN</span>
            </>
          )}
        </div>

        <h1 className="font-grotesk text-4xl sm:text-5xl font-bold text-white leading-tight">
          {eventDisplayName(event.slug, event.name)}
        </h1>

        <div className="flex flex-wrap items-center gap-4 text-sm text-zinc-400">
          <span
            className={`inline-flex items-center gap-1.5 text-[10px] font-bold uppercase tracking-widest px-2.5 py-1 rounded-full border ${STATUS_BADGE[status]}`}
          >
            {status === 'live' && (
              <span className="relative flex h-1.5 w-1.5">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
                <span className="relative inline-flex rounded-full h-1.5 w-1.5 bg-emerald-500" />
              </span>
            )}
            {status}
          </span>

          <span>{formatDateRange(event.start_date, event.end_date)}</span>

          {event.location && (
            <span>
              {event.country ? countryFlag(event.country) + ' ' : ''}{event.location}
            </span>
          )}

          {event.prize_pool && (
            <span className="text-white font-semibold">{formatPrize(event.prize_pool)}</span>
          )}

          {teamCount > 0 && (
            <span className="text-zinc-600">{teamCount} teams</span>
          )}
        </div>

        {teams.length > 0 && <TeamLogoStrip teams={teams} max={16} size="md" />}
      </div>
    </div>
  )
}
