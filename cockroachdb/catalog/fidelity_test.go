package catalog_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cockroachdb/analysis"
	"github.com/bytebase/omni/cockroachdb/catalog"
	"github.com/bytebase/omni/cockroachdb/driver"
	"github.com/bytebase/omni/cockroachdb/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/metadata"
)

// TestMain accounts for this package's tapes and removes the CockroachDB a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}

// fidelityCase is one schema that must survive a trip through CockroachDB:
// the server stores it, the driver syncs it, the catalog renders it, the
// server runs that rendering, and the catalog renders the same definition
// again. Every case lives in schemas of its own, so dropping them empties the
// database.
type fidelityCase struct {
	name string
	sql  string
}

func fidelityCases() []fidelityCase {
	return []fidelityCase{
		{
			name: "columns-keys-indexes",
			sql: `CREATE SCHEMA fid_cols;
COMMENT ON SCHEMA fid_cols IS 'columns';
CREATE TYPE fid_cols.mood AS ENUM ('happy', 'sad', 'it''s complicated');
CREATE TYPE fid_cols.address AS (street STRING, city VARCHAR(50), zip INT8);
CREATE TABLE fid_cols.customers (
    id INT8 PRIMARY KEY DEFAULT unique_rowid(),
    alt_id INT8 UNIQUE,
    email STRING NOT NULL,
    name VARCHAR(50) COLLATE "en_US" NOT NULL DEFAULT 'anon',
    mood fid_cols.mood DEFAULT 'happy',
    home fid_cols.address,
    tags STRING[] DEFAULT ARRAY[]::STRING[],
    balance DECIMAL(12, 2) DEFAULT 0 CHECK (balance >= 0),
    ratio FLOAT8,
    small INT2,
    secret STRING NOT VISIBLE,
    doc JSONB,
    created TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT customers_email_check CHECK (email != '')
);
COMMENT ON TABLE fid_cols.customers IS 'a customer';
COMMENT ON COLUMN fid_cols.customers.email IS 'login';
COMMENT ON CONSTRAINT customers_email_check ON fid_cols.customers IS 'not empty';
CREATE INDEX customers_name_idx ON fid_cols.customers (name DESC) STORING (balance);
COMMENT ON INDEX fid_cols.customers@customers_name_idx IS 'names';
CREATE INDEX customers_small_idx ON fid_cols.customers (small) WHERE small > 0;
CREATE INDEX customers_lower_email_idx ON fid_cols.customers (lower(email));
CREATE INVERTED INDEX customers_doc_idx ON fid_cols.customers (doc);
CREATE INDEX customers_ratio_idx ON fid_cols.customers (ratio) USING HASH WITH (bucket_count = 8);
CREATE INDEX customers_created_idx ON fid_cols.customers (created) NOT VISIBLE;
CREATE TABLE fid_cols.order_lines (
    order_id INT8 NOT NULL,
    line INT4 NOT NULL,
    customer_id INT8 NOT NULL REFERENCES fid_cols.customers (id) ON DELETE CASCADE ON UPDATE RESTRICT,
    amount DECIMAL(10, 2) NOT NULL,
    CONSTRAINT order_lines_pkey PRIMARY KEY (order_id ASC, line DESC),
    CONSTRAINT order_lines_amount_positive CHECK (amount > 0)
);
CREATE TABLE fid_cols.events (at TIMESTAMPTZ, note STRING);
CREATE TABLE fid_cols.custom_key (rowid_custom INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(), v INT8, PRIMARY KEY (rowid_custom));
CREATE TABLE fid_cols.sharded (id INT8 NOT NULL, small INT4, PRIMARY KEY (id) USING HASH WITH (bucket_count = 4));
`,
		},
		{
			name: "views-routines-sequences",
			sql: `CREATE SCHEMA fid_routines;
CREATE SEQUENCE fid_routines.counter AS INT8 INCREMENT BY 5 MINVALUE -100 MAXVALUE 1000 START WITH 10 CACHE 3;
CREATE SEQUENCE fid_routines.plain;
CREATE TABLE fid_routines.tickets (id INT8 PRIMARY KEY, n INT8);
CREATE SEQUENCE fid_routines.ticket_seq OWNED BY fid_routines.tickets.n;
CREATE TABLE fid_routines.items (
    id INT8 PRIMARY KEY DEFAULT nextval('fid_routines.plain'),
    name STRING NOT NULL,
    price DECIMAL(10, 2),
    cents INT8 AS ((price * 100)::INT8) STORED,
    doubled DECIMAL AS (price * 2) VIRTUAL
);
CREATE VIEW fid_routines.v_items AS SELECT id, name FROM fid_routines.items WHERE price > 0;
CREATE VIEW fid_routines.v_item_names AS SELECT name FROM fid_routines.v_items;
CREATE VIEW fid_routines.v_renamed (item_id, item_name) AS SELECT id, name FROM fid_routines.items;
CREATE MATERIALIZED VIEW fid_routines.item_totals AS SELECT name, sum(price) AS total FROM fid_routines.items GROUP BY name;
CREATE INDEX item_totals_name_idx ON fid_routines.item_totals (name);
CREATE FUNCTION fid_routines.add(a INT8, b INT8) RETURNS INT8 LANGUAGE SQL IMMUTABLE AS $$ SELECT a + b $$;
CREATE FUNCTION fid_routines.add(a STRING, b STRING) RETURNS STRING LANGUAGE SQL IMMUTABLE AS $$ SELECT a || b $$;
CREATE FUNCTION fid_routines.item_count() RETURNS INT8 LANGUAGE SQL STABLE AS $$ SELECT count(*) FROM fid_routines.v_items $$;
CREATE FUNCTION fid_routines.touch() RETURNS INT8 LANGUAGE PLPGSQL AS $$
BEGIN
    RETURN (SELECT count(*) FROM fid_routines.items);
END
$$;
CREATE PROCEDURE fid_routines.reset_prices() LANGUAGE SQL AS $$ UPDATE fid_routines.items SET price = 0 $$;
CREATE TABLE fid_routines.scored (id INT8 PRIMARY KEY, score INT8 DEFAULT fid_routines.add(1, 2));
CREATE FUNCTION fid_routines.a_constant() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
CREATE TABLE fid_routines.t1 (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_routines.a_constant());
CREATE FUNCTION fid_routines.z_read() RETURNS INT8 LANGUAGE SQL AS $$ SELECT count(*) FROM fid_routines.t1 $$;
CREATE TABLE fid_routines.t2 (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_routines.z_read());
CREATE FUNCTION fid_routines."Quoted"() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 2 $$;
CREATE TABLE fid_routines.t3 (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_routines."Quoted"());
CREATE TABLE fid_routines.z_rel (id INT8 PRIMARY KEY);
CREATE VIEW fid_routines.a_count AS SELECT count(*) AS n FROM fid_routines.z_rel;
CREATE VIEW fid_routines.z_cte AS WITH a_count AS (SELECT 1 AS id) SELECT id FROM a_count;
CREATE VIEW fid_routines.a_reads_z AS SELECT id FROM fid_routines.z_cte;
CREATE MATERIALIZED VIEW fid_routines.a_mcount AS SELECT count(*) AS n FROM fid_routines.z_rel;
CREATE MATERIALIZED VIEW fid_routines.m_renamed (item_count) AS SELECT count(*) FROM fid_routines.z_rel;
`,
		},
		{
			name: "cross-schema",
			sql: `CREATE SCHEMA fid_a;
CREATE SCHEMA "Fid Mixed";
CREATE TABLE fid_a.parents (id INT8 PRIMARY KEY);
CREATE TYPE "Fid Mixed"."Kind" AS ENUM ('first', 'second');
CREATE TABLE "Fid Mixed"."Children" (
    id INT8 PRIMARY KEY,
    "Parent Id" INT8 REFERENCES fid_a.parents (id),
    kind "Fid Mixed"."Kind" DEFAULT 'first'
);
CREATE FUNCTION fid_a.f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
CREATE FUNCTION fid_a.g() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_a.f() $$;
CREATE FUNCTION "Fid Mixed".f() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_a.g() $$;
CREATE FUNCTION fid_a.f(n INT8) RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_a.g() + n $$;
CREATE TABLE fid_a.uses_g (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_a.g());
CREATE FUNCTION fid_a.h(n INT8) RETURNS INT8 LANGUAGE SQL AS $$ SELECT n $$;
CREATE FUNCTION fid_a.k() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_a.h(1:::INT8) $$;
CREATE FUNCTION fid_a.h(s STRING) RETURNS STRING LANGUAGE SQL AS $$ SELECT fid_a.k()::STRING || s $$;
CREATE TABLE fid_a.uses_k (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_a.k());
CREATE FUNCTION fid_a.p() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
CREATE TABLE fid_a.t_p (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_a.p());
CREATE FUNCTION fid_a.p(n INT8) RETURNS INT8 LANGUAGE SQL AS $$ SELECT count(*) + n FROM fid_a.t_p $$;
CREATE FUNCTION fid_a.z(n INT8) RETURNS INT8 LANGUAGE SQL AS $$ SELECT n $$;
CREATE FUNCTION fid_a.a_calls_z() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_a.z(1:::INT8) $$;
CREATE FUNCTION fid_a.z(s STRING) RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_a.a_calls_z() $$;
CREATE TABLE fid_a.uses_a (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_a.a_calls_z());
CREATE FUNCTION fid_a.arr(a INT8[]) RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
CREATE FUNCTION fid_a.b_calls_arr() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_a.arr(ARRAY[1, 2]) $$;
CREATE FUNCTION fid_a.arr(a INT8, b INT8) RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_a.b_calls_arr() $$;
CREATE TABLE fid_a.uses_b (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_a.b_calls_arr());
CREATE VIEW "Fid Mixed".family AS SELECT p.id, c.id AS child FROM fid_a.parents AS p JOIN "Fid Mixed"."Children" AS c ON c."Parent Id" = p.id;
`,
		},
	}
}

