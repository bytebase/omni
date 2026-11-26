package catalog_test

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
	"google.golang.org/protobuf/proto"

	"github.com/bytebase/omni/cockroachdb/catalog"
	"github.com/bytebase/omni/cockroachdb/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/metadata"
)

// TestSimulation holds Exec to testdata/simulate.tsv: each statement of
// simulationStatements, run on each fixture, with what CockroachDB did and
// what Exec does. An outcome is the SQLSTATE the statement is refused with, or
// ok: and a hash of the DDL CockroachDB's SHOW CREATE writes for the database
// afterwards.
//
// CockroachDB's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned CockroachDB again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestSimulation(t *testing.T) {
	var rows []simulate.Generated
	for _, fx := range simulationFixtures {
		for _, s := range simulationStatements {
			rows = append(rows, simulate.Generated{Fixture: fx, Rule: s.rule, Statement: s.sql})
		}
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		// CockroachDB makes up names from the ones a statement holds (t_pkey,
		// t_a_idx, check_a, rowid), the same on every run.
		MadeUp: nil,
		Omni:   simulationOmni,
		Engine: simulationEngine,
	})
}

// simulationFixtures are the starting schemas, in testdata/fixtures: none,
// one table t, and objects, which holds a schema, an enum and a composite
// type, a sequence, tables with keys, checks, a unique constraint, an index,
// a foreign key and comments, a table with no primary key, a view and a
// materialized view.
var simulationFixtures = []string{"empty", "table", "objects"}

