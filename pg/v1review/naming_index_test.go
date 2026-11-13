package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

// keyFinding is a finding of an index or key naming rule as the tests
// compare it.
type keyFinding struct {
	Code    v1review.Code
	Table   string
	Name    string
	Pattern string
	Detail  string
	Text    string
}

func keyMismatch(table, name, pattern, text string) keyFinding {
	return keyFinding{Code: v1review.NameMismatch, Table: table, Name: name, Pattern: pattern, Text: text}
}

type keyRule func([]ast.Node, v1review.TemplateConfig, v1review.Schema) []v1review.Finding

func runKeyRule(t *testing.T, rule keyRule, sql string, cfg v1review.TemplateConfig, schema v1review.Schema) []keyFinding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var findings []keyFinding
	for _, f := range rule(list.Items, cfg, schema) {
		findings = append(findings, keyFinding{f.Code, f.Table, f.Name, f.Pattern, f.Detail, sql[f.Start:f.End]})
	}
	return findings
}

func snapshotIndex(name string, primary, unique bool, columns ...string) *metadata.IndexMetadata {
	return &metadata.IndexMetadata{Name: name, Primary: primary, Unique: unique, Expressions: columns}
}

// keySchema is the pre-change schema the index and key naming tests look
// names up in.
func keySchema(objectCaseSensitive bool) v1review.Schema {
	return v1review.Schema{
		ObjectCaseSensitive: objectCaseSensitive,
		Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadata.TableMetadata{
					{Name: "tech_book", Indexes: []*metadata.IndexMetadata{
						snapshotIndex("old_pk", true, true, "id", "name"),
						snapshotIndex("old_uk", false, true, "id", "name"),
						snapshotIndex("old_index", false, false, "id", "name"),
						snapshotIndex("Old_Index", false, false, "x"),
						snapshotIndex("Old_Uk", false, true, "y"),
						snapshotIndex("Old_Pk", true, true, "z"),
					}},
					{Name: "Book", Indexes: []*metadata.IndexMetadata{
						snapshotIndex("book_pk", true, true, "a"),
						snapshotIndex("book_uk", false, true, "b"),
						snapshotIndex("book_idx", false, false, "c"),
					}},
					{Name: "twin_a", Indexes: []*metadata.IndexMetadata{
						snapshotIndex("twin_idx", false, false, "a"),
						snapshotIndex("twin_uk", false, true, "a"),
						snapshotIndex("twin_pk", true, true, "a"),
					}},
					{Name: "twin_b", Indexes: []*metadata.IndexMetadata{
						snapshotIndex("twin_idx", false, false, "b"),
						snapshotIndex("twin_uk", false, true, "b"),
						snapshotIndex("twin_pk", true, true, "b"),
					}},
					{
						Name:       "parted",
						Partitions: []*metadata.TablePartitionMetadata{{Name: "parted_p0"}},
						Indexes: []*metadata.IndexMetadata{
							snapshotIndex("parted_pk", true, true, "k"),
							snapshotIndex("parted_uk", false, true, "k"),
							snapshotIndex("parted_idx", false, false, "k"),
						},
					},
				},
				MaterializedViews: []*metadata.MaterializedViewMetadata{
					{Name: "mv", Indexes: []*metadata.IndexMetadata{
						snapshotIndex("mv_idx", false, false, "m"),
						snapshotIndex("mv_uk", false, true, "m"),
					}},
				},
			},
			{
				Name: "other",
				Tables: []*metadata.TableMetadata{
					{Name: "tech_book", Indexes: []*metadata.IndexMetadata{
						snapshotIndex("old_index", false, true, "o"),
						snapshotIndex("old_uk", true, true, "o"),
						snapshotIndex("other_pk", true, true, "o"),
						snapshotIndex("other_uk", false, true, "o"),
						snapshotIndex("other_idx", false, false, "o"),
					}},
				},
			},
		}},
	}
}

