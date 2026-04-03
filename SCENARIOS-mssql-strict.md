# SCENARIOS: MSSQL L2 Strict — Option Validation

## Goal

Eliminate all "omni too permissive" oracle mismatches where omni accepts invalid option names that SQL Server 2022 rejects. Drive TestKeywordOracleOptionPositions to 0 mismatches.

## Reference

- Oracle test: `mssql/parser/oracle_test.go` → `TestKeywordOracleOptionPositions` (80 cases, 33 mismatches)
- Oracle: SQL Server 2022 via testcontainers (`SET PARSEONLY ON`)
- SqlScriptDOM OptionsHelper classes: `../SqlScriptDOM/SqlScriptDom/Parser/TSql/*.cs`
- Build check: `go build ./mssql/...`
- Full regression: `go test ./mssql/... -count=1 -short`

## Current state

0 mismatches. All 17 categories passing. TestKeywordOracleOptionPositions: 33 → 0 mismatches.

---

## Phase 1: Infrastructure (sequential)

### Section 1.1: Option validation framework

- [x] `optionSet` type defined for declaring valid option names per position
- [x] helper function `isValidOption(opts optionSet) bool` checks current token against valid set
- [x] helper function `expectOption(opts optionSet) (string, error)` consumes and validates
- [x] optionSet supports both keyword tokens and identifier strings (for unregistered option names)
- [x] enforcement test: `TestOptionValidation` verifies valid options accepted, invalid rejected
- [x] framework does not break any existing parser tests

Verification: `go build ./mssql/...` + `go test ./mssql/... -count=1 -short`

---

## Phase 2: SET and database options (independent files)

### Section 2.1: SET predicate options (declare_set.go)

- [x] `SET ANSI_NULLS ON` accepted (valid)
- [x] `SET QUOTED_IDENTIFIER OFF` accepted (valid)
- [x] `SET ARITHABORT ON` accepted (valid)
- [x] `SET SELECT ON` rejected (invalid — SELECT is not a SET option)
- [x] `SET FROM OFF` rejected (invalid)
- [x] valid option set matches SqlScriptDOM PredicateSetOptionsHelper (25 options)

Verification: oracle test `set_predicate/*` → 0 mismatches

### Section 2.2: ALTER DATABASE SET options (alter_objects.go)

- [x] `ALTER DATABASE db SET RECOVERY SIMPLE` accepted (valid)
- [x] `ALTER DATABASE db SET READ_ONLY` accepted (valid)
- [x] `ALTER DATABASE db SET ANSI_NULLS ON` accepted (valid)
- [x] `ALTER DATABASE db SET SELECT ON` rejected (invalid)
- [x] `ALTER DATABASE db SET FROM OFF` rejected (invalid)
- [x] valid option set matches SqlScriptDOM DatabaseOptionKindHelper

Verification: oracle test `db_set/*` → 0 mismatches

---

## Phase 3: Index and DDL options (independent files)

### Section 3.1: Index WITH options (create_index.go, alter_objects.go)

- [x] `WITH (FILLFACTOR = 80)` accepted (valid)
- [x] `WITH (PAD_INDEX = ON)` accepted (valid — PARSEONLY may reject due to missing table, but parse should succeed)
- [x] `WITH (SELECT = 1)` rejected (invalid)
- [x] `WITH (FROM = 1)` rejected (invalid)
- [x] valid option set matches SqlScriptDOM IndexOptionHelper + IndexStateOption

Verification: oracle test `index/*` → 0 mismatches

### Section 3.2: CREATE TABLE WITH options (create_table.go)

- [x] `WITH (MEMORY_OPTIMIZED = ON)` accepted (valid — PARSEONLY may reject due to engine, but parse should succeed)
- [x] `WITH (SELECT = ON)` rejected (invalid)
- [x] `WITH (FROM = ON)` rejected (invalid)
- [x] valid option set matches SqlScriptDOM TableOptionHelper

Verification: oracle test `create_table/*` → 0 mismatches

### Section 3.3: CREATE/ALTER PROC WITH options (create_proc.go)

- [x] `WITH RECOMPILE` accepted (valid)
- [x] `WITH ENCRYPTION` accepted (valid)
- [x] `WITH SELECT` rejected (invalid)
- [x] valid option set matches SqlScriptDOM ProcedureOptionHelper

Verification: oracle test `proc/*` → 0 mismatches

### Section 3.4: CREATE VIEW WITH options (create_view.go)

- [x] `WITH SCHEMABINDING` accepted (valid)
- [x] `WITH ENCRYPTION` accepted (valid)
- [x] `WITH SELECT` rejected (invalid)
- [x] valid option set matches SqlScriptDOM ViewOptionHelper

Verification: oracle test `view/*` → 0 mismatches

---

## Phase 4: Query options (select.go)

### Section 4.1: Query hints — OPTION clause (select.go)

