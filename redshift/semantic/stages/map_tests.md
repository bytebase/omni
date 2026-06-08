# Stage: map_tests

Goal:
- Map semantic scenarios to PostgreSQL regression tests.

Inputs:
- `pg/semantic/families/<family>.json`
- `../postgres/src/test/regress/sql/`

Required work:
- Search only under `../postgres/src/test/regress/sql/`.
- Map one or more regress SQL files to each scenario where coverage exists.
- Mark scenarios that lack regress coverage for manual review.

Output contract:
- Emit only JSON with keys: `family_id`, `regress_files`, `scenario_test_map`, `evidence`.

Restrictions:
- Do not search outside `../postgres/src/test/regress/sql/`.
