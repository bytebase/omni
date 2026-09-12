package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	pg "github.com/bytebase/omni/pg"
	nodes "github.com/bytebase/omni/pg/ast"
)

// liveSchemaDDL exercises the shapes a real snapshot carries that a
// hand-written fixture rarely does: a length-qualified varchar, a scaled
// numeric, an array, a jsonb default that is a cast, an enum type used as a
// column type, a nextval default naming a sequence, an identity column, a
// stored generated column, and a view over it all.
const liveSchemaDDL = `
CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy');
CREATE SEQUENCE ticket_seq INCREMENT 5 START 100;
CREATE TABLE authors (
	id        bigserial PRIMARY KEY,
	handle    varchar(50) NOT NULL UNIQUE,
	bio       text,
	joined    timestamptz NOT NULL DEFAULT now(),
	karma     numeric(10,2) NOT NULL DEFAULT 0.0,
	tags      text[],
	prefs     jsonb NOT NULL DEFAULT '{}'::jsonb,
	ext       uuid,
	state     mood NOT NULL DEFAULT 'ok',
	ticket    bigint NOT NULL DEFAULT nextval('ticket_seq'),
	is_active boolean NOT NULL DEFAULT true
);
CREATE TABLE posts (
	id        int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	author_id bigint NOT NULL REFERENCES authors(id),
	title     varchar(200) NOT NULL,
	body      text,
	words     int GENERATED ALWAYS AS (length(body)) STORED,
	created   date NOT NULL DEFAULT CURRENT_DATE
);
CREATE INDEX posts_author_idx ON posts (author_id);
CREATE VIEW active_authors AS SELECT id, handle, karma FROM authors WHERE is_active;
`

// TestContainerLoadMetadataFromLiveSchema loads a snapshot introspected from a
// real PostgreSQL schema. Every other LoadMetadata test writes its own
// snapshot, which only ever contains what the author expected the engine to
// emit; this one makes the loader answer to what PostgreSQL actually reports —
// its type vocabulary, its default expressions, its identity and generation
// markers.
func TestContainerLoadMetadataFromLiveSchema(t *testing.T) {
	ctr := startPGContainer(t)
	schema := ctr.freshSchema(t)
	ctr.execInSchema(t, schema, liveSchemaDDL)

	meta := liveSnapshot(t, ctr, schema)

	c := New()
	report, err := c.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	if len(report.Degraded) != 0 || len(report.Missing) != 0 {
		t.Fatalf("a snapshot taken from a live schema must install whole; degraded=%v missing=%v", report.Degraded, report.Missing)
	}

	// The engine's own type spelling must survive the round trip: whatever
	// format_type reported going in, FormatType reports coming out. A type the
	// loader cannot parse would otherwise install as a text stand-in and only
	// surface later as wrong lineage.
	for _, want := range meta.GetSchemas()[0].GetTables() {
		rel := c.GetRelation(schema, want.GetName())
		if rel == nil {
			t.Errorf("table %s did not install", want.GetName())
			continue
		}
		for _, wantCol := range want.GetColumns() {
			col := relationColumn(rel, wantCol.GetName())
			if col == nil {
				t.Errorf("%s.%s did not install", want.GetName(), wantCol.GetName())
				continue
			}
			if got := c.FormatType(col.TypeOID, col.TypeMod); got != wantCol.GetType() {
				t.Errorf("%s.%s type = %q, want the engine's %q", want.GetName(), wantCol.GetName(), got, wantCol.GetType())
			}
			if col.NotNull == wantCol.GetNullable() {
				t.Errorf("%s.%s nullable = %v, want %v", want.GetName(), wantCol.GetName(), !col.NotNull, wantCol.GetNullable())
			}
		}
	}

	// The loaded catalog must answer a real query, including through the view
	// and across the foreign key.
	stmts, err := pg.Parse(`
		SELECT a.handle, a.karma, p.title, p.words
		FROM posts p JOIN active_authors a ON a.id = p.author_id`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	stmt, ok := stmts[0].AST.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected a SelectStmt, got %T", stmts[0].AST)
	}
	c.SetSearchPath([]string{schema})
	query, err := c.AnalyzeSelectStmt(stmt)
	if err != nil {
		t.Fatalf("analyzing a query against the loaded catalog: %v", err)
	}
	var names []string
	for _, target := range query.TargetList {
		names = append(names, target.ResName)
	}
	if strings.Join(names, ",") != "handle,karma,title,words" {
		t.Errorf("query resolved to %v, want handle, karma, title, words", names)
	}
}

