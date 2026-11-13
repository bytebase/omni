package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingUniqueKey(t *testing.T) {
	template := v1review.TemplateConfig{Format: "^$|^uk_{{table}}_{{column_list}}$"}
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
		{name: "create unique index ok", sql: "CREATE UNIQUE INDEX uk_book_isbn ON book (isbn);", config: template},
		{
			name:   "create unique index anchors the statement",
			sql:    "CREATE UNIQUE INDEX\n  book_isbn\n  ON book (isbn);",
			config: template,
			want:   []anchored{mismatch("book", "book_isbn", "^$|^uk_book_isbn$", "CREATE UNIQUE INDEX\n  book_isbn\n  ON book (isbn)")},
		},
		{
			name:   "create table anchors the key definition",
			sql:    "CREATE TABLE book (\n  id INT,\n  isbn INT,\n  UNIQUE KEY by_isbn (isbn)\n    COMMENT 'x'\n);",
			config: template,
			want:   []anchored{mismatch("book", "by_isbn", "^$|^uk_book_isbn$", "UNIQUE KEY by_isbn (isbn)\n    COMMENT 'x'")},
		},
		{
			name:   "the name after UNIQUE wins over the constraint symbol",
			sql:    "CREATE TABLE book (isbn INT, CONSTRAINT uk1 UNIQUE KEY uk_book_isbn (isbn), CONSTRAINT uk_book_isbn UNIQUE uk2 (isbn));",
			config: template,
			want:   []anchored{mismatch("book", "uk2", "^$|^uk_book_isbn$", "CONSTRAINT uk_book_isbn UNIQUE uk2 (isbn)")},
		},
		{
			name:   "the constraint symbol names a key written without a name",
			sql:    "CREATE TABLE book (isbn INT, CONSTRAINT uk1 UNIQUE (isbn));",
			config: template,
			want:   []anchored{mismatch("book", "uk1", "^$|^uk_book_isbn$", "CONSTRAINT uk1 UNIQUE (isbn)")},
		},
		{
			name:   "unnamed key judged with an empty name",
			sql:    "CREATE TABLE book (isbn INT, UNIQUE (isbn));",
			config: v1review.TemplateConfig{Format: "^uk_{{table}}$"},
			want:   []anchored{mismatch("book", "", "^uk_book$", "UNIQUE (isbn)")},
		},
		{
			name:   "unnamed key passes a pattern that takes an empty name",
			sql:    "CREATE TABLE book (isbn INT, UNIQUE INDEX (isbn));",
			config: template,
		},
		{
			name:   "alter add anchors the statement",
			sql:    "ALTER TABLE book\n  ADD CONSTRAINT UNIQUE by_isbn (isbn);",
			config: template,
			want:   []anchored{mismatch("book", "by_isbn", "^$|^uk_book_isbn$", "ALTER TABLE book\n  ADD CONSTRAINT UNIQUE by_isbn (isbn)")},
		},
		{
			name:   "alter add unique with if not exists",
			sql:    "ALTER TABLE book ADD UNIQUE KEY IF NOT EXISTS k1 (isbn);",
			config: template,
			want:   []anchored{mismatch("book", "k1", "^$|^uk_book_isbn$", "ALTER TABLE book ADD UNIQUE KEY IF NOT EXISTS k1 (isbn)")},
		},
		{
			name:   "qualifier dropped, key parts joined without prefix or direction",
			sql:    "CREATE UNIQUE INDEX u1 ON Shop.Book (title(10) DESC, author ASC);",
			config: template,
			want:   []anchored{mismatch("Book", "u1", "^$|^uk_Book_title_author$", "CREATE UNIQUE INDEX u1 ON Shop.Book (title(10) DESC, author ASC)")},
		},
		{
			name:   "other keys of a create table not judged",
			sql:    "CREATE TABLE book (id INT UNIQUE, v VECTOR(3) NOT NULL, CONSTRAINT Bad1 PRIMARY KEY (id), INDEX Bad2 (id), FULLTEXT Bad3 (id), SPATIAL Bad4 (id), VECTOR Bad5 (v), CONSTRAINT Bad6 CHECK (id > 0));",
			config: template,
		},
		{
			name:   "other keys of an alter table not judged",
			sql:    "ALTER TABLE book ADD INDEX Bad1 (id), ADD FULLTEXT Bad2 (id), ADD SPATIAL Bad3 (id), ADD PRIMARY KEY (id), ADD VECTOR Bad4 (v), ADD COLUMN c INT UNIQUE;",
			config: template,
		},
		{
			name:   "other create index forms not judged",
			sql:    "CREATE INDEX Bad ON book (id); CREATE FULLTEXT INDEX Bad ON book (id); CREATE SPATIAL INDEX Bad ON book (id); CREATE VECTOR INDEX Bad ON book (v);",
			config: template,
		},
		{
			name:   "alter add parenthesized form not judged",
			sql:    "ALTER TABLE book ADD (isbn INT, UNIQUE Bad (isbn));",
			config: template,
		},
		{
			name:   "routine body not judged",
			sql:    "CREATE PROCEDURE pr() BEGIN CREATE UNIQUE INDEX Bad ON book (isbn); ALTER TABLE book RENAME INDEX old_uk TO Bad; END;",
			config: template,
			schema: indexSchema(false),
		},
		{
			name:   "rename takes its columns from the schema",
			sql:    "ALTER TABLE book\n  RENAME INDEX old_uk TO uk_book;",
			config: template,
			schema: indexSchema(false),
			want:   []anchored{mismatch("book", "uk_book", "^$|^uk_book_isbn$", "ALTER TABLE book\n  RENAME INDEX old_uk TO uk_book")},
		},
		{
			name:   "rename ok",
			sql:    "ALTER TABLE book RENAME KEY old_uk TO uk_book_isbn;",
			config: template,
			schema: indexSchema(false),
		},
		{
			name:   "rename of a primary key's index judged as the unique index it is",
			sql:    "ALTER TABLE book RENAME INDEX `primary` TO Bad;",
			config: template,
			schema: indexSchema(false),
			want:   []anchored{mismatch("book", "Bad", "^$|^uk_book_id$", "ALTER TABLE book RENAME INDEX `primary` TO Bad")},
		},
		{
			name:   "rename of an index that is not unique not judged",
			sql:    "ALTER TABLE book RENAME INDEX old_idx TO Bad;",
			config: template,
			schema: indexSchema(false),
		},
		{
			name:   "rename of an index the schema does not hold not judged",
			sql:    "ALTER TABLE book RENAME INDEX missing TO Bad; ALTER TABLE missing RENAME INDEX old_uk TO Bad;",
			config: template,
			schema: indexSchema(false),
		},
		{
			name:   "rename with no schema not judged",
			sql:    "ALTER TABLE book RENAME INDEX old_uk TO Bad;",
			config: template,
		},
		{
			name:   "rename finds the index whatever its case, the table by case where object names are case sensitive",
			sql:    "ALTER TABLE book RENAME INDEX OLD_UK TO Bad; ALTER TABLE BOOK RENAME INDEX old_uk TO Bad;",
			config: template,
			schema: indexSchema(true),
			want:   []anchored{mismatch("book", "Bad", "^$|^uk_book_isbn$", "ALTER TABLE book RENAME INDEX OLD_UK TO Bad")},
		},
		{
			name:   "rename finds the table whatever its case where they are not, and ignores the database qualifier",
			sql:    "ALTER TABLE other.BOOK RENAME INDEX old_uk TO Bad;",
			config: template,
			schema: indexSchema(false),
			want:   []anchored{mismatch("BOOK", "Bad", "^$|^uk_BOOK_isbn$", "ALTER TABLE other.BOOK RENAME INDEX old_uk TO Bad")},
		},
		{
			name:   "several keys of one alter, in written order",
			sql:    "ALTER TABLE book ADD UNIQUE u1 (a), RENAME INDEX old_uk TO u2, ADD INDEX Bad (b), ADD UNIQUE KEY u3 (c);",
			config: template,
			schema: indexSchema(false),
			want: []anchored{
				mismatch("book", "u1", "^$|^uk_book_a$", "ALTER TABLE book ADD UNIQUE u1 (a), RENAME INDEX old_uk TO u2, ADD INDEX Bad (b), ADD UNIQUE KEY u3 (c)"),
				mismatch("book", "u2", "^$|^uk_book_isbn$", "ALTER TABLE book ADD UNIQUE u1 (a), RENAME INDEX old_uk TO u2, ADD INDEX Bad (b), ADD UNIQUE KEY u3 (c)"),
				mismatch("book", "u3", "^$|^uk_book_c$", "ALTER TABLE book ADD UNIQUE u1 (a), RENAME INDEX old_uk TO u2, ADD INDEX Bad (b), ADD UNIQUE KEY u3 (c)"),
			},
		},
		{
			name:   "a name that keeps the pattern from compiling",
			sql:    "CREATE UNIQUE INDEX u1 ON `b(k` (a);",
			config: v1review.TemplateConfig{Format: "^uk_{{table}}$", MaxLength: 1},
			want: []anchored{{
				Code: v1review.PatternInvalid, Name: "u1", Table: "b(k",
				Detail: "error parsing regexp: missing closing ): `^uk_b(k$`",
				Text:   "CREATE UNIQUE INDEX u1 ON `b(k` (a)",
			}},
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (a INT, UNIQUE `唯一` (a), UNIQUE abcde (a));",
			config: v1review.TemplateConfig{MaxLength: 5},
			want:   []anchored{tooLong("book", "唯一", "UNIQUE `唯一` (a)")},
		},
		{
			name:   "mismatch before length",
			sql:    "ALTER TABLE book ADD UNIQUE toolong (a);",
			config: v1review.TemplateConfig{Format: "^uk_{{column_list}}$", MaxLength: 4},
			want: []anchored{
				mismatch("book", "toolong", "^uk_a$", "ALTER TABLE book ADD UNIQUE toolong (a)"),
				tooLong("book", "toolong", "ALTER TABLE book ADD UNIQUE toolong (a)"),
			},
		},
		{
			name:   "unanchored format matches anywhere in the name",
			sql:    "ALTER TABLE book ADD UNIQUE my_uk_a_x (a), ADD UNIQUE other (a);",
			config: v1review.TemplateConfig{Format: "uk_{{column_list}}"},
			want:   []anchored{mismatch("book", "other", "uk_a", "ALTER TABLE book ADD UNIQUE my_uk_a_x (a), ADD UNIQUE other (a)")},
		},
		{name: "empty config requires nothing", sql: "ALTER TABLE book ADD UNIQUE Bad (a);"},
		{
			name:   "findings in statement order",
			sql:    "CREATE UNIQUE INDEX u1 ON book (a);\nCREATE TABLE student (a INT, UNIQUE u2 (a));",
			config: template,
			want: []anchored{
				mismatch("book", "u1", "^$|^uk_book_a$", "CREATE UNIQUE INDEX u1 ON book (a)"),
				mismatch("student", "u2", "^$|^uk_student_a$", "UNIQUE u2 (a)"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NamingUniqueKey(parse(t, tt.sql), tt.config, tt.schema)
			if diff := cmp.Diff(tt.want, anchor(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
