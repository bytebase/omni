-- SELECT INTO Examples for Amazon Redshift
-- Based on AWS documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_SELECT_INTO.html

-- Basic SELECT INTO - Create a new table from an existing table
SELECT * INTO newevent FROM event;

-- SELECT INTO with specific columns
SELECT eventid, eventname, venueid
INTO event_backup
FROM event;

-- SELECT INTO with WHERE clause
SELECT *
INTO recent_events
FROM event
WHERE starttime > '2023-01-01';

-- SELECT INTO with temporary table
SELECT *
INTO TEMP TABLE temp_events
FROM event;

-- SELECT INTO with TEMPORARY keyword (alternative syntax)
SELECT *
INTO TEMPORARY TABLE temp_events2
FROM event;

-- SELECT INTO with TABLE keyword (optional)
SELECT *
INTO TABLE event_copy
FROM event;

-- SELECT INTO TEMP TABLE with TABLE keyword
SELECT *
INTO TEMP TABLE temp_events3
FROM event;

-- SELECT INTO with aggregation
SELECT username, lastname, SUM(pricepaid - commission) AS profit
INTO TEMP TABLE profits
FROM sales, users
WHERE sales.sellerid = users.userid
GROUP BY 1, 2
ORDER BY 3 DESC;

-- SELECT INTO with JOIN
SELECT e.eventid, e.eventname, v.venuename, v.venuecity
INTO venue_events
FROM event e
         JOIN venue v ON e.venueid = v.venueid;

-- SELECT INTO with multiple JOINs
SELECT u.username, e.eventname, s.pricepaid
INTO user_event_sales
FROM users u
         JOIN sales s ON u.userid = s.buyerid
         JOIN event e ON s.eventid = e.eventid;

-- SELECT INTO with LEFT JOIN
SELECT c.catname, COUNT(e.eventid) as event_count
INTO category_counts
FROM category c
         LEFT JOIN event e ON c.catid = e.catid
GROUP BY c.catname;

-- SELECT INTO with subquery in FROM clause
SELECT *
INTO high_price_sales
FROM (
         SELECT * FROM sales WHERE pricepaid > 100
     ) AS expensive_sales;

-- SELECT INTO with CTE (Common Table Expression)
WITH high_commission_sales AS (
    SELECT * FROM sales WHERE commission > 50
)
SELECT *
INTO temp_high_commission
FROM high_commission_sales;

-- SELECT INTO with multiple CTEs
WITH
    venue_stats AS (
        SELECT venueid, COUNT(*) as event_count
        FROM event
        GROUP BY venueid
    ),
    popular_venues AS (
        SELECT * FROM venue_stats WHERE event_count > 10
    )
SELECT v.venuename, vs.event_count
INTO popular_venue_report
FROM popular_venues vs
         JOIN venue v ON vs.venueid = v.venueid;

-- SELECT INTO with UNION
SELECT eventid, 'event' as source
INTO combined_ids
FROM event
UNION
SELECT listid, 'listing' as source
FROM listing;

-- SELECT INTO with UNION ALL
SELECT userid, username, 'buyer' as role
INTO temp_all_users
FROM users
WHERE userid IN (SELECT buyerid FROM sales)
UNION ALL
SELECT userid, username, 'seller' as role
FROM users
WHERE userid IN (SELECT sellerid FROM sales);

-- SELECT INTO with DISTINCT
SELECT DISTINCT venuecity, venuestate
INTO unique_locations
FROM venue;

-- SELECT INTO with TOP clause
SELECT TOP 100 *
INTO top_sales
FROM sales
ORDER BY pricepaid DESC;

-- SELECT INTO with window functions
SELECT
    salesid,
    pricepaid,
    ROW_NUMBER() OVER (ORDER BY pricepaid DESC) as price_rank,
        RANK() OVER (PARTITION BY eventid ORDER BY pricepaid DESC) as event_price_rank
INTO sales_rankings
FROM sales;

-- SELECT INTO with CASE expressions
SELECT
    eventid,
    eventname,
    CASE
        WHEN starttime < '2023-01-01' THEN 'past'
        WHEN starttime >= '2023-01-01' AND starttime < '2024-01-01' THEN 'current'
        ELSE 'future'
        END as event_period
INTO event_periods
FROM event;

-- SELECT INTO with calculated columns
SELECT
    salesid,
    pricepaid,
    commission,
    pricepaid - commission as net_price,
    (commission / NULLIF(pricepaid, 0)) * 100 as commission_percentage
INTO sales_analysis
FROM sales;

-- SELECT INTO with GROUP BY and HAVING
SELECT
    eventid,
    COUNT(*) as sale_count,
    SUM(pricepaid) as total_revenue
INTO event_revenue
FROM sales
GROUP BY eventid
HAVING COUNT(*) > 5;

-- SELECT INTO with LIMIT and OFFSET
SELECT *
INTO recent_subset
FROM event
ORDER BY starttime DESC
    LIMIT 1000
OFFSET 100;

-- SELECT INTO with schema qualification
SELECT *
INTO myschema.event_backup
FROM public.event;

-- SELECT INTO with complex expressions
SELECT
    DATE_TRUNC('month', saletime) as sale_month,
    COUNT(*) as sale_count,
    AVG(pricepaid) as avg_price,
    MIN(pricepaid) as min_price,
    MAX(pricepaid) as max_price