func relationColumn(rel *Relation, name string) *Column {
	for _, col := range rel.Columns {
		if col.Name == name {
			return col
		}
	}
	return nil
}

// liveSnapshot introspects one schema into a Bytebase-shaped snapshot. It asks
// for column types with format_type, as sync does, so the snapshot carries the
// engine's own spelling ("character varying(50)") rather than the unqualified
// name information_schema reports.
func liveSnapshot(t *testing.T, ctr *pgContainer, schema string) *metadata.DatabaseSchemaMetadata {
	t.Helper()
	// Sync empties the search path for the duration of the read so that every
	// type and default expression comes back schema-qualified. The loader
	// relies on that: an unqualified name in a default has no schema to
	// resolve against once the snapshot leaves the connection that produced it.
	txn, err := ctr.db.BeginTx(ctr.ctx, nil)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer txn.Rollback()
	if _, err := txn.ExecContext(ctr.ctx, `SELECT pg_catalog.set_config('search_path', '', true)`); err != nil {
		t.Fatalf("clear the search path: %v", err)
	}

	sm := &metadata.SchemaMetadata{Name: schema}

	tables := map[string]*metadata.TableMetadata{}
	rows, err := txn.QueryContext(ctr.ctx, `
		SELECT c.relname, a.attname,
		       pg_catalog.format_type(a.atttypid, a.atttypmod), a.attnotnull, pg_get_expr(d.adbin, d.adrelid), a.attnum,
		       a.attidentity::text, a.attgenerated::text
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		JOIN pg_attribute a ON a.attrelid = c.oid
		LEFT JOIN pg_attrdef d ON d.adrelid = c.oid AND d.adnum = a.attnum
		WHERE n.nspname = $1 AND c.relkind = 'r' AND a.attnum > 0 AND NOT a.attisdropped
		ORDER BY c.relname, a.attnum`, schema)
	if err != nil {
		t.Fatalf("introspect columns: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var table, name, typ, identity, generated string
		var notNull bool
		var def sql.NullString
		var attnum int32
		if err := rows.Scan(&table, &name, &typ, &notNull, &def, &attnum, &identity, &generated); err != nil {
			t.Fatalf("scan column: %v", err)
		}
		tbl, ok := tables[table]
		if !ok {
			tbl = &metadata.TableMetadata{Name: table}
			tables[table] = tbl
			sm.Tables = append(sm.Tables, tbl)
		}
		col := &metadata.ColumnMetadata{Name: name, Type: typ, Nullable: !notNull, Position: attnum}
		switch {
		case identity == "a":
			col.IsIdentity, col.IdentityGeneration = true, metadata.ColumnMetadata_ALWAYS
		case identity == "d":
			col.IsIdentity, col.IdentityGeneration = true, metadata.ColumnMetadata_BY_DEFAULT
		case generated == "s" && def.Valid:
			col.Generation = &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_STORED, Expression: def.String}
		case def.Valid:
			col.Default = def.String
		}
		tbl.Columns = append(tbl.Columns, col)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("introspect columns: %v", err)
	}

	for _, tbl := range sm.Tables {
		tbl.Indexes = liveIndexes(t, ctr, txn, schema, tbl.GetName())
	}
	for _, e := range ctr.queryEnumTypes(t, schema) {
		sm.EnumTypes = append(sm.EnumTypes, &metadata.EnumTypeMetadata{Name: e.name, Values: strings.Split(e.values, ",")})
	}
	sm.Sequences = liveSequences(t, ctr, txn, schema)
	for _, v := range ctr.queryViews(t, schema) {
		sm.Views = append(sm.Views, &metadata.ViewMetadata{Name: v.name, Definition: v.definition})
	}
	return &metadata.DatabaseSchemaMetadata{Name: "omni_test", Schemas: []*metadata.SchemaMetadata{sm}}
}

