import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useApi } from '../hooks/useApi'
import type { TournamentDetail, TournamentTeam, Match, PlayerTournamentStats } from '../types'
import type { BracketData } from '../services/api'
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
  // Tracks which tabs have been opened at least once — prevents re-fetching on tab switch.
  const [loadedTabs, setLoadedTabs] = useState<Set<TabId>>(new Set())

  const handleTabSelect = (next: TabId) => {
    setLoadedTabs(prev => new Set([...prev, next]))
    setTab(next)
  }

  // ── Tournament detail (always, needed for the page to render at all) ────────
  const { data, loading, error } = useApi<TournamentDetail>(
    `/api/v1/tournaments/slug/${slug}`,
    { enabled: !!slug }
  )

  const id = data?.tournament.id
  const hasId = !!id

  // ── Teams: eager — needed immediately for the EventHero logo strip ──────────
  const { data: teams, loading: teamsLoading } = useApi<TournamentTeam[]>(
    `/api/v1/tournaments/${id}/teams`,
    { enabled: hasId }
  )

  // ── Matches: deferred until Matches tab is first opened ────────────────────
  const { data: matches, loading: matchesLoading, error: matchesError } = useApi<Match[]>(
    `/api/v1/tournaments/${id}/matches`,
    { enabled: hasId && loadedTabs.has('matches') }
  )

  // ── Bracket: deferred until Bracket tab is first opened ────────────────────
  const { data: bracketData, loading: bracketLoading, error: bracketError } = useApi<BracketData>(
    `/api/v1/tournaments/${id}/bracket`,
    { enabled: hasId && loadedTabs.has('bracket') }
  )

  // ── Stats: deferred until Stats tab is first opened ────────────────────────
  const { data: stats, loading: statsLoading } = useApi<PlayerTournamentStats[]>(
    `/api/v1/tournaments/${id}/stats`,
    { enabled: hasId && loadedTabs.has('stats') }
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

  return (
    <div className="min-h-screen bg-[#0a0a0a]">
      <EventHero
        event={tournament}
        teamCount={team_count}
        teams={teams ?? []}
      />

      <EventTabs
        active={tab}
        onSelect={handleTabSelect}
        tournamentType={tournament.tournament_type}
        hasStats={true}
      />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {tab === 'overview' && (
          <EventOverview event={tournament} teamCount={team_count} />
        )}
        {tab === 'bracket' && (
          <EventBracket data={bracketData} loading={bracketLoading} error={bracketError} />
        )}
        {tab === 'matches' && (
          <EventMatches matches={matches} loading={matchesLoading} error={matchesError} />
        )}
        {tab === 'teams' && (
          <EventTeams teams={teams} loading={teamsLoading} />
        )}
        {tab === 'stats' && (
          <EventStats stats={stats} loading={statsLoading} />
        )}
      </div>
    </div>
  )
}
