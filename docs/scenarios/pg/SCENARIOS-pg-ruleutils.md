# PG Ruleutils Complete Expression Coverage

> Goal: Complete ruleutils.go to handle ALL PG expression types that can appear in DDL/view contexts. For each type: add analyzed struct (query.go), analyzer support (analyze.go), deparse case (ruleutils.go).
> Verification: CREATE VIEW AS SELECT with the expression type, verify GetViewDefinition round-trips correctly. Oracle tests verify on real PG.
> Reference: PG ruleutils.c get_rule_expr switch, PG transformExpr in parse_expr.c

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Core Expression Types (DDL-critical)

### 1.1 Array and Composite Subscripting

- [ ] SubscriptingRef analyzed type + deparse: array subscript `arr[1]` in view SELECT
- [ ] SubscriptingRef in column DEFAULT: `DEFAULT tags[1]`
- [ ] SubscriptingRef in CHECK: `CHECK (tags[1] IS NOT NULL)`
- [ ] Multi-dimensional array subscript: `arr[1][2]`
- [ ] Array slice: `arr[1:3]`

### 1.2 Named Arguments and Coercions

- [ ] NamedArgExpr: function call with `func(arg_name => value)` in view
- [ ] ArrayCoerceExpr: `ARRAY[1,2]::text[]` in view SELECT
- [ ] ArrayCoerceExpr in DEFAULT: `DEFAULT ARRAY[1,2,3]::text[]`
- [ ] ConvertRowtypeExpr: row type conversion between compatible composite types
- [ ] CoerceToDomain: implicit domain coercion in view referencing domain-typed column

### 1.3 Case Internal and Row Compare

- [ ] CaseTestExpr: internal CASE placeholder — generated inside CASE WHEN with domain coercion
- [ ] RowCompareExpr: `(a, b) < (c, d)` in view WHERE clause
- [ ] RowCompareExpr with different operators: `(a, b) = (c, d)` in CHECK
- [ ] Param: `$1` parameter placeholder — roundtrip in expression context
- [ ] SetToDefault: DEFAULT keyword in expression context

### 1.4 Sequence and Type Helpers

- [ ] NextValueExpr: nextval for IDENTITY columns — deparse as nextval('seq')
- [ ] Simple type coercion chain: value → CoerceToDomain → column — verify no garbage in deparse
- [ ] FieldStore: composite field assignment (if analyzable)
- [ ] Multiple coercion chain: int → numeric → domain — verify correct deparse

---

## Phase 2: Function-Related Expression Types

### 2.1 Grouping and Window Enhancements

- [ ] GroupingFunc: `GROUPING(col)` in view with GROUP BY GROUPING SETS
- [ ] GROUP BY GROUPING SETS (a, b) deparse in view definition
- [ ] GROUP BY ROLLUP (a, b) deparse in view definition
- [ ] GROUP BY CUBE (a, b) deparse in view definition
- [ ] GroupingFunc combined with CASE WHEN in SELECT

### 2.2 XML Expressions

- [ ] XmlExpr XMLCONCAT: `xmlconcat(a, b)` in view
- [ ] XmlExpr XMLELEMENT: `xmlelement(name "row", col1, col2)` in view
- [ ] XmlExpr XMLFOREST: `xmlforest(col1, col2)` in view
- [ ] XmlExpr XMLPARSE: `xmlparse(document '<doc/>')` in view
- [ ] XmlExpr XMLPI: `xmlpi(name "php")` in view
- [ ] XmlExpr XMLROOT: `xmlroot(xmldata, version '1.0')` in view
- [ ] XmlExpr XMLSERIALIZE: `xmlserialize(content xmlcol AS text)` in view
- [ ] XmlExpr IS DOCUMENT: `xmlcol IS DOCUMENT` in view

### 2.3 Table Functions in FROM

- [ ] XMLTABLE in view FROM clause
- [ ] TABLESAMPLE BERNOULLI in view FROM clause
- [ ] TABLESAMPLE SYSTEM in view FROM clause

---

## Phase 3: JSON Expressions (PG16+)

### 3.1 JSON Constructor and Predicate

- [ ] JsonConstructorExpr JSON_OBJECT: `JSON_OBJECT('key': value)` in view
- [ ] JsonConstructorExpr JSON_ARRAY: `JSON_ARRAY(1, 2, 3)` in view
- [ ] JsonIsPredicate: `col IS JSON` in view WHERE
- [ ] JsonIsPredicate: `col IS NOT JSON OBJECT` in view WHERE

### 3.2 JSON Query and Value

- [ ] JsonExpr JSON_VALUE: `JSON_VALUE(col, '$.key')` in view
- [ ] JsonExpr JSON_QUERY: `JSON_QUERY(col, '$.arr')` in view
- [ ] JsonExpr JSON_EXISTS: `JSON_EXISTS(col, '$.key')` in view WHERE
- [ ] JsonValueExpr: JSON path expression in JSON_TABLE context

---

## Phase 4: Query Feature Completeness

### 4.1 SELECT Clause Completeness

- [ ] FOR UPDATE clause in view (if supported)
- [ ] FOR SHARE clause preserved in deparse
- [ ] DISTINCT ON with expressions
- [ ] VALUES clause in FROM: `SELECT * FROM (VALUES (1,'a'),(2,'b')) AS t(id, name)`

### 4.2 Plan-Internal Types (Skip or Stub)

- [ ] SubPlan: stub that outputs placeholder (plan-internal, never in user DDL)
- [ ] AlternativeSubPlan: stub that outputs placeholder
- [ ] MergeSupportFunc: stub that outputs placeholder
- [ ] InferenceElem: stub for ON CONFLICT (if needed)
- [ ] PartitionBoundSpec: stub for partition bound deparse

### 4.3 Integration Verification

- [ ] Complex view with 5+ expression types round-trips correctly
- [ ] View with CTE + window function + COALESCE + array subscript
- [ ] View with GROUP BY ROLLUP + GROUPING() + CASE WHEN
- [ ] Existing oracle tests still pass (no regression)
- [ ] Migration DDL with CURRENT_TIMESTAMP, COALESCE, CASE WHEN all correct on PG
