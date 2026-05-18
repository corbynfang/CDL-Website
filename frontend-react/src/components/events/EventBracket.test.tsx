import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import EventBracket from './EventBracket'
import { sampleBracketData, emptyBracketData } from '../../test/fixtures/matches'

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

describe('EventBracket', () => {
  beforeEach(() => vi.mocked(useApi).mockReset())

  it('shows BracketSkeleton while loading', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ loading: true }) as ReturnType<typeof useApi<unknown>>)
    const { container } = wrap(<EventBracket tournamentId={1} />)
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument()
  })

  it('shows error message when useApi returns an error', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ error: 'Not found' }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventBracket tournamentId={1} />)
    expect(screen.getByText(/bracket data not available/i)).toBeInTheDocument()
  })

  it('shows empty state when data is null', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: null }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventBracket tournamentId={1} />)
    expect(screen.getByText(/no bracket matches/i)).toBeInTheDocument()
  })

  it('shows empty state when total_matches is 0', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: emptyBracketData }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventBracket tournamentId={1} />)
    expect(screen.getByText(/no bracket matches/i)).toBeInTheDocument()
  })

  it('renders bracket content when data loads successfully', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleBracketData }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventBracket tournamentId={1} />)
    expect(screen.getByText('Winners Bracket')).toBeInTheDocument()
  })

  it('renders BracketControls when there are multiple rounds', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleBracketData }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventBracket tournamentId={1} />)
    // sampleBracketData has winners_r1 + grand_finals — controls should show
    expect(screen.getByRole('button', { name: /all rounds/i })).toBeInTheDocument()
  })

  it('does not render BracketControls when there is only one round', () => {
    const singleRound = { ...sampleBracketData, bracket: { winners_r1: sampleBracketData.bracket.winners_r1 } }
    vi.mocked(useApi).mockReturnValue(makeApi({ data: singleRound }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventBracket tournamentId={1} />)
    expect(screen.queryByRole('button', { name: /all rounds/i })).not.toBeInTheDocument()
  })

  it('renders match team abbreviations in the bracket', () => {
    vi.mocked(useApi).mockReturnValue(makeApi({ data: sampleBracketData }) as ReturnType<typeof useApi<unknown>>)
    wrap(<EventBracket tournamentId={1} />)
    expect(screen.getAllByText('OTX').length).toBeGreaterThan(0)
  })
})
