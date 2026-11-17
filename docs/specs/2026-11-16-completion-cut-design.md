# The completion cut

2026-11-16 · Spike, one engine landed (MySQL)

Completion is the last capability of `docs/design/engine-code-to-omni.md`'s target state
without a contract. That document and `docs/engine-capability-guide.md` (under L6) both
state the cut — "omni emits parser-native token/rule/scope/intent signals; Bytebase resolves
live metadata into product `Candidate` values" — but neither says what the signals are, what
the boundary call is, or which engines already follow it. This document answers those three
questions, lands the mechanism for MySQL, measures it against Bytebase's own completion
corpus, and prices the rest of the family.

The design document and every `backend/plugin/...` path below live in the Bytebase
repository; every other path is relative to this one.

## The cut

**Omni emits what the text says. Bytebase emits what the metadata says.**

Omni answers, at a byte offset in a statement that may not parse:

- which grammar positions accept the cursor (`CandidateSet`: token and rule candidates),
- which relations the parser can see there, with the names and aliases as written
  (`Scope`), and which columns of them the text itself names (`Columns`, from a CTE's or a
  derived table's select list),
- the CTE definitions of the statement (`CTEs`),
- the select-list aliases in scope, which `ORDER BY`, `GROUP BY`, and `HAVING` may name
  (`SelectAliases`),
- the object class the grammar position and the qualifier imply (`Intent`),
- the partial identifier at the cursor (`Prefix`) and the statement it is in
  (`StatementLoc`).

Bytebase answers, with its store and its product types:

- which databases, tables, views and columns exist — `GetDatabaseMetadataFunc`, the
  metadata seam it already owns,
- the columns of a derived table or CTE whose body the parser cannot read statically
  (`SELECT * FROM t`): it runs a query span over live metadata,
- quoting for its editor, function signatures and doc comments, keyword filtering by scene
  (a read-only console must not offer `DROP`), priority, and the product `Candidate` type
  with its `CandidateType` names,
- case sensitivity, which for MySQL is a server variable (`lower_case_table_names`) and
  therefore metadata, not grammar.

What crosses is one struct of plain values, defined per engine as
`<engine>/parser.CompletionContext`, plus the engine's own `ObjectKind` and
`RangeReferenceKind` vocabularies. No catalog object, no metadata proto, no `storepb`
type, and no product type crosses. The names are `string`s as written, so a caller folds
case itself, and a kind the caller does not know is ignored rather than an error.

**Why this shape.** Bytebase's metadata is live and per-instance; omni's catalog is a
projection of it. Passing the projection across (`omni/<engine>/catalog.Catalog`) looks like
the same thing and is not: it makes Bytebase resolve its store into omni's mutable
simulation type on every keystroke, it puts the catalog's adapter surface on the boundary,
and it makes the suggestion depend on how faithfully the catalog was populated rather than
on narrowing a signal set against a lookup. Two engines do that today (§ Where the engines
stand), and the cost is visible in their adapters.

**Rejected: a shared `omni/completion` contract package now.** Four engines already have
this struct, and they differ where their engines differ: PostgreSQL has `Catalog` and
`Lateral`, SQL Server has four-part names and `VALUES`/table-variable references, Oracle has
packages, synonyms, database links, and a `RANGE`-family vocabulary. A shared type written
from MySQL alone would either carry fields three engines never set or force the others to
reshape. The path is `omni/analysis`'s: land the shape per engine, promote it once two
engines agree field for field. The parity evidence for that promotion is the field names
below, which follow the SQL Server and Oracle spellings so promotion is a rename, not a
redesign.

**Rejected: a callback-style contract** (Bytebase passes a `LookupColumns(database, table)`
func into omni). It hides the metadata behind a function so the boundary no longer shows
what omni needs, makes cost invisible to the caller, and makes omni's answer depend on call
order. Signals keep the caller in control of the lookup it pays for.

**Rejected: keeping MySQL's completion in Bytebase.** That is the status quo, and it is what
that document calls out: 1,267 lines of MySQL completion deciding which object class a
position accepts, which is grammar.

## Where the engines stand

Two boundary styles are in production in Bytebase today.

**Signals (the cut).** `pg`, `redshift`, `mssql`, and `oracle` expose
`<engine>/parser.CollectCompletion(sql, cursorOffset) *CompletionContext`; the Bytebase
adapter calls it, resolves names against `GetDatabaseMetadataFunc`, and builds
`base.Candidate`. `pg`, `redshift`, and `mssql` also re-export the struct at the engine root
(`<engine>/completion.go`) as type aliases.

**Metadata as catalog (not the cut).** `tidb` and `mongo` build their own
`omni/<engine>/catalog.Catalog` from Bytebase metadata per request and call
`omni/<engine>/completion.Complete`, which resolves the candidates inside omni
(`backend/plugin/parser/tidb/completion.go`, `.../mongodb/completion.go`). Those two
adapters are a cut each, not a port, and are the follow-up's first task.

**Nothing (the gap).** `mysql`, `mariadb`, `oceanbase`, and `tidb` parse with a completion
mode that emits only the token/rule candidates plus two byte positions
(`CTEPositions`, `SelectAliasPositions`) for Bytebase to re-parse. Bytebase then computes
scope, intent, CTE definitions, and select aliases itself, from a token scan plus a
re-parse of the statement fragment (`backend/plugin/parser/mysql/completion.go`). The
`<engine>/completion` package next to each of them resolves a catalog and is **not** on the
boundary: Bytebase imports it only for TiDB and Mongo.

So "which engines follow the cut" is: four parser layers, and two engine adapters that
violate it.

## The MySQL mechanism

`mysql/parser/complete_context.go` adds:

```go
func CollectCompletion(sql string, cursorOffset int) *CompletionContext
```

`sql` is the statement text around the cursor and `cursorOffset` is the cursor's byte offset
into it. Nothing else is passed: there is no scene, no default database, no metadata, and no
configuration. A scene flag (`SceneTypeQuery`) is a product decision; a default database is
metadata. `Prefix` and the returned `StatementLoc` let a caller reslice the text itself if
it needs offsets rather than names.

The reader is a token walk, not an AST walk, because the text is being typed: MySQL's `Parse`
is strict and returns no partial tree, and completion must survive `SELECT * FROM t1 JOIN |`.
It reuses `Tokenize` and `Split` (through `Collect`), reads only the statement holding the
cursor, and never parses a fragment twice. It is linear in the statement's tokens: depths and
parenthesis matches are computed once per token slice, and a group nested past the parser's own
`maxNestingDepth` is skipped rather than descended into.

The fields, and why each is on the boundary:

| field | what it carries | why it is grammar, not metadata |
|---|---|---|
| `Candidates` | token candidates (keywords) and rule candidates (`table_ref`, `columnref`, `func_name`, `charset`, `engine`, ...) | the grammar's own accept set at that position |
| `Scope.LocalReferences` / `OuterReferences` / `References` | the cursor query block's relations, the correlated outer blocks, and both | MySQL reads an unqualified name against the innermost block, then outward; a derived table's body sees no sibling unless `LATERAL` |
| `RangeReference` | `Kind`, `Database`, `Table`, `Alias`, `Columns`, `AliasColumns`, `Loc`, `BodyLoc`, `Lateral` | the reference is read off the text; `Columns` is filled only when the text names them (a CTE's column list or select list, a derived table's alias list or select list) |
| `CTEs` | each `WITH` definition with its columns and body extent | MySQL 8.0's `WITH [RECURSIVE]` grammar |
| `SelectAliases` | the aliases of the cursor query block's select list | whether an item declares an alias is grammar (`AS x`, a bare alias, a string alias) |
| `Intent.ObjectKinds` | the object classes the next name part may be, narrowed by the parts already typed | `db.t.|` accepts only a column; `t.|` accepts a table, a view, or a column; `|` accepts a database, a table, a view, or a column |
| `Intent.Qualifier` | the dot-separated parts before the cursor, flat | MySQL reads `a.b` as `table.column` in a select list and as `db.table` in a `FROM`, so the signal does not pretend to know which |
| `Prefix` | the identifier the cursor is inside, or ends | the caller filters with it; inside backticks it is the quoted text |
| `Scope.DMLTarget` | the table an `UPDATE`, `DELETE`, `INSERT`, or `REPLACE` writes | column completion in `SET`, `WHERE`, and the column list needs it |
| `StatementLoc` | the statement holding the cursor | a caller slicing by statement cannot rely on `;` in text it has not lexed |

`ObjectKind` covers MySQL's object classes, not the product's suggestion categories:
`database`, `table`, `view`, `column`, `index`, `function`, `procedure`, `trigger`, `event`,
`type`, `charset`, `engine`, `variable`, `selectalias`, `constraint`. `selectalias` exists because MySQL
lets `ORDER BY`, `GROUP BY`, and `HAVING` name a select-list alias; Bytebase offers those as
column candidates.

### What Bytebase's adapter does

The call, exactly:

```go
// Bytebase: backend/plugin/parser/mysql/completion.go, replacing the token
// scan and the fragment re-parse.
ctx := mysqlparser.CollectCompletion(sql, byteOffset) // sql = the caret's statement, byteOffset its offset in sql
```

and then, unchanged from today's switch: `ctx.Candidates` gives keywords and rule names,
`ctx.Intent.ObjectKinds` with `ctx.Intent.Qualifier` replaces `determineQualifier` and
`determineColumnRef`, `ctx.Scope.LocalReferences` with `ctx.Scope.OuterReferences` replaces
`collectLeadingTableReferences`/`collectRemainingTableReferences`,
`parseTableReferences`, and `extractFromItem`, `ctx.CTEs` replaces `fetchCommonTableExpression`,
`ctx.SelectAliases` replaces `fetchSelectItemAliases`, and `ctx.Prefix` replaces the
adapter's prefix extraction. The metadata calls, quoting, keyword filtering, and
`base.Candidate` construction stay. Derived tables and CTEs with no statically readable
columns keep today's query-span path.

The adapter also changes shape as `parser/redshift` did: `Completion` becomes one function
over the signals instead of a `Completer` with a token stream, and the file loses the
truncation workarounds (`collectLeadingTableReferences`'s paren-level stack,
`collectRemainingTableReferences`'s forward scan, the `parseTableReferences` prefix loop, and
`statementBoundedText`, all of which exist because the adapter re-parses fragments).

## The corpus

`mysql/parser/complete_context_test.go` holds 97 cases over
`sql` × cursor positions: select list, `WHERE`, `ON`, joins with and without aliases,
cross and `STRAIGHT_JOIN`, derived tables (with and without an alias column list, and
unterminated), CTEs (with a column list, from the body's select list, `RECURSIVE`, two in a
row, a cursor inside a body, and a `WITH` nested in a derived table shadowing the
statement's), correlated and `LATERAL` subqueries, sibling subqueries at the same depth,
`UPDATE` (single and join), `DELETE` (single, target list with a join `USING`, and the
`FROM ... USING` form), `INSERT` (target, `INSERT ... SELECT`), `JSON_TABLE` (with a plain
and a nested `COLUMNS` clause), `ORDER BY`/`GROUP BY` aliases versus `WHERE` (including
`WHERE` after a closed subquery), a set operation's global `ORDER BY` against its first
arm, `TABLE t` as a query primary, prefixes with and without backticks (including a quoted
name with a space and a non-ASCII name), index hints and `PARTITION` (which must not read
as aliases), statement boundaries, `UNION` arms, a parenthesized join (with and without an
alias), DDL positions (`type_name`, `index_ref`, `engine`, `charset`, `view_ref`,
`database_ref`, `constraint_ref`), `DROP INDEX ... ON`, an `EXPLAIN` wrapper around a
DML statement, a join condition's visible operands, an index hint's `FOR GROUP BY`, an ODBC
`{ OJ ... }` escape, a script with its own `DELIMITER`, and a cursor inside a text literal. A second test pins
the columns the reader reports for a reference — nested `WITH` bodies, a shadowing clause,
expression items that name no column, qualified names, and `JSON_TABLE`. A third walks every
cursor offset of eight scripts, including 200 nested parentheses and unterminated strings, and
a fourth runs 10,002 nested parentheses, `{ OJ ... }` escapes, `NESTED PATH ... COLUMNS (`
clauses, and `WITH` bodies through the reader: each recursion stops at `maxNestingDepth`
instead of the process stack.

Three positions the corpus pins as gaps in the *grammar*, not in the reader:

- `CHARACTER SET |` and `SET a.|` offer no rule today, so no object kinds. Both behave the
  same in Bytebase's current adapter.
- An aliased parenthesized join, `FROM (t1 JOIN t2 ON ...) j`, does not parse: MySQL accepts
  it. The reader still names `j` as a reference, but the grammar offers nothing at the caret.
  Fixing that is a parser change with a `compat.tsv` row, and it belongs in its own PR — this
  one does not move what the parser accepts.

## Parity

Bytebase's MySQL completion corpus (`backend/plugin/parser/mysql/test-data/test_completion.yaml`,
200 cases) replayed through a throwaway Bytebase worktree with `go mod replace` to this
tree. The harness runs Bytebase's `NewStandardCompleter(...).completion()` (falling back to
the tricky completer, as `Completion` does) and compares, per case, what the signals *can*
explain: scope, object kinds, qualifier, select aliases, CTE definitions, and prefix.

