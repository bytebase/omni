package catalog_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cockroachdb/catalog"
	"github.com/bytebase/omni/cockroachdb/driver"
	"github.com/bytebase/omni/cockroachdb/internal/testcontainer"
	"github.com/bytebase/omni/metadata"
)

// migrationCase is a source and a target schema, each written as CockroachDB
// DDL. The test syncs both from the server, generates the migration between
// the snapshots, runs it on the source, and requires the result to sync to
// the target. A column a target adds comes last, where CockroachDB adds it.
type migrationCase struct {
	name   string
	source string
	target string
}

func migrationCases() []migrationCase {
	return []migrationCase{
		{
			name: "tables",
			source: `CREATE SCHEMA fid_m;
CREATE TYPE fid_m.status AS ENUM ('new', 'done');
CREATE SEQUENCE fid_m.ids START WITH 10;
CREATE TABLE fid_m.customers (
    id INT8 PRIMARY KEY,
    email STRING NOT NULL,
    nick STRING,
    age INT4 DEFAULT 0,
    legacy STRING,
    CONSTRAINT customers_age_check CHECK (age >= 0)
);
COMMENT ON TABLE fid_m.customers IS 'old';
COMMENT ON COLUMN fid_m.customers.legacy IS 'going';
CREATE INDEX customers_nick_idx ON fid_m.customers (nick);
CREATE INDEX customers_email_idx ON fid_m.customers (email);
CREATE TABLE fid_m.orders (
    id INT8 NOT NULL,
    customer_id INT8 REFERENCES fid_m.customers (id),
    status fid_m.status DEFAULT 'new',
    CONSTRAINT orders_pkey PRIMARY KEY (id)
);
CREATE TABLE fid_m.gone (id INT8 PRIMARY KEY, customer_id INT8 REFERENCES fid_m.customers (id));
`,
			target: `CREATE SCHEMA fid_m;
CREATE SCHEMA fid_new;
CREATE TYPE fid_m.status AS ENUM ('new', 'paid', 'done', 'void');
CREATE SEQUENCE fid_m.ids START WITH 10 INCREMENT BY 2 CACHE 4;
CREATE TABLE fid_m.customers (
    id INT8 PRIMARY KEY,
    email STRING,
    nick STRING NOT VISIBLE,
    age INT8 DEFAULT 18 NOT NULL,
    joined TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT customers_age_check CHECK (age >= 18),
    CONSTRAINT customers_email_check CHECK (email != '')
);
COMMENT ON TABLE fid_m.customers IS 'new';
COMMENT ON COLUMN fid_m.customers.email IS 'login';
CREATE INDEX customers_nick_idx ON fid_m.customers (nick DESC);
CREATE UNIQUE INDEX customers_joined_key ON fid_m.customers (joined) STORING (email);
COMMENT ON INDEX fid_m.customers@customers_joined_key IS 'joins';
CREATE TABLE fid_m.orders (
    id INT8 NOT NULL,
    customer_id INT8,
    status fid_m.status DEFAULT 'paid',
    line INT8 NOT NULL DEFAULT 1,
    CONSTRAINT orders_pkey PRIMARY KEY (id, line),
    CONSTRAINT orders_customer_fk FOREIGN KEY (customer_id) REFERENCES fid_m.customers (id) ON DELETE CASCADE
);
CREATE TABLE fid_new.audit (at TIMESTAMPTZ, customer_id INT8 REFERENCES fid_m.customers (id), id INT8 DEFAULT nextval('fid_m.ids'));
`,
		},
		{
			name: "views-routines",
			source: `CREATE SCHEMA fid_v;
CREATE TABLE fid_v.items (id INT8 PRIMARY KEY, name STRING, price DECIMAL(10, 2), legacy INT8);
CREATE VIEW fid_v.priced AS SELECT id, name, legacy FROM fid_v.items WHERE price > 0;
CREATE VIEW fid_v.priced_names AS SELECT name FROM fid_v.priced;
CREATE VIEW fid_v.unchanged AS SELECT id FROM fid_v.items;
CREATE VIEW fid_v.dropped AS SELECT name FROM fid_v.items;
CREATE MATERIALIZED VIEW fid_v.totals AS SELECT name, sum(price) AS total FROM fid_v.items GROUP BY name;
CREATE INDEX totals_name_idx ON fid_v.totals (name);
CREATE FUNCTION fid_v.add(a INT8, b INT8) RETURNS INT8 LANGUAGE SQL IMMUTABLE AS $$ SELECT a + b $$;
CREATE FUNCTION fid_v.gone() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
CREATE FUNCTION fid_v.count_items() RETURNS INT8 LANGUAGE SQL STABLE AS $$ SELECT count(*) FROM fid_v.items $$;
`,
			target: `CREATE SCHEMA fid_v;
CREATE TABLE fid_v.items (id INT8 PRIMARY KEY, name STRING, price DECIMAL(10, 2));
CREATE VIEW fid_v.priced AS SELECT id, name FROM fid_v.items WHERE price > 1;
CREATE VIEW fid_v.priced_names AS SELECT name FROM fid_v.priced;
CREATE VIEW fid_v.unchanged AS SELECT id FROM fid_v.items;
CREATE MATERIALIZED VIEW fid_v.totals AS SELECT name, sum(price) AS total FROM fid_v.items GROUP BY name;
CREATE INDEX totals_total_idx ON fid_v.totals (total);
CREATE FUNCTION fid_v.add(a INT8, b INT8) RETURNS INT8 LANGUAGE SQL IMMUTABLE AS $$ SELECT a + b + 0 $$;
CREATE FUNCTION fid_v.count_items() RETURNS INT8 LANGUAGE SQL STABLE AS $$ SELECT count(*) FROM fid_v.priced $$;
CREATE FUNCTION fid_v.touch() RETURNS INT8 LANGUAGE PLPGSQL AS $$ BEGIN RETURN (SELECT count(*) FROM fid_v.items); END $$;
CREATE TABLE fid_v.scored (id INT8 PRIMARY KEY, score INT8 DEFAULT fid_v.add(1, 2));
CREATE VIEW fid_v.cheap AS SELECT p.id FROM fid_v.priced AS p JOIN fid_v.items AS i ON i.id = p.id;
`,
		},
	}
}

