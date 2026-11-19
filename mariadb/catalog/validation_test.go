package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// TestEnumDuplicateValueMessage holds 1291's message to the server's: the value
// named is the first spelling, which MariaDB compares case-insensitively.
// Measured on mariadb:11.8.8: `CREATE TABLE t (e ENUM('s','S'))` answers
// "Column 'e' has duplicated value 's' in ENUM", and a SET the same with SET.
func TestEnumDuplicateValueMessage(t *testing.T) {
	for _, tc := range []struct{ ddl, want string }{
		{"CREATE TABLE t (e ENUM('s', 'S'))", "Column 'e' has duplicated value 's' in ENUM"},
		{"CREATE TABLE t (e SET('a', 'b', 'A'))", "Column 'e' has duplicated value 'a' in SET"},
		{"CREATE TABLE t (e ENUM('a ', 'a'))", "Column 'e' has duplicated value 'a' in ENUM"},
		{"CREATE TABLE t (e ENUM('S ', 's'))", "Column 'e' has duplicated value 'S' in ENUM"},
	} {
		err := execErr(t, tc.ddl)
		catErr, ok := err.(*Error)
		if !ok || catErr.Code != 1291 {
			t.Fatalf("%s: expected 1291, got %v", tc.ddl, err)
		}
		if !strings.Contains(catErr.Message, tc.want) {
			t.Errorf("%s: Message = %q, want %q", tc.ddl, catErr.Message, tc.want)
		}
	}

	// Outside strict sql_mode, 1291 is a note and the table is created: measured
	// on mariadb:11.8.8 under `SET sql_mode = ''` and `'ANSI'`.
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; SET sql_mode = ''; CREATE TABLE t (e ENUM('s', 's'));")
	if c.GetDatabase("d").GetTable("t") == nil {
		t.Fatal("non-strict sql_mode rejected ENUM('s', 's')")
	}
}

// TestDiscardTablespaceFollowedByCopyAlterRollsBack ensures that even when a
// caller hands alterTable a multi-command AST pairing ATDiscardTablespace with
// a copy command (which the SQL parser already rejects as 1064), the statement
// answers 1814 and rolls back TablespaceDiscarded and Engine.
func TestDiscardTablespaceFollowedByCopyAlterRollsBack(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT PRIMARY KEY);")
	err := c.alterTable(&nodes.AlterTableStmt{
		Table: &nodes.TableRef{Name: "t"},
		Commands: []*nodes.AlterTableCmd{
			{Type: nodes.ATDiscardTablespace},
			{Type: nodes.ATTableOption, Option: &nodes.TableOption{Name: "ENGINE", Value: "MyISAM"}},
		},
	})
	catErr, ok := err.(*Error)
	if !ok || catErr.Code != 1814 {
		t.Fatalf("expected 1814, got %v", err)
	}
	tbl := c.GetDatabase("d").GetTable("t")
	if tbl.TablespaceDiscarded || tbl.Engine != "" {
		t.Fatalf("failed multi-command ALTER mutated table: discarded=%v engine=%q", tbl.TablespaceDiscarded, tbl.Engine)
	}
}

// TestCTASDeclaredColumnExpressions ensures CREATE TABLE ... AS SELECT allows
// declared column DEFAULT and CHECK expressions to reference columns projected
// by the SELECT while still rejecting unknown columns with 1054.
func TestCTASDeclaredColumnExpressions(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; CREATE TABLE a (d INT CHECK (d < b), e INT DEFAULT (b + 1)) AS SELECT 1 AS b;")

	err := execErr(t, "CREATE TABLE a (d INT DEFAULT missing) AS SELECT 1 AS b")
	catErr, ok := err.(*Error)
	if !ok || catErr.Code != 1054 {
		t.Fatalf("expected 1054 for missing column in CTAS DEFAULT, got %v", err)
	}
}
