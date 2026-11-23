-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestQuerySpanFidelity default#1
> raw
| fixture
| CREATE DATABASE shop;
| CREATE DATABASE hr;
| CREATE TABLE shop.customers (id INT PRIMARY KEY, name VARCHAR(20), region VARCHAR(10), email VARCHAR(40));
| CREATE TABLE shop.orders (id INT PRIMARY KEY, customer_id INT, amount DECIMAL(10,2), status VARCHAR(10), created_at DATE);
| CREATE TABLE shop.items (order_id INT, sku VARCHAR(10), qty INT);
| CREATE TABLE shop.`MixedCase` (`Id` INT, `CamelCol` VARCHAR(10), lower_col VARCHAR(10));
| CREATE TABLE shop.mixedcase (id INT, note VARCHAR(10));
| CREATE TABLE hr.staff (id INT PRIMARY KEY, name VARCHAR(20), salary INT, dept_id INT);
| CREATE TABLE hr.depts (id INT PRIMARY KEY, name VARCHAR(20));
| INSERT INTO shop.customers VALUES (1, 'ann', 'east', 'a@x'), (2, 'bob', 'west', 'b@x'), (3, 'cy', 'east', 'c@x');
| INSERT INTO shop.orders VALUES (1, 1, 10.00, 'open', '2026-01-01'), (2, 2, 20.00, 'paid', '2026-01-02'), (3, 1, 5.00, 'open', '2026-01-03');
| INSERT INTO shop.items VALUES (1, 'a', 1), (1, 'b', 2), (2, 'a', 3);
| INSERT INTO shop.`MixedCase` VALUES (1, 'x', 'y'), (2, 'z', 'w');
| INSERT INTO shop.mixedcase VALUES (1, 'n'), (2, 'm');
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
| {"databases":[{"name":"hr","tables":[{"name":"depts","columns":["id","name"]},{"name":"staff","columns":["id","name","salary","dept_id"]}],"views":[{"name":"Rich_Names","columns":["VipName"],"definition":"select `rich`.`name` AS `VipName` from `hr`.`rich`"},{"name":"rich","columns":["name","salary"],"definition":"select `hr`.`staff`.`name` AS `name`,`hr`.`staff`.`salary` AS `salary` from `hr`.`staff` where `hr`.`staff`.`salary` \u003e 120"}]},{"name":"shop","tables":[{"name":"MixedCase","columns":["Id","CamelCol","lower_col"]},{"name":"customers","columns":["id","name","region","email"]},{"name":"items","columns":["order_id","sku","qty"]},{"name":"mixedcase","columns":["id","note"]},{"name":"orders","columns":["id","customer_id","amount","status","created_at"]}],"views":[{"name":"all_ids","columns":["id"],"definition":"select `shop`.`orders`.`id` AS `id` from `shop`.`orders` union select `shop`.`customers`.`id` AS `id` from `shop`.`customers`"},{"name":"big_orders","columns":["id","doubled","status","customer_id"],"definition":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`amount` * 2 AS `doubled`,`shop`.`orders`.`status` AS `status`,`shop`.`orders`.`customer_id` AS `customer_id` from `shop`.`orders` where `shop`.`orders`.`amount` \u003e 5"},{"name":"customer_totals","columns":["customer_id","total","n"],"definition":"select `shop`.`orders`.`customer_id` AS `customer_id`,sum(`shop`.`orders`.`amount`) AS `total`,count(0) AS `n` from `shop`.`orders` group by `shop`.`orders`.`customer_id`"},{"name":"order_names","columns":["id","doubled","name"],"definition":"select `b`.`id` AS `id`,`b`.`doubled` AS `doubled`,`c`.`name` AS `name` from (`shop`.`big_orders` `b` join `shop`.`customers` `c` on(`c`.`id` = `b`.`customer_id`))"},{"name":"renamed","columns":["a","b"],"definition":"select `shop`.`customers`.`id` AS `a`,`shop`.`customers`.`name` AS `b` from `shop`.`customers`"},{"name":"staff_pay","columns":["id","salary","dept"],"definition":"select `s`.`id` AS `id`,`s`.`salary` AS `salary`,`d`.`name` AS `dept` from (`hr`.`staff` `s` join `hr`.`depts` `d` on(`d`.`id` = `s`.`dept_id`))"}]}]}
> raw
| span shop
| SELECT id, amount FROM orders
< raw
| {"columns":["id","amount"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`amount` AS `amount` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT o.id AS order_id, o.amount FROM orders o
< raw
| {"columns":["order_id","amount"],"resolved":"select `shop`.`o`.`id` AS `order_id`,`shop`.`o`.`amount` AS `amount` from `shop`.`orders` `o`","plan":["o\tALL\t"]}
> raw
| span shop
| SELECT id + customer_id, upper(status) FROM orders
< raw
| {"columns":["id + customer_id","upper(status)"],"resolved":"select `shop`.`orders`.`id` + `shop`.`orders`.`customer_id` AS `id + customer_id`,ucase(`shop`.`orders`.`status`) AS `upper(status)` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT CAST(amount AS CHAR) AS a, concat(name, '@', region) AS b FROM customers
< raw
| {"error":"1054"}
> raw
| span shop
| SELECT 1, 'a' AS b
< raw
| {"columns":["1","b"],"resolved":"select 1 AS `1`,'a' AS `b`","plan":["NULL\tNULL\tNo tables used"]}
> raw
| span shop
| SELECT count(*) FROM orders
< raw
| {"columns":["count(*)"],"resolved":"select count(0) AS `count(*)` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT DISTINCT region FROM customers
< raw
| {"columns":["region"],"resolved":"select distinct `shop`.`customers`.`region` AS `region` from `shop`.`customers`","plan":["customers\tALL\tUsing temporary"]}
> raw
| span shop
| SELECT name FROM customers ORDER BY region LIMIT 2
< raw
| {"columns":["name"],"resolved":"select `shop`.`customers`.`name` AS `name` from `shop`.`customers` order by `shop`.`customers`.`region` limit 2","plan":["customers\tALL\tUsing filesort"]}
> raw
| span shop
| SELECT CASE WHEN amount > 5 THEN status ELSE 'small' END AS label FROM orders
< raw
| {"columns":["label"],"resolved":"select case when `shop`.`orders`.`amount` \u003e 5 then `shop`.`orders`.`status` else 'small' end AS `label` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT coalesce(email, name) AS contact FROM customers
< raw
| {"columns":["contact"],"resolved":"select coalesce(`shop`.`customers`.`email`,`shop`.`customers`.`name`) AS `contact` from `shop`.`customers`","plan":["customers\tALL\t"]}
> raw
| span shop
| SELECT amount AS id, id AS amount FROM orders
< raw
| {"columns":["id","amount"],"resolved":"select `shop`.`orders`.`amount` AS `id`,`shop`.`orders`.`id` AS `amount` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT id   +   1, -amount, id IN (1, 2), NOT id FROM orders
< raw
| {"columns":["id   +   1","-amount","id IN (1, 2)","NOT id"],"resolved":"select `shop`.`orders`.`id` + 1 AS `id   +   1`,-`shop`.`orders`.`amount` AS `-amount`,`shop`.`orders`.`id` in (1,2) AS `id IN (1, 2)`,`shop`.`orders`.`id` = 0 AS `NOT id` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT 'x', 2.50, NULL, TRUE, _utf8mb4'y'
< raw
| {"columns":["x","2.50","NULL","TRUE","y"],"resolved":"select 'x' AS `x`,2.50 AS `2.50`,NULL AS `NULL`,1 AS `TRUE`,_utf8mb4'y' AS `y`","plan":["NULL\tNULL\tNo tables used"]}
> raw
| span shop
| TABLE customers
< raw
| {"error":"1064"}
> raw
| span shop
| VALUES (1, 2 + 3), (4, 5)
< raw
| {"columns":["1","2 + 3"],"resolved":"values (1,2 + 3),(4,5)","plan":["NULL\tNULL\tNo tables used"]}
> raw
| span shop
| SELECT * FROM orders
< raw
| {"columns":["id","customer_id","amount","status","created_at"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`customer_id` AS `customer_id`,`shop`.`orders`.`amount` AS `amount`,`shop`.`orders`.`status` AS `status`,`shop`.`orders`.`created_at` AS `created_at` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT o.* FROM orders o
< raw
| {"columns":["id","customer_id","amount","status","created_at"],"resolved":"select `shop`.`o`.`id` AS `id`,`shop`.`o`.`customer_id` AS `customer_id`,`shop`.`o`.`amount` AS `amount`,`shop`.`o`.`status` AS `status`,`shop`.`o`.`created_at` AS `created_at` from `shop`.`orders` `o`","plan":["o\tALL\t"]}
> raw
| span shop
| SELECT *, amount FROM orders
< raw
| {"columns":["id","customer_id","amount","status","created_at","amount"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`customer_id` AS `customer_id`,`shop`.`orders`.`amount` AS `amount`,`shop`.`orders`.`status` AS `status`,`shop`.`orders`.`created_at` AS `created_at`,`shop`.`orders`.`amount` AS `amount` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT c.*, o.amount FROM customers c JOIN orders o ON o.customer_id = c.id
< raw
| {"columns":["id","name","region","email","amount"],"resolved":"select `shop`.`c`.`id` AS `id`,`shop`.`c`.`name` AS `name`,`shop`.`c`.`region` AS `region`,`shop`.`c`.`email` AS `email`,`shop`.`o`.`amount` AS `amount` from `shop`.`customers` `c` join `shop`.`orders` `o` where `shop`.`c`.`id` = `shop`.`o`.`customer_id`","plan":["o\tALL\tUsing where","c\teq_ref\t"]}
> raw
| span shop
| SELECT * FROM (SELECT * FROM customers) c
< raw
| {"columns":["id","name","region","email"],"resolved":"select `shop`.`customers`.`id` AS `id`,`shop`.`customers`.`name` AS `name`,`shop`.`customers`.`region` AS `region`,`shop`.`customers`.`email` AS `email` from `shop`.`customers`","plan":["customers\tALL\t"]}
> raw
| span shop
| SELECT id FROM shop.orders
< raw
| {"columns":["id"],"resolved":"select `shop`.`orders`.`id` AS `id` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT shop.orders.amount, orders.status FROM shop.orders
< raw
| {"columns":["amount","status"],"resolved":"select `shop`.`orders`.`amount` AS `amount`,`shop`.`orders`.`status` AS `status` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT s.name, s.salary FROM hr.staff s
< raw
| {"columns":["name","salary"],"resolved":"select `hr`.`s`.`name` AS `name`,`hr`.`s`.`salary` AS `salary` from `hr`.`staff` `s`","plan":["s\tALL\t"]}
> raw
| span shop
| SELECT o.id, s.name FROM orders o JOIN hr.staff s ON s.id = o.customer_id
< raw
| {"columns":["id","name"],"resolved":"select `shop`.`o`.`id` AS `id`,`hr`.`s`.`name` AS `name` from `shop`.`orders` `o` join `hr`.`staff` `s` where `hr`.`s`.`id` = `shop`.`o`.`customer_id`","plan":["o\tALL\tUsing where","s\teq_ref\t"]}
> raw
| span shop
| SELECT ID, Amount FROM orders
< raw
| {"columns":["ID","Amount"],"resolved":"select `shop`.`orders`.`id` AS `ID`,`shop`.`orders`.`amount` AS `Amount` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT `id`, `AMOUNT` FROM `orders`
< raw
| {"columns":["id","AMOUNT"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`amount` AS `AMOUNT` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT O.ID FROM orders O
< raw
| {"columns":["ID"],"resolved":"select `shop`.`O`.`id` AS `ID` from `shop`.`orders` `O`","plan":["O\tALL\t"]}
> raw
| span shop
| SELECT shop.o.id, shop.orders.status FROM orders o JOIN shop.orders ON orders.id = o.id
< raw
| {"columns":["id","status"],"resolved":"select `shop`.`o`.`id` AS `id`,`shop`.`orders`.`status` AS `status` from `shop`.`orders` `o` join `shop`.`orders` where `shop`.`orders`.`id` = `shop`.`o`.`id`","plan":["o\tALL\t","orders\teq_ref\t"]}
> raw
| span shop
| SELECT `Id`, `CamelCol`, lower_col FROM `MixedCase`
< raw
| {"columns":["Id","CamelCol","lower_col"],"resolved":"select `shop`.`MixedCase`.`Id` AS `Id`,`shop`.`MixedCase`.`CamelCol` AS `CamelCol`,`shop`.`MixedCase`.`lower_col` AS `lower_col` from `shop`.`MixedCase`","plan":["MixedCase\tALL\t"]}
> raw
| span shop
| SELECT * FROM `MixedCase`
< raw
| {"columns":["Id","CamelCol","lower_col"],"resolved":"select `shop`.`MixedCase`.`Id` AS `Id`,`shop`.`MixedCase`.`CamelCol` AS `CamelCol`,`shop`.`MixedCase`.`lower_col` AS `lower_col` from `shop`.`MixedCase`","plan":["MixedCase\tALL\t"]}
> raw
| span shop
| SELECT * FROM mixedcase
< raw
| {"columns":["id","note"],"resolved":"select `shop`.`mixedcase`.`id` AS `id`,`shop`.`mixedcase`.`note` AS `note` from `shop`.`mixedcase`","plan":["mixedcase\tALL\t"]}
> raw
| span shop
| SELECT m.camelcol, m.LOWER_COL FROM `MixedCase` m
< raw
| {"columns":["camelcol","LOWER_COL"],"resolved":"select `shop`.`m`.`CamelCol` AS `camelcol`,`shop`.`m`.`lower_col` AS `LOWER_COL` from `shop`.`MixedCase` `m`","plan":["m\tALL\t"]}
> raw
| span shop
| SELECT orders.name FROM customers AS orders
< raw
| {"columns":["name"],"resolved":"select `shop`.`orders`.`name` AS `name` from `shop`.`customers` `orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT customers.amount FROM orders customers
< raw
| {"columns":["amount"],"resolved":"select `shop`.`customers`.`amount` AS `amount` from `shop`.`orders` `customers`","plan":["customers\tALL\t"]}
> raw
| span shop
| SELECT o.id FROM customers o JOIN orders c ON c.customer_id = o.id
< raw
| {"columns":["id"],"resolved":"select `shop`.`o`.`id` AS `id` from `shop`.`customers` `o` join `shop`.`orders` `c` where `shop`.`o`.`id` = `shop`.`c`.`customer_id`","plan":["c\tALL\tUsing where","o\teq_ref\t"]}
> raw
| span shop
| SELECT * FROM (SELECT amount AS id FROM orders) orders
< raw
| {"columns":["id"],"resolved":"select `shop`.`orders`.`amount` AS `id` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["id","name"],"resolved":"select `shop`.`o`.`id` AS `id`,`shop`.`c`.`name` AS `name` from `shop`.`orders` `o` join `shop`.`customers` `c` where `shop`.`c`.`id` = `shop`.`o`.`customer_id`","plan":["o\tALL\tUsing where","c\teq_ref\t"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["id","name"],"resolved":"select `shop`.`o`.`id` AS `id`,`shop`.`c`.`name` AS `name` from `shop`.`orders` `o` left join `shop`.`customers` `c` on(`shop`.`c`.`id` = `shop`.`o`.`customer_id` and `shop`.`o`.`customer_id` is not null) where 1","plan":["o\tALL\t","c\teq_ref\tUsing where"]}
> raw
| span shop
| SELECT o.id, c.name FROM customers c RIGHT JOIN orders o ON c.id = o.customer_id
< raw
| {"columns":["id","name"],"resolved":"select `shop`.`o`.`id` AS `id`,`shop`.`c`.`name` AS `name` from `shop`.`orders` `o` left join `shop`.`customers` `c` on(`shop`.`c`.`id` = `shop`.`o`.`customer_id` and `shop`.`o`.`customer_id` is not null) where 1","plan":["o\tALL\t","c\teq_ref\tUsing where"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o CROSS JOIN customers c
< raw
| {"columns":["id","name"],"resolved":"select `shop`.`o`.`id` AS `id`,`shop`.`c`.`name` AS `name` from `shop`.`orders` `o` join `shop`.`customers` `c`","plan":["o\tALL\t","c\tALL\tUsing join buffer (flat, BNL join)"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o, customers c WHERE c.id = o.customer_id
< raw
| {"columns":["id","name"],"resolved":"select `shop`.`o`.`id` AS `id`,`shop`.`c`.`name` AS `name` from `shop`.`orders` `o` join `shop`.`customers` `c` where `shop`.`c`.`id` = `shop`.`o`.`customer_id`","plan":["o\tALL\tUsing where","c\teq_ref\t"]}
> raw
| span shop
| SELECT o.id, c.name FROM orders o STRAIGHT_JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["id","name"],"resolved":"select `shop`.`o`.`id` AS `id`,`shop`.`c`.`name` AS `name` from `shop`.`orders` `o` straight_join `shop`.`customers` `c` where `shop`.`c`.`id` = `shop`.`o`.`customer_id`","plan":["o\tALL\tUsing where","c\teq_ref\t"]}
> raw
| span shop
| SELECT id FROM orders JOIN customers USING (id)
< raw
| {"columns":["id"],"resolved":"select `shop`.`orders`.`id` AS `id` from `shop`.`orders` join `shop`.`customers` where `shop`.`customers`.`id` = `shop`.`orders`.`id`","plan":["orders\tALL\t","customers\teq_ref\t"]}
> raw
| span shop
| SELECT * FROM orders JOIN customers USING (id)
< raw
| {"columns":["id","customer_id","amount","status","created_at","name","region","email"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`customer_id` AS `customer_id`,`shop`.`orders`.`amount` AS `amount`,`shop`.`orders`.`status` AS `status`,`shop`.`orders`.`created_at` AS `created_at`,`shop`.`customers`.`name` AS `name`,`shop`.`customers`.`region` AS `region`,`shop`.`customers`.`email` AS `email` from `shop`.`orders` join `shop`.`customers` where `shop`.`customers`.`id` = `shop`.`orders`.`id`","plan":["orders\tALL\t","customers\teq_ref\t"]}
> raw
| span shop
| SELECT * FROM orders o LEFT JOIN customers c USING (id)
< raw
| {"columns":["id","customer_id","amount","status","created_at","name","region","email"],"resolved":"select `shop`.`o`.`id` AS `id`,`shop`.`o`.`customer_id` AS `customer_id`,`shop`.`o`.`amount` AS `amount`,`shop`.`o`.`status` AS `status`,`shop`.`o`.`created_at` AS `created_at`,`shop`.`c`.`name` AS `name`,`shop`.`c`.`region` AS `region`,`shop`.`c`.`email` AS `email` from `shop`.`orders` `o` left join `shop`.`customers` `c` on(`shop`.`c`.`id` = `shop`.`o`.`id`) where 1","plan":["o\tALL\t","c\teq_ref\t"]}
> raw
| span shop
| SELECT * FROM orders NATURAL JOIN customers
< raw
| {"columns":["id","customer_id","amount","status","created_at","name","region","email"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`customer_id` AS `customer_id`,`shop`.`orders`.`amount` AS `amount`,`shop`.`orders`.`status` AS `status`,`shop`.`orders`.`created_at` AS `created_at`,`shop`.`customers`.`name` AS `name`,`shop`.`customers`.`region` AS `region`,`shop`.`customers`.`email` AS `email` from `shop`.`orders` join `shop`.`customers` where `shop`.`customers`.`id` = `shop`.`orders`.`id`","plan":["orders\tALL\t","customers\teq_ref\t"]}
> raw
| span shop
| SELECT id, name FROM orders NATURAL LEFT JOIN customers
< raw
| {"columns":["id","name"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`customers`.`name` AS `name` from `shop`.`orders` left join `shop`.`customers` on(`shop`.`customers`.`id` = `shop`.`orders`.`id`) where 1","plan":["orders\tALL\t","customers\teq_ref\t"]}
> raw
| span shop
| SELECT i.sku, o.status, c.region FROM items i JOIN orders o ON o.id = i.order_id JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":["sku","status","region"],"resolved":"select `shop`.`i`.`sku` AS `sku`,`shop`.`o`.`status` AS `status`,`shop`.`c`.`region` AS `region` from `shop`.`items` `i` join `shop`.`orders` `o` join `shop`.`customers` `c` where `shop`.`c`.`id` = `shop`.`o`.`customer_id` and `shop`.`o`.`id` = `shop`.`i`.`order_id`","plan":["i\tALL\tUsing where","o\teq_ref\tUsing where","c\teq_ref\t"]}
> raw
| span shop
| SELECT a.id, b.id FROM orders a JOIN orders b ON b.customer_id = a.customer_id
< raw
| {"columns":["id","id"],"resolved":"select `shop`.`a`.`id` AS `id`,`shop`.`b`.`id` AS `id` from `shop`.`orders` `a` join `shop`.`orders` `b` where `shop`.`b`.`customer_id` = `shop`.`a`.`customer_id`","plan":["a\tALL\t","b\tALL\tUsing where; Using join buffer (flat, BNL join)"]}
> raw
| span shop
| SELECT o.id FROM orders o JOIN (customers c JOIN items i ON i.order_id = c.id) ON c.id = o.customer_id
< raw
| {"columns":["id"],"resolved":"select `shop`.`o`.`id` AS `id` from `shop`.`orders` `o` join `shop`.`customers` `c` join `shop`.`items` `i` where `shop`.`c`.`id` = `shop`.`o`.`customer_id` and `shop`.`i`.`order_id` = `shop`.`o`.`customer_id`","plan":["o\tALL\tUsing where","c\teq_ref\t","i\tALL\tUsing where; Using join buffer (flat, BNL join)"]}
> raw
| span shop
| SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS who FROM orders o
< raw
| {"columns":["id","who"],"resolved":"/* select#1 */ select `shop`.`o`.`id` AS `id`,\u003cexpr_cache\u003e\u003c`shop`.`o`.`customer_id`\u003e((/* select#2 */ select `shop`.`c`.`name` from `shop`.`customers` `c` where `shop`.`c`.`id` = `shop`.`o`.`customer_id`)) AS `who` from `shop`.`orders` `o`","plan":["o\tALL\t","c\teq_ref\t"]}
> raw
| span shop
| SELECT id, (SELECT max(amount) FROM orders) AS top FROM customers
< raw
| {"columns":["id","top"],"resolved":"/* select#1 */ select `shop`.`customers`.`id` AS `id`,(/* select#2 */ select max(`shop`.`orders`.`amount`) from `shop`.`orders`) AS `top` from `shop`.`customers`","plan":["customers\tALL\t","orders\tALL\t"]}
> raw
| span shop
| SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE region = 'east')
< raw
| {"columns":["id"],"resolved":"select `shop`.`orders`.`id` AS `id` from `shop`.`customers` join `shop`.`orders` where `shop`.`customers`.`id` = `shop`.`orders`.`customer_id` and `shop`.`customers`.`region` = 'east'","plan":["orders\tALL\tUsing where","customers\teq_ref\tUsing where"]}
> raw
| span shop
| SELECT id FROM orders o WHERE EXISTS (SELECT 1 FROM items i WHERE i.order_id = o.id)
< raw
| {"columns":["id"],"resolved":"select `shop`.`o`.`id` AS `id` from `shop`.`orders` `o` semi join (`shop`.`items` `i`) where `shop`.`i`.`order_id` = `shop`.`o`.`id`","plan":["o\tALL\t","i\tALL\tUsing where; FirstMatch(o); Using join buffer (flat, BNL join)"]}
> raw
| span shop
| SELECT d.k, d.v FROM (SELECT id AS k, amount AS v FROM orders) d
< raw
| {"columns":["k","v"],"resolved":"select `shop`.`orders`.`id` AS `k`,`shop`.`orders`.`amount` AS `v` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT d.k FROM (SELECT id AS k FROM orders LIMIT 2) d
< raw
| {"columns":["k"],"resolved":"/* select#1 */ select `d`.`k` AS `k` from (/* select#2 */ select `shop`.`orders`.`id` AS `k` from `shop`.`orders` limit 2) `d`","plan":["\u003cderived2\u003e\tALL\t","orders\tALL\t"]}
> raw
| span shop
| SELECT x.total FROM (SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id) x
< raw
| {"columns":["total"],"resolved":"/* select#1 */ select `x`.`total` AS `total` from (/* select#2 */ select `shop`.`orders`.`customer_id` AS `customer_id`,sum(`shop`.`orders`.`amount`) AS `total` from `shop`.`orders` group by `shop`.`orders`.`customer_id`) `x`","plan":["\u003cderived2\u003e\tALL\t","orders\tALL\tUsing temporary; Using filesort"]}
> raw
| span shop
| SELECT d.a, d.b FROM (SELECT id, name FROM customers) AS d (a, b)
< raw
| {"columns":["a","b"],"resolved":"select `shop`.`customers`.`id` AS `a`,`shop`.`customers`.`name` AS `b` from `shop`.`customers`","plan":["customers\tALL\t"]}
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
| {"columns":["k"],"resolved":"select `shop`.`customers`.`id` AS `k` from `shop`.`customers`","plan":["customers\tALL\t"]}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers) SELECT id, name FROM c
< raw
| {"columns":["id","name"],"resolved":"with c as (select `shop`.`customers`.`id` AS `id`,`shop`.`customers`.`name` AS `name` from `shop`.`customers`)select `shop`.`customers`.`id` AS `id`,`shop`.`customers`.`name` AS `name` from `shop`.`customers`","plan":["customers\tALL\t"]}
> raw
| span shop
| WITH c (k, v) AS (SELECT id, amount FROM orders) SELECT k, v FROM c
< raw
| {"columns":["k","v"],"resolved":"with c(`k`,`v`) as (select `shop`.`orders`.`id` AS `k`,`shop`.`orders`.`amount` AS `v` from `shop`.`orders`)select `shop`.`orders`.`id` AS `k`,`shop`.`orders`.`amount` AS `v` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT * FROM b
< raw
| {"columns":["id","name"],"resolved":"with a as (select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`customer_id` AS `customer_id` from `shop`.`orders`), b as (select `shop`.`orders`.`id` AS `id`,`shop`.`c`.`name` AS `name` from `shop`.`orders` join `shop`.`customers` `c`)select `shop`.`orders`.`id` AS `id`,`shop`.`c`.`name` AS `name` from `shop`.`orders` join `shop`.`customers` `c` where `shop`.`c`.`id` = `shop`.`orders`.`customer_id`","plan":["orders\tALL\tUsing where","c\teq_ref\t"]}
> raw
| span shop
| WITH c AS (SELECT id FROM customers) SELECT x.id, y.id FROM c x JOIN c y ON y.id = x.id
< raw
| {"columns":["id","id"],"resolved":"with c as (select `shop`.`customers`.`id` AS `id` from `shop`.`customers`)select `shop`.`customers`.`id` AS `id`,`shop`.`customers`.`id` AS `id` from `shop`.`customers` join `shop`.`customers` where `shop`.`customers`.`id` = `shop`.`customers`.`id`","plan":["customers\tALL\t","customers\teq_ref\t"]}
> raw
| span shop
| WITH RECURSIVE r (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r
< raw
| {"columns":["n"],"resolved":"with recursive r(`n`) as (/* select#2 */ select 1 AS `n` union all /* select#3 */ select `r`.`n` + 1 AS `n + 1` from `r` where `r`.`n` \u003c 3)/* select#1 */ select `r`.`n` AS `n` from `r`","plan":["\u003cderived2\u003e\tALL\t","NULL\tNULL\tNo tables used","\u003cderived2\u003e\tALL\tUsing where","\u003cunion2,3\u003e\tALL\t"]}
> raw
| span shop
| WITH RECURSIVE r (id, name) AS (SELECT id, name FROM customers WHERE id = 1 UNION ALL SELECT c.id, c.name FROM customers c JOIN r ON c.id = r.id + 1) SELECT id, name FROM r
< raw
| {"columns":["id","name"],"resolved":"with recursive r(`id`,`name`) as (/* select#2 */ select `shop`.`customers`.`id` AS `id`,`shop`.`customers`.`name` AS `name` from `shop`.`customers` where 1 union all /* select#3 */ select `shop`.`c`.`id` AS `id`,`shop`.`c`.`name` AS `name` from `shop`.`customers` `c` join `r` where `shop`.`c`.`id` = `r`.`id` + 1)/* select#1 */ select `r`.`id` AS `id`,`r`.`name` AS `name` from `r`","plan":["\u003cderived2\u003e\tALL\t","customers\tconst\t","\u003cderived2\u003e\tALL\t","c\teq_ref\tUsing where","\u003cunion2,3\u003e\tALL\t"]}
> raw
| span shop
| WITH orders AS (SELECT id, name FROM customers) SELECT id, name FROM orders
< raw
| {"columns":["id","name"],"resolved":"with orders as (select `shop`.`customers`.`id` AS `id`,`shop`.`customers`.`name` AS `name` from `shop`.`customers`)select `shop`.`customers`.`id` AS `id`,`shop`.`customers`.`name` AS `name` from `shop`.`customers`","plan":["customers\tALL\t"]}
> raw
| span shop
| WITH c AS (SELECT id, amount FROM orders) SELECT id FROM c WHERE amount > (SELECT avg(amount) FROM c)
< raw
| {"columns":["id"],"resolved":"with c as (/* select#2 */ select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`amount` AS `amount` from `shop`.`orders`)/* select#1 */ select `shop`.`orders`.`id` AS `id` from `shop`.`orders` where `shop`.`orders`.`amount` \u003e (/* select#3 */ select avg(`shop`.`orders`.`amount`) from `shop`.`orders`)","plan":["orders\tALL\tUsing where","orders\tALL\t"]}
> raw
| span shop
| WITH RECURSIVE c AS (SELECT id FROM orders), orders AS (SELECT 1 AS id) SELECT id FROM c
< raw
| {"columns":["id"],"resolved":"with recursive c as (/* select#2 */ select `orders`.`id` AS `id` from `orders`), orders as (/* select#3 */ select 1 AS `id`)/* select#1 */ select 1 AS `id` from dual","plan":["\u003cderived3\u003e\tsystem\t","NULL\tNULL\tNo tables used"]}
> raw
| span shop
| SELECT id FROM orders UNION SELECT id FROM customers
< raw
| {"columns":["id"],"resolved":"/* select#1 */ select `shop`.`orders`.`id` AS `id` from `shop`.`orders` union /* select#2 */ select `shop`.`customers`.`id` AS `id` from `shop`.`customers`","plan":["orders\tALL\t","customers\tALL\t","\u003cunion1,2\u003e\tALL\t"]}
> raw
| span shop
| SELECT id, status FROM orders UNION ALL SELECT id, name FROM customers
< raw
| {"columns":["id","status"],"resolved":"/* select#1 */ select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`status` AS `status` from `shop`.`orders` union all /* select#2 */ select `shop`.`customers`.`id` AS `id`,`shop`.`customers`.`name` AS `name` from `shop`.`customers`","plan":["orders\tALL\t","customers\tALL\t"]}
> raw
| span shop
| SELECT name FROM customers INTERSECT SELECT status FROM orders
< raw
| {"columns":["name"],"resolved":"/* select#1 */ select `shop`.`customers`.`name` AS `name` from `shop`.`customers` intersect /* select#2 */ select `shop`.`orders`.`status` AS `status` from `shop`.`orders`","plan":["customers\tALL\t","orders\tALL\t","\u003cintersect1,2\u003e\tALL\t"]}
> raw
| span shop
| SELECT id FROM orders EXCEPT SELECT customer_id FROM orders
< raw
| {"columns":["id"],"resolved":"/* select#1 */ select `shop`.`orders`.`id` AS `id` from `shop`.`orders` except /* select#2 */ select `shop`.`orders`.`customer_id` AS `customer_id` from `shop`.`orders`","plan":["orders\tALL\t","orders\tALL\t","\u003cexcept1,2\u003e\tALL\t"]}
> raw
| span shop
| SELECT id AS k FROM orders UNION SELECT 1 ORDER BY k
< raw
| {"columns":["k"],"resolved":"/* select#1 */ select `shop`.`orders`.`id` AS `k` from `shop`.`orders` union /* select#2 */ select 1 AS `1` order by `k`","plan":["orders\tALL\t","NULL\tNULL\tNo tables used","\u003cunion1,2\u003e\tALL\tUsing filesort"]}
> raw
| span shop
| (SELECT id FROM orders LIMIT 1) UNION (SELECT id FROM customers LIMIT 1)
< raw
| {"columns":["id"],"resolved":"(/* select#1 */ select `shop`.`orders`.`id` AS `id` from `shop`.`orders` limit 1) union (/* select#2 */ select `shop`.`customers`.`id` AS `id` from `shop`.`customers` limit 1)","plan":["orders\tALL\t","customers\tALL\t","\u003cunion1,2\u003e\tALL\t"]}
> raw
| span shop
| SELECT u.id FROM (SELECT id FROM orders UNION SELECT id FROM customers) u
< raw
| {"columns":["id"],"resolved":"/* select#1 */ select `u`.`id` AS `id` from (/* select#2 */ select `shop`.`orders`.`id` AS `id` from `shop`.`orders` union /* select#3 */ select `shop`.`customers`.`id` AS `id` from `shop`.`customers`) `u`","plan":["\u003cderived2\u003e\tALL\t","orders\tALL\t","customers\tALL\t","\u003cunion2,3\u003e\tALL\t"]}
> raw
| span shop
| SELECT id, doubled FROM big_orders
< raw
| {"columns":["id","doubled"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`amount` * 2 AS `doubled` from `shop`.`orders` where `shop`.`orders`.`amount` \u003e 5","plan":["orders\tALL\tUsing where"],"views":["big_orders"]}
> raw
| span shop
| SELECT b.status, b.doubled FROM big_orders b
< raw
| {"columns":["status","doubled"],"resolved":"select `shop`.`orders`.`status` AS `status`,`shop`.`orders`.`amount` * 2 AS `doubled` from `shop`.`orders` where `shop`.`orders`.`amount` \u003e 5","plan":["orders\tALL\tUsing where"],"views":["b"]}
> raw
| span shop
| SELECT * FROM big_orders
< raw
| {"columns":["id","doubled","status","customer_id"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`amount` * 2 AS `doubled`,`shop`.`orders`.`status` AS `status`,`shop`.`orders`.`customer_id` AS `customer_id` from `shop`.`orders` where `shop`.`orders`.`amount` \u003e 5","plan":["orders\tALL\tUsing where"],"views":["big_orders"]}
> raw
| span shop
| SELECT customer_id, total, n FROM customer_totals
< raw
| {"columns":["customer_id","total","n"],"resolved":"/* select#1 */ select `customer_totals`.`customer_id` AS `customer_id`,`customer_totals`.`total` AS `total`,`customer_totals`.`n` AS `n` from `shop`.`customer_totals`","plan":["\u003cderived2\u003e\tALL\t","orders\tALL\tUsing temporary; Using filesort"],"views":["customer_totals"]}
> raw
| span shop
| SELECT * FROM order_names
< raw
| {"columns":["id","doubled","name"],"resolved":"select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`amount` * 2 AS `doubled`,`shop`.`c`.`name` AS `name` from `shop`.`orders` join `shop`.`customers` `c` where `shop`.`c`.`id` = `shop`.`orders`.`customer_id` and `shop`.`orders`.`amount` \u003e 5","plan":["orders\tALL\tUsing where","c\teq_ref\t"],"views":["b","order_names"]}
> raw
| span shop
| SELECT o.name, c.region FROM order_names o JOIN customers c ON c.name = o.name
< raw
| {"columns":["name","region"],"resolved":"select `shop`.`c`.`name` AS `name`,`shop`.`c`.`region` AS `region` from `shop`.`orders` join `shop`.`customers` `c` join `shop`.`customers` `c` where `shop`.`c`.`name` = `shop`.`c`.`name` and `shop`.`c`.`id` = `shop`.`orders`.`customer_id` and `shop`.`orders`.`amount` \u003e 5","plan":["orders\tALL\tUsing where","c\teq_ref\t","c\tALL\tUsing where; Using join buffer (flat, BNL join)"],"views":["b","o"]}
> raw
| span shop
| SELECT a, b FROM renamed
< raw
| {"columns":["a","b"],"resolved":"select `shop`.`customers`.`id` AS `a`,`shop`.`customers`.`name` AS `b` from `shop`.`customers`","plan":["customers\tALL\t"],"views":["renamed"]}
> raw
| span shop
| SELECT id FROM all_ids
< raw
| {"columns":["id"],"resolved":"/* select#1 */ select `all_ids`.`id` AS `id` from `shop`.`all_ids`","plan":["\u003cderived2\u003e\tALL\t","orders\tALL\t","customers\tALL\t","\u003cunion2,3\u003e\tALL\t"],"views":["all_ids"]}
> raw
| span shop
| SELECT * FROM staff_pay
< raw
| {"columns":["id","salary","dept"],"resolved":"select `hr`.`s`.`id` AS `id`,`hr`.`s`.`salary` AS `salary`,`d`.`name` AS `dept` from `hr`.`staff` `s` join `hr`.`depts` `d` where `d`.`id` = `hr`.`s`.`dept_id`","plan":["s\tALL\tUsing where","d\teq_ref\t"],"views":["staff_pay"]}
> raw
| span shop
| SELECT name FROM hr.rich
< raw
| {"columns":["name"],"resolved":"select `hr`.`staff`.`name` AS `name` from `hr`.`staff` where `hr`.`staff`.`salary` \u003e 120","plan":["staff\tALL\tUsing where"],"views":["rich"]}
> raw
| span shop
| SELECT VipName FROM hr.`Rich_Names`
< raw
| {"columns":["VipName"],"resolved":"select `hr`.`staff`.`name` AS `VipName` from `hr`.`staff` where `hr`.`staff`.`salary` \u003e 120","plan":["staff\tALL\tUsing where"],"views":["Rich_Names","rich"]}
> raw
| span shop
| SELECT b.id, t.total FROM big_orders b JOIN customer_totals t ON t.customer_id = b.customer_id
< raw
| {"columns":["id","total"],"resolved":"/* select#1 */ select `shop`.`orders`.`id` AS `id`,`t`.`total` AS `total` from `shop`.`orders` join `shop`.`customer_totals` `t` where `t`.`customer_id` = `shop`.`orders`.`customer_id` and `shop`.`orders`.`amount` \u003e 5","plan":["orders\tALL\tUsing where","\u003cderived3\u003e\tref\t","orders\tALL\tUsing temporary; Using filesort"],"views":["b","t"]}
> raw
| span shop
| SELECT d.doubled FROM (SELECT doubled FROM big_orders) d
< raw
| {"columns":["doubled"],"resolved":"select `shop`.`orders`.`amount` * 2 AS `doubled` from `shop`.`orders` where `shop`.`orders`.`amount` \u003e 5","plan":["orders\tALL\tUsing where"],"views":["big_orders"]}
> raw
| span shop
| WITH v AS (SELECT * FROM big_orders) SELECT doubled FROM v
< raw
| {"columns":["doubled"],"resolved":"with v as (select `shop`.`orders`.`id` AS `id`,`shop`.`orders`.`amount` * 2 AS `doubled`,`shop`.`orders`.`status` AS `status`,`shop`.`orders`.`customer_id` AS `customer_id` from `shop`.`orders`)select `shop`.`orders`.`amount` * 2 AS `doubled` from `shop`.`orders` where `shop`.`orders`.`amount` \u003e 5","plan":["orders\tALL\tUsing where"],"views":["big_orders"]}
> raw
| span shop
| SELECT id, row_number() OVER (ORDER BY amount) AS rn FROM orders
< raw
| {"columns":["id","rn"],"resolved":"select `shop`.`orders`.`id` AS `id`,row_number() over ( order by `shop`.`orders`.`amount`) AS `rn` from `shop`.`orders`","plan":["orders\tALL\tUsing temporary"]}
> raw
| span shop
| SELECT id, sum(amount) OVER (PARTITION BY customer_id) AS running FROM orders
< raw
| {"columns":["id","running"],"resolved":"select `shop`.`orders`.`id` AS `id`,sum(`shop`.`orders`.`amount`) over ( partition by `shop`.`orders`.`customer_id`) AS `running` from `shop`.`orders`","plan":["orders\tALL\tUsing temporary"]}
> raw
| span shop
| SELECT id, rank() OVER w AS r FROM orders WINDOW w AS (PARTITION BY status ORDER BY amount)
< raw
| {"columns":["id","r"],"resolved":"select `shop`.`orders`.`id` AS `id`,rank() over ( partition by `shop`.`orders`.`status` order by `shop`.`orders`.`amount`) AS `r` from `shop`.`orders`","plan":["orders\tALL\tUsing temporary"]}
> raw
| span shop
| SELECT id, lag(amount) OVER (ORDER BY id) AS prev FROM orders
< raw
| {"columns":["id","prev"],"resolved":"select `shop`.`orders`.`id` AS `id`,lag(`shop`.`orders`.`amount`,1) over ( order by `shop`.`orders`.`id`) AS `prev` from `shop`.`orders`","plan":["orders\tALL\tUsing temporary"]}
> raw
| span shop
| SELECT customer_id, sum(amount), count(*) FROM orders GROUP BY customer_id
< raw
| {"columns":["customer_id","sum(amount)","count(*)"],"resolved":"select `shop`.`orders`.`customer_id` AS `customer_id`,sum(`shop`.`orders`.`amount`) AS `sum(amount)`,count(0) AS `count(*)` from `shop`.`orders` group by `shop`.`orders`.`customer_id`","plan":["orders\tALL\tUsing temporary; Using filesort"]}
> raw
| span shop
| SELECT status, avg(amount) AS a FROM orders GROUP BY status HAVING count(*) > 1
< raw
| {"columns":["status","a"],"resolved":"select `shop`.`orders`.`status` AS `status`,avg(`shop`.`orders`.`amount`) AS `a` from `shop`.`orders` group by `shop`.`orders`.`status` having 1","plan":["orders\tALL\tUsing temporary; Using filesort"]}
> raw
| span shop
| SELECT region, group_concat(name ORDER BY name) AS names FROM customers GROUP BY region
< raw
| {"columns":["region","names"],"resolved":"select `shop`.`customers`.`region` AS `region`,group_concat(`shop`.`customers`.`name` order by `shop`.`customers`.`name` ASC separator ',') AS `names` from `shop`.`customers` group by `shop`.`customers`.`region`","plan":["customers\tALL\tUsing filesort"]}
> raw
| span shop
| SELECT max(amount) FROM orders WHERE status = 'open'
< raw
| {"columns":["max(amount)"],"resolved":"select max(`shop`.`orders`.`amount`) AS `max(amount)` from `shop`.`orders` where `shop`.`orders`.`status` = 'open'","plan":["orders\tALL\tUsing where"]}
> raw
| span shop
| SELECT count(DISTINCT customer_id) AS buyers FROM orders
< raw
| {"columns":["buyers"],"resolved":"select count(distinct `shop`.`orders`.`customer_id`) AS `buyers` from `shop`.`orders`","plan":["orders\tALL\t"]}
> raw
| span shop
| SELECT c.region, sum(o.amount) FROM customers c JOIN orders o ON o.customer_id = c.id GROUP BY c.region
< raw
| {"columns":["region","sum(o.amount)"],"resolved":"select `shop`.`c`.`region` AS `region`,sum(`shop`.`o`.`amount`) AS `sum(o.amount)` from `shop`.`customers` `c` join `shop`.`orders` `o` where `shop`.`c`.`id` = `shop`.`o`.`customer_id` group by `shop`.`c`.`region`","plan":["o\tALL\tUsing where; Using temporary; Using filesort","c\teq_ref\t"]}
> raw
| span shop
| SELECT customer_id, sum(amount) FROM orders GROUP BY customer_id WITH ROLLUP
< raw
| {"columns":["customer_id","sum(amount)"],"resolved":"select `shop`.`orders`.`customer_id` AS `customer_id`,sum(`shop`.`orders`.`amount`) AS `sum(amount)` from `shop`.`orders` group by `shop`.`orders`.`customer_id` with rollup","plan":["orders\tALL\tUsing filesort"]}
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
| {"columns":["name","salary"],"resolved":"select `hr`.`staff`.`name` AS `name`,`hr`.`staff`.`salary` AS `salary` from `hr`.`staff`","plan":["staff\tALL\t"]}
> raw
| span hr
| SELECT s.name, d.name AS dept FROM staff s JOIN depts d ON d.id = s.dept_id
< raw
| {"columns":["name","dept"],"resolved":"select `hr`.`s`.`name` AS `name`,`d`.`name` AS `dept` from `hr`.`staff` `s` join `hr`.`depts` `d` where `d`.`id` = `hr`.`s`.`dept_id`","plan":["s\tALL\tUsing where","d\teq_ref\t"]}
> raw
| span hr
| SELECT o.amount FROM shop.orders o
< raw
| {"columns":["amount"],"resolved":"select `shop`.`o`.`amount` AS `amount` from `shop`.`orders` `o`","plan":["o\tALL\t"]}
> raw
| span hr
| SELECT * FROM rich
< raw
| {"columns":["name","salary"],"resolved":"select `hr`.`staff`.`name` AS `name`,`hr`.`staff`.`salary` AS `salary` from `hr`.`staff` where `hr`.`staff`.`salary` \u003e 120","plan":["staff\tALL\tUsing where"],"views":["rich"]}
> raw
| span hr
| SELECT * FROM shop.staff_pay
< raw
| {"columns":["id","salary","dept"],"resolved":"select `hr`.`s`.`id` AS `id`,`hr`.`s`.`salary` AS `salary`,`d`.`name` AS `dept` from `hr`.`staff` `s` join `hr`.`depts` `d` where `d`.`id` = `hr`.`s`.`dept_id`","plan":["s\tALL\tUsing where","d\teq_ref\t"],"views":["staff_pay"]}
> raw
| span hr
| SELECT name FROM staff UNION SELECT name FROM shop.customers
< raw
| {"columns":["name"],"resolved":"/* select#1 */ select `hr`.`staff`.`name` AS `name` from `hr`.`staff` union /* select#2 */ select `shop`.`customers`.`name` AS `name` from `shop`.`customers`","plan":["staff\tALL\t","customers\tALL\t","\u003cunion1,2\u003e\tALL\t"]}
