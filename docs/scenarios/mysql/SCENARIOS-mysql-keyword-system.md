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

- [x] Register `day_hour` (reserved)
- [x] Register `day_microsecond` (reserved)
- [x] Register `day_minute` (reserved)
- [x] Register `day_second` (reserved)
- [x] Register `hour_microsecond` (reserved)
- [x] Register `hour_minute` (reserved)
- [x] Register `hour_second` (reserved)
- [x] Register `minute_microsecond` (reserved)
- [x] Register `minute_second` (reserved)
- [x] Register `second_microsecond` (reserved)
- [x] Register `year_month` (reserved)
- [x] Register `utc_date` (reserved)
- [x] Register `utc_time` (reserved)
- [x] Register `utc_timestamp` (reserved)
- [x] Register `maxvalue` (reserved)
- [x] Register `no_write_to_binlog` (reserved)
- [x] Register `io_after_gtids` (reserved)
- [x] Register `io_before_gtids` (reserved)
- [x] Register `sqlexception` (reserved)
- [x] Register `sqlstate` (reserved)
- [x] Register `sqlwarning` (reserved)
- [x] INTERVAL expressions still work with keyword tokens as unit names
- [x] UTC functions still work when lexed as keyword tokens

### 1.4 Classify Existing Keywords — Move to Reserved

Existing keywords in lexer.go that are MySQL 8.0 reserved but missing from reservedKeywords map.

- [x] Classify kwCROSS as reserved
- [x] Classify kwNATURAL as reserved
- [x] Classify kwUSING as reserved
- [x] Classify kwASC as reserved
- [x] Classify kwDESC as reserved
- [x] Classify kwTO as reserved
- [x] Classify kwDIV as reserved
- [x] Classify kwMOD as reserved
- [x] Classify kwXOR as reserved
- [x] Classify kwREGEXP as reserved
- [x] Classify kwBINARY as reserved
- [x] Classify kwINTERVAL as reserved
- [x] Classify kwMATCH as reserved
- [x] Classify kwCURRENT_DATE as reserved
- [x] Classify kwCURRENT_TIME as reserved
- [x] Classify kwCURRENT_TIMESTAMP as reserved
- [x] Classify kwCURRENT_USER as reserved
- [x] Classify kwDATABASE as reserved
- [x] Classify kwFUNCTION as reserved
- [x] Classify kwPROCEDURE as reserved
- [x] Classify kwTRIGGER as reserved
- [x] Classify kwPARTITION as reserved
- [x] Classify kwRANGE as reserved
- [x] Classify kwROW as reserved
- [x] Classify kwROWS as reserved
- [x] Classify kwOVER as reserved
- [x] Classify kwWINDOW as reserved
- [x] Classify kwFORCE as reserved
- [x] Classify kwCONVERT as reserved
- [x] Classify kwCAST as reserved
- [x] Classify kwWITH as reserved
- [x] Classify kwREPLACE as reserved
- [x] Classify kwIGNORE as reserved
- [x] Classify kwLOAD as reserved
- [x] Classify kwUSE as reserved
- [x] Classify kwKILL as reserved
- [x] Classify kwEXPLAIN as reserved
- [x] Classify kwSPATIAL as reserved
- [x] Classify kwFULLTEXT as reserved
- [x] Classify kwOUTFILE as reserved
- [x] All existing tests still pass (no behavior change yet — reservedKeywords map not enforced differently until Phase 2)

### 1.5 Classify Existing Keywords — Ambiguous Categories

Existing keywords that are MySQL 8.0 non-reserved with specific ambiguous classification.

