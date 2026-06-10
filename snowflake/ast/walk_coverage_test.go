package ast_test

// Walker coverage regression tests for node-bearing fields whose types are
// not themselves Node (helper structs like WhenClause / OrderItem /
// WindowSpec / CTE / SelectTarget), slices of pointers to Node structs
// ([]*ObjectName, []*ColumnDef, ...), and by-value Node structs
// (FuncCallExpr.Name). The generated walker used to skip ALL of these, so
// ast.Inspect never reached e.g. CASE WHEN conditions or OVER (PARTITION BY
// ... ORDER BY ...) expressions, silently breaking column collection in every
// walker-based consumer (analysis query-span included).
//
// These tests parse real SQL and assert that specific nodes are visited.
// They live in an external test package so they can use the parser.

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// walkSummary records everything one ast.Inspect traversal visited.
type walkSummary struct {
	tags     map[ast.NodeTag]int
	colRefs  map[string]bool // normalized dotted ColumnRef ("A", "S.MF")
	objNames map[string]bool // normalized dotted ObjectName ("MYSCHEMA.MYFUNC")
	literals map[string]bool // raw Literal source text ("10", "'JAN'")
}

// summarize parses sql and walks the resulting file, recording all visits.
func summarize(t *testing.T, sql string) *walkSummary {
	t.Helper()
	file, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse %q: %v", sql, err)
	}
	s := &walkSummary{
		tags:     map[ast.NodeTag]int{},
		colRefs:  map[string]bool{},
		objNames: map[string]bool{},
		literals: map[string]bool{},
	}
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		s.tags[n.Tag()]++
		switch x := n.(type) {
		case *ast.ColumnRef:
			parts := make([]string, len(x.Parts))
			for i, p := range x.Parts {
				parts[i] = p.Normalize()
			}
			s.colRefs[strings.Join(parts, ".")] = true
		case *ast.ObjectName:
			s.objNames[x.Normalize()] = true
		case *ast.Literal:
			s.literals[x.Value] = true
		}
		return true
	})
	return s
}

// walkCoverageCase asserts minimum visit evidence for one SQL statement.
type walkCoverageCase struct {
	name string
	sql  string
	cols []string            // ColumnRefs that must be visited (normalized)
	objs []string            // ObjectNames that must be visited (normalized)
	lits []string            // Literals that must be visited (raw text)
	tags map[ast.NodeTag]int // minimum visit count per tag
}

func runWalkCoverage(t *testing.T, cases []walkCoverageCase) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := summarize(t, c.sql)
			for _, col := range c.cols {
				if !s.colRefs[col] {
					t.Errorf("ColumnRef %q not visited by Walk; visited: %v", col, keys(s.colRefs))
				}
			}
			for _, obj := range c.objs {
				if !s.objNames[obj] {
					t.Errorf("ObjectName %q not visited by Walk; visited: %v", obj, keys(s.objNames))
				}
			}
			for _, lit := range c.lits {
				if !s.literals[lit] {
					t.Errorf("Literal %q not visited by Walk; visited: %v", lit, keys(s.literals))
				}
			}
			for tag, min := range c.tags {
				if s.tags[tag] < min {
					t.Errorf("tag %v visited %d times, want >= %d", tag, s.tags[tag], min)
				}
			}
		})
	}
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// ---------------------------------------------------------------------------
// Expression-level gaps (CASE, window functions, WITHIN GROUP, JSON literals)
// ---------------------------------------------------------------------------

