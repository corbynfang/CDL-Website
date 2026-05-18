import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import type { TournamentDetail, TournamentTeam } from '../types'
import BlobbyLoader from './loaders/BlobbyLoader'
import EventHero from './events/EventHero'
import EventTabs, { type TabId } from './events/EventTabs'
import EventOverview from './events/EventOverview'
import EventBracket from './events/EventBracket'
import EventMatches from './events/EventMatches'
import EventTeams from './events/EventTeams'
import EventStats from './events/EventStats'

export default function EventDetailPage() {
  const { slug } = useParams<{ slug: string }>()
  const [tab, setTab] = useState<TabId>('overview')

  const { data, loading, error } = useApi<TournamentDetail>(
    `/api/v1/tournaments/slug/${slug}`,
    { enabled: !!slug }
  )

  const { data: teams } = useApi<TournamentTeam[]>(
    `/api/v1/tournaments/${data?.tournament.id}/teams`,
    { enabled: !!data?.tournament.id }
  )

  if (loading) return <BlobbyLoader label="Loading event..." />

  if (error || !data) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-20 text-center">
        <p className="text-zinc-500 text-sm mb-4">Event not found.</p>
        <Link to="/events" className="text-xs uppercase tracking-widest text-zinc-600 hover:text-white transition-colors">
          ← Back to Events
        </Link>
      </div>
    )
  }

  const { tournament, team_count } = data
  const teamList = teams ?? []

  return (
    <div className="min-h-screen bg-[#0a0a0a]">
      <EventHero event={tournament} teamCount={team_count} teams={teamList} />

      <EventTabs
        active={tab}
        onSelect={setTab}
        tournamentType={tournament.tournament_type}
        hasStats={true}
      />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {tab === 'overview' && (
          <EventOverview event={tournament} teamCount={team_count} />
        )}
        {tab === 'bracket' && (
          <EventBracket tournamentId={tournament.id} />
        )}
        {tab === 'matches' && (
          <EventMatches tournamentId={tournament.id} />
        )}
        {tab === 'teams' && (
          <EventTeams tournamentId={tournament.id} />
        )}
        {tab === 'stats' && (
          <EventStats tournamentId={tournament.id} />
        )}
      </div>
    </div>
  )
}