func liveIndexes(t *testing.T, ctr *pgContainer, txn *sql.Tx, schema, table string) []*metadata.IndexMetadata {
	t.Helper()
	rows, err := txn.QueryContext(ctr.ctx, `
		SELECT i.relname, ix.indisunique, ix.indisprimary,
		       array_to_string(array_agg(a.attname ORDER BY k.ord), ',')
		FROM pg_class t
		JOIN pg_namespace n ON n.oid = t.relnamespace
		JOIN pg_index ix ON ix.indrelid = t.oid
		JOIN pg_class i ON i.oid = ix.indexrelid
		JOIN LATERAL unnest(ix.indkey) WITH ORDINALITY AS k(attnum, ord) ON true
		JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = k.attnum
		WHERE n.nspname = $1 AND t.relname = $2
		GROUP BY i.relname, ix.indisunique, ix.indisprimary
		ORDER BY i.relname`, schema, table)
	if err != nil {
		t.Fatalf("introspect indexes: %v", err)
	}
	defer rows.Close()
	var out []*metadata.IndexMetadata
	for rows.Next() {
		var name, columns string
		var unique, primary bool
		if err := rows.Scan(&name, &unique, &primary, &columns); err != nil {
			t.Fatalf("scan index: %v", err)
		}
		out = append(out, &metadata.IndexMetadata{
			Name: name, Unique: unique, Primary: primary, Type: "btree",
			Expressions: strings.Split(columns, ","),
		})
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("introspect indexes: %v", err)
	}
	return out
}

// liveSequences introspects sequences with their ownership. Sync records the
// owning table and column, and the loader needs them: a sequence that backs a
// serial or identity column is installed with that column, not on its own.
func liveSequences(t *testing.T, ctr *pgContainer, txn *sql.Tx, schema string) []*metadata.SequenceMetadata {
	t.Helper()
	rows, err := txn.QueryContext(ctr.ctx, `
		SELECT s.relname, pg_catalog.format_type(sq.seqtypid, NULL),
		       sq.seqstart, sq.seqmin, sq.seqmax, sq.seqincrement, sq.seqcycle,
		       COALESCE(dc.relname, ''), COALESCE(da.attname, '')
		FROM pg_class s
		JOIN pg_namespace sn ON sn.oid = s.relnamespace
		JOIN pg_sequence sq ON sq.seqrelid = s.oid
		LEFT JOIN pg_depend dep ON dep.objid = s.oid
			AND dep.classid = 'pg_class'::regclass AND dep.deptype IN ('a', 'i')
		LEFT JOIN pg_class dc ON dc.oid = dep.refobjid
		LEFT JOIN pg_attribute da ON da.attrelid = dc.oid AND da.attnum = dep.refobjsubid
		WHERE sn.nspname = $1 AND s.relkind = 'S'
		ORDER BY s.relname`, schema)
	if err != nil {
		t.Fatalf("introspect sequences: %v", err)
	}
	defer rows.Close()
	var out []*metadata.SequenceMetadata
	for rows.Next() {
		var name, dataType, ownerTable, ownerColumn string
		var start, minValue, maxValue, increment int64
		var cycle bool
		if err := rows.Scan(&name, &dataType, &start, &minValue, &maxValue, &increment, &cycle, &ownerTable, &ownerColumn); err != nil {
			t.Fatalf("scan sequence: %v", err)
		}
		out = append(out, &metadata.SequenceMetadata{
			Name: name, DataType: dataType,
			Start:       fmt.Sprintf("%d", start),
			MinValue:    fmt.Sprintf("%d", minValue),
			MaxValue:    fmt.Sprintf("%d", maxValue),
			Increment:   fmt.Sprintf("%d", increment),
			Cycle:       cycle,
			OwnerTable:  ownerTable,
			OwnerColumn: ownerColumn,
		})
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("introspect sequences: %v", err)
	}
	return out
}
