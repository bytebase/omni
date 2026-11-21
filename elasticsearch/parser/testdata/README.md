# Elasticsearch parser test data

## `compat.tsv`: the specification grammar's rows

`compat.tsv` holds a single-line REST request for each alternative of each rule
in `../../testdata/ElasticsearchConsole.g4`, then its near-misses (one token
dropped, the end cut off, a stray word inserted, or a keyword swapped), followed
by `suite:console` hand-added rows. Each row records omni's verdict and the
grammar recognizer's verdict. `TestCompatibility` (`../compat_test.go`) reads it
through `enginetest/compat`, regenerates the rows from the grammar, and asks the
recognizer again on every run with no server or network.

| Rows | Total | Agree | Only the grammar accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 1,208 | 1,208 | 0 | 0 |
| Near-misses | 1,825 | 1,825 | 0 | 0 |
| `suite:console` | 20 | 20 | 0 | 0 |
| **All** | **3,053** | **3,053** | **0** | **0** |

- **The grammar source** is Elastic's official
  [`elastic/elasticsearch-specification`](https://github.com/elastic/elasticsearch-specification)
  (Apache-2.0), `output/schema/schema.json` on the `8.19` branch at commit
  `fd5d063c4824de5d33618b7fca4df0e6c1b8489f` (matching the pinned
  `elasticsearch:8.19.25` server used by `elasticsearch/driver`).
  `../../testdata/pins.tsv` pins the commit URL, the SHA-256 of `schema.json`,
  and the `8.19` tag under role `default` (leaving role `driver` for the
  container tapes).
- **`ElasticsearchConsole.g4` is generated deterministically** by
  `../cmd/specgrammar` from `schema.json`. It covers all 566 official endpoints
  across 44 namespaces (561 distinct sampled URL paths, 471 non-GET official
  `(method, path)` pairs across `POST`, `PUT`, `DELETE`, and `HEAD`, 83
  query-parameter call shapes built from the spec's query parameter definitions,
  25 Dev Console relative paths without a leading `/`, and 8 special-character
  Elasticsearch paths with wildcards, commas, colons, and percent-encoding).
  Each path is emitted as a single literal token so `enginetest/grammar` renders
  single-line Console request lines without inserting spaces inside URLs.
- **Token forms** are in `../../testdata/tokens.tsv`, citing the six Kibana Dev
  Console HTTP methods (`GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`).

To regenerate the grammar after a specification pin bump, then re-record the
rows:

```bash
go run ./elasticsearch/parser/cmd/specgrammar -update-pin -cache /tmp/es-spec
OMNI_RECORD=1 go test ./elasticsearch/parser/ -run '^TestCompatibility$'
```

To verify committed answers read-only:

```bash
OMNI_RECORD=verify go test ./elasticsearch/parser/ -run '^TestCompatibility$'
```

After a parser fix, commit the rows that now agree:

```bash
go test ./elasticsearch/parser/ -run '^TestCompatibility$' -update-compat
```

## Who judges what

Elasticsearch has no SQL statement grammar, and its server does not parse Kibana
Dev Console scripts; its syntax rows cover single-line REST call shapes (method
names, endpoints, and query-string forms), not multiline request bodies (§2.2 of
`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`).

- **Grammar judge (`SourceGrammar`).** A recognizer over
  `ElasticsearchConsole.g4` accepts any request derived by `request` and rejects
  every other row. Because `request` includes `httpMethod consolePath` alongside
  the official per-method endpoint rules, a valid generic Console request with a
  swapped HTTP method (such as `PUT /_aliases` or `HEAD /_cat/aliases`) is
  accepted on both sides rather than forcing `elasticsearch/parser` to hardcode
  an endpoint allowlist, while malformed request lines (missing method, missing
  URL, unknown HTTP method, duplicate method, stray trailing tokens, same-line
  JSON bodies, comment-only URLs, and unterminated quotes) are rejected on both
  sides.
- **Baseline vs. final.** Before the `elasticsearch/parser` fixes, 1,892 of
  3,053 rows agreed (1,161 `omni-only` gaps: all 1,154 `!insert` near-misses
  such as `GET /_search zzz` and 7 `suite:console` rows such as
  `GET /_search {"query":{}}`, `GET //comment`, and
  `GET /_search?q="unterminated`). Tightening request-line URL parsing in
  `elasticsearch/parser/parser.go` closed all 1,161 gaps to 3,053 / 3,053
  agreeing.

## What the judge cannot decide

- **Multiline JSON/HJSON/NDJSON bodies and multi-request scripts.** `compat.tsv`
  is one statement per line; body parsing, comment stripping, nesting limits,
  and multi-request splitting are covered by `elasticsearch/parser` and
  `elasticsearch/parsertest`.
- **Lexical whitespace sensitivity.** `enginetest/grammar`'s recognizer uses a
  generic SQL lexer that ignores whitespace between tokens and skips only `--`
  and `/* ... */` comments (not `//` or `#`). Whitespace-sensitive checks such as
  `GET/_search` (missing separator), `GET\t/_search` (tab separator), inline
  `//` and `#` comments, and CRLF (`\r\n`) NDJSON bodies are verified directly
  in `elasticsearch/parser/parser_test.go`.
- **Arbitrary custom or plugin paths.** `elasticsearch/parser` accepts any valid
  non-whitespace URL path (including custom index names and OpenSearch
  `_plugins/...` paths), whereas `ElasticsearchConsole.g4` enumerates the pinned
  official Elasticsearch 8.19 specification paths plus representative relative
  and special-character paths.
