package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestP2CreateIndexPartitionAndStorageOptions(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		key   string
		value string
	}{
		{
			name:  "local_partitioned_index",
			sql:   "CREATE INDEX ix_sales ON sales (sale_date) LOCAL (PARTITION p1 TABLESPACE ts1) TABLESPACE users",
			key:   "LOCAL",
			value: "( PARTITION P1 TABLESPACE TS1 )",
		},
		{
			name:  "global_partitioned_index",
			sql:   "CREATE INDEX ix_sales ON sales (sale_date) GLOBAL PARTITION BY RANGE (sale_date) (PARTITION p1 VALUES LESS THAN (10)) ONLINE",
			key:   "GLOBAL",
			value: "PARTITION BY RANGE ( SALE_DATE ) ( PARTITION P1 VALUES LESS THAN ( 10 ) )",
		},
		{
			name:  "storage_clause",
			sql:   "CREATE INDEX ix_sales ON sales (sale_date) STORAGE (INITIAL 64K NEXT 64K) ONLINE",
			key:   "STORAGE",
			value: "( INITIAL 64 K NEXT 64 K )",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := parseCreateIndexForP2(t, tt.sql)
			assertCreateIndexOption(t, stmt, tt.key, tt.value)
		})
	}
}

func TestP2CreateIndexOptionBoundariesAndLoc(t *testing.T) {
	stmt := parseCreateIndexForP2(t, "CREATE INDEX ix_sales ON sales (sale_date) GLOBAL PARTITION BY HASH (sale_date) PARTITIONS 4 ONLINE")
	if !stmt.Global || !stmt.Online {
		t.Fatalf("expected GLOBAL and ONLINE fields, got global=%v online=%v", stmt.Global, stmt.Online)
	}
	opt := findCreateIndexOption(stmt, "GLOBAL")
	if opt == nil {
		t.Fatalf("expected GLOBAL option in %#v", stmt.Options)
	}
	if opt.Value != "PARTITION BY HASH ( SALE_DATE ) PARTITIONS 4" {
		t.Fatalf("GLOBAL value=%q", opt.Value)
	}
	if opt.Loc.Start <= stmt.Loc.Start || opt.Loc.End > stmt.Loc.End || opt.Loc.Start >= opt.Loc.End {
		t.Fatalf("option Loc=%+v is not inside stmt Loc=%+v", opt.Loc, stmt.Loc)
	}
}

func TestP2CreateIndexMalformedStorageOption(t *testing.T) {
	ParseShouldFail(t, "CREATE INDEX ix_sales ON sales (sale_date) STORAGE")
}

func parseCreateIndexForP2(t *testing.T, sql string) *ast.CreateIndexStmt {
	t.Helper()
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreateIndexStmt)
	if !ok {
		t.Fatalf("expected CreateIndexStmt, got %T", raw.Stmt)
	}
	return stmt
}

func assertCreateIndexOption(t *testing.T, stmt *ast.CreateIndexStmt, key, value string) {
	t.Helper()
	opt := findCreateIndexOption(stmt, key)
	if opt == nil {
		t.Fatalf("expected option %q in %#v", key, stmt.Options)
	}
	if opt.Value != value {
		t.Fatalf("option %q value=%q, want %q", key, opt.Value, value)
	}
}

func findCreateIndexOption(stmt *ast.CreateIndexStmt, key string) *ast.DDLOption {
	if stmt.Options == nil {
		return nil
	}
	for _, item := range stmt.Options.Items {
		opt, ok := item.(*ast.DDLOption)
		if ok && opt.Key == key {
			return opt
		}
	}
	return nil
}

func TestParseCreateIndex(t *testing.T) {
	p := newTestParser("INDEX idx_emp_name ON employees (last_name)")
	stmt, parseErr1 := p.parseCreateIndexStmt(0)
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	if stmt == nil {
		t.Fatal("expected CreateIndexStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "IDX_EMP_NAME" {
		t.Errorf("expected index name IDX_EMP_NAME, got %v", stmt.Name)
	}
	if stmt.Table == nil || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected table EMPLOYEES, got %v", stmt.Table)
	}
	if stmt.Columns == nil || stmt.Columns.Len() != 1 {
		t.Fatalf("expected 1 column, got %d", stmt.Columns.Len())
	}
	col0 := stmt.Columns.Items[0].(*ast.IndexColumn)
	cr := col0.Expr.(*ast.ColumnRef)
	if cr.Column != "LAST_NAME" {
		t.Errorf("expected column LAST_NAME, got %q", cr.Column)
	}
}

