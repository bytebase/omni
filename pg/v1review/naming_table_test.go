package v1review

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

// finding is a finding as the tests compare it: the code, the judged name,
// and the text its range anchors.
type finding struct {
	Code v1review.Code
	Name string
	Text string
}

func run(t *testing.T, sql string, cfg v1review.NamingConfig) []finding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := NamingTable(list.Items, cfg)
	if err != nil {
		t.Fatalf("NamingTable: %v", err)
	}
	var findings []finding
	for _, f := range got {
		findings = append(findings, finding{f.Code, f.Name, sql[f.Start:f.End]})
	}
	return findings
}

func TestNamingTable(t *testing.T) {
	pattern := v1review.NamingConfig{Pattern: "^[a-z]+(_[a-z]+)*$"}
	mismatch := func(name, text string) finding { return finding{v1review.NameMismatch, name, text} }
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []finding
	}{
		{name: "create ok", sql: "CREATE TABLE order_item (id int PRIMARY KEY);", config: pattern},
		{
			name:   "create mismatch",
			sql:    `CREATE TABLE "BadName" (id int PRIMARY KEY);`,
			config: pattern,
			want:   []finding{mismatch("BadName", `CREATE TABLE "BadName" (id int PRIMARY KEY)`)},
		},
		{
			name:   "unquoted mixed case folds to lower and passes",
			sql:    "CREATE TABLE OrderItem (id int PRIMARY KEY);",
			config: pattern,
		},
		{
			name:   "quoted mixed case keeps its case and fails",
			sql:    `CREATE TABLE "OrderItem" (id int PRIMARY KEY);`,
			config: pattern,
			want:   []finding{mismatch("OrderItem", `CREATE TABLE "OrderItem" (id int PRIMARY KEY)`)},
		},
		{name: "qualifier not judged", sql: `CREATE TABLE "Shop".order_item (id int PRIMARY KEY);`, config: pattern},
		{
			name:   "qualified name judged",
			sql:    `CREATE TABLE shop."OrderItem" (id int PRIMARY KEY);`,
			config: pattern,
			want:   []finding{mismatch("OrderItem", `CREATE TABLE shop."OrderItem" (id int PRIMARY KEY)`)},
		},
		{
			name:   "temporary table judged",
			sql:    `CREATE TEMPORARY TABLE "TmpStage" (id int);`,
			config: pattern,
			want:   []finding{mismatch("TmpStage", `CREATE TEMPORARY TABLE "TmpStage" (id int)`)},
		},
		{
			name:   "unlogged if not exists judged",
			sql:    `CREATE UNLOGGED TABLE IF NOT EXISTS "BadStage" (id int);`,
			config: pattern,
			want:   []finding{mismatch("BadStage", `CREATE UNLOGGED TABLE IF NOT EXISTS "BadStage" (id int)`)},
		},
		{
			name:   "partition child judged, parent not",
			sql:    `CREATE TABLE "BadPart" PARTITION OF "BadParent" FOR VALUES IN (1);`,
			config: pattern,
			want:   []finding{mismatch("BadPart", `CREATE TABLE "BadPart" PARTITION OF "BadParent" FOR VALUES IN (1)`)},
		},
		{name: "create table as not judged", sql: `CREATE TABLE "BadSel" AS SELECT 1;`, config: pattern},
		{name: "create table as table not judged", sql: `CREATE TABLE "BadSel" AS TABLE t;`, config: pattern},
		{name: "select into not judged", sql: `SELECT 1 INTO "BadSel";`, config: pattern},
		{name: "explain analyze create table as not judged", sql: `EXPLAIN ANALYZE CREATE TABLE "BadName" AS SELECT 1;`, config: pattern},
		{
			name:   "rename judges only the new name",
			sql:    `ALTER TABLE "BadOld" RENAME TO "NewName";`,
			config: pattern,
			want:   []finding{mismatch("NewName", `ALTER TABLE "BadOld" RENAME TO "NewName"`)},
		},
		{
			name:   "routine body not judged",
			sql:    `CREATE PROCEDURE pr() LANGUAGE SQL BEGIN ATOMIC ALTER TABLE t RENAME TO "Bad"; END;`,
			config: pattern,
		},
		{name: "create schema element not judged", sql: `CREATE SCHEMA s CREATE TABLE "Bad" (id int);`, config: pattern},
		{
			name:   "views, foreign tables, and columns not judged",
			sql:    `CREATE VIEW "BadView" AS SELECT 1; CREATE MATERIALIZED VIEW "BadMv" AS SELECT 1; CREATE FOREIGN TABLE "BadFt" (a int) SERVER srv; ALTER VIEW v RENAME TO "BadV"; ALTER MATERIALIZED VIEW mv RENAME TO "BadM"; ALTER FOREIGN TABLE ft RENAME TO "BadF"; ALTER TABLE t RENAME COLUMN a TO "BadCol";`,
			config: pattern,
		},
		{
			name:   "length counts bytes",
			sql:    `CREATE TABLE "数据表" (id int);`,
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []finding{{v1review.NameTooLong, "数据表", `CREATE TABLE "数据表" (id int)`}},
		},
		{
			name:   "name truncated to 63 bytes before judging",
			sql:    `CREATE TABLE "` + strings.Repeat("a", 70) + `" (id int);`,
			config: v1review.NamingConfig{MaxLength: 63},
		},
		{name: "length within limit", sql: "CREATE TABLE order_item (id int PRIMARY KEY);", config: v1review.NamingConfig{MaxLength: 10}},
		{
			name:   "pattern and length both reported",
			sql:    `CREATE TABLE "Toolong5" (id int);`,
			config: v1review.NamingConfig{Pattern: "^[a-z]+$", MaxLength: 4},
			want: []finding{
				{v1review.NameMismatch, "Toolong5", `CREATE TABLE "Toolong5" (id int)`},
				{v1review.NameTooLong, "Toolong5", `CREATE TABLE "Toolong5" (id int)`},
			},
		},
		{name: "empty config requires nothing", sql: `CREATE TABLE "OrderItem" (id int PRIMARY KEY);`},
		{
			name:   "findings in statement order",
			sql:    "CREATE TABLE order_item (id int PRIMARY KEY);\nCREATE TABLE \"OrderItem\" (id int PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("OrderItem", `CREATE TABLE "OrderItem" (id int PRIMARY KEY)`)},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, run(t, tt.sql, tt.config)); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNamingTableInvalidPattern(t *testing.T) {
	list, err := parser.Parse("CREATE TABLE t (id int);")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := NamingTable(list.Items, v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
