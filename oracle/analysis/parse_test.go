package analysis

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

func TestParse(t *testing.T) {
	list, err := Parse("SELECT 1 FROM DUAL;\nSELECT A FROM T")
	require.NoError(t, err)
	require.Len(t, list.Items, 2)
	second, ok := list.Items[1].(*ast.RawStmt)
	require.True(t, ok)
	require.Equal(t, 20, second.Stmt.(*ast.SelectStmt).Loc.Start)

	list, err = Parse("-- only a comment")
	require.NoError(t, err)
	require.Empty(t, list.Items)

	_, err = Parse("SELECT 1 FROM DUAL;\nSELEC 2")
	var parseErr *parser.ParseError
	require.ErrorAs(t, err, &parseErr)
	require.GreaterOrEqual(t, parseErr.Position, 20)
}

func TestParseLocatesEveryNodeInTheScript(t *testing.T) {
	// Each statement's own nodes, reached by ast.Walk or not, lie inside it:
	// a DELETE's target, an UNPIVOT's value column, TABLE()'s argument, a
	// multi-table INSERT's WHEN condition.
	for _, statement := range []string{
		"DELETE a WHERE b = 1",
		"SELECT * FROM a UNPIVOT (v FOR k IN (c))",
		"SELECT * FROM TABLE(1)",
		"INSERT WHEN 1 = 1 THEN INTO a INTO b SELECT 1 FROM dual",
	} {
		script := "SELECT 1 FROM DUAL;\n" + statement
		list, err := Parse(script)
		require.NoError(t, err, statement)
		require.Len(t, list.Items, 2, statement)
		start := len("SELECT 1 FROM DUAL;\n")
		for _, loc := range locsOf(reflect.ValueOf(list.Items[1]), map[uintptr]bool{}) {
			if loc.Start < 0 {
				continue
			}
			require.GreaterOrEqual(t, loc.Start, start, statement)
			require.LessOrEqual(t, loc.End, len(script), statement)
		}
	}
}

// locsOf returns every Loc in a tree, following each field, not only those
// ast.Walk visits.
func locsOf(v reflect.Value, seen map[uintptr]bool) []ast.Loc {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return nil
		}
		if v.Kind() == reflect.Pointer {
			if seen[v.Pointer()] {
				return nil
			}
			seen[v.Pointer()] = true
		}
		return locsOf(v.Elem(), seen)
	case reflect.Struct:
		if loc, ok := v.Interface().(ast.Loc); ok {
			return []ast.Loc{loc}
		}
		var locs []ast.Loc
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() {
				locs = append(locs, locsOf(v.Field(i), seen)...)
			}
		}
		return locs
	case reflect.Slice:
		var locs []ast.Loc
		for i := 0; i < v.Len(); i++ {
			locs = append(locs, locsOf(v.Index(i), seen)...)
		}
		return locs
	default:
		return nil
	}
}

func TestParseAdditionalScriptsAndProbes(t *testing.T) {
	t.Run("slash terminated PL/SQL script", func(t *testing.T) {
		list, err := Parse(`
CREATE TABLE AUDIT_LOG (
  LOG_ID NUMBER PRIMARY KEY
);

CREATE OR REPLACE TRIGGER audit_log_trigger
BEFORE INSERT ON AUDIT_LOG
FOR EACH ROW
BEGIN
  NULL;
END;
/
`)
		require.NoError(t, err)
		require.Len(t, list.Items, 2)
		require.IsType(t, &ast.CreateTableStmt{}, list.Items[0].(*ast.RawStmt).Stmt)
		require.IsType(t, &ast.CreateTriggerStmt{}, list.Items[1].(*ast.RawStmt).Stmt)
	})

	t.Run("skips SQL*Plus commands", func(t *testing.T) {
		list, err := Parse(`
SET DEFINE OFF
PROMPT setup

CREATE TABLE AUDIT_LOG (
  LOG_ID NUMBER PRIMARY KEY
);

SPOOL out.log
CREATE OR REPLACE TRIGGER audit_log_trigger
BEFORE INSERT ON AUDIT_LOG
FOR EACH ROW
BEGIN
  NULL;
END;
/
SPOOL OFF
`)
		require.NoError(t, err)
		require.Len(t, list.Items, 2)
	})

	t.Run("keeps remark column in CREATE TABLE", func(t *testing.T) {
		list, err := Parse(`
CREATE TABLE parser_regression_remark_column (
  id                  NUMBER NOT NULL
      CONSTRAINT parser_regression_remark_column_pk
          PRIMARY KEY,
  display_name        VARCHAR2(100),
  source              VARCHAR2(100),
  remark              VARCHAR2(100)
)`)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		require.IsType(t, &ast.CreateTableStmt{}, list.Items[0].(*ast.RawStmt).Stmt)
	})

	t.Run("trigger referencing aliases", func(t *testing.T) {
		list, err := Parse(`CREATE OR REPLACE TRIGGER trg
BEFORE INSERT OR UPDATE OF col1, col2 ON tbl
REFERENCING OLD AS o NEW AS n
FOR EACH ROW
WHEN (n.col1 > 0)
BEGIN
  :n.col2 := :o.col2 + 1;
END;`)
		require.NoError(t, err)
		require.Len(t, list.Items, 1)
		trigger, ok := list.Items[0].(*ast.RawStmt).Stmt.(*ast.CreateTriggerStmt)
		require.True(t, ok)
		require.NotNil(t, trigger.Referencing)
		require.Equal(t, "O", trigger.Referencing.OldAlias)
		require.Equal(t, "N", trigger.Referencing.NewAlias)
	})

	t.Run("interval partition with DATE literal and q-quote", func(t *testing.T) {
		for _, stmt := range []string{
			`UPDATE t1 SET (c1) = 1 WHERE c2 = 2;`,
			`SELECT q'\This is String\' FROM DUAL;`,
			`CREATE TABLE APP.EVENTS (TXN_DATE DATE, STATUS NUMBER)
PARTITION BY RANGE (TXN_DATE)
INTERVAL (NUMTODSINTERVAL(1,'DAY'))
(PARTITION P0 VALUES LESS THAN (DATE '2026-01-01'));`,
		} {
			list, err := Parse(stmt)
			require.NoError(t, err, stmt)
			require.NotEmpty(t, list.Items, stmt)
		}
	})
}

