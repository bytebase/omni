-- ABORT command examples from AWS Redshift documentation
-- ABORT stops the current transaction and discards all updates
-- It is equivalent to ROLLBACK

-- Example 1: Basic ABORT
-- Create a table
CREATE TABLE movie_gross(
    name varchar(30),
    gross bigint
);

-- Start a transaction and insert data
BEGIN;
INSERT INTO movie_gross VALUES ( 'Raiders of the Lost Ark', 23400000);
INSERT INTO movie_gross VALUES ( 'Star Wars', 10000000 );

-- Roll back the transaction using ABORT
ABORT;

-- Example 2: ABORT with WORK keyword
BEGIN;
INSERT INTO movie_gross VALUES ( 'The Empire Strikes Back', 20000000);
ABORT WORK;

-- Example 3: ABORT with TRANSACTION keyword
BEGIN;
INSERT INTO movie_gross VALUES ( 'Return of the Jedi', 15000000);
ABORT TRANSACTION;

-- All variations of ABORT
ABORT;
ABORT WORK;
ABORT TRANSACTION;