package parser

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// dropItem is an ALTER TABLE ... DROP action as the tests compare it: its
// kind, the name it drops, whether it carries IF EXISTS, and the text its Loc
// spans.
type dropItem struct {
	kind     string // "column", "constraint", or "period"
	name     string
	ifExists bool
	text     string
}

// dropListCases lists ALTER TABLE ... DROP lists with SQL Server 2022's
// verdict under SET PARSEONLY ON, which TestAlterTableDropListOracle checks.
// An accepted case lists the actions omni builds, one per item; a rejected
// one names the token omni's syntax error is at, or "end of input".
var dropListCases = []struct {
	name    string
	sql     string
	accepts bool
	want    []dropItem
	near    string
}{
	// An item may restate CONSTRAINT or COLUMN.
	{name: "constraint_restated", sql: "ALTER TABLE dbo.t DROP CONSTRAINT fk_p, CONSTRAINT ck_n", accepts: true, want: []dropItem{
		{"constraint", "fk_p", false, "fk_p"},
		{"constraint", "ck_n", false, "ck_n"},
	}},
	{name: "constraint_then_column", sql: "ALTER TABLE dbo.t DROP CONSTRAINT uk_code, COLUMN n", accepts: true, want: []dropItem{
		{"constraint", "uk_code", false, "uk_code"},
		{"column", "n", false, "n"},
	}},
	{name: "column_then_constraint", sql: "ALTER TABLE dbo.t DROP COLUMN a, CONSTRAINT ck_n", accepts: true, want: []dropItem{
		{"column", "a", false, "a"},
		{"constraint", "ck_n", false, "ck_n"},
	}},
	{name: "bracketed_names", sql: "ALTER TABLE dbo.t DROP CONSTRAINT [fk p], COLUMN [n]", accepts: true, want: []dropItem{
		{"constraint", "fk p", false, "[fk p]"},
		{"column", "n", false, "[n]"},
	}},

	// An item without a keyword drops what the item before it dropped; the
	// first one drops a constraint.
	{name: "column_list", sql: "ALTER TABLE dbo.t DROP COLUMN a, b", accepts: true, want: []dropItem{
		{"column", "a", false, "a"},
		{"column", "b", false, "b"},
	}},
	{name: "bare_list", sql: "ALTER TABLE dbo.t DROP fk_p, ck_n", accepts: true, want: []dropItem{
		{"constraint", "fk_p", false, "fk_p"},
		{"constraint", "ck_n", false, "ck_n"},
	}},
	{name: "bare_then_column", sql: "ALTER TABLE dbo.t DROP ck_n, COLUMN a", accepts: true, want: []dropItem{
		{"constraint", "ck_n", false, "ck_n"},
		{"column", "a", false, "a"},
	}},
	{name: "kinds_inherited", sql: "ALTER TABLE dbo.t DROP COLUMN a, b, CONSTRAINT ck_n, df_n, COLUMN c", accepts: true, want: []dropItem{
		{"column", "a", false, "a"},
		{"column", "b", false, "b"},
		{"constraint", "ck_n", false, "ck_n"},
		{"constraint", "df_n", false, "df_n"},
		{"column", "c", false, "c"},
	}},

	// IF EXISTS covers only the item that writes it, after its keyword.
	{name: "if_exists_per_item", sql: "ALTER TABLE dbo.t DROP COLUMN IF EXISTS a, b, COLUMN IF EXISTS c", accepts: true, want: []dropItem{
		{"column", "a", true, "a"},
		{"column", "b", false, "b"},
		{"column", "c", true, "c"},
	}},
	{name: "if_exists_restated_constraint", sql: "ALTER TABLE dbo.t DROP COLUMN a, CONSTRAINT IF EXISTS ck_n WITH (ONLINE = ON)", accepts: true, want: []dropItem{
		{"column", "a", false, "a"},
		{"constraint", "ck_n", true, "ck_n WITH (ONLINE = ON)"},
	}},
	{name: "if_exists_without_keyword", sql: "ALTER TABLE dbo.t DROP IF EXISTS ck_n", near: "IF"},
	{name: "if_exists_after_comma_without_keyword", sql: "ALTER TABLE dbo.t DROP COLUMN a, IF EXISTS b", near: "IF"},
	{name: "if_exists_twice_without_keyword", sql: "ALTER TABLE dbo.t DROP CONSTRAINT IF EXISTS a, IF EXISTS b", near: "IF"},
	{name: "if_exists_before_keyword", sql: "ALTER TABLE dbo.t DROP IF EXISTS COLUMN a", near: "IF"},

	// WITH ( ... ) follows any name COLUMN does not introduce.
	{name: "with_then_column", sql: "ALTER TABLE dbo.t DROP CONSTRAINT pk_t WITH (ONLINE = ON), COLUMN b", accepts: true, want: []dropItem{
		{"constraint", "pk_t", false, "pk_t WITH (ONLINE = ON)"},
		{"column", "b", false, "b"},
	}},
	{name: "with_on_each_constraint", sql: "ALTER TABLE dbo.t DROP CONSTRAINT pk_t WITH (ONLINE = ON), uk_code WITH (MAXDOP = 1)", accepts: true, want: []dropItem{
		{"constraint", "pk_t", false, "pk_t WITH (ONLINE = ON)"},
		{"constraint", "uk_code", false, "uk_code WITH (MAXDOP = 1)"},
	}},
	{name: "with_on_inherited_column", sql: "ALTER TABLE dbo.t DROP COLUMN a, b WITH (ONLINE = ON)", accepts: true, want: []dropItem{
		{"column", "a", false, "a"},
		{"column", "b", false, "b WITH (ONLINE = ON)"},
	}},
	{name: "with_after_column", sql: "ALTER TABLE dbo.t DROP COLUMN a WITH (ONLINE = ON)", near: "("},
	{name: "with_after_column_if_exists", sql: "ALTER TABLE dbo.t DROP COLUMN IF EXISTS a WITH (ONLINE = ON), b", near: "("},
	{name: "with_without_list", sql: "ALTER TABLE dbo.t DROP CONSTRAINT pk_t WITH", near: "end of input"},

	// PERIOD FOR SYSTEM_TIME is an item; a bare name after it is a period too.
	// PERIOD without FOR is a name.
	{name: "period", sql: "ALTER TABLE dbo.t DROP PERIOD FOR SYSTEM_TIME", accepts: true, want: []dropItem{
		{"period", "", false, "PERIOD FOR SYSTEM_TIME"},
	}},
	{name: "period_then_columns", sql: "ALTER TABLE dbo.t DROP PERIOD FOR SYSTEM_TIME, COLUMN vs, ve", accepts: true, want: []dropItem{
		{"period", "", false, "PERIOD FOR SYSTEM_TIME"},
		{"column", "vs", false, "vs"},
		{"column", "ve", false, "ve"},
	}},
	{name: "constraint_then_period", sql: "ALTER TABLE dbo.t DROP CONSTRAINT ck_n, PERIOD FOR SYSTEM_TIME", accepts: true, want: []dropItem{
		{"constraint", "ck_n", false, "ck_n"},
		{"period", "", false, "PERIOD FOR SYSTEM_TIME"},
	}},
	{name: "period_inherited", sql: "ALTER TABLE dbo.t DROP PERIOD FOR SYSTEM_TIME, ck_n WITH (ONLINE = ON)", accepts: true, want: []dropItem{
		{"period", "", false, "PERIOD FOR SYSTEM_TIME"},
		{"period", "ck_n", false, "ck_n WITH (ONLINE = ON)"},
	}},
	{name: "period_as_constraint_name", sql: "ALTER TABLE dbo.t DROP PERIOD", accepts: true, want: []dropItem{
		{"constraint", "PERIOD", false, "PERIOD"},
	}},
	{name: "period_as_column_name", sql: "ALTER TABLE dbo.t DROP COLUMN a, PERIOD", accepts: true, want: []dropItem{
		{"column", "a", false, "a"},
		{"column", "PERIOD", false, "PERIOD"},
	}},
	{name: "period_with", sql: "ALTER TABLE dbo.t DROP PERIOD FOR SYSTEM_TIME WITH (ONLINE = ON)", near: "("},
	{name: "period_without_system_time", sql: "ALTER TABLE dbo.t DROP CONSTRAINT ck_n, PERIOD FOR", near: "end of input"},
	{name: "period_without_for", sql: "ALTER TABLE dbo.t DROP PERIOD SYSTEM_TIME", near: "SYSTEM_TIME"},

	// Every item needs a name.
	{name: "no_item", sql: "ALTER TABLE dbo.t DROP", near: "end of input"},
	{name: "constraint_without_name", sql: "ALTER TABLE dbo.t DROP CONSTRAINT", near: "end of input"},
	{name: "column_without_name", sql: "ALTER TABLE dbo.t DROP COLUMN", near: "end of input"},
	{name: "trailing_comma", sql: "ALTER TABLE dbo.t DROP COLUMN a,", near: "end of input"},
	{name: "restated_column_without_name", sql: "ALTER TABLE dbo.t DROP CONSTRAINT ck_n, COLUMN", near: "end of input"},
	{name: "restated_constraint_without_name", sql: "ALTER TABLE dbo.t DROP COLUMN a, CONSTRAINT", near: "end of input"},
	{name: "empty_item", sql: "ALTER TABLE dbo.t DROP COLUMN a, , b", near: ","},
	{name: "keyword_twice", sql: "ALTER TABLE dbo.t DROP COLUMN COLUMN a", near: "COLUMN"},
	{name: "two_keywords", sql: "ALTER TABLE dbo.t DROP CONSTRAINT COLUMN a", near: "COLUMN"},
}

