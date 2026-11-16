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
- A materialized-view definition renders its indexes.
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
  addition.
- A named `NOT NULL` constraint is recorded as sync reports it, a check
  constraint on the column. `ALTER TABLE ... MODIFY (... CONSTRAINT name NOT
  NULL)` does not carry the name through the parser, so a migration replay keeps
  the nullability but not that name.
- The deparse is byte-for-byte the Bytebase renderer. The fixtures that pin it
  are the Bytebase fixtures.

## What the deparse cannot restore

- Routine and package bodies are written as the server stored them. They are
  text, not re-rendered, so a body that depends on editioning or on a session
  setting is restored only as text.
- Storage clauses, tablespace placement, compression, parallelism and other
  physical attributes are not in the metadata and are not emitted.
- Editioning views are not a distinct metadata kind; they render as ordinary
  views, without the editioning property.
- Object tables and XMLType tables are skipped: their DDL is their type.
- View column aliases are not emitted deliberately. `ALL_VIEWS.TEXT` does not
  carry the alias list, and writing the synced column names can produce fewer
  names than output columns for a view with an invisible column (ORA-01730).
- Collation is not emitted. Oracle's `USING_NLS_COMP` is the column default and
  writing it back fails on sessions with a different string-size setting.
- Packages, synonyms, database links and public synonyms are not rendered.

## Tests

- `catalog_test.go` replays the Bytebase definition, metadata and migration
  fixtures, checks definition round trips, and replays every generated migration
  through the catalog.
- `regressions_test.go` pins sequence alterations, `DEFAULT ON NULL`,
  cross-schema foreign keys, national character lengths, view triggers, bitmap
  indexes, case-sensitive key columns and constraint renames.
- `oracle_test.go` executes every definition and a quoted, case-sensitive
  fixture on the pinned Oracle in `oracle/internal/testcontainer` and reads the
  objects back. It is a tape test; record it with `OMNI_RECORD=1`.
