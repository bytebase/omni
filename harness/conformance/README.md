# Conformance harness

Scoreboard generator for omni engine parsers against pinned real engines. It
harvests an engine's own test corpus (pre-labeled statements), runs every
statement through the omni parser, adjudicates divergences against a live
engine container, and emits a committed deterministic scoreboard. The bar is
**GAP = 0 on the upstream lane**: omni must never falsely reject SQL the
engine accepts at the pinned version. OVER (omni accepts what the engine
rejects) is triaged, not chased. INDETERMINATE is surfaced in the scoreboard
and queued for triage — never silently dropped.

Design doc: `CLAUDE_BB/plans/2026-07-07-omni-full-compatibility-conformance-design.md`
(workspace-level, outside this repo).

## Quickstart

Label-only sweep (no container):

```sh
./fetch_corpus.sh                                    # sparse-clones pingcap/tidb v8.5.5 into corpus/ (gitignored)
go run . -omni-sha "$(git -C ../.. rev-parse HEAD)"
```

Adjudicated sweep (divergences arbitrated by a live TiDB):

```sh
./start_tidb.sh                                      # pingcap/tidb:v8.5.5 on :14001
# run the two export lines it prints (TIDB_DSN, TIDB_CONTAINER_DIGEST)
go run . -adjudicate -omni-sha "$(git -C ../.. rev-parse HEAD)"
docker rm -f tidb-conformance                        # cleanup
```

The module `replace`s omni to `../..`, so the parser under test is always the
working tree; `-omni-sha` is provenance recorded in the output, not what
selects the code.

## Output contract

- `scoreboards/<engine>.md` — **committed.** Deterministic (stable ordering,
  stable cluster keys): regenerating without an omni/corpus/container change
  is byte-identical, and the file's git history is the progress report.
- `out/<engine>.jsonl` — **gitignored artifact.** One meta line first (engine
  version, omni SHA, corpus tag, container digest, classifier version), then
  one row per statement with full provenance: source path/line, SQL, both
  verdicts, raw engine error code+message, classification. Reclassification
  is offline — rerun against the JSONL, no container re-run needed.

Clusters (family + normalized divergence key) are the unit of work; statement
counts are coverage context only.

## Verdicts and classification

Each row can carry two ground-truth signals: the upstream label (`expected`,
harvested from the corpus) and the container verdict (`engine_verdict`, from
adjudication). The container verdict takes precedence; when both are present
and disagree, the row is INDETERMINATE `label_container_disagree` (extraction
bug / stale label / context loss) rather than silently resolved.

- **GAP** — engine accepts, omni rejects. The hard-bar metric.
- **OVER** — engine rejects, omni accepts. Triage: structural over-accepts
  (wrong-AST smoke) are fixed on sight; pure leniency is parked and documented.
- **INDETERMINATE** — conflicting or unknown signals: label/container
  disagreement, infra errors, statements unsafe to execute. Manual queue.

`classifier_version` (meta line + scoreboard) bumps whenever classification
pipeline semantics change — `classify()`, cluster-key normalization — so
scoreboard diffs across classifier changes are not mistaken for parser
movement.

## The two diff axes (never move both at once)

- **Omni moves, container/corpus pinned:** any AGREE → GAP/OVER flip is an
  omni regression and **blocks**. Run the sweep before every omni→bytebase
  dep bump.
- **Container/corpus tag bumps, omni pinned (re-baseline):** all flips are
  engine-version deltas — a ranked burn-down list; nothing blocks.

## Adding an engine

Per-engine seam, three pieces:

1. **Loader** — corpus files → `CorpusEntry` (SQL, expected verdict,
   provenance, skip reason). See `extract_tidb.go`.
2. **Verdict fn** — `omniXVerdict(sql)`, panic-safe: a panic is a reject with
   a `PANIC:` marker, and itself a finding. See `omni_tidb.go`.
3. **Adjudicator + exec-error classifier** — container probe mapping engine
   errors to parse-accept/parse-reject. See `adjudicate_tidb.go`.

Hard-won rules:

- **Enumerate the engine's parse-abort error-code space from its source
  before trusting the classifier.** TiDB needed a 26-code lattice (yacc
  1064/1149, grammar-action aborts, ast-validator codes), not "1064 only".
  MariaDB codes semantic errors like 1911 as parsed. StarRocks discriminates
  parse-vs-analyzing on the *message* ("Getting syntax error"), not the code.
  Codes that also occur at runtime fail closed to INDETERMINATE.
- **Parameterize the exec-error classifier instead of cloning
  `applyContainerVerdict`.** Extract the shared shape when the second engine
  lands — MariaDB shares the MySQL driver, and ~70% of `adjudicate_tidb.go`
  is shareable.
- **Syntactic corpus matchers are name-only.** `isTestCaseSlice` matches the
  type identifier by name with no type resolution — re-verify the
  single-type assumption per corpus before reusing the pattern.

Engine-generic concerns that carry over as-is:

- **Unsafe-statement predicate:** parser corpora literally contain
  `SHUTDOWN` / `KILL` / `RESTART`; these must never reach the shared oracle
  (INDETERMINATE `unsafe_to_adjudicate`).
- **Fresh connection per row:** session state (`USE`, `SET sql_mode`) changes
  how *later* statements parse, so the pool is pinned to zero idle
  connections and every row gets a fresh session.

## Corpus lanes

`lane ∈ {upstream, generated}`. `upstream` = statements harvested from the
engine's own test suite — the headline numbers. `generated` = future variant
lane (mutations, sweeps) that never dilutes the headline. The scoreboard
headline is computed on the upstream lane only.