func TestParseCreateUniqueIndex(t *testing.T) {
	p := newTestParser("UNIQUE INDEX idx_emp_id ON hr.employees (employee_id)")
	stmt, parseErr2 := p.parseCreateIndexStmt(0)
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	if !stmt.Unique {
		t.Error("expected Unique to be true")
	}
	if stmt.Name == nil || stmt.Name.Name != "IDX_EMP_ID" {
		t.Errorf("expected index name IDX_EMP_ID, got %v", stmt.Name)
	}
	if stmt.Table == nil || stmt.Table.Schema != "HR" || stmt.Table.Name != "EMPLOYEES" {
		t.Errorf("expected table HR.EMPLOYEES, got %v", stmt.Table)
	}
}

func TestParseCreateBitmapIndex(t *testing.T) {
	p := newTestParser("BITMAP INDEX idx_status ON orders (status)")
	stmt, parseErr3 := p.parseCreateIndexStmt(0)
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	if !stmt.Bitmap {
		t.Error("expected Bitmap to be true")
	}
}

func TestParseCreateIndexMultiColumn(t *testing.T) {
	p := newTestParser("INDEX idx_multi ON t (a ASC, b DESC)")
	stmt, parseErr4 := p.parseCreateIndexStmt(0)
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	if stmt.Columns == nil || stmt.Columns.Len() != 2 {
		t.Fatalf("expected 2 columns, got %d", stmt.Columns.Len())
	}
	col0 := stmt.Columns.Items[0].(*ast.IndexColumn)
	if col0.Dir != ast.SORTBY_ASC {
		t.Errorf("expected ASC for col0, got %d", col0.Dir)
	}
	col1 := stmt.Columns.Items[1].(*ast.IndexColumn)
	if col1.Dir != ast.SORTBY_DESC {
		t.Errorf("expected DESC for col1, got %d", col1.Dir)
	}
}

func TestParseCreateIndexReverse(t *testing.T) {
	p := newTestParser("INDEX idx_rev ON t (a) REVERSE")
	stmt, parseErr5 := p.parseCreateIndexStmt(0)
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	if !stmt.Reverse {
		t.Error("expected Reverse to be true")
	}
}

func TestParseCreateIndexTablespace(t *testing.T) {
	p := newTestParser("INDEX idx_ts ON t (a) TABLESPACE users")
	stmt, parseErr6 := p.parseCreateIndexStmt(0)
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	if stmt.Tablespace != "USERS" {
		t.Errorf("expected tablespace USERS, got %q", stmt.Tablespace)
	}
}

func TestParseCreateIndexLocal(t *testing.T) {
	p := newTestParser("INDEX idx_local ON t (a) LOCAL")
	stmt, parseErr7 := p.parseCreateIndexStmt(0)
	if parseErr7 != nil {
		t.Fatalf("parse: %v", parseErr7)
	}
	if !stmt.Local {
		t.Error("expected Local to be true")
	}
}

func TestParseCreateIndexGlobal(t *testing.T) {
	p := newTestParser("INDEX idx_global ON t (a) GLOBAL")
	stmt, parseErr8 := p.parseCreateIndexStmt(0)
	if parseErr8 != nil {
		t.Fatalf("parse: %v", parseErr8)
	}
	if !stmt.Global {
		t.Error("expected Global to be true")
	}
}

func TestParseCreateIndexOnline(t *testing.T) {
	p := newTestParser("INDEX idx_online ON t (a) ONLINE")
	stmt, parseErr9 := p.parseCreateIndexStmt(0)
	if parseErr9 != nil {
		t.Fatalf("parse: %v", parseErr9)
	}
	if !stmt.Online {
		t.Error("expected Online to be true")
	}
}

func TestParseCreateIndexLoc(t *testing.T) {
	p := newTestParser("INDEX idx ON t (a)")
	stmt, parseErr10 := p.parseCreateIndexStmt(0)
	if parseErr10 != nil {
		t.Fatalf("parse: %v", parseErr10)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}
