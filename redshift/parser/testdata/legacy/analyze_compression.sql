-- ANALYZE COMPRESSION statement test cases
-- https://docs.aws.amazon.com/redshift/latest/dg/r_ANALYZE_COMPRESSION.html

-- Basic ANALYZE COMPRESSION statements
ANALYZE COMPRESSION;  -- Analyze all tables in current database
ANALYSE COMPRESSION;  -- Alternative spelling

-- ANALYZE COMPRESSION for specific table
ANALYZE COMPRESSION listing;
ANALYZE COMPRESSION sales;
ANALYZE COMPRESSION venue;
ANALYZE COMPRESSION event;

-- ANALYZE COMPRESSION for specific columns
ANALYZE COMPRESSION listing(listid);
ANALYZE COMPRESSION listing(listid, sellerid, eventid, dateid);
ANALYZE COMPRESSION sales(qtysold, commission, saletime);
ANALYZE COMPRESSION venue(venuename, venuecity, venuestate);

-- ANALYZE COMPRESSION with COMPROWS option
ANALYZE COMPRESSION listing COMPROWS 1000;        -- Minimum value
ANALYZE COMPRESSION listing COMPROWS 100000;      -- Default value
ANALYZE COMPRESSION listing COMPROWS 500000;      -- Custom sample size
ANALYZE COMPRESSION sales COMPROWS 1000000;       -- 1 million rows
ANALYZE COMPRESSION venue COMPROWS 1000000000;    -- Maximum value (1 billion)

-- ANALYZE COMPRESSION with columns and COMPROWS
ANALYZE COMPRESSION listing(listid, sellerid) COMPROWS 250000;
ANALYZE COMPRESSION sales(qtysold, commission, saletime) COMPROWS 750000;

-- Schema-qualified table names
ANALYZE COMPRESSION public.listing;
ANALYZE COMPRESSION public.sales COMPROWS 200000;
ANALYZE COMPRESSION public.venue(venuename, venuecity) COMPROWS 50000;

-- Real-world examples from documentation
ANALYZE COMPRESSION listing;
ANALYZE COMPRESSION sales(qtysold, commission, saletime);
ANALYZE COMPRESSION listing COMPROWS 500000;

-- Edge cases for COMPROWS
ANALYZE COMPRESSION listing COMPROWS 1000;        -- Minimum allowed
ANALYZE COMPRESSION listing COMPROWS 5000;        -- Below default (will be upgraded)
ANALYZE COMPRESSION listing COMPROWS 99999;       -- Just below default
ANALYZE COMPRESSION listing COMPROWS 100001;      -- Just above default
ANALYZE COMPRESSION listing COMPROWS 999999999;   -- Just below maximum
ANALYZE COMPRESSION listing COMPROWS 1000000000;  -- Maximum allowed

-- Complex examples combining features
ANALYZE COMPRESSION customer(c_custkey, c_name, c_address, c_city, c_nation, c_region) COMPROWS 250000;
ANALYZE COMPRESSION lineitem(l_orderkey, l_partkey, l_suppkey, l_quantity, l_extendedprice) COMPROWS 1000000;

-- Note: SORTKEY columns are automatically skipped during compression analysis
-- Note: Cannot specify multiple tables in one ANALYZE COMPRESSION statement