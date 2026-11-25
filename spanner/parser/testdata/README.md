# GoogleSQL parser test data

`spanner/parser` is Spanner's grammar; BigQuery's is `bigquery/parser`, a fork
judged by BigQuery itself (`../../../bigquery/parser/testdata/README.md`).

## `compat.tsv`: pinned Spanner emulator verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/GoogleSQL.g4`, then its near-misses, plus 734 hand-written
`suite:` rows, with omni's verdict (`spanner/parser.Parse`) and the pinned
Cloud Spanner emulator's (`../../testdata/pins.tsv`, `parser` role, judge
revision 5: `gcr.io/cloud-spanner-emulator/emulator:1.5.54`, through
`spanner/spanner-oracle`). `GoogleSQL.g4` is converted by `../cmd/tm2g4` from
`googlesql.tm`, the Textmapper grammar of the GoogleSQL release the emulator
builds against, beside it as published. `TestCompatibility`
(`../compat_test.go`) reads the file through `enginetest/compat` and starts no
server in offline replay.

| Rows | Total | Agree | Only Spanner accepts | Only omni accepts | Agreement |
| --- | ---: | ---: | ---: | ---: | ---: |
| One per alternative | 2,864 | 2,864 | 0 | 0 | 100.00% |
| Near-misses | 6,573 | 6,573 | 0 | 0 | 100.00% |
| Suites (`suite:*`) | 734 | 731 | 0 | 3 | 99.59% |
| **All** | **10,171** | **10,168** | **0** | **3** | **99.97%** |

The three rows only omni accepts are documented Spanner forms the emulator's
DDL parser lags: `ALTER SEARCH INDEX idx {ADD|DROP} STORED COLUMN c` (DDL-049)
and `CREATE PLACEMENT`'s `read_lease_regions` option (DDL-044).

Spanner's DDL parser checks option names and values while it parses, and the
harness reads one such refusal ("Error parsing Spanner DDL statement: ...")
as a rejection but two at once ("Errors parsing ... has errors:") as a
semantic error, which the judge accepts. So `CREATE PLACEMENT p OPTIONS (x =
1)` reads as accepted while `CREATE PLACEMENT p` does not; the rows avoid
statements with two such refusals.

The rows start from `unterminated_statement`, a statement of a script, so
they cover GoogleSQL's scripting and BigQuery statements too; Spanner refuses
them, and so does omni. `select_column_expr`'s `ATTACHED_ALIAS` alternative
gets no row: it refuses every input (`../../testdata/tokens.tsv`).

### Judge revision 5

Revision 4 judged the file with a BigQuery ∪ Spanner union judge, for the
union parser `spanner/parser` was until BigQuery forked off. Revision 5 is
the Spanner emulator alone again, with revision 3's reading. Rejudging moved
783 rows from `accept` to `reject`, the BigQuery forms only the union judge
accepted; omni rejects 781 of them. It also found two Spanner gaps the old
`DialectSpanner` mode had, both closed: bare `INSERT REPLACE` and `INSERT
UPDATE` as statements.

The sections below record the union era.

### Judge revision 4 (`PinRole: "parser"`, BigQuery ∪ Spanner union judge)

At judge revision 3 (Spanner emulator alone, plus bare `QUALIFY` table alias
closure in union mode), `compat.tsv` had **8,937 / 9,463 (94.44%)** agreement
with **0 `Only Spanner accepts`** and **526 `Only omni accepts`** (477
generated + 49 `suite:` rows). Judge revision 4 judges each statement as
`accept` iff either the pinned Spanner emulator `1.5.54` (rev 3 rules) OR live
BigQuery dry-run (`bytebase-dev:US`, matching `bigquery/parser/compat_test.go`)
accepts it, moving **709** rows from `reject` to `accept`:

- **506** previously `Only omni accepts` rows (457 generated + all 49 `suite:`
  BigQuery rows) agreed immediately (`accept` / `accept`), bringing the rev-4
  baseline to **9,240 / 9,463 (97.64%)** with **203 `Only union accepts`** and
  **20 `Only omni accepts`**.
