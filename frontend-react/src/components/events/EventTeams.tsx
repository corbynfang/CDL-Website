import { Link } from 'react-router-dom'
import type { TournamentTeam } from '../../types'
import { getTeamLogo } from '../../utils/assets'

interface Props {
  teams: TournamentTeam[] | null
  loading: boolean
}

export default function EventTeams({ teams, loading }: Props) {
  if (loading) {
    return (
      <div className="space-y-2 animate-pulse">
        {Array.from({ length: 8 }).map((_, i) => (
          <div key={i} className="h-14 rounded border border-zinc-800 bg-zinc-900/60" />
        ))}
      </div>
    )
  }

  if (!teams || teams.length === 0) {
    return <p className="text-center text-zinc-600 py-16 text-sm">No teams registered yet.</p>
  }

  const sorted = [...teams].sort((a, b) => {
    if (a.placement != null && b.placement != null) return a.placement - b.placement
    if (a.placement != null) return -1
    if (b.placement != null) return 1
    return 0
  })

  return (
    <div className="border border-[#1a1a1a] divide-y divide-[#1a1a1a]">
      <div className="grid grid-cols-12 px-4 py-2 text-[10px] uppercase tracking-widest text-zinc-700">
        <span className="col-span-1">#</span>
        <span className="col-span-7">Team</span>
        <span className="col-span-2 text-right">W</span>
        <span className="col-span-2 text-right">L</span>
      </div>

      {sorted.map((team, idx) => {
        const logo = getTeamLogo(team.name)
        const placement = team.placement ?? idx + 1
        return (
          <Link
            key={team.id}
            to={`/teams/${team.id}`}
            className="grid grid-cols-12 items-center px-4 py-3 hover:bg-[#111111] transition-colors group"
          >
            <span className="col-span-1 text-sm text-zinc-600 tabular-nums">{placement}</span>
            <div className="col-span-7 flex items-center gap-3">
              {logo ? (
                <img src={logo} alt={team.name} className="w-7 h-7 object-contain" />
              ) : (
                <div className="w-7 h-7 bg-zinc-800 rounded-full flex items-center justify-center text-[9px] font-mono text-zinc-500">
                  {(team.abbreviation ?? '?').slice(0, 2)}
                </div>
              )}
              <span className="text-sm text-zinc-300 group-hover:text-white transition-colors">{team.name}</span>
            </div>
            <span className="col-span-2 text-right text-sm text-zinc-400 tabular-nums">{team.matches_won ?? '—'}</span>
            <span className="col-span-2 text-right text-sm text-zinc-600 tabular-nums">{team.matches_lost ?? '—'}</span>
          </Link>
        )
      })}
    </div>
  )
}
