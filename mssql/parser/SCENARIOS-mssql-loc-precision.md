# MSSQL Loc Precision Scenarios

> Goal: every AST node used for source extraction has a valid byte range whose slice is the complete original SQL fragment for that node.
> Verification: parser tests compare `sql[Loc.Start:Loc.End]` with the exact expected fragment, and traversal guards assert valid child locations stay in range.
> Status: [ ] pending, [x] covered by current precision tests, [~] partially covered by representative or statement-level tests.

This matrix is stricter than `SCENARIOS-mssql-loc.md`: file-level migration from `p.pos()` to `p.prevEnd()` is not enough. A scenario only counts when the exact node slice is verified.

---

## Phase 1: Loc Access Contract

### 1.1 AST Loc Helpers

- [x] `Loc{Start: 1, End: 3}` is valid.
- [x] `Loc{Start: 1, End: -1}` is invalid.
- [x] Merging `{4,8}` with `{1,5}` returns `{1,8}`.
- [x] `NodeLoc` returns `BinaryExpr.Loc`.
- [x] `NodeLoc` returns `FuncCallExpr.Loc`.
- [x] `NodeLoc` returns `ColumnRef.Loc`.
- [x] `NodeLoc` returns `LikeExpr.Loc`.
- [x] Every AST struct with a `Loc Loc` field is covered by `NodeLoc` either directly or through the generic Loc-field fallback.
- [x] `SpanNodes` ignores nil and invalid nodes while spanning valid nodes.

### 1.2 Universal Loc Invariants

- [x] For every exact-slice scenario, `Loc.Start >= 0`.
- [x] For every exact-slice scenario, `Loc.End >= 0`.
- [x] For every exact-slice scenario, `Loc.Start <= Loc.End`.
- [x] For every exact-slice scenario, `Loc.End <= len(sql)`.
- [x] For every exact-slice scenario, each child node slice is contained inside its parent node slice.
- [x] Leading whitespace before a statement is not included in the statement node slice.
- [x] Trailing semicolon behavior is consistent across statement nodes.
- [x] Multi-line source preserves byte offsets across newline boundaries.
- [x] Unicode source keeps byte-offset columns, not rune columns.

---

## Phase 2: Core Expression Exact Slices

### 2.1 Comparison And Predicate Expressions

- [x] `SELECT c > 0` slices `BinaryExpr` as `c > 0`.
- [x] `SELECT c IN (1, 2)` slices `InExpr` as `c IN (1, 2)`.
- [x] `SELECT c BETWEEN 1 AND 2` slices `BetweenExpr` as `c BETWEEN 1 AND 2`.
- [x] `SELECT c LIKE 'x%'` slices `LikeExpr` as `c LIKE 'x%'`.
- [x] `SELECT c IS NOT NULL` slices `IsExpr` as `c IS NOT NULL`.
- [x] `SELECT c NOT IN (1, 2)` slices `InExpr` as `c NOT IN (1, 2)`.
- [x] `SELECT c NOT BETWEEN 1 AND 2` slices `BetweenExpr` as `c NOT BETWEEN 1 AND 2`.
- [x] `SELECT c NOT LIKE 'x%' ESCAPE '\'` slices `LikeExpr` as `c NOT LIKE 'x%' ESCAPE '\'`.
- [x] `SELECT * FROM t WHERE c = ANY (SELECT x FROM s)` slices `SubqueryComparisonExpr` as `c = ANY (SELECT x FROM s)`.
- [x] `SELECT c IN (SELECT x FROM s)` slices `InExpr` as `c IN (SELECT x FROM s)`.

### 2.2 Boolean And Arithmetic Expressions

- [x] `SELECT a + b * c` slices the outer `BinaryExpr` as `a + b * c`.
- [x] `SELECT (a = 1) AND (b = 2)` slices the outer `BinaryExpr` as `(a = 1) AND (b = 2)`.
- [x] `SELECT a OR b AND c` slices the outer `BinaryExpr` as `a OR b AND c`.
- [x] `SELECT a - b + c` slices the outer `BinaryExpr` as `a - b + c`.
- [x] `SELECT a / b % c` slices the outer `BinaryExpr` as `a / b % c`.
- [x] `SELECT -c` slices `UnaryExpr` as `-c`.
- [x] `SELECT NOT (c > 0)` slices `UnaryExpr` as `NOT (c > 0)`.
- [x] `SELECT (a + b)` slices `ParenExpr` as `(a + b)`.

### 2.3 Postfix Expressions

