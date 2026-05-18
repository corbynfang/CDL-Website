import { Link } from 'react-router-dom'
import { getTeamLogo } from '../../utils/assets'
import type { BracketMatch } from '../../services/api'

interface Props {
  match: BracketMatch
}

export default function BracketMatchCard({ match }: Props) {
  const complete = match.winner_id != null
  const team1Won = match.winner_id === match.team1_id
  const team2Won = match.winner_id === match.team2_id

  function TeamRow({
    name, abbr, logo, score, won,
  }: { name: string; abbr: string; logo: string; score: number; won: boolean }) {
    const logoSrc = getTeamLogo(name) ?? (logo || undefined)
    const dimmed  = complete && !won

    return (
      <div className={`relative flex items-center gap-2 pl-2.5 pr-3 py-2.5 ${dimmed ? 'opacity-40' : ''}`}>
        {won && <span className="absolute left-0 inset-y-0 w-0.5 bg-white rounded-full" />}

        {logoSrc ? (
          <img src={logoSrc} alt={abbr} className="w-6 h-6 object-contain flex-shrink-0" />
        ) : (
          <div className="w-6 h-6 rounded bg-zinc-800 flex items-center justify-center text-[9px] font-mono font-bold text-zinc-500 flex-shrink-0 leading-none">
            {abbr.slice(0, 2)}
          </div>
        )}

        <span className={`text-[11px] flex-1 truncate ${won ? 'text-white font-bold tracking-wide' : 'text-zinc-400'}`}>
          {abbr}
        </span>

        <span className={`text-[11px] font-mono font-bold tabular-nums min-w-[18px] text-right ${
          won ? 'text-white' : complete ? 'text-zinc-600' : 'text-zinc-700'
        }`}>
          {complete ? score : '–'}
        </span>
      </div>
    )
  }

  return (
    <Link
      to={`/matches/${match.id}`}
      className="block w-[220px] rounded border border-[#1e1e1e] bg-[#0f0f0f] hover:border-[#2e2e2e] hover:bg-[#141414] transition-all overflow-hidden"
    >
      <TeamRow name={match.team1_name} abbr={match.team1_abbr} logo={match.team1_logo} score={match.team1_score} won={team1Won} />
      <div className="h-px bg-[#1e1e1e]" />
      <TeamRow name={match.team2_name} abbr={match.team2_abbr} logo={match.team2_logo} score={match.team2_score} won={team2Won} />
    </Link>
  )
}
