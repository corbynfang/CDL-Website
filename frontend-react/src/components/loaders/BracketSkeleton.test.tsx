import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/react'
import BracketSkeleton from './BracketSkeleton'

describe('BracketSkeleton', () => {
  it('renders without crashing', () => {
    const { container } = render(<BracketSkeleton />)
    expect(container.firstChild).toBeInTheDocument()
  })

  it('renders 8 skeleton match cards', () => {
    const { container } = render(<BracketSkeleton />)
    // Each card is a div with rounded-xl class
    const cards = container.querySelectorAll('.rounded-xl')
    expect(cards.length).toBe(8)
  })

  it('has the animate-pulse class for loading shimmer', () => {
    const { container } = render(<BracketSkeleton />)
    expect(container.firstChild).toHaveClass('animate-pulse')
  })
})
