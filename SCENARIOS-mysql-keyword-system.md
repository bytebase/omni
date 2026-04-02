# MySQL Parser Keyword System Alignment Scenarios

> Goal: Fully align omni's keyword classification and identifier parsing with MySQL 8.0's grammar (sql_yacc.yy) — 6-category keyword system, 5 context-dependent identifier functions, zero eqFold workarounds for registered keywords
> Verification: `go test ./mysql/parser/... -count=1` + oracle corpus tests
> Reference sources: mysql-server sql/sql_yacc.yy (grammar), sql/lex.h (keyword definitions), MySQL 8.0 docs

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Keyword Infrastructure

Build the 6-category keyword classification system and register all MySQL 8.0 keywords.

### 1.1 Keyword Classification Data Structure

Define the 6-category system matching MySQL's sql_yacc.yy:

- [x] Keywords have a category attribute: reserved, unambiguous, ambiguous_1, ambiguous_2, ambiguous_3, ambiguous_4
- [x] `isReserved(tokenType)` returns true only for reserved keywords
- [x] `isIdentKeyword(tokenType)` returns true for all 5 non-reserved categories (unambiguous + ambiguous 1-4)
- [x] `isLabelKeyword(tokenType)` returns true for unambiguous + ambiguous_3 + ambiguous_4 (excludes ambiguous_1, ambiguous_2)
- [x] `isRoleKeyword(tokenType)` returns true for unambiguous + ambiguous_2 + ambiguous_4 (excludes ambiguous_1, ambiguous_3)
- [x] `isLvalueKeyword(tokenType)` returns true for unambiguous + ambiguous_1 + ambiguous_2 + ambiguous_3 (excludes ambiguous_4)
- [x] All existing tests still pass — classification is additive, does not change behavior yet

### 1.2 Register Reserved Keywords — Core SQL

Register MySQL 8.0 reserved keywords missing from lexer. These are words declared WITHOUT `<lexer.keyword>` type tag in sql_yacc.yy.

- [x] Register `accessible` (reserved)
- [x] Register `asensitive` (reserved)
- [x] Register `cube` (reserved)
- [x] Register `cume_dist` (reserved)
- [x] Register `dense_rank` (reserved)
- [x] Register `dual` (reserved)
- [x] Register `first_value` (reserved)
- [x] Register `grouping` (reserved)
- [x] Register `insensitive` (reserved)
- [x] Register `lag` (reserved)
- [x] Register `last_value` (reserved)
- [x] Register `lead` (reserved)
- [x] Register `nth_value` (reserved)
- [x] Register `ntile` (reserved)
- [x] Register `of` (reserved)
- [x] Register `optimizer_costs` (reserved)
- [x] Register `percent_rank` (reserved)
- [x] Register `rank` (reserved)
- [x] Register `row_number` (reserved)
- [x] Register `sensitive` (reserved)
- [x] Register `specific` (reserved)
- [x] Register `usage` (reserved)
- [x] Register `varying` (reserved)
- [x] All newly registered words lex as keyword tokens, not tokIDENT

### 1.3 Register Reserved Keywords — Compound Interval & Temporal

- [ ] Register `day_hour` (reserved)
- [ ] Register `day_microsecond` (reserved)
- [ ] Register `day_minute` (reserved)
- [ ] Register `day_second` (reserved)
- [ ] Register `hour_microsecond` (reserved)
- [ ] Register `hour_minute` (reserved)
- [ ] Register `hour_second` (reserved)
- [ ] Register `minute_microsecond` (reserved)
- [ ] Register `minute_second` (reserved)
- [ ] Register `second_microsecond` (reserved)
- [ ] Register `year_month` (reserved)
- [ ] Register `utc_date` (reserved)
- [ ] Register `utc_time` (reserved)
- [ ] Register `utc_timestamp` (reserved)
- [ ] Register `maxvalue` (reserved)
- [ ] Register `no_write_to_binlog` (reserved)
- [ ] Register `io_after_gtids` (reserved)
- [ ] Register `io_before_gtids` (reserved)
- [ ] Register `sqlexception` (reserved)
- [ ] Register `sqlstate` (reserved)
- [ ] Register `sqlwarning` (reserved)
- [ ] INTERVAL expressions still work with keyword tokens as unit names
- [ ] UTC functions still work when lexed as keyword tokens

