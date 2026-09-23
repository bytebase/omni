# docs/

`docs/plans/`, `docs/specs/`, `docs/scenarios/`, and `docs/migration/` are dated records of how an engine or layer was built. They keep the paths, batch numbers, and terminology that were accurate when written.

- Do not treat a plan or spec as current architecture guidance. Verify every path it names against the working tree, and follow the root [AGENTS.md](../AGENTS.md) for today's rules.
- `engine-capability-guide.md` and `PARSER-DEFENSE-MATRIX.md` are the living references for building a new engine or capability layer. Update them when a lesson changes, not per PR.
- A change under `docs/` runs every package in CI because some tests read their corpus from here; keep corpus edits in their own commit.