| signal | agreement |
|---|---|
| relation scope | 167/200 |
| object kinds | 197/198 |
| qualifier | 200/200 |
| select aliases | 200/200 |
| CTE definitions | 196/200 |
| prefix | 200/200 |

**What the harness cannot compare, and why.** It cannot compare the candidate lists
themselves without a database. Bytebase's candidates come from `GetDatabaseMetadataFunc`;
the harness fakes that with four tables and a view, so the candidate lists it produces prove
plumbing, not metadata. What it compares instead is the *input* to that resolution. The
candidates' *kinds* are covered indirectly (object kinds 197/198): if the signals name a kind
the adapter never asks the store for, that divergence shows up as a rule/intent mismatch.

The 33 scope and 4 CTE disagreements, classified:

- **20 — omni binds a reference Bytebase loses.** 15 of them are a statement truncated at the
  cursor (`FROM t1 JOIN |`, `FROM t1, |`, `WHERE |`, `SET c1 = |`, `CREATE INDEX i ON t1 (|`,
  `FROM t|`); the rest are complete statements whose references Bytebase never collects
  (`UPDATE t1 SET |`, `LIMIT |`, an `UPDATE` target a correlated subquery may name, a
  `DELETE FROM t1, t2` target list, `DROP INDEX ... ON t`).
  `parseTableReferences` re-parses the fragment with a shrinking token prefix and gives up;
  the signals read the tokens. This is the mechanism removing workarounds, not a behavior
  change any user asked for — but it *is* a behavior change, and the Bytebase PR that
  adopts the signals should show it in the corpus diff.
