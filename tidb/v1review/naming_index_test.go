package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// keySchema is the pre-change schema the index and unique key tests rename
// against. idx_shared is on two tables, a plain index first in the
// snapshot's order and a unique one second.
func keySchema(objectCaseSensitive bool) v1review.Schema {
	index := func(name string, unique, primary bool, columns ...string) *metadata.IndexMetadata {
		return &metadata.IndexMetadata{Name: name, Unique: unique, Primary: primary, Expressions: columns}
	}
	return v1review.Schema{
		ObjectCaseSensitive: objectCaseSensitive,
		Metadata: &metadata.DatabaseSchemaMetadata{Name: "shop", Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{Name: "tech_book", Indexes: []*metadata.IndexMetadata{
					index("PRIMARY", true, true, "id"),
					index("old_uk", true, false, "id", "name"),
					index("Old_Index", false, false, "id", "name"),
					index("idx_expr", false, false, "lower(`name`)"),
				}},
				{Name: "author", Indexes: []*metadata.IndexMetadata{index("idx_author", false, false, "born"), index("idx_shared", false, false, "x")}},
				{Name: "review", Indexes: []*metadata.IndexMetadata{index("idx_shared", true, false, "y")}},
			},
		}}},
	}
}

func TestNamingIndex(t *testing.T) {
	const format = "^$|^idx_{{table}}_{{column_list}}$"
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
			sql:    "CREATE TABLE book (id INT, name TEXT, INDEX idx_book_id_name (id, name), KEY (id));",
			config: template,
		},
		{
			name:   "create table anchors the constraint",
			sql:    "CREATE TABLE book (\n  id INT,\n  KEY idx_id (id) USING BTREE COMMENT 'c'\n);",
			config: template,
			want:   []located{mismatch("book", "idx_id", "^$|^idx_book_id$", "KEY idx_id (id) USING BTREE COMMENT 'c'")},
		},
		{
			name:   "name after the keyword wins over the constraint symbol",
			sql:    "CREATE TABLE book (id INT, CONSTRAINT idx_book_id KEY Named (id), CONSTRAINT Sym INDEX (id));",
			config: template,
			want: []located{
				mismatch("book", "Named", "^$|^idx_book_id$", "CONSTRAINT idx_book_id KEY Named (id)"),
				mismatch("book", "Sym", "^$|^idx_book_id$", "CONSTRAINT Sym INDEX (id)"),
			},
		},
		{
			name:   "unnamed index judged as the empty name",
			sql:    "CREATE TABLE book (id INT, KEY (id));",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}$", MaxLength: 1},
			want:   []located{mismatch("book", "", "^idx_book$", "KEY (id)")},
		},
		{
			name:   "qualifier dropped, prefix and order ignored, expression parts give no name",
			sql:    "CREATE TABLE Shop.Book (a INT, b TEXT, KEY Bad (a DESC, b(10), (lower(b))), KEY Expr ((a + 1)));",
			config: template,
			want: []located{
				mismatch("Book", "Bad", "^$|^idx_Book_a_b$", "KEY Bad (a DESC, b(10), (lower(b)))"),
				mismatch("Book", "Expr", "^$|^idx_Book_$", "KEY Expr ((a + 1))"),
			},
		},
		{
			name:   "other constraints and column keys not judged",
			sql:    "CREATE TABLE book (id INT PRIMARY KEY, a INT UNIQUE, b TEXT, c INT KEY, UNIQUE KEY Uk (a), FULLTEXT Ft (b), CONSTRAINT Fk FOREIGN KEY Fk2 (a) REFERENCES r (x), CONSTRAINT Chk CHECK (a > 0), PRIMARY KEY (id), VECTOR INDEX Vec ((VEC_COSINE_DISTANCE(b))));",
			config: template,
		},
		{
			name:   "alter add has no range",
			sql:    "ALTER TABLE shop.tech_book ADD INDEX Idx (id, name), ADD KEY idx_tech_book_id (id), ADD CONSTRAINT Sym INDEX (id), ADD KEY (id), ADD UNIQUE Uk (id), ADD FULLTEXT INDEX Ft (name), ADD COLUMN x INT, ADD (KEY Paren (id));",
			config: template,
			want: []located{
				mismatch("tech_book", "Idx", "^$|^idx_tech_book_id_name$", ""),
				mismatch("tech_book", "Sym", "^$|^idx_tech_book_id$", ""),
			},
		},
		{
			name:   "create index, every kind but unique",
			sql:    "CREATE INDEX idx_tech_book_id_name ON tech_book (id, name);\nCREATE INDEX Bad ON shop.tech_book (id, name);\nCREATE UNIQUE INDEX Uk ON tech_book (id);\nCREATE FULLTEXT INDEX Ft ON tech_book (name);\nCREATE SPATIAL INDEX Sp ON tech_book (name);\nCREATE VECTOR INDEX Vec ON tech_book ((VEC_COSINE_DISTANCE(name)));\nCREATE INDEX Fn ON tech_book ((lower(name)), id);",
			config: template,
			want: []located{
				mismatch("tech_book", "Bad", "^$|^idx_tech_book_id_name$", ""),
				mismatch("tech_book", "Ft", "^$|^idx_tech_book_name$", ""),
				mismatch("tech_book", "Sp", "^$|^idx_tech_book_name$", ""),
				mismatch("tech_book", "Vec", "^$|^idx_tech_book_$", ""),
				mismatch("tech_book", "Fn", "^$|^idx_tech_book_id$", ""),
			},
		},
		{
			name:   "rename judged against the found index's expressions",
			sql:    "ALTER TABLE tech_book RENAME INDEX old_index TO idx_tech_book_id_name, RENAME KEY OLD_INDEX TO Bad, RENAME INDEX idx_expr TO Expr;",
			config: template,
			schema: keySchema(false),
			want: []located{
				mismatch("tech_book", "Bad", "^$|^idx_tech_book_id_name$", ""),
				mismatch("tech_book", "Expr", "^$|^idx_tech_book_lower(`name`)$", ""),
			},
		},
		{
			name:   "rename of a unique, primary or missing index not judged",
			sql:    "ALTER TABLE tech_book RENAME INDEX old_uk TO Bad, RENAME INDEX `PRIMARY` TO Bad, RENAME INDEX nope TO Bad;",
			config: template,
			schema: keySchema(false),
		},
		{
			name:   "rename looks the index up across the schema, whatever table and database the statement names",
			sql:    "ALTER TABLE other.Tech_Book RENAME INDEX idx_author TO Bad;\nALTER TABLE nope RENAME INDEX old_index TO Bad;",
			config: template,
			schema: keySchema(true),
			want: []located{
				mismatch("Tech_Book", "Bad", "^$|^idx_Tech_Book_born$", ""),
				mismatch("nope", "Bad", "^$|^idx_nope_id_name$", ""),
			},
		},
		{
			name:   "rename of a name two tables hold takes the first in the snapshot's order",
			sql:    "ALTER TABLE review RENAME INDEX idx_shared TO Bad;",
			config: template,
			schema: keySchema(false),
			want:   []located{mismatch("review", "Bad", "^$|^idx_review_x$", "")},
		},
		{
			name:   "rename without a schema not judged",
			sql:    "ALTER TABLE tech_book RENAME INDEX old_index TO Bad;",
			config: template,
		},
		{
			name:   "length counts bytes, allows the limit itself, and follows the mismatch",
			sql:    "CREATE TABLE book (a INT, KEY `索引名` (a), KEY `索引` (a), KEY `索引ab` (a), KEY idx_book_a (a));",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}_{{column_list}}$", MaxLength: 8},
			want: []located{
				{Code: v1review.NameMismatch, Name: "索引名", Table: "book", Pattern: "^idx_book_a$", Text: "KEY `索引名` (a)"},
				{Code: v1review.NameTooLong, Name: "索引名", Table: "book", Pattern: "^idx_book_a$", Text: "KEY `索引名` (a)"},
				{Code: v1review.NameMismatch, Name: "索引", Table: "book", Pattern: "^idx_book_a$", Text: "KEY `索引` (a)"},
				{Code: v1review.NameMismatch, Name: "索引ab", Table: "book", Pattern: "^idx_book_a$", Text: "KEY `索引ab` (a)"},
				{Code: v1review.NameTooLong, Name: "idx_book_a", Table: "book", Pattern: "^idx_book_a$", Text: "KEY idx_book_a (a)"},
			},
		},
		{
			name:   "names go into the pattern unquoted",
			sql:    "CREATE TABLE `a.c` (x INT, KEY idx_abc_x (x), KEY `idx_a.c_x` (x));",
			config: v1review.TemplateConfig{Format: "^idx_{{table}}_{{column_list}}$"},
		},
		{
			name:   "format without tokens",
			sql:    "ALTER TABLE book ADD KEY other (a), ADD KEY my_idx (a);",
			config: v1review.TemplateConfig{Format: "idx"},
			want:   []located{mismatch("book", "other", "idx", "")},
		},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE INDEX Bad ON book (a); END;", config: template},
		{
			name:   "wrapped statements not judged",
			sql:    "EXPLAIN ALTER TABLE book ADD KEY Bad (a);\nPLAN REPLAYER DUMP EXPLAIN ALTER TABLE book ADD KEY Bad (a);",
			config: template,
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (a INT, KEY Bad (a));"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NamingIndex(parse(t, tt.sql), tt.config, tt.schema)
			if diff := cmp.Diff(tt.want, locate(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// A name that holds a regular expression metacharacter can keep its key's
// pattern from compiling: that key alone gets PatternInvalid, with nothing
// else judged for it.
func TestNamingIndexPatternInvalid(t *testing.T) {
	sql := "CREATE TABLE `bo(ok` (a INT, KEY ToolongName (a));\nCREATE INDEX Bad ON book (a);\nALTER TABLE `bo(ok` RENAME INDEX old_index TO Renamed;"
	got := NamingIndex(parse(t, sql), v1review.TemplateConfig{Format: "^idx_{{table}}$", MaxLength: 4}, keySchema(false))
	want := []v1review.Finding{
		{Code: v1review.PatternInvalid, Name: "ToolongName", Table: "bo(ok", Detail: "error parsing regexp: missing closing ): `^idx_bo(ok$`"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "book", Pattern: "^idx_book$"},
		{Code: v1review.PatternInvalid, Name: "Renamed", Table: "bo(ok", Detail: "error parsing regexp: missing closing ): `^idx_bo(ok$`"},
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
func TestNamingIndexHandBuiltTree(t *testing.T) {
	key := &ast.Constraint{Type: ast.ConstrIndex, Name: "Bad", Columns: []string{"a"}}
	parts := &ast.Constraint{Type: ast.ConstrIndex, Name: "Parts", IndexColumns: []*ast.IndexColumn{
		nil, {Expr: &ast.ColumnRef{Column: "b"}}, {Expr: &ast.FuncCallExpr{Name: "lower"}}, {Expr: &ast.ColumnRef{Column: "c"}},
	}}
	stmts := []ast.Node{
		nil,
		&ast.CreateTableStmt{Constraints: []*ast.Constraint{key}},
		&ast.CreateTableStmt{Table: &ast.TableRef{Name: "book"}, Constraints: []*ast.Constraint{nil, key, parts}},
		&ast.AlterTableStmt{Commands: []*ast.AlterTableCmd{{Type: ast.ATAddConstraint, Constraint: key}, {Type: ast.ATRenameIndex, Name: "old_index", NewName: "Bad"}}},
		&ast.AlterTableStmt{Table: &ast.TableRef{Name: "tech_book"}, Commands: []*ast.AlterTableCmd{
			nil, {Type: ast.ATAddConstraint}, {Type: ast.ATAddIndex}, {Type: ast.ATAddIndex, Constraint: key}, {Type: ast.ATAddIndex, Constraint: parts},
			{Type: ast.ATAddIndex, Constraint: &ast.Constraint{Type: ast.ConstrUnique, Name: "Uk", Columns: []string{"a"}}},
		}},
		&ast.CreateIndexStmt{IndexName: "Bad", Columns: []*ast.IndexColumn{{Expr: &ast.ColumnRef{Column: "a"}}}},
		&ast.CreateIndexStmt{IndexName: "Bad", Table: &ast.TableRef{Name: "author"}, Columns: []*ast.IndexColumn{nil, {Expr: &ast.ColumnRef{Column: "a"}}}},
	}
	got := NamingIndex(stmts, v1review.TemplateConfig{Format: "^idx_{{table}}_{{column_list}}$"}, keySchema(false))
	want := []v1review.Finding{
		{Code: v1review.NameMismatch, Name: "Bad", Table: "book", Pattern: "^idx_book_a$"},
		{Code: v1review.NameMismatch, Name: "Parts", Table: "book", Pattern: "^idx_book_b_c$"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "tech_book", Pattern: "^idx_tech_book_a$"},
		{Code: v1review.NameMismatch, Name: "Parts", Table: "tech_book", Pattern: "^idx_tech_book_b_c$"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "author", Pattern: "^idx_author_a$"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
