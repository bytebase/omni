package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// analyzeOracle is a PostgreSQL 18 server, the release pg_proc_generated.go
// is generated from, that the analyzer's results are compared with.
type analyzeOracle struct {
	db  *sql.DB
	ctx context.Context
}

var (
	analyzeOracleOnce    sync.Once
	analyzeOracleInst    *analyzeOracle
	analyzeOracleErr     error
	analyzeOracleSchemas atomic.Int64
)

func startAnalyzeOracle(t *testing.T) *analyzeOracle {
	t.Helper()
	analyzeOracleOnce.Do(func() {
		ctx := context.Background()
		container, err := tcpg.Run(ctx, "postgres:18-alpine",
			tcpg.WithDatabase("omni_analyze"),
			tcpg.WithUsername("postgres"),
			tcpg.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2)),
		)
		if err != nil {
			analyzeOracleErr = fmt.Errorf("container start: %w", err)
			return
		}
		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			analyzeOracleErr = fmt.Errorf("conn string: %w", err)
			return
		}
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			analyzeOracleErr = fmt.Errorf("db open: %w", err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
			analyzeOracleErr = fmt.Errorf("ping: %w", err)
			return
		}
		analyzeOracleInst = &analyzeOracle{db: db, ctx: ctx}
	})
	if analyzeOracleErr != nil {
		t.Fatalf("PostgreSQL 18 oracle required: %v", analyzeOracleErr)
	}
	return analyzeOracleInst
}

// session opens a connection with its own schema holding setup, first on the
// search path, and drops the schema when the test ends.
func (o *analyzeOracle) session(t *testing.T, setup string) *sql.Conn {
	t.Helper()
	conn, err := o.db.Conn(o.ctx)
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("analyze_%d", analyzeOracleSchemas.Add(1))
	t.Cleanup(func() {
		_, _ = conn.ExecContext(o.ctx, fmt.Sprintf("DROP SCHEMA %s CASCADE", schema))
		conn.Close()
	})
	if _, err := conn.ExecContext(o.ctx, fmt.Sprintf("CREATE SCHEMA %s; SET search_path TO %s", schema, schema)); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.ExecContext(o.ctx, setup); err != nil {
		t.Fatalf("PostgreSQL rejected setup: %v\n%s", err, setup)
	}
	return conn
}

// sqlState returns the SQLSTATE of an error from PostgreSQL or omni, or ""
// for no error.
func sqlState(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		return ""
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	var catErr *Error
	if errors.As(err, &catErr) {
		return catErr.Code
	}
	t.Fatalf("error without SQLSTATE: %v", err)
	return ""
}

const analyzeOracleSchema = `
	CREATE TABLE t (id int, masked text, x text);
	CREATE TABLE u (id int, y text);
	CREATE TABLE w (id int, x text, z int);
	CREATE TABLE p (a int, b text, c int);
	CREATE TABLE q (c int, b text, d int);
	CREATE VIEW cc AS SELECT x COLLATE "C" AS c, y AS d FROM t, u;
	CREATE VIEW cp AS SELECT y COLLATE "POSIX" AS p FROM u;
	CREATE FUNCTION os_sfunc(int, int) RETURNS int LANGUAGE sql AS 'select 1';
	CREATE AGGREGATE my_os(float8 ORDER BY int) (sfunc = os_sfunc, stype = int);
	CREATE FUNCTION two_sfunc(int, int, int) RETURNS int LANGUAGE sql AS 'select 1';
	CREATE AGGREGATE two_arg(int, int) (sfunc = two_sfunc, stype = int);
	CREATE FUNCTION poly_def(a int, b anyelement DEFAULT NULL::int) RETURNS anyelement LANGUAGE sql AS 'select $2';
`

