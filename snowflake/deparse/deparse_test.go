package deparse_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/bytebase/omni/snowflake/parser"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Round-trip helpers
// ---------------------------------------------------------------------------

// assertRoundTrip parses sql, deparses it, re-parses the output, and checks
// that the two ASTs are structurally equivalent (ignoring Loc fields).
func assertRoundTrip(t *testing.T, sql string) string {
	t.Helper()
	ast1, err := parser.Parse(sql)
	require.NoErrorf(t, err, "first parse failed for %q", sql)

	out, err := deparse.DeparseFile(ast1)
	require.NoErrorf(t, err, "deparse failed for %q", sql)

	ast2, err := parser.Parse(out)
	require.NoErrorf(t, err, "second parse failed: deparsed=%q (original=%q)", out, sql)

	stripped1 := stripLoc(ast1)
	stripped2 := stripLoc(ast2)
	require.Truef(t, reflect.DeepEqual(stripped1, stripped2),
		"AST mismatch after round-trip\n  original SQL:   %q\n  deparsed SQL:   %q",
		sql, out)
	return out
}

// stripLoc walks an *ast.File and zeroes out every Loc field so that
// DeepEqual ignores source positions.
func stripLoc(node *ast.File) *ast.File {
	if node == nil {
		return nil
	}
	// We use a simple recursive value-copier via reflect to zero Loc fields.
	v := reflect.ValueOf(node)
	cloned := deepStripLoc(v)
	return cloned.Interface().(*ast.File)
}

func deepStripLoc(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return v
		}
		newPtr := reflect.New(v.Type().Elem())
		newPtr.Elem().Set(deepStripLoc(v.Elem()))
		return newPtr
	case reflect.Struct:
		newStruct := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			fv := v.Field(i)
			if field.Name == "Loc" {
				// zero the Loc field
				newStruct.Field(i).Set(reflect.Zero(fv.Type()))
				continue
			}
			newStruct.Field(i).Set(deepStripLoc(fv))
		}
		return newStruct
	case reflect.Slice:
		if v.IsNil() {
			return v
		}
		newSlice := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			newSlice.Index(i).Set(deepStripLoc(v.Index(i)))
		}
		return newSlice
	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		inner := deepStripLoc(v.Elem())
		newIface := reflect.New(v.Type()).Elem()
		newIface.Set(inner)
		return newIface
	default:
		return v
	}
}

// ---------------------------------------------------------------------------
// SELECT tests
// ---------------------------------------------------------------------------

func TestDeparse_Select_Simple(t *testing.T) {
	assertRoundTrip(t, `SELECT 1`)
}