func TestMigrationFidelity(t *testing.T) {
	db := testcontainer.Open(t)
	ctx := context.Background()
	sync := func(t *testing.T) *metadata.DatabaseSchemaMetadata {
		t.Helper()
		synced, err := driver.SyncDatabase(ctx, db, testcontainer.Database(t))
		require.NoError(t, err)
		return filterTo(synced, caseSchemas(synced))
	}
	definition := func(t *testing.T, snapshot *metadata.DatabaseSchemaMetadata) string {
		t.Helper()
		def, err := catalog.DatabaseDefinition(snapshot, catalog.DefinitionOptions{})
		require.NoError(t, err)
		return def
	}
	for _, tc := range migrationCases() {
		t.Run(tc.name, func(t *testing.T) {
			execScript(t, db, tc.target)
			target := sync(t)
			execScript(t, db, dropScript(target))
			execScript(t, db, tc.source)
			source := sync(t)

			same, err := catalog.MetadataMigration(source, source, catalog.MetadataMigrationOptions{})
			require.NoError(t, err)
			require.Empty(t, same)

			migration, err := catalog.MetadataMigration(source, target, catalog.MetadataMigrationOptions{})
			require.NoError(t, err)
			t.Logf("migration:\n%s", migration)
			execScript(t, db, migration+";")
			migrated := sync(t)
			require.Equal(t, definition(t, target), definition(t, migrated), "migration:\n%s", migration)
			rest, err := catalog.MetadataMigration(migrated, target, catalog.MetadataMigrationOptions{})
			require.NoError(t, err)
			require.Empty(t, rest, "migration:\n%s", migration)
			execScript(t, db, dropScript(migrated))
			execScript(t, db, dropScript(target))
		})
	}
}