### 1.4 Classify Existing Keywords — Move to Reserved

Existing keywords in lexer.go that are MySQL 8.0 reserved but missing from reservedKeywords map.

- [ ] Classify kwCROSS as reserved
- [ ] Classify kwNATURAL as reserved
- [ ] Classify kwUSING as reserved
- [ ] Classify kwASC as reserved
- [ ] Classify kwDESC as reserved
- [ ] Classify kwTO as reserved
- [ ] Classify kwDIV as reserved
- [ ] Classify kwMOD as reserved
- [ ] Classify kwXOR as reserved
- [ ] Classify kwREGEXP as reserved
- [ ] Classify kwBINARY as reserved
- [ ] Classify kwINTERVAL as reserved
- [ ] Classify kwMATCH as reserved
- [ ] Classify kwCURRENT_DATE as reserved
- [ ] Classify kwCURRENT_TIME as reserved
- [ ] Classify kwCURRENT_TIMESTAMP as reserved
- [ ] Classify kwCURRENT_USER as reserved
- [ ] Classify kwDATABASE as reserved
- [ ] Classify kwFUNCTION as reserved
- [ ] Classify kwPROCEDURE as reserved
- [ ] Classify kwTRIGGER as reserved
- [ ] Classify kwPARTITION as reserved
- [ ] Classify kwRANGE as reserved
- [ ] Classify kwROW as reserved
- [ ] Classify kwROWS as reserved
- [ ] Classify kwOVER as reserved
- [ ] Classify kwWINDOW as reserved
- [ ] Classify kwFORCE as reserved
- [ ] Classify kwCONVERT as reserved
- [ ] Classify kwCAST as reserved
- [ ] Classify kwWITH as reserved
- [ ] Classify kwREPLACE as reserved
- [ ] Classify kwIGNORE as reserved
- [ ] Classify kwLOAD as reserved
- [ ] Classify kwUSE as reserved
- [ ] Classify kwKILL as reserved
- [ ] Classify kwEXPLAIN as reserved
- [ ] Classify kwSPATIAL as reserved
- [ ] Classify kwFULLTEXT as reserved
- [ ] Classify kwOUTFILE as reserved
- [ ] All existing tests still pass (no behavior change yet — reservedKeywords map not enforced differently until Phase 2)

### 1.5 Classify Existing Keywords — Ambiguous Categories

Existing keywords that are MySQL 8.0 non-reserved with specific ambiguous classification.

- [ ] Classify kwEXECUTE as ambiguous_1 (not label, not role)
- [ ] Classify kwBEGIN as ambiguous_2 (not label)
- [ ] Classify kwCOMMIT as ambiguous_2 (not label)
- [ ] Classify kwEND as ambiguous_2 (not label) — demote from current reserved status
- [ ] Classify kwCONTAINS as ambiguous_2 (not label)
- [ ] Classify kwDO as ambiguous_2 (not label)
- [ ] Classify kwFLUSH as ambiguous_2 (not label)
- [ ] Classify kwFOLLOWS as ambiguous_2 (not label)
- [ ] Classify kwPRECEDES as ambiguous_2 (not label)
- [ ] Classify kwPREPARE as ambiguous_2 (not label)
- [ ] Classify kwREPAIR as ambiguous_2 (not label)
- [ ] Classify kwRESET as ambiguous_2 (not label)
- [ ] Classify kwROLLBACK as ambiguous_2 (not label)
- [ ] Classify kwSAVEPOINT as ambiguous_2 (not label)
- [ ] Classify kwSIGNED as ambiguous_2 (not label)
- [ ] Classify kwSLAVE as ambiguous_2 (not label)
- [ ] Classify kwSTART as ambiguous_2 (not label)
- [ ] Classify kwSTOP as ambiguous_2 (not label)
- [ ] Classify kwTRUNCATE as ambiguous_2 (not label)
- [ ] Classify kwXA as ambiguous_2 (not label)
- [ ] Classify kwEVENT as ambiguous_3 (not role)
- [ ] Classify kwPROCESS as ambiguous_3 (not role)
- [ ] Classify kwRELOAD as ambiguous_3 (not role)
- [ ] Classify kwREPLICATION as ambiguous_3 (not role)
- [ ] Classify kwGLOBAL as ambiguous_4 (not lvalue)
- [ ] Classify kwSESSION as ambiguous_4 (not lvalue)
- [ ] Classify kwLOCAL as ambiguous_4 (not lvalue)
- [ ] All remaining existing keywords classified as unambiguous (default)

