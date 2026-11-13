package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestNamingAutoIncrementColumn(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^id$"}
	mismatch := func(table, name, text string) located {
		return located{Code: v1review.NameMismatch, Name: name, Table: table, Text: text}
	}
	const alter = "ALTER TABLE shop.book ADD COLUMN (AddA INT AUTO_INCREMENT, add_b INT), ADD AddC INT AUTO_INCREMENT, CHANGE COLUMN id NewD INT AUTO_INCREMENT, MODIFY COLUMN ModE INT AUTO_INCREMENT, RENAME COLUMN id TO NewF"
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []located
	}{
		{name: "create ok", sql: "CREATE TABLE book (id INT AUTO_INCREMENT, Other INT);", config: pattern},
		{
			name:   "create anchors the column",
			sql:    "CREATE TABLE Shop.Book (\n  Other INT,\n  Shop.Book.Seq INT AUTO_INCREMENT\n);",
			config: pattern,
			want:   []located{mismatch("Book", "Seq", "Shop.Book.Seq INT AUTO_INCREMENT")},
		},
		{name: "auto_random not judged", sql: "CREATE TABLE book (seq BIGINT PRIMARY KEY AUTO_RANDOM);", config: pattern},
		{
			name:   "alter judges added, changed and modified columns, not renamed ones",
			sql:    alter + ";",
			config: pattern,
			want: []located{
				mismatch("book", "AddA", alter),
				mismatch("book", "AddC", alter),
				mismatch("book", "NewD", alter),
				mismatch("book", "ModE", alter),
			},
		},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE TABLE book (seq INT AUTO_INCREMENT); END;", config: pattern},
		{
			name:   "wrapped statements not judged",
			sql:    "EXPLAIN ALTER TABLE book ADD COLUMN seq INT AUTO_INCREMENT;\nPLAN REPLAYER DUMP EXPLAIN ALTER TABLE book ADD COLUMN seq INT AUTO_INCREMENT;",
			config: pattern,
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (`数据列` INT AUTO_INCREMENT);\nCREATE TABLE other (`数据` INT AUTO_INCREMENT);",
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []located{{Code: v1review.NameTooLong, Name: "数据列", Table: "book", Text: "`数据列` INT AUTO_INCREMENT"}},
		},
		{
			name:   "pattern and length both reported",
			sql:    "CREATE TABLE book (seq_id INT AUTO_INCREMENT);",
			config: v1review.NamingConfig{Pattern: "^id$", MaxLength: 5},
			want: []located{
				{Code: v1review.NameMismatch, Name: "seq_id", Table: "book", Text: "seq_id INT AUTO_INCREMENT"},
				{Code: v1review.NameTooLong, Name: "seq_id", Table: "book", Text: "seq_id INT AUTO_INCREMENT"},
			},
		},
		{name: "unanchored pattern matches anywhere", sql: "CREATE TABLE book (book_id_seq INT AUTO_INCREMENT);", config: v1review.NamingConfig{Pattern: "id"}},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (seq INT AUTO_INCREMENT);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamingAutoIncrementColumn(parse(t, tt.sql), tt.config)
			if err != nil {
				t.Fatalf("NamingAutoIncrementColumn: %v", err)
			}
			if diff := cmp.Diff(tt.want, locate(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingAutoIncrementColumnInvalidPattern(t *testing.T) {
	if _, err := NamingAutoIncrementColumn(parse(t, "CREATE TABLE t (id INT AUTO_INCREMENT);"), v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}

// The parser builds no nil entries, so a hand-built tree pins that the walk
// skips them, and a tree without a table, rather than dereferencing them.
func TestNamingAutoIncrementColumnNilEntries(t *testing.T) {
	bad := &ast.ColumnDef{Name: "Bad", AutoIncrement: true}
	stmts := []ast.Node{
		nil,
		&ast.CreateTableStmt{Columns: []*ast.ColumnDef{bad}},
		&ast.CreateTableStmt{Table: &ast.TableRef{Name: "book"}, Columns: []*ast.ColumnDef{nil, bad}},
		&ast.AlterTableStmt{Commands: []*ast.AlterTableCmd{{Type: ast.ATAddColumn, Column: bad}}},
		&ast.AlterTableStmt{Table: &ast.TableRef{Name: "tech_book"}, Commands: []*ast.AlterTableCmd{
			nil,
			{Type: ast.ATAddColumn},
			{Type: ast.ATAddColumn, Columns: []*ast.ColumnDef{nil, bad}},
			{Type: ast.ATChangeColumn},
			{Type: ast.ATModifyColumn},
		}},
	}
	got, err := NamingAutoIncrementColumn(stmts, v1review.NamingConfig{Pattern: "^id$"})
	if err != nil {
		t.Fatalf("NamingAutoIncrementColumn: %v", err)
	}
	want := []v1review.Finding{
		{Code: v1review.NameMismatch, Name: "Bad", Table: "book"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "tech_book"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
