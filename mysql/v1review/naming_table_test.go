package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/parser"
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
		{name: "qualifier not judged", sql: "CREATE TABLE Shop.order_item (id INT PRIMARY KEY);", config: pattern},
		{
			name:   "qualified name judged",
			sql:    "CREATE TABLE shop.OrderItem (id INT PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("OrderItem", "shop.OrderItem")},
		},
		{
			name:   "temporary table judged",
			sql:    "CREATE TEMPORARY TABLE TmpStage (id INT);",
			config: pattern,
			want:   []finding{mismatch("TmpStage", "TmpStage")},
		},
		{
			name:   "create like",
			sql:    "CREATE TABLE BadCopy LIKE t;",
			config: pattern,
			want:   []finding{mismatch("BadCopy", "BadCopy")},
		},
		{
			name:   "create as select",
			sql:    "CREATE TABLE BadSel AS SELECT 1;",
			config: pattern,
			want:   []finding{mismatch("BadSel", "BadSel")},
		},
		{
			name:   "rename table judges the new name",
			sql:    "RENAME TABLE BadOld TO NewName;",
			config: pattern,
			want:   []finding{mismatch("NewName", "NewName")},
		},
		{
			name:   "alter rename to",
			sql:    "ALTER TABLE t RENAME TO NewName;",
			config: pattern,
			want:   []finding{mismatch("NewName", "NewName")},
		},
		{
			name:   "alter rename bare",
			sql:    "ALTER TABLE t RENAME NewName;",
			config: pattern,
			want:   []finding{mismatch("NewName", "NewName")},
		},
		{
			name:   "routine body",
			sql:    "CREATE PROCEDURE pr() BEGIN CREATE TABLE Bad (id INT); END;",
			config: pattern,
			want:   []finding{mismatch("Bad", "Bad")},
		},
		{
			name:   "length counts code points",
			sql:    "CREATE TABLE `数据表` (id INT);",
			config: v1review.NamingConfig{MaxLength: 2},
			want:   []finding{{v1review.NameTooLong, "数据表", "`数据表`"}},
		},
		{name: "length within limit", sql: "CREATE TABLE order_item (id INT PRIMARY KEY);", config: v1review.NamingConfig{MaxLength: 10}},
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
			name:   "findings in statement order",
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
