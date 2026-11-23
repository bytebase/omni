-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestQuerySpanFidelity default#1
> raw
| fixture
| CREATE DATABASE shop;
| CREATE DATABASE hr;
| CREATE TABLE shop.customers (id INT, name VARCHAR(20), region VARCHAR(10), email VARCHAR(40)) DUPLICATE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 1 PROPERTIES ("replication_num" = "1");
| CREATE TABLE shop.orders (id INT, customer_id INT, amount DECIMAL(10,2), status VARCHAR(10), created_at DATE) DUPLICATE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 1 PROPERTIES ("replication_num" = "1");
| CREATE TABLE shop.items (order_id INT, sku VARCHAR(10), qty INT) DUPLICATE KEY(order_id) DISTRIBUTED BY HASH(order_id) BUCKETS 1 PROPERTIES ("replication_num" = "1");
| CREATE TABLE shop.`MixedCase` (`Id` INT, `CamelCol` VARCHAR(10), lower_col VARCHAR(10)) DUPLICATE KEY(`Id`) DISTRIBUTED BY HASH(`Id`) BUCKETS 1 PROPERTIES ("replication_num" = "1");
| CREATE TABLE shop.mixedcase (id INT, note VARCHAR(10)) DUPLICATE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 1 PROPERTIES ("replication_num" = "1");
| CREATE TABLE hr.staff (id INT, name VARCHAR(20), salary INT, dept_id INT) DUPLICATE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 1 PROPERTIES ("replication_num" = "1");
| CREATE TABLE hr.depts (id INT, name VARCHAR(20)) DUPLICATE KEY(id) DISTRIBUTED BY HASH(id) BUCKETS 1 PROPERTIES ("replication_num" = "1");
| CREATE VIEW shop.big_orders AS SELECT id, amount * 2 AS doubled, status, customer_id FROM shop.orders WHERE amount > 5;
| CREATE VIEW shop.customer_totals AS SELECT customer_id, SUM(amount) AS total, COUNT(*) AS n FROM shop.orders GROUP BY customer_id;
| CREATE VIEW shop.order_names AS SELECT b.id, b.doubled, c.name FROM shop.big_orders b JOIN shop.customers c ON c.id = b.customer_id;
| CREATE VIEW shop.renamed (a, b) AS SELECT id, name FROM shop.customers;
| CREATE VIEW shop.all_ids AS SELECT id FROM shop.orders UNION SELECT id FROM shop.customers;
| CREATE VIEW shop.staff_pay AS SELECT s.id, s.salary, d.name AS dept FROM hr.staff s JOIN hr.depts d ON d.id = s.dept_id;
| CREATE VIEW hr.rich AS SELECT name, salary FROM hr.staff WHERE salary > 120;
| CREATE VIEW hr.`Rich_Names` AS SELECT name AS `VipName` FROM hr.rich;
| CREATE MATERIALIZED VIEW shop.order_sums BUILD DEFERRED REFRESH COMPLETE ON MANUAL DISTRIBUTED BY RANDOM BUCKETS 1 PROPERTIES ("replication_num" = "1") AS SELECT customer_id, SUM(amount) AS total FROM shop.orders GROUP BY customer_id;
< raw
| {"databases":[{"name":"hr","tables":[{"name":"depts","columns":["id","name"]},{"name":"staff","columns":["id","name","salary","dept_id"]}],"views":[{"name":"Rich_Names","columns":["VipName"],"definition":"SELECT `internal`.`hr`.`rich`.`name` AS `VipName` FROM `internal`.`hr`.`rich`"},{"name":"rich","columns":["name","salary"],"definition":"SELECT `internal`.`hr`.`staff`.`name`, `internal`.`hr`.`staff`.`salary` FROM `internal`.`hr`.`staff` WHERE `internal`.`hr`.`staff`.`salary` \u003e 120"}]},{"name":"shop","tables":[{"name":"MixedCase","columns":["Id","CamelCol","lower_col"]},{"name":"customers","columns":["id","name","region","email"]},{"name":"items","columns":["order_id","sku","qty"]},{"name":"mixedcase","columns":["id","note"]},{"name":"order_sums","columns":["customer_id","total"]},{"name":"orders","columns":["id","customer_id","amount","status","created_at"]}],"views":[{"name":"all_ids","columns":["id"],"definition":"SELECT `internal`.`shop`.`orders`.`id` FROM `internal`.`shop`.`orders` UNION SELECT `internal`.`shop`.`customers`.`id` FROM `internal`.`shop`.`customers`"},{"name":"big_orders","columns":["id","doubled","status","customer_id"],"definition":"SELECT `internal`.`shop`.`orders`.`id`, `internal`.`shop`.`orders`.`amount` * 2 AS `doubled`, `internal`.`shop`.`orders`.`status`, `internal`.`shop`.`orders`.`customer_id` FROM `internal`.`shop`.`orders` WHERE `internal`.`shop`.`orders`.`amount` \u003e 5"},{"name":"customer_totals","columns":["customer_id","total","n"],"definition":"SELECT `internal`.`shop`.`orders`.`customer_id`, SUM(`internal`.`shop`.`orders`.`amount`) AS `total`, COUNT(*) AS `n` FROM `internal`.`shop`.`orders` GROUP BY `internal`.`shop`.`orders`.`customer_id`"},{"name":"order_names","columns":["id","doubled","name"],"definition":"SELECT `internal`.`shop`.`b`.`id`, `internal`.`shop`.`b`.`doubled`, `internal`.`shop`.`c`.`name` FROM `internal`.`shop`.`big_orders` `b` JOIN `internal`.`shop`.`customers` `c` ON `internal`.`shop`.`c`.`id` = `internal`.`shop`.`b`.`customer_id`"},{"name":"renamed","columns":["a","b"],"definition":"SELECT `internal`.`shop`.`customers`.`id` AS `a`, `internal`.`shop`.`customers`.`name` AS `b` FROM `internal`.`shop`.`customers`"},{"name":"staff_pay","columns":["id","salary","dept"],"definition":"SELECT `internal`.`hr`.`s`.`id`, `internal`.`hr`.`s`.`salary`, `internal`.`hr`.`d`.`name` AS `dept` FROM `internal`.`hr`.`staff` `s` JOIN `internal`.`hr`.`depts` `d` ON `internal`.`hr`.`d`.`id` = `internal`.`hr`.`s`.`dept_id`"}]}]}
> raw
| span shop
| SELECT id, amount FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT o.id AS order_id, o.amount FROM orders o
< raw
| {"columns":[{"name":"order_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT id + customer_id, upper(status) FROM orders
< raw
| {"columns":[{"name":"id + customer_id"},{"name":"upper(status)"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT CAST(id AS VARCHAR) AS a, concat(name, '@', region) AS b FROM customers
< raw
| {"columns":[{"name":"a"},{"name":"b"}],"relations":["shop.customers"]}
> raw
| span shop
| SELECT 1, 'a' AS b
< raw
| {"columns":[{"name":"1"},{"name":"b"}]}
> raw
| span shop
| SELECT count(*) FROM orders
< raw
| {"columns":[{"name":"count(*)"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT DISTINCT region FROM customers
< raw
| {"columns":[{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"}],"relations":["shop.customers"]}
> raw
| span shop
| SELECT name FROM customers ORDER BY region LIMIT 2
< raw
| {"columns":[{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"}],"relations":["shop.customers"]}
> raw
| span shop
| SELECT CASE WHEN amount > 5 THEN status ELSE 'small' END AS label FROM orders
< raw
| {"columns":[{"name":"label"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT coalesce(email, name) AS contact FROM customers
< raw
| {"columns":[{"name":"contact"}],"relations":["shop.customers"]}
> raw
| span shop
| SELECT amount AS id, id AS amount FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT id   +   1, -amount, id IN (1, 2), NOT id FROM orders
< raw
| {"columns":[{"name":"id   +   1"},{"name":"-amount"},{"name":"id IN (1, 2)"},{"name":"NOT id"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT 'x', 2.50, NULL, TRUE, "y"
< raw
| {"columns":[{"name":"x"},{"name":"2.50"},{"name":"NULL"},{"name":"TRUE"},{"name":"y"}]}
> raw
| span shop
| SELECT (id), ((amount)) FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| VALUES (1, 2 + 3), (4, 5)
< raw
| {"columns":[{"name":"1"},{"name":"2 + 3"}]}
> raw
| span shop
| SELECT * FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT o.* FROM orders o
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"o","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"o","orgTable":"orders","orgName":"created_at"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT *, amount FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT c.*, o.amount FROM customers c JOIN orders o ON o.customer_id = c.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"relations":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM (SELECT * FROM customers) c
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"}],"relations":["shop.customers"]}
> raw
| span shop
| SELECT * EXCEPT (amount, status) FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT id FROM shop.orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT shop.orders.amount, orders.status FROM shop.orders
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT s.name, s.salary FROM hr.staff s
< raw
| {"columns":[{"name":"name","schema":"hr","table":"s","orgTable":"staff","orgName":"name"},{"name":"salary","schema":"hr","table":"s","orgTable":"staff","orgName":"salary"}],"relations":["hr.staff"]}
> raw
| span shop
| SELECT o.id, s.name FROM orders o JOIN hr.staff s ON s.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"hr","table":"s","orgTable":"staff","orgName":"name"}],"relations":["shop.orders","hr.staff"]}
> raw
| span shop
| SELECT id, amount FROM internal.shop.orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT o.id FROM internal.shop.orders o JOIN internal.hr.staff s ON s.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"relations":["shop.orders","hr.staff"]}
> raw
| span shop
| SELECT internal.shop.orders.amount FROM internal.shop.orders
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT * FROM internal.hr.rich
< raw
| {"columns":[{"name":"name","schema":"hr","table":"rich","orgTable":"rich","orgName":"name"},{"name":"salary","schema":"hr","table":"rich","orgTable":"rich","orgName":"salary"}],"relations":["hr.rich","hr.staff"]}
> raw
| span shop
| SELECT shop.o.id, internal.shop.o.amount FROM orders o
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT ID, Amount FROM orders
< raw
| {"columns":[{"name":"ID","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"Amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT `id`, `AMOUNT` FROM `orders`
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"AMOUNT","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT O.ID FROM orders O
< raw
| {"columns":[{"name":"ID","schema":"shop","table":"O","orgTable":"orders","orgName":"id"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT `Id`, `CamelCol`, lower_col FROM `MixedCase`
< raw
| {"columns":[{"name":"Id","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"Id"},{"name":"CamelCol","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"CamelCol"},{"name":"lower_col","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"lower_col"}],"relations":["shop.MixedCase"]}
> raw
| span shop
| SELECT * FROM `MixedCase`
< raw
| {"columns":[{"name":"Id","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"Id"},{"name":"CamelCol","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"CamelCol"},{"name":"lower_col","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"lower_col"}],"relations":["shop.MixedCase"]}
> raw
| span shop
| SELECT * FROM mixedcase
< raw
| {"columns":[{"name":"id","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"id"},{"name":"note","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"note"}],"relations":["shop.mixedcase"]}
> raw
| span shop
| SELECT m.camelcol, m.LOWER_COL FROM `MixedCase` m
< raw
| {"columns":[{"name":"camelcol","schema":"shop","table":"m","orgTable":"MixedCase","orgName":"CamelCol"},{"name":"LOWER_COL","schema":"shop","table":"m","orgTable":"MixedCase","orgName":"lower_col"}],"relations":["shop.MixedCase"]}
> raw
| span shop
| SELECT id FROM Orders
< raw
| {"error":"errCode = 2, detailMessage = Table [Orders] does not exist in database [shop]."}
> raw
| span shop
| SELECT orders.name FROM customers AS orders
< raw
| {"columns":[{"name":"name","schema":"shop","table":"orders","orgTable":"customers","orgName":"name"}],"relations":["shop.customers"]}
> raw
| span shop
| SELECT customers.amount FROM orders customers
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"customers","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT o.id FROM customers o JOIN orders c ON c.customer_id = o.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"customers","orgName":"id"}],"relations":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM (SELECT amount AS id FROM orders) orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT o.id, c.name FROM customers c RIGHT JOIN orders o ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o FULL OUTER JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o CROSS JOIN customers c
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT o.id FROM orders o LEFT SEMI JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT o.id FROM orders o LEFT ANTI JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT * FROM orders o LEFT SEMI JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"o","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"o","orgTable":"orders","orgName":"created_at"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT * FROM orders o RIGHT ANTI JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"}],"relations":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT o.id FROM orders o JOIN customers c USING (id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT id, name FROM orders JOIN customers USING (id)
< raw
| {"error":"errCode = 2, detailMessage = id is ambiguous: id#0, id#5."}
> raw
| span shop
| SELECT * FROM orders o JOIN customers c USING (id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"o","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"o","orgTable":"orders","orgName":"created_at"},{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT * FROM orders o LEFT JOIN customers c USING (id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"o","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"o","orgTable":"orders","orgName":"created_at"},{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT * FROM orders NATURAL JOIN customers
< raw
| {"error":"errCode = 2, detailMessage = \nmismatched input 'NATURAL' expecting {\u003cEOF\u003e, ';'}(line 1, pos 21)\n"}
> raw
| span shop
| SELECT i.sku, o.status, c.region FROM items i JOIN orders o ON o.id = i.order_id JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"sku","schema":"shop","table":"i","orgTable":"items","orgName":"sku"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"}],"relations":["shop.items","shop.orders","shop.customers"]}
> raw
| span shop
| SELECT a.id, b.id FROM orders a JOIN orders b ON b.customer_id = a.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"a","orgTable":"orders","orgName":"id"},{"name":"id","schema":"shop","table":"b","orgTable":"orders","orgName":"id"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT o.id FROM orders o JOIN [shuffle] customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"who"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT id, (SELECT max(amount) FROM orders) AS top FROM customers
< raw
| {"columns":[{"name":"id","schema":"shop","table":"customers","orgTable":"customers","orgName":"id"},{"name":"top"}],"relations":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'east')
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT id FROM orders o WHERE EXISTS (SELECT 1 FROM items i WHERE i.order_id = o.id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"relations":["shop.orders","shop.items"]}
> raw
| span shop
| SELECT d.k, d.v FROM (SELECT id AS k, amount AS v FROM orders) d
< raw
| {"columns":[{"name":"k","schema":"shop","table":"d","orgTable":"orders","orgName":"id"},{"name":"v","schema":"shop","table":"d","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT d.k FROM (SELECT id AS k FROM orders LIMIT 2) d
< raw
| {"columns":[{"name":"k","schema":"shop","table":"d","orgTable":"orders","orgName":"id"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT x.total FROM (SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id) x
< raw
| {"columns":[{"name":"total","table":"x"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT y.k FROM (SELECT x.k FROM (SELECT id AS k FROM customers) x) y
< raw
| {"columns":[{"name":"k","schema":"shop","table":"y","orgTable":"customers","orgName":"id"}],"relations":["shop.customers"]}
> raw
| span shop
| SELECT d.a FROM (SELECT 1 AS a, 2) d
< raw
| {"columns":[{"name":"a","table":"d"}]}
> raw
| span shop
| SELECT o.id, e.x FROM orders o LATERAL VIEW explode(array(1, 2)) e AS x
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"x","table":"e"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT o.id, l.n FROM orders o, LATERAL (SELECT count(*) AS n FROM items i WHERE i.order_id = o.id) l
< raw
| {"error":"errCode = 2, detailMessage = \nextraneous input 'LATERAL' expecting {'(', '{', '}', 'ACTIONS', 'AFTER', 'AGG_STATE', 'AGGREGATE', 'ALIAS', 'ANALYZED', 'ARRAY', 'AT', 'AUTHORS', 'AUTO_INCREMENT', 'ALWAYS', 'BACKENDS', 'BACKUP', 'BEGIN', 'BELONG', 'BIN', 'BITAND', 'BITMAP', 'BITMAP_EMPTY', 'BITMAP_UNION', 'BITOR', 'BITXOR', 'BLOB', 'BOOLEAN', 'BRANCH', 'BRIEF', 'BROKER', 'BUCKETS', 'BUILD', 'BUILTIN', 'BULK', 'CACHE', 'CACHED', 'CALL', 'CATALOG', 'CATALOGS', 'CHAIN', CHAR, 'CHARSET', 'CHECK', 'CLUSTER', 'CLUSTERS', 'COLLATION', 'COLLECT', 'COLOCATE', 'COLUMNS', 'COMMENT', 'COMMIT', 'COMMITTED', 'COMPACT', 'COMPLETE', 'COMPRESS_TYPE', 'COMPUTE', 'CONDITIONS', 'CONFIG', 'CONNECTION', 'CONNECTION_ID', 'CONSISTENT', 'CONSTRAINTS', 'CONVERT', 'CONVERT_LIGHT_SCHEMA_CHANGE_PROCESS', 'COPY', 'COUNT', 'CREATION', 'CRON', 'CURRENT_CATALOG', 'CURRENT_DATE', 'CURRENT_TIME', 'CURRENT_TIMESTAMP', 'CURRENT_USER', 'DATA', 'DATE', 'DATETIME', 'DATETIMEV2', 'DATEV2', 'DATETIMEV1', 'DATEV1', 'DAY', 'DAYS', 'DECIMAL', 'DECIMALV2', 'DECIMALV3', 'DEFERRED', 'DEMAND', 'DIAGNOSE', 'DIAGNOSIS', 'DISTINCTPC', 'DISTINCTPCSA', 'DO', 'DORIS_INTERNAL_TABLE_ID', 'DUAL', 'DYNAMIC', 'E', 'ENABLE', 'ENCRYPTION', 'ENCRYPTKEY', 'ENCRYPTKEYS', 'END', 'ENDS', 'ENGINE', 'ENGINES', 'ERRORS', 'ESCAPE', 'EVENTS', 'EVERY', 'EXCLUDE', 'EXPIRED', 'EXTERNAL', 'FAILED_LOGIN_ATTEMPTS', 'FAST', 'FEATURE', 'FIELDS', 'FILE', 'FILTER', 'FIRST', 'FORMAT', 'FREE', 'FRONTENDS', 'FUNCTION', 'GENERATED', 'GENERIC', 'GLOBAL', 'GRAPH', 'GROUPING', 'GROUPS', 'HASH', 'HDFS', 'HELP', 'HISTOGRAM', 'HLL_UNION', 'HOSTNAME', 'HOTSPOT', 'HOUR', 'HOURS', 'HUB', 'IDENTIFIED', 'IGNORE', 'IMMEDIATE', 'INCREMENTAL', 'INDEXES', 'INVERTED', 'IPV4', 'IPV6', 'IS_NOT_NULL_PRED', 'IS_NULL_PRED', 'ISNULL', 'ISOLATION', 'JOB', 'JOBS', 'JSON', 'JSONB', 'LABEL', 'LAST', 'LDAP', 'LDAP_ADMIN_PASSWORD', 'LESS', 'LEVEL', 'LINES', 'LINK', 'LOCAL', 'LOCALTIME', 'LOCALTIMESTAMP', 'LOCATION', 'LOCK', 'LOGICAL', 'MANUAL', 'MAP', 'MATCH_ALL', 'MATCH_ANY', 'MATCH_PHRASE', 'MATCH_PHRASE_EDGE', 'MATCH_PHRASE_PREFIX', 'MATCH_REGEXP', 'MATERIALIZED', 'MAX', 'MEMO', 'MERGE', 'MIGRATE', 'MIGRATIONS', 'MIN', 'MINUTE', 'MINUTES', 'MODIFY', 'MONTH', 'MTMV', 'NAME', 'NAMES', 'NEGATIVE', 'NEVER', 'NEXT', 'NGRAM_BF', 'NO', 'NON_NULLABLE', 'NULLS', 'OF', 'OFFSET', 'ONLY', 'OPEN', 'OPTIMIZED', 'PARAMETER', 'PARSED', 'PARTITIONS', 'PASSWORD', 'PASSWORD_EXPIRE', 'PASSWORD_HISTORY', 'PASSWORD_LOCK_TIME', 'PASSWORD_REUSE', 'PATH', 'PAUSE', 'PERCENT', 'PERIOD', 'PERMISSIVE', 'PHYSICAL', 'PI', 'PLAN', 'PRIVILEGES', 'PROCESS', 'PLUGIN', 'PLUGINS', 'POLICY', 'PROC', 'PROCESSLIST', 'PROFILE', 'PROPERTIES', 'PROPERTY', 'QUANTILE_STATE', 'QUANTILE_UNION', 'QUERY', 'QUOTA', 'RANDOM', 'RECENT', 'RECOVER', 'RECYCLE', 'REFRESH', 'REPEATABLE', 'REPLACE', 'REPLACE_IF_NOT_NULL', 'REPLAYER', 'REPOSITORIES', 'REPOSITORY', 'RESOURCE', 'RESOURCES', 'RESTORE', 'RESTRICTIVE', 'RESUME', 'RETAIN', 'RETENTION', 'RETURNS', 'REWRITTEN', 'RLIKE', 'ROLLBACK', 'ROLLUP', 'ROOT', 'ROTATE', 'ROUTINE', 'S3', 'SAMPLE', 'SCHEDULE', 'SCHEDULER', 'SCHEMA', 'SECOND', 'SERIALIZABLE', 'SESSION', 'SESSION_USER', 'SET_SESSION_VARIABLE', 'SHAPE', 'SKEW', 'SNAPSHOT', 'SNAPSHOTS', 'SONAME', 'SPLIT', 'SQL', 'STAGE', 'STAGES', 'START', 'STARTS', 'STATS', 'STATUS', 'STOP', 'STORAGE', 'STREAM', 'STREAMING', 'STRING', 'STRUCT', 'SUM', 'TABLES', 'TAG', 'TASK', 'TASKS', 'TDE', 'TEMPORARY', 'TEXT', 'THAN', 'TIME', 'TIMESTAMP', 'TRANSACTION', 'TREE', 'TRIGGERS', 'TRUNCATE', 'TYPE', 'TYPES', 'UNCOMMITTED', 'UNLOCK', 'UNSET', 'UP', 'USER', 'VALUE', 'VARCHAR', 'VARIABLE', 'VARIABLES', 'VARIANT', 'VAULT', 'VAULTS', 'VERBOSE', 'VERSION', 'VIEW', 'VIEWS', 'WARM', 'WARNINGS', 'WEEK', 'WORK', 'YEAR', '/*+', '*/', '/*', IDENTIFIER, BACKQUOTED_IDENTIFIER}(line 1, pos 32)\n"}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers) SELECT id, name FROM c
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"relations":["shop.customers"]}
> raw
| span shop
| WITH c (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM c
< raw
| {"columns":[{"name":"k","schema":"shop","table":"c","orgTable":"orders","orgName":"id"},{"name":"v","schema":"shop","table":"c","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span shop
| WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT * FROM b
< raw
| {"columns":[{"name":"id","schema":"shop","table":"b","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"b","orgTable":"customers","orgName":"name"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| WITH c AS (SELECT id FROM customers) SELECT x.id, y.id FROM c x JOIN c y ON y.id = x.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"x","orgTable":"customers","orgName":"id"},{"name":"id","schema":"shop","table":"y","orgTable":"customers","orgName":"id"}],"relations":["shop.customers"]}
> raw
| span shop
| WITH orders AS (SELECT id, name FROM customers) SELECT id, name FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"orders","orgTable":"customers","orgName":"name"}],"relations":["shop.customers"]}
> raw
| span shop
| WITH c AS (SELECT id, amount FROM orders) SELECT id FROM c WHERE amount > (SELECT avg(amount) FROM c)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"orders","orgName":"id"}],"relations":["shop.orders"]}
> raw
| span shop
| WITH RECURSIVE r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r
< raw
| {"error":"errCode = 2, detailMessage = \nno viable alternative at input 'WITH RECURSIVE r'(line 1, pos 15)\n"}
> raw
| span shop
| SELECT id FROM orders UNION SELECT id FROM customers
< raw
| {"columns":[{"name":"id"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT id, status FROM orders UNION ALL SELECT id, name FROM customers
< raw
| {"columns":[{"name":"id"},{"name":"status"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT name FROM customers INTERSECT SELECT status FROM orders
< raw
| {"columns":[{"name":"name"}],"relations":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM orders EXCEPT SELECT customer_id FROM orders
< raw
| {"columns":[{"name":"id"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT id FROM orders MINUS SELECT customer_id FROM orders
< raw
| {"columns":[{"name":"id"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT id AS k FROM orders UNION SELECT 1 ORDER BY k
< raw
| {"error":"errCode = 2, detailMessage = Unknown column 'k' in 'table list' in SORT clause"}
> raw
| span shop
| (SELECT id FROM orders LIMIT 1) UNION (SELECT id FROM customers LIMIT 1)
< raw
| {"columns":[{"name":"id"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| WITH c AS (SELECT id FROM orders) VALUES (0) UNION ALL SELECT id FROM c
< raw
| {"columns":[{"name":"0"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT u.id FROM (SELECT id FROM orders UNION SELECT id FROM customers) u
< raw
| {"columns":[{"name":"id","table":"u"}],"relations":["shop.orders","shop.customers"]}
> raw
| span shop
| SELECT id, doubled FROM big_orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"id"},{"name":"doubled","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"doubled"}],"relations":["shop.big_orders","shop.orders"]}
> raw
| span shop
| SELECT b.status, b.doubled FROM big_orders b
< raw
| {"columns":[{"name":"status","schema":"shop","table":"b","orgTable":"big_orders","orgName":"status"},{"name":"doubled","schema":"shop","table":"b","orgTable":"big_orders","orgName":"doubled"}],"relations":["shop.big_orders","shop.orders"]}
> raw
| span shop
| SELECT * FROM big_orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"id"},{"name":"doubled","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"doubled"},{"name":"status","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"status"},{"name":"customer_id","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"customer_id"}],"relations":["shop.big_orders","shop.orders"]}
> raw
| span shop
| SELECT customer_id, total, n FROM customer_totals
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"customer_totals","orgTable":"customer_totals","orgName":"customer_id"},{"name":"total","schema":"shop","table":"customer_totals","orgTable":"customer_totals","orgName":"total"},{"name":"n","schema":"shop","table":"customer_totals","orgTable":"customer_totals","orgName":"n"}],"relations":["shop.customer_totals","shop.orders"]}
> raw
| span shop
| SELECT * FROM order_names
< raw
| {"columns":[{"name":"id","schema":"shop","table":"order_names","orgTable":"order_names","orgName":"id"},{"name":"doubled","schema":"shop","table":"order_names","orgTable":"order_names","orgName":"doubled"},{"name":"name","schema":"shop","table":"order_names","orgTable":"order_names","orgName":"name"}],"relations":["shop.order_names","shop.big_orders","shop.orders","shop.customers"]}
> raw
| span shop
| SELECT o.name, c.region FROM order_names o JOIN customers c ON c.name = o.name
< raw
| {"columns":[{"name":"name","schema":"shop","table":"o","orgTable":"order_names","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"}],"relations":["shop.order_names","shop.big_orders","shop.orders","shop.customers"]}
> raw
| span shop
| SELECT a, b FROM renamed
< raw
| {"columns":[{"name":"a","schema":"shop","table":"renamed","orgTable":"renamed","orgName":"a"},{"name":"b","schema":"shop","table":"renamed","orgTable":"renamed","orgName":"b"}],"relations":["shop.renamed","shop.customers"]}
> raw
| span shop
| SELECT id FROM all_ids
< raw
| {"columns":[{"name":"id","schema":"shop","table":"all_ids","orgTable":"all_ids","orgName":"id"}],"relations":["shop.all_ids","shop.orders","shop.customers"]}
> raw
| span shop
| SELECT * FROM staff_pay
< raw
| {"columns":[{"name":"id","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"id"},{"name":"salary","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"salary"},{"name":"dept","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"dept"}],"relations":["shop.staff_pay","hr.staff","hr.depts"]}
> raw
| span shop
| SELECT name FROM hr.rich
< raw
| {"columns":[{"name":"name","schema":"hr","table":"rich","orgTable":"rich","orgName":"name"}],"relations":["hr.rich","hr.staff"]}
> raw
| span shop
| SELECT VipName FROM hr.`Rich_Names`
< raw
| {"columns":[{"name":"VipName","schema":"hr","table":"Rich_Names","orgTable":"Rich_Names","orgName":"VipName"}],"relations":["hr.Rich_Names","hr.rich","hr.staff"]}
> raw
| span shop
| SELECT b.id, t.total FROM big_orders b JOIN customer_totals t ON t.customer_id = b.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"b","orgTable":"big_orders","orgName":"id"},{"name":"total","schema":"shop","table":"t","orgTable":"customer_totals","orgName":"total"}],"relations":["shop.big_orders","shop.orders","shop.customer_totals"]}
> raw
| span shop
| SELECT d.doubled FROM (SELECT doubled FROM big_orders) d
< raw
| {"columns":[{"name":"doubled","schema":"shop","table":"d","orgTable":"big_orders","orgName":"doubled"}],"relations":["shop.big_orders","shop.orders"]}
> raw
| span shop
| WITH v AS (SELECT * FROM big_orders) SELECT doubled FROM v
< raw
| {"columns":[{"name":"doubled","schema":"shop","table":"v","orgTable":"big_orders","orgName":"doubled"}],"relations":["shop.big_orders","shop.orders"]}
> raw
| span shop
| SELECT customer_id, total FROM order_sums
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"order_sums","orgTable":"order_sums","orgName":"customer_id"},{"name":"total","schema":"shop","table":"order_sums","orgTable":"order_sums","orgName":"total"}],"relations":["shop.order_sums"]}
> raw
| span shop
| SELECT * FROM internal.shop.order_sums
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"order_sums","orgTable":"order_sums","orgName":"customer_id"},{"name":"total","schema":"shop","table":"order_sums","orgTable":"order_sums","orgName":"total"}],"relations":["shop.order_sums"]}
> raw
| span shop
| SELECT id, row_number() OVER (ORDER BY amount) AS rn FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"rn"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT id, sum(amount) OVER (PARTITION BY customer_id) AS running FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"running"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT id, lag(amount) OVER (ORDER BY id) AS prev FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"prev"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT id, rank() OVER (PARTITION BY status ORDER BY amount) FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"rank() OVER (PARTITION BY status ORDER BY amount)"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT id FROM orders QUALIFY row_number() OVER (PARTITION BY customer_id ORDER BY amount) = 1
< raw
| {"error":"errCode = 2, detailMessage = \nmismatched input 'row_number' expecting {\u003cEOF\u003e, ';'}(line 1, pos 30)\n"}
> raw
| span shop
| SELECT customer_id, sum(amount), count(*) FROM orders GROUP BY customer_id
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"sum(amount)"},{"name":"count(*)"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT status, avg(amount) AS a FROM orders GROUP BY status HAVING count(*) > 1
< raw
| {"columns":[{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"a"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT region, group_concat(name ORDER BY name) AS names FROM customers GROUP BY region
< raw
| {"columns":[{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"},{"name":"names"}],"relations":["shop.customers"]}
> raw
| span shop
| SELECT max(amount) FROM orders WHERE status = 'open'
< raw
| {"columns":[{"name":"max(amount)"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT count(DISTINCT customer_id) AS buyers FROM orders
< raw
| {"columns":[{"name":"buyers"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT c.region, sum(o.amount) FROM customers c JOIN orders o ON o.customer_id = c.id GROUP BY c.region
< raw
| {"columns":[{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"sum(o.amount)"}],"relations":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT customer_id, status, sum(amount) FROM orders GROUP BY ROLLUP (customer_id, status)
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"sum(amount)"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT count(*) FROM orders WHERE 1 = 0
< raw
| {"columns":[{"name":"count(*)"}],"relations":["shop.orders"]}
> raw
| span shop
| SELECT number FROM numbers("number" = "3")
< raw
| {"columns":[{"name":"number","table":"_table_valued_function_numbers","orgTable":"NumbersTableValuedFunction","orgName":"number"}]}
> raw
| span shop
| SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_SCHEMA = 'shop'
< raw
| {"columns":[{"name":"TABLE_NAME","schema":"information_schema","table":"tables","orgTable":"tables","orgName":"TABLE_NAME"}],"relations":["information_schema.tables"]}
> raw
| span shop
| SELECT nope FROM orders
< raw
| {"error":"errCode = 2, detailMessage = Unknown column 'nope' in 'table list' in PROJECT clause"}
> raw
| span shop
| SELECT id FROM no_such_table
< raw
| {"error":"errCode = 2, detailMessage = Table [no_such_table] does not exist in database [shop]."}
> raw
| span shop
| SELECT id FROM orders JOIN customers ON customers.id = orders.customer_id
< raw
| {"error":"errCode = 2, detailMessage = id is ambiguous: id#0, id#5."}
> raw
| span shop
| SELECT id FROM no_such_db.orders
< raw
| {"error":"errCode = 2, detailMessage = Database [no_such_db] does not exist."}
> raw
| span hr
| SELECT name, salary FROM staff
< raw
| {"columns":[{"name":"name","schema":"hr","table":"staff","orgTable":"staff","orgName":"name"},{"name":"salary","schema":"hr","table":"staff","orgTable":"staff","orgName":"salary"}],"relations":["hr.staff"]}
> raw
| span hr
| SELECT s.name, d.name AS dept FROM staff s JOIN depts d ON d.id = s.dept_id
< raw
| {"columns":[{"name":"name","schema":"hr","table":"s","orgTable":"staff","orgName":"name"},{"name":"dept","schema":"hr","table":"depts","orgTable":"depts","orgName":"name"}],"relations":["hr.staff","hr.depts"]}
> raw
| span hr
| SELECT o.amount FROM shop.orders o
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"relations":["shop.orders"]}
> raw
| span hr
| SELECT * FROM rich
< raw
| {"columns":[{"name":"name","schema":"hr","table":"rich","orgTable":"rich","orgName":"name"},{"name":"salary","schema":"hr","table":"rich","orgTable":"rich","orgName":"salary"}],"relations":["hr.rich","hr.staff"]}
> raw
| span hr
| SELECT * FROM shop.staff_pay
< raw
| {"columns":[{"name":"id","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"id"},{"name":"salary","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"salary"},{"name":"dept","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"dept"}],"relations":["shop.staff_pay","hr.staff","hr.depts"]}
> raw
| span hr
| SELECT name FROM staff UNION SELECT name FROM shop.customers
< raw
| {"columns":[{"name":"name"}],"relations":["hr.staff","shop.customers"]}
> raw
| span hr
| SELECT id FROM internal.shop.big_orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"id"}],"relations":["shop.big_orders","shop.orders"]}
