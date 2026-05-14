package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestP2CreateTableComplexOptions(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		key   string
		value string
	}{
		{
			name:  "external_table_clause",
			sql:   "CREATE TABLE ext_t (id NUMBER) ORGANIZATION EXTERNAL (TYPE ORACLE_LOADER DEFAULT DIRECTORY data_dir ACCESS PARAMETERS (FIELDS TERMINATED BY ',')) REJECT LIMIT UNLIMITED",
			key:   "ORGANIZATION EXTERNAL",
			value: "( TYPE ORACLE_LOADER DEFAULT DIRECTORY DATA_DIR ACCESS PARAMETERS ( FIELDS TERMINATED BY , ) )",
		},
		{
			name:  "lob_storage_clause",
			sql:   "CREATE TABLE docs (id NUMBER, doc CLOB) LOB (doc) STORE AS SECUREFILE (TABLESPACE lob_ts)",
			key:   "LOB",
			value: "( DOC ) STORE AS SECUREFILE ( TABLESPACE LOB_TS )",
		},
		{
			name:  "inmemory_clause",
			sql:   "CREATE TABLE hot_t (id NUMBER) INMEMORY MEMCOMPRESS FOR QUERY HIGH PRIORITY HIGH",
			key:   "INMEMORY",
			value: "MEMCOMPRESS FOR QUERY HIGH PRIORITY HIGH",
		},
		{
			name:  "ilm_clause",
			sql:   "CREATE TABLE lifecycle_t (id NUMBER) ILM ADD POLICY ROW STORE COMPRESS ADVANCED",
			key:   "ILM",
			value: "ADD POLICY ROW STORE COMPRESS ADVANCED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := parseCreateTableForP2(t, tt.sql)
			assertCreateTableOption(t, stmt, tt.key, tt.value)
		})
	}
}

func TestP2CreateTableOptionLoc(t *testing.T) {
	stmt := parseCreateTableForP2(t, "CREATE TABLE docs (id NUMBER, doc CLOB) LOB (doc) STORE AS (TABLESPACE lob_ts)")
	opt := findCreateTableOption(stmt, "LOB")
	if opt == nil {
		t.Fatalf("expected LOB option in %#v", stmt.Options)
	}
	if opt.Loc.Start <= stmt.Loc.Start || opt.Loc.End > stmt.Loc.End || opt.Loc.Start >= opt.Loc.End {
		t.Fatalf("option Loc=%+v is not inside stmt Loc=%+v", opt.Loc, stmt.Loc)
	}
}

func TestP2CreateTableMalformedComplexOption(t *testing.T) {
	ParseShouldFail(t, "CREATE TABLE docs (id NUMBER, doc CLOB) LOB STORE AS lob_seg")
}

func TestCreateTableIntervalPartitioning(t *testing.T) {
	ParseAndCheck(t, `CREATE TABLE t_interval_full (d DATE)
PARTITION BY RANGE (d)
INTERVAL (NUMTOYMINTERVAL(1,'MONTH'))
(
  PARTITION p1 VALUES LESS THAN (TO_DATE('2012-01-01', 'YYYY-MM-DD'))
)`)
}

func TestCreateTableIntervalPartitioningRequiresRangePartition(t *testing.T) {
	ParseShouldFail(t, "CREATE TABLE t_interval_min (d DATE) PARTITION BY RANGE (d) INTERVAL (NUMTOYMINTERVAL(1,'MONTH'))")
}

func TestCreateTableDateLiteralPartitionBound(t *testing.T) {
	ParseAndCheck(t, "CREATE TABLE t_date_bound (d DATE) PARTITION BY RANGE (d) (PARTITION p1 VALUES LESS THAN (DATE '2020-01-01'))")
}

func TestCreateTablePartitionStorageClause(t *testing.T) {
	tests := []string{
		`CREATE TABLE t_part_storage_initial (n NUMBER)
PARTITION BY RANGE (n)
(
  PARTITION p1 VALUES LESS THAN (100)
    STORAGE (INITIAL 8388608)
)`,
		`CREATE TABLE t_part_storage_full (n NUMBER)
PARTITION BY RANGE (n)
(
  PARTITION p1 VALUES LESS THAN (100)
    STORAGE (
      INITIAL 8388608
      NEXT 1048576
      MINEXTENTS 1
      MAXEXTENTS 2147483645
      BUFFER_POOL DEFAULT
    )
)`,
		`CREATE TABLE t_part_storage_attrs (txn_date DATE)
ROW STORE COMPRESS ADVANCED
TABLESPACE users
PCTFREE 10
NOLOGGING
PARTITION BY RANGE (txn_date)
INTERVAL (NUMTOYMINTERVAL(1,'MONTH'))
(
  PARTITION part_01 VALUES LESS THAN (DATE '2024-01-01')
    NOLOGGING
    COMPRESS FOR OLTP
    TABLESPACE users
    PCTFREE 10
    STORAGE (
      INITIAL 8388608
      NEXT 1048576
      MINEXTENTS 1
      MAXEXTENTS 2147483645
      BUFFER_POOL DEFAULT
    )
)`,
	}

	for _, sql := range tests {
		ParseAndCheck(t, sql)
	}
}

