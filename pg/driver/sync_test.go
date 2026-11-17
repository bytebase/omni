package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
	pganalysis "github.com/bytebase/omni/pg/analysis"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// TestSyncQueriesOrderTheirRows holds every query sync.go sends to an
// ORDER BY. A sync query's row order decides the recorded OIDs (the
// testcontainer numbers them as the test reads them) and the order of the
// metadata's slices, so an order the container picks records a different
// tape and golden on the next run, and re-recording cannot prove a change.
func TestSyncQueriesOrderTheirRows(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sync.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	var exprs []ast.Expr
	ast.Inspect(file, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ValueSpec:
			exprs = append(exprs, n.Values...)
		case *ast.AssignStmt:
			exprs = append(exprs, n.Rhs...)
		case *ast.ReturnStmt:
			exprs = append(exprs, n.Results...)
		}
		return true
	})
	for _, e := range exprs {
		text, ok := queryText(e)
		if !ok || !readQuery.MatchString(text) || !readsRows.MatchString(text) || ordersRows.MatchString(text) {
			continue
		}
		t.Errorf("the query at %s returns rows the container orders itself, which a recording cannot reproduce. Add an ORDER BY that makes the order a function of the data:\n%s",
			fset.Position(e.Pos()), strings.TrimSpace(text))
	}
}

var (
	readQuery  = regexp.MustCompile(`(?is)^\s*(select|with|table|values)\b`)
	readsRows  = regexp.MustCompile(`(?i)\bfrom\b`)
	ordersRows = regexp.MustCompile(`(?i)\border\s+by\b`)
)

// queryText is the text of a query an expression builds: its string literals,
// with an fmt.Sprintf call contributing its format string, so a query written
// as concatenation is read whole. It reports false for anything else.
func queryText(e ast.Expr) (string, bool) {
	switch v := e.(type) {
	case *ast.BasicLit:
		if v.Kind != token.STRING {
			return "", false
		}
		s, err := strconv.Unquote(v.Value)
		if err != nil {
			return "", false
		}
		return s, true
	case *ast.ParenExpr:
		return queryText(v.X)
	case *ast.BinaryExpr:
		if v.Op != token.ADD {
			return "", false
		}
		left, leftOK := queryText(v.X)
		right, rightOK := queryText(v.Y)
		if !leftOK || !rightOK {
			return "", false
		}
		return left + right, true
	case *ast.CallExpr:
		if len(v.Args) == 0 {
			return "", false
		}
		return queryText(v.Args[0])
	}
	return "", false
}

// versions are the pinned servers of every major release Bytebase
// supports, oldest first.
var versions = []string{"pg12", "pg13", "pg14", "pg15", "pg16", testcontainer.Default, testcontainer.PG18}

// openDriver returns a driver on t's database on the server of role: a
// replay of t's tape, or the pinned server in a recording.
func openDriver(t *testing.T, role string) *Driver {
	t.Helper()
	db := testcontainer.OpenRoleNumberingOIDs(t, role)
	t.Cleanup(func() { db.Close() })
	return &Driver{db: db, databaseName: testcontainer.Database(t)}
}

func mustExec(t *testing.T, db *sql.DB, script string) {
	t.Helper()
	if _, err := db.ExecContext(context.Background(), script); err != nil {
		t.Fatalf("%v", err)
	}
}

