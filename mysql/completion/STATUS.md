# Completion cut: status

2026-11-17 · Resume point. This branch's PR replaces #1291, which was parked after Codex's review
quota ran out; the code is identical.

The completion cut is "omni emits parser-native token/rule/scope/intent signals; Bytebase resolves
live metadata into product `Candidate` values" (`docs/specs/2026-11-16-completion-cut-design.md`).
The mechanism is implemented for MySQL, MariaDB, OceanBase, and TiDB, and every disagreement with
Bytebase's own completion corpus is classified. What is left is review capacity, not design.

## Where the PRs stand

| PR | engine | head at parking | state |
|---|---|---|---|
| replacement for #1291 | MySQL | this branch | the active PR; #1291 closed as superseded, 0 open review threads (79 resolved) |
| #1299 | MariaDB | `bcf7a157` | draft, CI green, 25–27 open threads |
| #1300 | OceanBase | `98a69482` | draft, CI green, 25–27 open threads |
| #1302 | TiDB | `6e3f51f6` | draft, CI green, 25–27 open threads |
| #1303 | Mongo (design only) | merged | describes what Mongo needs; no code |

`origin/main` at parking is `3d5e08b5`. None of these branches contains it; each needs
`origin/main` merged (never rebased) before merging.

The merge gate is CI green **and** Codex reporting no major issues on the current head. The last
Codex request on #1291 returned a usage-limit `failed` status, not a review; a failed review is not
a finding and must not be read as either a pass or a fail. Do not retry it until the limit resets.

## The cut (signals only)

`<engine>/parser.CollectCompletion(sql, cursorOffset) *CompletionContext` carries:

- `Candidates`: the grammar's token and rule accept-set at the cursor.
- `Scope.LocalReferences` / `OuterReferences` / `References`: the cursor query block's relations,
  the correlated outer blocks, both; `Scope.DMLTarget` for a write target.
- `RangeReference`: `Kind`, `Database`, `Table`, `Alias`, `Columns` (only when the text names
  them), `AliasColumns`, `Loc`, `BodyLoc`, `Lateral`. `BodyLoc.End == -1` means an unterminated
  body that runs to the end of the text.
- `CTEs`: the definitions in scope, with their columns.
- `SelectAliases`: the cursor block's select-list aliases.
- `Intent.ObjectKinds` and `Intent.Qualifier`.
- `Prefix`, read from the lexer token.
- `StatementLoc`.

Nothing else crosses: no `catalog.Catalog`, no metadata proto, no `storepb` type, no product type.
The reader is a token walk, linear in the statement's tokens (depths and matching parentheses are
computed once; group recursion stops at `maxNestingDepth`).

## The corpus and every disagreement

`mysql/parser/complete_context_test.go` holds 144 cases over `sql` × cursor positions, plus a
sweep of every cursor offset in twelve scripts, a nesting-limit test, and a test that a statement
nested in `WITH` bodies is read once per clause. Each fork's corpus is that corpus with its
dialect cases (133 for MariaDB, 134 for OceanBase, 133 for TiDB at parking, several parent fixes
behind), and the engine-rejected cases become comments: MariaDB and TiDB lose the `TABLE t` query
primary and the two INSERT row-alias cases, OceanBase only `TABLE t`.

Parity against Bytebase's own corpus (`backend/plugin/parser/<engine>/test-data/test_completion.yaml`,
200 cases, candidate lists excluded because the metadata fixture would prove the fixture and not the
metadata; what is compared is the input to resolution):

| signal | MySQL / MariaDB / OceanBase | TiDB |
|---|---|---|
| relation scope | 167/200 | 167/200 |
| object kinds | 197/198 | 197/198 |
| qualifier | 200/200 | 200/200 |
| select aliases | 199/200 | 199/200 |
| CTE definitions | 196/200 | 196/200 |
| prefix | 200/200 | 200/200 |

TiDB has no completion YAML in Bytebase and its adapter computed no signals of its own, so its
parity run is Bytebase's MySQL corpus through TiDB's reader; every statement in it is TiDB-accepted.

Every disagreement, classified:

- **17 — omni binds a reference Bytebase loses at a truncated statement** (`FROM t1 JOIN |`,
  `FROM t1, |`, `LIMIT |`, `SET c1 = |`, `CREATE INDEX i ON t1 (|`, `FROM t|`, …).
  `parseTableReferences` re-parses the fragment with a shrinking prefix and gives up; the reader
  reads the tokens.
- **2 — complete statements whose references Bytebase never collects** (`UPDATE t1 SET |`,
  `ALTER TABLE t1 ADD COLUMN c2 int REFERENCES |`).
- **1 — a references clause binds its own table** (`ALTER TABLE t1 ADD CONSTRAINT fk FOREIGN KEY
  (c1) REFERENCES t2(|)`: omni says `t2`, Bytebase keeps the altered `t1`).
- **10 — set operations.** Bytebase pools the first arm's tables into the caret's scope across
  `UNION`/`INTERSECT`/`EXCEPT`; omni keeps arms separate. Both approximate the engine; an arm sees
  only its own `FROM`, and a global `ORDER BY` names the first arm's output names.
- **3 — DML target placement.** Bytebase flattens the write target into `references`; omni reports
  it in `Scope.DMLTarget` (and as an outer reference from a subquery of an `UPDATE`/`DELETE`). The
  adapter unions them.
