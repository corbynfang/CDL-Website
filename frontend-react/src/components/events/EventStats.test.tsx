import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import EventStats from './EventStats'
import { sampleStats, highKDStat, lowKDStat, noPlayerStat } from '../../test/fixtures/stats'

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

describe('EventStats', () => {
  beforeEach(() => vi.mocked(useApi).mockReset())

  it('shows skeleton loaders while loading', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ loading: true }) as ReturnType<typeof useApi<unknown>>)
    const { container } = wrap(<EventStats tournamentId={1} />)
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument()
  })

  it('shows empty state when data is null', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: null }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    expect(screen.getByText('No stats available for this event.')).toBeInTheDocument()
  })

  it('shows empty state when stats array is empty', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    expect(screen.getByText('No stats available for this event.')).toBeInTheDocument()
  })

  it('renders player gamertags', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleStats }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    expect(screen.getByText('Scump')).toBeInTheDocument()
    expect(screen.getByText('Simp')).toBeInTheDocument()
  })

  it('renders K/D ratios formatted to 2 decimal places', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [highKDStat] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    expect(screen.getByText('1.45')).toBeInTheDocument()
  })

  it('renders kill counts', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [highKDStat] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    expect(screen.getByText('120')).toBeInTheDocument()
  })

  it('renders team abbreviations', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [lowKDStat] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    expect(screen.getByText('ATL')).toBeInTheDocument()
  })

  it('links each player row to /players/:id', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [highKDStat] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    const links = screen.getAllByRole('link')
    expect(links.some(l => l.getAttribute('href') === '/players/1')).toBe(true)
  })

  it('falls back to "Player #N" when player is undefined', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: [noPlayerStat] }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    expect(screen.getByText('Player #99')).toBeInTheDocument()
  })

  it('renders table header columns', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleStats }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    expect(screen.getByText('Player')).toBeInTheDocument()
    expect(screen.getByText('K')).toBeInTheDocument()
    expect(screen.getByText('K/D')).toBeInTheDocument()
  })

  it('renders row index numbers starting from 1', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleStats }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventStats tournamentId={1} />)
    expect(screen.getByText('1')).toBeInTheDocument()
    expect(screen.getByText('2')).toBeInTheDocument()
  })
})