### 1.6 Register Non-Reserved Keywords — Type & Spatial

Register MySQL 8.0 non-reserved keywords currently handled via eqFold.

- [ ] Register `geometry` as unambiguous
- [ ] Register `point` as unambiguous
- [ ] Register `linestring` as unambiguous
- [ ] Register `polygon` as unambiguous
- [ ] Register `multipoint` as unambiguous
- [ ] Register `multilinestring` as unambiguous
- [ ] Register `multipolygon` as unambiguous
- [ ] Register `geometrycollection` as unambiguous
- [ ] Register `serial` as unambiguous
- [ ] Register `national` as unambiguous
- [ ] Register `nchar` as unambiguous
- [ ] Register `nvarchar` as unambiguous
- [ ] Classify existing kwSIGNED as ambiguous_2 (per MySQL classification — already registered, needs reclassification only)
- [ ] Register `precision` as unambiguous
- [ ] Register `boolean` if not already registered (unambiguous)
- [ ] Register `srid` as unambiguous

### 1.7 Register Non-Reserved Keywords — DDL & DML Options

- [ ] Register `enforced` as unambiguous
- [ ] Register `less` as unambiguous
- [ ] Register `than` as unambiguous
- [ ] Register `subpartitions` as unambiguous
- [ ] Register `leaves` as unambiguous
- [ ] Register `parser` as unambiguous
- [ ] Register `compression` as unambiguous
- [ ] Register `insert_method` as unambiguous
- [ ] Register `action` as unambiguous
- [ ] Register `partial` as unambiguous
- [ ] Register `format` as unambiguous
- [ ] Register `xml` as unambiguous
- [ ] Register `concurrent` as unambiguous
- [ ] Register `work` as unambiguous
- [ ] Register `export` as unambiguous
- [ ] Register `upgrade` as unambiguous
- [ ] Register `fast` as unambiguous
- [ ] Register `medium` as unambiguous
- [ ] Register `changed` as unambiguous
- [ ] Register `code` as unambiguous

### 1.8 Register Non-Reserved Keywords — SHOW/SET/Grant/Auth

- [ ] Register `events` as unambiguous
- [ ] Register `indexes` as unambiguous
- [ ] Register `grants` as unambiguous
- [ ] Register `triggers` as unambiguous
- [ ] Register `schemas` as unambiguous
- [ ] Register `partitions` as unambiguous
- [ ] Register `hosts` as unambiguous
- [ ] Register `mutex` as unambiguous
- [ ] Register `profile` as unambiguous
- [ ] Register `replicas` as unambiguous
- [ ] Register `names` as unambiguous
- [ ] Register `account` as unambiguous
- [ ] Register `option` as unambiguous
- [ ] Register `proxy` as ambiguous_3 (per MySQL: not role)
- [ ] Register `routine` as unambiguous
- [ ] Register `expire` as unambiguous
- [ ] Register `never` as unambiguous
- [ ] Register `day` as unambiguous
- [ ] Register `history` as unambiguous
- [ ] Register `reuse` as unambiguous
- [ ] Register `optional` as unambiguous
- [ ] Register `x509` as unambiguous
- [ ] Register `issuer` as unambiguous
- [ ] Register `subject` as unambiguous
- [ ] Register `cipher` as unambiguous

### 1.9 Register Non-Reserved Keywords — Scheduling & Misc