func TestDefinitionFidelity(t *testing.T) {
	db := testcontainer.Open(t)
	ctx := context.Background()
	for _, tc := range fidelityCases() {
		t.Run(tc.name, func(t *testing.T) {
			execScript(t, db, tc.sql)
			stored, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
			require.NoError(t, err)
			stored = filterTo(stored, caseSchemas(stored))
			def, err := catalog.DatabaseDefinition(stored, catalog.DefinitionOptions{PrintHeader: true})
			require.NoError(t, err)
			// The case's own DDL, read as SDL, holds what its sync holds.
			sdl, err := catalog.DatabaseDefinition(stored, catalog.DefinitionOptions{SDLFormat: true})
			require.NoError(t, err)
			plan, err := catalog.BuildSDLMigrationPlan(sdl, tc.sql, catalog.SDLMigrationOptions{})
			require.NoError(t, err)
			require.Nil(t, plan, "SDL:\n%s\ndiffers from the case by:\n%s", sdl, plan.SQL())

			execScript(t, db, dropScript(stored))
			execScript(t, db, def)
			restored, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
			require.NoError(t, err)
			restoredDef, err := catalog.DatabaseDefinition(filterTo(restored, stored), catalog.DefinitionOptions{PrintHeader: true})
			require.NoError(t, err)
			require.Equal(t, def, restoredDef)
			execScript(t, db, dropScript(stored))
		})
	}
}

