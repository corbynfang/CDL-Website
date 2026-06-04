# Backend Architecture (`internal/`)

This is a Go backend over a PostgreSQL database hosted on Supabase. It serves a
read-only JSON API (Gin) for the CDL stats site.

> **Supabase is just hosted Postgres here.** There is no Supabase SDK, no
> PostgREST, no Supabase client library. The app connects with a plain Postgres
> connection string through Supabase's connection **pooler** and talks SQL via
> GORM. Anything you'd do against a normal Postgres, you do here.

## The four layers

Every layer has exactly one job. Do not collapse or skip them.

```
internal/handlers/    HTTP only — parse request, call service, return JSON   ("controllers")
internal/services/    Business logic — validation, coordination, assembly
internal/store/       Database queries — GORM / raw SQL, nothing else
internal/models/      Struct definitions only — no methods, no DB code
internal/database/    config.go — DB connection + AutoMigrate (the one place DB is opened)
internal/middleware/  Gin middleware — headers, rate limiting
```

Dependencies point **down only**:

```
handlers ──> services ──> store ──> gorm ──> Supabase Postgres
   │             │           │
   └─ models ────┴───────────┘   (every layer may read model structs)
```

Enforced invariants (all currently true — keep them true):

- A handler file calls **only its own service** (`tournaments.go` → `h.tournaments`).
- A service uses **only its own injected store(s)**; it never reaches into another domain's store.
- Stores never import services; services never import handlers.
- Handlers never touch a store or `gorm` directly — the **only** exception is
  `handlers.New()`, which is the wiring point (the composition root).
- Stores use the **injected** `*gorm.DB` (`s.db`), never the global `database.DB`.

## Domains

Each domain is one handler file ↔ one service ↔ one store, named the same:

| Handler (`handlers/`) | Service (`services/`)   | Store (`store/`)   |
| --------------------- | ----------------------- | ------------------ |
| `seasons.go`          | `SeasonService`         | `SeasonStore`      |
| `teams.go`            | `TeamService`           | `TeamStore` (+ `SeasonStore`) |
| `players.go`          | `PlayerService`         | `PlayerStore`      |
| `franchises.go`       | `FranchiseService`      | `FranchiseStore`   |
| `matches.go`          | `MatchService`          | `MatchStore`       |
| `tournaments.go`      | `TournamentService`     | `TournamentStore`  |
| `transfers.go`        | `TransferService`       | `TransferStore`    |
| `stats.go`            | `StatsService`          | `StatsStore`       |

`handlers.go` holds the shared base: the `Handler` struct, the `New()`
constructor, and HTTP helpers (`validateID`, `parsePagination`, `noCacheHeaders`).

**One deliberate asymmetry:** the standard service shape is *one store, no
cache*. `TeamService` is the exception — it injects a second store (`seasons`,
for resolving the latest season when building rosters) and keeps an in-memory
cache of the team list (hot, rarely changes). If you're learning the pattern,
treat every other service as the template and `TeamService` as the justified
special case.

## Request → data flow

Startup (`cmd/main.go`):

```
database.ConnectDatabase()      reads DATABASE_URL, gorm.Open(postgres, …), tunes the pool
database.AutoMigrate()          creates/updates tables, pg_trgm extension, indexes
h := handlers.New(database.DB)  builds every store + service from the one *gorm.DB
handlers.RegisterRoutes(api, h) binds routes to handler methods
```

Per request:

```
gin route → Handler method → Service → Store → gorm → Supabase Postgres
 (route)     (parse params)   (logic)  (SQL)   (pool)  (over the pooler)
```

`database.DB` is a **single shared connection pool**. `New()` hands that same
pool to all eight stores; each store wraps it (`gormSeasonStore{ db }`) and is
the only place SQL lives.

### Worked example — `GET /api/tournaments/5/bracket`

1. `handlers/tournaments.go` `GetTournamentBracket` — parses `id=5` via
   `validateID`, opens a 15s context.
2. → `h.tournaments.AssembleBracket(ctx, 5)` (`services/tournaments.go`) — the
   business logic: fetch tournament + matches, detect bracket format, normalize
   round keys, bucket matches into bracket / group stage.
3. → `ts.tournaments.GetByID` / `GetBracketMatches` (`store/tournament.go`) —
   the SQL, via the injected GORM pool.
4. → GORM → Supabase pooler → Postgres → rows flow back up → JSON out.

Each layer does one job: handler parses HTTP, service decides *what* and *how to
assemble*, store knows *how to fetch*.

## Query styles in the store layer

Stores mix two approaches by query complexity — both run through the same GORM
connection to the same Postgres; raw SQL is not a separate data path:

- **GORM ORM** for simple CRUD (`season.go`):
  ```go
  s.db.WithContext(ctx).Order("start_date DESC").Find(&seasons)
  s.db.WithContext(ctx).Where("is_active = ?", true).First(&season)
  ```
- **Raw SQL** (`.Raw(...).Scan(...)`) for joins/aggregations/era gating
  (`team.go`, `tournament.go`, `player.go`):
  ```go
  s.db.WithContext(ctx).Raw(`SELECT ... FROM teams t WHERE EXISTS (...)`, args).Scan(&out)
  ```

## Database connection

- `internal/database/config.go` is the **only** place the DB is opened.
- Connection string comes from `DATABASE_URL` (loaded from `.env` /
  `.env.railway`); it points at Supabase's pooler host
  (`*.pooler.supabase.com`).
- Pool tuning: `MaxIdleConns(5)`, `MaxOpenConns(25)`, `ConnMaxLifetime(1h)`,
  with a startup ping and ret/backoff on connect.
- `AutoMigrate()` keeps the schema in sync with `internal/models/` and creates
  the `pg_trgm` extension plus search/lookup indexes.

## API surface

Routes are registered in `handlers/routes.go` under the `/api` group (see
`cmd/main.go`). Everything is `GET` — this API is read-only.

```
/seasons            /seasons/:id            /seasons/active
/teams              /teams/:id              /teams/:id/players      /teams/:id/stats
/players            /players/:id            /players/:id/stats      /players/:id/kd
/players/:id/matches  /players/:id/franchise-career  /players/top-kd
/stats/all-kd-by-tournament
/matches/:id
/franchises         /franchises/:key
/tournaments        /tournaments/slug/:slug /tournaments/:id        /tournaments/:id/bracket
/tournaments/:id/matches  /tournaments/:id/teams  /tournaments/:id/stats
/transfers
```

## Conventions that keep this consistent

- Models live in `internal/models/`. Use `models.Match`, `models.Team`, etc.
  There is no `database.Match` — that pattern was removed.
- Handlers are methods on `*Handler`, never package-level functions.
- When adding a domain: create `models` structs → a `Store` interface +
  `gormXStore` → a `Service` that injects the store → a handler file of
  `*Handler` methods → wire it in `handlers.New()` → register routes. Keep the
  four names aligned (handler file, service, store, test file).
- Quick compile check: `go build ./...`.
