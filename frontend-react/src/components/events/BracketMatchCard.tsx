import { Link } from 'react-router-dom'
import { getTeamLogo } from '../../utils/assets'
import type { BracketMatch } from '../../services/api'

interface Props {
  match: BracketMatch
}

export default function BracketMatchCard({ match }: Props) {
  const team1Won = match.winner_id === match.team1_id
  const team2Won = match.winner_id === match.team2_id
  const complete = match.winner_id != null

  function TeamRow({
    name, abbr, logo, score, won,
  }: { name: string; abbr: string; logo: string; score: number; won: boolean }) {
    const logoSrc = getTeamLogo(name) ?? (logo || undefined)
    return (
      <div className={`flex items-center gap-2 px-3 py-2 ${complete && !won ? 'opacity-40' : ''}`}>
        {logoSrc ? (
          <img src={logoSrc} alt={name} className="w-5 h-5 object-contain flex-shrink-0" />
        ) : (
          <div className="w-5 h-5 rounded-full bg-zinc-800 flex items-center justify-center text-[9px] font-mono text-zinc-500">
            {abbr.slice(0, 2)}
          </div>
        )}
        <span className={`text-xs flex-1 truncate ${won ? 'text-white font-semibold' : 'text-zinc-400'}`}>
          {abbr}
        </span>
        <span className={`text-xs font-bold tabular-nums ${won ? 'text-white' : 'text-zinc-500'}`}>
          {score}
        </span>
      </div>
    )
  }

  return (
    <Link
      to={`/matches/${match.id}`}
      className="block rounded-lg border border-[#1a1a1a] bg-[#111111] hover:border-[#2a2a2a] hover:bg-[#161616] transition-all overflow-hidden"
    >
      <TeamRow name={match.team1_name} abbr={match.team1_abbr} logo={match.team1_logo} score={match.team1_score} won={team1Won} />
      <div className="h-px bg-[#1a1a1a]" />
      <TeamRow name={match.team2_name} abbr={match.team2_abbr} logo={match.team2_logo} score={match.team2_score} won={team2Won} />
    </Link>
  )
}
