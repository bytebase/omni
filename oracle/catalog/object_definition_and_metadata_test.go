package catalog

import (
	"context"
	"fmt"
	"slices"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
	"github.com/stretchr/testify/require"
)

func TestGetTableDefinition(t *testing.T) {
	tests := []struct {
		name  string
		table *metadatapb.TableMetadata
		want  string
	}{
		{
			// A table whose last column has DEFAULT + NOT NULL, followed by a
			// table-level FOREIGN KEY with no other indexes or constraints.
			// The FK must be separated from the column by a comma (ORA-02253).
			name: "foreign key after default not null column without other constraints",
			table: &metadatapb.TableMetadata{
				Name: "T1",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "C1", Type: "NUMBER", Nullable: false},
					{Name: "C2", Type: "NUMBER(1)", Default: "0", Nullable: false},
				},
				ForeignKeys: []*metadatapb.ForeignKeyMetadata{
					{
						Name:              "T1_FK",
						Columns:           []string{"C1"},
						ReferencedTable:   "T2",
						ReferencedColumns: []string{"C1"},
					},
				},
			},
			want: `CREATE TABLE "T1" (
  "C1" NUMBER NOT NULL,
  "C2" NUMBER(1) DEFAULT 0 NOT NULL,
  CONSTRAINT "T1_FK" FOREIGN KEY ("C1") REFERENCES "T2" ("C1")
);

`,
		},
		{
			name: "check constraint after column without other constraints",
			table: &metadatapb.TableMetadata{
				Name: "T1",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "C1", Type: "NUMBER", Nullable: false},
				},
				CheckConstraints: []*metadatapb.CheckConstraintMetadata{
					{Name: "T1_CK", Expression: `"C1" > 0`},
				},
			},
			want: `CREATE TABLE "T1" (
  "C1" NUMBER NOT NULL,
  CONSTRAINT "T1_CK" CHECK ("C1" > 0)
);

`,
		},
		{
			name: "primary key, check, and foreign key",
			table: &metadatapb.TableMetadata{
				Name: "T1",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "C1", Type: "NUMBER", Nullable: false},
					{Name: "C2", Type: "NUMBER(1)", Default: "0", Nullable: false},
				},
				Indexes: []*metadatapb.IndexMetadata{
					{
						Name:         "T1_PK",
						Expressions:  []string{"C1"},
						Primary:      true,
						Unique:       true,
						IsConstraint: true,
					},
				},
				CheckConstraints: []*metadatapb.CheckConstraintMetadata{
					{Name: "T1_CK", Expression: `"C2" IN (0, 1)`},
				},
				ForeignKeys: []*metadatapb.ForeignKeyMetadata{
					{
						Name:              "T1_FK",
						Columns:           []string{"C1"},
						ReferencedTable:   "T2",
						ReferencedColumns: []string{"C1"},
					},
				},
			},
			want: `CREATE TABLE "T1" (
  "C1" NUMBER NOT NULL,
  "C2" NUMBER(1) DEFAULT 0 NOT NULL,
  CONSTRAINT "T1_PK" PRIMARY KEY ("C1"),
  CONSTRAINT "T1_CK" CHECK ("C2" IN (0, 1)),
  CONSTRAINT "T1_FK" FOREIGN KEY ("C1") REFERENCES "T2" ("C1")
);

`,
		},
		{
			name: "object table has no definition",
			table: &metadatapb.TableMetadata{
				Name:     "OT",
				SkipDump: true,
				Columns:  []*metadatapb.ColumnMetadata{{Name: "X", Type: "NUMBER", Nullable: true}},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTableDefinition("", tt.table, nil)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetObjectDefinition(t *testing.T) {
	tests := []struct {
		name string
		get  func() (string, error)
		want string
	}{
		{
			name: "view",
			get: func() (string, error) {
				return GetViewDefinition("", &metadatapb.ViewMetadata{
					Name:       "DEPT_EMPLOYEE_COUNT",
					Definition: "SELECT D.ID AS DEPT_ID, COUNT(E.ID) AS EMP_COUNT\nFROM DEPARTMENTS D LEFT JOIN EMPLOYEES E ON E.DEPT_ID = D.ID\nGROUP BY D.ID",
				})
			},
			want: `CREATE VIEW "DEPT_EMPLOYEE_COUNT" AS SELECT D.ID AS DEPT_ID, COUNT(E.ID) AS EMP_COUNT
FROM DEPARTMENTS D LEFT JOIN EMPLOYEES E ON E.DEPT_ID = D.ID
GROUP BY D.ID;

`,
		},
		{
			name: "materialized view",
			get: func() (string, error) {
				return GetMaterializedViewDefinition("", &metadatapb.MaterializedViewMetadata{
					Name:       "PRODUCT_STATS",
					Definition: "SELECT PRODUCT_ID, COUNT(*) AS ORDER_COUNT FROM ORDERS GROUP BY PRODUCT_ID",
				})
			},
			want: `CREATE MATERIALIZED VIEW "PRODUCT_STATS" AS SELECT PRODUCT_ID, COUNT(*) AS ORDER_COUNT FROM ORDERS GROUP BY PRODUCT_ID;

`,
		},
		{
			name: "function gains the CREATE OR REPLACE prefix",
			get: func() (string, error) {
				return GetFunctionDefinition("", &metadatapb.FunctionMetadata{
					Name:       "CALCULATE_DISCOUNT",
					Definition: "FUNCTION CALCULATE_DISCOUNT(AMOUNT NUMBER)\nRETURN NUMBER\nIS\nBEGIN\n    RETURN AMOUNT * 0.1;\nEND;",
				})
			},
			want: `CREATE OR REPLACE FUNCTION CALCULATE_DISCOUNT(AMOUNT NUMBER)
RETURN NUMBER
IS
BEGIN
    RETURN AMOUNT * 0.1;
END;

`,
		},
		{
			name: "procedure gains the CREATE OR REPLACE prefix",
			get: func() (string, error) {
				return GetProcedureDefinition("", &metadatapb.ProcedureMetadata{
					Name:       "LOG_AUDIT",
					Definition: "PROCEDURE LOG_AUDIT(P_TABLE_NAME VARCHAR2)\nIS\nBEGIN\n    NULL;\nEND;",
				})
			},
			want: `CREATE OR REPLACE PROCEDURE LOG_AUDIT(P_TABLE_NAME VARCHAR2)
IS
BEGIN
    NULL;
END;

`,
		},
		{
			name: "view keeps its INSTEAD OF trigger",
			get: func() (string, error) {
				return GetViewDefinition("", &metadatapb.ViewMetadata{
					Name:       "EMPLOYEE_VIEW",
					Definition: "SELECT EMP_ID, EMAIL FROM EMPLOYEES",
					Triggers: []*metadatapb.TriggerMetadata{
						{
							Name: "EMPLOYEE_VIEW_INSERT_TRG",
							Body: "CREATE OR REPLACE TRIGGER EMPLOYEE_VIEW_INSERT_TRG\nINSTEAD OF INSERT ON EMPLOYEE_VIEW\nFOR EACH ROW\nBEGIN\n    NULL;\nEND;",
						},
					},
				})
			},
			want: `CREATE VIEW "EMPLOYEE_VIEW" AS SELECT EMP_ID, EMAIL FROM EMPLOYEES;

CREATE OR REPLACE TRIGGER EMPLOYEE_VIEW_INSERT_TRG
INSTEAD OF INSERT ON EMPLOYEE_VIEW
FOR EACH ROW
BEGIN
    NULL;
END;

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.get()
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestObjectDefinitionsRunOnOracle executes the single-object definitions
// produced by GetTableDefinition, GetViewDefinition, GetMaterializedViewDefinition,
// GetFunctionDefinition, and GetProcedureDefinition against the pinned Oracle
// server. The tape replays it in CI; record with OMNI_RECORD=1.
func TestObjectDefinitionsRunOnOracle(t *testing.T) {
	db := testcontainer.Open(t, "system")
	defer db.Close()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	require.NoError(t, err)
	defer conn.Close()

	const schema = "OMNI_OBJ_DEF"
	execScript(t, conn, fmt.Sprintf(`CREATE USER "%s" IDENTIFIED BY "%s" QUOTA UNLIMITED ON USERS`, schema, testcontainer.Password))
	execScript(t, conn, fmt.Sprintf(`GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE PROCEDURE, CREATE TRIGGER, CREATE MATERIALIZED VIEW TO "%s"`, schema))
	execScript(t, conn, fmt.Sprintf(`ALTER SESSION SET CURRENT_SCHEMA = "%s"`, schema))

	// Base tables referenced by the single-object definitions below.
	execScript(t, conn, `
CREATE TABLE "T2" ("C1" NUMBER PRIMARY KEY);
CREATE TABLE "DEPARTMENTS" ("ID" NUMBER PRIMARY KEY);
CREATE TABLE "EMPLOYEES" ("ID" NUMBER PRIMARY KEY, "EMP_ID" NUMBER, "DEPT_ID" NUMBER, "EMAIL" VARCHAR2(100));
CREATE TABLE "ORDERS" ("ID" NUMBER PRIMARY KEY, "PRODUCT_ID" NUMBER);
`)

	tableDDL, err := GetTableDefinition("", &metadatapb.TableMetadata{
		Name: "T1",
		Columns: []*metadatapb.ColumnMetadata{
			{Name: "C1", Type: "NUMBER", Nullable: false},
			{Name: "C2", Type: "NUMBER(1)", Default: "0", Nullable: false},
		},
		Indexes: []*metadatapb.IndexMetadata{
			{Name: "T1_PK", Expressions: []string{"C1"}, Primary: true, Unique: true, IsConstraint: true},
		},
		CheckConstraints: []*metadatapb.CheckConstraintMetadata{
			{Name: "T1_CK", Expression: `"C2" IN (0, 1)`},
		},
		ForeignKeys: []*metadatapb.ForeignKeyMetadata{
			{Name: "T1_FK", Columns: []string{"C1"}, ReferencedTable: "T2", ReferencedColumns: []string{"C1"}},
		},
	}, nil)
	require.NoError(t, err)
	execScript(t, conn, tableDDL)

	viewDDL, err := GetViewDefinition("", &metadatapb.ViewMetadata{
		Name:       "DEPT_EMPLOYEE_COUNT",
		Definition: "SELECT D.ID AS DEPT_ID, COUNT(E.ID) AS EMP_COUNT\nFROM DEPARTMENTS D LEFT JOIN EMPLOYEES E ON E.DEPT_ID = D.ID\nGROUP BY D.ID",
	})
	require.NoError(t, err)
	execScript(t, conn, viewDDL)

	triggerViewDDL, err := GetViewDefinition("", &metadatapb.ViewMetadata{
		Name:       "EMPLOYEE_VIEW",
		Definition: "SELECT EMP_ID, EMAIL FROM EMPLOYEES",
		Triggers: []*metadatapb.TriggerMetadata{{
			Name: "EMPLOYEE_VIEW_INSERT_TRG",
			Body: "CREATE OR REPLACE TRIGGER EMPLOYEE_VIEW_INSERT_TRG\nINSTEAD OF INSERT ON EMPLOYEE_VIEW\nFOR EACH ROW\nBEGIN\n    NULL;\nEND;",
		}},
	})
	require.NoError(t, err)
	execScript(t, conn, triggerViewDDL)

	mvDDL, err := GetMaterializedViewDefinition("", &metadatapb.MaterializedViewMetadata{
		Name:       "PRODUCT_STATS",
		Definition: "SELECT PRODUCT_ID, COUNT(*) AS ORDER_COUNT FROM ORDERS GROUP BY PRODUCT_ID",
	})
	require.NoError(t, err)
	execScript(t, conn, mvDDL)

	fnDDL, err := GetFunctionDefinition("", &metadatapb.FunctionMetadata{
		Name:       "CALCULATE_DISCOUNT",
		Definition: "FUNCTION CALCULATE_DISCOUNT(AMOUNT NUMBER)\nRETURN NUMBER\nIS\nBEGIN\n    RETURN AMOUNT * 0.1;\nEND;",
	})
	require.NoError(t, err)
	execScript(t, conn, fnDDL)

	procDDL, err := GetProcedureDefinition("", &metadatapb.ProcedureMetadata{
		Name:       "LOG_AUDIT",
		Definition: "PROCEDURE LOG_AUDIT(P_TABLE_NAME VARCHAR2)\nIS\nBEGIN\n    NULL;\nEND;",
	})
	require.NoError(t, err)
	execScript(t, conn, procDDL)

	got := readObjects(t, conn, schema)
	for _, want := range []string{
		"TABLE/T1",
		"INDEX/T1_PK",
		"VIEW/DEPT_EMPLOYEE_COUNT",
		"VIEW/EMPLOYEE_VIEW",
		"TRIGGER/EMPLOYEE_VIEW_INSERT_TRG",
		"MATERIALIZED VIEW/PRODUCT_STATS",
		"FUNCTION/CALCULATE_DISCOUNT",
		"PROCEDURE/LOG_AUDIT",
	} {
		require.True(t, slices.Contains(got, want), "expected %s in %v", want, got)
	}
}

func TestGetDatabaseMetadataObjectKinds(t *testing.T) {
	tests := []struct {
		name   string
		ddl    string
		verify func(*testing.T, *metadatapb.DatabaseSchemaMetadata)
	}{
		{
			name: "invisible_and_virtual_columns",
			ddl: `
CREATE TABLE T (
    A NUMBER,
    H NUMBER INVISIBLE DEFAULT 7,
    IHV INVISIBLE AS (A * 2) VIRTUAL
);
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				table := requireOracleTable(t, requireSingleOracleSchema(t, metadata), "T")
				require.Len(t, table.Columns, 3)
				require.False(t, table.Columns[0].IsInvisible)
				require.True(t, table.Columns[1].IsInvisible)
				require.Equal(t, "7", table.Columns[1].Default)
				require.True(t, table.Columns[2].IsInvisible)
				require.Equal(t, "A * 2", table.Columns[2].GetGeneration().GetExpression())
			},
		},
		{
			name: "table_constraints_indexes_view_sequence",
			ddl: `
CREATE TABLE DEPARTMENTS (
    DEPT_ID NUMBER PRIMARY KEY,
    DEPT_NAME VARCHAR2(100) NOT NULL UNIQUE,
    MANAGER_ID NUMBER,
    CONSTRAINT chk_dept_name CHECK (DEPT_NAME IS NOT NULL)
);

CREATE TABLE EMPLOYEES (
    EMP_ID NUMBER,
    DEPT_ID NUMBER,
    EMAIL VARCHAR2(100) UNIQUE,
    SALARY NUMBER(10,2) DEFAULT 0,
    CONSTRAINT pk_employees PRIMARY KEY (EMP_ID),
    CONSTRAINT fk_emp_dept FOREIGN KEY (DEPT_ID) REFERENCES DEPARTMENTS(DEPT_ID) ON DELETE SET NULL,
    CONSTRAINT chk_salary CHECK (SALARY >= 0)
);

CREATE INDEX idx_emp_dept ON EMPLOYEES(DEPT_ID);

CREATE VIEW ACTIVE_EMPLOYEES AS
SELECT EMP_ID, DEPT_ID, EMAIL
FROM EMPLOYEES
WHERE EMAIL IS NOT NULL;

CREATE SEQUENCE emp_seq START WITH 1 INCREMENT BY 1;
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				schemaMetadata := requireSingleOracleSchema(t, metadata)

				departments := requireOracleTable(t, schemaMetadata, "DEPARTMENTS")
				require.Len(t, departments.Columns, 3)
				require.Equal(t, "NUMBER", departments.Columns[0].Type)
				require.False(t, departments.Columns[0].Nullable)
				require.Equal(t, "VARCHAR2(100 BYTE)", departments.Columns[1].Type)
				require.False(t, departments.Columns[1].Nullable)
				requireOracleIndex(t, departments, "PK_DEPARTMENTS", []string{"DEPT_ID"}, true, true)
				requireOracleIndex(t, departments, "UK_DEPARTMENTS_DEPT_NAME", []string{"DEPT_NAME"}, false, true)
				requireOracleCheckConstraint(t, departments, "CHK_DEPT_NAME", "DEPT_NAME IS NOT NULL")

				employees := requireOracleTable(t, schemaMetadata, "EMPLOYEES")
				require.Len(t, employees.Columns, 4)
				require.Equal(t, "0", employees.Columns[3].Default)
				requireOracleIndex(t, employees, "PK_EMPLOYEES", []string{"EMP_ID"}, true, true)
				requireOracleIndex(t, employees, "UK_EMPLOYEES_EMAIL", []string{"EMAIL"}, false, true)
				requireOracleIndex(t, employees, "IDX_EMP_DEPT", []string{"DEPT_ID"}, false, false)
				requireOracleCheckConstraint(t, employees, "CHK_SALARY", "SALARY >= 0")
				requireOracleForeignKey(t, employees, "FK_EMP_DEPT", []string{"DEPT_ID"}, "DEPARTMENTS", []string{"DEPT_ID"}, "SET NULL")

				require.Len(t, schemaMetadata.Views, 1)
				require.Equal(t, "ACTIVE_EMPLOYEES", schemaMetadata.Views[0].Name)
				require.Contains(t, schemaMetadata.Views[0].Definition, "SELECT EMP_ID, DEPT_ID, EMAIL")

				require.Len(t, schemaMetadata.Sequences, 1)
				require.Equal(t, "EMP_SEQ", schemaMetadata.Sequences[0].Name)
				require.Equal(t, "1", schemaMetadata.Sequences[0].Start)
			},
		},
		{
			name: "create_schema_wrapper",
			ddl: `
CREATE SCHEMA AUTHORIZATION HR
CREATE TABLE EMPLOYEES (
    EMP_ID NUMBER PRIMARY KEY,
    EMAIL VARCHAR2(100)
)
CREATE VIEW EMPLOYEE_VIEW AS
SELECT EMP_ID, EMAIL FROM EMPLOYEES;
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				schemaMetadata := requireSingleOracleSchema(t, metadata)
				require.Equal(t, "HR", schemaMetadata.Name)
				employees := requireOracleTable(t, schemaMetadata, "EMPLOYEES")
				require.Len(t, employees.Columns, 2)
				requireOracleIndex(t, employees, "PK_EMPLOYEES", []string{"EMP_ID"}, true, true)
				view := requireOracleView(t, schemaMetadata, "EMPLOYEE_VIEW")
				require.Contains(t, view.Definition, "SELECT EMP_ID, EMAIL FROM EMPLOYEES")
			},
		},
		{
			name: "comments_on_views_and_materialized_views",
			ddl: `
CREATE TABLE PRODUCT_SALES (
    PRODUCT_ID NUMBER NOT NULL,
    CATEGORY VARCHAR2(50) NOT NULL,
    SALES_AMOUNT NUMBER(12,2) NOT NULL
);

CREATE VIEW PRODUCT_SALES_VIEW AS
SELECT PRODUCT_ID, CATEGORY, SALES_AMOUNT
FROM PRODUCT_SALES;

CREATE MATERIALIZED VIEW PRODUCT_SALES_MV
BUILD IMMEDIATE
REFRESH COMPLETE ON DEMAND
AS
SELECT PRODUCT_ID, CATEGORY, SUM(SALES_AMOUNT) AS TOTAL_REVENUE
FROM PRODUCT_SALES
GROUP BY PRODUCT_ID, CATEGORY;

COMMENT ON TABLE PRODUCT_SALES_VIEW IS 'Product sales view';
COMMENT ON MATERIALIZED VIEW PRODUCT_SALES_MV IS 'Product sales materialized view';
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				schemaMetadata := requireSingleOracleSchema(t, metadata)
				require.Equal(t, "Product sales view", requireOracleView(t, schemaMetadata, "PRODUCT_SALES_VIEW").Comment)
				require.Len(t, schemaMetadata.MaterializedViews, 1)
				require.Equal(t, "PRODUCT_SALES_MV", schemaMetadata.MaterializedViews[0].Name)
				require.Equal(t, "Product sales materialized view", schemaMetadata.MaterializedViews[0].Comment)
			},
		},
		{
			name: "index_visibility_and_function_based_bitmap",
			ddl: `
CREATE TABLE ORDERS (
    ID NUMBER PRIMARY KEY,
    STATUS VARCHAR2(20)
);

CREATE INDEX idx_orders_status_invisible ON ORDERS(STATUS) INVISIBLE;
CREATE BITMAP INDEX idx_orders_lower_status ON ORDERS(LOWER(STATUS));
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				orders := requireOracleTable(t, requireSingleOracleSchema(t, metadata), "ORDERS")
				invisibleIndex := requireOracleIndexMetadata(t, orders, "IDX_ORDERS_STATUS_INVISIBLE")
				require.False(t, invisibleIndex.Visible)
				functionBasedBitmapIndex := requireOracleIndexMetadata(t, orders, "IDX_ORDERS_LOWER_STATUS")
				require.Equal(t, "FUNCTION-BASED BITMAP", functionBasedBitmapIndex.Type)
				require.True(t, functionBasedBitmapIndex.Visible)
			},
		},
		{
			name: "instead_of_trigger_on_view",
			ddl: `
CREATE TABLE EMPLOYEES (
    EMP_ID NUMBER PRIMARY KEY,
    EMAIL VARCHAR2(100)
);

CREATE VIEW EMPLOYEE_VIEW AS
SELECT EMP_ID, EMAIL
FROM EMPLOYEES;

CREATE OR REPLACE TRIGGER employee_view_insert_trg
INSTEAD OF INSERT ON EMPLOYEE_VIEW
FOR EACH ROW
BEGIN
    NULL;
END;
/
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				schemaMetadata := requireSingleOracleSchema(t, metadata)
				view := requireOracleView(t, schemaMetadata, "EMPLOYEE_VIEW")
				require.Len(t, view.Triggers, 1)
				require.Equal(t, "EMPLOYEE_VIEW_INSERT_TRG", view.Triggers[0].Name)
				require.Nil(t, findOracleTable(schemaMetadata, "EMPLOYEE_VIEW"))
			},
		},
		{
			name: "materialized_view_with_index",
			ddl: `
CREATE TABLE PRODUCT_SALES (
    PRODUCT_ID NUMBER NOT NULL,
    CATEGORY VARCHAR2(50) NOT NULL,
    SALES_AMOUNT NUMBER(12,2) NOT NULL
);

CREATE MATERIALIZED VIEW PRODUCT_SALES_MV
BUILD IMMEDIATE
REFRESH COMPLETE ON DEMAND
AS
SELECT PRODUCT_ID, CATEGORY, SUM(SALES_AMOUNT) AS TOTAL_REVENUE
FROM PRODUCT_SALES
GROUP BY PRODUCT_ID, CATEGORY;

CREATE INDEX idx_mv_category ON PRODUCT_SALES_MV(CATEGORY);
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				schemaMetadata := requireSingleOracleSchema(t, metadata)
				requireOracleTable(t, schemaMetadata, "PRODUCT_SALES")
				require.Len(t, schemaMetadata.MaterializedViews, 1)
				require.Equal(t, "PRODUCT_SALES_MV", schemaMetadata.MaterializedViews[0].Name)
				require.Contains(t, schemaMetadata.MaterializedViews[0].Definition, "SUM(SALES_AMOUNT)")
				requireOracleMaterializedViewIndex(t, schemaMetadata.MaterializedViews[0], "IDX_MV_CATEGORY", []string{"CATEGORY"}, false, false)
				require.Nil(t, findOracleTable(schemaMetadata, "PRODUCT_SALES_MV"))
			},
		},
		{
			name: "alter_table_add_foreign_key",
			ddl: `
CREATE TABLE DEPARTMENTS (
    DEPT_ID NUMBER PRIMARY KEY,
    MANAGER_ID NUMBER
);

CREATE TABLE EMPLOYEES (
    EMP_ID NUMBER PRIMARY KEY
);

ALTER TABLE DEPARTMENTS ADD CONSTRAINT fk_dept_manager
    FOREIGN KEY (MANAGER_ID) REFERENCES EMPLOYEES(EMP_ID) ON DELETE SET NULL;
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				schemaMetadata := requireSingleOracleSchema(t, metadata)
				departments := requireOracleTable(t, schemaMetadata, "DEPARTMENTS")
				requireOracleForeignKey(t, departments, "FK_DEPT_MANAGER", []string{"MANAGER_ID"}, "EMPLOYEES", []string{"EMP_ID"}, "SET NULL")
			},
		},
		{
			name: "alter_table_add_column",
			ddl: `
CREATE TABLE EMPLOYEES (
    EMP_ID NUMBER PRIMARY KEY
);

ALTER TABLE EMPLOYEES ADD (
    EMAIL VARCHAR2(100) NOT NULL,
    STATUS VARCHAR2(20) DEFAULT 'ACTIVE'
);
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				employees := requireOracleTable(t, requireSingleOracleSchema(t, metadata), "EMPLOYEES")
				require.Len(t, employees.Columns, 3)
				require.Equal(t, "EMAIL", employees.Columns[1].Name)
				require.Equal(t, "VARCHAR2(100 BYTE)", employees.Columns[1].Type)
				require.False(t, employees.Columns[1].Nullable)
				require.Equal(t, "STATUS", employees.Columns[2].Name)
				require.Equal(t, "VARCHAR2(20 BYTE)", employees.Columns[2].Type)
				require.Equal(t, "'ACTIVE'", employees.Columns[2].Default)
			},
		},
		{
			name: "inline_constraints_and_omitted_reference_columns",
			ddl: `
CREATE TABLE DEPARTMENTS (
    DEPT_ID NUMBER PRIMARY KEY
);

CREATE TABLE EMPLOYEES (
    EMP_ID NUMBER PRIMARY KEY,
    SALARY NUMBER CHECK (SALARY > 0) CHECK (SALARY < 1000000),
    DEPT_ID NUMBER REFERENCES DEPARTMENTS,
    MANAGER_DEPT_ID NUMBER,
    CONSTRAINT fk_manager_dept FOREIGN KEY (MANAGER_DEPT_ID) REFERENCES DEPARTMENTS
);
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				employees := requireOracleTable(t, requireSingleOracleSchema(t, metadata), "EMPLOYEES")
				requireOracleCheckConstraint(t, employees, "CHK_EMPLOYEES_SALARY", "SALARY > 0")
				requireOracleCheckConstraint(t, employees, "CHK_EMPLOYEES_SALARY_2", "SALARY < 1000000")
				requireOracleForeignKey(t, employees, "FK_EMPLOYEES_DEPT_ID", []string{"DEPT_ID"}, "DEPARTMENTS", []string{"DEPT_ID"}, "")
				requireOracleForeignKey(t, employees, "FK_MANAGER_DEPT", []string{"MANAGER_DEPT_ID"}, "DEPARTMENTS", []string{"DEPT_ID"}, "")
			},
		},
		{
			name: "virtual_column",
			ddl: `
CREATE TABLE ORDER_ITEMS (
    QTY NUMBER,
    PRICE NUMBER,
    TOTAL AS (QTY * PRICE)
);
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				orderItems := requireOracleTable(t, requireSingleOracleSchema(t, metadata), "ORDER_ITEMS")
				require.Len(t, orderItems.Columns, 3)
				require.Equal(t, "TOTAL", orderItems.Columns[2].Name)
				require.Empty(t, orderItems.Columns[2].Type)
				require.Empty(t, orderItems.Columns[2].Default)
				require.Equal(t, metadatapb.GenerationMetadata_TYPE_VIRTUAL, orderItems.Columns[2].GetGeneration().GetType())
				require.Equal(t, "QTY * PRICE", orderItems.Columns[2].GetGeneration().GetExpression())
			},
		},
		{
			name: "escaped_string_literals_in_expressions",
			ddl: `
CREATE TABLE AUTHORS (
    NAME VARCHAR2(100),
    IS_REILLY AS (CASE WHEN NAME = 'O''Reilly' THEN 1 ELSE 0 END),
    CONSTRAINT chk_author_name CHECK (NAME <> 'O''Reilly')
);
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				authors := requireOracleTable(t, requireSingleOracleSchema(t, metadata), "AUTHORS")
				require.Len(t, authors.Columns, 2)
				require.Equal(t, "CASE WHEN NAME = 'O''Reilly' THEN 1 ELSE 0 END", authors.Columns[1].GetGeneration().GetExpression())
				requireOracleCheckConstraint(t, authors, "CHK_AUTHOR_NAME", "NAME <> 'O''Reilly'")
			},
		},
		{
			name: "default_on_null_column",
			ddl: `
CREATE TABLE ORDERS (
    ID NUMBER PRIMARY KEY,
    STATUS VARCHAR2(20) DEFAULT ON NULL 'PENDING'
);
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				orders := requireOracleTable(t, requireSingleOracleSchema(t, metadata), "ORDERS")
				require.Len(t, orders.Columns, 2)
				require.Equal(t, "STATUS", orders.Columns[1].Name)
				require.Equal(t, "'PENDING'", orders.Columns[1].Default)
				require.True(t, orders.Columns[1].DefaultOnNull)
			},
		},
		{
			name: "package_metadata",
			ddl: `
CREATE OR REPLACE PACKAGE financial_utils AS
    c_default_currency CONSTANT VARCHAR2(3) := 'USD';
    FUNCTION format_currency(p_amount NUMBER) RETURN VARCHAR2;
END financial_utils;
/

CREATE OR REPLACE PACKAGE BODY financial_utils AS
    FUNCTION format_currency(p_amount NUMBER) RETURN VARCHAR2 IS
    BEGIN
        RETURN TO_CHAR(p_amount, 'FM9999990.00');
    END format_currency;
END financial_utils;
/
`,
			verify: func(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) {
				schemaMetadata := requireSingleOracleSchema(t, metadata)
				require.Len(t, schemaMetadata.Packages, 1)
				require.Equal(t, "FINANCIAL_UTILS", schemaMetadata.Packages[0].Name)
				require.Contains(t, schemaMetadata.Packages[0].Definition, "c_default_currency")
				require.Contains(t, schemaMetadata.Packages[0].Definition, "FUNCTION format_currency")
				require.Contains(t, schemaMetadata.Packages[0].Definition, "PACKAGE BODY financial_utils")
				require.Contains(t, schemaMetadata.Packages[0].Definition, "RETURN TO_CHAR")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			metadata, err := GetDatabaseMetadata(tc.ddl)
			require.NoError(t, err)
			tc.verify(t, metadata)
		})
	}
}

