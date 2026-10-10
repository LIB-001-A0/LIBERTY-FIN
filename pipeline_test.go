package fin

import (
	"testing"
	"time"
)

func TestAnalyze(t *testing.T) {
	events := []Event{
		{ID: "1", Time: time.Unix(1, 0), Kind: "session_open", SessionID: "s", Source: "api", Target: "auth", Allowed: true},
		{ID: "2", Time: time.Unix(2, 0), Kind: "connection_refused", Source: "gateway", Target: "internal", Allowed: false},
		{ID: "3", Time: time.Unix(3, 0), Kind: "session_revoke", SessionID: "s", Allowed: true},
	}
	r, err := Analyze(events)
	if err != nil {
		t.Fatal(err)
	}
	if r.Sessions["s"] != "revoked" {
		t.Fatalf("session: %v", r.Sessions)
	}
	if len(r.Edges) != 2 || r.Edges[1].Status != "refused" {
		t.Fatalf("edges: %+v", r.Edges)
	}
	if r.Decisions[1].Action != "review_block" {
		t.Fatalf("decision: %+v", r.Decisions[1])
	}
}
func TestDuplicateID(t *testing.T) {
	_, err := Analyze([]Event{{ID: "x"}, {ID: "x"}})
	if err == nil {
		t.Fatal("expected duplicate ID error")
	}
}
func TestPromptInjection(t *testing.T) {
	d := Guard(Event{ID: "p", Kind: "prompt_injection", Allowed: false})
	if d.Action != "quarantine_review" {
		t.Fatalf("unexpected action: %s", d.Action)
	}
}