func TestDeparse_Select_StarFromTable(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t`)
}

func TestDeparse_Select_ColumnList(t *testing.T) {
	assertRoundTrip(t, `SELECT a, b, c FROM t`)
}

func TestDeparse_Select_Aliases(t *testing.T) {
	assertRoundTrip(t, `SELECT a AS x, b AS y FROM t`)
}

func TestDeparse_Select_QualifiedStar(t *testing.T) {
	assertRoundTrip(t, `SELECT t.* FROM t`)
}

func TestDeparse_Select_ExcludeStar(t *testing.T) {
	assertRoundTrip(t, `SELECT * EXCLUDE (a, b) FROM t`)
}

func TestDeparse_Select_ExcludeStarBare(t *testing.T) {
	// Bare single-column EXCLUDE deparses to the parenthesized form, which
	// re-parses to the same AST (single Exclude entry).
	assertRoundTrip(t, `SELECT * EXCLUDE a FROM t`)
}

func TestDeparse_Select_RenameStar(t *testing.T) {
	assertRoundTrip(t, `SELECT * RENAME a AS b FROM t`)
}

func TestDeparse_Select_RenameStarList(t *testing.T) {
	assertRoundTrip(t, `SELECT * RENAME (a AS b, c AS d) FROM t`)
}

func TestDeparse_Select_ExcludeAndRenameStar(t *testing.T) {
	assertRoundTrip(t, `SELECT * EXCLUDE x RENAME (a AS b, c AS d) FROM t`)
}

func TestDeparse_Select_QualifiedStarExcludeRename(t *testing.T) {
	assertRoundTrip(t, `SELECT t.* EXCLUDE a, u.* RENAME b AS c FROM t INNER JOIN u ON t.id = u.id`)
}

func TestDeparse_Select_OrderByAll(t *testing.T) {
	assertRoundTrip(t, `SELECT a, b FROM t ORDER BY ALL`)
}

func TestDeparse_Select_OrderByAllDesc(t *testing.T) {
	assertRoundTrip(t, `SELECT a, b FROM t ORDER BY ALL DESC`)
}

func TestDeparse_Select_OrderByAllNullsFirst(t *testing.T) {
	assertRoundTrip(t, `SELECT a, b FROM t ORDER BY ALL NULLS FIRST`)
}

func TestDeparse_Select_TrailingComma(t *testing.T) {
	// The trailing comma is normalized away on deparse; the AST (3 targets)
	// must round-trip identically.
	assertRoundTrip(t, `SELECT a, b, c, FROM t`)
}

func TestDeparse_Select_Where(t *testing.T) {
	assertRoundTrip(t, `SELECT a FROM t WHERE a > 1`)
}

func TestDeparse_Where_OuterJoinMarker(t *testing.T) {
	out := assertRoundTrip(t, `SELECT t1.c1, t2.c2 FROM t1, t2 WHERE t1.c1 = t2.c2(+)`)
	require.Contains(t, out, "t2.c2(+)")
}

func TestDeparse_Where_OuterJoinMarkerMultiple(t *testing.T) {
	// A space before the marker in the source canonicalizes to no-space on
	// deparse; the round-trip still re-parses to an equivalent AST.
	out := assertRoundTrip(t, `SELECT t1.c1, t2.c2 FROM t1, t2 WHERE t1.c1 = t2.c2 (+) AND t1.c3 = t2.c4 (+)`)
	require.Contains(t, out, "t2.c2(+)")
	require.Contains(t, out, "t2.c4(+)")
}

func TestDeparse_Select_GroupBy(t *testing.T) {
	assertRoundTrip(t, `SELECT a, COUNT(*) FROM t GROUP BY a`)
}

func TestDeparse_Select_GroupByAll(t *testing.T) {
	assertRoundTrip(t, `SELECT a, b, COUNT(*) FROM t GROUP BY ALL`)
}

func TestDeparse_Select_Having(t *testing.T) {
	assertRoundTrip(t, `SELECT a, COUNT(*) FROM t GROUP BY a HAVING COUNT(*) > 1`)
}

func TestDeparse_Select_OrderBy(t *testing.T) {
	assertRoundTrip(t, `SELECT a FROM t ORDER BY a DESC`)
}

func TestDeparse_Select_OrderByNulls(t *testing.T) {
	assertRoundTrip(t, `SELECT a FROM t ORDER BY a ASC NULLS FIRST`)
}

func TestDeparse_Select_LimitOffset(t *testing.T) {
	assertRoundTrip(t, `SELECT a FROM t LIMIT 10 OFFSET 5`)
}

func TestDeparse_Select_Qualify(t *testing.T) {
	assertRoundTrip(t, `SELECT a, ROW_NUMBER() OVER (PARTITION BY a ORDER BY b ASC) AS rn FROM t QUALIFY rn = 1`)
}

func TestDeparse_Select_Distinct(t *testing.T) {
	assertRoundTrip(t, `SELECT DISTINCT a, b FROM t`)
}

func TestDeparse_Select_Top(t *testing.T) {
	assertRoundTrip(t, `SELECT TOP 10 a FROM t`)
}

func TestDeparse_Select_WithCTE(t *testing.T) {
	assertRoundTrip(t, `WITH cte AS (SELECT 1 AS x) SELECT x FROM cte`)
}

func TestDeparse_Select_WithMultipleCTEs(t *testing.T) {
	assertRoundTrip(t, `WITH a AS (SELECT 1), b AS (SELECT 2) SELECT * FROM a, b`)
}

func TestDeparse_Select_WithCTEColumnsAlias(t *testing.T) {
	assertRoundTrip(t, `WITH cte (x, y) AS (SELECT 1, 2) SELECT x, y FROM cte`)
}

func TestDeparse_Select_Subquery(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM (SELECT a FROM t) AS sub`)
}

func TestDeparse_Select_Join(t *testing.T) {
	assertRoundTrip(t, `SELECT a.x, b.y FROM a JOIN b ON a.id = b.id`)
}

func TestDeparse_Select_LeftJoin(t *testing.T) {
	assertRoundTrip(t, `SELECT a.x, b.y FROM a LEFT JOIN b ON a.id = b.id`)
}

func TestDeparse_Select_CrossJoin(t *testing.T) {
	assertRoundTrip(t, `SELECT a.x FROM a CROSS JOIN b`)
}

func TestDeparse_Select_JoinUsing(t *testing.T) {
	assertRoundTrip(t, `SELECT a.x FROM a JOIN b USING (id)`)
}

func TestDeparse_Select_Fetch(t *testing.T) {
	assertRoundTrip(t, `SELECT a FROM t FETCH FIRST 5 ROWS ONLY`)
}

// ---------------------------------------------------------------------------
// SET operations
// ---------------------------------------------------------------------------

func TestDeparse_Union(t *testing.T) {
	assertRoundTrip(t, `SELECT 1 UNION SELECT 2`)
}

func TestDeparse_UnionAll(t *testing.T) {
	assertRoundTrip(t, `SELECT 1 UNION ALL SELECT 2`)
}

func TestDeparse_Intersect(t *testing.T) {
	assertRoundTrip(t, `SELECT 1 INTERSECT SELECT 1`)
}

func TestDeparse_Except(t *testing.T) {
	assertRoundTrip(t, `SELECT 1 EXCEPT SELECT 2`)
}

func TestDeparse_UnionByName(t *testing.T) {
	assertRoundTrip(t, `SELECT a FROM t UNION BY NAME SELECT a FROM s`)
}

// ---------------------------------------------------------------------------
// Expression tests
// ---------------------------------------------------------------------------

func TestDeparse_Expr_Arithmetic(t *testing.T) {
	assertRoundTrip(t, `SELECT a + b * c - d FROM t`)
}

