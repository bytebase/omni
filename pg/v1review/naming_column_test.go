package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

// columnFinding is a NamingColumn finding as the tests compare it.
type columnFinding struct {
	Code  v1review.Code
	Table string
	Name  string
	Text  string
}

func runNamingColumn(t *testing.T, sql string, cfg v1review.NamingConfig) []columnFinding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := NamingColumn(list.Items, cfg)
	if err != nil {
		t.Fatalf("NamingColumn: %v", err)
	}
	var findings []columnFinding
	for _, f := range got {
		if f.Pattern != "" || f.Detail != "" {
			t.Errorf("finding %+v: Pattern and Detail should be empty", f)
		}
		findings = append(findings, columnFinding{f.Code, f.Table, f.Name, sql[f.Start:f.End]})
	}
	return findings
}

func TestNamingColumn(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^[a-z]+(_[a-z]+)*$"}
	mismatch := func(table, name, text string) columnFinding {
		return columnFinding{v1review.NameMismatch, table, name, text}
	}
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []columnFinding
	}{
		{name: "create ok", sql: "CREATE TABLE book (id int, creator_id int);", config: pattern},
		{
			name:   "create judges every column in order",
			sql:    `CREATE TABLE book ("ID" int, name text, "creatorId" int);`,
			config: pattern,
			want: []columnFinding{
				mismatch("book", "ID", `CREATE TABLE book ("ID" int, name text, "creatorId" int)`),
				mismatch("book", "creatorId", `CREATE TABLE book ("ID" int, name text, "creatorId" int)`),
			},
		},
		{name: "unquoted mixed case folds to lower and passes", sql: "CREATE TABLE book (CreatorId int);", config: v1review.NamingConfig{Pattern: "^[a-z]+$"}},
		{
			name:   "table is named without its schema",
			sql:    `CREATE TABLE "Shop"."Book" ("Id" int);`,
			config: pattern,
			want:   []columnFinding{mismatch("Book", "Id", `CREATE TABLE "Shop"."Book" ("Id" int)`)},
		},
		{
			name:   "typed table and partition columns judged",
			sql:    `CREATE TABLE t OF typ ("A" WITH OPTIONS NOT NULL); CREATE TABLE c PARTITION OF p ("B" NOT NULL) FOR VALUES IN (1);`,
			config: pattern,
			want: []columnFinding{
				mismatch("t", "A", `CREATE TABLE t OF typ ("A" WITH OPTIONS NOT NULL)`),
				mismatch("c", "B", `CREATE TABLE c PARTITION OF p ("B" NOT NULL) FOR VALUES IN (1)`),
			},
		},
		{
			name:   "like clause and constraints are not columns",
			sql:    `CREATE TABLE t (LIKE "Other", "A" int, CONSTRAINT "Pk" PRIMARY KEY ("A"));`,
			config: pattern,
			want:   []columnFinding{mismatch("t", "A", `CREATE TABLE t (LIKE "Other", "A" int, CONSTRAINT "Pk" PRIMARY KEY ("A"))`)},
		},
		{
			name:   "add column judged, with and without the keyword",
			sql:    `ALTER TABLE s.book ADD COLUMN "A" int, ADD "B" int, ADD COLUMN IF NOT EXISTS c int, DROP COLUMN "D", ALTER COLUMN "E" TYPE int;`,
			config: pattern,
			want: []columnFinding{
				mismatch("book", "A", `ALTER TABLE s.book ADD COLUMN "A" int, ADD "B" int, ADD COLUMN IF NOT EXISTS c int, DROP COLUMN "D", ALTER COLUMN "E" TYPE int`),
				mismatch("book", "B", `ALTER TABLE s.book ADD COLUMN "A" int, ADD "B" int, ADD COLUMN IF NOT EXISTS c int, DROP COLUMN "D", ALTER COLUMN "E" TYPE int`),
			},
		},
		{
			name:   "add column judged on a foreign table and a type",
			sql:    `ALTER FOREIGN TABLE ft ADD COLUMN "A" int; ALTER TYPE ty ADD ATTRIBUTE "B" int;`,
			config: pattern,
			want: []columnFinding{
				mismatch("ft", "A", `ALTER FOREIGN TABLE ft ADD COLUMN "A" int`),
				mismatch("ty", "B", `ALTER TYPE ty ADD ATTRIBUTE "B" int`),
			},
		},
		{
			name:   "rename column judges only the new name",
			sql:    `ALTER TABLE book RENAME COLUMN "Old" TO "New"; ALTER TABLE book RENAME "Old" TO fine;`,
			config: pattern,
			want:   []columnFinding{mismatch("book", "New", `ALTER TABLE book RENAME COLUMN "Old" TO "New"`)},
		},
		{
			name:   "rename column judged on views and foreign tables",
			sql:    `ALTER VIEW v RENAME COLUMN a TO "A"; ALTER MATERIALIZED VIEW mv RENAME COLUMN a TO "B"; ALTER FOREIGN TABLE ft RENAME COLUMN a TO "C";`,
			config: pattern,
			want: []columnFinding{
				mismatch("v", "A", `ALTER VIEW v RENAME COLUMN a TO "A"`),
				mismatch("mv", "B", `ALTER MATERIALIZED VIEW mv RENAME COLUMN a TO "B"`),
				mismatch("ft", "C", `ALTER FOREIGN TABLE ft RENAME COLUMN a TO "C"`),
			},
		},
		{
			name:   "other renames not judged",
			sql:    `ALTER TABLE t RENAME TO "Bad"; ALTER TABLE t RENAME CONSTRAINT c TO "Bad"; ALTER TYPE ty RENAME ATTRIBUTE a TO "Bad"; ALTER INDEX i RENAME TO "Bad";`,
			config: pattern,
		},
		{
			name:   "foreign table, create table as, and view columns not judged",
			sql:    `CREATE FOREIGN TABLE ft ("A" int) SERVER srv; CREATE TABLE t AS SELECT 1 AS "A"; CREATE TABLE t ("A") AS SELECT 1; CREATE VIEW v ("A") AS SELECT 1;`,
			config: pattern,
		},
		{
			name:   "create schema element and routine body not judged",
			sql:    `CREATE SCHEMA s CREATE TABLE t ("A" int); CREATE PROCEDURE pr() LANGUAGE SQL BEGIN ATOMIC ALTER TABLE t RENAME COLUMN a TO "A"; END;`,
			config: pattern,
		},
		{
			name:   "length counts bytes",
			sql:    `CREATE TABLE t ("数据列" int, "abcdefgh" int);`,
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []columnFinding{{v1review.NameTooLong, "t", "数据列", `CREATE TABLE t ("数据列" int, "abcdefgh" int)`}},
		},
		{
			name:   "name truncated to 63 bytes before judging",
			sql:    `CREATE TABLE t ("` + strings.Repeat("a", 70) + `" int);`,
			config: v1review.NamingConfig{MaxLength: 63},
		},
		{
			name:   "pattern then length for each name",
			sql:    `ALTER TABLE t ADD COLUMN "Toolong5" int, ADD COLUMN "X" int;`,
			config: v1review.NamingConfig{Pattern: "^[a-z]+$", MaxLength: 4},
			want: []columnFinding{
				{v1review.NameMismatch, "t", "Toolong5", `ALTER TABLE t ADD COLUMN "Toolong5" int, ADD COLUMN "X" int`},
				{v1review.NameTooLong, "t", "Toolong5", `ALTER TABLE t ADD COLUMN "Toolong5" int, ADD COLUMN "X" int`},
				{v1review.NameMismatch, "t", "X", `ALTER TABLE t ADD COLUMN "Toolong5" int, ADD COLUMN "X" int`},
			},
		},
		{
			name:   "pattern matches anywhere unless anchored",
			sql:    `CREATE TABLE t ("Has_id_Inside" int, "NONE" int);`,
			config: v1review.NamingConfig{Pattern: "_id"},
			want:   []columnFinding{mismatch("t", "NONE", `CREATE TABLE t ("Has_id_Inside" int, "NONE" int)`)},
		},
		{name: "empty config requires nothing", sql: `CREATE TABLE t ("A" int);`},
		{
			name:   "anchor excludes a leading comment and the semicolon",
			sql:    "/* c */\nCREATE TABLE t (\n  \"A\" int\n);",
			config: pattern,
			want:   []columnFinding{mismatch("t", "A", "CREATE TABLE t (\n  \"A\" int\n)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runNamingColumn(t, tt.sql, tt.config)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingColumnInvalidPattern(t *testing.T) {
	list, err := parser.Parse("CREATE TABLE t (id int);")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := NamingColumn(list.Items, v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
