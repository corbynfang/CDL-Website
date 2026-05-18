import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import EventMatches from './EventMatches'
import { sampleMatches, winnersR1Match, noTypeMatch } from '../../test/fixtures/matches'

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

describe('EventMatches', () => {
  beforeEach(() => vi.mocked(useApi).mockReset())

  it('shows skeleton loaders while loading', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ loading: true, data: null }) as ReturnType<typeof useApi<unknown>>)
    const { container } = wrap(<EventMatches tournamentId={1} />)
    // MatchCardSkeleton has animate-pulse
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument()
  })

  it('shows empty state when data is null', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: null }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventMatches tournamentId={1} />)
    expect(screen.getByText('No matches recorded yet.')).toBeInTheDocument()
  })

  it('shows empty state when matches array is empty', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventMatches tournamentId={1} />)
    expect(screen.getByText('No matches recorded yet.')).toBeInTheDocument()
  })

  it('renders match team names when data loads', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleMatches }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventMatches tournamentId={1} />)
    expect(screen.getAllByText('OpTic Texas').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Atlanta FaZe').length).toBeGreaterThan(0)
  })

  it('renders match scores', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [winnersR1Match] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventMatches tournamentId={1} />)
    expect(screen.getByText('3')).toBeInTheDocument()
    expect(screen.getByText('0')).toBeInTheDocument()
  })

  it('groups matches by match_type', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleMatches }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventMatches tournamentId={1} />)
    // The group header <p> renders the round name — both appear at least once
    expect(screen.getAllByText('winners r1').length).toBeGreaterThan(0)
    expect(screen.getAllByText('grand finals').length).toBeGreaterThan(0)
  })

  it('groups matches without match_type under "Other"', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [noTypeMatch] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventMatches tournamentId={1} />)
    expect(screen.getByText('Other')).toBeInTheDocument()
  })

  it('links each match to /matches/:id', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [winnersR1Match] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventMatches tournamentId={1} />)
    const links = screen.getAllByRole('link')
    expect(links.some(l => l.getAttribute('href') === '/matches/1')).toBe(true)
  })
})
