//go:build scriptdom

// Run with: go test -tags scriptdom ./mssql/parser/ -run TestScriptDOMDiff
//
// This test compares omni's AST shape against SqlScriptDOM by shelling out
// to the .NET harness at harness/mssql-scriptdom. The harness parses with
// Microsoft.SqlServer.TransactSql.ScriptDom and emits a compact JSON shape.
// We normalize omni's AST into the same shape and diff with go-cmp.

package parser

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/google/go-cmp/cmp"
)

type scriptdomHarness struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
}

var (
	harnessInstance *scriptdomHarness
	harnessInitErr  error
	harnessOnce     sync.Once
)

func getHarness(t *testing.T) *scriptdomHarness {
	t.Helper()
	harnessOnce.Do(func() {
		_, thisFile, _, _ := runtime.Caller(0)
		// mssql/parser/scriptdom_harness_test.go → repo root is ../..
		repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
		projDir := filepath.Join(repoRoot, "harness", "mssql-scriptdom")
		if _, err := os.Stat(projDir); err != nil {
			harnessInitErr = fmt.Errorf("harness project not found at %s", projDir)
			return
		}

		// Ensure the harness is built.
		build := exec.Command("dotnet", "build", "-c", "Release", "-v", "minimal", "--nologo")
		build.Dir = projDir
		if out, err := build.CombinedOutput(); err != nil {
			harnessInitErr = fmt.Errorf("dotnet build failed: %v\n%s", err, out)
			return
		}

		dll := filepath.Join(projDir, "bin", "Release", "net8.0", "mssql-scriptdom-harness.dll")
		cmd := exec.Command("dotnet", dll)
		cmd.Env = append(os.Environ(), "MSSQL_HARNESS_LINE=1")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			harnessInitErr = err
			return
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			harnessInitErr = err
			return
		}
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			harnessInitErr = err
			return
		}
		harnessInstance = &scriptdomHarness{
			cmd:    cmd,
			stdin:  stdin,
			stdout: bufio.NewReader(stdout),
		}
	})
	if harnessInitErr != nil {
		t.Fatalf("harness init: %v", harnessInitErr)
	}
	return harnessInstance
}

