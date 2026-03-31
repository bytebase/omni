-- @name: grant select
-- @valid: true
-- @source: Oracle SQL Reference
GRANT SELECT ON employees TO hr_user

-- @name: grant multiple
-- @valid: true
-- @source: Oracle SQL Reference
GRANT SELECT, INSERT, UPDATE ON employees TO app_user

-- @name: grant with grant option
-- @valid: true
-- @source: Oracle SQL Reference
GRANT SELECT ON employees TO hr_manager WITH GRANT OPTION

-- @name: grant role
-- @valid: true
-- @source: Oracle SQL Reference
GRANT dba TO admin_user

-- @name: revoke select
-- @valid: true
-- @source: Oracle SQL Reference
REVOKE SELECT ON employees FROM hr_user

-- @name: create user
-- @valid: true
-- @source: Oracle SQL Reference
CREATE USER test_user IDENTIFIED BY password123

-- @name: alter user
-- @valid: true
-- @source: Oracle SQL Reference
ALTER USER test_user IDENTIFIED BY new_password

-- @name: drop user
-- @valid: true
-- @source: Oracle SQL Reference
DROP USER test_user CASCADE

-- @name: create role
-- @valid: true
-- @source: Oracle SQL Reference
CREATE ROLE app_role

-- @name: create profile
-- @valid: true
-- @source: Oracle SQL Reference
CREATE PROFILE secure_profile LIMIT
  FAILED_LOGIN_ATTEMPTS 3
  PASSWORD_LOCK_TIME 1

-- @name: alter session set
-- @valid: true
-- @source: Oracle SQL Reference
ALTER SESSION SET NLS_DATE_FORMAT = 'YYYY-MM-DD'

-- @name: set transaction
-- @valid: true
-- @source: Oracle SQL Reference
SET TRANSACTION READ ONLY

-- @name: commit
-- @valid: true
-- @source: Oracle SQL Reference
COMMIT

-- @name: rollback
-- @valid: true
-- @source: Oracle SQL Reference
ROLLBACK

-- @name: savepoint
-- @valid: true
-- @source: Oracle SQL Reference
SAVEPOINT sp1

-- @name: rollback to savepoint
-- @valid: true
-- @source: Oracle SQL Reference
ROLLBACK TO SAVEPOINT sp1

-- @name: explain plan
-- @valid: true
-- @source: Oracle SQL Reference
EXPLAIN PLAN FOR SELECT * FROM employees

-- @name: flashback table
-- @valid: true
-- @source: Oracle SQL Reference
FLASHBACK TABLE employees TO TIMESTAMP SYSTIMESTAMP - INTERVAL '1' HOUR

-- @name: purge recyclebin
-- @valid: true
-- @source: Oracle SQL Reference
PURGE RECYCLEBIN

-- @name: analyze table
-- @valid: true
-- @source: Oracle SQL Reference
ANALYZE TABLE employees COMPUTE STATISTICS

-- @name: create tablespace
-- @valid: true
-- @source: Oracle SQL Reference
CREATE TABLESPACE app_data
DATAFILE '/u01/app_data01.dbf' SIZE 100M
AUTOEXTEND ON NEXT 50M

-- @name: alter tablespace
-- @valid: true
-- @source: Oracle SQL Reference
ALTER TABLESPACE app_data ADD DATAFILE '/u01/app_data02.dbf' SIZE 100M

-- @name: create database link
-- @valid: true
-- @source: Oracle SQL Reference
CREATE DATABASE LINK remote_db
CONNECT TO remote_user IDENTIFIED BY remote_pass
USING 'remote_tns'