// TestAnalyzeViewsMatchPostgres creates each view in PostgreSQL and in omni
// and compares pg_get_viewdef(v, true) and the view's columns.
func TestAnalyzeViewsMatchPostgres(t *testing.T) {
	o := startAnalyzeOracle(t)
	for _, tc := range []struct{ name, query string }{
		{"agg_filter", `SELECT sum(id) FILTER (WHERE masked = 'a') AS s, count(*) FILTER (WHERE x IS NULL) AS c FROM t`},
		{"agg_order", `SELECT string_agg(x, ',' ORDER BY masked DESC NULLS LAST) AS s, array_agg(DISTINCT x ORDER BY x) AS a FROM t`},
		{"agg_order_nulls_first", `SELECT array_agg(id ORDER BY x NULLS FIRST, masked DESC NULLS FIRST) AS a FROM t`},
		{"agg_order_using", `SELECT array_agg(x ORDER BY masked USING >, id USING <) AS a FROM t`},
		{"agg_order_qualified_op", `SELECT array_agg(x ORDER BY masked USING OPERATOR(pg_catalog.>)) AS a FROM t`},
		{"agg_order_expr", `SELECT string_agg(x, ',' ORDER BY lower(masked), id + 1) AS s FROM t`},
		{"agg_order_filter", `SELECT string_agg(x, ',' ORDER BY masked) FILTER (WHERE id > 0) AS s FROM t`},
		{"within_group", `SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY id) AS p, mode() WITHIN GROUP (ORDER BY x) AS m FROM t`},
		{"within_group_filter", `SELECT percentile_disc(0.5) WITHIN GROUP (ORDER BY x DESC) FILTER (WHERE id > 0) AS p FROM t`},
		{"hypothetical", `SELECT rank(3) WITHIN GROUP (ORDER BY id) AS r FROM t`},
		{"natural", `SELECT * FROM t NATURAL JOIN u`},
		{"natural_two_columns", `SELECT * FROM p NATURAL JOIN q`},
		{"natural_left", `SELECT * FROM t NATURAL LEFT JOIN w`},
		{"natural_right", `SELECT * FROM t NATURAL RIGHT JOIN w`},
		{"natural_full", `SELECT * FROM t NATURAL FULL JOIN w`},
		{"natural_no_common", `SELECT * FROM u NATURAL JOIN (SELECT 1 AS k) s`},
		{"natural_left_no_common", `SELECT * FROM u NATURAL LEFT JOIN (SELECT 1 AS k) s`},
		{"natural_chain", `SELECT * FROM t NATURAL JOIN u NATURAL JOIN w`},
		{"natural_nested", `SELECT * FROM (p NATURAL JOIN q) NATURAL JOIN u`},
		{"natural_columns", `SELECT id, masked, y FROM t NATURAL JOIN u WHERE id > 1`},
		{"natural_then_on", `SELECT masked, y, z FROM t NATURAL JOIN u JOIN w ON y = w.x`},
		{"using_order", `SELECT * FROM p JOIN (SELECT b, d FROM q) s USING (b)`},
		{"using_full", `SELECT * FROM p FULL JOIN q USING (b, c)`},
		{"using_right", `SELECT c, a, d FROM p RIGHT JOIN q USING (c)`},
		{"join_unqualified", `SELECT y, masked FROM t JOIN u ON t.id = u.id`},
		{"using_unqualified", `SELECT y, id FROM t JOIN u USING (id)`},
		{"setop", `SELECT id, x FROM t UNION SELECT id, y FROM u`},
		{"builtin_defaults", `SELECT jsonb_path_exists('{}'::jsonb, '$."a"') AS e, jsonb_set('{}'::jsonb, '{a}', '1') AS s, random_normal() AS r`},
		{"named_args", `SELECT make_interval(days => 1) AS a, jsonb_path_exists('{}'::jsonb, '$."a"', silent => true) AS b, make_interval(1, mins => id) AS c FROM t`},
		{"named_args_reordered", `SELECT jsonb_path_exists(path => '$."a"', target => '{}'::jsonb) AS e`},
		{"user_ordered_set", `SELECT my_os(0.5) WITHIN GROUP (ORDER BY id) AS m FROM t`},
		{"setop_collation_nondefault", `SELECT c FROM cc UNION SELECT d FROM cc`},
		{"setop_collation_default_first", `SELECT d FROM cc UNION SELECT c FROM cc`},
		{"setop_collation_explicit", `SELECT c COLLATE "POSIX" AS c FROM cc UNION SELECT p FROM cp`},
		{"setop_collation_nested", `(SELECT c FROM cc UNION ALL SELECT p FROM cp) UNION ALL SELECT d FROM cc`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ddl := analyzeOracleSchema + "CREATE VIEW v AS " + tc.query + ";"
			conn := o.session(t, ddl)
			var pgDef string
			if err := conn.QueryRowContext(o.ctx, "SELECT pg_get_viewdef('v'::regclass, true)").Scan(&pgDef); err != nil {
				t.Fatal(err)
			}
			rows, err := conn.QueryContext(o.ctx, `SELECT attname, atttypid, attcollation FROM pg_attribute
				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum`)
			if err != nil {
				t.Fatal(err)
			}
			var pgCols []string
			for rows.Next() {
				var name string
				var typ, coll uint32
				if err := rows.Scan(&name, &typ, &coll); err != nil {
					t.Fatal(err)
				}
				pgCols = append(pgCols, fmt.Sprintf("%s %d %d", name, typ, coll))
			}
			rows.Close()

			c := New()
			if _, err := c.Exec(ddl, nil); err != nil {
				t.Fatalf("omni rejected what PostgreSQL accepts: %v", err)
			}
			def, err := c.GetViewDefinition("public", "v", true)
			if err != nil {
				t.Fatal(err)
			}
			if def != pgDef {
				t.Errorf("view definition differs\nomni:\n%s\npostgres:\n%s", def, pgDef)
			}
			_, rel, err := c.findRelation("public", "v")
			if err != nil {
				t.Fatal(err)
			}
			var cols []string
			for _, col := range rel.Columns {
				cols = append(cols, fmt.Sprintf("%s %d %d", col.Name, col.TypeOID, col.Collation))
			}
			if strings.Join(cols, ", ") != strings.Join(pgCols, ", ") {
				t.Errorf("view columns differ\nomni:     %v\npostgres: %v", cols, pgCols)
			}
		})
	}
}

