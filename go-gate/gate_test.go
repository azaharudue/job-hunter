package main

import (
	"testing"
	"time"
)

func TestVerifyRejectsAggregatorMirror(t *testing.T) {
	g := NewGate(testConfig())
	j := SampleJob("1", "LLM Engineer")
	j.SourceKind = SourceAggregator
	if ok, _ := g.Verify(j); ok {
		t.Fatal("aggregator mirror must be rejected")
	}
}

func TestVerifyRejectsServicesCompany(t *testing.T) {
	g := NewGate(testConfig())
	j := SampleJob("2", "Applied AI Engineer")
	j.CompanyType = CompanyService
	if ok, _ := g.Verify(j); ok {
		t.Fatal("services / consultancy must be rejected")
	}
}

func TestVerifyRejectsDeniedSector(t *testing.T) {
	g := NewGate(testConfig())
	j := SampleJob("3", "Data Scientist")
	j.Description = "work on our credit risk algorithms for a bank"
	if ok, _ := g.Verify(j); ok {
		t.Fatal("denied sector (banking/credit) must be rejected")
	}
}

func TestVerifyRejectsStalePosting(t *testing.T) {
	g := NewGate(testConfig())
	j := SampleJob("4", "Backend Engineer")
	j.PostedAt = time.Now().Add(-30 * 24 * time.Hour)
	if ok, _ := g.Verify(j); ok {
		t.Fatal("posting older than MaxAge must be rejected")
	}
}

func TestVerifyAcceptsGoodPostingAndScoresIt(t *testing.T) {
	g := NewGate(testConfig())
	j := SampleJob("5", "AI Fullstack Developer")
	j.SalaryMin = 90000
	j.Description = "product work, remote"
	ok, b := g.Verify(j)
	if !ok {
		t.Fatalf("good product posting should pass; got %+v", b)
	}
	if b.Total < 60 {
		t.Fatalf("expected a strong score, got %d (%+v)", b.Total, b)
	}
	// sanity: 100 is the ceiling
	if b.Total > 100 {
		t.Fatalf("score over 100: %d", b.Total)
	}
}

func TestReportLine(t *testing.T) {
	g := NewGate(testConfig())
	j := SampleJob("6", "RAG Engineer")
	ok, b := g.Verify(j)
	line := g.Report(j, b, ok)
	want := "stepstone | ok | Example GmbH |"
	if len(line) < len(want) || line[:len(want)] != want {
		t.Fatalf("unexpected report line: %q", line)
	}
}

func testConfig() Config {
	return Config{
		SalaryMin:       70000,
		Region:          "germany",
		WorkModel:       "remote",
		MaxAge:          7 * 24 * time.Hour,
		DenySectors:     []string{"banking", "lending", "credit", "insurance", "gambling"},
		DenyCompanyType: []CompanyType{CompanyService, CompanyAgency, CompanyResearch},
		PrimaryStack:    []string{"Python", "TypeScript", "React", "FastAPI", "Postgres"},
	}
}
