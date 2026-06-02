# Trino ANTLR4 Legacy Grammar Catalog (truth2)

> truth2 — legacy grammar baseline. Present-in-grammar for all entries.
> This is a HINT about scope, not an oracle.

## 1. Header

**License**: Apache License 2.0 (stated in both .g4 file headers).

**Grammar provenance**:
- `TrinoParser.g4` — parser grammar, tokenVocab = TrinoLexer; antlr-format directives present; ~1132 lines, 128 rules.
- `TrinoLexer.g4` — lexer grammar, caseInsensitive = true; ~397 lines; no lexer modes.
- The `SKIP_` token carries an inline comment: `// Missing SKIP in official g4` — indicating this grammar diverges from upstream Trino's published g4 by adding `SKIP_` and several punctuation tokens (`DOT_`, `COLON_`, `COMMA_`) not present in the official file. The `TRY_CAST_` keyword is also a local extension.
- The `COLON_` token literal is `'_:'` which appears to be a local bug/quirk vs the expected `':'`.

**Entry rules**:
- `parse` — top-level entry; accepts `statements* EOF`.
- `statements` — dispatches to one of: `singleStatement`, `standaloneExpression`, `standalonePathSpecification`, `standaloneType`, `standaloneRowPattern`, `standaloneFunctionSpecification`.
- `singleStatement` — `statement SEMICOLON_`.
- `standaloneExpression` — `expression SEMICOLON_`.
- `standalonePathSpecification` — `pathSpecification SEMICOLON_`.
- `standaloneType` — `type SEMICOLON_`.
- `standaloneRowPattern` — `rowPattern SEMICOLON_?`.
- `standaloneFunctionSpecification` — `functionSpecification SEMICOLON_`.

**Total parser rule count**: 128 rules (per header comment).

---

## 2. Statement Catalog

The `statement` rule is the central dispatch rule with 67 labeled alternatives listed below. Each entry gives: label, one-line syntax description, category.

