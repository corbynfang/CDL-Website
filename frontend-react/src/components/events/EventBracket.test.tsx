import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import EventBracket from './EventBracket'
import {
  sampleBracketData,
  emptyBracketData,
  cdlGroupBracketData,
  ewcBracketData,
  ewcNoPlayoffData,
} from '../../test/fixtures/matches'

vi.mock('../../utils/assets', () => ({
  getTeamLogo: vi.fn().mockReturnValue(null),
  getPlayerAvatar: vi.fn().mockReturnValue('/placeholder.png'),
}))

function wrap(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>)
}

describe('EventBracket', () => {
  it('shows BracketSkeleton while loading', () => {
    const { container } = wrap(<EventBracket data={null} loading={true} error={null} />)
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument()
  })

  it('shows error message when error is provided', () => {
    wrap(<EventBracket data={null} loading={false} error="Not found" />)
    expect(screen.getByText(/bracket data not available/i)).toBeInTheDocument()
  })

  it('shows empty state when data is null and not loading', () => {
    wrap(<EventBracket data={null} loading={false} error={null} />)
    expect(screen.getByText(/no bracket matches/i)).toBeInTheDocument()
  })

  it('shows empty state when total_matches is 0', () => {
    wrap(<EventBracket data={emptyBracketData} loading={false} error={null} />)
    expect(screen.getByText(/no bracket matches/i)).toBeInTheDocument()
  })

  it('renders bracket content when data is provided', () => {
    wrap(<EventBracket data={sampleBracketData} loading={false} error={null} />)
    expect(screen.getByText('Winners Bracket')).toBeInTheDocument()
  })

  it('renders BracketControls when there are multiple rounds', () => {
    wrap(<EventBracket data={sampleBracketData} loading={false} error={null} />)
    expect(screen.getByRole('button', { name: /all rounds/i })).toBeInTheDocument()
  })

  it('does not render BracketControls when there is only one round', () => {
    const singleRound = { ...sampleBracketData, bracket: { winners_r1: sampleBracketData.bracket.winners_r1 } }
    wrap(<EventBracket data={singleRound} loading={false} error={null} />)
    expect(screen.queryByRole('button', { name: /all rounds/i })).not.toBeInTheDocument()
  })

  it('renders match team abbreviations in the bracket', () => {
    wrap(<EventBracket data={sampleBracketData} loading={false} error={null} />)
    expect(screen.getAllByText('OTX').length).toBeGreaterThan(0)
  })

  // ── Format-aware behaviour ──────────────────────────────────────────────────

  it('standard CDL shows no tabs', () => {
    wrap(<EventBracket data={sampleBracketData} loading={false} error={null} />)
    expect(screen.queryByRole('button', { name: /bracket/i })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /group stage/i })).not.toBeInTheDocument()
  })

  it('CDL major group bracket shows Bracket and Group Stage tabs', () => {
    wrap(<EventBracket data={cdlGroupBracketData} loading={false} error={null} />)
    expect(screen.getByRole('button', { name: /^bracket$/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /^group stage$/i })).toBeInTheDocument()
  })

  it('CDL major group bracket defaults to Bracket tab', () => {
    wrap(<EventBracket data={cdlGroupBracketData} loading={false} error={null} />)
    expect(screen.getByText('Winners Bracket')).toBeInTheDocument()
    // Exact 'Round 1' is the group stage key label; 'Winners Round 1' is a different element
    expect(screen.queryByText('Round 1')).not.toBeInTheDocument()
  })

  it('CDL major group bracket switches to Group Stage tab on click', async () => {
    wrap(<EventBracket data={cdlGroupBracketData} loading={false} error={null} />)
    await userEvent.click(screen.getByRole('button', { name: /^group stage$/i }))
    expect(screen.getByText(/round 1/i)).toBeInTheDocument()
  })

  it('EWC shows no tabs — uses combined tree view', () => {
    wrap(<EventBracket data={ewcBracketData} loading={false} error={null} />)
    expect(screen.queryByRole('button', { name: /^bracket$/i })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /^group stage$/i })).not.toBeInTheDocument()
  })

  it('EWC renders match cards in combined tree', () => {
    wrap(<EventBracket data={ewcBracketData} loading={false} error={null} />)
    // Match abbreviations rendered by BracketMatchCard
    expect(screen.getAllByText('OTX').length).toBeGreaterThan(0)
  })

  it('EWC tree shows zoom controls', () => {
    wrap(<EventBracket data={ewcBracketData} loading={false} error={null} />)
    // BracketTree renders +/- zoom buttons
    expect(screen.getByRole('button', { name: '+' })).toBeInTheDocument()
  })

  it('EWC with only group stage data renders tree (no empty state)', () => {
    wrap(<EventBracket data={ewcNoPlayoffData} loading={false} error={null} />)
    // Should NOT show the generic "no bracket matches" empty state
    expect(screen.queryByText(/no bracket matches/i)).not.toBeInTheDocument()
    // Match cards from group stage are rendered
    expect(screen.getAllByText('OTX').length).toBeGreaterThan(0)
  })
})
