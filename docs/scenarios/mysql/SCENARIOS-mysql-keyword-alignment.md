# MySQL Parser Keyword Alignment Scenarios

> Goal: Fully align omni's MySQL parser keyword set with MySQL 8.0's 252 reserved + 482 non-reserved classification. Eliminate eqFold workarounds for all MySQL keywords.
> Verification: `go test ./mysql/parser/... -count=1` — every scenario is a compile+test assertion
> Reference sources: MySQL 8.0 keyword list (https://dev.mysql.com/doc/refman/8.0/en/keywords.html)

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: eqFold Cleanup

Remove redundant eqFold checks for keywords already registered in lexer.go. These words already lex as keyword tokens, so the `tokIDENT && eqFold(...)` fallback is dead code.

### 1.1 set_show.go eqFold Cleanup

- [ ] `processlist` — remove eqFold, use kwPROCESSLIST token check only
- [ ] `keys` — remove eqFold, use kwKEYS token check only
- [ ] `fields` — remove eqFold, use kwFIELDS token check only
- [ ] `source` — remove eqFold, use kwSOURCE token check only
- [ ] `connection` — remove eqFold, use kwCONNECTION token check only
- [ ] `password` (in set_show.go) — remove eqFold, use kwPASSWORD token check only
- [ ] All existing SHOW/SET tests still pass after cleanup

### 1.2 trigger.go eqFold Cleanup

- [ ] `each` — remove eqFold, use kwEACH token check only
- [ ] `follows` — remove eqFold, use kwFOLLOWS token check only
- [ ] `precedes` — remove eqFold, use kwPRECEDES token check only
- [ ] `slave` (in trigger.go) — remove eqFold, use kwSLAVE token check only
- [ ] All existing trigger tests still pass after cleanup

### 1.3 create_function.go eqFold Cleanup

- [ ] `deterministic` — remove eqFold, use kwDETERMINISTIC token check only
- [ ] `security` (in create_function.go) — remove eqFold, use kwSECURITY token check only
- [ ] `contains` — remove eqFold, use kwCONTAINS token check only
- [ ] `reads` — remove eqFold, use kwREADS token check only
- [ ] `data` — remove eqFold, use kwDATA token check only
- [ ] `modifies` — remove eqFold, use kwMODIFIES token check only
- [ ] All existing CREATE FUNCTION/PROCEDURE tests still pass

### 1.4 create_table.go & create_view.go & create_index.go eqFold Cleanup

- [ ] `always` — remove eqFold, use kwALWAYS token check only
- [ ] `security` (in create_view.go) — remove eqFold, use kwSECURITY token check only
- [ ] `VISIBLE` — remove eqFold, use kwVISIBLE token check only
- [ ] `INVISIBLE` — remove eqFold, use kwINVISIBLE token check only
- [ ] All existing CREATE TABLE/VIEW/INDEX tests still pass

### 1.5 grant.go eqFold Cleanup

- [ ] `privileges` — remove eqFold, use kwPRIVILEGES token check only
- [ ] `tables` — remove eqFold, use kwTABLES token check only
- [ ] `slave` (in grant.go) — remove eqFold, use kwSLAVE token check only
- [ ] `current` — remove eqFold, use kwCURRENT token check only
- [ ] `unbounded` — remove eqFold, use kwUNBOUNDED token check only
- [ ] All existing GRANT/REVOKE tests still pass

### 1.6 utility.go & replication.go & misc eqFold Cleanup

- [ ] `logs` — remove eqFold, use kwLOGS token check only
- [ ] `group` (in utility.go resource group) — remove eqFold, use kwGROUP token check only
- [ ] `user` (in replication.go) — remove eqFold, use kwUSER token check only
- [ ] `password` (in replication.go) — remove eqFold, use kwPASSWORD token check only
- [ ] `slave` (in alter_misc.go) — remove eqFold, use kwSLAVE token check only
- [ ] `security` (in alter_misc.go) — remove eqFold, use kwSECURITY token check only
- [ ] `columns` (in load_data.go) — remove eqFold, use kwCOLUMNS token check only
- [ ] `global`/`session`/`local` (in name.go) — fix to check keyword token, not eqFold on tokIDENT
- [ ] All existing utility/replication/misc tests still pass

---

## Phase 2: Non-Reserved Keyword Registration

Register MySQL 8.0 non-reserved keywords in lexer.go and migrate eqFold patterns to keyword token checks.

### 2.1 Type & Spatial Keywords

- [ ] Register `geometry` as kwGEOMETRY, migrate eqFold in type.go
- [ ] Register `point` as kwPOINT, migrate eqFold in type.go
- [ ] Register `linestring` as kwLINESTRING, migrate eqFold in type.go
- [ ] Register `polygon` as kwPOLYGON, migrate eqFold in type.go
- [ ] Register `multipoint` as kwMULTIPOINT, migrate eqFold in type.go
- [ ] Register `multilinestring` as kwMULTILINESTRING, migrate eqFold in type.go
- [ ] Register `multipolygon` as kwMULTIPOLYGON, migrate eqFold in type.go
- [ ] Register `geometrycollection` as kwGEOMETRYCOLLECTION, migrate eqFold in type.go
- [ ] Register `serial` as kwSERIAL, migrate eqFold in type.go
- [ ] Register `national` as kwNATIONAL, migrate eqFold in type.go
- [ ] Register `nchar` as kwNCHAR, migrate eqFold in type.go
- [ ] Register `nvarchar` as kwNVARCHAR, migrate eqFold in type.go
- [ ] Register `signed` as kwSIGNED, migrate eqFold in type.go, expr.go
- [ ] Register `precision` as kwPRECISION, migrate eqFold in type.go
- [ ] Register `long` as kwLONG, migrate eqFold in type.go
- [ ] Register `int1`/`int2`/`int3`/`int4`/`int8` as keywords, migrate eqFold in type.go
- [ ] Register `middleint` as kwMIDDLEINT, migrate eqFold in type.go
- [ ] Register `float4`/`float8` as keywords, migrate eqFold in type.go
- [ ] All data type tests still pass after migration
- [ ] `CREATE TABLE t (geometry INT)` still parses (geometry is non-reserved, valid as column name)

### 2.2 DDL Option Keywords

- [ ] Register `enforced` as kwENFORCED, migrate eqFold in create_table.go/alter_table.go
- [ ] Register `less` as kwLESS, migrate eqFold in create_table.go
- [ ] Register `than` as kwTHAN, migrate eqFold in create_table.go
- [ ] Register `subpartitions` as kwSUBPARTITIONS, migrate eqFold in create_table.go
- [ ] Register `leaves` as kwLEAVES, migrate eqFold in create_table.go
- [ ] Register `srid` as kwSRID, migrate eqFold in type.go/create_table.go
- [ ] Register `parser` as kwPARSER_KW, migrate eqFold in create_index.go
- [ ] Register `compression` as kwCOMPRESSION, migrate eqFold in create_table.go
- [ ] Register `insert_method` as kwINSERT_METHOD, migrate eqFold in create_table.go
- [ ] All CREATE TABLE/ALTER TABLE/CREATE INDEX tests still pass

### 2.3 SHOW/SET Option Keywords

- [ ] Register `events` as kwEVENTS, migrate eqFold in set_show.go
- [ ] Register `indexes` as kwINDEXES, migrate eqFold in set_show.go
- [ ] Register `grants` as kwGRANTS, migrate eqFold in set_show.go
- [ ] Register `triggers` as kwTRIGGERS, migrate eqFold in set_show.go
- [ ] Register `schemas` as kwSCHEMAS, migrate eqFold in set_show.go
- [ ] Register `partitions` as kwPARTITIONS_KW, migrate eqFold in set_show.go
- [ ] Register `hosts` as kwHOSTS, migrate eqFold in set_show.go
- [ ] Register `mutex` as kwMUTEX, migrate eqFold in set_show.go
- [ ] Register `profile` as kwPROFILE, migrate eqFold in set_show.go
- [ ] Register `replicas` as kwREPLICAS, migrate eqFold in set_show.go
- [ ] Register `format` as kwFORMAT, migrate eqFold in set_show.go
- [ ] Register `names` as kwNAMES, migrate eqFold in set_show.go
- [ ] Register `persist_only` as kwPERSIST_ONLY, migrate eqFold in set_show.go
- [ ] All SHOW/SET tests still pass

### 2.4 Grant & Auth Keywords

- [ ] Register `account` as kwACCOUNT, migrate eqFold in grant.go
- [ ] Register `option` as kwOPTION, migrate eqFold in grant.go
- [ ] Register `proxy` as kwPROXY, migrate eqFold in grant.go
- [ ] Register `routine` as kwROUTINE, migrate eqFold in grant.go
- [ ] Register `expire` as kwEXPIRE, migrate eqFold in grant.go
- [ ] Register `never` as kwNEVER, migrate eqFold in grant.go
- [ ] Register `day` as kwDAY, migrate eqFold in grant.go
- [ ] Register `history` as kwHISTORY, migrate eqFold in grant.go
- [ ] Register `reuse` as kwREUSE, migrate eqFold in grant.go
- [ ] Register `optional` as kwOPTIONAL, migrate eqFold in grant.go
- [ ] Register `x509` as kwX509, migrate eqFold in grant.go
- [ ] Register `issuer` as kwISSUER, migrate eqFold in grant.go
- [ ] Register `subject` as kwSUBJECT, migrate eqFold in grant.go
- [ ] Register `cipher` as kwCIPHER, migrate eqFold in grant.go
- [ ] Register `factor` as kwFACTOR, migrate eqFold in grant.go
- [ ] Register `initiate` as kwINITIATE, migrate eqFold in grant.go
- [ ] Register `registration` as kwREGISTRATION, migrate eqFold in grant.go
- [ ] Register `finish` as kwFINISH, migrate eqFold in grant.go
- [ ] Register `unregister` as kwUNREGISTER, migrate eqFold in grant.go
- [ ] Register `initial` as kwINITIAL, migrate eqFold in grant.go
- [ ] Register `authentication` as kwAUTHENTICATION, migrate eqFold in grant.go
- [ ] Register `admin` as kwADMIN, migrate eqFold in grant.go
- [ ] Register `client` as kwCLIENT, migrate eqFold in grant.go
- [ ] All GRANT/CREATE USER/ALTER USER tests still pass

### 2.5 Scheduling & Event Keywords

- [ ] Register `schedule` as kwSCHEDULE, migrate eqFold in utility.go
- [ ] Register `completion` as kwCOMPLETION, migrate eqFold in utility.go
- [ ] Register `preserve` as kwPRESERVE, migrate eqFold in utility.go
- [ ] Register `every` as kwEVERY, migrate eqFold in utility.go
- [ ] Register `starts` as kwSTARTS, migrate eqFold in utility.go
- [ ] Register `ends` as kwENDS, migrate eqFold in utility.go
- [ ] All event/scheduling tests still pass

### 2.6 Misc Non-Reserved Keywords

- [ ] Register `action` as kwACTION, migrate eqFold in create_table.go
- [ ] Register `value` as kwVALUE, migrate eqFold in signal.go, insert.go, expr.go
- [ ] Register `stacked` as kwSTACKED, migrate eqFold in signal.go
- [ ] Register `unknown` as kwUNKNOWN, migrate eqFold in create_function.go
- [ ] Register `code` as kwCODE, migrate eqFold in set_show.go
- [ ] Register `xml` as kwXML, migrate eqFold in set_show.go
- [ ] Register `concurrent` as kwCONCURRENT, migrate eqFold in utility.go
- [ ] Register `work` as kwWORK, migrate eqFold in utility.go, transaction.go
- [ ] Register `xid` as kwXID, migrate eqFold in transaction.go
- [ ] Register `export` as kwEXPORT, migrate eqFold in utility.go
- [ ] Register `upgrade` as kwUPGRADE, migrate eqFold in utility.go
- [ ] Register `fast`/`medium`/`changed` as keywords, migrate eqFold in utility.go
- [ ] Register `wait` as kwWAIT, migrate eqFold in utility.go
- [ ] Register `active`/`inactive` as keywords, migrate eqFold in utility.go
- [ ] Register `attribute` as kwATTRIBUTE, migrate eqFold in grant.go
- [ ] Register `secondary_load` as kwSECONDARY_LOAD, migrate eqFold in alter_table.go
- [ ] Register `secondary_unload` as kwSECONDARY_UNLOAD, migrate eqFold in alter_table.go
- [ ] All misc tests still pass

### 2.7 Tablespace & Resource Group Keywords

- [ ] Register `undofile` as kwUNDOFILE, migrate eqFold in utility.go
- [ ] Register `nodegroup` as kwNODEGROUP, migrate eqFold in utility.go
- [ ] Register `extent_size` as kwEXTENT_SIZE, migrate eqFold in utility.go
- [ ] Register `initial_size` as kwINITIAL_SIZE, migrate eqFold in utility.go
- [ ] Register `max_size` as kwMAX_SIZE, migrate eqFold in utility.go
- [ ] Register `file_block_size` as kwFILE_BLOCK_SIZE, migrate eqFold in utility.go
- [ ] Register `vcpu` as kwVCPU, migrate eqFold in utility.go
- [ ] Register `thread_priority` as kwTHREAD_PRIORITY, migrate eqFold in utility.go
- [ ] Register `description` as kwDESCRIPTION, migrate eqFold in utility.go
- [ ] Register `organization` as kwORGANIZATION, migrate eqFold in utility.go
- [ ] Register `reference` as kwREFERENCE, migrate eqFold in utility.go
- [ ] Register `definition` as kwDEFINITION, migrate eqFold in utility.go
- [ ] Register `name` as kwNAME, migrate eqFold in utility.go
- [ ] Register `system` as kwSYSTEM, migrate eqFold in utility.go
- [ ] All tablespace/resource group tests still pass

---

## Phase 3: Reserved Keyword Alignment

Register missing reserved keywords and expand reservedKeywords map to match MySQL 8.0's 252 reserved words.

### 3.1 Register Missing Reserved Keywords — Window Functions

These are MySQL 8.0 reserved words not in lexer.go at all. Register and add to reservedKeywords.

- [ ] Register `cume_dist` as kwCUME_DIST (reserved)
- [ ] Register `dense_rank` as kwDENSE_RANK (reserved)
- [ ] Register `first_value` as kwFIRST_VALUE (reserved)
- [ ] Register `lag` as kwLAG (reserved)
- [ ] Register `last_value` as kwLAST_VALUE (reserved)
- [ ] Register `lead` as kwLEAD (reserved)
- [ ] Register `nth_value` as kwNTH_VALUE (reserved)
- [ ] Register `ntile` as kwNTILE (reserved)
- [ ] Register `percent_rank` as kwPERCENT_RANK (reserved)
- [ ] Register `rank` as kwRANK (reserved)
- [ ] Register `row_number` as kwROW_NUMBER (reserved)
- [ ] Window function parsing still works (these were parsed as tokIDENT function calls; now they need explicit handling in parsePrimaryExpr or remain callable as keyword-named functions)
- [ ] `CREATE TABLE t (rank INT)` rejected — rank is reserved, must be backtick-quoted

### 3.2 Register Missing Reserved Keywords — Interval Compounds

- [ ] Register `day_hour` as kwDAY_HOUR (reserved)
- [ ] Register `day_microsecond` as kwDAY_MICROSECOND (reserved)
- [ ] Register `day_minute` as kwDAY_MINUTE (reserved)
- [ ] Register `day_second` as kwDAY_SECOND (reserved)
- [ ] Register `hour_microsecond` as kwHOUR_MICROSECOND (reserved)
- [ ] Register `hour_minute` as kwHOUR_MINUTE (reserved)
- [ ] Register `hour_second` as kwHOUR_SECOND (reserved)
- [ ] Register `minute_microsecond` as kwMINUTE_MICROSECOND (reserved)
- [ ] Register `minute_second` as kwMINUTE_SECOND (reserved)
- [ ] Register `second_microsecond` as kwSECOND_MICROSECOND (reserved)
- [ ] Register `year_month` as kwYEAR_MONTH (reserved)
- [ ] INTERVAL expressions still work with these as keyword tokens (parseIntervalExpr must accept keyword tokens as unit names)

### 3.3 Register Missing Reserved Keywords — Misc

- [ ] Register `accessible` as kwACCESSIBLE (reserved)
- [ ] Register `asensitive` as kwASENSITIVE (reserved)
- [ ] Register `cube` as kwCUBE (reserved)
- [ ] Register `dual` as kwDUAL (reserved)
- [ ] Register `grouping` as kwGROUPING (reserved)
- [ ] Register `insensitive` as kwINSENSITIVE (reserved)
- [ ] Register `io_after_gtids` as kwIO_AFTER_GTIDS (reserved)
- [ ] Register `io_before_gtids` as kwIO_BEFORE_GTIDS (reserved)
- [ ] Register `maxvalue` as kwMAXVALUE (reserved) — migrate eqFold in create_table.go partition definitions
- [ ] Register `no_write_to_binlog` as kwNO_WRITE_TO_BINLOG (reserved)
- [ ] Register `of` as kwOF (reserved) — migrate eqFold in select.go, expr.go
- [ ] Register `optimizer_costs` as kwOPTIMIZER_COSTS (reserved)
- [ ] Register `sensitive` as kwSENSITIVE (reserved)
- [ ] Register `specific` as kwSPECIFIC (reserved)
- [ ] Register `sqlexception` as kwSQLEXCEPTION (reserved)
- [ ] Register `sqlstate` as kwSQLSTATE (reserved)
- [ ] Register `sqlwarning` as kwSQLWARNING (reserved)
- [ ] Register `usage` as kwUSAGE (reserved)
- [ ] Register `utc_date` as kwUTC_DATE (reserved)
- [ ] Register `utc_time` as kwUTC_TIME (reserved)
- [ ] Register `utc_timestamp` as kwUTC_TIMESTAMP (reserved)
- [ ] Register `varying` as kwVARYING (reserved)
- [ ] Existing parsing of DUAL, MAXVALUE, GROUPING, SQLSTATE, SQLWARNING, SQLEXCEPTION, UTC functions still works after keyword registration
- [ ] `SELECT 1 FROM DUAL` still parses (DUAL now a keyword token, FROM clause must accept it)

### 3.4 Existing Keywords → reservedKeywords Map — Operators & Expressions

Add these already-registered keywords to the reservedKeywords map.

- [ ] Add kwCROSS to reservedKeywords
- [ ] Add kwNATURAL to reservedKeywords
- [ ] Add kwFULL to reservedKeywords
- [ ] Add kwUSING to reservedKeywords
- [ ] Add kwASC to reservedKeywords
- [ ] Add kwDESC to reservedKeywords
- [ ] Add kwTO to reservedKeywords
- [ ] Add kwDIV to reservedKeywords
- [ ] Add kwMOD to reservedKeywords
- [ ] Add kwXOR to reservedKeywords
- [ ] Add kwREGEXP to reservedKeywords
- [ ] Add kwRLIKE to reservedKeywords
- [ ] Add kwBINARY to reservedKeywords
- [ ] Add kwINTERVAL to reservedKeywords
- [ ] Add kwMATCH to reservedKeywords
- [ ] Add kwESCAPE to reservedKeywords
- [ ] All expression/join tests still pass

### 3.5 Existing Keywords → reservedKeywords Map — Temporal & System Functions

- [ ] Add kwCURRENT_DATE to reservedKeywords
- [ ] Add kwCURRENT_TIME to reservedKeywords
- [ ] Add kwCURRENT_TIMESTAMP to reservedKeywords
- [ ] Add kwCURRENT_USER to reservedKeywords
- [ ] Add kwLOCALTIME to reservedKeywords
- [ ] Add kwLOCALTIMESTAMP to reservedKeywords
- [ ] Add kwDATABASE to reservedKeywords
- [ ] Add kwSCHEMA to reservedKeywords
- [ ] Add kwDATABASES to reservedKeywords
- [ ] All temporal function and schema tests still pass

### 3.6 Existing Keywords → reservedKeywords Map — Statements & Clauses

- [ ] Add kwREPLACE to reservedKeywords
- [ ] Add kwIGNORE to reservedKeywords
- [ ] Add kwLOAD to reservedKeywords
- [ ] Add kwUSE to reservedKeywords
- [ ] Add kwKILL to reservedKeywords
- [ ] Add kwEXPLAIN to reservedKeywords
- [ ] Add kwDESCRIBE to reservedKeywords
- [ ] Add kwWITH to reservedKeywords
- [ ] Add kwREAD to reservedKeywords
- [ ] Add kwWRITE to reservedKeywords
- [ ] Add kwCALL to reservedKeywords
- [ ] Add kwFUNCTION to reservedKeywords
- [ ] Add kwPROCEDURE to reservedKeywords
- [ ] Add kwTRIGGER to reservedKeywords
- [ ] All statement-level tests still pass

### 3.7 Existing Keywords → reservedKeywords Map — DDL & Misc

- [ ] Add kwPARTITION to reservedKeywords
- [ ] Add kwRANGE to reservedKeywords
- [ ] Add kwROW to reservedKeywords
- [ ] Add kwROWS to reservedKeywords
- [ ] Add kwOVER to reservedKeywords
- [ ] Add kwWINDOW to reservedKeywords
- [ ] Add kwSPATIAL to reservedKeywords
- [ ] Add kwFULLTEXT to reservedKeywords
- [ ] Add kwFORCE to reservedKeywords
- [ ] Add kwSTRAIGHT_JOIN to reservedKeywords
- [ ] Add kwDISTINCTROW to reservedKeywords
- [ ] Add kwHIGH_PRIORITY to reservedKeywords
- [ ] Add kwLOW_PRIORITY to reservedKeywords
- [ ] Add kwDELAYED to reservedKeywords
- [ ] Add kwSQL_CALC_FOUND_ROWS to reservedKeywords
- [ ] Add kwOUTFILE to reservedKeywords
- [ ] Add kwCONVERT to reservedKeywords
- [ ] Add kwCAST to reservedKeywords
- [ ] All DDL tests still pass

### 3.8 Existing Keywords → reservedKeywords Map — Procedural SQL

- [ ] Add kwDECLARE to reservedKeywords
- [ ] Add kwCONDITION to reservedKeywords
- [ ] Add kwCURSOR to reservedKeywords
- [ ] Add kwCONTINUE to reservedKeywords
- [ ] Add kwEXIT to reservedKeywords
- [ ] Add kwFETCH to reservedKeywords
- [ ] Add kwDO to reservedKeywords
- [ ] Add kwELSEIF to reservedKeywords
- [ ] Add kwREPEAT to reservedKeywords
- [ ] Add kwUNTIL to reservedKeywords
- [ ] Add kwLOOP to reservedKeywords
- [ ] Add kwLEAVE to reservedKeywords
- [ ] Add kwITERATE to reservedKeywords
- [ ] Add kwRETURN to reservedKeywords
- [ ] Add kwSIGNAL to reservedKeywords
- [ ] Add kwRESIGNAL to reservedKeywords
- [ ] Add kwGET to reservedKeywords
- [ ] Add kwEACH to reservedKeywords
- [ ] Add kwLATERAL to reservedKeywords
- [ ] Add kwRECURSIVE to reservedKeywords
- [ ] Add kwEMPTY to reservedKeywords
- [ ] Add kwSSL to reservedKeywords
- [ ] Add kwREQUIRE to reservedKeywords
- [ ] All procedural SQL tests still pass

### 3.9 Existing Keywords → reservedKeywords Map — LOAD DATA & String

- [ ] Add kwLINES to reservedKeywords
- [ ] Add kwFIELDS to reservedKeywords
- [ ] Add kwTERMINATED to reservedKeywords
- [ ] Add kwENCLOSED to reservedKeywords
- [ ] Add kwESCAPED to reservedKeywords
- [ ] Add kwSTARTING to reservedKeywords
- [ ] Add kwOPTIONALLY to reservedKeywords
- [ ] Add kwBOTH to reservedKeywords
- [ ] Add kwLEADING to reservedKeywords
- [ ] Add kwTRAILING to reservedKeywords
- [ ] Add kwUNDO to reservedKeywords
- [ ] All LOAD DATA and string function tests still pass

### 3.10 Test Audit & Fix

After all reserved keywords are registered, audit and fix tests that use reserved words as unquoted identifiers.

- [ ] Run full test suite, collect all new failures
- [ ] Fix test SQL: quote reserved words used as identifiers with backticks
- [ ] Verify oracle corpus tests: `CREATE TABLE select (a INT)` now correctly rejected
- [ ] Verify oracle corpus tests: `CREATE TABLE t (select INT)` now correctly rejected
- [ ] Update oracle_corpus_test.go: change knownMismatch to false for reserved word tests
- [ ] Full test suite passes with zero regressions
- [ ] `CREATE TABLE t (rank INT)` rejected — rank is now reserved
- [ ] `` CREATE TABLE t (`rank` INT) `` accepted — backtick-quoted reserved word OK
- [ ] `CREATE TABLE t (status INT)` accepted — status is non-reserved, still valid
- [ ] Completeness audit: reservedKeywords map contains all MySQL 8.0 reserved words that are registered in lexer
- [ ] Completeness audit: no eqFold calls remain for strings that are registered keywords (zero dual-check patterns)
