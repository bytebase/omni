# Cassandra parser test data

## `compat.tsv`: Cassandra's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/Cql.g4`, then its near-misses: the statement with one token
dropped, the end cut off, a stray word inserted, or a keyword swapped. Each
row has omni's verdict and the pinned Apache Cassandra 5.0.9 server's
(`../../testdata/pins.tsv`). `TestCompatibility` (`../compat_test.go`) reads
the file through `enginetest/compat` and starts no server. Recorded against
`cassandra:5.0.9` (`OMNI_RECORD=1`) and verified byte-identical on a second
fresh container (`OMNI_RECORD=verify`).

### Before vs. after upgrading from `grammar` to `engine`

| Rows | Total | Agree | Only judge accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| **Before (`grammar` + 13 overrides)** | | | | |
| One per alternative | 904 | 904 | 0 | 0 |
| Near-misses | 2,037 | 2,036 | 1 | 0 |
| All | 2,941 | 2,940 | 1 | 0 |
| **After (`engine`, `cassandra:5.0.9`)** | | | | |
| One per alternative | 904 | 904 | 0 | 0 |
| Near-misses | 2,037 | 2,037 | 0 | 0 |
| All | 2,941 | 2,941 | 0 | 0 |

- **The grammar** is Apache Cassandra's own, `src/antlr/Cql.g`, `Parser.g` and
  `Lexer.g` at `cassandra-5.0.9` (`grammar` role in `../../testdata/pins.tsv`),
  beside it byte for byte under the Apache License (`Cassandra.LICENSE`,
  `Cassandra.NOTICE`); `TestGrammarVendorDigest` holds them to their SHA-256.
- **`Cql.g4` is converted from them** by `../cmd/antlr3to4`, since they are
  ANTLR 3. Each rule is under the file and line it came from, and what the
  conversion lost is noted beside it: actions, predicates, and alternatives an
  action always refuses (`COMPACT STORAGE` on a view, a bind marker as a
  keyspace, table, index or role name, a quoted user name).
  `TestCqlG4IsTheConversion` fails while `Cql.g4` is anything else. After a pin
  bump replaces the server image and the three `.g` files:

  ```bash
  go run ./cassandra/parser/cmd/antlr3to4
  OMNI_RECORD=1 go test -timeout 20m ./cassandra/parser/ -run '^TestCompatibility$'
  OMNI_RECORD=verify go test -timeout 20m ./cassandra/parser/ -run '^TestCompatibility$'
  ```
- **Token forms** are in `../../testdata/tokens.tsv`, cited to `Lexer.g`.
  Three rows pick what rows write: `cident` is a name, since the rule lists
  `""` first and every column would carry it, and `COLUMNFAMILY | TABLE` and
  `COLUMNFAMILIES | TABLES` are written by their current spelling. Rows of
  the form `(every site)` name optional clauses that get a row at every
  alternative naming them, not only at the one their own row reaches: IF
  EXISTS in `ALTER TABLE` and `ALTER TYPE` (on the object and on a column or
  field), `STATIC` and `MASKED WITH` on an added column, the `PERMISSION`
  noise word, and a selector's field or element. They added 38 rows, all
  agreeing with `cassandra:5.0.9` (2,979 / 2,979).

After a parser fix, commit the rows that now agree:

```bash
go test ./cassandra/parser/ -run '^TestCompatibility$' -update-compat
```

## How the engine judge works

`compatJudge` (`../compat_test.go`) starts the pinned `cassandra:5.0.9` image,
verifies both the running container's image digest and
`SELECT release_version FROM system.local`, and asks Cassandra about every
statement via a CQL native-protocol `PREPARE` frame
(`gocql.Session.StatementMetadata`). `PREPARE` runs Cassandra's ANTLR 3 lexer,
parser (`Cql.g`, `Parser.g`, `Lexer.g`), and `Parser.g`'s inline Java actions
without executing the statement, so no row mutates server state.

`cassandraVerdict` classifies the response without guessing:

1. **`err == nil`:** `PREPARE` parsed and prepared the statement -> **`accept`**.
2. **Not a `gocql.RequestError`** (transport error, timeout, closed connection):
   no verdict -> **fail the recording**.
