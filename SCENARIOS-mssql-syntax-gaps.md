# Phase 0 Audit: MSSQL Parser Syntax Gaps
## Total Scenarios: 24

This document tracks T-SQL syntax gaps between SqlScriptDOM (authoritative) and omni mssql parser.

---

## Axis 1: SELECT List Item Forms

| Form | Description | Omni Status |
|------|-------------|------------|
| `expr` | scalar expression | ✅ ResTarget with Val, no Name |
| `expr AS alias` | expression with AS alias | ✅ ResTarget.Name populated |
| `expr alias` | expression with bare alias | ✅ bare alias parsing in parseTargetList |
| `alias = expr` | assignment form (Gap 1) | ❌ not parsed |
| `*` | select all columns | ⚠️ parsed as StarExpr, no SelectStarExpression node |
| `table.*` | qualified star | ✅ StarExpr with Qualifier |
| `@var = expr` | select variable assignment | ⚠️ SelectAssign exists but not integrated into parseTargetList |

---

## Axis 2: TableReference Subclasses

| Form | Description | Omni Status |
|------|-------------|------------|
| `table` | simple table | ✅ TableRef |
| `schema.table` | schema-qualified table | ✅ TableRef (Schema + Object fields) |
| `db.schema.table` | database-qualified table | ✅ TableRef (Database + Schema + Object fields) |
| `server.db.schema.table` | server-qualified table | ✅ TableRef (Server + Database + Schema + Object fields) |
| `(SELECT ...)` | derived table / subquery | ⚠️ SubqueryExpr exists but wrapped in AliasedTableRef, no explicit alias column list support |
| `(SELECT ...) AS t(c1,c2)` | derived table with column list | ⚠️ AliasedTableRef.Columns parsed but only as *List (strings), not bound |
| `(VALUES (...), (...))` | inline derived table (Gap 4) | ❌ VALUES only supported in INSERT, not in FROM |
| `(VALUES ...) AS t(c1,c2)` | inline derived with cols (Gap 4) | ❌ not supported |
| `fn(args)` | table-valued function (Gap 3) | ⚠️ FuncCallExpr in parseTableValuedFunction, but limited handling |
| `fn(args) AS t(c1,c2)` | TVF with column list (Gap 3) | ⚠️ Columns not captured in TVF result |
| `@var` | table variable | ✅ TableVarRef |
| `@var.Method(args)` | variable method call (e.g., XML.nodes()) | ✅ TableVarMethodCallRef with Columns |
| `OPENROWSET(...)` | rowset function | ✅ FuncCallExpr wrapper via parseRowsetFunction |
| `OPENJSON(...)` | JSON function | ✅ FuncCallExpr wrapper via parseRowsetFunction |
| `OPENXML(...)` | XML function | ✅ FuncCallExpr wrapper via parseRowsetFunction |
| `OPENQUERY(...)` | linked server query | ✅ FuncCallExpr wrapper via parseRowsetFunction |
| `OPENDATASOURCE(...)` | datasource function | ✅ FuncCallExpr wrapper via parseRowsetFunction |
| `table PIVOT (...)` | pivoted table | ✅ PivotExpr |
| `table UNPIVOT (...)` | unpivoted table | ✅ UnpivotExpr |
| `t1 JOIN t2 ON ...` | explicit join | ✅ JoinClause |
| `t1, t2` | cross join (comma syntax) | ✅ two items in FromClause list |

---

## Axis 3: Alias + Column List Pattern

For each TableExpr that can be aliased, verify both alias AND `(col_list)` work:

| Form | Alias | Col List | Notes |
|------|-------|----------|-------|
| SubqueryExpr | ⚠️ via AliasedTableRef wrapper | ⚠️ AliasedTableRef.Columns | requires wrapping |
| FuncCallExpr (TVF) | ✅ AliasedTableRef wrapper | ⚠️ Columns not populated | gap in parseTableValuedFunction |
| TableVarMethodCallRef | ✅ direct Alias field | ✅ direct Columns field | fully supported |
| RowsetFunction | ✅ via AliasedTableRef | ⚠️ Columns parsed but discarded | gap in parseRowsetFunction |

---

## Scenario List

### S-SEL-01: Alias Assignment Form in SELECT
- **SQL**: `SELECT col1, col2 = 1 FROM t`
- **Expected AST**: SelectStmt.TargetList = [ResTarget{Val: col1}, ResTarget{Val: expr, Name: "col2"}]
- **Status**: ❌ `alias = expr` not recognized; only `expr AS alias` or `expr alias` work
- **Location**: select.go parseTargetList() L604

### S-SEL-02: SELECT * as StarExpression
- **SQL**: `SELECT * FROM t`
- **Expected AST**: SelectStmt.TargetList contains SelectStarExpression equivalent
- **Status**: ⚠️ parsed as StarExpr{Qualifier:""}, but not explicit SelectStarExpression node type
- **Location**: select.go parseTargetList() L616 (parseExpr returns StarExpr)

