-- The span.tsv corpus: one query per line, run in the session the last
-- "-- session:" line names on the fixture the last "-- fixture:" line names.
-- The session is the default dataset, omni_span, which the fixture builds.
-- See ../README.md.

-- fixture: shop
-- session: omni_span

-- Columns, aliases and expressions.
SELECT CustomerId, Name FROM Customers
SELECT c.CustomerId AS id, c.Email FROM Customers AS c
SELECT c.CustomerId FROM Customers c
SELECT Amount * 2, UPPER(Status) FROM Orders
SELECT CAST(Amount AS STRING) AS a, SAFE_CAST(Status AS INT64) FROM Orders
SELECT 1, 'a' AS b
SELECT COUNT(*) FROM Orders
SELECT DISTINCT Region FROM Customers
SELECT Name FROM Customers ORDER BY Region LIMIT 5
SELECT OrderId FROM Orders LIMIT 2 OFFSET 1
SELECT CONCAT(Name, Email) AS contact FROM Customers
SELECT IF(Amount > 10, Status, NULL) AS s FROM Orders
SELECT CASE WHEN Amount > 10 THEN 'big' ELSE Status END FROM Orders
SELECT COALESCE(Email, Name) FROM Customers
SELECT (Name), ((c.Email)) FROM Customers c
SELECT Name AS x, Email AS x FROM Customers
SELECT Name, Name FROM Customers
SELECT 1 AS f0_, 2
SELECT CustomerId + OrderId AS k, Amount FROM Orders WHERE Status = 'open'
SELECT EXTRACT(YEAR FROM Placed) AS y, FORMAT_DATE('%Y', Placed) FROM Orders
SELECT CURRENT_DATE() AS today, Name FROM Customers

-- Qualified names.
SELECT Name FROM omni_span.Customers
SELECT c.Name FROM omni_span.Customers AS c
SELECT Name FROM `omni_span.Customers`
SELECT Name FROM `omni_span`.`Customers`
SELECT Name FROM `bytebase-dev.omni_span.Customers`
SELECT Name FROM `bytebase-dev`.omni_span.Customers
SELECT Customers.Name FROM Customers
SELECT omni_span.Customers.Name FROM omni_span.Customers
SELECT c.Name, o.Amount FROM omni_span.Customers c JOIN Orders o USING (CustomerId)

-- Stars.
SELECT * FROM Orders
SELECT o.* FROM Orders o
SELECT *, Amount FROM Orders
SELECT c.*, o.Amount FROM Customers c JOIN Orders o ON o.CustomerId = c.CustomerId
SELECT * FROM (SELECT * FROM Customers) c
SELECT * EXCEPT (Email) FROM Customers
SELECT * REPLACE (UPPER(Name) AS Name) FROM Customers
SELECT c.* EXCEPT (Region), o.OrderId FROM Customers c JOIN Orders o USING (CustomerId)
SELECT * EXCEPT (Name) FROM (SELECT (Name), Email FROM Customers)
SELECT * EXCEPT (Items, Ship) REPLACE (Amount * 2 AS Amount) FROM Orders
SELECT * FROM Customers, Targets
SELECT * FROM Customers c JOIN Targets t USING (Region)
SELECT t.* FROM (SELECT Name, Email FROM Customers) t