- [x] `SELECT name COLLATE Latin1_General_CI_AS LIKE 'A%'` slices `LikeExpr` as `name COLLATE Latin1_General_CI_AS LIKE 'A%'`.
- [x] `SELECT dt AT TIME ZONE 'UTC'` slices `AtTimeZoneExpr` as `dt AT TIME ZONE 'UTC'`.
- [x] `SELECT name COLLATE database_default` slices `CollateExpr` as `name COLLATE database_default`.
- [x] `SELECT dt AT TIME ZONE 'UTC' AT TIME ZONE 'Pacific Standard Time'` slices the outer `AtTimeZoneExpr` as the full chain.

### 2.4 Function And Built-in Expressions

- [x] `SELECT SYSDATETIME()` slices `FuncCallExpr` as `SYSDATETIME()`.
- [x] `SELECT COUNT(*)` slices `FuncCallExpr` as `COUNT(*)`.
- [x] `SELECT STRING_AGG(name, ',') WITHIN GROUP (ORDER BY name)` slices `FuncCallExpr` as the full `WITHIN GROUP` call.
- [x] `SELECT SUM(x) OVER (PARTITION BY y)` slices `FuncCallExpr` as the full windowed call.
- [x] `SELECT COUNT(DISTINCT c)` slices `FuncCallExpr` as `COUNT(DISTINCT c)`.
- [x] `SELECT CAST(c AS int)` slices `CastExpr` as `CAST(c AS int)`.
- [x] `SELECT CONVERT(varchar(20), c, 126)` slices `ConvertExpr` as `CONVERT(varchar(20), c, 126)`.
- [x] `SELECT TRY_CAST(c AS int)` slices `TryCastExpr` as `TRY_CAST(c AS int)`.
- [x] `SELECT TRY_CONVERT(int, c)` slices `TryConvertExpr` as `TRY_CONVERT(int, c)`.
- [x] `SELECT COALESCE(a, b, c)` slices `CoalesceExpr` as `COALESCE(a, b, c)`.
- [x] `SELECT NULLIF(a, b)` slices `NullifExpr` as `NULLIF(a, b)`.
- [x] `SELECT IIF(c > 0, 1, 0)` slices `IifExpr` as `IIF(c > 0, 1, 0)`.

### 2.5 CASE, Subquery, And Specialty Expressions

- [x] `SELECT CASE WHEN c > 0 THEN 1 ELSE 0 END` slices `CaseExpr` as the full CASE expression.
- [x] In `CASE WHEN c > 0 THEN 1 END`, `CaseWhen` slices as `WHEN c > 0 THEN 1`.
- [x] `SELECT EXISTS (SELECT 1 FROM t)` slices `ExistsExpr` as `EXISTS (SELECT 1 FROM t)`.
- [x] `SELECT (SELECT 1)` slices `SubqueryExpr` as `(SELECT 1)`.
- [x] `SELECT CURRENT OF cur` slices `CurrentOfExpr` as `CURRENT OF cur`.
- [x] `SELECT geography::Point(1, 2, 4326)` slices `MethodCallExpr` as `geography::Point(1, 2, 4326)`.
- [x] `GROUP BY GROUPING SETS ((a), (b))` slices `GroupingSetsExpr` as the full grouping sets expression.
- [x] `GROUP BY ROLLUP (a, b)` slices `RollupExpr` as `ROLLUP (a, b)`.
- [x] `GROUP BY CUBE (a, b)` slices `CubeExpr` as `CUBE (a, b)`.

---

## Phase 3: Names, References, And Types

### 3.1 Column And Variable References

- [x] `SELECT c` slices `ColumnRef` as `c`.
- [x] `SELECT t.c` slices `ColumnRef` as `t.c`.
- [x] `SELECT dbo.t.c` slices `ColumnRef` as `dbo.t.c`.
- [x] `SELECT t.*` slices `StarExpr` as `t.*`.
- [x] `SELECT db.dbo.t.c` slices `ColumnRef` as `db.dbo.t.c`.
- [x] `SELECT server.db.dbo.t.c` slices `ColumnRef` as `server.db.dbo.t.c`.
- [x] `SELECT @v` slices `VariableRef` as `@v`.
- [x] `SELECT @@ROWCOUNT` slices `VariableRef` as `@@ROWCOUNT`.

### 3.2 Table References

- [x] `SELECT * FROM t` slices `TableRef` as `t`.
- [x] `SELECT * FROM dbo.t` slices `TableRef` as `dbo.t`.
- [x] `SELECT * FROM db.dbo.t` slices `TableRef` as `db.dbo.t`.
- [x] `SELECT * FROM @t` slices `TableVarRef` as `@t`.
- [x] `SELECT * FROM @x.nodes('/r') n(c)` slices `TableVarMethodCallRef` as `@x.nodes('/r') n(c)`.
- [x] `SELECT * FROM t AS x` slices `TableRef` as `t AS x`.
- [x] `SELECT * FROM OPENROWSET(...)` slices table-valued `FuncCallExpr` or rowset table ref as the full rowset function.

