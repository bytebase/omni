package parser

import (
	"testing"
	"time"
)

// TestPredictWithListMakesProgress: the loop over PREDICT's WITH list
// skipped commas and otherwise called parsePredictColumnDef, which returned
// nil without consuming anything at a token that starts no column
// definition. A literal or a nested '(' in the list (WITH (1)) spun Parse
// forever. Completion stops parsing a data type once the cursor is behind
// it, which left the '(' of NVARCHAR(10) for the loop: Collect spun at 76 of
// the 94 cursor offsets of the first statement below, from the opening
// parenthesis through NVARCHAR.
func TestPredictWithListMakesProgress(t *testing.T) {
	within := func(what string, fn func()) {
		t.Helper()
		done := make(chan struct{})
		go func() {
			defer close(done)
			fn()
		}()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("%s: no result within 10s", what)
		}
	}
	for _, sql := range []string{
		"PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (Score FLOAT, Label NVARCHAR(10)) AS p",
		"PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (Score DECIMAL(10, 2) NOT NULL)",
		"PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (Label NVARCHAR(10) COLLATE Latin1_General_CI_AS NULL)",
		"PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (",
	} {
		for off := 0; off <= len(sql); off++ {
			within(sql[:off]+"|"+sql[off:], func() { Collect(sql, off) })
		}
	}
	for _, sql := range []string{
		"PREDICT (MODEL = @m, DATA = t AS d) WITH (1)",
		"PREDICT (MODEL = @m, DATA = t AS d) WITH (Score FLOAT, (x))",
		"PREDICT (MODEL = @m, DATA = t AS d) WITH (Score FLOAT, 'a')",
		"PREDICT (MODEL = @m, DATA = t AS d) WITH (Score)",
		"PREDICT (MODEL = @m, DATA = t AS d) WITH (Score FLOAT",
		"PREDICT (MODEL = @m, DATA = t AS d",
		"PREDICT MODEL = @m",
	} {
		var err error
		within(sql, func() { _, err = Parse(sql) })
		if err == nil {
			t.Errorf("Parse(%q) = nil error, want a syntax error", sql)
		}
	}
}