-- Joins.
SELECT c.Name, o.Amount FROM Customers c JOIN Orders o ON c.CustomerId = o.CustomerId
SELECT c.Name, o.Amount FROM Customers c INNER JOIN Orders o ON c.CustomerId = o.CustomerId
SELECT c.Name, o.Amount FROM Customers c LEFT JOIN Orders o ON c.CustomerId = o.CustomerId
SELECT c.Name, o.Amount FROM Customers c LEFT OUTER JOIN Orders o ON c.CustomerId = o.CustomerId
SELECT c.Name, o.Amount FROM Customers c RIGHT JOIN Orders o ON c.CustomerId = o.CustomerId
SELECT c.Name, o.Amount FROM Customers c FULL JOIN Orders o ON c.CustomerId = o.CustomerId
SELECT c.Name, p.Name FROM Customers c CROSS JOIN Products p
SELECT c.Name, p.Name FROM Customers c, Products p
SELECT CustomerId, Name, Amount FROM Customers JOIN Orders USING (CustomerId)
SELECT * FROM Customers JOIN Orders USING (CustomerId)
SELECT * FROM Orders JOIN LineItems USING (CustomerId, OrderId)
SELECT CustomerId FROM Customers FULL JOIN Orders USING (CustomerId)
SELECT e.Name, m.Name AS manager FROM Employees e LEFT JOIN Employees m ON e.ManagerId = m.EmpId
SELECT l.Qty, p.Price, c.Region FROM LineItems l JOIN Products p ON l.ProductId = p.ProductId JOIN Customers c ON c.CustomerId = l.CustomerId
SELECT o.OrderId FROM Orders o JOIN LineItems l ON o.OrderId = l.OrderId AND o.CustomerId = l.CustomerId
SELECT c.Name FROM Customers c JOIN (Orders o JOIN LineItems l USING (OrderId)) ON c.CustomerId = o.CustomerId
SELECT a.Name, b.Name FROM Customers a JOIN Customers b ON a.Region = b.Region
SELECT Name, Goal FROM Customers JOIN Targets USING (Region)
SELECT c.Name, t.Goal FROM Customers c LEFT JOIN Targets t ON c.Region = t.Region WHERE t.Goal IS NULL

-- Subqueries.
SELECT Name, (SELECT MAX(Amount) FROM Orders) AS top FROM Customers
SELECT Name, (SELECT MAX(o.Amount) FROM Orders o WHERE o.CustomerId = c.CustomerId) AS top FROM Customers c
SELECT Name FROM Customers WHERE CustomerId IN (SELECT CustomerId FROM Orders)
SELECT Name FROM Customers c WHERE EXISTS (SELECT 1 FROM Orders o WHERE o.CustomerId = c.CustomerId)
SELECT Name FROM Customers c WHERE NOT EXISTS (SELECT 1 FROM Orders o WHERE o.CustomerId = c.CustomerId)
SELECT t.total FROM (SELECT CustomerId, SUM(Amount) AS total FROM Orders GROUP BY CustomerId) t
SELECT x FROM (SELECT Name AS x FROM Customers UNION ALL SELECT Status FROM Orders)
SELECT Name FROM Customers WHERE Region = (SELECT Region FROM Targets ORDER BY Goal DESC LIMIT 1)
SELECT ARRAY(SELECT OrderId FROM Orders o WHERE o.CustomerId = c.CustomerId) AS ids FROM Customers c
SELECT Name, ARRAY(SELECT AS STRUCT o.OrderId, o.Amount FROM Orders o WHERE o.CustomerId = c.CustomerId) AS orders FROM Customers c
SELECT (SELECT COUNT(*) FROM LineItems) AS n
SELECT * FROM (SELECT * FROM (SELECT Name, Region FROM Customers))
SELECT Name FROM Customers WHERE Region IN (SELECT Region FROM Targets WHERE Goal > (SELECT AVG(Goal) FROM Targets))
SELECT Name FROM Customers WHERE CustomerId = ANY (SELECT CustomerId FROM Orders)
SELECT c.Name, x.n FROM Customers c, (SELECT COUNT(*) AS n FROM Orders) x
SELECT EXISTS(SELECT 1 FROM Orders) AS has_orders

