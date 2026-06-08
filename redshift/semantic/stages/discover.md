# Stage: discover

Goal:
- Refine scenario coverage for one seeded family without creating new family IDs.

Inputs:
- `pg/semantic/families.seed.json`
- `pg/semantic/STATE.json`
- `pg/semantic/families/<family>.json`
- `../postgres`

Required work:
- Read the selected family scope from the seed file.
- Inspect `../postgres` for top-level semantic variants of that family.
- Add or refine `scenarios` inside `pg/semantic/families/<family>.json`.
- Record evidence paths under `evidence`.

Output contract:
- Emit only JSON with keys: `family_id`, `scenarios`, `evidence`, `scenario_count`.
- Do not emit prose outside the JSON payload.

Restrictions:
- Do not add new families.
- Do not use Omni `pg/catalog` code as evidence.
