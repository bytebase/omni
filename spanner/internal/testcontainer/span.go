package testcontainer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// The harness's span requests (spanner/spanner-oracle/span.go): a
// fixture's DDL, or a query on the last fixture built.
const (
	spanFixturePrefix = "-- omni:span-fixture\n"
	spanQueryPrefix   = "-- omni:span-query\n"
)

// A SpanAnswer is the harness's answer to a span request.
type SpanAnswer struct {
	// Error is why there is no answer: the gRPC code of a query the
	// emulator refuses, or fixture or infra when the harness failed.
	Error   string `json:"error"`
	Message string `json:"message"`

	// Schemas is a fixture's stored schema, read as Bytebase's sync reads it.
	Schemas []SpanSchema `json:"schemas"`

	// Columns is a query's output column names, from its result set
	// metadata; an anonymous column is "".
	Columns []string `json:"columns"`
	// Tables is the tables and views the query reads, as schema.name or, in
	// the default schema, name: each one without which the query fails.
	Tables []string `json:"tables"`
}

// A SpanSchema is one schema of a fixture; the default schema is "".
type SpanSchema struct {
	Name   string         `json:"name"`
	Tables []SpanRelation `json:"tables"`
	Views  []SpanRelation `json:"views"`
}

// A SpanRelation is a table or a view.
type SpanRelation struct {
	Name       string       `json:"name"`
	Definition string       `json:"definition"`
	Columns    []SpanColumn `json:"columns"`
}

// A SpanColumn is one column of a table or view.
type SpanColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// SpanFixture builds a fixture from its DDL and returns its stored schema.
func (o *Oracle) SpanFixture(ddl string) (SpanAnswer, error) {
	a, err := o.span(spanFixturePrefix + ddl)
	if err == nil && a.Error != "" {
		err = fmt.Errorf("%s: %s", a.Error, a.Message)
	}
	return a, err
}

// Span returns a query's span on the last fixture built. A query the
// emulator refuses is an answer with Error set.
func (o *Oracle) Span(query string) (SpanAnswer, error) {
	a, err := o.span(spanQueryPrefix + query)
	if err == nil && (a.Error == "fixture" || a.Error == "infra") {
		err = fmt.Errorf("%s: %s", a.Error, a.Message)
	}
	return a, err
}

func (o *Oracle) span(req string) (SpanAnswer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	answer, err := o.ask(ctx, req)
	if err != nil {
		return SpanAnswer{}, err
	}
	var a SpanAnswer
	if err := json.Unmarshal([]byte(answer), &a); err != nil {
		return SpanAnswer{}, fmt.Errorf("decode the harness's answer %q: %w", answer, err)
	}
	return a, nil
}