-- CTEs.
WITH big AS (SELECT CustomerId, Amount FROM Orders WHERE Amount > 100) SELECT * FROM big
WITH big AS (SELECT CustomerId, Amount FROM Orders) SELECT c.Name, b.Amount FROM Customers c JOIN big b USING (CustomerId)
WITH a AS (SELECT Name FROM Customers), b AS (SELECT Name FROM a) SELECT Name FROM b
WITH Orders AS (SELECT Name AS Status FROM Customers) SELECT Status FROM Orders
WITH x AS (SELECT Region FROM Customers) SELECT Region FROM x UNION ALL SELECT Region FROM Targets
WITH RECURSIVE r AS (SELECT EmpId, ManagerId FROM Employees WHERE ManagerId IS NULL UNION ALL SELECT e.EmpId, e.ManagerId FROM Employees e JOIN r ON e.ManagerId = r.EmpId) SELECT EmpId FROM r
WITH RECURSIVE n AS (SELECT 1 AS i UNION ALL SELECT i + 1 FROM n WHERE i < 3) SELECT i FROM n
WITH t AS (SELECT * FROM Products) SELECT t.Name, t.Price FROM t
WITH t (a, b) AS (SELECT Name, Email FROM Customers) SELECT a FROM t
WITH t AS (SELECT 1 AS one) SELECT one, Name FROM t, Customers
WITH s AS (SELECT CustomerId, SUM(Amount) AS total FROM Orders GROUP BY CustomerId) SELECT c.Name, s.total FROM s JOIN Customers c USING (CustomerId)
WITH a AS (SELECT * FROM Customers) SELECT * FROM a WHERE Region = 'eu'
WITH Customers AS (SELECT * FROM omni_span.Customers WHERE Region = 'eu') SELECT Name FROM Customers

-- Set operations.
SELECT Name FROM Customers UNION ALL SELECT Name FROM Products
SELECT Name FROM Customers UNION DISTINCT SELECT Name FROM Employees
SELECT CustomerId FROM Customers INTERSECT DISTINCT SELECT CustomerId FROM Orders
SELECT CustomerId FROM Customers EXCEPT DISTINCT SELECT CustomerId FROM Orders
SELECT Name, Email FROM Customers UNION ALL SELECT Status, CAST(Amount AS STRING) FROM Orders
SELECT ProductId, Name FROM Products UNION ALL SELECT ProductId, Name FROM Products
(SELECT Name FROM Customers) UNION ALL (SELECT Name FROM Products ORDER BY Name LIMIT 1)
SELECT Name FROM Customers UNION ALL SELECT Name FROM Products UNION ALL SELECT Name FROM Employees
SELECT * FROM Targets UNION ALL SELECT Region, Goal FROM Targets
SELECT Name, Region FROM Customers UNION ALL BY NAME SELECT Region, Name FROM Customers
SELECT Name FROM Customers UNION ALL CORRESPONDING SELECT Name FROM Employees
SELECT 1 AS a UNION ALL SELECT 2
SELECT Name FROM Customers FULL UNION ALL BY NAME SELECT Name, Salary FROM Employees
SELECT Name, Region FROM Customers LEFT UNION ALL BY NAME SELECT Name FROM Employees

-- Views.
SELECT * FROM CustomerOrders
SELECT Name, Amount FROM CustomerOrders WHERE OrderId > 1
SELECT * FROM BigOrders
SELECT b.Name FROM BigOrders b JOIN Customers c ON b.CustomerId = c.CustomerId
SELECT * FROM RegionGoals
SELECT g.Goal, c.Name FROM RegionGoals g JOIN Customers c ON c.Region = g.Region
SELECT COUNT(*) FROM CustomerOrders
SELECT Name FROM omni_span.BigOrders
WITH v AS (SELECT * FROM CustomerOrders) SELECT Name FROM v

-- Table functions.
SELECT * FROM OrdersOf(1)
SELECT OrderId FROM omni_span.OrdersOf(2)
SELECT c.Name, f.Amount FROM Customers c, OrdersOf(c.CustomerId) f

-- Wildcard tables.
SELECT EventId FROM `events_*`
SELECT * FROM `events_*`
SELECT Kind FROM `omni_span.events_*` WHERE _TABLE_SUFFIX = '2024'
SELECT _TABLE_SUFFIX AS year, Kind FROM `events_*`
SELECT e.Payload FROM `events_20*` AS e WHERE e._TABLE_SUFFIX BETWEEN '24' AND '25'
SELECT COUNT(*) FROM `events_*` WHERE Kind = 'click'
SELECT * FROM events_2024 UNION ALL SELECT * FROM events_2025
SELECT EventId FROM `events_2024`

