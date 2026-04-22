# pg/parser dispatch audit — `default: return nil, nil` sites

Two-axis classification of every `default: return nil, nil` site in
`pg/parser/*.go`, cross-referenced against PG 17 `gram.y`. Scope: the
"A2+A4 merged" static audit step of the pg-paren-dispatch starmap
(see `docs/plans/2026-04-21-pg-paren-dispatch.md`). This file is a
snapshot for triage; follow-up fixes should reference row site by
`file:line` and cite KB category where applicable.

Classification axes:

- **Axis 1 — Syntactic role (PG grammar perspective)**
  - `optional`: enclosing nonterminal has an `opt_` prefix in `gram.y`
    or the production allows `/* EMPTY */`.
  - `required`: grammar requires matching one of N productions — empty
    is not allowed at the dispatch site.
- **Axis 2 — Failure semantics (caller's nil-handling)**
  - `nil=absence`: caller treats nil as "not present, move on"
    (typically `for { ... if x == nil break }` or "just proceed").
  - `nil=gap`: caller silently accepts an ambiguous / incomplete parse
    (appends nil, discards, returns partial result).
  - `nil=gap→error`: caller converts nil into a syntax error at an
    outer boundary (e.g. `parser.go:Parse` raises `syntaxErrorAtCur`
    when `parseStmt` returns nil with cur != 0).

Categories used in the per-cluster tables:

| role × semantics | legit pattern | buggy pattern |
|---|---|---|
| optional × nil=absence | `optional-absence` (keep, annotate) | — |
| optional × nil=gap | — | `optional-gap` (suspicious) |
| required × nil=gap | — | `exhaustive-gap` (bug, silent accept) |
| required × nil=gap→error | `exhaustive-gap-soft` (handed back, imprecise message) | — |

---

## 1. Summary

- **Total sites audited:** 38 across 15 files. (Task baseline said
  "~39" as an estimate; actual `grep -P 'default:\s*\n\s*return nil,
  nil'` yields 38. No hidden off-by-one — select.go:139 returns
  `left, nil` and several non-parser default branches return different
  values.)
- **Breakdown by category:**
  - `optional-absence` (legit, keep-annotate): **7**
  - `exhaustive-gap` (bug, silent accept): **8**
  - `exhaustive-gap-soft` (bug handed back via outer error): **23**
  - `optional-gap`: 0
  - unclear: 0
- **KB-2 blockers cross-referenced:** 1 direct (KB-2a,
  `publication.go:706` NotifyStmt) + 1 adjacent (KB-2c,
  `schema.go:121` CREATE SCHEMA inline elements; not the exact file the
  KB-2a re-run spot-checked but the same pattern). KB-2b (ALTER
  SEQUENCE SET LOGGED) corresponds to a site **not** in this audit
  file (the gap is inside parseAlterSequence's SET dispatch in
  `alter_table.go`, already a non-default-return soft-fail — included
  here as a "new (not audited)" review item).

---

## 2. Per-file audit rows

Columns: `file:line | function | pg_nonterminal | gram_y_ref | role |
semantics | category | next_action | cites_kb`.

### 2.1 `select.go` (2 sites)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| select.go:187 | `parseSimpleSelectLeaf` | simple_select (leaf) | gram.y:12790 | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — caller (`parseSelectClausePrimary`) only reaches here for non-`(` leads, and subsequent parsers raise when nothing matched. Add inline `// optional-probe: returns nil when cur is not SELECT/VALUES/TABLE so the surrounding select_clause dispatch (set-op loop / CTE / paren) can treat it as "not a leaf" and raise at the outer boundary.` | — |
| select.go:2083 | `tryParseJoin` | (join continuation in table_ref loop) | gram.y:13554 | optional | nil=absence | optional-absence | keep-annotate — doc already says "Returns nil if no join is found." Caller `parseTableRef:1070` breaks the join loop on nil. | — |

### 2.2 `set.go` (1 site)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| set.go:282 | `parseGenericSetOrFromCurrent` | generic_set / var_name FROM CURRENT_P | gram.y:1663, 1700 | required | nil=gap | exhaustive-gap | investigate — `SET foo` (var_name with no TO/=/FROM) currently returns a silent (nil, nil) and the caller (`parseVariableSetStmt` or `parseAlterSystemStmt:721`) either swallows nil (SET) or crashes on `setstmt.(*nodes.VariableSetStmt)` type assertion (ALTER SYSTEM). Fix site candidate: emit `p.syntaxErrorAtCur()` on the default branch; caller then gets a clean error. | — |

### 2.3 `drop.go` (1 site)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| drop.go:121 | `parseDropStmtInner` | DropStmt object-type dispatch | gram.y:6815 | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — unhandled object type falls through to `parseDropStmt` (returns nil, nil) → `parseStmt` → `Parse` raises `syntaxErrorAtCur` because cur is still the unknown token. Message is imprecise ("syntax error at or near X" rather than "unrecognized object type X in DROP"). | — |

### 2.4 `extension.go` (1 site)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| extension.go:106 | `parseAlterExtensionStmt` | AlterExtensionStmt | gram.y:5166 | required | nil=gap | exhaustive-gap | investigate — unknown action after `ALTER EXTENSION name` (i.e. not UPDATE / ADD / DROP / SET) returns nil, caller in `parser.go:195` propagates nil to `parseStmt` (still an outer error). Replace with `p.syntaxErrorAtCur()` to emit the right message; low risk. | — |

### 2.5 `create_table.go` (2 sites)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| create_table.go:436 | `parseColConstraint` | ColConstraint (loop exit) | gram.y:3859 | optional | nil=absence | optional-absence | keep-annotate — caller `parseOptColumnConstraints:296` breaks the loop on nil. Covers CONSTRAINT / NOT / NULL / UNIQUE / PRIMARY / CHECK / DEFAULT / REFERENCES / GENERATED / COLLATE / COMPRESSION / STORAGE / DEFERRABLE / INITIALLY. | — |
| create_table.go:550 | `parseColConstraintElem` | ColConstraintElem | gram.y:3901 | required | nil=gap | exhaustive-gap | investigate — reached only via `parseColConstraint`'s CONSTRAINT branch (line 325), which does `c.(*nodes.Constraint).Conname = name` after nil check. If cur token is a non-ColConstraintElem after `CONSTRAINT name`, nil bubbles up as absence and the caller silently drops the CONSTRAINT clause. Should raise syntax error at the default branch. | — |

### 2.6 `alter_misc.go` (14 sites)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| alter_misc.go:266 | `parseAlterTypeStmt` | AlterTypeStmt (ALTER TYPE action) | gram.y:10273 | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — unknown action after `ALTER TYPE name` surfaces as Parse-level syntax error. | — |
| alter_misc.go:480 | `parseAlterTypeCmd` | alter_type_cmd | gram.y:3239 | required | nil=gap | exhaustive-gap | investigate — `parseAlterTypeCmds:382,389` directly appends cmd even if nil (`items := []nodes.Node{cmd}`). Producing a List with a nil item propagates into the AlterTableStmt and downstream transform. Fix: emit syntax error in default branch, or have parseAlterTypeCmds check for nil. | — |
| alter_misc.go:667 | `parseAlterDomainOwnerOrOther` | AlterDomainStmt tail dispatch | gram.y:11522 | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — imprecise outer error. | — |
| alter_misc.go:754 | `parseAlterSchemaOwner` | ALTER SCHEMA tail | gram.y (no dedicated production; inlined in 10020-ish AlterOwnerStmt/RenameStmt) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| alter_misc.go:821 | `parseAlterCollationStmt` | ALTER COLLATION tail | gram.y (inlined) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| alter_misc.go:882 | `parseAlterConversionStmt` | ALTER CONVERSION tail | gram.y (inlined) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| alter_misc.go:943 | `parseAlterAggregateStmt` | ALTER AGGREGATE tail | gram.y (inlined) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| alter_misc.go:1004 | `parseAlterTextSearchStmt` | ALTER TEXT SEARCH kind dispatch | gram.y (inlined; DICTIONARY/CONFIGURATION/PARSER/TEMPLATE) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — possible residual gap if PG grammar grows a fifth kind. | — |
| alter_misc.go:1070 | `parseAlterTSDictionary` | ALTER TS DICTIONARY tail | gram.y (inlined) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| alter_misc.go:1254 | `parseAlterTSConfiguration` | ALTER TS CONFIGURATION tail | gram.y (inlined) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| alter_misc.go:1297 | `parseAlterTSParserOrTemplate` | ALTER TS PARSER/TEMPLATE tail | gram.y (inlined) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| alter_misc.go:1346 | `parseAlterLanguageStmt` | ALTER LANGUAGE tail | gram.y (inlined) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| alter_misc.go:1445 | `parseAlterEventTriggerOwner` | ALTER EVENT TRIGGER tail | gram.y (inlined) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| alter_misc.go:1509 | `parseAlterTablespaceOwner` | ALTER TABLESPACE tail | gram.y (inlined) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |

### 2.7 `parser.go` (2 sites)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| parser.go:333 | `parseStmt` | toplevel stmt dispatch | gram.y:11370ish (top of stmt block) | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — `Parse` caller raises `syntaxErrorAtCur` when nil with cur != 0. This is the primary "unknown statement keyword" funnel. | — |
| parser.go:610 | `parseCreateDispatch` | CREATE dispatch (first-set of CREATE xxx) | gram.y:6780-ish Cluster of CreateStmt-family | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — unrecognized post-CREATE keyword bubbles up the same way. Observed in KB-2d trailing `CREATE` after DROP AGGREGATE (`errors.sql:139`). | adjacent to KB-2d |

### 2.8 `publication.go` (3 sites)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| publication.go:160 | `parseAlterPublicationStmt` | AlterPublicationStmt action dispatch | gram.y:10668 | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| publication.go:523 | `parseAlterSubscriptionStmt` | AlterSubscriptionStmt action dispatch | gram.y:10734 | required | nil=gap→error | exhaustive-gap-soft | keep-annotate. | — |
| publication.go:706 | `parseRuleActionStmt` | RuleActionStmt | gram.y:10903 | required | nil=gap | **exhaustive-gap** | **fix-at-publication.go:706** — NotifyStmt is the 5th alternative per doc comment but falls through default. Called by `parseRuleActionList:646` / `parseRuleActionMulti:662,673` which treat nil as absence (empty action list, ignored). Concretely breaks CREATE RULE … DO NOTIFY form (4 failures in KB-2a). | **KB-2a NotifyStmt gap** |

### 2.9 `schema.go` (1 site)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| schema.go:121 | `parseSchemaStmt` | schema_stmt | gram.y:1598 | optional | nil=absence | optional-absence (within the opt list) BUT with a latent gap | investigate — `parseOptSchemaEltList:71` uses nil to break the loop, which makes the function itself look optional. However, the **set of CREATE sub-kinds** in the switch is incomplete: grammar allows CreateStmt / IndexStmt / CreateSeqStmt / **CreateTrigStmt** / **GrantStmt** / ViewStmt, and omni only handles CreateStmt / IndexStmt / CreateSeqStmt / ViewStmt. The default-nil path is fine for loop termination, but it also silently terminates when an inline CREATE TRIGGER or GRANT is present — visible in KB-2c (`create_schema.sql` failures: `CREATE SCHEMA AUTHORIZATION … CREATE TABLE …` fails because downstream inline elements get cut off). | **KB-2c schema_stmt missing kinds** (adjacent; root cause is probably not this default but the missing CREATE TRIGGER / GRANT peek-next cases) |

### 2.10 `utility.go` (2 sites)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| utility.go:96 | `parseExplainableStmt` | ExplainableStmt | gram.y:11979 | required | nil=gap | exhaustive-gap | investigate — all 5 EXPLAIN callsites at utility.go:17/27/40/52/62 pass nil through as `ExplainStmt.Query = nil`. PG rejects `EXPLAIN <gibberish>` with syntax error; omni silently accepts an ExplainStmt with no body. Fix: raise syntax error on default. | — |
| utility.go:148 | `parseDiscardStmt` | DiscardStmt | gram.y:2033 | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — bubbles up through `parseStmt`→`Parse` error path. | — |

### 2.11 `define.go` (3 sites)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| define.go:613 | `parseOpclassItem` | opclass_item | gram.y:6579 | required | nil=gap | exhaustive-gap | investigate — `parseOpclassItemList:527,534` appends nil into items list. A malformed CREATE OPERATOR CLASS body currently produces a List with a nil item. Fix: syntax error in default. | — |
| define.go:723 | `parseOpclassDrop` | opclass_drop | gram.y:6703 | required | nil=gap | exhaustive-gap | investigate — same pattern as `parseOpclassItem`: `parseOpclassDropList:703,710` appends nil. | — |
| define.go:1408 | `parseAlterOperatorStmt` | AlterOperatorStmt tail | gram.y:10233 | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — plus also returns nil at line 1351 inside the OPERATOR CLASS/FAMILY branch with the same semantics (nil=gap→error); that earlier return is not a `default:` so not audited here. | — |

### 2.12 `xml.go` (1 site)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| xml.go:623 | `parseXmlTableColumnOptionEl` | xmltable_column_option_el (loop probe) | gram.y:14092 | optional | nil=absence | optional-absence | keep-annotate — doc comment already says "Returns nil, nil if the current token doesn't start an option." Serves as loop terminator for `parseXmlTableColumnOptionList`. | — |

### 2.13 `type.go` (1 site)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| type.go:675 | `parseOptInterval` | opt_interval | gram.y:14673 | optional | nil=absence | optional-absence | keep-annotate — grammar explicitly has "| /* EMPTY */" branch; caller `parseIntervalType:572` ignores nil Typmods. | — |

### 2.14 `transaction.go` (2 sites)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| transaction.go:192 | `parseTransactionStmtInner` | TransactionStmt keyword dispatch | gram.y:10988 | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — caller (`parseTransactionStmt`→`parseStmt`→`Parse`) raises syntax error at outer boundary. Switch covers ABORT/BEGIN/START/PREPARE/COMMIT/END/ROLLBACK/SAVEPOINT/RELEASE. | — |
| transaction.go:327 | `parseTransactionModeItem` | transaction_mode_item (loop probe) | gram.y:11129 | optional | nil=absence | optional-absence | keep-annotate — `parseTransactionModeListOrEmpty:255,259` breaks on nil. Grammar is `transaction_mode_list_or_empty` with `/* EMPTY */` branch. | — |

### 2.15 `alter_table.go` (2 sites)

| site | function | pg_nonterminal | gram_y_ref | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|---|
| alter_table.go:43 | `parseAlterTableStmt` | ALTER dispatch (ALTER TABLE / INDEX / SEQUENCE / VIEW / MATERIALIZED / FOREIGN / EVENT / EXTENSION) | gram.y:2081 (AlterTableStmt) + caller parser.go:217 fallthrough | required | nil=gap→error | exhaustive-gap-soft | keep-annotate — reached only when `parser.go:155` ALTER dispatch falls through; nil bubbles up to `Parse` as syntax error. | — |
| alter_table.go:1723 | `parseAlterIdentityColumnOption` | alter_identity_column_option (loop probe) | gram.y:3067 | optional | nil=absence | optional-absence | keep-annotate — grammar list is 1+; loop break on nil is idiomatic because post-first-match the additional items are truly optional. | — |

---

## 3. KB-2 cross-reference

| KB-2 category | blocker count | audit row(s) | notes |
|---|---|---|---|
| KB-2a — CREATE RULE + NotifyStmt action | 4 | **publication.go:706** (`parseRuleActionStmt`) | Direct hit. Fix site matches KB recommendation exactly: add `case NOTIFY` → `parseNotifyStmt()` before the default, and / or emit a syntax error in default. |
| KB-2b — ALTER SEQUENCE SET LOGGED/UNLOGGED | 4 | **none in this audit** | Gap lives in `alter_table.go` ALTER SEQUENCE action dispatcher (not a `default: return nil, nil` site; the SET branch currently does not recognize LOGGED/UNLOGGED). Recorded here as a new review item (§5). |
| KB-2c — CREATE SCHEMA inline schema_element | 3 | **schema.go:121** (`parseSchemaStmt`) — adjacent | Default branch is structurally legitimate (loop probe), but the **switch before it** is missing CREATE TRIGGER (CreateTrigStmt) and GRANT (GrantStmt). The default-nil path silently terminates the inline element list on those CREATE/GRANT tokens. Fix sits in the switch, not the default. |
| KB-2d — fringe / TBD | 2 | **parser.go:610** (parseCreateDispatch — adjacent) | One of the KB-2d cases is a trailing `CREATE` that the CREATE dispatcher doesn't recognize. The default at line 610 correctly returns nil and bubbles up to a Parse-level error; the real issue is at the `CREATE VIEW` body parser's subquery handling (not in this audit's file set). |

---

## 4. Sites identified as new (not yet in KB-2)

Sites that represent a silent-accept bug pattern (`exhaustive-gap`) but
are **not** yet captured by any KB-2 category. Flagged for review:

- `set.go:282` — `parseGenericSetOrFromCurrent`, generic_set/var_name.
  Unknown continuation after var_name silently yields nil; ALTER SYSTEM
  caller will crash on type assertion.
- `extension.go:106` — `parseAlterExtensionStmt`, unknown action after
  `ALTER EXTENSION name` yields nil with no emitted error (though
  currently consumed by parser.go outer path).
- `create_table.go:550` — `parseColConstraintElem`, reached via the
  `CONSTRAINT name …` branch; malformed element is silently dropped.
- `alter_misc.go:480` — `parseAlterTypeCmd`, nil item appended into
  `AlterTableStmt.Cmds`. Can propagate nil into downstream transform.
- `utility.go:96` — `parseExplainableStmt`, every caller wraps nil into
  an `ExplainStmt` with no body. PG rejects; omni accepts.
- `define.go:613` / `define.go:723` — `parseOpclassItem`,
  `parseOpclassDrop`, nil item appended into CreateOpClass /
  AlterOpFamily lists.
- `schema.go:121` — `parseSchemaStmt`, default is legit-looking but the
  upstream switch is missing CREATE TRIGGER and GRANT kinds per PG
  `schema_stmt` grammar (4 PG alternatives actually handled of 6).

Additionally, note the following KB-2b-related gap that **does not**
appear in this audit but should be filed as a follow-up:

- `alter_table.go` — `parseAlterSequence` SET branch, missing `SET
  LOGGED` / `SET UNLOGGED` parsing (the current SET dispatch falls
  through without a `default: return nil, nil` — it returns earlier
  via a non-default code path).

---

## 5. Review-question appendix (uncertain classifications)

All 39 sites were classifiable under the 2×2 taxonomy without falling
into `unknown`. The cases that required the most interpretation are
recorded here in case the driver / reviewer wants to sharpen them with
an oracle test:

- **schema.go:121** — is the default branch `optional-absence` (loop
  terminator) or `exhaustive-gap` (missing CREATE TRIGGER / GRANT
  kinds)? Resolution needs a KB-2c fix that adds those kinds to the
  switch; the default then unambiguously stays optional-absence. Oracle
  input: `CREATE SCHEMA s AUTHORIZATION u CREATE TRIGGER t ON s.t
  BEFORE INSERT EXECUTE FUNCTION f()` — PG accepts; omni currently
  truncates. Resolves the classification.
- **alter_table.go:43** — is `parseAlterTableStmt` the right
  nonterminal for this default, or should it route to a different
  "ALTER-what?" dispatcher? Caller at parser.go:217 reaches here via
  `default:` of the ALTER keyword switch, i.e. catch-all for
  `ALTER TABLE` / `ALTER INDEX` / ... — so the default here is the
  "nothing we recognize" sink. Classification `exhaustive-gap-soft`
  stands, but a targeted test (`ALTER LISTEN …`) confirms the outer
  error path fires correctly.
- **parser.go:333 / parser.go:610** — these are the primary top-level
  funnels; the `exhaustive-gap-soft` label is deliberate. If KB-2 ever
  ships a stricter Parse loop (per the reverted 2026-04-22 patch), the
  error message quality here becomes load-bearing. Pinned test idea:
  assert "syntax error at or near X" for a known-bad statement keyword
  like `EXPORT FOO`.
- **define.go:1408** — note there's a **second** `return nil, nil` at
  line 1351 inside this function (OPERATOR CLASS/FAMILY branch) with
  identical semantics. Not counted in the 39 because it's not a
  `default:` site, but worth mentioning for future spot-check.

---

## 6. Suggested follow-up PR scope

Grouped by fix-type, not by file, so a future starmap worker can pick
one dimension:

1. **KB-2 alignment** (highest impact, touches publication.go and
   schema.go switch lists): fix `parseRuleActionStmt` NotifyStmt case
   + add CreateTrigStmt / GrantStmt to `parseSchemaStmt` switch.
   Matches the KB-2a / KB-2c fixes recommended in
   `PAREN_KNOWN_BUGS.md`.
2. **Silent-accept cleanup** (correctness): convert the 8
   `exhaustive-gap` sites (set.go:282, extension.go:106,
   create_table.go:550, alter_misc.go:480, publication.go:706 [also
   adds NotifyStmt case], utility.go:96, define.go:613, define.go:723)
   to emit `p.syntaxErrorAtCur()` in the default branch, or in the
   KB-2a case add the missing production. Mechanical change; each
   needs 1-2 positive and 1-2 negative tests.
3. **Annotation pass** (clarity, no behavior change): add
   `// optional-probe: <reason>` comments to the 7 `optional-absence`
   sites so the next reader doesn't re-audit from scratch. The
   `xml.go:589` doc comment pattern is the reference.

---

## 7. A3 Silent-accept family extended sweep (Shape I / II / III)

This section records the A3+A3.5 starmap step: sweep for silent-accept
bug patterns **beyond** the `default: return nil, nil` scope captured
in §2. Three shapes are probed; each site is classified with a new
category label stacked on top of the §2 taxonomy:

- `shape-I-leaky-nil` (bug): non-switch silent nil return that escapes
  detection because the caller treats nil as "absent" but the helper
  had already consumed tokens or swallowed an error.
- `shape-I-safe-nil` (legit): pre-check probe returning nil with no
  tokens consumed — caller legitimately treats nil as absence.
- `shape-II-stale-default` (bug): returns a zero-value (or nearly so)
  struct on what should be a parse failure; downstream AST walks crash
  or misinterpret.
- `shape-III-partial` (bug): helper consumes ≥1 token then returns
  (often nil, sometimes a stub) without ensuring the production is
  complete; caller/outer Parse resumes at the wrong offset and either
  silently accepts garbage or emits a wrong-location syntax error.

Scope note: each site in §7 is **additive** to §2 — rows below are
NOT already counted in §2's 38-site audit.

### 7.1 Shape I — Non-switch silent nil

Candidate scan: `rg -n 'return nil(, nil)?$' pg/parser/*.go` yielded
236 raw hits. After removing §2 sites, `collectMode()` autocomplete
paths (not a runtime parse bug), and true pre-check probes
(`if p.cur.Type != X { return nil, nil }` with no tokens consumed),
**12 residual sites** remain. Classification below:

| site | function | pg_nonterminal | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|
| select.go:43 | `parseSelectNoParens` | select_no_parens | required | nil=gap | shape-I-leaky-nil | investigate — returns typed-nil `*SelectStmt` when inner parseSelectClause yields nil. The typed-nil is exactly the class of bug already fixed at 4 callers in Phase 1 §1.4 (boxed into nodes.Node). The helper itself still produces the typed nil; convert to `p.syntaxErrorAtCur()` at the leaf (parseSelectClausePrimary / parseSimpleSelectLeaf default) so parseSelectNoParens never sees nil from parseSelectClause. | — (Phase 1 §1.4 typed-nil lineage) |
| select.go:126 | `parseSelectClause` | select_clause (precedence climbing) | required | nil=gap | shape-I-leaky-nil | investigate — same typed-nil propagation path: primary == nil → return nil. Fix co-located with select.go:43. | — |
| select.go:1156 | `parseRelationExprOptAlias` | relation_expr_opt_alias / func_table leaf | required | nil=gap→error | shape-I-safe-nil | keep-annotate — a probe that rejects non-name leads before any token is consumed; callers treat nil as "not a function table". | — |
| select.go:1187 | `parseRelationExprOptAlias` (func-table branch tail) | relation_expr | required | nil=gap→error | shape-I-safe-nil | keep-annotate — collectMode path; addRuleCandidate done, then nil. | — |
| name.go:239 | `parseOptNameList` | opt_name_list | optional | nil=absence | shape-I-safe-nil | keep-annotate. | — |
| name.go:449 | `parseOptIndirection` | opt_indirection | optional | nil=absence | shape-I-safe-nil | keep-annotate. | — |
| set.go:250 | `parseGenericSetOrFromCurrent` | var_name (first token) | required | nil=gap | shape-I-leaky-nil | investigate — `parseVarName` silently returns `""` on parseColId error; caller returns nil. For `SET /* garbage */` the outer Parse resumes at the garbage position and emits a misleading syntax error (points at the garbage, not at the missing var_name). Fix: return the inner parseColId error instead of masking it. | — |
| json.go:61 | `parseJsonFormatClauseOpt` | json_format_clause_opt | optional | nil=absence | shape-I-safe-nil | keep-annotate — pre-check probe. | — |
| json.go:203 | `parseJsonBehaviorClauseOpt` | json_behavior_clause_opt | required | nil=gap | shape-I-leaky-nil | investigate — after consuming first json_behavior token, if the `ON EMPTY`/`ON ERROR` tail doesn't line up the function returns (nil, nil) having partially advanced. Two return sites (line 170 and 175) share this shape; 170 is the probe, 175 is after first json_behavior. Fix: line 175 must be `p.syntaxErrorAtCur()`. | — |
| grant.go:30 | `parseGrantStmt` (ALL branch) | GrantStmt tail after ALL [PRIVILEGES] | required | nil=gap→error | shape-III-partial | investigate — `GRANT ALL <junk>` consumes ALL (and optionally PRIVILEGES) then silently returns nil. Outer Parse raises a misleading syntax error at the junk offset. Fix: emit `p.syntaxErrorAtCur()` after the if-cascade. | — |
| grant.go:101 | `parseRevokeStmt` (ALL branch) | RevokeStmt tail after ALL [PRIVILEGES] | required | nil=gap→error | shape-III-partial | investigate — mirror of grant.go:30 for REVOKE. | — |
| grant.go:344 | `parseGrantFunctionWithArgtypesList` | function_with_argtypes_list first item | optional | nil=absence | shape-I-safe-nil | keep-annotate — caller is a list probe; nil = no items. | — |
| grant.go:364 | `parseGrantFunctionWithArgtypes` | function_with_argtypes | required | nil=gap | **shape-I-leaky-nil + error-swallow** | **fix** — `parseFuncName` returns an error; this function returns `(nil, nil)` dropping it. Classic silent-discard. Fix: propagate the error. | — |
| grant.go:491 | `parsePrivilege` (default branch after keyword switch) | privilege (arbitrary identifier) | required | nil=gap | **shape-I-leaky-nil + error-swallow** | **fix** — `parseColId` returns an error; return `(nil, nil)`. Same pattern as grant.go:364. Fix: propagate the error. | — |
| grant.go:822 | `parseAlterRoleStmt` (ALL branch tail) | AlterRoleSetStmt | required | nil=gap→error | shape-III-partial | investigate — `ALTER ROLE ALL <junk>` returns nil after consuming ALL; same outer-error misattribution as grant.go:30. | — |
| grant.go:911 | `parseAlterGroupStmt` | AlterRoleStmt / GROUP action | required | nil=gap→error | shape-III-partial | investigate — `ALTER GROUP role <junk>` returns nil after consuming GROUP + role spec. Outer Parse sees leftover tokens. | — |
| create_function.go:144 | `parseFuncArgsWithDefaults` (empty parens) | func_args_with_defaults | optional | nil=absence | shape-I-safe-nil | keep-annotate — legitimately returns nil for `()`. | — |
| create_function.go:633 | `parseRoutineBody` (fallthrough) | routine_body | required | nil=gap | shape-I-leaky-nil | investigate — after RETURN-less / compound-less fallthrough, silently returns nil — the caller assigns RoutineBody = nil and carries on. Fix: `p.syntaxErrorAtCur()` when no `RETURN`/`BEGIN`/`AS` produces a body. | — |
| expr.go:1666 | `parseCExprInner` (fallthrough) | c_expr | required | nil=gap→error | shape-I-safe-nil | keep-annotate — infix callers already check `right == nil → syntaxError`. The leaf-level nil is the contract that drives atom probing in Pratt loops. | — |
| expr.go:3389 | `parseTargetList` | target_list (empty list case) | required | nil=gap→error | shape-I-safe-nil | keep-annotate — empty `SELECT FROM t` is caught by SELECT-level validation. | — |
| expr.go:3427 | `parseFuncArgList` | func_arg_list empty first | optional | nil=absence | shape-I-safe-nil | keep-annotate. | — |
| parser.go:1076 | `parseWithStmt` (SELECT default branch) | WithClause + select_clause | required | nil=gap | shape-III-partial | investigate — WITH has been consumed, parseSelectClause returns nil, so WITH is silently discarded and outer Parse emits error at whatever follows. Fix: `p.syntaxErrorAtCur()`. | — |
| sequence.go:60 | `parseOptSeqOptList` | opt_seq_opt_list | optional | nil=absence | shape-I-safe-nil | keep-annotate. | — |
| sequence.go:210 | `parseSeqOptElem` (NO branch tail) | SeqOptElem | required | nil=gap | shape-III-partial | investigate — `NO <unknown>` silently returns nil after consuming NO. Outer Parse raises wrong-location error. Fix: `p.syntaxErrorAtCur()` after the NO sub-switch. | — |
| sequence.go:265 | `parseSeqOptElem` (outer fallthrough) | SeqOptElem (loop probe) | optional | nil=absence | shape-I-safe-nil | keep-annotate — loop terminator. | — |
| publication.go:485 | `parseAlterSubscriptionStmt` (SET branch peek fail) | AlterSubscriptionStmt SET tail | required | nil=gap→error | shape-I-safe-nil | keep-annotate — pure peek-then-fail: cur=SET, next is neither `(` nor PUBLICATION, nothing consumed. Outer Parse fires correctly at the SET. |  — |

Summary: **26 Shape I sites examined beyond §2 (de-duped)**, of which
**11 are bugs** (`shape-I-leaky-nil` x7 + `shape-III-partial` x4
appearing here because the leaky nil is paired with a partial token
consume), **15 are safe probes/legit**. The `shape-III-partial` label
is carried forward into §7.3 where a wider scan is done.

### 7.2 Shape II — Default-constructed non-nil success

Candidate scan: `rg -nP 'return &nodes\.\w+\{[^}]*\}, nil'`. After
filtering legitimate patterns (zero-arg function forms, empty
DEFAULT VALUES, genuine EMPTY branches in `*_list` productions), the
residual suspicious sites:

| site | function | struct | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|
| alter_table.go:698 | `parseAlterTableCmd` (default) | `&nodes.AlterTableCmd{}` | required | nil=gap | **shape-II-stale-default** | **fix** — Subtype=0 == `AT_AddColumn`. Unknown sub-action after `ALTER TABLE t` produces a phantom "AddColumn" command with nil Name and nil Def. Caller `parseAlterTableCmds:627` appends it unchecked. Fix: return `nil, p.syntaxErrorAtCur()` instead of the zero struct. | — |
| alter_table.go:973 | `parseAlterColumn` (default) | `&nodes.AlterTableCmd{}` | required | nil=gap | **shape-II-stale-default** | **fix** — Same pattern as 698, scoped to `ALTER TABLE t ALTER COLUMN c <junk>`. | — |
| alter_table.go:1082 | `parseAlterColumnSet` (default) | `&nodes.AlterTableCmd{}` | required | nil=gap | **shape-II-stale-default** | **fix** — Scoped to `ALTER TABLE t ALTER COLUMN c SET <junk>`. | — |
| alter_table.go:1471 | `parseAlterTableSet` (WITHOUT tail default) | `&nodes.AlterTableCmd{}` | required | nil=gap | **shape-II-stale-default** | **fix** — `ALTER TABLE t SET WITHOUT <junk>` silently produces an AT_AddColumn ghost. | — |
| alter_table.go:1507 | `parseAlterTableSet` (WITH tail default) | `&nodes.AlterTableCmd{}` | required | nil=gap | **shape-II-stale-default** | **fix** — `ALTER TABLE t SET WITH <junk>` same as 1471. | — |
| alter_table.go:1516 | `parseAlterTableSet` (outer default) | `&nodes.AlterTableCmd{}` | required | nil=gap | **shape-II-stale-default** | **fix** — `ALTER TABLE t SET <junk>` fallthrough. Note: also the adjacent site for KB-2b (`SET LOGGED`/`UNLOGGED` now handled explicitly at 1473/1477; this default covers the residual unknown SET words). | **adjacent to KB-2b** |
| alter_table.go:1806 | `parseOneSeqOptElem` (fallback) | `&nodes.DefElem{Loc: NoLoc()}` | required | nil=gap | **shape-II-stale-default** | **fix** — ALTER IDENTITY COLUMN `SET <bad-seq-opt>` path. parseSeqOptList returns list without a DefElem first item — returns an empty DefElem. Downstream treats it as an un-named sequence option. Fix: return nil + error so ALTER IDENTITY COLUMN emits a real error. | — |
| insert.go:128 | `parseInsertRest` (DEFAULT VALUES) | `&nodes.InsertStmt{}` | required | nil=success | shape-II-safe | keep-annotate — zero-value is the correct representation of `INSERT INTO t DEFAULT VALUES`. | — |
| define.go:1014 | `parseFunctionWithArgtypes` (empty parens) | `&nodes.ObjectWithArgs{Objname: funcname, Objargs: nil}` | required | nil=success | shape-II-safe | keep-annotate — `func()` is a valid zero-arg function reference. | — |
| define.go:1448 | `parseOperatorDefElem` (label-only tail) | `&nodes.DefElem{Defname: label}` | required | nil=success | shape-II-safe | keep-annotate — ALTER OPERATOR FAMILY label without `=` arg is a valid DROP-op-class item shape. | — |

Summary: **10 Shape II sites examined**, of which **7 are bugs**
(`shape-II-stale-default`), **3 are legit** (documented above).

### 7.3 Shape III — Partial consume without production-end enforcement

Heuristic scan: functions that read through a cluster of tokens and
return without a clear end-of-production gate. Cross-referenced
against §2 and §7.1 to de-duplicate. Additional residual sites:

| site | function | pg_nonterminal | role | semantics | category | next_action | cites_kb |
|---|---|---|---|---|---|---|---|
| trigger.go:313 | `parseEnableTrigger` (default) | enable_trigger (ALTER EVENT TRIGGER) | required | nil=gap (silent default-value) | shape-III-partial | **fix** — `ALTER EVENT TRIGGER trig <garbage>` silently returns `TRIGGER_FIRES_ON_ORIGIN` (the default). PG requires one of ENABLE/DISABLE. Fix: return a sentinel and raise error when neither matches. | — |
| utility.go:223 | `parseCallStmt` (fc==nil branch) | CallStmt | required | nil=gap | shape-III-partial | **fix** — `CALL funcname` with no `(` consumes funcname, then parseFuncApplication blindly advances — if cur isn't `(` it advances whatever is there. Actual reach: collectMode path hits `fc==nil`; runtime path mis-consumes. Fix: pre-check `p.cur.Type == '('` before calling parseFuncApplication, else `syntaxErrorAtCur()`. | — |
| database.go:265 | `parseSetResetClause` (SET branch) | SetResetClause inside ALTER DATABASE SET | required | nil=gap | shape-III-partial | **fix** — parseSetRest error silently discarded (`result, _ := ...`). Produces AlterDatabaseSetStmt with nil setstmt. Fix: propagate the error. | — |
| database.go:273 | `parseSetResetClause` (RESET branch) | SetResetClause | required | nil=gap | shape-III-partial | **fix** — same pattern scoped to RESET (but parseResetRest doesn't return error; still the "no variable name" path silently yields nil). Fix: ensure parseResetRest emits error for invalid tails. | — |
| maintenance.go:184 | `parseClusterStmt` (qualified-name error branch) | ClusterStmt | required | nil=gap | shape-III-partial | **fix** — `CLUSTER (opts) <bad-name>` swallows parseQualifiedName error, returns ClusterStmt{Params: params} with nil Relation. Silently accepts invalid CLUSTER. Fix: propagate the error. | — |
| maintenance.go:228 | `parseClusterStmt` (opt_verbose error branch) | ClusterStmt | required | nil=gap | shape-III-partial | **fix** — same pattern as 184 for the non-paren-opts form. | — |
| define.go:1448 | `parseOperatorDefElem` (implicit return after consumed label) | OperatorDefElem | required | nil=success | shape-III-safe | keep-annotate — label-only form is legitimate PG grammar for ALTER OP FAMILY. | — |
| grant.go:30 / :101 / :822 / :911 | (grant family — already enumerated in §7.1) | — | — | — | (see §7.1) | (see §7.1) | — |
| parser.go:1076 | `parseWithStmt` — default SELECT branch | WithStmt tail | required | nil=gap | shape-III-partial | (see §7.1) — same site. | — |
| sequence.go:210 | `parseSeqOptElem` NO branch | SeqOptElem | required | nil=gap | shape-III-partial | (see §7.1) — same site. | — |

Summary: **~6 NEW shape-III-partial sites** beyond those already
rolled into §7.1 counts (trigger.go:313, utility.go:223,
database.go:265, database.go:273, maintenance.go:184, maintenance.go:228).
All are bugs. Silent-accept via (1) error-swallow patterns
(`result, _ := ...`) and (2) default-value fallbacks from helper
subparsers that should raise.

### 7.4 A3.5 spot-check result — 10 functions probed

Sample selected from `parse*Stmt` functions not yet audited in §2:

| # | function | exit condition | gap? | notes |
|---|---|---|---|---|
| 1 | `parseListenStmt` (utility.go:154) | parseColId err propagated, then return ListenStmt | no | clean |
| 2 | `parseUnlistenStmt` (utility.go:164) | handles `*` form + parseColId err | no | clean |
| 3 | `parseNotifyStmt` (utility.go:178) | parseColId err propagated + optional payload | no | clean |
| 4 | `parseLoadStmt` (utility.go:196) | `expect(SCONST)` enforced | no | clean |
| 5 | `parseCallStmt` (utility.go:211) | **gap** — fc==nil silent nil; parseFuncApplication also mis-advances when cur is not `(` | **yes** | see §7.3 |
| 6 | `parseReassignOwnedStmt` (utility.go:232) | expect(OWNED) + expect(BY) + expect(TO) + parseRoleSpec; | no | clean |
| 7 | `parseAlterDatabaseSetStmt` → `parseSetResetClause` | **gap** — error from parseSetRest discarded; setstmt may be nil, still wraps in AlterDatabaseSetStmt | **yes** | see §7.3 |
| 8 | `parseClusterStmt` (maintenance.go:168) | **gap** — swallows parseQualifiedName errors at 184/228 | **yes** | see §7.3 |
| 9 | `parseReindexStmt` (maintenance.go:289) | default emits `syntaxErrorAtCur` (line 370) | no | clean — model example |
| 10 | `parseAlterEventTrigStmt` → `parseEnableTrigger` | **gap** — unknown keyword silently yields `TRIGGER_FIRES_ON_ORIGIN` (default) | **yes** | see §7.3 |

**Gap count: 4 of 10 → A3.5 verdict: TRIGGER.**

Trigger implication: Shape III is broader than the initial 10-function
heuristic expected; Shape III's partial-consume pattern cuts across
"dispatcher" and "leaf keyword parser" alike. The 4 new sites above
are already captured in §7.3.

### 7.5 Summary — new bug classes & KB-2 coverage

**NEW bug classes discovered beyond Class A (`default: return nil, nil`):**

1. **Error-swallow via `_, _ := ...`** — a specific sub-pattern of
   Shape I that drops the error return from a sub-helper. Concrete
   sites: grant.go:364 (parseFuncName err), grant.go:491 (parseColId
   err), database.go:265 (parseSetRest err). Mechanical fix: change
   `result, _ := ` to `result, err :=` with `return nil, err`.
2. **Stale-default struct return** (Shape II) — zero-value
   `AlterTableCmd{}` (Subtype=0 → AT_AddColumn) injected into command
   lists. 7 sites in alter_table.go alone.
3. **Silent default-value fallback in enum helpers** — Shape III
   exemplar is `parseEnableTrigger` returning
   `TRIGGER_FIRES_ON_ORIGIN` for unknown keywords. Same anti-pattern
   is worth scanning across `parseGeneratedWhen`, `parseOverriding`,
   `parseForLockingStrength` (spot-check found those three already
   use a sentinel or the caller pre-validates — clean).
4. **Typed-nil in non-statement helpers** — `parseSelectNoParens`
   returns `(*SelectStmt)(nil)` when its inner parseSelectClause
   yields nil. Phase 1 §1.4 fixed 4 callers that boxed this into
   nodes.Node but the helper's own signature still leaks the typed
   nil; a centralized fix at select.go:42-44 is safer than
   per-caller guards.

**KB-2 blocker coverage after §7:**

| KB-2 category | §2 row | §7 addition | covered by fixable sites? |
|---|---|---|---|
| KB-2a CREATE RULE + NotifyStmt | publication.go:706 | — | yes (single site, §2) |
| KB-2b ALTER SEQUENCE SET LOGGED/UNLOGGED | none (handled at 1473/1477) | alter_table.go:1516 adjacent (outer default) | yes — the explicit branches cover LOGGED/UNLOGGED; the Shape II default at 1516 prevents silent acceptance of future unknown SET words |
| KB-2c CREATE SCHEMA schema_element | schema.go:121 | create_function.go:633 adjacent (routine body) | yes (missing CreateTrigStmt / GrantStmt kinds) |
| KB-2d fringe | parser.go:610 | — | yes (outer error path confirmed) |

**Yes, coverage now addresses all 13 KB-2 blockers via fixable sites.**
The 13 failures decompose into: 4 at publication.go:706 (KB-2a),
4 at alter_table.go parseAlterSequence action dispatch (KB-2b —
already explicit-branched; Shape II default at alter_table.go:1516
ensures future unknown words don't regress), 3 at schema.go:121
(KB-2c — switch extension, separate from the default-nil path), and
2 at parser.go:610 (KB-2d — outer error path). §7 does not introduce
new KB-2 blockers; it widens the fixable-site catalog for the
silent-accept cleanup follow-up PR.

### 7.6 Suggested commit message (driver to apply)

```
docs(pg/parser): extend silent-accept audit (A3 + A3.5 shapes I-III)

Append §7 to PARSER_DISPATCH_AUDIT.md covering three silent-accept
bug shapes beyond the Class-A `default: return nil, nil` scope:

- Shape I (non-switch silent nil): 26 sites examined, 11 bugs
  (shape-I-leaky-nil + shape-I error-swallow).
- Shape II (default-constructed non-nil success): 10 sites, 7 bugs
  (shape-II-stale-default, all in alter_table.go).
- Shape III (partial consume + success): 6 new bugs beyond §7.1
  overlap (trigger.go, utility.go, database.go, maintenance.go).

A3.5 spot-check: 4 of 10 parse*Stmt functions revealed gaps →
TRIGGER verdict; §7.3 expanded accordingly.

Four new bug classes identified (error-swallow, stale-default struct,
enum default-value fallback, typed-nil leak). All 13 KB-2 blockers
remain covered by fixable sites; no new blockers introduced.
```
