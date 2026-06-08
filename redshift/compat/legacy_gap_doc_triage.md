# Redshift Legacy Gap Documentation Triage

This triage answers one question: when the legacy ANTLR Redshift parser accepts a statement and the current Omni Redshift fork rejects it, is the gap backed by Amazon Redshift documentation and worth supporting?

Snapshot source:

- Gap set: `EvaluateLegacyStatementParity("redshift/parser/testdata/legacy")`
- Initial legacy-accepts/Omni-rejects: 828 statements
- Current legacy-accepts/Omni-rejects after the Redshift documented-tail tranche: 9 statements
- Official command source: Amazon Redshift SQL command documentation

Resolved in the first implementation tranche:

- `SELECT TOP`
- `MINUS` as Redshift's synonym for `EXCEPT`
- `SELECT TOP ... INTO`
- `CREATE VIEW ... WITH NO SCHEMA BINDING`
- Promoted legacy files: `create_view.sql`, `select.sql`, `select_into.sql`

Resolved in the schema quota tranche:

- `CREATE SCHEMA ... QUOTA`
- `ALTER SCHEMA ... QUOTA`
- Promoted legacy files: `create_schema.sql`, `alter_schema.sql`

Resolved in the ALTER TABLE APPEND tranche:

- `ALTER TABLE ... APPEND FROM ... [IGNOREEXTRA | FILLTARGET]`
- Promoted legacy file: `alter_table_append.sql`

Resolved in the DROP SCHEMA external database tranche:

- `DROP SCHEMA ... DROP EXTERNAL DATABASE`
- Promoted legacy file: `drop_schema.sql`

Resolved in the ALTER MATERIALIZED VIEW tranche:

- `ALTER MATERIALIZED VIEW ... AUTO REFRESH`
- `ALTER MATERIALIZED VIEW ... ALTER DISTKEY|DISTSTYLE|SORTKEY`
- `ALTER MATERIALIZED VIEW ... ROW LEVEL SECURITY ...`
- Promoted legacy file: `alter_materialized_view.sql`

Resolved in the CREATE MATERIALIZED VIEW tranche:

- `CREATE MATERIALIZED VIEW ... BACKUP`
- `CREATE MATERIALIZED VIEW ... DISTSTYLE|DISTKEY|SORTKEY`
- `CREATE MATERIALIZED VIEW ... AUTO REFRESH`
- Promoted legacy file: `create_materialized_view.sql`

Resolved in the REFRESH MATERIALIZED VIEW tranche:

- `REFRESH MATERIALIZED VIEW ... [RESTRICT | CASCADE]`
- Promoted legacy file: `refresh_materialized_view.sql`

Resolved in the Redshift table/type tranche:

- `VARCHAR(MAX)` and Redshift character/binary/special types: `NVARCHAR`, `BPCHAR`, `VARBYTE`, `VARBINARY`, `BINARY VARYING`, `SUPER`, `GEOMETRY`, `GEOGRAPHY`, `HLLSKETCH`
- `ALTER TABLE ... ADD COLUMN ... ENCODE`, `ALTER COLUMN ... ENCODE`, `ALTER ENCODE AUTO`
- `ALTER TABLE ... MASKING ON|OFF FOR DATASHARES`
- Redshift external table maintenance: `SET TABLE PROPERTIES`, `SET LOCATION`, `SET FILE FORMAT`, `ADD/DROP PARTITION`
- Promoted legacy files: `alter_function.sql`, `alter_table.sql`, `create_table.sql`

Resolved in the Redshift database DDL tranche:

- `ALTER DATABASE ... COLLATE { CASE_SENSITIVE | CS | CASE_INSENSITIVE | CI }`
- `CREATE DATABASE ... COLLATE ...`
- `CREATE DATABASE ... ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }`
- Promoted legacy files: `alter_database.sql`, `create_database.sql`

Resolved in the Redshift role DDL tranche:

- `CREATE ROLE ... EXTERNALID`
- `ALTER ROLE ... [WITH] OWNER TO ...`
- `ALTER ROLE ... [WITH] EXTERNALID TO ...`
- `ALTER ROLE ... [WITH] RENAME TO ...`
- `DROP ROLE ... FORCE|RESTRICT`
- Statement-level rejects reduced from 570 to 427. No role fixture file was promoted because the remaining rejected cases are not the documented role syntax: `EXTERNALID ""` is a zero-length delimited identifier edge case, and `drop_role.sql` contains multiple `DROP ROLE` statements without semicolon separators.

Resolved in the Redshift user DDL tranche:

- `CREATE USER` Redshift options: `EXTERNALID`, `PASSWORD DISABLE`, `SYSLOG ACCESS`, `SESSION TIMEOUT`, `CONNECTION LIMIT UNLIMITED`, and external identity usernames containing `:`
- `ALTER USER` Redshift options: `EXTERNALID`, `PASSWORD DISABLE`, `SYSLOG ACCESS`, `SESSION TIMEOUT`, `RESET SESSION TIMEOUT`, `CONNECTION LIMIT UNLIMITED`, and user-level `SET`/`RESET` parameters
- Promoted legacy files: `create_user.sql`, `alter_user.sql`
- Statement-level rejects reduced from 427 to 394.

