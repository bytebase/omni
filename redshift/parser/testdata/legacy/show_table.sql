-- SHOW TABLE Examples
-- https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_TABLE.html

-- Example 1: Show table definition for sales table
SHOW TABLE sales;

-- Example 2: Show table definition with schema qualification
SHOW TABLE public.category;

-- Example 3: Show table after creating with primary key
-- First create a table
CREATE TABLE foo(a INT PRIMARY KEY, b INT);
-- Then show its definition
SHOW TABLE foo;

-- Example 4: Show table in specific schema
SHOW TABLE myschema.mytable;

-- Example 5: Simple table name without schema
SHOW TABLE customer;

-- Example 6: Table with mixed case name (requires quotes in actual usage)
SHOW TABLE "MixedCaseTable";

-- Note: SHOW TABLE displays the complete CREATE TABLE statement including:
-- - Column definitions with data types
-- - Column constraints (NOT NULL, DEFAULT, etc.)
-- - Table constraints (PRIMARY KEY, FOREIGN KEY, etc.)
-- - Encoding specifications
-- - Distribution style (DISTKEY)
-- - Sort keys (SORTKEY)
-- - Table properties