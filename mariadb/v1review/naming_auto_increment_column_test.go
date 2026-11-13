package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingAutoIncrementColumn(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^id$"}
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
		{name: "create ok", sql: "CREATE TABLE book (id INT AUTO_INCREMENT PRIMARY KEY, Name TEXT);", config: pattern},
		{
			name:   "create anchors the statement",
			sql:    "CREATE TABLE Shop.book (\n  Name TEXT,\n  a INT AUTO_INCREMENT PRIMARY KEY\n);",
			config: pattern,
			want:   []anchored{mismatch("book", "a", "CREATE TABLE Shop.book (\n  Name TEXT,\n  a INT AUTO_INCREMENT PRIMARY KEY\n)")},
		},
		{name: "other columns not judged", sql: "CREATE TABLE book (a INT, b INT DEFAULT 1);", config: pattern},
		{
			name:   "serial default value judged",
			sql:    "CREATE TABLE book (a BIGINT SERIAL DEFAULT VALUE);",
			config: pattern,
			want:   []anchored{mismatch("book", "a", "CREATE TABLE book (a BIGINT SERIAL DEFAULT VALUE)")},
		},
		{name: "serial type not judged", sql: "CREATE TABLE book (a SERIAL);", config: pattern},
		{
			name:   "alter add",
			sql:    "ALTER TABLE Shop.book ADD COLUMN a INT AUTO_INCREMENT;",
			config: pattern,
			want:   []anchored{mismatch("book", "a", "ALTER TABLE Shop.book ADD COLUMN a INT AUTO_INCREMENT")},
		},
		{
			name:   "alter add parenthesized, in written order",
			sql:    "ALTER TABLE book ADD (a INT AUTO_INCREMENT, b INT, c INT AUTO_INCREMENT);",
			config: pattern,
			want: []anchored{
				mismatch("book", "a", "ALTER TABLE book ADD (a INT AUTO_INCREMENT, b INT, c INT AUTO_INCREMENT)"),
				mismatch("book", "c", "ALTER TABLE book ADD (a INT AUTO_INCREMENT, b INT, c INT AUTO_INCREMENT)"),
			},
		},
		{
			name:   "alter modify",
			sql:    "ALTER TABLE book MODIFY COLUMN name INT AUTO_INCREMENT;",
			config: pattern,
			want:   []anchored{mismatch("book", "name", "ALTER TABLE book MODIFY COLUMN name INT AUTO_INCREMENT")},
		},
		{
			name:   "alter change judges the new name",
			sql:    "ALTER TABLE book CHANGE COLUMN id a INT AUTO_INCREMENT;",
			config: pattern,
			want:   []anchored{mismatch("book", "a", "ALTER TABLE book CHANGE COLUMN id a INT AUTO_INCREMENT")},
		},
		{name: "alter change to a matching name", sql: "ALTER TABLE book CHANGE COLUMN a id INT AUTO_INCREMENT;", config: pattern},
		{name: "alter without auto_increment not judged", sql: "ALTER TABLE book ADD a INT, MODIFY b INT, CHANGE c d INT;", config: pattern},
		{name: "rename column not judged", sql: "ALTER TABLE book RENAME COLUMN id TO a;", config: pattern},
		{name: "table auto_increment option not judged", sql: "ALTER TABLE book AUTO_INCREMENT = 100;", config: pattern},
		{
			name:   "routine body not judged",
			sql:    "CREATE PROCEDURE pr() BEGIN CREATE TABLE book (a INT AUTO_INCREMENT PRIMARY KEY); END;",
			config: pattern,
		},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE book (`数据` INT AUTO_INCREMENT PRIMARY KEY);",
			config: v1review.NamingConfig{MaxLength: 5},
			want:   []anchored{tooLong("book", "数据", "CREATE TABLE book (`数据` INT AUTO_INCREMENT PRIMARY KEY)")},
		},
		{name: "length within limit", sql: "CREATE TABLE book (`数据` INT AUTO_INCREMENT PRIMARY KEY);", config: v1review.NamingConfig{MaxLength: 6}},
		{
			name:   "mismatch before length",
			sql:    "ALTER TABLE book MODIFY `my``id` INT AUTO_INCREMENT;",
			config: v1review.NamingConfig{Pattern: "^id$", MaxLength: 4},
			want: []anchored{
				mismatch("book", "my`id", "ALTER TABLE book MODIFY `my``id` INT AUTO_INCREMENT"),
				tooLong("book", "my`id", "ALTER TABLE book MODIFY `my``id` INT AUTO_INCREMENT"),
			},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE book (a INT AUTO_INCREMENT PRIMARY KEY);"},
		{
			name:   "findings in statement order",
			sql:    "ALTER TABLE book ADD a INT AUTO_INCREMENT;\nCREATE TABLE student (b INT AUTO_INCREMENT PRIMARY KEY);",
			config: pattern,
			want: []anchored{
				mismatch("book", "a", "ALTER TABLE book ADD a INT AUTO_INCREMENT"),
				mismatch("student", "b", "CREATE TABLE student (b INT AUTO_INCREMENT PRIMARY KEY)"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamingAutoIncrementColumn(parse(t, tt.sql), tt.config)
			if err != nil {
				t.Fatalf("NamingAutoIncrementColumn: %v", err)
			}
			if diff := cmp.Diff(tt.want, anchor(tt.sql, got)); diff != "" {
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
