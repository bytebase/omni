-- DECLARE command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DECLARE.html
-- Reference: https://www.postgresql.org/docs/current/sql-declare.html

-- Basic DECLARE statements
DECLARE my_cursor CURSOR FOR SELECT * FROM sales;
declare my_cursor cursor for select * from sales;

-- DECLARE with quoted cursor names
DECLARE "my_cursor" CURSOR FOR SELECT * FROM sales;
DECLARE "My Cursor" CURSOR FOR SELECT * FROM sales;

-- DECLARE with BINARY option
DECLARE my_cursor BINARY CURSOR FOR SELECT * FROM sales;
DECLARE my_cursor BINARY CURSOR FOR SELECT qty, pricepaid FROM sales;

DECLARE my_cursor INSENSITIVE CURSOR FOR SELECT * FROM sales;

-- DECLARE with SCROLL options
DECLARE my_cursor SCROLL CURSOR FOR SELECT * FROM sales;
DECLARE my_cursor NO SCROLL CURSOR FOR SELECT * FROM sales;

-- DECLARE with HOLD options
DECLARE my_cursor CURSOR WITH HOLD FOR SELECT * FROM sales;
DECLARE my_cursor CURSOR WITHOUT HOLD FOR SELECT * FROM sales;

-- DECLARE with multiple options
DECLARE my_cursor BINARY SCROLL CURSOR FOR SELECT * FROM sales;
DECLARE my_cursor BINARY INSENSITIVE CURSOR FOR SELECT * FROM sales;
DECLARE my_cursor BINARY SCROLL CURSOR WITH HOLD FOR SELECT * FROM sales;
DECLARE my_cursor BINARY INSENSITIVE CURSOR WITHOUT HOLD FOR SELECT * FROM sales;
DECLARE my_cursor SCROLL CURSOR WITH HOLD FOR SELECT * FROM sales;
DECLARE my_cursor NO SCROLL CURSOR WITH HOLD FOR SELECT * FROM sales;
DECLARE my_cursor INSENSITIVE NO SCROLL CURSOR FOR SELECT * FROM sales;

-- DECLARE with complex SELECT statements
DECLARE sales_cursor CURSOR FOR
    SELECT s.salesid, s.listid, s.sellerid, s.buyerid, s.eventid, s.dateid, s.qtysold, s.pricepaid, s.commission
    FROM sales s
    WHERE s.pricepaid > 100
    ORDER BY s.pricepaid DESC;

-- DECLARE with JOIN queries
DECLARE event_sales_cursor CURSOR FOR
    SELECT e.eventname, s.qtysold, s.pricepaid
    FROM event e
    JOIN sales s ON e.eventid = s.eventid
    WHERE e.catid = 7
    ORDER BY s.pricepaid DESC;

-- DECLARE with aggregate functions
DECLARE sales_summary_cursor CURSOR FOR
    SELECT eventid, COUNT(*) as ticket_count, SUM(pricepaid) as total_sales
    FROM sales
    GROUP BY eventid
    HAVING SUM(pricepaid) > 1000
    ORDER BY total_sales DESC;

-- DECLARE with subqueries
DECLARE high_value_sales_cursor CURSOR FOR
    SELECT * FROM sales
    WHERE pricepaid > (SELECT AVG(pricepaid) FROM sales)
    ORDER BY pricepaid DESC;

-- DECLARE with CTE (Common Table Expression)
DECLARE cte_cursor CURSOR FOR
    WITH sales_summary AS (
        SELECT eventid, COUNT(*) as ticket_count, SUM(pricepaid) as total_sales
        FROM sales
        GROUP BY eventid
    )
    SELECT ss.eventid, ss.ticket_count, ss.total_sales, e.eventname
    FROM sales_summary ss
    JOIN event e ON ss.eventid = e.eventid
    WHERE ss.total_sales > 1000
    ORDER BY ss.total_sales DESC;

-- DECLARE with UNION
DECLARE union_cursor CURSOR FOR
    SELECT 'High Value' as category, salesid, pricepaid FROM sales WHERE pricepaid > 500
    UNION ALL
    SELECT 'Low Value' as category, salesid, pricepaid FROM sales WHERE pricepaid <= 500
    ORDER BY pricepaid DESC;

-- DECLARE with VALUES statement
DECLARE values_cursor CURSOR FOR VALUES (1, 'First'), (2, 'Second'), (3, 'Third');

-- DECLARE with TABLE statement
DECLARE table_cursor CURSOR FOR TABLE sales;

-- DECLARE with window functions
DECLARE window_cursor CURSOR FOR
    SELECT salesid, pricepaid,
           ROW_NUMBER() OVER (PARTITION BY eventid ORDER BY pricepaid DESC) as rank_in_event
    FROM sales
    WHERE qtysold > 1;