func TestDeparse_Expr_BooleanLogic(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE a = 1 AND b = 2 OR c = 3`)
}

func TestDeparse_Expr_UnaryMinus(t *testing.T) {
	assertRoundTrip(t, `SELECT -a FROM t`)
}

func TestDeparse_Expr_Not(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE NOT a = 1`)
}

func TestDeparse_Expr_Cast(t *testing.T) {
	assertRoundTrip(t, `SELECT CAST(a AS VARCHAR(100)) FROM t`)
}

func TestDeparse_Expr_CastColonColon(t *testing.T) {
	assertRoundTrip(t, `SELECT a::INT FROM t`)
}

func TestDeparse_Expr_TryCast(t *testing.T) {
	assertRoundTrip(t, `SELECT TRY_CAST(a AS INT) FROM t`)
}

func TestDeparse_Expr_CaseSimple(t *testing.T) {
	assertRoundTrip(t, `SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END FROM t`)
}

func TestDeparse_Expr_CaseSearched(t *testing.T) {
	assertRoundTrip(t, `SELECT CASE WHEN a > 1 THEN 'big' ELSE 'small' END FROM t`)
}

func TestDeparse_Expr_Iff(t *testing.T) {
	assertRoundTrip(t, `SELECT IFF(a > 0, 'positive', 'non-positive') FROM t`)
}

func TestDeparse_Expr_IsNull(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE a IS NULL`)
}

func TestDeparse_Expr_IsNotNull(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE a IS NOT NULL`)
}

func TestDeparse_Expr_Between(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE a BETWEEN 1 AND 10`)
}

func TestDeparse_Expr_NotBetween(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE a NOT BETWEEN 1 AND 10`)
}

func TestDeparse_Expr_In(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE a IN (1, 2, 3)`)
}

func TestDeparse_Expr_NotIn(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE a NOT IN (1, 2, 3)`)
}

func TestDeparse_Expr_Like(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE a LIKE '%foo%'`)
}

func TestDeparse_Expr_ILike(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE a ILIKE '%foo%'`)
}

func TestDeparse_Expr_FunctionCall(t *testing.T) {
	assertRoundTrip(t, `SELECT COUNT(*) FROM t`)
}

func TestDeparse_Expr_FunctionCallDistinct(t *testing.T) {
	assertRoundTrip(t, `SELECT COUNT(DISTINCT a) FROM t`)
}

func TestDeparse_Expr_WindowFunction(t *testing.T) {
	assertRoundTrip(t, `SELECT ROW_NUMBER() OVER (PARTITION BY a ORDER BY b ASC) FROM t`)
}

func TestDeparse_Expr_WindowFunctionFrame(t *testing.T) {
	assertRoundTrip(t, `SELECT SUM(a) OVER (ORDER BY b ASC ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t`)
}

func TestDeparse_Expr_JsonAccess(t *testing.T) {
	assertRoundTrip(t, `SELECT v:key FROM t`)
}

func TestDeparse_Expr_BracketAccess(t *testing.T) {
	assertRoundTrip(t, `SELECT v[0] FROM t`)
}

func TestDeparse_Expr_ArrayLiteral(t *testing.T) {
	assertRoundTrip(t, `SELECT [1, 2, 3]`)
}

func TestDeparse_Expr_Subquery(t *testing.T) {
	assertRoundTrip(t, `SELECT (SELECT MAX(a) FROM t) AS mx`)
}