- **4 CTEs.** Two are the boundary itself (`SELECT *` bodies have no statically readable columns;
  Bytebase resolves them with a query span). One is a body whose select list does not read. One is
  `RECURSIVE`, where Bytebase drops the columns and omni reads the anchor arm's.
- **1 object kind.** `DROP INDEX |`: Bytebase offers databases and tables because it cannot resolve
  an index without its table; omni says `index`.
- **1 select-alias position.** `SELECT c1, SUM(c2) OVER (ORDER BY |) FROM t2`: Bytebase treats
  every `ORDER BY` as alias-accepting; a window specification's `ORDER BY` cannot name a select
  alias, so omni says no.

No disagreement is unexplained.

## The fork deltas

The fork readers are byte-identical to the parent's body; only the file header and the corpus
differ. Verified dialect facts, each a recorded engine verdict:

- **MariaDB** (`parser/testdata/compat.tsv`): `UNION`/`INTERSECT`/`EXCEPT` are set operators;
  `MINUS` is only a function/label keyword; `RETURNING` is accepted on `DELETE`/`INSERT`/`REPLACE`
  and is reserved, so it is never read as an alias. MariaDB rejects `TABLE t` and rejects the
  INSERT row alias by name ("not supported in MariaDB"), so those corpus cases are comments.
- **OceanBase** (`parser/testdata/compat.tsv`): `UNION`/`INTERSECT`/`EXCEPT` accepted; the engine
  accepts `MINUS` as a set operator and this parser does not yet (`set_type/MINUS` is
  omni=reject, engine=accept), so a MINUS arm ends at the next query block and the corpus pins
  that with the row named. Rejects `TABLE t`; accepts the INSERT row alias. Hints are comments;
  unreserved keywords are aliases.
- **TiDB** (`parser/testdata/compat.tsv`): `UNION`/`INTERSECT`/`EXCEPT` accepted, no `MINUS` or
  `RETURNING`; `AS OF TIMESTAMP` is not an alias because `OF` is reserved; accepts `TABLE t`;
  rejects the INSERT row alias.
- **TiDB's boundary fix**: `tidb/completion` is deleted. Bytebase must resolve the signals itself;
  its pin must not move to this branch until its adapter does.

## What a resumed session must do first

1. **Port the parent's reader as it stands after this PR merges.** Copy `complete_context.go` and
   `complete_context_test.go` from the merged parent, keep each fork's header comment, re-splice the
   fork's dialect cases, and drop the cases the engine rejects (the corpus comments list them). The
   parent gained a lot after the forks were cut — a `CREATE DEFINER` clause before the object type,
   the first arm of a set operation after the statement's `WITH` and inside its own query block, an
   unfinished `REFERENCES` list binding the parent only from inside it, a `CASE` terminator not being
   an implicit alias, a join condition that stops narrowing at the next clause, a trigger's `NEW` and
   `OLD` rows visible in a query block of its body, an unterminated parenthesized join keeping its
   members, an INSERT row alias only after `VALUES` or `SET`, `END CASE` not swallowing the next
   statement, and a `WITH` body read once per clause — so a fork copied from an older head carries
   roughly a dozen defects the merged corpus already pins.
2. Re-run per fork: the corpus and nesting tests, the six mutation checks (the parent's list is
   in #1291's body), and the 200-case parity replay. Update the parity table if it moves.
3. **Work the fork review threads** (25–27 per PR at parking). They are the same classes of
   findings Codex raised on MySQL after the forks were cut; fix or refute each with a corpus case or
   evidence, reply, and resolve. One PR per engine, merged on its own gate.
4. Merge `origin/main` into each branch (never rebase, never force-push), then merge
   `gh pr merge --squash --delete-branch` once CI is green and Codex reports no major issues.
5. Only then bump Bytebase's pin and land its adapter, whose call is the shape TiDB's already is:

   ```go
   ctx := mysqlparser.CollectCompletion(sql, byteOffset) // or mariadbparser/oceanbaseparser/tidbparser
   ```

   then `ctx.Candidates`/`ctx.Intent` for keywords and object classes, `ctx.Scope` + `ctx.CTEs`
   resolved through `c.GetDatabaseMetadataFunc`, `ctx.SelectAliases`, `ctx.Prefix`. Derived tables
   and CTEs with no statically readable columns keep the query-span path over live metadata.

## One owner per branch

A second agent worked #1291's branch while this file's author did. The branch carries commits
neither could fully account for, the review replies were posted twice for one round, and a `git
checkout` in the shared worktree reverted uncommitted edits once. The rule for resuming: **one
owner per branch.** Split work by PR (or by fork) and never by file, and do not share a git
worktree between agents.

## Parking notes

- The branch was parked, then recreated as a fresh PR against the same change, because #1291's
  Codex review quota ran out: a new PR gets a new review run once the quota resets, and a `failed`
  review is neither a pass nor a finding. The cut is still not switched on in Bytebase, so nothing
  forces a premature merge.
- Do not retry Codex while its usage limit is exhausted. A `failed` review is not a finding.
- The Mongo design is merged (`docs/specs/2026-11-17-mongo-completion-cut-design.md`); it is its own
  cut and needs a `Signals` struct of shell contexts, not a copy of this reader.