-- DECLARE with CASE statements
DECLARE case_cursor CURSOR FOR
    SELECT salesid,
           CASE
               WHEN pricepaid > 1000 THEN 'Premium'
               WHEN pricepaid > 500 THEN 'Standard'
               ELSE 'Economy'
           END as price_category
    FROM sales;

-- DECLARE with EXISTS subquery
DECLARE exists_cursor CURSOR FOR
    SELECT e.eventname
    FROM event e
    WHERE EXISTS (
        SELECT 1 FROM sales s
        WHERE s.eventid = e.eventid AND s.pricepaid > 1000
    );

-- DECLARE with IN clause
DECLARE in_cursor CURSOR FOR
    SELECT * FROM sales
    WHERE eventid IN (SELECT eventid FROM event WHERE catid IN (6, 7, 8))
    ORDER BY pricepaid DESC;

-- DECLARE with BETWEEN
DECLARE between_cursor CURSOR FOR
    SELECT * FROM sales
    WHERE pricepaid BETWEEN 100 AND 500
    AND dateid BETWEEN 1900 AND 2000;

-- DECLARE with LIKE pattern matching
DECLARE like_cursor CURSOR FOR
    SELECT e.eventname, s.pricepaid
    FROM event e
    JOIN sales s ON e.eventid = s.eventid
    WHERE e.eventname LIKE '%Concert%'
    ORDER BY s.pricepaid DESC;

-- DECLARE with NULL handling
DECLARE null_cursor CURSOR FOR
    SELECT salesid, pricepaid, commission
    FROM sales
    WHERE commission IS NOT NULL
    ORDER BY commission DESC;

-- DECLARE with LIMIT and OFFSET
DECLARE limited_cursor CURSOR FOR
    SELECT * FROM sales
    ORDER BY pricepaid DESC
    LIMIT 100 OFFSET 50;

-- DECLARE with mathematical operations
DECLARE math_cursor CURSOR FOR
    SELECT salesid,
           pricepaid,
           commission,
           (pricepaid - commission) as net_amount,
           (commission / pricepaid * 100) as commission_percentage
    FROM sales
    WHERE pricepaid > 0 AND commission > 0;

-- DECLARE with string concatenation
DECLARE concat_cursor CURSOR FOR
    SELECT u.firstname || ' ' || u.lastname as full_name,
           s.pricepaid
    FROM users u
    JOIN sales s ON u.userid = s.buyerid
    WHERE s.pricepaid > 100;

-- Complex example with all options
DECLARE complex_cursor BINARY INSENSITIVE SCROLL CURSOR WITH HOLD FOR
    WITH high_value_events AS (
        SELECT eventid, AVG(pricepaid) as avg_price
        FROM sales
        GROUP BY eventid
        HAVING AVG(pricepaid) > 200
    ),
    event_details AS (
        SELECT e.eventid, e.eventname, e.catid, hve.avg_price
        FROM event e
        JOIN high_value_events hve ON e.eventid = hve.eventid
    )
    SELECT ed.eventname,
           ed.avg_price,
           c.catname,
           COUNT(s.salesid) as total_sales,
           SUM(s.pricepaid) as total_revenue,
           CASE
               WHEN ed.avg_price > 500 THEN 'Premium'
               WHEN ed.avg_price > 200 THEN 'Standard'
               ELSE 'Economy'
           END as event_tier
    FROM event_details ed
    JOIN category c ON ed.catid = c.catid
    JOIN sales s ON ed.eventid = s.eventid
    WHERE s.pricepaid IS NOT NULL
    GROUP BY ed.eventname, ed.avg_price, c.catname
    HAVING COUNT(s.salesid) > 5
    ORDER BY total_revenue DESC, ed.avg_price DESC
    LIMIT 50;

-- Examples from AWS Redshift documentation (if any specific examples exist)
-- Basic cursor for sales analysis
DECLARE sales_analysis_cursor CURSOR FOR
    SELECT eventid, SUM(qtysold) as total_quantity, SUM(pricepaid) as total_revenue
    FROM sales
    GROUP BY eventid
    ORDER BY total_revenue DESC;

-- Cursor for user analysis
DECLARE user_cursor CURSOR FOR
    SELECT userid, firstname, lastname, city, state
    FROM users
    WHERE state IN ('CA', 'NY', 'FL')
    ORDER BY lastname, firstname;

-- Mixed case variations (testing case insensitivity)
Declare My_Cursor Cursor For Select * From Sales;
DECLARE my_cursor CURSOR FOR select * from sales;
declare MY_CURSOR cursor for SELECT * FROM SALES;

-- Edge cases and boundary conditions
DECLARE edge_cursor CURSOR FOR SELECT 1;
DECLARE edge_cursor CURSOR FOR SELECT NULL;
DECLARE edge_cursor CURSOR FOR SELECT 'test' as col1, 123 as col2;

-- Cursor with parameter-like syntax (though parameters might not be directly supported in DECLARE)
DECLARE param_cursor CURSOR FOR
    SELECT * FROM sales WHERE pricepaid > 100;