import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import MatchDetail from './MatchDetail'
import {
  matchDetailFixture,
  matchDetailNoMapsFixture,
  matchDetailEmptyStatsFixture,
  matchDetailZeroDamageFixture,
} from '../test/fixtures/matchDetail'

vi.mock('../hooks/useApi', () => ({ useApi: vi.fn() }))

vi.mock('../utils/logoAssets', () => ({
  getTeamLogo: vi.fn().mockReturnValue(null),
}))

import { useApi } from '../hooks/useApi'
const mockUseApi = vi.mocked(useApi)

function renderMatch(id = '42') {
  return render(
    <MemoryRouter initialEntries={[`/matches/${id}`]}>
      <Routes>
        <Route path="/matches/:id" element={<MatchDetail />} />
      </Routes>
    </MemoryRouter>
  )
}

beforeEach(() => { mockUseApi.mockReset() })

// ── Contract shape ────────────────────────────────────────────────────────────

describe('MatchDetail — response contract rendering', () => {
  it('renders both team names', () => {
    mockUseApi.mockReturnValue({ data: matchDetailFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    // Team names appear in both the header and per-map scoreboards.
    expect(screen.getAllByText('OpTic Texas').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Atlanta FaZe').length).toBeGreaterThan(0)
  })

  it('renders team abbreviations', () => {
    mockUseApi.mockReturnValue({ data: matchDetailFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    expect(screen.getByText('OTX')).toBeInTheDocument()
    expect(screen.getByText('ATL')).toBeInTheDocument()
  })

  it('renders the series score', () => {
    mockUseApi.mockReturnValue({ data: matchDetailFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    // team1_score=3 and team2_score=1 both appear
    const threes = screen.getAllByText('3')
    expect(threes.length).toBeGreaterThan(0)
  })

  it('renders tournament name in breadcrumb', () => {
    mockUseApi.mockReturnValue({ data: matchDetailFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    expect(screen.getByText('CDL Major 1 2025')).toBeInTheDocument()
  })

  it('renders played map names', () => {
    mockUseApi.mockReturnValue({ data: matchDetailFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    expect(screen.getByText('Skyline')).toBeInTheDocument()
    expect(screen.getByText('Rewind')).toBeInTheDocument()
  })
})

// ── Empty states do not crash ─────────────────────────────────────────────────

describe('MatchDetail — empty states', () => {
  it('renders without crashing when maps array is empty', () => {
    mockUseApi.mockReturnValue({ data: matchDetailNoMapsFixture, loading: false, error: null, refetch: vi.fn() })
    expect(() => renderMatch('43')).not.toThrow()
    // Header still shows team names
    expect(screen.getByText('OpTic Texas')).toBeInTheDocument()
  })

  it('renders without crashing when player stats arrays are empty', () => {
    mockUseApi.mockReturnValue({ data: matchDetailEmptyStatsFixture, loading: false, error: null, refetch: vi.fn() })
    expect(() => renderMatch('44')).not.toThrow()
    expect(screen.getByText('Skyline')).toBeInTheDocument()
  })

  it('shows loading state', () => {
    mockUseApi.mockReturnValue({ data: null, loading: true, error: null, refetch: vi.fn() })
    renderMatch()
    expect(screen.getByText(/Loading match/i)).toBeInTheDocument()
  })

  it('shows "Match not found" on error', () => {
    mockUseApi.mockReturnValue({ data: null, loading: false, error: 'Not found', refetch: vi.fn() })
    renderMatch()
    expect(screen.getByText(/Match not found/i)).toBeInTheDocument()
  })

  it('shows "Match not found" when data is null and no error', () => {
    mockUseApi.mockReturnValue({ data: null, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    expect(screen.getByText(/Match not found/i)).toBeInTheDocument()
  })
})

// ── Player stats rendering ────────────────────────────────────────────────────

describe('MatchDetail — player stats in scoreboard', () => {
  it('renders player gamertags', () => {
    mockUseApi.mockReturnValue({ data: matchDetailFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    // Each player appears in both maps, so getAllByText is needed.
    expect(screen.getAllByText('Shotzzy').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Cellium').length).toBeGreaterThan(0)
  })

  it('renders "—" for damage when damage is 0', () => {
    mockUseApi.mockReturnValue({ data: matchDetailZeroDamageFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    const dashes = screen.getAllByText('—')
    expect(dashes.length).toBeGreaterThan(0)
  })

  it('renders "—" for bp_rating when bp_rating is 0', () => {
    mockUseApi.mockReturnValue({ data: matchDetailZeroDamageFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    // damage=0 and bp_rating=0 both render "—"
    const dashes = screen.getAllByText('—')
    expect(dashes.length).toBeGreaterThanOrEqual(2)
  })

  it('renders HP mode header column', () => {
    mockUseApi.mockReturnValue({ data: matchDetailFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    // "Hill" appears once per Hardpoint map scoreboard; use getAllByText when multiple HP maps are present.
    expect(screen.getAllByText('Hill').length).toBeGreaterThan(0)
  })

  it('renders SND mode header columns', () => {
    mockUseApi.mockReturnValue({ data: matchDetailFixture, loading: false, error: null, refetch: vi.fn() })
    renderMatch()
    // Column headers appear once per S&D map scoreboard; use getAllByText when multiple S&D maps are present.
    expect(screen.getAllByText('Plants').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Defuses').length).toBeGreaterThan(0)
    expect(screen.getAllByText('FB').length).toBeGreaterThan(0)
  })
})