// dropItems renders the ALTER TABLE ... DROP actions of the statement sql
// parses to.
func dropItems(t *testing.T, sql string, list *ast.List) []dropItem {
	t.Helper()
	if list == nil || len(list.Items) != 1 {
		t.Fatalf("Parse(%q) = %v, want one statement", sql, list)
	}
	stmt, ok := list.Items[0].(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("Parse(%q) = %T, want *ast.AlterTableStmt", sql, list.Items[0])
	}
	var got []dropItem
	for _, item := range stmt.Actions.Items {
		a := item.(*ast.AlterTableAction)
		d := dropItem{ifExists: a.IfExists, text: sql[a.Loc.Start:a.Loc.End]}
		switch a.Type {
		case ast.ATDropColumn:
			d.kind, d.name = "column", a.ColName
		case ast.ATDropConstraint:
			d.kind, d.name = "constraint", a.Constraint.Name
			if c := sql[a.Constraint.Loc.Start:a.Constraint.Loc.End]; !strings.HasPrefix(d.text, c) || strings.Contains(c, " WITH") {
				t.Errorf("Parse(%q): constraint Loc spans %q, want the name at the start of %q", sql, c, d.text)
			}
		case ast.ATDropPeriod:
			d.kind = "period"
			if a.Names != nil {
				d.name = a.Names.Items[0].(*ast.String).Str
			}
		default:
			t.Fatalf("Parse(%q): action type %d is not a drop", sql, a.Type)
		}
		got = append(got, d)
	}
	return got
}

