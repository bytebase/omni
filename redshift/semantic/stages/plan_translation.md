# Stage: plan_translation

Goal:
- Produce PG C translation targets for one family.

Inputs:
- `pg/semantic/families/<family>.json`
- `../postgres`

Required work:
- Build an ordered list of PG C functions that would need translation later.
- Classify each target as `dispatcher`, `analyzer`, `resolver`, `catalog_writer`, or `validator`.
- Keep the targets rooted in `../postgres`.

Output contract:
- Emit only JSON with keys: `family_id`, `translation_targets`, `evidence`.

Restrictions:
- Translation targets must point to PG C functions, not Omni files.
