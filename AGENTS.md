# Job Boards & Integration Notes

This registry lists job boards and portals the `job-hunter` agent searches. Respect rate limits; verify every lead on its source board before scoring.

## North America
- **Indeed** (indeed.com) — Global leader, primary integration target
- **LinkedIn Jobs** (linkedin.com/jobs)
- **Glassdoor** (glassdoor.com)
- **Dice** (dice.com) — Tech/IT focused
- **FlexJobs** (flexjobs.com) — Curated remote/flexible (paid)
- **BuiltIn** (builtin.com) — Tech startups
- **Wellfound** (wellfound.com) — Startups (formerly AngelList)
- **We Work Remotely** (weworkremotely.com) — Remote only
- **Remote OK** (remoteok.com) — Remote only
- **USAJobs** (usajobs.gov) — Federal government only

## Europe
- **StepStone** (stepstone.de)
- **XING** (xing.de)
- **kununu** (kununu.com/de/jobs) — DACH employer reviews + live job search (XING/onlyfy listings), salary data; category pages like kununu.com/de/jobs/d-it-softwareentwicklung. Kununu score + salary-check give company-size/culture signals directly.
- **TotalJobs** (totaljobs.com) — UK
- **Reed** (reed.co.uk) — UK
- **CV-Library** (cv-library.co.uk) — UK
- **Welcome to the Jungle** (welcome-to-the-jungle.com) — FR tech
- **Cadremploi** (cadremploi.fr) — FR executive
- **Finn.no** (fin.no) — NO
- **Job Bank** (jobbank.gc.ca) — CA gov

## Aggregators / boards used directly in the workflow
- **join.com** (join.com) — startup/scaleup product companies
- **jobijoba**, **jobborse24** — meta-search
- **Personio / Workable / Greenhouse company career pages** — direct product-company listings
- **freehire.me** (freehire.me) — Greenhouse-sourced aggregator, useful for discovering product companies running on Greenhouse
- **DataBerlin**, **RemoteRocketship** — curated tech boards

## Country → boards mapping (default: Germany)
```json
{
  "germany":  ["stepstone", "xing", "kununu", "indeed_de", "monster_de"],
  "usa":      ["indeed", "linkedin", "glassdoor", "dice", "monster"],
  "canada":   ["indeed", "linkedin", "workopolis", "eluta"],
  "uk":       ["totaljobs", "cv-library", "linkedin", "reed"],
  "australia":["seek", "indeed_au", "linkedin_au"],
  "france":   ["indeed_fr", "cadremploi", "welcome-to-the-jungle"],
  "netherlands":["stepstone_nl", "indeed_nl", "linkedin_nl"]
}
```

## Crawling priority
1. Indeed — highest volume, most consistent DOM
2. LinkedIn — professional quality, strong filtering
3. StepStone — European coverage
4. kununu — DACH reviews + live job search (XING/onlyfy)
5. We Work Remotely — remote-specific
6. Dice — tech talent pool

## Rate limiting & politeness
- 1 request per second per domain (`RATE_LIMIT_MS = 1000`)
- Respect `robots.txt` where applicable
- User-Agent: `job-hunter/1.0 (personal job tracking agent)`
- Max 3 retries with exponential backoff
- If a board blocks a fetch, log it and move on — no retry hammering