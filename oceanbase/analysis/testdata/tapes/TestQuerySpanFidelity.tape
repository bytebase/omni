-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestQuerySpanFidelity default#1
> raw
| fixture
| CREATE DATABASE shop;
| CREATE DATABASE hr;
| CREATE TABLE shop.customers (id INT PRIMARY KEY, name VARCHAR(20), region VARCHAR(10), email VARCHAR(40));
| CREATE TABLE shop.orders (id INT PRIMARY KEY, customer_id INT, amount DECIMAL(10,2), status VARCHAR(10), created_at DATE);
| CREATE TABLE shop.items (order_id INT, sku VARCHAR(10), qty INT);
| CREATE TABLE shop.`MixedCase` (`Id` INT, `CamelCol` VARCHAR(10), lower_col VARCHAR(10));
| CREATE TABLE hr.staff (id INT PRIMARY KEY, name VARCHAR(20), salary INT, dept_id INT);
| CREATE TABLE hr.depts (id INT PRIMARY KEY, name VARCHAR(20));
| INSERT INTO shop.customers VALUES (1, 'ann', 'east', 'a@x'), (2, 'bob', 'west', 'b@x'), (3, 'cy', 'east', 'c@x');
| INSERT INTO shop.orders VALUES (1, 1, 10.00, 'open', '2026-01-01'), (2, 2, 20.00, 'paid', '2026-01-02'), (3, 1, 5.00, 'open', '2026-01-03');
| INSERT INTO shop.items VALUES (1, 'a', 1), (1, 'b', 2), (2, 'a', 3);
| INSERT INTO shop.`MixedCase` VALUES (1, 'x', 'y'), (2, 'z', 'w');
| INSERT INTO hr.staff VALUES (1, 'dee', 100, 1), (2, 'eve', 200, 2), (3, 'fay', 150, 1);
| INSERT INTO hr.depts VALUES (1, 'ops'), (2, 'dev');
| CREATE VIEW shop.big_orders AS SELECT id, amount * 2 AS doubled, status, customer_id FROM shop.orders WHERE amount > 5;
| CREATE VIEW shop.customer_totals AS SELECT customer_id, SUM(amount) AS total, COUNT(*) AS n FROM shop.orders GROUP BY customer_id;
| CREATE VIEW shop.order_names AS SELECT b.id, b.doubled, c.name FROM shop.big_orders b JOIN shop.customers c ON c.id = b.customer_id;
| CREATE VIEW shop.renamed (a, b) AS SELECT id, name FROM shop.customers;
| CREATE VIEW shop.all_ids AS SELECT id FROM shop.orders UNION SELECT id FROM shop.customers;
| CREATE VIEW shop.staff_pay AS SELECT s.id, s.salary, d.name AS dept FROM hr.staff s JOIN hr.depts d ON d.id = s.dept_id;
| CREATE VIEW hr.rich AS SELECT name, salary FROM hr.staff WHERE salary > 120;
| CREATE VIEW hr.`Rich_Names` AS SELECT name AS `VipName` FROM hr.rich;
< raw
| {"databases":[{"name":"hr","tables":[{"name":"depts","columns":["id","name"]},{"name":"staff","columns":["id","name","salary","dept_id"]}],"views":[{"name":"rich","columns":["name","salary"],"definition":"select `hr`.`staff`.`name` AS `name`,`hr`.`staff`.`salary` AS `salary` from `hr`.`staff` where (`hr`.`staff`.`salary` \u003e 120)"},{"name":"rich_names","columns":["VipName"],"definition":"select `hr`.`rich`.`name` AS `VipName` from `hr`.`rich`"}]},{"name":"shop","tables":[{"name":"customers","columns":["id","name","region","email"]},{"name":"items","columns":["order_id","sku","qty"]},{"name":"mixedcase","columns":["Id","CamelCol","lower_col"]},{"name":"orders","columns":["id","customer_id","amount","status","created_at"]}],"views":[{"name":"all_ids","columns":["id"],"definition":"(select `shop`.`orders`.`id` AS `id` from `shop`.`orders`) union (select `shop`.`customers`.`id` AS `id` from `shop`.`customers`)"},{"name":"big_orders","columns":["id","doubled","status","customer_id"],"definition":"select `shop`.`orders`.`id` AS `id`,(`shop`.`orders`.`amount` * 2) AS `doubled`,`shop`.`orders`.`status` AS `status`,`shop`.`orders`.`customer_id` AS `customer_id` from `shop`.`orders` where (`shop`.`orders`.`amount` \u003e 5)"},{"name":"customer_totals","columns":["customer_id","total","n"],"definition":"select `shop`.`orders`.`customer_id` AS `customer_id`,sum(`shop`.`orders`.`amount`) AS `total`,count(*) AS `n` from `shop`.`orders` group by `shop`.`orders`.`customer_id`"},{"name":"order_names","columns":["id","doubled","name"],"definition":"select `shop`.`b`.`id` AS `id`,`shop`.`b`.`doubled` AS `doubled`,`c`.`name` AS `name` from (`shop`.`big_orders` `b` join `shop`.`customers` `c` on (`c`.`id` = `shop`.`b`.`customer_id`))"},{"name":"renamed","columns":["a","b"],"definition":"select `shop`.`customers`.`id` AS `a`,`shop`.`customers`.`name` AS `b` from `shop`.`customers`"},{"name":"staff_pay","columns":["id","salary","dept"],"definition":"select `s`.`id` AS `id`,`s`.`salary` AS `salary`,`d`.`name` AS `dept` from (`hr`.`staff` `s` join `hr`.`depts` `d` on (`d`.`id` = `s`.`dept_id`))"}]}]}
> raw
| span shop
| SELECT id, amount FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT o.id AS order_id, o.amount FROM orders o
< raw
| {"columns":[{"name":"order_id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id + customer_id, upper(status) FROM orders
< raw
| {"columns":[{"name":"id + customer_id"},{"name":"upper(status)"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT CAST(amount AS CHAR) AS a, concat(name, '@', region) AS b FROM customers
< raw
| {"error":"1054"}
> raw
| span shop
| SELECT 1, 'a' AS b
< raw
| {"columns":[{"name":"1"},{"name":"b"}]}
> raw
| span shop
| SELECT count(*) FROM orders
< raw
| {"columns":[{"name":"count(*)"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT DISTINCT region FROM customers
< raw
| {"columns":[{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"}],"tables":["shop.customers"]}
> raw
| span shop
| SELECT name FROM customers ORDER BY region LIMIT 2
< raw
| {"columns":[{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"}],"tables":["shop.customers"]}
> raw
| span shop
| SELECT CASE WHEN amount > 5 THEN status ELSE 'small' END AS label FROM orders
< raw
| {"columns":[{"name":"label"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT coalesce(email, name) AS contact FROM customers
< raw
| {"columns":[{"name":"contact"}],"tables":["shop.customers"]}
> raw
| span shop
| SELECT amount AS id, id AS amount FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id   +   1, -amount, id IN (1, 2), NOT id FROM orders
< raw
| {"columns":[{"name":"id   +   1"},{"name":"-amount"},{"name":"id IN (1, 2)"},{"name":"NOT id"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT 'x', 2.50, NULL, TRUE, _utf8mb4'y'
< raw
| {"columns":[{"name":"x"},{"name":"2.50"},{"name":"NULL"},{"name":"TRUE"},{"name":"y"}]}
> raw
| span shop
| TABLE customers
< raw
| {"error":"1064"}
> raw
| span shop
| VALUES ROW(1, 2 + 3), ROW(4, 5)
< raw
| {"columns":[{"name":"column_0","table":"VALUES_TABLE1","orgTable":"VALUES_TABLE1","orgName":"column_0"},{"name":"column_1","table":"VALUES_TABLE1","orgTable":"VALUES_TABLE1","orgName":"column_1"}]}
> raw
| span shop
| SELECT * FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT o.* FROM orders o
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"o","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"o","orgTable":"orders","orgName":"created_at"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT *, amount FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT c.*, o.amount FROM customers c JOIN orders o ON o.customer_id = c.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM (SELECT * FROM customers) c
< raw
| {"columns":[{"name":"id","table":"c","orgTable":"c","orgName":"id"},{"name":"name","table":"c","orgTable":"c","orgName":"name"},{"name":"region","table":"c","orgTable":"c","orgName":"region"},{"name":"email","table":"c","orgTable":"c","orgName":"email"}],"tables":["shop.customers"]}
> raw
| span shop
| SELECT id FROM shop.orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT shop.orders.amount, orders.status FROM shop.orders
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT s.name, s.salary FROM hr.staff s
< raw
| {"columns":[{"name":"name","schema":"hr","table":"s","orgTable":"staff","orgName":"name"},{"name":"salary","schema":"hr","table":"s","orgTable":"staff","orgName":"salary"}],"tables":["hr.staff"]}
> raw
| span shop
| SELECT o.id, s.name FROM orders o JOIN hr.staff s ON s.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"hr","table":"s","orgTable":"staff","orgName":"name"}],"tables":["hr.staff","shop.orders"]}
> raw
| span shop
| SELECT ID, Amount FROM orders
< raw
| {"columns":[{"name":"ID","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"Amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT `id`, `AMOUNT` FROM `orders`
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"AMOUNT","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT O.ID FROM orders O
< raw
| {"columns":[{"name":"ID","schema":"shop","table":"O","orgTable":"orders","orgName":"id"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT shop.o.id, shop.orders.status FROM orders o JOIN shop.orders ON orders.id = o.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT `Id`, `CamelCol`, lower_col FROM `MixedCase`
< raw
| {"columns":[{"name":"Id","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"Id"},{"name":"CamelCol","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"CamelCol"},{"name":"lower_col","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"lower_col"}],"tables":["shop.mixedcase"]}
> raw
| span shop
| SELECT * FROM `MixedCase`
< raw
| {"columns":[{"name":"Id","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"Id"},{"name":"CamelCol","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"CamelCol"},{"name":"lower_col","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"lower_col"}],"tables":["shop.mixedcase"]}
> raw
| span shop
| SELECT * FROM mixedcase
< raw
| {"columns":[{"name":"Id","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"Id"},{"name":"CamelCol","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"CamelCol"},{"name":"lower_col","schema":"shop","table":"mixedcase","orgTable":"mixedcase","orgName":"lower_col"}],"tables":["shop.mixedcase"]}
> raw
| span shop
| SELECT m.camelcol, m.LOWER_COL FROM `MixedCase` m
< raw
| {"columns":[{"name":"camelcol","schema":"shop","table":"m","orgTable":"mixedcase","orgName":"CamelCol"},{"name":"LOWER_COL","schema":"shop","table":"m","orgTable":"mixedcase","orgName":"lower_col"}],"tables":["shop.mixedcase"]}
> raw
| span shop
| SELECT orders.name FROM customers AS orders
< raw
| {"columns":[{"name":"name","schema":"shop","table":"orders","orgTable":"customers","orgName":"name"}],"tables":["shop.customers"]}
> raw
| span shop
| SELECT customers.amount FROM orders customers
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"customers","orgTable":"orders","orgName":"amount"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT o.id FROM customers o JOIN orders c ON c.customer_id = o.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"customers","orgName":"id"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM (SELECT amount AS id FROM orders) orders
< raw
| {"columns":[{"name":"id","table":"orders","orgTable":"orders","orgName":"id"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT o.id, c.name FROM customers c RIGHT JOIN orders o ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o CROSS JOIN customers c
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o STRAIGHT_JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM orders JOIN customers USING (id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM orders JOIN customers USING (id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"},{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"customers","orgTable":"customers","orgName":"email"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM orders o LEFT JOIN customers c USING (id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"o","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"o","orgTable":"orders","orgName":"created_at"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM orders NATURAL JOIN customers
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"},{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"customers","orgTable":"customers","orgName":"email"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id, name FROM orders NATURAL LEFT JOIN customers
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT i.sku, o.status, c.region FROM items i JOIN orders o ON o.id = i.order_id JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"sku","schema":"shop","table":"i","orgTable":"items","orgName":"sku"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"}],"tables":["shop.customers","shop.items","shop.orders"]}
> raw
| span shop
| SELECT a.id, b.id FROM orders a JOIN orders b ON b.customer_id = a.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"a","orgTable":"orders","orgName":"id"},{"name":"id","schema":"shop","table":"b","orgTable":"orders","orgName":"id"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT o.id FROM orders o JOIN (customers c JOIN items i ON i.order_id = c.id) ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"tables":["shop.customers","shop.items","shop.orders"]}
> raw
| span shop
| SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"who"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id, (SELECT max(amount) FROM orders) AS top FROM customers
< raw
| {"columns":[{"name":"id","schema":"shop","table":"customers","orgTable":"customers","orgName":"id"},{"name":"top"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'east')
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM orders o WHERE EXISTS (SELECT 1 FROM items i WHERE i.order_id = o.id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"tables":["shop.items","shop.orders"]}
> raw
| span shop
| SELECT d.k, d.v FROM (SELECT id AS k, amount AS v FROM orders) d
< raw
| {"columns":[{"name":"k","table":"d","orgTable":"d","orgName":"k"},{"name":"v","table":"d","orgTable":"d","orgName":"v"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT d.k FROM (SELECT id AS k FROM orders LIMIT 2) d
< raw
| {"columns":[{"name":"k","table":"d","orgTable":"d","orgName":"k"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT x.total FROM (SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id) x
< raw
| {"columns":[{"name":"total","table":"x","orgTable":"x","orgName":"total"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT d.a, d.b FROM (SELECT id, name FROM customers) AS d (a, b)
< raw
| {"columns":[{"name":"a","table":"d","orgTable":"d","orgName":"a"},{"name":"b","table":"d","orgTable":"d","orgName":"b"}],"tables":["shop.customers"]}
> raw
| span shop
| SELECT o.id, l.n FROM orders o, LATERAL (SELECT count(*) AS n FROM items i WHERE i.order_id = o.id) l
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"n","table":"l","orgTable":"l","orgName":"n"}],"tables":["shop.items","shop.orders"]}
> raw
| span shop
| SELECT o.id, l.sku FROM orders o JOIN LATERAL (SELECT sku FROM items i WHERE i.order_id = o.id LIMIT 1) l ON true
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"sku","table":"l","orgTable":"l","orgName":"sku"}],"tables":["shop.items","shop.orders"]}
> raw
| span shop
| SELECT y.k FROM (SELECT x.k FROM (SELECT id AS k FROM customers) x) y
< raw
| {"columns":[{"name":"k","table":"y","orgTable":"y","orgName":"k"}],"tables":["shop.customers"]}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers) SELECT id, name FROM c
< raw
| {"columns":[{"name":"id","table":"c","orgTable":"c","orgName":"id"},{"name":"name","table":"c","orgTable":"c","orgName":"name"}],"tables":["shop.customers"]}
> raw
| span shop
| WITH c (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM c
< raw
| {"columns":[{"name":"k","table":"c","orgTable":"c","orgName":"k"},{"name":"v","table":"c","orgTable":"c","orgName":"v"}],"tables":["shop.orders"]}
> raw
| span shop
| WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT * FROM b
< raw
| {"columns":[{"name":"id","table":"b","orgTable":"b","orgName":"id"},{"name":"name","table":"b","orgTable":"b","orgName":"name"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| WITH c AS (SELECT id FROM customers) SELECT x.id, y.id FROM c x JOIN c y ON y.id = x.id
< raw
| {"columns":[{"name":"id","table":"x","orgTable":"c","orgName":"id"},{"name":"id","table":"y","orgTable":"c","orgName":"id"}],"tables":["shop.customers"]}
> raw
| span shop
| WITH RECURSIVE r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r
< raw
| {"columns":[{"name":"n","table":"r","orgTable":"r","orgName":"n"}]}
> raw
| span shop
| WITH RECURSIVE r (id, name) AS (SELECT id, name FROM customers WHERE id = 1 UNION ALL SELECT c.id, c.name FROM customers c JOIN r ON c.id = r.id + 1) SELECT id, name FROM r
< raw
| {"columns":[{"name":"id","table":"r","orgTable":"r","orgName":"id"},{"name":"name","table":"r","orgTable":"r","orgName":"name"}],"tables":["shop.customers"]}
> raw
| span shop
| WITH orders AS (SELECT id, name FROM customers) SELECT id, name FROM orders
< raw
| {"columns":[{"name":"id","table":"orders","orgTable":"orders","orgName":"id"},{"name":"name","table":"orders","orgTable":"orders","orgName":"name"}],"tables":["shop.customers"]}
> raw
| span shop
| WITH c AS (SELECT id, amount FROM orders) SELECT id FROM c WHERE amount > (SELECT avg(amount) FROM c)
< raw
| {"columns":[{"name":"id","table":"c","orgTable":"c","orgName":"id"}],"tables":["shop.orders"]}
> raw
| span shop
| WITH RECURSIVE c AS (SELECT id FROM orders), orders AS (SELECT 1 AS id) SELECT id FROM c
< raw
| {"columns":[{"name":"id","table":"c","orgTable":"c","orgName":"id"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id FROM orders UNION SELECT id FROM customers
< raw
| {"columns":[{"name":"id"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id, status FROM orders UNION ALL SELECT id, name FROM customers
< raw
| {"columns":[{"name":"id"},{"name":"status"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name FROM customers INTERSECT SELECT status FROM orders
< raw
| {"columns":[{"name":"name"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM orders EXCEPT SELECT customer_id FROM orders
< raw
| {"columns":[{"name":"id"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id AS k FROM orders UNION SELECT 1 ORDER BY k
< raw
| {"columns":[{"name":"k"}],"tables":["shop.orders"]}
> raw
| span shop
| (SELECT id FROM orders LIMIT 1) UNION (SELECT id FROM customers LIMIT 1)
< raw
| {"columns":[{"name":"id"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT u.id FROM (SELECT id FROM orders UNION SELECT id FROM customers) u
< raw
| {"columns":[{"name":"id","table":"u","orgTable":"u","orgName":"id"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id, doubled FROM big_orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"id"},{"name":"doubled","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"doubled"}],"tables":["shop.big_orders"]}
> raw
| span shop
| SELECT b.status, b.doubled FROM big_orders b
< raw
| {"columns":[{"name":"status","schema":"shop","table":"b","orgTable":"big_orders","orgName":"status"},{"name":"doubled","schema":"shop","table":"b","orgTable":"big_orders","orgName":"doubled"}],"tables":["shop.big_orders"]}
> raw
| span shop
| SELECT * FROM big_orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"id"},{"name":"doubled","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"doubled"},{"name":"status","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"status"},{"name":"customer_id","schema":"shop","table":"big_orders","orgTable":"big_orders","orgName":"customer_id"}],"tables":["shop.big_orders"]}
> raw
| span shop
| SELECT customer_id, total, n FROM customer_totals
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"customer_totals","orgTable":"customer_totals","orgName":"customer_id"},{"name":"total","schema":"shop","table":"customer_totals","orgTable":"customer_totals","orgName":"total"},{"name":"n","schema":"shop","table":"customer_totals","orgTable":"customer_totals","orgName":"n"}],"tables":["shop.customer_totals"]}
> raw
| span shop
| SELECT * FROM order_names
< raw
| {"columns":[{"name":"id","schema":"shop","table":"order_names","orgTable":"order_names","orgName":"id"},{"name":"doubled","schema":"shop","table":"order_names","orgTable":"order_names","orgName":"doubled"},{"name":"name","schema":"shop","table":"order_names","orgTable":"order_names","orgName":"name"}],"tables":["shop.order_names"]}
> raw
| span shop
| SELECT o.name, c.region FROM order_names o JOIN customers c ON c.name = o.name
< raw
| {"columns":[{"name":"name","schema":"shop","table":"o","orgTable":"order_names","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"}],"tables":["shop.customers","shop.order_names"]}
> raw
| span shop
| SELECT a, b FROM renamed
< raw
| {"columns":[{"name":"a","schema":"shop","table":"renamed","orgTable":"renamed","orgName":"a"},{"name":"b","schema":"shop","table":"renamed","orgTable":"renamed","orgName":"b"}],"tables":["shop.renamed"]}
> raw
| span shop
| SELECT id FROM all_ids
< raw
| {"columns":[{"name":"id","schema":"shop","table":"all_ids","orgTable":"all_ids","orgName":"id"}],"tables":["shop.all_ids"]}
> raw
| span shop
| SELECT * FROM staff_pay
< raw
| {"columns":[{"name":"id","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"id"},{"name":"salary","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"salary"},{"name":"dept","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"dept"}],"tables":["shop.staff_pay"]}
> raw
| span shop
| SELECT name FROM hr.rich
< raw
| {"columns":[{"name":"name","schema":"hr","table":"rich","orgTable":"rich","orgName":"name"}],"tables":["hr.rich"]}
> raw
| span shop
| SELECT VipName FROM hr.`Rich_Names`
< raw
| {"columns":[{"name":"VipName","schema":"hr","table":"rich_names","orgTable":"rich_names","orgName":"VipName"}],"tables":["hr.rich_names"]}
> raw
| span shop
| SELECT b.id, t.total FROM big_orders b JOIN customer_totals t ON t.customer_id = b.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"b","orgTable":"big_orders","orgName":"id"},{"name":"total","schema":"shop","table":"t","orgTable":"customer_totals","orgName":"total"}],"tables":["shop.big_orders","shop.customer_totals"]}
> raw
| span shop
| SELECT d.doubled FROM (SELECT doubled FROM big_orders) d
< raw
| {"columns":[{"name":"doubled","table":"d","orgTable":"d","orgName":"doubled"}],"tables":["shop.big_orders"]}
> raw
| span shop
| WITH v AS (SELECT * FROM big_orders) SELECT doubled FROM v
< raw
| {"columns":[{"name":"doubled","table":"v","orgTable":"v","orgName":"doubled"}],"tables":["shop.big_orders"]}
> raw
| span shop
| SELECT id, row_number() OVER (ORDER BY amount) AS rn FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"rn"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, sum(amount) OVER (PARTITION BY customer_id) AS running FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"running"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, rank() OVER w AS r FROM orders WINDOW w AS (PARTITION BY status ORDER BY amount)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"r"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, lag(amount) OVER (ORDER BY id) AS prev FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"prev"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT customer_id, sum(amount), count(*) FROM orders GROUP BY customer_id
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"sum(amount)"},{"name":"count(*)"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT status, avg(amount) AS a FROM orders GROUP BY status HAVING count(*) > 1
< raw
| {"columns":[{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"a"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT region, group_concat(name ORDER BY name) AS names FROM customers GROUP BY region
< raw
| {"columns":[{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"},{"name":"names"}],"tables":["shop.customers"]}
> raw
| span shop
| SELECT max(amount) FROM orders WHERE status = 'open'
< raw
| {"columns":[{"name":"max(amount)"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT count(DISTINCT customer_id) AS buyers FROM orders
< raw
| {"columns":[{"name":"buyers"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT c.region, sum(o.amount) FROM customers c JOIN orders o ON o.customer_id = c.id GROUP BY c.region
< raw
| {"columns":[{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"sum(o.amount)"}],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT customer_id, sum(amount) FROM orders GROUP BY customer_id WITH ROLLUP
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"sum(amount)"}],"tables":["shop.orders"]}
> raw
| span shop
| SELECT nope FROM orders
< raw
| {"error":"1054"}
> raw
| span shop
| SELECT id FROM no_such_table
< raw
| {"error":"1146"}
> raw
| span shop
| SELECT id FROM orders JOIN customers ON customers.id = orders.customer_id
< raw
| {"error":"1052"}
> raw
| span hr
| SELECT name, salary FROM staff
< raw
| {"columns":[{"name":"name","schema":"hr","table":"staff","orgTable":"staff","orgName":"name"},{"name":"salary","schema":"hr","table":"staff","orgTable":"staff","orgName":"salary"}],"tables":["hr.staff"]}
> raw
| span hr
| SELECT s.name, d.name AS dept FROM staff s JOIN depts d ON d.id = s.dept_id
< raw
| {"columns":[{"name":"name","schema":"hr","table":"s","orgTable":"staff","orgName":"name"},{"name":"dept","schema":"hr","table":"d","orgTable":"depts","orgName":"name"}],"tables":["hr.depts","hr.staff"]}
> raw
| span hr
| SELECT o.amount FROM shop.orders o
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"tables":["shop.orders"]}
> raw
| span hr
| SELECT * FROM rich
< raw
| {"columns":[{"name":"name","schema":"hr","table":"rich","orgTable":"rich","orgName":"name"},{"name":"salary","schema":"hr","table":"rich","orgTable":"rich","orgName":"salary"}],"tables":["hr.rich"]}
> raw
| span hr
| SELECT * FROM shop.staff_pay
< raw
| {"columns":[{"name":"id","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"id"},{"name":"salary","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"salary"},{"name":"dept","schema":"shop","table":"staff_pay","orgTable":"staff_pay","orgName":"dept"}],"tables":["shop.staff_pay"]}
> raw
| span hr
| SELECT name FROM staff UNION SELECT name FROM shop.customers
< raw
| {"columns":[{"name":"name"}],"tables":["hr.staff","shop.customers"]}
