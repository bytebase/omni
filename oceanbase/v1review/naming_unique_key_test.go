package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

func TestNamingUniqueKey(t *testing.T) {
	template := v1review.TemplateConfig{Format: "^$|^uk_{{table}}_{{column_list}}$"}
	never := v1review.TemplateConfig{Format: "^never$"}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		schema v1review.Schema
		want   []judged
	}{
		{name: "create unique index ok", sql: "CREATE UNIQUE INDEX uk_book_id_name ON book (id, name);", config: template},
		{
			name:   "create unique index anchors the statement",
			sql:    "CREATE UNIQUE INDEX\n  book_id_name\n  ON shop.book (id, (name), (id + 1));",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "book_id_name", Table: "book", Pattern: "^$|^uk_book_id_name$",
				Text: "CREATE UNIQUE INDEX\n  book_id_name\n  ON shop.book (id, (name), (id + 1))",
			}},
		},
		{name: "plain index not judged", sql: "CREATE INDEX Bad ON book (id);", config: never},
		{name: "fulltext, spatial and vector indexes not judged", sql: "CREATE FULLTEXT INDEX Bad ON book (a); CREATE SPATIAL INDEX Bad ON book (a); CREATE VECTOR INDEX Bad ON book (a);", config: never},
		{
			name:   "create table anchors the definition",
			sql:    "CREATE TABLE shop.book (\n  id INT,\n  name VARCHAR(20),\n  UNIQUE KEY book_name (name)\n);",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "book_name", Table: "book", Pattern: "^$|^uk_book_name$",
				Text: "UNIQUE KEY book_name (name)",
			}},
		},
		{
			name:   "name after UNIQUE wins over the constraint symbol",
			sql:    "CREATE TABLE book (id INT, CONSTRAINT uk_book_id UNIQUE KEY x (id), CONSTRAINT y UNIQUE (id), CONSTRAINT UNIQUE z (id));",
			config: template,
			want: []judged{
				{Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^$|^uk_book_id$", Text: "CONSTRAINT uk_book_id UNIQUE KEY x (id)"},
				{Code: v1review.NameMismatch, Name: "y", Table: "book", Pattern: "^$|^uk_book_id$", Text: "CONSTRAINT y UNIQUE (id)"},
				{Code: v1review.NameMismatch, Name: "z", Table: "book", Pattern: "^$|^uk_book_id$", Text: "CONSTRAINT UNIQUE z (id)"},
			},
		},
		{
			name:   "unnamed unique key judged as the empty name",
			sql:    "CREATE TABLE book (id INT, UNIQUE (id));",
			config: v1review.TemplateConfig{Format: "^uk_{{table}}_{{column_list}}$"},
			want: []judged{{
				Code: v1review.NameMismatch, Table: "book", Pattern: "^uk_book_id$",
				Text: "UNIQUE (id)",
			}},
		},
		{
			name:   "definition takes parenthesized columns only when it has no bare one",
			sql:    "CREATE TABLE book (a INT, b INT, UNIQUE x (a(3), (b)), UNIQUE y ((a), (b + 1), (book.b)));",
			config: template,
			want: []judged{
				{Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^$|^uk_book_a$", Text: "UNIQUE x (a(3), (b))"},
				{Code: v1review.NameMismatch, Name: "y", Table: "book", Pattern: "^$|^uk_book_a_b$", Text: "UNIQUE y ((a), (b + 1), (book.b))"},
			},
		},
		{
			name:   "other definitions and column attributes not judged",
			sql:    "CREATE TABLE book (a INT UNIQUE, b INT UNIQUE KEY, PRIMARY KEY (a), INDEX Bad (a), FULLTEXT Bad2 (a), CONSTRAINT Bad3 FOREIGN KEY (a) REFERENCES p (id));",
			config: never,
		},
		{
			name:   "alter anchors the statement, in command order",
			sql:    "ALTER TABLE shop.book\n  ADD UNIQUE x (a), ADD INDEX Bad (a), ADD CONSTRAINT y UNIQUE INDEX (b), ADD UNIQUE KEY uk_book_c (c);",
			config: template,
			want: []judged{
				{
					Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^$|^uk_book_a$",
					Text: "ALTER TABLE shop.book\n  ADD UNIQUE x (a), ADD INDEX Bad (a), ADD CONSTRAINT y UNIQUE INDEX (b), ADD UNIQUE KEY uk_book_c (c)",
				},
				{
					Code: v1review.NameMismatch, Name: "y", Table: "book", Pattern: "^$|^uk_book_b$",
					Text: "ALTER TABLE shop.book\n  ADD UNIQUE x (a), ADD INDEX Bad (a), ADD CONSTRAINT y UNIQUE INDEX (b), ADD UNIQUE KEY uk_book_c (c)",
				},
			},
		},
		{
			name:   "rename judges the new name with the schema's columns",
			sql:    "ALTER TABLE Book RENAME INDEX OLD_UK TO uk_book;",
			config: template,
			schema: indexSchema(true),
			want: []judged{{
				Code: v1review.NameMismatch, Name: "uk_book", Table: "Book", Pattern: "^$|^uk_Book_isbn$",
				Text: "ALTER TABLE Book RENAME INDEX OLD_UK TO uk_book",
			}},
		},
		{name: "rename ok", sql: "ALTER TABLE Book RENAME KEY old_uk TO uk_Book_isbn;", config: template, schema: indexSchema(true)},
		{
			name:   "renamed primary key judged, being unique",
			sql:    "ALTER TABLE Book RENAME INDEX `PRIMARY` TO x;",
			config: template,
			schema: indexSchema(true),
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "Book", Pattern: "^$|^uk_Book_id$",
				Text: "ALTER TABLE Book RENAME INDEX `PRIMARY` TO x",
			}},
		},
		{name: "renamed plain index not judged", sql: "ALTER TABLE Book RENAME INDEX Old_Idx TO Bad;", config: never, schema: indexSchema(true)},
		{name: "renamed index missing from the schema not judged", sql: "ALTER TABLE Book RENAME INDEX nope TO Bad;", config: never, schema: indexSchema(true)},
		{name: "table of another case not found when case sensitive", sql: "ALTER TABLE book RENAME INDEX old_uk TO Bad;", config: never, schema: indexSchema(true)},
		{
			name:   "table of another case found when not case sensitive",
			sql:    "ALTER TABLE book RENAME INDEX old_uk TO x;",
			config: template,
			schema: indexSchema(false),
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^$|^uk_book_isbn$",
				Text: "ALTER TABLE book RENAME INDEX old_uk TO x",
			}},
		},
		{name: "rename with no schema not judged", sql: "ALTER TABLE Book RENAME INDEX old_uk TO Bad;", config: never},
		{
			name:   "routine body not judged",
			sql:    "CREATE PROCEDURE pr() BEGIN CREATE UNIQUE INDEX Bad ON book (a); ALTER TABLE book ADD UNIQUE Bad2 (a); END;",
			config: never,
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (id INT, UNIQUE `唯一键` (id), UNIQUE abcdefgh (id));",
			config: v1review.TemplateConfig{MaxLength: 8},
			want: []judged{{
				Code: v1review.NameTooLong, Name: "唯一键", Table: "book",
				Text: "UNIQUE `唯一键` (id)",
			}},
		},
		{name: "empty config requires nothing", sql: "CREATE UNIQUE INDEX Bad ON book (a);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judge(t, tt.sql, func(stmts []ast.Node) ([]v1review.Finding, error) {
				return NamingUniqueKey(stmts, tt.config, tt.schema), nil
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingUniqueKeyPatternInvalid(t *testing.T) {
	sql := "CREATE TABLE `a(b` (id INT, UNIQUE toolong (id), UNIQUE KEY (id));"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got := NamingUniqueKey(list.Items, v1review.TemplateConfig{Format: "^uk_{{table}}$", MaxLength: 4}, v1review.Schema{})
	want := []v1review.Finding{
		{
			Code: v1review.PatternInvalid, Name: "toolong", Table: "a(b",
			Detail: "error parsing regexp: missing closing ): `^uk_a(b$`",
			Start:  28, End: 28 + len("UNIQUE toolong (id)"),
		},
		{
			Code: v1review.PatternInvalid, Table: "a(b",
			Detail: "error parsing regexp: missing closing ): `^uk_a(b$`",
			Start:  49, End: 49 + len("UNIQUE KEY (id)"),
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