// simulationStatements are the DDL forms Bytebase's SDL migration writes and
// the ones users write by hand, on the names the fixtures hold and on names
// none holds, so each meets its object and its error cases.
var simulationStatements = []struct{ rule, sql string }{
	// Schemas.
	{"schema/create", "CREATE SCHEMA s"},
	{"schema/create", "CREATE SCHEMA IF NOT EXISTS s"},
	{"schema/create", "CREATE SCHEMA s2"},
	{"schema/create", `CREATE SCHEMA "Mixed Case"`},
	{"schema/alter", "ALTER SCHEMA s RENAME TO s2"},
	{"schema/alter", "ALTER SCHEMA s RENAME TO public"},
	{"schema/alter", "ALTER SCHEMA missing RENAME TO s2"},
	{"schema/drop", "DROP SCHEMA s"},
	{"schema/drop", "DROP SCHEMA s CASCADE"},
	{"schema/drop", "DROP SCHEMA IF EXISTS s"},
	{"schema/drop", "DROP SCHEMA IF EXISTS missing"},
	{"schema/drop", "DROP SCHEMA missing"},
	{"schema/drop", "DROP SCHEMA public"},
	{"schema/comment", "COMMENT ON SCHEMA s IS 'a schema'"},
	{"schema/comment", "COMMENT ON SCHEMA public IS 'the default'"},

	// Tables.
	{"table/create", "CREATE TABLE t (id INT8 PRIMARY KEY)"},
	{"table/create", "CREATE TABLE IF NOT EXISTS t (id INT8 PRIMARY KEY)"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, v STRING NOT NULL)"},
	{"table/create", "CREATE TABLE n (x INT8, y STRING)"},
	{"table/create", "CREATE TABLE n (x INT8 NOT NULL, y INT8 NOT NULL, PRIMARY KEY (x, y DESC))"},
	{"table/create", "CREATE TABLE n (id INT8, CONSTRAINT n_key PRIMARY KEY (id))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 DEFAULT 7, b STRING DEFAULT 'x', c TIMESTAMPTZ DEFAULT now())"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 CHECK (a > 0), b INT8, CHECK (b < 10), CONSTRAINT n_ab CHECK (a < b))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, e STRING UNIQUE, f STRING, UNIQUE (f))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, a INT8, INDEX (a), INDEX n_named (a DESC))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, price DECIMAL(10, 2), cents INT8 AS ((price * 100)::INT8) STORED, twice DECIMAL AS (price * 2) VIRTUAL)"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, secret STRING NOT VISIBLE, tags STRING[], doc JSONB, name VARCHAR(20) COLLATE \"en_US\")"},
	{"table/create", "CREATE TABLE n (id SERIAL PRIMARY KEY, v STRING)"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES t (id))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES t)"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8, CONSTRAINT n_t FOREIGN KEY (t_id) REFERENCES t (id) ON DELETE CASCADE ON UPDATE SET NULL)"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES missing (id))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, k_x INT8 REFERENCES k (x))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, m mood DEFAULT 'sad', p pair)"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, m missing_type)"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, v INT8 DEFAULT nextval('sq'))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, id INT8)"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 CHECK (missing > 0))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY, INDEX (missing))"},
	{"table/create", "CREATE TABLE n (id INT8 PRIMARY KEY) WITH (ttl_expire_after = '1 day')"},
	{"table/create", "CREATE TABLE s.n (id INT8 PRIMARY KEY)"},
	{"table/create", "CREATE TABLE s.u (id INT8 PRIMARY KEY)"},
	{"table/create", "CREATE TABLE missing.n (id INT8 PRIMARY KEY)"},
	{"table/create", "CREATE TABLE v (id INT8 PRIMARY KEY)"},
	{"table/create", "CREATE TABLE sq (id INT8 PRIMARY KEY)"},
	{"table/create", `CREATE TABLE "Mixed" ("Id" INT8 PRIMARY KEY, "Some Value" STRING)`},
	{"table/create", "CREATE TABLE n AS SELECT 1 AS x"},
	{"table/create", "CREATE TEMP TABLE n (id INT8 PRIMARY KEY)"},
	{"table/alter", "ALTER TABLE t RENAME TO t2"},
	{"table/alter", "ALTER TABLE t RENAME TO c"},
	{"table/alter", "ALTER TABLE IF EXISTS t RENAME TO t2"},
	{"table/alter", "ALTER TABLE missing RENAME TO t2"},
	{"table/alter", "ALTER TABLE IF EXISTS missing RENAME TO t2"},
	{"table/alter", "ALTER TABLE t SET SCHEMA s"},
	{"table/alter", "ALTER TABLE k SET SCHEMA s"},
	{"table/alter", "ALTER TABLE s.u RENAME TO u2"},
	{"table/drop", "DROP TABLE t"},
	{"table/drop", "DROP TABLE t CASCADE"},
	{"table/drop", "DROP TABLE IF EXISTS t"},
	{"table/drop", "DROP TABLE IF EXISTS missing"},
	{"table/drop", "DROP TABLE missing"},
	{"table/drop", "DROP TABLE c"},
	{"table/drop", "DROP TABLE c, t"},
	{"table/drop", "DROP TABLE k"},
	{"table/drop", "DROP TABLE s.u"},
	{"table/drop", "DROP TABLE v"},
	{"table/truncate", "TRUNCATE t"},
	{"table/comment", "COMMENT ON TABLE t IS 'changed'"},
	{"table/comment", "COMMENT ON TABLE t IS NULL"},
	{"table/comment", "COMMENT ON TABLE missing IS 'x'"},

	// Columns.
	{"column/add", "ALTER TABLE t ADD COLUMN n INT8"},
	{"column/add", "ALTER TABLE t ADD n INT8"},
	{"column/add", "ALTER TABLE t ADD COLUMN n INT8 NOT NULL DEFAULT 5"},
	{"column/add", "ALTER TABLE t ADD COLUMN n INT8 NOT NULL"},
	{"column/add", "ALTER TABLE t ADD COLUMN n STRING DEFAULT 'x' CHECK (n <> '')"},
	{"column/add", "ALTER TABLE t ADD COLUMN n STRING UNIQUE"},
	{"column/add", "ALTER TABLE t ADD COLUMN n INT8 REFERENCES t (id)"},
	{"column/add", "ALTER TABLE t ADD COLUMN n INT8 AS (a * 2) STORED"},
	{"column/add", "ALTER TABLE t ADD COLUMN n INT8 AS (a * 2) VIRTUAL"},
	{"column/add", "ALTER TABLE t ADD COLUMN n mood"},
	{"column/add", "ALTER TABLE t ADD COLUMN n INT8 NOT VISIBLE"},
	{"column/add", "ALTER TABLE t ADD COLUMN n INT8 DEFAULT nextval('sq')"},
	{"column/add", "ALTER TABLE t ADD COLUMN a INT8"},
	{"column/add", "ALTER TABLE t ADD COLUMN IF NOT EXISTS a INT8"},
	{"column/add", "ALTER TABLE t ADD COLUMN IF NOT EXISTS n INT8"},
	{"column/add", "ALTER TABLE t ADD COLUMN n INT8, ADD COLUMN o STRING"},
	{"column/add", "ALTER TABLE missing ADD COLUMN n INT8"},
	{"column/add", "ALTER TABLE k ADD COLUMN n INT8"},
	{"column/add", "ALTER TABLE v ADD COLUMN n INT8"},
	{"column/drop", "ALTER TABLE t DROP COLUMN b"},
	{"column/drop", "ALTER TABLE t DROP b"},
	{"column/drop", "ALTER TABLE t DROP COLUMN a"},
	{"column/drop", "ALTER TABLE t DROP COLUMN a CASCADE"},
	{"column/drop", "ALTER TABLE t DROP COLUMN m"},
	{"column/drop", "ALTER TABLE t DROP COLUMN id"},
	{"column/drop", "ALTER TABLE t DROP COLUMN missing"},
	{"column/drop", "ALTER TABLE t DROP COLUMN IF EXISTS missing"},
	{"column/drop", "ALTER TABLE t DROP COLUMN IF EXISTS b"},
	{"column/drop", "ALTER TABLE c DROP COLUMN t_id"},
	{"column/drop", "ALTER TABLE c DROP COLUMN note"},
	{"column/drop", "ALTER TABLE k DROP COLUMN y"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN a SET DEFAULT 9"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN b SET DEFAULT 'x'"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN a DROP DEFAULT"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN b DROP DEFAULT"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN missing SET DEFAULT 1"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN b SET NOT NULL"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN a DROP NOT NULL"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN id DROP NOT NULL"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN missing SET NOT NULL"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN b TYPE STRING(100)"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN b SET DATA TYPE VARCHAR(50)"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN a TYPE INT4"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN a TYPE STRING"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN missing TYPE STRING"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN b SET NOT VISIBLE"},
	{"column/alter", "ALTER TABLE t ALTER COLUMN b SET VISIBLE"},
	{"column/rename", "ALTER TABLE t RENAME COLUMN b TO z"},
	{"column/rename", "ALTER TABLE t RENAME b TO z"},
	{"column/rename", "ALTER TABLE t RENAME COLUMN a TO z"},
	{"column/rename", "ALTER TABLE t RENAME COLUMN id TO z"},
	{"column/rename", "ALTER TABLE t RENAME COLUMN b TO a"},
	{"column/rename", "ALTER TABLE t RENAME COLUMN missing TO z"},
	{"column/rename", "ALTER TABLE c RENAME COLUMN t_id TO parent"},
	{"column/comment", "COMMENT ON COLUMN t.b IS 'bee'"},
	{"column/comment", "COMMENT ON COLUMN t.a IS NULL"},
	{"column/comment", "COMMENT ON COLUMN t.missing IS 'x'"},

	// Constraints.
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_pos CHECK (a > -1)"},
	{"constraint/add", "ALTER TABLE t ADD CHECK (a < 1000)"},
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_pos CHECK (a > -1) NOT VALID"},
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_a_check CHECK (a > 1)"},
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_bad CHECK (missing > 0)"},
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_b_unique UNIQUE (b)"},
	{"constraint/add", "ALTER TABLE t ADD UNIQUE (a, b)"},
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_b_key UNIQUE (b)"},
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_bad UNIQUE (missing)"},
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_self FOREIGN KEY (a) REFERENCES t (id)"},
	{"constraint/add", "ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES t (id) ON DELETE CASCADE"},
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_bad FOREIGN KEY (a) REFERENCES missing (id)"},
	{"constraint/add", "ALTER TABLE t ADD CONSTRAINT t_bad FOREIGN KEY (a) REFERENCES t (b)"},
	{"constraint/add", "ALTER TABLE c ADD CONSTRAINT c_note CHECK (note <> '')"},
	{"constraint/add", "ALTER TABLE k ADD PRIMARY KEY (x)"},
	{"constraint/add", "ALTER TABLE t ADD PRIMARY KEY (a)"},
	{"constraint/drop", "ALTER TABLE t DROP CONSTRAINT t_a_check"},
	{"constraint/drop", "ALTER TABLE t DROP CONSTRAINT IF EXISTS t_a_check"},
	{"constraint/drop", "ALTER TABLE t DROP CONSTRAINT IF EXISTS missing"},
	{"constraint/drop", "ALTER TABLE t DROP CONSTRAINT missing"},
	{"constraint/drop", "ALTER TABLE t DROP CONSTRAINT t_b_key"},
	{"constraint/drop", "ALTER TABLE t DROP CONSTRAINT t_pkey"},
	{"constraint/drop", "ALTER TABLE c DROP CONSTRAINT c_t_id_fkey"},
	{"constraint/alter", "ALTER TABLE t RENAME CONSTRAINT t_a_check TO t_a_nonneg"},
	{"constraint/alter", "ALTER TABLE t RENAME CONSTRAINT t_b_key TO t_b_uniq"},
	{"constraint/alter", "ALTER TABLE t RENAME CONSTRAINT t_pkey TO t_key"},
	{"constraint/alter", "ALTER TABLE t RENAME CONSTRAINT missing TO z"},
	{"constraint/alter", "ALTER TABLE t VALIDATE CONSTRAINT t_a_check"},
	{"constraint/alter", "ALTER TABLE c VALIDATE CONSTRAINT c_t_id_fkey"},
	{"constraint/alter", "ALTER TABLE t ALTER PRIMARY KEY USING COLUMNS (a)"},
	{"constraint/alter", "ALTER TABLE k ALTER PRIMARY KEY USING COLUMNS (x)"},
	{"constraint/comment", "COMMENT ON CONSTRAINT t_a_check ON t IS 'non-negative'"},
	{"constraint/comment", "COMMENT ON CONSTRAINT t_pkey ON t IS 'the key'"},
	{"constraint/comment", "COMMENT ON CONSTRAINT missing ON t IS 'x'"},

	// Indexes.
	{"index/create", "CREATE INDEX ON t (b)"},
	{"index/create", "CREATE INDEX t_b_idx ON t (b)"},
	{"index/create", "CREATE INDEX t_ab ON t (a DESC, b)"},
	{"index/create", "CREATE UNIQUE INDEX t_b_uniq ON t (b)"},
	{"index/create", "CREATE INDEX t_b_store ON t (b) STORING (a)"},
	{"index/create", "CREATE INDEX t_b_part ON t (b) WHERE a > 0"},
	{"index/create", "CREATE INDEX t_lower ON t (lower(b))"},
	{"index/create", "CREATE INDEX t_hash ON t (b) USING HASH"},
	{"index/create", "CREATE INDEX t_hidden ON t (b) NOT VISIBLE"},
	{"index/create", "CREATE INVERTED INDEX t_doc ON t (doc)"},
	{"index/create", "CREATE INDEX t_a_idx ON t (b)"},
	{"index/create", "CREATE INDEX IF NOT EXISTS t_a_idx ON t (b)"},
	{"index/create", "CREATE INDEX IF NOT EXISTS t_b_idx ON t (b)"},
	{"index/create", "CREATE INDEX t_bad ON t (missing)"},
	{"index/create", "CREATE INDEX n_idx ON missing (a)"},
	{"index/create", "CREATE INDEX ON k (x)"},
	{"index/create", "CREATE INDEX ON s.u (v)"},
	{"index/create", "CREATE INDEX mv_b ON mv (b)"},
	{"index/create", "CREATE INDEX v_a ON v (a)"},
	{"index/create", "CREATE INDEX CONCURRENTLY t_b_conc ON t (b)"},
	{"index/alter", "ALTER INDEX t@t_a_idx RENAME TO t_a_idx2"},
	{"index/alter", "ALTER INDEX t_a_idx RENAME TO t_a_idx2"},
	{"index/alter", "ALTER INDEX t_a_idx RENAME TO t_b_key"},
	{"index/alter", "ALTER INDEX missing RENAME TO x"},
	{"index/alter", "ALTER INDEX IF EXISTS missing RENAME TO x"},
	{"index/alter", "ALTER INDEX t@t_a_idx NOT VISIBLE"},
	{"index/alter", "ALTER INDEX t@t_a_idx VISIBLE"},
	{"index/alter", "ALTER TABLE t RENAME INDEX t_a_idx TO t_a_idx2"},
	{"index/drop", "DROP INDEX t_a_idx"},
	{"index/drop", "DROP INDEX t@t_a_idx"},
	{"index/drop", "DROP INDEX IF EXISTS t_a_idx"},
	{"index/drop", "DROP INDEX IF EXISTS missing"},
	{"index/drop", "DROP INDEX missing"},
	{"index/drop", "DROP INDEX t_b_key"},
	{"index/drop", "DROP INDEX t_b_key CASCADE"},
	{"index/drop", "DROP INDEX t@t_pkey"},
	{"index/drop", "DROP INDEX CONCURRENTLY t_a_idx"},
	{"index/comment", "COMMENT ON INDEX t@t_a_idx IS 'by a'"},
	{"index/comment", "COMMENT ON INDEX t_a_idx IS 'by a'"},
	{"index/comment", "COMMENT ON INDEX missing IS 'x'"},

	// Sequences.
	{"sequence/create", "CREATE SEQUENCE sq"},
	{"sequence/create", "CREATE SEQUENCE IF NOT EXISTS sq"},
	{"sequence/create", "CREATE SEQUENCE sq2"},
	{"sequence/create", "CREATE SEQUENCE sq2 AS INT4 INCREMENT BY 2 MINVALUE 5 MAXVALUE 500 START WITH 7 CACHE 10"},
	{"sequence/create", "CREATE SEQUENCE sq2 INCREMENT BY -1"},
	{"sequence/create", "CREATE SEQUENCE sq2 CYCLE"},
	{"sequence/create", "CREATE SEQUENCE sq2 OWNED BY t.a"},
	{"sequence/create", "CREATE SEQUENCE sq2 OWNED BY t.missing"},
	{"sequence/create", "CREATE SEQUENCE s.sq2"},
	{"sequence/create", "CREATE SEQUENCE t"},
	{"sequence/create", "CREATE SEQUENCE sq2 MINVALUE 10 MAXVALUE 5"},
	{"sequence/alter", "ALTER SEQUENCE sq INCREMENT BY 5"},
	{"sequence/alter", "ALTER SEQUENCE sq MINVALUE 2 MAXVALUE 100"},
	{"sequence/alter", "ALTER SEQUENCE sq NO MAXVALUE"},
	{"sequence/alter", "ALTER SEQUENCE sq START WITH 50"},
	{"sequence/alter", "ALTER SEQUENCE sq RESTART WITH 50"},
	{"sequence/alter", "ALTER SEQUENCE sq CACHE 20"},
	{"sequence/alter", "ALTER SEQUENCE sq CYCLE"},
	{"sequence/alter", "ALTER SEQUENCE sq AS INT4"},
	{"sequence/alter", "ALTER SEQUENCE sq OWNED BY t.a"},
	{"sequence/alter", "ALTER SEQUENCE sq OWNED BY NONE"},
	{"sequence/alter", "ALTER SEQUENCE missing INCREMENT BY 5"},
	{"sequence/alter", "ALTER SEQUENCE IF EXISTS missing INCREMENT BY 5"},
	{"sequence/alter", "ALTER SEQUENCE sq RENAME TO sq2"},
	{"sequence/alter", "ALTER SEQUENCE sq RENAME TO t"},
	{"sequence/alter", "ALTER SEQUENCE sq SET SCHEMA s"},
	{"sequence/drop", "DROP SEQUENCE sq"},
	{"sequence/drop", "DROP SEQUENCE sq CASCADE"},
	{"sequence/drop", "DROP SEQUENCE IF EXISTS sq"},
	{"sequence/drop", "DROP SEQUENCE IF EXISTS missing"},
	{"sequence/drop", "DROP SEQUENCE missing"},
	{"sequence/drop", "DROP SEQUENCE t"},

	// Views.
	{"view/create", "CREATE VIEW v AS SELECT id FROM t"},
	{"view/create", "CREATE VIEW IF NOT EXISTS v AS SELECT id FROM t"},
	{"view/create", "CREATE VIEW v2 AS SELECT id, a FROM t"},
	{"view/create", "CREATE VIEW v2 AS SELECT * FROM t"},
	{"view/create", "CREATE VIEW v2 (x, y) AS SELECT id, a FROM t"},
	{"view/create", "CREATE VIEW v2 AS SELECT id, a + 1 AS a1, count(*) OVER () AS n FROM t WHERE a > 0"},
	{"view/create", "CREATE VIEW v2 AS SELECT id FROM v"},
	{"view/create", "CREATE VIEW v2 AS SELECT t.id, c.note FROM t JOIN c ON c.t_id = t.id"},
	{"view/create", "CREATE VIEW v2 AS SELECT missing FROM t"},
	{"view/create", "CREATE VIEW v2 AS SELECT 1 AS one"},
	{"view/create", "CREATE VIEW v2 AS SELECT id FROM missing"},
	{"view/create", "CREATE VIEW s.v2 AS SELECT id FROM t"},
	{"view/create", "CREATE VIEW t AS SELECT 1 AS one"},
	{"view/create", "CREATE OR REPLACE VIEW v AS SELECT id, a FROM t"},
	{"view/create", "CREATE OR REPLACE VIEW v AS SELECT id, a, b FROM t"},
	{"view/create", "CREATE OR REPLACE VIEW v AS SELECT id FROM t"},
	{"view/create", "CREATE OR REPLACE VIEW v AS SELECT a, id FROM t"},
	{"view/create", "CREATE OR REPLACE VIEW v2 AS SELECT id FROM t"},
	{"view/create", "CREATE MATERIALIZED VIEW mv2 AS SELECT id, a FROM t"},
	{"view/create", "CREATE MATERIALIZED VIEW mv AS SELECT id FROM t"},
	{"view/create", "CREATE MATERIALIZED VIEW IF NOT EXISTS mv AS SELECT id FROM t"},
	{"view/create", "CREATE MATERIALIZED VIEW mv2 AS SELECT id FROM t WITH NO DATA"},
	{"view/alter", "ALTER VIEW v RENAME TO v2"},
	{"view/alter", "ALTER VIEW v RENAME TO t"},
	{"view/alter", "ALTER VIEW missing RENAME TO v2"},
	{"view/alter", "ALTER VIEW IF EXISTS missing RENAME TO v2"},
	{"view/alter", "ALTER VIEW v SET SCHEMA s"},
	{"view/alter", "ALTER MATERIALIZED VIEW mv RENAME TO mv2"},
	{"view/drop", "DROP VIEW v"},
	{"view/drop", "DROP VIEW v CASCADE"},
	{"view/drop", "DROP VIEW IF EXISTS v"},
	{"view/drop", "DROP VIEW IF EXISTS missing"},
	{"view/drop", "DROP VIEW missing"},
	{"view/drop", "DROP VIEW t"},
	{"view/drop", "DROP VIEW mv"},
	{"view/drop", "DROP MATERIALIZED VIEW mv"},
	{"view/drop", "DROP MATERIALIZED VIEW v"},
	{"view/drop", "DROP MATERIALIZED VIEW IF EXISTS mv"},
	{"view/comment", "COMMENT ON VIEW v IS 'x'"},

	// Types.
	{"type/create", "CREATE TYPE mood AS ENUM ('a')"},
	{"type/create", "CREATE TYPE IF NOT EXISTS mood AS ENUM ('a')"},
	{"type/create", "CREATE TYPE color AS ENUM ('red', 'green', 'blue')"},
	{"type/create", "CREATE TYPE color AS ENUM ()"},
	{"type/create", "CREATE TYPE color AS ENUM ('red', 'red')"},
	{"type/create", "CREATE TYPE s.color AS ENUM ('red')"},
	{"type/create", "CREATE TYPE missing.color AS ENUM ('red')"},
	{"type/create", "CREATE TYPE point2 AS (x FLOAT8, y FLOAT8)"},
	{"type/create", "CREATE TYPE point2 AS (m mood, label STRING)"},
	{"type/create", "CREATE TYPE pair AS (x INT8)"},
	{"type/create", "CREATE TYPE t AS ENUM ('a')"},
	{"type/alter", "ALTER TYPE mood ADD VALUE 'meh'"},
	{"type/alter", "ALTER TYPE mood ADD VALUE 'meh' BEFORE 'sad'"},
	{"type/alter", "ALTER TYPE mood ADD VALUE 'meh' AFTER 'happy'"},
	{"type/alter", "ALTER TYPE mood ADD VALUE 'happy'"},
	{"type/alter", "ALTER TYPE mood ADD VALUE IF NOT EXISTS 'happy'"},
	{"type/alter", "ALTER TYPE mood ADD VALUE 'meh' BEFORE 'missing'"},
	{"type/alter", "ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'"},
	{"type/alter", "ALTER TYPE mood RENAME VALUE 'missing' TO 'blue'"},
	{"type/alter", "ALTER TYPE mood RENAME VALUE 'sad' TO 'happy'"},
	{"type/alter", "ALTER TYPE mood DROP VALUE 'sad'"},
	{"type/alter", "ALTER TYPE mood DROP VALUE 'happy'"},
	{"type/alter", "ALTER TYPE mood DROP VALUE 'missing'"},
	{"type/alter", "ALTER TYPE mood RENAME TO feeling"},
	{"type/alter", "ALTER TYPE mood RENAME TO pair"},
	{"type/alter", "ALTER TYPE mood SET SCHEMA s"},
	{"type/alter", "ALTER TYPE missing ADD VALUE 'x'"},
	{"type/alter", "ALTER TYPE pair ADD VALUE 'x'"},
	{"type/drop", "DROP TYPE mood"},
	{"type/drop", "DROP TYPE pair"},
	{"type/drop", "DROP TYPE IF EXISTS mood"},
	{"type/drop", "DROP TYPE IF EXISTS missing"},
	{"type/drop", "DROP TYPE missing"},
	{"type/drop", "DROP TYPE mood, pair"},
	{"type/drop", "DROP TYPE t"},
	{"type/comment", "COMMENT ON TYPE mood IS 'feelings'"},
	{"type/comment", "COMMENT ON TYPE pair IS 'two'"},
	{"type/comment", "COMMENT ON TYPE missing IS 'x'"},
}

