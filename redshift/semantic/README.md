# PostgreSQL Semantic Pipeline

This directory holds the Codex-based semantic extraction pipeline for PostgreSQL.

Scope:
- Source of truth is `../postgres`
- Only top-level statement families are in scope
- Focus is on statements with stable system catalog semantics
- Low-level parser groups such as `names`, `types`, `expressions`, and `functions` are out of scope

Primary artifacts:
- `families.seed.json`: controlled family universe
- `schema/`: JSON schemas for machine-written artifacts
- `STATE.json`: stage state for family/scenario execution
- `families/`: one semantic inventory file per family
- `index.json`: global summary index

Family categories:
- `ddl`
- `dml`
- `query`
- `set`
- `object_control`

The pipeline is designed for staged Codex execution:
1. `discover`
2. `map`
3. `trace_writes`
4. `trace_reads`
5. `map_tests`
6. `plan_translation`
7. `synthesize`