-- Quoted and mixed-case names.
SELECT `Id`, `CamelCol` FROM `MixedCase`
SELECT id, camelcol, LOWER_COL FROM MixedCase
SELECT m.camelCol FROM MixedCase m
SELECT * FROM MixedCase
SELECT `Customers`.`Name` FROM `Customers`
SELECT name FROM customers
SELECT NAME FROM Customers
SELECT c.NAME, C.email FROM Customers c

-- Aliases that shadow table names.
SELECT Orders.Name FROM Customers AS Orders
SELECT Customers.Amount FROM Orders AS Customers JOIN Customers AS c ON Customers.CustomerId = c.CustomerId
SELECT Products.Name FROM (SELECT Name FROM Employees) AS Products

-- Window functions and QUALIFY.
SELECT Name, ROW_NUMBER() OVER (ORDER BY CustomerId) AS rn FROM Customers
SELECT OrderId, SUM(Amount) OVER (PARTITION BY CustomerId) FROM Orders
SELECT Name, RANK() OVER (PARTITION BY Region ORDER BY Name) FROM Customers
SELECT Name, LAG(Name) OVER w AS prev FROM Customers WINDOW w AS (ORDER BY CustomerId)
SELECT Name FROM Customers QUALIFY ROW_NUMBER() OVER (PARTITION BY Region ORDER BY Name) = 1
SELECT OrderId, Amount FROM Orders WHERE TRUE QUALIFY RANK() OVER (ORDER BY Amount DESC) <= 3
SELECT FIRST_VALUE(Name) OVER (PARTITION BY Region ORDER BY CustomerId) AS first FROM Customers

-- Aggregates.
SELECT Region, COUNT(*) AS n FROM Customers GROUP BY Region
SELECT CustomerId, SUM(Amount), MAX(Status) FROM Orders GROUP BY CustomerId HAVING SUM(Amount) > 10
SELECT ARRAY_AGG(Name) FROM Customers
SELECT STRING_AGG(Status, ',') AS statuses FROM Orders
SELECT COUNT(DISTINCT Region) FROM Customers
SELECT ANY_VALUE(Email) AS e FROM Customers GROUP BY Region
SELECT Region, COUNT(*) FROM Customers GROUP BY 1
SELECT Region, Name, COUNT(*) FROM Customers GROUP BY ROLLUP (Region, Name)
SELECT Region, COUNT(*) FROM Customers GROUP BY ALL
SELECT ARRAY_AGG(STRUCT(OrderId, Amount) ORDER BY Amount LIMIT 3) AS top FROM Orders
SELECT APPROX_COUNT_DISTINCT(Email) FROM Customers
SELECT Region, LOGICAL_OR(Email IS NULL) FROM Customers GROUP BY Region
SELECT AVG(Amount) AS avg_amount, MIN(Placed) FROM Orders

-- UNNEST and arrays.
SELECT item FROM Orders, UNNEST(Items) AS item
SELECT o.OrderId, item FROM Orders o CROSS JOIN UNNEST(o.Items) AS item
SELECT item, off FROM Orders o, UNNEST(o.Items) AS item WITH OFFSET AS off
SELECT x FROM UNNEST([1, 2, 3]) AS x
SELECT * FROM UNNEST(['a', 'b'])
SELECT o.OrderId, (SELECT COUNT(*) FROM UNNEST(o.Items)) AS n FROM Orders o
SELECT Items[SAFE_OFFSET(0)] AS first FROM Orders
SELECT ARRAY_LENGTH(Items) FROM Orders
SELECT o.OrderId, i FROM Orders o, o.Items AS i
SELECT o.OrderId, i FROM Orders o LEFT JOIN UNNEST(o.Items) AS i
SELECT t.k, t.v FROM Products, UNNEST(Tags) AS t
SELECT p.Name, tag.k FROM Products p CROSS JOIN UNNEST(p.Tags) AS tag WHERE tag.v = 'x'
SELECT * FROM Products, UNNEST(Tags)
SELECT p.Name, (SELECT STRING_AGG(t.v) FROM UNNEST(p.Tags) t) AS vals FROM Products p
SELECT Items[OFFSET(0)], Tags[SAFE_ORDINAL(1)].k FROM Orders, Products
SELECT ARRAY_CONCAT(Items, ['x']) AS more FROM Orders
SELECT * FROM UNNEST([STRUCT(1 AS a, 'x' AS b), (2, 'y')])
SELECT a, b FROM UNNEST([STRUCT(1 AS a, 'x' AS b)])
SELECT x, y FROM UNNEST([1, 2]) AS x, UNNEST([3]) AS y
SELECT item FROM UNNEST((SELECT Items FROM Orders LIMIT 1)) AS item