- **10 — set operations.** Bytebase pools the first arm's tables into the caret's scope across
  `UNION`/`INTERSECT`/`EXCEPT`; omni keeps arms separate, so after `... UNION SELECT | FROM t2`
  omni reports `t2` and Bytebase reports `t1` and `t2`. Both are approximations of the
  engine's answer, which is that a set-operation arm sees only its own `FROM`, and
  `ORDER BY` after a `UNION` sees the output column names rather than any table. Deciding
  which approximation to keep is part of adopting the signals.
- **2 — DML target placement.** `INSERT INTO t1 SELECT | FROM t2` and one `WITH ... INSERT`
  case: Bytebase flattens the write target into `references`, omni keeps it in
  `Scope.DMLTarget`. The adapter unions the two; the field exists so a caller can tell.
- **1 — the reference the signals get right and Bytebase does not.** For
  `ALTER TABLE t1 ADD CONSTRAINT fk FOREIGN KEY (c1) REFERENCES t2(|)` Bytebase keeps `t1`
  (the altered table) in scope; omni names `t2`, which is whose columns the parenthesized
  list takes. The signals are right and Bytebase's adapter will have to follow.
- **4 CTEs.** Two are the boundary itself: `WITH x AS (SELECT * FROM t2)` gives omni no columns
  and Bytebase resolves them with a query span over live metadata. One is a body whose
  select list does not read (`SELECT a.|`), where omni still reports the CTE with no columns
  and Bytebase drops it. One is `RECURSIVE`, where Bytebase drops the columns and omni reads
  the anchor arm's select list.
