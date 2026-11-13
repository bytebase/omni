package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
	"github.com/bytebase/omni/v1review"
)

// located is a finding as the column, foreign key, and drop tests compare
// it, with Text the source its range anchors: empty when it has no range.
type located struct {
	Code    v1review.Code
	Name    string
	Table   string
	Pattern string
	Text    string
}

func locate(sql string, findings []v1review.Finding) []located {
	var out []located
	for _, f := range findings {
		out = append(out, located{f.Code, f.Name, f.Table, f.Pattern, sql[f.Start:f.End]})
	}
	return out
}

func parse(t *testing.T, sql string) []ast.Node {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return list.Items
}

func TestNamingColumn(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^[a-z]+(_[a-z]+)*$"}
	mismatch := func(table, name, text string) located {
		return located{Code: v1review.NameMismatch, Name: name, Table: table, Text: text}
	}
	const alter = "ALTER TABLE shop.book ADD COLUMN (AddA INT, add_b INT), ADD AddC INT, CHANGE COLUMN old_d NewD INT, MODIFY COLUMN ModE INT, RENAME COLUMN old_f TO NewF, DROP COLUMN DropG, ALTER COLUMN AlterH SET DEFAULT 1"
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []located
	}{
		{name: "create ok", sql: "CREATE TABLE book (id INT, creator_id INT);", config: pattern},
		{
			name:   "create anchors each column",
			sql:    "CREATE TABLE Shop.Book (\n  id INT,\n  creatorId INT,\n  UpdatedTs TIMESTAMP\n);",
			config: pattern,
			want: []located{
				mismatch("Book", "creatorId", "creatorId INT"),
				mismatch("Book", "UpdatedTs", "UpdatedTs TIMESTAMP"),
			},
		},
		{
			name:   "column qualifier not judged",
			sql:    "CREATE TABLE book (Book.creator_id INT, Shop.Book.BadCol INT);",
			config: pattern,
			want:   []located{mismatch("book", "BadCol", "Shop.Book.BadCol INT")},
		},
		{
			name:   "quoted name judged unquoted",
			sql:    "CREATE TABLE book (`creator id` INT, `a``b` INT);",
			config: pattern,
			want: []located{
				mismatch("book", "creator id", "`creator id` INT"),
				mismatch("book", "a`b", "`a``b` INT"),
			},
		},
		{
			name:   "create with select judges only defined columns",
			sql:    "CREATE TABLE book (BadCol INT) AS SELECT 1 AS OtherBad;\nCREATE TABLE copy LIKE book;",
			config: pattern,
			want:   []located{mismatch("book", "BadCol", "BadCol INT")},
		},
		{
			name:   "alter judges added, changed and renamed columns, not modified ones",
			sql:    alter + ";",
			config: pattern,
			want: []located{
				mismatch("book", "AddA", alter),
				mismatch("book", "AddC", alter),
				mismatch("book", "NewD", alter),
				mismatch("book", "NewF", alter),
			},
		},
		{name: "alter old names not judged", sql: "ALTER TABLE book CHANGE COLUMN OldName new_name INT, RENAME COLUMN OldOther TO new_other;", config: pattern},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE TABLE book (BadCol INT); END;", config: pattern},
		{
			name:   "wrapped statements not judged",
			sql:    "EXPLAIN ALTER TABLE book ADD COLUMN BadCol INT;\nPLAN REPLAYER DUMP EXPLAIN ALTER TABLE book ADD COLUMN BadCol INT;",
			config: pattern,
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (`数据列` INT, `数据` INT);",
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []located{{Code: v1review.NameTooLong, Name: "数据列", Table: "book", Text: "`数据列` INT"}},
		},
		{
			name:   "pattern and length both reported, per column",
			sql:    "CREATE TABLE book (Toolong5 INT, Other6 INT);",
			config: v1review.NamingConfig{Pattern: "^[a-z]+$", MaxLength: 6},
			want: []located{
				{Code: v1review.NameMismatch, Name: "Toolong5", Table: "book", Text: "Toolong5 INT"},
				{Code: v1review.NameTooLong, Name: "Toolong5", Table: "book", Text: "Toolong5 INT"},
				{Code: v1review.NameMismatch, Name: "Other6", Table: "book", Text: "Other6 INT"},
			},
		},
		{
			name:   "unanchored pattern matches anywhere",
			sql:    "CREATE TABLE book (BADcol INT, BAD INT);",
			config: v1review.NamingConfig{Pattern: "[a-z]+"},
			want:   []located{mismatch("book", "BAD", "BAD INT")},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (BadCol INT);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamingColumn(parse(t, tt.sql), tt.config)
			if err != nil {
				t.Fatalf("NamingColumn: %v", err)
			}
			if diff := cmp.Diff(tt.want, locate(tt.sql, got)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingColumnInvalidPattern(t *testing.T) {
	if _, err := NamingColumn(parse(t, "CREATE TABLE t (id INT);"), v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}

// The parser builds no nil entries, so a hand-built tree pins that the walk
// skips them, and a tree without a table, rather than dereferencing them.
func TestNamingColumnNilEntries(t *testing.T) {
	bad := &ast.ColumnDef{Name: "Bad"}
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
			{Type: ast.ATRenameColumn},
		}},
	}
	got, err := NamingColumn(stmts, v1review.NamingConfig{Pattern: "^[a-z]+$"})
	if err != nil {
		t.Fatalf("NamingColumn: %v", err)
	}
	want := []v1review.Finding{
		{Code: v1review.NameMismatch, Name: "Bad", Table: "book"},
		{Code: v1review.NameMismatch, Name: "Bad", Table: "tech_book"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings mismatch (-want +got):\n%s", diff)
	}
}
