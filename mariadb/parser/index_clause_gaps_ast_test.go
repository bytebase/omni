package parser

import (
	"testing"

	ast "github.com/bytebase/omni/mariadb/ast"
)

// TestIndexLockWaitAST: the wait modifier is kept, with its literal as written, and the
// clauses around it still land where a consumer reads them.
func TestIndexLockWaitAST(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		noWait  bool
		hasWait bool
		seconds string
	}{
		{"CREATE INDEX i ON t (a) WAIT 3 COMMENT 'c' ALGORITHM=INPLACE", false, true, "3"},
		{"CREATE INDEX i ON t (a) WAIT 1.5", false, true, "1.5"},
		{"CREATE INDEX i ON t (a) WAIT 0x10", false, true, "0x10"},
		{"CREATE UNIQUE INDEX i ON t (a) NOWAIT", true, false, ""},
		{"CREATE VECTOR INDEX i ON t (a) WAIT 1e3 M=6", false, true, "1e3"},
		{"CREATE INDEX i ON t (a)", false, false, ""},
	} {
		stmt := oneNode(t, tc.sql).(*ast.CreateIndexStmt)
		if stmt.NoWait != tc.noWait || stmt.HasWait != tc.hasWait || stmt.WaitSeconds != tc.seconds {
			t.Errorf("%q: NoWait=%v HasWait=%v WaitSeconds=%q", tc.sql, stmt.NoWait, stmt.HasWait, stmt.WaitSeconds)
		}
		if stmt.Loc.End != len(tc.sql) {
			t.Errorf("%q: Loc.End=%d, want %d", tc.sql, stmt.Loc.End, len(tc.sql))
		}
	}
	full := oneNode(t, "CREATE INDEX i ON t (a) WAIT 3 COMMENT 'c' ALGORITHM=INPLACE").(*ast.CreateIndexStmt)
	if len(full.Options) != 1 || full.Options[0].Name != "COMMENT" || full.Algorithm != "INPLACE" {
		t.Errorf("Options=%+v Algorithm=%q", full.Options, full.Algorithm)
	}

	drop := oneNode(t, "DROP INDEX IF EXISTS i ON t WAIT 1.5").(*ast.DropIndexStmt)
	if !drop.IfExists || !drop.HasWait || drop.WaitSeconds != "1.5" || drop.NoWait || drop.Table.Name != "t" {
		t.Errorf("DROP INDEX: %+v", drop)
	}
	if drop := oneNode(t, "DROP INDEX i ON t NOWAIT").(*ast.DropIndexStmt); !drop.NoWait || drop.HasWait {
		t.Errorf("DROP INDEX NOWAIT: %+v", drop)
	}
	if drop := oneNode(t, "DROP INDEX i ON t").(*ast.DropIndexStmt); drop.NoWait || drop.HasWait {
		t.Errorf("DROP INDEX: %+v", drop)
	}
}

// TestAddIfNotExistsAST: every kind of key keeps the clause and its own name, and the
// clause is read on no key that did not write it.
func TestAddIfNotExistsAST(t *testing.T) {
	for _, tc := range []struct {
		sql       string
		kind      ast.ConstraintType
		name      string
		indexName string
	}{
		{"ALTER TABLE t ADD INDEX IF NOT EXISTS n (a)", ast.ConstrIndex, "n", ""},
		{"ALTER TABLE t ADD KEY IF NOT EXISTS (a)", ast.ConstrIndex, "", ""},
		{"ALTER TABLE t ADD UNIQUE KEY IF NOT EXISTS n USING BTREE (a)", ast.ConstrUnique, "n", ""},
		{"ALTER TABLE t ADD CONSTRAINT c UNIQUE IF NOT EXISTS (a)", ast.ConstrUnique, "c", ""},
		{"ALTER TABLE t ADD CONSTRAINT c UNIQUE IF NOT EXISTS n (a)", ast.ConstrUnique, "n", ""},
		{"ALTER TABLE t ADD FULLTEXT INDEX IF NOT EXISTS n (a)", ast.ConstrFulltextIndex, "n", ""},
		{"ALTER TABLE t ADD SPATIAL IF NOT EXISTS n (a)", ast.ConstrSpatialIndex, "n", ""},
		{"ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a)", ast.ConstrPrimaryKey, "", ""},
		{"ALTER TABLE t ADD FOREIGN KEY IF NOT EXISTS n (a) REFERENCES r (id)", ast.ConstrForeignKey, "", "n"},
		{"ALTER TABLE t ADD CONSTRAINT c FOREIGN KEY IF NOT EXISTS n (a) REFERENCES r (id)", ast.ConstrForeignKey, "c", "n"},
		{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c CHECK (a > 0)", ast.ConstrCheck, "c", ""},
		{"ALTER TABLE t ADD VECTOR IF NOT EXISTS n (a)", ast.ConstrVector, "n", ""},
	} {
		c := oneNode(t, tc.sql).(*ast.AlterTableStmt).Commands[0].Constraint
		if c == nil {
			t.Errorf("%q: no constraint", tc.sql)
			continue
		}
		if !c.IfNotExists || c.Type != tc.kind || c.Name != tc.name || c.IndexName != tc.indexName {
			t.Errorf("%q: IfNotExists=%v Type=%v Name=%q IndexName=%q", tc.sql, c.IfNotExists, c.Type, c.Name, c.IndexName)
		}
		if c.Type == ast.ConstrCheck && c.Expr == nil {
			t.Errorf("%q: the CHECK lost its expression", tc.sql)
		}
		if c.Loc.End != len(tc.sql) {
			t.Errorf("%q: Loc.End=%d, want %d", tc.sql, c.Loc.End, len(tc.sql))
		}
	}

	list := oneNode(t, "ALTER TABLE t ADD (d INT, INDEX IF NOT EXISTS n (a), UNIQUE u (a))").(*ast.AlterTableStmt).Commands[0]
	if len(list.Columns) != 1 || len(list.Constraints) != 2 || !list.Constraints[0].IfNotExists || list.Constraints[1].IfNotExists {
		t.Errorf("the list: columns %+v, constraints %+v", list.Columns, list.Constraints)
	}
	plain := oneNode(t, "ALTER TABLE t ADD INDEX n (a)").(*ast.AlterTableStmt).Commands[0].Constraint
	if plain.IfNotExists {
		t.Errorf("a plain ADD INDEX reads as IF NOT EXISTS")
	}
}

// TestNamedScalarTypeSizeAST: the size UUID, INET4 and INET6 take is read and dropped.
func TestNamedScalarTypeSizeAST(t *testing.T) {
	sql := "CREATE TABLE t (a inet6(3), b UUID(3,2), c `INET4`(1.5), d inet6)"
	stmt := oneNode(t, sql).(*ast.CreateTableStmt)
	want := []struct{ name, text string }{{"INET6", "inet6(3)"}, {"UUID", "UUID(3,2)"}, {"INET4", "`INET4`(1.5)"}, {"INET6", "inet6"}}
	for i, w := range want {
		dt := stmt.Columns[i].TypeName
		if dt.Name != w.name || dt.Length != 0 {
			t.Errorf("column %d: Name=%q Length=%d, want %s and no length", i, dt.Name, dt.Length, w.name)
		}
		if got := sql[dt.Loc.Start:dt.Loc.End]; got != w.text {
			t.Errorf("column %d: Loc spans %q, want %q", i, got, w.text)
		}
	}
}