// Shape returns the SqlScriptDOM shape JSON for the given SQL.
func (h *scriptdomHarness) Shape(t *testing.T, sql string) map[string]any {
	t.Helper()
	h.mu.Lock()
	defer h.mu.Unlock()
	enc := base64.StdEncoding.EncodeToString([]byte(sql))
	if _, err := fmt.Fprintln(h.stdin, enc); err != nil {
		t.Fatalf("harness write: %v", err)
	}
	line, err := h.stdout.ReadBytes('\n')
	if err != nil {
		t.Fatalf("harness read: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(line, &out); err != nil {
		t.Fatalf("harness json: %v\nline=%s", err, line)
	}
	return out
}

// ---------- omni → shape normalization ----------

func omniShape(result *ast.List) map[string]any {
	stmts := make([]any, 0, len(result.Items))
	for _, it := range result.Items {
		stmts = append(stmts, omniStmt(it))
	}
	return map[string]any{"shape": map[string]any{"kind": "Script", "stmts": stmts}}
}

func omniStmt(n ast.Node) any {
	switch s := n.(type) {
	case *ast.SelectStmt:
		return map[string]any{"kind": "Select", "query": omniQuery(s)}
	case *ast.InsertStmt:
		return map[string]any{"kind": "Insert"}
	}
	return map[string]any{"kind": fmt.Sprintf("%T", n)}
}

func omniQuery(s *ast.SelectStmt) any {
	selectList := make([]any, 0)
	if s.TargetList != nil {
		for _, t := range s.TargetList.Items {
			selectList = append(selectList, omniSelectElement(t))
		}
	}
	var from []any
	if s.FromClause != nil {
		from = make([]any, 0, len(s.FromClause.Items))
		for _, tr := range s.FromClause.Items {
			from = append(from, omniTableRef(tr))
		}
	}
	q := map[string]any{"kind": "QuerySpec", "select_list": selectList}
	if from != nil {
		q["from"] = from
	}
	return q
}

func omniSelectElement(n ast.Node) any {
	rt, ok := n.(*ast.ResTarget)
	if !ok {
		return map[string]any{"kind": fmt.Sprintf("%T", n)}
	}
	switch v := rt.Val.(type) {
	case *ast.StarExpr:
		out := map[string]any{"kind": "SelectStar"}
		if v.Qualifier != "" {
			out["qualifier"] = v.Qualifier
		}
		return out
	case *ast.SelectAssign:
		out := map[string]any{
			"kind":      "SelectSetVariable",
			"variable":  v.Variable,
			"expr_kind": exprKind(v.Value),
		}
		return out
	default:
		out := map[string]any{
			"kind":      "SelectScalar",
			"expr_kind": exprKind(rt.Val),
		}
		if rt.Name != "" {
			out["alias"] = rt.Name
		}
		return out
	}
}

func exprKind(n ast.Node) string {
	// Map omni AST types to SqlScriptDOM-ish class names for the fields we diff on.
	switch n.(type) {
	case *ast.ColumnRef:
		return "ColumnReferenceExpression"
	case *ast.Literal:
		// SqlScriptDOM uses IntegerLiteral / StringLiteral subclasses, but for
		// our current scenarios the top-level kind "Literal" is enough and
		// the harness's raw output uses the specific subclass. To keep diffs
		// stable, callers should avoid asserting expr_kind for literals unless
		// they produce the same subclass name here.
		return "Literal"
	case *ast.BinaryExpr:
		return "BinaryExpression"
	case *ast.FuncCallExpr:
		return "FunctionCall"
	case *ast.SubqueryExpr:
		return "ScalarSubquery"
	}
	return fmt.Sprintf("%T", n)
}

func omniTableRef(n ast.Node) any {
	switch v := n.(type) {
	case *ast.TableRef:
		return map[string]any{"kind": "NamedTableReference"}
	case *ast.AliasedTableRef:
		return omniAliasedTable(v)
	case *ast.SubqueryExpr:
		return map[string]any{"kind": "QueryDerivedTable"}
	case *ast.JoinClause:
		left := omniTableRef(v.Left)
		right := omniTableRef(v.Right)
		if isUnqualifiedJoin(v.Type) {
			return map[string]any{
				"kind":      "UnqualifiedJoin",
				"join_type": unqualifiedJoinKind(v.Type),
				"left":      left,
				"right":     right,
			}
		}
		return map[string]any{
			"kind":      "QualifiedJoin",
			"join_type": qualifiedJoinKind(v.Type),
			"left":      left,
			"right":     right,
		}
	}
	return map[string]any{"kind": fmt.Sprintf("%T", n)}
}

func omniAliasedTable(n *ast.AliasedTableRef) any {
	kind := "NamedTableReference"
	// Identify inner kind to match SqlScriptDOM's table-ref hierarchy.
	switch inner := n.Table.(type) {
	case *ast.SubqueryExpr:
		kind = "QueryDerivedTable"
		_ = inner
	case *ast.ValuesClause:
		kind = "InlineDerivedTable"
	case *ast.FuncCallExpr:
		kind = classifyFuncCall(inner)
	}
	out := map[string]any{"kind": kind}
	if n.Alias != "" {
		out["alias"] = n.Alias
	}
	if cols := plainColumnNames(n.Columns); cols != nil {
		out["columns"] = cols
	} else if wc := withCols(n.Columns); wc != nil {
		out["with_cols"] = wc
	}
	return out
}

// classifyFuncCall returns SqlScriptDOM's table-ref subclass for a
// FuncCallExpr used as a table source.
func classifyFuncCall(fc *ast.FuncCallExpr) string {
	name := funcCallName(fc)
	switch strings.ToUpper(name) {
	case "OPENJSON":
		return "OpenJsonTableReference"
	case "OPENROWSET":
		return "BulkOpenRowset" // placeholder; scripters emit various kinds
	case "OPENQUERY":
		return "OpenQueryTableReference"
	case "OPENXML":
		return "OpenXmlTableReference"
	case "OPENDATASOURCE":
		return "SchemaObjectFunctionTableReference"
	}
	return "SchemaObjectFunctionTableReference"
}

func funcCallName(fc *ast.FuncCallExpr) string {
	if fc.Name != nil {
		return fc.Name.Object
	}
	return ""
}

func plainColumnNames(l *ast.List) []any {
	if l == nil || len(l.Items) == 0 {
		return nil
	}
	names := make([]any, 0, len(l.Items))
	for _, it := range l.Items {
		s, ok := it.(*ast.String)
		if !ok {
			return nil // not a plain name list (probably with_cols)
		}
		names = append(names, s.Str)
	}
	return names
}

func withCols(l *ast.List) []any {
	if l == nil || len(l.Items) == 0 {
		return nil
	}
	out := make([]any, 0, len(l.Items))
	for _, it := range l.Items {
		cd, ok := it.(*ast.ColumnDef)
		if !ok {
			return nil
		}
		entry := map[string]any{"name": cd.Name}
		if cd.DataType != nil {
			entry["type"] = cd.DataType.Name
		}
		out = append(out, entry)
	}
	return out
}

func isUnqualifiedJoin(jt ast.JoinType) bool {
	switch jt {
	case ast.JoinCross, ast.JoinCrossApply, ast.JoinOuterApply:
		return true
	}
	return false
}

func qualifiedJoinKind(jt ast.JoinType) string {
	switch jt {
	case ast.JoinInner:
		return "Inner"
	case ast.JoinLeft:
		return "LeftOuter"
	case ast.JoinRight:
		return "RightOuter"
	case ast.JoinFull:
		return "FullOuter"
	}
	return "Inner"
}

func unqualifiedJoinKind(jt ast.JoinType) string {
	switch jt {
	case ast.JoinCross:
		return "CrossJoin"
	case ast.JoinCrossApply:
		return "CrossApply"
	case ast.JoinOuterApply:
		return "OuterApply"
	}
	return "CrossJoin"
}

// ---------- the test ----------

type diffFixture struct {
	name string
	sql  string
	// Some scenarios (literals, function calls deep inside exprs) diverge in
	// detailed `expr_kind` between omni and ScriptDOM. Fields listed here are
	// redacted on both sides before comparison.
	ignoreFields []string
}

func TestScriptDOMDiff(t *testing.T) {
	h := getHarness(t)

	fixtures := []diffFixture{
		{name: "gap2-subquery-cols", sql: `SELECT * FROM (SELECT a, b FROM t) AS x(c1, c2)`},
		{name: "gap3-tvf-cols", sql: `SELECT * FROM fn(1, 2) AS t(a, b)`},
		{name: "gap4-values-in-from", sql: `SELECT * FROM (VALUES (1, 2), (3, 4)) AS t(a, b)`},
		{name: "gap1-alias-eq-expr", sql: `SELECT a = c FROM t`},
		{name: "sel-var-assign", sql: `SELECT @v = x FROM t`},
		{name: "openjson-with", sql: `SELECT * FROM OPENJSON(@j) WITH (id int '$.id', name nvarchar(50) '$.name') AS t`},
		{name: "subquery-alias-no-cols", sql: `SELECT * FROM (SELECT a FROM t) AS x`},
		{name: "plain-named-table", sql: `SELECT a FROM t`},
	}

	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			got, err := Parse(f.sql)
			if err != nil {
				t.Fatalf("omni parse: %v", err)
			}
			want := h.Shape(t, f.sql)
			gotShape := omniShape(got)
			if errs, ok := want["errors"]; ok && errs != nil {
				t.Fatalf("ScriptDOM rejected: %v", errs)
			}
			delete(want, "errors")
			if diff := cmp.Diff(want, gotShape); diff != "" {
				t.Errorf("shape mismatch for %q\n-scriptdom +omni\n%s", f.sql, diff)
			}
		})
	}
}