- [x] Classify kwEXECUTE as ambiguous_1 (not label, not role)
- [x] Classify kwBEGIN as ambiguous_2 (not label)
- [x] Classify kwCOMMIT as ambiguous_2 (not label)
- [x] Classify kwEND as ambiguous_2 (not label) — demote from current reserved status
- [x] Classify kwCONTAINS as ambiguous_2 (not label)
- [x] Classify kwDO as ambiguous_2 (not label)
- [x] Classify kwFLUSH as ambiguous_2 (not label)
- [x] Classify kwFOLLOWS as ambiguous_2 (not label)
- [x] Classify kwPRECEDES as ambiguous_2 (not label)
- [x] Classify kwPREPARE as ambiguous_2 (not label)
- [x] Classify kwREPAIR as ambiguous_2 (not label)
- [x] Classify kwRESET as ambiguous_2 (not label)
- [x] Classify kwROLLBACK as ambiguous_2 (not label)
- [x] Classify kwSAVEPOINT as ambiguous_2 (not label)
- [x] Classify kwSIGNED as ambiguous_2 (not label)
- [x] Classify kwSLAVE as ambiguous_2 (not label)
- [x] Classify kwSTART as ambiguous_2 (not label)
- [x] Classify kwSTOP as ambiguous_2 (not label)
- [x] Classify kwTRUNCATE as ambiguous_2 (not label)
- [x] Classify kwXA as ambiguous_2 (not label)
- [x] Classify kwEVENT as ambiguous_3 (not role)
- [x] Classify kwPROCESS as ambiguous_3 (not role)
- [x] Classify kwRELOAD as ambiguous_3 (not role)
- [x] Classify kwREPLICATION as ambiguous_3 (not role)
- [x] Classify kwGLOBAL as ambiguous_4 (not lvalue)
- [x] Classify kwSESSION as ambiguous_4 (not lvalue)
- [x] Classify kwLOCAL as ambiguous_4 (not lvalue)
- [x] All remaining existing keywords classified as unambiguous (default)

### 1.6 Register Non-Reserved Keywords — Type & Spatial

Register MySQL 8.0 non-reserved keywords currently handled via eqFold.

- [x] Register `geometry` as unambiguous
- [x] Register `point` as unambiguous
- [x] Register `linestring` as unambiguous
- [x] Register `polygon` as unambiguous
- [x] Register `multipoint` as unambiguous
- [x] Register `multilinestring` as unambiguous
- [x] Register `multipolygon` as unambiguous
- [x] Register `geometrycollection` as unambiguous
- [x] Register `serial` as unambiguous
- [x] Register `national` as unambiguous
- [x] Register `nchar` as unambiguous
- [x] Register `nvarchar` as unambiguous
- [x] Classify existing kwSIGNED as ambiguous_2 (per MySQL classification — already registered, needs reclassification only)
- [x] Register `precision` as unambiguous
- [x] Register `boolean` if not already registered (unambiguous)
- [x] Register `srid` as unambiguous

### 1.7 Register Non-Reserved Keywords — DDL & DML Options

- [x] Register `enforced` as unambiguous
- [x] Register `less` as unambiguous
- [x] Register `than` as unambiguous
- [x] Register `subpartitions` as unambiguous
- [x] Register `leaves` as unambiguous
- [x] Register `parser` as unambiguous
- [x] Register `compression` as unambiguous
- [x] Register `insert_method` as unambiguous
- [x] Register `action` as unambiguous
- [x] Register `partial` as unambiguous
- [x] Register `format` as unambiguous
- [x] Register `xml` as unambiguous
- [x] Register `concurrent` as unambiguous
- [x] Register `work` as unambiguous
- [x] Register `xid` as unambiguous
- [x] Register `export` as unambiguous
- [x] Register `upgrade` as unambiguous
- [x] Register `fast` as unambiguous
- [x] Register `medium` as unambiguous
- [x] Register `changed` as unambiguous
- [x] Register `code` as unambiguous

### 1.8 Register Non-Reserved Keywords — SHOW/SET/Grant/Auth

- [x] Register `events` as unambiguous
- [x] Register `indexes` as unambiguous
- [x] Register `grants` as unambiguous
- [x] Register `triggers` as unambiguous
- [x] Register `schemas` as unambiguous
- [x] Register `partitions` as unambiguous
- [x] Register `hosts` as unambiguous
- [x] Register `mutex` as unambiguous
- [x] Register `profile` as unambiguous
- [x] Register `replicas` as unambiguous
- [x] Register `names` as unambiguous
- [x] Register `account` as unambiguous
- [x] Register `option` as unambiguous
- [x] Register `proxy` as ambiguous_3 (per MySQL: not role)
- [x] Register `routine` as unambiguous
- [x] Register `expire` as unambiguous
- [x] Register `never` as unambiguous
- [x] Register `day` as unambiguous
- [x] Register `history` as unambiguous
- [x] Register `reuse` as unambiguous
- [x] Register `optional` as unambiguous
- [x] Register `x509` as unambiguous
- [x] Register `issuer` as unambiguous
- [x] Register `subject` as unambiguous
- [x] Register `cipher` as unambiguous

