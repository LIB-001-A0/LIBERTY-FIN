package fin

import "sort"

type Edge struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Status  string `json:"status"`
	EventID string `json:"event_id"`
}

// BuildGraph records observed connections, including refused and blocked ones.
// No scanning, tunneling, or active network calls are performed.
func BuildGraph(events []Event) ([]string, []Edge) {
	nodes := map[string]bool{}
	edges := []Edge{}
	for _, e := range events {
		if e.Source == "" || e.Target == "" {
			continue
		}
		nodes[e.Source] = true
		nodes[e.Target] = true
		status := "observed"
		if e.Kind == "connection_refused" {
			status = "refused"
		} else if !e.Allowed {
			status = "blocked"
		}
		edges = append(edges, Edge{From: e.Source, To: e.Target, Status: status, EventID: e.ID})
	}
	out := make([]string, 0, len(nodes))
	for n := range nodes {
		out = append(out, n)
	}
	sort.Strings(out)
	return out, edges
}