// TestScriptDOMRejectAlignment verifies that SQL which SqlScriptDOM rejects
// is ALSO rejected by omni. Pinned via the negative fixtures below so that
// comma-list strictness can't regress silently.
func TestScriptDOMRejectAlignment(t *testing.T) {
	h := getHarness(t)

	fixtures := []struct {
		name string
		sql  string
	}{
		// empty lists
		{"empty/alias-cols", `SELECT * FROM (SELECT 1 AS a) AS t()`},
		{"empty/cte-cols", `WITH c() AS (SELECT 1) SELECT * FROM c`},
		{"empty/view-cols", `CREATE VIEW v() AS SELECT 1`},
		// INSERT INTO t () — SqlScriptDOM accepts; parity with omni ensured by
		// commaListAllowEmpty at the call site (insert.go parseInsertStmt).
		{"empty/in-list", `SELECT * FROM t WHERE a IN ()`},
		{"empty/values-row", `INSERT INTO t VALUES ()`},
		{"empty/create-table-cols", `CREATE TABLE t ()`},
		{"empty/pivot-in", `SELECT * FROM t PIVOT (SUM(x) FOR y IN ()) p`},
		{"empty/unpivot-in", `SELECT * FROM t UNPIVOT (v FOR c IN ()) u`},
		{"empty/rollup", `SELECT COUNT(*) FROM t GROUP BY ROLLUP()`},
		{"empty/cube", `SELECT COUNT(*) FROM t GROUP BY CUBE()`},
		{"empty/partition-by", `SELECT *, ROW_NUMBER() OVER (PARTITION BY) FROM t`},
		{"empty/index-cols", `CREATE INDEX ix ON t()`},
		// trailing commas
		{"trail/alias-cols", `SELECT * FROM (SELECT 1 AS a) AS t(a,)`},
		{"trail/cte-cols", `WITH c(a,) AS (SELECT 1) SELECT * FROM c`},
		{"trail/cte-list", `WITH c AS (SELECT 1), SELECT 1`},
		{"trail/view-cols", `CREATE VIEW v(a,) AS SELECT 1`},
		{"trail/insert-cols", `INSERT INTO t (a,) VALUES (1)`},
		{"trail/values-row", `INSERT INTO t VALUES (1,)`},
		{"trail/in-list", `SELECT * FROM t WHERE a IN (1,)`},
		{"trail/func-args", `SELECT f(1,) FROM t`},
		{"trail/tvf-args", `SELECT * FROM dbo.fn(1,) AS t(a)`},
		{"trail/select-list", `SELECT * FROM t ORDER BY a,`},
		{"trail/rollup", `SELECT COUNT(*) FROM t GROUP BY ROLLUP(a,)`},
		{"trail/partition-by", `SELECT *, ROW_NUMBER() OVER (PARTITION BY a,) FROM t`},
		{"trail/pivot-in", `SELECT * FROM t PIVOT (SUM(x) FOR y IN ([a],)) p`},
		{"trail/unpivot-in", `SELECT * FROM t UNPIVOT (v FOR c IN (a,)) u`},
		{"trail/index-cols", `CREATE INDEX ix ON t(a,)`},
		{"trail/primary-key", `ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (a,b,)`},
		{"trail/constraint-cols", `CREATE TABLE t (a INT, b INT, CONSTRAINT pk PRIMARY KEY (a,b,))`},
		{"trail/table-hint", `SELECT * FROM t WITH (NOLOCK,)`},
	}

	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			// Verify SqlScriptDOM rejects (sanity on our expectation).
			want := h.Shape(t, f.sql)
			errs, _ := want["errors"].([]any)
			if len(errs) == 0 {
				t.Fatalf("fixture asserts reject but SqlScriptDOM accepted: %s", f.sql)
			}
			// Verify omni also rejects.
			if _, err := Parse(f.sql); err == nil {
				t.Errorf("omni accepted input that SqlScriptDOM rejects: %s", f.sql)
			}
		})
	}
}