func TestWalkCoverage_Expressions(t *testing.T) {
	runWalkCoverage(t, []walkCoverageCase{
		{
			name: "CaseExpr.Whens searched",
			sql:  "SELECT CASE WHEN wa > 0 THEN wb ELSE wc END FROM t",
			cols: []string{"WA", "WB", "WC"},
		},
		{
			name: "CaseExpr.Whens simple",
			sql:  "SELECT CASE sx WHEN sy THEN sz END FROM t",
			cols: []string{"SX", "SY", "SZ"},
		},
		{
			name: "FuncCallExpr.Over partition and order",
			sql:  "SELECT ROW_NUMBER() OVER (PARTITION BY px ORDER BY oy) FROM t",
			cols: []string{"PX", "OY"},
		},
		{
			name: "FuncCallExpr.Over frame bound offsets",
			sql:  "SELECT SUM(sx) OVER (ORDER BY oy ROWS BETWEEN 5 PRECEDING AND 3 FOLLOWING) FROM t",
			cols: []string{"SX", "OY"},
			lits: []string{"5", "3"},
		},
		{
			name: "FuncCallExpr.OrderBy within group",
			sql:  "SELECT LISTAGG(lx, ',') WITHIN GROUP (ORDER BY wy) FROM t",
			cols: []string{"LX", "WY"},
		},
		{
			name: "FuncCallExpr.Name qualified function name",
			sql:  "SELECT myschema.myfunc(fa) FROM t",
			cols: []string{"FA"},
			objs: []string{"MYSCHEMA.MYFUNC"},
		},
		{
			name: "JsonLiteralExpr.Pairs values",
			sql:  "SELECT {'k': jv} FROM t",
			cols: []string{"JV"},
		},
	})
}

// ---------------------------------------------------------------------------
// SELECT clause gaps (targets, WITH, GROUP BY, ORDER BY, FETCH)
// ---------------------------------------------------------------------------

func TestWalkCoverage_SelectClauses(t *testing.T) {
	runWalkCoverage(t, []walkCoverageCase{
		{
			name: "SelectStmt.Targets expressions",
			sql:  "SELECT ta + tb AS s FROM t",
			cols: []string{"TA", "TB"},
		},
		{
			name: "SelectStmt.With CTE body",
			sql:  "WITH c AS (SELECT ca FROM t) SELECT * FROM c",
			cols: []string{"CA"},
			tags: map[ast.NodeTag]int{ast.T_SelectStmt: 2},
		},
		{
			name: "SelectStmt.GroupBy items",
			sql:  "SELECT COUNT(*) FROM t GROUP BY ga, gb",
			cols: []string{"GA", "GB"},
		},
		{
			name: "SelectStmt.OrderBy items",
			sql:  "SELECT a FROM t ORDER BY ob DESC NULLS LAST",
			cols: []string{"OB"},
		},
		{
			name: "SelectStmt.Fetch count",
			sql:  "SELECT a FROM t FETCH FIRST 10 ROWS ONLY",
			lits: []string{"10"},
		},
		{
			name: "scalar subquery in select list reaches inner select",
			sql:  "SELECT (SELECT inner_col FROM u) FROM t",
			cols: []string{"INNER_COL"},
			tags: map[ast.NodeTag]int{ast.T_SelectStmt: 2},
		},
	})
}

// ---------------------------------------------------------------------------
// FROM-clause helper gaps (PIVOT IN values/order, UNPIVOT columns,
// MATCH_RECOGNIZE measures/define/order)
// ---------------------------------------------------------------------------

func TestWalkCoverage_FromClauses(t *testing.T) {
	runWalkCoverage(t, []walkCoverageCase{
		{
			name: "PivotInClause.Values",
			sql:  "SELECT * FROM monthly_sales PIVOT (SUM(amount) FOR month IN ('JAN', 'FEB'))",
			cols: []string{"AMOUNT", "MONTH"},
			tags: map[ast.NodeTag]int{ast.T_PivotValue: 2},
		},
		{
			name: "PivotInClause.OrderBy",
			sql:  "SELECT * FROM quarterly_sales PIVOT(SUM(amount) FOR quarter IN (ANY ORDER BY quarter)) ORDER BY empid",
			cols: []string{"QUARTER", "EMPID"},
		},
		{
			name: "UnpivotClause.Columns",
			sql:  "SELECT * FROM monthly_sales UNPIVOT (sales FOR month IN (jan, feb))",
			tags: map[ast.NodeTag]int{ast.T_UnpivotColumn: 2},
		},
		{
			name: "MatchRecognizeClause OrderBy/Measures/Define",
			sql: `SELECT * FROM stock_price_history MATCH_RECOGNIZE(
				PARTITION BY company
				ORDER BY price_date
				MEASURES FIRST(price_date) AS start_date
				ONE ROW PER MATCH
				PATTERN(dn up)
				DEFINE dn AS price < LAG(price), up AS price > LAG(price)
			)`,
			cols: []string{"COMPANY", "PRICE_DATE", "PRICE"},
		},
	})
}