func parseCreateTableForP2(t *testing.T, sql string) *ast.CreateTableStmt {
	t.Helper()
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	stmt, ok := raw.Stmt.(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("expected CreateTableStmt, got %T", raw.Stmt)
	}
	return stmt
}

func assertCreateTableOption(t *testing.T, stmt *ast.CreateTableStmt, key, value string) {
	t.Helper()
	opt := findCreateTableOption(stmt, key)
	if opt == nil {
		t.Fatalf("expected option %q in %#v", key, stmt.Options)
	}
	if opt.Value != value {
		t.Fatalf("option %q value=%q, want %q", key, opt.Value, value)
	}
}

func findCreateTableOption(stmt *ast.CreateTableStmt, key string) *ast.DDLOption {
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

// TestParseCreateTableSimple tests a basic CREATE TABLE with columns.
func TestParseCreateTableSimple(t *testing.T) {
	sql := `CREATE TABLE employees (
		id NUMBER(10),
		name VARCHAR2(100),
		salary NUMBER(12,2)
	)`
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
	raw := result.Items[0].(*ast.RawStmt)
	ct, ok := raw.Stmt.(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("expected CreateTableStmt, got %T", raw.Stmt)
	}
	if ct.Name.Name != "EMPLOYEES" {
		t.Errorf("expected table name EMPLOYEES, got %q", ct.Name.Name)
	}
	if ct.Columns.Len() != 3 {
		t.Fatalf("expected 3 columns, got %d", ct.Columns.Len())
	}

	col0 := ct.Columns.Items[0].(*ast.ColumnDef)
	if col0.Name != "ID" {
		t.Errorf("expected column name ID, got %q", col0.Name)
	}
	col1 := ct.Columns.Items[1].(*ast.ColumnDef)
	if col1.Name != "NAME" {
		t.Errorf("expected column name NAME, got %q", col1.Name)
	}
	col2 := ct.Columns.Items[2].(*ast.ColumnDef)
	if col2.Name != "SALARY" {
		t.Errorf("expected column name SALARY, got %q", col2.Name)
	}
}

// TestParseCreateTableSchemaQualified tests CREATE TABLE with schema prefix.
func TestParseCreateTableSchemaQualified(t *testing.T) {
	sql := `CREATE TABLE hr.employees (id NUMBER)`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)
	if ct.Name.Schema != "HR" {
		t.Errorf("expected schema HR, got %q", ct.Name.Schema)
	}
	if ct.Name.Name != "EMPLOYEES" {
		t.Errorf("expected table name EMPLOYEES, got %q", ct.Name.Name)
	}
}

// TestParseCreateTableNotNull tests CREATE TABLE with NOT NULL constraints.
func TestParseCreateTableNotNull(t *testing.T) {
	sql := `CREATE TABLE t (
		id NUMBER NOT NULL,
		name VARCHAR2(50) NULL
	)`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)
	col0 := ct.Columns.Items[0].(*ast.ColumnDef)
	if !col0.NotNull {
		t.Error("expected id to be NOT NULL")
	}
	col1 := ct.Columns.Items[1].(*ast.ColumnDef)
	if !col1.Null {
		t.Error("expected name to be explicitly NULL")
	}
}

