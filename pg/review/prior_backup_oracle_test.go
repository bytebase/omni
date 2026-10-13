package review

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/review"
)

// priorBackupRows gives the oracle's tables rows, so a backup query that
// changes data shows it.
const priorBackupRows = `
CREATE SCHEMA bbdataarchive;
INSERT INTO p VALUES (1), (2);
INSERT INTO t (id, p_id, n) VALUES (1, 1, 1), (2, 2, 2);
INSERT INTO r VALUES (1), (2);
INSERT INTO nokey VALUES (1), (2);
`

// TestPriorBackupAgainstPostgres builds each table's backup query as
// Bytebase does and runs it on PostgreSQL before the change, as Bytebase
// does. A PriorBackup finding on a statement must mean that query fails or
// changes data; where the rule follows the whole change (complete), every
// table whose query does so must have a finding.
func TestPriorBackupAgainstPostgres(t *testing.T) {
	db := startPostgres(t)
	tests := []struct {
		name     string
		change   string
		complete bool
	}{
		{"a column the change adds", "ALTER TABLE nokey ADD COLUMN b int;\nUPDATE nokey SET b = 1 WHERE b IS NULL;", true},
		{"a column the change renames to", "ALTER TABLE t RENAME COLUMN n TO m;\nDELETE FROM t WHERE m > 5;", true},
		{"a column the change adds, which the backup does not read", "ALTER TABLE nokey ADD COLUMN b int;\nUPDATE nokey SET b = 1 WHERE a = 1;", true},
		{"a column the change adds, through an alias", "ALTER TABLE t ADD COLUMN flag bool;\nUPDATE t AS x SET n = 1 WHERE x.flag;", true},
		{"a column the change adds to a USING table", "ALTER TABLE p ADD COLUMN gone bool;\nDELETE FROM t USING p WHERE t.p_id = p.id AND p.gone;", true},
		{"a column the change adds, in a subquery", "ALTER TABLE p ADD COLUMN gone bool;\nDELETE FROM t WHERE p_id IN (SELECT id FROM p WHERE gone);", false},
		{"whole-row and system columns", "ALTER TABLE nokey ADD COLUMN ctid2 int;\nDELETE FROM nokey WHERE nokey IS NOT NULL AND ctid IS NOT NULL;", true},
		{"a table the change creates", "CREATE TABLE n (id int);\nDELETE FROM n WHERE id = 1;", true},
		{"a table the change creates, in USING", "CREATE TABLE ids (id int);\nDELETE FROM t USING ids WHERE t.id = ids.id;", true},
		{"a table the change creates, in a subquery", "CREATE TABLE ids (id int);\nUPDATE t SET n = 1 WHERE id IN (SELECT id FROM ids);", true},
		{"a table the change creates, in the WITH clause", "CREATE TABLE ids (id int);\nWITH x AS (SELECT id FROM ids) DELETE FROM t WHERE id IN (SELECT id FROM x);", true},
		{"a table the change renames to, in FROM", "ALTER TABLE nokey RENAME TO nk;\nUPDATE t SET n = 1 FROM nk WHERE t.id = nk.a;", true},
		{"a function the change creates, named with its schema", "CREATE FUNCTION public.stale(int) RETURNS bool LANGUAGE sql AS 'SELECT true';\nDELETE FROM nokey WHERE public.stale(a);", true},
		{"a function the change creates, unqualified", "CREATE FUNCTION stale(int) RETURNS bool LANGUAGE sql AS 'SELECT true';\nDELETE FROM nokey WHERE stale(a);", false},
		{"a data-modifying WITH", "WITH d AS (DELETE FROM r RETURNING t_id) UPDATE t SET n = 5 WHERE id IN (SELECT t_id FROM d);", true},
		{"a later statement's own WITH", "DELETE FROM nokey WHERE a = 1;\nWITH x AS (SELECT 2 AS a) DELETE FROM nokey WHERE a IN (SELECT a FROM x);", true},
		{"a later statement's WITH the first one defines too", "WITH x AS (SELECT 1 AS a) DELETE FROM nokey WHERE a IN (SELECT a FROM x);\nWITH x AS (SELECT 2 AS a) DELETE FROM nokey WHERE a IN (SELECT a FROM x);", true},
		{"a WITH query shadowing a table", "WITH nokey AS (SELECT 1 AS a) DELETE FROM t WHERE id IN (SELECT a FROM nokey);", true},
		{"after SET search_path", "SET search_path = s, public;\nALTER TABLE t ADD COLUMN x int;\nDELETE FROM t WHERE x = 1;", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			conn, err := db.Conn(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			resetOracle(t, conn)
			if _, err := conn.ExecContext(ctx, priorBackupRows); err != nil {
				t.Fatal(err)
			}
			schema := syncSchema(t, conn)
			result, err := Review(ctx, tt.change, review.Options{
				Rules:  []review.Rule{review.PriorBackup},
				Change: review.Change{PriorBackup: true, MaxBackupSize: 1 << 20},
			}, []review.Target{{Schema: schema, SessionUser: "postgres"}})
			if err != nil {
				t.Fatalf("Review: %v", err)
			}
			groups := backupTruth(t, conn, tt.change, schema)
			for _, f := range result.Findings {
				if f.Rule == review.Syntax {
					t.Fatalf("syntax: %s", f.Message)
				}
				if f.Statement < 0 || strings.Contains(f.Message, "both UPDATE and DELETE") {
					// The size limit and Bytebase's own refusal are not the
					// query's to judge.
					continue
				}
				i := slices.IndexFunc(groups, func(g backupGroup) bool { return slices.Contains(g.statements, f.Statement) })
				if i < 0 || groups[i].outcome == "" {
					t.Errorf("reported what the server contradicts, the backup of the statement's table runs: %d %s", f.Statement, f.Message)
				}
			}
			if !tt.complete {
				return
			}
			for _, g := range groups {
				if g.outcome == "" {
					continue
				}
				if !slices.ContainsFunc(result.Findings, func(f review.Finding) bool { return slices.Contains(g.statements, f.Statement) }) {
					t.Errorf("missed: the backup of %s %s\nfindings: %v", g.table, g.outcome, result.Findings)
				}
			}
		})
	}
}

