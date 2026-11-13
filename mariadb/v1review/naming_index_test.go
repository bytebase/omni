package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

// indexSchema is the pre-change schema the index and unique key tests
// rename against: one plain, one unique and one primary index on book, a
// partitioned table, and a table and an index whose names carry upper case.
func indexSchema(objectCaseSensitive bool) v1review.Schema {
	return v1review.Schema{
		ObjectCaseSensitive: objectCaseSensitive,
		Metadata: &metadata.DatabaseSchemaMetadata{
			Name: "shop",
			Schemas: []*metadata.SchemaMetadata{{
				Tables: []*metadata.TableMetadata{
					{
						Name: "book",
						Indexes: []*metadata.IndexMetadata{
							{Name: "PRIMARY", Expressions: []string{"id"}, Unique: true, Primary: true},
							{Name: "old_uk", Expressions: []string{"isbn"}, Unique: true},
							{Name: "old_idx", Expressions: []string{"author", "title"}},
							{Name: "Old_Mixed", Expressions: []string{"price"}},
						},
					},
					{
						Name:       "sale",
						Indexes:    []*metadata.IndexMetadata{{Name: "old_idx", Expressions: []string{"day"}}},
						Partitions: []*metadata.TablePartitionMetadata{{Name: "p0"}},
					},
					{
						Name:    "Shelf",
						Indexes: []*metadata.IndexMetadata{{Name: "old_idx", Expressions: []string{"room"}}},
					},
				},
			}},
		},
	}
}

