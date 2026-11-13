package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/pg/internal/testcontainer"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// TestSimulation holds the catalog to testdata/simulate.tsv: every table,
// index, view and sequence statement PostgreSQL accepts in the parser's
// compat.tsv, run on each fixture, with what PostgreSQL did and what Exec
// does. An outcome is the SQLSTATE the statement is refused with, or ok: and
// a hash of the DDL that rebuilds what the schema holds afterwards
// (simulationForms).
//
// PostgreSQL's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned PostgreSQL again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestSimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), simulationFixtures, simulated)
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		// PostgreSQL makes up names from the ones a statement holds (a_pkey,
		// a_a_seq, a_d_idx1), the same on every run.
		MadeUp: nil,
		Omni:   simulationOmni,
		Engine: simulationEngine,
	})
}

// simulationFixtures are the starting schemas, in testdata/fixtures: one
// table, with a key and an index, under four namings. A compat.tsv statement
// writes a for every name it holds, so what a statement finds is set by what
// the fixture names a: nothing, the table, a column of it too, and an index.
// PostgreSQL keeps an index in the namespace of its table, so the index named
// a is on a table named t.
var simulationFixtures = []string{"one_table@missing", "one_table@table", "one_table@column", "one_table@index"}

// simulatedStatement matches the statements the file holds: DDL on tables,
// indexes, views, materialized views and sequences, and their comments,
// whose whole effect is in the schema a row reads back. Temporary tables are
// in no schema a row reads, so CREATE TEMP and its like are left out; so are
// schemas, types, routines, triggers, policies, roles and the objects of
// extensions, which a later fixture can take.
var simulatedStatement = regexp.MustCompile(`^(?:` +
	`CREATE (?:UNLOGGED )?TABLE |ALTER TABLE |DROP TABLE |TRUNCATE ` +
	`|CREATE (?:UNIQUE )?INDEX |ALTER INDEX |DROP INDEX ` +
	`|CREATE (?:OR REPLACE )?(?:RECURSIVE )?VIEW |ALTER VIEW |DROP VIEW ` +
	`|CREATE MATERIALIZED VIEW |ALTER MATERIALIZED VIEW |DROP MATERIALIZED VIEW |REFRESH MATERIALIZED VIEW ` +
	`|CREATE (?:UNLOGGED )?SEQUENCE |ALTER SEQUENCE |DROP SEQUENCE ` +
	`|COMMENT ON (?:TABLE|COLUMN|INDEX|VIEW|MATERIALIZED VIEW|SEQUENCE) ` +
	`)`)

func simulated(statement string) bool { return simulatedStatement.MatchString(statement) }

// simulationEngine asks PostgreSQL, or the tape. Each row is three requests,
// none of which reads omni: the schema public dropped and built again with
// the fixture, the statement, and after one PostgreSQL took, the schema's
// objects as DDL. The test's database holds nothing else.
func simulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	db := testcontainer.Open(t)
	ctx := context.Background()
	return func(fx simulate.Fixture, statement string) (simulate.Answer, error) {
		if _, err := db.ExecContext(ctx, "DROP SCHEMA IF EXISTS public CASCADE;\nCREATE SCHEMA public;\n"+fx.SQL); err != nil {
			return simulate.Answer{}, fmt.Errorf("build fixture %s: %w", fx.Name, err)
		}
		if _, err := db.ExecContext(ctx, statement); err != nil {
			var refusal *pgconn.PgError
			if !errors.As(err, &refusal) {
				return simulate.Answer{}, err
			}
			return simulate.Answer{Code: refusal.Code}, nil
		}
		return simulationForms(ctx, db)
	}
}

