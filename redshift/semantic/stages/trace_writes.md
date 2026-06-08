# Stage: trace_writes

Goal:
- Trace system catalog mutations caused by the mapped handler chain for one family.

Inputs:
- `pg/semantic/families/<family>.json`
- `../postgres`

Required work:
- Follow the handler chain through `../postgres/src/backend/...`.
- Identify catalog writes such as `CatalogTupleInsert`, `CatalogTupleUpdate`, `CatalogTupleDelete`, `recordDependencyOn`.
- Normalize affected catalog names and preserve trigger functions.

Output contract:
- Emit only JSON with keys: `family_id`, `writes`, `evidence`.
- Each write entry must include `catalog`, `operation`, `trigger`, and optional `condition`.

Restrictions:
- Use only evidence from `../postgres`.
