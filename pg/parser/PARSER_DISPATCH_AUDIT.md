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
