package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mssql/parser"
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
		{name: "create ok", sql: "CREATE TABLE order_item (id INT PRIMARY KEY);", config: pattern},
		{
			name:   "create mismatch",
			sql:    "CREATE TABLE OrderItem (id INT PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("OrderItem", "OrderItem")},
		},
		{name: "qualifier not judged", sql: "CREATE TABLE Dbo.order_item (id INT PRIMARY KEY);", config: pattern},
		{
			name:   "qualified name judged",
			sql:    "CREATE TABLE dbo.BadName (id INT PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("BadName", "dbo.BadName")},
		},
		{
			name:   "bracketed name judged unbracketed",
			sql:    "CREATE TABLE [Bad Name] (id INT);",
			config: pattern,
			want:   []finding{mismatch("Bad Name", "[Bad Name]")},
		},
		{
			name:   "temporary table judged with its prefix",
			sql:    "CREATE TABLE #TmpStage (id INT);",
			config: pattern,
			want:   []finding{mismatch("#TmpStage", "#TmpStage")},
		},
		{
			name:   "select into",
			sql:    "SELECT 1 AS id INTO BadSel;",
			config: pattern,
			want:   []finding{mismatch("BadSel", "BadSel")},
		},
		{
			name:   "create table as select",
			sql:    "CREATE TABLE BadCtas WITH (DISTRIBUTION = ROUND_ROBIN) AS SELECT 1 AS id;",
			config: pattern,
			want:   []finding{mismatch("BadCtas", "BadCtas")},
		},
		{
			name:   "create table as clone",
			sql:    "CREATE TABLE BadClone AS CLONE OF dbo.t;",
			config: pattern,
			want:   []finding{mismatch("BadClone", "BadClone")},
		},
		{
			name:   "sp_rename judges the new name",
			sql:    "EXEC sp_rename 'dbo.BadOld', 'NewName';",
			config: pattern,
			want:   []finding{mismatch("NewName", "EXEC sp_rename 'dbo.BadOld', 'NewName'")},
		},
		{name: "sp_rename of a column not judged", sql: "EXEC sp_rename 'dbo.t.n', 'BadCol', 'COLUMN';", config: pattern},
		{
			name:   "rename object",
			sql:    "RENAME OBJECT dbo.t TO NewName;",
			config: pattern,
			want:   []finding{mismatch("NewName", "RENAME OBJECT dbo.t TO NewName")},
		},
		{
			name:   "routine body",
			sql:    "CREATE PROCEDURE pr2 AS CREATE TABLE Bad (id INT);",
			config: pattern,
			want:   []finding{mismatch("Bad", "Bad")},
		},
		{
			name:   "length counts code points",
			sql:    "CREATE TABLE [数据表] (id INT);",
			config: v1review.NamingConfig{MaxLength: 2},
			want:   []finding{{v1review.NameTooLong, "数据表", "[数据表]"}},
		},
		{name: "length within limit", sql: "CREATE TABLE order_item (id INT PRIMARY KEY);", config: v1review.NamingConfig{MaxLength: 10}},
		{name: "multibyte length within limit", sql: "CREATE TABLE [数据表] (id INT);", config: v1review.NamingConfig{MaxLength: 3}},
		{
			name:   "pattern and length both reported",
			sql:    "CREATE TABLE Toolong5 (id INT);",
			config: v1review.NamingConfig{Pattern: "^[a-z]+$", MaxLength: 4},
			want: []finding{
				{v1review.NameMismatch, "Toolong5", "Toolong5"},
				{v1review.NameTooLong, "Toolong5", "Toolong5"},
			},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE OrderItem (id INT PRIMARY KEY);"},
		{
			name:   "second statement",
			sql:    "CREATE TABLE order_item (id INT PRIMARY KEY);\nCREATE TABLE OrderItem (id INT PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("OrderItem", "OrderItem")},
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
	list, err := parser.Parse("CREATE TABLE t (id INT);")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := NamingTable(list.Items, v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
