package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/snowflake/parser"
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
	file, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	got, err := NamingTable(file.Stmts, cfg)
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
		{
			name:   "a statement under EXPLAIN, which does not run it",
			sql:    "EXPLAIN CREATE TABLE OrderItem (id INT);",
			config: pattern,
		},
		{name: "create ok", sql: "CREATE TABLE order_item (id INT PRIMARY KEY);", config: pattern},
		{name: "identifier literal judged by the name it spells", sql: "CREATE TABLE IDENTIFIER('order_item') (id INT);", config: pattern},
		{name: "identifier literal qualifier not judged", sql: "CREATE TABLE IDENTIFIER('Shop.PUBLIC.order_item') (id INT);", config: pattern},
		{
			name:   "identifier literal mismatch",
			sql:    `CREATE TABLE identifier('"OrderItem"') (id INT);`,
			config: pattern,
			want:   []finding{mismatch("OrderItem", `identifier('"OrderItem"')`)},
		},
		{name: "identifier variable not judged", sql: "CREATE TABLE IDENTIFIER($tbl) (id INT);", config: pattern},
		{
			name:   "identifier literal on rename",
			sql:    "ALTER TABLE t RENAME TO IDENTIFIER('OrderItem');",
			config: pattern,
			want:   []finding{mismatch("OrderItem", "IDENTIFIER('OrderItem')")},
		},
		{
			name:   "create mismatch",
			sql:    "CREATE TABLE OrderItem (id INT PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("OrderItem", "OrderItem")},
		},
		{name: "qualifier not judged", sql: "CREATE TABLE Shop.order_item (id INT PRIMARY KEY);", config: pattern},
		{
			name:   "qualified name judged",
			sql:    "CREATE TABLE shop.public.OrderItem (id INT PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("OrderItem", "shop.public.OrderItem")},
		},
		{
			name:   "or replace",
			sql:    "CREATE OR REPLACE TABLE OrderItem (id INT);",
			config: pattern,
			want:   []finding{mismatch("OrderItem", "OrderItem")},
		},
		{
			name:   "temporary and transient tables judged",
			sql:    "CREATE TEMPORARY TABLE TmpStage (id INT);\nCREATE TRANSIENT TABLE TmpStage (id INT);",
			config: pattern,
			want: []finding{
				mismatch("TmpStage", "TmpStage"),
				mismatch("TmpStage", "TmpStage"),
			},
		},
		{
			name:   "create as select",
			sql:    "CREATE TABLE BadSel AS SELECT id FROM t;",
			config: pattern,
			want:   []finding{mismatch("BadSel", "BadSel")},
		},
		{
			name:   "create like and clone",
			sql:    "CREATE TABLE BadCopy LIKE t;\nCREATE TABLE BadClone CLONE t;",
			config: pattern,
			want: []finding{
				mismatch("BadCopy", "BadCopy"),
				mismatch("BadClone", "BadClone"),
			},
		},
		{
			name:   "alter rename judges only the new name",
			sql:    "ALTER TABLE BadOld RENAME TO NewName;\nALTER TABLE IF EXISTS t RENAME TO s.NewName;",
			config: pattern,
			want: []finding{
				mismatch("NewName", "NewName"),
				mismatch("NewName", "s.NewName"),
			},
		},
		{
			name:   "a swap, a dynamic table, and a view are not judged",
			sql:    "ALTER TABLE t SWAP WITH BadName;\nALTER DYNAMIC TABLE dt RENAME TO BadName;\nALTER VIEW v RENAME TO BadName;",
			config: pattern,
		},
		{
			name:   "names are judged as written, quotes left out",
			sql:    "CREATE TABLE \"OrderItem\" (id INT);\nCREATE TABLE \"order_item\" (id INT);\nCREATE TABLE ORDER_ITEM (id INT);",
			config: pattern,
			want: []finding{
				mismatch("OrderItem", `"OrderItem"`),
				mismatch("ORDER_ITEM", "ORDER_ITEM"),
			},
		},
		{
			name:   "length counts code points",
			sql:    "CREATE TABLE \"数据表\" (id INT);",
			config: v1review.NamingConfig{MaxLength: 2},
			want:   []finding{{v1review.NameTooLong, "数据表", `"数据表"`}},
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
		{
			name:   "inside a block",
			sql:    "BEGIN\n  CREATE TABLE Bad (id INT);\nEND",
			config: pattern,
			want:   []finding{mismatch("Bad", "Bad")},
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
	file, err := parser.Parse("CREATE TABLE t (id INT);")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := NamingTable(file.Stmts, v1review.NamingConfig{Pattern: "("}); err == nil {
		t.Fatal("invalid pattern: no error")
	}
}
