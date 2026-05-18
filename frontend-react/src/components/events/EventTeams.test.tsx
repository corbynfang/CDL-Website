import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import EventTeams from './EventTeams'
import { sampleTeams, opticTexas, atlantaFaze, bostonBreach } from '../../test/fixtures/teams'
import type { TournamentTeam } from '../../types'

vi.mock('../../utils/assets', () => ({
  getTeamLogo: vi.fn().mockReturnValue(null),
  getPlayerAvatar: vi.fn().mockReturnValue('/placeholder.png'),
}))

function wrap(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>)
}

describe('EventTeams', () => {
  it('shows skeleton loaders while loading', () => {
    const { container } = wrap(<EventTeams teams={null} loading={true} />)
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument()
  })

  it('shows empty state when teams is null', () => {
    wrap(<EventTeams teams={null} loading={false} />)
    expect(screen.getByText('No teams registered yet.')).toBeInTheDocument()
  })

  it('shows empty state when teams array is empty', () => {
    wrap(<EventTeams teams={[]} loading={false} />)
    expect(screen.getByText('No teams registered yet.')).toBeInTheDocument()
  })

  it('renders team names when data is provided', () => {
    wrap(<EventTeams teams={sampleTeams} loading={false} />)
    expect(screen.getByText('OpTic Texas')).toBeInTheDocument()
    expect(screen.getByText('Atlanta FaZe')).toBeInTheDocument()
  })

  it('renders W/L record for each team', () => {
    wrap(<EventTeams teams={[opticTexas]} loading={false} />)
    expect(screen.getByText('5')).toBeInTheDocument()
    expect(screen.getAllByText('1').length).toBeGreaterThan(0)
  })

  it('links each team to /teams/:id', () => {
    wrap(<EventTeams teams={[opticTexas]} loading={false} />)
    const links = screen.getAllByRole('link')
    expect(links.some(l => l.getAttribute('href') === '/teams/1')).toBe(true)
  })

  it('sorts teams by placement ascending', () => {
    const unsorted: TournamentTeam[] = [bostonBreach, opticTexas, atlantaFaze]
    wrap(<EventTeams teams={unsorted} loading={false} />)
    const rows = screen.getAllByRole('link')
    expect(rows[0].textContent).toContain('OpTic Texas')
  })

  it('renders placement numbers for each team', () => {
    wrap(<EventTeams teams={sampleTeams} loading={false} />)
    expect(screen.getAllByText('1').length).toBeGreaterThan(0)
    expect(screen.getAllByText('2').length).toBeGreaterThan(0)
    expect(screen.getAllByText('3').length).toBeGreaterThan(0)
  })

  it('renders table header columns', () => {
    wrap(<EventTeams teams={sampleTeams} loading={false} />)
    expect(screen.getByText('#')).toBeInTheDocument()
    expect(screen.getByText('Team')).toBeInTheDocument()
    expect(screen.getByText('W')).toBeInTheDocument()
    expect(screen.getByText('L')).toBeInTheDocument()
  })
})
