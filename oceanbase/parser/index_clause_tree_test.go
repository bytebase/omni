package parser

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// What the index_type clause leaves in the tree. A clause that parses into the
// wrong field is worse than one that is rejected: a CREATE INDEX whose index
// type never reaches the tree leaves a catalog that has forgotten it.
//
// The corpus in tidb/parser cannot see any of this, because it records only
// whether a statement parses.

// TestCreateIndexClausesInTree reads the index type back out of a CREATE INDEX
// statement, from both positions the grammar allows: in front of ON, and among
// the options after the column list.
func TestCreateIndexClausesInTree(t *testing.T) {
	for _, tc := range []struct {
		sql string
		// before is CreateIndexStmt.IndexType, the clause in front of ON.
		before string
		// options is the index type each USING option carries, in order.
		options []string
		// algorithm and lock are the clauses that close the statement.
		algorithm string
		lock      string
	}{
		{sql: "CREATE INDEX i ON t (a) USING BTREE", before: "", options: []string{"BTREE"}, algorithm: "", lock: ""},
		{sql: "CREATE INDEX i USING HASH ON t (a)", before: "HASH", options: nil, algorithm: "", lock: ""},
		{sql: "CREATE INDEX i USING BTREE ON t (a) USING HASH", before: "BTREE", options: []string{"HASH"}, algorithm: "", lock: ""},
		{sql: "CREATE INDEX i ON t (a) USING BTREE COMMENT 'c'", before: "", options: []string{"BTREE"}, algorithm: "", lock: ""},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneCreateIndex(t, tc.sql)
			if stmt.IndexType != tc.before {
				t.Errorf("IndexType = %q, want %q", stmt.IndexType, tc.before)
			}
			var got []string
			for _, opt := range stmt.Options {
				if opt.Name != "USING" {
					continue
				}
				s, ok := opt.Value.(*ast.StringLit)
				if !ok {
					t.Fatalf("USING option carries %T, want a string", opt.Value)
				}
				got = append(got, s.Value)
			}
			if len(got) != len(tc.options) {
				t.Fatalf("USING options = %q, want %q", got, tc.options)
			}
			for i := range got {
				if got[i] != tc.options[i] {
					t.Errorf("USING option %d = %q, want %q", i, got[i], tc.options[i])
				}
			}
			if stmt.Algorithm != tc.algorithm {
				t.Errorf("Algorithm = %q, want %q", stmt.Algorithm, tc.algorithm)
			}
			if stmt.Lock != tc.lock {
				t.Errorf("Lock = %q, want %q", stmt.Lock, tc.lock)
			}
		})
	}
}

// parseOneCreateIndex parses sql and returns the single CREATE INDEX it holds.
func parseOneCreateIndex(t *testing.T, sql string) *ast.CreateIndexStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) read %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("Parse(%q) read a %T, want a CreateIndexStmt", sql, list.Items[0])
	}
	return stmt
}

// TestDropIndexClausesInTree reads the name, the IF EXISTS and the closing
// clauses back out of a DROP INDEX statement. The name is the field at risk:
// a clause taken as the index name leaves a statement that parses and then
// drops the wrong index, or none.
func TestDropIndexClausesInTree(t *testing.T) {
	for _, tc := range []struct {
		sql       string
		name      string
		algorithm string
		lock      string
	}{
		{sql: "DROP INDEX idx ON t", name: "idx"},
		{sql: "DROP INDEX idx ON d.t", name: "idx"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.sql, err)
			}
			if len(list.Items) != 1 {
				t.Fatalf("Parse(%q) read %d statements, want 1", tc.sql, len(list.Items))
			}
			stmt, ok := list.Items[0].(*ast.DropIndexStmt)
			if !ok {
				t.Fatalf("Parse(%q) read a %T, want a DropIndexStmt", tc.sql, list.Items[0])
			}
			if stmt.Name != tc.name {
				t.Errorf("Name = %q, want %q", stmt.Name, tc.name)
			}
			if stmt.Algorithm != tc.algorithm {
				t.Errorf("Algorithm = %q, want %q", stmt.Algorithm, tc.algorithm)
			}
			if stmt.Lock != tc.lock {
				t.Errorf("Lock = %q, want %q", stmt.Lock, tc.lock)
			}
		})
	}
}
