package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/v1review"
)

// anchored is a finding as the column, foreign key and drop tests compare
// it: every field but the range, and the text the range anchors.
type anchored struct {
	Code    v1review.Code
	Name    string
	Table   string
	Pattern string
	Detail  string
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

func anchor(sql string, findings []v1review.Finding) []anchored {
	var got []anchored
	for _, f := range findings {
		got = append(got, anchored{f.Code, f.Name, f.Table, f.Pattern, f.Detail, sql[f.Start:f.End]})
	}
	return got
}

func TestNamingColumn(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^[a-z]+(_[a-z]+)*$"}
	mismatch := func(table, name, text string) anchored {
		return anchored{Code: v1review.NameMismatch, Name: name, Table: table, Text: text}
	}
	tooLong := func(table, name, text string) anchored {
		return anchored{Code: v1review.NameTooLong, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []anchored
	}{
		{name: "create ok", sql: "CREATE TABLE book (id INT, creator_id INT);", config: pattern},
		{
			name:   "create anchors the column definition",
			sql:    "CREATE TABLE book (\n  id INT,\n  creatorId INT NOT NULL\n);",
			config: pattern,
			want:   []anchored{mismatch("book", "creatorId", "creatorId INT NOT NULL")},
		},
		{
			name:   "table qualifier dropped",
			sql:    "CREATE TABLE Shop.book (creatorId INT);",
			config: pattern,
			want:   []anchored{mismatch("book", "creatorId", "creatorId INT")},
		},
		{
			name:   "create or replace temporary",
			sql:    "CREATE OR REPLACE TEMPORARY TABLE book (creatorId INT);",
			config: pattern,
			want:   []anchored{mismatch("book", "creatorId", "creatorId INT")},
		},
		{
			name:   "defined columns of create select judged, selected ones not",
			sql:    "CREATE TABLE book (creatorId INT) AS SELECT 1 AS otherId;",
			config: pattern,
			want:   []anchored{mismatch("book", "creatorId", "creatorId INT")},
		},
		{name: "create as select not judged", sql: "CREATE TABLE book AS SELECT 1 AS otherId;", config: pattern},
		{name: "create like not judged", sql: "CREATE TABLE book LIKE OtherBook;", config: pattern},
		{
			name:   "alter add anchors the statement",
			sql:    "ALTER TABLE book\n  ADD COLUMN contentString VARCHAR(255);",
			config: pattern,
			want:   []anchored{mismatch("book", "contentString", "ALTER TABLE book\n  ADD COLUMN contentString VARCHAR(255)")},
		},
		{
			name:   "alter add without the COLUMN keyword",
			sql:    "ALTER TABLE Shop.book ADD contentString VARCHAR(255) FIRST;",
			config: pattern,
			want:   []anchored{mismatch("book", "contentString", "ALTER TABLE Shop.book ADD contentString VARCHAR(255) FIRST")},
		},
		{
			name:   "alter add parenthesized, in written order",
			sql:    "ALTER TABLE book ADD COLUMN (personAge INT, title TEXT, contentString TEXT, PRIMARY KEY (personAge));",
			config: pattern,
			want: []anchored{
				mismatch("book", "personAge", "ALTER TABLE book ADD COLUMN (personAge INT, title TEXT, contentString TEXT, PRIMARY KEY (personAge))"),
				mismatch("book", "contentString", "ALTER TABLE book ADD COLUMN (personAge INT, title TEXT, contentString TEXT, PRIMARY KEY (personAge))"),
			},
		},
		{
			name:   "change judges the new name",
			sql:    "ALTER TABLE book CHANGE COLUMN OldName creatorId INT;",
			config: pattern,
			want:   []anchored{mismatch("book", "creatorId", "ALTER TABLE book CHANGE COLUMN OldName creatorId INT")},
		},
		{name: "change to a matching name", sql: "ALTER TABLE book CHANGE OldName creator_id INT;", config: pattern},
		{
			name:   "alter table qualifier dropped",
			sql:    "ALTER TABLE Shop.book RENAME COLUMN a TO bC, CHANGE d eF INT;",
			config: pattern,
			want: []anchored{
				mismatch("book", "bC", "ALTER TABLE Shop.book RENAME COLUMN a TO bC, CHANGE d eF INT"),
				mismatch("book", "eF", "ALTER TABLE Shop.book RENAME COLUMN a TO bC, CHANGE d eF INT"),
			},
		},
		{
			name:   "rename column judges the new name",
			sql:    "ALTER TABLE book RENAME COLUMN OldName TO creatorId;",
			config: pattern,
			want:   []anchored{mismatch("book", "creatorId", "ALTER TABLE book RENAME COLUMN OldName TO creatorId")},
		},
		{name: "rename column to an empty name not judged", sql: "ALTER TABLE book RENAME COLUMN a TO ``;", config: pattern},
		{
			name:   "empty name in a definition judged",
			sql:    "CREATE TABLE book (`` INT);",
			config: pattern,
			want:   []anchored{mismatch("book", "", "`` INT")},
		},
		{name: "modify not judged", sql: "ALTER TABLE book MODIFY COLUMN creatorId BIGINT;", config: pattern},
		{name: "drop column not judged", sql: "ALTER TABLE book DROP COLUMN creatorId;", config: pattern},
		{name: "rename table and index not judged", sql: "ALTER TABLE book RENAME TO NewBook, RENAME INDEX a TO BadIndex;", config: pattern},
		{
			name:   "several actions, in written order",
			sql:    "ALTER TABLE book ADD aB INT, DROP COLUMN cD, RENAME COLUMN x TO eF, CHANGE y gH INT;",
			config: pattern,
			want: []anchored{
				mismatch("book", "aB", "ALTER TABLE book ADD aB INT, DROP COLUMN cD, RENAME COLUMN x TO eF, CHANGE y gH INT"),
				mismatch("book", "eF", "ALTER TABLE book ADD aB INT, DROP COLUMN cD, RENAME COLUMN x TO eF, CHANGE y gH INT"),
				mismatch("book", "gH", "ALTER TABLE book ADD aB INT, DROP COLUMN cD, RENAME COLUMN x TO eF, CHANGE y gH INT"),
			},
		},
		{
			name:   "routine body not judged",
			sql:    "CREATE PROCEDURE pr() BEGIN CREATE TABLE book (creatorId INT); ALTER TABLE book ADD otherId INT; END;",
			config: pattern,
		},
		{
			name:   "name judged as parsed, quotes removed",
			sql:    "CREATE TABLE book (`creator``Id` INT, `creator_id` INT);",
			config: pattern,
			want:   []anchored{mismatch("book", "creator`Id", "`creator``Id` INT")},
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (`数据` INT, abcde INT, abcdef INT);",
			config: v1review.NamingConfig{MaxLength: 5},
			want: []anchored{
				tooLong("book", "数据", "`数据` INT"),
				tooLong("book", "abcdef", "abcdef INT"),
			},
		},
		{name: "length within limit", sql: "ALTER TABLE book ADD `数据` INT;", config: v1review.NamingConfig{MaxLength: 6}},
		{
			name:   "mismatch before length",
			sql:    "ALTER TABLE book RENAME COLUMN a TO Toolong;",
			config: v1review.NamingConfig{Pattern: "^[a-z]+$", MaxLength: 4},
			want: []anchored{
				mismatch("book", "Toolong", "ALTER TABLE book RENAME COLUMN a TO Toolong"),
				tooLong("book", "Toolong", "ALTER TABLE book RENAME COLUMN a TO Toolong"),
			},
		},
		{
			name:   "unanchored pattern matches anywhere in the name",
			sql:    "CREATE TABLE book (Creator_id INT, CREATOR INT);",
			config: v1review.NamingConfig{Pattern: "[a-z]+"},
			want:   []anchored{mismatch("book", "CREATOR", "CREATOR INT")},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (creatorId INT);"},
		{
			name:   "findings in statement order",
			sql:    "ALTER TABLE book ADD aB INT;\nCREATE TABLE student (cD INT);",
			config: pattern,
			want: []anchored{
				mismatch("book", "aB", "ALTER TABLE book ADD aB INT"),
				mismatch("student", "cD", "cD INT"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamingColumn(parse(t, tt.sql), tt.config)
			if err != nil {
				t.Fatalf("NamingColumn: %v", err)
			}
			if diff := cmp.Diff(tt.want, anchor(tt.sql, got)); diff != "" {
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