### 1.9 Register Non-Reserved Keywords — Scheduling & Misc

- [x] Register `schedule` as unambiguous
- [x] Register `completion` as unambiguous
- [x] Register `preserve` as unambiguous
- [x] Register `every` as unambiguous
- [x] Register `starts` as unambiguous
- [x] Register `ends` as unambiguous
- [x] Register `value` as unambiguous
- [x] Register `stacked` as unambiguous
- [x] Register `unknown` as unambiguous
- [x] Register `wait` as unambiguous
- [x] Register `active` as unambiguous
- [x] Register `inactive` as unambiguous
- [x] Register `attribute` as unambiguous
- [x] Register `admin` as unambiguous
- [x] Register `description` as unambiguous
- [x] Register `organization` as unambiguous
- [x] Register `reference` as unambiguous
- [x] Register `definition` as unambiguous
- [x] Register `name` as unambiguous
- [x] Register `system` as unambiguous
- [x] Register `rotate` as unambiguous
- [x] Register `keyring` as unambiguous
- [x] Register `tls` as unambiguous
- [x] Register `stream` as unambiguous
- [x] Register `generate` as unambiguous
- [x] Completeness: all MySQL 8.0 keywords that appear in omni's eqFold patterns are now registered

---

## Phase 2: Identifier Context Functions

Create context-dependent identifier parsing functions and migrate all call sites.

### 2.1 Identifier Function Variants

Create 5 identifier parsing functions matching MySQL's grammar hierarchy.

- [x] `parseIdent()` — accepts tokIDENT + all 5 non-reserved keyword categories (ident rule)
- [x] `parseLabelIdent()` — accepts tokIDENT + unambiguous + ambiguous_3 + ambiguous_4 (label_ident rule)
- [x] `parseRoleIdent()` — accepts tokIDENT + unambiguous + ambiguous_2 + ambiguous_4 (role_ident rule)
- [x] `parseLvalueIdent()` — accepts tokIDENT + unambiguous + ambiguous_1 + ambiguous_2 + ambiguous_3 (lvalue_ident rule)
- [x] `parseKeywordOrIdent()` — accepts tokIDENT + ANY keyword token (for option values, enum values, action words)
- [x] Existing `parseIdentifier()` becomes an alias for `parseIdent()` (gradual migration)
- [x] `parseTableRef()` and `parseColumnRef()` use `parseIdent()` internally
- [x] `isIdentToken()` updated to match `parseIdent()` semantics
- [x] All existing tests still pass after function creation (no call site changes yet)

### 2.2 Migrate General Ident Call Sites — select.go

- [x] CTE name uses parseIdent
- [x] SELECT alias (after AS) uses parseIdent
- [x] SELECT implicit alias uses parseIdent
- [x] JOIN USING column uses parseIdent
- [x] Derived table alias uses parseIdent
- [x] WINDOW name uses parseIdent
- [x] Index hint name uses parseIdent
- [x] INTO OUTFILE charset uses parseIdent
- [x] All existing SELECT tests still pass

### 2.3 Migrate General Ident Call Sites — DDL files

- [x] Column definition name uses parseIdent
- [x] Constraint name uses parseIdent
- [x] Index name uses parseIdent
- [x] Partition name uses parseIdent
- [x] CREATE DATABASE name uses parseIdent
- [x] CREATE VIEW column name uses parseIdent
- [x] Procedure/function parameter name uses parseIdent
- [x] Trigger name uses parseIdent
- [x] Event name uses parseIdent
- [x] All DDL tests still pass

### 2.4 Migrate General Ident Call Sites — DML & Other files

- [x] INSERT table alias uses parseIdent
- [x] DELETE table alias uses parseIdent
- [x] UPDATE SET target uses parseIdent (via parseColumnRef)
- [x] PREPARE/EXECUTE statement name uses parseIdent
- [x] SAVEPOINT name uses parseIdent
- [x] DECLARE variable/cursor name uses parseIdent
- [x] GRANT user/host name uses parseIdent
- [x] COLLATE collation name uses parseIdent
- [x] All DML tests still pass

### 2.5 Migrate Label Ident Call Sites

