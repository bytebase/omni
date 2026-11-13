package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

// located is a finding as the column, foreign key and drop tests compare
// it, with the text its range anchors in place of the range.
type located struct {
	Code    v1review.Code
	Name    string
	Table   string
	Pattern string
	Text    string
}

func parse(t *testing.T, sql string) []ast.Node {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return list.Items
}

func locate(sql string, got []v1review.Finding) []located {
	var findings []located
	for _, f := range got {
		findings = append(findings, located{f.Code, f.Name, f.Table, f.Pattern, sql[f.Start:f.End]})
	}
	return findings
}

func TestNamingColumn(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^[a-z]+(_[a-z]+)*$"}
	mismatch := func(table, name, text string) located {
		return located{Code: v1review.NameMismatch, Name: name, Table: table, Text: text}
	}
	tooLong := func(table, name, text string) located {
		return located{Code: v1review.NameTooLong, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []located
	}{
		{name: "create ok", sql: "CREATE TABLE book (id INT, creator_id INT);", config: pattern},
		{
			name:   "create anchors the column definition",
			sql:    "CREATE TABLE book (\n  id INT,\n  creatorId INT NOT NULL,\n  updaterId INT\n);",
			config: pattern,
			want: []located{
				mismatch("book", "creatorId", "creatorId INT NOT NULL"),
				mismatch("book", "updaterId", "updaterId INT"),
			},
		},
		{
			name:   "table printed without its database",
			sql:    "CREATE TABLE Shop.book (badCol INT);",
			config: pattern,
			want:   []located{mismatch("book", "badCol", "badCol INT")},
		},
		{
			name:   "quoted name judged without its quotes",
			sql:    "CREATE TABLE book (`bad col` INT, `a``b` INT);",
			config: pattern,
			want: []located{
				mismatch("book", "bad col", "`bad col` INT"),
				mismatch("book", "a`b", "`a``b` INT"),
			},
		},
		{
			name:   "create as select judges the written columns",
			sql:    "CREATE TABLE book (badCol INT) AS SELECT 1 AS OtherBad;",
			config: pattern,
			want:   []located{mismatch("book", "badCol", "badCol INT")},
		},
		{name: "create like has no columns", sql: "CREATE TABLE BadCopy LIKE OtherBad;", config: pattern},
		{
			name:   "add column anchors the statement",
			sql:    "ALTER TABLE book\n  ADD COLUMN contentString VARCHAR(255);",
			config: pattern,
			want:   []located{mismatch("book", "contentString", "ALTER TABLE book\n  ADD COLUMN contentString VARCHAR(255)")},
		},
		{
			name:   "add column list in written order",
			sql:    "ALTER TABLE book ADD COLUMN (personAge INT, ok_col INT, contentString INT);",
			config: pattern,
			want: []located{
				mismatch("book", "personAge", "ALTER TABLE book ADD COLUMN (personAge INT, ok_col INT, contentString INT)"),
				mismatch("book", "contentString", "ALTER TABLE book ADD COLUMN (personAge INT, ok_col INT, contentString INT)"),
			},
		},
		{
			name:   "rename column judges the new name",
			sql:    "ALTER TABLE book RENAME COLUMN BadOld TO newName;",
			config: pattern,
			want:   []located{mismatch("book", "newName", "ALTER TABLE book RENAME COLUMN BadOld TO newName")},
		},
		{
			name:   "change column judges the new name",
			sql:    "ALTER TABLE db.book CHANGE COLUMN BadOld newName INT;",
			config: pattern,
			want:   []located{mismatch("book", "newName", "ALTER TABLE db.book CHANGE COLUMN BadOld newName INT")},
		},
		{name: "modify column not judged", sql: "ALTER TABLE book MODIFY COLUMN BadName INT;", config: pattern},
		{name: "drop column not judged", sql: "ALTER TABLE book DROP COLUMN BadName;", config: pattern},
		{name: "rename table not judged", sql: "ALTER TABLE book RENAME TO BadName;", config: pattern},
		{
			name:   "several actions in written order",
			sql:    "ALTER TABLE book ADD aB INT, RENAME COLUMN x TO cD, CHANGE y eF INT;",
			config: pattern,
			want: []located{
				mismatch("book", "aB", "ALTER TABLE book ADD aB INT, RENAME COLUMN x TO cD, CHANGE y eF INT"),
				mismatch("book", "cD", "ALTER TABLE book ADD aB INT, RENAME COLUMN x TO cD, CHANGE y eF INT"),
				mismatch("book", "eF", "ALTER TABLE book ADD aB INT, RENAME COLUMN x TO cD, CHANGE y eF INT"),
			},
		},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE TABLE t (BadCol INT); ALTER TABLE t ADD BadTwo INT; END;", config: pattern},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (`数据列` INT, `数据` INT);",
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []located{tooLong("book", "数据列", "`数据列` INT")},
		},
		{name: "length within limit", sql: "CREATE TABLE book (abcdefgh INT);", config: v1review.NamingConfig{MaxLength: 8}},
		{
			name:   "pattern and length both reported, pattern first",
			sql:    "CREATE TABLE book (Toolong5 INT, alsolong INT);",
			config: v1review.NamingConfig{Pattern: "^[a-z]+$", MaxLength: 4},
			want: []located{
				mismatch("book", "Toolong5", "Toolong5 INT"),
				tooLong("book", "Toolong5", "Toolong5 INT"),
				tooLong("book", "alsolong", "alsolong INT"),
			},
		},
		{
			name:   "unanchored pattern matches anywhere",
			sql:    "CREATE TABLE book (BADid INT, BAD INT);",
			config: v1review.NamingConfig{Pattern: "[a-z]"},
			want:   []located{mismatch("book", "BAD", "BAD INT")},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (BadCol INT);"},
		{
			name:   "findings in statement order",
			sql:    "ALTER TABLE b ADD oneBad INT;\nCREATE TABLE a (twoBad INT);",
			config: pattern,
			want: []located{
				mismatch("b", "oneBad", "ALTER TABLE b ADD oneBad INT"),
				mismatch("a", "twoBad", "twoBad INT"),
			},
		},
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