func TestNamingIndex(t *testing.T) {
	template := v1review.TemplateConfig{Format: "^$|^idx_{{table}}_{{column_list}}$"}
	mismatch := func(table, name, pattern, text string) anchored {
		return anchored{Code: v1review.NameMismatch, Name: name, Table: table, Pattern: pattern, Text: text}
	}
	tooLong := func(table, name, text string) anchored {
		return anchored{Code: v1review.NameTooLong, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		schema v1review.Schema
		want   []anchored
	}{
		{name: "create index ok", sql: "CREATE INDEX idx_book_author_title ON book (author, title);", config: template},
		{
			name:   "create index anchors the statement",
			sql:    "CREATE INDEX\n  book_author\n  ON book (author);",
			config: template,
			want:   []anchored{mismatch("book", "book_author", "^$|^idx_book_author$", "CREATE INDEX\n  book_author\n  ON book (author)")},
		},
		{
			name:   "create table anchors the index definition",
			sql:    "CREATE TABLE book (\n  id INT,\n  author INT,\n  INDEX by_author (author)\n    COMMENT 'x'\n);",
			config: template,
			want:   []anchored{mismatch("book", "by_author", "^$|^idx_book_author$", "INDEX by_author (author)\n    COMMENT 'x'")},
		},
		{
			name:   "unnamed index judged with an empty name",
			sql:    "CREATE TABLE book (author INT, KEY (author));",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}$"},
			want:   []anchored{mismatch("book", "", "^idx_book$", "KEY (author)")},
		},
		{
			name:   "unnamed index passes a pattern that takes an empty name",
			sql:    "CREATE TABLE book (author INT, KEY (author));",
			config: template,
		},
		{
			name:   "alter add anchors the statement",
			sql:    "ALTER TABLE book\n  ADD INDEX by_author (author);",
			config: template,
			want:   []anchored{mismatch("book", "by_author", "^$|^idx_book_author$", "ALTER TABLE book\n  ADD INDEX by_author (author)")},
		},
		{
			name:   "alter add key with if not exists",
			sql:    "ALTER TABLE book ADD KEY IF NOT EXISTS k1 (author);",
			config: template,
			want:   []anchored{mismatch("book", "k1", "^$|^idx_book_author$", "ALTER TABLE book ADD KEY IF NOT EXISTS k1 (author)")},
		},
		{
			name:   "qualifier dropped, key parts joined without prefix or direction",
			sql:    "CREATE INDEX i1 ON Shop.Book (title(10) DESC, author ASC);",
			config: template,
			want:   []anchored{mismatch("Book", "i1", "^$|^idx_Book_title_author$", "CREATE INDEX i1 ON Shop.Book (title(10) DESC, author ASC)")},
		},
		{
			name:   "other keys of a create table not judged",
			sql:    "CREATE TABLE book (id INT, v VECTOR(3) NOT NULL, CONSTRAINT Bad1 PRIMARY KEY (id), UNIQUE Bad2 (id), FULLTEXT Bad3 (id), SPATIAL Bad4 (id), VECTOR Bad5 (v), CONSTRAINT Bad6 CHECK (id > 0));",
			config: template,
		},
		{
			name:   "other keys of an alter table not judged",
			sql:    "ALTER TABLE book ADD UNIQUE Bad1 (id), ADD FULLTEXT Bad2 (id), ADD SPATIAL Bad3 (id), ADD PRIMARY KEY (id), ADD VECTOR Bad4 (v);",
			config: template,
		},
		{name: "create unique index not judged", sql: "CREATE UNIQUE INDEX Bad ON book (id);", config: template},
		{name: "create fulltext index not judged", sql: "CREATE FULLTEXT INDEX Bad ON book (id);", config: template},
		{name: "create spatial index not judged", sql: "CREATE SPATIAL INDEX Bad ON book (id);", config: template},
		{
			name:   "create vector index judged",
			sql:    "CREATE VECTOR INDEX Bad ON book (v);",
			config: template,
			want:   []anchored{mismatch("book", "Bad", "^$|^idx_book_v$", "CREATE VECTOR INDEX Bad ON book (v)")},
		},
		{
			name:   "alter add parenthesized form not judged",
			sql:    "ALTER TABLE book ADD (author INT, INDEX Bad (author));",
			config: template,
		},
		{name: "drop index not judged", sql: "ALTER TABLE book DROP INDEX Bad; DROP INDEX Bad ON book;", config: template},
		{
			name:   "routine body not judged",
			sql:    "CREATE PROCEDURE pr() BEGIN CREATE INDEX Bad ON book (author); ALTER TABLE book RENAME INDEX old_idx TO Bad; END;",
			config: template,
			schema: indexSchema(false),
		},
		{
			name:   "rename takes its columns from the schema",
			sql:    "ALTER TABLE book\n  RENAME INDEX old_idx TO idx_book;",
			config: template,
			schema: indexSchema(false),
			want:   []anchored{mismatch("book", "idx_book", "^$|^idx_book_author_title$", "ALTER TABLE book\n  RENAME INDEX old_idx TO idx_book")},
		},
		{
			name:   "rename ok",
			sql:    "ALTER TABLE book RENAME KEY old_idx TO idx_book_author_title;",
			config: template,
			schema: indexSchema(false),
		},
		{
			name:   "rename of a unique index not judged",
			sql:    "ALTER TABLE book RENAME INDEX old_uk TO Bad, RENAME INDEX `PRIMARY` TO Bad2;",
			config: template,
			schema: indexSchema(false),
		},
		{
			name:   "rename of an index the schema does not hold not judged",
			sql:    "ALTER TABLE book RENAME INDEX missing TO Bad; ALTER TABLE missing RENAME INDEX old_idx TO Bad;",
			config: template,
			schema: indexSchema(false),
		},
		{
			name:   "rename with no schema not judged",
			sql:    "ALTER TABLE book RENAME INDEX old_idx TO Bad;",
			config: template,
		},
		{
			name:   "rename finds the index whatever its case",
			sql:    "ALTER TABLE book RENAME INDEX OLD_IDX TO Bad, RENAME INDEX old_mixed TO Bad2;",
			config: template,
			schema: indexSchema(true),
			want: []anchored{
				mismatch("book", "Bad", "^$|^idx_book_author_title$", "ALTER TABLE book RENAME INDEX OLD_IDX TO Bad, RENAME INDEX old_mixed TO Bad2"),
				mismatch("book", "Bad2", "^$|^idx_book_price$", "ALTER TABLE book RENAME INDEX OLD_IDX TO Bad, RENAME INDEX old_mixed TO Bad2"),
			},
		},
		{
			name:   "rename finds the table by case where object names are case sensitive",
			sql:    "ALTER TABLE BOOK RENAME INDEX old_idx TO Bad; ALTER TABLE shelf RENAME INDEX old_idx TO Bad; ALTER TABLE Shelf RENAME INDEX old_idx TO Bad;",
			config: template,
			schema: indexSchema(true),
			want:   []anchored{mismatch("Shelf", "Bad", "^$|^idx_Shelf_room$", "ALTER TABLE Shelf RENAME INDEX old_idx TO Bad")},
		},
		{
			name:   "rename finds the table whatever its case where they are not, and prints the name as written",
			sql:    "ALTER TABLE BOOK RENAME INDEX old_idx TO Bad; ALTER TABLE shelf RENAME INDEX old_idx TO Bad;",
			config: template,
			schema: indexSchema(false),
			want: []anchored{
				mismatch("BOOK", "Bad", "^$|^idx_BOOK_author_title$", "ALTER TABLE BOOK RENAME INDEX old_idx TO Bad"),
				mismatch("shelf", "Bad", "^$|^idx_shelf_room$", "ALTER TABLE shelf RENAME INDEX old_idx TO Bad"),
			},
		},
		{
			name:   "rename ignores the database qualifier",
			sql:    "ALTER TABLE other.book RENAME INDEX old_idx TO Bad;",
			config: template,
			schema: indexSchema(false),
			want:   []anchored{mismatch("book", "Bad", "^$|^idx_book_author_title$", "ALTER TABLE other.book RENAME INDEX old_idx TO Bad")},
		},
		{
			name:   "rename on a partitioned table, not on its partition's name",
			sql:    "ALTER TABLE sale RENAME INDEX old_idx TO Bad; ALTER TABLE p0 RENAME INDEX old_idx TO Bad;",
			config: template,
			schema: indexSchema(false),
			want:   []anchored{mismatch("sale", "Bad", "^$|^idx_sale_day$", "ALTER TABLE sale RENAME INDEX old_idx TO Bad")},
		},
		{
			name:   "several indexes of one alter, in written order",
			sql:    "ALTER TABLE book ADD INDEX i1 (a), RENAME INDEX old_idx TO i2, ADD UNIQUE Bad (b), ADD KEY i3 (c);",
			config: template,
			schema: indexSchema(false),
			want: []anchored{
				mismatch("book", "i1", "^$|^idx_book_a$", "ALTER TABLE book ADD INDEX i1 (a), RENAME INDEX old_idx TO i2, ADD UNIQUE Bad (b), ADD KEY i3 (c)"),
				mismatch("book", "i2", "^$|^idx_book_author_title$", "ALTER TABLE book ADD INDEX i1 (a), RENAME INDEX old_idx TO i2, ADD UNIQUE Bad (b), ADD KEY i3 (c)"),
				mismatch("book", "i3", "^$|^idx_book_c$", "ALTER TABLE book ADD INDEX i1 (a), RENAME INDEX old_idx TO i2, ADD UNIQUE Bad (b), ADD KEY i3 (c)"),
			},
		},
		{
			name:   "names judged as parsed, quotes removed",
			sql:    "ALTER TABLE `my book` ADD INDEX `idx_my book_a``b` (`a``b`);",
			config: template,
		},
		{
			name:   "a metacharacter in a name changes the pattern",
			sql:    "CREATE INDEX idx_book_a ON `b.k` (a);",
			config: template,
			want:   []anchored{mismatch("b.k", "idx_book_a", "^$|^idx_b.k_a$", "CREATE INDEX idx_book_a ON `b.k` (a)")},
		},
		{name: "a metacharacter in a name matches as one", sql: "CREATE INDEX idx_bXk_a ON `b.k` (a);", config: template},
		{
			name:   "a name that keeps the pattern from compiling",
			sql:    "CREATE INDEX i1 ON `b(k` (a);",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}$", MaxLength: 1},
			want: []anchored{{
				Code: v1review.PatternInvalid, Name: "i1", Table: "b(k",
				Detail: "error parsing regexp: missing closing ): `^idx_b(k$`",
				Text:   "CREATE INDEX i1 ON `b(k` (a)",
			}},
		},
		{
			name:   "a format that never compiles, once per index",
			sql:    "CREATE TABLE book (a INT, INDEX i1 (a), INDEX i2 (a));",
			config: v1review.TemplateConfig{Format: "(idx"},
			want: []anchored{
				{Code: v1review.PatternInvalid, Name: "i1", Table: "book", Detail: "error parsing regexp: missing closing ): `(idx`", Text: "INDEX i1 (a)"},
				{Code: v1review.PatternInvalid, Name: "i2", Table: "book", Detail: "error parsing regexp: missing closing ): `(idx`", Text: "INDEX i2 (a)"},
			},
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (a INT, INDEX `索引` (a), INDEX abcde (a));",
			config: v1review.TemplateConfig{MaxLength: 5},
			want:   []anchored{tooLong("book", "索引", "INDEX `索引` (a)")},
		},
		{
			name:   "mismatch before length",
			sql:    "ALTER TABLE book ADD INDEX toolong (a);",
			config: v1review.TemplateConfig{Format: "^idx_{{column_list}}$", MaxLength: 4},
			want: []anchored{
				mismatch("book", "toolong", "^idx_a$", "ALTER TABLE book ADD INDEX toolong (a)"),
				tooLong("book", "toolong", "ALTER TABLE book ADD INDEX toolong (a)"),
			},
		},
		{
			name:   "unanchored format matches anywhere in the name",
			sql:    "ALTER TABLE book ADD INDEX my_idx_a_x (a), ADD INDEX other (a);",
			config: v1review.TemplateConfig{Format: "idx_{{column_list}}"},
			want:   []anchored{mismatch("book", "other", "idx_a", "ALTER TABLE book ADD INDEX my_idx_a_x (a), ADD INDEX other (a)")},
		},
		{name: "empty config requires nothing", sql: "ALTER TABLE book ADD INDEX Bad (a);"},
		{
			name:   "findings in statement order",
			sql:    "CREATE INDEX i1 ON book (a);\nCREATE TABLE student (a INT, INDEX i2 (a));",
			config: template,
			want: []anchored{
				mismatch("book", "i1", "^$|^idx_book_a$", "CREATE INDEX i1 ON book (a)"),
				mismatch("student", "i2", "^$|^idx_student_a$", "INDEX i2 (a)"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NamingIndex(parse(t, tt.sql), tt.config, tt.schema)
			if diff := cmp.Diff(tt.want, anchor(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