- **1 object kind.** `DROP INDEX |`: Bytebase offers databases and tables because it cannot
  resolve an index without its table; omni says `index`. That is a product policy Bytebase
  keeps on top of the signal.

## The family

| engine | parser layer | work |
|---|---|---|
| `mariadb` | `complete.go` is byte-identical to MySQL's | copy `complete_context.go`, diff the dialect keyword lists, copy the corpus and add MariaDB's (`INTERSECT`/`EXCEPT`, `MINUS` in Oracle mode, and `RETURNING` on `INSERT`/`DELETE`) |
| `oceanbase` | identical | same copy, plus OceanBase's hints and generated-column forms |
| `tidb` | `complete.go` diverges | the copy, plus the adapter: `backend/plugin/parser/tidb/completion.go` (499 lines) must stop building a catalog and resolve signals instead |
| `mysql` | this PR | done |

Measured cost, so the rest is a bounded follow-up: this PR is 1,252 lines of reader and 492
lines of corpus for the first engine. MariaDB and OceanBase are that reader copied, with the
dialect keyword lists and a corpus per engine; TiDB is the reader plus a 499-line adapter
rewrite, because it is the one MySQL-family engine whose Bytebase adapter resolves metadata
inside omni today. Nothing in the reader is MySQL-specific except the keyword tables and the
dialect's DML forms, which is why the forks are copies rather than designs.

Out of the MySQL family: PostgreSQL, Redshift, SQL Server, and Oracle already expose the
signals and need nothing; what they need is the promotion of the four shapes into
`omni/completion` once a second engine's fields match MySQL's, and that is a rename.
MongoDB is its own cut (no SQL grammar; `db.collection` paths and shell methods), and it
also has the catalog adapter to unwind.

## Verdict

Land it, after review. The cut is ready for MySQL: the signals explain what Bytebase's
adapter computes for itself in 168/200 scopes, 200/200 qualifiers, 200/200 select-alias
positions, and 200/200 prefixes over Bytebase's own corpus, and every disagreement above is
understood and classified rather than unexplained. What is *not* ready is calling the family
done: MariaDB, OceanBase, and TiDB are a copy each, and the TiDB and MongoDB adapters still
resolve metadata inside omni, which is the cut the design document says Bytebase owns. Until
those land, omni has two completion boundaries in production and this one is the newer.
