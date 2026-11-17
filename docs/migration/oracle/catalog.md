# Oracle catalog

`oracle/catalog` holds the schema-conversion layer `oracle/driver` needs for a
`Dump`: metadata loading, deparse, diff and migration. It was ported from
Bytebase's `backend/plugin/schema/oracle` package, whose fixtures are now
`oracle/catalog/testdata`.

## API

- `New`, `LoadSQL`, `LoadSDL`, `Exec` and `LoadMetadata` build catalog state.
- `DatabaseDefinition` renders a synced snapshot as Oracle DDL.
- `Diff` compares two catalogs; `GenerateMigration` turns the diff into a
  `MigrationPlan` whose `SQL()` is the script.
- `Catalog.Metadata()` returns the snapshot the catalog holds.

Oracle has one schema per user, so a snapshot normally carries one schema. The
catalog stores the snapshot rather than re-executing every object definition;
`Exec` applies the DDL statements it understands, which is what the migration
replay tests exercise.

## Deliberate differences from `schema/oracle`

- Expressions taken from statement text keep it: a CHECK or virtual-column
  expression reads as written, with its spacing and quotes. The Bytebase
  extractor rebuilt these from the AST and lost both
  (`STATUSIN('A','B')` instead of `STATUS IN ('A', 'B')`), which broke a
  definition's round trip.
- `CREATE SEQUENCE` records `MINVALUE`, `MAXVALUE`, `CACHE` and `CYCLE` as well
  as `START WITH` and `INCREMENT BY`. The Bytebase extractor dropped the rest,
  so a deparse could not restore them. A changed sequence renders
  `ALTER SEQUENCE` instead of being silently skipped.
- `NVARCHAR2` and `NCHAR` keep the character length they declare. The Bytebase
  extractor doubled it, while sync reads `CHAR_LENGTH`; SQL-derived metadata and
  a synced snapshot therefore disagreed.
- Cross-schema foreign keys keep the referenced owner, so deparse and migration
  qualify the reference.
- A primary or unique constraint recreated under a new name is rendered as
  `RENAME CONSTRAINT`, not drop and create, so its incoming foreign keys stay
  attached.
- `DEFAULT ON NULL` is kept in CREATE/ADD and changed through MODIFY.
- Migration statements keep a bitmap index's `BITMAP` modifier and an
  invisible index's `INVISIBLE`, quote a case-sensitive key column, and recreate
  a changed view's triggers.
- Identity columns render `GENERATED ... AS IDENTITY` with their seed and
  increment, in a definition and in CREATE/ADD.
- A definition with mutually referencing tables creates the tables first and
  adds their foreign keys afterwards.
- A migration creates a key added to an existing table before a new table's
  foreign key that references it, and skips `CREATE USER` for the unnamed
  current schema.
- A changed view's comments use `COMMENT ON TABLE`, which is Oracle's form for
  a view, and a changed procedure is always `CREATE OR REPLACE`.
- Renaming a column rewrites identifier tokens in check expressions, not string
  literals.
- Identity options are read from the clause text, and identity alterations render
  through MODIFY.
- Views record the relations their query names, so a generated definition or
  migration orders a view after its dependencies.
- A view comment change alone renders the comment without dropping and
  recreating the view.
- A statement that names a second owner is refused instead of merging its
  objects into the schema already loaded.
- A foreign key keeps its `ON DELETE` action in a definition.
- Dropping the unnamed schema still drops its objects, and only suppresses
  `DROP USER`.
- A quoted column in a function-based index keeps its quotes.
- An existing table's foreign key waits for a key this migration adds, and a
  whole-schema drop detaches its foreign keys before dropping tables.
- Renaming a column updates both sides of every foreign key that names it.
- Executing a statement keeps schema metadata the extractor does not own, such
  as synonyms.
- Dropping a column removes the indexes, checks and foreign keys that name it,
  and renaming a table updates incoming foreign keys.
- A sequence definition renders `MINVALUE` and `CACHE`, and a quoted
  user-defined type keeps its spelling.
