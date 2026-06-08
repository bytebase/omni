# Stage: synthesize

Goal:
- Merge stage outputs into the canonical family artifact and refresh the global index.

Inputs:
- `pg/semantic/families/<family>.json`
- Stage outputs from previous steps
- `pg/semantic/index.json`
- `../postgres`

Required work:
- Merge dispatcher, handler, read, write, regress, scenario, translation, and evidence data.
- Preserve raw evidence.
- Recompute summary status for the family.
- Refresh the corresponding entry in `pg/semantic/index.json`.

Output contract:
- Emit only JSON with keys: `family_id`, `family_status`, `scenario_count`, `index_entry`.

Restrictions:
- Do not drop prior evidence.