// caseSchemas is the snapshot without the schemas every database holds.
func caseSchemas(snapshot *metadata.DatabaseSchemaMetadata) *metadata.DatabaseSchemaMetadata {
	wanted := &metadata.DatabaseSchemaMetadata{}
	for _, schema := range snapshot.Schemas {
		if strings.HasPrefix(strings.ToLower(schema.Name), "fid") {
			wanted.Schemas = append(wanted.Schemas, schema)
		}
	}
	return wanted
}

// filterTo keeps the schemas of synced that wanted names.
func filterTo(synced, wanted *metadata.DatabaseSchemaMetadata) *metadata.DatabaseSchemaMetadata {
	names := make(map[string]bool)
	for _, schema := range wanted.Schemas {
		names[schema.Name] = true
	}
	filtered := &metadata.DatabaseSchemaMetadata{Name: synced.Name, SearchPath: synced.SearchPath}
	for _, schema := range synced.Schemas {
		if names[schema.Name] {
			filtered.Schemas = append(filtered.Schemas, schema)
		}
	}
	return filtered
}

// dropScript drops the schemas a snapshot holds and everything in them.
func dropScript(snapshot *metadata.DatabaseSchemaMetadata) string {
	var names []string
	for _, schema := range snapshot.Schemas {
		names = append(names, schema.Name)
	}
	slices.Sort(names)
	var lines []string
	for _, name := range names {
		lines = append(lines, fmt.Sprintf(`DROP SCHEMA IF EXISTS "%s" CASCADE;`, strings.ReplaceAll(name, `"`, `""`)))
	}
	return strings.Join(lines, "\n")
}

// execScript sends a script one statement at a time, as CockroachDB's parser
// splits it.
func execScript(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	for _, stmt := range analysis.Split(script) {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := db.ExecContext(context.Background(), stmt); err != nil {
			t.Fatalf("%s\n%v", stmt, err)
		}
	}
}
