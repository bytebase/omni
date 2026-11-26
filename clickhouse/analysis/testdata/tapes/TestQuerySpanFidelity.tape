-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestQuerySpanFidelity default#1
> raw
| fixture
| CREATE DATABASE shop
| CREATE DATABASE hr
| CREATE TABLE shop.customers (id UInt32, name String, email String, country String, created Date) ENGINE = MergeTree ORDER BY id
| CREATE TABLE shop.orders (id UInt32, customer_id UInt32, amount Decimal(10, 2), status String, tags Array(String), created DateTime) ENGINE = MergeTree ORDER BY id
| CREATE TABLE shop.items (order_id UInt32, sku String, qty UInt32, price Float64) ENGINE = ReplacingMergeTree ORDER BY (order_id, sku)
| CREATE VIEW shop.big_orders AS SELECT id, customer_id, amount FROM shop.orders WHERE amount > 100
| CREATE VIEW shop.customer_totals AS SELECT c.id AS customer_id, c.name AS name, sum(o.amount) AS total FROM shop.customers AS c INNER JOIN shop.orders AS o ON o.customer_id = c.id GROUP BY c.id, c.name
| CREATE VIEW shop.vip AS SELECT customer_id, upper(name) AS name FROM shop.customer_totals WHERE total > 1000
| CREATE TABLE hr.employees (id UInt32, name String, salary UInt64, dept_id UInt32, manager_id UInt32) ENGINE = MergeTree ORDER BY id
| CREATE TABLE hr.depts (id UInt32, title String) ENGINE = MergeTree ORDER BY id
| CREATE VIEW hr.payroll AS SELECT e.name AS name, d.title AS dept, e.salary * 12 AS yearly FROM hr.employees AS e LEFT JOIN hr.depts AS d ON d.id = e.dept_id
< raw
| {"databases":[{"name":"hr","tables":[{"name":"depts","columns":["id","title"]},{"name":"employees","columns":["id","name","salary","dept_id","manager_id"]}],"views":[{"name":"payroll","columns":["name","dept","yearly"],"definition":"CREATE VIEW hr.payroll (`name` String, `dept` String, `yearly` UInt64) AS SELECT e.name AS name, d.title AS dept, e.salary * 12 AS yearly FROM hr.employees AS e LEFT JOIN hr.depts AS d ON d.id = e.dept_id"}]},{"name":"shop","tables":[{"name":"customers","columns":["id","name","email","country","created"]},{"name":"items","columns":["order_id","sku","qty","price"]},{"name":"orders","columns":["id","customer_id","amount","status","tags","created"]}],"views":[{"name":"big_orders","columns":["id","customer_id","amount"],"definition":"CREATE VIEW shop.big_orders (`id` UInt32, `customer_id` UInt32, `amount` Decimal(10, 2)) AS SELECT id, customer_id, amount FROM shop.orders WHERE amount \u003e 100"},{"name":"customer_totals","columns":["customer_id","name","total"],"definition":"CREATE VIEW shop.customer_totals (`customer_id` UInt32, `name` String, `total` Decimal(38, 2)) AS SELECT c.id AS customer_id, c.name AS name, sum(o.amount) AS total FROM shop.customers AS c INNER JOIN shop.orders AS o ON o.customer_id = c.id GROUP BY c.id, c.name"},{"name":"vip","columns":["customer_id","name"],"definition":"CREATE VIEW shop.vip (`customer_id` UInt32, `name` String) AS SELECT customer_id, upper(name) AS name FROM shop.customer_totals WHERE total \u003e 1000"}]}]}
> raw
| span shop
| SELECT id, name FROM customers
< raw
| {"columns":["id","name"],"lineage":[["shop.customers.id"],["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT * FROM customers
< raw
| {"columns":["id","name","email","country","created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT customers.* FROM customers
< raw
| {"columns":["id","name","email","country","created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT c.* FROM customers AS c
< raw
| {"columns":["id","name","email","country","created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT shop.customers.name FROM customers
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT c.name, c.email FROM customers AS c
< raw
| {"columns":["name","email"],"lineage":[["shop.customers.name"],["shop.customers.email"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT name AS customer, email AS mail FROM customers
< raw
| {"columns":["customer","mail"],"lineage":[["shop.customers.name"],["shop.customers.email"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT id FROM shop.orders
< raw
| {"columns":["id"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT amount, status FROM shop.orders WHERE status = 'paid'
< raw
| {"columns":["amount","status"],"lineage":[["shop.orders.amount"],["shop.orders.status"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT DISTINCT country FROM customers
< raw
| {"columns":["country"],"lineage":[["shop.customers.country"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT id, name FROM customers ORDER BY name LIMIT 10
< raw
| {"columns":["id","name"],"lineage":[["shop.customers.id"],["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT id FROM customers LIMIT 5 OFFSET 2
< raw
| {"columns":["id"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT name FROM customers WHERE id IN (1, 2, 3)
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT * FROM items
< raw
| {"columns":["order_id","sku","qty","price"],"lineage":[["shop.items.order_id"],["shop.items.sku"],["shop.items.qty"],["shop.items.price"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT o.* FROM shop.orders o
< raw
| {"columns":["id","customer_id","amount","status","tags","created"],"lineage":[["shop.orders.id"],["shop.orders.customer_id"],["shop.orders.amount"],["shop.orders.status"],["shop.orders.tags"],["shop.orders.created"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT `name`, `email` FROM `customers`
< raw
| {"columns":["name","email"],"lineage":[["shop.customers.name"],["shop.customers.email"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT 1
< raw
| {"columns":["1"],"lineage":[null]}
> raw
| span shop
| SELECT 1 AS one, 'x' AS letter
< raw
| {"columns":["one","letter"],"lineage":[null,null]}
> raw
| span shop
| SELECT now() AS ts
< raw
| {"columns":["ts"],"lineage":[null]}
> raw
| span shop
| SELECT count() FROM customers
< raw
| {"columns":["count()"],"lineage":[null],"tables":["shop.customers"]}
> raw
| span shop
| SELECT count(*) FROM orders
< raw
| {"columns":["count()"],"lineage":[null],"tables":["shop.orders"]}
> raw
| span shop
| SELECT amount + 1 FROM orders
< raw
| {"columns":["plus(amount, 1)"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT amount * 2, amount / 3, amount - 1 FROM orders
< raw
| {"columns":["multiply(amount, 2)","divide(amount, 3)","minus(amount, 1)"],"lineage":[["shop.orders.amount"],["shop.orders.amount"],["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id % 2 FROM orders
< raw
| {"columns":["modulo(id, 2)"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT -amount FROM orders
< raw
| {"columns":["negate(amount)"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT NOT id FROM orders
< raw
| {"columns":["not(id)"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id = 1 AND status != 'x' FROM orders
< raw
| {"columns":["and(equals(id, 1), notEquals(status, 'x'))"],"lineage":[["shop.orders.id","shop.orders.status"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id = 1 OR id = 2 OR id = 3 FROM orders
< raw
| {"columns":["or(equals(id, 1), equals(id, 2), equals(id, 3))"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id BETWEEN 1 AND 10 FROM orders
< raw
| {"columns":["and(greaterOrEquals(id, 1), lessOrEquals(id, 10))"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id NOT BETWEEN 1 AND 10 FROM orders
< raw
| {"columns":["or(less(id, 1), greater(id, 10))"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT CASE WHEN amount > 100 THEN 'big' ELSE 'small' END FROM orders
< raw
| {"columns":["multiIf(greater(amount, 100), 'big', 'small')"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT CASE status WHEN 'paid' THEN 1 ELSE 0 END FROM orders
< raw
| {"columns":["caseWithExpression(status, 'paid', 1, 0)"],"lineage":[["shop.orders.status"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT amount > 100 ? 'big' : 'small' FROM orders
< raw
| {"columns":["if(greater(amount, 100), 'big', 'small')"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT if(amount > 100, id, customer_id) FROM orders
< raw
| {"columns":["if(greater(amount, 100), id, customer_id)"],"lineage":[["shop.orders.amount","shop.orders.customer_id","shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT [id, customer_id] FROM orders
< raw
| {"columns":["[id, customer_id]"],"lineage":[["shop.orders.customer_id","shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT (id, status) FROM orders
< raw
| {"columns":["(id, status)"],"lineage":[["shop.orders.id","shop.orders.status"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT tags[1] FROM orders
< raw
| {"columns":["arrayElement(tags, 1)"],"lineage":[["shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT CAST(id AS String) FROM orders
< raw
| {"columns":["CAST(id, 'String')"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id::Int64 FROM orders
< raw
| {"columns":["CAST(id, 'Int64')"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT toString(id) || status FROM orders
< raw
| {"columns":["concat(toString(id), status)"],"lineage":[["shop.orders.id","shop.orders.status"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT concat(name, ' <', email, '>') AS contact FROM customers
< raw
| {"columns":["contact"],"lineage":[["shop.customers.email","shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT lower(name), upper(country) FROM customers
< raw
| {"columns":["lower(name)","upper(country)"],"lineage":[["shop.customers.name"],["shop.customers.country"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT length(name) > 3 FROM customers
< raw
| {"columns":["greater(length(name), 3)"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT name LIKE 'A%' FROM customers
< raw
| {"columns":["like(name, 'A%')"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT name NOT LIKE 'A%' FROM customers
< raw
| {"columns":["notLike(name, 'A%')"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT name ILIKE 'a%' FROM customers
< raw
| {"columns":["ilike(name, 'a%')"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT email IS NULL, email IS NOT NULL FROM customers
< raw
| {"columns":["isNull(email)","isNotNull(email)"],"lineage":[["shop.customers.email"],["shop.customers.email"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT id IN (1, 2) FROM customers
< raw
| {"columns":["in(id, (1, 2))"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT toYear(created), toMonth(created) FROM customers
< raw
| {"columns":["toYear(created)","toMonth(created)"],"lineage":[["shop.customers.created"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT EXTRACT(YEAR FROM created) FROM customers
< raw
| {"columns":["toYear(created)"],"lineage":[["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT created + INTERVAL 1 DAY FROM customers
< raw
| {"columns":["plus(created, toIntervalDay(1))"],"lineage":[["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT substring(name, 1, 3) FROM customers
< raw
| {"columns":["substring(name, 1, 3)"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT 1.5, 'it''s', NULL, true FROM customers
< raw
| {"columns":["1.5","'it\\'s'","NULL","true"],"lineage":[null,null,null,null],"tables":["shop.customers"]}
> raw
| span shop
| SELECT -1, 0x1F, 1e3 FROM customers
< raw
| {"columns":["-1","31","1000."],"lineage":[null,null,null],"tables":["shop.customers"]}
> raw
| span shop
| SELECT coalesce(email, name) FROM customers
< raw
| {"columns":["coalesce(email, name)"],"lineage":[["shop.customers.email","shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT multiIf(qty > 10, 'many', qty > 1, 'some', 'one') FROM items
< raw
| {"columns":["multiIf(greater(qty, 10), 'many', greater(qty, 1), 'some', 'one')"],"lineage":[["shop.items.qty"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT round(price * qty, 2) AS total FROM items
< raw
| {"columns":["total"],"lineage":[["shop.items.price","shop.items.qty"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT (price + 1) * qty FROM items
< raw
| {"columns":["multiply(plus(price, 1), qty)"],"lineage":[["shop.items.price","shop.items.qty"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT tuple(id, sku).1 FROM items
< raw
| {"error":"UNKNOWN_IDENTIFIER"}
> raw
| span shop
| SELECT tuple(*) FROM customers
< raw
| {"columns":["(id, name, email, country, created)"],"lineage":[["shop.customers.country","shop.customers.created","shop.customers.email","shop.customers.id","shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT tuple(customers.*) FROM customers
< raw
| {"columns":["(id, name, email, country, created)"],"lineage":[["shop.customers.country","shop.customers.created","shop.customers.email","shop.customers.id","shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT count(), sum(amount) FROM orders
< raw
| {"columns":["count()","sum(amount)"],"lineage":[null,["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT COUNT(id) FROM orders
< raw
| {"columns":["COUNT(id)"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT count(DISTINCT customer_id) FROM orders
< raw
| {"columns":["countDistinct(customer_id)"],"lineage":[["shop.orders.customer_id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT customer_id, sum(amount) FROM orders GROUP BY customer_id
< raw
| {"columns":["customer_id","sum(amount)"],"lineage":[["shop.orders.customer_id"],["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id HAVING total > 10
< raw
| {"columns":["customer_id","total"],"lineage":[["shop.orders.customer_id"],["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT status, max(amount), min(amount), avg(amount) FROM orders GROUP BY status
< raw
| {"columns":["status","max(amount)","min(amount)","avg(amount)"],"lineage":[["shop.orders.status"],["shop.orders.amount"],["shop.orders.amount"],["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT uniqExact(customer_id) FROM orders
< raw
| {"columns":["uniqExact(customer_id)"],"lineage":[["shop.orders.customer_id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT sumIf(amount, status = 'paid') FROM orders
< raw
| {"columns":["sumIf(amount, equals(status, 'paid'))"],"lineage":[["shop.orders.amount","shop.orders.status"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT quantile(0.5)(amount) FROM orders
< raw
| {"columns":["quantile(0.5)(amount)"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT groupArray(sku) FROM items
< raw
| {"columns":["groupArray(sku)"],"lineage":[["shop.items.sku"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT argMax(sku, price) FROM items
< raw
| {"columns":["argMax(sku, price)"],"lineage":[["shop.items.price","shop.items.sku"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT order_id, sum(price * qty) AS total FROM items GROUP BY order_id ORDER BY total DESC
< raw
| {"columns":["order_id","total"],"lineage":[["shop.items.order_id"],["shop.items.price","shop.items.qty"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT country, count() FROM customers GROUP BY country WITH TOTALS
< raw
| {"columns":["country","count()"],"lineage":[["shop.customers.country"],null],"tables":["shop.customers"]}
> raw
| span shop
| SELECT country, count() FROM customers GROUP BY ROLLUP(country)
< raw
| {"columns":["country","count()"],"lineage":[["shop.customers.country"],null],"tables":["shop.customers"]}
> raw
| span shop
| SELECT sum(amount) * 2 FROM orders
< raw
| {"columns":["multiply(sum(amount), 2)"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT status, count() FROM orders GROUP BY ALL
< raw
| {"columns":["status","count()"],"lineage":[["shop.orders.status"],null],"tables":["shop.orders"]}
> raw
| span shop
| SELECT c.name, o.amount FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.name, o.amount FROM customers c INNER JOIN orders o ON o.customer_id = c.id WHERE o.amount > 10
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.name, o.amount FROM customers AS c LEFT JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.name, o.amount FROM customers AS c RIGHT JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.name, o.amount FROM customers AS c FULL JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.name, o.amount FROM customers AS c CROSS JOIN orders AS o
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.name, o.amount FROM customers AS c, orders AS o WHERE o.customer_id = c.id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["id","name","email","country","created","o.id","customer_id","amount","status","tags","o.created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"],["shop.orders.id"],["shop.orders.customer_id"],["shop.orders.amount"],["shop.orders.status"],["shop.orders.tags"],["shop.orders.created"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.*, o.amount FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["id","name","email","country","created","amount"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["id"],"lineage":[["shop.customers.id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.id, o.id FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["id","o.id"],"lineage":[["shop.customers.id"],["shop.orders.id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT o.id, c.id FROM customers AS c JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["o.id","id"],"lineage":[["shop.orders.id"],["shop.customers.id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name, amount FROM customers JOIN orders ON orders.customer_id = customers.id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM orders JOIN items ON items.order_id = orders.id
< raw
| {"columns":["id","customer_id","amount","status","tags","created","order_id","sku","qty","price"],"lineage":[["shop.orders.id"],["shop.orders.customer_id"],["shop.orders.amount"],["shop.orders.status"],["shop.orders.tags"],["shop.orders.created"],["shop.items.order_id"],["shop.items.sku"],["shop.items.qty"],["shop.items.price"]],"tables":["shop.items","shop.orders"]}
> raw
| span shop
| SELECT o.id, i.sku, c.name FROM orders AS o JOIN items AS i ON i.order_id = o.id JOIN customers AS c ON c.id = o.customer_id
< raw
| {"columns":["o.id","sku","name"],"lineage":[["shop.orders.id"],["shop.items.sku"],["shop.customers.name"]],"tables":["shop.customers","shop.items","shop.orders"]}
> raw
| span shop
| SELECT * FROM orders AS o JOIN items AS i ON i.order_id = o.id JOIN customers AS c ON c.id = o.customer_id
< raw
| {"columns":["o.id","customer_id","amount","status","tags","o.created","order_id","sku","qty","price","c.id","name","email","country","c.created"],"lineage":[["shop.orders.id"],["shop.orders.customer_id"],["shop.orders.amount"],["shop.orders.status"],["shop.orders.tags"],["shop.orders.created"],["shop.items.order_id"],["shop.items.sku"],["shop.items.qty"],["shop.items.price"],["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers","shop.items","shop.orders"]}
> raw
| span shop
| SELECT id FROM orders AS o JOIN items AS i ON i.order_id = o.id JOIN customers AS c ON c.id = o.customer_id
< raw
| {"error":"AMBIGUOUS_IDENTIFIER"}
> raw
| span shop
| SELECT sku FROM orders AS o JOIN items AS i ON i.order_id = o.id JOIN customers AS c ON c.id = o.customer_id
< raw
| {"columns":["sku"],"lineage":[["shop.items.sku"]],"tables":["shop.customers","shop.items","shop.orders"]}
> raw
| span shop
| SELECT * FROM orders AS o JOIN (SELECT order_id AS id, sku FROM items) AS i USING (id)
< raw
| {"columns":["id","customer_id","amount","status","tags","created","sku"],"lineage":[["shop.orders.id"],["shop.orders.customer_id"],["shop.orders.amount"],["shop.orders.status"],["shop.orders.tags"],["shop.orders.created"],["shop.items.sku"]],"tables":["shop.items","shop.orders"]}
> raw
| span shop
| SELECT id, sku FROM orders AS o JOIN (SELECT order_id AS id, sku FROM items) AS i USING id
< raw
| {"columns":["id","sku"],"lineage":[["shop.orders.id"],["shop.items.sku"]],"tables":["shop.items","shop.orders"]}
> raw
| span shop
| SELECT a.name, b.name FROM customers AS a JOIN customers AS b ON a.id = b.id
< raw
| {"columns":["name","b.name"],"lineage":[["shop.customers.name"],["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT * FROM customers AS a, customers AS b
< raw
| {"columns":["id","name","email","country","created","b.id","b.name","b.email","b.country","b.created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"],["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT b.id FROM customers AS a JOIN customers AS b USING (id)
< raw
| {"columns":["b.id"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT c.name FROM customers AS c LEFT SEMI JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.name FROM customers AS c LEFT ANTI JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.name, o.amount FROM customers AS c ANY LEFT JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT o.id, i.qty FROM orders AS o ASOF JOIN items AS i ON i.order_id = o.id AND i.qty <= o.id
< raw
| {"columns":["id","qty"],"lineage":[["shop.orders.id"],["shop.items.qty"]],"tables":["shop.items","shop.orders"]}
> raw
| span shop
| SELECT c.name, o.amount FROM customers AS c GLOBAL JOIN orders AS o ON o.customer_id = c.id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT c.name, e.salary FROM customers AS c JOIN hr.employees AS e ON e.id = c.id
< raw
| {"columns":["name","salary"],"lineage":[["shop.customers.name"],["hr.employees.salary"]],"tables":["hr.employees","shop.customers"]}
> raw
| span shop
| SELECT x FROM (SELECT id AS x FROM customers)
< raw
| {"columns":["x"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT t.x FROM (SELECT id AS x FROM customers) AS t
< raw
| {"columns":["x"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT * FROM (SELECT id, name FROM customers) AS t
< raw
| {"columns":["id","name"],"lineage":[["shop.customers.id"],["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT total FROM (SELECT customer_id, sum(amount) AS total FROM orders GROUP BY customer_id)
< raw
| {"columns":["total"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT n FROM (SELECT upper(name) AS n FROM customers) WHERE n != ''
< raw
| {"columns":["n"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT a, b FROM (SELECT id, name FROM customers) AS t(a, b)
< raw
| {"columns":["a","b"],"lineage":[["shop.customers.id"],["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT y FROM (SELECT x AS y FROM (SELECT amount AS x FROM orders))
< raw
| {"columns":["y"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT * FROM (SELECT * FROM (SELECT id, email FROM customers))
< raw
| {"columns":["id","email"],"lineage":[["shop.customers.id"],["shop.customers.email"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT t.amount FROM (SELECT * FROM orders) AS t
< raw
| {"columns":["amount"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT `plus(id, 1)` FROM (SELECT id + 1 FROM customers)
< raw
| {"columns":["plus(id, 1)"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT s.cnt FROM (SELECT count() AS cnt FROM items) AS s
< raw
| {"columns":["cnt"],"lineage":[null],"tables":["shop.items"]}
> raw
| span shop
| SELECT `a.b` FROM (SELECT 1 AS `a.b`)
< raw
| {"columns":["a.b"],"lineage":[null]}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers) SELECT name FROM c
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers) SELECT * FROM c
< raw
| {"columns":["id","name"],"lineage":[["shop.customers.id"],["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers), o AS (SELECT customer_id, amount FROM orders) SELECT c.name, o.amount FROM c JOIN o ON o.customer_id = c.id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customers.name"],["shop.orders.amount"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| WITH big AS (SELECT * FROM orders WHERE amount > 100) SELECT id, amount FROM big
< raw
| {"columns":["id","amount"],"lineage":[["shop.orders.id"],["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| WITH t AS (SELECT sum(amount) AS total FROM orders) SELECT total * 2 FROM t
< raw
| {"columns":["multiply(total, 2)"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| WITH a AS (SELECT id FROM customers), b AS (SELECT id FROM a) SELECT id FROM b
< raw
| {"columns":["id"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| WITH c AS (SELECT name FROM customers) SELECT x.name FROM c AS x
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| WITH 10 AS threshold SELECT id FROM orders WHERE amount > threshold
< raw
| {"columns":["id"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| WITH 10 AS threshold SELECT threshold, id FROM orders
< raw
| {"columns":["threshold","id"],"lineage":[null,["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| WITH amount * 2 AS doubled SELECT doubled FROM orders
< raw
| {"columns":["doubled"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| WITH (SELECT max(amount) FROM orders) AS top SELECT id, top FROM orders
< raw
| {"columns":["id","top"],"lineage":[["shop.orders.id"],["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| WITH upper(name) AS n SELECT n, email FROM customers
< raw
| {"columns":["n","email"],"lineage":[["shop.customers.name"],["shop.customers.email"]],"tables":["shop.customers"]}
> raw
| span shop
| WITH x -> concat(x, email) AS f SELECT f(name) FROM customers
< raw
| {"columns":["f(name)"],"lineage":[["shop.customers.email","shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers) SELECT name FROM c UNION ALL SELECT sku FROM items
< raw
| {"columns":["name"],"lineage":[["shop.customers.name","shop.items.sku"]],"tables":["shop.customers","shop.items"]}
> raw
| span shop
| WITH customers AS (SELECT order_id AS id FROM items) SELECT id FROM customers
< raw
| {"columns":["id"],"lineage":[["shop.items.order_id"]],"tables":["shop.items"]}
> raw
| span shop
| WITH c AS (SELECT id, name FROM customers) SELECT count() FROM c WHERE name != ''
< raw
| {"columns":["count()"],"lineage":[null],"tables":["shop.customers"]}
> raw
| span shop
| WITH 1 AS x, c AS (SELECT x) SELECT * FROM c
< raw
| {"columns":["c.x"],"lineage":[null]}
> raw
| span shop
| WITH 2 AS k, c AS (SELECT id * k AS n FROM customers) SELECT n FROM c
< raw
| {"columns":["n"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| WITH 1 AS x SELECT * FROM (SELECT x)
< raw
| {"columns":["x"],"lineage":[null]}
> raw
| span shop
| WITH 3 AS k SELECT n FROM (SELECT id * k AS n FROM customers)
< raw
| {"columns":["n"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| WITH RECURSIVE r AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM r WHERE n < 3) SELECT n FROM r
< raw
| {"columns":["n"],"lineage":[["~"]]}
> raw
| span shop
| WITH RECURSIVE r AS (SELECT toString(id) AS a, email AS b, toUInt64(0) AS n FROM customers UNION ALL SELECT a, a AS b, n + 1 FROM r WHERE n < 1) SELECT b FROM r
< raw
| {"columns":["b"],"lineage":[["shop.customers.email","~"]],"tables":["shop.customers"]}
> raw
| span shop
| WITH RECURSIVE r AS (SELECT id, customer_id AS parent FROM orders UNION ALL SELECT o.id, r.parent FROM orders AS o JOIN r ON o.customer_id = r.id) SELECT * FROM r
< raw
| {"columns":["id","parent"],"lineage":[["shop.orders.id"],["shop.orders.customer_id","~"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id FROM customers UNION ALL SELECT id FROM orders
< raw
| {"columns":["id"],"lineage":[["shop.customers.id","shop.orders.id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name FROM customers UNION DISTINCT SELECT sku FROM items
< raw
| {"columns":["name"],"lineage":[["shop.customers.name","shop.items.sku"]],"tables":["shop.customers","shop.items"]}
> raw
| span shop
| SELECT id, name FROM customers UNION ALL SELECT id, status FROM orders
< raw
| {"columns":["id","name"],"lineage":[["shop.customers.id","shop.orders.id"],["shop.customers.name","shop.orders.status"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM customers INTERSECT SELECT customer_id FROM orders
< raw
| {"columns":["id"],"lineage":[["shop.customers.id","shop.orders.customer_id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM customers EXCEPT SELECT customer_id FROM orders
< raw
| {"columns":["id"],"lineage":[["shop.customers.id","shop.orders.customer_id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM customers UNION ALL SELECT id FROM orders UNION ALL SELECT order_id FROM items
< raw
| {"columns":["id"],"lineage":[["shop.customers.id","shop.items.order_id","shop.orders.id"]],"tables":["shop.customers","shop.items","shop.orders"]}
> raw
| span shop
| SELECT x FROM (SELECT id AS x FROM customers UNION ALL SELECT customer_id FROM orders)
< raw
| {"columns":["x"],"lineage":[["shop.customers.id","shop.orders.customer_id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM (SELECT name FROM customers UNION ALL SELECT status FROM orders) AS u
< raw
| {"columns":["name"],"lineage":[["shop.customers.name","shop.orders.status"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| (SELECT id FROM customers) UNION ALL (SELECT id FROM orders)
< raw
| {"columns":["id"],"lineage":[["shop.customers.id","shop.orders.id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT 1 AS n UNION ALL SELECT id FROM customers
< raw
| {"columns":["n"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT name FROM customers UNION ALL SELECT name FROM hr.employees
< raw
| {"columns":["name"],"lineage":[["hr.employees.name","shop.customers.name"]],"tables":["hr.employees","shop.customers"]}
> raw
| span shop
| SELECT (SELECT max(amount) FROM orders)
< raw
| {"columns":["_subquery_1"],"lineage":[["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT (SELECT max(amount) FROM orders) AS top, name FROM customers
< raw
| {"columns":["top","name"],"lineage":[["shop.orders.amount"],["shop.customers.name"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders)
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name, id IN (SELECT customer_id FROM orders) AS has_orders FROM customers
< raw
| {"columns":["name","has_orders"],"lineage":[["shop.customers.name"],["shop.customers.id","shop.orders.customer_id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name FROM customers WHERE id NOT IN (SELECT customer_id FROM orders WHERE amount > 10)
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name FROM customers WHERE EXISTS (SELECT 1 FROM orders)
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name, NOT EXISTS (SELECT 1 FROM orders WHERE amount > 0) AS none FROM customers
< raw
| {"columns":["name","none"],"lineage":[["shop.customers.name"],null],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name FROM customers WHERE id GLOBAL IN (SELECT customer_id FROM orders)
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| WITH ids AS (SELECT customer_id FROM orders) SELECT name FROM customers WHERE id IN ids
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| WITH ids AS (SELECT customer_id FROM orders) SELECT name, id IN ids AS has_orders FROM customers
< raw
| {"columns":["name","has_orders"],"lineage":[["shop.customers.name"],["shop.customers.id","shop.orders.customer_id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name FROM customers WHERE id IN nope
< raw
| {"error":"UNKNOWN_IDENTIFIER"}
> raw
| span shop
| SELECT id, (SELECT count() FROM items) AS n FROM orders
< raw
| {"columns":["id","n"],"lineage":[["shop.orders.id"],null],"tables":["shop.items","shop.orders"]}
> raw
| span shop
| SELECT name FROM customers WHERE id = (SELECT min(customer_id) FROM orders)
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM orders WHERE customer_id IN (SELECT id FROM customers WHERE country IN (SELECT country FROM customers WHERE id = 1))
< raw
| {"columns":["id"],"lineage":[["shop.orders.id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM big_orders
< raw
| {"columns":["id","customer_id","amount"],"lineage":[["shop.big_orders.id","shop.orders.id"],["shop.big_orders.customer_id","shop.orders.customer_id"],["shop.big_orders.amount","shop.orders.amount"]],"tables":["shop.big_orders","shop.orders"]}
> raw
| span shop
| SELECT id, amount FROM big_orders
< raw
| {"columns":["id","amount"],"lineage":[["shop.big_orders.id","shop.orders.id"],["shop.big_orders.amount","shop.orders.amount"]],"tables":["shop.big_orders","shop.orders"]}
> raw
| span shop
| SELECT amount * 2 FROM big_orders
< raw
| {"columns":["multiply(amount, 2)"],"lineage":[["shop.big_orders.amount","shop.orders.amount"]],"tables":["shop.big_orders","shop.orders"]}
> raw
| span shop
| SELECT * FROM customer_totals
< raw
| {"columns":["customer_id","name","total"],"lineage":[["shop.customer_totals.customer_id","shop.customers.id"],["shop.customer_totals.name","shop.customers.name"],["shop.customer_totals.total","shop.orders.amount"]],"tables":["shop.customer_totals","shop.customers","shop.orders"]}
> raw
| span shop
| SELECT name, total FROM customer_totals
< raw
| {"columns":["name","total"],"lineage":[["shop.customer_totals.name","shop.customers.name"],["shop.customer_totals.total","shop.orders.amount"]],"tables":["shop.customer_totals","shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM vip
< raw
| {"columns":["customer_id","name"],"lineage":[["shop.customer_totals.customer_id","shop.customers.id","shop.vip.customer_id"],["shop.customer_totals.name","shop.customers.name","shop.vip.name"]],"tables":["shop.customer_totals","shop.customers","shop.orders","shop.vip"]}
> raw
| span shop
| SELECT name FROM vip
< raw
| {"columns":["name"],"lineage":[["shop.customer_totals.name","shop.customers.name","shop.vip.name"]],"tables":["shop.customer_totals","shop.customers","shop.orders","shop.vip"]}
> raw
| span shop
| SELECT v.name, o.amount FROM vip AS v JOIN orders AS o ON o.customer_id = v.customer_id
< raw
| {"columns":["name","amount"],"lineage":[["shop.customer_totals.name","shop.customers.name","shop.vip.name"],["shop.orders.amount"]],"tables":["shop.customer_totals","shop.customers","shop.orders","shop.vip"]}
> raw
| span shop
| SELECT b.id FROM big_orders AS b JOIN items AS i ON i.order_id = b.id
< raw
| {"columns":["id"],"lineage":[["shop.big_orders.id","shop.orders.id"]],"tables":["shop.big_orders","shop.items","shop.orders"]}
> raw
| span shop
| SELECT count() FROM big_orders
< raw
| {"columns":["count()"],"lineage":[null],"tables":["shop.big_orders","shop.orders"]}
> raw
| span shop
| SELECT customer_id FROM big_orders UNION ALL SELECT customer_id FROM customer_totals
< raw
| {"columns":["customer_id"],"lineage":[["shop.big_orders.customer_id","shop.customer_totals.customer_id","shop.customers.id","shop.orders.customer_id"]],"tables":["shop.big_orders","shop.customer_totals","shop.customers","shop.orders"]}
> raw
| span shop
| WITH v AS (SELECT * FROM vip) SELECT name FROM v
< raw
| {"columns":["name"],"lineage":[["shop.customer_totals.name","shop.customers.name","shop.vip.name"]],"tables":["shop.customer_totals","shop.customers","shop.orders","shop.vip"]}
> raw
| span shop
| SELECT name FROM (SELECT * FROM customer_totals) WHERE total > 0
< raw
| {"columns":["name"],"lineage":[["shop.customer_totals.name","shop.customers.name"]],"tables":["shop.customer_totals","shop.customers","shop.orders"]}
> raw
| span shop
| SELECT * FROM hr.payroll
< raw
| {"columns":["name","dept","yearly"],"lineage":[["hr.employees.name","hr.payroll.name"],["hr.depts.title","hr.payroll.dept"],["hr.employees.salary","hr.payroll.yearly"]],"tables":["hr.depts","hr.employees","hr.payroll"]}
> raw
| span shop
| SELECT dept, sum(yearly) FROM hr.payroll GROUP BY dept
< raw
| {"columns":["dept","sum(yearly)"],"lineage":[["hr.depts.title","hr.payroll.dept"],["hr.employees.salary","hr.payroll.yearly"]],"tables":["hr.depts","hr.employees","hr.payroll"]}
> raw
| span shop
| SELECT id FROM customers WHERE id IN (SELECT customer_id FROM vip)
< raw
| {"columns":["id"],"lineage":[["shop.customers.id"]],"tables":["shop.customer_totals","shop.customers","shop.orders","shop.vip"]}
> raw
| span shop
| SELECT id, tag FROM orders ARRAY JOIN tags AS tag
< raw
| {"columns":["id","tag"],"lineage":[["shop.orders.id"],["shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, tags FROM orders ARRAY JOIN tags
< raw
| {"columns":["id","tags"],"lineage":[["shop.orders.id"],["shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, tag FROM orders LEFT ARRAY JOIN tags AS tag
< raw
| {"columns":["id","tag"],"lineage":[["shop.orders.id"],["shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT * FROM orders ARRAY JOIN tags
< raw
| {"columns":["id","customer_id","amount","status","tags","created"],"lineage":[["shop.orders.id"],["shop.orders.customer_id"],["shop.orders.amount"],["shop.orders.status"],["shop.orders.tags"],["shop.orders.created"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT * FROM orders ARRAY JOIN tags AS tag
< raw
| {"columns":["id","customer_id","amount","status","tags","created"],"lineage":[["shop.orders.id"],["shop.orders.customer_id"],["shop.orders.amount"],["shop.orders.status"],["shop.orders.tags"],["shop.orders.created"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT tag, n FROM orders ARRAY JOIN tags AS tag, arrayEnumerate(tags) AS n
< raw
| {"columns":["tag","n"],"lineage":[["shop.orders.tags"],["shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT arrayMap(x -> upper(x), tags) FROM orders
< raw
| {"columns":["arrayMap(lambda(tuple(x), upper(x)), tags)"],"lineage":[["shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT arrayFilter(x -> x != status, tags) FROM orders
< raw
| {"columns":["arrayFilter(lambda(tuple(x), notEquals(x, status)), tags)"],"lineage":[["shop.orders.status","shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT arrayMap(x -> x || status, tags) AS t FROM orders
< raw
| {"columns":["t"],"lineage":[["shop.orders.status","shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT length(tags), has(tags, 'gift') FROM orders
< raw
| {"columns":["length(tags)","has(tags, 'gift')"],"lineage":[["shop.orders.tags"],["shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT arrayJoin(tags) AS tag FROM orders
< raw
| {"columns":["tag"],"lineage":[["shop.orders.tags"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT arrayJoin([1, 2, 3]) AS n
< raw
| {"columns":["n"],"lineage":[null]}
> raw
| span shop
| SELECT groupArray(name) AS names FROM customers
< raw
| {"columns":["names"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT arraySum(x -> x * 2, [qty, order_id]) FROM items
< raw
| {"columns":["arraySum(lambda(tuple(x), multiply(x, 2)), [qty, order_id])"],"lineage":[["shop.items.order_id","shop.items.qty"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT id, row_number() OVER (ORDER BY id) FROM orders
< raw
| {"columns":["id","row_number() OVER (ORDER BY id ASC)"],"lineage":[["shop.orders.id"],["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, sum(amount) OVER (PARTITION BY customer_id) FROM orders
< raw
| {"columns":["id","sum(amount) OVER (PARTITION BY customer_id)"],"lineage":[["shop.orders.id"],["shop.orders.amount","shop.orders.customer_id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, sum(amount) OVER (PARTITION BY customer_id ORDER BY id) AS running FROM orders
< raw
| {"columns":["id","running"],"lineage":[["shop.orders.id"],["shop.orders.amount","shop.orders.customer_id","shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, rank() OVER (PARTITION BY status ORDER BY amount DESC) FROM orders
< raw
| {"columns":["id","rank() OVER (PARTITION BY status ORDER BY amount DESC)"],"lineage":[["shop.orders.id"],["shop.orders.amount","shop.orders.status"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, lagInFrame(amount) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM orders
< raw
| {"columns":["id","lagInFrame(amount) OVER (ORDER BY id ASC ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)"],"lineage":[["shop.orders.id"],["shop.orders.amount","shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id, count() OVER w FROM orders WINDOW w AS (PARTITION BY customer_id)
< raw
| {"columns":["id","count() OVER w"],"lineage":[["shop.orders.id"],["shop.orders.customer_id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT count(*) OVER (PARTITION BY email) FROM customers
< raw
| {"columns":["count() OVER (PARTITION BY email)"],"lineage":[["shop.customers.email"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT id FROM orders QUALIFY row_number() OVER (PARTITION BY customer_id ORDER BY amount DESC) = 1
< raw
| {"columns":["id"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT * EXCEPT (email) FROM customers
< raw
| {"columns":["id","name","country","created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT * EXCEPT email FROM customers
< raw
| {"columns":["id","name","country","created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT * EXCEPT (email, country) FROM customers
< raw
| {"columns":["id","name","created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT * REPLACE (upper(name) AS name) FROM customers
< raw
| {"columns":["id","name","email","country","created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT c.* EXCEPT (id) FROM customers AS c
< raw
| {"columns":["name","email","country","created"],"lineage":[["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT * APPLY (toString) FROM items
< raw
| {"columns":["toString(order_id)","toString(sku)","toString(qty)","toString(price)"],"lineage":[["shop.items.order_id"],["shop.items.sku"],["shop.items.qty"],["shop.items.price"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT * APPLY (x -> concat(toString(x), status)) FROM orders
< raw
| {"columns":["concat(toString(id), status)","concat(toString(customer_id), status)","concat(toString(amount), status)","concat(toString(status), status)","concat(toString(tags), status)","concat(toString(created), status)"],"lineage":[["shop.orders.id","shop.orders.status"],["shop.orders.customer_id","shop.orders.status"],["shop.orders.amount","shop.orders.status"],["shop.orders.status"],["shop.orders.status","shop.orders.tags"],["shop.orders.created","shop.orders.status"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT * EXCEPT (sku) APPLY (max) FROM items
< raw
| {"columns":["max(order_id)","max(qty)","max(price)"],"lineage":[["shop.items.order_id"],["shop.items.qty"],["shop.items.price"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT * EXCEPT ('mail') FROM customers
< raw
| {"columns":["id","name","country","created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT untuple((id, email)) FROM customers
< raw
| {"columns":["tupleElement((id, email), '1')","tupleElement((id, email), '2')"],"lineage":[["shop.customers.email","shop.customers.id"],["shop.customers.email","shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT COLUMNS('^c') FROM customers
< raw
| {"columns":["country","created"],"lineage":[["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT concat(COLUMNS('^(name|email)$')) FROM customers
< raw
| {"columns":["concat(name, email)"],"lineage":[["shop.customers.email","shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT amount * 2 AS doubled, doubled + 1 FROM orders
< raw
| {"columns":["doubled","plus(doubled, 1)"],"lineage":[["shop.orders.amount"],["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id + 1 AS id FROM orders
< raw
| {"columns":["id"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT name AS id, id FROM customers
< raw
| {"columns":["id","id"],"lineage":[["shop.customers.name"],["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT email AS e FROM customers WHERE e != ''
< raw
| {"columns":["e"],"lineage":[["shop.customers.email"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT upper(name) AS n FROM customers ORDER BY n
< raw
| {"columns":["n"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT country AS c, count() FROM customers GROUP BY c
< raw
| {"columns":["c","count()"],"lineage":[["shop.customers.country"],null],"tables":["shop.customers"]}
> raw
| span shop
| SELECT (amount + 1 AS a) * 2, a FROM orders
< raw
| {"columns":["multiply(a, 2)","a"],"lineage":[["shop.orders.amount"],["shop.orders.amount"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id AS x, x * 2 AS y, y + 1 FROM orders
< raw
| {"columns":["x","y","plus(y, 1)"],"lineage":[["shop.orders.id"],["shop.orders.id"],["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT customer_id AS amount, amount FROM orders
< raw
| {"columns":["amount","amount"],"lineage":[["shop.orders.customer_id"],["shop.orders.customer_id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT x FROM customers WHERE (email AS x) != ''
< raw
| {"columns":["x"],"lineage":[["shop.customers.email"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT toString(id) AS status, concat(status, '') AS s FROM orders SETTINGS prefer_column_name_to_alias = 1
< raw
| {"columns":["status","s"],"lineage":[["shop.orders.id"],["shop.orders.status"]],"tables":["shop.orders"]}
> raw
| span shop
| FROM customers SELECT name
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| FROM customers SELECT *
< raw
| {"columns":["id","name","email","country","created"],"lineage":[["shop.customers.id"],["shop.customers.name"],["shop.customers.email"],["shop.customers.country"],["shop.customers.created"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT name FROM customers PREWHERE id > 1
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT order_id, qty FROM items FINAL
< raw
| {"columns":["order_id","qty"],"lineage":[["shop.items.order_id"],["shop.items.qty"]],"tables":["shop.items"]}
> raw
| span shop
| SELECT id FROM orders LIMIT 1 BY customer_id
< raw
| {"columns":["id"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT id FROM orders ORDER BY id WITH FILL
< raw
| {"columns":["id"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT name FROM customers SETTINGS max_threads = 1
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT id FROM orders WHERE _part = 'all_0_0_0'
< raw
| {"columns":["id"],"lineage":[["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT _part, id FROM orders
< raw
| {"columns":["_part","id"],"lineage":[["shop.orders._part"],["shop.orders.id"]],"tables":["shop.orders"]}
> raw
| span shop
| SELECT name FROM customers WHERE id > 0 ORDER BY id LIMIT 3
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT id FROM customers LIMIT (SELECT count() FROM orders)
< raw
| {"columns":["id"],"lineage":[["shop.customers.id"]],"tables":["shop.customers","shop.orders"]}
> raw
| span shop
| SELECT id FROM customers LIMIT 1 OFFSET (SELECT count() FROM items)
< raw
| {"columns":["id"],"lineage":[["shop.customers.id"]],"tables":["shop.customers","shop.items"]}
> raw
| span shop
| SELECT DISTINCT ON (country) country, name FROM customers
< raw
| {"columns":["country","name"],"lineage":[["shop.customers.country"],["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT TOP 3 name FROM customers
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT number FROM numbers(10)
< raw
| {"columns":["number"],"lineage":[["~"]],"tables":["~"]}
> raw
| span shop
| SELECT number * 2 AS n FROM numbers(3)
< raw
| {"columns":["n"],"lineage":[["~"]],"tables":["~"]}
> raw
| span shop
| SELECT * FROM numbers(2)
< raw
| {"columns":["number"],"lineage":[["~"]],"tables":["~"]}
> raw
| span shop
| SELECT c.name, n.number FROM customers AS c, numbers(2) AS n
< raw
| {"columns":["name","number"],"lineage":[["shop.customers.name"],["~"]],"tables":["~","shop.customers"]}
> raw
| span shop
| SELECT email FROM merge('shop', '^customers$')
< raw
| {"columns":["email"],"lineage":[["~"]],"tables":["~"]}
> raw
| span shop
| SELECT * FROM merge('^(orders|customers)$')
< raw
| {"columns":["id","customer_id","amount","status","tags","name","email","country","created"],"lineage":[["~"],["~"],["~"],["~"],["~"],["~"],["~"],["~"],["~"]],"tables":["~"]}
> raw
| span shop
| SELECT name FROM system.tables
< raw
| {"columns":["name"],"lineage":[["system.tables.name"]],"tables":["system.tables"]}
> raw
| span shop
| SELECT database, name FROM system.tables WHERE database = 'shop'
< raw
| {"columns":["database","name"],"lineage":[["system.tables.database"],["system.tables.name"]],"tables":["system.tables"]}
> raw
| span shop
| SELECT name, type FROM system.columns WHERE table = 'customers'
< raw
| {"columns":["name","type"],"lineage":[["system.columns.name"],["system.columns.type"]],"tables":["system.columns"]}
> raw
| span shop
| SELECT count() FROM system.databases
< raw
| {"columns":["count()"],"lineage":[null],"tables":["system.databases"]}
> raw
| span shop
| SELECT name, salary FROM hr.employees
< raw
| {"columns":["name","salary"],"lineage":[["hr.employees.name"],["hr.employees.salary"]],"tables":["hr.employees"]}
> raw
| span shop
| SELECT * FROM hr.depts
< raw
| {"columns":["id","title"],"lineage":[["hr.depts.id"],["hr.depts.title"]],"tables":["hr.depts"]}
> raw
| span shop
| SELECT e.name, d.title FROM hr.employees AS e JOIN hr.depts AS d ON d.id = e.dept_id
< raw
| {"columns":["name","title"],"lineage":[["hr.employees.name"],["hr.depts.title"]],"tables":["hr.depts","hr.employees"]}
> raw
| span shop
| SELECT e.name, m.name FROM hr.employees AS e LEFT JOIN hr.employees AS m ON m.id = e.manager_id
< raw
| {"columns":["name","m.name"],"lineage":[["hr.employees.name"],["hr.employees.name"]],"tables":["hr.employees"]}
> raw
| span shop
| SELECT hr.employees.name FROM hr.employees
< raw
| {"columns":["name"],"lineage":[["hr.employees.name"]],"tables":["hr.employees"]}
> raw
| span shop
| SELECT nope FROM customers
< raw
| {"error":"UNKNOWN_IDENTIFIER"}
> raw
| span shop
| SELECT c.nope FROM customers AS c
< raw
| {"error":"UNKNOWN_IDENTIFIER"}
> raw
| span shop
| SELECT * FROM nope
< raw
| {"error":"UNKNOWN_TABLE"}
> raw
| span shop
| SELECT * FROM nodb.customers
< raw
| {"error":"UNKNOWN_DATABASE"}
> raw
| span shop
| SELECT id FROM customers WHERE nope = 1
< raw
| {"error":"UNKNOWN_IDENTIFIER"}
> raw
| span shop
| SELECT x.id FROM customers AS c
< raw
| {"error":"UNKNOWN_IDENTIFIER"}
> raw
| span shop
| SELECT * FROM customers WHERE
< raw
| {"error":"SYNTAX_ERROR"}
> raw
| span shop
| SELECT id, FROM customers
< raw
| {"columns":["id"],"lineage":[["shop.customers.id"]],"tables":["shop.customers"]}
> raw
| span shop
| SELECT name FROM customers WHERE id IN (SELECT nope FROM orders)
< raw
| {"error":"UNKNOWN_IDENTIFIER"}
> raw
| span shop
| SELECT name FROM hr.customers
< raw
| {"error":"UNKNOWN_TABLE"}
> raw
| span hr
| SELECT name FROM employees
< raw
| {"columns":["name"],"lineage":[["hr.employees.name"]],"tables":["hr.employees"]}
> raw
| span hr
| SELECT * FROM employees
< raw
| {"columns":["id","name","salary","dept_id","manager_id"],"lineage":[["hr.employees.id"],["hr.employees.name"],["hr.employees.salary"],["hr.employees.dept_id"],["hr.employees.manager_id"]],"tables":["hr.employees"]}
> raw
| span hr
| SELECT e.name, d.title FROM employees AS e JOIN depts AS d ON d.id = e.dept_id
< raw
| {"columns":["name","title"],"lineage":[["hr.employees.name"],["hr.depts.title"]],"tables":["hr.depts","hr.employees"]}
> raw
| span hr
| SELECT name, yearly FROM payroll
< raw
| {"columns":["name","yearly"],"lineage":[["hr.employees.name","hr.payroll.name"],["hr.employees.salary","hr.payroll.yearly"]],"tables":["hr.depts","hr.employees","hr.payroll"]}
> raw
| span hr
| SELECT name FROM shop.customers
< raw
| {"columns":["name"],"lineage":[["shop.customers.name"]],"tables":["shop.customers"]}
> raw
| span hr
| SELECT e.name, c.email FROM employees AS e JOIN shop.customers AS c ON c.name = e.name
< raw
| {"columns":["name","email"],"lineage":[["hr.employees.name"],["shop.customers.email"]],"tables":["hr.employees","shop.customers"]}
> raw
| span hr
| SELECT dept_id, avg(salary) FROM employees GROUP BY dept_id
< raw
| {"columns":["dept_id","avg(salary)"],"lineage":[["hr.employees.dept_id"],["hr.employees.salary"]],"tables":["hr.employees"]}
> raw
| span hr
| WITH d AS (SELECT id, title FROM depts) SELECT e.name, d.title FROM employees AS e JOIN d ON d.id = e.dept_id
< raw
| {"columns":["name","title"],"lineage":[["hr.employees.name"],["hr.depts.title"]],"tables":["hr.depts","hr.employees"]}
> raw
| span hr
| SELECT * FROM shop.vip
< raw
| {"columns":["customer_id","name"],"lineage":[["shop.customer_totals.customer_id","shop.customers.id","shop.vip.customer_id"],["shop.customer_totals.name","shop.customers.name","shop.vip.name"]],"tables":["shop.customer_totals","shop.customers","shop.orders","shop.vip"]}
> raw
| span hr
| SELECT name FROM customers
< raw
| {"error":"UNKNOWN_TABLE"}
> raw
| fixture
| CREATE DATABASE lab
| CREATE TABLE lab.events (id UInt32, attrs Nested(k String, v UInt32)) ENGINE = MergeTree ORDER BY id
< raw
| {"databases":[{"name":"lab","tables":[{"name":"events","columns":["id","attrs.k","attrs.v"]}],"views":null}]}
> raw
| span lab
| SELECT attrs.k FROM events
< raw
| {"columns":["attrs.k"],"lineage":[["lab.events.attrs.k"]],"tables":["lab.events"]}
> raw
| span lab
| SELECT events.attrs.v FROM events
< raw
| {"columns":["attrs.v"],"lineage":[["lab.events.attrs.v"]],"tables":["lab.events"]}
> raw
| span lab
| SELECT * FROM events
< raw
| {"columns":["id","attrs.k","attrs.v"],"lineage":[["lab.events.id"],["lab.events.attrs.k"],["lab.events.attrs.v"]],"tables":["lab.events"]}
> raw
| span lab
| SELECT id, k FROM events ARRAY JOIN attrs.k AS k
< raw
| {"columns":["id","k"],"lineage":[["lab.events.id"],["lab.events.attrs.k"]],"tables":["lab.events"]}
> raw
| span lab
| SELECT id, attrs.k FROM events ARRAY JOIN attrs
< raw
| {"columns":["id","attrs.k"],"lineage":[["lab.events.id"],["lab.events.attrs.k","lab.events.attrs.v"]],"tables":["lab.events"]}
