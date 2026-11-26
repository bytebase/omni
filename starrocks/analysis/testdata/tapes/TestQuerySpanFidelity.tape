-- pin: default starrocks/allin1-ubuntu:3.4.10 sha256:8bc803ea7a5bb490c77117a456b32b792ccae698b1f0e47dfcafdac41542566d 3.4.10

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
| CREATE MATERIALIZED VIEW shop.order_sums DISTRIBUTED BY HASH(customer_id) BUCKETS 1 REFRESH MANUAL PROPERTIES ("replication_num" = "1") AS SELECT customer_id, SUM(amount) AS total FROM shop.orders GROUP BY customer_id;
| INSERT INTO shop.customers VALUES (1, 'ann', 'east', 'a@x'), (2, 'bob', 'west', 'b@x'), (3, 'cy', 'east', 'c@x');
| INSERT INTO shop.orders VALUES (1, 1, 10.00, 'open', '2026-01-01'), (2, 2, 20.00, 'paid', '2026-01-02'), (3, 1, 5.00, 'open', '2026-01-03');
| INSERT INTO shop.items VALUES (1, 'a', 1), (1, 'b', 2), (2, 'a', 3);
| INSERT INTO shop.`MixedCase` VALUES (1, 'x', 'y'), (2, 'z', 'w');
| INSERT INTO shop.mixedcase VALUES (1, 'n'), (2, 'm');
| INSERT INTO hr.staff VALUES (1, 'dee', 100, 1), (2, 'eve', 200, 2), (3, 'fay', 150, 1);
| INSERT INTO hr.depts VALUES (1, 'ops'), (2, 'dev');
< raw
| {"databases":[{"name":"hr","tables":[{"name":"depts","columns":["id","name"]},{"name":"staff","columns":["id","name","salary","dept_id"]}],"views":[{"name":"Rich_Names","columns":["VipName"],"definition":"SELECT `hr`.`rich`.`name` AS `VipName`\nFROM `hr`.`rich`"},{"name":"rich","columns":["name","salary"],"definition":"SELECT `hr`.`staff`.`name`, `hr`.`staff`.`salary`\nFROM `hr`.`staff`\nWHERE `hr`.`staff`.`salary` \u003e 120"}]},{"name":"shop","tables":[{"name":"MixedCase","columns":["Id","CamelCol","lower_col"]},{"name":"customers","columns":["id","name","region","email"]},{"name":"items","columns":["order_id","sku","qty"]},{"name":"mixedcase","columns":["id","note"]},{"name":"orders","columns":["id","customer_id","amount","status","created_at"]}],"views":[{"name":"all_ids","columns":["id"],"definition":"SELECT `shop`.`orders`.`id`\nFROM `shop`.`orders` UNION SELECT `shop`.`customers`.`id`\nFROM `shop`.`customers`"},{"name":"big_orders","columns":["id","doubled","status","customer_id"],"definition":"SELECT `shop`.`orders`.`id`, `shop`.`orders`.`amount` * 2 AS `doubled`, `shop`.`orders`.`status`, `shop`.`orders`.`customer_id`\nFROM `shop`.`orders`\nWHERE `shop`.`orders`.`amount` \u003e 5"},{"name":"customer_totals","columns":["customer_id","total","n"],"definition":"SELECT `shop`.`orders`.`customer_id`, sum(`shop`.`orders`.`amount`) AS `total`, count(*) AS `n`\nFROM `shop`.`orders`\nGROUP BY `shop`.`orders`.`customer_id`"},{"name":"order_names","columns":["id","doubled","name"],"definition":"SELECT `shop`.`b`.`id`, `shop`.`b`.`doubled`, `shop`.`c`.`name`\nFROM `shop`.`big_orders` AS `b` INNER JOIN `shop`.`customers` AS `c` ON `shop`.`c`.`id` = `shop`.`b`.`customer_id`"},{"name":"renamed","columns":["a","b"],"definition":"SELECT `shop`.`customers`.`id`, `shop`.`customers`.`name`\nFROM `shop`.`customers`"},{"name":"staff_pay","columns":["id","salary","dept"],"definition":"SELECT `hr`.`s`.`id`, `hr`.`s`.`salary`, `hr`.`d`.`name` AS `dept`\nFROM `hr`.`staff` AS `s` INNER JOIN `hr`.`depts` AS `d` ON `hr`.`d`.`id` = `hr`.`s`.`dept_id`"}],"materializedViews":[{"name":"order_sums","columns":["customer_id","total"],"definition":"CREATE MATERIALIZED VIEW `order_sums` (`customer_id`, `total`)\nDISTRIBUTED BY HASH(`customer_id`) BUCKETS 1 \nREFRESH MANUAL\nPROPERTIES (\n\"replicated_storage\" = \"true\",\n\"replication_num\" = \"1\",\n\"storage_medium\" = \"HDD\"\n)\nAS SELECT `orders`.`customer_id`, sum(`orders`.`amount`) AS `total`\nFROM `shop`.`orders`\nGROUP BY `orders`.`customer_id`;"}]}]}
> raw
| span shop
| SELECT id, amount FROM orders
< raw
| {"columns":["id","amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 3:amount]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT o.id AS order_id, o.amount FROM orders o
< raw
| {"columns":["order_id","amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 3:amount]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT id + customer_id, upper(status) FROM orders
< raw
| {"columns":["id + customer_id","upper(status)"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:expr, 7:upper]","    - SCAN [orders] =\u003e [6:expr, 7:upper]","            6:expr := cast(1:id as bigint(20)) + cast(2:customer_id as bigint(20))","            7:upper := upper(4:status)"]}
> raw
| span shop
| SELECT CAST(id AS VARCHAR) AS a, concat(name, '@', region) AS b FROM customers
< raw
| {"columns":["a","b"],"tables":["shop.customers"],"plan":["- Output =\u003e [5:cast, 6:concat]","    - SCAN [customers] =\u003e [5:cast, 6:concat]","            5:cast := cast(1:id as varchar)","            6:concat := concat(2:name, '@', 3:region)"]}
> raw
| span shop
| SELECT 1, 'a' AS b
< raw
| {"columns":["1","b"],"plan":["EXECUTE IN FE","- Output =\u003e [2:expr, 3:expr]","    - VALUES =\u003e [2:expr, 3:expr]","            2:expr := 1","            3:expr := 'a'"]}
> raw
| span shop
| SELECT count(*) FROM orders
< raw
| {"columns":["count(*)"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:count]","    - AGGREGATE(GLOBAL) []","            6:count := count()","        - SCAN [orders] =\u003e [8:auto_fill_col]","                8:auto_fill_col := 1"]}
> raw
| span shop
| SELECT DISTINCT region FROM customers
< raw
| {"columns":["region"],"tables":["shop.customers"],"plan":["- Output =\u003e [3:region]","    - AGGREGATE(GLOBAL) [3:region]","        - SCAN [customers] =\u003e [3:region]"]}
> raw
| span shop
| SELECT name FROM customers ORDER BY region LIMIT 2
< raw
| {"columns":["name"],"tables":["shop.customers"],"plan":["- Output =\u003e [2:name]","    - TOP-2(FINAL)[3: region ASC NULLS FIRST] =\u003e [2:name]","        - TOP-2(PARTIAL)[3: region ASC NULLS FIRST]","            - SCAN [customers] =\u003e [2:name, 3:region]"]}
> raw
| span shop
| SELECT CASE WHEN amount > 5 THEN status ELSE 'small' END AS label FROM orders
< raw
| {"columns":["label"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:case]","    - SCAN [orders] =\u003e [6:case]","            6:case := if(3:amount \u003e 5, 4:status, 'small')"]}
> raw
| span shop
| SELECT coalesce(email, name) AS contact FROM customers
< raw
| {"columns":["contact"],"tables":["shop.customers"],"plan":["- Output =\u003e [5:coalesce]","    - SCAN [customers] =\u003e [5:coalesce]","            5:coalesce := coalesce(4:email, 2:name)"]}
> raw
| span shop
| SELECT amount AS id, id AS amount FROM orders
< raw
| {"columns":["id","amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [3:amount, 1:id]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT id   +   1, -amount, id IN (1, 2), NOT id FROM orders
< raw
| {"columns":["id + 1","-1 * amount","id IN (1, 2)","NOT id"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:expr, 7:expr, 8:expr, 9:expr]","    - SCAN [orders] =\u003e [6:expr, 7:expr, 8:expr, 9:expr]","            6:expr := cast(1:id as bigint(20)) + 1","            7:expr := -1 * cast(3:amount as decimal(10, 2))","            8:expr := 1:id IN (1, 2)","            9:expr := NOT cast(1:id as boolean)"]}
> raw
| span shop
| SELECT 'x', 2.50, NULL, TRUE, "y"
< raw
| {"columns":["'x'","2.50","NULL","TRUE","'y'"],"plan":["EXECUTE IN FE","- Output =\u003e [2:expr, 3:expr, 4:expr, 5:expr, 6:expr]","    - VALUES =\u003e [2:expr, 3:expr, 4:expr, 5:expr, 6:expr]","            2:expr := 'x'","            3:expr := 2.50","            4:expr := null","            5:expr := true","            6:expr := 'y'"]}
> raw
| span shop
| SELECT (id), ((amount)) FROM orders
< raw
| {"columns":["id","amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 3:amount]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| VALUES (1, 2 + 3), (4, 5)
< raw
| {"error":"1064 Getting syntax error at line 1, column 0. Detail message: Unexpected input 'VALUES', the most similar input is {'PAUSE', 'ALTER', 'ANALYZE', 'USE', 'DESC', '(', ';'}."}
> raw
| span shop
| SELECT * FROM orders
< raw
| {"columns":["id","customer_id","amount","status","created_at"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at]","    - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at]"]}
> raw
| span shop
| SELECT o.* FROM orders o
< raw
| {"columns":["id","customer_id","amount","status","created_at"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at]","    - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at]"]}
> raw
| span shop
| SELECT *, amount FROM orders
< raw
| {"columns":["id","customer_id","amount","status","created_at","amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at, 3:amount]","    - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at]"]}
> raw
| span shop
| SELECT c.*, o.amount FROM customers c JOIN orders o ON o.customer_id = c.id
< raw
| {"columns":["id","name","region","email","amount"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 2:name, 3:region, 4:email, 7:amount]","    - HASH/INNER JOIN [1:id = 6:customer_id] =\u003e [1:id, 2:name, 3:region, 4:email, 7:amount]","        - SCAN [customers] =\u003e [1:id, 2:name, 3:region, 4:email]","                predicate: 1:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [orders] =\u003e [6:customer_id, 7:amount]","                    predicate: 6:customer_id IS NOT NULL"]}
> raw
| span shop
| SELECT * FROM (SELECT * FROM customers) c
< raw
| {"columns":["id","name","region","email"],"tables":["shop.customers"],"plan":["- Output =\u003e [1:id, 2:name, 3:region, 4:email]","    - SCAN [customers] =\u003e [1:id, 2:name, 3:region, 4:email]"]}
> raw
| span shop
| SELECT * EXCEPT (amount, status) FROM orders
< raw
| {"error":"1064 Getting syntax error at line 1, column 17. Detail message: Unexpected input 'amount', the most similar input is {'WITH', 'SELECT', '('}."}
> raw
| span shop
| SELECT id FROM shop.orders
< raw
| {"columns":["id"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id]","    - SCAN [orders] =\u003e [1:id]"]}
> raw
| span shop
| SELECT shop.orders.amount, orders.status FROM shop.orders
< raw
| {"columns":["amount","status"],"tables":["shop.orders"],"plan":["- Output =\u003e [3:amount, 4:status]","    - SCAN [orders] =\u003e [3:amount, 4:status]"]}
> raw
| span shop
| SELECT s.name, s.salary FROM hr.staff s
< raw
| {"columns":["name","salary"],"tables":["hr.staff"],"plan":["- Output =\u003e [2:name, 3:salary]","    - SCAN [staff] =\u003e [2:name, 3:salary]"]}
> raw
| span shop
| SELECT o.id, s.name FROM orders o JOIN hr.staff s ON s.id = o.customer_id
< raw
| {"columns":["id","name"],"tables":["hr.staff","shop.orders"],"plan":["- Output =\u003e [1:id, 7:name]","    - HASH/INNER JOIN [6:id = 2:customer_id] =\u003e [1:id, 7:name]","        - SCAN [staff] =\u003e [6:id, 7:name]","                predicate: 6:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","                    predicate: 2:customer_id IS NOT NULL"]}
> raw
| span shop
| SELECT id, amount FROM default_catalog.shop.orders
< raw
| {"columns":["id","amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 3:amount]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT o.id FROM default_catalog.shop.orders o JOIN default_catalog.hr.staff s ON s.id = o.customer_id
< raw
| {"columns":["id"],"tables":["hr.staff","shop.orders"],"plan":["- Output =\u003e [1:id]","    - HASH/INNER JOIN [2:customer_id = 6:id] =\u003e [1:id]","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","                    predicate: 2:customer_id IS NOT NULL","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [staff] =\u003e [6:id]","                    predicate: 6:id IS NOT NULL"]}
> raw
| span shop
| SELECT default_catalog.shop.orders.amount FROM default_catalog.shop.orders
< raw
| {"columns":["amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [3:amount]","    - SCAN [orders] =\u003e [3:amount]"]}
> raw
| span shop
| SELECT * FROM default_catalog.hr.rich
< raw
| {"columns":["name","salary"],"tables":["hr.staff"],"views":["hr.rich"],"plan":["- Output =\u003e [2:name, 3:salary]","    - SCAN [staff] =\u003e [2:name, 3:salary]","            predicate: 3:salary \u003e 120"]}
> raw
| span shop
| SELECT shop.o.id, default_catalog.shop.o.amount FROM orders o
< raw
| {"columns":["id","amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 3:amount]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT ID, Amount FROM orders
< raw
| {"columns":["ID","Amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 3:amount]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT `id`, `AMOUNT` FROM `orders`
< raw
| {"columns":["id","AMOUNT"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 3:amount]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT O.ID FROM orders O
< raw
| {"columns":["ID"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id]","    - SCAN [orders] =\u003e [1:id]"]}
> raw
| span shop
| SELECT `Id`, `CamelCol`, lower_col FROM `MixedCase`
< raw
| {"columns":["Id","CamelCol","lower_col"],"tables":["shop.MixedCase"],"plan":["- Output =\u003e [1:Id, 2:CamelCol, 3:lower_col]","    - SCAN [MixedCase] =\u003e [1:Id, 2:CamelCol, 3:lower_col]"]}
> raw
| span shop
| SELECT * FROM `MixedCase`
< raw
| {"columns":["Id","CamelCol","lower_col"],"tables":["shop.MixedCase"],"plan":["- Output =\u003e [1:Id, 2:CamelCol, 3:lower_col]","    - SCAN [MixedCase] =\u003e [1:Id, 2:CamelCol, 3:lower_col]"]}
> raw
| span shop
| SELECT * FROM mixedcase
< raw
| {"columns":["id","note"],"tables":["shop.mixedcase"],"plan":["- Output =\u003e [1:id, 2:note]","    - SCAN [mixedcase] =\u003e [1:id, 2:note]"]}
> raw
| span shop
| SELECT m.camelcol, m.LOWER_COL FROM `MixedCase` m
< raw
| {"columns":["camelcol","LOWER_COL"],"tables":["shop.MixedCase"],"plan":["- Output =\u003e [2:CamelCol, 3:lower_col]","    - SCAN [MixedCase] =\u003e [2:CamelCol, 3:lower_col]"]}
> raw
| span shop
| SELECT id FROM Orders
< raw
| {"error":"5502 Getting analyzing error. Detail message: Unknown table 'shop.Orders'."}
> raw
| span shop
| SELECT orders.name FROM customers AS orders
< raw
| {"columns":["name"],"tables":["shop.customers"],"plan":["- Output =\u003e [2:name]","    - SCAN [customers] =\u003e [2:name]"]}
> raw
| span shop
| SELECT customers.amount FROM orders customers
< raw
| {"columns":["amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [3:amount]","    - SCAN [orders] =\u003e [3:amount]"]}
> raw
| span shop
| SELECT o.id FROM customers o JOIN orders c ON c.customer_id = o.id
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id]","    - HASH/INNER JOIN [1:id = 6:customer_id] =\u003e [1:id]","        - SCAN [customers] =\u003e [1:id]","                predicate: 1:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [orders] =\u003e [6:customer_id]","                    predicate: 6:customer_id IS NOT NULL"]}
> raw
| span shop
| SELECT * FROM (SELECT amount AS id FROM orders) orders
< raw
| {"columns":["id"],"tables":["shop.orders"],"plan":["- Output =\u003e [3:amount]","    - SCAN [orders] =\u003e [3:amount]"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["id","name"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 7:name]","    - HASH/INNER JOIN [6:id = 2:customer_id] =\u003e [1:id, 7:name]","        - SCAN [customers] =\u003e [6:id, 7:name]","                predicate: 6:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","                    predicate: 2:customer_id IS NOT NULL"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["id","name"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 7:name]","    - HASH/RIGHT OUTER JOIN [6:id = 2:customer_id] =\u003e [1:id, 7:name]","        - SCAN [customers] =\u003e [6:id, 7:name]","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]"]}
> raw
| span shop
| SELECT o.id, c.name FROM customers c RIGHT JOIN orders o ON c.id = o.customer_id
< raw
| {"columns":["id","name"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [5:id, 2:name]","    - HASH/RIGHT OUTER JOIN [1:id = 6:customer_id] =\u003e [2:name, 5:id]","        - SCAN [customers] =\u003e [1:id, 2:name]","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [orders] =\u003e [5:id, 6:customer_id]"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o FULL OUTER JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["id","name"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 7:name]","    - HASH/FULL OUTER JOIN [6:id = 2:customer_id] =\u003e [1:id, 7:name]","        - SCAN [customers] =\u003e [6:id, 7:name]","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o CROSS JOIN customers c
< raw
| {"columns":["id","name"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 7:name]","    - NESTLOOP/CROSS JOIN","        - SCAN [customers] =\u003e [7:name]","        - EXCHANGE(BROADCAST)","            - SCAN [orders] =\u003e [1:id]"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id
< raw
| {"columns":["id","name"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 7:name]","    - HASH/INNER JOIN [6:id = 2:customer_id] =\u003e [1:id, 7:name]","        - SCAN [customers] =\u003e [6:id, 7:name]","                predicate: 6:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","                    predicate: 2:customer_id IS NOT NULL"]}
> raw
| span shop
| SELECT o.id FROM orders o LEFT SEMI JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id]","    - HASH/LEFT SEMI JOIN [2:customer_id = 6:id] =\u003e [1:id]","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [customers] =\u003e [6:id]","                    predicate: 6:id IS NOT NULL"]}
> raw
| span shop
| SELECT o.id FROM orders o LEFT ANTI JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id]","    - HASH/LEFT ANTI JOIN [2:customer_id = 6:id] =\u003e [1:id]","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [customers] =\u003e [6:id]"]}
> raw
| span shop
| SELECT o.id FROM orders o JOIN customers c USING (id)
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id]","    - HASH/INNER JOIN [1:id = 6:id] =\u003e [1:id]","        - SCAN [orders] =\u003e [1:id]","                predicate: 1:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [customers] =\u003e [6:id]","                    predicate: 6:id IS NOT NULL"]}
> raw
| span shop
| SELECT id, name FROM orders JOIN customers USING (id)
< raw
| {"error":"1064 Getting analyzing error. Detail message: Column 'id' is ambiguous."}
> raw
| span shop
| SELECT * FROM orders o JOIN customers c USING (id)
< raw
| {"columns":["id","customer_id","amount","status","created_at","id","name","region","email"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at, 6:id, 7:name, 8:region, 9:email]","    - HASH/INNER JOIN [6:id = 1:id]","        - SCAN [customers] =\u003e [6:id, 7:name, 8:region, 9:email]","                predicate: 6:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [1]","            - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at]","                    predicate: 1:id IS NOT NULL"]}
> raw
| span shop
| SELECT * FROM orders o LEFT JOIN customers c USING (id)
< raw
| {"columns":["id","customer_id","amount","status","created_at","id","name","region","email"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at, 6:id, 7:name, 8:region, 9:email]","    - HASH/RIGHT OUTER JOIN [6:id = 1:id]","        - SCAN [customers] =\u003e [6:id, 7:name, 8:region, 9:email]","        - EXCHANGE(SHUFFLE) [1]","            - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at]"]}
> raw
| span shop
| SELECT * FROM orders NATURAL JOIN customers
< raw
| {"columns":["id","customer_id","amount","status","created_at","id","name","region","email"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at, 6:id, 7:name, 8:region, 9:email]","    - NESTLOOP/CROSS JOIN","        - SCAN [customers] =\u003e [6:id, 7:name, 8:region, 9:email]","        - EXCHANGE(BROADCAST)","            - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at]"]}
> raw
| span shop
| SELECT i.sku, o.status, c.region FROM items i JOIN orders o ON o.id = i.order_id JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["sku","status","region"],"tables":["shop.customers","shop.items","shop.orders"],"plan":["- Output =\u003e [2:sku, 7:status, 11:region]","    - HASH/INNER JOIN [5:customer_id = 9:id] =\u003e [2:sku, 11:region, 7:status]","        - EXCHANGE(SHUFFLE) [5]","            - HASH/INNER JOIN [4:id = 1:order_id] =\u003e [2:sku, 5:customer_id, 7:status]","                - SCAN [orders] =\u003e [4:id, 5:customer_id, 7:status]","                        predicate: 5:customer_id IS NOT NULL","                - EXCHANGE(SHUFFLE) [1]","                    - SCAN [items] =\u003e [1:order_id, 2:sku]","                            predicate: 1:order_id IS NOT NULL","        - EXCHANGE(SHUFFLE) [9]","            - SCAN [customers] =\u003e [9:id, 11:region]","                    predicate: 9:id IS NOT NULL"]}
> raw
| span shop
| SELECT a.id, b.id FROM orders a JOIN orders b ON b.customer_id = a.customer_id
< raw
| {"columns":["id","id"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 6:id]","    - HASH/INNER JOIN [2:customer_id = 7:customer_id] =\u003e [1:id, 6:id]","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","                    predicate: 2:customer_id IS NOT NULL","        - EXCHANGE(SHUFFLE) [7]","            - SCAN [orders] =\u003e [6:id, 7:customer_id]","                    predicate: 7:customer_id IS NOT NULL"]}
> raw
| span shop
| SELECT o.id FROM orders o JOIN [shuffle] customers c ON c.id = o.customer_id
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id]","    - HASH/INNER JOIN [2:customer_id = 6:id] =\u003e [1:id]","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","                    predicate: 2:customer_id IS NOT NULL","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [customers] =\u003e [6:id]","                    predicate: 6:id IS NOT NULL"]}
> raw
| span shop
| SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
< raw
| {"columns":["id","who"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 10:name]","    - HASH/RIGHT OUTER JOIN [6:id = 2:customer_id] =\u003e [1:id, 10:name, 13:subquery_assertion]","            10:name := 12:anyValue","            13:subquery_assertion := assert_true(11:countRows IS NULL OR 11:countRows \u003c= 1, 'correlate scalar subquery result must 1 row')","        - AGGREGATE(GLOBAL) [6:id]","                11:countRows := count(1)","                12:anyValue := any_value(7:name)","            - SCAN [customers] =\u003e [6:id, 7:name]","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]"]}
> raw
| span shop
| SELECT id, (SELECT max(amount) FROM orders) AS top FROM customers
< raw
| {"columns":["id","top"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 11:max]","    - NESTLOOP/CROSS JOIN =\u003e [1:id, 11:max]","            11:max := 10:max","        - SCAN [customers] =\u003e [1:id]","        - EXCHANGE(BROADCAST)","            - ASSERT LE 1","                - EXCHANGE(GATHER)","                    - AGGREGATE(GLOBAL) []","                            10:max := max(7:amount)","                        - SCAN [orders] =\u003e [7:amount]"]}
> raw
| span shop
| SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'east')
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id]","    - HASH/LEFT SEMI JOIN [2:customer_id = 6:id] =\u003e [1:id]","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [customers] =\u003e [6:id]","                    predicate: 6:id IS NOT NULL AND 8:region = 'east'"]}
> raw
| span shop
| SELECT id FROM orders o WHERE EXISTS (SELECT 1 FROM items i WHERE i.order_id = o.id)
< raw
| {"columns":["id"],"tables":["shop.items","shop.orders"],"plan":["- Output =\u003e [1:id]","    - HASH/LEFT SEMI JOIN [1:id = 6:order_id] =\u003e [1:id]","        - SCAN [orders] =\u003e [1:id]","        - EXCHANGE(SHUFFLE) [6]","            - SCAN [items] =\u003e [6:order_id]","                    predicate: 6:order_id IS NOT NULL"]}
> raw
| span shop
| SELECT d.k, d.v FROM (SELECT id AS k, amount AS v FROM orders) d
< raw
| {"columns":["k","v"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 3:amount]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT d.k FROM (SELECT id AS k FROM orders LIMIT 2) d
< raw
| {"columns":["k"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id]","    - LIMIT [2]","        - EXCHANGE(GATHER)","            - SCAN [orders] =\u003e [1:id]","                    limit: 2"]}
> raw
| span shop
| SELECT x.total FROM (SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id) x
< raw
| {"columns":["total"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:sum]","    - AGGREGATE(GLOBAL) [2:customer_id] =\u003e [6:sum]","            6:sum := sum(3:amount)","        - SCAN [orders] =\u003e [2:customer_id, 3:amount]"]}
> raw
| span shop
| SELECT y.k FROM (SELECT x.k FROM (SELECT id AS k FROM customers) x) y
< raw
| {"columns":["k"],"tables":["shop.customers"],"plan":["- Output =\u003e [1:id]","    - SCAN [customers] =\u003e [1:id]"]}
> raw
| span shop
| SELECT d.a FROM (SELECT 1 AS a, 2) d
< raw
| {"columns":["a"],"plan":["EXECUTE IN FE","- Output =\u003e [4:a]","    - VALUES =\u003e [4:a]","            4:a := 1"]}
> raw
| span shop
| SELECT o.id, u.x FROM orders o, unnest([1, 2]) AS u(x)
< raw
| {"columns":["id","x"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 6:x]","    - TABLE FUNCTION [unnest] =\u003e [{1,6}]","        - SCAN [orders] =\u003e [1:id, 7:expr]","                7:expr := 1, 2"]}
> raw
| span shop
| SELECT o.id, l.n FROM orders o, LATERAL (SELECT count(*) AS n FROM items i WHERE i.order_id = o.id) l
< raw
| {"error":"1064 Getting analyzing error. Detail message: Column '`o`.`id`' cannot be resolved."}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers) SELECT id, name FROM c
< raw
| {"columns":["id","name"],"tables":["shop.customers"],"plan":["- Output =\u003e [1:id, 2:name]","    - SCAN [customers] =\u003e [1:id, 2:name]"]}
> raw
| span shop
| WITH c (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM c
< raw
| {"columns":["k","v"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 3:amount]","    - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT * FROM b
< raw
| {"columns":["id","name"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 7:name]","    - HASH/INNER JOIN [6:id = 2:customer_id] =\u003e [1:id, 7:name]","        - SCAN [customers] =\u003e [6:id, 7:name]","                predicate: 6:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","                    predicate: 2:customer_id IS NOT NULL"]}
> raw
| span shop
| WITH c AS (SELECT id FROM customers) SELECT x.id, y.id FROM c x JOIN c y ON y.id = x.id
< raw
| {"columns":["id","id"],"tables":["shop.customers"],"plan":["- Output =\u003e [5:id, 9:id]","    - CTENoOp[1]","        - HASH/INNER JOIN [5:id = 9:id]","            - CTENoOp[1]","                - SCAN [customers] =\u003e [5:id]","                        predicate: 5:id IS NOT NULL","            - CTENoOp[1]","                - SCAN [customers] =\u003e [9:id]","                        predicate: 9:id IS NOT NULL"]}
> raw
| span shop
| WITH orders AS (SELECT id, name FROM customers) SELECT id, name FROM orders
< raw
| {"columns":["id","name"],"tables":["shop.customers"],"plan":["- Output =\u003e [1:id, 2:name]","    - SCAN [customers] =\u003e [1:id, 2:name]"]}
> raw
| span shop
| WITH c AS (SELECT id, amount FROM orders) SELECT id FROM c WHERE amount > (SELECT avg(amount) FROM c)
< raw
| {"columns":["id"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:id]","    - CTENoOp[1]","        - NESTLOOP/INNER JOIN [cast(8:amount as decimal(38, 8)) \u003e 16:avg] =\u003e [6:id]","            - CTENoOp[1]","                - SCAN [orders] =\u003e [6:id, 8:amount]","            - EXCHANGE(BROADCAST)","                - ASSERT LE 1","                    - AGGREGATE(GLOBAL) []","                            16:avg := avg(16:avg)","                        - EXCHANGE(GATHER)","                            - AGGREGATE(LOCAL) []","                                    16:avg := avg(13:amount)","                                - CTENoOp[1]","                                    - SCAN [orders] =\u003e [13:amount]"]}
> raw
| span shop
| WITH RECURSIVE r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r
< raw
| {"error":"1064 Getting syntax error at line 1, column 15. Detail message: No viable statement for input 'WITH RECURSIVE r'."}
> raw
| span shop
| SELECT id FROM orders UNION SELECT id FROM customers
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [10:id]","    - AGGREGATE(GLOBAL) [10:id]","        - EXCHANGE(SHUFFLE) [10]","            - AGGREGATE(LOCAL) [10:id]","                - UNION =\u003e [10: id]","                    - SCAN [orders] =\u003e [1:id]","                    - SCAN [customers] =\u003e [6:id]"]}
> raw
| span shop
| SELECT id, status FROM orders UNION ALL SELECT id, name FROM customers
< raw
| {"columns":["id","status"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [11:id, 12:cast]","    - UNION =\u003e [11: id, 12: cast]","        - SCAN [orders] =\u003e [1:id, 6:cast]","                6:cast := 4:status","        - SCAN [customers] =\u003e [7:id, 8:name]"]}
> raw
| span shop
| SELECT name FROM customers INTERSECT SELECT status FROM orders
< raw
| {"columns":["name"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [11:name]","    - INTERSECT =\u003e [11: name]","        - SCAN [customers] =\u003e [2:name]","        - SCAN [orders] =\u003e [10:cast]","                10:cast := 8:status"]}
> raw
| span shop
| SELECT id FROM orders EXCEPT SELECT customer_id FROM orders
< raw
| {"columns":["id"],"tables":["shop.orders"],"plan":["- Output =\u003e [11:id]","    - EXCEPT =\u003e [11: id]","        - SCAN [orders] =\u003e [1:id]","        - SCAN [orders] =\u003e [7:customer_id]"]}
> raw
| span shop
| SELECT id FROM orders MINUS SELECT customer_id FROM orders
< raw
| {"columns":["id"],"tables":["shop.orders"],"plan":["- Output =\u003e [11:id]","    - EXCEPT =\u003e [11: id]","        - SCAN [orders] =\u003e [1:id]","        - SCAN [orders] =\u003e [7:customer_id]"]}
> raw
| span shop
| SELECT id AS k FROM orders UNION SELECT 1 ORDER BY k
< raw
| {"columns":["k"],"tables":["shop.orders"],"plan":["- Output =\u003e [9:id]","    - SORT(FINAL)[9: id ASC NULLS FIRST]","        - SORT(PARTIAL)[9: id ASC NULLS FIRST]","            - AGGREGATE(GLOBAL) [9:id]","                - EXCHANGE(SHUFFLE) [9]","                    - AGGREGATE(LOCAL) [9:id]","                        - UNION =\u003e [9: id]","                            - SCAN [orders] =\u003e [1:id]","                            - VALUES =\u003e [8:cast]","                                    8:cast := 1"]}
> raw
| span shop
| (SELECT id FROM orders LIMIT 1) UNION (SELECT id FROM customers LIMIT 1)
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [10:id]","    - AGGREGATE(GLOBAL) [10:id]","        - EXCHANGE(SHUFFLE) [10]","            - AGGREGATE(LOCAL) [10:id]","                - UNION =\u003e [10: id]","                    - LIMIT [1]","                        - EXCHANGE(GATHER)","                            - SCAN [orders] =\u003e [1:id]","                                    limit: 1","                    - LIMIT [1]","                        - EXCHANGE(GATHER)","                            - SCAN [customers] =\u003e [6:id]","                                    limit: 1"]}
> raw
| span shop
| SELECT u.id FROM (SELECT id FROM orders UNION SELECT id FROM customers) u
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [10:id]","    - AGGREGATE(GLOBAL) [10:id]","        - EXCHANGE(SHUFFLE) [10]","            - AGGREGATE(LOCAL) [10:id]","                - UNION =\u003e [10: id]","                    - SCAN [orders] =\u003e [1:id]","                    - SCAN [customers] =\u003e [6:id]"]}
> raw
| span shop
| SELECT id, doubled FROM big_orders
< raw
| {"columns":["id","doubled"],"tables":["shop.orders"],"views":["shop.big_orders"],"plan":["- Output =\u003e [1:id, 6:expr]","    - SCAN [orders] =\u003e [1:id, 6:expr]","            6:expr := cast(3:amount as decimal(10, 2)) * 2","            predicate: 3:amount \u003e 5"]}
> raw
| span shop
| SELECT b.status, b.doubled FROM big_orders b
< raw
| {"columns":["status","doubled"],"tables":["shop.orders"],"views":["shop.big_orders"],"plan":["- Output =\u003e [4:status, 6:expr]","    - SCAN [orders] =\u003e [4:status, 6:expr]","            6:expr := cast(3:amount as decimal(10, 2)) * 2","            predicate: 3:amount \u003e 5"]}
> raw
| span shop
| SELECT * FROM big_orders
< raw
| {"columns":["id","doubled","status","customer_id"],"tables":["shop.orders"],"views":["shop.big_orders"],"plan":["- Output =\u003e [1:id, 6:expr, 4:status, 2:customer_id]","    - SCAN [orders] =\u003e [1:id, 2:customer_id, 4:status, 6:expr]","            6:expr := cast(3:amount as decimal(10, 2)) * 2","            predicate: 3:amount \u003e 5"]}
> raw
| span shop
| SELECT customer_id, total, n FROM customer_totals
< raw
| {"columns":["customer_id","total","n"],"tables":["shop.orders"],"views":["shop.customer_totals"],"plan":["- Output =\u003e [2:customer_id, 6:sum, 7:count]","    - AGGREGATE(GLOBAL) [2:customer_id]","            6:sum := sum(3:amount)","            7:count := count()","        - SCAN [orders] =\u003e [2:customer_id, 3:amount]"]}
> raw
| span shop
| SELECT * FROM order_names
< raw
| {"columns":["id","doubled","name"],"tables":["shop.customers","shop.orders"],"views":["shop.big_orders","shop.order_names"],"plan":["- Output =\u003e [1:id, 6:expr, 8:name]","    - HASH/INNER JOIN [7:id = 2:customer_id] =\u003e [1:id, 6:expr, 8:name]","        - SCAN [customers] =\u003e [7:id, 8:name]","                predicate: 7:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id, 6:expr]","                    6:expr := cast(3:amount as decimal(10, 2)) * 2","                    predicate: 2:customer_id IS NOT NULL AND 3:amount \u003e 5"]}
> raw
| span shop
| SELECT o.name, c.region FROM order_names o JOIN customers c ON c.name = o.name
< raw
| {"columns":["name","region"],"tables":["shop.customers","shop.orders"],"views":["shop.big_orders","shop.order_names"],"plan":["- Output =\u003e [8:name, 13:region]","    - HASH/INNER JOIN [12:name = 8:name] =\u003e [13:region, 8:name]","        - EXCHANGE(SHUFFLE) [12]","            - SCAN [customers] =\u003e [12:name, 13:region]","                    predicate: 12:name IS NOT NULL","        - EXCHANGE(SHUFFLE) [8]","            - HASH/INNER JOIN [7:id = 2:customer_id] =\u003e [8:name]","                - SCAN [customers] =\u003e [7:id, 8:name]","                        predicate: 8:name IS NOT NULL","                - EXCHANGE(SHUFFLE) [2]","                    - SCAN [orders] =\u003e [2:customer_id]","                            predicate: 2:customer_id IS NOT NULL AND 3:amount \u003e 5"]}
> raw
| span shop
| SELECT a, b FROM renamed
< raw
| {"columns":["a","b"],"tables":["shop.customers"],"views":["shop.renamed"],"plan":["- Output =\u003e [1:id, 2:name]","    - SCAN [customers] =\u003e [1:id, 2:name]"]}
> raw
| span shop
| SELECT id FROM all_ids
< raw
| {"columns":["id"],"tables":["shop.customers","shop.orders"],"views":["shop.all_ids"],"plan":["- Output =\u003e [10:id]","    - AGGREGATE(GLOBAL) [10:id]","        - EXCHANGE(SHUFFLE) [10]","            - AGGREGATE(LOCAL) [10:id]","                - UNION =\u003e [10: id]","                    - SCAN [orders] =\u003e [1:id]","                    - SCAN [customers] =\u003e [6:id]"]}
> raw
| span shop
| SELECT * FROM staff_pay
< raw
| {"columns":["id","salary","dept"],"tables":["hr.depts","hr.staff"],"views":["shop.staff_pay"],"plan":["- Output =\u003e [1:id, 3:salary, 6:name]","    - HASH/INNER JOIN [5:id = 4:dept_id] =\u003e [1:id, 3:salary, 6:name]","        - SCAN [depts] =\u003e [5:id, 6:name]","                predicate: 5:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [4]","            - SCAN [staff] =\u003e [1:id, 3:salary, 4:dept_id]","                    predicate: 4:dept_id IS NOT NULL"]}
> raw
| span shop
| SELECT name FROM hr.rich
< raw
| {"columns":["name"],"tables":["hr.staff"],"views":["hr.rich"],"plan":["- Output =\u003e [2:name]","    - SCAN [staff] =\u003e [2:name]","            predicate: 3:salary \u003e 120"]}
> raw
| span shop
| SELECT VipName FROM hr.`Rich_Names`
< raw
| {"columns":["VipName"],"tables":["hr.staff"],"views":["hr.Rich_Names","hr.rich"],"plan":["- Output =\u003e [2:name]","    - SCAN [staff] =\u003e [2:name]","            predicate: 3:salary \u003e 120"]}
> raw
| span shop
| SELECT b.id, t.total FROM big_orders b JOIN customer_totals t ON t.customer_id = b.customer_id
< raw
| {"columns":["id","total"],"tables":["shop.orders"],"views":["shop.big_orders","shop.customer_totals"],"plan":["- Output =\u003e [1:id, 12:sum]","    - HASH/INNER JOIN [8:customer_id = 2:customer_id] =\u003e [1:id, 12:sum]","        - EXCHANGE(SHUFFLE) [8]","            - AGGREGATE(GLOBAL) [8:customer_id]","                    12:sum := sum(9:amount)","                - SCAN [orders] =\u003e [8:customer_id, 9:amount]","                        predicate: 8:customer_id IS NOT NULL","        - EXCHANGE(SHUFFLE) [2]","            - SCAN [orders] =\u003e [1:id, 2:customer_id]","                    predicate: 2:customer_id IS NOT NULL AND 3:amount \u003e 5"]}
> raw
| span shop
| SELECT d.doubled FROM (SELECT doubled FROM big_orders) d
< raw
| {"columns":["doubled"],"tables":["shop.orders"],"views":["shop.big_orders"],"plan":["- Output =\u003e [6:expr]","    - SCAN [orders] =\u003e [6:expr]","            6:expr := cast(3:amount as decimal(10, 2)) * 2","            predicate: 3:amount \u003e 5"]}
> raw
| span shop
| WITH v AS (SELECT * FROM big_orders) SELECT doubled FROM v
< raw
| {"columns":["doubled"],"tables":["shop.orders"],"views":["shop.big_orders"],"plan":["- Output =\u003e [6:expr]","    - SCAN [orders] =\u003e [6:expr]","            6:expr := cast(3:amount as decimal(10, 2)) * 2","            predicate: 3:amount \u003e 5"]}
> raw
| span shop
| SELECT customer_id, total FROM order_sums
< raw
| {"columns":["customer_id","total"],"tables":["shop.order_sums","shop.orders"],"plan":["EXECUTE IN FE","- Output =\u003e [1:customer_id, 2:total]","    - EMPTY"]}
> raw
| span shop
| SELECT * FROM default_catalog.shop.order_sums
< raw
| {"columns":["customer_id","total"],"tables":["shop.order_sums","shop.orders"],"plan":["EXECUTE IN FE","- Output =\u003e [1:customer_id, 2:total]","    - EMPTY"]}
> raw
| span shop
| SELECT id, row_number() OVER (ORDER BY amount) AS rn FROM orders
< raw
| {"columns":["id","rn"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 6:row_number()]","    - ANALYTIC [ partition by), order by (3: amount ASC NULLS FIRST)]","            6:row_number() := row_number() ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW","        - SORT(FINAL)[3: amount ASC NULLS FIRST]","            - SORT(PARTIAL)[3: amount ASC NULLS FIRST]","                - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT id, sum(amount) OVER (PARTITION BY customer_id) AS running FROM orders
< raw
| {"columns":["id","running"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 6:sum(3: amount)]","    - ANALYTIC [ partition by (2:customer_id), order by)]","            6:sum(3: amount) := sum(3:amount) RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW","        - SORT(GLOBAL)[2: customer_id ASC NULLS FIRST]","            - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount]"]}
> raw
| span shop
| SELECT id, lag(amount) OVER (ORDER BY id) AS prev FROM orders
< raw
| {"columns":["id","prev"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 6:lag(3: amount, 1, null)]","    - ANALYTIC [ partition by), order by (1: id ASC NULLS FIRST)]","            6:lag(3: amount, 1, null) := lag(3:amount, 1, null) ROWS BETWEEN UNBOUNDED PRECEDING AND 1 PRECEDING","        - SORT(FINAL)[1: id ASC NULLS FIRST]","            - SORT(PARTIAL)[1: id ASC NULLS FIRST]","                - SCAN [orders] =\u003e [1:id, 3:amount]"]}
> raw
| span shop
| SELECT id, rank() OVER (PARTITION BY status ORDER BY amount) FROM orders
< raw
| {"columns":["id","rank() OVER (PARTITION BY status ORDER BY amount ASC )"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 6:rank()]","    - ANALYTIC [ partition by (4:status), order by (3: amount ASC NULLS FIRST)]","            6:rank() := rank() RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW","        - SORT(GLOBAL)[4: status ASC NULLS FIRST, 3: amount ASC NULLS FIRST]","            - SCAN [orders] =\u003e [1:id, 3:amount, 4:status]"]}
> raw
| span shop
| SELECT id, sum(amount) OVER (PARTITION BY customer_id, status ORDER BY id DESC, amount) FROM orders
< raw
| {"columns":["id","sum(amount) OVER (PARTITION BY customer_id, status ORDER BY id DESC, amount ASC )"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 6:sum(3: amount)]","    - ANALYTIC [ partition by (2:customer_id, 4:status), order by (1: id DESC NULLS LAST, 3: amount ASC NULLS FIRST)]","            6:sum(3: amount) := sum(3:amount) RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW","        - SORT(GLOBAL)[2: customer_id ASC NULLS FIRST, 4: status ASC NULLS FIRST, 1: id DESC NULLS LAST, 3: amount ASC NULLS FIRST]","            - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount, 4:status]"]}
> raw
| span shop
| SELECT d.`rank() OVER (PARTITION BY status ORDER BY amount ASC )` FROM (SELECT rank() OVER (PARTITION BY status ORDER BY amount) FROM orders) d
< raw
| {"columns":["rank() OVER (PARTITION BY status ORDER BY amount ASC )"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:rank()]","    - ANALYTIC [ partition by (4:status), order by (3: amount ASC NULLS FIRST)]","            6:rank() := rank() RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW","        - SORT(GLOBAL)[4: status ASC NULLS FIRST, 3: amount ASC NULLS FIRST]","            - SCAN [orders] =\u003e [3:amount, 4:status]"]}
> raw
| span shop
| SELECT id FROM orders QUALIFY row_number() OVER (PARTITION BY customer_id ORDER BY amount) = 1
< raw
| {"columns":["id"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id]","    - PREDICATE [6: row_number() = 1]","        - ANALYTIC [ partition by (2:customer_id), order by (3: amount ASC NULLS FIRST)]","                6:row_number() := row_number() ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW","            - SORT(GLOBAL)[2: customer_id ASC NULLS FIRST, 3: amount ASC NULLS FIRST]","                - SORT(PARTIAL)[2: customer_id ASC NULLS FIRST, 3: amount ASC NULLS FIRST]","                    - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount]"]}
> raw
| span shop
| SELECT customer_id, sum(amount), count(*) FROM orders GROUP BY customer_id
< raw
| {"columns":["customer_id","sum(amount)","count(*)"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 6:sum, 7:count]","    - AGGREGATE(GLOBAL) [2:customer_id]","            6:sum := sum(3:amount)","            7:count := count()","        - SCAN [orders] =\u003e [2:customer_id, 3:amount]"]}
> raw
| span shop
| SELECT status, avg(amount) AS a FROM orders GROUP BY status HAVING count(*) > 1
< raw
| {"columns":["status","a"],"tables":["shop.orders"],"plan":["- Output =\u003e [4:status, 6:avg]","    - AGGREGATE(GLOBAL) [4:status] =\u003e [4:status, 6:avg]","            6:avg := avg(3:amount)","            7:count := count()","            predicate: 7:count \u003e 1","        - SCAN [orders] =\u003e [3:amount, 4:status]"]}
> raw
| span shop
| SELECT region, group_concat(name ORDER BY name) AS names FROM customers GROUP BY region
< raw
| {"columns":["region","names"],"tables":["shop.customers"],"plan":["- Output =\u003e [3:region, 5:group_concat]","    - AGGREGATE(GLOBAL) [3:region]","            5:group_concat := group_concat(2:name, ',', 2:name)","        - SCAN [customers] =\u003e [2:name, 3:region]"]}
> raw
| span shop
| SELECT max(amount) FROM orders WHERE status = 'open'
< raw
| {"columns":["max(amount)"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:max]","    - AGGREGATE(GLOBAL) []","            6:max := max(3:amount)","        - SCAN [orders] =\u003e [3:amount]","                predicate: 4:status = 'open'"]}
> raw
| span shop
| SELECT count(DISTINCT customer_id) AS buyers FROM orders
< raw
| {"columns":["buyers"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:count]","    - AGGREGATE(GLOBAL) []","            6:count := count(2:customer_id)","        - SCAN [orders] =\u003e [2:customer_id]"]}
> raw
| span shop
| SELECT c.region, sum(o.amount) FROM customers c JOIN orders o ON o.customer_id = c.id GROUP BY c.region
< raw
| {"columns":["region","sum(o.amount)"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [3:region, 10:sum]","    - AGGREGATE(GLOBAL) [3:region]","            10:sum := sum(10:sum)","        - EXCHANGE(SHUFFLE) [3]","            - AGGREGATE(LOCAL) [3:region]","                    10:sum := sum(7:amount)","                - HASH/INNER JOIN [1:id = 6:customer_id] =\u003e [3:region, 7:amount]","                    - SCAN [customers] =\u003e [1:id, 3:region]","                            predicate: 1:id IS NOT NULL","                    - EXCHANGE(SHUFFLE) [6]","                        - SCAN [orders] =\u003e [6:customer_id, 7:amount]","                                predicate: 6:customer_id IS NOT NULL"]}
> raw
| span shop
| SELECT customer_id, status, sum(amount) FROM orders GROUP BY ROLLUP (customer_id, status)
< raw
| {"columns":["customer_id","status","sum(amount)"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 4:status, 6:sum]","    - AGGREGATE(GLOBAL) [2:customer_id, 4:status, 7:GROUPING_ID] =\u003e [2:customer_id, 4:status, 6:sum]","            6:sum := sum(6:sum)","        - EXCHANGE(SHUFFLE) [2, 4, 7]","            - AGGREGATE(LOCAL) [2:customer_id, 4:status, 7:GROUPING_ID]","                    6:sum := sum(3:amount)","                - REPEAT [[], [2:customer_id], [2:customer_id, 4:status]] =\u003e [2:customer_id, 4:status, 7:GROUPING_ID]","                    - SCAN [orders] =\u003e [2:customer_id, 3:amount, 4:status]"]}
> raw
| span shop
| SELECT count(*) FROM orders WHERE 1 = 0
< raw
| {"columns":["count(*)"],"tables":["shop.orders"],"plan":["- Output =\u003e [6:count]","    - AGGREGATE(GLOBAL) []","            6:count := count()","        - EMPTY"]}
> raw
| span shop
| SELECT * FROM orders PIVOT (sum(amount) FOR status IN ('open', 'paid'))
< raw
| {"columns":["id","customer_id","created_at","'open'","'paid'"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 2:customer_id, 5:created_at, 8:sum, 9:sum]","    - AGGREGATE(GLOBAL) [1:id, 2:customer_id, 5:created_at]","            8:sum := sum(6:case)","            9:sum := sum(7:case)","        - SCAN [orders] =\u003e [1:id, 2:customer_id, 5:created_at, 6:case, 7:case]","                6:case := if(4:status = 'open', 3:amount, null)","                7:case := if(4:status = 'paid', 3:amount, null)"]}
> raw
| span shop
| SELECT * FROM orders PIVOT (sum(amount) AS total, count(*) AS n FOR status IN ('open' AS o, 'paid' AS p))
< raw
| {"columns":["id","customer_id","created_at","o_total","o_n","p_total","p_n"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 2:customer_id, 5:created_at, 10:sum, 11:count, 12:sum, 13:count]","    - AGGREGATE(GLOBAL) [1:id, 2:customer_id, 5:created_at]","            10:sum := sum(6:case)","            11:count := count(7:case)","            12:sum := sum(8:case)","            13:count := count(9:case)","        - SCAN [orders] =\u003e [1:id, 2:customer_id, 5:created_at, 6:case, 7:case, 8:case, 9:case]","                6:case := if(14:expr, 3:amount, null)","                7:case := if(14:expr, 1, null)","                8:case := if(15:expr, 3:amount, null)","                9:case := if(15:expr, 1, null)","                14:expr := 4:status = 'open', ","                15:expr := 4:status = 'paid', "]}
> raw
| span shop
| SELECT * FROM orders PIVOT (SUM(amount), max(id) FOR status IN ('open', 'paid'))
< raw
| {"columns":["customer_id","created_at","'open'_sum(shop.orders.amount)","'open'_max(shop.orders.id)","'paid'_sum(shop.orders.amount)","'paid'_max(shop.orders.id)"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 5:created_at, 10:sum, 11:max, 12:sum, 13:max]","    - AGGREGATE(GLOBAL) [2:customer_id, 5:created_at]","            10:sum := sum(6:case)","            11:max := max(7:case)","            12:sum := sum(8:case)","            13:max := max(9:case)","        - SCAN [orders] =\u003e [2:customer_id, 5:created_at, 6:case, 7:case, 8:case, 9:case]","                6:case := if(14:expr, 3:amount, null)","                7:case := if(14:expr, 1:id, null)","                8:case := if(15:expr, 3:amount, null)","                9:case := if(15:expr, 1:id, null)","                14:expr := 4:status = 'open', ","                15:expr := 4:status = 'paid', "]}
> raw
| span shop
| SELECT * FROM orders PIVOT (sum(amount) FOR (status, customer_id) IN (('open', 1) AS o1, ('paid', 2)))
< raw
| {"columns":["id","created_at","o1","{'paid',2}"],"tables":["shop.orders"],"plan":["- Output =\u003e [1:id, 5:created_at, 8:sum, 9:sum]","    - AGGREGATE(GLOBAL) [1:id, 5:created_at]","            8:sum := sum(6:case)","            9:sum := sum(7:case)","        - SCAN [orders] =\u003e [1:id, 5:created_at, 6:case, 7:case]","                6:case := if(4:status = 'open' AND 2:customer_id = 1, 3:amount, null)","                7:case := if(4:status = 'paid' AND 2:customer_id = 2, 3:amount, null)"]}
> raw
| span shop
| SELECT * FROM (SELECT customer_id, status, amount FROM orders) t PIVOT (sum(amount) FOR customer_id IN (1, 2))
< raw
| {"columns":["status","1","2"],"tables":["shop.orders"],"plan":["- Output =\u003e [4:status, 8:sum, 9:sum]","    - AGGREGATE(GLOBAL) [4:status]","            8:sum := sum(6:case)","            9:sum := sum(7:case)","        - SCAN [orders] =\u003e [4:status, 6:case, 7:case]","                6:case := if(2:customer_id = 1, 3:amount, null)","                7:case := if(2:customer_id = 2, 3:amount, null)"]}
> raw
| span shop
| SELECT * FROM (SELECT customer_id, status, amount FROM orders) t PIVOT (sum(amount * 2) AS d FOR status IN ('open', 'paid'))
< raw
| {"columns":["customer_id","amount","'open'","'paid'"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 3:amount, 8:sum, 9:sum]","    - AGGREGATE(GLOBAL) [2:customer_id, 3:amount]","            8:sum := sum(6:case)","            9:sum := sum(7:case)","        - SCAN [orders] =\u003e [2:customer_id, 3:amount, 6:case, 7:case]","                6:case := if(4:status = 'open', 11:multiply, null)","                7:case := if(4:status = 'paid', 11:multiply, null)","                10:cast := cast(3:amount as decimal(10, 2)), ","                11:multiply := 10:cast * 2, "]}
> raw
| span shop
| SELECT * FROM (SELECT customer_id, status FROM orders) t PIVOT (count(*) AS n FOR status IN ('open', 'paid'))
< raw
| {"columns":["customer_id","'open'","'paid'"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 8:count, 9:count]","    - AGGREGATE(GLOBAL) [2:customer_id]","            8:count := count(6:case)","            9:count := count(7:case)","        - SCAN [orders] =\u003e [2:customer_id, 6:case, 7:case]","                6:case := if(4:status = 'open', 1, null)","                7:case := if(4:status = 'paid', 1, null)"]}
> raw
| span shop
| SELECT * FROM (SELECT id, status, created_at FROM orders) t PIVOT (count(id) FOR created_at IN (DATE '2026-01-01', '2026-01-02'))
< raw
| {"columns":["status","'2026-01-01'","'2026-01-02'"],"tables":["shop.orders"],"plan":["- Output =\u003e [4:status, 8:count, 9:count]","    - AGGREGATE(GLOBAL) [4:status]","            8:count := count(6:case)","            9:count := count(7:case)","        - SCAN [orders] =\u003e [4:status, 6:case, 7:case]","                6:case := if(5:created_at = 2026-01-01, 1:id, null)","                7:case := if(5:created_at = 2026-01-02, 1:id, null)"]}
> raw
| span shop
| SELECT * FROM (SELECT id, amount FROM orders) t PIVOT (count(id) FOR amount IN (10.00, 20, NULL))
< raw
| {"columns":["10.00","20","NULL"],"tables":["shop.orders"],"plan":["- Output =\u003e [9:count, 10:count, 11:count]","    - AGGREGATE(GLOBAL) []","            9:count := count(6:case)","            10:count := count(7:case)","            11:count := count(null)","        - SCAN [orders] =\u003e [6:case, 7:case]","                6:case := if(3:amount = 10.00, 1:id, null)","                7:case := if(3:amount = 20, 1:id, null)"]}
> raw
| span shop
| SELECT * FROM orders o JOIN customers c ON c.id = o.customer_id PIVOT (sum(o.amount) FOR region IN ('east', 'west'))
< raw
| {"columns":["id","customer_id","status","created_at","id","name","email","'east'","'west'"],"tables":["shop.customers","shop.orders"],"plan":["- Output =\u003e [1:id, 2:customer_id, 4:status, 5:created_at, 6:id, 7:name, 9:email, 12:sum, 13:sum]","    - AGGREGATE(GLOBAL) [1:id, 2:customer_id, 4:status, 5:created_at, 6:id, 7:name, 9:email]","            12:sum := sum(10:case)","            13:sum := sum(11:case)","        - HASH/INNER JOIN [6:id = 2:customer_id] =\u003e [1:id, 2:customer_id, 4:status, 5:created_at, 6:id, 7:name, 9:email, 10:case, 11:case]","                10:case := if(8:region = 'east', 3:amount, null)","                11:case := if(8:region = 'west', 3:amount, null)","            - SCAN [customers] =\u003e [6:id, 7:name, 8:region, 9:email]","                    predicate: 6:id IS NOT NULL","            - EXCHANGE(SHUFFLE) [2]","                - SCAN [orders] =\u003e [1:id, 2:customer_id, 3:amount, 4:status, 5:created_at]","                        predicate: 2:customer_id IS NOT NULL"]}
> raw
| span shop
| SELECT customer_id, o AS open_total FROM orders PIVOT (sum(amount) FOR status IN ('open' AS o, 'paid' AS p))
< raw
| {"columns":["customer_id","open_total"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 8:sum]","    - AGGREGATE(GLOBAL) [1:id, 2:customer_id, 5:created_at] =\u003e [2:customer_id, 8:sum]","            8:sum := sum(6:case)","        - SCAN [orders] =\u003e [1:id, 2:customer_id, 5:created_at, 6:case]","                6:case := if(4:status = 'open', 3:amount, null)"]}
> raw
| span shop
| SELECT t.customer_id, `'open'` FROM (SELECT customer_id, status, amount FROM orders) t PIVOT (sum(amount) FOR status IN ('open', 'paid')) WHERE customer_id = 1
< raw
| {"columns":["customer_id","'open'"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 8:sum]","    - AGGREGATE(GLOBAL) [2:customer_id]","            8:sum := sum(6:case)","        - SCAN [orders] =\u003e [2:customer_id, 6:case]","                6:case := if(4:status = 'open', 3:amount, null)","                predicate: 2:customer_id = 1"]}
> raw
| span shop
| SELECT customer_id, sum(o) FROM (SELECT customer_id, status, amount FROM orders) t PIVOT (sum(amount) FOR status IN ('open' AS o)) GROUP BY customer_id
< raw
| {"columns":["customer_id","sum(o)"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 8:sum]","    - AGGREGATE(GLOBAL) [2:customer_id] =\u003e [2:customer_id, 8:sum]","            7:sum := sum(6:case)","            8:sum := 7:sum","        - SCAN [orders] =\u003e [2:customer_id, 6:case]","                6:case := if(4:status = 'open', 3:amount, null)"]}
> raw
| span shop
| SELECT * FROM (SELECT customer_id, status, amount FROM orders) t PIVOT (sum(amount) FOR status IN ('open' AS 'Open Orders'))
< raw
| {"columns":["customer_id","Open Orders"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 7:sum]","    - AGGREGATE(GLOBAL) [2:customer_id]","            7:sum := sum(6:case)","        - SCAN [orders] =\u003e [2:customer_id, 6:case]","                6:case := if(4:status = 'open', 3:amount, null)"]}
> raw
| span shop
| SELECT s.c FROM (SELECT customer_id AS c FROM orders PIVOT (sum(amount) FOR status IN ('open'))) s
< raw
| {"columns":["c"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id]","    - AGGREGATE(GLOBAL) [1:id, 2:customer_id, 5:created_at] =\u003e [2:customer_id]","        - SCAN [orders] =\u003e [1:id, 2:customer_id, 5:created_at]"]}
> raw
| span shop
| SELECT * FROM customer_totals PIVOT (sum(total) FOR customer_id IN (1, 2))
< raw
| {"columns":["n","1","2"],"tables":["shop.orders"],"views":["shop.customer_totals"],"plan":["- Output =\u003e [7:count, 10:sum, 11:sum]","    - AGGREGATE(GLOBAL) [7:count]","            10:sum := sum(10:sum)","            11:sum := sum(11:sum)","        - EXCHANGE(SHUFFLE) [7]","            - AGGREGATE(LOCAL) [7:count]","                    10:sum := sum(8:case)","                    11:sum := sum(9:case)","                - AGGREGATE(GLOBAL) [2:customer_id] =\u003e [7:count, 8:case, 9:case]","                        6:sum := sum(3:amount)","                        7:count := count()","                        8:case := if(2:customer_id = 1, 6:sum, null)","                        9:case := if(2:customer_id = 2, 6:sum, null)","                    - SCAN [orders] =\u003e [2:customer_id, 3:amount]"]}
> raw
| span shop
| SELECT * FROM orders o PIVOT (sum(o.amount), count(id) FOR status IN ('open'))
< raw
| {"columns":["customer_id","created_at","'open'_sum(shop.o.amount)","'open'_count(shop.o.id)"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 5:created_at, 8:sum, 9:count]","    - AGGREGATE(GLOBAL) [2:customer_id, 5:created_at]","            8:sum := sum(6:case)","            9:count := count(7:case)","        - SCAN [orders] =\u003e [2:customer_id, 5:created_at, 6:case, 7:case]","                6:case := if(10:expr, 3:amount, null)","                7:case := if(10:expr, 1:id, null)","                10:expr := 4:status = 'open', "]}
> raw
| span shop
| SELECT * FROM (SELECT customer_id, status, amount FROM orders) t PIVOT (sum((amount)), max(customer_id + 1) FOR status IN ('open'))
< raw
| {"columns":["customer_id","'open'_sum(t.amount)","'open'_max(t.customer_id + 1)"],"tables":["shop.orders"],"plan":["- Output =\u003e [2:customer_id, 8:sum, 9:max]","    - AGGREGATE(GLOBAL) [2:customer_id]","            8:sum := sum(6:case)","            9:max := max(7:case)","        - SCAN [orders] =\u003e [2:customer_id, 6:case, 7:case]","                6:case := if(10:expr, 3:amount, null)","                7:case := if(10:expr, cast(2:customer_id as bigint(20)) + 1, null)","                10:expr := 4:status = 'open', "]}
> raw
| span shop
| SELECT * FROM hr.staff PIVOT (sum(salary), count(*) FOR dept_id IN (1, 2))
< raw
| {"columns":["id","name","1_sum(hr.staff.salary)","1_count(*)","2_sum(hr.staff.salary)","2_count(*)"],"tables":["hr.staff"],"plan":["- Output =\u003e [1:id, 2:name, 9:sum, 10:count, 11:sum, 12:count]","    - AGGREGATE(GLOBAL) [1:id, 2:name]","            9:sum := sum(5:case)","            10:count := count(6:case)","            11:sum := sum(7:case)","            12:count := count(8:case)","        - SCAN [staff] =\u003e [1:id, 2:name, 5:case, 6:case, 7:case, 8:case]","                5:case := if(13:expr, 3:salary, null)","                6:case := if(13:expr, 1, null)","                7:case := if(14:expr, 3:salary, null)","                8:case := if(14:expr, 1, null)","                13:expr := 4:dept_id = 1, ","                14:expr := 4:dept_id = 2, "]}
> raw
| span shop
| WITH c AS (SELECT status, amount FROM orders) SELECT * FROM c PIVOT (sum(amount), min(amount) FOR status IN ('open'))
< raw
| {"columns":["'open'_sum(c.amount)","'open'_min(c.amount)"],"tables":["shop.orders"],"plan":["- Output =\u003e [7:sum, 8:min]","    - AGGREGATE(GLOBAL) []","            7:sum := sum(6:case)","            8:min := min(6:case)","        - SCAN [orders] =\u003e [6:case]","                6:case := if(4:status = 'open', 3:amount, null)"]}
> raw
| span shop
| SELECT * FROM (SELECT id, created_at FROM orders) t PIVOT (count(id) FOR created_at IN (DATETIME '2026-01-01 00:00:00'))
< raw
| {"columns":["'2026-01-01 00:00:00'"],"tables":["shop.orders"],"plan":["- Output =\u003e [7:count]","    - AGGREGATE(GLOBAL) []","            7:count := count(if(5:created_at = 2026-01-01, 1:id, null))","        - SCAN [orders] =\u003e [1:id, 5:created_at]"]}
> raw
| span shop
| SELECT * FROM orders PIVOT (sum(amount) FOR status IN ('open')) WHERE status = 'open'
< raw
| {"error":"1064 Getting analyzing error. Detail message: Column 'status' cannot be resolved."}
> raw
| span shop
| SELECT amount FROM orders PIVOT (sum(amount) FOR status IN ('open'))
< raw
| {"error":"1064 Getting analyzing error. Detail message: Column 'amount' cannot be resolved."}
> raw
| span shop
| SELECT * FROM orders PIVOT (sum(nope) FOR status IN ('open'))
< raw
| {"error":"1064 Getting analyzing error. Detail message: Column 'nope' cannot be resolved."}
> raw
| span shop
| SELECT * FROM orders PIVOT (sum(amount) FOR nope IN ('open'))
< raw
| {"error":"1064 Getting analyzing error. Detail message: Column 'nope' cannot be resolved."}
> raw
| span shop
| SELECT generate_series FROM TABLE(generate_series(1, 3))
< raw
| {"columns":["generate_series"],"plan":["- Output =\u003e [2:generate_series]","    - TABLE FUNCTION [generate_series] =\u003e [{2}]","        - VALUES =\u003e [3:expr, 4:expr, 5:expr]","                3:expr := 1","                4:expr := 3","                5:expr := 1"]}
> raw
| span shop
| SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_SCHEMA = 'shop'
< raw
| {"columns":["TABLE_NAME"],"tables":["information_schema.tables"],"plan":["- Output =\u003e [3:TABLE_NAME]","    - SCHEMA-SCAN [tables] =\u003e [3:TABLE_NAME]","            predicate: 2:TABLE_SCHEMA = 'shop'"]}
> raw
| span shop
| SELECT nope FROM orders
< raw
| {"error":"1064 Getting analyzing error. Detail message: Column 'nope' cannot be resolved."}
> raw
| span shop
| SELECT id FROM no_such_table
< raw
| {"error":"5502 Getting analyzing error. Detail message: Unknown table 'shop.no_such_table'."}
> raw
| span shop
| SELECT id FROM orders JOIN customers ON customers.id = orders.customer_id
< raw
| {"error":"1064 Getting analyzing error. Detail message: Column 'id' is ambiguous."}
> raw
| span shop
| SELECT id FROM no_such_db.orders
< raw
| {"error":"5501 Getting analyzing error. Detail message: Unknown database 'no_such_db'."}
> raw
| span hr
| SELECT name, salary FROM staff
< raw
| {"columns":["name","salary"],"tables":["hr.staff"],"plan":["- Output =\u003e [2:name, 3:salary]","    - SCAN [staff] =\u003e [2:name, 3:salary]"]}
> raw
| span hr
| SELECT s.name, d.name AS dept FROM staff s JOIN depts d ON d.id = s.dept_id
< raw
| {"columns":["name","dept"],"tables":["hr.depts","hr.staff"],"plan":["- Output =\u003e [2:name, 6:name]","    - HASH/INNER JOIN [4:dept_id = 5:id] =\u003e [2:name, 6:name]","        - EXCHANGE(SHUFFLE) [4]","            - SCAN [staff] =\u003e [2:name, 4:dept_id]","                    predicate: 4:dept_id IS NOT NULL","        - EXCHANGE(SHUFFLE) [5]","            - SCAN [depts] =\u003e [5:id, 6:name]","                    predicate: 5:id IS NOT NULL"]}
> raw
| span hr
| SELECT o.amount FROM shop.orders o
< raw
| {"columns":["amount"],"tables":["shop.orders"],"plan":["- Output =\u003e [3:amount]","    - SCAN [orders] =\u003e [3:amount]"]}
> raw
| span hr
| SELECT * FROM rich
< raw
| {"columns":["name","salary"],"tables":["hr.staff"],"views":["hr.rich"],"plan":["- Output =\u003e [2:name, 3:salary]","    - SCAN [staff] =\u003e [2:name, 3:salary]","            predicate: 3:salary \u003e 120"]}
> raw
| span hr
| SELECT * FROM shop.staff_pay
< raw
| {"columns":["id","salary","dept"],"tables":["hr.depts","hr.staff"],"views":["shop.staff_pay"],"plan":["- Output =\u003e [1:id, 3:salary, 6:name]","    - HASH/INNER JOIN [5:id = 4:dept_id] =\u003e [1:id, 3:salary, 6:name]","        - SCAN [depts] =\u003e [5:id, 6:name]","                predicate: 5:id IS NOT NULL","        - EXCHANGE(SHUFFLE) [4]","            - SCAN [staff] =\u003e [1:id, 3:salary, 4:dept_id]","                    predicate: 4:dept_id IS NOT NULL"]}
> raw
| span hr
| SELECT name FROM staff UNION SELECT name FROM shop.customers
< raw
| {"columns":["name"],"tables":["hr.staff","shop.customers"],"plan":["- Output =\u003e [9:name]","    - AGGREGATE(GLOBAL) [9:name]","        - EXCHANGE(SHUFFLE) [9]","            - AGGREGATE(LOCAL) [9:name]","                - UNION =\u003e [9: name]","                    - SCAN [staff] =\u003e [2:name]","                    - SCAN [customers] =\u003e [6:name]"]}
> raw
| span hr
| SELECT id FROM default_catalog.shop.big_orders
< raw
| {"columns":["id"],"tables":["shop.orders"],"views":["shop.big_orders"],"plan":["- Output =\u003e [1:id]","    - SCAN [orders] =\u003e [1:id]","            predicate: 3:amount \u003e 5"]}