- [x] `OPTION (RECOMPILE)` accepted (valid)
- [x] `OPTION (MAXDOP 1)` accepted (valid)
- [x] `OPTION (SELECT)` rejected (invalid)
- [x] `OPTION (FROM)` rejected (invalid)
- [x] valid hint set matches SqlScriptDOM optimizer hint helpers

Verification: oracle test `query_hint/*` → 0 mismatches

### Section 4.2: FOR XML/JSON options (select.go)

- [x] `FOR XML RAW, ELEMENTS` accepted (valid)
- [x] `FOR XML RAW, ROOT('r')` accepted (valid)
- [x] `FOR XML RAW, SELECT` rejected (invalid)
- [x] `FOR JSON PATH, ROOT('r')` accepted (valid)
- [x] `FOR JSON PATH, WITHOUT_ARRAY_WRAPPER` accepted (valid)
- [x] `FOR JSON PATH, SELECT` rejected (invalid)
- [x] valid option sets match SqlScriptDOM XmlForClauseOptionsHelper / JsonForClauseOptionsHelper

Verification: oracle test `for_xml/*` + `for_json/*` → 0 mismatches

### Section 4.3: Cursor options (cursor.go)

- [x] `CURSOR FAST_FORWARD FOR` accepted (valid)
- [x] `CURSOR SCROLL FOR` accepted (valid)
- [x] `CURSOR STATIC FOR` accepted (valid)
- [x] `CURSOR SELECT FOR` rejected (invalid)
- [x] valid option set matches SqlScriptDOM CursorOptionsHelper (13 options)

Verification: oracle test `cursor/*` → 0 mismatches

---

## Phase 5: Backup, restore, bulk operations (independent files)

### Section 5.1: Backup WITH options (backup_restore.go)

- [x] `WITH COMPRESSION` accepted (valid)
- [x] `WITH INIT` accepted (valid)
- [x] `WITH SELECT` rejected (invalid)
- [x] `WITH FROM` rejected (invalid)
- [x] valid option set matches SqlScriptDOM BackupOptionsNoValueHelper + BackupOptionsWithValueHelper

Verification: oracle test `backup/*` → 0 mismatches

### Section 5.2: Restore WITH options (backup_restore.go)

- [x] `WITH NORECOVERY` accepted (valid)
- [x] `WITH REPLACE` accepted (valid)
- [x] `WITH SELECT` rejected (invalid)
- [x] `WITH FROM` rejected (invalid)
- [x] valid option set matches SqlScriptDOM RestoreOptionNoValueHelper + RestoreOptionWithValueHelper

Verification: oracle test `restore/*` → 0 mismatches

### Section 5.3: Bulk insert WITH options (bulk_insert.go)

- [x] `WITH (FIELDTERMINATOR = ',')` accepted (valid)
- [x] `WITH (ROWTERMINATOR = '\n')` accepted (valid)
- [x] `WITH (SELECT = 1)` rejected (invalid)
- [x] `WITH (FROM = 1)` rejected (invalid)
- [x] valid option set matches SqlScriptDOM BulkInsertFlagOptionsHelper + IntOptionHelper + StringOptionHelper

Verification: oracle test `bulk_insert/*` → 0 mismatches

---

## Phase 6: Server and HA options (independent files)

### Section 6.1: Fulltext index options (fulltext.go)

- [x] `WITH (CHANGE_TRACKING = AUTO)` accepted (valid)
- [x] `WITH (SELECT = ON)` rejected (invalid)
- [x] valid option set matches SqlScriptDOM fulltext-related helpers

Verification: oracle test `fulltext/*` → 0 mismatches

### Section 6.2: Service broker options (service_broker.go)

- [x] `CREATE MESSAGE TYPE msg VALIDATION = NONE` accepted (valid)
- [x] `CREATE SERVICE svc ON QUEUE dbo.q (SELECT)` rejected (invalid — SELECT is not a contract name pattern)
- [x] valid patterns match SqlScriptDOM service broker grammar

Verification: oracle test `broker/*` → 0 mismatches

### Section 6.3: Availability group options (availability.go)

- [x] `WITH (AUTOMATED_BACKUP_PREFERENCE = SECONDARY)` accepted (valid)
- [x] `WITH (SELECT = ON)` rejected (invalid)
- [x] valid option set matches SqlScriptDOM AvailabilityReplicaOptionsHelper

Verification: oracle test `ag/*` → 0 mismatches

### Section 6.4: Endpoint options (endpoint.go)

- [x] `STATE = STARTED AS TCP (LISTENER_PORT = 5022)` accepted (valid)
- [x] `SELECT = STARTED AS TCP` rejected (invalid — SELECT is not an endpoint option)
- [x] valid option set matches SqlScriptDOM EndpointProtocolOptionsHelper

Verification: oracle test `endpoint/*` → 0 mismatches

---

## Proof

### Section-local proof
Each section: its oracle test subcases show 0 mismatches.

### Global proof
After all sections: `TestKeywordOracleOptionPositions` → 0 mismatches total. `go test ./mssql/... -count=1 -short` all green.