-- STRUCT fields and SELECT AS.
SELECT Ship.City, Ship.Zip FROM Orders
SELECT o.Ship.City AS city FROM Orders o
SELECT Ship FROM Orders
SELECT Ship.* FROM Orders
SELECT o.Ship.*, o.OrderId FROM Orders o
SELECT s.a, s.b FROM (SELECT STRUCT(CustomerId AS a, Name AS b) AS s FROM Customers)
SELECT STRUCT(CustomerId, Name).Name AS n FROM Customers
SELECT t.s.x FROM (SELECT STRUCT(Amount AS x) AS s FROM Orders) AS t
SELECT v.CustomerId FROM (SELECT AS STRUCT CustomerId, Name FROM Customers) AS v
SELECT AS STRUCT CustomerId, Name FROM Customers
SELECT AS VALUE STRUCT(CustomerId AS id) FROM Customers
SELECT AS VALUE Name FROM Customers
SELECT e.OrderId FROM Customers c, UNNEST(ARRAY(SELECT AS STRUCT o.OrderId FROM Orders o WHERE o.CustomerId = c.CustomerId)) AS e
SELECT STRUCT(Name, Email) AS contact FROM Customers
SELECT [Name, Email] AS arr FROM Customers
SELECT (SELECT AS STRUCT Name, Email) AS s FROM Customers
SELECT c FROM Customers c
SELECT TO_JSON_STRING(c) FROM Customers c
SELECT ANY_VALUE(Ship).City FROM Orders

-- PIVOT and UNPIVOT.
SELECT * FROM Sales PIVOT (SUM(Revenue) FOR Quarter IN ('Q1', 'Q2'))
SELECT Region, Q1 FROM Sales PIVOT (SUM(Revenue) AS Q FOR Quarter IN ('Q1' AS Q1, 'Q2' AS Q2))
SELECT * FROM (SELECT Region, Quarter, Revenue FROM Sales) PIVOT (MAX(Revenue) FOR Quarter IN ('Q1'))
SELECT * FROM Targets UNPIVOT (val FOR col IN (Goal))
SELECT Region, val FROM (SELECT Region, Revenue, Revenue * 2 AS Double FROM Sales) UNPIVOT (val FOR kind IN (Revenue, Double))

-- TABLESAMPLE and time travel.
SELECT Name FROM Customers TABLESAMPLE SYSTEM (10 PERCENT)
SELECT o.Amount FROM Orders o TABLESAMPLE SYSTEM (5 PERCENT)
SELECT Name FROM Customers FOR SYSTEM_TIME AS OF CURRENT_TIMESTAMP()
SELECT c.Name FROM Customers FOR SYSTEM TIME AS OF CURRENT_TIMESTAMP() AS c

