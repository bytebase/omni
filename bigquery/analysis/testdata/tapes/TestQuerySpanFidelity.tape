-- pin: span gcp:bigquery:bytebase-dev:US sha256:fef5b54092b50edbb5d7c363182e0cd578a2147d527595e968a26f5131429231 2026-11-26

== TestQuerySpanFidelity span#1
> raw
| {"op":"fixture","sql":"CREATE TABLE omni_span.Customers (CustomerId INT64 NOT NULL, Name STRING, Email STRING, Region STRING);\nCREATE TABLE omni_span.Orders (CustomerId INT64 NOT NULL, OrderId INT64 NOT NULL, Amount FLOAT64, Status STRING, Items ARRAY\u003cSTRING\u003e, Ship STRUCT\u003cCity STRING, Zip STRING\u003e, Placed DATE);\nCREATE TABLE omni_span.LineItems (CustomerId INT64, OrderId INT64, LineId INT64, ProductId INT64, Qty INT64);\nCREATE TABLE omni_span.Products (ProductId INT64, Name STRING, Price NUMERIC, Tags ARRAY\u003cSTRUCT\u003ck STRING, v STRING\u003e\u003e);\nCREATE TABLE omni_span.Employees (EmpId INT64, ManagerId INT64, Name STRING, Salary INT64);\nCREATE TABLE omni_span.MixedCase (Id INT64, CamelCol STRING, lower_col STRING);\nCREATE TABLE omni_span.Targets (Region STRING, Goal INT64);\nCREATE TABLE omni_span.events_2024 (EventId INT64, Kind STRING, Payload STRING);\nCREATE TABLE omni_span.events_2025 (EventId INT64, Kind STRING, Payload STRING);\nCREATE TABLE omni_span.Sales (Region STRING, Quarter STRING, Revenue INT64);\nCREATE VIEW omni_span.CustomerOrders AS SELECT c.CustomerId, c.Name, o.OrderId, o.Amount FROM omni_span.Customers AS c JOIN omni_span.Orders AS o ON c.CustomerId = o.CustomerId;\nCREATE VIEW omni_span.BigOrders AS SELECT co.CustomerId, co.Name, co.Amount FROM omni_span.CustomerOrders AS co WHERE co.Amount \u003e 100;\nCREATE VIEW omni_span.RegionGoals AS SELECT c.Region, t.Goal FROM omni_span.Customers AS c JOIN omni_span.Targets AS t ON c.Region = t.Region;\nCREATE TABLE FUNCTION omni_span.OrdersOf(cid INT64) AS SELECT OrderId, Amount FROM omni_span.Orders WHERE CustomerId = cid;"}
< raw
| {"relations":[{"name":"BigOrders","view":true,"definition":"SELECT co.CustomerId, co.Name, co.Amount FROM omni_span.CustomerOrders AS co WHERE co.Amount \u003e 100","columns":["CustomerId","Name","Amount"]},{"name":"CustomerOrders","view":true,"definition":"SELECT c.CustomerId, c.Name, o.OrderId, o.Amount FROM omni_span.Customers AS c JOIN omni_span.Orders AS o ON c.CustomerId = o.CustomerId","columns":["CustomerId","Name","OrderId","Amount"]},{"name":"Customers","columns":["CustomerId","Name","Email","Region"]},{"name":"Employees","columns":["EmpId","ManagerId","Name","Salary"]},{"name":"LineItems","columns":["CustomerId","OrderId","LineId","ProductId","Qty"]},{"name":"MixedCase","columns":["Id","CamelCol","lower_col"]},{"name":"Orders","columns":["CustomerId","OrderId","Amount","Status","Items","Ship","Placed"]},{"name":"Products","columns":["ProductId","Name","Price","Tags"]},{"name":"RegionGoals","view":true,"definition":"SELECT c.Region, t.Goal FROM omni_span.Customers AS c JOIN omni_span.Targets AS t ON c.Region = t.Region","columns":["Region","Goal"]},{"name":"Sales","columns":["Region","Quarter","Revenue"]},{"name":"Targets","columns":["Region","Goal"]},{"name":"events_2024","columns":["EventId","Kind","Payload"]},{"name":"events_2025","columns":["EventId","Kind","Payload"]}]}
> raw
| {"op":"query","sql":"SELECT CustomerId, Name FROM Customers"}
< raw
| {"columns":["CustomerId","Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT c.CustomerId AS id, c.Email FROM Customers AS c"}
< raw
| {"columns":["id","Email"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT c.CustomerId FROM Customers c"}
< raw
| {"columns":["CustomerId"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Amount * 2, UPPER(Status) FROM Orders"}
< raw
| {"columns":["f0_","f1_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT CAST(Amount AS STRING) AS a, SAFE_CAST(Status AS INT64) FROM Orders"}
< raw
| {"columns":["a","f0_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT 1, 'a' AS b"}
< raw
| {"columns":["f0_","b"]}
> raw
| {"op":"query","sql":"SELECT COUNT(*) FROM Orders"}
< raw
| {"columns":["f0_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT DISTINCT Region FROM Customers"}
< raw
| {"columns":["Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers ORDER BY Region LIMIT 5"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT OrderId FROM Orders LIMIT 2 OFFSET 1"}
< raw
| {"columns":["OrderId"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT CONCAT(Name, Email) AS contact FROM Customers"}
< raw
| {"columns":["contact"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT IF(Amount \u003e 10, Status, NULL) AS s FROM Orders"}
< raw
| {"columns":["s"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT CASE WHEN Amount \u003e 10 THEN 'big' ELSE Status END FROM Orders"}
< raw
| {"columns":["f0_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT COALESCE(Email, Name) FROM Customers"}
< raw
| {"columns":["f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT (Name), ((c.Email)) FROM Customers c"}
< raw
| {"columns":["Name","Email"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name AS x, Email AS x FROM Customers"}
< raw
| {"columns":["x","x_1"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name, Name FROM Customers"}
< raw
| {"columns":["Name","Name_1"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT 1 AS f0_, 2"}
< raw
| {"columns":["f0_","f1_"]}
> raw
| {"op":"query","sql":"SELECT CustomerId + OrderId AS k, Amount FROM Orders WHERE Status = 'open'"}
< raw
| {"columns":["k","Amount"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT EXTRACT(YEAR FROM Placed) AS y, FORMAT_DATE('%Y', Placed) FROM Orders"}
< raw
| {"columns":["y","f0_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT CURRENT_DATE() AS today, Name FROM Customers"}
< raw
| {"columns":["today","Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM omni_span.Customers"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT c.Name FROM omni_span.Customers AS c"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM `omni_span.Customers`"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM `omni_span`.`Customers`"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM `bytebase-dev.omni_span.Customers`"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM `bytebase-dev`.omni_span.Customers"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Customers.Name FROM Customers"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT omni_span.Customers.Name FROM omni_span.Customers"}
< raw
| {"error":"invalidQuery","message":"Unrecognized name: omni_span at [1:8]"}
> raw
| {"op":"query","sql":"SELECT c.Name, o.Amount FROM omni_span.Customers c JOIN Orders o USING (CustomerId)"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM Orders"}
< raw
| {"columns":["CustomerId","OrderId","Amount","Status","Items","Ship","Placed"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT o.* FROM Orders o"}
< raw
| {"columns":["CustomerId","OrderId","Amount","Status","Items","Ship","Placed"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT *, Amount FROM Orders"}
< raw
| {"columns":["CustomerId","OrderId","Amount","Status","Items","Ship","Placed","Amount_1"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT c.*, o.Amount FROM Customers c JOIN Orders o ON o.CustomerId = c.CustomerId"}
< raw
| {"columns":["CustomerId","Name","Email","Region","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM (SELECT * FROM Customers) c"}
< raw
| {"columns":["CustomerId","Name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT * EXCEPT (Email) FROM Customers"}
< raw
| {"columns":["CustomerId","Name","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT * REPLACE (UPPER(Name) AS Name) FROM Customers"}
< raw
| {"columns":["CustomerId","Name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT c.* EXCEPT (Region), o.OrderId FROM Customers c JOIN Orders o USING (CustomerId)"}
< raw
| {"columns":["CustomerId","Name","Email","OrderId"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * EXCEPT (Name) FROM (SELECT (Name), Email FROM Customers)"}
< raw
| {"columns":["Email"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT * EXCEPT (Items, Ship) REPLACE (Amount * 2 AS Amount) FROM Orders"}
< raw
| {"columns":["CustomerId","OrderId","Amount","Status","Placed"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM Customers, Targets"}
< raw
| {"columns":["CustomerId","Name","Email","Region","Region_1","Goal"],"tables":["omni_span.Customers","omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT * FROM Customers c JOIN Targets t USING (Region)"}
< raw
| {"columns":["Region","CustomerId","Name","Email","Goal"],"tables":["omni_span.Customers","omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT t.* FROM (SELECT Name, Email FROM Customers) t"}
< raw
| {"columns":["Name","Email"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT c.Name, o.Amount FROM Customers c JOIN Orders o ON c.CustomerId = o.CustomerId"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT c.Name, o.Amount FROM Customers c INNER JOIN Orders o ON c.CustomerId = o.CustomerId"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT c.Name, o.Amount FROM Customers c LEFT JOIN Orders o ON c.CustomerId = o.CustomerId"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT c.Name, o.Amount FROM Customers c LEFT OUTER JOIN Orders o ON c.CustomerId = o.CustomerId"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT c.Name, o.Amount FROM Customers c RIGHT JOIN Orders o ON c.CustomerId = o.CustomerId"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT c.Name, o.Amount FROM Customers c FULL JOIN Orders o ON c.CustomerId = o.CustomerId"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT c.Name, p.Name FROM Customers c CROSS JOIN Products p"}
< raw
| {"columns":["Name","Name_1"],"tables":["omni_span.Customers","omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT c.Name, p.Name FROM Customers c, Products p"}
< raw
| {"columns":["Name","Name_1"],"tables":["omni_span.Customers","omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT CustomerId, Name, Amount FROM Customers JOIN Orders USING (CustomerId)"}
< raw
| {"columns":["CustomerId","Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM Customers JOIN Orders USING (CustomerId)"}
< raw
| {"columns":["CustomerId","Name","Email","Region","OrderId","Amount","Status","Items","Ship","Placed"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM Orders JOIN LineItems USING (CustomerId, OrderId)"}
< raw
| {"columns":["CustomerId","OrderId","Amount","Status","Items","Ship","Placed","LineId","ProductId","Qty"],"tables":["omni_span.LineItems","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT CustomerId FROM Customers FULL JOIN Orders USING (CustomerId)"}
< raw
| {"columns":["CustomerId"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT e.Name, m.Name AS manager FROM Employees e LEFT JOIN Employees m ON e.ManagerId = m.EmpId"}
< raw
| {"columns":["Name","manager"],"tables":["omni_span.Employees"]}
> raw
| {"op":"query","sql":"SELECT l.Qty, p.Price, c.Region FROM LineItems l JOIN Products p ON l.ProductId = p.ProductId JOIN Customers c ON c.CustomerId = l.CustomerId"}
< raw
| {"columns":["Qty","Price","Region"],"tables":["omni_span.Customers","omni_span.LineItems","omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT o.OrderId FROM Orders o JOIN LineItems l ON o.OrderId = l.OrderId AND o.CustomerId = l.CustomerId"}
< raw
| {"columns":["OrderId"],"tables":["omni_span.LineItems","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT c.Name FROM Customers c JOIN (Orders o JOIN LineItems l USING (OrderId)) ON c.CustomerId = o.CustomerId"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.LineItems","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT a.Name, b.Name FROM Customers a JOIN Customers b ON a.Region = b.Region"}
< raw
| {"columns":["Name","Name_1"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name, Goal FROM Customers JOIN Targets USING (Region)"}
< raw
| {"columns":["Name","Goal"],"tables":["omni_span.Customers","omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT c.Name, t.Goal FROM Customers c LEFT JOIN Targets t ON c.Region = t.Region WHERE t.Goal IS NULL"}
< raw
| {"columns":["Name","Goal"],"tables":["omni_span.Customers","omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT Name, (SELECT MAX(Amount) FROM Orders) AS top FROM Customers"}
< raw
| {"columns":["Name","top"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name, (SELECT MAX(o.Amount) FROM Orders o WHERE o.CustomerId = c.CustomerId) AS top FROM Customers c"}
< raw
| {"columns":["Name","top"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers WHERE CustomerId IN (SELECT CustomerId FROM Orders)"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers c WHERE EXISTS (SELECT 1 FROM Orders o WHERE o.CustomerId = c.CustomerId)"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers c WHERE NOT EXISTS (SELECT 1 FROM Orders o WHERE o.CustomerId = c.CustomerId)"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT t.total FROM (SELECT CustomerId, SUM(Amount) AS total FROM Orders GROUP BY CustomerId) t"}
< raw
| {"columns":["total"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT x FROM (SELECT Name AS x FROM Customers UNION ALL SELECT Status FROM Orders)"}
< raw
| {"columns":["x"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers WHERE Region = (SELECT Region FROM Targets ORDER BY Goal DESC LIMIT 1)"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT ARRAY(SELECT OrderId FROM Orders o WHERE o.CustomerId = c.CustomerId) AS ids FROM Customers c"}
< raw
| {"columns":["ids"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name, ARRAY(SELECT AS STRUCT o.OrderId, o.Amount FROM Orders o WHERE o.CustomerId = c.CustomerId) AS orders FROM Customers c"}
< raw
| {"columns":["Name","orders"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT (SELECT COUNT(*) FROM LineItems) AS n"}
< raw
| {"columns":["n"],"tables":["omni_span.LineItems"]}
> raw
| {"op":"query","sql":"SELECT * FROM (SELECT * FROM (SELECT Name, Region FROM Customers))"}
< raw
| {"columns":["Name","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers WHERE Region IN (SELECT Region FROM Targets WHERE Goal \u003e (SELECT AVG(Goal) FROM Targets))"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers WHERE CustomerId = ANY (SELECT CustomerId FROM Orders)"}
< raw
| {"error":"invalidQuery","message":"Quantified comparisons (= ANY) are not supported at [1:45]"}
> raw
| {"op":"query","sql":"SELECT c.Name, x.n FROM Customers c, (SELECT COUNT(*) AS n FROM Orders) x"}
< raw
| {"columns":["Name","n"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT EXISTS(SELECT 1 FROM Orders) AS has_orders"}
< raw
| {"columns":["has_orders"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"WITH big AS (SELECT CustomerId, Amount FROM Orders WHERE Amount \u003e 100) SELECT * FROM big"}
< raw
| {"columns":["CustomerId","Amount"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"WITH big AS (SELECT CustomerId, Amount FROM Orders) SELECT c.Name, b.Amount FROM Customers c JOIN big b USING (CustomerId)"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"WITH a AS (SELECT Name FROM Customers), b AS (SELECT Name FROM a) SELECT Name FROM b"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"WITH Orders AS (SELECT Name AS Status FROM Customers) SELECT Status FROM Orders"}
< raw
| {"columns":["Status"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"WITH x AS (SELECT Region FROM Customers) SELECT Region FROM x UNION ALL SELECT Region FROM Targets"}
< raw
| {"columns":["Region"],"tables":["omni_span.Customers","omni_span.Targets"]}
> raw
| {"op":"query","sql":"WITH RECURSIVE r AS (SELECT EmpId, ManagerId FROM Employees WHERE ManagerId IS NULL UNION ALL SELECT e.EmpId, e.ManagerId FROM Employees e JOIN r ON e.ManagerId = r.EmpId) SELECT EmpId FROM r"}
< raw
| {"columns":["EmpId"],"tables":["omni_span.Employees"]}
> raw
| {"op":"query","sql":"WITH RECURSIVE n AS (SELECT 1 AS i UNION ALL SELECT i + 1 FROM n WHERE i \u003c 3) SELECT i FROM n"}
< raw
| {"columns":["i"]}
> raw
| {"op":"query","sql":"WITH t AS (SELECT * FROM Products) SELECT t.Name, t.Price FROM t"}
< raw
| {"columns":["Name","Price"],"tables":["omni_span.Products"]}
> raw
| {"op":"query","sql":"WITH t (a, b) AS (SELECT Name, Email FROM Customers) SELECT a FROM t"}
< raw
| {"error":"invalidQuery","message":"Syntax error: Expected \")\" but got identifier \"a\" at [1:9]"}
> raw
| {"op":"query","sql":"WITH t AS (SELECT 1 AS one) SELECT one, Name FROM t, Customers"}
< raw
| {"columns":["one","Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"WITH s AS (SELECT CustomerId, SUM(Amount) AS total FROM Orders GROUP BY CustomerId) SELECT c.Name, s.total FROM s JOIN Customers c USING (CustomerId)"}
< raw
| {"columns":["Name","total"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"WITH a AS (SELECT * FROM Customers) SELECT * FROM a WHERE Region = 'eu'"}
< raw
| {"columns":["CustomerId","Name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"WITH Customers AS (SELECT * FROM omni_span.Customers WHERE Region = 'eu') SELECT Name FROM Customers"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers UNION ALL SELECT Name FROM Products"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers UNION DISTINCT SELECT Name FROM Employees"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Employees"]}
> raw
| {"op":"query","sql":"SELECT CustomerId FROM Customers INTERSECT DISTINCT SELECT CustomerId FROM Orders"}
< raw
| {"columns":["CustomerId"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT CustomerId FROM Customers EXCEPT DISTINCT SELECT CustomerId FROM Orders"}
< raw
| {"columns":["CustomerId"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name, Email FROM Customers UNION ALL SELECT Status, CAST(Amount AS STRING) FROM Orders"}
< raw
| {"columns":["Name","Email"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT ProductId, Name FROM Products UNION ALL SELECT ProductId, Name FROM Products"}
< raw
| {"columns":["ProductId","Name"],"tables":["omni_span.Products"]}
> raw
| {"op":"query","sql":"(SELECT Name FROM Customers) UNION ALL (SELECT Name FROM Products ORDER BY Name LIMIT 1)"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers UNION ALL SELECT Name FROM Products UNION ALL SELECT Name FROM Employees"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Employees","omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT * FROM Targets UNION ALL SELECT Region, Goal FROM Targets"}
< raw
| {"columns":["Region","Goal"],"tables":["omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT Name, Region FROM Customers UNION ALL BY NAME SELECT Region, Name FROM Customers"}
< raw
| {"columns":["Name","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers UNION ALL CORRESPONDING SELECT Name FROM Employees"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Employees"]}
> raw
| {"op":"query","sql":"SELECT 1 AS a UNION ALL SELECT 2"}
< raw
| {"columns":["a"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers FULL UNION ALL BY NAME SELECT Name, Salary FROM Employees"}
< raw
| {"columns":["Name","Salary"],"tables":["omni_span.Customers","omni_span.Employees"]}
> raw
| {"op":"query","sql":"SELECT Name, Region FROM Customers LEFT UNION ALL BY NAME SELECT Name FROM Employees"}
< raw
| {"columns":["Name","Region"],"tables":["omni_span.Customers","omni_span.Employees"]}
> raw
| {"op":"query","sql":"SELECT * FROM CustomerOrders"}
< raw
| {"columns":["CustomerId","Name","OrderId","Amount"],"tables":["omni_span.CustomerOrders","omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name, Amount FROM CustomerOrders WHERE OrderId \u003e 1"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.CustomerOrders","omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM BigOrders"}
< raw
| {"columns":["CustomerId","Name","Amount"],"tables":["omni_span.BigOrders","omni_span.CustomerOrders","omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT b.Name FROM BigOrders b JOIN Customers c ON b.CustomerId = c.CustomerId"}
< raw
| {"columns":["Name"],"tables":["omni_span.BigOrders","omni_span.CustomerOrders","omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM RegionGoals"}
< raw
| {"columns":["Region","Goal"],"tables":["omni_span.Customers","omni_span.RegionGoals","omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT g.Goal, c.Name FROM RegionGoals g JOIN Customers c ON c.Region = g.Region"}
< raw
| {"columns":["Goal","Name"],"tables":["omni_span.Customers","omni_span.RegionGoals","omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT COUNT(*) FROM CustomerOrders"}
< raw
| {"columns":["f0_"],"tables":["omni_span.CustomerOrders","omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name FROM omni_span.BigOrders"}
< raw
| {"columns":["Name"],"tables":["omni_span.BigOrders","omni_span.CustomerOrders","omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"WITH v AS (SELECT * FROM CustomerOrders) SELECT Name FROM v"}
< raw
| {"columns":["Name"],"tables":["omni_span.CustomerOrders","omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM OrdersOf(1)"}
< raw
| {"error":"invalidQuery","message":"Table-valued function not found: OrdersOf at [1:15]"}
> raw
| {"op":"query","sql":"SELECT OrderId FROM omni_span.OrdersOf(2)"}
< raw
| {"columns":["OrderId"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT c.Name, f.Amount FROM Customers c, OrdersOf(c.CustomerId) f"}
< raw
| {"error":"invalidQuery","message":"Table-valued function not found: OrdersOf at [1:43]"}
> raw
| {"op":"query","sql":"SELECT EventId FROM `events_*`"}
< raw
| {"columns":["EventId"],"tables":["omni_span.events_*"]}
> raw
| {"op":"query","sql":"SELECT * FROM `events_*`"}
< raw
| {"columns":["EventId","Kind","Payload"],"tables":["omni_span.events_*"]}
> raw
| {"op":"query","sql":"SELECT Kind FROM `omni_span.events_*` WHERE _TABLE_SUFFIX = '2024'"}
< raw
| {"columns":["Kind"],"tables":["omni_span.events_*"]}
> raw
| {"op":"query","sql":"SELECT _TABLE_SUFFIX AS year, Kind FROM `events_*`"}
< raw
| {"columns":["year","Kind"],"tables":["omni_span.events_*"]}
> raw
| {"op":"query","sql":"SELECT e.Payload FROM `events_20*` AS e WHERE e._TABLE_SUFFIX BETWEEN '24' AND '25'"}
< raw
| {"columns":["Payload"],"tables":["omni_span.events_20*"]}
> raw
| {"op":"query","sql":"SELECT COUNT(*) FROM `events_*` WHERE Kind = 'click'"}
< raw
| {"columns":["f0_"],"tables":["omni_span.events_*"]}
> raw
| {"op":"query","sql":"SELECT * FROM events_2024 UNION ALL SELECT * FROM events_2025"}
< raw
| {"columns":["EventId","Kind","Payload"],"tables":["omni_span.events_2024","omni_span.events_2025"]}
> raw
| {"op":"query","sql":"SELECT EventId FROM `events_2024`"}
< raw
| {"columns":["EventId"],"tables":["omni_span.events_2024"]}
> raw
| {"op":"query","sql":"SELECT `Id`, `CamelCol` FROM `MixedCase`"}
< raw
| {"columns":["Id","CamelCol"],"tables":["omni_span.MixedCase"]}
> raw
| {"op":"query","sql":"SELECT id, camelcol, LOWER_COL FROM MixedCase"}
< raw
| {"columns":["id","camelcol","LOWER_COL"],"tables":["omni_span.MixedCase"]}
> raw
| {"op":"query","sql":"SELECT m.camelCol FROM MixedCase m"}
< raw
| {"columns":["camelCol"],"tables":["omni_span.MixedCase"]}
> raw
| {"op":"query","sql":"SELECT * FROM MixedCase"}
< raw
| {"columns":["Id","CamelCol","lower_col"],"tables":["omni_span.MixedCase"]}
> raw
| {"op":"query","sql":"SELECT `Customers`.`Name` FROM `Customers`"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT name FROM customers"}
< raw
| {"error":"notFound","message":"Not found: Table bytebase-dev:omni_span.customers was not found in location US"}
> raw
| {"op":"query","sql":"SELECT NAME FROM Customers"}
< raw
| {"columns":["NAME"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT c.NAME, C.email FROM Customers c"}
< raw
| {"columns":["NAME","email"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Orders.Name FROM Customers AS Orders"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Customers.Amount FROM Orders AS Customers JOIN Customers AS c ON Customers.CustomerId = c.CustomerId"}
< raw
| {"columns":["Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Products.Name FROM (SELECT Name FROM Employees) AS Products"}
< raw
| {"columns":["Name"],"tables":["omni_span.Employees"]}
> raw
| {"op":"query","sql":"SELECT Name, ROW_NUMBER() OVER (ORDER BY CustomerId) AS rn FROM Customers"}
< raw
| {"columns":["Name","rn"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT OrderId, SUM(Amount) OVER (PARTITION BY CustomerId) FROM Orders"}
< raw
| {"columns":["OrderId","f0_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name, RANK() OVER (PARTITION BY Region ORDER BY Name) FROM Customers"}
< raw
| {"columns":["Name","f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name, LAG(Name) OVER w AS prev FROM Customers WINDOW w AS (ORDER BY CustomerId)"}
< raw
| {"columns":["Name","prev"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers QUALIFY ROW_NUMBER() OVER (PARTITION BY Region ORDER BY Name) = 1"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT OrderId, Amount FROM Orders WHERE TRUE QUALIFY RANK() OVER (ORDER BY Amount DESC) \u003c= 3"}
< raw
| {"columns":["OrderId","Amount"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT FIRST_VALUE(Name) OVER (PARTITION BY Region ORDER BY CustomerId) AS first FROM Customers"}
< raw
| {"columns":["first"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Region, COUNT(*) AS n FROM Customers GROUP BY Region"}
< raw
| {"columns":["Region","n"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT CustomerId, SUM(Amount), MAX(Status) FROM Orders GROUP BY CustomerId HAVING SUM(Amount) \u003e 10"}
< raw
| {"columns":["CustomerId","f0_","f1_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT ARRAY_AGG(Name) FROM Customers"}
< raw
| {"columns":["f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT STRING_AGG(Status, ',') AS statuses FROM Orders"}
< raw
| {"columns":["statuses"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT COUNT(DISTINCT Region) FROM Customers"}
< raw
| {"columns":["f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT ANY_VALUE(Email) AS e FROM Customers GROUP BY Region"}
< raw
| {"columns":["e"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Region, COUNT(*) FROM Customers GROUP BY 1"}
< raw
| {"columns":["Region","f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Region, Name, COUNT(*) FROM Customers GROUP BY ROLLUP (Region, Name)"}
< raw
| {"columns":["Region","Name","f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Region, COUNT(*) FROM Customers GROUP BY ALL"}
< raw
| {"columns":["Region","f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT ARRAY_AGG(STRUCT(OrderId, Amount) ORDER BY Amount LIMIT 3) AS top FROM Orders"}
< raw
| {"columns":["top"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT APPROX_COUNT_DISTINCT(Email) FROM Customers"}
< raw
| {"columns":["f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Region, LOGICAL_OR(Email IS NULL) FROM Customers GROUP BY Region"}
< raw
| {"columns":["Region","f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT AVG(Amount) AS avg_amount, MIN(Placed) FROM Orders"}
< raw
| {"columns":["avg_amount","f0_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT item FROM Orders, UNNEST(Items) AS item"}
< raw
| {"columns":["item"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT o.OrderId, item FROM Orders o CROSS JOIN UNNEST(o.Items) AS item"}
< raw
| {"columns":["OrderId","item"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT item, off FROM Orders o, UNNEST(o.Items) AS item WITH OFFSET AS off"}
< raw
| {"columns":["item","off"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT x FROM UNNEST([1, 2, 3]) AS x"}
< raw
| {"columns":["x"]}
> raw
| {"op":"query","sql":"SELECT * FROM UNNEST(['a', 'b'])"}
< raw
| {"columns":["f0_"]}
> raw
| {"op":"query","sql":"SELECT o.OrderId, (SELECT COUNT(*) FROM UNNEST(o.Items)) AS n FROM Orders o"}
< raw
| {"columns":["OrderId","n"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Items[SAFE_OFFSET(0)] AS first FROM Orders"}
< raw
| {"columns":["first"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT ARRAY_LENGTH(Items) FROM Orders"}
< raw
| {"columns":["f0_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT o.OrderId, i FROM Orders o, o.Items AS i"}
< raw
| {"columns":["OrderId","i"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT o.OrderId, i FROM Orders o LEFT JOIN UNNEST(o.Items) AS i"}
< raw
| {"columns":["OrderId","i"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT t.k, t.v FROM Products, UNNEST(Tags) AS t"}
< raw
| {"columns":["k","v"],"tables":["omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT p.Name, tag.k FROM Products p CROSS JOIN UNNEST(p.Tags) AS tag WHERE tag.v = 'x'"}
< raw
| {"columns":["Name","k"],"tables":["omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT * FROM Products, UNNEST(Tags)"}
< raw
| {"columns":["ProductId","Name","Price","Tags","k","v"],"tables":["omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT p.Name, (SELECT STRING_AGG(t.v) FROM UNNEST(p.Tags) t) AS vals FROM Products p"}
< raw
| {"columns":["Name","vals"],"tables":["omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT Items[OFFSET(0)], Tags[SAFE_ORDINAL(1)].k FROM Orders, Products"}
< raw
| {"columns":["f0_","k"],"tables":["omni_span.Orders","omni_span.Products"]}
> raw
| {"op":"query","sql":"SELECT ARRAY_CONCAT(Items, ['x']) AS more FROM Orders"}
< raw
| {"columns":["more"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM UNNEST([STRUCT(1 AS a, 'x' AS b), (2, 'y')])"}
< raw
| {"columns":["a","b"]}
> raw
| {"op":"query","sql":"SELECT a, b FROM UNNEST([STRUCT(1 AS a, 'x' AS b)])"}
< raw
| {"columns":["a","b"]}
> raw
| {"op":"query","sql":"SELECT x, y FROM UNNEST([1, 2]) AS x, UNNEST([3]) AS y"}
< raw
| {"columns":["x","y"]}
> raw
| {"op":"query","sql":"SELECT item FROM UNNEST((SELECT Items FROM Orders LIMIT 1)) AS item"}
< raw
| {"columns":["item"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Ship.City, Ship.Zip FROM Orders"}
< raw
| {"columns":["City","Zip"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT o.Ship.City AS city FROM Orders o"}
< raw
| {"columns":["city"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Ship FROM Orders"}
< raw
| {"columns":["Ship"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Ship.* FROM Orders"}
< raw
| {"columns":["City","Zip"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT o.Ship.*, o.OrderId FROM Orders o"}
< raw
| {"columns":["City","Zip","OrderId"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT s.a, s.b FROM (SELECT STRUCT(CustomerId AS a, Name AS b) AS s FROM Customers)"}
< raw
| {"columns":["a","b"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT STRUCT(CustomerId, Name).Name AS n FROM Customers"}
< raw
| {"columns":["n"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT t.s.x FROM (SELECT STRUCT(Amount AS x) AS s FROM Orders) AS t"}
< raw
| {"columns":["x"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT v.CustomerId FROM (SELECT AS STRUCT CustomerId, Name FROM Customers) AS v"}
< raw
| {"columns":["CustomerId"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT AS STRUCT CustomerId, Name FROM Customers"}
< raw
| {"columns":["CustomerId","Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT AS VALUE STRUCT(CustomerId AS id) FROM Customers"}
< raw
| {"columns":["id"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT AS VALUE Name FROM Customers"}
< raw
| {"error":"invalidQuery","message":"No support for value tables that are not of type struct. Type is STRING"}
> raw
| {"op":"query","sql":"SELECT e.OrderId FROM Customers c, UNNEST(ARRAY(SELECT AS STRUCT o.OrderId FROM Orders o WHERE o.CustomerId = c.CustomerId)) AS e"}
< raw
| {"columns":["OrderId"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT STRUCT(Name, Email) AS contact FROM Customers"}
< raw
| {"columns":["contact"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT [Name, Email] AS arr FROM Customers"}
< raw
| {"columns":["arr"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT (SELECT AS STRUCT Name, Email) AS s FROM Customers"}
< raw
| {"columns":["s"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT c FROM Customers c"}
< raw
| {"columns":["c"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT TO_JSON_STRING(c) FROM Customers c"}
< raw
| {"columns":["f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT ANY_VALUE(Ship).City FROM Orders"}
< raw
| {"columns":["City"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM Sales PIVOT (SUM(Revenue) FOR Quarter IN ('Q1', 'Q2'))"}
< raw
| {"columns":["Region","Q1","Q2"],"tables":["omni_span.Sales"]}
> raw
| {"op":"query","sql":"SELECT Region, Q1 FROM Sales PIVOT (SUM(Revenue) AS Q FOR Quarter IN ('Q1' AS Q1, 'Q2' AS Q2))"}
< raw
| {"error":"invalidQuery","message":"Unrecognized name: Q1 at [1:16]"}
> raw
| {"op":"query","sql":"SELECT * FROM (SELECT Region, Quarter, Revenue FROM Sales) PIVOT (MAX(Revenue) FOR Quarter IN ('Q1'))"}
< raw
| {"columns":["Region","Q1"],"tables":["omni_span.Sales"]}
> raw
| {"op":"query","sql":"SELECT * FROM Targets UNPIVOT (val FOR col IN (Goal))"}
< raw
| {"columns":["Region","val","col"],"tables":["omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT Region, val FROM (SELECT Region, Revenue, Revenue * 2 AS Double FROM Sales) UNPIVOT (val FOR kind IN (Revenue, Double))"}
< raw
| {"columns":["Region","val"],"tables":["omni_span.Sales"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers TABLESAMPLE SYSTEM (10 PERCENT)"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT o.Amount FROM Orders o TABLESAMPLE SYSTEM (5 PERCENT)"}
< raw
| {"columns":["Amount"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers FOR SYSTEM_TIME AS OF CURRENT_TIMESTAMP()"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT c.Name FROM Customers FOR SYSTEM TIME AS OF CURRENT_TIMESTAMP() AS c"}
< raw
| {"error":"invalidQuery","message":"Syntax error: Expected end of input but got keyword AS at [1:72]"}
> raw
| {"op":"query","sql":"FROM Customers"}
< raw
| {"columns":["CustomerId","Name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e SELECT Name"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e WHERE Region = 'eu' |\u003e SELECT Name, Email"}
< raw
| {"columns":["Name","Email"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Orders |\u003e AGGREGATE SUM(Amount) AS total GROUP BY CustomerId"}
< raw
| {"columns":["CustomerId","total"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"FROM Orders |\u003e EXTEND Amount * 2 AS double"}
< raw
| {"columns":["CustomerId","OrderId","Amount","Status","Items","Ship","Placed","double"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"FROM Orders |\u003e EXTEND Amount * 2 AS double |\u003e SELECT OrderId, double"}
< raw
| {"columns":["OrderId","double"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e SET Name = UPPER(Name)"}
< raw
| {"columns":["CustomerId","Name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e DROP Email"}
< raw
| {"columns":["CustomerId","Name","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e RENAME Name AS full_name"}
< raw
| {"columns":["CustomerId","full_name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e AS c |\u003e SELECT c.Name"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers c |\u003e JOIN Orders o ON c.CustomerId = o.CustomerId |\u003e SELECT c.Name, o.Amount"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e JOIN Orders USING (CustomerId) |\u003e SELECT Name, Amount"}
< raw
| {"columns":["Name","Amount"],"tables":["omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e ORDER BY Name |\u003e LIMIT 3"}
< raw
| {"columns":["CustomerId","Name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e SELECT Name |\u003e UNION ALL (SELECT Name FROM Products)"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers","omni_span.Products"]}
> raw
| {"op":"query","sql":"FROM Orders |\u003e AGGREGATE COUNT(*) AS n GROUP BY Status |\u003e WHERE n \u003e 1"}
< raw
| {"columns":["Status","n"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e LEFT JOIN Targets USING (Region) |\u003e SELECT Name, Goal"}
< raw
| {"columns":["Name","Goal"],"tables":["omni_span.Customers","omni_span.Targets"]}
> raw
| {"op":"query","sql":"FROM Orders |\u003e CROSS JOIN UNNEST(Items) AS item |\u003e SELECT OrderId, item"}
< raw
| {"columns":["OrderId","item"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"FROM Sales |\u003e PIVOT (SUM(Revenue) FOR Quarter IN ('Q1', 'Q2'))"}
< raw
| {"columns":["Region","Q1","Q2"],"tables":["omni_span.Sales"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e WINDOW ROW_NUMBER() OVER (ORDER BY CustomerId) AS rn"}
< raw
| {"columns":["CustomerId","Name","Email","Region","rn"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e EXTEND ROW_NUMBER() OVER (ORDER BY CustomerId) AS rn |\u003e SELECT Name, rn"}
< raw
| {"columns":["Name","rn"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e DISTINCT"}
< raw
| {"columns":["CustomerId","Name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e SELECT Region |\u003e DISTINCT"}
< raw
| {"columns":["Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"WITH t AS (SELECT * FROM Orders) FROM t |\u003e SELECT OrderId"}
< raw
| {"columns":["OrderId"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM Customers |\u003e WHERE Region = 'eu'"}
< raw
| {"columns":["CustomerId","Name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name, Email FROM Customers |\u003e SELECT Name"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e TABLESAMPLE SYSTEM (10 PERCENT) |\u003e SELECT Name"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM CustomerOrders |\u003e SELECT Name"}
< raw
| {"columns":["Name"],"tables":["omni_span.CustomerOrders","omni_span.Customers","omni_span.Orders"]}
> raw
| {"op":"query","sql":"FROM `events_*` |\u003e SELECT Kind"}
< raw
| {"columns":["Kind"],"tables":["omni_span.events_*"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e SELECT * EXCEPT (Email)"}
< raw
| {"columns":["CustomerId","Name","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e AGGREGATE COUNT(*) GROUP BY Region"}
< raw
| {"columns":["Region","f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"FROM Customers |\u003e CALL OrdersOf(1)"}
< raw
| {"error":"invalidQuery","message":"Table-valued function not found: OrdersOf at [1:24]"}
> raw
| {"op":"query","sql":"SELECT table_name FROM omni_span.INFORMATION_SCHEMA.TABLES"}
< raw
| {"columns":["table_name"],"tables":["omni_span.INFORMATION_SCHEMA.TABLES"]}
> raw
| {"op":"query","sql":"SELECT column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = 'Orders'"}
< raw
| {"columns":["column_name","data_type"],"tables":["omni_span.INFORMATION_SCHEMA.COLUMNS"]}
> raw
| {"op":"query","sql":"SELECT word FROM `bigquery-public-data.samples.shakespeare`"}
< raw
| {"columns":["word"],"tables":["samples.shakespeare"]}
> raw
| {"op":"query","sql":"SELECT word, corpus FROM `bigquery-public-data`.samples.shakespeare LIMIT 1"}
< raw
| {"columns":["word","corpus"],"tables":["samples.shakespeare"]}
> raw
| {"op":"query","sql":"SELECT s.word, c.Name FROM `bigquery-public-data.samples.shakespeare` s, Customers c"}
< raw
| {"columns":["word","Name"],"tables":["omni_span.Customers","samples.shakespeare"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers WHERE Name LIKE 'a%' AND Email IS NOT NULL"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT SAFE_DIVIDE(Amount, 2) AS half FROM Orders"}
< raw
| {"columns":["half"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT * FROM Customers WHERE CustomerId BETWEEN 1 AND 10"}
< raw
| {"columns":["CustomerId","Name","Email","Region"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers ORDER BY 1 DESC"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT IFNULL(Email, 'none') AS email FROM Customers"}
< raw
| {"columns":["email"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers WHERE STARTS_WITH(Email, 'x') LIMIT 1"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT GENERATE_UUID() AS id, Name FROM Customers"}
< raw
| {"columns":["id","Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT ARRAY(SELECT x FROM UNNEST(Items) AS x WHERE x != '') AS clean FROM Orders"}
< raw
| {"columns":["clean"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Name, (SELECT Goal FROM Targets t WHERE t.Region = c.Region) AS goal FROM Customers c"}
< raw
| {"columns":["Name","goal"],"tables":["omni_span.Customers","omni_span.Targets"]}
> raw
| {"op":"query","sql":"SELECT DATE_TRUNC(Placed, MONTH) AS month, SUM(Amount) FROM Orders GROUP BY month"}
< raw
| {"columns":["month","f0_"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT REGEXP_EXTRACT(Email, r'@(.*)') AS domain FROM Customers"}
< raw
| {"columns":["domain"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT Name || ' \u003c' || Email || '\u003e' AS addr FROM Customers"}
< raw
| {"columns":["addr"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT `Name` FROM Customers WHERE `Region` = 'eu'"}
< raw
| {"columns":["Name"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT 1 FROM Customers"}
< raw
| {"columns":["f0_"],"tables":["omni_span.Customers"]}
> raw
| {"op":"query","sql":"SELECT NULL AS n FROM Orders"}
< raw
| {"columns":["n"],"tables":["omni_span.Orders"]}
> raw
| {"op":"query","sql":"SELECT Nope FROM Customers"}
< raw
| {"error":"invalidQuery","message":"Unrecognized name: Nope; Did you mean Name? at [1:8]"}
> raw
| {"op":"query","sql":"SELECT Name FROM Missing"}
< raw
| {"error":"notFound","message":"Not found: Table bytebase-dev:omni_span.Missing was not found in location US"}
> raw
| {"op":"query","sql":"SELECT x.Name FROM Customers c"}
< raw
| {"error":"invalidQuery","message":"Unrecognized name: x at [1:8]"}
> raw
| {"op":"query","sql":"SELECT Name FROM Customers, Products"}
< raw
| {"error":"invalidQuery","message":"Column name Name is ambiguous at [1:8]"}
> raw
| {"op":"query","sql":"SELECT * FROM Customers JOIN Orders ON TRUE JOIN Products USING (Name)"}
< raw
| {"columns":["Name","CustomerId","Email","Region","CustomerId_1","OrderId","Amount","Status","Items","Ship","Placed","ProductId","Price","Tags"],"tables":["omni_span.Customers","omni_span.Orders","omni_span.Products"]}
