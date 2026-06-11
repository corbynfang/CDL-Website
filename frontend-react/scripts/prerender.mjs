import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const DIST = path.join(__dirname, '..', 'dist')
const API_BASE = process.env.PRERENDER_API_BASE ?? 'https://cdlytics.com'
const SITE_URL = 'https://cdlytics.com'

if (!fs.existsSync(DIST)) {
  console.error('[prerender] dist/ not found — run `bun run build` first')
  process.exit(1)
}

const template = fs.readFileSync(path.join(DIST, 'index.html'), 'utf8')

function injectMeta(html, { title, description, canonical, ogType = 'website', ldJson }) {
  const fullTitle = title ? `${title} | CDLytics` : 'CDLytics — CDL Stats & Analytics'
  const canonicalUrl = `${SITE_URL}${canonical}`

  const metaBlock = [
    `<title>${fullTitle}</title>`,
    `<meta name="description" content="${description}">`,
    `<link rel="canonical" href="${canonicalUrl}">`,
    `<meta property="og:title" content="${fullTitle}">`,
    `<meta property="og:description" content="${description}">`,
    `<meta property="og:url" content="${canonicalUrl}">`,
    `<meta property="og:type" content="${ogType}">`,
    `<meta name="twitter:title" content="${fullTitle}">`,
    `<meta name="twitter:description" content="${description}">`,
    ldJson ? `<script type="application/ld+json">${JSON.stringify(ldJson)}</script>` : '',
  ].filter(Boolean).join('\n    ')

  return html
    .replace(/<title>.*?<\/title>/, '')
    .replace(/<meta name="title"[^>]*>/g, '')
    .replace(/<meta name="description"[^>]*>/g, '')
    .replace(/<meta name="keywords"[^>]*>/g, '')
    .replace(/<meta property="og:[^"]*"[^>]*>/g, '')
    .replace(/<meta name="twitter:[^"]*"[^>]*>/g, '')
    .replace(/<link rel="canonical"[^>]*>/g, '')
    .replace('</head>', `    ${metaBlock}\n  </head>`)
}

function writeRoute(route, html) {
  const dir = path.join(DIST, route === '/' ? '.' : route)
  fs.mkdirSync(dir, { recursive: true })
  fs.writeFileSync(path.join(dir, 'index.html'), html)
}

async function fetchJson(url) {
  const res = await fetch(url)
  if (!res.ok) throw new Error(`${res.status} ${url}`)
  return res.json()
}

const FAQ_ITEMS = [
  { q: 'What is CDLytics?', a: 'CDLytics is an independent statistics database for the Call of Duty League (CDL). It tracks player K/D ratios, match results, tournament brackets, team rosters, and transfer history across every CDL season.' },
  { q: 'Which CDL seasons does CDLytics cover?', a: 'CDLytics covers all CDL seasons from launch through the current season, including Black Ops Cold War, Vanguard, Modern Warfare II, Modern Warfare III, and Black Ops 6.' },
  { q: 'How are K/D ratios calculated?', a: 'K/D ratios are calculated from official match data. Overall K/D is kills divided by deaths across all maps played. Mode-specific K/D (Hardpoint, Search & Destroy, Control) is calculated from maps of that mode only.' },
  { q: 'Who has the highest K/D ratio in CDL history?', a: 'K/D rankings change each season. Visit the Stats page to see the current leaderboard filtered by season.' },
  { q: 'Is CDLytics affiliated with Activision or the CDL?', a: 'No. CDLytics is an independent fan-built project and is not affiliated with Activision, the Call of Duty League, or any listed team.' },
]

const faqSchema = {
  '@context': 'https://schema.org',
  '@type': 'FAQPage',
  mainEntity: FAQ_ITEMS.map(({ q, a }) => ({
    '@type': 'Question',
    name: q,
    acceptedAnswer: { '@type': 'Answer', text: a },
  })),
}

