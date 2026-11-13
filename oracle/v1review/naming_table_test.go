package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oracle/parser"
	"github.com/bytebase/omni/v1review"
)

// finding is a finding as the tests compare it: the code, the judged name,
// and the text its range anchors: the whole CREATE TABLE, or the RENAME
// clause of an ALTER TABLE.
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
	pattern := v1review.NamingConfig{Pattern: "^[A-Z][A-Z0-9_]*$"}
	mismatch := func(name, text string) finding { return finding{v1review.NameMismatch, name, text} }
	for _, tt := range []struct {
		name   string
		sql    string
		config v1review.NamingConfig
		want   []finding
	}{
		{name: "create ok", sql: "CREATE TABLE order_item (id NUMBER PRIMARY KEY);", config: pattern},
		{
			name:   "create mismatch",
			sql:    `CREATE TABLE "orderItem" (id NUMBER PRIMARY KEY);`,
			config: pattern,
			want:   []finding{mismatch("orderItem", `CREATE TABLE "orderItem" (id NUMBER PRIMARY KEY)`)},
		},
		{name: "qualifier not judged", sql: `CREATE TABLE "shop".order_item (id NUMBER PRIMARY KEY);`, config: pattern},
		{
			name:   "qualified name judged",
			sql:    `CREATE TABLE shop."orderItem" (id NUMBER PRIMARY KEY);`,
			config: pattern,
			want:   []finding{mismatch("orderItem", `CREATE TABLE shop."orderItem" (id NUMBER PRIMARY KEY)`)},
		},
		{
			name:   "global temporary table judged",
			sql:    `CREATE GLOBAL TEMPORARY TABLE "tmpStage" (id NUMBER);`,
			config: pattern,
			want:   []finding{mismatch("tmpStage", `CREATE GLOBAL TEMPORARY TABLE "tmpStage" (id NUMBER)`)},
		},
		{
			name:   "private temporary table judged",
			sql:    "CREATE PRIVATE TEMPORARY TABLE ora$ptt_stage (id NUMBER);",
			config: pattern,
			want:   []finding{mismatch("ORA$PTT_STAGE", "CREATE PRIVATE TEMPORARY TABLE ora$ptt_stage (id NUMBER)")},
		},
		{
			name:   "create as select",
			sql:    `CREATE TABLE "badSel" AS SELECT 1 FROM dual;`,
			config: pattern,
			want:   []finding{mismatch("badSel", `CREATE TABLE "badSel" AS SELECT 1 FROM dual`)},
		},
		{name: "rename statement not judged", sql: `RENAME "badOld" TO "newName";`, config: pattern},
		{name: "explain plan not judged", sql: `EXPLAIN PLAN FOR CREATE TABLE "badName" AS SELECT 1 AS id FROM dual`, config: pattern},
		{
			name:   "create schema nested table not judged",
			sql:    `CREATE SCHEMA AUTHORIZATION shop CREATE TABLE "badNested" (id NUMBER)`,
			config: pattern,
		},
		{
			name:   "alter rename to",
			sql:    `ALTER TABLE t RENAME TO "newName";`,
			config: pattern,
			want:   []finding{mismatch("newName", `RENAME TO "newName"`)},
		},
		{
			name:   "unquoted name folds to upper case",
			sql:    "CREATE TABLE bad_name (id NUMBER);",
			config: v1review.NamingConfig{Pattern: "^[a-z_]+$"},
			want:   []finding{mismatch("BAD_NAME", "CREATE TABLE bad_name (id NUMBER)")},
		},
		{
			name:   "quoted name kept as written",
			sql:    `CREATE TABLE "bad_name" (id NUMBER);`,
			config: v1review.NamingConfig{Pattern: "^[a-z_]+$"},
		},
		{
			name:   "length counts bytes",
			sql:    `CREATE TABLE "数据表" (id NUMBER);`,
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []finding{{v1review.NameTooLong, "数据表", `CREATE TABLE "数据表" (id NUMBER)`}},
		},
		{name: "length within limit", sql: "CREATE TABLE order_item (id NUMBER PRIMARY KEY);", config: v1review.NamingConfig{MaxLength: 10}},
		{
			name:   "pattern and length both reported",
			sql:    "CREATE TABLE toolong5 (id NUMBER);",
			config: v1review.NamingConfig{Pattern: "^[A-Z]+$", MaxLength: 4},
			want: []finding{
				{v1review.NameMismatch, "TOOLONG5", "CREATE TABLE toolong5 (id NUMBER)"},
				{v1review.NameTooLong, "TOOLONG5", "CREATE TABLE toolong5 (id NUMBER)"},
			},
		},
		{name: "empty config requires nothing", sql: `CREATE TABLE "orderItem" (id NUMBER PRIMARY KEY);`},
		{
			name:   "findings in statement order",
			sql:    "CREATE TABLE order_item (id NUMBER PRIMARY KEY);\nCREATE TABLE \"orderItem\" (id NUMBER PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("orderItem", `CREATE TABLE "orderItem" (id NUMBER PRIMARY KEY)`)},
		},
		{
			name:   "multi-line create anchors the statement",
			sql:    "CREATE TABLE\n  \"orderItem\" (id NUMBER);",
			config: pattern,
			want:   []finding{mismatch("orderItem", "CREATE TABLE\n  \"orderItem\" (id NUMBER)")},
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
	list, err := parser.Parse("CREATE TABLE t (id NUMBER);")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := NamingTable(list.Items, v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