// simulationFormsQuery writes each object of the schema public as the DDL
// that builds it, statements ending in ";\n", in name order:
//
//   - a sequence, but one an identity column makes, with its options;
//   - a table, unlogged or partitioned, as CREATE TABLE with its local
//     columns (type, collation, default or generation, identity, NOT NULL),
//     parents and partition key, or as a partition of its parent; then its
//     local constraints, the indexes no constraint and no parent index
//     makes, its comments, and the sequences its columns own;
//   - a view or materialized view, with its options, query, indexes and
//     comments.
//
// Storage, statistics, ownership, privileges, replica identity and row
// security are not written, so a statement that changes only those stores
// the schema it found. Every name is the stored one, quoted as needed.
const simulationFormsQuery = `WITH rel AS (
	SELECT c.oid, c.relname, c.relkind, c.relpersistence, c.relispartition, c.relispopulated, c.relpartbound,
		coalesce(' WITH (' || array_to_string(c.reloptions, ', ') || ')', '') AS options
	FROM pg_catalog.pg_class c
	WHERE c.relnamespace = 'public'::pg_catalog.regnamespace AND c.relkind IN ('r', 'p', 'v', 'm', 'S')
		AND NOT EXISTS (SELECT 1 FROM pg_catalog.pg_depend d
			WHERE d.classid = 'pg_catalog.pg_class'::pg_catalog.regclass AND d.objid = c.oid AND d.deptype = 'i')
), comments AS (
	SELECT r.oid, coalesce((SELECT format(E'COMMENT ON %s %I IS %L;\n',
				CASE r.relkind WHEN 'v' THEN 'VIEW' WHEN 'm' THEN 'MATERIALIZED VIEW' WHEN 'S' THEN 'SEQUENCE' ELSE 'TABLE' END,
				r.relname, pg_catalog.obj_description(r.oid, 'pg_class'))
			WHERE pg_catalog.obj_description(r.oid, 'pg_class') IS NOT NULL), '')
		|| coalesce((SELECT string_agg(format(E'COMMENT ON COLUMN %I.%I IS %L;\n', r.relname, a.attname, pg_catalog.col_description(r.oid, a.attnum)), '' ORDER BY a.attnum)
			FROM pg_catalog.pg_attribute a WHERE a.attrelid = r.oid AND a.attnum > 0 AND NOT a.attisdropped AND pg_catalog.col_description(r.oid, a.attnum) IS NOT NULL), '')
		|| coalesce((SELECT string_agg(format(E'COMMENT ON INDEX %I IS %L;\n', i.relname, pg_catalog.obj_description(i.oid, 'pg_class')), '' ORDER BY i.relname)
			FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
			WHERE x.indrelid = r.oid AND pg_catalog.obj_description(i.oid, 'pg_class') IS NOT NULL), '') AS text
	FROM rel r
), indexes AS (
	SELECT r.oid, coalesce(string_agg(pg_catalog.pg_get_indexdef(x.indexrelid) || E';\n', '' ORDER BY i.relname), '') AS text
	FROM rel r JOIN pg_catalog.pg_index x ON x.indrelid = r.oid JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
	WHERE NOT EXISTS (SELECT 1 FROM pg_catalog.pg_constraint k WHERE k.conindid = x.indexrelid AND k.contype IN ('p', 'u', 'x') AND k.conrelid = r.oid)
		AND NOT EXISTS (SELECT 1 FROM pg_catalog.pg_inherits h WHERE h.inhrelid = x.indexrelid)
	GROUP BY r.oid
)
SELECT kind, relname, text FROM (
	SELECT 'SEQUENCE' AS kind, r.relname, format(E'CREATE %sSEQUENCE %I AS %s INCREMENT BY %s MINVALUE %s MAXVALUE %s START WITH %s CACHE %s%s;\n',
			CASE r.relpersistence WHEN 'u' THEN 'UNLOGGED ' ELSE '' END, r.relname, pg_catalog.format_type(q.seqtypid, NULL),
			q.seqincrement, q.seqmin, q.seqmax, q.seqstart, q.seqcache, CASE WHEN q.seqcycle THEN ' CYCLE' ELSE '' END) || m.text AS text
		FROM rel r JOIN pg_catalog.pg_sequence q ON q.seqrelid = r.oid JOIN comments m ON m.oid = r.oid
		WHERE r.relkind = 'S'
	UNION ALL
	SELECT 'TABLE', r.relname,
		CASE WHEN r.relispartition THEN format(E'CREATE TABLE %I PARTITION OF %I %s', r.relname,
				(SELECT p.relname FROM pg_catalog.pg_inherits h JOIN pg_catalog.pg_class p ON p.oid = h.inhparent WHERE h.inhrelid = r.oid),
				pg_catalog.pg_get_expr(r.relpartbound, r.oid))
			ELSE format('CREATE %sTABLE %I (%s)', CASE r.relpersistence WHEN 'u' THEN 'UNLOGGED ' ELSE '' END, r.relname,
				coalesce((SELECT string_agg(format('%I %s', a.attname, pg_catalog.format_type(a.atttypid, a.atttypmod))
						|| CASE WHEN a.attcollation <> 0 AND a.attcollation <> t.typcollation
							THEN ' COLLATE ' || (SELECT format('%I', l.collname) FROM pg_catalog.pg_collation l WHERE l.oid = a.attcollation) ELSE '' END
						|| CASE WHEN a.attgenerated = 's' THEN ' GENERATED ALWAYS AS (' || pg_catalog.pg_get_expr(d.adbin, d.adrelid) || ') STORED'
							WHEN d.adbin IS NOT NULL THEN ' DEFAULT ' || pg_catalog.pg_get_expr(d.adbin, d.adrelid) ELSE '' END
						|| CASE a.attidentity WHEN 'a' THEN ' GENERATED ALWAYS AS IDENTITY' WHEN 'd' THEN ' GENERATED BY DEFAULT AS IDENTITY' ELSE '' END
						|| CASE WHEN a.attnotnull THEN ' NOT NULL' ELSE '' END, ', ' ORDER BY a.attnum)
					FROM pg_catalog.pg_attribute a JOIN pg_catalog.pg_type t ON t.oid = a.atttypid
					LEFT JOIN pg_catalog.pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
					WHERE a.attrelid = r.oid AND a.attnum > 0 AND NOT a.attisdropped AND a.attislocal), ''))
				|| coalesce((SELECT ' INHERITS (' || string_agg(format('%I', p.relname), ', ' ORDER BY h.inhseqno) || ')'
					FROM pg_catalog.pg_inherits h JOIN pg_catalog.pg_class p ON p.oid = h.inhparent WHERE h.inhrelid = r.oid), '')
		END
		|| CASE WHEN r.relkind = 'p' THEN ' PARTITION BY ' || pg_catalog.pg_get_partkeydef(r.oid) ELSE '' END
		|| r.options || E';\n'
		|| coalesce((SELECT string_agg(format(E'ALTER TABLE %I ADD CONSTRAINT %I %s;\n', r.relname, k.conname, pg_catalog.pg_get_constraintdef(k.oid)), '' ORDER BY k.conname)
			FROM pg_catalog.pg_constraint k WHERE k.conrelid = r.oid AND k.conislocal AND k.contype IN ('c', 'f', 'p', 'u', 'x')), '')
		|| coalesce(x.text, '') || m.text
		|| coalesce((SELECT string_agg(format(E'ALTER SEQUENCE %I OWNED BY %I.%I;\n', s.relname, r.relname, a.attname), '' ORDER BY s.relname)
			FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_class s ON s.oid = d.objid
			JOIN pg_catalog.pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
			WHERE d.classid = 'pg_catalog.pg_class'::pg_catalog.regclass AND d.refobjid = r.oid AND d.deptype = 'a' AND s.relkind = 'S'), '')
		FROM rel r JOIN comments m ON m.oid = r.oid LEFT JOIN indexes x ON x.oid = r.oid
		WHERE r.relkind IN ('r', 'p')
	UNION ALL
	SELECT CASE r.relkind WHEN 'v' THEN 'VIEW' ELSE 'MATERIALIZED VIEW' END, r.relname,
		format(E'CREATE %s %I%s AS %s%s;\n', CASE r.relkind WHEN 'v' THEN 'VIEW' ELSE 'MATERIALIZED VIEW' END, r.relname, r.options,
			rtrim(btrim(pg_catalog.pg_get_viewdef(r.oid)), ';'),
			CASE WHEN r.relkind = 'm' AND NOT r.relispopulated THEN ' WITH NO DATA' ELSE '' END)
		|| coalesce(x.text, '') || m.text
		FROM rel r JOIN comments m ON m.oid = r.oid LEFT JOIN indexes x ON x.oid = r.oid
		WHERE r.relkind IN ('v', 'm')
) o ORDER BY relname COLLATE "C", kind`

