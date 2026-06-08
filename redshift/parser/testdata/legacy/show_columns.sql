-- SHOW COLUMNS command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_COLUMNS.html

-- Basic syntax with fully qualified table name
SHOW COLUMNS FROM TABLE sample_data_dev.tickit.event;

-- AWS Glue Data Catalog example with row limit
SHOW COLUMNS FROM TABLE awsdatacatalog.batman.nation LIMIT 2;

-- Basic table reference (single name)
SHOW COLUMNS FROM TABLE event;
SHOW COLUMNS FROM TABLE sales;
SHOW COLUMNS FROM TABLE venue;
SHOW COLUMNS FROM TABLE listing;
SHOW COLUMNS FROM TABLE category;
SHOW COLUMNS FROM TABLE date;
SHOW COLUMNS FROM TABLE users;

-- Database and table name (two-part name)
SHOW COLUMNS FROM TABLE tickit.event;
SHOW COLUMNS FROM TABLE tickit.sales;
SHOW COLUMNS FROM TABLE tickit.venue;
SHOW COLUMNS FROM TABLE public.listing;
SHOW COLUMNS FROM TABLE public.category;

-- Full three-part names (database.schema.table)
SHOW COLUMNS FROM TABLE mydb.public.sales;
SHOW COLUMNS FROM TABLE mydb.public.venue;
SHOW COLUMNS FROM TABLE sample_data_dev.public.users;
SHOW COLUMNS FROM TABLE sample_data_dev.tickit.event;
SHOW COLUMNS FROM TABLE sample_data_dev.tickit.listing;

-- With LIKE pattern matching
SHOW COLUMNS FROM TABLE event LIKE 'event%';
SHOW COLUMNS FROM TABLE sales LIKE '%id';
SHOW COLUMNS FROM TABLE venue LIKE 'venue_%';
SHOW COLUMNS FROM TABLE users LIKE 'user%';
SHOW COLUMNS FROM TABLE listing LIKE '%price%';
SHOW COLUMNS FROM TABLE category LIKE 'cat%';

-- With LIKE pattern on fully qualified tables
SHOW COLUMNS FROM TABLE tickit.event LIKE '%date%';
SHOW COLUMNS FROM TABLE tickit.sales LIKE 'sale%';
SHOW COLUMNS FROM TABLE sample_data_dev.tickit.venue LIKE '%name%';

-- With LIMIT clause
SHOW COLUMNS FROM TABLE event LIMIT 5;
SHOW COLUMNS FROM TABLE sales LIMIT 10;
SHOW COLUMNS FROM TABLE venue LIMIT 3;
SHOW COLUMNS FROM TABLE listing LIMIT 100;
SHOW COLUMNS FROM TABLE category LIMIT 1000;

-- With LIMIT on fully qualified tables
SHOW COLUMNS FROM TABLE tickit.event LIMIT 5;
SHOW COLUMNS FROM TABLE tickit.sales LIMIT 10;
SHOW COLUMNS FROM TABLE sample_data_dev.tickit.venue LIMIT 20;

-- Combining LIKE and LIMIT
SHOW COLUMNS FROM TABLE event LIKE 'event%' LIMIT 5;
SHOW COLUMNS FROM TABLE sales LIKE '%id' LIMIT 10;
SHOW COLUMNS FROM TABLE venue LIKE 'venue_%' LIMIT 3;
SHOW COLUMNS FROM TABLE tickit.listing LIKE '%price%' LIMIT 100;
SHOW COLUMNS FROM TABLE sample_data_dev.tickit.category LIKE 'cat%' LIMIT 50;

-- Edge cases and special patterns
SHOW COLUMNS FROM TABLE sales LIKE '_id';  -- Single character wildcard
SHOW COLUMNS FROM TABLE sales LIKE '%_id'; -- Ends with _id
SHOW COLUMNS FROM TABLE sales LIKE 'sale_';
SHOW COLUMNS FROM TABLE event LIKE '%time%';
SHOW COLUMNS FROM TABLE users LIKE 'user_%_%'; -- Multiple wildcards

-- Maximum limit value (10000)
SHOW COLUMNS FROM TABLE event LIMIT 10000;
SHOW COLUMNS FROM TABLE tickit.sales LIMIT 10000;

-- Zero limit (should return no rows)
SHOW COLUMNS FROM TABLE event LIMIT 0;

-- External tables (if supported)
SHOW COLUMNS FROM TABLE external_schema.external_table;
SHOW COLUMNS FROM TABLE awsdatacatalog.external_db.external_table;

-- Mixed case identifiers
SHOW COLUMNS FROM TABLE MyDatabase.MySchema.MyTable;
SHOW COLUMNS FROM TABLE "MixedCase"."Table";

-- Special characters in identifiers (quoted)
SHOW COLUMNS FROM TABLE "special-db"."special-schema"."special-table";
SHOW COLUMNS FROM TABLE "db.with.dots"."schema"."table";

-- Complex LIKE patterns
SHOW COLUMNS FROM TABLE sales LIKE 's%s';
SHOW COLUMNS FROM TABLE sales LIKE '%a%e%';
SHOW COLUMNS FROM TABLE venue LIKE 'v____e%'; -- Four underscores
SHOW COLUMNS FROM TABLE listing LIKE 'l%ing%';

-- AWS Glue Data Catalog patterns
SHOW COLUMNS FROM TABLE awsdatacatalog.mydb.mytable;
SHOW COLUMNS FROM TABLE awsdatacatalog.mydb.mytable LIKE 'col%';
SHOW COLUMNS FROM TABLE awsdatacatalog.mydb.mytable LIMIT 5;
SHOW COLUMNS FROM TABLE awsdatacatalog.mydb.mytable LIKE '%id' LIMIT 10;