// TestParseCreateTableColumnConstraints tests column-level constraints.
func TestParseCreateTableColumnConstraints(t *testing.T) {
	sql := `CREATE TABLE t (
		id NUMBER PRIMARY KEY,
		email VARCHAR2(100) UNIQUE,
		age NUMBER CHECK (age > 0),
		dept_id NUMBER REFERENCES departments(id)
	)`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)

	if ct.Columns.Len() != 4 {
		t.Fatalf("expected 4 columns, got %d", ct.Columns.Len())
	}

	// id PRIMARY KEY
	col0 := ct.Columns.Items[0].(*ast.ColumnDef)
	if col0.Constraints == nil || col0.Constraints.Len() != 1 {
		t.Fatalf("expected 1 constraint on id, got %d", col0.Constraints.Len())
	}
	cc0 := col0.Constraints.Items[0].(*ast.ColumnConstraint)
	if cc0.Type != ast.CONSTRAINT_PRIMARY {
		t.Errorf("expected PRIMARY constraint, got %d", cc0.Type)
	}

	// email UNIQUE
	col1 := ct.Columns.Items[1].(*ast.ColumnDef)
	cc1 := col1.Constraints.Items[0].(*ast.ColumnConstraint)
	if cc1.Type != ast.CONSTRAINT_UNIQUE {
		t.Errorf("expected UNIQUE constraint, got %d", cc1.Type)
	}

	// age CHECK
	col2 := ct.Columns.Items[2].(*ast.ColumnDef)
	cc2 := col2.Constraints.Items[0].(*ast.ColumnConstraint)
	if cc2.Type != ast.CONSTRAINT_CHECK {
		t.Errorf("expected CHECK constraint, got %d", cc2.Type)
	}
	if cc2.Expr == nil {
		t.Error("expected CHECK expression, got nil")
	}

	// dept_id REFERENCES
	col3 := ct.Columns.Items[3].(*ast.ColumnDef)
	cc3 := col3.Constraints.Items[0].(*ast.ColumnConstraint)
	if cc3.Type != ast.CONSTRAINT_FOREIGN {
		t.Errorf("expected FOREIGN constraint, got %d", cc3.Type)
	}
	if cc3.RefTable == nil || cc3.RefTable.Name != "DEPARTMENTS" {
		t.Errorf("expected RefTable DEPARTMENTS, got %v", cc3.RefTable)
	}
	if cc3.RefColumns == nil || cc3.RefColumns.Len() != 1 {
		t.Fatal("expected 1 ref column")
	}
}

// TestParseCreateTableNamedConstraint tests named column constraints.
func TestParseCreateTableNamedConstraint(t *testing.T) {
	sql := `CREATE TABLE t (
		id NUMBER CONSTRAINT pk_t PRIMARY KEY
	)`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)
	col0 := ct.Columns.Items[0].(*ast.ColumnDef)
	cc := col0.Constraints.Items[0].(*ast.ColumnConstraint)
	if cc.Name != "PK_T" {
		t.Errorf("expected constraint name PK_T, got %q", cc.Name)
	}
	if cc.Type != ast.CONSTRAINT_PRIMARY {
		t.Errorf("expected PRIMARY constraint, got %d", cc.Type)
	}
}

// TestParseCreateTableDefaults tests column defaults.
func TestParseCreateTableDefaults(t *testing.T) {
	sql := `CREATE TABLE t (
		status VARCHAR2(10) DEFAULT 'ACTIVE',
		created_at DATE DEFAULT SYSDATE,
		counter NUMBER DEFAULT 0
	)`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)

	col0 := ct.Columns.Items[0].(*ast.ColumnDef)
	if col0.Default == nil {
		t.Error("expected default for status")
	}
	col1 := ct.Columns.Items[1].(*ast.ColumnDef)
	if col1.Default == nil {
		t.Error("expected default for created_at")
	}
	col2 := ct.Columns.Items[2].(*ast.ColumnDef)
	if col2.Default == nil {
		t.Error("expected default for counter")
	}
}

// TestParseCreateTableAsSelect tests CREATE TABLE AS SELECT (CTAS).
func TestParseCreateTableAsSelect(t *testing.T) {
	sql := `CREATE TABLE t2 AS SELECT id, name FROM t1 WHERE id > 10`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)
	if ct.Name.Name != "T2" {
		t.Errorf("expected table name T2, got %q", ct.Name.Name)
	}
	if ct.AsQuery == nil {
		t.Fatal("expected AsQuery (CTAS)")
	}
	sel, ok := ct.AsQuery.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt, got %T", ct.AsQuery)
	}
	if sel.TargetList.Len() != 2 {
		t.Errorf("expected 2 select targets, got %d", sel.TargetList.Len())
	}
}

// TestParseCreateGlobalTemporaryTable tests CREATE GLOBAL TEMPORARY TABLE.
func TestParseCreateGlobalTemporaryTable(t *testing.T) {
	sql := `CREATE GLOBAL TEMPORARY TABLE session_data (
		id NUMBER,
		data VARCHAR2(100)
	) ON COMMIT PRESERVE ROWS`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)
	if !ct.Global {
		t.Error("expected Global flag to be true")
	}
	if ct.OnCommit != "PRESERVE ROWS" {
		t.Errorf("expected ON COMMIT PRESERVE ROWS, got %q", ct.OnCommit)
	}
}

// TestParseCreateGlobalTempDeleteRows tests ON COMMIT DELETE ROWS.
func TestParseCreateGlobalTempDeleteRows(t *testing.T) {
	sql := `CREATE GLOBAL TEMPORARY TABLE tmp (id NUMBER) ON COMMIT DELETE ROWS`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)
	if !ct.Global {
		t.Error("expected Global flag to be true")
	}
	if ct.OnCommit != "DELETE ROWS" {
		t.Errorf("expected ON COMMIT DELETE ROWS, got %q", ct.OnCommit)
	}
}