// simulationForms reads the schema public's objects as simulationFormsQuery
// writes them.
func simulationForms(ctx context.Context, db *sql.DB) (simulate.Answer, error) {
	rows, err := db.QueryContext(ctx, simulationFormsQuery)
	if err != nil {
		return simulate.Answer{}, fmt.Errorf("read the stored schema: %w", err)
	}
	defer rows.Close()
	var forms []simulate.Form
	for rows.Next() {
		var f simulate.Form
		if err := rows.Scan(&f.Kind, &f.Name, &f.Text); err != nil {
			return simulate.Answer{}, err
		}
		forms = append(forms, f)
	}
	return simulate.Answer{Forms: forms}, rows.Err()
}

// simulationOmni runs the row in the catalog. The catalog stores what
// PostgreSQL stores when Diff finds nothing between it and a catalog loaded
// from PostgreSQL's forms: what the forms do not write, such as storage or
// ownership, does not count.
func simulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	c, err := LoadSQL(fx.SQL)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	results, err := c.Exec(statement, nil)
	var syntax *pgparser.ParseError
	switch {
	case errors.As(err, &syntax):
		code := syntax.Code
		if code == "" {
			code = "42601"
		}
		return simulate.Omni{Code: code, Why: err.Error()}
	case err != nil:
		return simulate.Omni{Code: simulate.NoCode, Why: err.Error()}
	}
	for _, r := range results {
		var refusal *Error
		switch {
		case r.Error == nil:
		case errors.As(r.Error, &refusal) && refusal.Code != "":
			return simulate.Omni{Code: refusal.Code, Why: r.Error.Error()}
		default:
			return simulate.Omni{Code: simulate.NoCode, Why: r.Error.Error()}
		}
	}
	if engine.Refused() {
		return simulate.Omni{}
	}
	want, err := loadForms(engine.Forms)
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what PostgreSQL stores: %v", err)}
	}
	diff := Diff(c, want)
	if diff.IsEmpty() {
		return simulate.Omni{Same: true}
	}
	return simulate.Omni{Why: "to store PostgreSQL's schema the catalog would run\n" + GenerateMigration(c, want, diff).SQL()}
}