- A changed sequence start uses `RESTART START WITH`, and returning to the
  default increment emits `INCREMENT BY 1`.
- Dropping tables resolves a same-schema foreign key against the owning
  schema's table node.
- An implicit foreign key to an owner outside the snapshot is kept and rendered
  without a column list.
- A virtual column declared without a datatype stays untyped, so Oracle infers
  it from the expression.
- `Exec` applies `ALTER SEQUENCE`, and a named NOT NULL check disappears with
  `MODIFY ... NULL`.
- Owner-qualified dependencies connect to the same table nodes as the objects
  themselves, and a migration qualifies the index name it creates.
- A function declaration is recognized through whitespace between `CREATE`,
  `OR REPLACE` and `FUNCTION`.
- Expression comparison recognizes Oracle's alternative-quoted literals.
- `Exec` applies `DROP IDENTITY`, and `NOCACHE` keeps its explicit state through
  extraction, definitions and migrations.
- A materialized-view definition renders its indexes, except Oracle's
  `SYS_C_SNAP$` container index, which the view recreates and whose object id
  differs between snapshots.
- Replacing a key detaches and restores the unchanged incoming foreign keys.
- Catalog snapshots are copied with a standard-library helper, so production
  catalog code imports no protobuf runtime.
- Only the incoming foreign keys detached by a key replacement are restored.
- Renaming a column rewrites function-based index and virtual-column
  expressions.
- Dropping a named NOT NULL check makes its column nullable when no other check
  still requires it.
- `ALTER TABLE` validates the statement's owner, and renaming a table leaves
  foreign keys owned by another schema alone.
- An identity column renders `ON NULL` when its metadata requests it.
- A table definition renders its table and column comments.
- A constraint column that is reserved or not a simple identifier is quoted in
  a migration.
- `DROP PRIMARY KEY`/`DROP UNIQUE` without a constraint name removes the
  matching key, and identity `ON NULL` changes render through MODIFY.
- View dependencies exclude names introduced by a CTE while keeping the tables
  its query reads; quoted CTE names stay case-sensitive.
- A non-identity `DEFAULT ON NULL` change does not emit identity DDL.
- `DROP` validates the statement's owner before removing an object.
- `CREATE TABLE IF NOT EXISTS` leaves an existing table unchanged.
- Function-based indexes are emitted after the functions they call.
- A self-referencing foreign key's referenced columns follow a column rename, and
  a renamed identifier that needs quotes is written quoted.
- `DROP UNIQUE (columns)` removes only the matching unique key.
- View definitions keep `WITH READ ONLY` and `WITH CHECK OPTION`.
- A sequence definition keeps explicit zero starts and bounds.
- An added foreign key waits for a referenced column rebuild as well as a key
  addition, and a key replacement plus a column rebuild detaches and restores
  an unchanged foreign key once.
- An identity column's backing-sequence default is not treated as an ordinary
  default during a MODIFY.
- A column rename leaves an external foreign key's referenced columns alone.
- A parenthesized `DROP (A, B)` drops every listed column.
- `CREATE TABLE AS SELECT` is rejected instead of installing an empty table.
- Renaming a table rewrites the target relation of its triggers.
- SQL identity clauses keep `ON NULL`.
- Dropping a column removes only local foreign-key references.
- CTE exclusion is scoped to its query block.
- Trigger comparison keeps line-comment boundaries.
- Synced and extracted constraint indexes compare by their logical attributes.
- A self-referencing foreign key is detached around a key replacement.
- A renamed key column stays a raw constraint name, while a quoted reference is
  matched case-sensitively.
- Temporary tables are rejected rather than installed as permanent.
- `CREATE OR REPLACE TRIGGER` replaces the previous target, and a cascading
  table drop removes incoming local foreign keys.
- A table-cycle fallback still orders views by their dependencies.
- Sequence owners are validated before creation and alteration.
- Function declarations are normalized by tokenized prefixes, including
  `CREATE`/`OR REPLACE` on separate lines and `EDITIONABLE`.
- Triggers extracted from SQL carry their event and timing, and trigger bodies
  compare without schema qualification or `OR REPLACE` spelling differences.
