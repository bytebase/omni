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
// column comments, triggers, constraints and constraint indexes.
type dependencyOracleCase struct {
	name string
	sql  string
}

func dependencyOracleCases() []dependencyOracleCase {
	return []dependencyOracleCase{
		// DROP COLUMN checks and drops the dependents of that column only.
		{"drop column beside a check on another column", `
CREATE TABLE t (id int PRIMARY KEY, n int CHECK (n > 0));
ALTER TABLE t DROP COLUMN id;`},
		{"drop column takes its check", `
CREATE TABLE t (id int PRIMARY KEY, n int CHECK (n > 0));
ALTER TABLE t DROP COLUMN n;
ALTER TABLE t ADD CONSTRAINT t_n_check CHECK (id > 0);`},
		{"drop column takes a check on two columns", `
CREATE TABLE t (a int, b int, CHECK (a < b));
ALTER TABLE t DROP COLUMN a;`},
		{"key dropped with its column", `
CREATE TABLE t (id int PRIMARY KEY, code text UNIQUE, n int);
ALTER TABLE t DROP COLUMN id;
ALTER TABLE t DROP COLUMN code;`},
		{"drop column takes expression and partial indexes", `
CREATE TABLE t (a int, name text);
CREATE INDEX t_lower ON t (lower(name));
CREATE INDEX t_part ON t (a) WHERE name IS NOT NULL;
ALTER TABLE t DROP COLUMN name;`},
		{"drop column takes its serial sequence", `
CREATE TABLE t (a int, b serial);
ALTER TABLE t DROP COLUMN b;`},
		{"drop column refuses a generated column", `
CREATE TABLE t (a int, b int GENERATED ALWAYS AS (a * 2) STORED);
ALTER TABLE t DROP COLUMN a;`},
		{"drop column cascades to a generated column", `
CREATE TABLE t (a int, b int GENERATED ALWAYS AS (a * 2) STORED, c int);
ALTER TABLE t DROP COLUMN a CASCADE;`},
		{"drop generated column then its source", `
CREATE TABLE t (a int, b int GENERATED ALWAYS AS (a * 2) STORED);
ALTER TABLE t DROP COLUMN b;
ALTER TABLE t DROP COLUMN a;`},
		{"drop column ignores a view on other columns", `
CREATE TABLE t (a int, b int);
CREATE VIEW v AS SELECT a FROM t;
ALTER TABLE t DROP COLUMN b;`},
		{"drop column cascades to the views reading it", `
CREATE TABLE t (a int, b int);
CREATE VIEW v AS SELECT a FROM t;
CREATE VIEW vv AS SELECT a FROM v;
CREATE VIEW w AS SELECT b FROM t;
ALTER TABLE t DROP COLUMN a;
ALTER TABLE t DROP COLUMN a CASCADE;`},
		{"view reads a column in ORDER BY", `
CREATE TABLE t (a int, b int);
CREATE VIEW v AS SELECT a FROM t ORDER BY b;
ALTER TABLE t DROP COLUMN b;`},
		{"view reads a column in HAVING", `
CREATE TABLE t (a int, b int);
CREATE VIEW v AS SELECT a, count(*) FROM t GROUP BY a HAVING max(b) > 0;
ALTER TABLE t DROP COLUMN b;`},
		{"view reads a column in JOIN ON", `
CREATE TABLE u (id int, name text);
CREATE TABLE o (id int, uid int);
CREATE VIEW v AS SELECT u.name FROM u JOIN o ON o.uid = u.id;
ALTER TABLE o DROP COLUMN uid;
ALTER TABLE o DROP COLUMN id;`},
		{"view reads a column in a correlated subquery", `
CREATE TABLE u (id int, name text);
CREATE TABLE o (id int, uid int);
CREATE VIEW v AS SELECT u.name FROM u WHERE EXISTS (SELECT 1 FROM o WHERE o.uid = u.id);
ALTER TABLE u DROP COLUMN id;
ALTER TABLE o DROP COLUMN id;`},
		{"view reads a column left of IN", `
CREATE TABLE u (id int, name text);
CREATE TABLE o (id int, uid int);
CREATE VIEW v AS SELECT u.name FROM u WHERE u.id IN (SELECT uid FROM o);
ALTER TABLE u DROP COLUMN id;`},
		{"view reads a column through a join", `
CREATE TABLE t (a int, x int);
CREATE TABLE u (b int, x int);
CREATE VIEW v AS SELECT a FROM t JOIN u ON true;
CREATE VIEW w AS SELECT x FROM t JOIN u USING (x);
CREATE VIEW z AS SELECT b FROM (t JOIN u USING (x)) j;
ALTER TABLE t DROP COLUMN a;
ALTER TABLE u DROP COLUMN b;
ALTER TABLE t DROP COLUMN x;
ALTER TABLE u DROP COLUMN x;
ALTER TABLE t ALTER COLUMN a TYPE bigint;
DROP VIEW v, z;
ALTER TABLE t DROP COLUMN a;
ALTER TABLE u ALTER COLUMN b TYPE bigint;`},
		{"view reads every column of a subquery", `
CREATE TABLE t (a int, b int);
CREATE VIEW v AS SELECT x FROM (SELECT a AS x, b AS y FROM t) s;
ALTER TABLE t DROP COLUMN b;`},
		{"drop column refuses a foreign key on it", `
CREATE TABLE t (id int PRIMARY KEY, n int);
CREATE TABLE c (tid int REFERENCES t (id));
ALTER TABLE t DROP COLUMN id;
ALTER TABLE t DROP COLUMN n;
ALTER TABLE t DROP COLUMN id CASCADE;`},
		{"drop column of a self-referencing key", `
CREATE TABLE t (id int PRIMARY KEY, parent int REFERENCES t (id));
ALTER TABLE t DROP COLUMN id;
ALTER TABLE t DROP COLUMN parent;`},
		{"drop column refuses a trigger on it", `
CREATE TABLE t (a int, b int);
CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$;
CREATE TRIGGER tg BEFORE UPDATE OF b ON t FOR EACH ROW EXECUTE FUNCTION f();
ALTER TABLE t DROP COLUMN b;
ALTER TABLE t DROP COLUMN a;`},
		{"drop columns one after another", `
CREATE TABLE t (a int, b int, c int, d int, e int);
CREATE VIEW v AS SELECT c FROM t;
CREATE UNIQUE INDEX t_e ON t (e);
COMMENT ON COLUMN t.b IS 'b';
COMMENT ON COLUMN t.c IS 'c';
COMMENT ON COLUMN t.e IS 'e';
ALTER TABLE t DROP COLUMN a;
ALTER TABLE t DROP COLUMN b;
ALTER TABLE t DROP COLUMN c;
ALTER TABLE t DROP COLUMN d;
ALTER TABLE t DROP COLUMN e;`},
		{"drop column of an inheritance parent", `
CREATE TABLE p (id int, x int, y int);
CREATE TABLE c (y int) INHERITS (p);
ALTER TABLE c DROP COLUMN x;
ALTER TABLE c DROP COLUMN y;
ALTER TABLE p DROP COLUMN x;
ALTER TABLE p DROP COLUMN y;
ALTER TABLE c DROP COLUMN y;`},
		{"drop column of a partitioned table", `
CREATE TABLE p (id int, x int, y int) PARTITION BY RANGE (id);
CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10);
CREATE TABLE p2 (id int, x int, y int);
ALTER TABLE p ATTACH PARTITION p2 FOR VALUES FROM (10) TO (20);
ALTER TABLE p1 DROP COLUMN x;
ALTER TABLE p DROP COLUMN x;`},

		{"drop column refused by a view on an inheritance child", `
CREATE TABLE p (id int, a int);
CREATE TABLE c () INHERITS (p);
CREATE VIEW v AS SELECT a FROM c;
ALTER TABLE p DROP COLUMN a;`},
		{"drop column refused by a view on a partition", `
CREATE TABLE p (id int, a int) PARTITION BY RANGE (id);
CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10);
CREATE VIEW v AS SELECT a FROM p1;
ALTER TABLE p DROP COLUMN a;
ALTER TABLE p DROP COLUMN a CASCADE;`},
		{"drop column of a partitioned table with an index and a trigger", `
CREATE TABLE p (id int, a int, b int) PARTITION BY RANGE (id);
CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10);
CREATE INDEX p_a ON p (a);
CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$;
CREATE TRIGGER tg BEFORE UPDATE OF b ON p FOR EACH ROW EXECUTE FUNCTION f();
ALTER TABLE p DROP COLUMN a;
ALTER TABLE p DROP COLUMN b;
ALTER TABLE p DROP COLUMN b CASCADE;`},
		{"drop column through an inheritance tree", `
CREATE TABLE p (id int, a int);
CREATE TABLE c () INHERITS (p);
CREATE TABLE g () INHERITS (c);
CREATE TABLE q (id int, a int);
CREATE TABLE q1 () INHERITS (q);
CREATE TABLE q2 () INHERITS (q1);
CREATE VIEW v AS SELECT a FROM q2;
ALTER TABLE q DROP COLUMN a;
ALTER TABLE p DROP COLUMN a;`},
		{"drop column through sub-partitions", `
CREATE TABLE p (id int, a int, b int) PARTITION BY RANGE (id);
CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10) PARTITION BY RANGE (b);
CREATE TABLE p11 PARTITION OF p1 FOR VALUES FROM (0) TO (5);
ALTER TABLE p DROP COLUMN a;`},
		{"drop column cascades to a partition trigger's WHEN", `
CREATE TABLE p (id int, a int, b int) PARTITION BY RANGE (id);
CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$;
CREATE TRIGGER tg BEFORE INSERT ON p FOR EACH ROW WHEN (NEW.b > 0) EXECUTE FUNCTION f();
CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10);
ALTER TABLE p DROP COLUMN b;
ALTER TABLE p DROP COLUMN b CASCADE;`},
		{"drop column from only the parent", `
CREATE TABLE p (id int, a int);
CREATE TABLE c () INHERITS (p);
CREATE TABLE p2 (id int, a int);
CREATE TABLE c2 (a int) INHERITS (p2);
ALTER TABLE ONLY p DROP COLUMN a;
ALTER TABLE ONLY p2 DROP COLUMN a;
ALTER TABLE c DROP COLUMN a;
ALTER TABLE c2 DROP COLUMN a;
CREATE TABLE q (id int, a int) PARTITION BY RANGE (id);
CREATE TABLE q1 PARTITION OF q FOR VALUES FROM (0) TO (10);
ALTER TABLE ONLY q DROP COLUMN a;
CREATE TABLE r (id int, a int) PARTITION BY RANGE (id);
ALTER TABLE ONLY r DROP COLUMN a;`},
		{"drop column of a sub-partition key", `
CREATE TABLE p (a int, b int) PARTITION BY RANGE (a);
CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10) PARTITION BY RANGE (b);
ALTER TABLE p DROP COLUMN b;
ALTER TABLE p DROP COLUMN b CASCADE;
ALTER TABLE p DROP COLUMN a;`},
		{"drop referenced column cascades to partition foreign keys", `
CREATE TABLE t (id int PRIMARY KEY, k int);
CREATE TABLE p (id int, tid int REFERENCES t (id)) PARTITION BY RANGE (id);
CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10);
ALTER TABLE t DROP COLUMN id;
ALTER TABLE t DROP COLUMN id CASCADE;`},
		// DROP TABLE takes the table's own constraints and columns with it.
		{"drop table with checks, generated columns and a self key", `
CREATE TABLE t (id int PRIMARY KEY, n int CHECK (n > 0), m int GENERATED ALWAYS AS (n * 2) STORED, parent int REFERENCES t (id));
DROP TABLE t;`},
		{"drop table with a column trigger", `
CREATE TABLE t (a int, b int);
CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$;
CREATE TRIGGER tg BEFORE UPDATE OF b ON t FOR EACH ROW EXECUTE FUNCTION f();
DROP TABLE t;`},
		{"drop table cascades to a view on a generated column", `
CREATE TABLE t (a int, b int GENERATED ALWAYS AS (a * 2) STORED);
CREATE VIEW v AS SELECT b FROM t;
DROP TABLE t;
DROP TABLE t CASCADE;`},

		// ALTER COLUMN TYPE rebuilds constraints but refuses views,
		// generated columns and triggers reading the column.
		{"alter type of a referenced key", `
CREATE TABLE t (id int PRIMARY KEY);
CREATE TABLE c (tid int REFERENCES t (id));
ALTER TABLE t ALTER COLUMN id TYPE bigint;`},
		{"alter column type", `
CREATE TABLE t (id int, n int CHECK (n > 0), m int, g int GENERATED ALWAYS AS (m * 2) STORED, k int, w int);
CREATE VIEW v AS SELECT k FROM t;
CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$;
CREATE TRIGGER tg BEFORE UPDATE OF w ON t FOR EACH ROW EXECUTE FUNCTION f();
ALTER TABLE t ALTER COLUMN id TYPE bigint;
ALTER TABLE t ALTER COLUMN n TYPE bigint;
ALTER TABLE t ALTER COLUMN m TYPE bigint;
ALTER TABLE t ALTER COLUMN k TYPE bigint;
ALTER TABLE t ALTER COLUMN w TYPE bigint;`},

		// A CHECK reading one column is named after it.
		{"check constraint names", `
CREATE TABLE t (c int CHECK (c > 0), d int, CHECK (d > 0), CHECK (c < d));
ALTER TABLE t DROP CONSTRAINT t_c_check;
ALTER TABLE t DROP CONSTRAINT t_d_check;
ALTER TABLE t DROP CONSTRAINT t_check;`},

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

// pgSchemaState lists the relations, columns with their comments, triggers
// and constraints of schema on PostgreSQL, sorted.
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
SELECT 'trigger ' || r.relname || '.' || tg.tgname
FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
     JOIN pg_namespace n ON n.oid = r.relnamespace
WHERE n.nspname = $1 AND NOT tg.tgisinternal
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
		for _, trig := range c.triggersByRel[rel.OID] {
			state = append(state, "trigger "+name+"."+trig.Name)
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