-- Pipe syntax.
FROM Customers
FROM Customers |> SELECT Name
FROM Customers |> WHERE Region = 'eu' |> SELECT Name, Email
FROM Orders |> AGGREGATE SUM(Amount) AS total GROUP BY CustomerId
FROM Orders |> EXTEND Amount * 2 AS double
FROM Orders |> EXTEND Amount * 2 AS double |> SELECT OrderId, double
FROM Customers |> SET Name = UPPER(Name)
FROM Customers |> DROP Email
FROM Customers |> RENAME Name AS full_name
FROM Customers |> AS c |> SELECT c.Name
FROM Customers c |> JOIN Orders o ON c.CustomerId = o.CustomerId |> SELECT c.Name, o.Amount
FROM Customers |> JOIN Orders USING (CustomerId) |> SELECT Name, Amount
FROM Customers |> ORDER BY Name |> LIMIT 3
FROM Customers |> SELECT Name |> UNION ALL (SELECT Name FROM Products)
FROM Orders |> AGGREGATE COUNT(*) AS n GROUP BY Status |> WHERE n > 1
FROM Customers |> LEFT JOIN Targets USING (Region) |> SELECT Name, Goal
FROM Orders |> CROSS JOIN UNNEST(Items) AS item |> SELECT OrderId, item
FROM Sales |> PIVOT (SUM(Revenue) FOR Quarter IN ('Q1', 'Q2'))
FROM Customers |> WINDOW ROW_NUMBER() OVER (ORDER BY CustomerId) AS rn
FROM Customers |> EXTEND ROW_NUMBER() OVER (ORDER BY CustomerId) AS rn |> SELECT Name, rn
FROM Customers |> DISTINCT
FROM Customers |> SELECT Region |> DISTINCT
WITH t AS (SELECT * FROM Orders) FROM t |> SELECT OrderId
SELECT * FROM Customers |> WHERE Region = 'eu'
SELECT Name, Email FROM Customers |> SELECT Name
FROM Customers |> TABLESAMPLE SYSTEM (10 PERCENT) |> SELECT Name
FROM CustomerOrders |> SELECT Name
FROM `events_*` |> SELECT Kind
FROM Customers |> SELECT * EXCEPT (Email)
FROM Customers |> AGGREGATE COUNT(*) GROUP BY Region
FROM Customers |> CALL OrdersOf(1)
FROM Customers c |> SET Name = Email |> SELECT c.Name, Name
FROM Customers c |> DROP Email |> SELECT c.Email, *
FROM Customers c |> RENAME Name AS n |> SELECT c.Name, n
FROM Customers c |> EXTEND UPPER(Name) AS u |> SELECT c.Email, u
FROM Customers c |> WHERE Region = 'eu' |> SELECT c.*
FROM Customers |> SELECT Customers.Name
FROM Customers, Targets |> SELECT Customers.Name, Targets.Goal
FROM Customers |> SET Name = Email, Email = Name
FROM Customers |> EXTEND Email AS Contact |> DROP Email |> RENAME Contact AS Email
FROM Customers |> DROP email, region
FROM MixedCase |> RENAME camelcol AS Renamed
FROM Customers |> SET Name = Email |> AS c |> SELECT c.Name
FROM Customers |> JOIN Orders USING (CustomerId)
FROM Customers c |> JOIN Orders o ON c.CustomerId = o.CustomerId
FROM Customers c |> JOIN Orders o USING (CustomerId) |> SELECT c.CustomerId, o.CustomerId, CustomerId
FROM Customers |> AS c |> JOIN Orders AS o ON c.CustomerId = o.CustomerId |> SELECT c.Name, o.OrderId
FROM Orders o |> CROSS JOIN UNNEST(o.Items) AS item WITH OFFSET |> SELECT item, offset
FROM Products |> JOIN UNNEST(Tags) AS tag |> SELECT Name, tag.k
FROM Orders |> AGGREGATE SUM(Amount) AS total, COUNT(*) GROUP BY ROLLUP(CustomerId, Status)
FROM Orders |> AGGREGATE COUNT(*) GROUP BY GROUPING SETS ((Status), (Status, CustomerId))
FROM Orders |> AGGREGATE MAX(Amount) AS top GROUP AND ORDER BY Status DESC, CustomerId
FROM Orders |> AGGREGATE SUM(Amount) AS total GROUP BY CustomerId AS cid |> JOIN Customers ON cid = Customers.CustomerId |> SELECT Name, total
FROM Orders |> EXTEND Ship.City AS city |> AGGREGATE COUNT(*) AS n GROUP BY city
FROM Customers |> AGGREGATE COUNT(DISTINCT Region) AS regions
FROM Customers |> SELECT Name, Region |> UNION ALL BY NAME (SELECT Region, Name FROM Customers)
FROM Customers |> SELECT Name |> UNION DISTINCT (SELECT Name FROM Products), (SELECT Name FROM Employees)
FROM Customers |> SELECT CustomerId |> INTERSECT DISTINCT (SELECT CustomerId FROM Orders)
FROM Customers |> SELECT * |> SELECT Name
FROM Customers |> SELECT * REPLACE (Email AS Name) |> SELECT Name
FROM Customers |> SELECT * EXCEPT (Email) |> EXTEND 1 AS one
FROM Customers |> EXTEND Name AS Name2 |> SELECT * EXCEPT (Name)
SELECT * FROM Customers |> DROP Email |> SET Name = LOWER(Name)
SELECT c.Name, o.Amount FROM Customers c JOIN Orders o USING (CustomerId) |> EXTEND Amount * 2 AS twice
WITH big AS (FROM Orders |> WHERE Amount > 10 |> SELECT CustomerId, Amount) SELECT c.Name, big.Amount FROM big JOIN Customers c USING (CustomerId)
SELECT * FROM (FROM Customers |> RENAME Email AS mail |> DROP Name)
FROM Customers |> EXTEND (SELECT MAX(Amount) FROM Orders o WHERE o.CustomerId = Customers.CustomerId) AS top
FROM CustomerOrders |> RENAME Amount AS total
FROM Customers |> SELECT Name |> UNION ALL (SELECT Name FROM Products) |> SELECT Name
SELECT c.Name FROM (FROM Customers AS c)
(FROM Customers AS c |> WHERE Region = 'eu') |> SELECT c.Email
(FROM Customers |> AS c) |> SELECT c.Name
FROM (FROM Customers AS c |> SET Name = Email) |> SELECT c.Name, Name
SELECT c.Name FROM (FROM Customers AS c) AS s