// TestAnalyzeErrorsMatchPostgres checks that omni accepts or rejects each
// query as PostgreSQL does, with the same SQLSTATE.
func TestAnalyzeErrorsMatchPostgres(t *testing.T) {
	o := startAnalyzeOracle(t)
	conn := o.session(t, analyzeOracleSchema)
	c := New()
	if _, err := c.Exec(analyzeOracleSchema, nil); err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{
		`SELECT lower(x) FILTER (WHERE true) FROM t`,
		`SELECT lower(x ORDER BY x) FROM t`,
		`SELECT lower(DISTINCT x) FROM t`,
		`SELECT lower(x) OVER () FROM t`,
		`SELECT lower(x) WITHIN GROUP (ORDER BY x) FROM t`,
		`SELECT ntile() WITHIN GROUP (ORDER BY 3) OVER () FROM t`,
		`SELECT row_number() FROM t`,
		`SELECT row_number() FILTER (WHERE true) OVER () FROM t`,
		`SELECT sum(id ORDER BY id) OVER () FROM t`,
		`SELECT sum(id) FILTER (WHERE id > 0) OVER () FROM t`,
		`SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY id) OVER () FROM t`,
		`SELECT percentile_cont(0.5, 1.0) FROM t`,
		`SELECT string_agg(x) WITHIN GROUP (ORDER BY ',') FROM t`,
		`SELECT count(*) WITHIN GROUP (ORDER BY id) FROM t`,
		`SELECT jsonb_path_exists('{}'::jsonb, '$.a', '{}', false, 1)`,
		`SELECT * FROM t JOIN u USING (id, id)`,
		`SELECT * FROM t JOIN u USING (nope)`,
		`SELECT * FROM t JOIN u USING (masked)`,
		`SELECT * FROM (t JOIN u ON true) NATURAL JOIN u AS u2`,
		`SELECT * FROM u AS u2 NATURAL JOIN (t JOIN u ON true)`,
		`SELECT * FROM (t JOIN u ON true) JOIN u AS u2 USING (y)`,
		`SELECT id FROM t JOIN u ON t.id = u.id`,
		`SELECT id FROM t NATURAL JOIN u`,
		`SELECT masked FROM (t JOIN u USING (id)) JOIN w USING (id)`,
		`SELECT 1 FROM t JOIN u ON t.id = u.id JOIN w ON y = w.x`,
		`SELECT 1 FROM (t JOIN u USING (id)) JOIN w ON masked = w.x`,
		`SELECT 1 FROM (t JOIN u USING (id)) JOIN w ON id = w.z`,
		`SELECT 1 FROM t NATURAL JOIN u JOIN w ON y = w.x AND z > 0`,
		`SELECT 1 FROM t, u JOIN w ON y = w.x`,
		`SELECT 1 FROM t JOIN u ON y = (SELECT max(x) FROM w WHERE w.id = u.id)`,
		`SELECT pg_logical_slot_get_changes('slot', NULL, NULL)`,
		`SELECT pg_logical_slot_get_changes('slot', NULL)`,
		`SELECT poly_def(1)`,
		`SELECT poly_def(1, 'x'::text)`,
		`SELECT jsonb_path_exists('{}'::jsonb, '$.a', silent => true)`,
		`SELECT jsonb_path_exists(target => '{}'::jsonb, path => '$.a')`,
		`SELECT jsonb_path_exists('{}'::jsonb, silent => true)`,
		`SELECT jsonb_path_exists('{}'::jsonb, '$.a', target => '{}')`,
		`SELECT make_interval(days => 1, days => 2)`,
		`SELECT make_interval(days => 1, 2)`,
		`SELECT make_interval(nope => 1)`,
		`SELECT format(formatstr => 'x')`,
		`SELECT percentile_cont(fraction => 0.5) WITHIN GROUP (ORDER BY id) FROM t`,
		`SELECT my_os(0.5, id) FROM t`,
		`SELECT two_arg(id) WITHIN GROUP (ORDER BY id) FROM t`,
		`SELECT two_arg(id, id) FROM t`,
		`SELECT c FROM cc UNION SELECT p FROM cp`,
		`SELECT c FROM cc EXCEPT ALL SELECT p FROM cp`,
		`SELECT c FROM cc INTERSECT SELECT d FROM cc`,
		`SELECT c FROM cc UNION ALL SELECT p FROM cp`,
		`SELECT c COLLATE "POSIX" FROM cc UNION ALL SELECT d COLLATE "C" FROM cc`,
		`SELECT c COLLATE "C" FROM cc UNION SELECT p FROM cp`,
	} {
		t.Run(query, func(t *testing.T) {
			// EXPLAIN analyzes and plans the query without running it.
			rows, pgErr := conn.QueryContext(o.ctx, "EXPLAIN "+query)
			if pgErr == nil {
				rows.Close()
			}
			_, err := analyzeWithOptions(t, c, query, AnalyzeOptions{})
			if got, want := sqlState(t, err), sqlState(t, pgErr); got != want {
				t.Errorf("omni: %q (%v), postgres: %q (%v)", got, err, want, pgErr)
			}
		})
	}
}