- [x] BEGIN...END block label uses parseLabelIdent
- [x] LEAVE label uses parseLabelIdent
- [x] ITERATE label uses parseLabelIdent
- [x] `CREATE TABLE begin (a INT)` accepted — BEGIN is ambiguous_2 (allowed in ident, not in label)
- [x] `label1: BEGIN ... END label1` with label1=`begin` rejected — BEGIN not allowed as label
- [x] All compound statement tests still pass

### 2.6 Migrate Role Ident Call Sites

- [ ] GRANT WITH ROLE role_name uses parseRoleIdent
- [ ] `CREATE ROLE event` rejected — EVENT is ambiguous_3 (not allowed as role)
- [ ] `CREATE ROLE begin` accepted — BEGIN is ambiguous_2 (allowed as role)
- [ ] All GRANT/ROLE tests still pass

### 2.7 Migrate Lvalue Ident Call Sites

- [x] SET variable name uses parseLvalueIdent
- [x] RESET PERSIST variable name uses parseLvalueIdent
- [x] `SET global = 1` rejected — GLOBAL is ambiguous_4 (not allowed as lvalue)
- [x] `SET begin = 1` accepted — BEGIN is ambiguous_2 (allowed as lvalue)
- [x] All SET tests still pass

### 2.8 Migrate Any-Keyword Call Sites — DDL Options

Option values, enum values, and action words that should accept ANY keyword including reserved.

- [x] ALGORITHM = value (UNDEFINED/MERGE/TEMPTABLE/DEFAULT/INSTANT/INPLACE/COPY) uses parseKeywordOrIdent
- [x] LOCK = value (DEFAULT/NONE/SHARED/EXCLUSIVE) uses parseKeywordOrIdent
- [x] SQL SECURITY value (DEFINER/INVOKER) uses parseKeywordOrIdent
- [x] MATCH type (FULL/PARTIAL/SIMPLE) uses parseKeywordOrIdent
- [x] USING index type (BTREE/HASH) uses parseKeywordOrIdent
- [x] RETURNS type for loadable function (STRING/INTEGER/REAL) uses parseKeywordOrIdent
- [x] LANGUAGE value (SQL) uses parseKeywordOrIdent
- [x] RESOURCE GROUP TYPE value (SYSTEM/USER) uses parseKeywordOrIdent
- [x] consumeOptionValue fallback uses parseKeywordOrIdent
- [x] All DDL option tests still pass

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

- [x] `eqFold("geometry")` → kwGEOMETRY token check
- [x] `eqFold("point")` → kwPOINT token check
- [x] `eqFold("linestring")` → kwLINESTRING token check
- [x] `eqFold("polygon")` → kwPOLYGON token check
- [x] `eqFold("multipoint")` → kwMULTIPOINT token check
- [x] `eqFold("multilinestring")` → kwMULTILINESTRING token check
- [x] `eqFold("multipolygon")` → kwMULTIPOLYGON token check
- [x] `eqFold("geometrycollection")` → kwGEOMETRYCOLLECTION token check
- [x] `eqFold("serial")` → kwSERIAL token check
- [x] `eqFold("national")` → kwNATIONAL token check
- [x] `eqFold("nchar")` → kwNCHAR token check
- [x] `eqFold("nvarchar")` → kwNVARCHAR token check
- [x] `eqFold("signed")` → kwSIGNED token check
- [x] `eqFold("precision")` → kwPRECISION token check
- [x] `eqFold("long")` → kwLONG token check — NOTE: int1-int8, middleint, float4/float8 are NOT registered keywords; eqFold correctly remains for these type aliases
- [x] `eqFold("int1")`...`eqFold("int8")` → keyword token checks — N/A: not registered keywords, eqFold retained
- [x] `eqFold("middleint")` → kwMIDDLEINT token check — N/A: not a registered keyword, eqFold retained
- [x] `eqFold("float4")`/`eqFold("float8")` → keyword token checks — N/A: not registered keywords, eqFold retained
- [x] `eqFold("srid")` → kwSRID token check
- [x] Zero eqFold calls remain in type.go for registered keywords
- [x] All data type tests still pass

### 3.2 eqFold Migration — create_table.go

