package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

// indexSchema is the pre-change schema the index naming tests rename in.
func indexSchema(objectCaseSensitive bool) v1review.Schema {
	return v1review.Schema{
		ObjectCaseSensitive: objectCaseSensitive,
		Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
			{
				Tables: []*metadata.TableMetadata{
					{
						Name: "Book",
						Indexes: []*metadata.IndexMetadata{
							{Name: "PRIMARY", Expressions: []string{"id"}, Unique: true, Primary: true},
							{Name: "old_uk", Expressions: []string{"isbn"}, Unique: true},
							{Name: "Old_Idx", Expressions: []string{"author_id", "title"}},
							{Name: "old_expr", Expressions: []string{"(lower(`title`))"}},
						},
						Partitions: []*metadata.TablePartitionMetadata{{Name: "p0"}},
					},
				},
			},
			{Name: "lib", Tables: []*metadata.TableMetadata{
				{Name: "shelf", Indexes: []*metadata.IndexMetadata{{Name: "old_row", Expressions: []string{"row_no"}}}},
			}},
		}},
	}
}

func TestNamingIndex(t *testing.T) {
	template := v1review.TemplateConfig{Format: "^$|^idx_{{table}}_{{column_list}}$"}
	never := v1review.TemplateConfig{Format: "^never$"}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		schema v1review.Schema
		want   []judged
	}{
		{name: "create index ok", sql: "CREATE INDEX idx_book_id_name ON book (id, name);", config: template},
		{
			name:   "create index anchors the statement",
			sql:    "CREATE INDEX\n  book_id_name\n  ON shop.book (id, name);",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "book_id_name", Table: "book", Pattern: "^$|^idx_book_id_name$",
				Text: "CREATE INDEX\n  book_id_name\n  ON shop.book (id, name)",
			}},
		},
		{
			name:   "create index takes bare and parenthesized columns, not other expressions",
			sql:    "CREATE INDEX x ON book (a, (b), (a + 1), ((c)), (book.d), e(3) DESC);",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^$|^idx_book_a_b_d_e$",
				Text: "CREATE INDEX x ON book (a, (b), (a + 1), ((c)), (book.d), e(3) DESC)",
			}},
		},
		{
			name:   "create vector index judged",
			sql:    "CREATE VECTOR INDEX x ON book (v);",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^$|^idx_book_v$",
				Text: "CREATE VECTOR INDEX x ON book (v)",
			}},
		},
		{name: "create unique index not judged", sql: "CREATE UNIQUE INDEX Bad ON book (id);", config: never},
		{name: "create fulltext index not judged", sql: "CREATE FULLTEXT INDEX Bad ON book (title);", config: never},
		{name: "create spatial index not judged", sql: "CREATE SPATIAL INDEX Bad ON book (g);", config: never},
		{
			name:   "create table anchors the definition",
			sql:    "CREATE TABLE shop.book (\n  id INT,\n  name VARCHAR(20),\n  INDEX book_name (name)\n);",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "book_name", Table: "book", Pattern: "^$|^idx_book_name$",
				Text: "INDEX book_name (name)",
			}},
		},
		{
			name:   "unnamed index judged as the empty name",
			sql:    "CREATE TABLE book (id INT, KEY (id));",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}_{{column_list}}$"},
			want: []judged{{
				Code: v1review.NameMismatch, Table: "book", Pattern: "^idx_book_id$",
				Text: "KEY (id)",
			}},
		},
		{
			name:   "definition takes parenthesized columns only when it has no bare one",
			sql:    "CREATE TABLE book (a INT, b INT, INDEX x (a(3), (b)), INDEX y ((a), (b + 1), (book.b)), INDEX z (((a))));",
			config: template,
			want: []judged{
				{Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^$|^idx_book_a$", Text: "INDEX x (a(3), (b))"},
				{Code: v1review.NameMismatch, Name: "y", Table: "book", Pattern: "^$|^idx_book_a_b$", Text: "INDEX y ((a), (b + 1), (book.b))"},
				{Code: v1review.NameMismatch, Name: "z", Table: "book", Pattern: "^$|^idx_book_$", Text: "INDEX z (((a)))"},
			},
		},
		{
			name:   "other definitions and column attributes not judged",
			sql:    "CREATE TABLE book (a INT KEY, b INT UNIQUE KEY, PRIMARY KEY (a), UNIQUE KEY Bad (a), FULLTEXT Bad2 (a), SPATIAL INDEX Bad3 (a), VECTOR INDEX Bad4 (a), CONSTRAINT Bad5 FOREIGN KEY (a) REFERENCES p (id), CONSTRAINT Bad6 CHECK (a > 0));",
			config: never,
		},
		{
			name:   "alter anchors the statement",
			sql:    "ALTER TABLE shop.book\n  ADD INDEX book_name (name);",
			config: template,
			want: []judged{{
				Code: v1review.NameMismatch, Name: "book_name", Table: "book", Pattern: "^$|^idx_book_name$",
				Text: "ALTER TABLE shop.book\n  ADD INDEX book_name (name)",
			}},
		},
		{
			name:   "alter adds several",
			sql:    "ALTER TABLE book ADD KEY x (a), ADD UNIQUE Bad (a), ADD INDEX idx_book_b (b), ADD FULLTEXT INDEX Bad2 (c), ADD INDEX (c);",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}_{{column_list}}$"},
			want: []judged{
				{
					Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^idx_book_a$",
					Text: "ALTER TABLE book ADD KEY x (a), ADD UNIQUE Bad (a), ADD INDEX idx_book_b (b), ADD FULLTEXT INDEX Bad2 (c), ADD INDEX (c)",
				},
				{
					Code: v1review.NameMismatch, Table: "book", Pattern: "^idx_book_c$",
					Text: "ALTER TABLE book ADD KEY x (a), ADD UNIQUE Bad (a), ADD INDEX idx_book_b (b), ADD FULLTEXT INDEX Bad2 (c), ADD INDEX (c)",
				},
			},
		},
		{name: "drop index not judged", sql: "ALTER TABLE book DROP INDEX Bad;", config: never},
		{
			name:   "rename judges the new name with the schema's columns",
			sql:    "ALTER TABLE Book RENAME INDEX Old_Idx TO idx_book;",
			config: template,
			schema: indexSchema(true),
			want: []judged{{
				Code: v1review.NameMismatch, Name: "idx_book", Table: "Book", Pattern: "^$|^idx_Book_author_id_title$",
				Text: "ALTER TABLE Book RENAME INDEX Old_Idx TO idx_book",
			}},
		},
		{name: "rename ok", sql: "ALTER TABLE Book RENAME KEY Old_Idx TO idx_Book_author_id_title;", config: template, schema: indexSchema(true)},
		{
			name:   "renamed index found whatever its case",
			sql:    "ALTER TABLE Book RENAME INDEX OLD_IDX TO x;",
			config: template,
			schema: indexSchema(true),
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "Book", Pattern: "^$|^idx_Book_author_id_title$",
				Text: "ALTER TABLE Book RENAME INDEX OLD_IDX TO x",
			}},
		},
		{name: "table of another case not found when case sensitive", sql: "ALTER TABLE book RENAME INDEX Old_Idx TO Bad;", config: never, schema: indexSchema(true)},
		{
			name:   "table of another case found when not case sensitive, printed as written",
			sql:    "ALTER TABLE book RENAME INDEX Old_Idx TO x;",
			config: template,
			schema: indexSchema(false),
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "book", Pattern: "^$|^idx_book_author_id_title$",
				Text: "ALTER TABLE book RENAME INDEX Old_Idx TO x",
			}},
		},
		{
			name:   "renamed index's expressions go in as recorded",
			sql:    "ALTER TABLE Book RENAME INDEX old_expr TO x;",
			config: v1review.TemplateConfig{Format: "^idx_{{column_list}}$"},
			schema: indexSchema(true),
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "Book", Pattern: "^idx_(lower(`title`))$",
				Text: "ALTER TABLE Book RENAME INDEX old_expr TO x",
			}},
		},
		{name: "renamed unique key not judged", sql: "ALTER TABLE Book RENAME INDEX old_uk TO Bad;", config: never, schema: indexSchema(true)},
		{name: "renamed primary key not judged", sql: "ALTER TABLE Book RENAME INDEX `PRIMARY` TO Bad;", config: never, schema: indexSchema(true)},
		{name: "renamed index missing from the schema not judged", sql: "ALTER TABLE Book RENAME INDEX nope TO Bad;", config: never, schema: indexSchema(true)},
		{name: "renamed index on a missing table not judged", sql: "ALTER TABLE nope RENAME INDEX Old_Idx TO Bad;", config: never, schema: indexSchema(true)},
		{name: "partition name does not find its table's index", sql: "ALTER TABLE p0 RENAME INDEX Old_Idx TO Bad;", config: never, schema: indexSchema(true)},
		{name: "table in a named schema not found", sql: "ALTER TABLE lib.shelf RENAME INDEX old_row TO Bad;", config: never, schema: indexSchema(true)},
		{
			name:   "database qualifier ignored in the lookup",
			sql:    "ALTER TABLE elsewhere.Book RENAME INDEX Old_Idx TO x;",
			config: template,
			schema: indexSchema(true),
			want: []judged{{
				Code: v1review.NameMismatch, Name: "x", Table: "Book", Pattern: "^$|^idx_Book_author_id_title$",
				Text: "ALTER TABLE elsewhere.Book RENAME INDEX Old_Idx TO x",
			}},
		},
		{name: "rename with no schema not judged", sql: "ALTER TABLE Book RENAME INDEX Old_Idx TO Bad;", config: never},
		{
			name:   "index the script created is not in the schema",
			sql:    "CREATE INDEX idx_Book_id ON Book (id);\nALTER TABLE Book RENAME INDEX idx_Book_id TO Bad;",
			config: template,
			schema: indexSchema(true),
		},
		{
			name:   "rename and add in one statement, in command order",
			sql:    "ALTER TABLE Book RENAME INDEX Old_Idx TO x, ADD INDEX y (id), RENAME INDEX old_uk TO Bad;",
			config: template,
			schema: indexSchema(true),
			want: []judged{
				{
					Code: v1review.NameMismatch, Name: "x", Table: "Book", Pattern: "^$|^idx_Book_author_id_title$",
					Text: "ALTER TABLE Book RENAME INDEX Old_Idx TO x, ADD INDEX y (id), RENAME INDEX old_uk TO Bad",
				},
				{
					Code: v1review.NameMismatch, Name: "y", Table: "Book", Pattern: "^$|^idx_Book_id$",
					Text: "ALTER TABLE Book RENAME INDEX Old_Idx TO x, ADD INDEX y (id), RENAME INDEX old_uk TO Bad",
				},
			},
		},
		{
			name:   "routine body not judged",
			sql:    "CREATE PROCEDURE pr() BEGIN CREATE INDEX Bad ON book (a); ALTER TABLE book ADD INDEX Bad2 (a); END;",
			config: never,
		},
		{
			name:   "value not quoted in the pattern",
			sql:    "CREATE TABLE `a.b` (id INT, INDEX idx_axb (id), INDEX idx_ab (id));",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}$"},
			want: []judged{{
				Code: v1review.NameMismatch, Name: "idx_ab", Table: "a.b", Pattern: "^idx_a.b$",
				Text: "INDEX idx_ab (id)",
			}},
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (id INT, INDEX `索引名` (id), INDEX abcdefgh (id));",
			config: v1review.TemplateConfig{MaxLength: 8},
			want: []judged{{
				Code: v1review.NameTooLong, Name: "索引名", Table: "book",
				Text: "INDEX `索引名` (id)",
			}},
		},
		{
			name:   "pattern and length both reported",
			sql:    "CREATE INDEX toolong ON book (a);",
			config: v1review.TemplateConfig{Format: "^idx_", MaxLength: 4},
			want: []judged{
				{Code: v1review.NameMismatch, Name: "toolong", Table: "book", Pattern: "^idx_", Text: "CREATE INDEX toolong ON book (a)"},
				{Code: v1review.NameTooLong, Name: "toolong", Table: "book", Text: "CREATE INDEX toolong ON book (a)"},
			},
		},
		{name: "empty config requires nothing", sql: "CREATE INDEX Bad ON book (a);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judge(t, tt.sql, func(stmts []ast.Node) ([]v1review.Finding, error) {
				return NamingIndex(stmts, tt.config, tt.schema), nil
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// An index whose expanded pattern does not compile is reported once, with
// the compile error, and is not judged for length; the next one is judged
// on its own.
func TestNamingIndexPatternInvalid(t *testing.T) {
	sql := "CREATE TABLE `a(b` (id INT, INDEX toolong (id));\n" +
		"ALTER TABLE Book RENAME INDEX old_expr TO toolong;\n" +
		"CREATE INDEX toolong ON ok (id);"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got := NamingIndex(list.Items, v1review.TemplateConfig{Format: "^idx_{{table}}_{{column_list}}$", MaxLength: 4}, indexSchema(true))
	definition := "INDEX toolong (id)"
	rename := "ALTER TABLE Book RENAME INDEX old_expr TO toolong"
	create := "CREATE INDEX toolong ON ok (id)"
	want := []v1review.Finding{
		{
			Code: v1review.PatternInvalid, Name: "toolong", Table: "a(b",
			Detail: "error parsing regexp: missing closing ): `^idx_a(b_id$`",
			Start:  28, End: 28 + len(definition),
		},
		{Code: v1review.NameMismatch, Name: "toolong", Table: "Book", Pattern: "^idx_Book_(lower(`title`))$", Start: 49, End: 49 + len(rename)},
		{Code: v1review.NameTooLong, Name: "toolong", Table: "Book", Start: 49, End: 49 + len(rename)},
		{Code: v1review.NameMismatch, Name: "toolong", Table: "ok", Pattern: "^idx_ok_id$", Start: 100, End: 100 + len(create)},
		{Code: v1review.NameTooLong, Name: "toolong", Table: "ok", Start: 100, End: 100 + len(create)},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
