import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { getKdColorClass } from './kdUtils'

// ── getKdColorClass unit tests ─────────────────────────────────────────────────

describe('getKdColorClass', () => {
  it.each([
    [null,      'text-[#a3a3a3]'],
    [undefined, 'text-[#a3a3a3]'],
    [1.0,       'text-[#a3a3a3]'],
    [0.99,      'text-red-400'],
    [0,         'text-red-400'],
    [1.01,      'text-green-400'],
    [2.0,       'text-green-400'],
  ] as const)('getKdColorClass(%s) → %s', (kd, expected) => {
    expect(getKdColorClass(kd)).toBe(expected)
  })
})

// ── K/D color rendering tests ─────────────────────────────────────────────────

function KdValue({ kd }: { kd: number | null }) {
  return (
    <span data-testid="kd" className={getKdColorClass(kd)}>
      {kd?.toFixed(2) ?? '—'}
    </span>
  )
}

describe('K/D color rendering', () => {
  it('applies a green class when kd > 1', () => {
    render(<KdValue kd={1.5} />)
    expect(screen.getByTestId('kd').className).toContain('green')
  })

  it('applies a red class when kd < 1', () => {
    render(<KdValue kd={0.75} />)
    expect(screen.getByTestId('kd').className).toContain('red')
  })

  it('applies neither green nor red when kd === 1', () => {
    render(<KdValue kd={1.0} />)
    const el = screen.getByTestId('kd')
    expect(el.className).not.toContain('green')
    expect(el.className).not.toContain('red')
  })

  it('applies neither green nor red when kd is null', () => {
    render(<KdValue kd={null} />)
    const el = screen.getByTestId('kd')
    expect(el.className).not.toContain('green')
    expect(el.className).not.toContain('red')
  })
})