// simulationEngine asks CockroachDB, or the tape. Each row is a conversation
// that reads nothing of omni: the database sim dropped and made again with the
// fixture, the statement, and after one CockroachDB took, SHOW CREATE ALL
// SCHEMAS, TYPES and TABLES. They run on one connection, which USE leaves in
// sim.
func simulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	db := testcontainer.Open(t)
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return func(fx simulate.Fixture, statement string) (simulate.Answer, error) {
		if _, err := conn.ExecContext(ctx, "USE defaultdb;\nDROP DATABASE IF EXISTS sim CASCADE;\nCREATE DATABASE sim;\nUSE sim;\n"+fx.SQL); err != nil {
			return simulate.Answer{}, fmt.Errorf("build fixture %s: %w", fx.Name, err)
		}
		if _, err := conn.ExecContext(ctx, statement); err != nil {
			var refusal *pgconn.PgError
			if !errors.As(err, &refusal) {
				return simulate.Answer{}, err
			}
			return simulate.Answer{Code: refusal.Code}, nil
		}
		var forms []simulate.Form
		for _, kind := range []string{"SCHEMAS", "TYPES", "TABLES"} {
			read, err := showCreate(ctx, conn, kind)
			if err != nil {
				return simulate.Answer{}, err
			}
			forms = append(forms, read...)
		}
		return simulate.Answer{Forms: forms}, nil
	}
}