func TestNamingIndex(t *testing.T) {
	format := v1review.TemplateConfig{Format: "^idx_{{table}}_{{column_list}}$"}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		schema v1review.Schema
		want   []keyFinding
	}{
		{
			name:   "create index ok",
			sql:    "CREATE INDEX idx_tech_book_id_name ON tech_book (id, name);",
			config: format,
		},
		{
			name:   "create index mismatch",
			sql:    "CREATE INDEX tech_book_id_name ON tech_book (id, name);",
			config: format,
			want:   []keyFinding{keyMismatch("tech_book", "tech_book_id_name", "^idx_tech_book_id_name$", "CREATE INDEX tech_book_id_name ON tech_book (id, name)")},
		},
		{
			name:   "unique and unnamed indexes not judged",
			sql:    "CREATE UNIQUE INDEX bad ON t (a); CREATE INDEX ON t (a); CREATE INDEX CONCURRENTLY IF NOT EXISTS bad2 ON t (a);",
			config: format,
			want:   []keyFinding{keyMismatch("t", "bad2", "^idx_t_a$", "CREATE INDEX CONCURRENTLY IF NOT EXISTS bad2 ON t (a)")},
		},
		{
			name:   "table loses its schema, columns keep their stored spelling, expressions are left out",
			sql:    `CREATE INDEX i ON "S"."Tbl" USING btree (A DESC, lower(b), "C" NULLS FIRST, (d + 1)) INCLUDE (e) WHERE a > 0;`,
			config: format,
			want: []keyFinding{keyMismatch("Tbl", "i", "^idx_Tbl_a_C$",
				`CREATE INDEX i ON "S"."Tbl" USING btree (A DESC, lower(b), "C" NULLS FIRST, (d + 1)) INCLUDE (e) WHERE a > 0`)},
		},
		{
			name:   "index over expressions only has an empty column list",
			sql:    "CREATE INDEX idx_t_ ON t (lower(a)); CREATE INDEX idx_t_a ON t (lower(a));",
			config: format,
			want:   []keyFinding{keyMismatch("t", "idx_t_a", "^idx_t_$", "CREATE INDEX idx_t_a ON t (lower(a))")},
		},
		{
			name:   "create index does not consult the schema",
			sql:    "CREATE INDEX bad ON tech_book (id);",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("tech_book", "bad", "^idx_tech_book_id$", "CREATE INDEX bad ON tech_book (id)")},
		},
		{
			name:   "renamed index takes the snapshot's table and columns",
			sql:    "ALTER INDEX old_index RENAME TO idx_tech_book_id_name; ALTER INDEX IF EXISTS old_index RENAME TO idx_tech_book;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("tech_book", "idx_tech_book", "^idx_tech_book_id_name$", "ALTER INDEX IF EXISTS old_index RENAME TO idx_tech_book")},
		},
		{
			name:   "renamed unique, primary, and missing indexes not judged",
			sql:    "ALTER INDEX old_uk RENAME TO bad; ALTER INDEX old_pk RENAME TO bad; ALTER INDEX nope RENAME TO bad;",
			config: format,
			schema: keySchema(true),
		},
		{
			name:   "rename without a snapshot not judged",
			sql:    "ALTER INDEX old_index RENAME TO bad;",
			config: format,
		},
		{
			name:   "rename looks in public whatever schema the statement names",
			sql:    "ALTER INDEX other.old_index RENAME TO bad; ALTER INDEX other.other_idx RENAME TO bad2; ALTER INDEX other_idx RENAME TO bad3;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("tech_book", "bad", "^idx_tech_book_id_name$", "ALTER INDEX other.old_index RENAME TO bad")},
		},
		{
			name:   "old index name is case sensitive",
			sql:    `ALTER INDEX "Old_Index" RENAME TO bad; ALTER INDEX "OLD_INDEX" RENAME TO bad2; ALTER INDEX OLD_INDEX RENAME TO bad3;`,
			config: format,
			schema: keySchema(false),
			want: []keyFinding{
				keyMismatch("tech_book", "bad", "^idx_tech_book_x$", `ALTER INDEX "Old_Index" RENAME TO bad`),
				keyMismatch("tech_book", "bad3", "^idx_tech_book_id_name$", "ALTER INDEX OLD_INDEX RENAME TO bad3"),
			},
		},
		{
			name:   "table keeps the snapshot's spelling",
			sql:    "ALTER INDEX book_idx RENAME TO bad;",
			config: format,
			schema: keySchema(false),
			want:   []keyFinding{keyMismatch("Book", "bad", "^idx_Book_c$", "ALTER INDEX book_idx RENAME TO bad")},
		},
		{
			name:   "partitioned table's index",
			sql:    "ALTER INDEX parted_idx RENAME TO bad;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("parted", "bad", "^idx_parted_k$", "ALTER INDEX parted_idx RENAME TO bad")},
		},
		{
			name:   "materialized view's index",
			sql:    "ALTER INDEX mv_idx RENAME TO bad; ALTER INDEX mv_uk RENAME TO bad2;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("mv", "bad", "^idx_mv_m$", "ALTER INDEX mv_idx RENAME TO bad")},
		},
		{
			name:   "materialized view's index name follows object case sensitivity",
			sql:    `ALTER INDEX "MV_IDX" RENAME TO bad;`,
			config: format,
			schema: keySchema(false),
			want:   []keyFinding{keyMismatch("mv", "bad", "^idx_mv_m$", `ALTER INDEX "MV_IDX" RENAME TO bad`)},
		},
		{
			name:   "materialized view's index name case sensitive when objects are",
			sql:    `ALTER INDEX "MV_IDX" RENAME TO bad;`,
			config: format,
			schema: keySchema(true),
		},
		{
			name:   "two tables hold the name, the first in snapshot order is taken",
			sql:    "ALTER INDEX twin_idx RENAME TO bad;",
			config: format,
			schema: keySchema(true),
			want:   []keyFinding{keyMismatch("twin_a", "bad", "^idx_twin_a_a$", "ALTER INDEX twin_idx RENAME TO bad")},
		},
		{
			name:   "other renames and alters not judged",
			sql:    "ALTER TABLE old_index RENAME TO bad; ALTER TABLE tech_book RENAME CONSTRAINT old_index TO bad; ALTER INDEX old_index SET (fillfactor = 50); ALTER TABLE tech_book ADD CONSTRAINT bad UNIQUE (id);",
			config: format,
			schema: keySchema(true),
		},
		{
			name:   "create schema element and routine body not judged",
			sql:    "CREATE SCHEMA s CREATE TABLE t (a int) CREATE INDEX bad ON t (a); CREATE PROCEDURE pr() LANGUAGE SQL BEGIN ATOMIC CREATE INDEX bad ON t (a); END;",
			config: format,
		},
		{
			name:   "unquoted names fold to lower",
			sql:    "CREATE INDEX IDX_Book_Id ON Book (Id);",
			config: format,
		},
		{
			name:   "value with a metacharacter changes the pattern",
			sql:    `CREATE INDEX idx_abc_x ON "a.c" (x); CREATE INDEX "idx_a+c_x" ON "a+c" (x);`,
			config: format,
			want:   []keyFinding{keyMismatch("a+c", "idx_a+c_x", "^idx_a+c_x$", `CREATE INDEX "idx_a+c_x" ON "a+c" (x)`)},
		},
		{
			name:   "value that keeps the pattern from compiling",
			sql:    `CREATE INDEX "` + strings.Repeat("k", 20) + `" ON "a(" (x); CREATE INDEX nope ON t (x);`,
			config: v1review.TemplateConfig{Format: "^idx_{{table}}$", MaxLength: 10},
			want: []keyFinding{
				{
					Code: v1review.PatternInvalid, Table: "a(", Name: strings.Repeat("k", 20),
					Detail: "error parsing regexp: missing closing ): `^idx_a($`",
					Text:   `CREATE INDEX "` + strings.Repeat("k", 20) + `" ON "a(" (x)`,
				},
				keyMismatch("t", "nope", "^idx_t$", "CREATE INDEX nope ON t (x)"),
			},
		},
		{
			name:   "token with no value is left as written",
			sql:    "CREATE INDEX idx_t ON t (x);",
			config: v1review.TemplateConfig{Format: "^idx_{{referencing_table}}$"},
			want:   []keyFinding{keyMismatch("t", "idx_t", "^idx_{{referencing_table}}$", "CREATE INDEX idx_t ON t (x)")},
		},
		{
			name:   "length counts bytes",
			sql:    `CREATE INDEX "索引名" ON t (x); CREATE INDEX abcdefgh ON t (x);`,
			config: v1review.TemplateConfig{MaxLength: 8},
			want:   []keyFinding{{Code: v1review.NameTooLong, Table: "t", Name: "索引名", Text: `CREATE INDEX "索引名" ON t (x)`}},
		},
		{
			name:   "name truncated to 63 bytes before judging",
			sql:    `CREATE INDEX "` + strings.Repeat("a", 70) + `" ON t (x);`,
			config: v1review.TemplateConfig{Format: "^a{63}$", MaxLength: 63},
		},
		{
			name:   "pattern then length",
			sql:    "ALTER INDEX old_index RENAME TO toolong;",
			config: v1review.TemplateConfig{Format: "^idx$", MaxLength: 4},
			schema: keySchema(true),
			want: []keyFinding{
				keyMismatch("tech_book", "toolong", "^idx$", "ALTER INDEX old_index RENAME TO toolong"),
				{Code: v1review.NameTooLong, Table: "tech_book", Name: "toolong", Text: "ALTER INDEX old_index RENAME TO toolong"},
			},
		},
		{name: "empty config requires nothing", sql: "CREATE INDEX anything ON t (x);"},
		{
			name:   "anchor excludes a leading comment and the semicolon",
			sql:    "-- c\nCREATE INDEX\n  bad ON t (x);",
			config: format,
			want:   []keyFinding{keyMismatch("t", "bad", "^idx_t_x$", "CREATE INDEX\n  bad ON t (x)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runKeyRule(t, NamingIndex, tt.sql, tt.config, tt.schema)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
