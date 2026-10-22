package review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// coverageSchema has one object of each kind the statement forms name,
// and no trigger, rule, partition, or routine, which would end the walk
// at DML.
func coverageSchema() *metadata.DatabaseSchemaMetadata {
	cols := func(names ...string) []*metadata.ColumnMetadata {
		var out []*metadata.ColumnMetadata
		for _, n := range names {
			out = append(out, &metadata.ColumnMetadata{Name: n})
		}
		return out
	}
	return &metadata.DatabaseSchemaMetadata{
		Name:       "db",
		SearchPath: `"$user", public`,
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadata.TableMetadata{
					{
						Name:             "t",
						Columns:          cols("id", "code", "n"),
						Indexes:          []*metadata.IndexMetadata{{Name: "t_pkey", Primary: true, Unique: true, IsConstraint: true}, {Name: "t_n_idx"}},
						CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "t_n_check"}},
					},
					{Name: "nokey", Columns: cols("a")},
					{Name: "parted", Columns: cols("k")},
				},
				Views:             []*metadata.ViewMetadata{{Name: "tv"}},
				MaterializedViews: []*metadata.MaterializedViewMetadata{{Name: "mv"}},
				Sequences:         []*metadata.SequenceMetadata{{Name: "q"}},
				EnumTypes:         []*metadata.EnumTypeMetadata{{Name: "mood"}},
				CompositeTypes:    []*metadata.CompositeTypeMetadata{{Name: "pair"}},
			},
			{Name: "s", Tables: []*metadata.TableMetadata{{Name: "st", Columns: cols("a")}}},
		},
	}
}