// TestAlterTableDropList checks the actions omni builds for an ALTER TABLE
// ... DROP list, and where it rejects one.
func TestAlterTableDropList(t *testing.T) {
	for _, tc := range dropListCases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if !tc.accepts {
				if err == nil {
					t.Fatalf("Parse(%q) succeeded, want a syntax error", tc.sql)
				}
				want := fmt.Sprintf("at or near %q", tc.near)
				if tc.near == "end of input" {
					want = "at end of input"
				}
				if !strings.Contains(err.Error(), want) {
					t.Errorf("Parse(%q) error = %q, want it %s", tc.sql, err.Error(), want)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			if got := dropItems(t, tc.sql, list); !slices.Equal(got, tc.want) {
				t.Errorf("Parse(%q) actions:\n got %+v\nwant %+v", tc.sql, got, tc.want)
			}
		})
	}
}

// TestAlterTableDropListOracle checks accept/reject parity with SQL Server
// for every ALTER TABLE ... DROP list in dropListCases.
func TestAlterTableDropListOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range dropListCases {
		t.Run(tc.name, func(t *testing.T) {
			ssErr, err := oracle.parseError(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}
			ssAccepts := ssErr == nil
			if ssAccepts != tc.accepts {
				t.Errorf("table says SQL Server %s but it %s (%v):\n  SQL: %s",
					boolToAcceptReject(tc.accepts), boolToAcceptReject(ssAccepts), ssErr, tc.sql)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil
			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s (%v), omni %s (%v)\n  SQL: %s",
					boolToAcceptReject(ssAccepts), ssErr, boolToAcceptReject(omniAccepts), omniErr, tc.sql)
			}
		})
	}
}