// ---------------------------------------------------------------------------
// DML gaps (UPDATE SET, MERGE WHEN, INSERT ALL branches, SET variables)
// ---------------------------------------------------------------------------

func TestWalkCoverage_DML(t *testing.T) {
	runWalkCoverage(t, []walkCoverageCase{
		{
			name: "UpdateStmt.Sets value expressions",
			sql:  "UPDATE t SET c1 = ux + 1 WHERE id = 5",
			cols: []string{"UX", "ID"},
			objs: []string{"C1"},
		},
		{
			name: "MergeStmt.Whens conditions, sets, and insert values",
			sql: `MERGE INTO t USING s ON t.id = s.id
				WHEN MATCHED AND s.mf > 0 THEN UPDATE SET c = s.mv
				WHEN NOT MATCHED THEN INSERT (a) VALUES (s.mw)`,
			cols: []string{"S.MF", "S.MV", "S.MW"},
		},
		{
			name: "InsertMultiStmt.Branches targets and values",
			sql: `INSERT ALL
				INTO t1 VALUES (1, 'a')
				INTO t2 VALUES (2, 'b')
				SELECT * FROM src`,
			objs: []string{"T1", "T2"},
			lits: []string{"1", "2", "a", "b"},
		},
		{
			name: "SetStmt.Vars values",
			sql:  "SET V1 = 10",
			lits: []string{"10"},
		},
	})
}

// ---------------------------------------------------------------------------
// DDL gaps (column defs, constraints, indexes, clones, tags, policies,
// routine signatures, copy options, stream time travel, ...)
// ---------------------------------------------------------------------------