// showCreate reads one SHOW CREATE ALL answer, a statement or a few to a row,
// in the order CockroachDB writes them, which builds the database.
func showCreate(ctx context.Context, conn *sql.Conn, kind string) ([]simulate.Form, error) {
	rows, err := conn.QueryContext(ctx, "SELECT create_statement FROM [SHOW CREATE ALL "+kind+"]")
	if err != nil {
		return nil, fmt.Errorf("SHOW CREATE ALL %s: %w", kind, err)
	}
	defer rows.Close()
	var forms []simulate.Form
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			return nil, err
		}
		name := ""
		if m := createdName.FindStringSubmatch(text); m != nil {
			name = m[1]
		}
		forms = append(forms, simulate.Form{Kind: strings.TrimSuffix(kind, "S"), Name: name, Text: text})
	}
	return forms, rows.Err()
}

var createdName = regexp.MustCompile(`^CREATE (?:SCHEMA|TYPE|TABLE|VIEW|MATERIALIZED VIEW|SEQUENCE) (\S+)`)

// unwritten is what the catalog holds without what SHOW CREATE does not
// write: the column a sequence is owned by, and the indexes of a
// materialized view.
func unwritten(meta *metadata.DatabaseSchemaMetadata) *metadata.DatabaseSchemaMetadata {
	meta = proto.Clone(meta).(*metadata.DatabaseSchemaMetadata)
	for _, s := range meta.Schemas {
		for _, sequence := range s.Sequences {
			sequence.OwnerTable, sequence.OwnerColumn = "", ""
		}
		for _, view := range s.MaterializedViews {
			view.Indexes = nil
		}
	}
	return meta
}

