# Stage: trace_reads

Goal:
- Trace mandatory semantic catalog reads for one family.

Inputs:
- `pg/semantic/families/<family>.json`
- `../postgres`

Required work:
- Trace analyzer, resolver, and lookup paths in `../postgres`.
- Identify semantic reads such as `SearchSysCache`, `SearchSysCacheList`, `heap_open`, namespace resolution, relation lookup, type lookup, and operator lookup.
- Record lookup keys and purpose.

Output contract:
- Emit only JSON with keys: `family_id`, `reads`, `evidence`.
- Each read entry must include `catalog`, `lookup_key`, and `purpose`.

Restrictions:
- Do not mix in Omni `pg/catalog` behavior as evidence.
