import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import BracketControls from './BracketControls'

describe('BracketControls', () => {
  it('always renders an "All Rounds" button', () => {
    render(<BracketControls rounds={[]} active={null} onSelect={vi.fn()} />)
    expect(screen.getByRole('button', { name: /all rounds/i })).toBeInTheDocument()
  })

  it('renders a button for each round in sorted order', () => {
    render(<BracketControls rounds={['grand_finals', 'winners_r1']} active={null} onSelect={vi.fn()} />)
    expect(screen.getByRole('button', { name: /winners round 1/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /grand finals/i })).toBeInTheDocument()
  })

  it('renders with an empty rounds array without crashing', () => {
    render(<BracketControls rounds={[]} active={null} onSelect={vi.fn()} />)
    expect(screen.getByRole('button', { name: /all rounds/i })).toBeInTheDocument()
  })

  it('calls onSelect with null when "All Rounds" is clicked', async () => {
    const onSelect = vi.fn()
    render(<BracketControls rounds={['winners_r1']} active="winners_r1" onSelect={onSelect} />)
    await userEvent.click(screen.getByRole('button', { name: /all rounds/i }))
    expect(onSelect).toHaveBeenCalledWith(null)
  })

  it('calls onSelect with the round key when a round button is clicked', async () => {
    const onSelect = vi.fn()
    render(<BracketControls rounds={['winners_r1', 'grand_finals']} active={null} onSelect={onSelect} />)
    await userEvent.click(screen.getByRole('button', { name: /grand finals/i }))
    expect(onSelect).toHaveBeenCalledWith('grand_finals')
  })

  it('formats unknown round keys into human-readable labels', () => {
    render(<BracketControls rounds={['custom_round']} active={null} onSelect={vi.fn()} />)
    expect(screen.getByRole('button', { name: /custom round/i })).toBeInTheDocument()
  })

  it('renders all CDL bracket round types', () => {
    const rounds = ['winners_r1', 'winners_r2', 'winners_finals', 'elim_r1', 'elim_finals', 'grand_finals']
    render(<BracketControls rounds={rounds} active={null} onSelect={vi.fn()} />)
    expect(screen.getByRole('button', { name: /winners round 1/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /winners round 2/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /grand finals/i })).toBeInTheDocument()
  })
})