| # | Label | Syntax description | Category |
|---|-------|--------------------|----------|
| 1 | `statementDefault` | Bare `rootQuery` (SELECT/WITH query as a statement) | DQL-wrapper |
| 2 | `use` | `USE schema` — set current schema | admin-session |
| 3 | `use` | `USE catalog.schema` — set current catalog and schema (same label, second alternative) | admin-session |
| 4 | `createCatalog` | `CREATE CATALOG [IF NOT EXISTS] name USING connector [COMMENT …] [AUTHORIZATION …] [WITH …]` | DDL |
| 5 | `dropCatalog` | `DROP CATALOG [IF EXISTS] name [CASCADE\|RESTRICT]` | DDL |
| 6 | `createSchema` | `CREATE SCHEMA [IF NOT EXISTS] qualifiedName [AUTHORIZATION …] [WITH …]` | DDL |
| 7 | `dropSchema` | `DROP SCHEMA [IF EXISTS] qualifiedName [CASCADE\|RESTRICT]` | DDL |
| 8 | `renameSchema` | `ALTER SCHEMA qualifiedName RENAME TO identifier` | DDL |
| 9 | `setSchemaAuthorization` | `ALTER SCHEMA qualifiedName SET AUTHORIZATION principal` | DDL |
| 10 | `createTableAsSelect` | `CREATE [OR REPLACE] TABLE [IF NOT EXISTS] qualifiedName [cols] [COMMENT …] [WITH …] AS (rootQuery) [WITH [NO] DATA]` | DDL |
| 11 | `createTable` | `CREATE [OR REPLACE] TABLE [IF NOT EXISTS] qualifiedName (tableElement, …) [COMMENT …] [WITH …]` | DDL |
| 12 | `dropTable` | `DROP TABLE [IF EXISTS] qualifiedName` | DDL |
| 13 | `insertInto` | `INSERT INTO qualifiedName [columnAliases] rootQuery` | DML |
| 14 | `delete` | `DELETE FROM qualifiedName [WHERE booleanExpression]` | DML |
| 15 | `truncateTable` | `TRUNCATE TABLE qualifiedName` | DML |
| 16 | `commentTable` | `COMMENT ON TABLE qualifiedName IS (string \| NULL)` | DDL |
| 17 | `commentView` | `COMMENT ON VIEW qualifiedName IS (string \| NULL)` | DDL |
| 18 | `commentColumn` | `COMMENT ON COLUMN qualifiedName IS (string \| NULL)` | DDL |
| 19 | `renameTable` | `ALTER TABLE [IF EXISTS] from RENAME TO to` | DDL |
| 20 | `addColumn` | `ALTER TABLE [IF EXISTS] tableName ADD COLUMN [IF NOT EXISTS] columnDefinition` | DDL |
| 21 | `renameColumn` | `ALTER TABLE [IF EXISTS] tableName RENAME COLUMN [IF EXISTS] from TO to` | DDL |
| 22 | `dropColumn` | `ALTER TABLE [IF EXISTS] tableName DROP COLUMN [IF EXISTS] column` | DDL |
| 23 | `setColumnType` | `ALTER TABLE [IF EXISTS] tableName ALTER COLUMN columnName SET DATA TYPE type` | DDL |
| 24 | `setTableAuthorization` | `ALTER TABLE tableName SET AUTHORIZATION principal` | DDL |
| 25 | `setTableProperties` | `ALTER TABLE tableName SET PROPERTIES propertyAssignments` | DDL |
| 26 | `tableExecute` | `ALTER TABLE tableName EXECUTE procedureName [(args)] [WHERE …]` | DML |
| 27 | `analyze` | `ANALYZE qualifiedName [WITH properties]` | utility |
| 28 | `createMaterializedView` | `CREATE [OR REPLACE] MATERIALIZED VIEW [IF NOT EXISTS] qualifiedName [GRACE PERIOD interval] [COMMENT …] [WITH …] AS rootQuery` | DDL |
| 29 | `createView` | `CREATE [OR REPLACE] VIEW qualifiedName [COMMENT …] [SECURITY DEFINER\|INVOKER] AS rootQuery` | DDL |
| 30 | `refreshMaterializedView` | `REFRESH MATERIALIZED VIEW qualifiedName` | DML |
| 31 | `dropMaterializedView` | `DROP MATERIALIZED VIEW [IF EXISTS] qualifiedName` | DDL |
| 32 | `renameMaterializedView` | `ALTER MATERIALIZED VIEW [IF EXISTS] from RENAME TO to` | DDL |
| 33 | `setMaterializedViewProperties` | `ALTER MATERIALIZED VIEW qualifiedName SET PROPERTIES propertyAssignments` | DDL |
| 34 | `dropView` | `DROP VIEW [IF EXISTS] qualifiedName` | DDL |
| 35 | `renameView` | `ALTER VIEW from RENAME TO to` | DDL |
| 36 | `setViewAuthorization` | `ALTER VIEW from SET AUTHORIZATION principal` | DDL |
| 37 | `call` | `CALL qualifiedName(callArgument, …)` — stored procedure call | utility |
| 38 | `createFunction` | `CREATE [OR REPLACE] functionSpecification` — inline SQL function | DDL |
| 39 | `dropFunction` | `DROP FUNCTION [IF EXISTS] functionDeclaration` | DDL |
| 40 | `createRole` | `CREATE ROLE name [WITH ADMIN grantor] [IN catalog]` | DCL |
| 41 | `dropRole` | `DROP ROLE name [IN catalog]` | DCL |
| 42 | `grantRoles` | `GRANT roles TO principal(s) [WITH ADMIN OPTION] [GRANTED BY grantor] [IN catalog]` | DCL |
| 43 | `revokeRoles` | `REVOKE [ADMIN OPTION FOR] roles FROM principal(s) [GRANTED BY grantor] [IN catalog]` | DCL |
| 44 | `setRole` | `SET ROLE (ALL \| NONE \| identifier) [IN catalog]` | DCL |
| 45 | `grant` | `GRANT privilege(s)\|ALL PRIVILEGES ON [SCHEMA\|TABLE] qualifiedName TO principal [WITH GRANT OPTION]` | DCL |
| 46 | `deny` | `DENY privilege(s)\|ALL PRIVILEGES ON [SCHEMA\|TABLE] qualifiedName TO principal` | DCL |
| 47 | `revoke` | `REVOKE [GRANT OPTION FOR] privilege(s)\|ALL PRIVILEGES ON [SCHEMA\|TABLE] qualifiedName FROM principal` | DCL |
| 48 | `showGrants` | `SHOW GRANTS [ON TABLE? qualifiedName]` | utility |
| 49 | `explain` | `EXPLAIN [(options)] statement` — shows query plan | utility |
| 50 | `explainAnalyze` | `EXPLAIN ANALYZE [VERBOSE] statement` — executes and shows actual plan | utility |
| 51 | `showCreateTable` | `SHOW CREATE TABLE qualifiedName` | utility |
| 52 | `showCreateSchema` | `SHOW CREATE SCHEMA qualifiedName` | utility |
| 53 | `showCreateView` | `SHOW CREATE VIEW qualifiedName` | utility |
| 54 | `showCreateMaterializedView` | `SHOW CREATE MATERIALIZED VIEW qualifiedName` | utility |
| 55 | `showTables` | `SHOW TABLES [(FROM\|IN) qualifiedName] [LIKE pattern [ESCAPE escape]]` | utility |
| 56 | `showSchemas` | `SHOW SCHEMAS [(FROM\|IN) identifier] [LIKE pattern [ESCAPE escape]]` | utility |
| 57 | `showCatalogs` | `SHOW CATALOGS [LIKE pattern [ESCAPE escape]]` | utility |
| 58 | `showColumns` | `SHOW COLUMNS (FROM\|IN) qualifiedName? [LIKE pattern [ESCAPE escape]]` (also triggered by `DESCRIBE` and `DESC`) | utility |
| 59 | `showStats` | `SHOW STATS FOR qualifiedName` | utility |
| 60 | `showStatsForQuery` | `SHOW STATS FOR (rootQuery)` | utility |
| 61 | `showRoles` | `SHOW [CURRENT] ROLES [(FROM\|IN) identifier]` | utility |
| 62 | `showRoleGrants` | `SHOW ROLE GRANTS [(FROM\|IN) identifier]` | utility |
| 63 | `showFunctions` | `SHOW FUNCTIONS [(FROM\|IN) qualifiedName] [LIKE pattern [ESCAPE escape]]` | utility |
| 64 | `showSession` | `SHOW SESSION [LIKE pattern [ESCAPE escape]]` | utility |
| 65 | `setSessionAuthorization` | `SET SESSION AUTHORIZATION authorizationUser` | admin-session |
| 66 | `resetSessionAuthorization` | `RESET SESSION AUTHORIZATION` | admin-session |
| 67 | `setSession` | `SET SESSION qualifiedName = expression` — set session property | admin-session |
| 68 | `resetSession` | `RESET SESSION qualifiedName` | admin-session |
| 69 | `startTransaction` | `START TRANSACTION [transactionMode, …]` | transaction |
| 70 | `commit` | `COMMIT [WORK]` | transaction |
| 71 | `rollback` | `ROLLBACK [WORK]` | transaction |
| 72 | `prepare` | `PREPARE identifier FROM statement` | prepared |
| 73 | `deallocate` | `DEALLOCATE PREPARE identifier` | prepared |
| 74 | `execute` | `EXECUTE identifier [USING expression, …]` | prepared |
| 75 | `executeImmediate` | `EXECUTE IMMEDIATE string [USING expression, …]` | prepared |
| 76 | `describeInput` | `DESCRIBE INPUT identifier` | prepared |
| 77 | `describeOutput` | `DESCRIBE OUTPUT identifier` | prepared |
| 78 | `setPath` | `SET PATH pathSpecification` | admin-session |
| 79 | `setTimeZone` | `SET TIME ZONE (LOCAL \| expression)` | admin-session |
| 80 | `update` | `UPDATE qualifiedName SET col=expr, … [WHERE …]` | DML |
| 81 | `merge` | `MERGE INTO qualifiedName [AS alias] USING relation ON expr mergeCase+` | DML |