func requireSingleOracleSchema(t *testing.T, metadata *metadatapb.DatabaseSchemaMetadata) *metadatapb.SchemaMetadata {
	t.Helper()
	require.Len(t, metadata.Schemas, 1)
	return metadata.Schemas[0]
}

func requireOracleTable(t *testing.T, schemaMetadata *metadatapb.SchemaMetadata, name string) *metadatapb.TableMetadata {
	t.Helper()
	table := findOracleTable(schemaMetadata, name)
	require.NotNil(t, table)
	return table
}

func requireOracleView(t *testing.T, schemaMetadata *metadatapb.SchemaMetadata, name string) *metadatapb.ViewMetadata {
	t.Helper()
	for _, view := range schemaMetadata.Views {
		if view.Name == name {
			return view
		}
	}
	t.Fatalf("view %q not found", name)
	return nil
}

func findOracleTable(schemaMetadata *metadatapb.SchemaMetadata, name string) *metadatapb.TableMetadata {
	for _, table := range schemaMetadata.Tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

func requireOracleIndex(t *testing.T, table *metadatapb.TableMetadata, name string, expressions []string, primary, unique bool) {
	t.Helper()
	index := requireOracleIndexMetadata(t, table, name)
	require.Equal(t, expressions, index.Expressions)
	require.Equal(t, primary, index.Primary)
	require.Equal(t, unique, index.Unique)
}

func requireOracleIndexMetadata(t *testing.T, table *metadatapb.TableMetadata, name string) *metadatapb.IndexMetadata {
	t.Helper()
	for _, index := range table.Indexes {
		if index.Name == name {
			return index
		}
	}
	t.Fatalf("index %q not found in table %q", name, table.Name)
	return nil
}

func requireOracleMaterializedViewIndex(t *testing.T, mv *metadatapb.MaterializedViewMetadata, name string, expressions []string, primary, unique bool) {
	t.Helper()
	for _, index := range mv.Indexes {
		if index.Name != name {
			continue
		}
		require.Equal(t, expressions, index.Expressions)
		require.Equal(t, primary, index.Primary)
		require.Equal(t, unique, index.Unique)
		return
	}
	t.Fatalf("index %q not found in materialized view %q", name, mv.Name)
}

func requireOracleCheckConstraint(t *testing.T, table *metadatapb.TableMetadata, name, expression string) {
	t.Helper()
	for _, check := range table.CheckConstraints {
		if check.Name != name {
			continue
		}
		require.Equal(t, expression, check.Expression)
		return
	}
	t.Fatalf("check constraint %q not found in table %q", name, table.Name)
}

func requireOracleForeignKey(t *testing.T, table *metadatapb.TableMetadata, name string, columns []string, referencedTable string, referencedColumns []string, onDelete string) {
	t.Helper()
	for _, fk := range table.ForeignKeys {
		if fk.Name != name {
			continue
		}
		require.Equal(t, columns, fk.Columns)
		require.Equal(t, referencedTable, fk.ReferencedTable)
		require.Equal(t, referencedColumns, fk.ReferencedColumns)
		require.Equal(t, onDelete, fk.OnDelete)
		return
	}
	t.Fatalf("foreign key %q not found in table %q", name, table.Name)
}