- [x] `eqFold("enforced")` → kwENFORCED token check
- [x] `eqFold("less")` → kwLESS token check
- [x] `eqFold("than")` → kwTHAN token check
- [x] `eqFold("maxvalue")` → kwMAXVALUE token check
- [x] `eqFold("subpartitions")` → kwSUBPARTITIONS token check
- [x] `eqFold("leaves")` → kwLEAVES token check
- [x] `eqFold("action")` → kwACTION token check
- [x] `eqFold("partial")` → kwPARTIAL token check
- [x] All table option eqFold patterns migrated to keyword tokens where applicable — remaining eqFold in create_table.go is for non-keyword option strings (engine_attribute, secondary_engine_attribute, key_block_size, etc.)
- [x] Zero eqFold calls for registered keywords remain in create_table.go
- [x] All CREATE TABLE tests still pass

### 3.3 eqFold Migration — grant.go

- [x] `eqFold("account")` → kwACCOUNT token check
- [x] `eqFold("option")` → kwOPTION token check
- [x] `eqFold("proxy")` → kwPROXY token check — NOTE: 2 remaining eqFold for "PROXY" compare against extracted string array, not current token
- [x] `eqFold("routine")` → kwROUTINE token check
- [x] `eqFold("expire")` → kwEXPIRE token check
- [x] `eqFold("never")` → kwNEVER token check
- [x] `eqFold("day")` → kwDAY token check
- [x] `eqFold("history")` → kwHISTORY token check
- [x] `eqFold("reuse")` → kwREUSE token check
- [x] `eqFold("x509")` → kwX509 token check
- [x] `eqFold("issuer")` → kwISSUER token check
- [x] `eqFold("subject")` → kwSUBJECT token check
- [x] `eqFold("cipher")` → kwCIPHER token check
- [x] `eqFold("attribute")` → kwATTRIBUTE token check
- [x] All remaining grant.go eqFold patterns for registered keywords migrated — `eqFold("tablespace")` fixed to kwTABLESPACE token check
- [x] Zero eqFold calls for registered keywords remain in grant.go — remaining eqFold is for non-keyword strings (max_queries_per_hour, factor, initiate, registration, etc.) and extracted-string PROXY comparison
- [x] All GRANT/USER tests still pass

### 3.4 eqFold Migration — utility.go

- [x] `eqFold("schedule")` → kwSCHEDULE token check
- [x] `eqFold("completion")` → kwCOMPLETION token check
- [x] `eqFold("preserve")` → kwPRESERVE token check
- [x] `eqFold("every")` → kwEVERY token check
- [x] `eqFold("starts")` → kwSTARTS token check
- [x] `eqFold("ends")` → kwENDS token check
- [x] `eqFold("rotate")` → kwROTATE token check
- [x] `eqFold("keyring")` → kwKEYRING token check
- [x] `eqFold("tls")` → kwTLS token check
- [x] `eqFold("concurrent")` → kwCONCURRENT token check
- [x] `eqFold("work")` → kwWORK token check
- [x] `eqFold("export")` → kwEXPORT token check
- [x] `eqFold("upgrade")` → kwUPGRADE token check
- [x] `eqFold("fast")`/`eqFold("medium")`/`eqFold("changed")` → keyword token checks
- [x] All remaining utility.go eqFold patterns for registered keywords migrated — `eqFold("channel")` fixed to kwCHANNEL token check
- [x] Zero eqFold calls for registered keywords remain in utility.go — remaining eqFold is for non-keyword compound option strings (histogram, buckets, autoextend_size, initial_size, etc.)
- [x] All utility tests still pass

### 3.5 eqFold Migration — set_show.go

- [x] `eqFold("events")` → kwEVENTS token check
- [x] `eqFold("indexes")` → kwINDEXES token check
- [x] `eqFold("grants")` → kwGRANTS token check
- [x] `eqFold("triggers")` → kwTRIGGERS token check
- [x] `eqFold("schemas")` → kwSCHEMAS token check
- [x] `eqFold("partitions")` → kwPARTITIONS token check
- [x] `eqFold("hosts")` → kwHOSTS token check
- [x] `eqFold("mutex")` → kwMUTEX token check
- [x] `eqFold("profile")` → kwPROFILE token check
- [x] `eqFold("replicas")` → kwREPLICAS token check
- [x] `eqFold("format")` → kwFORMAT token check
- [x] `eqFold("names")` → kwNAMES token check
- [x] `eqFold("code")` → kwCODE token check
- [x] `eqFold("xml")` → kwXML token check
- [x] Zero eqFold calls for registered keywords remain in set_show.go — remaining eqFold is for persist_only (not a keyword) and isIdentLike helper (used for non-keyword profile option names)
- [x] All SHOW/SET tests still pass

