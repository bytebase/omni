-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestQuerySpanFidelity default#1
> raw
| fixture
| CREATE SCHEMA sales;
| CREATE SCHEMA hr;
| CREATE SCHEMA "Mixed";
| CREATE TABLE public.customers (id int PRIMARY KEY, name text, email text, region text);
| CREATE TABLE public.orders (id int PRIMARY KEY, customer_id int, amount numeric, status text, created_at timestamp);
| CREATE TABLE public.items (order_id int, sku text, qty int, price numeric);
| CREATE TABLE public.accounts (id bigint, customer_id int, iban text);
| CREATE TABLE sales.orders (id int, region text, total numeric);
| CREATE TABLE sales.targets (region text, goal numeric);
| CREATE TABLE hr.employees (id int, name text, ssn text, manager_id int, salary numeric, dept text);
| CREATE TABLE "Mixed"."CamelTable" ("Id" int, "MixedCol" text, lower_col text);
| CREATE VIEW public.customer_emails AS SELECT id, email FROM public.customers;
| CREATE VIEW public.order_totals AS SELECT customer_id, sum(amount) AS total, count(*) AS n FROM public.orders GROUP BY customer_id;
| CREATE VIEW public.vip AS SELECT e.id, e.email, t.total FROM public.customer_emails e JOIN public.order_totals t ON t.customer_id = e.id WHERE t.total > 1000;
| CREATE VIEW public.labeled AS SELECT id, upper(name) AS label, region FROM public.customers;
| CREATE VIEW hr.directory AS SELECT id, name, dept FROM hr.employees;
| CREATE VIEW public.staff_orders AS SELECT e.name, o.amount FROM hr.employees e JOIN public.orders o ON o.customer_id = e.id;
| CREATE VIEW "Mixed"."CamelView" AS SELECT "Id", "MixedCol" AS "Renamed" FROM "Mixed"."CamelTable";
| CREATE MATERIALIZED VIEW public.region_counts AS SELECT region, count(*) AS n FROM public.customers GROUP BY region;
| CREATE VIEW public.renamed (cid, mail) AS SELECT id, email FROM public.customers;
| CREATE VIEW "Mixed"."Upper View" AS SELECT "Id" AS "ID", lower_col AS "Lower" FROM "Mixed"."CamelTable";
< raw
| {"database":"db","schemas":[{"name":"Mixed","tables":[{"name":"CamelTable","columns":[{"name":"Id","type":"bigint"},{"name":"MixedCol","type":"text"},{"name":"lower_col","type":"text"},{"name":"rowid","type":"bigint","hidden":true}]}],"views":[{"name":"CamelView","columns":[{"name":"Id","type":"bigint","origin":"db.Mixed.CamelTable.Id"},{"name":"Renamed","type":"text","origin":"db.Mixed.CamelTable.MixedCol"}],"definition":"SELECT \"Id\", \"MixedCol\" AS \"Renamed\" FROM db.\"Mixed\".\"CamelTable\""},{"name":"Upper View","columns":[{"name":"ID","type":"bigint","origin":"db.Mixed.CamelTable.Id"},{"name":"Lower","type":"text","origin":"db.Mixed.CamelTable.lower_col"}],"definition":"SELECT \"Id\" AS \"ID\", lower_col AS \"Lower\" FROM db.\"Mixed\".\"CamelTable\""}]},{"name":"hr","tables":[{"name":"employees","columns":[{"name":"id","type":"bigint"},{"name":"name","type":"text"},{"name":"ssn","type":"text"},{"name":"manager_id","type":"bigint"},{"name":"salary","type":"numeric"},{"name":"dept","type":"text"},{"name":"rowid","type":"bigint","hidden":true}]}],"views":[{"name":"directory","columns":[{"name":"id","type":"bigint","origin":"db.hr.employees.id"},{"name":"name","type":"text","origin":"db.hr.employees.name"},{"name":"dept","type":"text","origin":"db.hr.employees.dept"}],"definition":"SELECT id, name, dept FROM db.hr.employees"}]},{"name":"public","tables":[{"name":"accounts","columns":[{"name":"id","type":"bigint"},{"name":"customer_id","type":"bigint"},{"name":"iban","type":"text"},{"name":"rowid","type":"bigint","hidden":true}]},{"name":"customers","columns":[{"name":"id","type":"bigint"},{"name":"name","type":"text"},{"name":"email","type":"text"},{"name":"region","type":"text"}]},{"name":"items","columns":[{"name":"order_id","type":"bigint"},{"name":"sku","type":"text"},{"name":"qty","type":"bigint"},{"name":"price","type":"numeric"},{"name":"rowid","type":"bigint","hidden":true}]},{"name":"orders","columns":[{"name":"id","type":"bigint"},{"name":"customer_id","type":"bigint"},{"name":"amount","type":"numeric"},{"name":"status","type":"text"},{"name":"created_at","type":"timestamp without time zone"}]}],"views":[{"name":"customer_emails","columns":[{"name":"id","type":"bigint","origin":"db.public.customers.id"},{"name":"email","type":"text","origin":"db.public.customers.email"}],"definition":"SELECT id, email FROM db.public.customers"},{"name":"labeled","columns":[{"name":"id","type":"bigint","origin":"db.public.customers.id"},{"name":"label","type":"text"},{"name":"region","type":"text","origin":"db.public.customers.region"}],"definition":"SELECT id, upper(name) AS label, region FROM db.public.customers"},{"name":"order_totals","columns":[{"name":"customer_id","type":"bigint","origin":"db.public.orders.customer_id"},{"name":"total","type":"numeric"},{"name":"n","type":"bigint"}],"definition":"SELECT customer_id, sum(amount) AS total, count(*) AS n FROM db.public.orders GROUP BY customer_id"},{"name":"renamed","columns":[{"name":"cid","type":"bigint","origin":"db.public.customers.id"},{"name":"mail","type":"text","origin":"db.public.customers.email"}],"definition":"SELECT id, email FROM db.public.customers"},{"name":"staff_orders","columns":[{"name":"name","type":"text","origin":"db.hr.employees.name"},{"name":"amount","type":"numeric","origin":"db.public.orders.amount"}],"definition":"SELECT e.name, o.amount FROM db.hr.employees AS e JOIN db.public.orders AS o ON o.customer_id = e.id"},{"name":"vip","columns":[{"name":"id","type":"bigint","origin":"db.public.customers.id"},{"name":"email","type":"text","origin":"db.public.customers.email"},{"name":"total","type":"numeric"}],"definition":"SELECT e.id, e.email, t.total FROM db.public.customer_emails AS e JOIN db.public.order_totals AS t ON t.customer_id = e.id WHERE t.total \u003e 1000"}],"materializedViews":[{"name":"region_counts","columns":[{"name":"region","type":"text"},{"name":"n","type":"bigint"},{"name":"rowid","type":"bigint","hidden":true}],"definition":"SELECT region, count(*) AS n FROM db.public.customers GROUP BY region"}]},{"name":"sales","tables":[{"name":"orders","columns":[{"name":"id","type":"bigint"},{"name":"region","type":"text"},{"name":"total","type":"numeric"},{"name":"rowid","type":"bigint","hidden":true}]},{"name":"targets","columns":[{"name":"region","type":"text"},{"name":"goal","type":"numeric"},{"name":"rowid","type":"bigint","hidden":true}]}]}]}
> raw
| span public
| SELECT id, name FROM customers
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT c.id AS customer, c.email FROM customers c
< raw
| {"columns":[{"name":"customer","source":"db.public.customers.id"},{"name":"email","source":"db.public.customers.email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT id + customer_id, upper(status) FROM orders
< raw
| {"columns":[{"name":"?column?"},{"name":"upper"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT amount * 2 AS doubled, status || '!' AS shout FROM orders
< raw
| {"columns":[{"name":"doubled"},{"name":"shout"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT CAST(amount AS text) AS a, amount::int FROM orders
< raw
| {"columns":[{"name":"a"},{"name":"amount"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT 1, 'a' AS b
< raw
| {"columns":[{"name":"?column?"},{"name":"b"}]}
> raw
| span public
| SELECT DISTINCT region FROM customers
< raw
| {"columns":[{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT name FROM customers ORDER BY region LIMIT 5
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT id FROM orders OFFSET 1 ROWS FETCH FIRST 2 ROWS ONLY
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT CASE WHEN amount > 10 THEN status ELSE 'small' END AS s FROM orders
< raw
| {"columns":[{"name":"s"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT coalesce(email, name) FROM customers
< raw
| {"columns":[{"name":"coalesce"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT id FROM customers WHERE email LIKE '%@example.com'
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT public.customers.name FROM public.customers
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| TABLE customers
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT (c).name, (c).* FROM customers c
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"},{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT (ROW(id, email)).* FROM customers
< raw
| {"columns":[{"name":"?column?","source":"db.public.customers.id"},{"name":"?column?","source":"db.public.customers.email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT ((id, email)).* FROM customers
< raw
| {"columns":[{"name":"?column?","source":"db.public.customers.id"},{"name":"?column?","source":"db.public.customers.email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT b FROM (SELECT (ROW(id, email)).* FROM customers) AS q(a, b)
< raw
| {"columns":[{"name":"b","source":"db.public.customers.email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT (((id, email) AS k, v)).* FROM customers
< raw
| {"columns":[{"name":"k","source":"db.public.customers.id"},{"name":"v","source":"db.public.customers.email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT ((SELECT id, email FROM customers LIMIT 1)).*
< raw
| {"error":"42601"}
> raw
| span public
| SELECT row_to_json(c) FROM customers c
< raw
| {"columns":[{"name":"row_to_json"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT * FROM ONLY customers
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT a, b FROM customers AS c(a, b)
< raw
| {"columns":[{"name":"a","source":"db.public.customers.id"},{"name":"b","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT * FROM customers
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT c.* FROM customers c
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT *, id AS again FROM orders
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"again","source":"db.public.orders.id"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM orders o, customers c WHERE o.customer_id = c.id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT s.* FROM (SELECT id, email FROM customers) s
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"email","source":"db.public.customers.email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT * FROM items
< raw
| {"columns":[{"name":"order_id","source":"db.public.items.order_id"},{"name":"sku","source":"db.public.items.sku"},{"name":"qty","source":"db.public.items.qty"},{"name":"price","source":"db.public.items.price"}],"tables":["db.public.items"]}
> raw
| span public
| SELECT rowid, sku FROM items
< raw
| {"columns":[{"name":"rowid","source":"db.public.items.rowid"},{"name":"sku","source":"db.public.items.sku"}],"tables":["db.public.items"]}
> raw
| span public
| SELECT i.*, a.* FROM items i JOIN accounts a ON a.id = i.order_id
< raw
| {"columns":[{"name":"order_id","source":"db.public.items.order_id"},{"name":"sku","source":"db.public.items.sku"},{"name":"qty","source":"db.public.items.qty"},{"name":"price","source":"db.public.items.price"},{"name":"id","source":"db.public.accounts.id"},{"name":"customer_id","source":"db.public.accounts.customer_id"},{"name":"iban","source":"db.public.accounts.iban"}],"tables":["db.public.accounts","db.public.items"]}
> raw
| span public
| SELECT * FROM items NATURAL JOIN accounts
< raw
| {"columns":[{"name":"order_id","source":"db.public.items.order_id"},{"name":"sku","source":"db.public.items.sku"},{"name":"qty","source":"db.public.items.qty"},{"name":"price","source":"db.public.items.price"},{"name":"id","source":"db.public.accounts.id"},{"name":"customer_id","source":"db.public.accounts.customer_id"},{"name":"iban","source":"db.public.accounts.iban"}],"tables":["db.public.accounts","db.public.items"]}
> raw
| span public
| SELECT * FROM items JOIN accounts USING (rowid)
< raw
| {"error":"42703"}
> raw
| span public
| SELECT s.* FROM (SELECT * FROM hr.employees) s
< raw
| {"columns":[{"name":"id","source":"db.hr.employees.id"},{"name":"name","source":"db.hr.employees.name"},{"name":"ssn","source":"db.hr.employees.ssn"},{"name":"manager_id","source":"db.hr.employees.manager_id"},{"name":"salary","source":"db.hr.employees.salary"},{"name":"dept","source":"db.hr.employees.dept"}],"tables":["db.hr.employees"]}
> raw
| span public
| TABLE items
< raw
| {"columns":[{"name":"order_id","source":"db.public.items.order_id"},{"name":"sku","source":"db.public.items.sku"},{"name":"qty","source":"db.public.items.qty"},{"name":"price","source":"db.public.items.price"}],"tables":["db.public.items"]}
> raw
| span public
| SELECT o.id, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT o.id, c.name FROM orders o RIGHT JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT o.id, c.name FROM orders o FULL JOIN customers c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT o.id, c.name FROM orders o CROSS JOIN customers c
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM orders JOIN items ON items.order_id = orders.id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"order_id","source":"db.public.items.order_id"},{"name":"sku","source":"db.public.items.sku"},{"name":"qty","source":"db.public.items.qty"},{"name":"price","source":"db.public.items.price"}],"tables":["db.public.items","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers JOIN accounts ON accounts.customer_id = customers.id
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"id","source":"db.public.accounts.id"},{"name":"customer_id","source":"db.public.accounts.customer_id"},{"name":"iban","source":"db.public.accounts.iban"}],"tables":["db.public.accounts","db.public.customers"]}
> raw
| span public
| SELECT a.id, b.id FROM customers a JOIN customers b ON a.region = b.region
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"id","source":"db.public.customers.id"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT * FROM (orders o JOIN customers c ON c.id = o.customer_id) j
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT j.name FROM (orders o JOIN customers c ON c.id = o.customer_id) AS j
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers JOIN orders USING (id) JOIN items ON items.order_id = id
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"order_id","source":"db.public.items.order_id"},{"name":"sku","source":"db.public.items.sku"},{"name":"qty","source":"db.public.items.qty"},{"name":"price","source":"db.public.items.price"}],"tables":["db.public.customers","db.public.items","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers JOIN orders USING (id) JOIN accounts USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"customer_id","source":"db.public.accounts.customer_id"},{"name":"iban","source":"db.public.accounts.iban"}],"tables":["db.public.accounts","db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers FULL JOIN orders USING (id) FULL JOIN accounts USING (id)
< raw
| {"columns":[{"name":"id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"customer_id","source":"db.public.accounts.customer_id"},{"name":"iban","source":"db.public.accounts.iban"}],"tables":["db.public.accounts","db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers JOIN customer_emails USING (id, email)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"email","source":"db.public.customers.email"},{"name":"name","source":"db.public.customers.name"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customer_emails","db.public.customers"]}
> raw
| span public
| SELECT * FROM customers NATURAL FULL JOIN customer_emails
< raw
| {"columns":[{"name":"id"},{"name":"email"},{"name":"name","source":"db.public.customers.name"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customer_emails","db.public.customers"]}
> raw
| span public
| SELECT * FROM (SELECT * FROM customers) s JOIN orders USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM orders JOIN items ON true JOIN customers USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"order_id","source":"db.public.items.order_id"},{"name":"sku","source":"db.public.items.sku"},{"name":"qty","source":"db.public.items.qty"},{"name":"price","source":"db.public.items.price"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers","db.public.items","db.public.orders"]}
> raw
| span public
| SELECT id, name, amount FROM customers JOIN orders USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"amount","source":"db.public.orders.amount"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers JOIN orders USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers LEFT JOIN orders USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers RIGHT JOIN orders USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers FULL JOIN orders USING (id)
< raw
| {"columns":[{"name":"id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers JOIN accounts USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"customer_id","source":"db.public.accounts.customer_id"},{"name":"iban","source":"db.public.accounts.iban"}],"tables":["db.public.accounts","db.public.customers"]}
> raw
| span public
| SELECT id FROM accounts JOIN customers USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.accounts.id"}],"tables":["db.public.accounts","db.public.customers"]}
> raw
| span public
| SELECT customers.id, orders.id FROM customers JOIN orders USING (id)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"id","source":"db.public.orders.id"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM customers NATURAL JOIN customer_emails
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"email","source":"db.public.customers.email"},{"name":"name","source":"db.public.customers.name"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customer_emails","db.public.customers"]}
> raw
| span public
| SELECT * FROM orders NATURAL JOIN items
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"order_id","source":"db.public.items.order_id"},{"name":"sku","source":"db.public.items.sku"},{"name":"qty","source":"db.public.items.qty"},{"name":"price","source":"db.public.items.price"}],"tables":["db.public.items","db.public.orders"]}
> raw
| span public
| SELECT EXISTS (SELECT 1 FROM orders) AS any
< raw
| {"columns":[{"name":"any"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT id FROM customers WHERE id = ANY (SELECT customer_id FROM orders)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT orders.name FROM (SELECT * FROM customers) AS orders
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT c.id, t.x FROM customers c, LATERAL unnest(ARRAY[c.name, c.email]) AS t(x)
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"x"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT * FROM generate_series(1, 3) AS g
< raw
| {"columns":[{"name":"g"}]}
> raw
| span public
| SELECT g.n, c.name FROM generate_series(1, 3) AS g(n) JOIN customers c ON c.id = g.n
< raw
| {"columns":[{"name":"n"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT * FROM ROWS FROM (generate_series(1, 2), unnest(ARRAY['a', 'b'])) AS r(n, s)
< raw
| {"columns":[{"name":"n"},{"name":"s"}]}
> raw
| span public
| SELECT s.id FROM (SELECT id FROM customers UNION SELECT id FROM orders ORDER BY 1 LIMIT 2) s
< raw
| {"columns":[{"name":"id"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT (SELECT max(amount) FROM orders) AS top
< raw
| {"columns":[{"name":"top"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT name, (SELECT count(*) FROM orders o WHERE o.customer_id = c.id) AS n FROM customers c
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"},{"name":"n"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders)
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT name FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT id FROM (SELECT id, name FROM customers) AS s
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT x FROM (SELECT id AS x FROM customers) s
< raw
| {"columns":[{"name":"x","source":"db.public.customers.id"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT a, b FROM (SELECT id, email FROM customers) AS s(a, b)
< raw
| {"columns":[{"name":"a","source":"db.public.customers.id"},{"name":"b","source":"db.public.customers.email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT t.n FROM (SELECT amount + 1 AS n FROM orders) t
< raw
| {"columns":[{"name":"n"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT * FROM (SELECT * FROM (SELECT id, region FROM customers) a) b
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT c.name, l.amount FROM customers c, LATERAL (SELECT amount FROM orders o WHERE o.customer_id = c.id) l
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"},{"name":"amount","source":"db.public.orders.amount"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT c.name, l.* FROM customers c CROSS JOIN LATERAL (SELECT max(amount) AS m, min(status) FROM orders o WHERE o.customer_id = c.id) l
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"},{"name":"m"},{"name":"min"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT c.name, l.amount FROM customers c LEFT JOIN LATERAL (SELECT amount FROM orders o WHERE o.customer_id = c.id LIMIT 1) l ON true
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"},{"name":"amount","source":"db.public.orders.amount"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| WITH x AS (SELECT id, email FROM customers) SELECT * FROM x a JOIN x b USING (id)
< raw
| {"columns":[{"name":"id"},{"name":"email"},{"name":"email"}],"tables":["db.public.customers"]}
> raw
| span public
| WITH "Cte" AS (SELECT id, name FROM customers) SELECT "Cte".name FROM "Cte"
< raw
| {"columns":[{"name":"name"}],"tables":["db.public.customers"]}
> raw
| span public
| WITH c AS (SELECT id, email FROM customers) SELECT * FROM (WITH c AS (SELECT id, status FROM orders) SELECT * FROM c) inner_q
< raw
| {"columns":[{"name":"id"},{"name":"status"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| WITH orders AS (SELECT * FROM customers) SELECT o.name, (SELECT count(*) FROM orders) FROM orders o
< raw
| {"columns":[{"name":"name"},{"name":"count"}],"tables":["db.public.customers"]}
> raw
| span public
| WITH a AS (SELECT id, status FROM orders), b AS (SELECT * FROM a UNION ALL SELECT id, name FROM customers) SELECT status FROM b
< raw
| {"columns":[{"name":"status"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| WITH c AS (SELECT id, email FROM customers) SELECT * FROM c
< raw
| {"columns":[{"name":"id"},{"name":"email"}],"tables":["db.public.customers"]}
> raw
| span public
| WITH c AS (SELECT id, email FROM customers) SELECT email FROM c WHERE id > 1
< raw
| {"columns":[{"name":"email"}],"tables":["db.public.customers"]}
> raw
| span public
| WITH c(k, v) AS (SELECT id, email FROM customers) SELECT k, v FROM c
< raw
| {"columns":[{"name":"k"},{"name":"v"}],"tables":["db.public.customers"]}
> raw
| span public
| WITH a AS (SELECT id, name FROM customers), b AS (SELECT a.name FROM a) SELECT * FROM b
< raw
| {"columns":[{"name":"name"}],"tables":["db.public.customers"]}
> raw
| span public
| WITH t AS (SELECT customer_id, sum(amount) AS s FROM orders GROUP BY customer_id) SELECT c.name, t.s FROM customers c JOIN t ON t.customer_id = c.id
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"},{"name":"s"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| WITH customers AS (SELECT id, status AS name FROM orders) SELECT name FROM customers
< raw
| {"columns":[{"name":"name"}],"tables":["db.public.orders"]}
> raw
| span public
| WITH c AS MATERIALIZED (SELECT id, email FROM customers) SELECT email FROM c
< raw
| {"columns":[{"name":"email"}],"tables":["db.public.customers"]}
> raw
| span public
| WITH unused AS (SELECT ssn FROM hr.employees) SELECT id FROM customers
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"}],"tables":["db.hr.employees","db.public.customers"]}
> raw
| span public
| WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT n FROM r
< raw
| {"columns":[{"name":"n"}]}
> raw
| span public
| WITH RECURSIVE chain AS (SELECT id, manager_id, name FROM hr.employees WHERE manager_id IS NULL UNION ALL SELECT e.id, e.manager_id, e.name FROM hr.employees e JOIN chain ON e.manager_id = chain.id) SELECT * FROM chain
< raw
| {"columns":[{"name":"id"},{"name":"manager_id"},{"name":"name"}],"tables":["db.hr.employees"]}
> raw
| span public
| WITH RECURSIVE chain AS (SELECT id, name FROM hr.employees UNION SELECT e.id, e.name FROM hr.employees e JOIN chain c ON c.id = e.manager_id) SELECT name FROM chain
< raw
| {"columns":[{"name":"name"}],"tables":["db.hr.employees"]}
> raw
| span public
| WITH RECURSIVE c AS (SELECT id FROM customers UNION ALL SELECT id FROM customers LIMIT (SELECT count(*) FROM orders)) SELECT * FROM c
< raw
| {"columns":[{"name":"id"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT id, name FROM customers UNION SELECT id, status FROM orders
< raw
| {"columns":[{"name":"id"},{"name":"name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT id, name FROM customers UNION ALL SELECT id, status FROM orders
< raw
| {"columns":[{"name":"id"},{"name":"name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT email FROM customers INTERSECT SELECT status FROM orders
< raw
| {"columns":[{"name":"email","source":"db.public.customers.email"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT region FROM customers EXCEPT SELECT status FROM orders
< raw
| {"columns":[{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT id FROM customers UNION SELECT id FROM orders UNION SELECT order_id FROM items
< raw
| {"columns":[{"name":"id"}],"tables":["db.public.customers","db.public.items","db.public.orders"]}
> raw
| span public
| SELECT * FROM (SELECT id, email FROM customers UNION SELECT id, status FROM orders) u
< raw
| {"columns":[{"name":"id"},{"name":"email"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT name FROM customers UNION SELECT 'none' ORDER BY 1
< raw
| {"columns":[{"name":"name"}],"tables":["db.public.customers"]}
> raw
| span public
| VALUES (1, 'a'), (2, 'b')
< raw
| {"columns":[{"name":"column1"},{"name":"column2"}]}
> raw
| span public
| SELECT * FROM (VALUES (1, 'a')) AS v(n, s)
< raw
| {"columns":[{"name":"n"},{"name":"s"}]}
> raw
| span public
| SELECT * FROM renamed
< raw
| {"columns":[{"name":"cid","source":"db.public.customers.id"},{"name":"mail","source":"db.public.customers.email"}],"tables":["db.public.renamed"]}
> raw
| span public
| SELECT mail FROM renamed WHERE cid > 0
< raw
| {"columns":[{"name":"mail","source":"db.public.customers.email"}],"tables":["db.public.renamed"]}
> raw
| span public
| SELECT * FROM "Mixed"."Upper View"
< raw
| {"columns":[{"name":"ID","source":"db.Mixed.CamelTable.Id"},{"name":"Lower","source":"db.Mixed.CamelTable.lower_col"}],"tables":["db.Mixed.Upper View"]}
> raw
| span public
| SELECT "ID", "Lower" FROM "Mixed"."Upper View" u
< raw
| {"columns":[{"name":"ID","source":"db.Mixed.CamelTable.Id"},{"name":"Lower","source":"db.Mixed.CamelTable.lower_col"}],"tables":["db.Mixed.Upper View"]}
> raw
| span public
| SELECT * FROM customer_emails
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"email","source":"db.public.customers.email"}],"tables":["db.public.customer_emails"]}
> raw
| span public
| SELECT email FROM customer_emails WHERE id = 1
< raw
| {"columns":[{"name":"email","source":"db.public.customers.email"}],"tables":["db.public.customer_emails"]}
> raw
| span public
| SELECT * FROM order_totals
< raw
| {"columns":[{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"total"},{"name":"n"}],"tables":["db.public.order_totals"]}
> raw
| span public
| SELECT customer_id, total FROM order_totals
< raw
| {"columns":[{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"total"}],"tables":["db.public.order_totals"]}
> raw
| span public
| SELECT * FROM vip
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"email","source":"db.public.customers.email"},{"name":"total"}],"tables":["db.public.vip"]}
> raw
| span public
| SELECT v.email, o.status FROM vip v JOIN orders o ON o.customer_id = v.id
< raw
| {"columns":[{"name":"email","source":"db.public.customers.email"},{"name":"status","source":"db.public.orders.status"}],"tables":["db.public.orders","db.public.vip"]}
> raw
| span public
| SELECT label FROM labeled
< raw
| {"columns":[{"name":"label"}],"tables":["db.public.labeled"]}
> raw
| span public
| SELECT * FROM staff_orders
< raw
| {"columns":[{"name":"name","source":"db.hr.employees.name"},{"name":"amount","source":"db.public.orders.amount"}],"tables":["db.public.staff_orders"]}
> raw
| span public
| SELECT * FROM hr.directory
< raw
| {"columns":[{"name":"id","source":"db.hr.employees.id"},{"name":"name","source":"db.hr.employees.name"},{"name":"dept","source":"db.hr.employees.dept"}],"tables":["db.hr.directory"]}
> raw
| span public
| SELECT d.name, e.ssn FROM hr.directory d JOIN hr.employees e USING (id)
< raw
| {"columns":[{"name":"name","source":"db.hr.employees.name"},{"name":"ssn","source":"db.hr.employees.ssn"}],"tables":["db.hr.directory","db.hr.employees"]}
> raw
| span public
| SELECT * FROM region_counts
< raw
| {"columns":[{"name":"region","source":"db.public.region_counts.region"},{"name":"n","source":"db.public.region_counts.n"}],"tables":["db.public.region_counts"]}
> raw
| span public
| SELECT r.region, r.n, c.name FROM region_counts r JOIN customers c USING (region)
< raw
| {"columns":[{"name":"region","source":"db.public.region_counts.region"},{"name":"n","source":"db.public.region_counts.n"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.region_counts"]}
> raw
| span public
| SELECT rowid, region FROM public.region_counts
< raw
| {"columns":[{"name":"rowid","source":"db.public.region_counts.rowid"},{"name":"region","source":"db.public.region_counts.region"}],"tables":["db.public.region_counts"]}
> raw
| span public
| SELECT "Id", "MixedCol", lower_col FROM "Mixed"."CamelTable"
< raw
| {"columns":[{"name":"Id","source":"db.Mixed.CamelTable.Id"},{"name":"MixedCol","source":"db.Mixed.CamelTable.MixedCol"},{"name":"lower_col","source":"db.Mixed.CamelTable.lower_col"}],"tables":["db.Mixed.CamelTable"]}
> raw
| span public
| SELECT * FROM "Mixed"."CamelTable"
< raw
| {"columns":[{"name":"Id","source":"db.Mixed.CamelTable.Id"},{"name":"MixedCol","source":"db.Mixed.CamelTable.MixedCol"},{"name":"lower_col","source":"db.Mixed.CamelTable.lower_col"}],"tables":["db.Mixed.CamelTable"]}
> raw
| span public
| SELECT t."MixedCol" AS "Out" FROM "Mixed"."CamelTable" t
< raw
| {"columns":[{"name":"Out","source":"db.Mixed.CamelTable.MixedCol"}],"tables":["db.Mixed.CamelTable"]}
> raw
| span public
| SELECT * FROM "Mixed"."CamelView"
< raw
| {"columns":[{"name":"Id","source":"db.Mixed.CamelTable.Id"},{"name":"Renamed","source":"db.Mixed.CamelTable.MixedCol"}],"tables":["db.Mixed.CamelView"]}
> raw
| span public
| SELECT "Renamed" FROM "Mixed"."CamelView"
< raw
| {"columns":[{"name":"Renamed","source":"db.Mixed.CamelTable.MixedCol"}],"tables":["db.Mixed.CamelView"]}
> raw
| span public
| SELECT ID, NAME FROM CUSTOMERS
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT "id", "name" FROM "customers"
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT id AS "Weird Name", name AS "x.y" FROM customers
< raw
| {"columns":[{"name":"Weird Name","source":"db.public.customers.id"},{"name":"x.y","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT orders.id FROM customers AS orders
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT customers.amount FROM orders AS customers
< raw
| {"columns":[{"name":"amount","source":"db.public.orders.amount"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT o.id FROM orders o JOIN customers orders ON orders.id = o.customer_id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT customers.name FROM customers JOIN (SELECT id, status AS name FROM orders) AS customers2 ON customers2.id = customers.id
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT id, row_number() OVER (ORDER BY amount) FROM orders
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"row_number"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT id, sum(amount) OVER (PARTITION BY customer_id) AS running FROM orders
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"running"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT name, rank() OVER w FROM customers WINDOW w AS (PARTITION BY region ORDER BY name)
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"},{"name":"rank"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT id, lag(status) OVER (ORDER BY created_at) AS prev FROM orders
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"prev"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT DISTINCT ON (region) region, name FROM customers ORDER BY region, name
< raw
| {"columns":[{"name":"region","source":"db.public.customers.region"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT max(id) AS id FROM customers
< raw
| {"columns":[{"name":"id"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT count(*) FROM orders
< raw
| {"columns":[{"name":"count"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT customer_id, sum(amount), max(status) FROM orders GROUP BY customer_id
< raw
| {"columns":[{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"sum"},{"name":"max"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT region, count(DISTINCT email) FROM customers GROUP BY region HAVING count(*) > 1
< raw
| {"columns":[{"name":"region","source":"db.public.customers.region"},{"name":"count"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT string_agg(name, ',' ORDER BY name) FROM customers
< raw
| {"columns":[{"name":"string_agg"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT customer_id FROM orders GROUP BY customer_id
< raw
| {"columns":[{"name":"customer_id","source":"db.public.orders.customer_id"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT count(*) FILTER (WHERE status = 'paid') AS paid FROM orders
< raw
| {"columns":[{"name":"paid"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT *, c.name FROM customers c CROSS JOIN unnest(ARRAY['a']) AS tags
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"tags"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT c.name, upper(c.email) AS e, * FROM customers c CROSS JOIN unnest(ARRAY['a']) AS tags
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"},{"name":"e"},{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"tags"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT id, *, region FROM customers, unnest(ARRAY['a']) AS tags
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"tags"},{"name":"region","source":"db.public.customers.region"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT *, s.email FROM (SELECT id, upper(email) AS email FROM customers) s CROSS JOIN unnest(ARRAY['a']) AS tags
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"email"},{"name":"tags"},{"name":"email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT *, customer_id FROM orders JOIN items ON items.order_id = orders.id CROSS JOIN unnest(ARRAY['a']) AS tags
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"},{"name":"order_id","source":"db.public.items.order_id"},{"name":"sku","source":"db.public.items.sku"},{"name":"qty","source":"db.public.items.qty"},{"name":"price","source":"db.public.items.price"},{"name":"tags"},{"name":"customer_id","source":"db.public.orders.customer_id"}],"tables":["db.public.items","db.public.orders"]}
> raw
| span public
| SELECT id, name FROM customers AS OF SYSTEM TIME '-1ms'
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"}],"tablesUnstated":true}
> raw
| span public
| SELECT c.name, o.amount FROM customers AS c JOIN orders AS o ON o.customer_id = c.id AS OF SYSTEM TIME '-1ms'
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"},{"name":"amount","source":"db.public.orders.amount"}],"tablesUnstated":true}
> raw
| span public
| SELECT name FROM customers@customers_pkey
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT c.email FROM customers@{FORCE_INDEX=customers_pkey} AS c WHERE c.id > 1
< raw
| {"columns":[{"name":"email","source":"db.public.customers.email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT * FROM orders@{NO_FULL_SCAN} WHERE id = 1
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"customer_id","source":"db.public.orders.customer_id"},{"name":"amount","source":"db.public.orders.amount"},{"name":"status","source":"db.public.orders.status"},{"name":"created_at","source":"db.public.orders.created_at"}],"tables":["db.public.orders"]}
> raw
| span public
| SELECT o.id, c.name FROM orders AS o INNER HASH JOIN customers AS c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT o.id, c.name FROM orders AS o LEFT LOOKUP JOIN customers AS c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT name FROM customers ORDER BY PRIMARY KEY customers
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT * FROM customers WITH ORDINALITY
< raw
| {"columns":[{"name":"id","source":"db.public.customers.id"},{"name":"name","source":"db.public.customers.name"},{"name":"email","source":"db.public.customers.email"},{"name":"region","source":"db.public.customers.region"},{"name":"ordinality"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT IF(amount > 10, status, 'small') AS s, IFNULL(email, status) FROM orders JOIN customers ON customers.id = orders.customer_id
< raw
| {"columns":[{"name":"s"},{"name":"coalesce"}],"tables":["db.public.customers","db.public.orders"]}
> raw
| span public
| SELECT * FROM [SELECT id, email FROM customers] AS s
< raw
| {"columns":[{"name":"id"},{"name":"email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT s.email FROM [SELECT id, email FROM customers] AS s WHERE s.id > 1
< raw
| {"columns":[{"name":"email"}],"tables":["db.public.customers"]}
> raw
| span public
| SELECT * FROM [SHOW TABLES]
< raw
| {"columns":[{"name":"schema_name"},{"name":"table_name"},{"name":"type"},{"name":"owner"},{"name":"estimated_row_count"},{"name":"locality"}]}
> raw
| span public
| SELECT column_name FROM [SHOW COLUMNS FROM customers]
< raw
| {"columns":[{"name":"column_name"}]}
> raw
| span public
| SHOW COLUMNS FROM customers
< raw
| {"columns":[{"name":"column_name","source":"db.information_schema.columns.column_name"},{"name":"data_type","source":"db.information_schema.columns.crdb_sql_type"},{"name":"is_nullable"},{"name":"column_default","source":"db.information_schema.columns.column_default"},{"name":"generation_expression","source":"db.information_schema.columns.generation_expression"},{"name":"indices"},{"name":"is_hidden"}],"tablesUnstated":true}
> raw
| span public
| UPSERT INTO customers (id, name) VALUES (1, 'a')
< raw
| {"tablesUnstated":true}
> raw
| span public
| UPSERT INTO customers (id, name) SELECT id, status FROM orders
< raw
| {"tablesUnstated":true}
> raw
| span public
| EXPLAIN (OPT) SELECT name FROM customers
< raw
| {"columns":[{"name":"info"}],"tablesUnstated":true}
> raw
| span public
| EXPLAIN (VEC) SELECT o.id FROM orders AS o JOIN customers AS c ON c.id = o.customer_id
< raw
| {"columns":[{"name":"info"}],"tablesUnstated":true}
> raw
| span public
| EXPLAIN ANALYZE SELECT name FROM customers WHERE region = 'x'
< raw
| {"columns":[{"name":"info"}],"tablesUnstated":true}
> raw
| span sales,public
| SELECT id, region FROM orders
< raw
| {"columns":[{"name":"id","source":"db.sales.orders.id"},{"name":"region","source":"db.sales.orders.region"}],"tables":["db.sales.orders"]}
> raw
| span sales,public
| SELECT * FROM orders
< raw
| {"columns":[{"name":"id","source":"db.sales.orders.id"},{"name":"region","source":"db.sales.orders.region"},{"name":"total","source":"db.sales.orders.total"}],"tables":["db.sales.orders"]}
> raw
| span sales,public
| SELECT o.total, t.goal FROM orders o JOIN targets t USING (region)
< raw
| {"columns":[{"name":"total","source":"db.sales.orders.total"},{"name":"goal","source":"db.sales.targets.goal"}],"tables":["db.sales.orders","db.sales.targets"]}
> raw
| span sales,public
| SELECT id, amount FROM public.orders
< raw
| {"columns":[{"name":"id","source":"db.public.orders.id"},{"name":"amount","source":"db.public.orders.amount"}],"tables":["db.public.orders"]}
> raw
| span sales,public
| SELECT name FROM customers
< raw
| {"columns":[{"name":"name","source":"db.public.customers.name"}],"tables":["db.public.customers"]}
> raw
| span hr,public
| SELECT name, ssn FROM employees
< raw
| {"columns":[{"name":"name","source":"db.hr.employees.name"},{"name":"ssn","source":"db.hr.employees.ssn"}],"tables":["db.hr.employees"]}
> raw
| span hr,public
| SELECT * FROM directory
< raw
| {"columns":[{"name":"id","source":"db.hr.employees.id"},{"name":"name","source":"db.hr.employees.name"},{"name":"dept","source":"db.hr.employees.dept"}],"tables":["db.hr.directory"]}
> raw
| span hr,public
| SELECT e.name, m.name AS manager FROM employees e LEFT JOIN employees m ON m.id = e.manager_id
< raw
| {"columns":[{"name":"name","source":"db.hr.employees.name"},{"name":"manager","source":"db.hr.employees.name"}],"tables":["db.hr.employees"]}
> raw
| span hr,public
| SELECT d.name, c.email FROM directory d JOIN customers c ON c.id = d.id
< raw
| {"columns":[{"name":"name","source":"db.hr.employees.name"},{"name":"email","source":"db.public.customers.email"}],"tables":["db.hr.directory","db.public.customers"]}
