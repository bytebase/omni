-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| {"databases":[{"name":"hr","tables":[{"name":"depts","columns":["id","name"]},{"name":"staff","columns":["id","name","salary","dept_id"]}],"views":[{"name":"Rich_Names","columns":["VipName"],"definition":"SELECT `name` AS `VipName` FROM `hr`.`rich`"},{"name":"rich","columns":["name","salary"],"definition":"SELECT `name` AS `name`,`salary` AS `salary` FROM `hr`.`staff` WHERE `salary`\u003e120"}]},{"name":"shop","tables":[{"name":"MixedCase","columns":["Id","CamelCol","lower_col"]},{"name":"customers","columns":["id","name","region","email"]},{"name":"items","columns":["order_id","sku","qty"]},{"name":"orders","columns":["id","customer_id","amount","status","created_at"]}],"views":[{"name":"all_ids","columns":["id"],"definition":"SELECT `id` AS `id` FROM `shop`.`orders` UNION SELECT `id` AS `id` FROM `shop`.`customers`"},{"name":"big_orders","columns":["id","doubled","status","customer_id"],"definition":"SELECT `id` AS `id`,`amount`*2 AS `doubled`,`status` AS `status`,`customer_id` AS `customer_id` FROM `shop`.`orders` WHERE `amount`\u003e5"},{"name":"customer_totals","columns":["customer_id","total","n"],"definition":"SELECT `customer_id` AS `customer_id`,SUM(`amount`) AS `total`,COUNT(1) AS `n` FROM `shop`.`orders` GROUP BY `customer_id`"},{"name":"order_names","columns":["id","doubled","name"],"definition":"SELECT `b`.`id` AS `id`,`b`.`doubled` AS `doubled`,`c`.`name` AS `name` FROM `shop`.`big_orders` AS `b` JOIN `shop`.`customers` AS `c` ON `c`.`id`=`b`.`customer_id`"},{"name":"renamed","columns":["a","b"],"definition":"SELECT `id` AS `id`,`name` AS `name` FROM `shop`.`customers`"},{"name":"staff_pay","columns":["id","salary","dept"],"definition":"SELECT `s`.`id` AS `id`,`s`.`salary` AS `salary`,`d`.`name` AS `dept` FROM `hr`.`staff` AS `s` JOIN `hr`.`depts` AS `d` ON `d`.`id`=`s`.`dept_id`"}]}]}
> raw
| span shop
| SELECT id, amount FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT o.id AS order_id, o.amount FROM orders o
< raw
| {"columns":[{"name":"order_id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT id + customer_id, upper(status) FROM orders
< raw
| {"columns":[{"name":"id + customer_id"},{"name":"upper(status)"}],"plan":["Projection\troot\t\tplus(shop.orders.id, shop.orders.customer_id)-\u003eColumn#6, upper(shop.orders.status)-\u003eColumn#7","└─TableReader\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT CAST(amount AS CHAR) AS a, concat(name, '@', region) AS b FROM customers
< raw
| {"error":"1054"}
> raw
| span shop
| SELECT 1, 'a' AS b
< raw
| {"columns":[{"name":"1"},{"name":"b"}],"plan":["Projection\troot\t\t1-\u003eColumn#1, a-\u003eColumn#2","└─TableDual\troot\t\trows:1"]}
> raw
| span shop
| SELECT count(*) FROM orders
< raw
| {"columns":[{"name":"count(*)"}],"plan":["HashAgg\troot\t\tfuncs:count(Column#7)-\u003eColumn#6","└─TableReader\troot\t\tdata:HashAgg","  └─HashAgg\tcop[tikv]\t\tfuncs:count(1)-\u003eColumn#7","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT DISTINCT region FROM customers
< raw
| {"columns":[{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"}],"plan":["HashAgg\troot\t\tgroup by:shop.customers.region, funcs:firstrow(shop.customers.region)-\u003eshop.customers.region","└─TableReader\troot\t\tdata:HashAgg","  └─HashAgg\tcop[tikv]\t\tgroup by:shop.customers.region, ","    └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| SELECT name FROM customers ORDER BY region LIMIT 2
< raw
| {"columns":[{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"}],"plan":["Projection\troot\t\tshop.customers.name","└─TopN\troot\t\tshop.customers.region, offset:0, count:2","  └─TableReader\troot\t\tdata:TopN","    └─TopN\tcop[tikv]\t\tshop.customers.region, offset:0, count:2","      └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| SELECT CASE WHEN amount > 5 THEN status ELSE 'small' END AS label FROM orders
< raw
| {"columns":[{"name":"label"}],"plan":["Projection\troot\t\tcase(gt(shop.orders.amount, 5), shop.orders.status, small)-\u003eColumn#6","└─TableReader\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT coalesce(email, name) AS contact FROM customers
< raw
| {"columns":[{"name":"contact"}],"plan":["Projection\troot\t\tcoalesce(shop.customers.email, shop.customers.name)-\u003eColumn#5","└─TableReader\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| SELECT amount AS id, id AS amount FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"plan":["Projection\troot\t\tshop.orders.amount, shop.orders.id","└─TableReader\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT id   +   1, -amount, id IN (1, 2), NOT id FROM orders
< raw
| {"columns":[{"name":"id   +   1"},{"name":"-amount"},{"name":"id IN (1, 2)"},{"name":"NOT id"}],"plan":["Projection\troot\t\tplus(shop.orders.id, 1)-\u003eColumn#6, unaryminus(shop.orders.amount)-\u003eColumn#7, in(shop.orders.id, 1, 2)-\u003eColumn#8, not(shop.orders.id)-\u003eColumn#9","└─TableReader\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT 'x', 2.50, NULL, TRUE, _utf8mb4'y'
< raw
| {"columns":[{"name":"x"},{"name":"2.50"},{"name":"NULL"},{"name":"TRUE"},{"name":"y"}],"plan":["Projection\troot\t\tx-\u003eColumn#1, 2.50-\u003eColumn#2, \u003cnil\u003e-\u003eColumn#3, 1-\u003eColumn#4, y-\u003eColumn#5","└─TableDual\troot\t\trows:1"]}
> raw
| span shop
| TABLE customers
< raw
| {"columns":[{"name":"id","schema":"shop","table":"customers","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"customers","orgTable":"customers","orgName":"email"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| VALUES ROW(1, 2 + 3), ROW(4, 5)
< raw
| {"error":"1051"}
> raw
| span shop
| SELECT * FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT o.* FROM orders o
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"o","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"o","orgTable":"orders","orgName":"created_at"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT *, amount FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"plan":["Projection\troot\t\tshop.orders.id, shop.orders.customer_id, shop.orders.amount, shop.orders.status, shop.orders.created_at, shop.orders.amount","└─TableReader\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT c.*, o.amount FROM customers c JOIN orders o ON o.customer_id = c.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"plan":["Projection\troot\t\tshop.customers.id, shop.customers.name, shop.customers.region, shop.customers.email, shop.orders.amount","└─HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","  ├─TableReader(Build)\troot\t\tdata:Selection","  │ └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","  │   └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo","  └─TableReader(Probe)\troot\t\tdata:TableFullScan","    └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.customers,shop.orders"}
> raw
| span shop
| SELECT * FROM (SELECT * FROM customers) c
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| SELECT id FROM shop.orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT shop.orders.amount, orders.status FROM shop.orders
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT s.name, s.salary FROM hr.staff s
< raw
| {"columns":[{"name":"name","schema":"hr","table":"s","orgTable":"staff","orgName":"name"},{"name":"salary","schema":"hr","table":"s","orgTable":"staff","orgName":"salary"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:s\tkeep order:false, stats:pseudo"],"tables":"hr.staff"}
> raw
| span shop
| SELECT o.id, s.name FROM orders o JOIN hr.staff s ON s.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"hr","table":"s","orgTable":"staff","orgName":"name"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, hr.staff.id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","│   └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:s\tkeep order:false, stats:pseudo"],"tables":"shop.orders,hr.staff"}
> raw
| span shop
| SELECT ID, Amount FROM orders
< raw
| {"columns":[{"name":"ID","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"Amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT `id`, `AMOUNT` FROM `orders`
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"AMOUNT","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT O.ID FROM orders O
< raw
| {"columns":[{"name":"ID","schema":"shop","table":"O","orgTable":"orders","orgName":"id"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:O\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT shop.o.id, shop.orders.status FROM orders o JOIN shop.orders ON orders.id = o.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"}],"plan":["MergeJoin\troot\t\tinner join, left key:shop.orders.id, right key:shop.orders.id","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:true, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:true, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT `Id`, `CamelCol`, lower_col FROM `MixedCase`
< raw
| {"columns":[{"name":"Id","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"Id"},{"name":"CamelCol","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"CamelCol"},{"name":"lower_col","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"lower_col"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:MixedCase\tkeep order:false, stats:pseudo"],"tables":"shop.mixedcase"}
> raw
| span shop
| SELECT * FROM `MixedCase`
< raw
| {"columns":[{"name":"Id","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"Id"},{"name":"CamelCol","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"CamelCol"},{"name":"lower_col","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"lower_col"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:MixedCase\tkeep order:false, stats:pseudo"],"tables":"shop.mixedcase"}
> raw
| span shop
| SELECT * FROM mixedcase
< raw
| {"columns":[{"name":"Id","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"Id"},{"name":"CamelCol","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"CamelCol"},{"name":"lower_col","schema":"shop","table":"MixedCase","orgTable":"MixedCase","orgName":"lower_col"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:MixedCase\tkeep order:false, stats:pseudo"],"tables":"shop.mixedcase"}
> raw
| span shop
| SELECT m.camelcol, m.LOWER_COL FROM `MixedCase` m
< raw
| {"columns":[{"name":"camelcol","schema":"shop","table":"m","orgTable":"MixedCase","orgName":"CamelCol"},{"name":"LOWER_COL","schema":"shop","table":"m","orgTable":"MixedCase","orgName":"lower_col"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:m\tkeep order:false, stats:pseudo"],"tables":"shop.mixedcase"}
> raw
| span shop
| SELECT orders.name FROM customers AS orders
< raw
| {"columns":[{"name":"name","schema":"shop","table":"orders","orgTable":"customers","orgName":"name"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| SELECT customers.amount FROM orders customers
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"customers","orgTable":"orders","orgName":"amount"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT o.id FROM customers o JOIN orders c ON c.customer_id = o.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"customers","orgName":"id"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","│   └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo"],"tables":"shop.customers,shop.orders"}
> raw
| span shop
| SELECT * FROM (SELECT amount AS id FROM orders) orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","│   └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"plan":["HashJoin\troot\t\tleft outer join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT o.id, c.name FROM customers c RIGHT JOIN orders o ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"plan":["Projection\troot\t\tshop.orders.id, shop.customers.name","└─HashJoin\troot\t\tright outer join, equal:[eq(shop.customers.id, shop.orders.customer_id)]","  ├─TableReader(Build)\troot\t\tdata:TableFullScan","  │ └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo","  └─TableReader(Probe)\troot\t\tdata:TableFullScan","    └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.customers,shop.orders"}
> raw
| span shop
| SELECT o.id, c.name FROM orders o CROSS JOIN customers c
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"plan":["HashJoin\troot\t\tCARTESIAN inner join","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","│   └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT o.id, c.name FROM orders o STRAIGHT_JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","│   └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT id FROM orders JOIN customers USING (id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"plan":["MergeJoin\troot\t\tinner join, left key:shop.orders.id, right key:shop.customers.id","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:true, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:true, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT * FROM orders JOIN customers USING (id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"},{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"customers","orgTable":"customers","orgName":"email"}],"plan":["MergeJoin\troot\t\tinner join, left key:shop.orders.id, right key:shop.customers.id","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:true, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:true, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT * FROM orders o LEFT JOIN customers c USING (id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"o","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"o","orgTable":"orders","orgName":"created_at"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"c","orgTable":"customers","orgName":"email"}],"plan":["MergeJoin\troot\t\tleft outer join, left key:shop.orders.id, right key:shop.customers.id","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:true, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:true, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT * FROM orders NATURAL JOIN customers
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"amount","schema":"shop","table":"orders","orgTable":"orders","orgName":"amount"},{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"created_at","schema":"shop","table":"orders","orgTable":"orders","orgName":"created_at"},{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"},{"name":"email","schema":"shop","table":"customers","orgTable":"customers","orgName":"email"}],"plan":["MergeJoin\troot\t\tinner join, left key:shop.orders.id, right key:shop.customers.id","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:true, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:true, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT id, name FROM orders NATURAL LEFT JOIN customers
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"}],"plan":["MergeJoin\troot\t\tleft outer join, left key:shop.orders.id, right key:shop.customers.id","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:true, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:true, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT i.sku, o.status, c.region FROM items i JOIN orders o ON o.id = i.order_id JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"sku","schema":"shop","table":"i","orgTable":"items","orgName":"sku"},{"name":"status","schema":"shop","table":"o","orgTable":"orders","orgName":"status"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo","└─HashJoin(Probe)\troot\t\tinner join, equal:[eq(shop.items.order_id, shop.orders.id)]","  ├─TableReader(Build)\troot\t\tdata:Selection","  │ └─Selection\tcop[tikv]\t\tnot(isnull(shop.items.order_id))","  │   └─TableFullScan\tcop[tikv]\ttable:i\tkeep order:false, stats:pseudo","  └─TableReader(Probe)\troot\t\tdata:Selection","    └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","      └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo"],"tables":"shop.items,shop.orders,shop.customers"}
> raw
| span shop
| SELECT a.id, b.id FROM orders a JOIN orders b ON b.customer_id = a.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"a","orgTable":"orders","orgName":"id"},{"name":"id","schema":"shop","table":"b","orgTable":"orders","orgName":"id"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.orders.customer_id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","│   └─TableFullScan\tcop[tikv]\ttable:b\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","    └─TableFullScan\tcop[tikv]\ttable:a\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT o.id FROM orders o JOIN (customers c JOIN items i ON i.order_id = c.id) ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.customers.id, shop.items.order_id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(shop.items.order_id))","│   └─TableFullScan\tcop[tikv]\ttable:i\tkeep order:false, stats:pseudo","└─HashJoin(Probe)\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","  ├─TableReader(Build)\troot\t\tdata:TableFullScan","  │ └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo","  └─TableReader(Probe)\troot\t\tdata:Selection","    └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","      └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers,shop.items"}
> raw
| span shop
| SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"},{"name":"who"}],"plan":["HashJoin\troot\t\tleft outer join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT id, (SELECT max(amount) FROM orders) AS top FROM customers
< raw
| {"columns":[{"name":"id","schema":"shop","table":"customers","orgTable":"customers","orgName":"id"},{"name":"top"}],"plan":["Projection\troot\t\tshop.customers.id, 20.00-\u003eColumn#16","└─TableReader\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers,shop.orders"}
> raw
| span shop
| SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'east')
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.customers.id, shop.orders.customer_id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\teq(shop.customers.region, \"east\")","│   └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT id FROM orders o WHERE EXISTS (SELECT 1 FROM items i WHERE i.order_id = o.id)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"o","orgTable":"orders","orgName":"id"}],"plan":["HashJoin\troot\t\tsemi join, equal:[eq(shop.orders.id, shop.items.order_id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(shop.items.order_id))","│   └─TableFullScan\tcop[tikv]\ttable:i\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.items"}
> raw
| span shop
| SELECT d.k, d.v FROM (SELECT id AS k, amount AS v FROM orders) d
< raw
| {"columns":[{"name":"k","schema":"shop","table":"d","orgTable":"orders","orgName":"id"},{"name":"v","schema":"shop","table":"d","orgTable":"orders","orgName":"amount"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT d.k FROM (SELECT id AS k FROM orders LIMIT 2) d
< raw
| {"columns":[{"name":"k","schema":"shop","table":"d","orgTable":"orders","orgName":"id"}],"plan":["Limit\troot\t\toffset:0, count:2","└─TableReader\troot\t\tdata:Limit","  └─Limit\tcop[tikv]\t\toffset:0, count:2","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT x.total FROM (SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id) x
< raw
| {"columns":[{"name":"total","schema":"shop","table":"x"}],"plan":["HashAgg\troot\t\tgroup by:shop.orders.customer_id, funcs:sum(Column#7)-\u003eColumn#6","└─TableReader\troot\t\tdata:HashAgg","  └─HashAgg\tcop[tikv]\t\tgroup by:shop.orders.customer_id, funcs:sum(shop.orders.amount)-\u003eColumn#7","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT d.a, d.b FROM (SELECT id, name FROM customers) AS d (a, b)
< raw
| {"error":"1064"}
> raw
| span shop
| SELECT o.id, l.n FROM orders o, LATERAL (SELECT count(*) AS n FROM items i WHERE i.order_id = o.id) l
< raw
| {"error":"1064"}
> raw
| span shop
| SELECT o.id, l.sku FROM orders o JOIN LATERAL (SELECT sku FROM items i WHERE i.order_id = o.id LIMIT 1) l ON true
< raw
| {"error":"1064"}
> raw
| span shop
| SELECT y.k FROM (SELECT x.k FROM (SELECT id AS k FROM customers) x) y
< raw
| {"columns":[{"name":"k","schema":"shop","table":"y","orgTable":"customers","orgName":"id"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers) SELECT id, name FROM c
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"c","orgTable":"customers","orgName":"name"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| WITH c (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM c
< raw
| {"columns":[{"name":"k","schema":"shop","table":"c","orgTable":"orders","orgName":"id"},{"name":"v","schema":"shop","table":"c","orgTable":"orders","orgName":"amount"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT * FROM b
< raw
| {"columns":[{"name":"id","schema":"shop","table":"b","orgTable":"orders","orgName":"id"},{"name":"name","schema":"shop","table":"b","orgTable":"customers","orgName":"name"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","│   └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| WITH c AS (SELECT id FROM customers) SELECT x.id, y.id FROM c x JOIN c y ON y.id = x.id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"x","orgTable":"customers","orgName":"id"},{"name":"id","schema":"shop","table":"y","orgTable":"customers","orgName":"id"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.customers.id, shop.customers.id)]","├─Selection(Build)\troot\t\tnot(isnull(shop.customers.id))","│ └─CTEFullScan\troot\tCTE:c AS y\tdata:CTE_0","└─Selection(Probe)\troot\t\tnot(isnull(shop.customers.id))","  └─CTEFullScan\troot\tCTE:c AS x\tdata:CTE_0","CTE_0\troot\t\tNon-Recursive CTE","└─TableReader(Seed Part)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| WITH RECURSIVE r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r
< raw
| {"columns":[{"name":"n","schema":"shop","table":"r"}],"plan":["CTEFullScan\troot\tCTE:r\tdata:CTE_0","CTE_0\troot\t\tRecursive CTE","├─Projection(Seed Part)\troot\t\t1-\u003eColumn#2","│ └─TableDual\troot\t\trows:1","└─Projection(Recursive Part)\troot\t\tcast(plus(Column#3, 1), bigint BINARY)-\u003eColumn#5","  └─Selection\troot\t\tlt(Column#3, 3)","    └─CTETable\troot\t\tScan on CTE_0"]}
> raw
| span shop
| WITH RECURSIVE r (id, name) AS (SELECT id, name FROM customers WHERE id = 1 UNION ALL SELECT c.id, c.name FROM customers c JOIN r ON c.id = r.id + 1) SELECT id, name FROM r
< raw
| {"columns":[{"name":"id","schema":"shop","table":"r","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"r","orgTable":"customers","orgName":"name"}],"plan":["CTEFullScan\troot\tCTE:r\tdata:CTE_0","CTE_0\troot\t\tRecursive CTE","├─Point_Get(Seed Part)\troot\ttable:customers\thandle:1","└─IndexHashJoin(Recursive Part)\troot\t\tinner join, inner:TableReader, outer key:Column#23, inner key:shop.customers.id, equal cond:eq(Column#23, shop.customers.id)","  ├─Projection(Build)\troot\t\tplus(shop.customers.id, 1)-\u003eColumn#23","  │ └─CTETable\troot\t\tScan on CTE_0","  └─TableReader(Probe)\troot\t\tdata:TableRangeScan","    └─TableRangeScan\tcop[tikv]\ttable:c\trange: decided by [Column#23], keep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| WITH orders AS (SELECT id, name FROM customers) SELECT id, name FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"customers","orgName":"id"},{"name":"name","schema":"shop","table":"orders","orgTable":"customers","orgName":"name"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| WITH c AS (SELECT id, amount FROM orders) SELECT id FROM c WHERE amount > (SELECT avg(amount) FROM c)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"orders","orgName":"id"}],"plan":["HashJoin\troot\t\tCARTESIAN inner join, other cond:gt(shop.orders.amount, Column#10)","├─Selection(Build)\troot\t\tnot(isnull(Column#10))","│ └─MaxOneRow\troot\t\t","│   └─HashAgg\troot\t\tfuncs:avg(shop.orders.amount)-\u003eColumn#10","│     └─CTEFullScan\troot\tCTE:c\tdata:CTE_0","└─Selection(Probe)\troot\t\tnot(isnull(shop.orders.amount))","  └─CTEFullScan\troot\tCTE:c\tdata:CTE_0","CTE_0\troot\t\tNon-Recursive CTE","└─TableReader(Seed Part)\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tor(not(isnull(shop.orders.amount)), 1)","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| WITH RECURSIVE c AS (SELECT id FROM orders), orders AS (SELECT 1 AS id) SELECT id FROM c
< raw
| {"columns":[{"name":"id","schema":"shop","table":"c","orgTable":"orders","orgName":"id"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT id FROM orders UNION SELECT id FROM customers
< raw
| {"columns":[{"name":"id"}],"plan":["HashAgg\troot\t\tgroup by:Column#10, funcs:firstrow(Column#10)-\u003eColumn#10","└─Union\troot\t\t","  ├─TableReader\troot\t\tdata:TableFullScan","  │ └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","  └─TableReader\troot\t\tdata:TableFullScan","    └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT id, status FROM orders UNION ALL SELECT id, name FROM customers
< raw
| {"columns":[{"name":"id"},{"name":"status"}],"plan":["Union\troot\t\t","├─Projection\troot\t\tshop.orders.id-\u003eColumn#10, cast(shop.orders.status, varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)-\u003eColumn#11","│ └─TableReader\troot\t\tdata:TableFullScan","│   └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","└─TableReader\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT name FROM customers INTERSECT SELECT status FROM orders
< raw
| {"columns":[{"name":"name","schema":"shop","table":"customers","orgTable":"customers","orgName":"name"}],"plan":["HashJoin\troot\t\tsemi join, equal:[nulleq(shop.customers.name, shop.orders.status)]","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","└─HashAgg(Probe)\troot\t\tgroup by:shop.customers.name, funcs:firstrow(shop.customers.name)-\u003eshop.customers.name","  └─TableReader\troot\t\tdata:HashAgg","    └─HashAgg\tcop[tikv]\t\tgroup by:shop.customers.name, ","      └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers,shop.orders"}
> raw
| span shop
| SELECT id FROM orders EXCEPT SELECT customer_id FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"}],"plan":["HashJoin\troot\t\tanti semi join, equal:[nulleq(shop.orders.id, shop.orders.customer_id)]","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT id AS k FROM orders UNION SELECT 1 ORDER BY k
< raw
| {"columns":[{"name":"k"}],"plan":["Sort\troot\t\tColumn#7","└─HashAgg\troot\t\tgroup by:Column#7, funcs:firstrow(Column#7)-\u003eColumn#7","  └─Union\troot\t\t","    ├─Projection\troot\t\tcast(shop.orders.id, bigint BINARY)-\u003eColumn#7","    │ └─TableReader\troot\t\tdata:TableFullScan","    │   └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","    └─Projection\troot\t\t1-\u003eColumn#7","      └─TableDual\troot\t\trows:1"],"tables":"shop.orders"}
> raw
| span shop
| (SELECT id FROM orders LIMIT 1) UNION (SELECT id FROM customers LIMIT 1)
< raw
| {"columns":[{"name":"id"}],"plan":["HashAgg\troot\t\tgroup by:Column#10, funcs:firstrow(Column#10)-\u003eColumn#10","└─Union\troot\t\t","  ├─Limit\troot\t\toffset:0, count:1","  │ └─TableReader\troot\t\tdata:Limit","  │   └─Limit\tcop[tikv]\t\toffset:0, count:1","  │     └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","  └─Limit\troot\t\toffset:0, count:1","    └─TableReader\troot\t\tdata:Limit","      └─Limit\tcop[tikv]\t\toffset:0, count:1","        └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT u.id FROM (SELECT id FROM orders UNION SELECT id FROM customers) u
< raw
| {"columns":[{"name":"id","schema":"shop","table":"u"}],"plan":["HashAgg\troot\t\tgroup by:Column#10, funcs:firstrow(Column#10)-\u003eColumn#10","└─Union\troot\t\t","  ├─TableReader\troot\t\tdata:TableFullScan","  │ └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","  └─TableReader\troot\t\tdata:TableFullScan","    └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.orders,shop.customers"}
> raw
| span shop
| SELECT id, doubled FROM big_orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"big_orders","orgTable":"orders","orgName":"id"},{"name":"doubled","schema":"shop","table":"big_orders","orgName":"doubled"}],"plan":["Projection\troot\t\tshop.orders.id-\u003eColumn#1, mul(shop.orders.amount, 2)-\u003eColumn#6","└─TableReader\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tgt(shop.orders.amount, 5)","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.big_orders"}
> raw
| span shop
| SELECT b.status, b.doubled FROM big_orders b
< raw
| {"columns":[{"name":"status","schema":"shop","table":"b","orgTable":"orders","orgName":"status"},{"name":"doubled","schema":"shop","table":"b","orgName":"doubled"}],"plan":["Projection\troot\t\tshop.orders.status-\u003eColumn#4, mul(shop.orders.amount, 2)-\u003eColumn#6","└─TableReader\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tgt(shop.orders.amount, 5)","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.big_orders"}
> raw
| span shop
| SELECT * FROM big_orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"big_orders","orgTable":"orders","orgName":"id"},{"name":"doubled","schema":"shop","table":"big_orders","orgName":"doubled"},{"name":"status","schema":"shop","table":"big_orders","orgTable":"orders","orgName":"status"},{"name":"customer_id","schema":"shop","table":"big_orders","orgTable":"orders","orgName":"customer_id"}],"plan":["Projection\troot\t\tshop.orders.id-\u003eColumn#1, mul(shop.orders.amount, 2)-\u003eColumn#6, shop.orders.status-\u003eColumn#4, shop.orders.customer_id-\u003eColumn#2","└─TableReader\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tgt(shop.orders.amount, 5)","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.big_orders"}
> raw
| span shop
| SELECT customer_id, total, n FROM customer_totals
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"customer_totals","orgTable":"orders","orgName":"customer_id"},{"name":"total","schema":"shop","table":"customer_totals","orgName":"total"},{"name":"n","schema":"shop","table":"customer_totals","orgName":"n"}],"plan":["Projection\troot\t\tshop.orders.customer_id, Column#6, Column#7","└─HashAgg\troot\t\tgroup by:shop.orders.customer_id, funcs:sum(Column#8)-\u003eColumn#6, funcs:count(Column#9)-\u003eColumn#7, funcs:firstrow(shop.orders.customer_id)-\u003eshop.orders.customer_id","  └─TableReader\troot\t\tdata:HashAgg","    └─HashAgg\tcop[tikv]\t\tgroup by:shop.orders.customer_id, funcs:sum(shop.orders.amount)-\u003eColumn#8, funcs:count(1)-\u003eColumn#9","      └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.customer_totals"}
> raw
| span shop
| SELECT * FROM order_names
< raw
| {"columns":[{"name":"id","schema":"shop","table":"order_names","orgTable":"orders","orgName":"id"},{"name":"doubled","schema":"shop","table":"order_names","orgName":"doubled"},{"name":"name","schema":"shop","table":"order_names","orgTable":"customers","orgName":"name"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(Column#2, shop.customers.id)]","├─Projection(Build)\troot\t\tshop.orders.id-\u003eColumn#1, mul(shop.orders.amount, 2)-\u003eColumn#6, shop.orders.customer_id-\u003eColumn#2","│ └─TableReader\troot\t\tdata:Selection","│   └─Selection\tcop[tikv]\t\tgt(shop.orders.amount, 5), not(isnull(shop.orders.customer_id))","│     └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.order_names"}
> raw
| span shop
| SELECT o.name, c.region FROM order_names o JOIN customers c ON c.name = o.name
< raw
| {"columns":[{"name":"name","schema":"shop","table":"o","orgTable":"customers","orgName":"name"},{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.customers.name, shop.customers.name)]","├─HashJoin(Build)\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","│ ├─TableReader(Build)\troot\t\tdata:Selection","│ │ └─Selection\tcop[tikv]\t\tgt(shop.orders.amount, 5), not(isnull(shop.orders.customer_id))","│ │   └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","│ └─TableReader(Probe)\troot\t\tdata:Selection","│   └─Selection\tcop[tikv]\t\tnot(isnull(shop.customers.name))","│     └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tnot(isnull(shop.customers.name))","    └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo"],"tables":"shop.order_names,shop.customers"}
> raw
| span shop
| SELECT a, b FROM renamed
< raw
| {"columns":[{"name":"a","schema":"shop","table":"renamed","orgTable":"customers","orgName":"id"},{"name":"b","schema":"shop","table":"renamed","orgTable":"customers","orgName":"name"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.renamed"}
> raw
| span shop
| SELECT id FROM all_ids
< raw
| {"columns":[{"name":"id","schema":"shop","table":"all_ids","orgName":"id"}],"plan":["HashAgg\troot\t\tgroup by:Column#10, funcs:firstrow(Column#10)-\u003eColumn#10","└─Union\troot\t\t","  ├─TableReader\troot\t\tdata:TableFullScan","  │ └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","  └─TableReader\troot\t\tdata:TableFullScan","    └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.all_ids"}
> raw
| span shop
| SELECT * FROM staff_pay
< raw
| {"columns":[{"name":"id","schema":"shop","table":"staff_pay","orgTable":"staff","orgName":"id"},{"name":"salary","schema":"shop","table":"staff_pay","orgTable":"staff","orgName":"salary"},{"name":"dept","schema":"shop","table":"staff_pay","orgTable":"depts","orgName":"dept"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(hr.staff.dept_id, hr.depts.id)]","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:d\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tnot(isnull(hr.staff.dept_id))","    └─TableFullScan\tcop[tikv]\ttable:s\tkeep order:false, stats:pseudo"],"tables":"shop.staff_pay"}
> raw
| span shop
| SELECT name FROM hr.rich
< raw
| {"columns":[{"name":"name","schema":"hr","table":"rich","orgTable":"staff","orgName":"name"}],"plan":["TableReader\troot\t\tdata:Projection","└─Projection\tcop[tikv]\t\thr.staff.name","  └─Selection\tcop[tikv]\t\tgt(hr.staff.salary, 120)","    └─TableFullScan\tcop[tikv]\ttable:staff\tkeep order:false, stats:pseudo"],"tables":"hr.rich"}
> raw
| span shop
| SELECT VipName FROM hr.`Rich_Names`
< raw
| {"columns":[{"name":"VipName","schema":"hr","table":"Rich_Names","orgTable":"staff","orgName":"VipName"}],"plan":["TableReader\troot\t\tdata:Projection","└─Projection\tcop[tikv]\t\thr.staff.name","  └─Selection\tcop[tikv]\t\tgt(hr.staff.salary, 120)","    └─TableFullScan\tcop[tikv]\ttable:staff\tkeep order:false, stats:pseudo"],"tables":"hr.rich_names"}
> raw
| span shop
| SELECT b.id, t.total FROM big_orders b JOIN customer_totals t ON t.customer_id = b.customer_id
< raw
| {"columns":[{"name":"id","schema":"shop","table":"b","orgTable":"orders","orgName":"id"},{"name":"total","schema":"shop","table":"t","orgName":"total"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.orders.customer_id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tgt(shop.orders.amount, 5), not(isnull(shop.orders.customer_id))","│   └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo","└─HashAgg(Probe)\troot\t\tgroup by:shop.orders.customer_id, funcs:sum(Column#14)-\u003eColumn#12, funcs:firstrow(shop.orders.customer_id)-\u003eshop.orders.customer_id","  └─TableReader\troot\t\tdata:HashAgg","    └─HashAgg\tcop[tikv]\t\tgroup by:shop.orders.customer_id, funcs:sum(shop.orders.amount)-\u003eColumn#14","      └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","        └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.big_orders,shop.customer_totals"}
> raw
| span shop
| SELECT d.doubled FROM (SELECT doubled FROM big_orders) d
< raw
| {"columns":[{"name":"doubled","schema":"shop","table":"d","orgName":"doubled"}],"plan":["Projection\troot\t\tmul(shop.orders.amount, 2)-\u003eColumn#6","└─TableReader\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tgt(shop.orders.amount, 5)","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.big_orders"}
> raw
| span shop
| WITH v AS (SELECT * FROM big_orders) SELECT doubled FROM v
< raw
| {"columns":[{"name":"doubled","schema":"shop","table":"v","orgName":"doubled"}],"plan":["Projection\troot\t\tmul(shop.orders.amount, 2)-\u003eColumn#16","└─TableReader\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tgt(shop.orders.amount, 5)","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.big_orders"}
> raw
| span shop
| SELECT id, row_number() OVER (ORDER BY amount) AS rn FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"rn"}],"plan":["Projection\troot\t\tshop.orders.id, Column#7-\u003eColumn#9","└─Window\troot\t\trow_number()-\u003eColumn#7 over(order by shop.orders.amount rows between current row and current row)","  └─Sort\troot\t\tshop.orders.amount","    └─TableReader\troot\t\tdata:TableFullScan","      └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT id, sum(amount) OVER (PARTITION BY customer_id) AS running FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"running"}],"plan":["Projection\troot\t\tshop.orders.id, Column#7-\u003eColumn#9","└─Shuffle\troot\t\texecution info: concurrency:5, data sources:[TableReader]","  └─Window\troot\t\tsum(shop.orders.amount)-\u003eColumn#7 over(partition by shop.orders.customer_id)","    └─Sort\troot\t\tshop.orders.customer_id","      └─ShuffleReceiver\troot\t\t","        └─TableReader\troot\t\tdata:TableFullScan","          └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT id, rank() OVER w AS r FROM orders WINDOW w AS (PARTITION BY status ORDER BY amount)
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"r"}],"plan":["Projection\troot\t\tshop.orders.id, Column#7-\u003eColumn#9","└─Shuffle\troot\t\texecution info: concurrency:5, data sources:[TableReader]","  └─Window\troot\t\trank()-\u003eColumn#7 over(partition by shop.orders.status order by shop.orders.amount)","    └─Sort\troot\t\tshop.orders.status, shop.orders.amount","      └─ShuffleReceiver\troot\t\t","        └─TableReader\troot\t\tdata:TableFullScan","          └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT id, lag(amount) OVER (ORDER BY id) AS prev FROM orders
< raw
| {"columns":[{"name":"id","schema":"shop","table":"orders","orgTable":"orders","orgName":"id"},{"name":"prev"}],"plan":["Projection\troot\t\tshop.orders.id, Column#7-\u003eColumn#9","└─Window\troot\t\tlag(shop.orders.amount)-\u003eColumn#7 over(order by shop.orders.id)","  └─TableReader\troot\t\tdata:TableFullScan","    └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:true, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT customer_id, sum(amount), count(*) FROM orders GROUP BY customer_id
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"orders","orgTable":"orders","orgName":"customer_id"},{"name":"sum(amount)"},{"name":"count(*)"}],"plan":["Projection\troot\t\tshop.orders.customer_id, Column#6, Column#7","└─HashAgg\troot\t\tgroup by:shop.orders.customer_id, funcs:sum(Column#8)-\u003eColumn#6, funcs:count(Column#9)-\u003eColumn#7, funcs:firstrow(shop.orders.customer_id)-\u003eshop.orders.customer_id","  └─TableReader\troot\t\tdata:HashAgg","    └─HashAgg\tcop[tikv]\t\tgroup by:shop.orders.customer_id, funcs:sum(shop.orders.amount)-\u003eColumn#8, funcs:count(1)-\u003eColumn#9","      └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT status, avg(amount) AS a FROM orders GROUP BY status HAVING count(*) > 1
< raw
| {"columns":[{"name":"status","schema":"shop","table":"orders","orgTable":"orders","orgName":"status"},{"name":"a"}],"plan":["Projection\troot\t\tshop.orders.status, Column#6-\u003eColumn#9","└─Selection\troot\t\tgt(Column#7, 1)","  └─HashAgg\troot\t\tgroup by:shop.orders.status, funcs:avg(Column#10, Column#11)-\u003eColumn#6, funcs:count(Column#12)-\u003eColumn#7, funcs:firstrow(shop.orders.status)-\u003eshop.orders.status","    └─TableReader\troot\t\tdata:HashAgg","      └─HashAgg\tcop[tikv]\t\tgroup by:shop.orders.status, funcs:count(shop.orders.amount)-\u003eColumn#10, funcs:sum(shop.orders.amount)-\u003eColumn#11, funcs:count(1)-\u003eColumn#12","        └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT region, group_concat(name ORDER BY name) AS names FROM customers GROUP BY region
< raw
| {"columns":[{"name":"region","schema":"shop","table":"customers","orgTable":"customers","orgName":"region"},{"name":"names"}],"plan":["Projection\troot\t\tshop.customers.region, Column#5","└─HashAgg\troot\t\tgroup by:shop.customers.region, funcs:group_concat(shop.customers.name order by shop.customers.name separator \",\")-\u003eColumn#5, funcs:firstrow(shop.customers.region)-\u003eshop.customers.region","  └─TableReader\troot\t\tdata:TableFullScan","    └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"shop.customers"}
> raw
| span shop
| SELECT max(amount) FROM orders WHERE status = 'open'
< raw
| {"columns":[{"name":"max(amount)"}],"plan":["StreamAgg\troot\t\tfuncs:max(shop.orders.amount)-\u003eColumn#6","└─TopN\troot\t\tshop.orders.amount:desc, offset:0, count:1","  └─TableReader\troot\t\tdata:TopN","    └─TopN\tcop[tikv]\t\tshop.orders.amount:desc, offset:0, count:1","      └─Selection\tcop[tikv]\t\teq(shop.orders.status, \"open\"), not(isnull(shop.orders.amount))","        └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT count(DISTINCT customer_id) AS buyers FROM orders
< raw
| {"columns":[{"name":"buyers"}],"plan":["HashAgg\troot\t\tfuncs:count(distinct shop.orders.customer_id)-\u003eColumn#6","└─TableReader\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span shop
| SELECT c.region, sum(o.amount) FROM customers c JOIN orders o ON o.customer_id = c.id GROUP BY c.region
< raw
| {"columns":[{"name":"region","schema":"shop","table":"c","orgTable":"customers","orgName":"region"},{"name":"sum(o.amount)"}],"plan":["Projection\troot\t\tshop.customers.region, Column#10","└─HashAgg\troot\t\tgroup by:shop.customers.region, funcs:sum(shop.orders.amount)-\u003eColumn#10, funcs:firstrow(shop.customers.region)-\u003eshop.customers.region","  └─Projection\troot\t\tshop.customers.region, shop.orders.amount","    └─HashJoin\troot\t\tinner join, equal:[eq(shop.orders.customer_id, shop.customers.id)]","      ├─TableReader(Build)\troot\t\tdata:TableFullScan","      │ └─TableFullScan\tcop[tikv]\ttable:c\tkeep order:false, stats:pseudo","      └─TableReader(Probe)\troot\t\tdata:Selection","        └─Selection\tcop[tikv]\t\tnot(isnull(shop.orders.customer_id))","          └─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo"],"tables":"shop.customers,shop.orders"}
> raw
| span shop
| SELECT customer_id, sum(amount) FROM orders GROUP BY customer_id WITH ROLLUP
< raw
| {"columns":[{"name":"customer_id","schema":"shop","table":"ex_orders","orgTable":"orders","orgName":"customer_id"},{"name":"sum(amount)"}],"plan":["Projection\troot\t\tColumn#6-\u003eColumn#9, Column#8","└─HashAgg\troot\t\tgroup by:Column#6, gid, funcs:sum(shop.orders.amount)-\u003eColumn#8, funcs:firstrow(Column#6)-\u003eColumn#6","  └─Expand\troot\t\tlevel-projection:[shop.orders.amount, \u003cnil\u003e-\u003eColumn#6, 0-\u003egid],[shop.orders.amount, Column#6, 1-\u003egid]; schema: [shop.orders.amount,Column#6,gid]","    └─Projection\troot\t\tshop.orders.amount, shop.orders.customer_id-\u003eColumn#6","      └─TableReader\troot\t\tdata:TableFullScan","        └─TableFullScan\tcop[tikv]\ttable:orders\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
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
| {"columns":[{"name":"name","schema":"hr","table":"staff","orgTable":"staff","orgName":"name"},{"name":"salary","schema":"hr","table":"staff","orgTable":"staff","orgName":"salary"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:staff\tkeep order:false, stats:pseudo"],"tables":"hr.staff"}
> raw
| span hr
| SELECT s.name, d.name AS dept FROM staff s JOIN depts d ON d.id = s.dept_id
< raw
| {"columns":[{"name":"name","schema":"hr","table":"s","orgTable":"staff","orgName":"name"},{"name":"dept","schema":"hr","table":"d","orgTable":"depts","orgName":"name"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(hr.staff.dept_id, hr.depts.id)]","├─TableReader(Build)\troot\t\tdata:Selection","│ └─Selection\tcop[tikv]\t\tnot(isnull(hr.staff.dept_id))","│   └─TableFullScan\tcop[tikv]\ttable:s\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:TableFullScan","  └─TableFullScan\tcop[tikv]\ttable:d\tkeep order:false, stats:pseudo"],"tables":"hr.staff,hr.depts"}
> raw
| span hr
| SELECT o.amount FROM shop.orders o
< raw
| {"columns":[{"name":"amount","schema":"shop","table":"o","orgTable":"orders","orgName":"amount"}],"plan":["TableReader\troot\t\tdata:TableFullScan","└─TableFullScan\tcop[tikv]\ttable:o\tkeep order:false, stats:pseudo"],"tables":"shop.orders"}
> raw
| span hr
| SELECT * FROM rich
< raw
| {"columns":[{"name":"name","schema":"hr","table":"rich","orgTable":"staff","orgName":"name"},{"name":"salary","schema":"hr","table":"rich","orgTable":"staff","orgName":"salary"}],"plan":["TableReader\troot\t\tdata:Selection","└─Selection\tcop[tikv]\t\tgt(hr.staff.salary, 120)","  └─TableFullScan\tcop[tikv]\ttable:staff\tkeep order:false, stats:pseudo"],"tables":"hr.rich"}
> raw
| span hr
| SELECT * FROM shop.staff_pay
< raw
| {"columns":[{"name":"id","schema":"shop","table":"staff_pay","orgTable":"staff","orgName":"id"},{"name":"salary","schema":"shop","table":"staff_pay","orgTable":"staff","orgName":"salary"},{"name":"dept","schema":"shop","table":"staff_pay","orgTable":"depts","orgName":"dept"}],"plan":["HashJoin\troot\t\tinner join, equal:[eq(hr.staff.dept_id, hr.depts.id)]","├─TableReader(Build)\troot\t\tdata:TableFullScan","│ └─TableFullScan\tcop[tikv]\ttable:d\tkeep order:false, stats:pseudo","└─TableReader(Probe)\troot\t\tdata:Selection","  └─Selection\tcop[tikv]\t\tnot(isnull(hr.staff.dept_id))","    └─TableFullScan\tcop[tikv]\ttable:s\tkeep order:false, stats:pseudo"],"tables":"shop.staff_pay"}
> raw
| span hr
| SELECT name FROM staff UNION SELECT name FROM shop.customers
< raw
| {"columns":[{"name":"name"}],"plan":["HashAgg\troot\t\tgroup by:Column#9, funcs:firstrow(Column#9)-\u003eColumn#9","└─Union\troot\t\t","  ├─TableReader\troot\t\tdata:TableFullScan","  │ └─TableFullScan\tcop[tikv]\ttable:staff\tkeep order:false, stats:pseudo","  └─TableReader\troot\t\tdata:TableFullScan","    └─TableFullScan\tcop[tikv]\ttable:customers\tkeep order:false, stats:pseudo"],"tables":"hr.staff,shop.customers"}
