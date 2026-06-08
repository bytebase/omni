-- ANALYZE statement test cases
-- https://docs.aws.amazon.com/redshift/latest/dg/r_ANALYZE.html

-- Basic ANALYZE statements
ANALYZE;  -- Analyze all tables in the current database
ANALYSE;  -- Alternative spelling

-- ANALYZE with VERBOSE option
ANALYZE VERBOSE;
ANALYSE VERBOSE;

-- ANALYZE specific table
ANALYZE listing;
ANALYZE sales;
ANALYZE venue;
ANALYZE event;

-- ANALYZE specific table with VERBOSE
ANALYZE VERBOSE listing;
ANALYZE VERBOSE sales;

-- ANALYZE specific columns
ANALYZE listing(listid);
ANALYZE listing(listid, sellerid, eventid);
ANALYZE venue(venueid, venuename);
ANALYZE sales(salesid, listid, sellerid, buyerid);

-- ANALYZE with PREDICATE COLUMNS option
ANALYZE venue PREDICATE COLUMNS;
ANALYZE listing PREDICATE COLUMNS;
ANALYZE sales PREDICATE COLUMNS;

-- ANALYZE with ALL COLUMNS option (default behavior)
ANALYZE venue ALL COLUMNS;
ANALYZE listing ALL COLUMNS;
ANALYZE sales ALL COLUMNS;

-- ANALYZE with VERBOSE and column options
ANALYZE VERBOSE venue PREDICATE COLUMNS;
ANALYZE VERBOSE listing ALL COLUMNS;

-- Schema-qualified table names
ANALYZE public.sales;
ANALYZE public.venue(venueid, venuename);
ANALYZE public.listing PREDICATE COLUMNS;

-- Multiple tables cannot be specified (these would be separate statements)
ANALYZE listing;
ANALYZE sales;
ANALYZE venue;

-- Real-world examples from documentation
ANALYZE VERBOSE listing(listid, sellerid, eventid, dateid);
ANALYZE venue(venueid, venuename, venuecity, venuestate);
ANALYZE sales(salesid, listid, sellerid, buyerid, eventid, dateid, qtysold, pricepaid, commission);