// TestAnalyzeDDLMatchesPostgres checks that omni accepts or rejects each DDL
// statement as PostgreSQL does: a view column with no collation, which a
// UNION ALL of conflicting collations yields, and a hypothetical-set
// aggregate that is not ordered-set.
func TestAnalyzeDDLMatchesPostgres(t *testing.T) {
	o := startAnalyzeOracle(t)
	for _, ddl := range []string{
		`CREATE VIEW v AS SELECT c FROM cc UNION ALL SELECT p FROM cp`,
		`CREATE VIEW v AS SELECT c AS renamed FROM cc UNION ALL SELECT p FROM cp`,
		`CREATE VIEW v (renamed) AS SELECT c FROM cc UNION ALL SELECT p FROM cp`,
		`CREATE VIEW v AS SELECT c FROM cc UNION ALL SELECT d FROM cc`,
		`CREATE AGGREGATE my_bad(int) (sfunc = two_sfunc, stype = int, hypothetical)`,
		`CREATE AGGREGATE my_hyp(int ORDER BY int) (sfunc = os_sfunc, stype = int, hypothetical)`,
	} {
		t.Run(ddl, func(t *testing.T) {
			conn := o.session(t, analyzeOracleSchema)
			_, pgErr := conn.ExecContext(o.ctx, ddl)

			c := New()
			if _, err := c.Exec(analyzeOracleSchema, nil); err != nil {
				t.Fatal(err)
			}
			results, err := c.Exec(ddl, &ExecOptions{ContinueOnError: true})
			if err == nil && len(results) > 0 {
				err = results[0].Error
			}
			if got, want := sqlState(t, err), sqlState(t, pgErr); got != want {
				t.Errorf("omni: %q (%v), postgres: %q (%v)", got, err, want, pgErr)
			}
		})
	}
}

