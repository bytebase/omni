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
CREATE TABLE s.ids (id int);
CREATE SEQUENCE public.old_seq;
CREATE TABLE public.pt (id int) PARTITION BY RANGE (id);
CREATE TABLE public.pt1 PARTITION OF public.pt FOR VALUES FROM (0) TO (10);
INSERT INTO public.pt VALUES (1);
DO $$ BEGIN CREATE ROLE prior_backup_old SUPERUSER; EXCEPTION WHEN duplicate_object THEN NULL; END $$;
CREATE FUNCTION public.old(int) RETURNS bool LANGUAGE sql AS 'SELECT true';
INSERT INTO p VALUES (1), (2);
INSERT INTO t (id, p_id, n) VALUES (1, 1, 1), (2, 2, 2);
INSERT INTO r VALUES (1), (2);
INSERT INTO nokey VALUES (1), (2);
`

// TestPriorBackupAgainstPostgres builds each table's backup query as
// Bytebase does and runs it on PostgreSQL before the change, as Bytebase
// does. A PriorBackup finding on a statement must mean that query fails,
// changes data, backs up rows the table does not have, or misses rows a
// statement of the table selects; where the rule follows the whole change
// (complete), every table whose query does so must have a finding.
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
		{"a data-modifying WITH of a later statement", "DELETE FROM t WHERE id = 1;\nWITH d AS (DELETE FROM r RETURNING t_id) DELETE FROM t WHERE id IN (SELECT t_id FROM d);", true},
		{"a later statement's own WITH", "DELETE FROM nokey WHERE a = 1;\nWITH x AS (SELECT 2 AS a) DELETE FROM nokey WHERE a IN (SELECT a FROM x);", true},
		// The backup reads the second statement's rows through the first x,
		// which here misses them, and elsewhere may not: not judged.
		{"a later statement's WITH the first one defines otherwise", "WITH x AS (SELECT 1 AS a) DELETE FROM nokey WHERE a IN (SELECT a FROM x);\nWITH x AS (SELECT 2 AS a) DELETE FROM nokey WHERE a IN (SELECT a FROM x);", false},
		{"a later statement's WITH the first one defines alike", "WITH x AS (SELECT 1 AS a) DELETE FROM nokey WHERE a IN (SELECT a FROM x);\nWITH x AS (SELECT 1 AS a) DELETE FROM nokey WHERE a IN (SELECT a FROM x);", true},
		{"a WITH query shadowing a table", "WITH nokey AS (SELECT 1 AS a) DELETE FROM t WHERE id IN (SELECT a FROM nokey);", true},
		{"after SET search_path", "SET search_path = s, public;\nALTER TABLE t ADD COLUMN x int;\nDELETE FROM t WHERE x = 1;", true},
		{"a WITH query named like the table", "WITH nokey AS (SELECT 1 AS a) DELETE FROM nokey WHERE a = 1;", true},
		{"a WITH query named like the table, through an alias", "WITH nokey AS (SELECT 9 AS a) DELETE FROM nokey AS x WHERE x.a = 9;", true},
		{"aliases that differ between statements", "UPDATE nokey AS x SET a = 1 WHERE x.a = 1;\nUPDATE nokey AS y SET a = 2 WHERE y.a = 2;", true},
		{"an alias, then none", "DELETE FROM nokey AS x WHERE x.a = 1;\nDELETE FROM nokey WHERE a = 2;", true},
		{"none, then an alias", "DELETE FROM nokey WHERE a = 1;\nDELETE FROM nokey AS x WHERE x.a = 2;", true},
		{"the same alias", "UPDATE nokey AS x SET a = 1 WHERE x.a = 1;\nUPDATE nokey AS x SET a = 2 WHERE x.a = 2;", true},
		{"a column named with its schema", "ALTER TABLE public.p ADD COLUMN gone bool;\nDELETE FROM t USING public.p WHERE t.p_id = public.p.id AND public.p.gone;", true},
		{"a sequence the change creates", "CREATE SEQUENCE ids;\nDELETE FROM t USING ids WHERE false;", true},
		{"a table the change creates, of a name a schema off the path holds", "CREATE TABLE public.ids (id int);\nDELETE FROM t USING ids WHERE t.id = ids.id;", true},
		{"a view column the change renames to", "ALTER VIEW tv RENAME COLUMN code TO c;\nDELETE FROM t USING tv WHERE tv.c = t.code;", true},
		{"a routine the change renames to", "ALTER FUNCTION public.old(int) RENAME TO stale;\nDELETE FROM nokey WHERE public.stale(a);", true},
		{"a column the change adds to a table of that name in another schema", "ALTER TABLE s.t ADD COLUMN b int;\nDELETE FROM public.t AS x WHERE x.b;", false},
		{"a column of a table the change recreates", "DROP TABLE nokey;\nCREATE TABLE nokey (b int);\nDELETE FROM nokey WHERE b = 1;", true},
		{"a routine the change creates in another schema", "CREATE FUNCTION s.stale(int) RETURNS bool LANGUAGE sql AS 'SELECT true';\nDELETE FROM nokey WHERE public.stale(a);", false},
		{"a table the first SET search_path moves", "UPDATE t SET n = 1 WHERE id = 1;\nSET search_path = s, public;\nDELETE FROM t WHERE id = 2;", true},
		{"a table the first SET search_path moves, through an alias", "UPDATE t AS x SET n = 1 WHERE x.id = 1;\nSET search_path = s, public;\nDELETE FROM t AS y WHERE y.id = 2;", true},
		{"WHERE CURRENT OF", "UPDATE t SET n = 1 WHERE CURRENT OF c;", true},
		{"a sequence the change renames to", "ALTER SEQUENCE old_seq RENAME TO ids;\nDELETE FROM t USING ids WHERE false;", true},
		{"a column the change adds, under TABLESAMPLE", "ALTER TABLE p ADD COLUMN gone bool;\nDELETE FROM t USING p TABLESAMPLE SYSTEM (100) WHERE p.gone;", true},
		{"a column the change adds, as a field of a row", "ALTER TABLE p ADD COLUMN gone bool;\nDELETE FROM t USING p WHERE (p).gone;", true},
		{"a relation named pg_ in a schema of the target", "CREATE TABLE public.pg_ids (id int);\nDELETE FROM t USING public.pg_ids WHERE false;", true},
		{"SET ROLE of a role the change creates", "CREATE ROLE prior_backup_worker;\nSET ROLE prior_backup_worker;\nUPDATE public.t SET n = 1 WHERE id = 1;", true},
		{"a column list of an alias", "ALTER TABLE p RENAME COLUMN id TO a;\nDELETE FROM t USING p AS x(a) WHERE x.a = t.p_id;", true},
		{"a column the change adds, past a column list of an alias", "ALTER TABLE p ADD COLUMN gone bool;\nDELETE FROM t USING p AS x(i) WHERE x.gone;", true},
		{"a column of a table CREATE TABLE AS recreates", "DROP TABLE nokey;\nCREATE TABLE nokey AS SELECT 1 AS b;\nDELETE FROM nokey WHERE b = 1;", true},
		{"a column named with the database", "ALTER TABLE public.p ADD COLUMN gone bool;\nDELETE FROM t USING omni_review.public.p WHERE omni_review.public.p.gone;", true},
		{"a column of an aliased join", "ALTER TABLE p ADD COLUMN gone bool;\nDELETE FROM t USING (p JOIN nokey ON true) AS j WHERE j.gone;", true},
		{"a column the change adds to a partitioned table, through a partition", "ALTER TABLE pt ADD COLUMN flag bool;\nDELETE FROM pt1 WHERE pt1.flag;", true},
		{"SET ROLE of a role the change recreates", "DROP ROLE prior_backup_old;\nCREATE ROLE prior_backup_old SUPERUSER;\nSET ROLE prior_backup_old;\nUPDATE public.t SET n = 1 WHERE id = 1;", true},
		{"a column a join's USING names", "ALTER TABLE p ADD COLUMN gone bool;\nALTER TABLE nokey ADD COLUMN gone bool;\nDELETE FROM t USING p JOIN nokey USING (gone) WHERE false;", true},
		{"an alias's column list longer than the relation", "ALTER TABLE p ADD COLUMN gone bool;\nDELETE FROM t USING p AS x(i, g) WHERE false;", true},
		{"a column of a table CREATE TABLE LIKE recreates", "DROP TABLE nokey;\nCREATE TABLE nokey (LIKE p);\nDELETE FROM nokey WHERE id = 1;", true},
		{"a column of the table beside a VALUES list", "ALTER TABLE nokey ADD COLUMN b bool;\nDELETE FROM nokey USING (VALUES (1)) AS q(x) WHERE b;", true},
		{"a relation after SET ROLE DEFAULT", "CREATE TABLE public.ids (id int);\nSET role TO DEFAULT;\nDELETE FROM t USING ids WHERE t.id = ids.id;", true},
		{"a column of a table a rename gives the name", "DROP TABLE nokey;\nALTER TABLE t RENAME TO nokey;\nDELETE FROM nokey WHERE code = 'x';", true},
		{"a column of a table EXPLAIN ANALYZE recreates", "DROP TABLE nokey;\nEXPLAIN ANALYZE CREATE TABLE nokey AS SELECT 1 AS b;\nDELETE FROM nokey WHERE b = 1;", true},
		{"a column a copied WITH query reads", "ALTER TABLE p ADD COLUMN gone int;\nWITH x AS (SELECT p.gone FROM p) DELETE FROM t WHERE id IN (SELECT gone FROM x);", true},
		{"a column beside a join's USING", "ALTER TABLE nokey ADD COLUMN b bool;\nDELETE FROM nokey USING p JOIN t USING (id) WHERE b;", true},
		{"a column a function of the FROM list reads", "ALTER TABLE p ADD COLUMN gone int;\nDELETE FROM t USING p, generate_series(1, p.gone) AS g WHERE false;", true},
		{"a relation after SET search_path FROM CURRENT", "CREATE TABLE public.ids (id int);\nSET search_path FROM CURRENT;\nDELETE FROM t USING ids WHERE t.id = ids.id;", true},
		{"a table RESET search_path moves", "SET search_path = s, public;\nUPDATE t SET id = 1 WHERE id = 1;\nRESET search_path;\nUPDATE t SET n = 1 WHERE id = 1;", true},
		{"a first alias the later statement's table answers to", "DELETE FROM nokey AS nokey WHERE a = 1;\nDELETE FROM nokey WHERE a = 2;", true},
		{"no first alias, then the table's name as an alias", "DELETE FROM nokey WHERE a = 1;\nDELETE FROM nokey AS nokey WHERE a = 2;", true},
		{"a column of a view the change recreates", "DROP VIEW tv;\nCREATE VIEW tv AS SELECT 1 AS b;\nDELETE FROM t USING tv WHERE tv.b = t.id;", true},
		{"a column an outer source has, named in a subquery", "ALTER TABLE p ADD COLUMN gone bool;\nDELETE FROM t USING p WHERE EXISTS (SELECT 1 WHERE p.gone);", true},
		{"a column a copied WITH query orders by", "ALTER TABLE p ADD COLUMN gone int;\nWITH x AS (SELECT p.id FROM p ORDER BY p.gone) DELETE FROM t WHERE id IN (SELECT id FROM x);", true},
		{"an output a copied WITH query orders by", "ALTER TABLE p ADD COLUMN gone int;\nWITH y AS (SELECT p.id AS gone FROM p ORDER BY gone) DELETE FROM nokey WHERE a IN (SELECT gone FROM y);", true},
		{"a join alias's column list longer than the join", "ALTER TABLE p ADD COLUMN gone bool;\nDELETE FROM t USING (p JOIN nokey ON true) AS j(i, a, g) WHERE false;", true},
		{"a table the change creates earlier on the path", "SET search_path = s, public;\nCREATE TABLE s.p (id int);\nUPDATE p SET id = 1 WHERE id = 2;", true},
		{"a search path set_config sets", "SELECT pg_catalog.set_config('search_path', 's, public', false);\nUPDATE t SET id = 1 WHERE id = 2;", true},
		// The first CTE stands for p in the second arm, which reads other
		// rows than the statement's, or not: not judged.
		{"a later source a WITH query of the first statement shadows", "WITH p AS (SELECT 999 AS id) DELETE FROM t WHERE id = 1;\nDELETE FROM t USING p WHERE t.p_id = p.id;", false},
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
			syncPartitions(t, conn, schema)
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
				i := slices.IndexFunc(groups, func(g backupGroup) bool { return slices.Contains(g.statements, f.Statement) && g.outcome != "" })
				if i < 0 {
					i = slices.IndexFunc(groups, func(g backupGroup) bool { return slices.Contains(g.statements, f.Statement) })
				}
				if i < 0 || groups[i].outcome == "" {
					t.Errorf("reported what the server contradicts, the backup of the statement's table runs: %d %s", f.Statement, f.Message)
				}
			}
			if !tt.complete {
				return
			}
			for _, g := range groups {
				// A table's query that the prepended SET fails is the SET's
				// to report.
				if g.outcome == "" || strings.HasPrefix(g.outcome, "prepend fails") && g.table != prependedSet {
					continue
				}
				if !slices.ContainsFunc(result.Findings, func(f review.Finding) bool { return slices.Contains(g.statements, f.Statement) }) {
					t.Errorf("missed: the backup of %s %s\nfindings: %v", g.table, g.outcome, result.Findings)
				}
			}
		})
	}
}

// syncPartitions lists each synced table's partitions under it, as
// Bytebase syncs them, with their own partitions.
func syncPartitions(t *testing.T, conn *sql.Conn, schema *metadata.DatabaseSchemaMetadata) {
	t.Helper()
	ctx := context.Background()
	var walk func(parent string) []*metadata.TablePartitionMetadata
	walk = func(parent string) []*metadata.TablePartitionMetadata {
		rows, err := conn.QueryContext(ctx, `SELECT c.oid::regclass::text, c.relname FROM pg_inherits i JOIN pg_class c ON c.oid = i.inhrelid
			WHERE i.inhparent = $1::regclass AND c.relispartition ORDER BY 2`, parent)
		if err != nil {
			t.Fatal(err)
		}
		type child struct{ ref, name string }
		var children []child
		for rows.Next() {
			var c child
			if err := rows.Scan(&c.ref, &c.name); err != nil {
				t.Fatal(err)
			}
			children = append(children, c)
		}
		rows.Close()
		var out []*metadata.TablePartitionMetadata
		for _, c := range children {
			out = append(out, &metadata.TablePartitionMetadata{Name: c.name, Subpartitions: walk(c.ref)})
		}
		return out
	}
	for _, s := range schema.GetSchemas() {
		for _, tb := range s.GetTables() {
			tb.Partitions = walk(fmt.Sprintf("%q.%q", s.GetName(), tb.GetName()))
		}
	}
}

// prependedSet names the group of the SET Bytebase runs before the
// backup's queries.
const prependedSet = "the prepended SET"

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
	prependAt := -1
	// changes is the table each UPDATE or DELETE changes, as the server
	// resolves it after the change's SETs before it, 0 when unknown.
	changes := make(map[int]int64)
	ctx := context.Background()
	etx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	engineKnown := true
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
		// The statements before an UPDATE or DELETE run, as the change
		// runs them, so the server resolves its table as the change does.
		// After one fails, the transaction is aborted and the server's
		// resolution is not asked again.
		switch node.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
		default:
			if engineKnown {
				if _, err := etx.ExecContext(ctx, change[seg.ByteStart:seg.ByteEnd]); err != nil {
					engineKnown = false
				}
			}
		}
		var rv *ast.RangeVar
		switch v := node.(type) {
		case *ast.VariableSetStmt:
			if prepend == "" && (strings.EqualFold(v.Name, "role") || strings.EqualFold(v.Name, "search_path")) {
				prepend = strings.TrimSuffix(strings.TrimSpace(change[seg.ByteStart:seg.ByteEnd]), ";") + ";"
				prependAt = i
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
		if engineKnown {
			name := fmt.Sprintf("%q", rv.Relname)
			if rv.Schemaname != "" {
				name = fmt.Sprintf("%q.%q", rv.Schemaname, rv.Relname)
			}
			var oid sql.NullInt64
			if err := etx.QueryRowContext(ctx, "SELECT to_regclass($1)::oid::int8", name).Scan(&oid); err != nil {
				engineKnown = false
			} else {
				changes[i] = oid.Int64
			}
		}
		// Bytebase refuses the statement before it runs any query.
		if names := dataModifyingCTEs(withOf(node)); len(names) > 0 {
			unresolved = append(unresolved, backupGroup{table: rv.Relname, statements: []int{i}, outcome: "fails in Bytebase: a data-modifying WITH query"})
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
	if err := etx.Rollback(); err != nil {
		t.Fatal(err)
	}
	text := func(loc ast.Loc) string {
		if loc.Start < 0 || loc.End <= loc.Start || loc.End > len(change) {
			return ""
		}
		return strings.TrimSpace(change[loc.Start:loc.End])
	}
	groups := unresolved
	if prepend != "" && len(order) > 0 {
		// The SET Bytebase runs before the queries, on its own.
		outcome := runBackup(t, conn, prepend, "SELECT 1", "", "", nil)
		groups = append(groups, backupGroup{table: prependedSet, statements: []int{prependAt}, outcome: outcome})
	}
	for n, key := range order {
		items := byTable[key]
		var b strings.Builder
		backup := fmt.Sprintf("\"bbdataarchive\".\"bk_%d\"", n)
		fmt.Fprintf(&b, "CREATE TABLE %s AS\n", backup)
		// Every arm selects through the first statement's alias.
		firstTarget, _, _ := dmlParts(items[0].node)
		alias := aliasOf(firstTarget)
		if w := withOf(items[0].node); w != nil {
			b.WriteString(text(w.Loc) + "\n")
		}
		g := backupGroup{table: key}
		// owns are the statements' own queries of their rows, with their
		// own WITH clause and alias, which the backup must hold.
		var owns []string
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
			if alias != "" {
				fmt.Fprintf(&b, "  SELECT \"%s\".* ", alias)
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
			own := ""
			if w := withOf(it.node); w != nil {
				own = text(w.Loc) + " "
			}
			sel := key
			if a := aliasOf(rv); a != "" {
				sel = fmt.Sprintf("%q", a)
			}
			own += "SELECT " + sel + ".* FROM " + text(rv.Loc)
			if from != nil && len(from.Items) > 0 {
				own += ", " + text(ast.ListSpan(from))
			}
			if where != nil {
				own += " WHERE " + text(ast.NodeLoc(where))
			}
			owns = append(owns, own)
		}
		b.WriteString(";")
		g.outcome = runBackup(t, conn, prepend, b.String(), backup, key, owns)
		if g.outcome == "" {
			var oid sql.NullInt64
			if err := conn.QueryRowContext(ctx, "SELECT to_regclass($1)::oid::int8", key).Scan(&oid); err != nil {
				t.Fatal(err)
			}
			for _, i := range g.statements {
				if c := changes[i]; c != 0 && c != oid.Int64 {
					g.outcome = fmt.Sprintf("backs up %s for statement %d, which changes another table", key, i)
				}
			}
		}
		groups = append(groups, g)
	}
	return groups
}

// runBackup runs one backup query in a transaction it rolls back, and
// says how it went: "" when it ran, left every table's rows as they were,
// backed up only rows the table has, and every row each of the
// statements' own queries selects, where those run.
func runBackup(t *testing.T, conn *sql.Conn, prepend, query, backup, table string, owns []string) string {
	t.Helper()
	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	if prepend != "" {
		if _, err := tx.ExecContext(ctx, prepend); err != nil {
			return "prepend fails: " + err.Error()
		}
	}
	if backup == "" {
		return ""
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
	var foreign int
	if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM (SELECT * FROM "+backup+" EXCEPT ALL SELECT * FROM "+table+") x").Scan(&foreign); err != nil {
		return "backs up rows of another shape than " + table + ": " + err.Error()
	}
	if foreign > 0 {
		return "backs up rows " + table + " does not have"
	}
	for k, own := range owns {
		if _, err := tx.ExecContext(ctx, "SAVEPOINT own"); err != nil {
			t.Fatal(err)
		}
		var missed int
		if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM (("+own+") EXCEPT (SELECT * FROM "+backup+")) x").Scan(&missed); err != nil {
			if _, err := tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT own"); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if missed > 0 {
			return fmt.Sprintf("misses rows its statement %d selects", k)
		}
	}
	return ""
}
