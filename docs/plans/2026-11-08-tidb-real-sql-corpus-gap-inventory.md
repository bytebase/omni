# TiDB real-SQL corpus gap inventory (§3.4 sample 1)

Date: 2026-11-08. The first measurement of §3.4 sample 1, "the rest of pingcap's
parser tests" (`tidb/parser/testdata/pingcap_parser_rest_corpus.tsv`). The
`parser.y` completeness corpus (§3.3,
`2026-11-07-tidb-parser-y-gap-inventory.md`) is closed; this is the next grammar
work list, the forms that only fail when combined the way people write them.

## The measurement

885 statements, each the `src` of a pingcap `parser_test.go` case the first
extraction skipped, or the statement an `ast/*_test.go` restore case parses
(a fragment inside a full-statement template). Each was measured against omni and
against TiDB v8.5.5 via `PREPARE`.

| | first | after FK fix | after var fix | after OPTIMIZE fix |
|---|---:|---:|---:|---:|
| both accept | 869 | 874 | 879 | 880 |
| gap: TiDB accepts, omni rejects | 11 | 6 | 1 | 0 |
| **over-acceptance: TiDB rejects, omni accepts** | **0** | **0** | **0** | **0** |
| both reject | 5 | 5 | 5 | 5 |
| no verdict | 0 | 0 | 0 | 0 |

**COMPLETE: §3.4 sample 1 is fully closed.** omni parses **880 of the 880
statements TiDB accepts — 100.0%** (was 869/880, 98.8%, at first measurement),
with **zero over-acceptances** and zero remaining gaps. The 5 both-reject rows
are statements the server itself refuses (`AFFINITY 1` with a non-string, an
unterminated `/*` comment, a resource group option with a bad value); they are
not work.

## The gaps, grouped

### A. Foreign-key reference lists with index parts — 5 rows — ✅ CLOSED

**Closed** by `fix(tidb): read index-part lists in a REFERENCES column list`.
TiDB's `REFERENCES t (...)` list is an `IndexPartSpecificationList`, not a plain
column list: a column may carry a prefix length and the part may be an
expression (the FK's own `FOREIGN KEY (...)` list omni already read this way).
omni accepted only plain names in the REFERENCES list. MySQL, MariaDB and
OceanBase take only plain names there and even refuse a functional key part, so
this is a TiDB dialect delta, not ported to the forks.

```
CREATE TABLE child (id INT, parent_id INT, INDEX par_ind (parent_id), FOREIGN KEY (parent_id) REFERENCES parent(id,hello(12)) ON DELETE CASCADE)
CREATE TABLE child (id INT, parent_id INT, INDEX par_ind (parent_id), FOREIGN KEY (parent_id) REFERENCES parent(id(8),hello(12)) ON DELETE CASCADE)
CREATE TABLE child (id INT, parent_id INT, INDEX par_ind (parent_id), FOREIGN KEY (parent_id) REFERENCES parent((id+1)))
CREATE TABLE child (id INT, parent_id INT, CONSTRAINT FOREIGN KEY (parent_id(2),hello(4)) REFERENCES parent((id+1)) ON DELETE CASCADE)
CREATE TABLE child (id INT, parent_id INT, CONSTRAINT FOREIGN KEY (parent_id) REFERENCES parent((id+1)) ON DELETE CASCADE ON UPDATE RESTRICT)
```

### B. User- and system-variable names as quoted strings — 5 rows — ✅ CLOSED

**Closed** by `fix(tidb): read quoted user- and system-variable names`. A user
variable `@name` and a system variable `@@scope.name` take a quoted string or
quoted identifier as the name, not only a bare word. The lexer now folds a
single/double/backtick-quoted user-variable name into one token (TiDB's
`singleAtIdentifier`, which the `@'host'` of a username also uses), and
`parseSysVarName` reads a quoted system-variable name component. It also closed
`SET @' ' = 1` in the single-form pingcap corpus.

```
select @'a':=1
select @`a``b`=4
select @"aBC">1
select @@session.'C'
select @@local."aBc"
```

### C. `OPTIMIZE TABLES` plural — 1 row — ✅ CLOSED

**Closed** by `fix(tidb,mysql,mariadb,oceanbase): accept OPTIMIZE TABLES and
require TABLE or TABLES`.

```
optimize tables t1, t2
```

`OPTIMIZE`'s `TableOrTables` is `TABLE | TABLES` and is required; omni matched
only the singular and made it optional, so it missed `OPTIMIZE TABLES` and
over-accepted `OPTIMIZE t` with neither. Measured against TiDB v8.5.5, MySQL 8.0,
MariaDB 11.8 and OceanBase: all four accept `OPTIMIZE TABLES t` and answer 1064 to
`OPTIMIZE t`, so the fix is the same across the family (fixing the latent
`OPTIMIZE t` over-acceptance in each).

## Next

Only cluster C remains (1 row, `OPTIMIZE TABLES`). Close it as §3.2's rows were
closed — a small PR with §3.5's tree checks, zero new over-acceptances — then
re-record this corpus with `-update` and read the diff. Samples 2 (TiDB
integrationtest) and 3 (Bytebase's own TiDB test files) are later sessions and
will extend this list.
