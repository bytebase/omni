package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/pg/plpgsql/ast"
)

// parseOK is a test helper that verifies a PL/pgSQL body parses without error.
func parseOK(t *testing.T, body string) *ast.PLBlock {
	t.Helper()
	p := newParser(body)
	result, err := p.parse()
	if err != nil {
		t.Fatalf("Parse(%q): %v", body, err)
	}
	return result
}

// parseErr is a test helper that verifies parsing fails with an error containing wantContains.
func parseErr(t *testing.T, body string, wantContains string) {
	t.Helper()
	p := newParser(body)
	_, err := p.parse()
	if err == nil {
		t.Fatalf("Parse(%q): expected error containing %q, got nil", body, wantContains)
	}
	if !strings.Contains(err.Error(), wantContains) {
		t.Fatalf("Parse(%q): error %q does not contain %q", body, err.Error(), wantContains)
	}
}