// backupGroup is the statements the backup copies into one query, the
// table's, and how running that query went: "" when it ran and changed
// nothing.
type backupGroup struct {
	table      string
	statements []int
	outcome    string
}

// backupTruth builds the backup's queries for a change and runs each, with
// the first SET role or search_path of the change before it, in a
// transaction it rolls back.
//
// bytebase: backend/plugin/parser/pg/backup.go — prepareTransformation,
// generateSQLForTable; backend/runner/taskrun/database_migrate_executor.go —
// getPrependStatements
func backupTruth(t *testing.T, conn *sql.Conn, change string, schema *metadata.DatabaseSchemaMetadata) []backupGroup {
	t.Helper()
	type item struct {
		index int
		node  ast.Node
	}
	path := backupSearchPath(schema.GetSearchPath())
	prepend := ""
	byTable := make(map[string][]item)
	var order []string
	var unresolved []backupGroup
	index := 0
	for _, seg := range pg.Split(change) {
		if seg.Empty() {
			continue
		}
		i := index
		index++
		list, err := parser.ParseRange(change, seg.ByteStart, seg.ByteEnd)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		node := list.Items[0].(*ast.RawStmt).Stmt
		var rv *ast.RangeVar
		switch v := node.(type) {
		case *ast.VariableSetStmt:
			if prepend == "" && (strings.EqualFold(v.Name, "role") || strings.EqualFold(v.Name, "search_path")) {
				prepend = strings.TrimSuffix(strings.TrimSpace(change[seg.ByteStart:seg.ByteEnd]), ";") + ";"
			}
			if p, ok := backupSetSearchPath(v); ok {
				path = p
			}
			continue
		case *ast.UpdateStmt:
			rv = v.Relation
		case *ast.DeleteStmt:
			rv = v.Relation
		default:
			continue
		}
		table := rv.Schemaname
		if table == "" {
			table = searchBackupRelation(schema, path, rv.Relname)
			if table == "" {
				unresolved = append(unresolved, backupGroup{table: rv.Relname, statements: []int{i}, outcome: "fails in Bytebase: the table is not in the synced schema"})
				continue
			}
		}
		key := fmt.Sprintf("%q.%q", table, rv.Relname)
		if _, ok := byTable[key]; !ok {
			order = append(order, key)
		}
		byTable[key] = append(byTable[key], item{i, node})
	}
	text := func(loc ast.Loc) string {
		if loc.Start < 0 || loc.End <= loc.Start || loc.End > len(change) {
			return ""
		}
		return strings.TrimSpace(change[loc.Start:loc.End])
	}
	groups := unresolved
	for n, key := range order {
		items := byTable[key]
		var b strings.Builder
		fmt.Fprintf(&b, "CREATE TABLE \"bbdataarchive\".\"bk_%d\" AS\n", n)
		if w := withOf(items[0].node); w != nil {
			b.WriteString(text(w.Loc) + "\n")
		}
		g := backupGroup{table: key}
		for j, it := range items {
			g.statements = append(g.statements, it.index)
			if j > 0 {
				b.WriteString("\n  UNION\n")
			}
			var rv *ast.RangeVar
			var from *ast.List
			var where ast.Node
			switch v := it.node.(type) {
			case *ast.UpdateStmt:
				rv, from, where = v.Relation, v.FromClause, v.WhereClause
			case *ast.DeleteStmt:
				rv, from, where = v.Relation, v.UsingClause, v.WhereClause
			}
			if rv.Alias != nil && rv.Alias.Aliasname != "" {
				fmt.Fprintf(&b, "  SELECT \"%s\".* ", rv.Alias.Aliasname)
			} else {
				fmt.Fprintf(&b, "  SELECT %s.* ", key)
			}
			b.WriteString("FROM " + text(rv.Loc))
			if from != nil && len(from.Items) > 0 {
				b.WriteString(", " + text(ast.ListSpan(from)))
			}
			if where != nil {
				b.WriteString(" WHERE " + text(ast.NodeLoc(where)))
			}
		}
		b.WriteString(";")
		g.outcome = runBackup(t, conn, prepend, b.String())
		groups = append(groups, g)
	}
	return groups
}

// runBackup runs one backup query in a transaction it rolls back, and
// says how it went: "" when it ran and left every table's rows as they
// were.
func runBackup(t *testing.T, conn *sql.Conn, prepend, query string) string {
	t.Helper()
	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	if prepend != "" {
		if _, err := tx.ExecContext(ctx, prepend); err != nil {
			return "fails: " + err.Error()
		}
	}
	counts := func() map[string]int {
		out := make(map[string]int)
		rows, err := tx.QueryContext(ctx, `SELECT format('%I.%I', n.nspname, c.relname) FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
			WHERE c.relkind = 'r' AND n.nspname NOT LIKE 'pg\_%' AND n.nspname NOT IN ('information_schema', 'bbdataarchive')`)
		if err != nil {
			t.Fatal(err)
		}
		var names []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				t.Fatal(err)
			}
			names = append(names, name)
		}
		rows.Close()
		for _, name := range names {
			var n int
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM "+name).Scan(&n); err != nil {
				t.Fatal(err)
			}
			out[name] = n
		}
		return out
	}
	before := counts()
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return "fails: " + err.Error()
	}
	after := counts()
	for name, n := range before {
		if after[name] != n {
			return "changes the rows of " + name
		}
	}
	return ""
}
