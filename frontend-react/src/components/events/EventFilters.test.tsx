import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import EventFilters from './EventFilters'

const defaultFilters = { game: '', type: '', status: '' }

describe('EventFilters', () => {
  it('renders all three selects', () => {
    render(<EventFilters filters={defaultFilters} onChange={vi.fn()} />)
    const selects = screen.getAllByRole('combobox')
    expect(selects).toHaveLength(3)
  })

  it('shows all game options', () => {
    render(<EventFilters filters={defaultFilters} onChange={vi.fn()} />)
    expect(screen.getByRole('option', { name: 'All Games' })).toBeInTheDocument()
    expect(screen.getByRole('option', { name: 'Black Ops 6' })).toBeInTheDocument()
    expect(screen.getByRole('option', { name: 'Modern Warfare III' })).toBeInTheDocument()
  })

  it('shows all type options', () => {
    render(<EventFilters filters={defaultFilters} onChange={vi.fn()} />)
    expect(screen.getByRole('option', { name: 'All Types' })).toBeInTheDocument()
    expect(screen.getByRole('option', { name: 'Major' })).toBeInTheDocument()
    expect(screen.getByRole('option', { name: 'Championship' })).toBeInTheDocument()
    expect(screen.getByRole('option', { name: 'Qualifier' })).toBeInTheDocument()
  })

  it('shows all status options', () => {
    render(<EventFilters filters={defaultFilters} onChange={vi.fn()} />)
    expect(screen.getByRole('option', { name: 'Upcoming' })).toBeInTheDocument()
    expect(screen.getByRole('option', { name: 'Live' })).toBeInTheDocument()
    expect(screen.getByRole('option', { name: 'Completed' })).toBeInTheDocument()
  })

  it('calls onChange with updated game when game select changes', async () => {
    const onChange = vi.fn()
    render(<EventFilters filters={defaultFilters} onChange={onChange} />)
    const [gameSelect] = screen.getAllByRole('combobox')
    await userEvent.selectOptions(gameSelect, 'BO6')
    expect(onChange).toHaveBeenCalledWith({ game: 'BO6', type: '', status: '' })
  })

  it('calls onChange with updated type when type select changes', async () => {
    const onChange = vi.fn()
    render(<EventFilters filters={defaultFilters} onChange={onChange} />)
    const [, typeSelect] = screen.getAllByRole('combobox')
    await userEvent.selectOptions(typeSelect, 'major_tournament')
    expect(onChange).toHaveBeenCalledWith({ game: '', type: 'major_tournament', status: '' })
  })

  it('calls onChange with updated status when status select changes', async () => {
    const onChange = vi.fn()
    render(<EventFilters filters={defaultFilters} onChange={onChange} />)
    const [,, statusSelect] = screen.getAllByRole('combobox')
    await userEvent.selectOptions(statusSelect, 'upcoming')
    expect(onChange).toHaveBeenCalledWith({ game: '', type: '', status: 'upcoming' })
  })

  it('reflects current filter values in the selects', () => {
    render(<EventFilters filters={{ game: 'MW3', type: 'qualifier', status: 'completed' }} onChange={vi.fn()} />)
    const [gameSelect, typeSelect, statusSelect] = screen.getAllByRole('combobox') as HTMLSelectElement[]
    expect(gameSelect.value).toBe('MW3')
    expect(typeSelect.value).toBe('qualifier')
    expect(statusSelect.value).toBe('completed')
  })
})
