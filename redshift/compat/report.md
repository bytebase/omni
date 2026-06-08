# Redshift Compatibility Report

## AWS Command Coverage

Source: https://docs.aws.amazon.com/redshift/latest/dg/c_SQL_commands.html

| Status | Count |
|---|---:|
| supported_runtime | 27 |
| supported_parse | 91 |
| explicit_unsupported | 6 |
| not_relevant | 0 |

## Legacy Corpus

| Metric | Count |
|---|---:|
| total files | 115 |
| passing files | 109 |
| expected failure files | 6 |
| new failure files | 0 |
| promoted files | 0 |

## Legacy Statement Parity

| Metric | Count |
|---|---:|
| total statements | 5075 |
| both accept | 5066 |
| both reject | 0 |
| legacy accepts, omni rejects | 9 |
| legacy rejects, omni accepts | 0 |

### Legacy Accepts, Omni Rejects

- alter_role.sql[28]: ALTER ROLE role1 EXTERNALID TO "";
- create_role.sql[36]: -- Empty external ID (edge case) CREATE ROLE empty_external_role EXTERNALID "";
- drop_materialized_view.sql[114]: -- DROP MATERIALIZED VIEW statements without semicolons (should also be valid) DROP MATERIALIZED VIEW test_schema.test_mv DROP MATERIALIZED VIEW IF EXISTS test_...
- drop_role.sql[66]: -- DROP ROLE without semicolon (should also work) DROP ROLE sample_role1 DROP ROLE sample_role1 FORCE DROP ROLE sample_role1 RESTRICT -- Edge cases and boundary...
- grant.sql[63]: GRANT ALL ON TABLE customers TO user1 WITH GRANT OPTION, ROLE customer_admin, GROUP support_team;
- grant.sql[87]: GRANT USAGE ON SCHEMA "" TO edge_case_user;
- revoke.sql[163]: REVOKE USAGE ON SCHEMA "" FROM edge_case_user;
- revoke.sql[199]: -- Conditional revocations based on object existence REVOKE SELECT ON TABLE IF EXISTS optional_table FROM optional_user;
- revoke.sql[200]: REVOKE USAGE ON SCHEMA IF EXISTS optional_schema FROM optional_role;

## Runtime Semantics

| Metric | Count |
|---|---:|
| total checks | 8 |
| passed checks | 8 |
| failed checks | 0 |

## Reference Redshift

Status: skipped (REDSHIFT_COMPAT_DSN is not set)

## Command Details