func TestDeparse_Expr_Exists(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM t WHERE EXISTS (SELECT 1 FROM s WHERE s.id = t.id)`)
}

func TestDeparse_Expr_Paren(t *testing.T) {
	assertRoundTrip(t, `SELECT (a + b) * c FROM t`)
}

func TestDeparse_Expr_Concat(t *testing.T) {
	assertRoundTrip(t, `SELECT a || b FROM t`)
}

func TestDeparse_Expr_QuotedIdent(t *testing.T) {
	assertRoundTrip(t, `SELECT "My Column" FROM "My Table"`)
}

func TestDeparse_Expr_ThreePartName(t *testing.T) {
	assertRoundTrip(t, `SELECT * FROM db.schema.table`)
}

func TestDeparse_Expr_StringLiteral(t *testing.T) {
	assertRoundTrip(t, `SELECT 'hello world'`)
}

func TestDeparse_Expr_StringWithEscape(t *testing.T) {
	assertRoundTrip(t, `SELECT 'it''s a test'`)
}

func TestDeparse_Expr_NullLiteral(t *testing.T) {
	assertRoundTrip(t, `SELECT NULL`)
}

func TestDeparse_Expr_BoolLiterals(t *testing.T) {
	assertRoundTrip(t, `SELECT TRUE, FALSE`)
}

// ---------------------------------------------------------------------------
// INSERT tests
// ---------------------------------------------------------------------------

func TestDeparse_Insert_Values(t *testing.T) {
	assertRoundTrip(t, `INSERT INTO t (a, b) VALUES (1, 2)`)
}

func TestDeparse_Insert_MultiRow(t *testing.T) {
	assertRoundTrip(t, `INSERT INTO t (a, b) VALUES (1, 2), (3, 4), (5, 6)`)
}

func TestDeparse_Insert_Select(t *testing.T) {
	assertRoundTrip(t, `INSERT INTO t SELECT a, b FROM s`)
}

func TestDeparse_Insert_Overwrite(t *testing.T) {
	assertRoundTrip(t, `INSERT OVERWRITE INTO t VALUES (1)`)
}

func TestDeparse_Insert_Multi_All(t *testing.T) {
	assertRoundTrip(t, `INSERT ALL INTO t1 VALUES (1) INTO t2 VALUES (2) SELECT * FROM dual`)
}

func TestDeparse_Insert_Multi_First(t *testing.T) {
	assertRoundTrip(t, `INSERT FIRST WHEN id = 1 THEN INTO t1 VALUES (id) ELSE INTO t2 VALUES (id) SELECT id FROM src`)
}

// ---------------------------------------------------------------------------
// UPDATE tests
// ---------------------------------------------------------------------------

func TestDeparse_Update_Simple(t *testing.T) {
	assertRoundTrip(t, `UPDATE t SET a = 1 WHERE id = 42`)
}

func TestDeparse_Update_MultiSet(t *testing.T) {
	assertRoundTrip(t, `UPDATE t SET a = 1, b = 2`)
}

func TestDeparse_Update_WithAlias(t *testing.T) {
	assertRoundTrip(t, `UPDATE t AS tbl SET tbl.a = 1`)
}

func TestDeparse_Update_FromClause(t *testing.T) {
	assertRoundTrip(t, `UPDATE t SET t.a = s.a FROM s WHERE t.id = s.id`)
}

// ---------------------------------------------------------------------------
// DELETE tests
// ---------------------------------------------------------------------------

func TestDeparse_Delete_Simple(t *testing.T) {
	assertRoundTrip(t, `DELETE FROM t WHERE id = 1`)
}

func TestDeparse_Delete_WithAlias(t *testing.T) {
	assertRoundTrip(t, `DELETE FROM t AS tbl WHERE tbl.a = 1`)
}

func TestDeparse_Delete_Using(t *testing.T) {
	assertRoundTrip(t, `DELETE FROM t USING s WHERE t.id = s.id`)
}

// ---------------------------------------------------------------------------
// MERGE tests
// ---------------------------------------------------------------------------

func TestDeparse_Merge_BasicUpdate(t *testing.T) {
	assertRoundTrip(t, `MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.v = s.v`)
}

func TestDeparse_Merge_InsertWhenNotMatched(t *testing.T) {
	assertRoundTrip(t, `MERGE INTO target t USING source s ON t.id = s.id WHEN NOT MATCHED THEN INSERT (id, v) VALUES (s.id, s.v)`)
}

func TestDeparse_Merge_Delete(t *testing.T) {
	assertRoundTrip(t, `MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED AND s.active = 0 THEN DELETE`)
}

func TestDeparse_Merge_UpdateAllByName(t *testing.T) {
	out := assertRoundTrip(t, `MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN UPDATE ALL BY NAME`)
	require.Contains(t, out, "UPDATE ALL BY NAME")
}

func TestDeparse_Merge_InsertAllByName(t *testing.T) {
	out := assertRoundTrip(t, `MERGE INTO target t USING source s ON t.id = s.id WHEN NOT MATCHED THEN INSERT ALL BY NAME`)
	require.Contains(t, out, "INSERT ALL BY NAME")
}

func TestDeparse_Merge_BothAllByName(t *testing.T) {
	assertRoundTrip(t, `MERGE INTO target t USING source s ON t.id = s.id WHEN MATCHED THEN UPDATE ALL BY NAME WHEN NOT MATCHED THEN INSERT ALL BY NAME`)
}

// ---------------------------------------------------------------------------
// CREATE TABLE tests
// ---------------------------------------------------------------------------

func TestDeparse_CreateTable_Simple(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT, name VARCHAR(100))`)
}

func TestDeparse_CreateTable_OrReplace(t *testing.T) {
	assertRoundTrip(t, `CREATE OR REPLACE TABLE t (id INT)`)
}

// OR ALTER must round-trip distinctly from OR REPLACE: a deparse that emitted
// OR REPLACE (or dropped the modifier) would change the reparsed AST and fail
// the DeepEqual in assertRoundTrip.
func TestDeparse_CreateTable_OrAlter(t *testing.T) {
	out := assertRoundTrip(t, `CREATE OR ALTER TABLE t (id INT)`)
	if !strings.Contains(out, "OR ALTER") {
		t.Errorf("deparsed SQL missing OR ALTER: %q", out)
	}
}

func TestDeparse_CreateDatabase_OrAlter(t *testing.T) {
	assertRoundTrip(t, `CREATE OR ALTER DATABASE db1`)
}

func TestDeparse_CreateSchema_OrAlter(t *testing.T) {
	assertRoundTrip(t, `CREATE OR ALTER SCHEMA s1`)
}

func TestDeparse_CreateView_OrAlter(t *testing.T) {
	assertRoundTrip(t, `CREATE OR ALTER VIEW v2 (one) AS SELECT a FROM my_table`)
}

func TestDeparse_CreateTable_Transient(t *testing.T) {
	assertRoundTrip(t, `CREATE TRANSIENT TABLE t (id INT)`)
}

