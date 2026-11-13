package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

// keySchema is a pre-change schema holding one index of each kind on Book,
// whose name differs from a statement's `book` only by case.
func keySchema(objectCaseSensitive bool) v1review.Schema {
	index := func(name string, unique, primary bool, expressions ...string) *metadata.IndexMetadata {
		return &metadata.IndexMetadata{Name: name, Unique: unique, Primary: primary, Expressions: expressions}
	}
	return v1review.Schema{
		ObjectCaseSensitive: objectCaseSensitive,
		Metadata: &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name: "Book",
					Indexes: []*metadata.IndexMetadata{
						index("PRIMARY", true, true, "id"),
						index("Old_Uk", true, false, "id", "name"),
						index("Old_Idx", false, false, "id", "name"),
						index("old_fulltext", false, false, "body"),
						index("old_expr", false, false, "(`id` + 1)"),
					},
					Partitions: []*metadata.TablePartitionMetadata{{Name: "p0"}},
				},
			},
		}}},
	}
}

func TestNamingIndex(t *testing.T) {
	format := v1review.TemplateConfig{Format: "^$|^idx_{{table}}_{{column_list}}$"}
	strict := v1review.TemplateConfig{Format: "^idx_{{table}}_{{column_list}}$"}
	mismatch := func(table, name, pattern, text string) located {
		return located{Code: v1review.NameMismatch, Name: name, Table: table, Pattern: pattern, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		schema v1review.Schema
		want   []located
	}{
		{name: "create index ok", sql: "CREATE INDEX idx_book_id_name ON book (id, name);", config: strict},
		{
			name:   "create index anchors the statement",
			sql:    "CREATE INDEX\n  bad\n  ON shop.book (id(10) DESC, name);",
			config: strict,
			want:   []located{mismatch("book", "bad", "^idx_book_id_name$", "CREATE INDEX\n  bad\n  ON shop.book (id(10) DESC, name)")},
		},
		{name: "create unique, fulltext and spatial index not judged", sql: "CREATE UNIQUE INDEX bad ON t (a); CREATE FULLTEXT INDEX bad ON t (a); CREATE SPATIAL INDEX bad ON t (a);", config: strict},
		{
			name:   "create index counts a parenthesized column, not an expression",
			sql:    "CREATE INDEX bad ON t (a, (b), (t.c), (a + 1));",
			config: strict,
			want:   []located{mismatch("t", "bad", "^idx_t_a_b_c$", "CREATE INDEX bad ON t (a, (b), (t.c), (a + 1))")},
		},
		{name: "create table ok", sql: "CREATE TABLE book (id INT, name TEXT, INDEX idx_book_id_name (id, name), KEY idx_book_id (id));", config: strict},
		{name: "unnamed matches an empty alternative", sql: "CREATE TABLE book (id INT, INDEX (id));", config: format},
		{
			name:   "unnamed judged as the empty name",
			sql:    "CREATE TABLE book (id INT, INDEX (id));",
			config: strict,
			want:   []located{mismatch("book", "", "^idx_book_id$", "INDEX (id)")},
		},
		{
			name:   "create table anchors the definition",
			sql:    "CREATE TABLE shop.book (\n  id INT,\n  KEY bad\n    USING BTREE (id) COMMENT 'x',\n  INDEX bad2 (id)\n);",
			config: strict,
			want: []located{
				mismatch("book", "bad", "^idx_book_id$", "KEY bad\n    USING BTREE (id) COMMENT 'x'"),
				mismatch("book", "bad2", "^idx_book_id$", "INDEX bad2 (id)"),
			},
		},
		{
			name:   "a parenthesized column counts only when no part is a plain column",
			sql:    "CREATE TABLE t (a INT, b INT, INDEX i1 ((a)), INDEX i2 (a, (b)), INDEX i3 ((a + 1)), INDEX i4 ((t.a), (b)));",
			config: strict,
			want: []located{
				mismatch("t", "i1", "^idx_t_a$", "INDEX i1 ((a))"),
				mismatch("t", "i2", "^idx_t_a$", "INDEX i2 (a, (b))"),
				mismatch("t", "i3", "^idx_t_$", "INDEX i3 ((a + 1))"),
				mismatch("t", "i4", "^idx_t_a_b$", "INDEX i4 ((t.a), (b))"),
			},
		},
		{
			name:   "other constraints not judged",
			sql:    "CREATE TABLE t (a INT UNIQUE KEY, PRIMARY KEY (a), UNIQUE KEY bad (a), FULLTEXT bad2 (a), SPATIAL INDEX bad3 (a), CONSTRAINT bad4 FOREIGN KEY bad5 (a) REFERENCES p (id), CONSTRAINT bad6 CHECK (a > 0));",
			config: strict,
		},
		{
			name:   "alter add anchors the statement",
			sql:    "ALTER TABLE shop.book\n  ADD INDEX idx_book_id (id),\n  ADD KEY bad (id, name),\n  ADD UNIQUE bad2 (id);",
			config: strict,
			want: []located{mismatch("book", "bad", "^idx_book_id_name$",
				"ALTER TABLE shop.book\n  ADD INDEX idx_book_id (id),\n  ADD KEY bad (id, name),\n  ADD UNIQUE bad2 (id)")},
		},
		{name: "alter add column list index not judged", sql: "ALTER TABLE t ADD COLUMN (z INT, INDEX bad (z)); ALTER TABLE t ADD (INDEX bad (z));", config: strict},
		{name: "alter drop index not judged", sql: "ALTER TABLE t DROP INDEX bad;", config: strict},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE INDEX bad ON t (a); END;", config: strict},
		{
			name:   "rename takes the columns from the schema",
			sql:    "ALTER TABLE book\n  RENAME INDEX old_idx TO idx_book_id_name,\n  RENAME KEY OLD_IDX TO bad;",
			config: strict,
			schema: keySchema(false),
			want: []located{mismatch("book", "bad", "^idx_book_id_name$",
				"ALTER TABLE book\n  RENAME INDEX old_idx TO idx_book_id_name,\n  RENAME KEY OLD_IDX TO bad")},
		},
		{name: "rename of a unique or primary index not judged", sql: "ALTER TABLE book RENAME INDEX old_uk TO bad, RENAME INDEX `PRIMARY` TO bad;", config: strict, schema: keySchema(false)},
		{name: "rename of a missing index not judged", sql: "ALTER TABLE book RENAME INDEX nope TO bad;", config: strict, schema: keySchema(false)},
		{name: "rename on a missing table not judged", sql: "ALTER TABLE nope RENAME INDEX old_idx TO bad;", config: strict, schema: keySchema(false)},
		{name: "rename on a partition's name not judged", sql: "ALTER TABLE p0 RENAME INDEX old_idx TO bad;", config: strict, schema: keySchema(false)},
		{name: "rename without a schema not judged", sql: "ALTER TABLE book RENAME INDEX old_idx TO bad;", config: strict},
		{name: "rename misses a table whose case differs where that matters", sql: "ALTER TABLE book RENAME INDEX old_idx TO bad;", config: strict, schema: keySchema(true)},
		{
			name:   "rename finds an index whose case differs even so",
			sql:    "ALTER TABLE Book RENAME INDEX OLD_IDX TO bad;",
			config: strict,
			schema: keySchema(true),
			want:   []located{mismatch("Book", "bad", "^idx_Book_id_name$", "ALTER TABLE Book RENAME INDEX OLD_IDX TO bad")},
		},
		{
			name:   "rename ignores the database the statement names",
			sql:    "ALTER TABLE elsewhere.book RENAME INDEX old_idx TO bad;",
			config: strict,
			schema: keySchema(false),
			want:   []located{mismatch("book", "bad", "^idx_book_id_name$", "ALTER TABLE elsewhere.book RENAME INDEX old_idx TO bad")},
		},
		{
			name:   "rename of a fulltext index judged, of one the script created not",
			sql:    "CREATE INDEX idx_book_name ON book (name); ALTER TABLE book RENAME INDEX idx_book_name TO bad, RENAME INDEX old_fulltext TO bad2;",
			config: strict,
			schema: keySchema(false),
			want: []located{mismatch("book", "bad2", "^idx_book_body$",
				"ALTER TABLE book RENAME INDEX idx_book_name TO bad, RENAME INDEX old_fulltext TO bad2")},
		},
		{
			name:   "rename puts a stored expression into the pattern as it is",
			sql:    "ALTER TABLE book RENAME INDEX old_expr TO bad;",
			config: v1review.TemplateConfig{Format: "^idx_{{column_list}}$", MaxLength: 2},
			schema: keySchema(false),
			want:   []located{mismatch("book", "bad", "^idx_(`id` + 1)$", "ALTER TABLE book RENAME INDEX old_expr TO bad"), {Code: v1review.NameTooLong, Name: "bad", Table: "book", Text: "ALTER TABLE book RENAME INDEX old_expr TO bad"}},
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE t (a INT, INDEX `索引名` (a), INDEX `索引` (a));",
			config: v1review.TemplateConfig{MaxLength: 8},
			want:   []located{{Code: v1review.NameTooLong, Name: "索引名", Table: "t", Text: "INDEX `索引名` (a)"}},
		},
		{
			name:   "pattern and length both reported, pattern first",
			sql:    "CREATE INDEX toolong ON t (a);",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}$", MaxLength: 4},
			want: []located{
				mismatch("t", "toolong", "^idx_t$", "CREATE INDEX toolong ON t (a)"),
				{Code: v1review.NameTooLong, Name: "toolong", Table: "t", Text: "CREATE INDEX toolong ON t (a)"},
			},
		},
		{
			name:   "a name is not quoted into the pattern",
			sql:    "CREATE TABLE `a.c` (x INT, INDEX idx_abc (x), INDEX idx_a (x));",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}$"},
			want:   []located{mismatch("a.c", "idx_a", "^idx_a.c$", "INDEX idx_a (x)")},
		},
		{
			name:   "a name that breaks the pattern",
			sql:    "CREATE TABLE `a(` (x INT, `b)` INT, INDEX toolong (x), INDEX i2 (`b)`));",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}{{column_list}}$", MaxLength: 4},
			want: []located{
				{Code: v1review.PatternInvalid, Name: "toolong", Table: "a(", Text: "INDEX toolong (x)"},
				mismatch("a(", "i2", "^idx_a(b)$", "INDEX i2 (`b)`)"),
			},
		},
		{
			name:   "a format without tokens",
			sql:    "ALTER TABLE t ADD INDEX idx1 (a), ADD INDEX bad (a);",
			config: v1review.TemplateConfig{Format: "idx"},
			want:   []located{mismatch("t", "bad", "idx", "ALTER TABLE t ADD INDEX idx1 (a), ADD INDEX bad (a)")},
		},
		{name: "empty config requires nothing", sql: "CREATE INDEX Bad ON t (a);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NamingIndex(parse(t, tt.sql), tt.config, tt.schema)
			if diff := cmp.Diff(tt.want, locate(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingIndexPatternInvalid(t *testing.T) {
	sql := "CREATE INDEX idx1 ON t (a);"
	got := NamingIndex(parse(t, sql), v1review.TemplateConfig{Format: "^idx_{{table}}($"}, v1review.Schema{})
	if len(got) != 1 || got[0].Code != v1review.PatternInvalid {
		t.Fatalf("findings = %+v, want one PatternInvalid", got)
	}
	if want := "error parsing regexp: missing closing ): `^idx_t($`"; got[0].Detail != want {
		t.Errorf("Detail = %q, want %q", got[0].Detail, want)
	}
}
