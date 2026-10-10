package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/example/fin-domain/fin"
	"os"
)

func main() {
	input := flag.String("input", "examples/events.json", "offline JSON event dataset")
	flag.Parse()
	f, err := os.Open(*input)
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	var events []fin.Event
	if err := dec.Decode(&events); err != nil {
		fatal(err)
	}
	report, err := fin.Analyze(events)
	if err != nil {
		fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		fatal(err)
	}
}
func fatal(err error) { fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1) }