- [ ] Register `schedule` as unambiguous
- [ ] Register `completion` as unambiguous
- [ ] Register `preserve` as unambiguous
- [ ] Register `every` as unambiguous
- [ ] Register `starts` as unambiguous
- [ ] Register `ends` as unambiguous
- [ ] Register `value` as unambiguous
- [ ] Register `stacked` as unambiguous
- [ ] Register `unknown` as unambiguous
- [ ] Register `wait` as unambiguous
- [ ] Register `active` as unambiguous
- [ ] Register `inactive` as unambiguous
- [ ] Register `attribute` as unambiguous
- [ ] Register `admin` as unambiguous
- [ ] Register `description` as unambiguous
- [ ] Register `organization` as unambiguous
- [ ] Register `reference` as unambiguous
- [ ] Register `definition` as unambiguous
- [ ] Register `name` as unambiguous
- [ ] Register `system` as unambiguous
- [ ] Register `rotate` as unambiguous
- [ ] Register `keyring` as unambiguous
- [ ] Register `tls` as unambiguous
- [ ] Register `stream` as unambiguous
- [ ] Register `generate` as unambiguous
- [ ] Completeness: all MySQL 8.0 keywords that appear in omni's eqFold patterns are now registered

---

## Phase 2: Identifier Context Functions

Create context-dependent identifier parsing functions and migrate all call sites.

### 2.1 Identifier Function Variants

Create 5 identifier parsing functions matching MySQL's grammar hierarchy.

- [ ] `parseIdent()` — accepts tokIDENT + all 5 non-reserved keyword categories (ident rule)
- [ ] `parseLabelIdent()` — accepts tokIDENT + unambiguous + ambiguous_3 + ambiguous_4 (label_ident rule)
- [ ] `parseRoleIdent()` — accepts tokIDENT + unambiguous + ambiguous_2 + ambiguous_4 (role_ident rule)
- [ ] `parseLvalueIdent()` — accepts tokIDENT + unambiguous + ambiguous_1 + ambiguous_2 + ambiguous_3 (lvalue_ident rule)
- [ ] `parseKeywordOrIdent()` — accepts tokIDENT + ANY keyword token (for option values, enum values, action words)
- [ ] Existing `parseIdentifier()` becomes an alias for `parseIdent()` (gradual migration)
- [ ] `parseTableRef()` and `parseColumnRef()` use `parseIdent()` internally
- [ ] `isIdentToken()` updated to match `parseIdent()` semantics
- [ ] All existing tests still pass after function creation (no call site changes yet)

### 2.2 Migrate General Ident Call Sites — select.go

- [ ] CTE name uses parseIdent
- [ ] SELECT alias (after AS) uses parseIdent
- [ ] SELECT implicit alias uses parseIdent
- [ ] JOIN USING column uses parseIdent
- [ ] Derived table alias uses parseIdent
- [ ] WINDOW name uses parseIdent
- [ ] Index hint name uses parseIdent
- [ ] INTO OUTFILE charset uses parseIdent
- [ ] All existing SELECT tests still pass

### 2.3 Migrate General Ident Call Sites — DDL files

- [ ] Column definition name uses parseIdent
- [ ] Constraint name uses parseIdent
- [ ] Index name uses parseIdent
- [ ] Partition name uses parseIdent
- [ ] CREATE DATABASE name uses parseIdent
- [ ] CREATE VIEW column name uses parseIdent
- [ ] Procedure/function parameter name uses parseIdent
- [ ] Trigger name uses parseIdent
- [ ] Event name uses parseIdent
- [ ] All DDL tests still pass

### 2.4 Migrate General Ident Call Sites — DML & Other files

- [ ] INSERT table alias uses parseIdent
- [ ] DELETE table alias uses parseIdent
- [ ] UPDATE SET target uses parseIdent (via parseColumnRef)
- [ ] PREPARE/EXECUTE statement name uses parseIdent
- [ ] SAVEPOINT name uses parseIdent
- [ ] DECLARE variable/cursor name uses parseIdent
- [ ] GRANT user/host name uses parseIdent
- [ ] COLLATE collation name uses parseIdent
- [ ] All DML tests still pass

### 2.5 Migrate Label Ident Call Sites

- [ ] BEGIN...END block label uses parseLabelIdent
- [ ] LEAVE label uses parseLabelIdent
- [ ] ITERATE label uses parseLabelIdent
- [ ] `CREATE TABLE begin (a INT)` accepted — BEGIN is ambiguous_2 (allowed in ident, not in label)
- [ ] `label1: BEGIN ... END label1` with label1=`begin` rejected — BEGIN not allowed as label
- [ ] All compound statement tests still pass

