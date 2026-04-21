# PAREN_AUDIT.json schema

**Machine-readable mirror of `PAREN_AUDIT.md`.** Source of truth for the
Phase 4 §5.3 two-proof-bar text (caller-context paragraphs + ≥5 pinned
tests for ambiguity-present sites; ≥1 test-file citation or pgregress
corpus entry for non-ambiguous sites). The markdown cluster tables in
`PAREN_AUDIT.md` carry the terse technique/status summaries — consult
this JSON for full `proof_notes`.

Enforced by `TestPARENAuditLint` in `paren_audit_lint_test.go`
(default build tag — runs on every PR via `.github/workflows/ci.yml`).

## Top-level shape

One JSON array of row objects. Each row describes one `(` or `)`
dispatch site in `pg/parser/*.go`.

```json
[
  {
    "site": "select.go:166",
    "function": "parseSelectClausePrimary",
    ...
  },
  ...
]
```

## Row fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `site` | string | yes | Original audit coordinate `<file>:<line>`. Line numbers are **stable row IDs**, NOT live code coords — see the note at the top of `PAREN_AUDIT.md`. For current location, grep `function`. |
| `function` | string | yes | Enclosing Go function name. Used for cross-reference against live code in the lint test. |
| `nonterminals` | array of string | yes | PG 17 grammar nonterminals this dispatch site covers (e.g. `["select_with_parens", "joined_table"]`). Zero-length allowed only for unconditional opens/closes with no grammar dispatch (rare). |
| `ambiguity_present` | bool | yes | True iff more than one PG nonterminal reaches this site and they share a `(` prefix (the site has to disambiguate). False for unconditional `expect('(')` / `expect(')')` and T1 peeks on optional lists. |
| `current_technique` | string | no | Which of the T1–T8 techniques from `docs/plans/2026-04-21-pg-paren-dispatch.md` §3 is applied. May be null for `none` (unconditional expect) or for blocked rows. |
| `pg_reference` | string | yes | `gram.y:<line>` or `gram.y:<start>-<end>` pointing at the PG 17 grammar production this site targets. Used to anchor future grammar-drift triage. |
| `aligned` | enum string | yes | One of `yes`, `no`, `blocked`, `unclear`. `aligned=yes` rows MUST have non-empty `proof_notes` (lint-enforced). |
| `blocked_by` | string or null | no | If `aligned=blocked`, names the upstream dependency (e.g. `"pg-nonterminal-alignment"`, `"pg-first-sets"`). Null otherwise. |
| `cluster` | string | yes | One of `C1`..`C5` matching the §4 cluster tables in the audit. Used for triage slicing. |
| `priority` | enum string | yes | `high`, `med`, `low`. Set at Phase 1 scope-lock; do not modify without an explicit review. |
| `proof_notes` | string | yes for `aligned=yes` | Free-form proof text. For ambiguity-present sites, must carry the §5.3 two-proof bar: (a) caller-context paragraph + (b) ≥5 pinned test citations. For non-ambiguous sites, must carry at least one test-file citation or pgregress entry. Empty string fails the lint gate. |
| `suspicion_notes` | string or null | no | Free-form "this could regress if …" notes. Null when the site is fully locked down. |

## Allowed values

### `aligned`

- `yes` — omni's `(` dispatch matches PG 17's grammar reduction for every reachable caller. Requires non-empty `proof_notes`.
- `no` — known divergence. `proof_notes` should describe the bug. Rare; all such rows should have an open scenario in `SCENARIOS-pg-paren-dispatch.md`.
- `blocked` — cannot align without an upstream change. Must set `blocked_by`.
- `unclear` — audit could not determine alignment. Should be re-audited before being closed.

### `cluster`

- `C1` — `table_ref` / `select_with_parens` / `joined_table` (primarily `select.go`).
- `C2` — expression parens (primarily `expr.go`).
- `C3` — type parens (primarily `type.go`).
- `C4` — DDL element lists (`create_table.go`, `create_index.go`, `define.go`).
- `C5` — utility statements (heterogeneous tail: `copy.go`, `grant.go`, `extension.go`, etc.).

### `priority`

- `high` — ambiguity-present and routinely exercised by user SQL.
- `med` — grammar-disambiguated but worth citation for regression fence.
- `low` — unconditional or T1-peek-on-optional; low regression risk.

## Invariants enforced by `TestPARENAuditLint`

1. **Schema completeness.** Every row has all required fields populated with the right types.
2. **Proof fence.** Every `aligned: yes` row has non-empty `proof_notes`. Deliberate regression test: removing a proof note makes the lint fail immediately.
3. **Allowed-value discipline.** `aligned` ∈ {yes, no, blocked, unclear}; `cluster` ∈ {C1..C5}; `priority` ∈ {high, med, low}.
4. **Site identifier uniqueness.** The `site` field is unique across rows (no two rows share a file:line audit coordinate).
5. **Code-vs-audit drift detection.** Every `p.cur.Type == '('` / `p.cur.Type == ')'` site in `pg/parser/*.go` (non-test) must resolve to a matching `(file, function)` pair in the audit. Drift is scoped via a baseline allowlist (`PAREN_AUDIT_DRIFT_BASELINE.txt`) so historical rename-drift captured at Phase 5 commit time doesn't block CI; **new** dispatch sites added after baseline fail the lint.

## Adding a new dispatch site

When a parser change introduces a new `(` / `)` dispatch site:

1. Add a row to `PAREN_AUDIT.json` with the full schema fields populated.
2. Mirror the terse summary into `PAREN_AUDIT.md` under the right cluster table.
3. Add at least one test (dedicated `*_test.go` for ambiguity-present sites; pgregress corpus entry acceptable for low-priority non-ambiguous sites).
4. Run `go test -run TestPARENAuditLint ./pg/parser/... -count=1` locally — it should pass.
5. If the lint gate flags the new site, either (a) you forgot to add the audit row, or (b) your dispatch uses an unusual pattern (`expect('(')` / `match('(')`) that the lint scanner skips. For (b), add the row anyway — the audit is the canonical inventory regardless of the specific Go pattern used.

## Re-baselining `PAREN_AUDIT_DRIFT_BASELINE.txt`

The baseline file records `(file, function)` pairs where the code has
drifted away from the audit's original function names (e.g. when a
function was refactored/renamed after the Phase 1 audit). It's a
one-way allowlist: discovered drift that is NOT in the baseline fails
the lint.

To re-baseline (rare — done when the audit rows themselves are updated
to match current code):

```bash
go test -run TestPARENAuditLint ./pg/parser/... -count=1 -v 2>&1 \
  | grep "UNDOCUMENTED" \
  > pg/parser/PAREN_AUDIT_DRIFT_BASELINE.txt
# review the diff and commit
```

The right fix is usually to update the audit row's `function` field to
match current code; re-baselining should only be done when that update
is intentionally deferred.