// TestParseCreateTableTableConstraints tests table-level constraints.
func TestParseCreateTableTableConstraints(t *testing.T) {
	sql := `CREATE TABLE orders (
		id NUMBER,
		customer_id NUMBER,
		order_date DATE,
		CONSTRAINT pk_orders PRIMARY KEY (id),
		CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers(id),
		CONSTRAINT chk_date CHECK (order_date IS NOT NULL),
		UNIQUE (customer_id, order_date)
	)`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)

	if ct.Columns.Len() != 3 {
		t.Fatalf("expected 3 columns, got %d", ct.Columns.Len())
	}
	if ct.Constraints == nil || ct.Constraints.Len() != 4 {
		t.Fatalf("expected 4 table constraints, got %d", ct.Constraints.Len())
	}

	// PRIMARY KEY
	tc0 := ct.Constraints.Items[0].(*ast.TableConstraint)
	if tc0.Name != "PK_ORDERS" {
		t.Errorf("expected constraint name PK_ORDERS, got %q", tc0.Name)
	}
	if tc0.Type != ast.CONSTRAINT_PRIMARY {
		t.Errorf("expected PRIMARY constraint, got %d", tc0.Type)
	}
	if tc0.Columns.Len() != 1 {
		t.Errorf("expected 1 column in PK, got %d", tc0.Columns.Len())
	}

	// FOREIGN KEY
	tc1 := ct.Constraints.Items[1].(*ast.TableConstraint)
	if tc1.Name != "FK_CUSTOMER" {
		t.Errorf("expected constraint name FK_CUSTOMER, got %q", tc1.Name)
	}
	if tc1.Type != ast.CONSTRAINT_FOREIGN {
		t.Errorf("expected FOREIGN constraint, got %d", tc1.Type)
	}
	if tc1.RefTable == nil || tc1.RefTable.Name != "CUSTOMERS" {
		t.Errorf("expected RefTable CUSTOMERS, got %v", tc1.RefTable)
	}

	// CHECK
	tc2 := ct.Constraints.Items[2].(*ast.TableConstraint)
	if tc2.Type != ast.CONSTRAINT_CHECK {
		t.Errorf("expected CHECK constraint, got %d", tc2.Type)
	}
	if tc2.Expr == nil {
		t.Error("expected CHECK expression")
	}

	// UNIQUE (unnamed)
	tc3 := ct.Constraints.Items[3].(*ast.TableConstraint)
	if tc3.Type != ast.CONSTRAINT_UNIQUE {
		t.Errorf("expected UNIQUE constraint, got %d", tc3.Type)
	}
	if tc3.Columns.Len() != 2 {
		t.Errorf("expected 2 columns in UNIQUE, got %d", tc3.Columns.Len())
	}
}

// TestParseCreateTableOrReplace tests CREATE OR REPLACE TABLE (23c).
func TestParseCreateTableOrReplace(t *testing.T) {
	sql := `CREATE OR REPLACE TABLE t (id NUMBER)`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)
	if !ct.OrReplace {
		t.Error("expected OrReplace to be true")
	}
}

// TestParseCreateTableMultipleStatements tests CREATE TABLE followed by another statement.
func TestParseCreateTableMultipleStatements(t *testing.T) {
	sql := `CREATE TABLE t (id NUMBER); SELECT 1 FROM dual`
	result := ParseAndCheck(t, sql)
	if result.Len() != 2 {
		t.Fatalf("expected 2 statements, got %d", result.Len())
	}
	raw0 := result.Items[0].(*ast.RawStmt)
	if _, ok := raw0.Stmt.(*ast.CreateTableStmt); !ok {
		t.Errorf("expected CreateTableStmt, got %T", raw0.Stmt)
	}
	raw1 := result.Items[1].(*ast.RawStmt)
	if _, ok := raw1.Stmt.(*ast.SelectStmt); !ok {
		t.Errorf("expected SelectStmt, got %T", raw1.Stmt)
	}
}

// TestParseCreateTablePrivateTemporary tests CREATE PRIVATE TEMPORARY TABLE.
func TestParseCreateTablePrivateTemporary(t *testing.T) {
	sql := `CREATE PRIVATE TEMPORARY TABLE ora$ptt_tmp (id NUMBER) ON COMMIT PRESERVE ROWS`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)
	if !ct.Private {
		t.Error("expected Private flag to be true")
	}
}

// TestParseCreateTableTablespace tests CREATE TABLE with TABLESPACE.
func TestParseCreateTableTablespace(t *testing.T) {
	sql := `CREATE TABLE t (id NUMBER) TABLESPACE users`
	result := ParseAndCheck(t, sql)
	raw := result.Items[0].(*ast.RawStmt)
	ct := raw.Stmt.(*ast.CreateTableStmt)
	if ct.Tablespace != "USERS" {
		t.Errorf("expected tablespace USERS, got %q", ct.Tablespace)
	}
}