// TestSyncDatabase syncs testdata/corpus.sql on each pinned server and holds
// the result to testdata/sync/<role>.json.
func TestSyncDatabase(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	for _, role := range versions {
		t.Run(role, func(t *testing.T) {
			d := openDriver(t, role)
			mustExec(t, d.db, string(corpus))
			got, err := d.SyncDatabase(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			golden := filepath.Join("testdata", "sync", role+".json")
			if *update {
				b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(golden, append(b, '\n'), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want := readGolden(t, golden)
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("SyncDatabase on %s differs from %s (-want +got):\n%s", role, golden, diff)
			}
		})
	}
}

func readGolden(t *testing.T, path string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(b, m); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return m
}

// TestSyncInstance reads each pinned server's version, databases and roles,
// and holds them to testdata/sync/instance.json.
func TestSyncInstance(t *testing.T) {
	got := map[string]json.RawMessage{}
	for _, role := range versions {
		t.Run(role, func(t *testing.T) {
			d := openDriver(t, role)
			mustExec(t, d.db, `
CREATE ROLE sync_nologin NOLOGIN;
CREATE ROLE sync_admin SUPERUSER NOINHERIT CREATEROLE CREATEDB REPLICATION BYPASSRLS LOGIN VALID UNTIL '2030-06-01 12:00:00+00';
CREATE ROLE sync_plain LOGIN;`)
			instance, err := d.SyncInstance(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			databases := make([]string, 0, len(instance.Databases))
			for _, db := range instance.Databases {
				databases = append(databases, fmt.Sprintf("%s %s %s %s", db.Name, db.CharacterSet, db.Collation, db.Owner))
			}
			roles := make([]string, 0, len(instance.Roles))
			for _, r := range instance.Roles {
				roles = append(roles, r.Name+": "+*r.Attribute)
			}
			b, err := json.MarshalIndent(map[string]any{"version": instance.Version, "databases": databases, "roles": roles}, "  ", "  ")
			if err != nil {
				t.Fatal(err)
			}
			got[role] = b
		})
	}
	golden := filepath.Join("testdata", "sync", "instance.json")
	b, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	b = append(b, '\n')
	if *update {
		if err := os.WriteFile(golden, b, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), string(b)); diff != "" {
		t.Errorf("SyncInstance differs from %s (-want +got):\n%s", golden, diff)
	}
}

func TestSyncDatabaseNotFound(t *testing.T) {
	d := openDriver(t, testcontainer.Default)
	d.databaseName = "no_such_database"
	if _, err := d.SyncDatabase(context.Background()); !errors.Is(err, driver.ErrNotFound) {
		t.Fatalf("SyncDatabase of a missing database: %v, want driver.ErrNotFound", err)
	}
}

func TestIsAtLeastMajorVersion(t *testing.T) {
	for _, tc := range []struct {
		version string
		want    bool
	}{
		{"hello", true},
		{"", true},
		{"9.6.24", false},
		{"9.624", false},
		{"10.0.0", true},
		{"16.0.0", true},
		{"16.14", true},
		{"v12", true},
	} {
		if got := isAtLeastMajorVersion(tc.version, 10); got != tc.want {
			t.Errorf("isAtLeastMajorVersion(%q, 10) = %v, want %v", tc.version, got, tc.want)
		}
	}
}

// findTable returns schema.table from m, failing t when it is missing.
func findTable(t *testing.T, m *metadatapb.DatabaseSchemaMetadata, schema, table string) *metadatapb.TableMetadata {
	t.Helper()
	for _, s := range m.Schemas {
		if s.Name != schema {
			continue
		}
		for _, tb := range s.Tables {
			if tb.Name == table {
				return tb
			}
		}
	}
	t.Fatalf("%s.%s not found in the synced metadata", schema, table)
	return nil
}

func columnDefaults(table *metadatapb.TableMetadata) map[string]string {
	m := map[string]string{}
	for _, c := range table.Columns {
		m[c.Name] = c.Default
	}
	return m
}

// TestSyncColumnDefaults: a default names its types, sequences and functions
// with their schema, since the sync clears the search path first.
func TestSyncColumnDefaults(t *testing.T) {
	d := openDriver(t, testcontainer.Default)
	mustExec(t, d.db, `
CREATE TYPE test_size AS ENUM ('small', 'medium', 'large');
CREATE SEQUENCE test_sequence START 1;
CREATE SCHEMA custom_schema;
CREATE TYPE custom_schema.status_type AS ENUM ('active', 'inactive');
CREATE SEQUENCE custom_schema.my_sequence START 100;
CREATE FUNCTION custom_schema.get_prefix() RETURNS text AS $$ BEGIN RETURN 'p'; END $$ LANGUAGE plpgsql IMMUTABLE;
CREATE TABLE test_defaults (
    id integer,
    col_null_default text DEFAULT NULL,
    col_string_literal varchar(50) DEFAULT 'hello',
    col_numeric integer DEFAULT 42,
    col_boolean boolean DEFAULT true,
    col_function timestamp DEFAULT now(),
    col_expression integer DEFAULT (10 + 20),
    col_serial serial,
    col_enum test_size DEFAULT 'medium',
    col_sequence integer DEFAULT nextval('test_sequence'),
    col_array integer[] DEFAULT '{1,2,3}',
    col_json jsonb DEFAULT '{"key": "value"}',
    col_binary bytea DEFAULT '\xDEADBEEF',
    col_cross_enum custom_schema.status_type DEFAULT 'active',
    col_cross_sequence integer DEFAULT nextval('custom_schema.my_sequence'),
    col_cross_function text DEFAULT custom_schema.get_prefix()
);`)
	m, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	got := columnDefaults(findTable(t, m, "public", "test_defaults"))
	want := map[string]string{
		"id":                 "",
		"col_null_default":   "",
		"col_string_literal": "'hello'::character varying",
		"col_numeric":        "42",
		"col_boolean":        "true",
		"col_function":       "now()",
		"col_expression":     "(10 + 20)",
		"col_serial":         "nextval('public.test_defaults_col_serial_seq'::regclass)",
		"col_enum":           "'medium'::public.test_size",
		"col_sequence":       "nextval('public.test_sequence'::regclass)",
		"col_array":          "'{1,2,3}'::integer[]",
		"col_json":           `'{"key": "value"}'::jsonb`,
		"col_binary":         `'\xdeadbeef'::bytea`,
		"col_cross_enum":     "'active'::custom_schema.status_type",
		"col_cross_sequence": "nextval('custom_schema.my_sequence'::regclass)",
		"col_cross_function": "custom_schema.get_prefix()",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("column defaults (-want +got):\n%s", diff)
	}
}

// TestSyncSequenceOwner: a sequence's owner column is the one its attnum
// names, even after a column before it was dropped.
func TestSyncSequenceOwner(t *testing.T) {
	d := openDriver(t, testcontainer.Default)
	mustExec(t, d.db, `
CREATE TABLE seq_owner (removed text, id serial, payload text);
ALTER TABLE seq_owner DROP COLUMN removed;
CREATE TABLE identity_owner (gone int, other int, id int GENERATED ALWAYS AS IDENTITY);
ALTER TABLE identity_owner DROP COLUMN gone;`)
	m, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, s := range m.Schemas {
		for _, seq := range s.Sequences {
			got[seq.Name] = seq.OwnerTable + "." + seq.OwnerColumn
		}
	}
	want := map[string]string{
		"seq_owner_id_seq":      "seq_owner.id",
		"identity_owner_id_seq": "identity_owner.id",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("sequence owners (-want +got):\n%s", diff)
	}
}

// TestSyncConstraintComments reads COMMENT ON CONSTRAINT into a check
// constraint's and a foreign key's comment, which PostgreSQL keeps in
// pg_description beside the constraint rather than in it.
func TestSyncConstraintComments(t *testing.T) {
	d := openDriver(t, testcontainer.Default)
	mustExec(t, d.db, `
CREATE TABLE parent (
    id integer PRIMARY KEY,
    CONSTRAINT parent_id_check CHECK (id > 0)
);
CREATE TABLE child (
    id integer CONSTRAINT child_id_check CHECK (id > 0),
    parent_id integer,
    CONSTRAINT child_uncommented CHECK (id <> 42),
    CONSTRAINT child_parent_fk FOREIGN KEY (parent_id) REFERENCES parent (id)
);
COMMENT ON CONSTRAINT child_id_check ON child IS 'a child has a positive id';
COMMENT ON CONSTRAINT child_parent_fk ON child IS 'the parent of a child';
COMMENT ON CONSTRAINT parent_id_check ON parent IS 'a parent''s id is positive
and the note keeps its line break';`)
	m, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	comments := func(table *metadatapb.TableMetadata) []string {
		var out []string
		for _, check := range table.CheckConstraints {
			out = append(out, "check "+check.Name+": "+check.Comment)
		}
		for _, fk := range table.ForeignKeys {
			out = append(out, "fk "+fk.Name+": "+fk.Comment)
		}
		return out
	}
	got := map[string][]string{
		"parent": comments(findTable(t, m, "public", "parent")),
		"child":  comments(findTable(t, m, "public", "child")),
	}
	want := map[string][]string{
		"parent": {"check parent_id_check: a parent's id is positive\nand the note keeps its line break"},
		"child":  {"check child_id_check: a child has a positive id", "check child_uncommented: ", "fk child_parent_fk: the parent of a child"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("constraint comments (-want +got):\n%s", diff)
	}
}

// TestSyncCompositeTypes pins composite types exactly: the definition
// round trip compares two syncs with each other, which a sync wrong the same
// way on both sides passes.
func TestSyncCompositeTypes(t *testing.T) {
	d := openDriver(t, testcontainer.Default)
	mustExec(t, d.db, `
CREATE TYPE zz_base AS (street text COLLATE "C", city varchar(50));
COMMENT ON TYPE zz_base IS 'base address type';
COMMENT ON COLUMN zz_base.street IS 'street line';
CREATE TYPE aa_nested AS (home zz_base, homes zz_base[], tags text[]);
CREATE TYPE empty_type AS ();
CREATE SCHEMA geo;
CREATE TYPE geo.point2 AS (lat numeric(9,6), lng numeric(9,6));
CREATE TYPE wrapper AS (p geo.point2);
CREATE TABLE plain_table (id int);
CREATE SCHEMA locale;
CREATE COLLATION locale.mycoll (locale = 'C');
CREATE TYPE coll_qualified AS (v text COLLATE locale.mycoll);`)
	m, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	got := map[string][]string{}
	for _, s := range m.Schemas {
		for _, c := range s.CompositeTypes {
			attrs := []string{"comment=" + c.Comment}
			for _, a := range c.Attributes {
				attrs = append(attrs, fmt.Sprintf("%s %s collate=%s comment=%s", a.Name, a.Type, a.Collation, a.Comment))
			}
			got[s.Name+"."+c.Name] = attrs
		}
	}
	want := map[string][]string{
		"public.aa_nested":      {"comment=", "home public.zz_base collate= comment=", "homes public.zz_base[] collate= comment=", "tags text[] collate= comment="},
		"public.coll_qualified": {"comment=", "v text collate=locale.mycoll comment="},
		"public.empty_type":     {"comment="},
		"public.wrapper":        {"comment=", "p geo.point2 collate= comment="},
		"public.zz_base":        {"comment=base address type", `street text collate="C" comment=street line`, "city character varying(50) collate= comment="},
		"geo.point2":            {"comment=", "lat numeric(9,6) collate= comment=", "lng numeric(9,6) collate= comment="},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("composite types (-want +got):\n%s", diff)
	}
}

// legacyListColumnQuery is the information_schema.columns query the
// pg_catalog one replaced; for a user with every privilege the two must
// return the same rows.
var legacyListColumnQuery = `
SELECT
	cols.table_schema,
	cols.table_name,
	cols.column_name,
	cols.data_type,
	cols.character_maximum_length,
	cols.numeric_precision,
	cols.numeric_scale,
	cols.datetime_precision,
	cols.ordinal_position,
	cols.column_default,
	pa.attgenerated AS generation_kind,
	cols.generation_expression,
	cols.is_nullable,
	cols.collation_name,
	cols.udt_schema,
	cols.udt_name,
	cols.identity_generation,
	pg_catalog.col_description(c.oid, cols.ordinal_position::int) as column_comment
FROM INFORMATION_SCHEMA.COLUMNS AS cols
LEFT JOIN pg_namespace n ON n.nspname = cols.table_schema
LEFT JOIN pg_class c ON c.relname = cols.table_name AND c.relnamespace = n.oid AND c.relkind IN ('r', 'v', 'm', 'p', 'f')
LEFT JOIN pg_attribute pa ON pa.attrelid = c.oid AND pa.attname = cols.column_name AND NOT pa.attisdropped` + fmt.Sprintf(`
WHERE cols.table_schema NOT IN (%s)
  AND cols.table_schema NOT LIKE 'pg_temp%%'
  AND cols.table_schema NOT LIKE 'pg_toast%%'
ORDER BY cols.table_schema, cols.table_name, cols.ordinal_position;`, pganalysis.SystemSchemaWhereClause)

// queryRows returns a query's rows as text, sorted by its first columns.
func queryRows(t *testing.T, q interface {
	Query(string, ...any) (*sql.Rows, error)
}, query string, sortBy int) [][]string {
	t.Helper()
	rows, err := q.Query(query)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}
	var out [][]string
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		ptrs := make([]any, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			t.Fatal(err)
		}
		row := make([]string, len(columns))
		for i, v := range values {
			row[i] = v.String
			if !v.Valid {
				row[i] = "<null>"
			}
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	// The queries order exotic names differently: information_schema sorts
	// varchar in the database collation, pg_catalog sorts name in C.
	slices.SortFunc(out, func(a, b []string) int {
		return strings.Compare(strings.Join(a[:sortBy], "\x00"), strings.Join(b[:sortBy], "\x00"))
	})
	return out
}

func TestListColumnQueryMatchesInformationSchema(t *testing.T) {
	d := openDriver(t, testcontainer.Default)
	mustExec(t, d.db, `
CREATE SCHEMA app;
CREATE TYPE app.mood AS ENUM ('happy', 'sad');
CREATE TYPE public.shirt_size AS ENUM ('s', 'm', 'l');
CREATE DOMAIN app.positive_int AS integer NOT NULL CHECK (VALUE > 0);
CREATE DOMAIN public.us_zip AS varchar(10);
CREATE DOMAIN app.int_array AS integer[];
CREATE SEQUENCE app.seq1;
CREATE TABLE app.kitchen_sink (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    id2 integer GENERATED BY DEFAULT AS IDENTITY,
    c_smallint smallint,
    c_int integer NOT NULL DEFAULT 42,
    c_bigint bigint DEFAULT nextval('app.seq1'),
    c_numeric numeric,
    c_numeric_p numeric(8),
    c_numeric_p0 numeric(8, 0),
    c_numeric_ps numeric(10, 2),
    c_real real,
    c_double double precision,
    c_money money,
    c_char char,
    c_char_n char(5),
    c_varchar varchar,
    c_varchar_n varchar(20) DEFAULT 'hello',
    c_text text COLLATE "C",
    c_bytea bytea,
    c_bool boolean DEFAULT true,
    c_date date,
    c_time time,
    c_time_p time(3),
    c_timetz time with time zone,
    c_timetz_p time(5) with time zone,
    c_ts timestamp,
    c_ts_p timestamp(2),
    c_tstz timestamp with time zone DEFAULT now(),
    c_tstz_p timestamp(4) with time zone,
    c_interval interval,
    c_interval_p interval(3),
    c_uuid uuid,
    c_json json,
    c_jsonb jsonb DEFAULT '{}'::jsonb,
    c_xml xml,
    c_inet inet,
    c_cidr cidr,
    c_macaddr macaddr,
    c_bit bit,
    c_bit_n bit(3),
    c_varbit varbit,
    c_varbit_n varbit(8),
    c_int_array integer[],
    c_text_matrix text[][],
    c_varchar_array varchar(5)[],
    c_enum app.mood DEFAULT 'happy',
    c_enum_public public.shirt_size,
    c_enum_array app.mood[],
    c_domain app.positive_int,
    c_domain_varchar public.us_zip DEFAULT 'x',
    c_domain_array app.int_array,
    c_point point,
    c_tsvector tsvector,
    c_int4range int4range,
    c_oid oid,
    c_name name,
    c_internal_char "char",
    c_generated integer GENERATED ALWAYS AS (c_int * 2) STORED,
    c_drop_me text
);
ALTER TABLE app.kitchen_sink DROP COLUMN c_drop_me;
COMMENT ON COLUMN app.kitchen_sink.c_int IS 'int comment';
CREATE TABLE app.parted (id integer, region text NOT NULL) PARTITION BY LIST (region);
CREATE TABLE app.parted_us PARTITION OF app.parted FOR VALUES IN ('us');
CREATE VIEW app.v_sink AS SELECT id, c_int, c_text FROM app.kitchen_sink;
CREATE TABLE public.plain (a serial, b text);`)
	txn, err := d.db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer txn.Rollback()
	// The sync clears the search path before reading columns, so that
	// defaults name their schemas.
	if err := setTxSearchPath(txn, ""); err != nil {
		t.Fatal(err)
	}
	legacy := queryRows(t, txn, legacyListColumnQuery, 2)
	current := queryRows(t, txn, getListColumnQuery(true), 2)
	if len(legacy) < 50 {
		t.Fatalf("information_schema returned %d columns, want the kitchen sink's 50 and more", len(legacy))
	}
	if diff := cmp.Diff(legacy, current); diff != "" {
		t.Errorf("column query rows (-information_schema +pg_catalog):\n%s", diff)
	}
}

// legacyListForeignTableQuery is the information_schema.foreign_tables query
// the pg_catalog one replaced.
const legacyListForeignTableQuery = `SELECT
	foreign_table.foreign_table_schema,
	foreign_table.foreign_table_name,
	foreign_table.foreign_server_catalog,
	foreign_table.foreign_server_name
FROM information_schema.foreign_tables AS foreign_table
ORDER BY foreign_table.foreign_table_schema, foreign_table.foreign_table_name;`

const foreignTableSetup = `
CREATE EXTENSION postgres_fdw;
CREATE SERVER loopback_server FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'localhost', dbname 'postgres');
CREATE SCHEMA fdw_schema;
CREATE FOREIGN TABLE fdw_schema.remote_orders (
    id integer NOT NULL,
    note varchar(30)
) SERVER loopback_server OPTIONS (schema_name 'public', table_name 'orders');
CREATE FOREIGN TABLE public.remote_items (
    item_id bigint
) SERVER loopback_server OPTIONS (schema_name 'public', table_name 'items');
`

func TestListForeignTableQueryMatchesInformationSchema(t *testing.T) {
	d := openDriver(t, testcontainer.Default)
	mustExec(t, d.db, foreignTableSetup)
	legacy := queryRows(t, d.db, legacyListForeignTableQuery, 2)
	current := queryRows(t, d.db, getListForeignTableQuery(), 2)
	if len(legacy) != 2 {
		t.Fatalf("information_schema returned %d foreign tables, want 2", len(legacy))
	}
	if diff := cmp.Diff(legacy, current); diff != "" {
		t.Errorf("foreign table query rows (-information_schema +pg_catalog):\n%s", diff)
	}
}

// openUnprivileged returns a driver whose one connection runs as a new
// role that holds pg_read_all_stats and no privilege on any table:
// information_schema hides every table from it, pg_catalog does not.
func openUnprivileged(t *testing.T, setup string) *Driver {
	t.Helper()
	d := openDriver(t, testcontainer.Default)
	mustExec(t, d.db, setup+`
CREATE ROLE limited LOGIN;
GRANT pg_read_all_stats TO limited;`)
	d.db.SetMaxOpenConns(1)
	d.db.SetMaxIdleConns(1)
	mustExec(t, d.db, "SET ROLE limited")
	return d
}

func countRows(t *testing.T, db *sql.DB, query string) int {
	t.Helper()
	var n int
	if err := db.QueryRow(query).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func TestSyncColumnsWithoutTablePrivilege(t *testing.T) {
	d := openUnprivileged(t, `
CREATE SCHEMA restricted;
CREATE TABLE restricted.secret_table (
    id integer PRIMARY KEY,
    name varchar(50) NOT NULL DEFAULT 'anon',
    amount numeric(10, 2),
    CONSTRAINT positive_amount CHECK (amount >= 0)
);
CREATE INDEX idx_secret_name ON restricted.secret_table (name);`)
	if n := countRows(t, d.db, `SELECT count(*) FROM information_schema.columns WHERE table_schema = 'restricted'`); n != 0 {
		t.Fatalf("information_schema.columns shows the unprivileged role %d columns, want 0", n)
	}
	m, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	table := findTable(t, m, "restricted", "secret_table")
	if len(table.Indexes) == 0 || len(table.CheckConstraints) == 0 {
		t.Errorf("secret_table synced with %d indexes and %d checks, want both", len(table.Indexes), len(table.CheckConstraints))
	}
	var got []string
	for _, c := range table.Columns {
		got = append(got, fmt.Sprintf("%s %s nullable=%v default=%s", c.Name, c.Type, c.Nullable, c.Default))
	}
	want := []string{
		"id integer nullable=false default=",
		"name character varying(50) nullable=false default='anon'::character varying",
		"amount numeric(10,2) nullable=true default=",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("columns (-want +got):\n%s", diff)
	}
}

func TestSyncForeignTablesWithoutTablePrivilege(t *testing.T) {
	d := openUnprivileged(t, foreignTableSetup)
	if n := countRows(t, d.db, `SELECT count(*) FROM information_schema.foreign_tables`); n != 0 {
		t.Fatalf("information_schema.foreign_tables shows the unprivileged role %d tables, want 0", n)
	}
	m, err := d.SyncDatabase(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, s := range m.Schemas {
		for _, e := range s.ExternalTables {
			got = append(got, fmt.Sprintf("%s.%s server=%s database=%s", s.Name, e.Name, e.ExternalServerName, e.ExternalDatabaseName))
			for _, c := range e.Columns {
				got = append(got, fmt.Sprintf("  %s %s nullable=%v", c.Name, c.Type, c.Nullable))
			}
		}
	}
	want := []string{
		"fdw_schema.remote_orders server=loopback_server database=" + d.databaseName,
		"  id integer nullable=false",
		"  note character varying(30) nullable=true",
		"public.remote_items server=loopback_server database=" + d.databaseName,
		"  item_id bigint nullable=true",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("foreign tables (-want +got):\n%s", diff)
	}
}
