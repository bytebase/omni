# Stage: map

Goal:
- Map each top-level statement family to PostgreSQL dispatcher and primary handler chains.

Inputs:
- `pg/semantic/families/<family>.json`
- `../postgres`

Required work:
- Find dispatch roots such as `ProcessUtility`, `standard_ProcessUtility`, `parse_analyze_fixedparams`, and `transform*Stmt`.
- Trace from the AST node to the primary handler chain in `../postgres/src/backend/...`.
- Update `dispatcher`, `handlers`, and `evidence`.

Output contract:
- Emit only JSON with keys: `family_id`, `dispatcher`, `handlers`, `evidence`.
- Use path-like evidence rooted in `../postgres`.

Restrictions:
- Do not cite Omni code as semantic source.
- Stay at top-level statement semantics.