// loadForms builds a catalog from PostgreSQL's forms. The forms are in name
// order, which is not the order they build in: sequences go first, then each
// statement in turn, again and again while one that failed for want of
// another object now runs, as a partition does before its parent.
func loadForms(forms []simulate.Form) (*Catalog, error) {
	c := New()
	var first, rest []string
	for _, f := range forms {
		stmts := strings.SplitAfter(f.Text, ";\n")
		if stmts[len(stmts)-1] == "" {
			stmts = stmts[:len(stmts)-1]
		}
		if f.Kind == "SEQUENCE" {
			first = append(first, stmts...)
		} else {
			rest = append(rest, stmts...)
		}
	}
	for _, s := range first {
		if err := execOne(c, s); err != nil {
			return nil, fmt.Errorf("%s: %w", strings.TrimSpace(s), err)
		}
	}
	for len(rest) > 0 {
		var failed []string
		var last error
		for _, s := range rest {
			if err := execOne(c, s); err != nil {
				failed = append(failed, s)
				last = fmt.Errorf("%s: %w", strings.TrimSpace(s), err)
			}
		}
		if len(failed) == len(rest) {
			return nil, last
		}
		rest = failed
	}
	return c, nil
}

// execOne runs one statement in the catalog and returns its error.
func execOne(c *Catalog, stmt string) error {
	results, err := c.Exec(stmt, nil)
	if err != nil {
		return err
	}
	for _, r := range results {
		if r.Error != nil {
			return r.Error
		}
	}
	return nil
}