// simulationOmni runs the row in the catalog. The catalog stores what
// CockroachDB stores when the migration from it to CockroachDB's schema, read
// back with LoadSDL, has no statement: what that migration does not write
// does not count.
func simulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	c := catalog.New()
	results, err := c.Exec(fx.SQL)
	if err == nil {
		for _, r := range results {
			if r.Error != nil {
				err = fmt.Errorf("line %d: %s: %w", r.Line, r.SQL, r.Error)
			}
		}
	}
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	results, err = c.Exec(statement)
	if err != nil {
		var syntax *pgconn.PgError
		if errors.As(err, &syntax) && syntax.Code != "" {
			return simulate.Omni{Code: syntax.Code, Why: err.Error()}
		}
		return simulate.Omni{Code: "42601", Why: err.Error()}
	}
	for _, r := range results {
		var refusal *catalog.Error
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
	var text strings.Builder
	for _, f := range engine.Forms {
		text.WriteString(f.Text)
		text.WriteString("\n")
	}
	want, err := catalog.LoadSDL(text.String())
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what CockroachDB stores: %v", err)}
	}
	plan, err := catalog.MetadataMigrationPlan(unwritten(c.Metadata()), want, catalog.MetadataMigrationOptions{})
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("no migration to CockroachDB's schema: %v", err)}
	}
	if len(plan.Ops) == 0 {
		return simulate.Omni{Same: true}
	}
	return simulate.Omni{Why: "to store CockroachDB's schema the catalog would run\n" + plan.SQL()}
}
