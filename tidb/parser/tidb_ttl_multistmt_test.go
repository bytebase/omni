package parser

import (
	"strings"
	"testing"

	ast "github.com/bytebase/omni/tidb/ast"
)

// TestParseTTLInMultiStatement pins a parser panic observed when TTL
// appears in any segment except the first of a multi-statement input.
//
// Bug: parseTableOption for TTL indexed `p.lexer.input[exprStart:exprEnd]`
// directly. Loc values include the absolute baseOffset added by
// parseSingle for non-leading segments, but `lexer.input` is the
// segment-local text. When baseOffset > 0 the slice bounds exceed the
// segment length and the runtime panics:
//
//	runtime error: slice bounds out of range [:105] with length 80
//
// The fix replaces the direct slice with p.inputText(start, end), which
// subtracts baseOffset before slicing — the same helper every other
// body-capture site in the parser uses.
func TestParseTTLInMultiStatement(t *testing.T) {
	sql := `CREATE DATABASE dummy; CREATE TABLE t (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR;`

	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	// Find the CREATE TABLE statement and verify the TTL option captured
	// the expression text (not garbage / truncated).
	var ttlValue string
	for _, stmt := range list.Items {
		ct, ok := stmt.(*ast.CreateTableStmt)
		if !ok {
			continue
		}
		for _, opt := range ct.Options {
			if opt.Name == "TTL" {
				ttlValue = opt.Value
			}
		}
	}
	if ttlValue == "" {
		t.Fatal("TTL option not captured on CREATE TABLE in multi-statement input")
	}
	// The captured text should begin with the column reference and end
	// with the unit; we don't assert exact whitespace because inputText
	// trims leading/trailing ws via TrimSpace elsewhere.
	if !strings.Contains(ttlValue, "created_at") || !strings.Contains(strings.ToUpper(ttlValue), "INTERVAL 1 YEAR") {
		t.Errorf("TTL value captured incorrectly: %q", ttlValue)
	}
}
