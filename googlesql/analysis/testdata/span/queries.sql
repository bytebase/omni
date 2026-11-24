-- The span.tsv corpus: one query per line, run in the session the last
-- "-- session:" line names on the fixture the last "-- fixture:" line names.
-- Spanner has no session schema: "default" is the default schema, "".
-- See ../README.md.

-- fixture: shop
-- session: default

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
SELECT * EXCEPT (Name) FROM (SELECT (Name), Email FROM Customers)

-- Stars.
SELECT * FROM Orders
SELECT o.* FROM Orders o
SELECT *, Amount FROM Orders
SELECT c.*, o.Amount FROM Customers c JOIN Orders o ON o.CustomerId = c.CustomerId
SELECT * FROM (SELECT * FROM Customers) c
SELECT * EXCEPT (Email) FROM Customers
SELECT * REPLACE (UPPER(Name) AS Name) FROM Customers
SELECT c.* EXCEPT (Region), o.OrderId FROM Customers c JOIN Orders o USING (CustomerId)

-- Joins.
SELECT c.Name, o.Amount FROM Customers c JOIN Orders o ON c.CustomerId = o.CustomerId
SELECT c.Name, o.Amount FROM Customers c LEFT JOIN Orders o ON c.CustomerId = o.CustomerId
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
SELECT o.OrderId FROM Orders o JOIN @{JOIN_METHOD=HASH_JOIN} LineItems l ON o.OrderId = l.OrderId

-- Interleaved tables and indexes.
SELECT l.LineId, o.Status FROM Orders o JOIN LineItems l ON l.CustomerId = o.CustomerId AND l.OrderId = o.OrderId
SELECT LineId FROM LineItems
SELECT OrderId, Status FROM Orders@{FORCE_INDEX=OrdersByStatus} WHERE Status = 'open'
SELECT Status FROM Orders@{FORCE_INDEX=_BASE_TABLE}

-- Subqueries.
SELECT Name, (SELECT MAX(Amount) FROM Orders) AS top FROM Customers
SELECT Name, (SELECT MAX(o.Amount) FROM Orders o WHERE o.CustomerId = c.CustomerId) AS top FROM Customers c
SELECT Name FROM Customers WHERE CustomerId IN (SELECT CustomerId FROM Orders)
SELECT Name FROM Customers c WHERE EXISTS (SELECT 1 FROM Orders o WHERE o.CustomerId = c.CustomerId)
SELECT t.total FROM (SELECT CustomerId, SUM(Amount) AS total FROM Orders GROUP BY CustomerId) t
SELECT x FROM (SELECT Name AS x FROM Customers UNION ALL SELECT Status FROM Orders)
SELECT Name FROM Customers WHERE Region = (SELECT Region FROM sales.Targets ORDER BY Goal DESC LIMIT 1)
SELECT ARRAY(SELECT OrderId FROM Orders o WHERE o.CustomerId = c.CustomerId) AS ids FROM Customers c
SELECT Name, ARRAY(SELECT AS STRUCT o.OrderId, o.Amount FROM Orders o WHERE o.CustomerId = c.CustomerId) AS orders FROM Customers c

-- CTEs.
WITH big AS (SELECT CustomerId, Amount FROM Orders WHERE Amount > 100) SELECT * FROM big
WITH big AS (SELECT CustomerId, Amount FROM Orders) SELECT c.Name, b.Amount FROM Customers c JOIN big b USING (CustomerId)
WITH a AS (SELECT Name FROM Customers), b AS (SELECT Name FROM a) SELECT Name FROM b
WITH Orders AS (SELECT Name AS Status FROM Customers) SELECT Status FROM Orders
WITH x AS (SELECT Region FROM Customers) SELECT Region FROM x UNION ALL SELECT Region FROM sales.Targets
WITH RECURSIVE r AS (SELECT EmpId, ManagerId FROM Employees WHERE ManagerId IS NULL UNION ALL SELECT e.EmpId, e.ManagerId FROM Employees e JOIN r ON e.ManagerId = r.EmpId) SELECT EmpId FROM r
WITH RECURSIVE n AS (SELECT 1 AS i UNION ALL SELECT i + 1 FROM n WHERE i < 3) SELECT i FROM n

