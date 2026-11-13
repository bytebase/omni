package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestNamingUniqueKey(t *testing.T) {
	const format = "^$|^uk_{{table}}_{{column_list}}$"
	template := v1review.TemplateConfig{Format: format}
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
		{
			name:   "create table ok",
			sql:    "CREATE TABLE book (id INT, name TEXT, UNIQUE INDEX uk_book_id_name (id, name), UNIQUE KEY (id), UNIQUE (name));",
			config: template,
		},
		{
			name:   "create table anchors the constraint",
			sql:    "CREATE TABLE book (\n  id INT,\n  UNIQUE KEY uk_id (id) USING BTREE COMMENT 'c'\n);",
			config: template,
			want:   []located{mismatch("book", "uk_id", "^$|^uk_book_id$", "UNIQUE KEY uk_id (id) USING BTREE COMMENT 'c'")},
		},
		{
			name:   "constraint symbol wins over the name after unique",
			sql:    "CREATE TABLE book (id INT, CONSTRAINT Sym UNIQUE KEY uk_book_id (id), CONSTRAINT Only UNIQUE (id), UNIQUE Named (id));",
			config: template,
			want: []located{
				mismatch("book", "Sym", "^$|^uk_book_id$", "CONSTRAINT Sym UNIQUE KEY uk_book_id (id)"),
				mismatch("book", "Only", "^$|^uk_book_id$", "CONSTRAINT Only UNIQUE (id)"),
				mismatch("book", "Named", "^$|^uk_book_id$", "UNIQUE Named (id)"),
			},
		},
		{
			name:   "unnamed key judged as the empty name",
			sql:    "CREATE TABLE book (id INT, UNIQUE (id));",
			config: v1review.TemplateConfig{Format: "^uk_{{table}}$", MaxLength: 1},
			want:   []located{mismatch("book", "", "^uk_book$", "UNIQUE (id)")},
		},
		{
			name:   "qualifier dropped, prefix and order ignored, expression parts give no name",
			sql:    "CREATE TABLE Shop.Book (a INT, b TEXT, UNIQUE Bad (a DESC, b(10), (lower(b))), UNIQUE Expr ((a + 1)));",
			config: template,
			want: []located{
				mismatch("Book", "Bad", "^$|^uk_Book_a_b$", "UNIQUE Bad (a DESC, b(10), (lower(b)))"),
				mismatch("Book", "Expr", "^$|^uk_Book_$", "UNIQUE Expr ((a + 1))"),
			},
		},
		{
			name:   "other constraints and column keys not judged",
			sql:    "CREATE TABLE book (id INT PRIMARY KEY, a INT UNIQUE, b TEXT UNIQUE KEY, KEY Idx (a), FULLTEXT Ft (b), CONSTRAINT Fk FOREIGN KEY Fk2 (a) REFERENCES r (x), CONSTRAINT Chk CHECK (a > 0), PRIMARY KEY (id));",
			config: template,
		},
		{
			name:   "alter add has no range",
			sql:    "ALTER TABLE shop.tech_book ADD UNIQUE Uk (id, name), ADD UNIQUE KEY uk_tech_book_id (id), ADD CONSTRAINT Sym UNIQUE INDEX uk_tech_book_id (id), ADD UNIQUE (id), ADD KEY Idx (id), ADD COLUMN x INT UNIQUE, ADD (UNIQUE Paren (id));",
			config: template,
			want: []located{
				mismatch("tech_book", "Uk", "^$|^uk_tech_book_id_name$", ""),
				mismatch("tech_book", "Sym", "^$|^uk_tech_book_id$", ""),
			},
		},
		{
			name:   "create unique index only",
			sql:    "CREATE UNIQUE INDEX uk_tech_book_id_name ON tech_book (id, name);\nCREATE UNIQUE INDEX Bad ON shop.tech_book (id, name);\nCREATE INDEX Idx ON tech_book (id);\nCREATE FULLTEXT INDEX Ft ON tech_book (name);\nCREATE UNIQUE INDEX Fn ON tech_book ((lower(name)), id);",
			config: template,
			want: []located{
				mismatch("tech_book", "Bad", "^$|^uk_tech_book_id_name$", ""),
				mismatch("tech_book", "Fn", "^$|^uk_tech_book_id$", ""),
			},
		},
		{
			name:   "rename judged against the found index's expressions, a primary key's index too",
			sql:    "ALTER TABLE tech_book RENAME INDEX old_uk TO uk_tech_book_id_name, RENAME KEY OLD_UK TO Bad, RENAME INDEX `primary` TO Pk;",
			config: template,
			schema: keySchema(false),
			want: []located{
				mismatch("tech_book", "Bad", "^$|^uk_tech_book_id_name$", ""),
				mismatch("tech_book", "Pk", "^$|^uk_tech_book_id$", ""),
			},
		},
		{
			name:   "rename of a plain or missing index not judged",
			sql:    "ALTER TABLE tech_book RENAME INDEX old_index TO Bad, RENAME INDEX idx_author TO Bad, RENAME INDEX nope TO Bad;",
			config: template,
			schema: keySchema(false),
		},
		{
			name:   "rename looks the index up across the schema, whatever table and database the statement names",
			sql:    "ALTER TABLE other.Author RENAME INDEX old_uk TO Bad;",
			config: template,
			schema: keySchema(true),
			want:   []located{mismatch("Author", "Bad", "^$|^uk_Author_id_name$", "")},
		},
		{
			name:   "rename of a name two tables hold takes the first in the snapshot's order, here a plain index",
			sql:    "ALTER TABLE review RENAME INDEX idx_shared TO Bad;",
			config: template,
			schema: keySchema(false),
		},
		{
			name:   "rename without a schema not judged",
			sql:    "ALTER TABLE tech_book RENAME INDEX old_uk TO Bad;",
			config: template,
		},
		{
			name:   "length counts bytes, allows the limit itself, and follows the mismatch",
			sql:    "CREATE TABLE book (a INT, UNIQUE `唯一键` (a), UNIQUE `唯一` (a), UNIQUE `唯一ab` (a), UNIQUE uk_book_a (a));",
			config: v1review.TemplateConfig{Format: "^uk_{{table}}_{{column_list}}$", MaxLength: 8},
			want: []located{
				{Code: v1review.NameMismatch, Name: "唯一键", Table: "book", Pattern: "^uk_book_a$", Text: "UNIQUE `唯一键` (a)"},
				{Code: v1review.NameTooLong, Name: "唯一键", Table: "book", Pattern: "^uk_book_a$", Text: "UNIQUE `唯一键` (a)"},
				{Code: v1review.NameMismatch, Name: "唯一", Table: "book", Pattern: "^uk_book_a$", Text: "UNIQUE `唯一` (a)"},
				{Code: v1review.NameMismatch, Name: "唯一ab", Table: "book", Pattern: "^uk_book_a$", Text: "UNIQUE `唯一ab` (a)"},
				{Code: v1review.NameTooLong, Name: "uk_book_a", Table: "book", Pattern: "^uk_book_a$", Text: "UNIQUE uk_book_a (a)"},
			},
		},
		{
			name:   "format without tokens",
			sql:    "ALTER TABLE book ADD UNIQUE other (a), ADD UNIQUE my_uk (a);",
			config: v1review.TemplateConfig{Format: "uk"},
			want:   []located{mismatch("book", "other", "uk", "")},
		},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE UNIQUE INDEX Bad ON book (a); END;", config: template},
		{
			name:   "wrapped statements not judged",
			sql:    "EXPLAIN ALTER TABLE book ADD UNIQUE Bad (a);\nPLAN REPLAYER DUMP EXPLAIN ALTER TABLE book ADD UNIQUE Bad (a);",
			config: template,
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (a INT, UNIQUE Bad (a));"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NamingUniqueKey(parse(t, tt.sql), tt.config, tt.schema)
			if diff := cmp.Diff(tt.want, locate(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// A name that holds a regular expression metacharacter can keep its key's
// pattern from compiling: that key alone gets PatternInvalid, with nothing
// else judged for it.
func TestNamingUniqueKeyPatternInvalid(t *testing.T) {
	sql := "CREATE TABLE `bo(ok` (a INT, UNIQUE ToolongName (a));\nCREATE UNIQUE INDEX Bad ON book (a);\nALTER TABLE `bo(ok` RENAME INDEX old_uk TO Renamed;"
	got := NamingUniqueKey(parse(t, sql), v1review.TemplateConfig{Format: "^uk_{{table}}$", MaxLength: 4}, keySchema(false))
	want := []v1review.Finding{
		{Code: v1review.PatternInvalid, Name: "ToolongName", Table: "bo(ok", Detail: "error parsing regexp: missing closing ): `^uk_bo(ok$`"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "book", Pattern: "^uk_book$"},
		{Code: v1review.PatternInvalid, Name: "Renamed", Table: "bo(ok", Detail: "error parsing regexp: missing closing ): `^uk_bo(ok$`"},
	}
	for i := range got {
		got[i].Start, got[i].End = 0, 0
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}

// The parser builds no nil entries, no ATAddIndex command, and no key
// constraint without Columns, so a hand-built tree pins that the walk skips
// the first, judges the second, and reads the key parts of the third.
func TestNamingUniqueKeyHandBuiltTree(t *testing.T) {
	key := &ast.Constraint{Type: ast.ConstrUnique, Name: "Bad", Columns: []string{"a"}}
	parts := &ast.Constraint{Type: ast.ConstrUnique, Name: "Parts", IndexColumns: []*ast.IndexColumn{
		nil, {Expr: &ast.ColumnRef{Column: "b"}}, {Expr: &ast.FuncCallExpr{Name: "lower"}}, {Expr: &ast.ColumnRef{Column: "c"}},
	}}
	stmts := []ast.Node{
		nil,
		&ast.CreateTableStmt{Constraints: []*ast.Constraint{key}},
		&ast.CreateTableStmt{Table: &ast.TableRef{Name: "book"}, Constraints: []*ast.Constraint{nil, key, parts}},
		&ast.AlterTableStmt{Commands: []*ast.AlterTableCmd{{Type: ast.ATAddConstraint, Constraint: key}, {Type: ast.ATRenameIndex, Name: "old_uk", NewName: "Bad"}}},
		&ast.AlterTableStmt{Table: &ast.TableRef{Name: "tech_book"}, Commands: []*ast.AlterTableCmd{
			nil, {Type: ast.ATAddConstraint}, {Type: ast.ATAddIndex}, {Type: ast.ATAddIndex, Constraint: key}, {Type: ast.ATAddIndex, Constraint: parts},
			{Type: ast.ATAddIndex, Constraint: &ast.Constraint{Type: ast.ConstrIndex, Name: "Idx", Columns: []string{"a"}}},
		}},
		&ast.CreateIndexStmt{Unique: true, IndexName: "Bad", Columns: []*ast.IndexColumn{{Expr: &ast.ColumnRef{Column: "a"}}}},
		&ast.CreateIndexStmt{Unique: true, IndexName: "Bad", Table: &ast.TableRef{Name: "author"}, Columns: []*ast.IndexColumn{nil, {Expr: &ast.ColumnRef{Column: "a"}}}},
	}
	got := NamingUniqueKey(stmts, v1review.TemplateConfig{Format: "^uk_{{table}}_{{column_list}}$"}, keySchema(false))
	want := []v1review.Finding{
		{Code: v1review.NameMismatch, Name: "Bad", Table: "book", Pattern: "^uk_book_a$"},
		{Code: v1review.NameMismatch, Name: "Parts", Table: "book", Pattern: "^uk_book_b_c$"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "tech_book", Pattern: "^uk_tech_book_a$"},
		{Code: v1review.NameMismatch, Name: "Parts", Table: "tech_book", Pattern: "^uk_tech_book_b_c$"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "author", Pattern: "^uk_author_a$"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
