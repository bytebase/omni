# MySQL Routine-Body Grammar-Driven Parse (Design Doc, rev 3)

Date: 2026-04-20
Author: junyi
Status: draft — rev 4, post second codex review, pre-implementation
Changelog:
- rev 2: dropped Split revert (Codex P0), reordered commits to green-at-every-
  step (P1), refocused audit on compound-side gaps (P1), added
  `catalog/table.go` to cascade (P2), tightened MySQL-parity wording (P2).
- rev 3: scope locked to **B2 = Layer 1a (grammar) + label matching +
  DECLARE ordering**. Items 4–8 deferred to `mysql/semantic` PR.
- rev 4: (Codex P0) corrected DECLARE phase tracker — handler may follow
  var/condition directly without requiring a cursor. (Codex P1) corrected
  label matching to **case-insensitive** after MySQL 8.0 container
  verification. (Codex P1) spelled out `DO`-token handling for the two
  event sites in commit 3. (Codex P1) named concrete corpus seeds.
  (Codex P2) acknowledged existing fixtures don't need repair.

## Problem

`parseCreate{Function,Trigger,Event}Stmt` and `parseAlterEventStmt` consume
the routine body via a hand-rolled raw-token depth scanner that only tracks
`BEGIN`/`END` pairs, then stores the body as an opaque `string`. All compound
constructs inside the body (`IF ... END IF`, `CASE`, `WHILE`, `LOOP`, `REPEAT`,
nested blocks, cursors, handlers, labels) have to be "balanced" by text-level
heuristics (prev-token-is-`END`, next-char-is-`(`, next-word-is-`EXISTS`,
etc.).

Every time a real-world procedure pattern hits a corner the heuristics don't
cover, body termination is wrong and the trailing tokens leak out to be
re-parsed as a new statement — `syntax error at or near "IF"` / `"WHILE"`.
Fix #97 handled one class (inner `END`); the `if(x) then`, `WHILE(cond) DO`,
`IF EXISTS (subquery) THEN` flow-control, and comment-between-keywords
patterns all remain.

MySQL's own parser is yacc over the grammar: `IF_SYM` at `sp_proc_stmt`
position is compound IF; `IF_SYM` at `expr` position is the ternary function.
Container-verified on MySQL 8.0: whitespace is irrelevant.

omni's recursive-descent parser already has the grammar machinery —
`parseCompoundStmtOrStmt` dispatches `kwIF` → `parseIfStmt`; `parseStmt` →
`parseExpr` handles `kwIF` in expression position as a function. The only
reason the grammar is not used for routine bodies is the historical shortcut
of the raw scanner.

## Goal

Replace routine/trigger/event body scanning with grammar-driven parse. Delete
the text-level heuristics that only existed to prop up the scanner. Enforce
two MySQL-level static constraints that `parseCompoundStmtOrStmt` today
leaves lax: **label matching** and **DECLARE ordering**. Make body-inside
content available as an AST so follow-up semantic work has a real tree.

## Scope (B2)

**In scope** — category numbers refer to the full MySQL CREATE-time static
validation list:

| # | Constraint | Implementation |
|---|---|---|
| 1 | Grammar | switch 4 body sites to `parseCompoundStmtOrStmt`; delete raw scanners |
| 2 | Label matching (start label ≡ end label; end label only if start label present) | validate in `parseBeginEndBlock`, `parseWhileStmt`, `parseLoopStmt`, `parseRepeatStmt` |
| 3 | DECLARE ordering inside `BEGIN ... END` (var/condition → cursor → handler → stmts) | phase tracker in `parseBeginEndBlock` |

**Out of scope — explicit follow-up PR** (`mysql/semantic: static validation
for stored routine bodies`):

| # | Constraint | Reason to defer |
|---|---|---|
| 4 | Symbol resolution: variable / cursor refs must be declared | needs scope chain / symbol table — an independent analyzer layer |
| 5 | Label resolution for `LEAVE` / `ITERATE` | same infra |
| 6 | Duplicate declaration detection (var, cursor, handler, label within scope) | depends on #4 infra |
| 7 | Handler condition uniqueness within one `DECLARE ... HANDLER FOR` | small, but naturally paired with #6 |
| 8 | Function RETURN coverage (all paths end with `RETURN`) | requires a small CFG/flow pass |

**Out of scope — not planned** in either PR:

- `sql_mode`-sensitive parsing (ANSI_QUOTES, NO_BACKSLASH_ESCAPES,
  PIPES_AS_CONCAT, IGNORE_SPACE, HIGH_NOT_PRECEDENCE). Separate bigger
  project; requires lexer restructure.