- Initial migrations defer function-based indexes until their functions exist.
- An identity column that becomes ordinary applies its new default after
  `DROP IDENTITY`.
- `ALTER TRIGGER ... ENABLE/DISABLE` is kept in the stored declaration.
- Dropping an independent CHECK does not clear a separate column NOT NULL.
- Dependent virtual columns are dropped before the columns they reference.
- A self-referencing foreign key the diff removes is not restored.
- The migration fallback targets materialized-view nodes and skips comment-only
  views and materialized views.
- Trigger comparison normalizes only the declaration header, not literals in
  the body.
- Materialized views keep their `BUILD`/`REFRESH` policy in the definition.
- An independent CHECK does not clear a column's NOT NULL when it is dropped.
- A retained virtual column is altered before a column it used to reference is
  dropped.
- An index's target table owner is validated, inline `DISABLE` triggers can be
  re-enabled, and non-table triggers are rejected.
- Constraint column names are matched exactly during renames and drops, and
  `MODIFY ... NULL` leaves an independent CHECK in place.
- Diff matches an unnamed synced schema to the same owner-named schema.
- A virtual column's function is ordered before its table in both definitions
  and initial migrations.
- Trigger disable state is tracked by an explicit statement, not body text.
- Routine owners are validated, colliding constraint renames fall back to
  drop/create, and deferrable foreign keys are rejected.
- The cycle fallback still creates a virtual column's function before its table.
- Exec resolves an unnamed synced schema to its database owner, trigger owners
  are validated, and new source columns are added before a virtual-column
  expression change.
- Dropping a primary key restores only the nullability it implied, and
  standalone `RENAME` moves the table.
- `DROP UNIQUE (columns)` matches raw column names exactly.
- A qualified trigger target follows a table rename and keeps its owner.
- `CREATE USER` for a generated schema includes `NO AUTHENTICATION`.
- An inline primary key's explicit NOT NULL survives dropping the key.
- Comment owners are validated, `ALTER INDEX` rename/visibility is applied, and
  constraint enable/disable is rejected.
- A view's called functions are ordered before the view.
- A renamed trigger target is quoted for a case-sensitive table name.
- Dropping a local column leaves a foreign key to an external same-named
  table alone.
- Identifier renames skip Oracle alternative-quoted literals.
- Deferrable primary and unique constraints are rejected.
- `DEFAULT ON NULL` is removed before relaxing NOT NULL, and ordinary-to-identity
  conversions are reported unsupported rather than emitted as invalid DDL.
- Constraint column names compare exactly, whitespace included.
- DEFAULT ON NULL implies NOT NULL during extraction and MODIFY.
- Dependency scanning skips alternative-quoted literals and function calls, and
  dropping or renaming a constraint leaves a standalone index that shares its
  name.
- Primary-key nullability provenance follows a renamed column.
- Incoming-key lookup resolves an unnamed synced schema to its owner.
- A literal default containing `ISEQ$$_` is kept.
- Identifier renames leave function names alone.
- Primary-key and named-NOT-NULL provenance move with a renamed table.
- Function references allow whitespace before the call parenthesis.
- Dropping a table clears its provenance, `ALTER TRIGGER ... RENAME` moves the
  trigger, and `SET UNUSED` is rejected.
- A qualified statement establishes the owner for later validation.
- Function references keep scanning after a non-call occurrence of the name.
- Dependency scanning skips line and block comments.
- An early virtual-column alteration detaches and restores its dependent
  indexes and constraints.
- `DROP IDENTITY` runs before relaxing an identity column's NOT NULL.
- `LoadMetadata`/`SetMetadata` reset provenance maps.
- `DROP PRIMARY KEY ... CASCADE` removes incoming foreign keys.
- A sequence restart is emitted only when the old start value is known.
- A column alteration needed by a new view runs before that view.
- An early virtual-column alteration excludes dependents the diff already drops.
- A materialized view keeps its refresh schedule (`START WITH`/`NEXT`).
- A cascading key drop removes only foreign keys on the dropped key's columns.
- Dropping a column clears its nullability provenance, including reused names.
- Trigger-target scanning accepts whitespace after `ON`.
- Disabled constraints are rejected at creation, but literal text containing
  `DISABLE` does not count.
