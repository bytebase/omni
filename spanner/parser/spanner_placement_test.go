package parser

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/spanner/ast"
)

// Structural tests for CREATE/ALTER/DROP PLACEMENT and CREATE/ALTER MODEL's
// OPTIONS.
// Accept/reject parity with the pinned emulator is in the suite:spanner-ddl
// rows of testdata/compat.tsv.

func TestCreatePlacement(t *testing.T) {
	sql := "CREATE PLACEMENT IF NOT EXISTS us_placement OPTIONS (instance_partition = 'us-partition', default_leader = 'us-east1')"
	n := parseDDL(t, sql)
	s, ok := n.(*ast.CreatePlacementStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.CreatePlacementStmt", n)
	}
	if !s.IfNotExists || s.Name.String() != "us_placement" || len(s.Options) != 2 {
		t.Errorf("got IfNotExists=%v Name=%q opts=%d, want true/us_placement/2", s.IfNotExists, s.Name.String(), len(s.Options))
	}
	if s.Loc.Start != 0 || s.Loc.End != len(sql) {
		t.Errorf("Loc = %+v, want {0 %d}", s.Loc, len(sql))
	}
}

func TestCreatePlacementQuotedProto(t *testing.T) {
	// A quoted 'PROTO' is a string, unlike options_entry's bare PROTO keyword.
	parseDDL(t, "CREATE PLACEMENT p OPTIONS (instance_partition = 'PROTO')")
}

func TestDropPlacement(t *testing.T) {
	sql := "DROP PLACEMENT IF EXISTS us_placement"
	n := parseDDL(t, sql)
	d, ok := n.(*ast.DropPlacementStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.DropPlacementStmt", n)
	}
	if !d.IfExists || d.Name.String() != "us_placement" {
		t.Errorf("got IfExists=%v Name=%q, want true/us_placement", d.IfExists, d.Name.String())
	}
	if d.Loc.Start != 0 || d.Loc.End != len(sql) {
		t.Errorf("Loc = %+v, want {0 %d}", d.Loc, len(sql))
	}
}

func TestCreateModelRemoteOptions(t *testing.T) {
	sql := "CREATE MODEL m INPUT (a INT64) OUTPUT (b INT64) REMOTE OPTIONS (endpoints = ['//a', '//b'], default_batch_size = -1)"
	n := parseDDL(t, sql)
	m, ok := n.(*ast.CreateModelStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.CreateModelStmt", n)
	}
	if !m.Remote || len(m.Options) != 2 || m.Options[0].Name != "endpoints" {
		t.Errorf("got Remote=%v opts=%d, want true/2 starting with endpoints", m.Remote, len(m.Options))
	}
	if m.Loc.End != len(sql) {
		t.Errorf("Loc.End = %d, want %d", m.Loc.End, len(sql))
	}
}

// TestSpannerOptionsRefused holds the options Spanner's DDL parser refuses,
// each with the error at the option it names.
func TestSpannerOptionsRefused(t *testing.T) {
	for _, tc := range []struct {
		sql, at string
	}{
		{"CREATE MODEL m REMOTE OPTIONS (x = 1)", "x = 1"},
		{"CREATE MODEL m REMOTE OPTIONS (endpoint = '//a', endpoint = '//b')", "endpoint = '//b'"},
		{"CREATE MODEL m REMOTE OPTIONS (endpoint = 1)", "endpoint = 1"},
		{"CREATE MODEL m REMOTE OPTIONS (endpoints = [])", "endpoints = []"},
		{"CREATE MODEL m REMOTE OPTIONS (default_batch_size = 'x')", "default_batch_size = 'x'"},
		{"CREATE MODEL m REMOTE OPTIONS (endpoint += '//a')", "endpoint += '//a'"},
		{"CREATE PLACEMENT p OPTIONS (instance_partition = 1)", "instance_partition = 1"},
		{"CREATE PLACEMENT p OPTIONS (instance_partition = 'x', x = 1)", "x = 1"},
		{"CREATE PLACEMENT p OPTIONS (instance_partition = PROTO)", "instance_partition = PROTO"},
		{"CREATE MODEL m REMOTE OPTIONS (endpoint = proto)", "endpoint = proto"},
	} {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) = %v, want a *ParseError", tc.sql, err)
			continue
		}
		if got := tc.sql[pe.Position:pe.End]; got != tc.at {
			t.Errorf("Parse(%q) error at %q, want %q", tc.sql, got, tc.at)
		}
	}
	for _, sql := range []string{
		"CREATE PLACEMENT p",
		"CREATE PLACEMENT p OPTIONS (default_leader = 'us-east1')",
		"CREATE PLACEMENT p OPTIONS (instance_partition = NULL)",
		"CREATE OR REPLACE PLACEMENT p OPTIONS (instance_partition = 'x')",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted, want an error", sql)
		}
	}
}

func TestAlterPlacement(t *testing.T) {
	sql := "ALTER PLACEMENT IF EXISTS p SET OPTIONS (default_leader = 'us-east1')"
	n := parseDDL(t, sql)
	a, ok := n.(*ast.AlterPlacementStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.AlterPlacementStmt", n)
	}
	if !a.IfExists || a.Name.String() != "p" || len(a.SetOptions) != 1 {
		t.Errorf("got IfExists=%v Name=%q opts=%d, want true/p/1", a.IfExists, a.Name.String(), len(a.SetOptions))
	}
	if a.Loc.Start != 0 || a.Loc.End != len(sql) {
		t.Errorf("Loc = %+v, want {0 %d}", a.Loc, len(sql))
	}
	// SET OPTIONS is optional, and instance_partition is not required.
	if b := parseDDL(t, "ALTER PLACEMENT p").(*ast.AlterPlacementStmt); len(b.SetOptions) != 0 {
		t.Errorf("bare ALTER PLACEMENT SetOptions = %d, want 0", len(b.SetOptions))
	}
}

func TestAlterModel(t *testing.T) {
	sql := "ALTER MODEL TextEmbedder SET OPTIONS (default_batch_size = 5)"
	n := parseDDL(t, sql)
	a, ok := n.(*ast.AlterModelStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.AlterModelStmt", n)
	}
	if a.IfExists || a.Name.String() != "TextEmbedder" || len(a.SetOptions) != 1 {
		t.Errorf("got IfExists=%v Name=%q opts=%d, want false/TextEmbedder/1", a.IfExists, a.Name.String(), len(a.SetOptions))
	}
	if a.Loc.Start != 0 || a.Loc.End != len(sql) {
		t.Errorf("Loc = %+v, want {0 %d}", a.Loc, len(sql))
	}
}

func TestAlterModelAndPlacementRefused(t *testing.T) {
	for _, sql := range []string{
		"ALTER MODEL m SET OPTIONS (x = 1)",
		"ALTER MODEL m SET OPTIONS (endpoint = PROTO)",
		"ALTER MODEL m SET (endpoint = '//x')",
		"ALTER MODEL m RENAME TO n",
		"ALTER PLACEMENT p SET OPTIONS ()",
		"ALTER PLACEMENT p SET OPTIONS (default_leader = 1)",
		"ALTER PLACEMENT p OPTIONS (default_leader = 'x')",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) accepted, want an error", sql)
		}
	}
}
