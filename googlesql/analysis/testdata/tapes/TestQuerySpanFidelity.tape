-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54 judge:2

== TestQuerySpanFidelity default#1
> query
| -- omni:span-fixture
| CREATE TABLE Customers (CustomerId INT64 NOT NULL, Name STRING(100), Email STRING(100), Region STRING(20)) PRIMARY KEY (CustomerId);
| CREATE TABLE Orders (CustomerId INT64 NOT NULL, OrderId INT64 NOT NULL, Amount FLOAT64, Status STRING(20), Items ARRAY<STRING(50)>) PRIMARY KEY (CustomerId, OrderId), INTERLEAVE IN PARENT Customers ON DELETE CASCADE;
| CREATE TABLE LineItems (CustomerId INT64 NOT NULL, OrderId INT64 NOT NULL, LineId INT64 NOT NULL, ProductId INT64, Qty INT64) PRIMARY KEY (CustomerId, OrderId, LineId), INTERLEAVE IN PARENT Orders;
| CREATE TABLE Products (ProductId INT64 NOT NULL, Name STRING(100), Price NUMERIC) PRIMARY KEY (ProductId);
| CREATE TABLE Employees (EmpId INT64 NOT NULL, ManagerId INT64, Name STRING(100), Salary INT64) PRIMARY KEY (EmpId);
| CREATE TABLE `MixedCase` (`Id` INT64 NOT NULL, `CamelCol` STRING(10), `lower_col` STRING(10)) PRIMARY KEY (`Id`);
| CREATE INDEX OrdersByStatus ON Orders(Status);
| CREATE SCHEMA sales;
| CREATE TABLE sales.Targets (Region STRING(20) NOT NULL, Goal INT64) PRIMARY KEY (Region);
| CREATE VIEW CustomerOrders SQL SECURITY INVOKER AS SELECT c.CustomerId, c.Name, o.OrderId, o.Amount FROM Customers AS c JOIN Orders AS o ON c.CustomerId = o.CustomerId;
| CREATE VIEW BigOrders SQL SECURITY INVOKER AS SELECT co.CustomerId, co.Name, co.Amount FROM CustomerOrders AS co WHERE co.Amount > 100;
| CREATE VIEW sales.RegionGoals SQL SECURITY INVOKER AS SELECT c.Region, t.Goal FROM Customers AS c JOIN sales.Targets AS t ON c.Region = t.Region;
| CREATE CHANGE STREAM OrderStream FOR Orders;
< result
< column "verdict" "JSON"
< row
: string
| {"schemas":[{"name":"","tables":[{"name":"Customers","columns":[{"name":"CustomerId","type":"INT64"},{"name":"Name","type":"STRING(100)"},{"name":"Email","type":"STRING(100)"},{"name":"Region","type":"STRING(20)"}]},{"name":"Employees","columns":[{"name":"EmpId","type":"INT64"},{"name":"ManagerId","type":"INT64"},{"name":"Name","type":"STRING(100)"},{"name":"Salary","type":"INT64"}]},{"name":"LineItems","columns":[{"name":"CustomerId","type":"INT64"},{"name":"OrderId","type":"INT64"},{"name":"LineId","type":"INT64"},{"name":"ProductId","type":"INT64"},{"name":"Qty","type":"INT64"}]},{"name":"MixedCase","columns":[{"name":"Id","type":"INT64"},{"name":"CamelCol","type":"STRING(10)"},{"name":"lower_col","type":"STRING(10)"}]},{"name":"Orders","columns":[{"name":"CustomerId","type":"INT64"},{"name":"OrderId","type":"INT64"},{"name":"Amount","type":"FLOAT64"},{"name":"Status","type":"STRING(20)"},{"name":"Items","type":"ARRAY\u003cSTRING(50)\u003e"}]},{"name":"Products","columns":[{"name":"ProductId","type":"INT64"},{"name":"Name","type":"STRING(100)"},{"name":"Price","type":"NUMERIC"}]}],"views":[{"name":"BigOrders","definition":"SELECT co.CustomerId, co.Name, co.Amount FROM CustomerOrders AS co WHERE co.Amount \u003e 100","columns":[{"name":"CustomerId","type":"INT64"},{"name":"Name","type":"STRING(0)"},{"name":"Amount","type":"FLOAT64"}]},{"name":"CustomerOrders","definition":"SELECT c.CustomerId, c.Name, o.OrderId, o.Amount FROM Customers AS c JOIN Orders AS o ON c.CustomerId = o.CustomerId","columns":[{"name":"CustomerId","type":"INT64"},{"name":"Name","type":"STRING(0)"},{"name":"OrderId","type":"INT64"},{"name":"Amount","type":"FLOAT64"}]}]},{"name":"sales","tables":[{"name":"Targets","columns":[{"name":"Region","type":"STRING(20)"},{"name":"Goal","type":"INT64"}]}],"views":[{"name":"RegionGoals","definition":"SELECT c.Region, t.Goal FROM Customers AS c JOIN sales.Targets AS t ON c.Region = t.Region","columns":[{"name":"Region","type":"STRING(0)"},{"name":"Goal","type":"INT64"}]}]}]}
> query
| -- omni:span-query
| SELECT CustomerId, Name FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT c.CustomerId AS id, c.Email FROM Customers AS c
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["id","Email"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT c.CustomerId FROM Customers c
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT Amount * 2, UPPER(Status) FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["",""],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT CAST(Amount AS STRING) AS a, SAFE_CAST(Status AS INT64) FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["a",""],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT 1, 'a' AS b
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["","b"]}
> query
| -- omni:span-query
| SELECT COUNT(*) FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":[""],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT DISTINCT Region FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Region"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT Name FROM Customers ORDER BY Region LIMIT 5
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT OrderId FROM Orders LIMIT 2 OFFSET 1
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["OrderId"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT CONCAT(Name, Email) AS contact FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["contact"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT IF(Amount > 10, Status, NULL) AS s FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["s"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT CASE WHEN Amount > 10 THEN 'big' ELSE Status END FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":[""],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT COALESCE(Email, Name) FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":[""],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT (Name), ((c.Email)) FROM Customers c
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Email"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT * EXCEPT (Name) FROM (SELECT (Name), Email FROM Customers)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Email"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT * FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","OrderId","Amount","Status","Items"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT o.* FROM Orders o
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","OrderId","Amount","Status","Items"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT *, Amount FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","OrderId","Amount","Status","Items","Amount"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT c.*, o.Amount FROM Customers c JOIN Orders o ON o.CustomerId = c.CustomerId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name","Email","Region","Amount"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT * FROM (SELECT * FROM Customers) c
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name","Email","Region"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT * EXCEPT (Email) FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name","Region"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT * REPLACE (UPPER(Name) AS Name) FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name","Email","Region"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT c.* EXCEPT (Region), o.OrderId FROM Customers c JOIN Orders o USING (CustomerId)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name","Email","OrderId"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT c.Name, o.Amount FROM Customers c JOIN Orders o ON c.CustomerId = o.CustomerId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Amount"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT c.Name, o.Amount FROM Customers c LEFT JOIN Orders o ON c.CustomerId = o.CustomerId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Amount"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT c.Name, o.Amount FROM Customers c RIGHT JOIN Orders o ON c.CustomerId = o.CustomerId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Amount"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT c.Name, o.Amount FROM Customers c FULL JOIN Orders o ON c.CustomerId = o.CustomerId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Amount"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT c.Name, p.Name FROM Customers c CROSS JOIN Products p
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Name"],"tables":["Customers","Products"]}
> query
| -- omni:span-query
| SELECT c.Name, p.Name FROM Customers c, Products p
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Name"],"tables":["Customers","Products"]}
> query
| -- omni:span-query
| SELECT CustomerId, Name, Amount FROM Customers JOIN Orders USING (CustomerId)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name","Amount"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT * FROM Customers JOIN Orders USING (CustomerId)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name","Email","Region","OrderId","Amount","Status","Items"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT * FROM Orders JOIN LineItems USING (CustomerId, OrderId)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","OrderId","Amount","Status","Items","LineId","ProductId","Qty"],"tables":["LineItems","Orders"]}
> query
| -- omni:span-query
| SELECT CustomerId FROM Customers FULL JOIN Orders USING (CustomerId)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT e.Name, m.Name AS manager FROM Employees e LEFT JOIN Employees m ON e.ManagerId = m.EmpId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","manager"],"tables":["Employees"]}
> query
| -- omni:span-query
| SELECT l.Qty, p.Price, c.Region FROM LineItems l JOIN Products p ON l.ProductId = p.ProductId JOIN Customers c ON c.CustomerId = l.CustomerId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Qty","Price","Region"],"tables":["Customers","LineItems","Products"]}
> query
| -- omni:span-query
| SELECT o.OrderId FROM Orders o JOIN @{JOIN_METHOD=HASH_JOIN} LineItems l ON o.OrderId = l.OrderId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["OrderId"],"tables":["LineItems","Orders"]}
> query
| -- omni:span-query
| SELECT l.LineId, o.Status FROM Orders o JOIN LineItems l ON l.CustomerId = o.CustomerId AND l.OrderId = o.OrderId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["LineId","Status"],"tables":["LineItems","Orders"]}
> query
| -- omni:span-query
| SELECT LineId FROM LineItems
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["LineId"],"tables":["LineItems"]}
> query
| -- omni:span-query
| SELECT OrderId, Status FROM Orders@{FORCE_INDEX=OrdersByStatus} WHERE Status = 'open'
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["OrderId","Status"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT Status FROM Orders@{FORCE_INDEX=_BASE_TABLE}
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Status"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT Name, (SELECT MAX(Amount) FROM Orders) AS top FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","top"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT Name, (SELECT MAX(o.Amount) FROM Orders o WHERE o.CustomerId = c.CustomerId) AS top FROM Customers c
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","top"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT Name FROM Customers WHERE CustomerId IN (SELECT CustomerId FROM Orders)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT Name FROM Customers c WHERE EXISTS (SELECT 1 FROM Orders o WHERE o.CustomerId = c.CustomerId)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT t.total FROM (SELECT CustomerId, SUM(Amount) AS total FROM Orders GROUP BY CustomerId) t
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["total"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT x FROM (SELECT Name AS x FROM Customers UNION ALL SELECT Status FROM Orders)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["x"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT Name FROM Customers WHERE Region = (SELECT Region FROM sales.Targets ORDER BY Goal DESC LIMIT 1)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers","sales.Targets"]}
> query
| -- omni:span-query
| SELECT ARRAY(SELECT OrderId FROM Orders o WHERE o.CustomerId = c.CustomerId) AS ids FROM Customers c
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["ids"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT Name, ARRAY(SELECT AS STRUCT o.OrderId, o.Amount FROM Orders o WHERE o.CustomerId = c.CustomerId) AS orders FROM Customers c
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","orders"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| WITH big AS (SELECT CustomerId, Amount FROM Orders WHERE Amount > 100) SELECT * FROM big
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Amount"],"tables":["Orders"]}
> query
| -- omni:span-query
| WITH big AS (SELECT CustomerId, Amount FROM Orders) SELECT c.Name, b.Amount FROM Customers c JOIN big b USING (CustomerId)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Amount"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| WITH a AS (SELECT Name FROM Customers), b AS (SELECT Name FROM a) SELECT Name FROM b
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers"]}
> query
| -- omni:span-query
| WITH Orders AS (SELECT Name AS Status FROM Customers) SELECT Status FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Status"],"tables":["Customers"]}
> query
| -- omni:span-query
| WITH x AS (SELECT Region FROM Customers) SELECT Region FROM x UNION ALL SELECT Region FROM sales.Targets
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Region"],"tables":["Customers","sales.Targets"]}
> query
| -- omni:span-query
| WITH RECURSIVE r AS (SELECT EmpId, ManagerId FROM Employees WHERE ManagerId IS NULL UNION ALL SELECT e.EmpId, e.ManagerId FROM Employees e JOIN r ON e.ManagerId = r.EmpId) SELECT EmpId FROM r
< result
< column "verdict" "JSON"
< row
: string
| {"error":"InvalidArgument","message":"RECURSIVE is not supported in the WITH clause [at 1:1] WITH RECURSIVE r AS (SELECT EmpId, ManagerId FROM Employees WHERE ManagerId I... ^"}
> query
| -- omni:span-query
| WITH RECURSIVE n AS (SELECT 1 AS i UNION ALL SELECT i + 1 FROM n WHERE i < 3) SELECT i FROM n
< result
< column "verdict" "JSON"
< row
: string
| {"error":"InvalidArgument","message":"RECURSIVE is not supported in the WITH clause [at 1:1] WITH RECURSIVE n AS (SELECT 1 AS i UNION ALL SELECT i + 1 FROM n WHERE i \u003c 3)... ^"}
> query
| -- omni:span-query
| SELECT Name FROM Customers UNION ALL SELECT Name FROM Products
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers","Products"]}
> query
| -- omni:span-query
| SELECT Name FROM Customers UNION DISTINCT SELECT Name FROM Employees
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers","Employees"]}
> query
| -- omni:span-query
| SELECT CustomerId FROM Customers INTERSECT DISTINCT SELECT CustomerId FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT CustomerId FROM Customers EXCEPT DISTINCT SELECT CustomerId FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT Name, Email FROM Customers UNION ALL SELECT Status, CAST(Amount AS STRING) FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Email"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT * FROM Products UNION ALL SELECT ProductId, Name, Price FROM Products
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["ProductId","Name","Price"],"tables":["Products"]}
> query
| -- omni:span-query
| (SELECT Name FROM Customers) UNION ALL (SELECT Name FROM Products ORDER BY Name LIMIT 1)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers","Products"]}
> query
| -- omni:span-query
| SELECT * FROM CustomerOrders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name","OrderId","Amount"],"tables":["Customers","Orders","CustomerOrders"]}
> query
| -- omni:span-query
| SELECT Name, Amount FROM CustomerOrders WHERE OrderId > 1
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name","Amount"],"tables":["Customers","Orders","CustomerOrders"]}
> query
| -- omni:span-query
| SELECT * FROM BigOrders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","Name","Amount"],"tables":["Customers","Orders","BigOrders","CustomerOrders"]}
> query
| -- omni:span-query
| SELECT b.Name FROM BigOrders b JOIN Customers c ON b.CustomerId = c.CustomerId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers","Orders","BigOrders","CustomerOrders"]}
> query
| -- omni:span-query
| SELECT * FROM sales.RegionGoals
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Region","Goal"],"tables":["Customers","sales.Targets","sales.RegionGoals"]}
> query
| -- omni:span-query
| SELECT g.Goal, c.Name FROM sales.RegionGoals g JOIN Customers c ON c.Region = g.Region
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Goal","Name"],"tables":["Customers","sales.Targets","sales.RegionGoals"]}
> query
| -- omni:span-query
| SELECT COUNT(*) FROM CustomerOrders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":[""],"tables":["Customers","Orders","CustomerOrders"]}
> query
| -- omni:span-query
| SELECT Region FROM sales.Targets
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Region"],"tables":["sales.Targets"]}
> query
| -- omni:span-query
| SELECT `Id`, `CamelCol` FROM `MixedCase`
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Id","CamelCol"],"tables":["MixedCase"]}
> query
| -- omni:span-query
| SELECT id, camelcol, LOWER_COL FROM mixedcase
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["id","camelcol","LOWER_COL"],"tables":["MixedCase"]}
> query
| -- omni:span-query
| SELECT m.camelCol FROM MIXEDCASE m
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["camelCol"],"tables":["MixedCase"]}
> query
| -- omni:span-query
| SELECT * FROM MixedCase
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Id","CamelCol","lower_col"],"tables":["MixedCase"]}
> query
| -- omni:span-query
| SELECT `customers`.`name` FROM `customers`
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["name"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT customers.name FROM customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["name"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT Orders.Name FROM Customers AS Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT Customers.Amount FROM Orders AS Customers JOIN Customers AS c ON Customers.CustomerId = c.CustomerId
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Amount"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT Products.Name FROM (SELECT Name FROM Employees) AS Products
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Employees"]}
> query
| -- omni:span-query
| SELECT Name, ROW_NUMBER() OVER (ORDER BY CustomerId) AS rn FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"error":"InvalidArgument","message":"Analytic functions not supported [at 1:14] SELECT Name, ROW_NUMBER() OVER (ORDER BY CustomerId) AS rn FROM Customers              ^"}
> query
| -- omni:span-query
| SELECT OrderId, SUM(Amount) OVER (PARTITION BY CustomerId) FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"error":"InvalidArgument","message":"Analytic functions not supported [at 1:17] SELECT OrderId, SUM(Amount) OVER (PARTITION BY CustomerId) FROM Orders                 ^"}
> query
| -- omni:span-query
| SELECT Name, RANK() OVER (PARTITION BY Region ORDER BY Name) FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"error":"InvalidArgument","message":"Analytic functions not supported [at 1:14] SELECT Name, RANK() OVER (PARTITION BY Region ORDER BY Name) FROM Customers              ^"}
> query
| -- omni:span-query
| SELECT Region, COUNT(*) AS n FROM Customers GROUP BY Region
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Region","n"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT CustomerId, SUM(Amount), MAX(Status) FROM Orders GROUP BY CustomerId HAVING SUM(Amount) > 10
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId","",""],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT ARRAY_AGG(Name) FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":[""],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT STRING_AGG(Status, ',') AS statuses FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["statuses"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT COUNT(DISTINCT Region) FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":[""],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT ANY_VALUE(Email) AS e FROM Customers GROUP BY Region
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["e"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT item FROM Orders, UNNEST(Items) AS item
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["item"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT o.OrderId, item FROM Orders o CROSS JOIN UNNEST(o.Items) AS item
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["OrderId","item"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT item, off FROM Orders o, UNNEST(o.Items) AS item WITH OFFSET AS off
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["item","off"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT x FROM UNNEST([1, 2, 3]) AS x
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["x"]}
> query
| -- omni:span-query
| SELECT * FROM UNNEST(['a', 'b'])
< result
< column "verdict" "JSON"
< row
: string
| {"columns":[""]}
> query
| -- omni:span-query
| SELECT o.OrderId, (SELECT COUNT(*) FROM UNNEST(o.Items)) AS n FROM Orders o
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["OrderId","n"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT Items[SAFE_OFFSET(0)] AS first FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["first"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT ARRAY_LENGTH(Items) FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"columns":[""],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT s.a, s.b FROM (SELECT STRUCT(CustomerId AS a, Name AS b) AS s FROM Customers)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["a","b"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT STRUCT(CustomerId, Name).Name AS n FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["n"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT t.s.x FROM (SELECT STRUCT(Amount AS x) AS s FROM Orders) AS t
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["x"],"tables":["Orders"]}
> query
| -- omni:span-query
| SELECT v.CustomerId FROM (SELECT AS STRUCT CustomerId, Name FROM Customers) AS v
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["CustomerId"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT AS STRUCT CustomerId, Name FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"error":"Unimplemented","message":"Unsupported query shape: A struct value cannot be returned as a column value. Rewrite the query to flatten the struct fields in the result."}
> query
| -- omni:span-query
| SELECT AS VALUE STRUCT(CustomerId AS id) FROM Customers
< result
< column "verdict" "JSON"
< row
: string
| {"error":"Unimplemented","message":"Unsupported query shape: A struct value cannot be returned as a column value. Rewrite the query to flatten the struct fields in the result."}
> query
| -- omni:span-query
| SELECT e.OrderId FROM Customers c, UNNEST(ARRAY(SELECT AS STRUCT o.OrderId FROM Orders o WHERE o.CustomerId = c.CustomerId)) AS e
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["OrderId"],"tables":["Customers","Orders"]}
> query
| -- omni:span-query
| SELECT Name FROM Customers TABLESAMPLE BERNOULLI (10 PERCENT)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Name"],"tables":["Customers"]}
> query
| -- omni:span-query
| SELECT o.Amount FROM Orders o TABLESAMPLE RESERVOIR (5 ROWS)
< result
< column "verdict" "JSON"
< row
: string
| {"columns":["Amount"],"tables":["Orders"]}
