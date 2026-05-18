-- seed-tournament-metadata.sql
-- Run this ONCE in the Supabase SQL editor AFTER deploying the Go code that adds
-- the is_lan, country, and logo_url columns to the tournaments table.
--
-- Updates are slug-based so they're safe to re-run (idempotent).
-- All qualifiers, minors, season_summary, and unknown types are set to online (is_lan=false).

-- ── Mark all qualifiers / minors / artifacts as online ───────────────────────
UPDATE tournaments SET is_lan = false
WHERE tournament_type IN ('qualifier', 'minor_tournament', 'season_summary', 'unknown');

-- ── Black Ops Cold War 2021 ───────────────────────────────────────────────────
-- Stages 1-3 were online-only
UPDATE tournaments SET is_lan = false
WHERE slug IN (
    'cdl-major-1-tournament-2021',
    'cdl-major-2-tournament-2021',
    'cdl-major-3-tournament-2021'
);

-- Stage 4 Major — first LAN since March 2020
UPDATE tournaments
SET location = 'Dallas, Texas', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-4-tournament-2021';

-- Stage 5 Major
UPDATE tournaments
SET location = 'Arlington, Texas', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-5-tournament-2021';

-- CDL Championship 2021
UPDATE tournaments
SET location = 'Los Angeles, California', country = 'USA', is_lan = true, prize_pool = 2500000
WHERE slug = 'cdl-championship-2021';

-- ── Vanguard 2022 ─────────────────────────────────────────────────────────────
UPDATE tournaments
SET location = 'Arlington, Texas', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-1-tournament-2022';

UPDATE tournaments
SET location = 'Prior Lake, Minnesota', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-2-tournament-2022';

UPDATE tournaments
SET location = 'Toronto, Ontario', country = 'CAN', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-3-tournament-2022';

UPDATE tournaments
SET location = 'Brooklyn, New York', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-4-tournament-2022';

UPDATE tournaments
SET location = 'Los Angeles, California', country = 'USA', is_lan = true, prize_pool = 2550000
WHERE slug = 'cdl-championship-2022';

-- ── Modern Warfare II 2023 ────────────────────────────────────────────────────
-- Kickoff event (LAN, same location as Major 1)
UPDATE tournaments
SET location = 'Raleigh, North Carolina', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-1-kickoff-2023';

UPDATE tournaments
SET location = 'Raleigh, North Carolina', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-1-tournament-2023';

UPDATE tournaments
SET location = 'Boston, Massachusetts', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-2-tournament-2023';

UPDATE tournaments
SET location = 'Arlington, Texas', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-3-tournament-2023';

UPDATE tournaments
SET location = 'Columbus, Ohio', country = 'USA', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-4-tournament-2023';

UPDATE tournaments
SET location = 'Toronto, Ontario', country = 'CAN', is_lan = true, prize_pool = 500000
WHERE slug = 'cdl-major-5-tournament-2023';

UPDATE tournaments
SET location = 'Las Vegas, Nevada', country = 'USA', is_lan = true, prize_pool = 2380000
WHERE slug = 'cdl-championship-2023';

-- ── Modern Warfare III 2024 ───────────────────────────────────────────────────
UPDATE tournaments
SET location = 'Boston, Massachusetts', country = 'USA', is_lan = true, prize_pool = 375000
WHERE slug = 'cdl-major-1-tournament-2024';

UPDATE tournaments
SET location = 'Fort Lauderdale, Florida', country = 'USA', is_lan = true, prize_pool = 375000
WHERE slug = 'cdl-major-2-tournament-2024';

UPDATE tournaments
SET location = 'Toronto, Ontario', country = 'CAN', is_lan = true, prize_pool = 375000
WHERE slug = 'cdl-major-3-tournament-2024';

UPDATE tournaments
SET location = 'Burbank, California', country = 'USA', is_lan = true, prize_pool = 375000
WHERE slug = 'cdl-major-4-tournament-2024';

UPDATE tournaments
SET location = 'Allen, Texas', country = 'USA', is_lan = true, prize_pool = 2000000
WHERE slug = 'cdl-league-championship-2024';

-- ── Black Ops 6 2025 ─────────────────────────────────────────────────────────
UPDATE tournaments
SET location = 'Madrid', country = 'ESP', is_lan = true, prize_pool = 375000
WHERE slug = 'cdl-major-1-tournament-2025';

UPDATE tournaments
SET location = 'Allen, Texas', country = 'USA', is_lan = true, prize_pool = 375000
WHERE slug = 'cdl-major-2-tournament-2025';

UPDATE tournaments
SET location = 'Boca Raton, Florida', country = 'USA', is_lan = true, prize_pool = 375000
WHERE slug = 'cdl-major-3-tournament-2025';

UPDATE tournaments
SET location = 'Dallas, Texas', country = 'USA', is_lan = true, prize_pool = 375000
WHERE slug = 'cdl-major-4-tournament-2025';

UPDATE tournaments
SET location = 'Kitchener, Ontario', country = 'CAN', is_lan = true, prize_pool = 2000000
WHERE slug = 'cdl-league-championship-2025';

-- ── EWC ──────────────────────────────────────────────────────────────────────
UPDATE tournaments
SET location = 'Riyadh', country = 'SAU', is_lan = true, prize_pool = 1800000
WHERE slug = 'esports-world-cup-2025';

UPDATE tournaments
SET location = 'Riyadh', country = 'SAU', is_lan = true
WHERE slug = 'esports-world-cup-2024';

-- Verify the update
SELECT slug, name, location, country, is_lan, prize_pool
FROM tournaments
WHERE tournament_type NOT IN ('season_summary', 'unknown', 'qualifier', 'minor_tournament')
ORDER BY start_date DESC;