- Privilege checks. omni is not an execution engine.
- Schema reference validation (table/column existence). MySQL itself defers
  this to runtime; omni does the same.

## Architectural change

### AST shape

Four statement types carry a body: `CreateFunctionStmt` (covers both FUNCTION
and PROCEDURE via `IsProcedure`), `CreateTriggerStmt`, `CreateEventStmt`,
`AlterEventStmt`. All four gain:

```go
type CreateFunctionStmt struct {
    ...
    Body     ast.Node  // parsed sp_proc_stmt (compound or simple statement)
    BodyText string    // raw source bytes of body, captured via Loc range
}
```

`Body ast.Node` is the authoritative representation; `BodyText` is the raw
source bytes preserved verbatim from the input segment. The existing
`Body string` field is **renamed to `BodyText`**; the new `Body` field is an
`ast.Node`. Renaming (rather than keeping `Body string` + adding `BodyAST`)
forces every existing call site to make an explicit choice between "I want
text" and "I want AST", and makes future grep audits work correctly.

Per MySQL grammar a routine body is exactly one `sp_proc_stmt`, so `Body` is
a single `ast.Node`. For `BEGIN...END` bodies this is a `*ast.BeginEndBlock`
whose `Stmts` holds the list; for single-statement bodies (`RETURN expr`,
`IF ... END IF` without BEGIN, `CALL foo()`, etc.) it is the corresponding
node directly.

### Parser flow

The grammar-parse block is three lines plus site-specific prelude:

```go
bodyStart := p.pos()
body, err := p.parseCompoundStmtOrStmt()
if err != nil {
    return nil, err
}
bodyEnd := p.pos()
stmt.Body = body
stmt.BodyText = p.inputText(bodyStart, bodyEnd)
```

Call sites (four, not five — `CREATE FUNCTION` and `CREATE PROCEDURE` share
`parseCreateFunctionStmt`):

1. `parseCreateFunctionStmt` (FUNCTION + PROCEDURE) — call after parsing
   characteristics. No DO token. Body is always present.
2. `parseCreateTriggerStmt` — call after parsing `FOR EACH ROW` and the
   optional `FOLLOWS/PRECEDES` clause. No DO token. Body is always present.
3. `parseCreateEventStmt` — body is introduced by the mandatory `DO`
   keyword. The DO consumption stays in place; grammar-parse block runs
   after `if p.cur.Type == kwDO { p.advance() }`. Body is always present
   once DO is consumed.
4. `parseAlterEventStmt` — body is introduced by the **optional** `DO`
   keyword. Grammar-parse block is gated by
   `if p.cur.Type == kwDO { p.advance(); ...grammar-parse block... }`. If
   DO is absent the alter statement has no body (and the surrounding
   characteristic-only form is already handled). This was the site Fix #97
   missed.

### Static constraint enforcement (new)

**Label matching** — applied in `parseBeginEndBlock`, `parseWhileStmt`,
`parseLoopStmt`, `parseRepeatStmt`. Verified against MySQL 8.0 oracle
(2026-04-20):

- Match is **case-insensitive** (`myLabel` matches `MYLABEL`, confirmed in
  container).
- If start label is present and end label is present, they must match
  (case-insensitive equality).
- If start label is absent and end label is present, reject — MySQL errors
  with `End-label X without match` (ERR 1310). omni's message:
  `end label "X" without matching begin label`.
- Both absent, or start-only, are both valid.
- Same rule applies to `WHILE`, `LOOP`, `REPEAT`.

**DECLARE ordering** — applied in `parseBeginEndBlock`. State machine
(verified against MySQL 8.0 oracle):

```
states: init → saw_var → saw_cursor → saw_handler → saw_stmt
allowed moves per incoming construct:
  DECLARE var | condition   : state ∈ {init, saw_var}                         → saw_var
  DECLARE cursor            : state ∈ {init, saw_var, saw_cursor}             → saw_cursor
  DECLARE handler           : state ∈ {init, saw_var, saw_cursor, saw_handler} → saw_handler
  regular statement         : any state                                       → saw_stmt
  any DECLARE when state is saw_stmt                                          → reject
```

**Key:** a `DECLARE HANDLER` may follow `DECLARE VAR`/`DECLARE CONDITION`
directly without any cursor in between; this is legal MySQL (container-
verified). Only cursors must not follow handlers, and variables/conditions
must not follow either cursors or handlers.

Errors raised (mirroring MySQL errors 1337/1338 in meaning; omni uses its
own style text):

