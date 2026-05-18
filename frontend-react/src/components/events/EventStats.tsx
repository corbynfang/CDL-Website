import { Link } from 'react-router-dom'
import type { PlayerTournamentStats } from '../../types'
import { getPlayerAvatar } from '../../utils/assets'

interface Props {
  stats: PlayerTournamentStats[] | null
  loading: boolean
}

const kdColor = (kd: number) => {
  if (kd >= 1.3) return 'text-green-400'
  if (kd < 1.0) return 'text-zinc-600'
  return 'text-white'
}

export default function EventStats({ stats, loading }: Props) {
  if (loading) {
    return (
      <div className="space-y-2 animate-pulse">
        {Array.from({ length: 10 }).map((_, i) => (
          <div key={i} className="h-12 rounded border border-zinc-800 bg-zinc-900/60" />
        ))}
      </div>
    )
  }

  if (!stats || stats.length === 0) {
    return <p className="text-center text-zinc-600 py-16 text-sm">No stats available for this event.</p>
  }

  return (
    <div className="border border-[#1a1a1a] divide-y divide-[#1a1a1a]">
      <div className="grid grid-cols-12 px-4 py-2 text-[10px] uppercase tracking-widest text-zinc-700">
        <span className="col-span-1">#</span>
        <span className="col-span-5">Player</span>
        <span className="col-span-2 text-center hidden sm:block">Maps</span>
        <span className="col-span-2 text-center">K</span>
        <span className="col-span-2 text-right">K/D</span>
      </div>

      {stats.map((s, idx) => {
        const gamertag = s.player?.gamertag ?? `Player #${s.player_id}`
        const teamAbbr = s.team?.abbreviation ?? ''
        const avatar = getPlayerAvatar(gamertag)

        return (
          <Link
            key={s.id}
            to={`/players/${s.player_id}`}
            className="group grid grid-cols-12 items-center px-4 py-3 hover:bg-[#111111] transition-colors"
          >
            <span className="col-span-1 text-xs text-zinc-600 tabular-nums">{idx + 1}</span>

            <div className="col-span-5 flex items-center gap-2">
              <img
                src={avatar}
                alt={gamertag}
                className="w-7 h-7 rounded-full object-cover opacity-80"
              />
              <div className="min-w-0">
                <p className="text-sm text-zinc-300 group-hover:text-white transition-colors truncate">{gamertag}</p>
                {teamAbbr && (
                  <p className="text-[10px] text-zinc-600 uppercase tracking-wider">{teamAbbr}</p>
                )}
              </div>
            </div>

            <span className="col-span-2 text-center text-sm text-zinc-500 tabular-nums hidden sm:block">
              {s.overall_maps || '—'}
            </span>
            <span className="col-span-2 text-center text-sm text-zinc-400 tabular-nums">
              {s.total_kills}
            </span>
            <span className={`col-span-2 text-right text-sm font-bold tabular-nums ${kdColor(s.kd_ratio)}`}>
              {s.kd_ratio.toFixed(2)}
            </span>
          </Link>
        )
      })}
    </div>
  )
}