// HYBRID must round-trip distinctly: a deparse that dropped the modifier would
// change the reparsed AST (Hybrid false vs true) and fail the DeepEqual.
func TestDeparse_CreateTable_Hybrid(t *testing.T) {
	out := assertRoundTrip(t, `CREATE HYBRID TABLE t (id INT)`)
	if !strings.Contains(out, "HYBRID TABLE") {
		t.Errorf("deparsed SQL missing HYBRID TABLE: %q", out)
	}
}

func TestDeparse_CreateTable_HybridOrReplace(t *testing.T) {
	out := assertRoundTrip(t, `CREATE OR REPLACE HYBRID TABLE t (id INT)`)
	if !strings.Contains(out, "OR REPLACE HYBRID TABLE") {
		t.Errorf("deparsed SQL missing OR REPLACE HYBRID TABLE: %q", out)
	}
}

// Inline INDEX elements round-trip: name + parenthesized column list.
func TestDeparse_CreateTable_HybridIndexSingle(t *testing.T) {
	out := assertRoundTrip(t,
		`CREATE HYBRID TABLE t (id INT, full_name VARCHAR(255), INDEX index_full_name (full_name))`)
	if !strings.Contains(out, "INDEX index_full_name (full_name)") {
		t.Errorf("deparsed SQL missing INDEX element: %q", out)
	}
}

func TestDeparse_CreateTable_HybridIndexMulti(t *testing.T) {
	out := assertRoundTrip(t,
		`CREATE HYBRID TABLE t (a INT, b INT, c INT, INDEX idx_abc (a, b, c))`)
	if !strings.Contains(out, "INDEX idx_abc (a, b, c)") {
		t.Errorf("deparsed SQL missing multi-column INDEX element: %q", out)
	}
}

// Full corpus example_01 shape round-trips (AUTOINCREMENT/PK/UNIQUE/VARIANT +
// inline INDEX in one statement).
func TestDeparse_CreateTable_HybridCorpus01(t *testing.T) {
	assertRoundTrip(t, `CREATE HYBRID TABLE mytable (customer_id INT AUTOINCREMENT PRIMARY KEY, full_name VARCHAR(255), email VARCHAR(255) UNIQUE, extended_customer_info VARIANT, INDEX index_full_name (full_name))`)
}

// INDEX element alongside an out-of-line constraint round-trips.
func TestDeparse_CreateTable_HybridIndexWithConstraint(t *testing.T) {
	assertRoundTrip(t, `CREATE OR REPLACE HYBRID TABLE ht2 (c1 INT PRIMARY KEY, c2 VARCHAR(10), INDEX idx_c2 (c2))`)
}

func TestDeparse_CreateTable_Temporary(t *testing.T) {
	assertRoundTrip(t, `CREATE TEMPORARY TABLE t (id INT)`)
}

func TestDeparse_CreateTable_IfNotExists(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE IF NOT EXISTS t (id INT)`)
}

func TestDeparse_CreateTable_NotNull(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT NOT NULL)`)
}

func TestDeparse_CreateTable_Default(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT DEFAULT 0)`)
}

func TestDeparse_CreateTable_Identity(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT IDENTITY(1, 1))`)
}

func TestDeparse_CreateTable_IdentityOrder(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT IDENTITY ORDER)`)
}

func TestDeparse_CreateTable_PrimaryKey(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR)`)
}

func TestDeparse_CreateTable_ForeignKey(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT FOREIGN KEY REFERENCES other (id))`)
}

func TestDeparse_CreateTable_Like(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t LIKE other`)
}

func TestDeparse_CreateTable_CTAS(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t AS SELECT * FROM s`)
}

func TestDeparse_CreateTable_CTASWithColumns(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (a, b) AS SELECT x, y FROM s`)
}

func TestDeparse_CreateTable_ClusterBy(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT) CLUSTER BY (id)`)
}

func TestDeparse_CreateTable_ClusterByLinear(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT) CLUSTER BY LINEAR (id)`)
}

func TestDeparse_CreateTable_Clone(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t2 CLONE t1`)
}

func TestDeparse_CreateTable_CloneAtTimestamp(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t2 CLONE t1 AT (TIMESTAMP => '2024-01-01')`)
}

func TestDeparse_CreateTable_CloneBeforeStatement(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t2 CLONE t1 BEFORE (STATEMENT => '8e5d0ca1-e866')`)
}

func TestDeparse_CreateTable_TableConstraint_PK(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT, name VARCHAR, PRIMARY KEY (id))`)
}

func TestDeparse_CreateTable_TableConstraint_Named(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT, CONSTRAINT pk_t PRIMARY KEY (id))`)
}

// ---------------------------------------------------------------------------
// ALTER TABLE tests
// ---------------------------------------------------------------------------

func TestDeparse_AlterTable_Rename(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t1 RENAME TO t2`)
}

func TestDeparse_AlterTable_SwapWith(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t1 SWAP WITH t2`)
}

func TestDeparse_AlterTable_AddColumn(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t ADD COLUMN c INT`)
}

func TestDeparse_AlterTable_DropColumn(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t DROP COLUMN c`)
}

func TestDeparse_AlterTable_RenameColumn(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t RENAME COLUMN old_name TO new_name`)
}

func TestDeparse_AlterTable_AlterColumnType(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t ALTER COLUMN c SET DATA TYPE VARCHAR(200)`)
}

func TestDeparse_AlterTable_AlterColumnNotNull(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t ALTER COLUMN c SET NOT NULL`)
}

