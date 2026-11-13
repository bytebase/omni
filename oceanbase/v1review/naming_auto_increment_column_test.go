package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

func TestNamingAutoIncrementColumn(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^id$"}
	mismatch := func(table, name, text string) judged {
		return judged{Code: v1review.NameMismatch, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []judged
	}{
		{name: "create ok", sql: "CREATE TABLE book (id INT AUTO_INCREMENT);", config: pattern},
		{
			name:   "create anchors the statement",
			sql:    "CREATE TABLE shop.book (\n  name INT,\n  a INT AUTO_INCREMENT\n);",
			config: pattern,
			want:   []judged{mismatch("book", "a", "CREATE TABLE shop.book (\n  name INT,\n  a INT AUTO_INCREMENT\n)")},
		},
		{name: "other columns not judged", sql: "CREATE TABLE book (a INT, b INT DEFAULT 1);", config: pattern},
		{name: "serial not judged", sql: "CREATE TABLE book (a SERIAL);", config: pattern},
		{
			name:   "add column",
			sql:    "ALTER TABLE book ADD COLUMN a INT AUTO_INCREMENT;",
			config: pattern,
			want:   []judged{mismatch("book", "a", "ALTER TABLE book ADD COLUMN a INT AUTO_INCREMENT")},
		},
		{
			name:   "add column list",
			sql:    "ALTER TABLE book ADD (a INT AUTO_INCREMENT, b INT, id INT AUTO_INCREMENT);",
			config: pattern,
			want:   []judged{mismatch("book", "a", "ALTER TABLE book ADD (a INT AUTO_INCREMENT, b INT, id INT AUTO_INCREMENT)")},
		},
		{
			name:   "modify column",
			sql:    "ALTER TABLE book MODIFY COLUMN a INT AUTO_INCREMENT;",
			config: pattern,
			want:   []judged{mismatch("book", "a", "ALTER TABLE book MODIFY COLUMN a INT AUTO_INCREMENT")},
		},
		{
			name:   "change column judges the new name",
			sql:    "ALTER TABLE book CHANGE COLUMN id a INT AUTO_INCREMENT;",
			config: pattern,
			want:   []judged{mismatch("book", "a", "ALTER TABLE book CHANGE COLUMN id a INT AUTO_INCREMENT")},
		},
		{name: "change column to a matching name", sql: "ALTER TABLE book CHANGE COLUMN a id INT AUTO_INCREMENT;", config: pattern},
		{name: "rename column not judged", sql: "ALTER TABLE book RENAME COLUMN id TO a;", config: pattern},
		{name: "table option not judged", sql: "ALTER TABLE book AUTO_INCREMENT = 5;", config: pattern},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE TABLE book (a INT AUTO_INCREMENT); END;", config: pattern},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (`数据列` INT AUTO_INCREMENT);",
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []judged{{Code: v1review.NameTooLong, Name: "数据列", Table: "book", Text: "CREATE TABLE book (`数据列` INT AUTO_INCREMENT)"}},
		},
		{
			name:   "pattern and length both reported",
			sql:    "ALTER TABLE book MODIFY Toolong5 INT AUTO_INCREMENT;",
			config: v1review.NamingConfig{Pattern: "^id$", MaxLength: 4},
			want: []judged{
				{Code: v1review.NameMismatch, Name: "Toolong5", Table: "book", Text: "ALTER TABLE book MODIFY Toolong5 INT AUTO_INCREMENT"},
				{Code: v1review.NameTooLong, Name: "Toolong5", Table: "book", Text: "ALTER TABLE book MODIFY Toolong5 INT AUTO_INCREMENT"},
			},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (a INT AUTO_INCREMENT);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judge(t, tt.sql, func(stmts []ast.Node) ([]v1review.Finding, error) {
				return NamingAutoIncrementColumn(stmts, tt.config)
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingAutoIncrementColumnInvalidPattern(t *testing.T) {
	list, err := parser.Parse("CREATE TABLE t (id INT AUTO_INCREMENT);")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := NamingAutoIncrementColumn(list.Items, v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
