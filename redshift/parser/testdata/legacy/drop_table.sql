-- DROP TABLE Test Cases for Amazon Redshift
-- Based on AWS documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_TABLE.html

-- Basic DROP TABLE
DROP TABLE feedback;

-- DROP TABLE with schema qualification
DROP TABLE public.feedback;
DROP TABLE schema1.table1;

-- DROP TABLE with IF EXISTS
DROP TABLE IF EXISTS feedback;
DROP TABLE IF EXISTS public.feedback;
DROP TABLE IF EXISTS schema1.table1;

-- Multiple tables in single statement
DROP TABLE feedback, buyers;
DROP TABLE table1, table2, table3;
DROP TABLE public.feedback, public.buyers;
DROP TABLE schema1.table1, schema2.table2;

-- DROP TABLE with CASCADE (drops dependent objects)
DROP TABLE feedback CASCADE;
DROP TABLE public.feedback CASCADE;

-- DROP TABLE with RESTRICT (prevents drop if dependencies exist - default behavior)
DROP TABLE feedback RESTRICT;
DROP TABLE public.feedback RESTRICT;

-- Combination of IF EXISTS with CASCADE/RESTRICT
DROP TABLE IF EXISTS feedback CASCADE;
DROP TABLE IF EXISTS feedback RESTRICT;
DROP TABLE IF EXISTS public.feedback CASCADE;
DROP TABLE IF EXISTS public.feedback RESTRICT;

-- Multiple tables with IF EXISTS
DROP TABLE IF EXISTS feedback, buyers;
DROP TABLE IF EXISTS table1, table2, table3;
DROP TABLE IF EXISTS public.feedback, public.buyers;

-- Multiple tables with CASCADE/RESTRICT
DROP TABLE feedback, buyers CASCADE;
DROP TABLE feedback, buyers RESTRICT;
DROP TABLE table1, table2, table3 CASCADE;
DROP TABLE table1, table2, table3 RESTRICT;

-- Complex combinations
DROP TABLE IF EXISTS feedback, buyers CASCADE;
DROP TABLE IF EXISTS feedback, buyers RESTRICT;
DROP TABLE IF EXISTS public.feedback, public.buyers CASCADE;
DROP TABLE IF EXISTS schema1.table1, schema2.table2 RESTRICT;

-- Tables with various identifier patterns
DROP TABLE "quoted_table";
DROP TABLE public."quoted_table";
DROP TABLE IF EXISTS "quoted_table";
DROP TABLE "table with spaces";
DROP TABLE public."table with spaces";

-- Mixed quoted and unquoted identifiers
DROP TABLE feedback, "quoted_table";
DROP TABLE public.feedback, public."quoted_table";
DROP TABLE IF EXISTS feedback, "quoted_table" CASCADE;

-- Real-world examples based on AWS documentation scenarios
-- Example 1: Basic table creation and drop
-- CREATE TABLE feedback(a int);
DROP TABLE feedback;

-- Example 2: Multiple table drop
-- CREATE TABLE feedback(a int);
-- CREATE TABLE buyers(a int);
DROP TABLE feedback, buyers;

-- Example 3: Drop with dependencies (CASCADE scenario)
-- CREATE TABLE feedback(a int);
-- CREATE VIEW feedback_view AS SELECT * FROM feedback;
DROP TABLE feedback CASCADE;

-- Example 4: Conditional drop to avoid errors
DROP TABLE IF EXISTS feedback;
DROP TABLE IF EXISTS nonexistent_table;

-- Schema-qualified examples
DROP TABLE public.sales_data;
DROP TABLE analytics.user_metrics;
DROP TABLE staging.temp_data;

-- Complex schema scenarios
DROP TABLE IF EXISTS public.sales_data, analytics.user_metrics CASCADE;
DROP TABLE public.old_table, public.legacy_table, public.archive_table RESTRICT;

-- Edge cases with identifier names
DROP TABLE table123;
DROP TABLE _private_table;
DROP TABLE table_with_underscores;
DROP TABLE TABLE_UPPERCASE;
DROP TABLE "123_table";
DROP TABLE "table-with-dashes";

-- Multiple schemas in single statement
DROP TABLE schema1.table1, schema2.table2, schema3.table3;
DROP TABLE IF EXISTS dev.temp_table, staging.temp_table, prod.temp_table CASCADE;