- `variable or condition declaration after cursor or handler declaration`
  (MySQL ERR 1337)
- `cursor declaration after handler declaration` (MySQL ERR 1338)
- `DECLARE after regular statement` — MySQL returns ERR 1064 here (generic
  syntax error); omni gives the specific form since the parser knows the
  context.

Phase is per-block: each nested `BEGIN...END` has its own state machine.

### Things deleted

- `mysql/parser/parser.go::consumeRoutineBody`
- `mysql/parser/split.go::findCompoundBodyEnd`
- The raw `bodyStart / depth / for loop / kwBEGIN++ / kwEND--` block in each
  of the four call sites

### Split

**Not modified.** Split's compound-depth tracking (including the
`IF EXISTS (subquery)` heuristic added in PR #97) is required for the
default-`;` delimiter case, where Split must avoid splitting at intra-body
`;`. Since Split runs before the parser, grammar-driven body parse cannot
rescue a wrongly split segment — Split must hand the parser a complete
procedure body. PR #97's heuristic stays.

## Behavior matrix

Container-verified against MySQL 8.0 (2026-04-20). All rows must pass with
new implementation.

**Grammar (Layer 1 / category 1):**

| Pattern | MySQL | New omni |
|---|---|---|
| `if(x) then ... end if` (no space, stmt ctx) | compound IF | `parseIfStmt` |
| `IF (x) THEN ... END IF` (space, stmt ctx) | compound IF | `parseIfStmt` |
| `SET @r = IF(a,1,0)` (no space, expr ctx) | function | `parseExpr` → function |
| `SET @r = IF (a,1,0)` (space, expr ctx) | function | `parseExpr` → function |
| `IF EXISTS (SELECT 1) THEN ... END IF` | compound IF w/ EXISTS predicate | `parseIfStmt` |
| `DROP TABLE IF EXISTS t` inside body | DDL modifier | `parseDropStmt` |
| `CREATE TABLE IF NOT EXISTS t(...)` inside body | DDL modifier | `parseCreateTableStmt` |
| `BEGIN /* c */ IF x THEN ... END IF; END` | compound IF | tokenizer skips comment, `parseIfStmt` runs |
| `DECLARE EXIT HANDLER FOR cond IF ... END IF` | single-stmt handler body | `parseDeclareHandlerStmt` → `parseCompoundStmtOrStmt` |
| `SET r = CASE x WHEN 1 THEN 'a' END` (expr CASE) | CASE expression | `parseExpr` |
| `CASE x WHEN 1 THEN SET @y=1; END CASE` (stmt CASE) | compound CASE | `parseCaseStmt` |
| nested `BEGIN ... IF ... WHILE ... END WHILE; END IF; END` | nested | grammar recursion |
| `ALTER EVENT e DO BEGIN IF x THEN ... END IF; END` | compound IF | `parseAlterEventStmt` |
| syntax error in body (`IF x THEN_MISSING`) | rejected at CREATE | rejected at `parseCompoundStmtOrStmt` |
| unknown table/column ref in body | accepted (deferred to runtime) | accepted |

**Label matching (category 2):** Case-insensitive, MySQL 8.0 container-
verified.

| Pattern | MySQL | New omni |
|---|---|---|
| `lbl: BEGIN ... END lbl` | ✓ | ✓ |
| `lbl: BEGIN ... END` (end label omitted) | ✓ | ✓ |
| `BEGIN ... END` (both absent) | ✓ | ✓ |
| `myLabel: BEGIN ... END MYLABEL` (different case) | ✓ | ✓ (case-insensitive match) |
| `BEGIN ... END lbl` (end without start) | ✗ (ERR 1310 `End-label lbl without match`) | ✗ `end label "lbl" without matching begin label` |
| `foo: BEGIN ... END bar` (mismatch) | ✗ (ERR 1310) | ✗ `end label "bar" does not match begin label "foo"` |
| `lbl: WHILE x DO ... END WHILE lbl` | ✓ | ✓ |
| `lbl1: WHILE x DO ... END WHILE lbl2` | ✗ | ✗ same error |
| same for `LOOP`, `REPEAT` | as above | as above |

**DECLARE ordering (category 3):** State machine per revised design above.
MySQL 8.0 container-verified.