Resolved in the Redshift ANALYZE maintenance tranche:

- `ANALYZE ... PREDICATE COLUMNS`
- `ANALYZE ... ALL COLUMNS`
- `ANALYZE COMPRESSION [table[(columns)]] [COMPROWS numrows]`
- Promoted legacy files: `analyze.sql`, `analyze_compression.sql`
- Statement-level rejects reduced from 394 to 356.

Resolved in the Redshift USE tranche:

- `USE database`
- Promoted legacy file: `use.sql`
- Statement-level rejects reduced from 356 to 316.

Resolved in the Redshift COPY tranche:

- Documented Redshift COPY authorization/data-source/format/conversion/load options including `ACCESS_KEY_ID`, `SECRET_ACCESS_KEY`, `SESSION_TOKEN`, `MANIFEST`, `REGION`, `JSON`, `AVRO`, `FIXEDWIDTH`, `IGNOREHEADER`, `DATEFORMAT`, `TIMEFORMAT`, `MAXERROR`, `READRATIO`, `COMPUPDATE`, `STATUPDATE`, compression formats, and boolean conversion/load flags.
- Promoted legacy file: `copy.sql`
- Statement-level rejects reduced from 145 to 111.

Resolved in the Redshift UNLOAD tranche:

- Documented Redshift UNLOAD format/output/security/file options including `PARTITION BY`, `INCLUDE`, `MANIFEST VERBOSE`, `HEADER`, `FIXEDWIDTH`, `ENCRYPTED AUTO`, `KMS_KEY_ID`, `ADDQUOTES`, `ALLOWOVERWRITE`, `CLEANPATH`, `PARALLEL`, `MAXFILESIZE`, `ROWGROUPSIZE`, `REGION`, and `EXTENSION`.
- Promoted legacy file: `unload.sql`
- Statement-level rejects reduced from 111 to 90.

Resolved in the Redshift SHOW variable tranche:

- `SHOW current_user`
- `SHOW session_user`
- Promoted legacy file: `show.sql`
- Statement-level rejects reduced from 90 to 87.

Resolved in the Redshift documented-tail tranche:

- Redshift cloud-integration utility DDL: `CREATE EXTERNAL FUNCTION`, `CREATE EXTERNAL MODEL`, and `CREATE [OR REPLACE] LIBRARY`
- Redshift utility command: `CANCEL process_id [message]`
- Redshift procedure option: `CREATE PROCEDURE ... NONATOMIC`
- Redshift CTAS temp shorthand: `CREATE TABLE #name AS ...`
- Redshift DELETE optional `FROM`: `DELETE table_name`
- Promoted legacy files: `create_external_function.sql`, `create_external_model.sql`, `c_create_commands.sql`, `create_library.sql`, `cancel.sql`, `create_procedure.sql`, `create_table_as.sql`, `delete.sql`
- Statement-level rejects reduced from 87 to 9.

## Verdict Summary

Most gaps are real Redshift syntax, not legacy-parser over-acceptance. The right response is not to blindly make all 828 pass with rich semantics; instead, support them by tier:

- **Rich support:** parse into useful AST/runtime classification because Bytebase/Omni consumers can use the structure.
- **Parse-only support:** accept syntax and classify conservatively, but do not attempt query span or changed-resource semantics beyond top-level object names.
- **Defer/manual review:** documented but low value for Bytebase workflows, or corpus looks malformed.

## Triage Table