// walkForms are the statement forms of the PG spec's §5.2 table, each on
// coverageSchema. A form the table lists must lower: one the parser gives
// in a shape the lowering does not expect would end the walk unseen.
var walkForms = []walktest.Form{
	// Schemas.
	{SQL: "CREATE SCHEMA s2;"},
	{SQL: "CREATE SCHEMA IF NOT EXISTS s;"},
	{SQL: "DROP SCHEMA s CASCADE;"},
	{SQL: "DROP SCHEMA IF EXISTS s2;"},
	{SQL: "ALTER SCHEMA s RENAME TO s3;"},
	// Tables.
	{SQL: "CREATE TABLE x (a int PRIMARY KEY, b text UNIQUE, c int REFERENCES t (id), CHECK (a > 0));"},
	{SQL: "CREATE TABLE IF NOT EXISTS t (a int);"},
	{SQL: "CREATE TEMP TABLE tmp (a int);"},
	{SQL: "CREATE TABLE x (LIKE t);"},
	{SQL: "CREATE TABLE x (b int) INHERITS (nokey);"},
	{SQL: "CREATE TABLE x2 (k int) PARTITION BY RANGE (k);"},
	{SQL: "CREATE TABLE x3 (k int) PARTITION BY RANGE (k);\nCREATE TABLE x3_1 PARTITION OF x3 FOR VALUES FROM (0) TO (10);"},
	{SQL: "CREATE TABLE x AS SELECT 1 AS a;"},
	{SQL: "SELECT 1 AS a INTO x;"},
	{SQL: "CREATE MATERIALIZED VIEW mv2 AS SELECT 1 AS a;"},
	// Indexes, views, sequences, types.
	{SQL: "CREATE INDEX i ON t (n);"},
	{SQL: "CREATE UNIQUE INDEX IF NOT EXISTS i2 ON t (code);"},
	{SQL: "CREATE INDEX ON t (n);"},
	{SQL: "CREATE VIEW v AS SELECT 1 AS a;"},
	{SQL: "CREATE OR REPLACE VIEW tv AS SELECT code FROM t;"},
	{SQL: "CREATE SEQUENCE q2;"},
	{SQL: "CREATE TYPE e AS ENUM ('a');"},
	{SQL: "CREATE TYPE c AS (a int);"},
	{SQL: "CREATE DOMAIN d AS int;"},
	// ALTER TABLE.
	{SQL: "ALTER TABLE t ADD COLUMN z int;"},
	{SQL: "ALTER TABLE t ADD COLUMN IF NOT EXISTS z int;"},
	{SQL: "ALTER TABLE t DROP COLUMN n;"},
	{SQL: "ALTER TABLE t ALTER COLUMN n SET NOT NULL;"},
	{SQL: "ALTER TABLE t ALTER COLUMN n TYPE bigint;"},
	{SQL: "ALTER TABLE t ADD CONSTRAINT c1 CHECK (n > 0);"},
	{SQL: "ALTER TABLE t DROP CONSTRAINT t_n_check;"},
	{SQL: "ALTER TABLE t VALIDATE CONSTRAINT t_n_check;"},
	{SQL: "ALTER TABLE t OWNER TO bob;"},
	{SQL: "ALTER TABLE t SET (fillfactor = 70);"},
	{SQL: "ALTER TABLE t ADD COLUMN y int, DROP COLUMN n;"},
	// Renames and moves.
	{SQL: "ALTER TABLE t RENAME COLUMN n TO m;"},
	{SQL: "ALTER TABLE t RENAME CONSTRAINT t_n_check TO c2;"},
	{SQL: "ALTER TABLE t RENAME TO t2;"},
	{SQL: "ALTER TABLE t SET SCHEMA s;"},
	{SQL: "ALTER INDEX t_n_idx RENAME TO i3;"},
	{SQL: "ALTER SEQUENCE q RENAME TO q3;"},
	{SQL: "ALTER VIEW tv RENAME TO tv2;"},
	{SQL: "ALTER MATERIALIZED VIEW mv RENAME TO mv3;"},
	{SQL: "ALTER VIEW tv SET SCHEMA s;"},
	{SQL: "ALTER SEQUENCE q RESTART;"},
	// Drops.
	{SQL: "DROP TABLE nokey;"},
	{SQL: "DROP TABLE IF EXISTS x, y CASCADE;"},
	{SQL: "DROP VIEW tv;"},
	{SQL: "DROP MATERIALIZED VIEW mv;"},
	{SQL: "DROP INDEX t_n_idx;"},
	{SQL: "DROP INDEX CONCURRENTLY t_n_idx;"},
	{SQL: "DROP SEQUENCE q;"},
	{SQL: "DROP TYPE public.mood;"},
	{SQL: "DROP TYPE IF EXISTS mood;"},
	{SQL: "DROP DOMAIN IF EXISTS public.dd;"},
	// Statements that change no name.
	{SQL: "COMMENT ON TABLE t IS 'x';"},
	{SQL: "COMMENT ON COLUMN t.n IS 'x';"},
	{SQL: "COMMENT ON SCHEMA s IS 'x';"},
	{SQL: "GRANT SELECT ON t TO PUBLIC;"},
	{SQL: "TRUNCATE nokey;"},
	{SQL: "LOCK TABLE t;"},
	{SQL: "INSERT INTO nokey VALUES (1);"},
	{SQL: "UPDATE nokey SET a = 1;"},
	{SQL: "DELETE FROM nokey;"},
	{SQL: "SELECT 1;"},
	// Session and transactions.
	{SQL: "SET search_path = public;"},
	{SQL: "RESET search_path;"},
	{SQL: "SET ROLE bob;\nRESET ROLE;"},
	{SQL: "SELECT pg_catalog.set_config('search_path', '', false);"},
	{SQL: "BEGIN;\nSAVEPOINT a;\nRELEASE SAVEPOINT a;\nCOMMIT;"},
	// Routines and triggers: recorded, then the walk goes on.
	{SQL: "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT 1';"},
	{SQL: "CREATE TRIGGER tr BEFORE INSERT ON nokey FOR EACH ROW EXECUTE FUNCTION f();"},
	// A refresh runs stored code: lowered, and the walk ends after it.
	{SQL: "REFRESH MATERIALIZED VIEW mv;", Last: true},
}

func TestWalkThroughLowersEveryForm(t *testing.T) {
	walktest.Coverage(t, walkForms, func(t *testing.T, sql string) (n, lowered int, last bool) {
		t.Helper()
		stmts, syntax := parse(sql, statementRanges(sql))
		if syntax != nil {
			t.Fatalf("syntax: %s\n%s", syntax.Message, sql)
		}
		w := newWalker(stmts, review.Target{Schema: coverageSchema(), SessionUser: "u"})
		l := walk.NewLedger(w.pg)
		for i := range stmts {
			s := w.lower(i, l)
			if s == nil {
				return len(stmts), i, false
			}
			if s.Last {
				return len(stmts), i + 1, true
			}
			if s.Apply != nil {
				s.Apply(l)
			}
		}
		return len(stmts), len(stmts), false
	})
}