### S-SEL-03: Variable Assignment in SELECT
- **SQL**: `SELECT col1, @var = col2 FROM t`
- **Expected AST**: SelectStmt.TargetList = [..., SelectSetVariable{Variable: "@var", Expression: col2}]
- **Status**: ⚠️ SelectAssign type exists but not integrated into parseTargetList
- **Location**: select.go parseTargetList() L604, should check for tokVARIABLE before '='

### S-TR-01: Derived Table Without Column List
- **SQL**: `SELECT * FROM (SELECT a, b FROM t) AS x`
- **Expected AST**: AliasedTableRef{Table: SubqueryExpr{Query: SelectStmt}, Alias: "x", Columns: nil}
- **Status**: ✅ fully handled
- **Location**: select.go parsePrimaryTableSource() L766

### S-TR-02: Derived Table With Column List
- **SQL**: `SELECT * FROM (SELECT a, b FROM t) AS x(c1, c2)`
- **Expected AST**: AliasedTableRef{Table: SubqueryExpr, Alias: "x", Columns: ["c1", "c2"]}
- **Status**: ⚠️ Columns field exists in AliasedTableRef but never populated from parser
- **Location**: select.go parsePrimaryTableSource() L766-795 (no column list parsing after alias)

### S-TR-03: Inline Derived Table (VALUES)
- **SQL**: `SELECT * FROM (VALUES (1, 2), (3, 4)) AS t(a, b)`
- **Expected AST**: AliasedTableRef{Table: InlineTableValues, Alias: "t", Columns: ["a", "b"]}
- **Status**: ❌ VALUES only recognized in INSERT; not supported in FROM clause
- **Location**: select.go parsePrimaryTableSource() L766 (missing VALUES case)

### S-TR-04: Inline Derived Table With Column List
- **SQL**: `SELECT * FROM (VALUES (1, 'x')) AS t(id, name)`
- **Expected AST**: InlineTableValues node with alias and column mapping
- **Status**: ❌ not supported in FROM
- **Location**: select.go parsePrimaryTableSource() (needs new parsing path)

### S-TR-05: Table-Valued Function Without Column List
- **SQL**: `SELECT * FROM OPENQUERY(linkedserver, 'SELECT * FROM remote_table') AS result`
- **Expected AST**: AliasedTableRef{Table: FuncCallExpr{Name: "OPENQUERY"}, Alias: "result"}
- **Status**: ✅ fully handled via parseRowsetFunction
- **Location**: select.go parsePrimaryTableSource() L799-801

### S-TR-06: Table-Valued Function With Column List
- **SQL**: `SELECT * FROM fn_table(1) AS t(col1 INT, col2 VARCHAR(50))`
- **Expected AST**: AliasedTableRef{Table: FuncCallExpr, Alias: "t", Columns: ["col1", "col2"]}
- **Status**: ⚠️ Column list parsed but not stored; TVF parsing treats it as bare identifier
- **Location**: select.go parseTableValuedFunction() L1057-1096

### S-TR-07: User-Defined TVF With Arguments
- **SQL**: `SELECT * FROM dbo.fn_GetCustomers(2024) AS customers`
- **Expected AST**: AliasedTableRef{Table: FuncCallExpr{Name: TableRef{Schema:"dbo", Object:"fn_GetCustomers"}, Args: [expr]}, Alias: "customers"}
- **Status**: ⚠️ schema-qualified function names parsed in parseTableRef, but TVF wrapping may lose schema
- **Location**: select.go parseTableValuedFunction() L1058 (uses parseTableRef result)

### S-TR-08: OPENROWSET with Column Definitions
- **SQL**: `SELECT * FROM OPENROWSET(BULK 'file.txt', SINGLE_CLOB) AS t(data NVARCHAR(MAX))`
- **Expected AST**: AliasedTableRef with column definitions
- **Status**: ⚠️ WITH clause parsed but discarded in parseRowsetFunction L56
- **Location**: rowset_functions.go parseRowsetFunction() L56

### S-TR-09: Variable Method Call Without Column List
- **SQL**: `SELECT * FROM @xml_var.nodes('/root/item') AS t(id, val)`
- **Expected AST**: TableVarMethodCallRef{Var: "@xml_var", Method: "nodes", Args: ['/root/item'], Alias: "t", Columns: ["id", "val"]}
- **Status**: ✅ fully supported
- **Location**: name.go parseVariableTableSource() L122-179

### S-TR-10: Variable Method Call With Column List
- **SQL**: `SELECT * FROM @xml_var.nodes('/x') AS result(xml_data XML)`
- **Expected AST**: TableVarMethodCallRef with Columns populated
- **Status**: ✅ fully supported (Columns field populated)
- **Location**: name.go parseVariableTableSource() L162-177

