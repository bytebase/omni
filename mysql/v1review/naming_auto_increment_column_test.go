package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingAutoIncrementColumn(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^id$"}
	mismatch := func(table, name, text string) located {
		return located{Code: v1review.NameMismatch, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []located
	}{
		{name: "create ok", sql: "CREATE TABLE book (id INT AUTO_INCREMENT);", config: pattern},
		{name: "other columns not judged", sql: "CREATE TABLE book (id INT AUTO_INCREMENT, BadName INT);", config: pattern},
		{
			name:   "create anchors the statement, not the column",
			sql:    "CREATE TABLE Shop.book (\n  name INT,\n  a INT AUTO_INCREMENT\n);",
			config: pattern,
			want:   []located{mismatch("book", "a", "CREATE TABLE Shop.book (\n  name INT,\n  a INT AUTO_INCREMENT\n)")},
		},
		{
			name:   "serial default value counts",
			sql:    "CREATE TABLE book (a BIGINT SERIAL DEFAULT VALUE);",
			config: pattern,
			want:   []located{mismatch("book", "a", "CREATE TABLE book (a BIGINT SERIAL DEFAULT VALUE)")},
		},
		{name: "serial type does not count", sql: "CREATE TABLE book (a SERIAL);", config: pattern},
		{
			name:   "add column",
			sql:    "ALTER TABLE book ADD COLUMN a INT AUTO_INCREMENT;",
			config: pattern,
			want:   []located{mismatch("book", "a", "ALTER TABLE book ADD COLUMN a INT AUTO_INCREMENT")},
		},
		{
			name:   "add column list",
			sql:    "ALTER TABLE book ADD COLUMN (b INT, a INT AUTO_INCREMENT);",
			config: pattern,
			want:   []located{mismatch("book", "a", "ALTER TABLE book ADD COLUMN (b INT, a INT AUTO_INCREMENT)")},
		},
		{
			name:   "modify column",
			sql:    "ALTER TABLE db.book MODIFY COLUMN name INT AUTO_INCREMENT;",
			config: pattern,
			want:   []located{mismatch("book", "name", "ALTER TABLE db.book MODIFY COLUMN name INT AUTO_INCREMENT")},
		},
		{
			name:   "change column judges the new name",
			sql:    "ALTER TABLE book CHANGE COLUMN id a INT AUTO_INCREMENT;",
			config: pattern,
			want:   []located{mismatch("book", "a", "ALTER TABLE book CHANGE COLUMN id a INT AUTO_INCREMENT")},
		},
		{name: "change column to a good name", sql: "ALTER TABLE book CHANGE COLUMN a id INT AUTO_INCREMENT;", config: pattern},
		{name: "rename column not judged", sql: "ALTER TABLE book RENAME COLUMN id TO a;", config: pattern},
		{name: "table auto_increment option not judged", sql: "ALTER TABLE book AUTO_INCREMENT = 5;", config: pattern},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr() BEGIN CREATE TABLE t (a INT AUTO_INCREMENT); END;", config: pattern},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (`数据列` INT AUTO_INCREMENT);",
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []located{{Code: v1review.NameTooLong, Name: "数据列", Table: "book", Text: "CREATE TABLE book (`数据列` INT AUTO_INCREMENT)"}},
		},
		{
			name:   "pattern and length both reported, pattern first",
			sql:    "ALTER TABLE book MODIFY Toolong5 INT AUTO_INCREMENT;",
			config: v1review.NamingConfig{Pattern: "^id$", MaxLength: 4},
			want: []located{
				mismatch("book", "Toolong5", "ALTER TABLE book MODIFY Toolong5 INT AUTO_INCREMENT"),
				{Code: v1review.NameTooLong, Name: "Toolong5", Table: "book", Text: "ALTER TABLE book MODIFY Toolong5 INT AUTO_INCREMENT"},
			},
		},
		{name: "unanchored pattern matches anywhere", sql: "CREATE TABLE book (book_id INT AUTO_INCREMENT);", config: v1review.NamingConfig{Pattern: "id"}},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (a INT AUTO_INCREMENT);"},
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
	if _, err := NamingAutoIncrementColumn(parse(t, "CREATE TABLE t (id INT);"), v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