### 2.6 Migrate Role Ident Call Sites

- [ ] GRANT WITH ROLE role_name uses parseRoleIdent
- [ ] `CREATE ROLE event` rejected — EVENT is ambiguous_3 (not allowed as role)
- [ ] `CREATE ROLE begin` accepted — BEGIN is ambiguous_2 (allowed as role)
- [ ] All GRANT/ROLE tests still pass

### 2.7 Migrate Lvalue Ident Call Sites

- [ ] SET variable name uses parseLvalueIdent
- [ ] RESET PERSIST variable name uses parseLvalueIdent
- [ ] `SET global = 1` rejected — GLOBAL is ambiguous_4 (not allowed as lvalue)
- [ ] `SET begin = 1` accepted — BEGIN is ambiguous_2 (allowed as lvalue)
- [ ] All SET tests still pass

### 2.8 Migrate Any-Keyword Call Sites — DDL Options

Option values, enum values, and action words that should accept ANY keyword including reserved.

- [ ] ALGORITHM = value (UNDEFINED/MERGE/TEMPTABLE/DEFAULT/INSTANT/INPLACE/COPY) uses parseKeywordOrIdent
- [ ] LOCK = value (DEFAULT/NONE/SHARED/EXCLUSIVE) uses parseKeywordOrIdent
- [ ] SQL SECURITY value (DEFINER/INVOKER) uses parseKeywordOrIdent
- [ ] MATCH type (FULL/PARTIAL/SIMPLE) uses parseKeywordOrIdent
- [ ] USING index type (BTREE/HASH) uses parseKeywordOrIdent
- [ ] RETURNS type for loadable function (STRING/INTEGER/REAL) uses parseKeywordOrIdent
- [ ] LANGUAGE value (SQL) uses parseKeywordOrIdent
- [ ] RESOURCE GROUP TYPE value (SYSTEM/USER) uses parseKeywordOrIdent
- [ ] consumeOptionValue fallback uses parseKeywordOrIdent
- [ ] All DDL option tests still pass

### 2.9 Migrate Any-Keyword Call Sites — Replication & Utility

- [ ] ALTER INSTANCE action words use parseKeywordOrIdent — fixes infinite loop bug
- [ ] Replication source option values (ON/OFF/STREAM/GENERATE) use parseKeywordOrIdent — fixes ON keyword bug
- [ ] Replication source option names use parseKeywordOrIdent
- [ ] Replication filter type names use parseKeywordOrIdent
- [ ] START REPLICA UNTIL type names use parseKeywordOrIdent
- [ ] FLUSH/RESET option words use parseKeywordOrIdent
- [ ] EXPLAIN FORMAT value uses parseKeywordOrIdent
- [ ] SERVER OPTIONS keyword names use parseKeywordOrIdent
- [ ] HELP topic uses parseKeywordOrIdent
- [ ] `ALTER INSTANCE ROTATE INNODB MASTER KEY` parses correctly (no infinite loop)
- [ ] `CHANGE REPLICATION SOURCE TO REQUIRE_TABLE_PRIMARY_KEY_CHECK = ON` parses correctly
- [ ] All replication and utility tests still pass

### 2.10 Migrate Any-Keyword Call Sites — Expressions

- [ ] EXTRACT unit uses parseKeywordOrIdent (accepts DAY, HOUR, etc. as keyword tokens)
- [ ] INTERVAL unit uses parseKeywordOrIdent (accepts compound units like DAY_HOUR as keyword tokens)
- [ ] INTERVAL unit validation still works with keyword tokens (not just strings)
- [ ] All expression tests still pass

---

## Phase 3: eqFold Migration

Replace all eqFold string matching with keyword token matching, file by file. Each scenario is: the eqFold call is replaced with a keyword token check, and tests still pass.

### 3.1 eqFold Migration — type.go