### 3.6 eqFold Migration — replication.go & trigger.go & signal.go

- [x] All replication.go eqFold patterns for registered keywords migrated — `eqFold("UNTIL")` fixed to kwUNTIL token check
- [x] `eqFold("stream")` → kwSTREAM token check
- [x] `eqFold("generate")` → kwGENERATE token check
- [x] All trigger.go eqFold patterns for registered keywords migrated (interval units in event scheduling)
- [x] All signal.go eqFold patterns migrated (`value`, `stacked`)
- [x] Zero eqFold calls for registered keywords remain in these files — remaining eqFold in replication.go is for non-keyword compound strings (IGNORE_SERVER_IDS, REPLICATE_REWRITE_DB, SQL_AFTER_MTS_GAPS, SOURCE_LOG_FILE, DEFAULT_AUTH, PLUGIN_DIR, IO_THREAD, SQL_THREAD)
- [x] All replication/trigger/signal tests still pass

### 3.7 eqFold Migration — expr.go & compound.go & remaining files

- [x] All expr.go eqFold patterns for registered keywords migrated — zero eqFold in expr.go
- [x] All compound.go eqFold patterns migrated — zero eqFold in compound.go
- [x] All create_function.go remaining eqFold patterns migrated — zero eqFold in create_function.go
- [x] All create_index.go remaining eqFold patterns migrated — remaining eqFold is for non-keyword option strings (KEY_BLOCK_SIZE, ENGINE_ATTRIBUTE, SECONDARY_ENGINE_ATTRIBUTE)
- [x] All create_view.go remaining eqFold patterns migrated — zero eqFold in create_view.go
- [x] All alter_table.go remaining eqFold patterns migrated — remaining eqFold is for non-keyword strings (secondary_load, secondary_unload)
- [x] All alter_misc.go remaining eqFold patterns migrated — zero eqFold in alter_misc.go
- [x] All load_data.go remaining eqFold patterns migrated — zero eqFold in load_data.go
- [x] All transaction.go eqFold patterns migrated — zero eqFold in transaction.go
- [x] All insert.go eqFold patterns migrated — zero eqFold in insert.go
- [x] All select.go eqFold patterns migrated — zero eqFold in select.go
- [x] All stmt.go eqFold patterns migrated — `eqFold("reference")` fixed to kwREFERENCE token check; remaining eqFold is for non-keyword GROUP_REPLICATION and extracted-string current_user comparison
- [x] Zero eqFold calls for registered keywords remain across entire parser (95 eqFold calls remain, all for non-keyword strings)
- [x] All tests still pass

### 3.8 Completeness Audit

- [x] Every MySQL 8.0 keyword in sql_yacc.yy that appears in omni's parser is registered as a keyword token
- [x] Every registered keyword has the correct 6-category classification matching sql_yacc.yy
- [x] Zero eqFold calls remain for strings that are registered keywords — verified: 95 eqFold calls remain, all for non-keyword strings
- [x] eqFold calls only remain for: (a) non-keyword compound option strings (key_block_size, max_rows, replication options, etc.), (b) option-name dispatch patterns in create_table.go where post-parse string matching is used, (c) `@@`-prefixed variable scope parsing in name.go (lexer emits `@@global.var` as single token)
- [~] Oracle corpus: `CREATE TABLE select (a INT)` correctly rejected — known mismatch: omni accepts, MySQL rejects (strictness enforcement pending)
- [~] Oracle corpus: `CREATE TABLE t (select INT)` correctly rejected — known mismatch: omni accepts, MySQL rejects (strictness enforcement pending)
- [~] Oracle corpus: `CREATE TABLE t (rank INT)` correctly rejected (rank is reserved) — rank registered as reserved; strictness enforcement pending
- [x] Oracle corpus: `CREATE TABLE t (status INT)` correctly accepted (status is non-reserved)
- [x] Oracle corpus: `CREATE TABLE begin (a INT)` correctly accepted (begin is ambiguous_2, allowed as ident) — verified in compare_test.go
- [x] Full test suite passes with zero regressions
