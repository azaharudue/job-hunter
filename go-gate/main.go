package main

import (
	"log"
	"net/http"
	"os"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

const schema = `
	schema {
		query: Query
		mutation: Mutation
	}

	type Query {
		recent(limit: Int = 10): [Job!]!
	}

	type Mutation {
		# verify() replays a posting through the gate and returns the verdict.
		# This is the "never trust the aggregator" check, in code.
		verify(id: ID!): Verdict!
	}

	type Job {
		id: ID!
		title: String!
		company: String!
		board: String!
		location: String!
		salaryMin: Int!
		salaryMax: Int!
		skills: [String!]!
		ageHours: Int!
	}

	type Verdict {
		accepted: Boolean!
		reason: String!
		score: Score!
		report: String!
	}

	type Score {
		stack: Int!
		salary: Int!
		company: Int!
		region: Int!
		culture: Int!
		total: Int!
	}
`

type jobResolver struct{ j Job }

func (r *jobResolver) ID() graphql.ID   { return graphql.ID(r.j.ID) }
func (r *jobResolver) Title() string    { return r.j.Title }
func (r *jobResolver) Company() string  { return r.j.Company }
func (r *jobResolver) Board() string    { return r.j.Board }
func (r *jobResolver) Location() string { return r.j.Location }
func (r *jobResolver) SalaryMin() int32 { return int32(r.j.SalaryMin) }
func (r *jobResolver) SalaryMax() int32 { return int32(r.j.SalaryMax) }
func (r *jobResolver) Skills() []string { return r.j.Skills }
func (r *jobResolver) AgeHours() int32  { return int32(time.Since(r.j.PostedAt).Hours()) }

type verdictResolver struct {
	accepted bool
	reason   string
	b        ScoreBreakdown
	report   string
}

func (r *verdictResolver) Accepted() bool { return r.accepted }
func (r *verdictResolver) Reason() string { return r.reason }
func (r *verdictResolver) Score() *scoreResolver {
	return &scoreResolver{r.b}
}
func (r *verdictResolver) Report() string { return r.report }

type scoreResolver struct{ b ScoreBreakdown }

func (r *scoreResolver) Stack() int32   { return int32(r.b.Stack) }
func (r *scoreResolver) Salary() int32  { return int32(r.b.Salary) }
func (r *scoreResolver) Company() int32 { return int32(r.b.Company) }
func (r *scoreResolver) Region() int32  { return int32(r.b.Region) }
func (r *scoreResolver) Culture() int32 { return int32(r.b.Culture) }
func (r *scoreResolver) Total() int32   { return int32(r.b.Total) }

type root struct {
	store *MemoryStore
	gate  *Gate
}

func (r *root) Recent(args struct{ Limit int32 }) []*jobResolver {
	out := []*jobResolver{}
	limit := int(args.Limit)
	if limit <= 0 || limit > len(r.store.List()) {
		limit = len(r.store.List())
	}
	for _, j := range r.store.List()[:limit] {
		out = append(out, &jobResolver{j})
	}
	return out
}

func (r *root) Verify(args struct{ ID graphql.ID }) (*verdictResolver, error) {
	id := string(args.ID)
	var found *Job
	for i := range r.store.List() {
		if r.store.List()[i].ID == id {
			found = &r.store.List()[i]
			break
		}
	}
	if found == nil {
		return nil, &notFoundError{id}
	}
	accepted, b := r.gate.Verify(*found)
	reason := "accepted"
	if !accepted {
		reason = "rejected by gate"
	}
	report := r.gate.Report(*found, b, accepted)
	return &verdictResolver{accepted, reason, b, report}, nil
}

type notFoundError struct{ id string }

func (e *notFoundError) Error() string { return "job not found: " + e.id }

func main() {
	cfg := Config{
		SalaryMin:       70000,
		Region:          "germany",
		WorkModel:       "remote",
		MaxAge:          7 * 24 * time.Hour,
		DenySectors:     []string{"banking", "lending", "credit", "insurance", "gambling"},
		DenyCompanyType: []CompanyType{CompanyService, CompanyAgency, CompanyResearch},
		PrimaryStack:    []string{"Python", "TypeScript", "React", "FastAPI", "Postgres"},
	}

	store := NewMemoryStore()
	for i := 0; i < 6; i++ {
		j := SampleJob(string(rune('A'+i)), []string{
			"LLM Engineer",
			"AI Fullstack Developer",
			"Backend Engineer",
			"Applied AI Engineer",
			"RAG Engineer",
			"DevOps Engineer",
		}[i])
		store.Add(j)
	}
	// one rejected-with-reason example (aggregator mirror)
	store.Add(Job{
		ID:         "X",
		Title:      "AI Engineer",
		Company:    "Fake Aggregator Ltd",
		Board:      "aggregator",
		SourceKind: SourceAggregator,
		Location:   "Berlin",
		PostedAt:   time.Now().Add(-2 * time.Hour),
		SalaryMin:  90000,
		Skills:     []string{"Python"},
	})

	gate := NewGate(cfg)
	root := &root{store, gate}

	s := graphql.MustParseSchema(schema, root)
	http.Handle("/graphql", &relay.Handler{Schema: s})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	log.Printf("job-hunter gate listening on :%s (GraphQL at /graphql)", addr)
	log.Fatal(http.ListenAndServe(":"+addr, nil))
}
