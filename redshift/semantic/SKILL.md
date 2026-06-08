# PostgreSQL Semantic Extraction Skill

You are operating the PostgreSQL semantic extraction pipeline.

Source of truth:
- PostgreSQL source tree: `../postgres`
- Regression tests: `../postgres/src/test/regress/sql/`

Out of scope:
- Omni `pg/catalog/*.go` as semantic evidence
- Non-top-level parser groups such as `names`, `types`, `expressions`, `functions`, `xml`, and `json`

Primary inputs:
- `pg/semantic/families.seed.json`
- `pg/semantic/STATE.json`
- `pg/semantic/families/<family>.json`
- `pg/semantic/stages/<stage>.md`

Primary outputs:
- `pg/semantic/families/<family>.json`
- `pg/semantic/index.json`
- `pg/semantic/logs/`

Execution rules:
1. Pick the next actionable family and stage from `pg/semantic/STATE.json`.
2. Only work on top-level statement families in the seed universe.
3. Never invent new family IDs. Discovery may add scenarios only.
4. Every claim must point to evidence in `../postgres`.
5. Emit only structured JSON or deterministic stage markers required by the active stage contract.

Stage order:
1. `discover`
2. `map`
3. `trace_writes`
4. `trace_reads`
5. `map_tests`
6. `plan_translation`
7. `synthesize`