### 3.3 Data Types

- [x] `CREATE TABLE t (c int)` slices `DataType` as `int`.
- [x] `CREATE TABLE t (c varchar(50))` slices `DataType` as `varchar(50)`.
- [x] `CREATE TABLE t (c decimal(10, 2))` slices `DataType` as `decimal(10, 2)`.
- [x] `CREATE TABLE t (c nvarchar(max))` slices `DataType` as `nvarchar(max)`.
- [x] `CREATE TABLE t (c dbo.MyType)` slices `DataType` as `dbo.MyType`.

---

## Phase 4: CREATE TABLE Metadata Extraction

### 4.1 Column Definitions

- [x] `CREATE TABLE t (c int)` slices `ColumnDef` as `c int`.
- [x] `CREATE TABLE t (c int NULL)` slices `NullableSpec` as `NULL`.
- [x] `CREATE TABLE t (c int NOT NULL)` slices `NullableSpec` as `NOT NULL`.
- [x] `CREATE TABLE t (id int IDENTITY(1,1))` slices `IdentitySpec` as `IDENTITY(1,1)`.
- [x] `CREATE TABLE t (c AS (a + b))` slices `ComputedColumnDef.Expr` as `a + b`.
- [x] `CREATE TABLE t (c AS (a + b) PERSISTED)` slices `ComputedColumnDef` as `AS (a + b) PERSISTED`.

### 4.2 DEFAULT And CHECK Constraints

- [x] `CREATE TABLE t (d datetime2 DEFAULT SYSDATETIME())` slices default `FuncCallExpr` as `SYSDATETIME()`.
- [x] `CREATE TABLE t (e int DEFAULT dbo.f())` slices default `FuncCallExpr` as `dbo.f()`.
- [x] `CREATE TABLE t (c int DEFAULT COUNT(*))` slices default `FuncCallExpr` as `COUNT(*)`.
- [x] `CREATE TABLE t (c int, CONSTRAINT ck CHECK (c > 0))` slices CHECK expression as `c > 0`.
- [x] `CREATE TABLE t (c int CHECK (c IN (1, 2)))` slices CHECK expression as `c IN (1, 2)`.
- [x] `CREATE TABLE t (c int CONSTRAINT df DEFAULT (0))` slices `ConstraintDef.Expr` as `0`.
- [x] `CREATE TABLE t (c int CONSTRAINT ck CHECK (c BETWEEN 1 AND 10))` slices `ConstraintDef.Expr` as `c BETWEEN 1 AND 10`.

### 4.3 Table Constraints And Options

- [x] Primary key constraint slices `ConstraintDef` as the full `CONSTRAINT ... PRIMARY KEY (...)` fragment.
- [x] Unique constraint slices `ConstraintDef` as the full `CONSTRAINT ... UNIQUE (...)` fragment.
- [x] Foreign key constraint slices `ConstraintDef` as the full `FOREIGN KEY ... REFERENCES ...` fragment.
- [x] Edge constraint slices `ConstraintDef` as the full edge constraint fragment.
- [x] Memory optimized table option slices `TableOption` as `MEMORY_OPTIMIZED = ON`.
- [x] System versioning option slices `TableOption` as the full `SYSTEM_VERSIONING = ON (...)` fragment.
- [x] Inline index slices `InlineIndexDef` as the full `INDEX ...` fragment.
- [x] Inline index column slices `IndexColumn` as the full column item.

---

## Phase 5: Index, Spatial, Fulltext, And Search Metadata

### 5.1 Index Definitions

- [x] `CREATE INDEX ix ON t (c)` slices `CreateIndexStmt` as the full statement.
- [x] `CREATE INDEX ix ON t (c DESC)` slices `IndexColumn` as `c DESC`.
- [x] `CREATE INDEX ix ON t (c) INCLUDE (d)` slices included `IndexColumn` as `d`.
- [x] `CREATE INDEX ix ON t (c) WHERE c > 0` slices filter predicate as `c > 0`.
- [x] `CREATE INDEX ix ON t (c) WITH (FILLFACTOR = 90)` preserves exact option text.
- [x] `ALTER INDEX ix ON t REBUILD WITH (ONLINE = ON)` preserves exact option text.

### 5.2 Spatial Index Definitions

- [x] `CREATE SPATIAL INDEX ... WITH (BOUNDING_BOX = (...))` preserves `BOUNDING_BOX=(...)`.
- [x] `CREATE SPATIAL INDEX ... WITH (GRIDS = (...))` preserves `GRIDS=(...)`.
- [x] `CREATE SPATIAL INDEX ... WITH (CELLS_PER_OBJECT = 16)` preserves `CELLS_PER_OBJECT=16`.
- [x] Spatial index column expression slices as the geometry/geography column reference.
- [x] Spatial index `USING GEOMETRY_GRID` slices the tessellation scheme correctly.
- [x] Spatial index `USING GEOGRAPHY_AUTO_GRID` slices the tessellation scheme correctly.