// dropProbeTable is the table TestAlterTableDropListEffectOracle drops from:
// columns, three named constraints, and a SYSTEM_TIME period.
const dropProbeTable = `CREATE TABLE dbo.drop_probe (
	id INT CONSTRAINT pk_probe PRIMARY KEY,
	a INT,
	b INT,
	c INT CONSTRAINT ck_probe CHECK (c > 0),
	d INT CONSTRAINT df_probe DEFAULT 1,
	vs DATETIME2 GENERATED ALWAYS AS ROW START NOT NULL,
	ve DATETIME2 GENERATED ALWAYS AS ROW END NOT NULL,
	PERIOD FOR SYSTEM_TIME (vs, ve)
)`

// dropProbeObjects lists what dbo.drop_probe has: "column x",
// "constraint x", and "period".
const dropProbeObjects = `SELECT 'column ' + name FROM sys.columns WHERE object_id = OBJECT_ID('dbo.drop_probe')
UNION ALL SELECT 'constraint ' + name FROM sys.objects WHERE parent_object_id = OBJECT_ID('dbo.drop_probe') AND type IN ('PK', 'C', 'D')
UNION ALL SELECT 'period' FROM sys.periods WHERE object_id = OBJECT_ID('dbo.drop_probe')`

// TestAlterTableDropListEffectOracle runs ALTER TABLE ... DROP lists on SQL
// Server and checks that omni's actions say what SQL Server does: which kind
// a name without a keyword drops, and which items IF EXISTS covers. From the
// actions it predicts what the statement drops, or that it fails on a name
// that does not exist, and compares that with the outcome. Each statement
// runs in a transaction on a fresh dbo.drop_probe and is rolled back.
func TestAlterTableDropListEffectOracle(t *testing.T) {
	oracle := startParserOracle(t)

	tests := []struct {
		sql   string
		fails bool // SQL Server fails: an item names something the table lacks
	}{
		{sql: "ALTER TABLE dbo.drop_probe DROP COLUMN a, b"},
		{sql: "ALTER TABLE dbo.drop_probe DROP ck_probe, COLUMN a"},
		{sql: "ALTER TABLE dbo.drop_probe DROP CONSTRAINT ck_probe, df_probe, COLUMN c, d"},
		{sql: "ALTER TABLE dbo.drop_probe DROP COLUMN a, CONSTRAINT ck_probe, df_probe"},
		// ck_probe follows COLUMN a, so it names a column.
		{sql: "ALTER TABLE dbo.drop_probe DROP COLUMN a, ck_probe", fails: true},
		{sql: "ALTER TABLE dbo.drop_probe DROP COLUMN IF EXISTS nope, b"},
		{sql: "ALTER TABLE dbo.drop_probe DROP CONSTRAINT IF EXISTS nope1, CONSTRAINT IF EXISTS ck_probe, COLUMN IF EXISTS nope2"},
		// IF EXISTS covers a, not nope.
		{sql: "ALTER TABLE dbo.drop_probe DROP COLUMN IF EXISTS a, nope", fails: true},
		{sql: "ALTER TABLE dbo.drop_probe DROP CONSTRAINT IF EXISTS ck_probe, nope", fails: true},
		{sql: "ALTER TABLE dbo.drop_probe DROP PERIOD FOR SYSTEM_TIME, COLUMN vs, ve"},
		// ck_probe follows PERIOD FOR SYSTEM_TIME, so it names a period.
		{sql: "ALTER TABLE dbo.drop_probe DROP PERIOD FOR SYSTEM_TIME, ck_probe", fails: true},
	}

	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			items := dropItems(t, tc.sql, list)

			before, gone, ssErr, err := oracle.dropEffect(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}
			if (ssErr != nil) != tc.fails {
				t.Errorf("table says fails=%v, SQL Server returned %v\n  SQL: %s", tc.fails, ssErr, tc.sql)
			}

			// What omni's actions say the statement does to the table.
			has := make(map[string]bool, len(before))
			for _, o := range before {
				has[o] = true
			}
			var want []string
			wantFail := false
			for _, it := range items {
				o := it.kind + " " + it.name
				if it.kind == "period" && it.name == "" {
					o = "period"
				}
				switch {
				case has[o]:
					has[o] = false
					want = append(want, o)
				case !it.ifExists:
					wantFail = true
				}
			}
			if wantFail != (ssErr != nil) {
				t.Errorf("omni's actions %+v say fails=%v, SQL Server returned %v\n  SQL: %s", items, wantFail, ssErr, tc.sql)
			}
			if ssErr != nil {
				return
			}
			slices.Sort(want)
			if !slices.Equal(gone, want) {
				t.Errorf("SQL Server dropped %v, omni's actions %+v drop %v\n  SQL: %s", gone, items, want, tc.sql)
			}
		})
	}
}

