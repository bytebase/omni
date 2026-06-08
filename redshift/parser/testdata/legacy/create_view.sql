-- CREATE VIEW test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_VIEW.html

-- Basic view creation (from AWS documentation)
CREATE VIEW myevent AS
SELECT eventname FROM event
WHERE eventname = 'LeAnn Rimes';

-- OR REPLACE view (from AWS documentation)
CREATE OR REPLACE VIEW myuser AS
SELECT lastname FROM users;

-- View with WITH NO SCHEMA BINDING (from AWS documentation)
CREATE VIEW myevent AS
SELECT eventname FROM public.event
WITH NO SCHEMA BINDING;

-- View with explicit column list
CREATE VIEW event_summary (event_name, venue_name) AS
SELECT eventname, venuename
FROM event
WHERE eventid < 100;

-- OR REPLACE with column list
CREATE OR REPLACE VIEW user_info (last_name, first_name) AS
SELECT lastname, firstname
FROM users
WHERE userid < 500;

-- View with column list and WITH NO SCHEMA BINDING
CREATE VIEW sales_summary (product, total_sales, avg_price) AS
SELECT productname, SUM(quantity), AVG(pricepaid)
FROM sales
GROUP BY productname
WITH NO SCHEMA BINDING;

-- Complex view with JOINs
CREATE VIEW event_sales AS
SELECT
    e.eventname,
    e.starttime,
    v.venuename,
    COUNT(*) as tickets_sold,
    SUM(s.pricepaid) as total_revenue
FROM event e
JOIN venue v ON e.venueid = v.venueid
JOIN sales s ON e.eventid = s.eventid
GROUP BY e.eventname, e.starttime, v.venuename;

-- View with WHERE clause and aggregations
CREATE VIEW high_value_sales AS
SELECT
    buyerid,
    COUNT(*) as purchase_count,
    SUM(pricepaid) as total_spent,
    AVG(pricepaid) as avg_purchase
FROM sales
WHERE pricepaid > 100
GROUP BY buyerid
HAVING SUM(pricepaid) > 1000;

-- View with subquery
CREATE VIEW venue_capacity AS
SELECT
    v.venuename,
    v.venuecity,
    v.venueseats,
    avg_sales.avg_tickets_per_event
FROM venue v
JOIN (
    SELECT
        venueid,
        AVG(tickets_sold) as avg_tickets_per_event
    FROM (
        SELECT
            e.venueid,
            e.eventid,
            COUNT(*) as tickets_sold
        FROM event e
        JOIN sales s ON e.eventid = s.eventid
        GROUP BY e.venueid, e.eventid
    ) event_sales
    GROUP BY venueid
) avg_sales ON v.venueid = avg_sales.venueid;

-- View with CASE statement
CREATE VIEW user_category AS
SELECT
    userid,
    firstname,
    lastname,
    CASE
        WHEN total_purchases >= 10 THEN 'VIP'
        WHEN total_purchases >= 5 THEN 'Regular'
        ELSE 'Occasional'
    END as customer_type
FROM (
    SELECT
        u.userid,
        u.firstname,
        u.lastname,
        COUNT(s.saleid) as total_purchases
    FROM users u
    LEFT JOIN sales s ON u.userid = s.buyerid
    GROUP BY u.userid, u.firstname, u.lastname
);

-- View with window functions
CREATE VIEW sales_ranking AS
SELECT
    saleid,
    buyerid,
    eventid,
    pricepaid,
    RANK() OVER (PARTITION BY eventid ORDER BY pricepaid DESC) as price_rank,
    ROW_NUMBER() OVER (ORDER BY pricepaid DESC) as overall_rank
FROM sales;

-- View with UNION
CREATE VIEW all_names AS
SELECT 'User' as type, firstname as name, lastname
FROM users
UNION ALL
SELECT 'Venue' as type, venuename as name, venuecity
FROM venue;

-- View with date functions and filtering
CREATE VIEW recent_events AS
SELECT
    eventid,
    eventname,
    starttime,
    EXTRACT(YEAR FROM starttime) as event_year,
    EXTRACT(MONTH FROM starttime) as event_month
FROM event
WHERE starttime >= '2023-01-01'
ORDER BY starttime;

-- View with string functions
CREATE VIEW formatted_users AS
SELECT
    userid,
    UPPER(firstname) as first_name_upper,
    LOWER(lastname) as last_name_lower,
    firstname || ' ' || lastname as full_name