INTO monthly_sales_stats
FROM sales
GROUP BY 1
ORDER BY 1;

-- SELECT INTO with EXCLUDE (if supported)
SELECT * EXCLUDE (commission, saletime)
INTO sales_simplified
FROM sales;

-- SELECT INTO with LISTAGG
SELECT
    venueid,
    LISTAGG(eventname, ', ') WITHIN GROUP (ORDER BY eventname) as events_list
INTO venue_event_list
FROM event
GROUP BY venueid;

-- SELECT INTO with GROUPING SETS
SELECT
    venuecity,
    venuestate,
    COUNT(*) as venue_count
INTO venue_grouping_stats
FROM venue
GROUP BY GROUPING SETS ((venuecity), (venuestate), (venuecity, venuestate));

-- SELECT INTO with ROLLUP
SELECT
    DATE_PART('year', saletime) as sale_year,
    DATE_PART('month', saletime) as sale_month,
    SUM(pricepaid) as total_sales
INTO sales_rollup
FROM sales
GROUP BY ROLLUP(1, 2);

-- SELECT INTO with CUBE
SELECT
    venuecity,
    venuestate,
    COUNT(*) as venue_count
INTO venue_cube_stats
FROM venue
GROUP BY CUBE(venuecity, venuestate);

-- SELECT INTO with EXISTS subquery
SELECT *
INTO venues_with_events
FROM venue v
WHERE EXISTS (
    SELECT 1 FROM event e WHERE e.venueid = v.venueid
);

-- SELECT INTO with NOT IN subquery
SELECT *
INTO users_without_sales
FROM users
WHERE userid NOT IN (
    SELECT DISTINCT buyerid FROM sales
    UNION
    SELECT DISTINCT sellerid FROM sales
);

-- SELECT INTO with correlated subquery
SELECT
    u.userid,
    u.username,
    (SELECT COUNT(*) FROM sales s WHERE s.buyerid = u.userid) as purchase_count
INTO user_purchase_counts
FROM users u;

-- SELECT INTO with INTERSECT
SELECT venueid
INTO venues_with_both
FROM event
WHERE catid = 1
INTERSECT
SELECT venueid
FROM event
WHERE catid = 2;

-- SELECT INTO with EXCEPT
SELECT userid
INTO buyers_not_sellers
FROM users
WHERE userid IN (SELECT DISTINCT buyerid FROM sales)
EXCEPT
SELECT userid
FROM users
WHERE userid IN (SELECT DISTINCT sellerid FROM sales);

-- SELECT INTO with complex date operations
SELECT
    eventid,
    eventname,
    starttime,
    DATEADD(hour, 3, starttime) as endtime_estimate,
    DATEDIFF(day, CURRENT_DATE, starttime) as days_until_event
INTO event_schedule
FROM event
WHERE starttime > CURRENT_DATE;

-- SELECT INTO with string operations
SELECT
    userid,
    username,
    UPPER(firstname) as firstname_upper,
    LOWER(lastname) as lastname_lower,
    firstname || ' ' || lastname as full_name,
    LENGTH(email) as email_length
INTO user_details
FROM users;

-- SELECT INTO with QUALIFY clause (if supported in Redshift)
SELECT
    salesid,
    eventid,
    pricepaid,
    ROW_NUMBER() OVER (PARTITION BY eventid ORDER BY pricepaid DESC) as rn
INTO top_sales_per_event
FROM sales
         QUALIFY rn <= 3;

-- SELECT INTO with recursive CTE (if supported)
WITH RECURSIVE event_hierarchy AS (
    SELECT eventid, eventname, 1 as level
    FROM event
    WHERE eventid = 1
    UNION ALL
    SELECT e.eventid, e.eventname, eh.level + 1
    FROM event e
             JOIN event_hierarchy eh ON e.eventid = eh.eventid + 1
)
SELECT *
INTO event_tree
FROM event_hierarchy;

-- Complex real-world example: Sales performance analysis
WITH monthly_sales AS (
    SELECT
        DATE_TRUNC('month', saletime) as month,
    sellerid,
    COUNT(*) as sale_count,
    SUM(pricepaid) as total_revenue,
    AVG(commission) as avg_commission
FROM sales
GROUP BY 1, 2
    ),
    seller_ranks AS (
SELECT
    month,
    sellerid,
    sale_count,
    total_revenue,
    avg_commission,
    RANK() OVER (PARTITION BY month ORDER BY total_revenue DESC) as revenue_rank,
    PERCENT_RANK() OVER (PARTITION BY month ORDER BY sale_count DESC) as sale_count_percentile
FROM monthly_sales
    )
SELECT
    sr.*,
    u.username,
    u.firstname,
    u.lastname,
    CASE
        WHEN revenue_rank <= 10 THEN 'Top Performer'
        WHEN sale_count_percentile <= 0.25 THEN 'High Volume'
        ELSE 'Standard'
        END as seller_category
INTO TEMP TABLE seller_performance_report
FROM seller_ranks sr
    JOIN users u ON sr.sellerid = u.userid
WHERE sr.month >= DATE_TRUNC('month', DATEADD(month, -12, CURRENT_DATE))
ORDER BY sr.month DESC, sr.revenue_rank;