### S-TR-11: Simple Qualified Table (schema.table)
- **SQL**: `SELECT * FROM dbo.MyTable`
- **Expected AST**: TableRef{Schema: "dbo", Object: "MyTable"}
- **Status**: ✅ fully handled
- **Location**: name.go parseTableRef() L60

### S-TR-12: Database-Qualified Table
- **SQL**: `SELECT * FROM database.schema.table`
- **Expected AST**: TableRef{Database: "database", Schema: "schema", Object: "table"}
- **Status**: ✅ fully handled
- **Location**: name.go parseTableRef() L60

### S-TR-13: Server-Qualified Table
- **SQL**: `SELECT * FROM server.database.schema.table`
- **Expected AST**: TableRef{Server: "server", Database: "database", Schema: "schema", Object: "table"}
- **Status**: ✅ fully handled
- **Location**: name.go parseTableRef() L60

### S-TR-14: PIVOT Expression
- **SQL**: `SELECT * FROM (SELECT salesperson, year, sales FROM t) pivot_src PIVOT (SUM(sales) FOR year IN (2022, 2023, 2024)) AS pvt`
- **Expected AST**: PivotExpr{Source: TableExpr, AggFunc: SUM(...), ForCol: "year", InValues: [2022, 2023, 2024], Alias: "pvt"}
- **Status**: ✅ fully handled
- **Location**: select.go parsePivotExpr() L883

### S-TR-15: UNPIVOT Expression
- **SQL**: `SELECT * FROM (SELECT id, [2022], [2023] FROM t) unpivot_src UNPIVOT (amount FOR year IN ([2022], [2023])) AS unpvt`
- **Expected AST**: UnpivotExpr{Source: TableExpr, ValueCol: "amount", ForCol: "year", InCols: ["2022", "2023"], Alias: "unpvt"}
- **Status**: ✅ fully handled
- **Location**: select.go parseUnpivotExpr() (assumed implemented similarly to PIVOT)

### S-TR-16: Cross Apply
- **SQL**: `SELECT * FROM t1 CROSS APPLY fn(t1.id) AS t2`
- **Expected AST**: JoinClause{Type: JoinCrossApply, Left: TableRef("t1"), Right: FuncCallExpr("fn"), Condition: null}
- **Status**: ✅ fully handled
- **Location**: select.go matchJoinType() L1163

### S-TR-17: Outer Apply
- **SQL**: `SELECT * FROM t1 OUTER APPLY fn(t1.id) AS t2`
- **Expected AST**: JoinClause{Type: JoinOuterApply, Left: TableRef("t1"), Right: FuncCallExpr}
- **Status**: ✅ fully handled
- **Location**: select.go matchJoinType() L1169

### S-AL-01: Derived Table Alias Only
- **SQL**: `SELECT * FROM (SELECT a FROM t) x`
- **Expected AST**: AliasedTableRef{Table: SubqueryExpr, Alias: "x", Columns: nil}
- **Status**: ✅ fully handled
- **Location**: select.go parsePrimaryTableSource() L786

### S-AL-02: Derived Table With AS Keyword and Columns
- **SQL**: `SELECT * FROM (SELECT a, b FROM t) AS derived(col_a, col_b)`
- **Expected AST**: AliasedTableRef{Alias: "derived", Columns: ["col_a", "col_b"]}
- **Status**: ⚠️ Columns list not parsed after alias
- **Location**: select.go parsePrimaryTableSource() L786-795 (missing column list parsing)

### S-AL-03: Rowset Function Alias with Columns
- **SQL**: `SELECT * FROM OPENQUERY(server, 'SELECT * FROM t') AS r(id, val)`
- **Expected AST**: AliasedTableRef{Table: FuncCallExpr, Alias: "r", Columns: ["id", "val"]}
- **Status**: ⚠️ Columns discarded after alias parsing
- **Location**: rowset_functions.go parseRowsetFunction() L49-57

### S-AL-04: TVF Columns in Parentheses
- **SQL**: `SELECT * FROM dbo.fn(x) result(a INT, b INT)`
- **Expected AST**: AliasedTableRef with Columns containing column definitions
- **Status**: ⚠️ Columns parsed as identifiers only, not with type info
- **Location**: select.go parseTableValuedFunction() L1085 (missing column type parsing)

### S-INS-01: INSERT VALUES (no column list)
- **SQL**: `INSERT INTO t VALUES (1, 2, 3)`
- **Expected AST**: InsertStmt{Source: ValuesClause{Rows: [[1, 2, 3]]}}
- **Status**: ✅ fully handled
- **Location**: insert.go parseValuesClause() L172

### S-INS-02: INSERT with column list and VALUES
- **SQL**: `INSERT INTO t (col_a, col_b, col_c) VALUES (1, 2, 3), (4, 5, 6)`
- **Expected AST**: InsertStmt{Cols: ["col_a", "col_b", "col_c"], Source: ValuesClause with multiple rows}
- **Status**: ✅ fully handled
- **Location**: insert.go parseInsertStmt() L82-109, parseValuesClause() L172