**Total statement-rule labeled alternatives: 81**

Note: `showColumns` appears three times in the grammar (`SHOW COLUMNS`, `DESCRIBE`, `DESC`) but all share the same label — counted as one unique label with 3 grammar alternatives.

---

## 3. Query / SELECT Rules

| Rule | Description |
|------|-------------|
| `rootQuery` | Top-level query wrapper: optional `withFunction` (inline SQL functions) then `query`. |
| `withFunction` | `WITH functionSpecification, …` — inline function definitions preceding a query (Trino extension). |
| `query` | Optional `with` CTE clause then `queryNoWith`. |
| `with` | `WITH [RECURSIVE] namedQuery, …` — CTE list. |
| `namedQuery` | `name [(columns)] AS (query)` — a single CTE definition. |
| `queryNoWith` | Core SELECT body: `queryTerm [ORDER BY …] [OFFSET …] [LIMIT … \| FETCH FIRST … ONLY\|WITH TIES]`. |
| `limitRowCount` | `ALL` or a `rowCount`; used for LIMIT clause. |
| `rowCount` | `INTEGER_VALUE` or `?` (positional parameter); used in LIMIT/OFFSET/FETCH. |
| `queryTerm` | Labeled alternatives: `queryTermDefault` (single primary), `setOperation` (INTERSECT / UNION / EXCEPT with optional setQuantifier). |
| `queryPrimary` | Labeled alternatives: `queryPrimaryDefault` (querySpecification), `table` (TABLE qualifiedName), `inlineTable` (VALUES expr, …), `subquery` (parenthesized queryNoWith). |
| `querySpecification` | Full SELECT: `SELECT [setQuantifier] selectItem, … [FROM relation, …] [WHERE …] [GROUP BY …] [HAVING …] [WINDOW …]`. |
| `setQuantifier` | `DISTINCT` or `ALL`. |
| `selectItem` | Labeled alternatives: `selectSingle` (expression [AS alias]), `selectAll` (qualified.* [AS (cols)] or bare *). |
| `as_column_alias` | `[AS] column_alias` — optional alias marker. |
| `column_alias` | Single `identifier` used as a column alias. |
| `sortItem` | `expression [ASC\|DESC] [NULLS FIRST\|LAST]`. |
| `groupBy` | `[setQuantifier] groupingElement, …`. |
| `groupingElement` | Labeled alternatives: `singleGroupingSet` (groupingSet), `rollup` (ROLLUP(…)), `cube` (CUBE(…)), `multipleGroupingSets` (GROUPING SETS(…)). |
| `groupingSet` | `(expr, …)` or bare `expr`. |
| `windowDefinition` | `name AS (windowSpecification)` — named window definition. |
| `windowSpecification` | Optional existing-window-name, PARTITION BY, ORDER BY, and windowFrame. |
| `over` | `OVER (windowName \| (windowSpecification))` — window clause on a function call. |
| `windowFrame` | Optional MEASURES, frameExtent, AFTER MATCH skipTo, INITIAL/SEEK, PATTERN, SUBSET, DEFINE — MATCH_RECOGNIZE frame extension. |
| `frameExtent` | RANGE/ROWS/GROUPS with single or BETWEEN frameBound pair. |
| `frameBound` | Labeled alternatives: `unboundedFrame` (UNBOUNDED PRECEDING/FOLLOWING), `currentRowBound` (CURRENT ROW), `boundedFrame` (expr PRECEDING/FOLLOWING). |
| `relation` | Labeled alternatives: `joinRelation` (CROSS JOIN / typed JOIN / NATURAL JOIN), `relationDefault` (sampledRelation). |
| `joinType` | INNER (optional), or LEFT/RIGHT/FULL [OUTER]. |
| `joinCriteria` | `ON booleanExpression` or `USING (identifier, …)`. |
| `sampledRelation` | `patternRecognition [TABLESAMPLE sampleType (percentage)]`. |
| `sampleType` | `BERNOULLI` or `SYSTEM`. |
| `aliasedRelation` | `relationPrimary [AS? identifier [columnAliases]]`. |
| `relationPrimary` | Labeled alternatives: `tableName` (qualifiedName [queryPeriod]), `subqueryRelation` ((query)), `unnest` (UNNEST(…) [WITH ORDINALITY]), `lateral` (LATERAL (query)), `tableFunctionInvocation` (TABLE(tableFunctionCall)), `parenthesizedRelation` ((relation)). |
| `queryPeriod` | `FOR rangeType AS OF end` — time-travel clause (FOR TIMESTAMP\|VERSION AS OF expr). |
| `rangeType` | `TIMESTAMP` or `VERSION`. |
| `tableFunctionCall` | `qualifiedName(tableFunctionArgument, … [COPARTITION copartitionTables, …])` — polymorphic table function invocation. |
| `tableFunctionArgument` | Optional named (`identifier =>`), then `tableArgument`, `descriptorArgument`, or `expression`. |
| `tableArgument` | `TABLE(qualifiedName\|query) [PARTITION BY …] [PRUNE/KEEP WHEN EMPTY] [ORDER BY …]`. |
| `tableArgumentRelation` | Labeled alternatives: `tableArgumentTable` (TABLE(qualifiedName)), `tableArgumentQuery` (TABLE(query)). |
| `descriptorArgument` | `DESCRIPTOR(field, …)` or `CAST(NULL AS DESCRIPTOR)`. |
| `descriptorField` | `identifier [type]`. |
| `copartitionTables` | `(qualifiedName, qualifiedName, …)` — list of co-partitioned table names. |
| `patternRecognition` | `aliasedRelation [MATCH_RECOGNIZE (…)]` — MATCH_RECOGNIZE clause wrapper. |
| `measureDefinition` | `expression AS identifier` — MEASURES clause entry. |
| `rowsPerMatch` | `ONE ROW PER MATCH` or `ALL ROWS PER MATCH [emptyMatchHandling]`. |
| `emptyMatchHandling` | `SHOW EMPTY MATCHES`, `OMIT EMPTY MATCHES`, or `WITH UNMATCHED ROWS`. |
| `skipTo` | `SKIP TO (NEXT ROW \| [FIRST\|LAST] identifier \| PAST LAST ROW)` — AFTER MATCH skip strategy. |
| `subsetDefinition` | `name = (identifier, …)` — SUBSET union definition. |
| `variableDefinition` | `identifier AS expression` — DEFINE clause entry. |
| `rowPattern` | Labeled alternatives: `quantifiedPrimary` (patternPrimary [patternQuantifier]), `patternConcatenation` (rowPattern rowPattern), `patternAlternation` (rowPattern \| rowPattern). |
| `patternPrimary` | Labeled alternatives: `patternVariable` (identifier), `emptyPattern` (()), `patternPermutation` (PERMUTE(…)), `groupedPattern` ((rowPattern)), `partitionStartAnchor` (^), `partitionEndAnchor` ($), `excludedPattern` ({- rowPattern -}). |
| `patternQuantifier` | Labeled alternatives: `zeroOrMoreQuantifier` (* [?]), `oneOrMoreQuantifier` (+ [?]), `zeroOrOneQuantifier` (? [?]), `rangeQuantifier` ({n} [?] or {m,n} [?]). |

