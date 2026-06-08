-- START TRANSACTION test cases for Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_START_TRANSACTION.html
-- START TRANSACTION is a synonym of BEGIN

-- Basic START TRANSACTION
START TRANSACTION;

-- START TRANSACTION with isolation levels
-- Note: All isolation levels map to SERIALIZABLE in Redshift
START TRANSACTION ISOLATION LEVEL SERIALIZABLE;
START TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
START TRANSACTION ISOLATION LEVEL READ COMMITTED;
START TRANSACTION ISOLATION LEVEL REPEATABLE READ;

-- START TRANSACTION with read/write modes
START TRANSACTION READ WRITE;
START TRANSACTION READ ONLY;

-- START TRANSACTION with both isolation level and read/write mode
START TRANSACTION ISOLATION LEVEL SERIALIZABLE READ WRITE;
START TRANSACTION ISOLATION LEVEL SERIALIZABLE READ ONLY;
START TRANSACTION ISOLATION LEVEL READ COMMITTED READ WRITE;
START TRANSACTION ISOLATION LEVEL READ COMMITTED READ ONLY;
START TRANSACTION ISOLATION LEVEL READ UNCOMMITTED READ WRITE;
START TRANSACTION ISOLATION LEVEL READ UNCOMMITTED READ ONLY;
START TRANSACTION ISOLATION LEVEL REPEATABLE READ READ WRITE;
START TRANSACTION ISOLATION LEVEL REPEATABLE READ READ ONLY;

-- Order variations (read/write mode before isolation level)
START TRANSACTION READ WRITE ISOLATION LEVEL SERIALIZABLE;
START TRANSACTION READ ONLY ISOLATION LEVEL READ COMMITTED;

-- Case variations (SQL is case-insensitive)
start transaction;
Start Transaction;
START transaction;
start TRANSACTION;

-- With isolation level case variations
start transaction isolation level serializable;
START TRANSACTION ISOLATION level READ UNCOMMITTED;
Start Transaction Isolation Level Read Committed;

-- With read/write mode case variations
start transaction read write;
START TRANSACTION READ only;
Start Transaction Read Write;

-- Complex case variations
start transaction isolation level serializable read write;
START TRANSACTION ISOLATION LEVEL REPEATABLE READ READ ONLY;
Start Transaction Read Only Isolation Level Read Uncommitted;

-- Example usage in a transaction block
START TRANSACTION;
-- Insert some data
-- UPDATE operations
-- SELECT queries
COMMIT;

-- Example with read only mode
START TRANSACTION READ ONLY;
-- SELECT queries only
-- No INSERT/UPDATE/DELETE allowed
COMMIT;

-- Example with explicit isolation level
START TRANSACTION ISOLATION LEVEL SERIALIZABLE;
-- Operations with serializable isolation
-- This is the default in Redshift
END;

-- Example combining with other transaction control
START TRANSACTION READ WRITE;
-- Some operations
ROLLBACK;

-- Note: These are equivalent in Redshift:
-- START TRANSACTION;
-- BEGIN;
-- BEGIN TRANSACTION;
-- BEGIN WORK;