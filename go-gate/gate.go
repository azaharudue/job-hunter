package main

import (
	"fmt"
	"strings"
	"time"
)

// CompanyType models the head-filter in agent/job-hunter.md: product companies
// only, no consultancies/agencies/science labs.
type CompanyType string

const (
	CompanyProduct  CompanyType = "product"
	CompanyService  CompanyType = "service"
	CompanyAgency   CompanyType = "agency"
	CompanyStartup  CompanyType = "startup"  // startup but product-focused
	CompanyResearch CompanyType = "research" // research/teaching institution
)

// SourceKind distinguishes the real board from an aggregator mirror. The gate
// only ever trusts SourceBoard URLs.
type SourceKind string

const (
	SourceBoard      SourceKind = "board"
	SourceAggregator SourceKind = "aggregator"
)

// Job is one verified posting as extracted from a source board.
type Job struct {
	ID          string
	URL         string
	Title       string
	Company     string
	CompanyType CompanyType
	Location    string
	Board       string // source board short id, e.g. stepstone
	SourceKind  SourceKind
	Description string
	PostedAt    time.Time
	Deadline    *time.Time
	SalaryMin   int
	SalaryMax   int
	Skills      []string
	IsGerman    bool
}

// Config is the scoring + filtering profile. Mirrors the keys in
// profiles/PROFILE.md (parameters, not hardcoded facts).
type Config struct {
	SalaryMin       int
	Region          string
	WorkModel       string
	MaxAge          time.Duration
	DenySectors     []string
	DenyCompanyType []CompanyType
	AllowSectors    []string
	SeniorityMix    string
	PrimaryStack    []string
}

// ScoreBreakdown is why a job scored its number. We keep it visible so the
// user can override a heuristic with judgment.
type ScoreBreakdown struct {
	Stack   int
	Salary  int
	Company int
	Region  int
	Culture int
	Total   int
}

// Gate is the verification + scoring engine. Stateless by design: every call
// filters through the same rules, so re-runs are idempotent.
type Gate struct {
	cfg Config
}

func NewGate(cfg Config) *Gate { return &Gate{cfg: cfg} }

// Verify decides whether a posting passes the hard filters, and if so how
// good a fit it is. The bool indicates "acceptable"; the breakdown explains
// the decision either way.
func (g *Gate) Verify(j Job) (bool, ScoreBreakdown) {
	var b ScoreBreakdown

	if j.SourceKind != SourceBoard {
		return false, b
	}
	if j.CompanyType != "" && g.deniedCompany(j.CompanyType) {
		return false, b
	}
	if g.deniedSector(j.Description) {
		return false, b
	}
	if j.PostedAt.IsZero() || time.Since(j.PostedAt) > g.cfg.MaxAge {
		return false, b
	}

	b.Stack = g.stackScore(j.Skills)
	b.Salary = g.salaryScore(j.SalaryMin, j.SalaryMax)
	b.Company = g.companyScore(j.CompanyType)
	b.Region = g.regionScore(j.Location)
	b.Culture = g.cultureScore(j.Description)

	if b.Stack == 0 && b.Salary == 0 && b.Company == 0 {
		// truly nothing matched: still return accepted=false to avoid an
		// empty zero-justification card in the brief
		return false, b
	}

	total := b.Stack + b.Salary + b.Company + b.Region + b.Culture
	if total > 100 {
		total = 100
	}
	b.Total = total
	return true, b
}

func (g *Gate) deniedCompany(t CompanyType) bool {
	for _, d := range g.cfg.DenyCompanyType {
		if d == t {
			return true
		}
	}
	return false
}

func (g *Gate) deniedSector(text string) bool {
	lower := strings.ToLower(text)
	for _, d := range g.cfg.DenySectors {
		if d == "" {
			continue
		}
		if strings.Contains(lower, strings.ToLower(d)) {
			return true
		}
	}
	return false
}

func (g *Gate) stackScore(skills []string) int {
	if len(g.cfg.PrimaryStack) == 0 {
		return 40
	}
	score := 0
	for _, s := range skills {
		for _, p := range g.cfg.PrimaryStack {
			if strings.EqualFold(s, p) {
				score += 20
			}
		}
	}
	if score > 40 {
		return 40
	}
	return score
}

func (g *Gate) salaryScore(min, max int) int {
	if min >= g.cfg.SalaryMin {
		return 20
	}
	if max >= g.cfg.SalaryMin {
		return 15
	}
	if max > 0 && max >= g.cfg.SalaryMin*7/10 {
		return 10
	}
	return 5
}

func (g *Gate) companyScore(t CompanyType) int {
	switch t {
	case CompanyProduct, CompanyStartup:
		return 10
	case CompanyService:
		return 5
	default:
		return 0
	}
}

func (g *Gate) regionScore(loc string) int {
	if g.cfg.Region == "" {
		return 15
	}
	if strings.Contains(strings.ToLower(loc), strings.ToLower(g.cfg.Region)) {
		return 15
	}
	if strings.Contains(strings.ToLower(loc), "remote") || strings.Contains(strings.ToLower(loc), "home") {
		return 10
	}
	return 5
}

func (g *Gate) cultureScore(desc string) int {
	lower := strings.ToLower(desc)
	score := 0
	if strings.Contains(lower, "remote") || strings.Contains(lower, "flexible") {
		score += 5
	}
	if strings.Contains(lower, "angel") || strings.Contains(lower, "known") {
		// a hint of "no investors / famous" culture
		score += 5
	}
	return score
}

// Report is the one-line telemetry entry appended to DAILY-JOB-HUNT-LOG.md.
func (g *Gate) Report(j Job, b ScoreBreakdown, accepted bool) string {
	status := "ok"
	if !accepted {
		status = "rejected"
	}
	return fmt.Sprintf("%s | %s | %s | %d/100 | %s %s", j.Board, status, j.Company, b.Total, j.Title, j.URL)
}