| Pattern | MySQL | New omni |
|---|---|---|
| `BEGIN DECLARE x INT; SELECT 1; END` | ✓ | ✓ |
| `BEGIN DECLARE x INT; DECLARE cur CURSOR FOR SELECT 1; DECLARE CONTINUE HANDLER FOR NOT FOUND SET @a=1; SELECT 1; END` | ✓ | ✓ |
| `BEGIN DECLARE x INT; DECLARE CONTINUE HANDLER FOR NOT FOUND SET @a=1; SELECT 1; END` (var → handler, skip cursor) | ✓ (verified) | ✓ |
| `BEGIN DECLARE dup_key CONDITION FOR SQLSTATE '23000'; DECLARE EXIT HANDLER FOR dup_key SET @err=1; SELECT 1; END` (condition → handler, skip cursor) | ✓ (verified) | ✓ |
| `BEGIN DECLARE cur CURSOR FOR SELECT 1; DECLARE x INT; END` (var after cursor) | ✗ (ERR 1337) | ✗ `variable or condition declaration after cursor or handler declaration` |
| `BEGIN DECLARE CONTINUE HANDLER FOR NOT FOUND SET @a=1; DECLARE cur CURSOR FOR SELECT 1; END` (cursor after handler) | ✗ (ERR 1338 `Cursor declaration after handler declaration`) | ✗ `cursor declaration after handler declaration` |
| `BEGIN SELECT 1; DECLARE x INT; END` (DECLARE after stmt) | ✗ (ERR 1064 generic) | ✗ `DECLARE after regular statement` |
| nested block resets phase: `BEGIN DECLARE x INT; BEGIN DECLARE y INT; END; END` | ✓ (inner scope) | ✓ (inner `parseBeginEndBlock` has its own state machine) |

## Audit

Audit is **done inside this PR** as realworld-corpus test code. It is not a
separate PR and not a standalone investigation; the tests land in this PR
and must pass by the last commit.

Audit deliverables in this PR:

