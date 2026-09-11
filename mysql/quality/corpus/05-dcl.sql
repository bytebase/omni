-- @name: grant select
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
GRANT SELECT ON mydb.* TO 'appuser'@'localhost'

-- @name: grant multiple privileges
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
GRANT SELECT, INSERT, UPDATE ON mydb.employees TO 'appuser'@'localhost'

-- @name: grant all privileges
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
GRANT ALL PRIVILEGES ON *.* TO 'admin'@'localhost'

-- @name: grant with grant option
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
GRANT SELECT ON mydb.* TO 'lead'@'localhost' WITH GRANT OPTION

-- @name: grant role
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
GRANT app_role TO 'user1'@'localhost'

-- @name: revoke select
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
REVOKE SELECT ON mydb.* FROM 'appuser'@'localhost'

-- @name: revoke all
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
REVOKE ALL PRIVILEGES ON *.* FROM 'appuser'@'localhost'

-- @name: create user
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
CREATE USER 'testuser'@'localhost' IDENTIFIED BY 'password123'

-- @name: create user if not exists
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
CREATE USER IF NOT EXISTS 'testuser'@'localhost' IDENTIFIED BY 'password123'

-- @name: drop user
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DROP USER 'testuser'@'localhost'

-- @name: drop user if exists
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DROP USER IF EXISTS 'testuser'@'localhost'

-- @name: alter user
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
ALTER USER 'testuser'@'localhost' IDENTIFIED BY 'newpassword'

-- @name: create role
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
CREATE ROLE app_role

-- @name: drop role
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DROP ROLE app_role

-- @name: set password
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SET PASSWORD FOR 'testuser'@'localhost' = 'newpass'

-- @name: create trigger before insert
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
CREATE TRIGGER trg_audit BEFORE INSERT ON employees FOR EACH ROW SET NEW.created_at = NOW()

-- @name: create trigger after update
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
CREATE TRIGGER trg_update AFTER UPDATE ON employees FOR EACH ROW SET NEW.updated_at = NOW()

-- @name: drop trigger
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DROP TRIGGER trg_audit

-- @name: drop trigger if exists
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DROP TRIGGER IF EXISTS trg_audit

-- @name: drop procedure
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DROP PROCEDURE IF EXISTS my_proc

-- @name: drop function
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DROP FUNCTION IF EXISTS my_func

-- @name: call procedure
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
CALL my_proc()

-- @name: call procedure with args
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
CALL my_proc(1, 'test')

-- @name: do expression
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DO 1

-- @name: do sleep
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DO SLEEP(1)

-- @name: load data infile
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
LOAD DATA INFILE '/tmp/data.csv' INTO TABLE employees

-- @name: aurora load data from s3
-- @valid: true
-- @source: Aurora MySQL User Guide (AuroraMySQL.Integrating.LoadFromS3); rejected by stock MySQL 8.0
LOAD DATA FROM S3 's3://mybucket/data.txt' INTO TABLE t1 FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n'

-- @name: aurora load data from s3 manifest
-- @valid: true
-- @source: Aurora MySQL User Guide (AuroraMySQL.Integrating.LoadFromS3); rejected by stock MySQL 8.0
LOAD DATA FROM S3 MANIFEST 's3-us-west-2://aurora-bucket/customer.manifest' INTO TABLE CUSTOMER FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n' (ID, FIRSTNAME, LASTNAME, EMAIL)

-- @name: aurora load xml from s3
-- @valid: true
-- @source: Aurora MySQL User Guide (AuroraMySQL.Integrating.LoadFromS3); rejected by stock MySQL 8.0
LOAD XML FROM S3 's3://mybucket/data.xml' INTO TABLE table1 (column1, @var1) SET table_column2 = @var1/100

-- @name: aurora load data local from s3 rejected
-- @valid: false
-- @source: Aurora MySQL User Guide: LOCAL cannot be used with an S3 source
LOAD DATA LOCAL FROM S3 's3://mybucket/data.txt' INTO TABLE t1
