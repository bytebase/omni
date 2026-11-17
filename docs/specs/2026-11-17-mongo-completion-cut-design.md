# The Mongo completion cut

2026-11-17 · Design, not implemented

The completion cut in [`docs/specs/2026-11-16-completion-cut-design.md`](2026-11-16-completion-cut-design.md)
is "omni emits parser-native token/rule/scope/intent signals; Bytebase resolves live metadata into
product `Candidate` values". MySQL landed it in #1291, and MariaDB, OceanBase, and TiDB follow as
copies of that reader. Mongo does not, and this note says why it is its own cut and what it needs.

## Where Mongo stands

`mongo/completion` is one package whose `Complete(input, cursorOffset, *catalog.Catalog) []Candidate`
does three jobs at once:

- it lexes the mongosh input up to the caret and classifies the position into one of twelve
  contexts — top level, `db.|`, `db.<collection>.|`, `db[|`, inside call arguments, a document key,
  a query operator, an aggregation stage, a cursor chain, a `show` target, `rs.|`, `sh.|` — in
  `context.go`;
- it returns the static vocabulary each context accepts — collection methods, cursor methods,
  aggregation stages, query operators, BSON helpers, `show` targets, `db`/`rs`/`sh` methods — from
  `candidates.go`;
- it reads collection names from `catalog.Catalog.Collections()` in the two contexts that take one
  (`db.|` and `db[|`), and wraps everything in omni's own `Candidate`/`CandidateType`.

Bytebase's `backend/plugin/parser/mongodb/completion.go` (159 lines) builds that catalog from
`GetDatabaseMetadataFunc` on every request, calls `Complete`, formats collection candidates for
bracket notation, and converts to `base.Candidate`. That is the same boundary violation TiDB had:
the caller's metadata becomes an omni catalog projection, and omni owns the resolution step the
cut assigns to Bytebase.

## Why it is not a copy of the MySQL reader

The MySQL reader is a token walk over a SQL statement: it reads query blocks, `FROM` lists, CTE
definitions, select lists, and the SQL object classes a dotted name can address. Mongo has none of
those. Its surface is a shell language of paths, method chains, document literals, aggregation
pipelines, and query operators, and the object classes are shell positions, not SQL kinds:
`collection` is a position under `db.`, not a relation in a query block; a document key is a field
of whatever collection the chain names; `$group` is a stage of a pipeline array.

So the shared shape cannot be reached by renaming. The signal vocabulary is twelve contexts, a
dotted path, and whether the caret is inside a document, an argument list, or a pipeline — not
`Scope`/`CTEs`/`SelectAliases`. Writing MySQL's struct here would mean fields Mongo never sets and
a promotion into `omni/completion` that has to be redone once a second shell-like engine appears.

## What it needs

1. **A signal struct, `mongo/completion.Signals(input string, cursorOffset int)`** — no catalog.
   - `Context` (the twelve above), so the caller knows which vocabulary and which metadata lookup
     applies.
   - `Prefix`: the identifier, string, or operator the caret is in or ends, read from the lexer.
   - `Path`: the dot-separated and bracket-indexed names typed before the caret, flat, as written
     (`db.users` → `{"db", "users"}`), so the caller can tell `db.|` from `db.users.|`.
   - `InDocument` / `InPipeline`: whether the caret is inside `{ ... }` or `[ ... ]` at that path,
     which decides between a document key, a query operator, and an aggregation stage.
   - `Rules`: the rule names the position accepts (`collection`, `db_method`, `collection_method`,
     `cursor_method`, `agg_stage`, `query_operator`, `bson_helper`, `show_target`, `rs_method`,
     `sh_method`, `document_key`), the way MySQL's `CandidateSet` names its rules. The vocabulary
     behind each rule stays in omni as grammar, in one table per rule.
2. **Bytebase resolves the two metadata rules**: `collection` through the `GetDatabaseMetadataFunc`
   seam it already calls, and `document_key` through the field names its metadata carries for the
   collection the path names. The shell vocabulary is grammar, so it comes back with the signal the
   way SQL keywords do, and the adapter only decides which rules to ask the store about.
3. **The package splits**: `Signals`, the vocabulary tables, and `context.go` stay in omni;
   `Candidate`, `CandidateType`, prefix filtering, and the catalog parameter go, exactly as
   `tidb/completion` goes in the TiDB PR. `mongo/catalog` stays — review, the query path, and the
   driver use it.
4. **A corpus per engine**: an omni corpus of shell positions × scripts (paths, chained calls,
   document literals with nested keys, pipelines, operator lists, `show` targets, bracket notation,
   unterminated strings) and a replay of Bytebase's `test-data/test_completion.yaml`. That file has
   11 cases today, so the parity run is per caret position over those 11 scripts, not the 200-case
   SQL corpus: candidate lists are excluded because the fixture's metadata is small, and what is
   compared is the context, path, prefix, and rule set.
5. **The design record gets one sentence**: the cut's sentence about promotion into `omni/completion`
   ("promote it once two engines agree field for field") assumes SQL object classes. Mongo's signals
   would promote into a shell vocabulary, not that contract, and the document should say so before
   someone tries.

## Cost

`mongo/completion` is 625 non-test lines and 581 test lines today. The split above keeps roughly
350 lines (contexts, vocabulary, prefix and path reading) and moves the rest (`Candidate`,
filtering, catalog use) into Bytebase's 159-line adapter, which grows to resolve two rules. The
corpus is the new work: there is no SQL-shaped corpus to copy, so the shell positions have to be
enumerated first.

That is a cut of its own, not a session's follow-up to the MySQL family, which is why MongoDB is
left with this statement instead of a port.