func TestWalkCoverage_DDL(t *testing.T) {
	runWalkCoverage(t, []walkCoverageCase{
		{
			name: "CreateTableStmt.Columns and Constraints with FK references",
			sql:  "CREATE TABLE t (c1 INT DEFAULT 7, c2 INT, CONSTRAINT fk FOREIGN KEY (c2) REFERENCES p (py))",
			lits: []string{"7"},
			objs: []string{"P"},
			tags: map[ast.NodeTag]int{ast.T_ColumnDef: 2, ast.T_TableConstraint: 1, ast.T_TypeName: 2},
		},
		{
			name: "ColumnDef.InlineConstraint references",
			sql:  "CREATE TABLE t (customer_id INT FOREIGN KEY REFERENCES customers (id))",
			objs: []string{"CUSTOMERS"},
		},
		{
			name: "CreateTableStmt.Indexes (hybrid table)",
			sql:  "CREATE HYBRID TABLE t (id INT, full_name VARCHAR(255), INDEX index_full_name (full_name))",
			tags: map[ast.NodeTag]int{ast.T_TableIndex: 1},
		},
		{
			name: "CreateTableStmt.Clone source",
			sql:  "CREATE TABLE t1 CLONE t2",
			objs: []string{"T2"},
		},
		{
			name: "CreateViewStmt.RowPolicy",
			sql:  "CREATE VIEW v WITH ROW ACCESS POLICY my_policy ON (col1, col2) AS SELECT 1",
			objs: []string{"MY_POLICY"},
		},
		{
			name: "AlterTableStmt.Actions add column with default",
			sql:  "ALTER TABLE t ADD COLUMN nc INT DEFAULT 3",
			lits: []string{"3"},
			tags: map[ast.NodeTag]int{ast.T_ColumnDef: 1},
		},
		{
			name: "AlterTableAction.UnsetTags",
			sql:  "ALTER TABLE t UNSET TAG (my_tag, other_tag)",
			objs: []string{"MY_TAG", "OTHER_TAG"},
		},
		{
			name: "AlterViewStmt set tag assignments",
			sql:  "ALTER VIEW v SET TAG (env = 'prod')",
			objs: []string{"ENV"},
		},
		{
			name: "CopyIntoTableStmt.Options nested file format group",
			sql:  "COPY INTO t FROM @s FILE_FORMAT = (TYPE = CSV SKIP_HEADER = 1 FIELD_DELIMITER = ',')",
			tags: map[ast.NodeTag]int{ast.T_CopyOption: 3},
		},
		{
			name: "CreateRoutineStmt argument and return types",
			sql:  "CREATE FUNCTION multiply(a NUMBER, b NUMBER) RETURNS NUMBER AS 'a * b'",
			tags: map[ast.NodeTag]int{ast.T_TypeName: 3},
		},
		{
			name: "AlterRoutineStmt.ArgTypes signature",
			sql:  "ALTER FUNCTION f(NUMBER) RENAME TO g",
			tags: map[ast.NodeTag]int{ast.T_TypeName: 1},
		},
		{
			name: "CommentStmt.Signature types",
			sql:  "COMMENT ON FUNCTION f(INT, STRING) IS 'fn'",
			tags: map[ast.NodeTag]int{ast.T_TypeName: 2},
		},
		{
			name: "CreateStreamStmt.TimeTravel value",
			sql:  "CREATE STREAM s ON TABLE t AT (TIMESTAMP => TO_TIMESTAMP(40*365*86400))",
			lits: []string{"40", "365", "86400"},
		},
		{
			name: "CreatePolicyStmt.Args",
			sql:  "CREATE MASKING POLICY mp AS (val string) RETURNS string -> val",
			cols: []string{"VAL"},
			tags: map[ast.NodeTag]int{ast.T_PolicyArg: 1},
		},
		{
			name: "CreateSemanticViewStmt.Sections",
			sql:  "CREATE SEMANTIC VIEW sv TABLES (orders) METRICS (orders.total AS SUM(amount))",
			tags: map[ast.NodeTag]int{ast.T_SemanticViewSection: 2},
		},
		{
			name: "CreateIntegrationStmt.Triggers (resource monitor)",
			sql:  "CREATE RESOURCE MONITOR rm WITH CREDIT_QUOTA = 100 TRIGGERS ON 75 PERCENT DO NOTIFY",
			tags: map[ast.NodeTag]int{ast.T_ResourceMonitorTrigger: 1},
		},
		{
			name: "CreateReplicationGroupStmt.Options",
			sql:  "CREATE REPLICATION GROUP rg OBJECT_TYPES = DATABASES ALLOWED_ACCOUNTS = org.acct",
			tags: map[ast.NodeTag]int{ast.T_GroupOption: 1},
		},
		{
			name: "AlterExternalTableStmt.Files literals",
			sql:  "ALTER EXTERNAL TABLE et ADD FILES ('p/f1.parquet', 'p/f2.parquet')",
			lits: []string{"p/f1.parquet", "p/f2.parquet"},
		},
	})
}

// TestWalkCoverage_PostOrderBalance verifies the widened traversal still
// delivers balanced pre/post events (every recursion ends with Visit(nil)).
func TestWalkCoverage_PostOrderBalance(t *testing.T) {
	sql := `WITH c AS (SELECT ca FROM t)
		SELECT CASE WHEN a > 0 THEN b END,
		       ROW_NUMBER() OVER (PARTITION BY x ORDER BY y ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)
		FROM c GROUP BY g ORDER BY o FETCH FIRST 1 ROWS ONLY`
	file, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	pre, post := 0, 0
	v := &balanceVisitor{pre: &pre, post: &post}
	ast.Walk(v, file)
	if pre != post {
		t.Errorf("unbalanced traversal: %d pre-order visits, %d post-order visits", pre, post)
	}
	if pre == 0 {
		t.Error("no nodes visited")
	}
}

type balanceVisitor struct {
	pre, post *int
}

func (v *balanceVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		*v.post++
		return nil
	}
	*v.pre++
	return v
}