- All **223** remaining union-mode gaps were then closed in `DialectDefault`
  (**9,240 → 9,463 / 9,463 = 100.00%**) while keeping both `bigquery/parser`
  (**8,837 / 8,837 = 100.00%**) and `spanner/parser`
  (**8,837 / 8,837 = 100.00%**) at 100.00%:
  - **203 `Only union accepts` closed:** Enabled full BigQuery syntax in union
    mode (`!p.isSpanner()`) — `FROM`-first pipe queries, `MATCH_RECOGNIZE` table
    clauses, `CREATE MODEL` BigQuery ML clauses, `CREATE/DROP EXTERNAL
    TABLE/SCHEMA`, `CREATE/DROP CONNECTION`, `CREATE/DROP CONSTANT`,
    `CREATE/ALTER/DROP PRIVILEGE RESTRICTION`, `ALTER SEARCH|VECTOR INDEX`
    (preserving Spanner DDL-049 `{ADD|DROP} STORED COLUMN`), `ALTER [ALL] ROW
    ACCESS POLICY/POLICIES`, untyped `CREATE TABLE` columns, `COPY`/`CLONE`
    time-travel clauses, `CREATE INDEX` table aliases, adjacent string literal
    concatenation in `GRANT`, `DROP INDEX ... ON <table>`, `CREATE OR REPLACE
    SEQUENCE`, and BigQuery `CREATE PROPERTY GRAPH` modifiers/optional `NODE
    TABLES`.
  - **20 `Only omni accepts` closed:** Rejected forms refused by both Spanner
    and BigQuery — pipe operators on `DESCRIBE` (2 rows), non-assignment
    `ANALYZE OPTIONS(...)` (1 row), trailing Spanner `PRIMARY KEY(...)` or
    `ROW DELETION POLICY(...)` on columnless `CREATE TABLE` (13 rows), and
    non-identifier/non-parenthesized index keys on Spanner-only
    `NULL_FILTERED` / `, INTERLEAVE IN` indexes (4 rows).

### Judge revision 3 (`PinRole: "parser"`)

At judge revision 2 (`default` pin, after `grammar_refusals.go` moved 90 rows
to `reject`), `compat.tsv` had **7,909 / 9,463 (83.58%)** agreement with
**1,103 `Only Spanner accepts`** and **451 `Only omni accepts`**. Judge
revision 3 moves 1,185 rows from `accept` to `reject` (**1,102** false
`Only Spanner accepts` resolved, leaving only `TABLE a() QUALIFY` where
`QUALIFY` is a BigQuery query clause in union mode and a non-reserved
identifier in Spanner; **+83** `Only omni accepts`, 76 generated + 7 suite,
bringing agreement to **8,928 / 9,463 = 94.35%**):

- **Parse-error fallback (`Statement not supported: <Kind>Stmt`):** When
  GoogleSQL's `ParseStatement` succeeds on an unsupported statement kind, it
  reports the AST node name ending in `Statement` or `Assignment`
  (`SubpipelineStatement`, `MergeStatement`, `BeginStatement`,
  `SingleAssignment`, etc. — kept `accept`). When `ParseStatement` fails on an
  unsupported statement kind, GoogleSQL's fallback path
  (`GetNextStatementKind`) replaces the syntax error with the
  `ResolvedNodeKind` name ending in `Stmt` (`SubpipelineStmt`, `MergeStmt`,
  `AuxLoadDataStmt`, `BeginStmt`, `UndropStmt`, etc.). Judge rev 3 classifies
  `...Stmt` as `reject`.
- **Lookahead tokenizer and dynamic grammar-action refusals:** Rejects
  `EXCEPT must be followed by ALL, DISTINCT, or "("` and `DYNAMIC
  LABEL|PROPERTIES cannot be used more than once in a single element table`
  (`MakeSyntaxError` in `googlesql.tm:2862`).
- **`CREATE VIEW` missing `SQL SECURITY`:** Re-asks with a `SQL SECURITY
  INVOKER` stand-in before `AS` so broken view bodies (`CREATE VIEW a AS NULL`,
  `CREATE VIEW a AS SELECT`) reject on syntax while valid view bodies
  (`CREATE VIEW a AS SELECT NULL`) stay accepted.

Most of the 526 `Only omni accepts` rows are intentional BigQuery-only syntax
that `spanner/parser` accepts in union mode (`DialectDefault`) — scripting
(`DECLARE`, `IF`, `CASE`, `BEGIN ... END`, etc.), BigQuery DDL, `QUALIFY`,
`IS [NOT] DISTINCT FROM`, `|> MATCH_RECOGNIZE(...)`, `LOAD DATA` with
`PARTITIONS(...)` / `TEMP TABLE`, etc. (The 451 rows present at judge rev 2 are
catalogued below.) Eight genuinely invalid union-mode forms exposed by judge
rev 3 were closed in `DialectDefault` without removing any BigQuery syntax
(**8,928 → 8,936 / 9,463 = 94.43%**):

