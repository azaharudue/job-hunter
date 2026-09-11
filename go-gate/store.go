package main

import (
	"fmt"
	"time"
)

// MemoryStore keeps a small list of recently seen postings. The production
// variant swaps this for Postgres (see README).
type MemoryStore struct {
	jobs []Job
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{} }

func (m *MemoryStore) Add(j Job) {
	m.jobs = append([]Job{j}, m.jobs...) // newest first
}

func (m *MemoryStore) List() []Job { return m.jobs }

func (m *MemoryStore) String() string {
	return fmt.Sprintf("%d postings", len(m.jobs))
}

// SampleJob builds a posting the way the crawler would after pulling a board
// page and extracting meta tags. In production this comes from the real board.
func SampleJob(id, title string) Job {
	return Job{
		ID:          id,
		Title:       title,
		Company:     "Example GmbH",
		CompanyType: CompanyProduct,
		Board:       "stepstone",
		SourceKind:  SourceBoard,
		Location:    "Berlin",
		PostedAt:    time.Now().Add(-6 * time.Hour),
		SalaryMin:   80000,
		SalaryMax:   95000,
		Skills:      []string{"Python", "FastAPI", "React", "TypeScript"},
		Description: "product work, remote-friendly, no investors",
	}
}