// TestAnalyzeViewDependenciesMatchPostgres checks that a view blocks DROP
// COLUMN of each column it reads only through the new clauses.
func TestAnalyzeViewDependenciesMatchPostgres(t *testing.T) {
	o := startAnalyzeOracle(t)
	for _, tc := range []struct{ name, query, drop string }{
		{"agg_filter", `SELECT sum(id) FILTER (WHERE masked = 'a') FROM t`, `ALTER TABLE t DROP COLUMN masked`},
		{"agg_order", `SELECT string_agg(x, ',' ORDER BY masked) FROM t`, `ALTER TABLE t DROP COLUMN masked`},
		{"within_group", `SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY id) FROM t`, `ALTER TABLE t DROP COLUMN id`},
		{"natural_left_input", `SELECT y FROM t NATURAL JOIN u`, `ALTER TABLE t DROP COLUMN id`},
		{"natural_right_input", `SELECT y FROM t NATURAL JOIN u`, `ALTER TABLE u DROP COLUMN id`},
		{"full_using_right", `SELECT id FROM t FULL JOIN u USING (id)`, `ALTER TABLE u DROP COLUMN id`},
		{"unrelated_column", `SELECT sum(id) FILTER (WHERE masked = 'a') FROM t`, `ALTER TABLE t DROP COLUMN x`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ddl := analyzeOracleSchema + "CREATE VIEW v AS " + tc.query + ";"
			conn := o.session(t, ddl)
			_, pgErr := conn.ExecContext(o.ctx, tc.drop)

			c := New()
			if _, err := c.Exec(ddl, nil); err != nil {
				t.Fatal(err)
			}
			results, err := c.Exec(tc.drop, &ExecOptions{ContinueOnError: true})
			if err == nil && len(results) > 0 {
				err = results[0].Error
			}
			if got, want := sqlState(t, err), sqlState(t, pgErr); got != want {
				t.Errorf("omni: %q (%v), postgres: %q (%v)", got, err, want, pgErr)
			}
		})
	}
}

// TestBuiltinProcExtrasMatchPostgres checks the default counts and parameter
// names of every built-in function, and builtinOrderedSetAggKinds, against
// pg_proc and pg_aggregate.
func TestBuiltinProcExtrasMatchPostgres(t *testing.T) {
	o := startAnalyzeOracle(t)
	byOID := make(map[uint32]*BuiltinProc, len(BuiltinProcs))
	for i := range BuiltinProcs {
		byOID[BuiltinProcs[i].OID] = &BuiltinProcs[i]
	}

	rows, err := o.db.QueryContext(o.ctx, `SELECT oid, pronargdefaults, array_to_string(proargnames, ',') FROM pg_proc WHERE oid < 16384`)
	if err != nil {
		t.Fatal(err)
	}
	seen := 0
	for rows.Next() {
		var oid uint32
		var n int16
		var names sql.NullString
		if err := rows.Scan(&oid, &n, &names); err != nil {
			t.Fatal(err)
		}
		p := byOID[oid]
		if p == nil {
			continue
		}
		seen++
		if p.NArgDefaults != n {
			t.Errorf("%s (OID %d): NArgDefaults %d, pg_proc.pronargdefaults %d", p.Name, oid, p.NArgDefaults, n)
		}
		if (p.ArgNames != nil) != names.Valid || strings.Join(p.ArgNames, ",") != names.String {
			t.Errorf("%s (OID %d): ArgNames %q, pg_proc.proargnames %q", p.Name, oid, p.ArgNames, names.String)
		}
	}
	rows.Close()
	if seen == 0 {
		t.Fatal("no built-in function matched pg_proc")
	}

	rows, err = o.db.QueryContext(o.ctx, `SELECT aggfnoid::oid, aggkind FROM pg_aggregate WHERE aggfnoid::oid < 16384`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var oid uint32
		var kind string
		if err := rows.Scan(&oid, &kind); err != nil {
			t.Fatal(err)
		}
		p := byOID[oid]
		if p == nil {
			continue
		}
		if got := builtinAggKind(p); string(got) != kind {
			t.Errorf("%s (OID %d): aggregate kind %q, pg_aggregate.aggkind %q", p.Name, oid, got, kind)
		}
	}
}