| Command | Status | Parse |
|---|---|---|
| ABORT | supported_parse | ok |
| ALTER DATABASE | supported_parse | ok |
| ALTER DATASHARE | supported_parse | ok |
| ALTER DEFAULT PRIVILEGES | supported_parse | ok |
| ALTER EXTERNAL SCHEMA | supported_parse | ok |
| ALTER EXTERNAL VIEW | supported_parse | ok |
| ALTER FUNCTION | supported_parse | ok |
| ALTER GROUP | supported_parse | ok |
| ALTER IDENTITY PROVIDER | supported_parse | ok |
| ALTER MASKING POLICY | supported_parse | ok |
| ALTER MATERIALIZED VIEW | supported_parse | ok |
| ALTER RLS POLICY | supported_parse | ok |
| ALTER ROLE | supported_parse | ok |
| ALTER PROCEDURE | supported_parse | ok |
| ALTER SCHEMA | supported_parse | ok |
| ALTER SYSTEM | supported_parse | ok |
| ALTER TABLE | supported_runtime | ok |
| ALTER TABLE APPEND | supported_parse | ok |
| ALTER TEMPLATE | explicit_unsupported | error |
| ALTER USER | supported_parse | ok |
| ANALYZE | supported_parse | ok |
| ANALYZE COMPRESSION | supported_parse | ok |
| ATTACH MASKING POLICY | supported_parse | ok |
| ATTACH RLS POLICY | supported_parse | ok |
| BEGIN | supported_parse | ok |
| CALL | supported_parse | ok |
| CANCEL | supported_parse | ok |
| CLOSE | supported_parse | ok |
| COMMENT | supported_runtime | ok |
| COMMIT | supported_parse | ok |
| COPY | supported_runtime | ok |
| CREATE DATABASE | supported_parse | ok |
| CREATE DATASHARE | supported_parse | ok |
| CREATE EXTERNAL FUNCTION | supported_parse | ok |
| CREATE EXTERNAL MODEL | supported_parse | ok |
| CREATE EXTERNAL SCHEMA | supported_parse | ok |
| CREATE EXTERNAL TABLE | supported_parse | ok |
| CREATE EXTERNAL VIEW | supported_parse | ok |
| CREATE FUNCTION | supported_parse | ok |
| CREATE GROUP | supported_parse | ok |
| CREATE IDENTITY PROVIDER | supported_parse | ok |
| CREATE LIBRARY | supported_parse | ok |
| CREATE MASKING POLICY | supported_parse | ok |
| CREATE MATERIALIZED VIEW | supported_parse | ok |
| CREATE MODEL | supported_parse | ok |
| CREATE PROCEDURE | supported_parse | ok |
| CREATE RLS POLICY | supported_parse | ok |
| CREATE ROLE | supported_parse | ok |
| CREATE SCHEMA | supported_runtime | ok |
| CREATE TABLE | supported_runtime | ok |
| CREATE TABLE AS | supported_parse | ok |
| CREATE TEMPLATE | explicit_unsupported | error |
| CREATE USER | supported_parse | ok |
| CREATE VIEW | supported_parse | ok |
| DEALLOCATE | supported_parse | ok |
| DECLARE | supported_parse | ok |
| DELETE | supported_runtime | ok |
| DESC DATASHARE | supported_parse | ok |
| DESC IDENTITY PROVIDER | supported_parse | ok |
| DETACH MASKING POLICY | supported_parse | ok |
| DETACH RLS POLICY | supported_parse | ok |
| DROP DATABASE | supported_parse | ok |
| DROP DATASHARE | supported_parse | ok |
| DROP EXTERNAL VIEW | supported_parse | ok |
| DROP FUNCTION | supported_parse | ok |
| DROP GROUP | supported_parse | ok |
| DROP IDENTITY PROVIDER | supported_parse | ok |
| DROP LIBRARY | supported_parse | ok |
| DROP MASKING POLICY | supported_parse | ok |
| DROP MODEL | supported_parse | ok |
| DROP MATERIALIZED VIEW | supported_parse | ok |
| DROP PROCEDURE | supported_parse | ok |
| DROP RLS POLICY | supported_parse | ok |
| DROP ROLE | supported_parse | ok |
| DROP SCHEMA | supported_runtime | ok |
| DROP TABLE | supported_runtime | ok |
| DROP TEMPLATE | explicit_unsupported | error |
| DROP USER | supported_parse | ok |
| DROP VIEW | supported_parse | ok |
| END | supported_parse | ok |
| EXECUTE | supported_parse | ok |
| EXPLAIN | supported_runtime | ok |
| FETCH | supported_parse | ok |
| GRANT | supported_parse | ok |
| INSERT | supported_runtime | ok |
| INSERT (external table) | supported_parse | ok |
| LOCK | supported_parse | ok |
| MERGE | supported_runtime | ok |
| PREPARE | supported_parse | ok |
| REFRESH MATERIALIZED VIEW | supported_parse | ok |
| RESET | supported_parse | ok |
| REVOKE | supported_parse | ok |
| ROLLBACK | supported_parse | ok |
| SELECT | supported_runtime | ok |
| SELECT INTO | supported_runtime | ok |
| SET | supported_runtime | ok |
| SET SESSION AUTHORIZATION | supported_parse | ok |
| SET SESSION CHARACTERISTICS | supported_parse | ok |
| SHOW | supported_runtime | ok |
| SHOW COLUMN GRANTS | explicit_unsupported | error |
| SHOW COLUMNS | supported_runtime | ok |
| SHOW CONSTRAINTS | explicit_unsupported | error |
| SHOW EXTERNAL TABLE | supported_runtime | ok |
| SHOW DATABASES | supported_runtime | ok |
| SHOW FUNCTIONS | supported_parse | ok |
| SHOW GRANTS | supported_runtime | ok |
| SHOW MODEL | supported_parse | ok |
| SHOW DATASHARES | supported_runtime | ok |
| SHOW PARAMETERS | supported_parse | ok |
| SHOW POLICIES | supported_parse | ok |
| SHOW PROCEDURE | supported_parse | ok |
| SHOW PROCEDURES | supported_parse | ok |
| SHOW SCHEMAS | supported_runtime | ok |
| SHOW TABLE | supported_runtime | ok |
| SHOW TABLES | supported_runtime | ok |
| SHOW TEMPLATE | explicit_unsupported | error |
| SHOW TEMPLATES | supported_parse | ok |
| SHOW VIEW | supported_runtime | ok |
| START TRANSACTION | supported_parse | ok |
| TRUNCATE | supported_runtime | ok |
| UNLOAD | supported_runtime | ok |
| UPDATE | supported_runtime | ok |
| USE | supported_parse | ok |
| VACUUM | supported_parse | ok |
