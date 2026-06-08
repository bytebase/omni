-- Amazon Redshift INSERT statement test cases
-- Based on: https://docs.aws.amazon.com/redshift/latest/dg/r_INSERT.html

-- Basic INSERT with explicit column list
INSERT INTO event (eventid, eventname, starttime, venueid)
VALUES (5000, 'Salome', '2008-01-05 14:00', 3);

-- INSERT with all columns (implicit column list)
INSERT INTO event
VALUES (5000, 3, 5, 'Salome', '2008-01-05 14:00', 'Opera', 'Strauss');

-- Multiple rows INSERT
INSERT INTO event (eventid, eventname, starttime, venueid)
VALUES
    (5000, 'Salome', '2008-01-05 14:00', 3),
    (5001, 'La Traviata', '2008-01-05 20:00', 4),
    (5002, 'Tosca', '2008-01-06 14:00', 5);

-- INSERT with DEFAULT values
INSERT INTO users (userid, username, firstname, lastname, city, state, email, phone)
VALUES (1001, 'user1001', DEFAULT, DEFAULT, 'New York', 'NY', DEFAULT, DEFAULT);

-- INSERT with NULL values
INSERT INTO venue (venueid, venuename, venuecity, venuestate, venueseats)
VALUES (500, NULL, 'New York', 'NY', NULL);

-- INSERT with SELECT statement
INSERT INTO event_backup
SELECT * FROM event WHERE starttime > '2008-01-01';

-- INSERT with specific columns from SELECT
INSERT INTO category_stage (catid, catname)
SELECT catid, catname FROM category WHERE catgroup = 'Concerts';

-- INSERT with complex SELECT (joins, aggregates)
INSERT INTO sales_summary (salesid, total_price, commission)
SELECT
    salesid,
    pricepaid * qtysold as total_price,
    commission
FROM sales
WHERE saletime >= '2008-01-01'
ORDER BY salesid;

-- INSERT with SELECT and GROUP BY
INSERT INTO venue_events_count (venueid, event_count)
SELECT
    venueid,
    COUNT(*) as event_count
FROM event
GROUP BY venueid
HAVING COUNT(*) > 10;

-- INSERT with SELECT and JOIN
INSERT INTO user_sales (userid, username, total_sales)
SELECT
    u.userid,
    u.username,
    SUM(s.pricepaid) as total_sales
FROM users u
JOIN sales s ON u.userid = s.buyerid
GROUP BY u.userid, u.username;

-- INSERT with date literals
INSERT INTO date_table (dateid, caldate, day, week, month, qtr, year)
VALUES (20080105, '2008-01-05', 'SA', 1, 'JAN', 1, 2008);

-- INSERT with timestamp literals
INSERT INTO listing (listid, sellerid, eventid, dateid, numtickets, priceperticket, totalprice, listtime)
VALUES (50000, 1001, 5000, 20080105, 4, 25.00, 100.00, '2008-01-05 07:45:00');

-- INSERT with boolean literals
INSERT INTO customer_preferences (customerid, email_opt_in, sms_opt_in, active)
VALUES (1001, TRUE, FALSE, TRUE);

-- INSERT with numeric literals (integers and decimals)
INSERT INTO product (productid, name, price, quantity, weight)
VALUES (2001, 'Widget Pro', 29.99, 150, 2.5);

-- INSERT with string literals and special characters
INSERT INTO venue (venueid, venuename, venuecity, venuestate, venueseats)
VALUES (600, 'O''Brien Theater', 'Boston', 'MA', 2500);

-- INSERT with DEFAULT VALUES clause
INSERT INTO category DEFAULT VALUES;

-- INSERT with schema-qualified table name
INSERT INTO public.event (eventid, eventname, starttime, venueid)
VALUES (6000, 'Carmen', '2008-02-05 19:00', 8);

-- INSERT with quoted identifiers
INSERT INTO "Event" ("EventID", "EventName", "StartTime", "VenueID")
VALUES (7000, 'Phantom', '2008-03-05 20:00', 10);

-- INSERT with subquery in VALUES
INSERT INTO sales (salesid, listid, sellerid, buyerid, eventid, dateid, qtysold, pricepaid, commission, saletime)
VALUES (
    100000,
    (SELECT MAX(listid) FROM listing WHERE eventid = 5000),
    1001,
    2001,
    5000,
    20080105,
    2,
    50.00,
    7.50,
    '2008-01-05 09:30:00'
);

-- INSERT with arithmetic expressions in SELECT
INSERT INTO sales_metrics (metric_date, avg_price, total_revenue, ticket_count)
SELECT
    DATE(saletime) as metric_date,
    AVG(pricepaid) as avg_price,
    SUM(pricepaid * qtysold) as total_revenue,
    SUM(qtysold) as ticket_count
FROM sales
WHERE saletime >= '2008-01-01'
GROUP BY DATE(saletime);

-- INSERT with CASE expression in SELECT
INSERT INTO user_category (userid, username, user_type)
SELECT
    userid,
    username,
    CASE
        WHEN likemusicals = TRUE THEN 'Musical Fan'
        WHEN likesports = TRUE THEN 'Sports Fan'
        ELSE 'General'
    END as user_type
FROM users;

