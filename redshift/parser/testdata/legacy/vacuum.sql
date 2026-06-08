-- VACUUM command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_VACUUM_command.html

-- Basic VACUUM commands
VACUUM;
vacuum;
VACUUM sales;
vacuum "sales_table";
vacuum public.sales;

-- VACUUM FULL (default option)
VACUUM FULL;
VACUUM FULL sales;
vacuum full customer_table;
VACUUM FULL sales TO 75 PERCENT;
vacuum full sales to 90 percent;
VACUUM FULL sales BOOST;
vacuum full sales to 80 percent boost;

-- VACUUM SORT ONLY
VACUUM SORT ONLY;
vacuum sort only;
VACUUM SORT ONLY sales;
vacuum sort only customer_table;
VACUUM SORT ONLY sales TO 75 PERCENT;
vacuum sort only sales to 50 percent;
VACUUM SORT ONLY sales BOOST;
vacuum sort only sales to 95 percent boost;

-- VACUUM DELETE ONLY
VACUUM DELETE ONLY;
vacuum delete only;
VACUUM DELETE ONLY sales;
vacuum delete only inventory_table;
VACUUM DELETE ONLY sales TO 75 PERCENT;
vacuum delete only sales to 60 percent;
VACUUM DELETE ONLY sales BOOST;
vacuum delete only sales to 85 percent boost;

-- VACUUM REINDEX (requires table name)
VACUUM REINDEX listing;
vacuum reindex customer_table;
VACUUM REINDEX listing TO 80 PERCENT;
vacuum reindex listing to 90 percent;
VACUUM REINDEX listing BOOST;
vacuum reindex listing to 75 percent boost;

-- VACUUM RECLUSTER (requires table name)
VACUUM RECLUSTER listing;
vacuum recluster sales_fact;
VACUUM RECLUSTER listing TO 80 PERCENT;
vacuum recluster listing to 95 percent;
VACUUM RECLUSTER listing BOOST;
vacuum recluster listing to 70 percent boost;

-- Edge cases and variations
VACUUM TO 100 PERCENT;
vacuum to 0 percent;
VACUUM BOOST;
vacuum boost;

-- Mixed case and quoted identifiers
VACUUM "MixedCase_Table";
vacuum "table with spaces";
VACUUM FULL "schema"."table";
vacuum sort only "MyTable" to 75 percent boost;

-- Real-world examples from AWS documentation
-- Example 1: Run VACUUM on the entire database
VACUUM;

-- Example 2: Run VACUUM on specific table with default threshold
VACUUM sales;

-- Example 3: Run VACUUM with custom threshold
VACUUM sales TO 75 PERCENT;

-- Example 4: Sort only operation
VACUUM SORT ONLY sales TO 75 PERCENT;

-- Example 5: Delete only operation
VACUUM DELETE ONLY sales TO 75 PERCENT;

-- Example 6: Reindex interleaved table
VACUUM REINDEX listing;

-- Example 7: Recluster table with boost
VACUUM RECLUSTER listing BOOST;

-- Comments and formatting tests
VACUUM; -- Basic vacuum
VACUUM /* inline comment */ sales;
VACUUM
    FULL
    sales
    TO
    80
    PERCENT;

-- Multiple VACUUM commands
VACUUM sales;
VACUUM customer_table;
VACUUM inventory TO 90 PERCENT BOOST;