- [ ] `eqFold("geometry")` → kwGEOMETRY token check
- [ ] `eqFold("point")` → kwPOINT token check
- [ ] `eqFold("linestring")` → kwLINESTRING token check
- [ ] `eqFold("polygon")` → kwPOLYGON token check
- [ ] `eqFold("multipoint")` → kwMULTIPOINT token check
- [ ] `eqFold("multilinestring")` → kwMULTILINESTRING token check
- [ ] `eqFold("multipolygon")` → kwMULTIPOLYGON token check
- [ ] `eqFold("geometrycollection")` → kwGEOMETRYCOLLECTION token check
- [ ] `eqFold("serial")` → kwSERIAL token check
- [ ] `eqFold("national")` → kwNATIONAL token check
- [ ] `eqFold("nchar")` → kwNCHAR token check
- [ ] `eqFold("nvarchar")` → kwNVARCHAR token check
- [ ] `eqFold("signed")` → kwSIGNED token check
- [ ] `eqFold("precision")` → kwPRECISION token check
- [ ] `eqFold("long")` → kwLONG token check
- [ ] `eqFold("int1")`...`eqFold("int8")` → keyword token checks
- [ ] `eqFold("middleint")` → kwMIDDLEINT token check
- [ ] `eqFold("float4")`/`eqFold("float8")` → keyword token checks
- [ ] `eqFold("srid")` → kwSRID token check
- [ ] Zero eqFold calls remain in type.go for registered keywords
- [ ] All data type tests still pass

### 3.2 eqFold Migration — create_table.go

- [ ] `eqFold("enforced")` → kwENFORCED token check
- [ ] `eqFold("less")` → kwLESS token check
- [ ] `eqFold("than")` → kwTHAN token check
- [ ] `eqFold("maxvalue")` → kwMAXVALUE token check
- [ ] `eqFold("subpartitions")` → kwSUBPARTITIONS token check
- [ ] `eqFold("leaves")` → kwLEAVES token check
- [ ] `eqFold("action")` → kwACTION token check
- [ ] `eqFold("partial")` → kwPARTIAL token check
- [ ] All table option eqFold patterns migrated to keyword tokens where applicable
- [ ] Zero eqFold calls for registered keywords remain in create_table.go
- [ ] All CREATE TABLE tests still pass

### 3.3 eqFold Migration — grant.go

- [ ] `eqFold("account")` → kwACCOUNT token check
- [ ] `eqFold("option")` → kwOPTION token check
- [ ] `eqFold("proxy")` → kwPROXY token check
- [ ] `eqFold("routine")` → kwROUTINE token check
- [ ] `eqFold("expire")` → kwEXPIRE token check
- [ ] `eqFold("never")` → kwNEVER token check
- [ ] `eqFold("day")` → kwDAY token check
- [ ] `eqFold("history")` → kwHISTORY token check
- [ ] `eqFold("reuse")` → kwREUSE token check
- [ ] `eqFold("x509")` → kwX509 token check
- [ ] `eqFold("issuer")` → kwISSUER token check
- [ ] `eqFold("subject")` → kwSUBJECT token check
- [ ] `eqFold("cipher")` → kwCIPHER token check
- [ ] `eqFold("attribute")` → kwATTRIBUTE token check
- [ ] All remaining grant.go eqFold patterns for registered keywords migrated
- [ ] Zero eqFold calls for registered keywords remain in grant.go
- [ ] All GRANT/USER tests still pass

### 3.4 eqFold Migration — utility.go

- [ ] `eqFold("schedule")` → kwSCHEDULE token check
- [ ] `eqFold("completion")` → kwCOMPLETION token check
- [ ] `eqFold("preserve")` → kwPRESERVE token check
- [ ] `eqFold("every")` → kwEVERY token check
- [ ] `eqFold("starts")` → kwSTARTS token check
- [ ] `eqFold("ends")` → kwENDS token check
- [ ] `eqFold("rotate")` → kwROTATE token check
- [ ] `eqFold("keyring")` → kwKEYRING token check
- [ ] `eqFold("tls")` → kwTLS token check
- [ ] `eqFold("concurrent")` → kwCONCURRENT token check
- [ ] `eqFold("work")` → kwWORK token check
- [ ] `eqFold("export")` → kwEXPORT token check
- [ ] `eqFold("upgrade")` → kwUPGRADE token check
- [ ] `eqFold("fast")`/`eqFold("medium")`/`eqFold("changed")` → keyword token checks
- [ ] All remaining utility.go eqFold patterns for registered keywords migrated
- [ ] Zero eqFold calls for registered keywords remain in utility.go
- [ ] All utility tests still pass

