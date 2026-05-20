import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import PlayerDetail from './PlayerDetail'
import {
  playerFixture,
  playerKDFixture,
  playerKDNoControlFixture,
  playerMatchesFixture,
  playerMatchesEmptyFixture,
  playerMatchesNullKDFixture,
} from '../test/fixtures/playerDetail'

// Mock useApi — each test controls what each URL returns.
vi.mock('../hooks/useApi', () => ({ useApi: vi.fn() }))

// Mock asset helpers — tests don't need real image files on disk.
vi.mock('../utils/avatarAssets', () => ({
  getPlayerAvatar: vi.fn().mockReturnValue('/placeholder.png'),
}))

import { useApi } from '../hooks/useApi'
const mockUseApi = vi.mocked(useApi)

const careerDefault = { player_id: 1, gamertag: 'Shotzzy', franchises: [] }

function setupMocks({
  player = playerFixture,
  stats  = playerKDFixture,
  matches = playerMatchesFixture,
  career  = careerDefault,
} = {}) {
  mockUseApi.mockImplementation((url: string) => {
    if (url.includes('/kd'))               return { data: stats,   loading: false, error: null, refetch: vi.fn() }
    if (url.includes('/matches'))          return { data: matches, loading: false, error: null, refetch: vi.fn() }
    if (url.includes('/franchise-career')) return { data: career,  loading: false, error: null, refetch: vi.fn() }
    return                                        { data: player,  loading: false, error: null, refetch: vi.fn() }
  })
}

function renderPlayer(id = '1') {
  return render(
    <MemoryRouter initialEntries={[`/players/${id}`]}>
      <Routes>
        <Route path="/players/:id" element={<PlayerDetail />} />
      </Routes>
    </MemoryRouter>
  )
}

beforeEach(() => { mockUseApi.mockReset() })

// ── Hero section ──────────────────────────────────────────────────────────────

describe('PlayerDetail — hero section', () => {
  it('renders the player gamertag', () => {
    setupMocks()
    renderPlayer()
    expect(screen.getByText('Shotzzy')).toBeInTheDocument()
  })

  it('renders country when present', () => {
    setupMocks()
    renderPlayer()
    expect(screen.getByText('US')).toBeInTheDocument()
  })

  it('renders Active badge for active players', () => {
    setupMocks()
    renderPlayer()
    expect(screen.getByText('Active')).toBeInTheDocument()
  })

  it('shows loading state while data is fetching', () => {
    mockUseApi.mockReturnValue({ data: null, loading: true, error: null, refetch: vi.fn() })
    renderPlayer()
    expect(screen.getByText(/Loading player data/i)).toBeInTheDocument()
  })

  it('shows not found when player is null and error is set', () => {
    mockUseApi.mockReturnValue({ data: null, loading: false, error: 'Not found', refetch: vi.fn() })
    renderPlayer()
    expect(screen.getByText(/Player not found/i)).toBeInTheDocument()
  })
})

// ── K/D stats section ─────────────────────────────────────────────────────────

describe('PlayerDetail — K/D stats section', () => {
  it('renders all four mode labels', () => {
    setupMocks()
    renderPlayer()
    expect(screen.getByText('Overall')).toBeInTheDocument()
    expect(screen.getByText('Hardpoint')).toBeInTheDocument()
    expect(screen.getByText('Search & Destroy')).toBeInTheDocument()
    expect(screen.getByText('Control')).toBeInTheDocument()
  })

  it('renders avg_kd formatted to 2 decimal places', () => {
    setupMocks()
    renderPlayer()
    expect(screen.getByText('1.32')).toBeInTheDocument()
  })

  it('renders total kills count', () => {
    setupMocks()
    renderPlayer()
    expect(screen.getByText('500')).toBeInTheDocument()
  })
})

// ── Last 5 Matches ────────────────────────────────────────────────────────────

describe('PlayerDetail — Last 5 Matches tab', () => {
  it('renders at most 5 match rows when more than 5 exist', () => {
    setupMocks() // fixture has 6 total matches across 2 events
    renderPlayer()
    // Each rendered match row shows "vs <abbr>" — count those
    const rows = screen.getAllByText(/^vs /)
    expect(rows.length).toBeLessThanOrEqual(5)
  })

  it('renders exactly 5 rows when 6 matches are available', () => {
    setupMocks()
    renderPlayer()
    const rows = screen.getAllByText(/^vs /)
    expect(rows.length).toBe(5)
  })

  it('shows "No matches available" when events array is empty', () => {
    setupMocks({ matches: playerMatchesEmptyFixture })
    renderPlayer()
    expect(screen.getByText('No matches available')).toBeInTheDocument()
  })

  it('renders W for winning matches', () => {
    setupMocks()
    renderPlayer()
    // result column shows just the first character ("W" or "L")
    const wLabels = screen.getAllByText('W')
    expect(wLabels.length).toBeGreaterThan(0)
  })

  it('renders L for losing matches', () => {
    setupMocks()
    renderPlayer()
    const lLabels = screen.getAllByText('L')
    expect(lLabels.length).toBeGreaterThan(0)
  })
})

// ── K/D null handling (documents bugs) ───────────────────────────────────────

describe('PlayerDetail — null K/D display (current behavior)', () => {
  it('renders "0.00" for null kd instead of "—" (known bug)', () => {
    // Bug: PlayerDetail.tsx line 247 — null kd falls through to "0.00" not "—"
    setupMocks({ matches: playerMatchesNullKDFixture })
    renderPlayer()
    const kdValues = screen.getAllByText('0.00')
    expect(kdValues.length).toBeGreaterThan(0)
  })
})

// ── Control stats with zero value ─────────────────────────────────────────────

describe('PlayerDetail — Control K/D when no Control maps played', () => {
  it('renders 0.00 when backend returns control_kd_ratio = 0', () => {
    // Backend contract: returns 0 (not null) when ctlMapsTotal == 0.
    // Component renders 0.00 — there is no "—" for missing Control data.
    setupMocks({ stats: playerKDNoControlFixture })
    renderPlayer()
    const values = screen.getAllByText('0.00')
    expect(values.length).toBeGreaterThan(0)
  })
})