const STATIC = [
  {
    route: '/',
    title: 'CDL Stats & Analytics',
    description: 'The independent source for Call of Duty League statistics. Search player K/D ratios, tournament brackets, team rosters, and transfer history across every CDL season.',
    ldJson: faqSchema,
  },
  {
    route: '/players',
    title: 'CDL Players',
    description: 'Browse every Call of Duty League player — gamertags, countries, roles, and active status. Click through to view individual K/D stats and match history.',
  },
  {
    route: '/teams',
    title: 'CDL Teams',
    description: 'Every Call of Duty League franchise and team — rosters, season history, and performance stats. Filter by CDL season.',
  },
  {
    route: '/events',
    title: 'CDL Events & Tournaments',
    description: 'All Call of Duty League tournaments — Majors, Championship events, and more. Browse brackets, results, and player stats for every CDL event.',
  },
  {
    route: '/stats',
    title: 'CDL K/D Rankings',
    description: 'Call of Duty League K/D leaderboards by season. See which players have the highest kill/death ratios across Hardpoint, Search & Destroy, and Control.',
  },
  {
    route: '/transfers',
    title: 'CDL Transfers & Roster Moves',
    description: 'Complete history of Call of Duty League player transfers, signings, and roster moves. Track which players moved between teams and when.',
  },
]

async function main() {
  let done = 0

  for (const page of STATIC) {
    const html = injectMeta(template, { ...page, canonical: page.route })
    writeRoute(page.route, html)
    done++
  }
  console.log(`[prerender] ${done} static routes done`)

  try {
    const [players, teams, tournaments] = await Promise.all([
      fetchJson(`${API_BASE}/api/v1/players?limit=500&page=1`),
      fetchJson(`${API_BASE}/api/v1/teams`),
      fetchJson(`${API_BASE}/api/v1/tournaments`),
    ])

    for (const player of players.data ?? []) {
      const desc = `${player.gamertag} CDL statistics — match history, event stats, and career K/D breakdown.`
      const ldJson = {
        '@context': 'https://schema.org',
        '@type': 'Person',
        name: player.gamertag,
        url: `${SITE_URL}/players/${player.id}`,
      }
      const html = injectMeta(template, {
        title: `${player.gamertag} CDL Stats`,
        description: desc,
        canonical: `/players/${player.id}`,
        ogType: 'profile',
        ldJson,
      })
      writeRoute(`/players/${player.id}`, html)
      done++
    }

    for (const team of (teams ?? []).filter(t => t.is_cdl_franchise)) {
      const desc = `${team.name} CDL roster, match history, and season stats on CDLytics.`
      const html = injectMeta(template, {
        title: `${team.name} CDL Stats`,
        description: desc,
        canonical: `/teams/${team.id}`,
      })
      writeRoute(`/teams/${team.id}`, html)
      done++
    }

    for (const event of tournaments ?? []) {
      const year = event.start_date ? new Date(event.start_date).getFullYear() : ''
      const desc = `${event.name}${year ? ` (${year})` : ''} — CDL tournament results, bracket, match scores, and player stats.`
      const ldJson = {
        '@context': 'https://schema.org',
        '@type': 'SportsEvent',
        name: event.name,
        startDate: event.start_date,
        url: `${SITE_URL}/events/${event.slug}`,
        sport: 'Call of Duty',
        organizer: { '@type': 'Organization', name: 'Call of Duty League' },
      }
      const html = injectMeta(template, {
        title: `${event.name} Results & Stats`,
        description: desc,
        canonical: `/events/${event.slug}`,
        ogType: 'article',
        ldJson,
      })
      writeRoute(`/events/${event.slug}`, html)
      done++
    }

    console.log(`[prerender] ${done} total routes written`)
  } catch (e) {
    console.warn(`[prerender] API unreachable — only static routes written:`, e.message)
    console.log(`[prerender] ${done} routes written (static only)`)
  }
}

main().catch(e => { console.error(e); process.exit(1) })
