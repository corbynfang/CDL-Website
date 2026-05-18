interface Filters {
  game: string
  type: string
  status: string
}

interface Props {
  filters: Filters
  onChange: (f: Filters) => void
}

const GAME_OPTIONS = [
  { value: '', label: 'All Games' },
  { value: 'BO6', label: 'Black Ops 6' },
  { value: 'MW3', label: 'Modern Warfare III' },
  { value: 'MW2', label: 'Modern Warfare II' },
  { value: 'VG', label: 'Vanguard' },
  { value: 'CW', label: 'Cold War' },
]

const TYPE_OPTIONS = [
  { value: '', label: 'All Types' },
  { value: 'major_tournament', label: 'Major' },
  { value: 'championship', label: 'Championship' },
  { value: 'international_major', label: 'International' },
  { value: 'kickoff', label: 'Kickoff' },
  { value: 'qualifier', label: 'Qualifier' },
]

const STATUS_OPTIONS = [
  { value: '', label: 'All' },
  { value: 'upcoming', label: 'Upcoming' },
  { value: 'live', label: 'Live' },
  { value: 'completed', label: 'Completed' },
]

const SELECT_CLASS =
  'bg-[#111111] border border-[#1a1a1a] px-3 py-2 text-xs text-[#a3a3a3] focus:outline-none focus:border-[#2a2a2a] uppercase tracking-wider'

export default function EventFilters({ filters, onChange }: Props) {
  return (
    <div className="flex flex-wrap gap-2">
      <select
        value={filters.game}
        onChange={e => onChange({ ...filters, game: e.target.value })}
        className={SELECT_CLASS}
      >
        {GAME_OPTIONS.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
      </select>

      <select
        value={filters.type}
        onChange={e => onChange({ ...filters, type: e.target.value })}
        className={SELECT_CLASS}
      >
        {TYPE_OPTIONS.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
      </select>

      <select
        value={filters.status}
        onChange={e => onChange({ ...filters, status: e.target.value })}
        className={SELECT_CLASS}
      >
        {STATUS_OPTIONS.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
      </select>
    </div>
  )
}
