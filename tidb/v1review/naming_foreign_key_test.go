package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestNamingForeignKey(t *testing.T) {
	const format = "^$|^fk_{{referencing_table}}_{{referencing_column}}_{{referenced_table}}_{{referenced_column}}$"
	template := v1review.TemplateConfig{Format: format}
	mismatch := func(table, name, pattern, text string) located {
		return located{Code: v1review.NameMismatch, Name: name, Table: table, Pattern: pattern, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.TemplateConfig
		want   []located
	}{
		{
			name:   "create ok",
			sql:    "CREATE TABLE book (id INT, author_id INT, FOREIGN KEY fk_book_author_id_author_id (author_id) REFERENCES author (id));",
			config: template,
		},
		{
			name:   "create anchors the constraint",
			sql:    "CREATE TABLE book (\n  id INT,\n  author_id INT,\n  FOREIGN KEY fk_book_author_id (author_id) REFERENCES author (id) ON DELETE CASCADE\n);",
			config: template,
			want: []located{mismatch("book", "fk_book_author_id", "^$|^fk_book_author_id_author_id$",
				"FOREIGN KEY fk_book_author_id (author_id) REFERENCES author (id) ON DELETE CASCADE")},
		},
		{
			name:   "constraint symbol wins over the name after foreign key",
			sql:    "CREATE TABLE book (author_id INT, CONSTRAINT fk1 FOREIGN KEY fk_book_author_id_author_id (author_id) REFERENCES author (id));",
			config: template,
			want: []located{mismatch("book", "fk1", "^$|^fk_book_author_id_author_id$",
				"CONSTRAINT fk1 FOREIGN KEY fk_book_author_id_author_id (author_id) REFERENCES author (id)")},
		},
		{
			name:   "unnamed key judged as the empty name",
			sql:    "CREATE TABLE book (author_id INT, FOREIGN KEY (author_id) REFERENCES author (id));",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$", MaxLength: 1},
			want:   []located{mismatch("book", "", "^fk_book$", "FOREIGN KEY (author_id) REFERENCES author (id)")},
		},
		{
			name:   "qualifiers dropped and columns joined",
			sql:    "CREATE TABLE Shop.Book (a INT, b INT, CONSTRAINT Bad FOREIGN KEY (a, b) REFERENCES Lib.Author (x, y));",
			config: template,
			want: []located{mismatch("Book", "Bad", "^$|^fk_Book_a_b_Author_x_y$",
				"CONSTRAINT Bad FOREIGN KEY (a, b) REFERENCES Lib.Author (x, y)")},
		},
		{
			name:   "column references and other constraints not judged",
			sql:    "CREATE TABLE book (author_id INT REFERENCES author (id), CONSTRAINT Bad UNIQUE (author_id), KEY Worse (author_id), CONSTRAINT Chk CHECK (author_id > 0));",
			config: template,
		},
		{
			name:   "alter has no range",
			sql:    "ALTER TABLE shop.tech_book ADD CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES author (id), ADD FOREIGN KEY Other (a) REFERENCES r (b), ADD CONSTRAINT Uk UNIQUE (a), ADD KEY Idx (a);",
			config: template,
			want: []located{
				mismatch("tech_book", "fk_author_id", "^$|^fk_tech_book_author_id_author_id$", ""),
				mismatch("tech_book", "Other", "^$|^fk_tech_book_a_r_b$", ""),
			},
		},
		{
			name:   "length counts bytes and follows the mismatch",
			sql:    "CREATE TABLE book (a INT, CONSTRAINT `外键约` FOREIGN KEY (a) REFERENCES r (b), CONSTRAINT `外键` FOREIGN KEY (a) REFERENCES r (b), CONSTRAINT fk_book_a_r_b FOREIGN KEY (a) REFERENCES r (b));",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}_{{referencing_column}}_{{referenced_table}}_{{referenced_column}}$", MaxLength: 8},
			want: []located{
				{Code: v1review.NameMismatch, Name: "外键约", Table: "book", Pattern: "^fk_book_a_r_b$", Text: "CONSTRAINT `外键约` FOREIGN KEY (a) REFERENCES r (b)"},
				{Code: v1review.NameTooLong, Name: "外键约", Table: "book", Pattern: "^fk_book_a_r_b$", Text: "CONSTRAINT `外键约` FOREIGN KEY (a) REFERENCES r (b)"},
				{Code: v1review.NameMismatch, Name: "外键", Table: "book", Pattern: "^fk_book_a_r_b$", Text: "CONSTRAINT `外键` FOREIGN KEY (a) REFERENCES r (b)"},
				{Code: v1review.NameTooLong, Name: "fk_book_a_r_b", Table: "book", Pattern: "^fk_book_a_r_b$", Text: "CONSTRAINT fk_book_a_r_b FOREIGN KEY (a) REFERENCES r (b)"},
			},
		},
		{
			name:   "names go into the pattern unquoted",
			sql:    "CREATE TABLE `a.c` (x INT, CONSTRAINT fk_abc FOREIGN KEY (x) REFERENCES r (b), CONSTRAINT `fk_a.c` FOREIGN KEY (x) REFERENCES r (b));",
			config: v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$"},
		},
		{
			name:   "format without tokens",
			sql:    "ALTER TABLE book ADD CONSTRAINT other FOREIGN KEY (a) REFERENCES r (b), ADD CONSTRAINT my_fk FOREIGN KEY (a) REFERENCES r (b);",
			config: v1review.TemplateConfig{Format: "fk"},
			want:   []located{mismatch("book", "other", "fk", "")},
		},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN ALTER TABLE book ADD CONSTRAINT Bad FOREIGN KEY (a) REFERENCES r (b); END;", config: template},
		{
			name:   "wrapped statements not judged",
			sql:    "EXPLAIN ALTER TABLE book ADD CONSTRAINT Bad FOREIGN KEY (a) REFERENCES r (b);\nPLAN REPLAYER DUMP EXPLAIN ALTER TABLE book ADD CONSTRAINT Bad FOREIGN KEY (a) REFERENCES r (b);",
			config: template,
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (a INT, CONSTRAINT Bad FOREIGN KEY (a) REFERENCES r (b));"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamingForeignKey(parse(t, tt.sql), tt.config)
			if err != nil {
				t.Fatalf("NamingForeignKey: %v", err)
			}
			if diff := cmp.Diff(tt.want, locate(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// A name that holds a regular expression metacharacter can keep its key's
// pattern from compiling: that key alone gets PatternInvalid, with nothing
// else judged for it.
func TestNamingForeignKeyPatternInvalid(t *testing.T) {
	sql := "CREATE TABLE `bo(ok` (a INT, CONSTRAINT ToolongName FOREIGN KEY (a) REFERENCES r (b));\nCREATE TABLE book (a INT, CONSTRAINT Bad FOREIGN KEY (a) REFERENCES r (b));"
	got, err := NamingForeignKey(parse(t, sql), v1review.TemplateConfig{Format: "^fk_{{referencing_table}}$", MaxLength: 4})
	if err != nil {
		t.Fatalf("NamingForeignKey: %v", err)
	}
	want := []v1review.Finding{
		{Code: v1review.PatternInvalid, Name: "ToolongName", Table: "bo(ok", Detail: "error parsing regexp: missing closing ): `^fk_bo(ok$`"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "book", Pattern: "^fk_book$"},
	}
	for i := range got {
		got[i].Start, got[i].End = 0, 0
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}

// The parser builds no nil entries, so a hand-built tree pins that the walk
// skips them, and a tree without a table, rather than dereferencing them.
func TestNamingForeignKeyNilEntries(t *testing.T) {
	key := &ast.Constraint{Type: ast.ConstrForeignKey, Name: "Bad", Columns: []string{"a"}, RefColumns: []string{"b"}}
	stmts := []ast.Node{
		nil,
		&ast.CreateTableStmt{Constraints: []*ast.Constraint{key}},
		&ast.CreateTableStmt{Table: &ast.TableRef{Name: "book"}, Constraints: []*ast.Constraint{nil, key}},
		&ast.AlterTableStmt{Commands: []*ast.AlterTableCmd{{Type: ast.ATAddConstraint, Constraint: key}}},
		&ast.AlterTableStmt{Table: &ast.TableRef{Name: "tech_book"}, Commands: []*ast.AlterTableCmd{nil, {Type: ast.ATAddConstraint}, {Type: ast.ATAddConstraint, Constraint: key}}},
	}
	got, err := NamingForeignKey(stmts, v1review.TemplateConfig{Format: "^fk_{{referencing_table}}_{{referencing_column}}_{{referenced_table}}_{{referenced_column}}$"})
	if err != nil {
		t.Fatalf("NamingForeignKey: %v", err)
	}
	want := []v1review.Finding{
		{Code: v1review.NameMismatch, Name: "Bad", Table: "book", Pattern: "^fk_book_a__b$"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "tech_book", Pattern: "^fk_tech_book_a__b$"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
