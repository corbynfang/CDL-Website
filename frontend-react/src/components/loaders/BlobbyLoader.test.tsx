import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import BlobbyLoader from './BlobbyLoader'

describe('BlobbyLoader', () => {
  it('renders the default label', () => {
    render(<BlobbyLoader />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('renders a custom label', () => {
    render(<BlobbyLoader label="Loading event..." />)
    expect(screen.getByText('Loading event...')).toBeInTheDocument()
  })

  it('does not crash with an empty label', () => {
    render(<BlobbyLoader label="" />)
    // Empty label still renders the blob div without error
    expect(document.querySelector('[style]')).toBeInTheDocument()
  })

  it('renders one blob div and one label paragraph', () => {
    const { container } = render(<BlobbyLoader label="Please wait" />)
    expect(container.querySelector('[style]')).toBeInTheDocument()
    expect(screen.getByText('Please wait')).toBeInTheDocument()
  })
})