FROM users
WHERE firstname IS NOT NULL AND lastname IS NOT NULL;

-- View with NULL handling
CREATE VIEW complete_venues AS
SELECT
    venueid,
    venuename,
    COALESCE(venuecity, 'Unknown City') as city,
    COALESCE(venuestate, 'Unknown State') as state,
    CASE
        WHEN venuecity IS NOT NULL AND venuestate IS NOT NULL THEN 'Complete'
        WHEN venuecity IS NOT NULL OR venuestate IS NOT NULL THEN 'Partial'
        ELSE 'Incomplete'
    END as address_completeness
FROM venue;

-- View with mathematical operations
CREATE VIEW sales_with_commission AS
SELECT
    saleid,
    eventid,
    pricepaid,
    commission * 0.01 as commission_rate,
    pricepaid * commission * 0.01 as commission_amount,
    pricepaid - (pricepaid * commission * 0.01) as net_amount
FROM sales
WHERE commission > 0;

-- View with EXISTS clause
CREATE VIEW events_with_sales AS
SELECT
    eventid,
    eventname,
    starttime,
    venueid
FROM event e
WHERE EXISTS (
    SELECT 1
    FROM sales s
    WHERE s.eventid = e.eventid
);

-- View with NOT EXISTS clause
CREATE VIEW events_without_sales AS
SELECT
    eventid,
    eventname,
    starttime,
    venueid
FROM event e
WHERE NOT EXISTS (
    SELECT 1
    FROM sales s
    WHERE s.eventid = e.eventid
);

-- View with IN clause
CREATE VIEW premium_category_events AS
SELECT
    eventid,
    eventname,
    catid,
    starttime
FROM event
WHERE catid IN (1, 2, 3);

-- View with BETWEEN clause
CREATE VIEW mid_range_sales AS
SELECT
    saleid,
    buyerid,
    eventid,
    pricepaid,
    saledate
FROM sales
WHERE pricepaid BETWEEN 50 AND 500;

-- View with LIKE pattern matching
CREATE VIEW rock_events AS
SELECT
    eventid,
    eventname,
    catid,
    starttime
FROM event
WHERE eventname LIKE '%Rock%' OR eventname LIKE '%Concert%';

-- View with multiple JOINs
CREATE VIEW detailed_sales AS
SELECT
    s.saleid,
    s.saledate,
    s.pricepaid,
    u.firstname,
    u.lastname,
    e.eventname,
    e.starttime,
    v.venuename,
    v.venuecity,
    c.catname
FROM sales s
JOIN users u ON s.buyerid = u.userid
JOIN event e ON s.eventid = e.eventid
JOIN venue v ON e.venueid = v.venueid
JOIN category c ON e.catid = c.catid;

-- View with LEFT JOIN
CREATE VIEW user_sales_summary AS
SELECT
    u.userid,
    u.firstname,
    u.lastname,
    COUNT(s.saleid) as total_purchases,
    COALESCE(SUM(s.pricepaid), 0) as total_spent
FROM users u
LEFT JOIN sales s ON u.userid = s.buyerid
GROUP BY u.userid, u.firstname, u.lastname;

-- View with RIGHT JOIN
CREATE VIEW event_sales_summary AS
SELECT
    e.eventid,
    e.eventname,
    COUNT(s.saleid) as tickets_sold,
    COALESCE(SUM(s.pricepaid), 0) as revenue
FROM sales s
RIGHT JOIN event e ON s.eventid = e.eventid
GROUP BY e.eventid, e.eventname;

-- View with FULL OUTER JOIN
CREATE VIEW complete_user_event_matrix AS
SELECT
    COALESCE(u.userid, 0) as user_id,
    COALESCE(u.firstname, 'Unknown') as first_name,
    COALESCE(e.eventid, 0) as event_id,
    COALESCE(e.eventname, 'No Event') as event_name,
    s.pricepaid
FROM users u
FULL OUTER JOIN sales s ON u.userid = s.buyerid
FULL OUTER JOIN event e ON s.eventid = e.eventid;

-- View with aggregate functions
CREATE VIEW category_statistics AS
SELECT
    c.catid,
    c.catname,
    COUNT(e.eventid) as total_events,
    COUNT(s.saleid) as total_sales,
    MIN(s.pricepaid) as min_price,
    MAX(s.pricepaid) as max_price,
    AVG(s.pricepaid) as avg_price,
    SUM(s.pricepaid) as total_revenue