### 3.5 eqFold Migration — set_show.go

- [ ] `eqFold("events")` → kwEVENTS token check
- [ ] `eqFold("indexes")` → kwINDEXES token check
- [ ] `eqFold("grants")` → kwGRANTS token check
- [ ] `eqFold("triggers")` → kwTRIGGERS token check
- [ ] `eqFold("schemas")` → kwSCHEMAS token check
- [ ] `eqFold("partitions")` → kwPARTITIONS token check
- [ ] `eqFold("hosts")` → kwHOSTS token check
- [ ] `eqFold("mutex")` → kwMUTEX token check
- [ ] `eqFold("profile")` → kwPROFILE token check
- [ ] `eqFold("replicas")` → kwREPLICAS token check
- [ ] `eqFold("format")` → kwFORMAT token check
- [ ] `eqFold("names")` → kwNAMES token check
- [ ] `eqFold("code")` → kwCODE token check
- [ ] `eqFold("xml")` → kwXML token check
- [ ] Zero eqFold calls for registered keywords remain in set_show.go
- [ ] All SHOW/SET tests still pass

### 3.6 eqFold Migration — replication.go & trigger.go & signal.go

- [ ] All replication.go eqFold patterns for registered keywords migrated
- [ ] `eqFold("stream")` → kwSTREAM token check
- [ ] `eqFold("generate")` → kwGENERATE token check
- [ ] All trigger.go eqFold patterns for registered keywords migrated (interval units in event scheduling)
- [ ] All signal.go eqFold patterns migrated (`value`, `stacked`)
- [ ] Zero eqFold calls for registered keywords remain in these files
- [ ] All replication/trigger/signal tests still pass

### 3.7 eqFold Migration — expr.go & compound.go & remaining files

- [ ] All expr.go eqFold patterns for registered keywords migrated
- [ ] All compound.go eqFold patterns migrated
- [ ] All create_function.go remaining eqFold patterns migrated
- [ ] All create_index.go remaining eqFold patterns migrated
- [ ] All create_view.go remaining eqFold patterns migrated
- [ ] All alter_table.go remaining eqFold patterns migrated
- [ ] All alter_misc.go remaining eqFold patterns migrated
- [ ] All load_data.go remaining eqFold patterns migrated
- [ ] All transaction.go eqFold patterns migrated
- [ ] All insert.go eqFold patterns migrated
- [ ] All select.go eqFold patterns migrated (e.g., `of` in FOR UPDATE OF)
- [ ] All stmt.go eqFold patterns migrated (e.g., `reference` in CREATE SRS)
- [ ] Zero eqFold calls for registered keywords remain across entire parser
- [ ] All tests still pass

### 3.8 Completeness Audit

- [ ] Every MySQL 8.0 keyword in sql_yacc.yy that appears in omni's parser is registered as a keyword token
- [ ] Every registered keyword has the correct 6-category classification matching sql_yacc.yy
- [ ] Zero eqFold calls remain for strings that are registered keywords
- [ ] eqFold calls only remain for: (a) non-keyword compound option strings (key_block_size, max_rows, replication options, etc.), (b) option-name dispatch patterns in create_table.go where post-parse string matching is used, (c) `@@`-prefixed variable scope parsing in name.go (lexer emits `@@global.var` as single token)
- [ ] Oracle corpus: `CREATE TABLE select (a INT)` correctly rejected
- [ ] Oracle corpus: `CREATE TABLE t (select INT)` correctly rejected
- [ ] Oracle corpus: `CREATE TABLE t (rank INT)` correctly rejected (rank is reserved)
- [ ] Oracle corpus: `CREATE TABLE t (status INT)` correctly accepted (status is non-reserved)
- [ ] Oracle corpus: `CREATE TABLE begin (a INT)` correctly accepted (begin is ambiguous_2, allowed as ident)
- [ ] Full test suite passes with zero regressions