- Truncated `|> EXPORT` / `|> EXPORT AS` without `DATA` (2 rows; `googlesql.tm:4273`)
- Mismatched end label `a : BEGIN END zzz` (1 row; `googlesql.tm:1242-1249`)
- Empty `()` column list in `LOAD DATA INTO a() FROM FILES()` (1 row; `googlesql.tm:3295`)
- Bare `|> INSERT REPLACE a` / `|> INSERT UPDATE a` without `OR` in pipe `INSERT` (4 rows; standalone `INSERT REPLACE|UPDATE INTO t ...` remains accepted in union mode)

## Historical notes (judge revision 2)

The sections below document the judge rev 2 audit before judge rev 3 classified
`Statement not supported: <Kind>Stmt` fallbacks, lookahead `EXCEPT`, duplicate
`DYNAMIC LABEL|PROPERTIES`, and masked `CREATE VIEW` bodies as `reject`.
(`SELECT NULL.*`, `|> EXTEND NULL.*`, and `DROP MODEL a` were also closed in
the parser prior to judge rev 3.)

## Most of those families are not parse verdicts either

For `MERGE`, `LOAD DATA`, `BEGIN … END`, `IMPORT`, `MODULE`, `EXPORT DATA`,
`EXPORT MODEL`, `SHOW`, `EXECUTE IMMEDIATE`, `UNDROP` and `CLONE DATA` the
emulator answers `Statement not supported: <Kind>Stmt` from statement dispatch,
BEFORE it parses the body, and the judge reads that as an acceptance. It says so
for `MERGE a USING a ON NULL WHEN MATCHED THEN DELETION` and
`MERGE a USING a ON NULL WHEN MATCHED THEN DELETE` alike, for
`LOAD DATA INTO FROM FILES()` and `SHOW a zzz`, and for `UNDROP CONNECTION`
with no name. So those rows say `accept` whatever the text is, and a near-miss
in them (a dropped or swapped token) cannot be closed by matching the cell: the
cell is not a verdict. `UNDROP` adds a grammar reason of its own —
`undrop_statement`'s action allows only the SCHEMA kind and reports
"UNDROP <kind> is not supported" for the rest, so the `UNDROP CONNECTION …`
rows the generator wrote are text the grammar refuses.

The `TABLE` clause is the exception, and the work here: the emulator parses a
TABLE statement in full before it answers `TABLE statements are not supported`,
so it rejects `TABLE a zzz`, `TABLE a AS x`, `TABLE a PIVOT(…)` and `TABLE a(`
on syntax and accepts the `tvf_argument` shapes. Its rows are a real verdict,
and omni now reads `TABLE a`, `TABLE a(…)` with DESCRIPTOR / MODEL / CONNECTION /
INPUT TABLE / `name => TABLE t` arguments, and `TABLE a WHERE …` — 55 rows. The
same change closed `SHOW` (its three forms, including `SHOW [MATERIALIZED]
VIEWS … FROM … LIKE`, 11), `IMPORT MODULE|PROTO` with `AS`/`INTO` aliases and
`OPTIONS` (9), `MODULE <path> [OPTIONS]` (5), `EXPORT TABLE METADATA FROM` and
`EXPORT DATA AS GRAPH` (3), and the empty `EXCEPTION WHEN ERROR THEN END`
handler that `statement_list`'s empty alternative allows (2): 85 rows in all,
none of them a cell the recording could not judge. One TABLE near-miss is
left: `TABLE a() QUALIFY`, where QUALIFY is the clause word and omni reads no
alias there.

## `SET` is judged on one alternative only

`set_statement` splits by target, and the recording parses exactly one target:
`SET a = NULL` and its near-misses (`SET a = NULL zzz` rejects) are a parse
verdict, and omni agrees on all of them. The other two are not. `SET @@ a = NULL
zzz` and `SET TRANSACTION READ ONLY zzz` are recorded `accept` while their
valid forms are too, so the inserted token changed nothing: the emulator
answers `Statement not supported` for a system-variable or transaction-mode SET
before it parses the assignment. The 10 rows left therefore say `accept` for
text the grammar refuses (`SET TRANSACTION ONLY`, `SET TRANSACTION READ`, `SET
@@ a =`), and matching them would mean accepting what `set_statement` does not
derive.

## The graph rows are a parse verdict