- A cascading key drop includes self-referencing foreign keys.
- Renaming a disabled trigger moves its stored disable statement.
- Deferrable CHECKs are rejected, and table-wide trigger enable/disable is
  rejected.
- View dependency lookups resolve an empty owner to the view's schema.
- Quoted function names are recognized as virtual-column dependencies.
- Materialized-view refresh schedules render without an explicit method/mode.
- Renaming a view rewrites its attached trigger targets.
- A constraint that reuses an existing index is rejected.
- A definition orders views through the schema being defined when a
  dependency leaves its owner empty, and `CASCADE` on an unnamed `DROP PRIMARY
  KEY` removes only the foreign keys on that key's columns.
- `CREATE SCHEMA AUTHORIZATION` validates the owner against the schema already
  loaded, and a standalone `RENAME` moves a sequence, view or materialized
  view, with the triggers attached to a view, while an object the catalog does
  not hold is rejected.
- A materialized view's refresh schedule renders even when it names no refresh
  method or mode, and a quoted function name counts as a call in a dependency.
- A materialized view's function-based index is written after the functions it
  calls, and a child foreign key detaches and returns around a rebuilt virtual
  column's key.
- A virtual column rebuilt for a source datatype change is dropped and added
  again with its target expression, comment and dependents, whether or not its
  own expression changed.
- A function-based index detaches around an ordinary datatype change, and a
  retained function-based index is rebuilt after this migration replaces its
  function, which Oracle disables with it.
- A procedure whose declaration splits `CREATE` from `OR REPLACE` with any SQL
  whitespace is still altered, not created again.
- A virtual column detaches when either its old or its target expression reads a
  datatype-changing column, and the foreign-key cycle fallback writes a function
  before a view that calls it.
- `ALTER TABLE ... MODIFY CONSTRAINT` enable/disable is rejected like the
  `ENABLE`/`DISABLE` form, and `CREATE SEQUENCE IF NOT EXISTS` leaves an
  existing sequence unchanged.
- A trigger declaration's `DISABLE` clause is removed with whatever whitespace
  precedes it, so a later enable leaves no stale keyword.
- A constraint backed by an index the statement names (`USING INDEX
  [schema.]index`) is rejected: the metadata holds the constraint's index, not
  the separate index it reuses, so a definition would create both.
- A cascading composite-key drop compares the complete referenced column set,
  so a key sharing one column survives.
- Renaming a qualified trigger rewrites the declaration while keeping its
  owner.
- Function-based materialized-view indexes are created after their functions.
- An identity alteration omits `START WITH` unless the seed changed.
- A retained virtual column is dropped and recreated around a datatype change
  to a column it reads.
- A cascading key drop compares the complete referenced column set.
- Function-based materialized-view indexes are restored when the view is
  recreated.
- An identity alteration omits `START WITH` unless the seed changed.
- Retained virtual columns are detached as a batch around a source datatype
  change, including comment-only column changes, and a shared dependent is
  dropped and restored once.
- `DROP ... KEEP INDEX` keeps the backing index as a standalone index.
- The deparse is byte-for-byte the Bytebase renderer. The fixtures that pin it
  are the Bytebase fixtures.

## What this catalog restores, and what it does not

It restores, from a synced snapshot:

- Sequences, with `START WITH`, `INCREMENT BY`, `MINVALUE`, `MAXVALUE`,
  `CACHE` and `CYCLE`, created before the tables whose defaults call them.
  Oracle's dictionary does not carry the declared `START WITH`;
  `oracle/driver`'s sync reads every other option, so a dump of its snapshot
  restarts a sequence at its minimum value.
- Tables, with columns, declared lengths, defaults (including `DEFAULT ON
  NULL`), nullability, identity generation with seed and increment, virtual
  columns, table and column comments, primary keys, unique and check
  constraints, foreign keys with their `ON DELETE` action, and standalone
  indexes (unique, bitmap, invisible, function-based). Tables that reference
  each other are created first and their foreign keys added afterwards.
