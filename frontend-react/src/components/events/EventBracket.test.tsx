import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import EventBracket from './EventBracket'
import { sampleBracketData, emptyBracketData } from '../../test/fixtures/matches'

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
})
