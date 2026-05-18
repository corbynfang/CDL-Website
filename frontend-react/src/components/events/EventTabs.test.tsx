import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import EventTabs from './EventTabs'

describe('EventTabs', () => {
  const noop = vi.fn()

  it('always renders Overview, Matches, Teams', () => {
    render(<EventTabs active="overview" onSelect={noop} tournamentType="qualifier" hasStats={false} />)
    expect(screen.getByRole('button', { name: /overview/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /matches/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /teams/i })).toBeInTheDocument()
  })

  it('does not render Bracket tab for qualifier tournament type', () => {
    render(<EventTabs active="overview" onSelect={noop} tournamentType="qualifier" hasStats={false} />)
    expect(screen.queryByRole('button', { name: /bracket/i })).not.toBeInTheDocument()
  })

  it('renders Bracket tab for major_tournament type', () => {
    render(<EventTabs active="overview" onSelect={noop} tournamentType="major_tournament" hasStats={false} />)
    expect(screen.getByRole('button', { name: /bracket/i })).toBeInTheDocument()
  })

  it('renders Bracket tab for championship type', () => {
    render(<EventTabs active="overview" onSelect={noop} tournamentType="championship" hasStats={false} />)
    expect(screen.getByRole('button', { name: /bracket/i })).toBeInTheDocument()
  })

  it('renders Bracket tab for kickoff type', () => {
    render(<EventTabs active="overview" onSelect={noop} tournamentType="kickoff" hasStats={false} />)
    expect(screen.getByRole('button', { name: /bracket/i })).toBeInTheDocument()
  })

  it('renders Stats tab when hasStats is true', () => {
    render(<EventTabs active="overview" onSelect={noop} tournamentType="qualifier" hasStats={true} />)
    expect(screen.getByRole('button', { name: /stats/i })).toBeInTheDocument()
  })

  it('does not render Stats tab when hasStats is false', () => {
    render(<EventTabs active="overview" onSelect={noop} tournamentType="qualifier" hasStats={false} />)
    expect(screen.queryByRole('button', { name: /stats/i })).not.toBeInTheDocument()
  })

  it('calls onSelect with the clicked tab id', async () => {
    const onSelect = vi.fn()
    render(<EventTabs active="overview" onSelect={onSelect} tournamentType="major_tournament" hasStats={true} />)
    await userEvent.click(screen.getByRole('button', { name: /matches/i }))
    expect(onSelect).toHaveBeenCalledWith('matches')
  })

  it('calls onSelect with bracket when bracket tab is clicked', async () => {
    const onSelect = vi.fn()
    render(<EventTabs active="overview" onSelect={onSelect} tournamentType="major_tournament" hasStats={false} />)
    await userEvent.click(screen.getByRole('button', { name: /bracket/i }))
    expect(onSelect).toHaveBeenCalledWith('bracket')
  })
})
