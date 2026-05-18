import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/react'
import MatchCardSkeleton from './MatchCardSkeleton'

describe('MatchCardSkeleton', () => {
  it('renders without crashing', () => {
    const { container } = render(<MatchCardSkeleton />)
    expect(container.firstChild).toBeInTheDocument()
  })

  it('has the animate-pulse class for loading shimmer', () => {
    const { container } = render(<MatchCardSkeleton />)
    expect(container.firstChild).toHaveClass('animate-pulse')
  })

  it('renders two team avatar placeholder circles', () => {
    const { container } = render(<MatchCardSkeleton />)
    const circles = container.querySelectorAll('.rounded-full')
    expect(circles.length).toBe(2)
  })
})
