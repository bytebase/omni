package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingUniqueKey(t *testing.T) {
	format := v1review.TemplateConfig{Format: "^$|^uk_{{table}}_{{column_list}}$"}
	strict := v1review.TemplateConfig{Format: "^uk_{{table}}_{{column_list}}$"}
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
		{name: "create unique index ok", sql: "CREATE UNIQUE INDEX uk_book_id_name ON book (id, name);", config: strict},
		{
			name:   "create unique index anchors the statement",
			sql:    "CREATE UNIQUE INDEX\n  bad\n  ON shop.book (id, (name), (id + 1));",
			config: strict,
			want:   []located{mismatch("book", "bad", "^uk_book_id_name$", "CREATE UNIQUE INDEX\n  bad\n  ON shop.book (id, (name), (id + 1))")},
		},
		{name: "create ordinary, fulltext and spatial index not judged", sql: "CREATE INDEX bad ON t (a); CREATE FULLTEXT INDEX bad ON t (a); CREATE SPATIAL INDEX bad ON t (a);", config: strict},
		{name: "unnamed matches an empty alternative", sql: "CREATE TABLE book (id INT, UNIQUE KEY (id));", config: format},
		{
			name:   "unnamed judged as the empty name",
			sql:    "CREATE TABLE book (id INT, UNIQUE (id));",
			config: strict,
			want:   []located{mismatch("book", "", "^uk_book_id$", "UNIQUE (id)")},
		},
		{
			name:   "the index name wins over the constraint symbol",
			sql:    "CREATE TABLE book (\n  id INT,\n  CONSTRAINT uk_book_id UNIQUE KEY bad (id),\n  CONSTRAINT bad2\n    UNIQUE (id),\n  CONSTRAINT bad3 UNIQUE INDEX uk_book_id (id)\n);",
			config: strict,
			want: []located{
				mismatch("book", "bad", "^uk_book_id$", "CONSTRAINT uk_book_id UNIQUE KEY bad (id)"),
				mismatch("book", "bad2", "^uk_book_id$", "CONSTRAINT bad2\n    UNIQUE (id)"),
			},
		},
		{
			name:   "a parenthesized column counts only when no part is a plain column",
			sql:    "CREATE TABLE t (a INT, b INT, UNIQUE u1 ((a)), UNIQUE u2 (a, (b)));",
			config: strict,
			want: []located{
				mismatch("t", "u1", "^uk_t_a$", "UNIQUE u1 ((a))"),
				mismatch("t", "u2", "^uk_t_a$", "UNIQUE u2 (a, (b))"),
			},
		},
		{
			name:   "other constraints and a column's unique attribute not judged",
			sql:    "CREATE TABLE t (a INT UNIQUE KEY, PRIMARY KEY (a), INDEX bad (a), FULLTEXT bad2 (a), CONSTRAINT bad3 FOREIGN KEY bad4 (a) REFERENCES p (id));",
			config: strict,
		},
		{
			name:   "alter add anchors the statement",
			sql:    "ALTER TABLE shop.book\n  ADD UNIQUE uk_book_id (id),\n  ADD CONSTRAINT bad UNIQUE KEY (id, name),\n  ADD INDEX bad2 (id);",
			config: strict,
			want: []located{mismatch("book", "bad", "^uk_book_id_name$",
				"ALTER TABLE shop.book\n  ADD UNIQUE uk_book_id (id),\n  ADD CONSTRAINT bad UNIQUE KEY (id, name),\n  ADD INDEX bad2 (id)")},
		},
		{name: "alter add column list unique key not judged", sql: "ALTER TABLE t ADD COLUMN (z INT, UNIQUE bad (z));", config: strict},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE UNIQUE INDEX bad ON t (a); END;", config: strict},
		{
			name:   "rename takes the columns from the schema",
			sql:    "ALTER TABLE book RENAME INDEX old_uk TO uk_book_id_name, RENAME KEY OLD_UK TO bad;",
			config: strict,
			schema: keySchema(false),
			want: []located{mismatch("book", "bad", "^uk_book_id_name$",
				"ALTER TABLE book RENAME INDEX old_uk TO uk_book_id_name, RENAME KEY OLD_UK TO bad")},
		},
		{
			name:   "rename of the primary key judged",
			sql:    "ALTER TABLE book RENAME INDEX `PRIMARY` TO bad;",
			config: strict,
			schema: keySchema(false),
			want:   []located{mismatch("book", "bad", "^uk_book_id$", "ALTER TABLE book RENAME INDEX `PRIMARY` TO bad")},
		},
		{name: "rename of an ordinary index not judged", sql: "ALTER TABLE book RENAME INDEX old_idx TO bad;", config: strict, schema: keySchema(false)},
		{name: "rename of a missing index not judged", sql: "ALTER TABLE book RENAME INDEX nope TO bad;", config: strict, schema: keySchema(false)},
		{name: "rename on a partition's name not judged", sql: "ALTER TABLE p0 RENAME INDEX old_uk TO bad;", config: strict, schema: keySchema(false)},
		{name: "rename without a schema not judged", sql: "ALTER TABLE book RENAME INDEX old_uk TO bad;", config: strict},
		{name: "rename misses a table whose case differs where that matters", sql: "ALTER TABLE book RENAME INDEX old_uk TO bad;", config: strict, schema: keySchema(true)},
		{
			name:   "rename finds an index whose case differs even so",
			sql:    "ALTER TABLE elsewhere.Book RENAME INDEX OLD_UK TO bad;",
			config: strict,
			schema: keySchema(true),
			want:   []located{mismatch("Book", "bad", "^uk_Book_id_name$", "ALTER TABLE elsewhere.Book RENAME INDEX OLD_UK TO bad")},
		},
		{
			name:   "length counts bytes",
			sql:    "ALTER TABLE t ADD UNIQUE `唯一键` (a), ADD UNIQUE `唯一` (a);",
			config: v1review.TemplateConfig{MaxLength: 8},
			want:   []located{{Code: v1review.NameTooLong, Name: "唯一键", Table: "t", Text: "ALTER TABLE t ADD UNIQUE `唯一键` (a), ADD UNIQUE `唯一` (a)"}},
		},
		{
			name:   "a name that breaks the pattern",
			sql:    "CREATE TABLE `a(` (x INT, UNIQUE toolong (x));",
			config: v1review.TemplateConfig{Format: "^uk_{{table}}$", MaxLength: 4},
			want:   []located{{Code: v1review.PatternInvalid, Name: "toolong", Table: "a(", Text: "UNIQUE toolong (x)"}},
		},
		{name: "empty config requires nothing", sql: "CREATE UNIQUE INDEX Bad ON t (a);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NamingUniqueKey(parse(t, tt.sql), tt.config, tt.schema)
			if diff := cmp.Diff(tt.want, locate(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
