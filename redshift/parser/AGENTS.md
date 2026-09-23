# redshift/parser

`redshift/parser` is a fork of `pg/parser`. Follow [pg/parser/AGENTS.md](../../pg/parser/AGENTS.md): the FIRST-set discipline, the oracle probe templates, and the backtracking rules apply here unchanged, and this package has its own `first_sets.go` and `first_set_oracle_test.go` that must be kept in step with the parent when a probe or predicate changes.

Redshift-only material stays here because tests read it: `PAREN_AUDIT*.md`, `PARSER_DISPATCH_AUDIT.md`, and `testdata/`. The legacy-parser gap triage lives in `../compat/`.