FROM category c
LEFT JOIN event e ON c.catid = e.catid
LEFT JOIN sales s ON e.eventid = s.eventid
GROUP BY c.catid, c.catname;

-- View with DISTINCT
CREATE VIEW unique_venues_per_city AS
SELECT DISTINCT
    venuecity,
    venuestate,
    COUNT(DISTINCT venueid) as venue_count
FROM venue
WHERE venuecity IS NOT NULL
GROUP BY venuecity, venuestate;

-- OR REPLACE with complex query and WITH NO SCHEMA BINDING
CREATE OR REPLACE VIEW comprehensive_sales_analysis AS
SELECT
    sales_month,
    sales_year,
    category,
    total_events,
    total_sales,
    total_revenue,
    avg_price_per_ticket,
    RANK() OVER (PARTITION BY sales_year ORDER BY total_revenue DESC) as monthly_revenue_rank
FROM (
    SELECT
        EXTRACT(MONTH FROM e.starttime) as sales_month,
        EXTRACT(YEAR FROM e.starttime) as sales_year,
        c.catname as category,
        COUNT(DISTINCT e.eventid) as total_events,
        COUNT(s.saleid) as total_sales,
        SUM(s.pricepaid) as total_revenue,
        AVG(s.pricepaid) as avg_price_per_ticket
    FROM public.event e
    JOIN public.category c ON e.catid = c.catid
    JOIN public.sales s ON e.eventid = s.eventid
    WHERE e.starttime >= '2023-01-01'
    GROUP BY
        EXTRACT(MONTH FROM e.starttime),
        EXTRACT(YEAR FROM e.starttime),
        c.catname
) monthly_stats
WITH NO SCHEMA BINDING;

-- View with nested CASE statements
CREATE VIEW customer_segmentation AS
SELECT
    userid,
    firstname,
    lastname,
    total_spent,
    purchase_count,
    CASE
        WHEN total_spent >= 1000 THEN
            CASE
                WHEN purchase_count >= 10 THEN 'Premium Frequent'
                ELSE 'Premium Occasional'
            END
        WHEN total_spent >= 500 THEN 'Standard'
        WHEN total_spent >= 100 THEN 'Basic'
        ELSE 'Minimal'
    END as customer_segment
FROM (
    SELECT
        u.userid,
        u.firstname,
        u.lastname,
        COALESCE(SUM(s.pricepaid), 0) as total_spent,
        COUNT(s.saleid) as purchase_count
    FROM users u
    LEFT JOIN sales s ON u.userid = s.buyerid
    GROUP BY u.userid, u.firstname, u.lastname
);

-- Simple view for basic functionality testing
CREATE VIEW simple_event_list AS
SELECT eventid, eventname
FROM event;

-- View with schema-qualified table names
CREATE VIEW public_event_details AS
SELECT
    public.event.eventid,
    public.event.eventname,
    public.venue.venuename
FROM public.event
JOIN public.venue ON public.event.venueid = public.venue.venueid;

-- View with quoted identifiers
CREATE VIEW "Special View Name" AS
SELECT
    "eventid" as "Event ID",
    "eventname" as "Event Name"
FROM "event"
WHERE "eventid" < 100;

-- View with column aliases using AS keyword
CREATE VIEW event_venue_info AS
SELECT
    e.eventid AS event_identifier,
    e.eventname AS event_title,
    e.starttime AS event_datetime,
    v.venuename AS venue_title,
    v.venuecity AS venue_location
FROM event AS e
JOIN venue AS v ON e.venueid = v.venueid;

-- View with column aliases without AS keyword
CREATE VIEW sales_info AS
SELECT
    saleid sale_identifier,
    buyerid buyer_identifier,
    eventid event_identifier,
    pricepaid ticket_price,
    saledate purchase_date
FROM sales;

-- View with ORDER BY
CREATE VIEW top_events_by_revenue AS
SELECT
    e.eventid,
    e.eventname,
    SUM(s.pricepaid) as total_revenue,
    COUNT(*) as tickets_sold
FROM event e
JOIN sales s ON e.eventid = s.eventid
GROUP BY e.eventid, e.eventname
ORDER BY total_revenue DESC
LIMIT 10;

-- View with multiple ORDER BY columns
CREATE VIEW events_by_date_and_name AS
SELECT
    eventid,
    eventname,
    starttime,
    venueid
FROM event
ORDER BY starttime ASC, eventname DESC;