# ✅ job-hunter — a production AI-agent pipeline for job hunting

**A self-contained agent that finds, verifies, scores, and packages job leads — run it once a day, or on demand.**

`job-hunter` is a real, deployed AI agent built on agent-native workflows. Every day it crawls job boards and company career pages, **proves each posting is real** (not an aggregator fake), scores it against a candidate profile (0–100), and emits a browsable HTML briefing you can open in any browser. It has run headlessly on a schedule, rejected hundreds of fakes, and produced a ranked, verifiable shortlist — collecting **zero fabricated facts** along the way.

It is designed to be *one agent in a fleet*: deterministic, self-verifying, measurable, and cheap to extend.

---

## Why this exists (the problem)

Most AI job-hunt tools dump unvetted postings. Aggregator boards mirror foreign roles with fake local locations and invented salary bands; consultancies and staffing firms flood "AI Engineer" searches; and half the fresh listings are 3-month-old reposts.

`job-hunter` treats job hunting as an **engineering pipeline with a hard verification gate**, not a search query:

```
search  →  verify (source board, country, salary, posting date, company type)  →  score  →  brief
```

---

## Architecture

```
job-hunter/agent/job-hunter.md      # the agent: role targeting, hard filters, scoring rubric
job-hunter/profiles/PROFILE.md      # candidate profile (skills, region, salary, sector allow-list)
job-hunter/AGENTS.md                # board registry (Indeed, LinkedIn, StepStone, kununu, …)
job-hunter/joblists/                # output: DAILY-JOB-BRIEFING-*.html + DAILY-JOB-HUNT-LOG.md
```

| Component | Responsibility |
|---|---|
| **Role tiering** | Tier 1 (Applied/GenAI/LLM/RAG/AI-Product Engineer) weighted over Tier 2 (Full-Stack/Backend AI); consultancies, Data-Sci-only and Staff/Principal roles excluded |
| **Hard filters** | Geography, work model (remote-first), salary floor, **product-companies only**, sector allow-list (no banking/insurance/riba), seniority mix 70/20/10, freshness < 7 days |
| **Verification gate** | Every lead is re-opened on its **source board** (not the aggregator). Country, salary band, posting date, deadline and company type confirmed before scoring; rejections logged with reasons |
| **Scoring rubric** | +40 stack/tier match · +20 salary · +15 location/work-model · +15 AI relevance · +10 company fit |
| **Deliverable** | Self-contained `DAILY-JOB-BRIEFING-YYYY-MM-DD.html` — inline CSS, summary table of the 30 most recent verified leads, per-job fit/gap cards, verification notes, relocation section |
| **Telemetry** | One-line daily append to `DAILY-JOB-HUNT-LOG.md`: boards searched, leads verified, top lead + score, fakes rejected |

### Anti-fake measures (what "production-grade" means here)
- Aggregator mirrors are never trusted — the source board page is fetched and checked.
- Recruiters are allowed only when the hiring company is a **named, product-based** client ("via \<agency\>" noted).
- Every posting is scored against the **profile only**; no invented skills, salaries, company sizes, or culture claims.
- Blocked boards are logged and skipped — no retry hammering, no rate-limit abuse.

---

## Deployment (already proven)

Runs three ways, same agent:

1. **Scheduled headless** — a cron / Windows Task Scheduler job invokes `opencode run --agent job-hunter --auto "Run your daily job-hunt routine now."` → zero-interaction daily briefings.
2. **On-demand CLI** — `/job-hunt` inside opencode, Claude Code, or Codex boots the full routine now.
3. **Agent-to-agent** — because it is a plain markdown agent spec, it is *tool-agnostic*: the same file loads in opencode agents, Claude Code skills/commands, and Codex `AGENTS.md`.

---

## Install

### opencode
```
mkdir -p .opencode/agent .opencode/command
cp agent/job-hunter.md .opencode/agent/job-hunter.md
```

Optionally add a slash command:

```markdown
# .opencode/command/job-hunt.md
---
description: Run the job-hunter agent now for fresh job updates. Usage: /job-hunt
agent: job-hunter
---
Run your daily job-hunt routine now: pull the latest postings, verify them on source
boards, score against the profile, and produce the DAILY-JOB-BRIEFING HTML. Report
what changed since the last run.
```

### Claude Code
```
mkdir -p .claude/commands
cp agent/job-hunter.md .claude/skills/job-hunter/SKILL.md
# or the body as .claude/commands/job-hunt.md for a slash command
```

### Codex (OpenAI)
The agent instruction body drops into your project's `AGENTS.md` — codex reads it
per-run. Headless scheduling works the same way:

```
codex exec "Run your daily job-hunt routine now."
```

---

## Configure

Edit `profiles/PROFILE.md` — one file controls targeting:

| Key | Example |
|---|---|
| `skills` | Python, FastAPI, React, TypeScript, Azure, Docker, Kubernetes, PostgreSQL, LLM/RAG/agentic |
| `region` | Germany, then EU |
| `work_model` | 100% remote, max 50% office |
| `salary_min` | €70,000–€90,000 |
| `sector_allow` | software products, healthcare, climate/energy, compliance, education |
| `seniority_mix` | 70% mid · 20% senior · 10% exceptional entry |

Boards are registered in `AGENTS.md` (kununu included for company-size/culture/salary signals).

---

## Repository layout

```
job-hunter-agent/
├── README.md              ← you are here
├── agent/
│   └── job-hunter.md      ← canonical agent spec (the deliverable)
├── profiles/
│   └── PROFILE.md         ← candidate profile (one-file targeting)
├── AGENTS.md              ← job-board registry + crawling notes
└── scripts/
    ├── job-hunt-daily.sh  # cron/UNIX runner
    └── job-hunt-daily.bat # Windows Task Scheduler runner
```

## License

MIT — use it, fork it, teach your own agent fleet to do its job.

---

*Built by an engineer who believed a job search should be as reproducible as a CI pipeline.*