The graph sub-language is the largest family the grammar derives and omni did
not read: `graph_expression`'s endpoint (`IS [NOT] SOURCE|DESTINATION [OF]`) and
label (`IS [NOT] LABELED`) predicates, the braced graph subqueries wherever a
query can appear (`SELECT ARRAY { … }`, `SELECT EXISTS { … }`, `SELECT VALUE {
… }`, `expr IN { … }`, and a CALL's inline subquery), and the operator, pattern
and CALL alternatives inside them. Because a graph query is reached only through
an expression, every one of those rows was a `reject` at the `{`: 208 rows moved
on the one change.

Unlike the pipe and dispatch families, the emulator judges this one: inside the
same rules it rejects `SELECT EXISTS { MATCH }`, `{ MATCH - MATCH }`, `{ - MATCH
- }`, `{ SHORTEST - }`, `{ - [ ] - }` and `{ < - [ ] - }` (the grammar's
`ValidateNoWhitespace` on an edge pattern's delimiters), `{ (a)* }` (a quantifier
on a node pattern), `{(@1 -) }` (a hint opening a path pattern), `{ -@1 - }` (a
hint between two edges) and `{(COST NULL) }`. Text the grammar derives is
accepted and text it refuses is rejected, so the cells are verdicts. The
`COST` filler tail is the one place the recording disagrees with
`googlesql.tm` — `opt_graph_cost/COST_expression` is recorded `reject` — so omni
reads the word as an element named `COST` and not as the tail, which is what the
grammar's own shift preference and the cell agree on.

The rows left in this family are not verdicts: `DYNAMIC LABEL(NULL) DYNAMIC
LABEL(NULL)` and its `drop-10` twin are recorded `accept` with and without the
closing `)`, and `PROPERTIES ALL COLUMNS EXCEPT a))` is recorded `accept` with a
stray `)`, so the emulator answered something other than a syntax error for a
`CREATE PROPERTY GRAPH` whose `DYNAMIC LABEL` support its DDL parser reports
later (4 rows).

`SELECT NULL.*` (3 rows, and the two pipe `|> EXTEND NULL.*` rows) is a real
verdict the grammar derives from `select_column_dot_star`'s
`expression_higher_prec_than_and "." "*"`; omni's dot-star is path-only, as it
was before either family, so they stay. `DROP MODEL a` is the same shape of real
verdict outside every family here: `MODEL` is one alternative of the grammar's
`schema_object_kind`, and the emulator's DDL parser reads it and reports the
missing model later.

## The keyword rows

`common_keyword_as_identifier` is a word the grammar takes as a name, and omni
reserved six of them: `MATCH_RECOGNIZE` (which the legacy grammar lists through
the placeholder token `KW_MATCH_RECOGNIZE_NONRESERVED`, so its reserved split
cannot see the word), `QUALIFY`, `NOTHING`, `SHORTEST` (missing from the older
grammar's list), and `SAFE_CAST` / `REPLACE_FIELDS` (read only as functions).
The pinned grammar lists all six, and the recording agrees: `SELECT
MATCH_RECOGNIZE`, `SELECT NOTHING`, `SELECT QUALIFY`, `SELECT SHORTEST`, `SELECT
SAFE_CAST` and `SELECT REPLACE_FIELDS` are accepted, each also with a trailing
alias. 12 rows closed. `SELECT EXCEPT` is left: `EXCEPT` is a reserved word, so
the grammar cannot derive that statement, and the emulator answered an empty
select list — an analysis error, not a syntax error — which the judge reads as
an acceptance.

## Rows only omni accepts, measured (475)

Every row whose cells were `accept`/`reject` was measured at the pin: the
emulator's gRPC code (always `InvalidArgument`), the position it reports, and
a control — the form beside it that the same server takes. The recognizer over
the pinned grammar (`enginetest/grammar`'s, the one the rows are generated
from) derives 473 of the 475, near-misses included, so omni's reading is the
grammar's and the emulator's refusal is its dialect. The other two are the
`ALTER SEARCH INDEX … STORED
COLUMN` rows, which Spanner's DDL reference quotes (`truth1/spanner/ddl.md:1006`)
while both the pinned grammar and the emulator's own DDL parser lag it: that
reference, not the grammar, is their authority. The `CREATE PROPERTY GRAPH g
OPTIONS(…) NODE TABLES (T)` row, which neither the grammar (`googlesql.tm:2666`
ends the rule with `opt_options_list`) nor Spanner's graph DDL reference puts
before the content, was omni's and is closed.

Positions are the emulator's own, 1-based columns on these single-line ASCII
statements, so the byte is the column minus one. A query's position is the
`[at 1:N]` GoogleSQL's parser puts on its message; a DDL position is the
`column N` of an `Error parsing Spanner DDL statement:`, which is why the DDL
groups carry no `Syntax error` prefix.

### Closed (1)

`CREATE PROPERTY GRAPH g OPTIONS(x=true) NODE TABLES (T)` — the emulator
answers `Expected "." or keyword NODE but got keyword OPTIONS` at 1:25, and
takes `CREATE PROPERTY GRAPH g NODE TABLES (T) OPTIONS(x=true)` (the control).
The pinned grammar's rule ends `NODE TABLES element_table_list
edge_table_clause? opt_options_list`, Spanner's reference writes
`property_graph_content [OPTIONS (key=value, …)]`, and the legacy ANTLR
`GoogleSQLParser.g4`, which omni had ported, is the only place the list leads.
`TestCreatePropertyGraph_OptionsAfterContent` holds both positions to the
parser, and the `suite:graph-query` row is the cell this closes.

### Measured (475)

| Cause | Rows | The emulator's answer, and the control |
| --- | --- | --- |
| BigQuery scripting statements | 121 | `Syntax error: Unexpected keyword <word>` at 1:1 (`DECLARE` 57, `IF` 12, `CASE` 11, `BREAK` 8, `RAISE` 5, `CONTINUE` 4, `REPEAT` 4, `FOR` 3, `ITERATE` 3, `LEAVE` 3, `LOOP` 3, `RETURN` 3, `WHILE` 3), and `Unexpected "label_1"` at 1:1 for the two labelled blocks: the emulator's query parser has no script statement. `SELECT IF(NULL, 1, 2)` and `SELECT CASE WHEN NULL THEN NULL ELSE NULL END` run, so the words are expressions it has and statements it does not |
| The `QUALIFY` clause | 8 | `QUALIFY is not supported` at 1:17, 24, 25 or 26 (6; the action at `googlesql.tm:6378` with the feature off), and `Expected end of input but got keyword NULL` at 1:21 for the two reserved-keyword forms (`SELECT NULL QUALIFY NULL`, which the emulator reads as `QUALIFY`-the-alias). `SELECT QUALIFY` runs: the word is a name it takes |
| `IS [NOT] DISTINCT FROM` | 4 | `IS DISTINCT FROM is not supported` at 1:13 or 1:20, another action with the feature off. `SELECT NULL IS NULL` runs |
| An alias on a parenthesized outer query | 2 | `Alias not allowed on parenthesized outer query` at 1:15 for `(SELECT NULL) AS a`. `SELECT * FROM (SELECT NULL) AS a` runs: the same alias under `FROM` is derived |
| A dashed path | 2 | `Unexpected "-"` at 1:10 for `DELETE a - a`. `SELECT a - a` runs, and so does `DELETE FROM a WHERE FALSE` |
| Words the emulator reserves | 4 | `Unexpected keyword GRAPH_TABLE` / `LATERAL` at 1:8, `IGNORE` at 1:13 (`SELECT NULL.IGNORE`): the grammar lists them as names (`common_keyword_as_identifier`, `function_name_from_keyword`). `SELECT MATCH_RECOGNIZE` runs, a word the grammar also takes as a name and the emulator does not reserve |
| `CREATE TABLE` with no element list | 55 | `Expecting '(' but found '<word>'` at column 14 when the text ends (`CREATE TABLE a`, 20) or 16 at the first clause word (`PRIMARY`, `AS`, `LIKE`, `CLONE`, `COPY`, `PARTITION`, `CLUSTER`, `DEFAULT`, `ROW`, `OPTIONS`; 35). `CREATE TABLE a(a INT64) PRIMARY KEY(a)` applies |
| An element list with no primary key | 8 | `Must specify either table or column PRIMARY KEY`, with no position and no `Syntax error` phrase, under the DDL prefix: this is a check of the parsed statement, and the judge reads the prefix as a rejection. `CREATE TABLE a(a INT64) PRIMARY KEY(a)` applies, and `CREATE TABLE a()` and `CREATE TABLE a(a a,)` are the refused texts |
| An index key that is not a column name | 39 | `Expecting ')' but found 'NULL'` at 21 (33; `CREATE INDEX a ON a(NULL)`), at 28 after `UNIQUE`/`SEARCH` (3), at 35 after `NULL_FILTERED` (2), and `found 'ALL'` at 38 for `CREATE SEARCH INDEX my_index ON ds.t(ALL COLUMNS)`: Spanner reads a key as a name and not the grammar's `column_ordering_and_options_expr`. `CREATE INDEX a ON a(a)` and `CREATE INDEX a ON a(a) STORING(b)` are read (`Table not found` is an acceptance) |
| A materialized, approx or recursive view, or `REPLICA OF` | 11 | `Encountered 'MATERIALIZED'`, `'APPROX'` or `'RECURSIVE'` while parsing `create_statement` at column 8. `CREATE VIEW a AS SELECT NULL` parses (it is refused later for a missing `SQL SECURITY`, which the judge reads as an acceptance) |
| A view column list or options | 6 | `Expecting 'AS' but found '('` at column 14 (5; `CREATE VIEW a(a) AS …`) and `found 'OPTIONS'` at 15 (1). `CREATE VIEW a AS SELECT NULL` parses |
| `CREATE FUNCTION` without `RETURNS` or a body | 21 | `Encountered 'EOF' while parsing: create_function_statement` at 19 (11; `CREATE FUNCTION a()`), 26 (`REMOTE`), 37 (`NOT DETERMINISTIC`) or 40 (`SQL SECURITY INVOKER`); `Expecting '(' but found ''a''` at 24 for the string body (3); and the two body-only `suite:bq-ddl` rows, which the emulator rewrites to `CREATE FUNCTION … RETURNS  SQL SECURITY INVOKER AS (…)` and refuses as `Expected end of input but got keyword SECURITY` at 1:36 and 1:46. `CREATE FUNCTION f() RETURNS INT64 AS (1)` and `CREATE FUNCTION f(x INT64) RETURNS INT64 AS (x)` apply |
| A function parameter or return type the emulator does not know | 19 | `Encountered '<word>' while parsing: function_data_type` at 20, 21, 22, 23 or 29: `a` and `zzz` are the grammar's `type` (`path_expression`) and no built-in type this parser has. `CREATE FUNCTION f(x INT64) RETURNS INT64 AS (x)` applies |
| A function's determinism or language clause | 11 | `Encountered 'DETERMINISTIC'` (2 at 21 and 1 at 45), `'IMMUTABLE'`, `'STABLE'` and `'VOLATILE'` at 21 (5 there in all), and `LANGUAGE <name>` at 28, 30 (3) or 58 (`LANGUAGE a`, `LANGUAGE js`). The `RETURNS … AS (…)` form applies without them |
| `AGGREGATE`, `TEMP` or `TABLE FUNCTION` | 4 | `Encountered 'AGGREGATE'` / `'TEMP'` while parsing `create_statement` at column 8 (3), and `Expecting '(' but found 'a'` at 23 for `CREATE TABLE FUNCTION a()`: object kinds this server does not have. `CREATE FUNCTION f() RETURNS INT64 AS (1)` applies |
| `CREATE PROCEDURE` | 27 | `Encountered 'PROCEDURE' while parsing: create_statement` at column 8, whatever follows it. `CREATE FUNCTION f() RETURNS INT64 AS (1)` applies |
| `CREATE [ROW ACCESS] POLICY` | 13 | `Encountered 'ROW' while parsing: create_statement` at column 8. `CREATE TABLE a(a INT64) PRIMARY KEY(a)` applies |
| `GRANT` and `REVOKE` | 28 | `Encountered 'ALL'`, `'ALTER'`, `'a'`, `'SEQUENCE'` or `'PRIVILEGES'` while parsing `grant_statement` at column 7 (23), `Encountered 'a' while parsing: target_type` at 17 for `GRANT SELECT ON a TO 'a'` (1), and `Encountered 'ALL' while parsing: revoke_statement` at column 8 (4): Spanner's grant is `GRANT <privilege> ON TABLE <path> TO ROLE <role>`, not the grammar's `ON <path> TO '<grantee>'`. `GRANT SELECT ON TABLE a TO ROLE r` and `REVOKE SELECT ON TABLE a FROM ROLE r` apply |
| `ANALYZE` with a table, column list or options | 35 | `Expecting 'EOF' but found '<word>'` at column 9, from `ANALYZE a` and `ANALYZE a(a)` to `ANALYZE OPTIONS(a = PROTO)`. `ANALYZE` alone applies; `ANALYZE TABLE a` is refused the same way |
| An empty `OPTIONS()` list | 6 | `Encountered ')' while parsing: identifier` at 25 (`CREATE SCHEMA a OPTIONS()`), 27 (`CREATE SEQUENCE a OPTIONS()`) and 29 (`CREATE FUNCTION a() OPTIONS()`): the emulator's list takes one option at least. `CREATE SEQUENCE a OPTIONS(sequence_kind='bit_reversed_positive')` applies |
| `ALTER` on a generic entity | 8 | `Encountered 'a' while parsing: alter_statement` at column 7 for `ALTER a SET AS 'a'`, `ALTER a SET AS JSON 'a'`, `ALTER a SET OPTIONS()` and `ALTER a a SET AS 'a'`: the action list follows an object kind, not a bare path. `ALTER TABLE a ADD COLUMN b INT64` applies |
| `ALTER [SEARCH\|VECTOR] INDEX` | 3 | `Encountered 'SEARCH'` at column 7 (2) and `Encountered 'ON'` at 24 (`ALTER VECTOR INDEX idx ON ds.t REBUILD`). `ALTER INDEX a ADD STORED COLUMN c` applies, which is the action the two `suite:bq-ddl` rows name (`truth1/spanner/ddl.md:1006`). The pinned grammar's `alter_index_action` is BigQuery's and lacks `STORED`, and the emulator's DDL parser has no `ALTER SEARCH INDEX` at all, so Spanner's DDL reference is the authority for those two rows |
| `DROP` | 20 | `Encountered '<word>' while parsing: drop_statement` at column 6 (14; `ALL`, `APPROX`, `DATABASE`, `EXTERNAL`, `MATERIALIZED`, `PROCEDURE`, `PROJECT`, `ROW`, `SNAPSHOT`, `a`) or 8 (6), and `Expecting 'EOF' but found 'ON'` at 21 or 23 for `DROP [SEARCH\|VECTOR] INDEX … ON …` and `DROP TABLE FUNCTION …`. `DROP TABLE a` and `DROP INDEX a` are read (`Table`/`Index not found` is an acceptance) |
| `CREATE` on a generic entity | 8 | `Encountered 'PROJECT'` / `'a'` while parsing `create_statement` at column 8 (7; `CREATE a a`, `CREATE a a AS 'a'`, `CREATE PROJECT a`) and `Expecting 'GROUP' but found 'a'` at 17 (`CREATE LOCALITY a`). `CREATE TABLE a(a INT64) PRIMARY KEY(a)` applies |
| `CREATE SNAPSHOT TABLE` | 3 | `Encountered 'SNAPSHOT' while parsing: create_statement` at column 8. `CREATE TABLE a(a INT64) PRIMARY KEY(a)` applies |
| `CAPACITY`, `RESERVATION`, `OR REPLACE TABLE FUNCTION` | 3 | `Encountered 'CAPACITY'` / `'RESERVATION'` at column 8, `'TABLE'` at 19. `CREATE TABLE a(a INT64) PRIMARY KEY(a)` applies |
| `TRUNCATE` | 4 | `Encountered 'TRUNCATE' while parsing: ddl_statement` at 1:1: the harness routes the first word to the DDL parser, which has no such statement. `DELETE FROM a WHERE FALSE` runs |
| `RENAME` without `TABLE` | 2 | `Expecting 'TABLE' but found 'a'` at column 8. `RENAME TABLE a TO b` is read (`Table not found` is an acceptance) |

The 42 `suite:` rows among these are `suite:scripting` (21), `suite:bq-ddl`
(20) and `suite:select` (1); the closed row was `suite:graph-query`'s. No cell
here is the emulator's silence: each group's control is a statement it answered.

## The pipe rows are not a parse verdict

GoogleSQL's pipe syntax is a family the grammar reaches three ways:
`sql_statement_body_no_pipe_suffix → subpipeline_statement →
subpipeline_no_parens` (a pipe suffix alone, `|> AGGREGATE`),
`query → query pipe_and_pipe_operator*` (`SELECT … |> WHERE …`), and
`sql_statement_body_maybe_pipe_suffix pipe_and_pipe_operator*` (`DESCRIBE … |>
…`). The first is a complete statement, not a fragment — the bare suffix is an
alternative of the start rule — so omni parses it (`SubpipelineStmt` in
`spanner/ast`) rather than refusing it.

A recording cannot judge this family. The emulator answers a statement
containing a pipe operator as one it does not run, without a syntax error, so
`compat.tsv` records `accept` for every one of them — including text the grammar
refuses, such as `|> ALL`, `|> FULL`, `|> WHERE`, or a pipe statement with a
required token dropped (`|> ORDER BY`). Of the 1,165 pipe rows, 413 are valid
per the grammar and now agree; the other 752 are the only-Spanner rows left, and
the grammar, not the recording, is what refuses them. Two of those are the
select-list `expr.*` form over a non-path expression (`|> EXTEND NULL.*`): the
grammar derives it, and omni's dot-star support is path-only, as it was before
this family.

Closing a pipe row therefore means matching the grammar, not the cell: the cell
says `accept` for valid and invalid text alike. A later recording that reads a
real parse verdict for a pipe statement is what would change these rows.

## Suites

The hand-written corpora that held `Parse` to a live emulator are `suite:`
rows, judged by the same recording: the fixtures of the `*_oracle_test.go`
files that compared `Parse` of a statement with the emulator's verdict on it.

| Suite | Rows | Agree | Only union accepts | Only omni accepts |
| --- | ---: | ---: | ---: | ---: |
| `suite:select` | 95 | 95 | 0 | 0 |
| `suite:pivot-unpivot` | 49 | 49 | 0 | 0 |
| `suite:dml` | 68 | 68 | 0 | 0 |
| `suite:ddl` | 91 | 91 | 0 | 0 |
| `suite:spanner-ddl` | 81 | 81 | 0 | 0 |
| `suite:utility` | 67 | 67 | 0 | 0 |
| `suite:scripting` | 47 | 47 | 0 | 0 |
| `suite:graph-query` | 86 | 86 | 0 | 0 |
| `suite:bq-ddl` | 24 | 24 | 0 | 0 |
| `suite:export-load-clone` | 18 | 18 | 0 | 0 |

All 626 `suite:` rows agree (**626 / 626 = 100.00%**) under the union judge. A
statement added to a suite is a row with `?` in both verdict cells, judged by a
recording:

```text
?	?	-	suite:select	<statement>
```

The other oracle tests compare something else with the emulator's verdict (an
expression or a type parsed alone, a Spanner docs block with its comments and
line breaks, which a row cannot hold) and replay tapes (`tapes/`): the
harness's answer for each statement they send, recorded with `OMNI_RECORD=1`.

## What the judge needs

- **The harness decides.** `spanner/spanner-oracle` sends each statement to
  the emulator by its first word: a query runs, DML runs in a transaction it
  rolls back, and DDL is applied to a database of its own
  (`../../spanner-oracle/README.md`). Its verdict is the judge's
  (`spannerVerdict` in `../compat_test.go`).
- **A rejection is the parser's error:** `Syntax error:` from GoogleSQL's
  parser, also where the emulator parses SQL inside DDL (`CREATE PROPERTY
  GRAPH`, a function's body); a message GoogleSQL's grammar writes when an
  action refuses input, as `QUALIFY is not supported` where the emulator
  leaves the feature off (`grammar_refusals.go` in the harness, generated
  from `googlesql.tm` by `../cmd/tm2g4`); `Error parsing Spanner DDL
  statement:` from the emulator's own DDL parser; or, at judge revision 3,
  GoogleSQL's `Statement not supported: <Kind>Stmt` parse-error fallback,
  lookahead `EXCEPT` and duplicate `DYNAMIC LABEL|PROPERTIES` syntax refusals,
  and syntax errors exposed after injecting a `SQL SECURITY INVOKER` stand-in
  on `CREATE [OR REPLACE] VIEW` missing `SQL SECURITY`. Every other answer is
  an acceptance, post-parse `Statement not supported: <ASTKind>Statement` /
  `Assignment` and `Unimplemented` (a query that returns a `STRUCT` column)
  included.
- **Every statement sees an empty schema.** The emulator answers
  `CREATE PROPERTY GRAPH a` by the names the schema holds before it parses the
  rest, so the harness applies DDL to a database it replaces after each DDL
  statement that applies.
- **One statement, no `;`.** The emulator's DDL path refuses a terminator.

## Recording and verification commands

Offline replay (CI default, zero Docker, zero network):

```bash
go test ./spanner/parser/ -run '^TestCompatibility$'
```

Live recording against pinned `gcr.io/cloud-spanner-emulator/emulator:1.5.54`
(about a minute for all 10,171 rows):

```bash
OMNI_RECORD=1 go test -v -timeout 60m ./spanner/parser/ -run '^TestCompatibility$'
```

Read-only live verification against the pinned emulator (re-judges every row,
writes nothing, fails on any diff):

```bash
OMNI_RECORD=verify go test -v -timeout 60m ./spanner/parser/ -run '^TestCompatibility$'
```
