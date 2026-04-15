package deparse_test

import (
	"reflect"
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

func TestDeparse_Select_Where(t *testing.T) {
	assertRoundTrip(t, `SELECT a FROM t WHERE a > 1`)
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

// ---------------------------------------------------------------------------
// CREATE TABLE tests
// ---------------------------------------------------------------------------

func TestDeparse_CreateTable_Simple(t *testing.T) {
	assertRoundTrip(t, `CREATE TABLE t (id INT, name VARCHAR(100))`)
}

func TestDeparse_CreateTable_OrReplace(t *testing.T) {
	assertRoundTrip(t, `CREATE OR REPLACE TABLE t (id INT)`)
}

func TestDeparse_CreateTable_Transient(t *testing.T) {
	assertRoundTrip(t, `CREATE TRANSIENT TABLE t (id INT)`)
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
