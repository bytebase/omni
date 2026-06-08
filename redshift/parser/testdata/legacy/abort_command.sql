-- ABORT command test cases based on AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ABORT.html

-- Basic ABORT syntax
ABORT;

-- ABORT with WORK keyword
ABORT WORK;

-- ABORT with TRANSACTION keyword
ABORT TRANSACTION;

-- Example from AWS documentation: Create table and abort transaction
CREATE TABLE movie_gross(
    name varchar(30),
    gross bigint
);

BEGIN;

INSERT INTO movie_gross VALUES ( 'Raiders of the Lost Ark', 23400000);
INSERT INTO movie_gross VALUES ( 'Star Wars', 10000000 );

-- Abort the transaction (all inserts will be rolled back)
ABORT;

-- Test ABORT after multiple operations
BEGIN;
CREATE TABLE test_abort (id int, name varchar(50));
INSERT INTO test_abort VALUES (1, 'test1');
INSERT INTO test_abort VALUES (2, 'test2');
UPDATE test_abort SET name = 'updated' WHERE id = 1;
DELETE FROM test_abort WHERE id = 2;
ABORT;

-- ABORT with WORK in a transaction
BEGIN TRANSACTION;
INSERT INTO movie_gross VALUES ('The Empire Strikes Back', 20000000);
ABORT WORK;

-- ABORT with TRANSACTION keyword in a transaction
START TRANSACTION;
DELETE FROM movie_gross WHERE name = 'Star Wars';
ABORT TRANSACTION;

-- Multiple ABORTs in sequence (each outside transaction)
ABORT;
ABORT WORK;
ABORT TRANSACTION;

-- ABORT in nested transactions with savepoints
BEGIN;
INSERT INTO movie_gross VALUES ('Return of the Jedi', 15000000);
SAVEPOINT sp1;
INSERT INTO movie_gross VALUES ('The Phantom Menace', 10000000);
SAVEPOINT sp2;
INSERT INTO movie_gross VALUES ('Attack of the Clones', 8000000);
ABORT;

-- Mixed transaction control statements
BEGIN;
INSERT INTO movie_gross VALUES ('Revenge of the Sith', 12000000);
-- Can use ABORT instead of ROLLBACK
ABORT;

-- Verify ABORT works the same as ROLLBACK
BEGIN WORK;
CREATE TEMP TABLE abort_test (col1 int);
INSERT INTO abort_test VALUES (1), (2), (3);
ABORT WORK;

-- ABORT after DDL operations
BEGIN;
CREATE TABLE abort_ddl_test (id int, data text);
ALTER TABLE abort_ddl_test ADD COLUMN created_at timestamp;
CREATE INDEX idx_abort_test ON abort_ddl_test(id);
ABORT TRANSACTION;

-- Clean up
DROP TABLE IF EXISTS movie_gross;
DROP TABLE IF EXISTS test_abort;
DROP TABLE IF EXISTS abort_test;
DROP TABLE IF EXISTS abort_ddl_test;