-- INFORMATION_SCHEMA and other datasets.
SELECT table_name FROM omni_span.INFORMATION_SCHEMA.TABLES
SELECT column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = 'Orders'
SELECT word FROM `bigquery-public-data.samples.shakespeare`
SELECT word, corpus FROM `bigquery-public-data`.samples.shakespeare LIMIT 1
SELECT s.word, c.Name FROM `bigquery-public-data.samples.shakespeare` s, Customers c

-- Misc expressions.
SELECT Name FROM Customers WHERE Name LIKE 'a%' AND Email IS NOT NULL
SELECT SAFE_DIVIDE(Amount, 2) AS half FROM Orders
SELECT * FROM Customers WHERE CustomerId BETWEEN 1 AND 10
SELECT Name FROM Customers ORDER BY 1 DESC
SELECT IFNULL(Email, 'none') AS email FROM Customers
SELECT Name FROM Customers WHERE STARTS_WITH(Email, 'x') LIMIT 1
SELECT GENERATE_UUID() AS id, Name FROM Customers
SELECT ARRAY(SELECT x FROM UNNEST(Items) AS x WHERE x != '') AS clean FROM Orders
SELECT Name, (SELECT Goal FROM Targets t WHERE t.Region = c.Region) AS goal FROM Customers c
SELECT DATE_TRUNC(Placed, MONTH) AS month, SUM(Amount) FROM Orders GROUP BY month
SELECT REGEXP_EXTRACT(Email, r'@(.*)') AS domain FROM Customers
SELECT Name || ' <' || Email || '>' AS addr FROM Customers
SELECT `Name` FROM Customers WHERE `Region` = 'eu'
SELECT 1 FROM Customers
SELECT NULL AS n FROM Orders

-- Statements a dry run refuses.
SELECT Nope FROM Customers
SELECT Name FROM Missing
SELECT x.Name FROM Customers c
SELECT Name FROM Customers, Products
SELECT * FROM Customers JOIN Orders ON TRUE JOIN Products USING (Name)
