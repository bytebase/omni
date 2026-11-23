-- Two databases: shop, which the sessions start in, and archive, which shop
-- reads across. No table has a key, so the optimizer cannot drop a join from
-- the plan and the showplan names every table a query reads. Batches are
-- separated by GO lines.
CREATE DATABASE shop
GO
CREATE DATABASE archive
GO
USE archive
GO
CREATE TABLE dbo.old_orders (id int, customer_id int, amount decimal(10,2))
GO
CREATE VIEW dbo.recent AS SELECT o.id, o.amount FROM dbo.old_orders AS o WHERE o.id > 100
GO
USE shop
GO
CREATE SCHEMA hr
GO
CREATE TABLE dbo.customers (id int, name nvarchar(50), email nvarchar(100), region varchar(10))
CREATE TABLE dbo.orders (id int, customer_id int, amount decimal(10,2), status varchar(10), created datetime2)
CREATE TABLE dbo.order_items (order_id int, line int, product_id int, qty int, price decimal(10,2))
CREATE TABLE dbo.products (id int, name nvarchar(50), category varchar(20))
CREATE TABLE dbo.[Order Details] ([Order ID] int, [Unit Price] money, Qty int)
CREATE TABLE dbo.MixedCase (Id int, CamelName nvarchar(20))
CREATE TABLE hr.employees (id int, name nvarchar(50), salary decimal(10,2), manager_id int, dept_id int, ssn char(11))
CREATE TABLE hr.departments (id int, name nvarchar(50))
CREATE TABLE hr.customers (id int, nickname nvarchar(20))
CREATE TABLE dbo.tickets (ticket_no int IDENTITY(1, 1), subject nvarchar(50))
GO
CREATE VIEW dbo.order_totals AS SELECT o.id, o.customer_id, o.amount * 1.1 AS gross, c.name AS customer_name FROM dbo.orders AS o JOIN dbo.customers AS c ON c.id = o.customer_id
GO
CREATE VIEW dbo.big_orders AS SELECT t.id, t.customer_name, t.gross FROM dbo.order_totals AS t WHERE t.gross > 100
GO
CREATE VIEW hr.staff (staff_id, staff_name, dept) AS SELECT e.id, e.name, d.name FROM hr.employees AS e LEFT JOIN hr.departments AS d ON d.id = e.dept_id
GO
CREATE VIEW dbo.all_people AS SELECT c.name FROM dbo.customers AS c UNION ALL SELECT e.name FROM hr.employees AS e
GO
CREATE VIEW dbo.v_products AS SELECT * FROM dbo.products
GO
CREATE VIEW dbo.archived AS SELECT o.id, o.amount FROM archive.dbo.old_orders AS o
GO
CREATE VIEW dbo.[Vip Names] AS SELECT c.name AS VipName FROM dbo.customers AS c WHERE c.region = 'vip'
GO