func TestParseQuerySpanNodeProbes(t *testing.T) {
	tests := []struct {
		name          string
		statement     string
		expectedNodes []string
	}{
		{"select columns and table ref", "SELECT A, T.B, PUBLIC.T.C FROM T", []string{"SelectStmt", "ResTarget", "ColumnRef", "TableRef"}},
		{"join using", "SELECT * FROM T T1 JOIN T T2 USING(A)", []string{"SelectStmt", "JoinClause", "TableRef"}},
		{"derived table", "SELECT * FROM (SELECT A FROM T) DT", []string{"SelectStmt", "SubqueryRef", "TableRef"}},
		{"scalar and in subqueries", "SELECT (SELECT MAX(A) FROM T) AS M FROM T WHERE B IN (SELECT C FROM T)", []string{"SelectStmt", "SubqueryExpr", "FuncCallExpr", "InExpr"}},
		{"cte and set operation", "WITH T1(D, C) AS (SELECT A, B FROM T UNION ALL SELECT C, D FROM T) SELECT * FROM T1", []string{"WithClause", "CTE", "SelectStmt", "TableRef"}},
		{"database link", "SELECT * FROM SCHEMA1.LT1@REMOTE", []string{"SelectStmt", "TableRef", "ObjectName"}},
		{"json table", "SELECT JT.ID FROM T, JSON_TABLE(T.J, '$' COLUMNS (ID NUMBER PATH '$.id')) JT", []string{"SelectStmt", "JsonTableRef", "JsonTableColumn", "TableRef"}},
		{"xml table", "SELECT XT.ID FROM T, XMLTABLE('/root' PASSING T.X COLUMNS ID NUMBER PATH 'id') XT", []string{"SelectStmt", "XmlTableRef", "XmlTableColumn", "TableRef"}},
		{"containers", "SELECT A FROM CONTAINERS(T)", []string{"SelectStmt", "ContainersExpr", "ObjectName"}},
		{"inline external table", "SELECT EXT.A FROM EXTERNAL ((A NUMBER, B VARCHAR2(20)) TYPE ORACLE_LOADER DEFAULT DIRECTORY D ACCESS PARAMETERS (FIELDS TERMINATED BY ',') LOCATION ('x.csv')) EXT", []string{"SelectStmt", "InlineExternalTable", "ColumnDef"}},
		{"pivot", "SELECT * FROM (SELECT A, B FROM T) PIVOT (COUNT(*) FOR B IN (1 AS ONE))", []string{"SelectStmt", "PivotClause", "SubqueryRef"}},
		{"unpivot", "SELECT * FROM (SELECT A, B FROM T) UNPIVOT (VAL FOR COL IN (A AS 'A', B AS 'B'))", []string{"SelectStmt", "UnpivotClause", "SubqueryRef"}},
		{"model", "SELECT * FROM T MODEL DIMENSION BY (A) MEASURES (B) RULES (B[1] = 1)", []string{"SelectStmt", "ModelClause", "ModelRule"}},
		{"lateral inline view", "SELECT * FROM T, LATERAL (SELECT A FROM DUAL) L", []string{"SelectStmt", "LateralRef", "TableRef"}},
		{"hierarchical query", "SELECT A FROM T START WITH A = 1 CONNECT BY PRIOR A = B", []string{"SelectStmt", "HierarchicalClause", "UnaryExpr", "BinaryExpr"}},
		{"analytic function", "SELECT SUM(A) OVER (PARTITION BY B ORDER BY C) FROM T", []string{"SelectStmt", "FuncCallExpr", "WindowSpec", "SortBy"}},
		{"explain root", "EXPLAIN PLAN FOR SELECT * FROM T", []string{"ExplainPlanStmt", "SelectStmt"}},
		{"dml roots", "INSERT INTO T(A) SELECT A FROM T2", []string{"InsertStmt", "SelectStmt", "TableRef"}},
		{"ddl root", "CREATE TABLE T(A NUMBER)", []string{"CreateTableStmt", "ColumnDef"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.statement)
			require.NoError(t, err)
			require.Len(t, list.Items, 1)
			raw, ok := list.Items[0].(*ast.RawStmt)
			require.True(t, ok)

			seen := make(map[string]bool)
			ast.Inspect(raw, func(node ast.Node) bool {
				if node == nil {
					return false
				}
				typ := reflect.TypeOf(node)
				if typ.Kind() == reflect.Pointer {
					typ = typ.Elem()
				}
				seen[typ.Name()] = true
				return true
			})
			for _, want := range tc.expectedNodes {
				require.Contains(t, seen, want)
			}
		})
	}
}