### 5.3 Fulltext And Search

- [x] `CREATE FULLTEXT INDEX ON t(c) KEY INDEX pk` slices `CreateFulltextIndexStmt` as the full statement.
- [x] Fulltext indexed column slices as the full column item.
- [x] `CONTAINS(c, 'term')` slices `FullTextPredicate` as `CONTAINS(c, 'term')`.
- [x] `FREETEXT(c, 'term')` slices `FullTextPredicate` as `FREETEXT(c, 'term')`.
- [x] `CONTAINSTABLE(t, c, 'term')` slices `FullTextTableRef` as the full table-valued search function.
- [x] `SEMANTICSIMILARITYTABLE(t, c, source_key)` slices `SemanticTableRef` as the full table-valued semantic function.

---

## Phase 6: SELECT And DML Clause Extraction

### 6.1 SELECT Clauses

- [x] `SELECT a AS x` slices `ResTarget` as `a AS x`.
- [x] `SELECT @v = a` slices `SelectAssign` as `@v = a`.
- [x] `SELECT TOP (10) a FROM t` slices `TopClause` as `TOP (10)`.
- [x] `SELECT a FROM t WHERE a > 0` slices `WhereClause` expression as `a > 0`.
- [x] `SELECT a FROM t ORDER BY a DESC` slices `OrderByItem` as `a DESC`.
- [x] `SELECT a FROM t OFFSET 10 ROWS FETCH NEXT 5 ROWS ONLY` slices `FetchClause` as `FETCH NEXT 5 ROWS ONLY`.
- [x] `SELECT a FROM t FOR XML PATH` slices `ForClause` as `FOR XML PATH`.
- [x] `WITH cte AS (SELECT 1) SELECT * FROM cte` slices `CommonTableExpr` as `cte AS (SELECT 1)`.
- [x] `WITH XMLNAMESPACES (...) SELECT ...` slices `XmlNamespaceDecl` as the full namespace declaration.

### 6.2 JOIN, WINDOW, GROUPING

- [x] `SELECT * FROM a JOIN b ON a.id = b.id` slices `JoinClause` as `a JOIN b ON a.id = b.id`.
- [x] `SELECT * FROM a CROSS APPLY fn(a.id)` slices `JoinClause` as `a CROSS APPLY fn(a.id)`.
- [x] `SELECT SUM(x) OVER (PARTITION BY y ORDER BY z ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)` slices `OverClause` as the full OVER clause.
- [x] Window frame slices `WindowFrame` as `ROWS BETWEEN 1 PRECEDING AND CURRENT ROW`.
- [x] Window bound slices `WindowBound` as `1 PRECEDING`.
- [x] Named window definition slices `WindowDef` as the full `WINDOW name AS (...)` fragment.

### 6.3 INSERT, UPDATE, DELETE, MERGE

- [x] `INSERT INTO t (a) VALUES (1)` slices `InsertStmt` as the full statement.
- [x] `VALUES (1), (2)` slices `ValuesClause` as the full values fragment.
- [x] `UPDATE t SET a = 1 WHERE b = 2` slices `UpdateStmt` as the full statement.
- [x] `DELETE FROM t WHERE a = 1` slices `DeleteStmt` as the full statement.
- [x] `OUTPUT inserted.a INTO audit` slices `OutputClause` as the full OUTPUT fragment.
- [x] `MERGE ... WHEN MATCHED THEN UPDATE ...` slices each `MergeWhenClause` as its full WHEN clause.
- [x] `MergeUpdateAction` slices as the full update action.
- [x] `MergeInsertAction` slices as the full insert action.
- [x] `MergeDeleteAction` slices as `DELETE`.

---

## Phase 7: Top-Level Statement And Public API Exactness

### 7.1 Statement Ranges

- [x] Every DML statement scenario has a valid top-level statement Loc.
- [x] Every DDL statement scenario has a valid top-level statement Loc.
- [x] Every security/admin statement scenario has a valid top-level statement Loc.
- [x] Multi-statement SQL produces non-overlapping statement ranges.
- [x] `GO` batch separators do not corrupt following statement ranges.
- [x] Empty statements between semicolons do not produce invalid Loc for real statements.

### 7.2 Public API Ranges

- [x] Public `mssql.Parse` returns `Statement.Text` matching `sql[ByteStart:ByteEnd]`.
- [x] Public `mssql.Parse` returns line/column start matching `ByteStart`.
- [x] Public `mssql.Parse` returns line/column end matching `ByteEnd`.
- [x] Public `mssql.Parse` keeps byte-based columns for tabs and Unicode content.