3. **`ErrCodeSyntax` (`0x2000`):** ANTLR lexer/parser error, `SyntaxException`,
   or `addRecognitionError` in `Parser.g` -> **`reject`**.
4. **`ErrCodeInvalid` (`0x2200`):** Cassandra uses `InvalidRequestException`
   both inside `Parser.g` inline Java actions during parsing (e.g.
   `CreateTableStatement.Raw.addColumn`, `extendClusteringOrder`,
   `FunctionResource.functionFromCql`) and in post-parse `prepare()` /
   `validate()` / schema-lookup checks (e.g. `"No keyspace has been
   specified"`, `"ALTER [ROLE|USER] can't be empty"`, `"Global TTL on the BATCH
   statement is not supported."`). To separate the two without matching message
   strings, the judge sends a tail probe `stmt + "\n; )"` (the leading newline
   closes any `--` or `//` comment; `;` ends `cqlStatement` in `query :
   cqlStatement (';')* EOF`; `)` cannot appear before `EOF`):
   - Probe returns **`ErrCodeSyntax` (`0x2000`)**: parsing reached `)` after
     `cqlStatement`, proving `stmt` parsed cleanly to `EOF` and the `0x2200`
     came after parsing -> **`accept`**.
   - Probe returns **`ErrCodeInvalid` (`0x2200`)**: parsing aborted inside a
     `Parser.g` action before reaching `)` -> **`reject`**.
   - Any other probe outcome (success, non-`RequestError`, or another code): no
     verdict -> **fail the recording**.
5. **`ErrCodeUnauthorized` (`0x2100`), `ErrCodeConfig` (`0x2300`),
   `ErrCodeAlreadyExists` (`0x2400`):** probed with `stmt + "\n; )"`; only a
   probe `ErrCodeSyntax` (`0x2000`) proves the statement parsed cleanly before
   the post-parse error -> **`accept`**. Any other probe result fails the
   recording.
6. **Any other Cassandra error code** (`ErrCodeServer`, `ErrCodeProtocol`,
   `ErrCodeOverloaded`, timeouts, `ErrCodeUnprepared`, etc.): no verdict ->
   **fail the recording**.

## What the upgrade settled and surfaced

- **`DESCRIBE SCHEMA zzz` settled (`reject` / `reject`).** Cassandra 5.0.9
  lexes `SCHEMA` as `K_SCHEMA` (`Lexer.g:90`), not `K_KEYSPACE`, and rejects
  `DESCRIBE SCHEMA zzz` with `0x2000` (`line 1:16 mismatched input 'zzz'
  expecting EOF`), confirming omni's parser.
- **All 13 former `grammar-overrides.tsv` rows confirmed `reject`.** Ten return
  `0x2000` directly and three (`CREATE TABLE` duplicate column, `CREATE TABLE`
  duplicate `CLUSTERING ORDER BY`, `CREATE MATERIALIZED VIEW` duplicate
  `CLUSTERING ORDER BY`) throw `0x2200` inside `Parser.g` actions and are
  classified `reject` by the tail probe. `grammar-overrides.tsv` is no longer
  needed and was removed.
- **6 real `cassandra/parser` disagreements surfaced and closed
  (2,941 / 2,941 = 100.00%).** All 6 were unqualified function names in
  `functionResource` (`LIST ALL ON FUNCTION a(...)`, 5 alternatives and 1
  near-miss):
  - `LIST ALL ON FUNCTION a()` (2 rows)
  - `LIST ALL ON FUNCTION a(zzz)` (1 near-miss row; `zzz` parses as a UDT name)
  - `LIST ALL ON FUNCTION a(ASCII)` (2 rows)
  - `LIST ALL ON FUNCTION a(ASCII, ASCII)` (1 row)

  In `Parser.g:1191`, `FunctionResource.functionFromCql($fn.s.keyspace, ...)`
  unconditionally rejects an unqualified function name during parsing with
  `0x2200` (`"In this context function name must be explictly qualified by a
  keyspace"`), regardless of session keyspace. `cassandra/parser` (`auth.go`,
  `parseResource`) now requires an explicit keyspace qualifier on `FUNCTION`
  resources, matching Cassandra 5.0.9 on all 2,941 rows.
