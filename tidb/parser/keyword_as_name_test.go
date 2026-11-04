package parser

import (
	"fmt"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// The words this engine registers as unreserved keywords, used as ordinary
// names.
//
// Registering a word as a keyword token is not free: every place that reads a
// name by testing for tokIDENT alone stops accepting it. That is how the SHOW
// batch below broke `DELETE FROM t config WHERE config.id=1`, which TiDB
// v8.5.5 parses -- TableAsNameOpt reads an Identifier, and an Identifier is
// any unreserved keyword. The parser has isIdentToken for exactly this, and
// the implicit-alias branch of parseDeleteStmt was testing the raw token
// instead.
//
// So this test sweeps every name position a statement has, with words from
// the batch in each of them. It is the gate the next batch has to pass: a
// keyword added without routing its name positions through isIdentToken
// fails here rather than in a customer's DELETE.
//
// Every case was measured against TiDB v8.5.5 on 2026-11-04. The server
// accepts all of them.

// showKeywordNames are the words tidb/parser/lexer.go registers for SHOW.
// They are TiDBKeyword or UnReservedKeyword upstream (parser.y:6962, 7353),
// so every one is an Identifier and may name anything.
var showKeywordNames = []string{
	"bindings", "binding_cache", "builtins", "backups", "restores",
	"column_stats_usage", "config", "distribution", "distributions",
	"histograms_in_flight", "job", "jobs", "next_row_id", "session_states",
	"stats_buckets", "stats_extended", "stats_healthy", "stats_histograms",
	"stats_locked", "stats_meta", "stats_topn",
}

// keywordNamePositions are the name positions each word is tried in. %s is
// the name; a template may use it more than once.
var keywordNamePositions = []string{
	// The alias positions, with and without AS. The implicit forms are the
	// ones a raw tokIDENT test breaks, because the explicit ones go through
	// parseIdent after consuming AS.
	"SELECT * FROM t %s",
	"SELECT * FROM t AS %s",
	"DELETE FROM t %s WHERE %s.id=1",
	"DELETE FROM t AS %s WHERE %s.id=1",
	"UPDATE t %s SET a=1",
	"UPDATE t AS %s SET a=1",
	"SELECT a %s FROM t",
	"SELECT a AS %s FROM t",
	"SELECT * FROM (SELECT 1) %s",
	"WITH %s AS (SELECT 1) SELECT * FROM %s",
	// Object names.
	"CREATE TABLE %s (a INT)",
	"CREATE TABLE t2 (%s INT)",
	"CREATE INDEX %s ON t (a)",
	"CREATE TABLE t2 (a INT, INDEX %s (a))",
	"ALTER TABLE t ADD COLUMN %s INT",
	"ALTER TABLE t ADD CONSTRAINT %s UNIQUE (a)",
	"CREATE VIEW %s AS SELECT 1",
	"CREATE DATABASE %s",
	"CREATE USER %s",
	"DROP TABLE %s",
	"RENAME TABLE %s TO t2",
	"ANALYZE TABLE %s",
	"USE %s",
	"SAVEPOINT %s",
	"PREPARE %s FROM 'select 1'",
	"SELECT * FROM t PARTITION (%s)",
	// Column references.
	"SELECT * FROM t WHERE %s = 1",
	"SELECT * FROM t ORDER BY %s",
	"SELECT * FROM t GROUP BY %s",
	"INSERT INTO t (%s) VALUES (1)",
}

// TestShowKeywordsStayUsableAsNames runs every word through every position.
func TestShowKeywordsStayUsableAsNames(t *testing.T) {
	for _, pos := range keywordNamePositions {
		for _, name := range showKeywordNames {
			sql := fmt.Sprintf(strings.ReplaceAll(pos, "%s", "%[1]s"), name)
			t.Run(sql, func(t *testing.T) {
				if _, err := Parse(sql); err != nil {
					t.Errorf("Parse(%q): %v\n\n%q is an unreserved keyword, so it is an Identifier and may be used as a name here. A name position that reads the raw tokIDENT instead of isIdentToken rejects it.", sql, err, name)
				}
			})
		}
	}
}

// TestDeleteAliasTakesAnUnreservedKeyword pins the position the SHOW keyword
// batch actually broke, on its own, so the reason stays legible if the sweep
// above is ever trimmed. `DELETE FROM t config` names an alias `config`;
// reading it as a table with nothing after it left `config` unconsumed and
// the statement failed.
func TestDeleteAliasTakesAnUnreservedKeyword(t *testing.T) {
	for _, tc := range []struct{ sql, alias string }{
		{"DELETE FROM t config WHERE config.id=1", "config"},
		{"DELETE FROM t AS config WHERE config.id=1", "config"},
		{"DELETE FROM t job WHERE job.id=1", "job"},
		{"DELETE FROM t bindings WHERE bindings.id=1", "bindings"},
		{"DELETE FROM t stats_meta WHERE stats_meta.id=1", "stats_meta"},
		// A plain identifier was never affected; it is here as the control.
		{"DELETE FROM t d WHERE d.id=1", "d"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.sql, err)
			}
			if len(list.Items) != 1 {
				t.Fatalf("Parse(%q) returned %d statements, want 1", tc.sql, len(list.Items))
			}
			stmt, ok := list.Items[0].(*nodes.DeleteStmt)
			if !ok {
				t.Fatalf("Parse(%q) returned %T, want *ast.DeleteStmt", tc.sql, list.Items[0])
			}
			if len(stmt.Tables) != 1 {
				t.Fatalf("got %d tables, want 1", len(stmt.Tables))
			}
			ref, ok := stmt.Tables[0].(*nodes.TableRef)
			if !ok {
				t.Fatalf("table is %T, want *ast.TableRef", stmt.Tables[0])
			}
			// The word is the alias of `t`, not a table of its own and not a
			// token left for whatever follows.
			if ref.Name != "t" {
				t.Errorf("table = %q, want t", ref.Name)
			}
			if ref.Alias != tc.alias {
				t.Errorf("alias = %q, want %q", ref.Alias, tc.alias)
			}
		})
	}
}
