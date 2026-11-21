# MongoDB parser test data

## `compat.tsv`: the `mongosh` call-shape rows

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/MongoShell.g4`, then its near-misses: the statement with one
token dropped, the end cut off, a stray word inserted, or a keyword swapped.
Each row has omni's verdict and the docs recognizer's. `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat` using pin role
`mongosh` from `../../testdata/pins.tsv` (keeping the `default` `mongo:8.0.38`
pin untouched for `mongo/driver`'s tapes). It starts no server, regenerates the
rows from the grammar, and asks the recognizer again on every run.

| Rows | Total | Agree | Only the docs accept | Only omni accepts | Undocumented |
| --- | --- | --- | --- | --- | --- |
| One per alternative | 569 | 569 | 0 | 0 | 0 |
| Near-misses | 845 | 20 | 0 | 0 | 825 |
| All | 1,414 | 589 | 0 | 0 | 825 |

- **The grammar** is derived directly from MongoDB's official
  `mongodb-js/mongosh` repository (Apache-2.0) at release tag `v2.14.3`
  (`82c228b597c3221de59bdbfaff4c346622780cdd`) by
  `../cmd/mongoshgrammar`. The converter reads the public TypeScript class and
  method definitions across `packages/shell-api/src/` (`ShellApi`,
  `ShellConfig`, `Database`, `Collection`, cursor hierarchy, `Explainable`,
  `Bulk`, `BulkFindOp`, `PlanCache`, `ReplicaSet`, `Shard`, `Streams`,
  `StreamProcessor`, `Mongo`, `KeyVault`, `ClientEncryption`, `Session`) and
  `packages/shell-bson/src/shell-bson.ts` (`ShellBsonBase`, static BSON helper
  methods, and `EJSON`), emitting `../../testdata/MongoShell.g4` with a
  `file:line` citation on every rule and pinning the exact SHA-256 of the 21
  vendor files read under the `mongosh` role in `../../testdata/pins.tsv`.
- **Token forms** are in `../../testdata/tokens.tsv`, covering identifiers,
  strings, signed and unsigned numbers, regex literals, and `$`-prefixed
  MongoDB operators.

To regenerate `MongoShell.g4` after a `mongosh` pin bump, then re-record the
rows:

```bash
go run ./mongo/parser/cmd/mongoshgrammar -cache /tmp/mongosh -tag v2.14.3 -commit 82c228b597c3221de59bdbfaff4c346622780cdd -update-pin
OMNI_RECORD=1 go test ./mongo/parser/ -run '^TestCompatibility$'
```

After a parser fix, commit the rows that now agree:

```bash
go test ./mongo/parser/ -run '^TestCompatibility$' -update-compat
```

## Who judges what

MongoDB's server (`mongod`) does not parse `mongosh` syntax, so `compat.tsv` is
judged by the `docs` recognizer over `MongoShell.g4` (§2.2 and §2.3 of
`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`):

- **A statement the grammar derives is accepted.** All 569 alternatives and all
  20 derived near-misses agree with `mongo/parser` (down from 90 baseline gaps:
  86 alternatives and 4 near-misses).
- **A near-miss the grammar does not derive is `undocumented`, not `reject`.**
  Because `MongoShell.g4` is a call-shape subset of `mongosh` rather than a
  complete ECMAScript grammar, `undocumented` rows pass whatever omni says so
  omni is never pushed to reject valid JavaScript or arbitrary user method
  calls outside the documented call shapes.
- **Gaps closed from the initial baseline (90 → 0):**
  1. Direct `@directShellCommand` syntax without parentheses (`use <db>`,
     `exit [code]`, `quit [code]`, `it`, `cls`) and function-call `show(...)`.
  2. Extended `show` forms from `Mongo.show()`: `show log <name>`,
     `show automationNotices`, and `show nonGenuineMongoDBCheck`.
  3. Built-in shell namespace calls `config.<method>(...)` (`ShellConfig`) and
     `EJSON.<method>(...)` (`ShellBson`).
  4. Standalone `db.getCollection(...)` and method/collection chaining on
     `db.<method>(...)` results (`db.aggregate(...).<cursor>()`,
     `db.watch().<changeStream>()`, `db.getSiblingDB(...).<dbMethod>()`,
     `db.getSiblingDB(...).<coll>.<collMethod>()`).
  5. Static helper calls on BSON keywords (`ObjectId.createFromHexString(...)`,
     `ObjectId.isValid(...)`, `Timestamp.fromInt(...)`, etc.), which
     previously failed because `parseValue()` required `(` immediately after
     BSON keyword tokens.

## What the judge leaves out

- **Arbitrary JavaScript syntax:** variable declarations (`var`/`let`/`const`),
  control flow (`for`/`if`/`while`), arrow/function declarations, and binary
  operators are intentionally outside `MongoShell.g4`'s call-shape scope.
- **`new` constructor expressions:** `mongo/parser` rejects the `new` keyword
  by design (pinned by `mongo/parsertest/TestNewKeywordError` and
  `TestNewKeywordErrorMessage`), so `MongoShell.g4` emits call-form helper and
  connection invocations without `new`.
- **Runtime arity and argument types:** `mongosh` validates method argument
  counts and types at runtime in TypeScript (`assertArgsDefinedType`), not
  during JS parsing, so `MongoShell.g4` models argument lists uniformly as
  `argList?`.