1. **Realworld corpus test** (`mysql/parser/routine_body_realworld_test.go`).
   Seeds (all easy to obtain — public sample databases + the existing
   in-repo samples):
   - **Existing in-repo seeds**: all procedures/functions/triggers/events
     in `mysql/parser/split_realworld_test.go` (`citycount`, `CalcIncome`,
     `trg_audit`).
   - **MySQL official sakila sample DB** (https://dev.mysql.com/doc/sakila/):
     procedures `film_in_stock`, `film_not_in_stock`, `rewards_report`;
     function `inventory_held_by_customer`, `inventory_in_stock`,
     `get_customer_balance`; triggers `customer_create_date`,
     `payment_date`, `ins_film`, `upd_film`, `del_film`, `rental_date`.
   - **MySQL official employees sample DB**: no stored programs in the
     base schema, skip.
   - **Canonical flow-control patterns**: synthetic procedures covering
     every row in the behavior matrix (`if(x) then`, `IF EXISTS (subq)`,
     labeled loops with LEAVE/ITERATE, nested CASE, cursor+handler, etc.).
   - **User-reported pattern**: the `if(quiz) then` / `WHILE(cond) DO`
     shape (anonymized version of the `totem-dev__*` failures) — one
     synthetic procedure covering the general shape, not the specific
     file.

   Each case asserts: `Parse()` succeeds, `stmt.Body != nil`, and the
   expected top-level node type.

2. **Round-trip test**: `stmt.BodyText` fed back through
   `parseCompoundStmtOrStmt` must produce a structurally equivalent AST.
3. **Spot-check for dynamic SQL**: one procedure covering `PREPARE` /
   `EXECUTE` / `DEALLOCATE` / `GET DIAGNOSTICS` / `SIGNAL` to confirm
   `parseStmt`'s existing dispatch works in a compound-body context.
4. **Coverage matrix tests** (from the tables above): one positive + one
   negative case per row where applicable.

Any parser gap surfaced by these tests is fixed in a preceding commit of
this PR. If a gap is too large to fix in-PR (e.g., an entire new statement
kind missing from `parseStmt`), the specific case is gated with `t.Skip`
plus a tracked follow-up issue, and the audit is considered complete for
the rest of the corpus.

## Cascade: downstream call sites

`grep -rn "\.Body\b" mysql/` finds:

| File | Usage | Change |
|---|---|---|
| `mysql/ast/parsenodes.go` (4 structs) | AST field decl | rename `Body` → `BodyText`; add `Body ast.Node` |
| `mysql/ast/outfuncs.go` (4 sites) | prints `:body %q` | rename `n.Body` → `n.BodyText`; follow-up PR adds `:body-stmt <subtree>` |
| `mysql/catalog/routinecmds.go` (2 sites) | populates `Routine.Body` string | `stmt.BodyText` |
| `mysql/catalog/eventcmds.go` (3 sites) | populates `Event.Body` string | `stmt.BodyText` |
| `mysql/catalog/triggercmds.go` (2 sites) | populates `Trigger.Body` string | `stmt.BodyText` |
| `mysql/catalog/table.go` (~line 105) | `Trigger.Body` storage struct field | unchanged (catalog-side field; populator at `triggercmds.go:56` switches to `stmt.BodyText`) |
| `mysql/parser/compare_test.go` (3 sites) | test helper reads `stmt.Body` | `stmt.BodyText` |
| `mysql/parser/alter_misc.go` (1 site) | assigns `stmt.Body = p.inputText(...)` inside old scanner | replaced by grammar-parse block |

Catalog code stays string-based in this PR. The `scenarios_bug_queue/c11.md`
trigger field validator follow-up will switch the catalog layer to consume
`stmt.Body ast.Node`, in combination with the `mysql/semantic` PR that
provides #4–#8 validation.

## Commit sequence

Every commit compiles and the test suite passes.

1. `feat(mysql/parser): realworld routine-body corpus`. Adds the new test
   file; existing scanner still in use; tests pass against current behavior
   for the subset the scanner happens to handle; documented failures for
   the subset that doesn't are gated with `t.Skip` carrying TODO markers
   referencing this plan.

2. `refactor(mysql/ast): rename Body→BodyText, add Body ast.Node`.
   Mechanical rename of the `Body string` field in all four AST structs to
   `BodyText`. Add a new `Body ast.Node` field (defaults to nil). Update
   every read site in `mysql/ast/outfuncs.go`, `mysql/catalog/{routine,
   event,trigger}cmds.go`, `mysql/parser/{compare_test,alter_misc,trigger,
   create_function}.go`. Parser still populates `BodyText` via the scanner;
   `Body` stays nil. Behavior unchanged, tests green.

3. `feat(mysql/parser): parse routine body via grammar`. Substantive change:
   - Delete `consumeRoutineBody`, `findCompoundBodyEnd`, all four raw
     scanners.
   - Insert the three-line grammar-parse pattern in each of the four
     sites.
   - `stmt.Body` is now populated; `stmt.BodyText` is captured via Loc
     range.
   - Remove `t.Skip` markers added in commit 1; corpus tests now all pass.
   - Unit tests assert `stmt.Body` is the expected node type per matrix
     row.
   - Round-trip test passes.

4. `feat(mysql/parser): enforce label matching in compound statements`.
   Adds the start-label ≡ end-label check to `parseBeginEndBlock`,
   `parseWhileStmt`, `parseLoopStmt`, `parseRepeatStmt`. Adds positive +
   negative test cases per matrix.

5. `feat(mysql/parser): enforce DECLARE ordering in BEGIN...END`. Adds
   phase tracker to `parseBeginEndBlock`. Adds positive + negative test
   cases per matrix.

Commits 4 and 5 are independent and can be reordered or split further.

## Risks

- **Audit misses a gap.** Realworld corpus is finite. If user scripts hit a
  compound construct the parser doesn't handle, CREATE PROCEDURE fails
  where it used to blindly "succeed" (storing garbage as text). This is a
  correctness improvement, but may generate noisy bug reports. Mitigated
  by extensive corpus + cheap per-gap fixes.
- **AST rename cascades outside `mysql/`.** Grep confirmed all consumers
  live inside `mysql/*`. Any non-mysql consumer will break at compile;
  fixed in commit 2.
- **`BodyText` round-trip not byte-identical.** Loc-range capture gives
  exact source bytes. Round-trip test enforces.
- **New static-constraint errors break existing tests.** Codex review (rev
  4) searched existing fixtures: all labeled fixtures already match
  (`myblock`/`myblock`, `lbl`/`lbl`, `label1`/`label1`,
  `outer_loop`/`outer_loop` in `compare_test.go` / `split_realworld_test.go`
  / wt fixtures), and all DECLARE-order fixtures are valid. Commits 4 and 5
  are therefore expected to be **additions only**, not fixture repairs.
  New negative fixtures added per matrix row.

## Follow-ups

- `mysql/semantic: static validation for stored routine bodies` — items
  #4–#8 (symbol resolution, duplicate detection, handler-condition
  uniqueness, function RETURN coverage). Takes the `ast.Node` `Body` field
  as input.
- `c11.md` trigger field validator — uses `Body ast.Node` + the
  `mysql/semantic` work.
- `sql_mode`-sensitive parsing (ANSI_QUOTES, etc.). Independent project.
- AST-driven deparse for procedures (replaces `BodyText`-based
  reconstruction).
- Completion inside routine bodies (consumes `Body ast.Node`).
