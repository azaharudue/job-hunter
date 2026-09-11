# go-gate — the verification gate, rewritten in Go

This is the core of the agent in `agent/job-hunter.md`, compiled into a small
Go service. It answers one question per posting, exactly like the agent does:

> Is this a real, fresh, sector-allowed product-company posting — and how good
> a fit is it?

It is a deliberately small demonstration: one crate of domain logic (`gate.go`),
a swap-able store, a GraphQL API, and tests. No ORM, no config framework, no
crate soup.

## Why it exists

The Markdown agent is portable and readable, but it is a *spec*. This is the
same rules as *code*: typed, testable, and re-runnable without an LLM in the
loop. It demonstrates how the agent's prose decisions map onto a compilable
artifact — and how the "never trust the aggregator" rule becomes a hard
`SourceKind` check instead of a nudge in a prompt.

## Files

| File | Role |
|---|---|
| `gate.go` | `Config`, `Job`, `ScoreBreakdown`, `Gate.Verify`. All hard filters and scoring weights. |
| `store.go` | `MemoryStore` (postings). Swap point for Postgres. |
| `main.go` | `net/http` server, GraphQL schema + resolvers (`/graphql`, `/healthz`). |
| `gate_test.go` | Table-driven tests: aggregator, services, denied sector, stale posting, scoring ceiling. |

## Rules mapped to code

| Agent rule | Code |
|---|---|
| Never trust aggregator mirrors | `Verify` returns `false` unless `SourceKind == SourceBoard` |
| Product companies only | `DenyCompanyType` (service, agency, research) |
| Sector allow/deny list | `DenySectors` substring check on the description |
| Max 7 days old | `PostedAt` within `Config.MaxAge` |
| Scoring rubric | `Stack +40, Salary +20, Company +10, Region +15, Culture +5` (capped at 100) |
| "Don't invent facts" | unstated salary/location score as `0`/partial — never guessed |

## Run it

```bash
cd go-gate
go run .            # serves :8080 (override with PORT=...)
```

```graphql
# reject a mirror with a logged reason
mutation { verify(id: "X") { accepted reason report } }
# score a real posting
mutation { verify(id: "A") { accepted score { total stack salary } report } }
# list what's in the store
query { recent(limit: 5) { id title company salaryMin } }
```

Tests:

```bash
go test ./...
```

## Deploying to GCP (as-is, no changes needed)

The service has no external state in `MemoryStore`, so it deploys unchanged:

```bash
gcloud run deploy job-hunter-gate --source . --region europe-west1 --allow-unauthenticated
```

With `Dockerfile` and a Postgres swap of `MemoryStore` (same `Add`/`List`
contract) it becomes a multi-instance service. Go + Cloud Run gives cold
starts in ~1s and no capacity planning — the trade was the point: the gate
should be *cheap to run twice a day*, not a job for a whole team.

## What's intentionally missing

- No Postgres yet — `MemoryStore` holds the interface, a `postgresStore`
  implementing it is the next 30 lines.
- No auth — the API is localhost/private-network by default.
- No crawler — this crate *judges* postings; fetching is the agent's job
  (that's the whole "never invent facts" boundary).