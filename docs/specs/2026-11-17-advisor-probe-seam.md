# The advisor probe seam: one interface, one capability

2026-11-17 · Proposed

Bytebase `345007290e`, omni `9bda507ad`.

V1's live review rules reach the target database through `advisor.Context.Driver *sql.DB`. After a
cutover there is no `*sql.DB`, so the rules need a seam. Six cutover proofs met it: oracle reused
pg's probe verbatim, mariadb re-ran mysql's shape, tidb measured what its absence costs, and pg,
mysql and oceanbase each built their own — three shapes for one question. The [dispatch
plan](2026-11-17-cutover-dispatch.md#5-what-the-measurements-say-the-plan-got-wrong)
concluded "budget two" without saying which two. This note says which two, from the rules' own
requirements: a **generic statement probe** and **the engine's row estimate**, both in
`backend/plugin/db`, with one `Context.Prober` field replacing `Context.Driver`. The
engine-specific handle every mysql-family proof invented (`DriverHandle any` plus a per-engine
interface asserted in the rule) is not needed and is the shape that makes a missing capability
silent.

The load-bearing finding is narrow, and it is what makes one seam enough:

> If the answer is a statement a rule can write, it is a probe. If the answer is the driver's own
> judgment about the statement, it is a capability.

Twelve of the sixteen rule files that probe ask a question a rule can phrase as SQL. The other
four — the MySQL family's two row-limit rules on each of mysql and mariadb — ask for the driver's
rewrite of the statement into the query its planner can count (`AffectedRowsQuery`), in the plan
format its session variable selects, with a fallback when the plan carries no count. That is engine
judgment, and it is the only thing the seam must ask the driver for.

## 1. What the rules ask, rule by rule

Sixteen files across six engines. The last column is what the seam must offer, in the shape this
note proposes.

| Engine | Rule | File | What it asks | What it calls today | Seam |
| --- | --- | --- | --- | --- | --- |
| `MYSQL` | `STATEMENT_AFFECTED_ROW_LIMIT` | `advisor/mysql/rule_statement_affected_row_limit.go` | how many rows this UPDATE/DELETE changes | `mysqldriver.ExplainJSON` + `mysqldriver.EstimateAffectedRows` over `*sql.DB` | capability (`RowEstimator`) |
| `MYSQL` | `STATEMENT_INSERT_ROW_LIMIT` | `advisor/mysql/rule_insert_row_limit.go` | how many rows this INSERT adds | same | capability (`RowEstimator`) |
| `MYSQL` | `STATEMENT_DML_DRY_RUN` | `advisor/mysql/rule_statement_dml_dry_run.go` | does the server accept this plan | `advisor.Query("EXPLAIN …")` | probe |
| `MARIADB` | `STATEMENT_AFFECTED_ROW_LIMIT` | `advisor/mariadb/rule_statement_affected_row_limit.go` | same as MySQL, MariaDB plan reader | `ExplainJSON` + `EstimateAffectedRowsMariaDB` | capability (`RowEstimator`) |
| `MARIADB` | `STATEMENT_INSERT_ROW_LIMIT` | `advisor/mariadb/rule_insert_row_limit.go` | same | same | capability (`RowEstimator`) |
| `MARIADB` | `STATEMENT_DML_DRY_RUN` | `advisor/mariadb/rule_statement_dml_dry_run.go` | does the server accept this plan | `advisor.Query("EXPLAIN …")` | probe |
| `TIDB` | `STATEMENT_DML_DRY_RUN` | `advisor/tidb/advisor_statement_dml_dry_run.go` | does the plan hold, including `BATCH … DRY RUN`, which refuses a transaction | `advisor.Query` + `queryInAutoCommit` over `*sql.DB` | probe, with `Autocommit` |
| `OCEANBASE` | `STATEMENT_AFFECTED_ROW_LIMIT` | `advisor/oceanbase/advisor_statement_affected_row_limit.go` | what this JSON plan says | `advisor.Query("EXPLAIN format=json …")` + `getEstimatedRowsFromJSON` | probe, read in the rule |
| `OCEANBASE` | `STATEMENT_INSERT_ROW_LIMIT` | `advisor/oceanbase/advisor_insert_row_limit.go` | what this JSON plan says | same | probe, read in the rule |
| `OCEANBASE` | `STATEMENT_DML_DRY_RUN` | `advisor/oceanbase/rule_statement_dml_dry_run.go` | does the server accept this plan | `advisor.Query("EXPLAIN …")` | probe |
| `OCEANBASE` | `STATEMENT_DISALLOW_OFFLINE_DDL` | `advisor/oceanbase/advisor_disallow_offline_ddl.go` | one column's `INFORMATION_SCHEMA.COLUMNS.EXTRA` | `Driver.QueryRowContext(… = ? AND … = ? AND … = ?)` | probe, with `Arguments` |
| `ORACLE` | `STATEMENT_DML_DRY_RUN` | `advisor/oracle/rule_statement_dml_dry_run.go` | does `EXPLAIN PLAN FOR` succeed | `advisor.Query` | probe |
| `POSTGRES` | `STATEMENT_AFFECTED_ROW_LIMIT` | `advisor/pg/advisor_statement_affected_row_limit.go` | what this JSON plan says, under the session's settings | `advisor.Query(EXPLAIN (FORMAT JSON) …, PreExecutions: session settings)` | probe, with `PreExecutions` |
| `POSTGRES` | `STATEMENT_INSERT_ROW_LIMIT` | `advisor/pg/advisor_insert_row_limit.go` | same, counting inserted rows | same | probe, with `PreExecutions` |
| `POSTGRES` | `STATEMENT_DML_DRY_RUN` | `advisor/pg/advisor_statement_dml_dry_run.go` | does the server accept this plan | `advisor.Query` | probe |
| `POSTGRES` | `STATEMENT_OBJECT_OWNER_CHECK` | `advisor/pg/advisor_statement_object_owner_check.go` | the current role | `Driver.QueryRowContext("SELECT CURRENT_USER")` | probe |

Read as questions rather than calls, the sixteen files ask five things:

1. **Does the server accept this statement?** six rules, one statement each (`EXPLAIN`,
   `EXPLAIN PLAN FOR`, TiDB's `EXPLAIN` and its `DRY RUN`).
2. **What does the plan say?** four rules (pg ×2, OceanBase ×2). The reading is a pure function of
   the plan text (`pganalysis.EstimatedAffectedRowsFromExplainJSON`, …), and omni owns it.
3. **How many rows will this statement change?** four rules (mysql, mariadb, ×2 each). The query to
   plan is not the statement: MySQL plans `AffectedRowsQuery(stmt, text)`, under a session
   variable, with a tabular re-plan when the JSON plan carries no count.
4. **What is one column's catalogue extra?** one rule, one parameterized `SELECT`.
5. **Who is the current user?** one rule, one `SELECT`.

Plus one property rather than a question: TiDB's `BATCH … DRY RUN` **cannot run inside a
transaction**, so one probe must run outside one.

Four of the five questions and the transaction property are SQL a rule can write. Only question 3
is engine judgement. That ratio is the whole argument, and it is why the three proofs' handles
disagree: each modeled an engine *method* instead of asking what the rule needs.

## 2. The three shapes the proofs grew

### pg: a generic probe, and it was enough

`prober.go` (new, 46 lines) + `Context.Prober` (+11/−1) + five rule files (+22/−16, +12/−5 of it
the object-owner rule). The probe is:

```go
type Prober interface {
	Probe(ctx context.Context, engine storepb.Engine, qCtx QueryContext, statement string) ([]any, error)
}
```

`ProberOf(driver)` hands back the driver itself when it probes through its own session, and
`SQLProber` wraps a `*sql.DB` for tests. The adapter's `Probe` builds `BEGIN; <settings>;
<statement>; ROLLBACK;` and runs it on one omni session, then reshapes the result into
`[]any{columnNames, columnTypeNames, rows}` — 75 lines, because the pg adapter also replays tenant
mode inside the transaction.

Four rules, four probes, no engine method in sight. The proof's verdict: the generic probe is the
seam, and `plancheck` was *not* the forced interface the prototype predicted.

### mysql: the probe, plus an untyped handle

`Context` gained three fields: `Driver *sql.DB` (kept), `Prober Prober` (from pg, **unused by this
path**), and `DriverHandle any`. Each rule asserts its own interface against the untyped field:

```go
driver, hasDriver := checkCtx.DriverHandle.(mysqlProber)
if !hasDriver && checkCtx.Driver != nil {
	driver, hasDriver = sqlProber{db: checkCtx.Driver}, true
}
...
if hasDriver { … }   // and when it is not: no probes, no advice
```

`mysqlProber` is two methods — `ExplainJSON(ctx, statement) (string, error)` and
`EstimateAffectedRows(ctx, stmt ast.Node, query, plan string) (int64, error)` — declared in the
rule package, naming an omni AST type. `sqlProber` is the legacy `*sql.DB` bridge. Four rule files:
+73/−20, and the package no longer imports `db/mysql` at the call site but still owns the same
question twice.

Two things fall out of that code, and both are costs no line count shows:

- **The `any` field moves the check into sixteen rules.** A signature drift in `mysqlProber` is a
  failed assertion, and a failed assertion is `nil, nil`.
- **The legacy bridge hides the transition.** `hasDriver` is true for as long as the driver still
  holds a `*sql.DB`; the day it stops, the rule goes quiet instead of wrong. The tidb proof
  measured exactly that on its own rule: *"with a `*sql.DB`, **1 EXPLAIN issued, 1 advice**; with
  the adapter's nil, **0 queries, 0 advices**."*

### oceanbase: the handle, alone, with a three-method interface

`Context.DriverHandle any` with no `Prober` at all, and one interface for four rules:

```go
type liveProber interface {
	Explain(ctx context.Context, statement string) error
	ExplainJSON(ctx context.Context, statement string) (string, error)
	ColumnExtra(ctx context.Context, database, table, column string) (string, error)
}
func liveProberFor(checkCtx advisor.Context) liveProber {
	if p, ok := checkCtx.DriverHandle.(liveProber); ok { return p }
	if checkCtx.Driver != nil { return sqlLiveProber{checkCtx.Driver} }
	return nil
}
```

A 57-line helper file, a second `*sql.DB` bridge, six changed files. Two of the three methods are
the generic probe wearing engine names: `Explain` is `Probe("EXPLAIN " + statement)` and
`ExplainJSON` is `Probe("EXPLAIN FORMAT=JSON " + statement)` plus a join. The third is one
parameterized `SELECT`.

### The mariadb proof re-ran mysql's shape

The same four files, 29 of the 73 lines duplicated because both engine values live in `db/mysql`,
and the same silent fallback. That proof's finding is the one this note is built around: with the
rule's probes behind a type assertion against an untyped handle, a driver that cannot answer
produces **no advice at all** — indistinguishable, in the review UI and in a plan check, from a
statement that passed.

## 3. The decision

Two interfaces in `backend/plugin/db`, one `db.Prober` value on `advisor.Context`:

```go
// backend/plugin/db/probe.go
package db

// ProbeContext is how a rule's probe runs.
type ProbeContext struct {
	// PreExecutions run before the statement in the same session, in order.
	PreExecutions []string
	// Arguments bind to the statement's placeholders, in the engine's
	// placeholder syntax.
	Arguments []any
	// Autocommit runs the statement outside a transaction, for an engine whose
	// probe refuses one. PreExecutions is invalid with it.
	Autocommit bool
}

// ProbeResult is one statement's answer.
type ProbeResult struct {
	Columns     []string
	ColumnTypes []string
	Rows        [][]any
}

// Plan is the result read as one JSON plan: every string cell of every row, in
// order, for an engine that splits its plan across rows.
func (r ProbeResult) Plan() (string, error)

// Prober runs a rule's probe statements against the target database.
type Prober interface {
	// Probe runs pCtx.PreExecutions, then statement, in one session whose
	// transaction is rolled back, and returns the statement's columns and rows.
	// With pCtx.Autocommit it runs the statement alone, outside a transaction.
	Probe(ctx context.Context, pCtx ProbeContext, statement string) (ProbeResult, error)
}

// RowEstimator is the capability a rule needs when the estimate is the
// driver's own rewrite of the statement rather than a plan the rule can read.
type RowEstimator interface {
	// CountAffectedRows returns the planner's estimate of the rows statement
	// modifies: exactly one INSERT, UPDATE or DELETE.
	CountAffectedRows(ctx context.Context, statement string) (int64, error)
}

// ProberOf returns the driver's probe, or nil when the driver has none.
func ProberOf(driver Driver) Prober

// NewSQLProber returns a probe over a *sql.DB, for tests and for a caller that
// holds only the SQL handle.
func NewSQLProber(db *sql.DB) Prober
```

```go
// backend/plugin/advisor/advisor.go
type Context struct {
	…
	// Prober runs the live probes the rules need. Nil means the review has no
	// live connection; a rule that needs one reports that it was not checked
	// (see advisor.UnsupportedError) rather than returning no advice.
	Prober db.Prober
	// Driver is gone: every rule reads Prober.
}
```

```go
// backend/plugin/advisor/prober.go
package advisor

// UnsupportedError reports that a rule needed a live capability the connection
// does not offer. A rule returns it instead of no advice: a review that skips a
// rule must not look like one that passed it.
type UnsupportedError struct {
	Engine storepb.Engine
	Need   string // "probe the database", "estimate affected rows"
}
func (e *UnsupportedError) Error() string

// ProbeFor returns the prober a live rule runs through, or an UnsupportedError.
func ProbeFor(checkCtx Context) (db.Prober, error)

// Estimate is a row-limit rule's one call: the engine's estimate, shaped as the
// rule reports it ({Rows} or {CountError}).
func Estimate(ctx context.Context, engine storepb.Engine, p db.Prober, statement string) (v1review.Probe, error)
```

Both implementations of `Probe` are shared, not per engine:

```go
// backend/plugin/db/util/probe.go
// ProbeSQL is advisor.Query moved down: it runs the pre-executions and the
// statement in one transaction that is rolled back.
func ProbeSQL(ctx context.Context, connection *sql.DB, pCtx db.ProbeContext, statement string) (db.ProbeResult, error)

// ProbeSession is the pg proof's adapter method with the pg parts lifted out: it
// opens one omni session, runs BEGIN, the pre-executions and the statement, then
// ROLLBACK.
func ProbeSession(ctx context.Context, querier driver.Querier, pCtx db.ProbeContext, statement string) (db.ProbeResult, error)
```

An engine keeps only what is its own: pg's tenant-mode prelude (the legacy driver already stores
`connectionCtx.TenantMode`; the adapter stores `Options.TenantMode`), and the read of a plan into a
count (`pganalysis`, `oceanbase/analysis`), which is already omni's.

### Why the seam lives in `backend/plugin/db`

**It is Bytebase's seam, and omni needs no change to carry it.** The two things it needs already
exist on omni's side: `Querier.Session`/`Session.Query` (with `QueryOptions.Explain` when a probe
wants a plan) and `CountAffectedRows`. The engine-specific methods the proofs reached for —
`Explain`, `ColumnExtra`, `DryRunAutoCommit` — also exist, and this seam reaches them with a probe
instead of an interface. What the proofs had to name is product vocabulary: `advisor.Context`,
`v1review.Probe`, `storepb.Advice`, and the engine enum. A seam in omni would have to hold those
types or a translation of them — the proofs' own signatures show the cost:
`Probe(ctx, engine storepb.Engine, qCtx, statement)` takes a product enum that two proofs ignore and
one passes as `MYSQL` for a MariaDB statement. So the cutover implements this in Bytebase and opens
no omni PR for it.

The dispatch plan called it `advisor.Prober`, and the field is advisor's. The **type** is not: the
engines that must implement it before their cutover are drivers, and a driver that imports
`backend/plugin/advisor` inverts the dependency (today `db` imports nothing from `advisor`, and
every engine package is below it). `db` is also where the second interface already lives —
`CountAffectedRows` is a method on eight Bytebase drivers and eight omni drivers today — and where
`SQLDB`/`SQLDBProvider` goes in the same wave. `advisor.Context.Prober db.Prober` is the one-line
import, and `advisor` stays the layer that decides what to ask.

### Why the capability is `CountAffectedRows` and not a new method

Every engine that needs an estimate already answers it:

| Side | Implementations |
| --- | --- |
| Bytebase `db.Driver` | mysql (both engine values), oceanbase, tidb, pg, oracle, mssql, cockroachdb, redshift |
| omni `<engine>/driver` | mysql, mariadb, oceanbase, tidb, pg, oracle, mssql, cockroachdb |

`db/mysql.Driver.CountAffectedRows` is the rule's current computation, verbatim: parse one DML,
`AffectedRowsQuery`, `EXPLAIN FORMAT=JSON` with `explain_json_format_version = 1` and error 1193
tolerated, read the plan, fall back to `EXPLAIN FORMAT=TRADITIONAL`, cap by the statement's `LIMIT`
and count the DML targets. `db/oceanbase` and `db/tidb` answer their own plans the same way. The
seam needs no new driver code for any of them; it needs a *name* for the shape.

The same name retires plancheck's eight concrete assertions
(`statement_report_executor.go:168-220`): the type it asserts for is
`interface{ CountAffectedRows(context.Context, string) (int64, error) }`, which is `db.RowEstimator`.
That change still carries the nine blank parser imports of
[dispatch §5](2026-11-17-cutover-dispatch.md#planchecks-assertions-were-hiding-parser-registration),
because the assertions were also the parser registration — but it no longer needs a second seam,
and it removes the pin that keeps `*mysqldriver.Driver` and its six siblings alive.

### Why the other engine-specific needs are probe options

- **`Arguments`** covers OceanBase's `ColumnExtra`. The rule writes the `SELECT` it already runs
  (`WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = ?`) and passes the values; no string
  formatting, no engine SQL in `db`.
- **`Autocommit`** covers TiDB's `BATCH … DRY RUN`. The property is "this statement must not run
  inside a transaction", which any connection can honour; if one does not (a driver that always
  wraps), the server refuses the statement and the rule reports a failed probe for that statement.
- **`PreExecutions`** covers pg's session settings (`SET ROLE`, `SET search_path` collected from the
  script's earlier statements by `liveStatements.probeEach`). Already in `advisor.QueryContext`.
- **`TenantMode` is not in the context.** The pg legacy driver stores `connectionCtx.TenantMode` and
  the adapter stores `Options.TenantMode`, so the owner lookup and `SET LOCAL ROLE` belong to the pg
  prober. That also removes `advisor.Query`'s only use of the engine enum — the parameter every
  rule passes and one of them (mariadb, `storepb.Engine_MYSQL` for a MariaDB statement) passes
  wrongly.

### Rejected alternatives

| Alternative | Why not |
| --- | --- |
| `Context.DriverHandle any` + an interface per rule (mysql, mariadb, oceanbase) | the assertion is the check, and a failed assertion is `nil, nil`; sixteen rules each re-derive it; the legacy bridge keeps it passing until the cutover, then it goes silent (§2) |
| Four capabilities (`Probe`, `RowEstimator`, `ColumnExtraLookup`, `DryRunAutoCommitter`) | three of them are statements; an interface per statement is the shape that produced three variants. Options are mechanical, and their failure is a per-statement probe error, not a skip |
| One interface with all four methods | every engine without a probing rule would implement methods it cannot answer, returning an error to satisfy a compile — a capability advertised as present |
| Keep `Context.Driver` and add `Prober` (the pg proof's tree) | two shapes for the length of the migration; that is exactly how three variants got written, and the rules that keep `Driver` are the ones that then skip |
| `EstimateAffectedRows` taking the omni AST node (the mysql proof's signature) | not contract-neutral: `omni/mysql/ast.Node` and `omni/mariadb/ast.Node` are different types, so the interface cannot be shared, and the rule re-derives what the driver already does |
| Put the seam in omni | nothing in omni changes: `Session.Query`, `CountAffectedRows`, `Explain`, `ColumnExtra` and `DryRunAutoCommit` already exist. The only names the proofs needed are product names |

## 4. What a rule does when the capability is absent

Four states, four answers. The first two are the mariadb finding: neither has an answer, and today
both are `nil, nil` — no advice, indistinguishable from a pass. The last two are the answers a
connection that *can* probe gives, and they do not change.

| State | Today | With this seam |
| --- | --- | --- |
| The review has a connection, the engine cannot answer this question | `return nil, nil` (affected-row), or skip the probes and let omni judge an empty map (insert-row) | `return nil, advisor.UnsupportedError{Engine, Need}` |
| The review has no connection (`Context.Prober == nil`) | `return nil, nil` | same error, `Need: "probe the database"` |
| The connection can answer, the server refuses one statement | `v1review.Probe{Error}` → per-statement advice | unchanged |
| The connection can answer, the budget is spent | `ExplainBudget.AppendSkippedAdvice` warning | unchanged |

`SQLReviewCheck` already turns a rule's error into a visible advice
(`sql_review.go:176-182`):

```text
ERROR  STATEMENT_AFFECTED_ROW_LIMIT  Rule check failed: MYSQL: the live connection cannot estimate
                                     affected rows, so this rule did not run
```

So the rule's whole obligation is one line — `if err != nil { return nil, err }` — and the framework
does the rest. Nothing new is needed in the advice pipeline, and the review stays usable: the other
rules still report, and this one says plainly that it did not run.

Three things keep the state from happening at all:

1. **A compile-time assertion per adapter.** `var _ db.Prober = (*Driver)(nil)` and, where the
   rules need it, `var _ db.RowEstimator = (*Driver)(nil)`. A cutover that forgets `Probe` does not
   compile.
2. **One contract test per engine, beside the other driver tests.** It builds the engine's driver
   against its pinned server and asserts the capabilities its rules need (`Probe` for the six
   probe engines; `RowEstimator` for mysql, mariadb, oceanbase, tidb, pg, oracle, mssql,
   cockroachdb and redshift through plancheck). The engine list is the table in §5; the test is the
   place it is checked rather than asserted.
3. **The static yaml harness stops pretending.** `RunSQLReviewRuleTest` builds
   `Context{Driver: nil}` (`utils_for_tests.go:414`), so today's fixtures for the four live rule
   types pass because nothing ran: `mysql/test/statement_affected_row_limit.yaml` holds two
   statements and expects no advice. Those fixtures are the silent failure in miniature. With the
   seam they assert the loud advice instead — one case per engine, which is the absence-path test —
   and the rules' real coverage is the live tests that open a connection
   (`mysql_rules_test.go:231`'s `checkTestRowLimitRule`, `pg/advisor_*_test.go`,
   `oceanbase/advisor_*_test.go`, `tidb/advisor_statement_dml_dry_run_test.go`).

## 5. What a driver implements per engine

`Probe` is new on both sides and is the same ~10 lines after `util.ProbeSession`; a capability is
already there or is not needed.

| Engine | `Probe` (legacy driver → adapter) | `RowEstimator` | Engine-specific part of the probe | Rules |
| --- | --- | --- | --- | --- |
| `POSTGRES` | `util.ProbeSQL` → `util.ProbeSession` | yes (unused by rules) | tenant-mode prelude inside the transaction | 4 |
| `COCKROACHDB` | same as pg | yes | none | 0 (plancheck) |
| `MYSQL` | `util.ProbeSQL` → `util.ProbeSession` | yes, both sides | none | 3 |
| `MARIADB` | the `db/mysql` package, branch on `dbType` | yes (MariaDB parser) | none | 3 |
| `TIDB` | `util.ProbeSQL` with `Autocommit` → `util.ProbeSession` | yes | `BATCH … DRY RUN` runs outside a transaction | 1 |
| `OCEANBASE` | `util.ProbeSQL` → `util.ProbeSession` | yes (unused by rules) | none | 4 |
| `ORACLE` | `util.ProbeSQL` → `util.ProbeSession` | yes (plancheck) | none | 1 |
| `MSSQL` | not required (its advisors are static AST rules) | yes (plancheck) | — | 0 |
| cassandra, redis, spanner, elasticsearch, trino, clickhouse, hive, doris, starrocks, mongo, cosmosdb, dynamodb, snowflake, bigquery, databricks | — | no | — | 0 |

`ProberOf(driver)` returns nil for the last row, and that is not a defect: no rule of those engines
probes. It becomes one only if an engine gains a probing rule without gaining a prober — which is
what §4's compile-time assertion and contract test prevent.

For the four engines whose rules read a plan rather than an estimate (pg ×2, OceanBase ×2), the
count stays in the rule: `res, err := prober.Probe(...)` then
`oceanbaseparser.EstimateAffectedRowsFromExplainJSON(plan)`. Reading a plan is omni's pure function,
not the driver's judgement, and moving it into the driver would be the same mistake in reverse.

## 6. Migration order

The seam is one change and lands in **wave 0**, not alongside pg's cutover, because every rule that
reads `Context.Driver` must move at once or the tree carries two shapes and the rules left on
`Driver` are the silent ones. No engine cuts over in this change.

**Wave 0 — the seam itself (no engine cut over, one change).**

1. `backend/plugin/db/probe.go`: `ProbeContext`, `ProbeResult`, `Prober`, `RowEstimator`,
   `ProberOf`, `NewSQLProber`.
2. `backend/plugin/db/util/probe.go`: `ProbeSQL` (today's `advisor.Query` moved verbatim, including
   the tenant-mode branch, which moves to the pg driver's `Probe`) and `ProbeSession` (the pg
   proof's adapter method, with the tenant-mode prelude as the caller's callback).
3. `Probe` on the six engine values in five legacy packages: pg (+tenant prelude, ~15 lines),
   mysql/mariadb (~5), tidb (+`Autocommit`, ~10), oceanbase (~5), oracle (~5).
4. `advisor.Context`: `Driver *sql.DB` → `Prober db.Prober` (+11/−1); `advisor/prober.go`
   (`UnsupportedError`, `ProbeFor`, `Estimate`); `advisor.Query` and `QueryContext` deleted.
5. The sixteen rule files and their tests; the static yaml fixtures for the four live rule types
   become absence-path cases (§4).
6. `RowEstimator` is defined here and first consumed in **wave 2** by plancheck
   (cockroachdb's change), exactly as the dispatch orders it — the definition is cheap, the
   assertion table's nine parser imports are what waits.

Estimated: +290/−190 across ~24 files. The proofs measured the shapes this replaces, and none of
them composes: `prober.go` 46 lines, `Context` +11/−1, pg's five rule files +22/−16, mysql's four
+73/−20 (29 of them paid twice by mariadb), OceanBase six advisor files ~108 changed lines with a
57-line helper. Three shapes, about 300 changed lines, and the rules that read the untyped handle
still skip when it is absent.

**Wave 3 — pg is the first cutover that lands the seam's omni side.** `util.ProbeSession` replaces
the 75-line inline method; the four pg rules do not change again, because wave 0 already moved
them. `oracle` follows and adds only its `Probe`. `tidb` and `oceanbase` then add `Probe` to their
adapters (TiDB keeps `Autocommit`, OceanBase's offline-DDL rule keeps `Arguments`). `mysql` and
`mariadb` last: their adapters already implement `CountAffectedRows`, so the estimator needs no new
code, and the cutover deletes the legacy `Probe` methods with the packages.

The order inside wave 3 does not change: pg (measured, 658 adapter lines), oracle (390, reuse),
tidb/oceanbase, mysql/mariadb (706, worst). What changes is that the seam is no longer a wave-3
item, so no engine's cutover carries a product interface.

## 7. Costs and deltas this change owns

| Item | Cost | Notes |
| --- | --- | --- |
| `db/probe.go` | ~85 lines | types, two constructors, `Plan` |
| `db/util/probe.go` | ~110 lines | `ProbeSQL` is `advisor.Query` moved (−45 there); `ProbeSession` is the pg proof's method, shared |
| legacy `Probe` × 6 | ~40 lines | deleted with each package at cutover |
| `advisor.Context` + `prober.go` | +11/−1, +55 | `RowEstimator` needs no driver change |
| 16 rule files | ~+60/−80 | mysql/mariadb's four shrink by ~20 each: the estimate is one call, and `db/mysql` leaves the advisor package |
| `advisor.Query`/`QueryContext` | −45 | no caller left |
| static fixtures | 4 yaml files → absence-path cases | they asserted silence; the live tests are the coverage |
| plancheck (wave 2) | +24/−62 plus the nine imports | unchanged from the dispatch's estimate; now one interface named in `db` |

Four behaviour deltas, each small and each named:

- **A row-limit estimate that fails is now `CountError`.** Today an `EXPLAIN` refusal is
  `Probe{Error}` ("dry runs failed: …") and an unreadable plan is `Probe{CountError}`
  (`code.Internal`). `CountAffectedRows` returns one error for both, so every failure becomes
  `CountError`: same advice level, same statement and server text, different code. When
  `STATEMENT_DML_DRY_RUN` is configured it reports the refusal separately; a rule that must keep the
  distinction can probe `EXPLAIN` on the failure path only.
- **`advisor.Query`'s `[]any` becomes `db.ProbeResult`.** `pg/utils.go`'s `getExplainPlan` (three
  type assertions) and `oceanbase/common.go`'s `getEstimatedRowsFromJSON` (three) collapse into
  `ProbeResult.Plan()`; their error text ("expected []any for row but got …") disappears.
- **OceanBase's offline-DDL lookup uses the rule's `ctx`.** Today it passes
  `context.Background()` (`advisor_disallow_offline_ddl.go:55`), so a cancelled review still queries.
- **MariaDB's dry run stops naming MySQL.** It passes `storepb.Engine_MYSQL` to
  `advisor.Query` today; the probe has no engine parameter.

## 8. What would falsify this

- **A rule that needs the driver to judge something the rule cannot phrase as SQL, other than the
  estimate.** The MySQL-family estimate is the one such case in sixteen files; a second would mean
  a second capability, and the seam grows a name rather than a shape. The test is whether the
  answer depends on a plan format, a session variable, or a rewrite of the statement — if it does
  not, it is a probe.
- **A driver whose probe must run more than one statement to answer one question.** Then
  `ProbeContext.PreExecutions` is setup, not the answer, and the answer's statement is one. If a
  real case appears, the probe grows a `Script` variant; nothing in this note assumes it.
- **A capability that cannot answer over `*sql.DB`.** The seam's migration value rests on
  `ProbeSQL` being able to serve every rule that does not need `RowEstimator`. If a future
  capability cannot, that engine's rules must move with its cutover, which is a wave-order cost, not
  a shape change.

## 9. References

- [Cutover dispatch](2026-11-17-cutover-dispatch.md) — §5 for the two-seam budget this replaces,
  §7 for the waves, §9 for what stays product-side.
- [Cutover open decisions](2026-11-17-cutover-open-decisions.md) — §1 (`QueryConn`'s `*sql.Conn`)
  and §2 (`db.SQLDB`) are the other two product seams; §5 is plancheck's assertion table.
- [Driver method parity](2026-11-17-driver-method-parity.md#getdb-replacements) — the list this
  note draws on: `Explain`, `ExplainJSON`, `EstimateAffectedRows`, `ColumnExtra`,
  `DryRunAutoCommit`, `CurrentUser` all exist on omni's side.
- Proof reports: pg (`8664fa7b1`, six majors), mysql (`6dee22d57`, two pinned images), mariadb
  (`ad3d0f433`), oceanbase, tidb (`v8.5.5`) — each in its session's report.
