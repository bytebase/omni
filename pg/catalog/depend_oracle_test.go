package catalog

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

// dependencyOracleCase is a script that changes objects others depend on.
// It runs statement by statement on PostgreSQL 17 and on the catalog: each
// statement must be accepted by both, or rejected by both with the same
// SQLSTATE, and the schema must end with the same relations, columns,
// column comments, constraints and constraint indexes.
type dependencyOracleCase struct {
	name string
	sql  string
}

func dependencyOracleCases() []dependencyOracleCase {
	return []dependencyOracleCase{
		// An index and the constraint it backs share their name.
		{"renamed index renames its constraint", `
CREATE TABLE t (id int PRIMARY KEY, code text CONSTRAINT t_code_key UNIQUE);
ALTER INDEX t_code_key RENAME TO t_code_uk;
ALTER INDEX t_pkey RENAME TO t_pk;
ALTER TABLE t DROP CONSTRAINT t_code_uk;`},
		{"renamed constraint renames its index", `
CREATE TABLE t (id int PRIMARY KEY, code text CONSTRAINT t_code_key UNIQUE);
ALTER TABLE t RENAME CONSTRAINT t_code_key TO t_code_uk;
DROP INDEX t_code_key;
ALTER INDEX t_code_uk RENAME TO t_code_uk2;
ALTER TABLE t RENAME CONSTRAINT t_pkey TO t_pk;`},
		{"renamed exclusion constraint and index", `
CREATE TABLE t (r int4range, EXCLUDE USING gist (r WITH &&));
ALTER INDEX t_r_excl RENAME TO t_ex;
ALTER TABLE t RENAME CONSTRAINT t_ex TO t_ex2;`},
		{"renamed partition key indexes", `
CREATE TABLE p (id int, PRIMARY KEY (id)) PARTITION BY RANGE (id);
CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10);
ALTER INDEX p1_pkey RENAME TO p1_pk;
ALTER INDEX p_pkey RENAME TO p_pk;`},
		{"rename collisions", `
CREATE TABLE t (id int PRIMARY KEY, code text CONSTRAINT t_code_key UNIQUE, CONSTRAINT foo CHECK (id > 0));
CREATE INDEX t_id_idx ON t (id);
CREATE SEQUENCE s;
CREATE TABLE u (x int);
ALTER INDEX t_code_key RENAME TO foo;
ALTER INDEX t_code_key RENAME TO s;
ALTER INDEX t_code_key RENAME TO t_id_idx;
ALTER TABLE t RENAME CONSTRAINT t_code_key TO u;
ALTER TABLE t RENAME CONSTRAINT t_code_key TO t_pkey;
ALTER TABLE t RENAME CONSTRAINT t_code_key TO foo;
ALTER TABLE t RENAME CONSTRAINT foo TO t_pkey;
ALTER TABLE t RENAME CONSTRAINT nope TO bar;
ALTER INDEX t_id_idx RENAME TO t_id_idx2;`},
	}
}

func TestDependencyPG17Oracle(t *testing.T) {
	oracle := startLoaderCompatOracle(t)
	for i, tc := range dependencyOracleCases() {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			results, err := c.Exec(tc.sql, &ExecOptions{ContinueOnError: true})
			if err != nil {
				t.Fatalf("catalog Exec: %v", err)
			}

			schema := fmt.Sprintf("dep_oracle_%d", i)
			conn, err := oracle.db.Conn(oracle.ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			for _, q := range []string{
				fmt.Sprintf(`DROP SCHEMA IF EXISTS %q CASCADE`, schema),
				fmt.Sprintf(`CREATE SCHEMA %q`, schema),
				fmt.Sprintf(`SET search_path TO %q`, schema),
			} {
				if _, err := conn.ExecContext(oracle.ctx, q); err != nil {
					t.Fatalf("PG setup: %v", err)
				}
			}
			t.Cleanup(func() {
				_, _ = oracle.db.ExecContext(oracle.ctx, fmt.Sprintf(`DROP SCHEMA IF EXISTS %q CASCADE`, schema))
			})

			for _, r := range results {
				_, pgErr := conn.ExecContext(oracle.ctx, r.SQL)
				want, got := pgSQLState(pgErr), catalogSQLState(r.Error)
				if want != got {
					t.Errorf("statement %d %q:\n  PG17:    %s %v\n  catalog: %s %v",
						r.Index, strings.TrimSpace(r.SQL), want, pgErr, got, r.Error)
				}
			}

			pgState, err := pgSchemaState(oracle, conn, schema)
			if err != nil {
				t.Fatal(err)
			}
			if got := catalogSchemaState(c); strings.Join(got, "\n") != strings.Join(pgState, "\n") {
				t.Errorf("schema state differs\nPG17:\n  %s\ncatalog:\n  %s",
					strings.Join(pgState, "\n  "), strings.Join(got, "\n  "))
			}
		})
	}
}

// pgSQLState returns "ok" for nil, else the SQLSTATE of a PostgreSQL error.
func pgSQLState(err error) string {
	if err == nil {
		return "ok"
	}
	var pe *pgconn.PgError
	if errors.As(err, &pe) {
		return pe.Code
	}
	return "?"
}

// catalogSQLState returns "ok" for nil, else the SQLSTATE of a catalog error.
func catalogSQLState(err error) string {
	if err == nil {
		return "ok"
	}
	var ce *Error
	if errors.As(err, &ce) {
		return ce.Code
	}
	return "?"
}

// pgSchemaState lists the relations, columns with their comments, and
// constraints of schema on PostgreSQL, sorted.
func pgSchemaState(oracle *loaderCompatOracle, conn *sql.Conn, schema string) ([]string, error) {
	const q = `
SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
            WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
            ELSE 'relation ' || c.relname END
FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
UNION ALL
SELECT 'column ' || c.relname || '.' || a.attname
       || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
     JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
      AND a.attnum > 0 AND NOT a.attisdropped
UNION ALL
SELECT 'constraint ' || r.relname || '.' || con.conname
       || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
     JOIN pg_namespace n ON n.oid = r.relnamespace
     LEFT JOIN pg_class i ON i.oid = con.conindid
WHERE n.nspname = $1`
	rows, err := conn.QueryContext(oracle.ctx, q, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var state []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return nil, err
		}
		state = append(state, line)
	}
	sort.Strings(state)
	return state, rows.Err()
}

// catalogSchemaState lists the catalog's public schema the way pgSchemaState
// lists PostgreSQL's.
func catalogSchemaState(c *Catalog) []string {
	schema := c.schemaByName["public"]
	var state []string
	for name, rel := range schema.Relations {
		if rel.RelKind == 'c' {
			continue
		}
		state = append(state, "relation "+name)
		for _, col := range rel.Columns {
			line := "column " + name + "." + col.Name
			if text, ok := c.GetComment('r', rel.OID, col.AttNum); ok {
				line += " comment " + text
			}
			state = append(state, line)
		}
		for _, con := range c.consByRel[rel.OID] {
			line := "constraint " + name + "." + con.Name
			if con.Type == ConstraintPK || con.Type == ConstraintUnique || con.Type == ConstraintExclude {
				if idx := c.indexes[con.IndexOID]; idx != nil {
					line += " index " + idx.Name
				}
			}
			state = append(state, line)
		}
	}
	for name := range schema.Indexes {
		state = append(state, "index "+name)
	}
	for name := range schema.Sequences {
		state = append(state, "sequence "+name)
	}
	sort.Strings(state)
	return state
}