// dropEffect creates dbo.drop_probe, runs stmt, and rolls both back. It
// returns what the table had, what stmt dropped (sorted), and stmt's error.
// The last result reports a failure of the oracle itself.
func (o *parserOracle) dropEffect(stmt string) (before, gone []string, stmtErr, err error) {
	tx, err := o.conn.BeginTx(o.ctx, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("begin: %w", err)
	}
	// A failed DROP may end the transaction on the server; rolling back
	// what is left is all that matters here.
	defer tx.Rollback() //nolint:errcheck

	if _, err := tx.ExecContext(o.ctx, dropProbeTable); err != nil {
		return nil, nil, nil, fmt.Errorf("create dbo.drop_probe: %w", err)
	}
	objects := func() ([]string, error) {
		rows, err := tx.QueryContext(o.ctx, dropProbeObjects)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var out []string
		for rows.Next() {
			var s string
			if err := rows.Scan(&s); err != nil {
				return nil, err
			}
			out = append(out, s)
		}
		return out, rows.Err()
	}
	if before, err = objects(); err != nil {
		return nil, nil, nil, fmt.Errorf("list objects: %w", err)
	}
	if _, stmtErr = tx.ExecContext(o.ctx, stmt); stmtErr != nil {
		return before, nil, stmtErr, nil
	}
	after, err := objects()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("list objects: %w", err)
	}
	for _, s := range before {
		if !slices.Contains(after, s) {
			gone = append(gone, s)
		}
	}
	slices.Sort(gone)
	return before, gone, nil, nil
}