| Gap family | Count | Documentation | Verdict | Suggested support |
|---|---:|---|---|---|
| `REVOKE` privilege variants | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_REVOKE.html | Resolved for documented Redshift forms including column privileges, `DATASHARE`, `MODEL`, explicit `ROLE`, system permissions, ASSUMEROLE, external objects, copy jobs, and scoped permissions. Three `revoke.sql` fixture edge cases remain: `SCHEMA ""` and non-documented `IF EXISTS` object targets. | Rich enough to parse grantees, object class, and privilege list; no query span. |
| `VACUUM` Redshift modes | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_VACUUM_command.html | Resolved. Redshift modes, thresholds, and `BOOST` now parse as maintenance utility syntax. | Parse-only support in `VacuumStmt.Options`; no Bytebase resource semantics. |
| `CREATE EXTERNAL MODEL` | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_create_external_model.html | Resolved. Current Bedrock-style syntax and legacy S3-style fixture syntax parse as Redshift object DDL. | Parse-only, explicitly no query span. |
| Additional `COPY` options | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_COPY.html | Resolved. Documented Redshift COPY data loading options now parse; target table and top-level options are retained. | Parse-only/rich enough for target table and top-level options. |
| `CREATE EXTERNAL FUNCTION` | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_FUNCTION.html | Resolved. Lambda UDF syntax parses as Redshift object DDL with function name and raw options retained. | Parse-only plus function name/signature. |
| Additional `UNLOAD` options | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_UNLOAD.html | Resolved. Documented Redshift UNLOAD options now parse and are retained as top-level option nodes. | Parse-only. |
| `GRANT` privilege variants | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html | Resolved for documented Redshift forms including `DATASHARE`, `MODEL`, explicit `ROLE`, and `CREATE MODEL`. Two `grant.sql` fixture edge cases remain: grant option placed before more grantees, and `SCHEMA ""`. | Rich enough to parse grantees, object class, and privilege list; no query span. |
| Duplicate `CREATE EXTERNAL FUNCTION` corpus bucket | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_FUNCTION.html | Resolved by promoting the mixed `c_create_commands.sql` corpus after external function/model and OR REPLACE object DDL support. | Merge into external-function parse-only work. |
| `CANCEL` | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_CANCEL.html | Resolved. `CANCEL process_id [message]` parses as a Redshift utility statement. | Parse-only utility. |
| `SHOW current_user/session_user` | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW.html | Resolved. Redshift SHOW variable forms for reserved SQL value-function names now parse as `VariableShowStmt`. | Parse-only/read-only utility. |
| `ALTER DEFAULT PRIVILEGES` Redshift role/procedure forms | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_DEFAULT_PRIVILEGES.html | Resolved. `ROLE`/`GROUP` grantees and `ON PROCEDURES` now parse. | Parse-only privilege AST. |
| `CREATE VIEW ... WITH NO SCHEMA BINDING` | 3 | https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_VIEW.html | Real Redshift view option. Needed for schema review/query span metadata. | Rich support: view name, query, no-schema-binding flag. |
| `CREATE PROCEDURE ... NONATOMIC` | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_PROCEDURE.html | Resolved. `NONATOMIC` is retained as a procedure option. | Parse-only procedure header; keep body split handling. |
| `SELECT TOP`, `MINUS` | 3 | https://docs.aws.amazon.com/redshift/latest/dg/r_SELECT_synopsis.html | Real Redshift SELECT syntax. Needed for SQL Editor, query span, masking. | Rich SELECT support. |
| `CREATE LIBRARY` | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_LIBRARY.html | Resolved. `CREATE OR REPLACE LIBRARY` now uses the Redshift object DDL path. | Parse-only plus library name. |
| `DELETE table` without `FROM` | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_DELETE.html | Resolved. Redshift's optional `FROM` form now parses as `DeleteStmt`. | Rich DML support. |
| `DROP MATERIALIZED VIEW` malformed no-semicolon corpus case | 1 | https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_MATERIALIZED_VIEW.html | Command is real, but the captured gap looks like two statements concatenated without a separator. | Manual review before parser work. |
| `SELECT TOP ... INTO` | 1 | https://docs.aws.amazon.com/redshift/latest/dg/r_SELECT_INTO.html | Real combination of documented SELECT TOP and SELECT INTO behavior. Needed. | Rich SELECT/SELECT INTO support. |
| `CREATE TABLE #temp AS` | 0 | https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_TABLE_AS.html | Resolved. `#`-prefixed CTAS targets parse as temporary tables. | Parse as temporary CTAS. |

## Recommended Implementation Order

1. **SELECT/query correctness first**
   - `SELECT TOP`, `MINUS`, `SELECT TOP ... INTO`
   - `CREATE VIEW ... WITH NO SCHEMA BINDING`
   - Reason: these affect SQL Editor, masking/query-span, and user-visible query behavior.

2. **Remaining schema/table DDL**
   - Documented table DDL gaps are resolved; only non-documented fixture edge cases remain.
   - Reason: table DDL feeds changed-resource extraction and schema review; the main Redshift `ALTER TABLE` and type tranche is now covered.

3. **Privilege and identity DDL**
   - Permission fixture edge cases only; documented `GRANT`/`REVOKE` forms are resolved.
   - Reason: common migration/admin scripts; rich enough AST avoids treating permission scripts as unknown blobs.

4. **Data movement and maintenance**
   - COPY/UNLOAD documented options are resolved.
   - Reason: important script acceptance, but mostly parse-only for Bytebase semantics.

5. **Cloud-integration DDL and utility tail**
   - Documented cloud-integration DDL and utility command gaps are resolved.
   - Reason: documented scripts should parse, but rich runtime semantics are low value.

## Cases Not To Promote Blindly

- `DROP MATERIALIZED VIEW` has one gap that appears to concatenate two statements without a semicolon. Verify against the fixture before adding grammar just for that exact text.
- Role DDL's documented syntax is now covered. Do not promote `create_role.sql`, `alter_role.sql`, or `drop_role.sql` just to accept `EXTERNALID ""` or adjacent `DROP ROLE` statements without semicolon separators unless a real Redshift reference run proves those cases are accepted.
- `grant.sql` and `revoke.sql` remaining gaps are not documented privilege syntax: zero-length delimited schema identifiers, `WITH GRANT OPTION` before additional grantees, and `IF EXISTS` object targets.
- External function/model/library statements are documented and now parse, but their option payloads are AWS-service integrations. Prefer raw option capture after the object name over deeply typed semantics.
- Maintenance commands should not be allowed to masquerade as read-only SQL Editor queries just because they parse.