-- Set operations.
SELECT Name FROM Customers UNION ALL SELECT Name FROM Products
SELECT Name FROM Customers UNION DISTINCT SELECT Name FROM Employees
SELECT CustomerId FROM Customers INTERSECT DISTINCT SELECT CustomerId FROM Orders
SELECT CustomerId FROM Customers EXCEPT DISTINCT SELECT CustomerId FROM Orders
SELECT Name, Email FROM Customers UNION ALL SELECT Status, CAST(Amount AS STRING) FROM Orders
SELECT * FROM Products UNION ALL SELECT ProductId, Name, Price FROM Products
(SELECT Name FROM Customers) UNION ALL (SELECT Name FROM Products ORDER BY Name LIMIT 1)

-- Views.
SELECT * FROM CustomerOrders
SELECT Name, Amount FROM CustomerOrders WHERE OrderId > 1
SELECT * FROM BigOrders
SELECT b.Name FROM BigOrders b JOIN Customers c ON b.CustomerId = c.CustomerId
SELECT * FROM sales.RegionGoals
SELECT g.Goal, c.Name FROM sales.RegionGoals g JOIN Customers c ON c.Region = g.Region
SELECT COUNT(*) FROM CustomerOrders
SELECT Region FROM sales.Targets

-- Quoted and mixed-case names.
SELECT `Id`, `CamelCol` FROM `MixedCase`
SELECT id, camelcol, LOWER_COL FROM mixedcase
SELECT m.camelCol FROM MIXEDCASE m
SELECT * FROM MixedCase
SELECT `customers`.`name` FROM `customers`
SELECT customers.name FROM customers

-- Aliases that shadow table names.
SELECT Orders.Name FROM Customers AS Orders
SELECT Customers.Amount FROM Orders AS Customers JOIN Customers AS c ON Customers.CustomerId = c.CustomerId
SELECT Products.Name FROM (SELECT Name FROM Employees) AS Products

-- Window functions.
SELECT Name, ROW_NUMBER() OVER (ORDER BY CustomerId) AS rn FROM Customers
SELECT OrderId, SUM(Amount) OVER (PARTITION BY CustomerId) FROM Orders
SELECT Name, RANK() OVER (PARTITION BY Region ORDER BY Name) FROM Customers

-- Aggregates.
SELECT Region, COUNT(*) AS n FROM Customers GROUP BY Region
SELECT CustomerId, SUM(Amount), MAX(Status) FROM Orders GROUP BY CustomerId HAVING SUM(Amount) > 10
SELECT ARRAY_AGG(Name) FROM Customers
SELECT STRING_AGG(Status, ',') AS statuses FROM Orders
SELECT COUNT(DISTINCT Region) FROM Customers
SELECT ANY_VALUE(Email) AS e FROM Customers GROUP BY Region

-- UNNEST and arrays.
SELECT item FROM Orders, UNNEST(Items) AS item
SELECT o.OrderId, item FROM Orders o CROSS JOIN UNNEST(o.Items) AS item
SELECT item, off FROM Orders o, UNNEST(o.Items) AS item WITH OFFSET AS off
SELECT x FROM UNNEST([1, 2, 3]) AS x
SELECT * FROM UNNEST(['a', 'b'])
SELECT o.OrderId, (SELECT COUNT(*) FROM UNNEST(o.Items)) AS n FROM Orders o
SELECT Items[SAFE_OFFSET(0)] AS first FROM Orders
SELECT ARRAY_LENGTH(Items) FROM Orders

-- STRUCT fields and SELECT AS.
SELECT s.a, s.b FROM (SELECT STRUCT(CustomerId AS a, Name AS b) AS s FROM Customers)
SELECT STRUCT(CustomerId, Name).Name AS n FROM Customers
SELECT t.s.x FROM (SELECT STRUCT(Amount AS x) AS s FROM Orders) AS t
SELECT v.CustomerId FROM (SELECT AS STRUCT CustomerId, Name FROM Customers) AS v
SELECT AS STRUCT CustomerId, Name FROM Customers
SELECT AS VALUE STRUCT(CustomerId AS id) FROM Customers
SELECT e.OrderId FROM Customers c, UNNEST(ARRAY(SELECT AS STRUCT o.OrderId FROM Orders o WHERE o.CustomerId = c.CustomerId)) AS e

-- TABLESAMPLE.
SELECT Name FROM Customers TABLESAMPLE BERNOULLI (10 PERCENT)
SELECT o.Amount FROM Orders o TABLESAMPLE RESERVOIR (5 ROWS)
