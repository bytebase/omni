package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

// judged is a finding as the tests compare it: its fields, with the text
// its range anchors in place of the range.
type judged struct {
	Code    v1review.Code
	Name    string
	Table   string
	Pattern string
	Text    string
}

// judge parses sql and returns what rule finds in its statements.
func judge(t *testing.T, sql string, rule func([]ast.Node) ([]v1review.Finding, error)) []judged {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := rule(list.Items)
	if err != nil {
		t.Fatalf("rule: %v", err)
	}
	var findings []judged
	for _, f := range got {
		findings = append(findings, judged{f.Code, f.Name, f.Table, f.Pattern, sql[f.Start:f.End]})
	}
	return findings
}

func TestNamingColumn(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^[a-z]+(_[a-z]+)*$"}
	mismatch := func(table, name, text string) judged {
		return judged{Code: v1review.NameMismatch, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []judged
	}{
		{name: "create ok", sql: "CREATE TABLE book (id INT, creator_id INT);", config: pattern},
		{
			name:   "create anchors the column definition",
			sql:    "CREATE TABLE book (\n  id INT,\n  creatorId INT NOT NULL,\n  UpdatedTs TIMESTAMP\n);",
			config: pattern,
			want: []judged{
				mismatch("book", "creatorId", "creatorId INT NOT NULL"),
				mismatch("book", "UpdatedTs", "UpdatedTs TIMESTAMP"),
			},
		},
		{
			name:   "table printed without its database",
			sql:    "CREATE TABLE shop.Book (Id INT);",
			config: pattern,
			want:   []judged{mismatch("Book", "Id", "Id INT")},
		},
		{
			name:   "quoted name judged without its quotes",
			sql:    "CREATE TABLE book (`creator id` INT, `a``b` INT);",
			config: pattern,
			want: []judged{
				mismatch("book", "creator id", "`creator id` INT"),
				mismatch("book", "a`b", "`a``b` INT"),
			},
		},
		{
			name:   "columns beside a query are judged, the query's are not",
			sql:    "CREATE TABLE book (Id INT) AS SELECT 1 AS Other;",
			config: pattern,
			want:   []judged{mismatch("book", "Id", "Id INT")},
		},
		{name: "create as select not judged", sql: "CREATE TABLE book AS SELECT 1 AS Other;", config: pattern},
		{name: "create like not judged", sql: "CREATE TABLE book LIKE Other;", config: pattern},
		{
			name:   "add column anchors the statement",
			sql:    "ALTER TABLE book\n  ADD COLUMN contentString VARCHAR(255);",
			config: pattern,
			want:   []judged{mismatch("book", "contentString", "ALTER TABLE book\n  ADD COLUMN contentString VARCHAR(255)")},
		},
		{
			name:   "add column list",
			sql:    "ALTER TABLE book ADD (personAge INT, ok_name INT, contentString INT);",
			config: pattern,
			want: []judged{
				mismatch("book", "personAge", "ALTER TABLE book ADD (personAge INT, ok_name INT, contentString INT)"),
				mismatch("book", "contentString", "ALTER TABLE book ADD (personAge INT, ok_name INT, contentString INT)"),
			},
		},
		{
			name:   "rename column judges the new name",
			sql:    "ALTER TABLE book RENAME COLUMN BadOld TO newName;",
			config: pattern,
			want:   []judged{mismatch("book", "newName", "ALTER TABLE book RENAME COLUMN BadOld TO newName")},
		},
		{
			name:   "change column judges the new name",
			sql:    "ALTER TABLE book CHANGE COLUMN BadOld newName INT;",
			config: pattern,
			want:   []judged{mismatch("book", "newName", "ALTER TABLE book CHANGE COLUMN BadOld newName INT")},
		},
		{name: "modify column not judged", sql: "ALTER TABLE book MODIFY COLUMN BadName INT;", config: pattern},
		{name: "drop column not judged", sql: "ALTER TABLE book DROP COLUMN BadName;", config: pattern},
		{name: "rename index not judged", sql: "ALTER TABLE book RENAME INDEX a TO BadName;", config: pattern},
		{
			name:   "commands in order",
			sql:    "ALTER TABLE shop.Book RENAME COLUMN a TO B, ADD C INT, CHANGE d E INT;",
			config: pattern,
			want: []judged{
				mismatch("Book", "B", "ALTER TABLE shop.Book RENAME COLUMN a TO B, ADD C INT, CHANGE d E INT"),
				mismatch("Book", "C", "ALTER TABLE shop.Book RENAME COLUMN a TO B, ADD C INT, CHANGE d E INT"),
				mismatch("Book", "E", "ALTER TABLE shop.Book RENAME COLUMN a TO B, ADD C INT, CHANGE d E INT"),
			},
		},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE TABLE book (Bad INT); END;", config: pattern},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (`数据列` INT, abcdefgh INT);",
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []judged{{Code: v1review.NameTooLong, Name: "数据列", Table: "book", Text: "`数据列` INT"}},
		},
		{
			name:   "pattern and length both reported",
			sql:    "CREATE TABLE book (Toolong5 INT);",
			config: v1review.NamingConfig{Pattern: "^[a-z]+$", MaxLength: 4},
			want: []judged{
				{Code: v1review.NameMismatch, Name: "Toolong5", Table: "book", Text: "Toolong5 INT"},
				{Code: v1review.NameTooLong, Name: "Toolong5", Table: "book", Text: "Toolong5 INT"},
			},
		},
		{
			name:   "unanchored pattern matches anywhere",
			sql:    "CREATE TABLE book (BADid INT, BAD INT);",
			config: v1review.NamingConfig{Pattern: "[a-z]"},
			want:   []judged{{Code: v1review.NameMismatch, Name: "BAD", Table: "book", Text: "BAD INT"}},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (Bad INT);"},
		{
			name:   "findings in statement order",
			sql:    "ALTER TABLE a ADD X INT;\nCREATE TABLE b (Y INT);",
			config: pattern,
			want: []judged{
				mismatch("a", "X", "ALTER TABLE a ADD X INT"),
				mismatch("b", "Y", "Y INT"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judge(t, tt.sql, func(stmts []ast.Node) ([]v1review.Finding, error) { return NamingColumn(stmts, tt.config) })
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingColumnInvalidPattern(t *testing.T) {
	list, err := parser.Parse("CREATE TABLE t (id INT);")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := NamingColumn(list.Items, v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