func TestDeparse_AlterTable_AlterColumnDropNotNull(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t ALTER COLUMN c DROP NOT NULL`)
}

func TestDeparse_AlterTable_DropPK(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t DROP PRIMARY KEY`)
}

func TestDeparse_AlterTable_AddPK(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t ADD PRIMARY KEY (id)`)
}

func TestDeparse_AlterTable_DropConstraint(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t DROP CONSTRAINT fk_t_other`)
}

func TestDeparse_AlterTable_ClusterBy(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t CLUSTER BY (a, b)`)
}

func TestDeparse_AlterTable_DropClusterKey(t *testing.T) {
	assertRoundTrip(t, `ALTER TABLE t DROP CLUSTERING KEY`)
}

// ---------------------------------------------------------------------------
// DATABASE / SCHEMA DDL tests
// ---------------------------------------------------------------------------

func TestDeparse_CreateDatabase(t *testing.T) {
	assertRoundTrip(t, `CREATE DATABASE mydb`)
}

func TestDeparse_CreateDatabase_Transient(t *testing.T) {
	assertRoundTrip(t, `CREATE TRANSIENT DATABASE mydb`)
}

func TestDeparse_CreateDatabase_OrReplace(t *testing.T) {
	assertRoundTrip(t, `CREATE OR REPLACE DATABASE mydb`)
}

func TestDeparse_AlterDatabase_Rename(t *testing.T) {
	assertRoundTrip(t, `ALTER DATABASE db1 RENAME TO db2`)
}

func TestDeparse_AlterDatabase_SwapWith(t *testing.T) {
	assertRoundTrip(t, `ALTER DATABASE db1 SWAP WITH db2`)
}

func TestDeparse_DropDatabase(t *testing.T) {
	assertRoundTrip(t, `DROP DATABASE mydb`)
}

func TestDeparse_DropDatabase_IfExists(t *testing.T) {
	assertRoundTrip(t, `DROP DATABASE IF EXISTS mydb`)
}

func TestDeparse_UndropDatabase(t *testing.T) {
	assertRoundTrip(t, `UNDROP DATABASE mydb`)
}

func TestDeparse_CreateSchema(t *testing.T) {
	assertRoundTrip(t, `CREATE SCHEMA myschema`)
}

func TestDeparse_CreateSchema_ManagedAccess(t *testing.T) {
	assertRoundTrip(t, `CREATE SCHEMA myschema WITH MANAGED ACCESS`)
}

func TestDeparse_AlterSchema_Rename(t *testing.T) {
	assertRoundTrip(t, `ALTER SCHEMA s1 RENAME TO s2`)
}

func TestDeparse_AlterSchema_EnableManagedAccess(t *testing.T) {
	assertRoundTrip(t, `ALTER SCHEMA s1 ENABLE MANAGED ACCESS`)
}

func TestDeparse_DropSchema(t *testing.T) {
	assertRoundTrip(t, `DROP SCHEMA myschema`)
}

func TestDeparse_UndropSchema(t *testing.T) {
	assertRoundTrip(t, `UNDROP SCHEMA myschema`)
}

// ---------------------------------------------------------------------------
// DROP / UNDROP (non-database/schema)
// ---------------------------------------------------------------------------

func TestDeparse_DropTable(t *testing.T) {
	assertRoundTrip(t, `DROP TABLE t`)
}

func TestDeparse_DropTable_IfExists(t *testing.T) {
	assertRoundTrip(t, `DROP TABLE IF EXISTS t`)
}

func TestDeparse_DropTable_Cascade(t *testing.T) {
	assertRoundTrip(t, `DROP TABLE t CASCADE`)
}

func TestDeparse_DropView(t *testing.T) {
	assertRoundTrip(t, `DROP VIEW v`)
}

func TestDeparse_DropMaterializedView(t *testing.T) {
	assertRoundTrip(t, `DROP MATERIALIZED VIEW mv`)
}

func TestDeparse_UndropTable(t *testing.T) {
	assertRoundTrip(t, `UNDROP TABLE t`)
}

func TestDeparse_DropObjectVariants(t *testing.T) {
	// Each newly-supported DROP object variant must round-trip exactly via
	// DropObjectKind.String(); multi-word kinds reproduce the spaced spelling.
	cases := []string{
		`DROP ALERT a`,
		`DROP CONNECTION c`,
		`DROP FAILOVER GROUP fg`,
		`DROP INTEGRATION i`,
		`DROP MANAGED ACCOUNT ma`,
		`DROP MASKING POLICY mp`,
		`DROP NETWORK POLICY np`,
		`DROP REPLICATION GROUP rg`,
		`DROP RESOURCE MONITOR rm`,
		`DROP ROW ACCESS POLICY rap`,
		`DROP SESSION POLICY sp`,
		`DROP SHARE s`,
		`DROP USER u`,
		`DROP CONNECTION IF EXISTS c`,
		`DROP MASKING POLICY IF EXISTS mp`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertRoundTrip(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// VIEW tests
// ---------------------------------------------------------------------------

func TestDeparse_CreateView_Simple(t *testing.T) {
	assertRoundTrip(t, `CREATE VIEW v AS SELECT * FROM t`)
}

func TestDeparse_CreateView_OrReplace(t *testing.T) {
	assertRoundTrip(t, `CREATE OR REPLACE VIEW v AS SELECT a, b FROM t`)
}

func TestDeparse_CreateView_Secure(t *testing.T) {
	assertRoundTrip(t, `CREATE SECURE VIEW v AS SELECT * FROM t`)
}

func TestDeparse_CreateView_WithColumns(t *testing.T) {
	assertRoundTrip(t, `CREATE VIEW v (a, b) AS SELECT x, y FROM t`)
}

func TestDeparse_AlterView_Rename(t *testing.T) {
	assertRoundTrip(t, `ALTER VIEW v1 RENAME TO v2`)
}

func TestDeparse_AlterView_SetSecure(t *testing.T) {
	assertRoundTrip(t, `ALTER VIEW v SET SECURE`)
}

func TestDeparse_AlterView_UnsetSecure(t *testing.T) {
	assertRoundTrip(t, `ALTER VIEW v UNSET SECURE`)
}

func TestDeparse_CreateMaterializedView_Simple(t *testing.T) {
	assertRoundTrip(t, `CREATE MATERIALIZED VIEW mv AS SELECT * FROM t`)
}

func TestDeparse_CreateMaterializedView_ClusterBy(t *testing.T) {
	assertRoundTrip(t, `CREATE MATERIALIZED VIEW mv CLUSTER BY (a) AS SELECT a, b FROM t`)
}

func TestDeparse_CreateInteractiveMaterializedView(t *testing.T) {
	assertRoundTrip(t, `CREATE INTERACTIVE MATERIALIZED VIEW mv AS SELECT * FROM t`)
}

func TestDeparse_CreateInteractiveMaterializedView_OrReplaceIfNotExists(t *testing.T) {
	assertRoundTrip(t, `CREATE OR REPLACE INTERACTIVE MATERIALIZED VIEW IF NOT EXISTS mv AS SELECT a FROM t`)
}

func TestDeparse_AlterMaterializedView_Rename(t *testing.T) {
	assertRoundTrip(t, `ALTER MATERIALIZED VIEW mv RENAME TO mv2`)
}

func TestDeparse_AlterMaterializedView_Suspend(t *testing.T) {
	assertRoundTrip(t, `ALTER MATERIALIZED VIEW mv SUSPEND`)
}

func TestDeparse_AlterMaterializedView_Resume(t *testing.T) {
	assertRoundTrip(t, `ALTER MATERIALIZED VIEW mv RESUME`)
}

func TestDeparse_AlterMaterializedView_DropClusterKey(t *testing.T) {
	assertRoundTrip(t, `ALTER MATERIALIZED VIEW mv DROP CLUSTERING KEY`)
}

// ---------------------------------------------------------------------------
// DeparseFile with multiple statements
// ---------------------------------------------------------------------------

func TestDeparse_DeparseFile_MultipleStatements(t *testing.T) {
	sql := `CREATE TABLE t (id INT);
SELECT * FROM t;
DROP TABLE t`
	file, err := parser.Parse(sql)
	require.NoError(t, err)
	out, err := deparse.DeparseFile(file)
	require.NoError(t, err)
	require.Contains(t, out, "CREATE TABLE")
	require.Contains(t, out, "SELECT")
	require.Contains(t, out, "DROP TABLE")
}

// ---------------------------------------------------------------------------
// Error handling
// ---------------------------------------------------------------------------

func TestDeparse_UnsupportedNode(t *testing.T) {
	// A File node is not a statement — calling Deparse on it directly
	// should return an error, not panic.
	f := &ast.File{}
	_, err := deparse.Deparse(f)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported")
}

// ---------------------------------------------------------------------------
// $-reference deparse (positional $N + named $var, optional qualifier).
// Round-trip is the structural correctness gate for the DollarRef node.
// ---------------------------------------------------------------------------

func TestDeparse_DollarPositional(t *testing.T) {
	assertRoundTrip(t, `SELECT $1, $2 FROM t`)
}

func TestDeparse_DollarQualified(t *testing.T) {
	out := assertRoundTrip(t, `SELECT d.$1 FROM t AS d`)
	require.Contains(t, out, "d.$1")
}

func TestDeparse_DollarNamedVar(t *testing.T) {
	assertRoundTrip(t, `SELECT 2 * $min FROM t`)
}

func TestDeparse_DollarColonCast(t *testing.T) {
	assertRoundTrip(t, `SELECT $1:num::NUMBER AS num FROM t`)
}

func TestDeparse_DollarFuncArg(t *testing.T) {
	out := assertRoundTrip(t, `SELECT TO_DATE($1) FROM t`)
	require.Contains(t, out, "$1")
}

// ---------------------------------------------------------------------------
// Named (keyword) function arguments  name => value
// ---------------------------------------------------------------------------

func TestDeparse_NamedArgSingle(t *testing.T) {
	out := assertRoundTrip(t, `SELECT f(LOCATION => '@s') FROM t`)
	require.Contains(t, out, "LOCATION => '@s'")
}

func TestDeparse_NamedArgMultiple(t *testing.T) {
	out := assertRoundTrip(t, `SELECT f(A => 1, B => 'x', C => TRUE) FROM t`)
	require.Contains(t, out, "A => 1")
	require.Contains(t, out, "B => 'x'")
	require.Contains(t, out, "C => TRUE")
}

func TestDeparse_NamedArgMixedPositionalThenNamed(t *testing.T) {
	out := assertRoundTrip(t, `SELECT f(1, x, TYPE => 'STREAMING') FROM t`)
	require.Contains(t, out, "f(1, x, TYPE => 'STREAMING')")
}

func TestDeparse_NamedArgValueIsFuncCall(t *testing.T) {
	out := assertRoundTrip(t, `SELECT f(TS => TO_TIMESTAMP_TZ('x', 'fmt')) FROM t`)
	require.Contains(t, out, "TS => TO_TIMESTAMP_TZ('x', 'fmt')")
}

func TestDeparse_NamedArgValueIsCast(t *testing.T) {
	out := assertRoundTrip(t, `SELECT f(N => $1::NUMBER) FROM t`)
	require.Contains(t, out, "N => $1::NUMBER")
}

// Table-function named args (INFER_SCHEMA / DATA_SOURCE / FLATTEN). The TABLE()
// wrapper is normalized away by the table-function deparse, so we assert the
// named `=>` argument survives rather than a full structural round-trip.
func TestDeparse_NamedArgInferSchema(t *testing.T) {
	f, err := parser.Parse(`SELECT * FROM TABLE(INFER_SCHEMA(LOCATION => '@s', FILE_FORMAT => 'fmt'))`)
	require.NoError(t, err)
	out, err := deparse.DeparseFile(f)
	require.NoError(t, err)
	require.Contains(t, out, "LOCATION => '@s'")
	require.Contains(t, out, "FILE_FORMAT => 'fmt'")
}

func TestDeparse_NamedArgFlatten(t *testing.T) {
	f, err := parser.Parse(`SELECT * FROM persons p, LATERAL FLATTEN(INPUT => p.c, PATH => 'contact') f`)
	require.NoError(t, err)
	out, err := deparse.DeparseFile(f)
	require.NoError(t, err)
	require.Contains(t, out, "INPUT => p.c")
	require.Contains(t, out, "PATH => 'contact'")
}

// ---------------------------------------------------------------------------
// CLONE ... { AT | BEFORE } ( kind => value ) time-travel
// ---------------------------------------------------------------------------

func TestDeparse_CloneAtStringValue(t *testing.T) {
	out := assertRoundTrip(t, `CREATE TABLE t2 CLONE t1 AT (TIMESTAMP => '2024-01-01')`)
	require.Contains(t, out, "AT (TIMESTAMP => '2024-01-01')")
}

func TestDeparse_CloneAtFuncValue(t *testing.T) {
	out := assertRoundTrip(t, `CREATE TABLE t2 CLONE t1 AT (TIMESTAMP => TO_TIMESTAMP_TZ('04/05/2013 01:02:03', 'mm/dd/yyyy hh24:mi:ss'))`)
	require.Contains(t, out, "AT (TIMESTAMP => TO_TIMESTAMP_TZ(")
}

func TestDeparse_CloneBeforeStatementValue(t *testing.T) {
	out := assertRoundTrip(t, `CREATE TABLE t2 CLONE t1 BEFORE (STATEMENT => '8e5d0ca1-e866-44fa-843b-5e6ad35e8bb7')`)
	require.Contains(t, out, "BEFORE (STATEMENT => '8e5d0ca1-e866-44fa-843b-5e6ad35e8bb7')")
}

// ---------------------------------------------------------------------------
// VALUES table source + $N table ref + ->> result-pipe (gap-from-values)
// ---------------------------------------------------------------------------

func TestDeparse_ValuesTableSource(t *testing.T) {
	out := assertRoundTrip(t, `SELECT * FROM (VALUES (1, 'one'), (2, 'two'), (3, 'three'))`)
	require.Contains(t, out, "(VALUES (1, 'one'), (2, 'two'), (3, 'three'))")
}

func TestDeparse_ValuesTableSourceDerivedColumnList(t *testing.T) {
	out := assertRoundTrip(t, `SELECT c1, c2 FROM (VALUES (1, 'one'), (2, 'two')) AS v1 (c1, c2)`)
	require.Contains(t, out, "AS v1 (c1, c2)")
}

func TestDeparse_ValuesTableSourceJoin(t *testing.T) {
	assertRoundTrip(t,
		`SELECT v1.$2, v2.$2 FROM (VALUES (1, 'one'), (2, 'two')) AS v1 `+
			`INNER JOIN (VALUES (1, 'One'), (3, 'three')) AS v2 WHERE v2.$1 = v1.$1`)
}

func TestDeparse_DollarTableRef(t *testing.T) {
	out := assertRoundTrip(t, `SELECT * FROM $1`)
	require.Contains(t, out, "FROM $1")
}

func TestDeparse_DollarTableRefAlias(t *testing.T) {
	out := assertRoundTrip(t, `SELECT d.x FROM $1 AS d`)
	require.Contains(t, out, "FROM $1 AS d")
}

func TestDeparse_ResultPipeGeneral(t *testing.T) {
	out := assertRoundTrip(t, `SELECT 1 ->> SELECT * FROM $1`)
	require.Contains(t, out, "->>")
	require.Contains(t, out, "FROM $1")
}
