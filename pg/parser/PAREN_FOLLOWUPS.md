# pg-paren-dispatch — post-merge follow-ups

Starmap completed 2026-04-21. Final state: 223 [x], 3 [~] (documented partials), 2 [ ] (accepted coverage debt in §3.1).

The following are post-merge items surfaced by Codex review at each phase. None blocked merging; all are tracked here for future work.

## Tracked parser bugs (oracle-discovered)

Discovered by Phase 2/3 oracle + fuzz. See `PAREN_KNOWN_BUGS.md` for full entries.

- **PAREN-KB-1** — `(T JOIN U)` accepted without ON/USING. Fix in `parseJoinedTable` / `parseJoinQual`. Pinned by §2.7 `t.Skip` + fuzz known-mismatches.
- **PAREN-KB-2** — `Parse` returns only first RawStmt for multi-statement input like `SELECT * FROM (SELECT 1) SELECT 1`. Fix in `parser.go:Parse`. Top-level statement-list issue, out of `parenBeginsSubquery` scope.
- **PAREN-KB-3** — `LATERAL ()` accepted with empty body; PG rejects. Fix in `parseSelectWithParens` nil-body check. Pinned by §3.2 `t.Skip`.

## CI hardening follow-ups

Identified by Codex Phase 5 review.

- **paren-oracle.yml path filter** — currently narrow (`pg/parser/**`). Extend to include `go.mod` / `go.sum` and any shared harness paths to close the blind spot.
- **Nightly fuzz cron concurrency** — add `concurrency: cancel-in-progress: true` to paren-fuzz-nightly so back-to-back cron runs don't stack.
- **Drift-detection smoke test** — the `paren_audit_lint_test.go` drift-detection path isn't exercised in commit transcript (only the proof-fence path is). Add a documented smoke test for the code-vs-audit drift branch.

## Audit debt

- **PAREN_AUDIT_DRIFT_BASELINE.txt 43-pair allowlist** — retained to preserve lint-pass today given function-name drift since Phase 1 scope-lock. Next audit-rebase pass should either (a) update each baseline entry into a proper audit row, or (b) remove entries that no longer correspond to a live dispatch site. Do NOT grow this file; treat as a one-time acknowledgment.

## §3.1 coverage debt (2 scenarios)

- `IN` in CHECK constraint: `CREATE TABLE t (x int CHECK (x IN (1,2,3)))`
- `IN` inside subquery's WHERE: `SELECT (SELECT 1 WHERE x IN (SELECT y FROM t))`

Both parse through the same `parseInExpr` / `parseAExpr(0)` code path already covered by existing §3.1 oracle tests — they're context-variant probes, not distinct dispatch logic. Acceptable to close the starmap without them per the plan's §5.3 bar; track here for thoroughness if a future pass wants full scenario count.
