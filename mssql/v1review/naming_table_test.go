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
			name:   "create mismatch anchors the statement",
			sql:    "CREATE TABLE OrderItem (id INT PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("OrderItem", "CREATE TABLE OrderItem (id INT PRIMARY KEY)")},
		},
		{name: "qualifier not judged", sql: "CREATE TABLE Dbo.order_item (id INT PRIMARY KEY);", config: pattern},
		{
			name:   "bracketed name judged unbracketed",
			sql:    "CREATE TABLE [Bad Name] (id INT);",
			config: pattern,
			want:   []finding{mismatch("Bad Name", "CREATE TABLE [Bad Name] (id INT)")},
		},
		{
			name:   "temporary table judged with its prefix",
			sql:    "CREATE TABLE #TmpStage (id INT);",
			config: pattern,
			want:   []finding{mismatch("#TmpStage", "CREATE TABLE #TmpStage (id INT)")},
		},
		{name: "select into not judged", sql: "SELECT 1 AS id INTO BadSel;", config: pattern},
		{name: "create table as select not judged", sql: "CREATE TABLE BadCtas WITH (DISTRIBUTION = ROUND_ROBIN) AS SELECT 1 AS id;", config: pattern},
		{name: "create table as clone not judged", sql: "CREATE TABLE BadClone AS CLONE OF dbo.t;", config: pattern},
		{name: "rename object not judged", sql: "RENAME OBJECT dbo.t TO NewName;", config: pattern},
		{name: "routine body not judged", sql: "CREATE PROCEDURE pr2 AS CREATE TABLE Bad (id INT);", config: pattern},
		{name: "if block not judged", sql: "IF 1 = 1 BEGIN CREATE TABLE BadIf (id INT); END", config: pattern},
		{
			name:   "sp_rename judges the new name",
			sql:    "EXEC sp_rename 'dbo.BadOld', 'NewName';",
			config: pattern,
			want:   []finding{mismatch("NewName", "EXEC sp_rename 'dbo.BadOld', 'NewName'")},
		},
		{
			name:   "sp_rename of a column judged",
			sql:    "EXEC sp_rename 'dbo.t.n', 'BadCol', 'COLUMN';",
			config: pattern,
			want:   []finding{mismatch("BadCol", "EXEC sp_rename 'dbo.t.n', 'BadCol', 'COLUMN'")},
		},
		{
			name:   "sp_rename of a database judged",
			sql:    "EXEC sp_rename 'olddb', 'BadDb', 'DATABASE';",
			config: pattern,
			want:   []finding{mismatch("BadDb", "EXEC sp_rename 'olddb', 'BadDb', 'DATABASE'")},
		},
		{
			name:   "sp_rename second argument read positionally",
			sql:    "EXEC sp_rename @newname = 'new_name', @objname = 'Sales.Territory';",
			config: pattern,
			want:   []finding{mismatch("Sales.Territory", "EXEC sp_rename @newname = 'new_name', @objname = 'Sales.Territory'")},
		},
		{
			name:   "sp_rename with a variable old name judged",
			sql:    "EXEC sp_rename @objname = @old, @newname = 'BadVar';",
			config: pattern,
			want:   []finding{mismatch("BadVar", "EXEC sp_rename @objname = @old, @newname = 'BadVar'")},
		},
		{
			name:   "sp_rename with a four-part old name judged",
			sql:    "EXEC sp_rename 'srv.db.dbo.t', 'BadFour';",
			config: pattern,
			want:   []finding{mismatch("BadFour", "EXEC sp_rename 'srv.db.dbo.t', 'BadFour'")},
		},
		{
			name:   "sp_rename new name loses one pair of enclosing quotes",
			sql:    "EXEC sp_rename 'dbo.t', '''Quoted''';",
			config: pattern,
			want:   []finding{mismatch("Quoted", "EXEC sp_rename 'dbo.t', '''Quoted'''")},
		},
		{name: "sp_rename with a variable new name not judged", sql: "EXEC sp_rename 'dbo.t', @new;", config: pattern},
		{name: "sys.sp_rename not judged", sql: "EXEC sys.sp_rename 'dbo.t', 'BadSys';", config: pattern},
		{
			name:   "length counts bytes",
			sql:    "CREATE TABLE [数据表] (id INT);",
			config: v1review.NamingConfig{MaxLength: 8},
			want:   []finding{{v1review.NameTooLong, "数据表", "CREATE TABLE [数据表] (id INT)"}},
		},
		{name: "length within limit", sql: "CREATE TABLE order_item (id INT PRIMARY KEY);", config: v1review.NamingConfig{MaxLength: 10}},
		{name: "multibyte length within limit", sql: "CREATE TABLE [数据表] (id INT);", config: v1review.NamingConfig{MaxLength: 9}},
		{
			name:   "pattern and length both reported",
			sql:    "CREATE TABLE Toolong5 (id INT);",
			config: v1review.NamingConfig{Pattern: "^[a-z]+$", MaxLength: 4},
			want: []finding{
				{v1review.NameMismatch, "Toolong5", "CREATE TABLE Toolong5 (id INT)"},
				{v1review.NameTooLong, "Toolong5", "CREATE TABLE Toolong5 (id INT)"},
			},
		},
		{name: "empty config requires nothing", sql: "CREATE TABLE OrderItem (id INT PRIMARY KEY);"},
		{
			name:   "second statement",
			sql:    "CREATE TABLE order_item (id INT PRIMARY KEY);\nCREATE TABLE OrderItem (id INT PRIMARY KEY);",
			config: pattern,
			want:   []finding{mismatch("OrderItem", "CREATE TABLE OrderItem (id INT PRIMARY KEY)")},
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
