import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import EventTeams from './EventTeams'
import { sampleTeams, opticTexas, atlantaFaze, bostonBreach } from '../../test/fixtures/teams'
import type { TournamentTeam } from '../../types'

vi.mock('../../utils/assets', () => ({
  getTeamLogo: vi.fn().mockReturnValue(null),
  getPlayerAvatar: vi.fn().mockReturnValue('/placeholder.png'),
}))

vi.mock('../../hooks/useApi')
import { useApi } from '../../hooks/useApi'

function wrap(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>)
}

function makeApi(overrides: Partial<ReturnType<typeof useApi>>): ReturnType<typeof useApi<unknown>> {
  return { data: null, loading: false, error: null, refetch: vi.fn(), ...overrides }
}

describe('EventTeams', () => {
  beforeEach(() => vi.mocked(useApi).mockReset())

  it('shows skeleton loaders while loading', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ loading: true }) as ReturnType<typeof useApi<unknown>>)
    const { container } = wrap(<EventTeams tournamentId={1} />)
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument()
  })

  it('shows empty state when data is null', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: null }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventTeams tournamentId={1} />)
    expect(screen.getByText('No teams registered yet.')).toBeInTheDocument()
  })

  it('shows empty state when teams array is empty', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventTeams tournamentId={1} />)
    expect(screen.getByText('No teams registered yet.')).toBeInTheDocument()
  })

  it('renders team names when data loads', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleTeams }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventTeams tournamentId={1} />)
    expect(screen.getByText('OpTic Texas')).toBeInTheDocument()
    expect(screen.getByText('Atlanta FaZe')).toBeInTheDocument()
  })

  it('renders W/L record for each team', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [opticTexas] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventTeams tournamentId={1} />)
    // wins=5, losses=1 — use getAllByText since "1" also appears as placement
    expect(screen.getByText('5')).toBeInTheDocument()
    expect(screen.getAllByText('1').length).toBeGreaterThan(0)
  })

  it('links each team to /teams/:id', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [opticTexas] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventTeams tournamentId={1} />)
    const links = screen.getAllByRole('link')
    expect(links.some(l => l.getAttribute('href') === '/teams/1')).toBe(true)
  })

  it('sorts teams by placement ascending', () => {
    const unsorted: TournamentTeam[] = [bostonBreach, opticTexas, atlantaFaze]
    vi.mocked(useApi).mockReturnValue(makeApi({ data: unsorted }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventTeams tournamentId={1} />)
    const rows = screen.getAllByRole('link')
    // First link should be placement 1 (OpTic Texas)
    expect(rows[0].textContent).toContain('OpTic Texas')
  })

  it('renders placement numbers for each team', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleTeams }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventTeams tournamentId={1} />)
    // Placements 1, 2, 3 appear — use getAllByText since numbers also appear in W/L columns
    expect(screen.getAllByText('1').length).toBeGreaterThan(0)
    expect(screen.getAllByText('2').length).toBeGreaterThan(0)
    expect(screen.getAllByText('3').length).toBeGreaterThan(0)
  })

  it('renders table header columns', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleTeams }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventTeams tournamentId={1} />)
    expect(screen.getByText('#')).toBeInTheDocument()
    expect(screen.getByText('Team')).toBeInTheDocument()
    expect(screen.getByText('W')).toBeInTheDocument()
    expect(screen.getByText('L')).toBeInTheDocument()
  })
})