-- INSERT with window functions in SELECT
INSERT INTO sales_ranking (salesid, rank_by_price, total_sales)
SELECT
    salesid,
    RANK() OVER (ORDER BY pricepaid DESC) as rank_by_price,
    SUM(pricepaid) OVER () as total_sales
FROM sales
WHERE saletime >= '2008-01-01';

-- INSERT with LIMIT clause in SELECT
INSERT INTO top_venues (venueid, venuename, total_events)
SELECT
    v.venueid,
    v.venuename,
    COUNT(e.eventid) as total_events
FROM venue v
JOIN event e ON v.venueid = e.venueid
GROUP BY v.venueid, v.venuename
ORDER BY total_events DESC
LIMIT 10;

-- INSERT with CTE (Common Table Expression)
INSERT INTO monthly_sales (month, year, total_sales, event_count)
WITH monthly_data AS (
    SELECT
        EXTRACT(MONTH FROM saletime) as month,
        EXTRACT(YEAR FROM saletime) as year,
        SUM(pricepaid) as total_sales,
        COUNT(DISTINCT eventid) as event_count
    FROM sales
    GROUP BY EXTRACT(MONTH FROM saletime), EXTRACT(YEAR FROM saletime)
)
SELECT month, year, total_sales, event_count
FROM monthly_data
WHERE year = 2008;

-- INSERT with multiple CTEs
INSERT INTO venue_performance (venueid, venuename, total_sales, avg_price)
WITH venue_sales AS (
    SELECT
        v.venueid,
        v.venuename,
        SUM(s.pricepaid) as total_sales,
        AVG(s.pricepaid) as avg_price
    FROM venue v
    JOIN event e ON v.venueid = e.venueid
    JOIN sales s ON e.eventid = s.eventid
    GROUP BY v.venueid, v.venuename
),
top_venues AS (
    SELECT *
    FROM venue_sales
    WHERE total_sales > 10000
)
SELECT venueid, venuename, total_sales, avg_price
FROM top_venues;

-- INSERT with UNION in SELECT
INSERT INTO all_names (id, name, type)
SELECT userid as id, username as name, 'user' as type FROM users
UNION ALL
SELECT venueid as id, venuename as name, 'venue' as type FROM venue;

-- INSERT with very long value lists (testing parser performance)
INSERT INTO large_table (id, value)
VALUES
    (1, 'value1'),
    (2, 'value2'),
    (3, 'value3'),
    (4, 'value4'),
    (5, 'value5'),
    (6, 'value6'),
    (7, 'value7'),
    (8, 'value8'),
    (9, 'value9'),
    (10, 'value10');

-- INSERT with mixed data types
INSERT INTO mixed_types (
    int_col,
    bigint_col,
    decimal_col,
    real_col,
    double_col,
    char_col,
    varchar_col,
    date_col,
    timestamp_col,
    boolean_col
) VALUES (
    42,
    9223372036854775807,
    123.45,
    3.14159,
    2.71828182845905,
    'A',
    'Hello, World!',
    '2008-01-05',
    '2008-01-05 14:30:00',
    TRUE
);

-- INSERT with escape sequences in strings
INSERT INTO text_data (id, content)
VALUES
    (1, 'Line 1\nLine 2'),
    (2, 'Tab\tseparated'),
    (3, 'Quote: ''escaped'''),
    (4, 'Backslash: \\');

-- INSERT with unicode characters
INSERT INTO international (id, city, greeting)
VALUES
    (1, 'Tokyo', '東京'),
    (2, 'Paris', 'Bonjour'),
    (3, 'Moscow', 'Москва');

-- INSERT with computed columns in SELECT
INSERT INTO calculated_values (id, base_value, computed_value)
SELECT
    id,
    value,
    value * 1.1 + 10 as computed_value
FROM source_table
WHERE value > 0;

-- INSERT with correlated subquery
INSERT INTO user_first_sale (userid, username, first_sale_date)
SELECT
    u.userid,
    u.username,
    (SELECT MIN(saletime) FROM sales s WHERE s.buyerid = u.userid) as first_sale_date
FROM users u
WHERE EXISTS (SELECT 1 FROM sales s WHERE s.buyerid = u.userid);

-- INSERT with SELECT DISTINCT
INSERT INTO unique_cities (city, state)
SELECT DISTINCT venuecity, venuestate
FROM venue
WHERE venuestate IS NOT NULL
ORDER BY venuecity, venuestate;

-- INSERT with outer join
INSERT INTO user_sales_summary (userid, username, total_spent, purchase_count)
SELECT
    u.userid,
    u.username,
    COALESCE(SUM(s.pricepaid), 0) as total_spent,
    COUNT(s.salesid) as purchase_count
FROM users u
LEFT JOIN sales s ON u.userid = s.buyerid
GROUP BY u.userid, u.username;

-- INSERT with full outer join
INSERT INTO venue_event_matrix (venueid, eventid, has_event)
SELECT
    COALESCE(v.venueid, e.venueid) as venueid,
    e.eventid,
    CASE WHEN e.eventid IS NOT NULL THEN TRUE ELSE FALSE END as has_event
FROM venue v
FULL OUTER JOIN event e ON v.venueid = e.venueid;