- Views, with `WITH READ ONLY` and `WITH CHECK OPTION`, their comments and
  triggers, ordered after the relations they read and the functions they call;
  and materialized views, with the `BUILD`/`REFRESH` policy a snapshot carries,
  their comments and their own indexes. Oracle's `SYS_C_SNAP$` container index is not written:
  the view makes Oracle create it again, and its object id differs between
  snapshots. `oracle/driver`'s sync reads `ALL_MVIEWS.QUERY` only, so a dump of
  its snapshot restores the view without the policy.
- Functions, procedures and triggers, whose bodies are written as stored text.
- Changes, as a migration: a key recreated under a new name becomes `RENAME
  CONSTRAINT`, so incoming foreign keys stay attached; a renamed table moves
  its triggers and incoming foreign keys; replaced keys detach and restore the
  foreign keys they carry; and the statements are ordered around the
  dependencies above.

The gaps, each pinned by a `TestDefinitionGaps` subtest that fails when it
closes:

- **Routine and package bodies are stored text, not re-rendered.** A body that
  depends on editioning or a session setting is restored only as that text.
- **Storage clauses, tablespace placement, compression, parallelism and other
  physical attributes** are not in the metadata and are not emitted, including
  a constraint's `USING INDEX` properties.
- **Partitioning** is neither synced nor emitted, so a partitioned table
  restores as a single heap table.
- **External tables** are not carried either: `ORGANIZATION EXTERNAL` and its
  location are dropped, so an external table restores as an empty heap table.
- **Editioning views** are not a distinct metadata kind. `ALL_VIEWS.TEXT`
  carries only the query, so an editioning view renders as an ordinary view,
  and the parser does not accept the `EDITIONING` keyword besides.
- **Object tables and XMLType tables** are skipped: their DDL is their type,
  which is not a synced object. Sync marks them `SkipDump`, and DDL that names
  the type (`CREATE TABLE ... OF`) is refused rather than installed empty.
- **View column aliases** are not emitted deliberately. `ALL_VIEWS.TEXT` does
  not carry the alias list, and writing the synced column names can produce
  fewer names than output columns for a view with an invisible column
  (ORA-01730).
- **Collation** is not emitted. Oracle's `USING_NLS_COMP` is the column default
  and writing it back fails on sessions with a different string-size setting.
- **Sequence options the metadata cannot represent.** A column carries an
  identity's generation, seed and increment; its backing sequence's `CACHE`,
  `CYCLE` and bounds have no field, so a definition drops them. A standalone
  sequence's declared `START WITH` is not in the dictionary at all, so a dump
  restarts it at `MINVALUE`.
- **Packages, synonyms, public synonyms and database links** are kept in the
  snapshot but not rendered, so a restored schema is missing them.
- **A named `NOT NULL` constraint's name.** `ALTER TABLE ... MODIFY (...
  CONSTRAINT name NOT NULL)` does not carry the name through the parser, so a
  replay keeps the nullability and not the name. A named `NOT NULL` constraint
  sync reports is recorded as a check constraint on the column.
- **Privileges, roles and profiles** are not in the metadata, so a restored
  schema grants nothing.

## Tests

- `catalog_test.go` replays the Bytebase definition, metadata and migration
  fixtures, checks definition round trips, and replays every generated migration
  through the catalog.
- `regressions_test.go` pins sequence alterations, `DEFAULT ON NULL`,
  cross-schema foreign keys, national character lengths, view triggers, bitmap
  indexes, case-sensitive key columns and constraint renames.
- `gaps_test.go` pins each gap listed here, one subtest per gap. A subtest
  fails when its gap closes, so the case moves into the fidelity tests instead.
- `oracle_test.go` executes every definition and a quoted, case-sensitive
  fixture on the pinned Oracle in `oracle/internal/testcontainer` and reads the
  objects back. It is a tape test; record it with `OMNI_RECORD=1`.