---

## 4. Expression Rules

| Rule | Description |
|------|-------------|
| `expression` | Top-level expression; simply delegates to `booleanExpression`. |
| `booleanExpression` | Labeled alternatives: `predicated` (valueExpression [predicate_]), `logicalNot` (NOT booleanExpression), `and` (booleanExpression AND booleanExpression), `or` (booleanExpression OR booleanExpression). |
| `predicate_` | Predicate suffixes (workaround for ANTLR4 issue #780). Labeled alternatives: `comparison` (comparisonOperator valueExpression), `quantifiedComparison` (op quantifier (query)), `between` ([NOT] BETWEEN lower AND upper), `inList` ([NOT] IN (expr, …)), `inSubquery` ([NOT] IN (query)), `like` ([NOT] LIKE pattern [ESCAPE escape]), `nullPredicate` (IS [NOT] NULL), `distinctFrom` (IS [NOT] DISTINCT FROM valueExpression). |
| `valueExpression` | Labeled alternatives: `valueExpressionDefault` (primaryExpression), `atTimeZone` (valueExpression AT timeZoneSpecifier), `arithmeticUnary` (- \| + valueExpression), `arithmeticBinary` (left * / % right, or left + - right), `concatenation` (left \|\| right). |
| `primaryExpression` | Large labeled rule for primary expressions. Labeled alternatives: |
| | `nullLiteral` — `NULL` |
| | `intervalLiteral` — `interval` |
| | `typeConstructor` — `identifier string_` or `DOUBLE PRECISION string_` |
| | `numericLiteral` — `number` |
| | `booleanLiteral` — `booleanValue` (TRUE/FALSE) |
| | `stringLiteral` — `string_` |
| | `binaryLiteral` — `BINARY_LITERAL` |
| | `parameter` — `?` |
| | `position` — `POSITION(valueExpression IN valueExpression)` |
| | `rowConstructor` — `(expr, expr, …)` or `ROW(expr, …)` |
| | `listagg` — `LISTAGG([DISTINCT] expr [, sep] [ON OVERFLOW …]) WITHIN GROUP (ORDER BY …) [FILTER …]` |
| | `functionCall` — `[processingMode] qualifiedName([DISTINCT\|ALL] expr, … [ORDER BY …]) [FILTER] [nullTreatment] [OVER]` (two alternatives: with `*` argument and without) |
| | `measure` — `identifier OVER` — MATCH_RECOGNIZE measure reference |
| | `lambda` — `identifier -> expr` or `(identifiers) -> expr` |
| | `subqueryExpression` — `(query)` |
| | `exists` — `EXISTS(query)` |
| | `simpleCase` — `CASE expr WHEN … [ELSE …] END` |
| | `searchedCase` — `CASE WHEN … [ELSE …] END` |
| | `cast` — `CAST(expr AS type)` or `TRY_CAST(expr AS type)` |
| | `arrayConstructor` — `ARRAY[expr, …]` |
| | `subscript` — `primaryExpression[valueExpression]` |
| | `columnReference` — bare `identifier` |
| | `dereference` — `primaryExpression.identifier` |
| | `specialDateTimeFunction` — CURRENT_DATE, CURRENT_TIME[(p)], CURRENT_TIMESTAMP[(p)], LOCALTIME[(p)], LOCALTIMESTAMP[(p)] |
| | `currentUser` — `CURRENT_USER` |
| | `currentCatalog` — `CURRENT_CATALOG` |
| | `currentSchema` — `CURRENT_SCHEMA` |
| | `currentPath` — `CURRENT_PATH` |
| | `trim` — `TRIM([LEADING\|TRAILING\|BOTH] [trimChar FROM] trimSource)` or `TRIM(source, char)` |
| | `substring` — `SUBSTRING(expr FROM expr [FOR expr])` |
| | `normalize` — `NORMALIZE(expr [, normalForm])` |
| | `extract` — `EXTRACT(field FROM expr)` |
| | `parenthesizedExpression` — `(expr)` |
| | `groupingOperation` — `GROUPING(qualifiedName, …)` |
| | `jsonExists` — `JSON_EXISTS(jsonPathInvocation [errorBehavior ON ERROR])` |
| | `jsonValue` — `JSON_VALUE(jsonPathInvocation [RETURNING type] [ON EMPTY] [ON ERROR])` |
| | `jsonQuery` — `JSON_QUERY(jsonPathInvocation [RETURNING type] [WRAPPER] [QUOTES] [ON EMPTY] [ON ERROR])` |
| | `jsonObject` — `JSON_OBJECT([members] [NULL\|ABSENT ON NULL] [UNIQUE KEYS] [RETURNING type])` |
| | `jsonArray` — `JSON_ARRAY([values] [NULL\|ABSENT ON NULL] [RETURNING type])` |
| `jsonPathInvocation` | `jsonValueExpression, path [PASSING jsonArgument, …]` — JSON path with optional PASSING arguments. |
| `jsonValueExpression` | `expression [FORMAT jsonRepresentation]` — expression optionally with JSON format. |
| `jsonRepresentation` | `JSON [ENCODING UTF8\|UTF16\|UTF32]` — JSON format with optional encoding; TODO comment in grammar for implementation-defined options. |
| `jsonArgument` | `jsonValueExpression AS identifier` — named JSON path argument. |
| `jsonExistsErrorBehavior` | `TRUE \| FALSE \| UNKNOWN \| ERROR` — ON ERROR behavior for JSON_EXISTS. |
| `jsonValueBehavior` | `ERROR \| NULL \| DEFAULT expression` — ON EMPTY/ERROR behavior for JSON_VALUE. |
| `jsonQueryWrapperBehavior` | `WITHOUT [ARRAY] \| WITH [CONDITIONAL\|UNCONDITIONAL] [ARRAY]` — WRAPPER clause for JSON_QUERY. |
| `jsonQueryBehavior` | `ERROR \| NULL \| EMPTY (ARRAY\|OBJECT)` — ON EMPTY/ERROR behavior for JSON_QUERY. |
| `jsonObjectMember` | `[KEY] expr VALUE jsonValueExpression` or `expr : jsonValueExpression` — key-value pair in JSON_OBJECT. |
| `processingMode` | `RUNNING` or `FINAL` — MATCH_RECOGNIZE processing mode for function calls. |
| `nullTreatment` | `IGNORE NULLS` or `RESPECT NULLS` — null handling for window functions. |
| `comparisonOperator` | One of `=`, `<>`, `!=`, `<`, `<=`, `>`, `>=`. |
| `comparisonQuantifier` | `ALL`, `SOME`, or `ANY`. |
| `booleanValue` | `TRUE` or `FALSE`. |
| `whenClause` | `WHEN condition THEN result` — used in CASE expressions. |
| `filter` | `FILTER (WHERE booleanExpression)` — aggregate filter clause. |
| `timeZoneSpecifier` | Labeled: `timeZoneInterval` (TIME ZONE interval), `timeZoneString` (TIME ZONE string_). |
| `trimsSpecification` | `LEADING`, `TRAILING`, or `BOTH` — TRIM direction. |
| `listAggOverflowBehavior` | `ERROR` or `TRUNCATE [string] listaggCountIndication` — LISTAGG overflow handling. |
| `listaggCountIndication` | `(WITH \| WITHOUT) COUNT` — whether to include count in LISTAGG truncation. |
| `normalForm` | `NFD`, `NFC`, `NFKD`, or `NFKC` — Unicode normal form for NORMALIZE. |
| `updateAssignment` | `identifier = expression` — single SET assignment in UPDATE. |
| `mergeCase` | Labeled alternatives: `mergeUpdate` (WHEN MATCHED … THEN UPDATE SET …), `mergeDelete` (WHEN MATCHED … THEN DELETE), `mergeInsert` (WHEN NOT MATCHED … THEN INSERT … VALUES …). |
| `explainOption` | Labeled alternatives: `explainFormat` (FORMAT TEXT\|GRAPHVIZ\|JSON), `explainType` (TYPE LOGICAL\|DISTRIBUTED\|VALIDATE\|IO). |
| `transactionMode` | Labeled alternatives: `isolationLevel` (ISOLATION LEVEL levelOfIsolation), `transactionAccessMode` (READ ONLY\|WRITE). |
| `levelOfIsolation` | Labeled alternatives: `readUncommitted`, `readCommitted`, `repeatableRead`, `serializable`. |

---

## 5. Type / Literal / Misc Rules

### Types

| Rule | Description |
|------|-------------|
| `type` | Labeled alternatives: `rowType` (ROW(rowField, …)), `intervalType` (INTERVAL field [TO field]), `dateTimeType` (TIMESTAMP/TIME [(p)] [WITHOUT\|WITH TIME ZONE]), `doublePrecisionType` (DOUBLE PRECISION), `legacyArrayType` (ARRAY\<type\>), `legacyMapType` (MAP\<keyType, valueType\>), `arrayType` (type ARRAY [size]), `genericType` (identifier [(typeParameter, …)]). |
| `rowField` | Bare `type` or named `identifier type`. |
| `typeParameter` | `INTEGER_VALUE` or `type`. |

### Identifiers

| Rule | Description |
|------|-------------|
| `qualifiedName` | `identifier (. identifier)*` — dotted-path name. |
| `identifier` | Labeled alternatives: `unquotedIdentifier` (IDENTIFIER_ or nonReserved keyword), `quotedIdentifier` (double-quoted), `backQuotedIdentifier` (backtick-quoted), `digitIdentifier` (starts with digit). |
| `nonReserved` | Flat list of ~150 keyword tokens usable as identifiers; all Trino-specific reserved-word exemptions are here (e.g., ARRAY, MAP, SHOW, GRANT, MERGE, MATCH_RECOGNIZE, JSON, …). |

### Strings / Numbers / Literals

| Rule | Description |
|------|-------------|
| `string_` | Labeled alternatives: `basicStringLiteral` (single-quoted STRING_), `unicodeStringLiteral` (U&'…' [UESCAPE '…']). |
| `number` | Labeled alternatives: `decimalLiteral` ([−] DECIMAL_VALUE_), `doubleLiteral` ([−] DOUBLE_VALUE_), `integerLiteral` ([−] INTEGER_VALUE_). |
| `interval` | `INTERVAL [sign] string_ from [TO to]` — interval literal. |
| `intervalField` | `YEAR \| MONTH \| DAY \| HOUR \| MINUTE \| SECOND`. |
| `booleanValue` | `TRUE` or `FALSE`. |

### DCL / Privilege / Role / Principal

| Rule | Description |
|------|-------------|
| `privilege` | One of `CREATE`, `SELECT`, `DELETE`, `INSERT`, `UPDATE`. |
| `roles` | `identifier (, identifier)*` — role name list. |
| `grantor` | Labeled: `specifiedPrincipal` (principal), `currentUserGrantor` (CURRENT_USER), `currentRoleGrantor` (CURRENT_ROLE). |
| `principal` | Labeled: `unspecifiedPrincipal` (identifier), `userPrincipal` (USER identifier), `rolePrincipal` (ROLE identifier). |
| `authorizationUser` | Labeled: `identifierUser` (identifier), `stringUser` (string_). |

### Table / Column / Property

| Rule | Description |
|------|-------------|
| `tableElement` | `columnDefinition` or `likeClause`. |
| `columnDefinition` | `identifier type [NOT NULL] [COMMENT string] [WITH properties]`. |
| `likeClause` | `LIKE qualifiedName [(INCLUDING\|EXCLUDING) PROPERTIES]`. |
| `properties` | `(propertyAssignments)`. |
| `propertyAssignments` | `property (, property)*`. |
| `property` | `identifier = propertyValue`. |
| `propertyValue` | Labeled: `defaultPropertyValue` (DEFAULT), `nonDefaultPropertyValue` (expression). |

### Function / Routine Rules

| Rule | Description |
|------|-------------|
| `functionSpecification` | `FUNCTION functionDeclaration returnsClause routineCharacteristic* controlStatement`. |
| `functionDeclaration` | `qualifiedName([parameterDeclaration, …])` — function signature. |
| `parameterDeclaration` | Optional `identifier` then `type` — parameter declaration. |
| `returnsClause` | `RETURNS type`. |
| `routineCharacteristic` | Labeled alternatives: `languageCharacteristic` (LANGUAGE identifier), `deterministicCharacteristic` ([NOT] DETERMINISTIC), `returnsNullOnNullInputCharacteristic` (RETURNS NULL ON NULL INPUT), `calledOnNullInputCharacteristic` (CALLED ON NULL INPUT), `securityCharacteristic` (SECURITY DEFINER\|INVOKER), `commentCharacteristic` (COMMENT string_). |
| `controlStatement` | Labeled alternatives: `returnStatement` (RETURN), `assignmentStatement` (SET id = expr), `simpleCaseStatement` (CASE expr … END CASE), `searchedCaseStatement` (CASE … END CASE), `ifStatement` (IF … THEN … [ELSEIF] [ELSE] END IF), `iterateStatement` (ITERATE label), `leaveStatement` (LEAVE label), `compoundStatement` (BEGIN … END), `loopStatement` ([label:] LOOP … END LOOP), `whileStatement` ([label:] WHILE … DO … END WHILE), `repeatStatement` ([label:] REPEAT … UNTIL … END REPEAT). |
| `caseStatementWhenClause` | `WHEN expression THEN sqlStatementList` — for case control statements. |
| `elseIfClause` | `ELSEIF expression THEN sqlStatementList`. |
| `elseClause` | `ELSE sqlStatementList`. |
| `variableDeclaration` | `DECLARE identifier(s) type [DEFAULT valueExpression]`. |
| `sqlStatementList` | `(controlStatement SEMICOLON_)+`. |

### Call / Path Arguments

| Rule | Description |
|------|-------------|
| `callArgument` | Labeled: `positionalArgument` (expression), `namedArgument` (identifier => expression). |
| `pathElement` | Labeled: `qualifiedArgument` (identifier.identifier), `unqualifiedArgument` (identifier). |
| `pathSpecification` | `pathElement (, pathElement)*` — SET PATH argument. |

---

## 6. Lexer Summary

### Token categories

**Keywords** (reserved and non-reserved combined): approximately 220 keyword tokens, all case-insensitive (`caseInsensitive = true` in lexer options). Reserved keywords (those NOT listed in `nonReserved`) include: ALTER, AND, AS, BETWEEN, BY, CASE, CAST, CONSTRAINT, CROSS, CUBE, DEALLOCATE, DELETE, DESCRIBE, DISTINCT, DROP, ELSE, END, ESCAPE, EXCEPT, EXECUTE, EXISTS, EXTRACT, FALSE, FETCH, FOR, FROM, FULL, GROUP, GROUPING, HAVING, IN, INNER, INSERT, INTERSECT, INTO, IS, JOIN, LATERAL, LEFT, LIKE, MERGE, NATURAL, NORMALIZE, NOT, NULL, NULLIF, ON, OR, ORDER, OUTER, PREPARE, RECURSIVE, RIGHT, ROLLUP, ROW, SELECT, SKIP, SUBSTRING, TABLE, THEN, TRIM, TRUE, UNION, UNNEST, UESCAPE, USING, VALUES, WHEN, WHERE, WITH, and the compound keywords (CURRENT_DATE, CURRENT_TIME, CURRENT_TIMESTAMP, LOCALTIME, LOCALTIMESTAMP, CURRENT_USER, CURRENT_CATALOG, CURRENT_SCHEMA, CURRENT_PATH, CURRENT_ROLE, JSON_ARRAY, JSON_EXISTS, JSON_OBJECT, JSON_QUERY, JSON_TABLE, JSON_VALUE, MATCH_RECOGNIZE, TRY_CAST, LISTAGG, NORMALIZE).

**Non-reserved keywords** (~150 tokens listed in the `nonReserved` parser rule): can be used as unquoted identifiers. Examples: ABSENT, ADD, ADMIN, AFTER, ALL, ANALYZE, ANY, ARRAY, ASC, AT, AUTHORIZATION, BEGIN, BERNOULLI, BOTH, CALL, CALLED, CASCADE, CATALOG, CATALOGS, COLUMN, COLUMNS, COMMENT, COMMIT, COMMITTED, CONDITIONAL, COPARTITION, COUNT, CURRENT, DATA, DATE, DAY, DECLARE, DEFAULT, DEFINE, DEFINER, DENY, DESC, DESCRIPTOR, DETERMINISTIC, DISTRIBUTED, DO, DOUBLE, ELSEIF, EMPTY, ENCODING, ERROR, EXCLUDING, EXPLAIN, FETCH, FILTER, FINAL, FIRST, FOLLOWING, FORMAT, FUNCTION, FUNCTIONS, GRACE, GRANT, GRANTED, GRANTS, GRAPHVIZ, GROUPS, HOUR, IF, IGNORE, IMMEDIATE, INCLUDING, INITIAL, INPUT, INTERVAL, INVOKER, IO, ITERATE, ISOLATION, JSON, KEEP, KEY, KEYS, LANGUAGE, LAST, LATERAL, LEADING, LEAVE, LEVEL, LIMIT, LOCAL, LOGICAL, LOOP, MAP, MATCH, MATCHED, MATCHES, MATCH_RECOGNIZE, MATERIALIZED, MEASURES, MERGE, MINUTE, MONTH, NESTED, NEXT, NFC, NFD, NFKC, NFKD, NO, NONE, NULLIF, NULLS, OBJECT, OF, OFFSET, OMIT, ONE, ONLY, OPTION, ORDINALITY, OUTPUT, OVER, OVERFLOW, PARTITION, PARTITIONS, PASSING, PAST, PATH, PATTERN, PER, PERIOD, PERMUTE, PLAN, POSITION, PRECEDING, PRECISION, PRIVILEGES, PROPERTIES, PRUNE, QUOTES, RANGE, READ, REFRESH, RENAME, REPEAT, REPEATABLE, REPLACE, RESET, RESPECT, RESTRICT, RETURN, RETURNING, RETURNS, REVOKE, ROLE, ROLES, ROLLBACK, ROW, ROWS, RUNNING, SCALAR, SCHEMA, SCHEMAS, SECOND, SECURITY, SEEK, SERIALIZABLE, SESSION, SET, SETS, SHOW, SOME, START, STATS, SUBSET, SUBSTRING, SYSTEM, TABLES, TABLESAMPLE, TEXT, TEXT_STRING (='STRING'), TIES, TIME, TIMESTAMP, TO, TRAILING, TRANSACTION, TRUNCATE, TRY_CAST, TYPE, UNBOUNDED, UNCOMMITTED, UNCONDITIONAL, UNIQUE, UNKNOWN, UNMATCHED, UNTIL, UPDATE, USE, USER, UTF16, UTF32, UTF8, VALIDATE, VALUE, VERBOSE, VERSION, VIEW, WHILE, WINDOW, WITHIN, WITHOUT, WORK, WRAPPER, WRITE, YEAR, ZONE.

**Comparison / arithmetic operators**:
`=`, `<>`, `!=`, `<`, `<=`, `>`, `>=`, `+`, `-`, `*`, `/`, `%`, `||` (CONCAT_), `?`.

**Punctuation / delimiters**:
- `SEMICOLON_` = `;`
- `DOT_` = `.` (added locally, not in official g4)
- `COLON_` = `'_:'` (local addition; literal appears to be a bug — likely intended as `:`)
- `COMMA_` = `,` (added locally)
- `LPAREN_` = `(`, `RPAREN_` = `)`
- `LSQUARE_` = `[`, `RSQUARE_` = `]`
- `LCURLY_` = `{`, `RCURLY_` = `}`
- `LCURLYHYPHEN_` = `{-`, `RCURLYHYPHEN_` = `-}` (for MATCH_RECOGNIZE excluded patterns)
- `LARROW_` = `<-`, `RARROW_` = `->`, `RDOUBLEARROW_` = `=>`
- `VBAR_` = `|`, `DOLLAR_` = `$`, `CARET_` = `^`

**Lexer modes**: None. Single default mode throughout.

### Literals / identifiers

| Token | Pattern | Notes |
|-------|---------|-------|
| `STRING_` | `'...'` with `''` escape | Single-quoted string literal |
| `UNICODE_STRING_` | `U&'...' [UESCAPE '...']` | Unicode string literal |
| `BINARY_LITERAL_` | `X'...'` | Hex binary literal; content validated at AST construction time |
| `INTEGER_VALUE_` | `[0-9]+` | Integer literal |
| `DECIMAL_VALUE_` | `digits.digits` or `.digits` | Decimal literal |
| `DOUBLE_VALUE_` | `digits[.digits]E[+-]digits` | Double/float literal with exponent |
| `IDENTIFIER_` | `[A-Za-z_][A-Za-z0-9_]*` | Regular unquoted identifier |
| `DIGIT_IDENTIFIER_` | `[0-9][A-Za-z0-9_]+` | Identifier starting with a digit |
| `QUOTED_IDENTIFIER_` | `"..."` with `""` escape | Double-quoted identifier |
| `BACKQUOTED_IDENTIFIER_` | `` `...` `` with ``` `` ``` escape | Backtick-quoted identifier |

### Comments / whitespace

| Token | Pattern | Channel |
|-------|---------|---------|
| `SIMPLE_COMMENT_` | `--` to end of line | HIDDEN |
| `BRACKETED_COMMENT_` | `/* ... */` (non-greedy) | HIDDEN |
| `WS_` | `[ \r\n\t]+` | HIDDEN |
| `UNRECOGNIZED_` | `.` (catch-all) | default (for delimiter splitting) |

### Naming convention

All tokens use a trailing underscore suffix (e.g., `SELECT_`, `LPAREN_`, `INTEGER_VALUE_`, `IDENTIFIER_`). This is a deliberate convention in this grammar variant to avoid conflicts with ANTLR4 reserved names and Go/target-language keywords. Fragment rules (`EXPONENT_`, `DIGIT_`, `LETTER_`) also carry the suffix.

### Local divergences from official Trino g4

- `SKIP_` keyword added with comment `// Missing SKIP in official g4`.
- `DOT_`, `COLON_`, `COMMA_` punctuation tokens added locally (official g4 inlines these as literals in parser rules).
- `TRY_CAST_` is a compound keyword token (official g4 treats it as `TRY` `(` `CAST` `...` `)` at the parser level).
- `COLON_` literal is `'_:'` — likely a copy/paste artifact; semantically should be `':'`.
- `TEXT_STRING_` maps to the keyword `'STRING'` (the Trino `STRING` type alias), which shadows the token name to avoid collision with the `STRING_` literal token.
