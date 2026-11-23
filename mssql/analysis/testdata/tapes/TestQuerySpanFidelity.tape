-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestQuerySpanFidelity default#1
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> query
| SELECT name FROM sys.databases WHERE database_id > 4 AND name <> N'omni_TestQuerySpanFidelity' ORDER BY name
< result
< column "name" "NVARCHAR"
> exec
| -- Two databases: shop, which the sessions start in, and archive, which shop
| -- reads across. No table has a key, so the optimizer cannot drop a join from
| -- the plan and the showplan names every table a query reads. Batches are
| -- separated by GO lines.
| CREATE DATABASE shop
|
< ok 0 -
> exec
|
| CREATE DATABASE archive
|
< ok 0 -
> exec
|
| USE archive
|
< ok 0 -
> exec
|
| CREATE TABLE dbo.old_orders (id int, customer_id int, amount decimal(10,2))
|
< ok 0 -
> exec
|
| CREATE VIEW dbo.recent AS SELECT o.id, o.amount FROM dbo.old_orders AS o WHERE o.id > 100
|
< ok 0 -
> exec
|
| USE shop
|
< ok 0 -
> exec
|
| CREATE SCHEMA hr
|
< ok 0 -
> exec
|
| CREATE TABLE dbo.customers (id int, name nvarchar(50), email nvarchar(100), region varchar(10))
| CREATE TABLE dbo.orders (id int, customer_id int, amount decimal(10,2), status varchar(10), created datetime2)
| CREATE TABLE dbo.order_items (order_id int, line int, product_id int, qty int, price decimal(10,2))
| CREATE TABLE dbo.products (id int, name nvarchar(50), category varchar(20))
| CREATE TABLE dbo.[Order Details] ([Order ID] int, [Unit Price] money, Qty int)
| CREATE TABLE dbo.MixedCase (Id int, CamelName nvarchar(20))
| CREATE TABLE hr.employees (id int, name nvarchar(50), salary decimal(10,2), manager_id int, dept_id int, ssn char(11))
| CREATE TABLE hr.departments (id int, name nvarchar(50))
| CREATE TABLE hr.customers (id int, nickname nvarchar(20))
|
< ok 0 -
> exec
|
| CREATE VIEW dbo.order_totals AS SELECT o.id, o.customer_id, o.amount * 1.1 AS gross, c.name AS customer_name FROM dbo.orders AS o JOIN dbo.customers AS c ON c.id = o.customer_id
|
< ok 0 -
> exec
|
| CREATE VIEW dbo.big_orders AS SELECT t.id, t.customer_name, t.gross FROM dbo.order_totals AS t WHERE t.gross > 100
|
< ok 0 -
> exec
|
| CREATE VIEW hr.staff (staff_id, staff_name, dept) AS SELECT e.id, e.name, d.name FROM hr.employees AS e LEFT JOIN hr.departments AS d ON d.id = e.dept_id
|
< ok 0 -
> exec
|
| CREATE VIEW dbo.all_people AS SELECT c.name FROM dbo.customers AS c UNION ALL SELECT e.name FROM hr.employees AS e
|
< ok 0 -
> exec
|
| CREATE VIEW dbo.v_products AS SELECT * FROM dbo.products
|
< ok 0 -
> exec
|
| CREATE VIEW dbo.archived AS SELECT o.id, o.amount FROM archive.dbo.old_orders AS o
|
< ok 0 -
> exec
|
| CREATE VIEW dbo.[Vip Names] AS SELECT c.name AS VipName FROM dbo.customers AS c WHERE c.region = 'vip'
|
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> query
| SELECT name FROM sys.databases WHERE database_id > 4 AND name <> N'omni_TestQuerySpanFidelity' ORDER BY name
< result
< column "name" "NVARCHAR"
< row
: string
| archive
< row
: string
| shop
> query
| SELECT s.name, o.name, o.type, c.name, m.definition
| FROM [archive].sys.schemas AS s
| LEFT JOIN [archive].sys.objects AS o ON o.schema_id = s.schema_id AND o.type IN ('U', 'V') AND o.is_ms_shipped = 0
| LEFT JOIN [archive].sys.columns AS c ON c.object_id = o.object_id
| LEFT JOIN [archive].sys.sql_modules AS m ON m.object_id = o.object_id
| WHERE s.schema_id < 16384 AND s.name NOT IN ('guest', 'INFORMATION_SCHEMA', 'sys')
| ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
< column "definition" "NVARCHAR"
< row
: string
| dbo
: string
| old_orders
: string
| U 
: string
| id
: null
< row
: string
| dbo
: string
| old_orders
: string
| U 
: string
| customer_id
: null
< row
: string
| dbo
: string
| old_orders
: string
| U 
: string
| amount
: null
< row
: string
| dbo
: string
| recent
: string
| V 
: string
| id
: string
|
| CREATE VIEW dbo.recent AS SELECT o.id, o.amount FROM dbo.old_orders AS o WHERE o.id > 100
|
< row
: string
| dbo
: string
| recent
: string
| V 
: string
| amount
: string
|
| CREATE VIEW dbo.recent AS SELECT o.id, o.amount FROM dbo.old_orders AS o WHERE o.id > 100
|
> exec
| USE [archive]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.recent', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| archive
: string
| dbo
: string
| old_orders
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> query
| SELECT s.name, o.name, o.type, c.name, m.definition
| FROM [shop].sys.schemas AS s
| LEFT JOIN [shop].sys.objects AS o ON o.schema_id = s.schema_id AND o.type IN ('U', 'V') AND o.is_ms_shipped = 0
| LEFT JOIN [shop].sys.columns AS c ON c.object_id = o.object_id
| LEFT JOIN [shop].sys.sql_modules AS m ON m.object_id = o.object_id
| WHERE s.schema_id < 16384 AND s.name NOT IN ('guest', 'INFORMATION_SCHEMA', 'sys')
| ORDER BY s.name, o.name, c.column_id
< result
< column "name" "NVARCHAR"
< column "name" "NVARCHAR"
< column "type" "CHAR"
< column "name" "NVARCHAR"
< column "definition" "NVARCHAR"
< row
: string
| dbo
: string
| all_people
: string
| V 
: string
| name
: string
|
| CREATE VIEW dbo.all_people AS SELECT c.name FROM dbo.customers AS c UNION ALL SELECT e.name FROM hr.employees AS e
|
< row
: string
| dbo
: string
| archived
: string
| V 
: string
| id
: string
|
| CREATE VIEW dbo.archived AS SELECT o.id, o.amount FROM archive.dbo.old_orders AS o
|
< row
: string
| dbo
: string
| archived
: string
| V 
: string
| amount
: string
|
| CREATE VIEW dbo.archived AS SELECT o.id, o.amount FROM archive.dbo.old_orders AS o
|
< row
: string
| dbo
: string
| big_orders
: string
| V 
: string
| id
: string
|
| CREATE VIEW dbo.big_orders AS SELECT t.id, t.customer_name, t.gross FROM dbo.order_totals AS t WHERE t.gross > 100
|
< row
: string
| dbo
: string
| big_orders
: string
| V 
: string
| customer_name
: string
|
| CREATE VIEW dbo.big_orders AS SELECT t.id, t.customer_name, t.gross FROM dbo.order_totals AS t WHERE t.gross > 100
|
< row
: string
| dbo
: string
| big_orders
: string
| V 
: string
| gross
: string
|
| CREATE VIEW dbo.big_orders AS SELECT t.id, t.customer_name, t.gross FROM dbo.order_totals AS t WHERE t.gross > 100
|
< row
: string
| dbo
: string
| customers
: string
| U 
: string
| id
: null
< row
: string
| dbo
: string
| customers
: string
| U 
: string
| name
: null
< row
: string
| dbo
: string
| customers
: string
| U 
: string
| email
: null
< row
: string
| dbo
: string
| customers
: string
| U 
: string
| region
: null
< row
: string
| dbo
: string
| MixedCase
: string
| U 
: string
| Id
: null
< row
: string
| dbo
: string
| MixedCase
: string
| U 
: string
| CamelName
: null
< row
: string
| dbo
: string
| Order Details
: string
| U 
: string
| Order ID
: null
< row
: string
| dbo
: string
| Order Details
: string
| U 
: string
| Unit Price
: null
< row
: string
| dbo
: string
| Order Details
: string
| U 
: string
| Qty
: null
< row
: string
| dbo
: string
| order_items
: string
| U 
: string
| order_id
: null
< row
: string
| dbo
: string
| order_items
: string
| U 
: string
| line
: null
< row
: string
| dbo
: string
| order_items
: string
| U 
: string
| product_id
: null
< row
: string
| dbo
: string
| order_items
: string
| U 
: string
| qty
: null
< row
: string
| dbo
: string
| order_items
: string
| U 
: string
| price
: null
< row
: string
| dbo
: string
| order_totals
: string
| V 
: string
| id
: string
|
| CREATE VIEW dbo.order_totals AS SELECT o.id, o.customer_id, o.amount * 1.1 AS gross, c.name AS customer_name FROM dbo.orders AS o JOIN dbo.customers AS c ON c.id = o.customer_id
|
< row
: string
| dbo
: string
| order_totals
: string
| V 
: string
| customer_id
: string
|
| CREATE VIEW dbo.order_totals AS SELECT o.id, o.customer_id, o.amount * 1.1 AS gross, c.name AS customer_name FROM dbo.orders AS o JOIN dbo.customers AS c ON c.id = o.customer_id
|
< row
: string
| dbo
: string
| order_totals
: string
| V 
: string
| gross
: string
|
| CREATE VIEW dbo.order_totals AS SELECT o.id, o.customer_id, o.amount * 1.1 AS gross, c.name AS customer_name FROM dbo.orders AS o JOIN dbo.customers AS c ON c.id = o.customer_id
|
< row
: string
| dbo
: string
| order_totals
: string
| V 
: string
| customer_name
: string
|
| CREATE VIEW dbo.order_totals AS SELECT o.id, o.customer_id, o.amount * 1.1 AS gross, c.name AS customer_name FROM dbo.orders AS o JOIN dbo.customers AS c ON c.id = o.customer_id
|
< row
: string
| dbo
: string
| orders
: string
| U 
: string
| id
: null
< row
: string
| dbo
: string
| orders
: string
| U 
: string
| customer_id
: null
< row
: string
| dbo
: string
| orders
: string
| U 
: string
| amount
: null
< row
: string
| dbo
: string
| orders
: string
| U 
: string
| status
: null
< row
: string
| dbo
: string
| orders
: string
| U 
: string
| created
: null
< row
: string
| dbo
: string
| products
: string
| U 
: string
| id
: null
< row
: string
| dbo
: string
| products
: string
| U 
: string
| name
: null
< row
: string
| dbo
: string
| products
: string
| U 
: string
| category
: null
< row
: string
| dbo
: string
| v_products
: string
| V 
: string
| id
: string
|
| CREATE VIEW dbo.v_products AS SELECT * FROM dbo.products
|
< row
: string
| dbo
: string
| v_products
: string
| V 
: string
| name
: string
|
| CREATE VIEW dbo.v_products AS SELECT * FROM dbo.products
|
< row
: string
| dbo
: string
| v_products
: string
| V 
: string
| category
: string
|
| CREATE VIEW dbo.v_products AS SELECT * FROM dbo.products
|
< row
: string
| dbo
: string
| Vip Names
: string
| V 
: string
| VipName
: string
|
| CREATE VIEW dbo.[Vip Names] AS SELECT c.name AS VipName FROM dbo.customers AS c WHERE c.region = 'vip'
|
< row
: string
| hr
: string
| customers
: string
| U 
: string
| id
: null
< row
: string
| hr
: string
| customers
: string
| U 
: string
| nickname
: null
< row
: string
| hr
: string
| departments
: string
| U 
: string
| id
: null
< row
: string
| hr
: string
| departments
: string
| U 
: string
| name
: null
< row
: string
| hr
: string
| employees
: string
| U 
: string
| id
: null
< row
: string
| hr
: string
| employees
: string
| U 
: string
| name
: null
< row
: string
| hr
: string
| employees
: string
| U 
: string
| salary
: null
< row
: string
| hr
: string
| employees
: string
| U 
: string
| manager_id
: null
< row
: string
| hr
: string
| employees
: string
| U 
: string
| dept_id
: null
< row
: string
| hr
: string
| employees
: string
| U 
: string
| ssn
: null
< row
: string
| hr
: string
| staff
: string
| V 
: string
| staff_id
: string
|
| CREATE VIEW hr.staff (staff_id, staff_name, dept) AS SELECT e.id, e.name, d.name FROM hr.employees AS e LEFT JOIN hr.departments AS d ON d.id = e.dept_id
|
< row
: string
| hr
: string
| staff
: string
| V 
: string
| staff_name
: string
|
| CREATE VIEW hr.staff (staff_id, staff_name, dept) AS SELECT e.id, e.name, d.name FROM hr.employees AS e LEFT JOIN hr.departments AS d ON d.id = e.dept_id
|
< row
: string
| hr
: string
| staff
: string
| V 
: string
| dept
: string
|
| CREATE VIEW hr.staff (staff_id, staff_name, dept) AS SELECT e.id, e.name, d.name FROM hr.employees AS e LEFT JOIN hr.departments AS d ON d.id = e.dept_id
|
> exec
| USE [shop]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.all_people', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| hr
: string
| employees
> exec
| USE [shop]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.archived', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| archive
: string
| dbo
: string
| old_orders
> exec
| USE [shop]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.big_orders', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| order_totals
> exec
| USE [shop]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.order_totals', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| USE [shop]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.v_products', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| products
> exec
| USE [shop]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.Vip Names', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| USE [shop]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'hr.staff', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| hr
: string
| departments
< row
: string
| shop
: string
| hr
: string
| employees
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, amount FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, amount FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, amount FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x06701B425542EC46" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, amount FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT o.id AS order_id, o.amount FROM orders o', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| order_id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT o.id AS order_id, o.amount FROM orders o
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT o.id AS order_id, o.amount FROM orders o" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x61EF5BA38BAC6972" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT o.id AS order_id, o.amount FROM orders o
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id + customer_id, UPPER(status) FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: null
: null
: null
: null
: null
: null
< row
: null
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id + customer_id, UPPER(status) FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id + customer_id, UPPER(status) FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032832" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x6280ACC12BA5309B" QueryPlanHash="0x60822993C95D1A43" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="[shop].[dbo].[orders].[id]+[shop].[dbo].[orders].[customer_id]"><Arithmetic Operation="ADD"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="upper([shop].[dbo].[orders].[status])"><Intrinsic FunctionName="upper"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></Identifier></ScalarOperator></Intrinsic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id + customer_id, UPPER(status) FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT CAST(amount AS varchar(20)) AS a, TRY_CAST(status AS int) AS b FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| a
: null
: null
: null
: null
: null
< row
: string
| b
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT CAST(amount AS varchar(20)) AS a, TRY_CAST(status AS int) AS b FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT CAST(amount AS varchar(20)) AS a, TRY_CAST(status AS int) AS b FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032832" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xC97185C321F9AF7B" QueryPlanHash="0x60822993C95D1A43" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><Warnings><PlanAffectingConvert ConvertIssue="Cardinality Estimate" Expression="CONVERT(varchar(20),[shop].[dbo].[orders].[amount],0)"></PlanAffectingConvert><PlanAffectingConvert ConvertIssue="Cardinality Estimate" Expression="TRY_CAST([shop].[dbo].[orders].[status] AS int)"></PlanAffectingConvert></Warnings><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="25" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CONVERT(varchar(20),[shop].[dbo].[orders].[amount],0)"><Convert DataType="varchar" Length="20" Style="0" Implicit="0"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="TRY_CAST([shop].[dbo].[orders].[status] AS int)"><Convert DataType="int" Style="0" Implicit="0"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="25" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT CAST(amount AS varchar(20)) AS a, TRY_CAST(status AS int) AS b FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT 1 AS one, ''a'' AS b', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| one
: null
: null
: null
: null
: null
< row
: string
| b
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT 1 AS one, 'a' AS b
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT 1 AS one, &apos;a&apos; AS b" StatementId="1" StatementCompId="1" StatementType="SELECT WITHOUT QUERY" RetrievedFromCache="false"></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT 1 AS one, 'a' AS b
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT x = id, y = amount FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| x
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| y
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT x = id, y = amount FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT x = id, y = amount FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x06701B425542EC46" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT x = id, y = amount FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT amount AS id, id AS amount FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT amount AS id, id AS amount FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT amount AS id, id AS amount FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x69A1DD5FA97D77B7" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT amount AS id, id AS amount FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT ''lit'' AS id FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT 'lit' AS id FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT &apos;lit&apos; AS id FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032832" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x42EC44A6A8C31F3A" QueryPlanHash="0x60822993C95D1A43" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="12" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="&apos;lit&apos;"><Const ConstValue="&apos;lit&apos;"></Const></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="9" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT 'lit' AS id FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT NULL AS n, GETDATE() AS now', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| n
: null
: null
: null
: null
: null
< row
: string
| now
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT NULL AS n, GETDATE() AS now
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT NULL AS n, GETDATE() AS now" StatementId="1" StatementCompId="1" StatementType="SELECT WITHOUT QUERY" RetrievedFromCache="false"></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT NULL AS n, GETDATE() AS now
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT CASE WHEN region = ''vip'' THEN name ELSE email END AS contact FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| contact
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT CASE WHEN region = 'vip' THEN name ELSE email END AS contact FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT CASE WHEN region = &apos;vip&apos; THEN name ELSE email END AS contact FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032832" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x444EB6DD48754D28" QueryPlanHash="0x15336BBCA9669F07" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="111" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CASE WHEN [shop].[dbo].[customers].[region]=&apos;vip&apos; THEN [shop].[dbo].[customers].[name] ELSE [shop].[dbo].[customers].[email] END"><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="&apos;vip&apos;"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator></Then><Else><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></Identifier></ScalarOperator></Else></IF></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="170" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT CASE WHEN region = 'vip' THEN name ELSE email END AS contact FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT COALESCE(name, email) AS c, ISNULL(region, ''none'') AS r FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| c
: null
: null
: null
: null
: null
< row
: string
| r
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT COALESCE(name, email) AS c, ISNULL(region, 'none') AS r FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT COALESCE(name, email) AS c, ISNULL(region, &apos;none&apos;) AS r FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032832" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x58D3F2457AA10A90" QueryPlanHash="0x15336BBCA9669F07" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="118" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CASE WHEN [shop].[dbo].[customers].[name] IS NOT NULL THEN [shop].[dbo].[customers].[name] ELSE [shop].[dbo].[customers].[email] END"><IF><Condition><ScalarOperator><Compare CompareOp="IS NOT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator></Then><Else><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></Identifier></ScalarOperator></Else></IF></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="isnull([shop].[dbo].[customers].[region],&apos;none&apos;)"><Intrinsic FunctionName="isnull"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="&apos;none&apos;"></Const></ScalarOperator></Intrinsic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="170" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT COALESCE(name, email) AS c, ISNULL(region, 'none') AS r FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT IIF(amount > 10, status, NULL) AS s FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| s
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT IIF(amount > 10, status, NULL) AS s FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT IIF(amount &gt; 10, status, NULL) AS s FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032832" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x629D2EAAA35D1CFA" QueryPlanHash="0x60822993C95D1A43" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="16" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CASE WHEN [shop].[dbo].[orders].[amount]&gt;(10.00) THEN [shop].[dbo].[orders].[status] ELSE NULL END"><IF><Condition><ScalarOperator><Compare CompareOp="GT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(10.00)"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></Identifier></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Else></IF></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="25" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT IIF(amount > 10, status, NULL) AS s FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT CONCAT(name, ''-'', email) AS ne, name + email AS ne2 FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| ne
: null
: null
: null
: null
: null
< row
: string
| ne2
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT CONCAT(name, '-', email) AS ne, name + email AS ne2 FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT CONCAT(name, &apos;-&apos;, email) AS ne, name + email AS ne2 FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032832" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x22A1571FB347A3CF" QueryPlanHash="0x15336BBCA9669F07" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="314" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="concat([shop].[dbo].[customers].[name],N&apos;-&apos;,[shop].[dbo].[customers].[email])"><Intrinsic FunctionName="concat"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="N&apos;-&apos;"></Const></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></Identifier></ScalarOperator></Intrinsic></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="[shop].[dbo].[customers].[name]+[shop].[dbo].[customers].[email]"><Arithmetic Operation="ADD"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></Identifier></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="163" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT CONCAT(name, '-', email) AS ne, name + email AS ne2 FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id AS [select], name AS "from" FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| select
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| from
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id AS [select], name AS "from" FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id AS [select], name AS &quot;from&quot; FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xF7B23518122D6AF7" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id AS [select], name AS "from" FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT DISTINCT region FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| region
: string
| shop
: string
| dbo
: string
| customers
: string
| region
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT DISTINCT region FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT DISTINCT region FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146444" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x4BCB98D44705AC9E" QueryPlanHash="0x92C8C35FF451D818" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Sort" LogicalOp="Distinct Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100016" AvgRowSize="16" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="1"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="16" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT DISTINCT region FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT TOP 5 id, amount FROM orders ORDER BY amount DESC', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT TOP 5 id, amount FROM orders ORDER BY amount DESC
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT TOP 5 id, amount FROM orders ORDER BY amount DESC" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146444" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x7837D0D3D0E0B623" QueryPlanHash="0x4A34000932D75565" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="16" SerialDesiredMemory="24" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Sort" LogicalOp="TopN Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.00010002" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><TopSort Distinct="0" Rows="5"><OrderBy><OrderByColumn Ascending="0"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></TopSort></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT TOP 5 id, amount FROM orders ORDER BY amount DESC
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT TOP (10) PERCENT name FROM customers ORDER BY name', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT TOP (10) PERCENT name FROM customers ORDER BY name
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT TOP (10) PERCENT name FROM customers ORDER BY name" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146445" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xEF45AED261423D53" QueryPlanHash="0xD4BFD411E0B37980" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Top" LogicalOp="Top" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0146445" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><Top RowCount="0" IsPercent="1" WithTies="0"><TopExpression><ScalarOperator ScalarString="(1.0000000000000000e+001)"><Const ConstValue="(1.0000000000000000e+001)"></Const></ScalarOperator></TopExpression><RelOp NodeId="1" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100061" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></Top></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT TOP (10) PERCENT name FROM customers ORDER BY name
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id FROM orders ORDER BY id OFFSET 5 ROWS FETCH NEXT 5 ROWS ONLY', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id FROM orders ORDER BY id OFFSET 5 ROWS FETCH NEXT 5 ROWS ONLY
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id FROM orders ORDER BY id OFFSET 5 ROWS FETCH NEXT 5 ROWS ONLY" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146446" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x1F1BA412C965255A" QueryPlanHash="0x9BF4CFED8AF7B331" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="16" SerialDesiredMemory="24" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Top" LogicalOp="Top" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-07" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0146446" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><Top RowCount="0" IsPercent="0" WithTies="0"><OffsetExpression><ScalarOperator ScalarString="(5)"><Const ConstValue="(5)"></Const></ScalarOperator></OffsetExpression><TopExpression><ScalarOperator ScalarString="(5)"><Const ConstValue="(5)"></Const></ScalarOperator></TopExpression><RelOp NodeId="1" PhysicalOp="Sort" LogicalOp="TopN Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100011" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><TopSort Distinct="0" Rows="10"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></TopSort></RelOp></Top></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id FROM orders ORDER BY id OFFSET 5 ROWS FETCH NEXT 5 ROWS ONLY
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id FROM orders WITH (NOLOCK)', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id FROM orders WITH (NOLOCK)
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id FROM orders WITH (NOLOCK)" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xEFA15ED25A4F9F7C" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id FROM orders WITH (NOLOCK)
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id FROM dbo.orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id FROM dbo.orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id FROM dbo.orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x3BB2392D98C6046B" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id FROM dbo.orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id FROM shop.dbo.orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id FROM shop.dbo.orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id FROM shop.dbo.orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x3BB2392D98C6046B" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id FROM shop.dbo.orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id FROM shop..orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id FROM shop..orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id FROM shop..orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x3BB2392D98C6046B" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id FROM shop..orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT orders.id, dbo.orders.amount FROM dbo.orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT orders.id, dbo.orders.amount FROM dbo.orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT orders.id, dbo.orders.amount FROM dbo.orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x06701B425542EC46" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT orders.id, dbo.orders.amount FROM dbo.orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| email
: string
| shop
: string
| dbo
: string
| customers
: string
| email
: null
< row
: string
| region
: string
| shop
: string
| dbo
: string
| customers
: string
| region
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xDD63C9CBA14D04B0" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="174" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT c.* FROM customers c', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| email
: string
| shop
: string
| dbo
: string
| customers
: string
| email
: null
< row
: string
| region
: string
| shop
: string
| dbo
: string
| customers
: string
| region
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT c.* FROM customers c
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT c.* FROM customers c" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x7C2A72F2C23C9694" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="174" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="email"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="email"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT c.* FROM customers c
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT *, id AS again FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| email
: string
| shop
: string
| dbo
: string
| customers
: string
| email
: null
< row
: string
| region
: string
| shop
: string
| dbo
: string
| customers
: string
| region
: null
< row
: string
| again
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT *, id AS again FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT *, id AS again FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x4C43D91115BA5C9B" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="174" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="email"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT *, id AS again FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| customer_id
: string
| shop
: string
| dbo
: string
| orders
: string
| customer_id
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
< row
: string
| status
: string
| shop
: string
| dbo
: string
| orders
: string
| status
: null
< row
: string
| created
: string
| shop
: string
| dbo
: string
| orders
: string
| created
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243596" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xE1F7C5A1A03A6146" QueryPlanHash="0x1B3BE2B12828744F" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="48" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177904" AvgRowSize="93" EstimatedTotalSubtreeCost="0.0243596" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="created"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="41" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="created"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="created"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT o.*, c.name FROM orders o JOIN customers c ON c.id = o.customer_id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM orders o JOIN customers c ON c.id = o.customer_id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| customer_id
: string
| shop
: string
| dbo
: string
| orders
: string
| customer_id
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
< row
: string
| status
: string
| shop
: string
| dbo
: string
| orders
: string
| status
: null
< row
: string
| created
: string
| shop
: string
| dbo
: string
| orders
: string
| created
: null
< row
: string
| id
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| email
: string
| shop
: string
| dbo
: string
| customers
: string
| email
: null
< row
: string
| region
: string
| shop
: string
| dbo
: string
| customers
: string
| region
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM orders o JOIN customers c ON c.id = o.customer_id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM orders o JOIN customers c ON c.id = o.customer_id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243596" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x079674C6D7E31A8C" QueryPlanHash="0x1B3BE2B12828744F" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="48" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177904" AvgRowSize="207" EstimatedTotalSubtreeCost="0.0243596" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="created"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="email"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="41" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="created"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="created"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="174" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="email"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="email"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM orders o JOIN customers c ON c.id = o.customer_id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT o.*, c.* FROM orders o, customers c WHERE c.id = o.customer_id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| customer_id
: string
| shop
: string
| dbo
: string
| orders
: string
| customer_id
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
< row
: string
| status
: string
| shop
: string
| dbo
: string
| orders
: string
| status
: null
< row
: string
| created
: string
| shop
: string
| dbo
: string
| orders
: string
| created
: null
< row
: string
| id
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| email
: string
| shop
: string
| dbo
: string
| customers
: string
| email
: null
< row
: string
| region
: string
| shop
: string
| dbo
: string
| customers
: string
| region
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT o.*, c.* FROM orders o, customers c WHERE c.id = o.customer_id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT o.*, c.* FROM orders o, customers c WHERE c.id = o.customer_id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243596" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xC83847CECA634AED" QueryPlanHash="0x1B3BE2B12828744F" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="48" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177904" AvgRowSize="207" EstimatedTotalSubtreeCost="0.0243596" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="created"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="email"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="41" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="created"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="created"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="174" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="email"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="email"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT o.*, c.* FROM orders o, customers c WHERE c.id = o.customer_id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM hr.employees', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| hr
: string
| employees
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| hr
: string
| employees
: string
| name
: null
< row
: string
| salary
: string
| shop
: string
| hr
: string
| employees
: string
| salary
: null
< row
: string
| manager_id
: string
| shop
: string
| hr
: string
| employees
: string
| manager_id
: null
< row
: string
| dept_id
: string
| shop
: string
| hr
: string
| employees
: string
| dept_id
: null
< row
: string
| ssn
: string
| shop
: string
| hr
: string
| employees
: string
| ssn
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM hr.employees
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM hr.employees" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x9CBC07D5AC7EDCF7" QueryPlanHash="0xEC3C9F9DBC9AF923" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="93" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="salary"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="manager_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="ssn"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="salary"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="manager_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="ssn"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM hr.employees
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| hr
: string
| employees
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT TOP 3 * FROM orders ORDER BY amount DESC', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| customer_id
: string
| shop
: string
| dbo
: string
| orders
: string
| customer_id
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
< row
: string
| status
: string
| shop
: string
| dbo
: string
| orders
: string
| status
: null
< row
: string
| created
: string
| shop
: string
| dbo
: string
| orders
: string
| created
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT TOP 3 * FROM orders ORDER BY amount DESC
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT TOP 3 * FROM orders ORDER BY amount DESC" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146444" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xF1855C5BA5C269D6" QueryPlanHash="0x94147CF428922869" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="16" SerialDesiredMemory="24" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Sort" LogicalOp="TopN Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100041" AvgRowSize="41" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><TopSort Distinct="0" Rows="3"><OrderBy><OrderByColumn Ascending="0"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="41" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></TopSort></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT TOP 3 * FROM orders ORDER BY amount DESC
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT o.id, c.name FROM orders o INNER JOIN customers c ON o.customer_id = c.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT o.id, c.name FROM orders o INNER JOIN customers c ON o.customer_id = c.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT o.id, c.name FROM orders o INNER JOIN customers c ON o.customer_id = c.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243434" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x00D02CA320F6602A" QueryPlanHash="0x1B3BE2B12828744F" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177742" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0243434" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[customers].[id] as [c].[id]=[shop].[dbo].[orders].[customer_id] as [o].[customer_id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT o.id, c.name FROM orders o INNER JOIN customers c ON o.customer_id = c.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON o.customer_id = c.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON o.customer_id = c.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON o.customer_id = c.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657068" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x511ADD6D9B9D0BB2" QueryPlanHash="0x53C1CF2444F7212B" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Outer Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="65" EstimatedTotalSubtreeCost="0.00657068" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT o.id, c.name FROM orders o LEFT JOIN customers c ON o.customer_id = c.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT o.id, c.name FROM orders o RIGHT OUTER JOIN customers c ON o.customer_id = c.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT o.id, c.name FROM orders o RIGHT OUTER JOIN customers c ON o.customer_id = c.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT o.id, c.name FROM orders o RIGHT OUTER JOIN customers c ON o.customer_id = c.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657068" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x75BC8670B2054BEC" QueryPlanHash="0x8DBEB4C0994882C3" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Outer Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="65" EstimatedTotalSubtreeCost="0.00657068" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT o.id, c.name FROM orders o RIGHT OUTER JOIN customers c ON o.customer_id = c.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT o.id, c.name FROM orders o FULL JOIN customers c ON o.customer_id = c.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT o.id, c.name FROM orders o FULL JOIN customers c ON o.customer_id = c.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT o.id, c.name FROM orders o FULL JOIN customers c ON o.customer_id = c.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0131417" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xBA61E05B454B53C0" QueryPlanHash="0xCA18024ABE718A4D" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="32" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Concatenation" LogicalOp="Concatenation" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-07" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0131417" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><Concat><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Nested Loops" LogicalOp="Left Outer Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="65" EstimatedTotalSubtreeCost="0.00657068" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp><RelOp NodeId="4" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="65" EstimatedTotalSubtreeCost="0.00657078" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ScalarOperator ScalarString="NULL"><Const ConstValue="NULL"></Const></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="5" PhysicalOp="Nested Loops" LogicalOp="Left Anti Semi Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="61" EstimatedTotalSubtreeCost="0.00657068" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="6" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="7" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></ComputeScalar></RelOp></Concat></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT o.id, c.name FROM orders o FULL JOIN customers c ON o.customer_id = c.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT o.id, p.name FROM orders o CROSS JOIN products p', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| products
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT o.id, p.name FROM orders o CROSS JOIN products p
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT o.id, p.name FROM orders o CROSS JOIN products p" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657038" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x8D8C42DC9B53A4D5" QueryPlanHash="0xFCA52044F94BFB03" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="65" EstimatedTotalSubtreeCost="0.00657038" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT o.id, p.name FROM orders o CROSS JOIN products p
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
< row
: string
| shop
: string
| dbo
: string
| products
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT e.name, m.name AS manager FROM hr.employees e LEFT JOIN hr.employees m ON m.id = e.manager_id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| hr
: string
| employees
: string
| name
: null
< row
: string
| manager
: string
| shop
: string
| hr
: string
| employees
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT e.name, m.name AS manager FROM hr.employees e LEFT JOIN hr.employees m ON m.id = e.manager_id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT e.name, m.name AS manager FROM hr.employees e LEFT JOIN hr.employees m ON m.id = e.manager_id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657068" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xFF2083A05C7F20B8" QueryPlanHash="0x084ADEEE6336B9B7" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Outer Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="113" EstimatedTotalSubtreeCost="0.00657068" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[m]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[hr].[employees].[id] as [m].[id]=[shop].[hr].[employees].[manager_id] as [e].[manager_id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[m]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="manager_id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="manager_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="manager_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[m]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[m]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[m]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[m]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[m]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT e.name, m.name AS manager FROM hr.employees e LEFT JOIN hr.employees m ON m.id = e.manager_id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| hr
: string
| employees
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT c.name, p.name AS product, i.qty FROM customers c JOIN orders o ON o.customer_id = c.id JOIN order_items i ON i.order_id = o.id JOIN products p ON p.id = i.product_id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| product
: string
| shop
: string
| dbo
: string
| products
: string
| name
: null
< row
: string
| qty
: string
| shop
: string
| dbo
: string
| order_items
: string
| qty
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT c.name, p.name AS product, i.qty FROM customers c JOIN orders o ON o.customer_id = c.id JOIN order_items i ON i.order_id = o.id JOIN products p ON p.id = i.product_id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT c.name, p.name AS product, i.qty FROM customers c JOIN orders o ON o.customer_id = c.id JOIN order_items i ON i.order_id = o.id JOIN products p ON p.id = i.product_id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.066499" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x3C0FF480AC857017" QueryPlanHash="0xAE9065C168602233" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="96" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="3072" SerialDesiredMemory="3168" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0178065" AvgRowSize="117" EstimatedTotalSubtreeCost="0.066499" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="qty"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="0.333333" Output="0.333333"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177742" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0454065" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="qty"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="0.333333" Output="0.333333"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="order_id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[order_items].[order_id] as [i].[order_id]=[shop].[dbo].[orders].[id] as [o].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="order_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="4" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177769" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0243461" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="order_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="qty"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="0.333333" Output="0.333333"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="product_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[order_items].[product_id] as [i].[product_id]=[shop].[dbo].[products].[id] as [p].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="product_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="5" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="19" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="order_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="product_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="qty"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="order_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="product_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="qty"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="6" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[products]" Alias="[p]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></Hash></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT c.name, p.name AS product, i.qty FROM customers c JOIN orders o ON o.customer_id = c.id JOIN order_items i ON i.order_id = o.id JOIN products p ON p.id = i.product_id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| order_items
< row
: string
| shop
: string
| dbo
: string
| orders
< row
: string
| shop
: string
| dbo
: string
| products
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT o.id, c.name FROM orders o JOIN (customers c JOIN hr.customers h ON h.id = c.id) ON c.id = o.customer_id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT o.id, c.name FROM orders o JOIN (customers c JOIN hr.customers h ON h.id = c.id) ON c.id = o.customer_id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT o.id, c.name FROM orders o JOIN (customers c JOIN hr.customers h ON h.id = c.id) ON c.id = o.customer_id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0454011" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x8A103699AEB48A45" QueryPlanHash="0xBFFB0318EFDEA244" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="64" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="2048" SerialDesiredMemory="2112" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177716" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0454011" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="0.5" Output="0.5"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" Column="id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[hr].[customers].[id] as [h].[id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177742" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0243434" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="0.5" Output="0.5"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="4" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT o.id, c.name FROM orders o JOIN (customers c JOIN hr.customers h ON h.id = c.id) ON c.id = o.customer_id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
< row
: string
| shop
: string
| hr
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT c.name, x.amount FROM customers c CROSS APPLY (SELECT TOP 1 o.amount FROM orders o WHERE o.customer_id = c.id ORDER BY o.amount DESC) x', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT c.name, x.amount FROM customers c CROSS APPLY (SELECT TOP 1 o.amount FROM orders o WHERE o.customer_id = c.id ORDER BY o.amount DESC) x
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT c.name, x.amount FROM customers c CROSS APPLY (SELECT TOP 1 o.amount FROM orders o WHERE o.customer_id = c.id ORDER BY o.amount DESC) x" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0179322" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x5C820F79BD2BC601" QueryPlanHash="0xC119894561FFA3EB" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="16" SerialDesiredMemory="24" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="70" EstimatedTotalSubtreeCost="0.0179322" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><NestedLoops Optimized="0"><OuterReferences><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></OuterReferences><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Sort" LogicalOp="TopN Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100016" AvgRowSize="16" EstimatedTotalSubtreeCost="0.0146449" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><TopSort Distinct="0" Rows="1"><OrderBy><OrderByColumn Ascending="0"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></TopSort></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT c.name, x.amount FROM customers c CROSS APPLY (SELECT TOP 1 o.amount FROM orders o WHERE o.customer_id = c.id ORDER BY o.amount DESC) x
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT c.name, x.n FROM customers c OUTER APPLY (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| n
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT c.name, x.n FROM customers c OUTER APPLY (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT c.name, x.n FROM customers c OUTER APPLY (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657206" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xDD71AEA9072BF16F" QueryPlanHash="0x737532E279E9DAF4" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Outer Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="65" EstimatedTotalSubtreeCost="0.00657206" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Column="Expr1005"></ColumnReference></OutputList><NestedLoops Optimized="0"><OuterReferences><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></OuterReferences><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="11" EstimatedTotalSubtreeCost="0.00328478" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1005"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1005"></ColumnReference><ScalarOperator ScalarString="[Expr1004]"><Identifier><ColumnReference Column="Expr1004"></ColumnReference></Identifier></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="3" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="11" EstimatedTotalSubtreeCost="0.00328468" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="CONVERT_IMPLICIT(int,[Expr1008],0)"><Convert DataType="int" Style="0" Implicit="1"><ScalarOperator><Identifier><ColumnReference Column="Expr1008"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="4" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="11" EstimatedTotalSubtreeCost="0.00328468" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1008"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1008"></ColumnReference><ScalarOperator ScalarString="Count(*)"><Aggregate AggType="countstar" Distinct="0"></Aggregate></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="5" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp></ComputeScalar></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT c.name, x.n FROM customers c OUTER APPLY (SELECT COUNT(*) AS n FROM orders o WHERE o.customer_id = c.id) x
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS cname FROM orders o', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| cname
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS cname FROM orders o
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS cname FROM orders o" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657254" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x57DE6108D674C8C6" QueryPlanHash="0x7B4D8BE32C5D117D" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="32" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="65" EstimatedTotalSubtreeCost="0.00657254" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Column="Expr1005"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1005"></ColumnReference><ScalarOperator ScalarString="[Expr1007]"><Identifier><ColumnReference Column="Expr1007"></ColumnReference></Identifier></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Nested Loops" LogicalOp="Left Outer Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="65" EstimatedTotalSubtreeCost="0.00657244" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Column="Expr1007"></ColumnReference></OutputList><NestedLoops Optimized="0"><OuterReferences><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OuterReferences><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="3" PhysicalOp="Assert" LogicalOp="Assert" EstimateRows="1" EstimateIO="0" EstimateCPU="4.8e-07" AvgRowSize="61" EstimatedTotalSubtreeCost="0.00328516" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1007"></ColumnReference></OutputList><Assert StartupExpression="0"><RelOp NodeId="4" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="69" EstimatedTotalSubtreeCost="0.00328468" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1006"></ColumnReference><ColumnReference Column="Expr1007"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1006"></ColumnReference><ScalarOperator ScalarString="Count(*)"><Aggregate AggType="countstar" Distinct="0"></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1007"></ColumnReference><ScalarOperator ScalarString="ANY([shop].[dbo].[customers].[name] as [c].[name])"><Aggregate Distinct="0" AggType="ANY"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="5" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[customers].[id] as [c].[id]=[shop].[dbo].[orders].[customer_id] as [o].[customer_id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></StreamAggregate></RelOp><Predicate><ScalarOperator ScalarString="CASE WHEN [Expr1006]&gt;(1) THEN (0) ELSE NULL END"><IF><Condition><ScalarOperator><Compare CompareOp="GT"><ScalarOperator><Identifier><ColumnReference Column="Expr1006"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Const ConstValue="(0)"></Const></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Else></IF></ScalarOperator></Predicate></Assert></RelOp></NestedLoops></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, (SELECT name FROM customers c WHERE c.id = o.customer_id) AS cname FROM orders o
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT (SELECT COUNT(*) FROM orders) AS n', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| n
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT (SELECT COUNT(*) FROM orders) AS n
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT (SELECT COUNT(*) FROM orders) AS n" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032843" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xD34B14A0BE4D39FD" QueryPlanHash="0xF5703E9D599A22D3" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032843" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1005"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1005"></ColumnReference><ScalarOperator ScalarString="[Expr1003]"><Identifier><ColumnReference Column="Expr1003"></ColumnReference></Identifier></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032842" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CONVERT_IMPLICIT(int,[Expr1008],0)"><Convert DataType="int" Style="0" Implicit="1"><ScalarOperator><Identifier><ColumnReference Column="Expr1008"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032842" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1008"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1008"></ColumnReference><ScalarOperator ScalarString="Count(*)"><Aggregate AggType="countstar" Distinct="0"></Aggregate></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="9" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT (SELECT COUNT(*) FROM orders) AS n
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM (SELECT id, amount FROM orders) AS t', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| amount
: string
| shop
: string
| dbo
: string
| orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM (SELECT id, amount FROM orders) AS t
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM (SELECT id, amount FROM orders) AS t" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x34F2B333F744CEEE" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM (SELECT id, amount FROM orders) AS t
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT t.x FROM (SELECT id AS x FROM orders) t', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| x
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT t.x FROM (SELECT id AS x FROM orders) t
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT t.x FROM (SELECT id AS x FROM orders) t" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xCFFDC5820041AF57" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT t.x FROM (SELECT id AS x FROM orders) t
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT a, b FROM (SELECT id, name FROM customers) AS t(a, b)', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| a
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| b
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT a, b FROM (SELECT id, name FROM customers) AS t(a, b)
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT a, b FROM (SELECT id, name FROM customers) AS t(a, b)" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xE7F7D4B5334DEB19" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT a, b FROM (SELECT id, name FROM customers) AS t(a, b)
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM (SELECT * FROM (SELECT id, name FROM customers) a) b', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM (SELECT * FROM (SELECT id, name FROM customers) a) b
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM (SELECT * FROM (SELECT id, name FROM customers) a) b" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x3DA02563F2DA71BD" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM (SELECT * FROM (SELECT id, name FROM customers) a) b
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders)', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders)
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders)" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657068" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xF3145AC063959006" QueryPlanHash="0xA0EC104BE9123A12" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Semi Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="61" EstimatedTotalSubtreeCost="0.00657068" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[customers].[id]=[shop].[dbo].[orders].[customer_id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders)
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT name FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT name FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT name FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657068" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xBE0A05EA3FC63201" QueryPlanHash="0xA0EC104BE9123A12" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Semi Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="61" EstimatedTotalSubtreeCost="0.00657068" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT name FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT name FROM customers WHERE id = (SELECT MAX(customer_id) FROM orders)', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT name FROM customers WHERE id = (SELECT MAX(customer_id) FROM orders)
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT name FROM customers WHERE id = (SELECT MAX(customer_id) FROM orders)" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657178" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xA6A341AE1477C8D3" QueryPlanHash="0x788F90D339180A6C" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="61" EstimatedTotalSubtreeCost="0.00657178" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[Expr1006]=[shop].[dbo].[customers].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Column="Expr1006"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032842" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1006"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1006"></ColumnReference><ScalarOperator ScalarString="MAX([shop].[dbo].[orders].[customer_id])"><Aggregate Distinct="0" AggType="MAX"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></StreamAggregate></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT name FROM customers WHERE id = (SELECT MAX(customer_id) FROM orders)
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, CASE WHEN EXISTS (SELECT 1 FROM order_items i WHERE i.order_id = o.id) THEN ''y'' ELSE ''n'' END AS has_items FROM orders o', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| has_items
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, CASE WHEN EXISTS (SELECT 1 FROM order_items i WHERE i.order_id = o.id) THEN 'y' ELSE 'n' END AS has_items FROM orders o
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, CASE WHEN EXISTS (SELECT 1 FROM order_items i WHERE i.order_id = o.id) THEN &apos;y&apos; ELSE &apos;n&apos; END AS has_items FROM orders o" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657096" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xF0D3DEE3247D7610" QueryPlanHash="0x05228F3F0DA32598" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="15" EstimatedTotalSubtreeCost="0.00657096" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Column="Expr1005"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1005"></ColumnReference><ScalarOperator ScalarString="CASE WHEN [Expr1006] THEN &apos;y&apos; ELSE &apos;n&apos; END"><IF><Condition><ScalarOperator><Identifier><ColumnReference Column="Expr1006"></ColumnReference></Identifier></ScalarOperator></Condition><Then><ScalarOperator><Const ConstValue="&apos;y&apos;"></Const></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="&apos;n&apos;"></Const></ScalarOperator></Else></IF></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Nested Loops" LogicalOp="Left Semi Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="12" EstimatedTotalSubtreeCost="0.00657086" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Column="Expr1006"></ColumnReference></OutputList><NestedLoops Optimized="0"><DefinedValues><DefinedValue><ColumnReference Column="Expr1006"></ColumnReference></DefinedValue></DefinedValues><OuterReferences><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></OuterReferences><ProbeColumn><ColumnReference Column="Expr1006"></ColumnReference></ProbeColumn><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[order_items].[order_id] as [i].[order_id]=[shop].[dbo].[orders].[id] as [o].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Alias="[i]" Column="order_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></NestedLoops></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, CASE WHEN EXISTS (SELECT 1 FROM order_items i WHERE i.order_id = o.id) THEN 'y' ELSE 'n' END AS has_items FROM orders o
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| order_items
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT v.k FROM (VALUES (1), (2)) AS v(k)', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| k
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT v.k FROM (VALUES (1), (2)) AS v(k)
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT v.k FROM (VALUES (1), (2)) AS v(k)" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="2.157e-06" StatementEstRows="2" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x6BCF24DD589CF0B1" QueryPlanHash="0xD3F9F22C1BAD3015" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Constant Scan" LogicalOp="Constant Scan" EstimateRows="2" EstimateIO="0" EstimateCPU="2.157e-06" AvgRowSize="11" EstimatedTotalSubtreeCost="2.157e-06" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1002"></ColumnReference></OutputList><ConstantScan><Values><Row><ScalarOperator ScalarString="(1)"><Const ConstValue="(1)"></Const></ScalarOperator></Row><Row><ScalarOperator ScalarString="(2)"><Const ConstValue="(2)"></Const></ScalarOperator></Row></Values></ConstantScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT v.k FROM (VALUES (1), (2)) AS v(k)
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT u.x FROM (SELECT id AS x FROM orders UNION ALL SELECT id FROM customers) u', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| x
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT u.x FROM (SELECT id AS x FROM orders UNION ALL SELECT id FROM customers) u
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT u.x FROM (SELECT id AS x FROM orders UNION ALL SELECT id FROM customers) u" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0065664" StatementEstRows="2" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x0908BDF48A3C6FC4" QueryPlanHash="0xBD9CDDF1353635C1" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Concatenation" LogicalOp="Concatenation" EstimateRows="2" EstimateIO="0" EstimateCPU="2e-07" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0065664" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1006"></ColumnReference></OutputList><Concat><DefinedValues><DefinedValue><ColumnReference Column="Union1006"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Concat></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT u.x FROM (SELECT id AS x FROM orders UNION ALL SELECT id FROM customers) u
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'WITH c AS (SELECT id, name FROM customers) SELECT * FROM c', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| WITH c AS (SELECT id, name FROM customers) SELECT * FROM c
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="WITH c AS (SELECT id, name FROM customers) SELECT * FROM c" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x64C2A067434387C3" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| WITH c AS (SELECT id, name FROM customers) SELECT * FROM c
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'WITH c(x, y) AS (SELECT id, name FROM customers) SELECT x, y FROM c', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| x
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| y
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| WITH c(x, y) AS (SELECT id, name FROM customers) SELECT x, y FROM c
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="WITH c(x, y) AS (SELECT id, name FROM customers) SELECT x, y FROM c" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x64C2A067434387C3" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| WITH c(x, y) AS (SELECT id, name FROM customers) SELECT x, y FROM c
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT name, id FROM b', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT name, id FROM b
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT name, id FROM b" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243434" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x1A34EB2D8957962F" QueryPlanHash="0x1B3BE2B12828744F" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177742" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0243434" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| WITH a AS (SELECT id, customer_id FROM orders), b AS (SELECT a.id, c.name FROM a JOIN customers c ON c.id = a.customer_id) SELECT name, id FROM b
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'WITH c AS (SELECT id FROM orders) SELECT x.id FROM c x JOIN c y ON x.id = y.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| WITH c AS (SELECT id FROM orders) SELECT x.id FROM c x JOIN c y ON x.id = y.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="WITH c AS (SELECT id FROM orders) SELECT x.id FROM c x JOIN c y ON x.id = y.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243408" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xC6D1D4FD5F9DC7DC" QueryPlanHash="0x0AB88DAD93FF5FB4" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177716" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0243408" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[id]=[shop].[dbo].[orders].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" TableReferenceId="2" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" TableReferenceId="1" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| WITH c AS (SELECT id FROM orders) SELECT x.id FROM c x JOIN c y ON x.id = y.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'WITH orders AS (SELECT id, name FROM customers) SELECT * FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| WITH orders AS (SELECT id, name FROM customers) SELECT * FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="WITH orders AS (SELECT id, name FROM customers) SELECT * FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x64C2A067434387C3" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| WITH orders AS (SELECT id, name FROM customers) SELECT * FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'WITH r AS (SELECT id, manager_id, name FROM hr.employees WHERE manager_id IS NULL UNION ALL SELECT e.id, e.manager_id, e.name FROM hr.employees e JOIN r ON e.manager_id = r.id) SELECT * FROM r', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: null
: null
: null
: null
: null
< row
: string
| manager_id
: null
: null
: null
: null
: null
< row
: string
| name
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| WITH r AS (SELECT id, manager_id, name FROM hr.employees WHERE manager_id IS NULL UNION ALL SELECT e.id, e.manager_id, e.name FROM hr.employees e JOIN r ON e.manager_id = r.id) SELECT * FROM r
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="WITH r AS (SELECT id, manager_id, name FROM hr.employees WHERE manager_id IS NULL UNION ALL SELECT e.id, e.manager_id, e.name FROM hr.employees e JOIN r ON e.manager_id = r.id) SELECT * FROM r" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00688768" StatementEstRows="4" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xE5A9BC87D270AAC6" QueryPlanHash="0x4D26CE879151BC55" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="32" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Index Spool" LogicalOp="Lazy Spool" EstimateRows="4" EstimateIO="0" EstimateCPU="1e-08" AvgRowSize="69" EstimatedTotalSubtreeCost="0.00688768" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1014"></ColumnReference><ColumnReference Column="Recr1008"></ColumnReference><ColumnReference Column="Recr1009"></ColumnReference><ColumnReference Column="Recr1010"></ColumnReference></OutputList><Spool Stack="1"><RelOp NodeId="1" PhysicalOp="Concatenation" LogicalOp="Concatenation" EstimateRows="4" EstimateIO="0" EstimateCPU="2e-09" AvgRowSize="69" EstimatedTotalSubtreeCost="0.00680759" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1014"></ColumnReference><ColumnReference Column="Recr1008"></ColumnReference><ColumnReference Column="Recr1009"></ColumnReference><ColumnReference Column="Recr1010"></ColumnReference></OutputList><Concat><DefinedValues><DefinedValue><ColumnReference Column="Expr1014"></ColumnReference><ColumnReference Column="Expr1011"></ColumnReference><ColumnReference Column="Expr1013"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Column="Recr1008"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Column="Recr1009"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="manager_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="manager_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Column="Recr1010"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="69" EstimatedTotalSubtreeCost="2e-08" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1011"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="manager_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1011"></ColumnReference><ScalarOperator ScalarString="(0)"><Const ConstValue="(0)"></Const></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="manager_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="manager_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[hr].[employees].[manager_id] IS NULL"><Compare CompareOp="IS"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="manager_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></ComputeScalar></RelOp><RelOp NodeId="4" PhysicalOp="Assert" LogicalOp="Assert" EstimateRows="4" EstimateIO="0" EstimateCPU="1.68e-07" AvgRowSize="69" EstimatedTotalSubtreeCost="0.00352401" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1013"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="manager_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></OutputList><Assert StartupExpression="0"><RelOp NodeId="5" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="4" EstimateIO="0" EstimateCPU="1.68e-07" AvgRowSize="69" EstimatedTotalSubtreeCost="0.00352401" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1013"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="manager_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><OuterReferences><ColumnReference Column="Expr1013"></ColumnReference><ColumnReference Column="Recr1005"></ColumnReference><ColumnReference Column="Recr1006"></ColumnReference><ColumnReference Column="Recr1007"></ColumnReference></OuterReferences><RelOp NodeId="6" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="69" EstimatedTotalSubtreeCost="2e-08" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1013"></ColumnReference><ColumnReference Column="Recr1005"></ColumnReference><ColumnReference Column="Recr1006"></ColumnReference><ColumnReference Column="Recr1007"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1013"></ColumnReference><ScalarOperator ScalarString="[Expr1012]+(1)"><Arithmetic Operation="ADD"><ScalarOperator><Identifier><ColumnReference Column="Expr1012"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="7" PhysicalOp="Table Spool" LogicalOp="Lazy Spool" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="69" EstimatedTotalSubtreeCost="2e-08" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1012"></ColumnReference><ColumnReference Column="Recr1005"></ColumnReference><ColumnReference Column="Recr1006"></ColumnReference><ColumnReference Column="Recr1007"></ColumnReference></OutputList><Spool PrimaryNodeId="0" Stack="1"></Spool></RelOp></ComputeScalar></RelOp><RelOp NodeId="8" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0035219" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="3" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="manager_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="manager_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[hr].[employees].[manager_id] as [e].[manager_id]=[Recr1005]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="manager_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Column="Recr1005"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></NestedLoops></RelOp><Predicate><ScalarOperator ScalarString="CASE WHEN [Expr1013]&gt;(100) THEN (0) ELSE NULL END"><IF><Condition><ScalarOperator><Compare CompareOp="GT"><ScalarOperator><Identifier><ColumnReference Column="Expr1013"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(100)"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Const ConstValue="(0)"></Const></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Else></IF></ScalarOperator></Predicate></Assert></RelOp></Concat></RelOp></Spool></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| WITH r AS (SELECT id, manager_id, name FROM hr.employees WHERE manager_id IS NULL UNION ALL SELECT e.id, e.manager_id, e.name FROM hr.employees e JOIN r ON e.manager_id = r.id) SELECT * FROM r
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| hr
: string
| employees
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'WITH r(n, label) AS (SELECT 1, CAST(name AS nvarchar(100)) FROM products UNION ALL SELECT n + 1, label FROM r WHERE n < 3) SELECT n, label FROM r', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| n
: null
: null
: null
: null
: null
< row
: string
| label
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| WITH r(n, label) AS (SELECT 1, CAST(name AS nvarchar(100)) FROM products UNION ALL SELECT n + 1, label FROM r WHERE n < 3) SELECT n, label FROM r
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="WITH r(n, label) AS (SELECT 1, CAST(name AS nvarchar(100)) FROM products UNION ALL SELECT n + 1, label FROM r WHERE n &lt; 3) SELECT n, label FROM r" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00329146" StatementEstRows="4" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xF8EABAA5E9AA0B28" QueryPlanHash="0xC515F92589417CA3" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="32" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Index Spool" LogicalOp="Lazy Spool" EstimateRows="4" EstimateIO="0" EstimateCPU="1e-08" AvgRowSize="115" EstimatedTotalSubtreeCost="0.00329146" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1013"></ColumnReference><ColumnReference Column="Recr1008"></ColumnReference><ColumnReference Column="Recr1009"></ColumnReference></OutputList><Spool Stack="1"><RelOp NodeId="1" PhysicalOp="Concatenation" LogicalOp="Concatenation" EstimateRows="4" EstimateIO="0" EstimateCPU="2e-09" AvgRowSize="115" EstimatedTotalSubtreeCost="0.00328987" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1013"></ColumnReference><ColumnReference Column="Recr1008"></ColumnReference><ColumnReference Column="Recr1009"></ColumnReference></OutputList><Concat><DefinedValues><DefinedValue><ColumnReference Column="Expr1013"></ColumnReference><ColumnReference Column="Expr1010"></ColumnReference><ColumnReference Column="Expr1012"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Column="Recr1008"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference><ColumnReference Column="Expr1007"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Column="Recr1009"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference><ColumnReference Column="Recr1006"></ColumnReference></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="115" EstimatedTotalSubtreeCost="2e-08" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1010"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1010"></ColumnReference><ScalarOperator ScalarString="(0)"><Const ConstValue="(0)"></Const></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="3" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="115" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="(1)"><Const ConstValue="(1)"></Const></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="CONVERT(nvarchar(100),[shop].[dbo].[products].[name],0)"><Convert DataType="nvarchar" Length="200" Style="0" Implicit="0"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="4" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[products]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></ComputeScalar></RelOp><RelOp NodeId="5" PhysicalOp="Assert" LogicalOp="Assert" EstimateRows="4" EstimateIO="0" EstimateCPU="1.68e-07" AvgRowSize="115" EstimatedTotalSubtreeCost="6.665e-06" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1012"></ColumnReference><ColumnReference Column="Expr1007"></ColumnReference><ColumnReference Column="Recr1006"></ColumnReference></OutputList><Assert StartupExpression="0"><RelOp NodeId="6" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="4" EstimateIO="0" EstimateCPU="1.68e-07" AvgRowSize="115" EstimatedTotalSubtreeCost="6.665e-06" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1012"></ColumnReference><ColumnReference Column="Expr1007"></ColumnReference><ColumnReference Column="Recr1006"></ColumnReference></OutputList><NestedLoops Optimized="0"><OuterReferences><ColumnReference Column="Expr1012"></ColumnReference><ColumnReference Column="Recr1005"></ColumnReference><ColumnReference Column="Recr1006"></ColumnReference></OuterReferences><RelOp NodeId="7" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="115" EstimatedTotalSubtreeCost="2e-08" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1012"></ColumnReference><ColumnReference Column="Recr1005"></ColumnReference><ColumnReference Column="Recr1006"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1012"></ColumnReference><ScalarOperator ScalarString="[Expr1011]+(1)"><Arithmetic Operation="ADD"><ScalarOperator><Identifier><ColumnReference Column="Expr1011"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="8" PhysicalOp="Table Spool" LogicalOp="Lazy Spool" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="115" EstimatedTotalSubtreeCost="2e-08" Parallel="0" EstimateRebinds="4" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1011"></ColumnReference><ColumnReference Column="Recr1005"></ColumnReference><ColumnReference Column="Recr1006"></ColumnReference></OutputList><Spool PrimaryNodeId="0" Stack="1"></Spool></RelOp></ComputeScalar></RelOp><RelOp NodeId="9" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="11" EstimatedTotalSubtreeCost="6.477e-06" Parallel="0" EstimateRebinds="3" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1007"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1007"></ColumnReference><ScalarOperator ScalarString="[Recr1005]+(1)"><Arithmetic Operation="ADD"><ScalarOperator><Identifier><ColumnReference Column="Recr1005"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="10" PhysicalOp="Filter" LogicalOp="Filter" EstimateRows="1" EstimateIO="0" EstimateCPU="4.8e-07" AvgRowSize="9" EstimatedTotalSubtreeCost="6.077e-06" Parallel="0" EstimateRebinds="3" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList></OutputList><Filter StartupExpression="1"><RelOp NodeId="11" PhysicalOp="Constant Scan" LogicalOp="Constant Scan" EstimateRows="1" EstimateIO="0" EstimateCPU="1.157e-06" AvgRowSize="9" EstimatedTotalSubtreeCost="4.157e-06" Parallel="0" EstimateRebinds="0" EstimateRewinds="3" EstimatedExecutionMode="Row"><OutputList></OutputList><ConstantScan></ConstantScan></RelOp><Predicate><ScalarOperator ScalarString="[Recr1005]&lt;(3)"><Compare CompareOp="LT"><ScalarOperator><Identifier><ColumnReference Column="Recr1005"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(3)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></Filter></RelOp></ComputeScalar></RelOp></NestedLoops></RelOp><Predicate><ScalarOperator ScalarString="CASE WHEN [Expr1012]&gt;(100) THEN (0) ELSE NULL END"><IF><Condition><ScalarOperator><Compare CompareOp="GT"><ScalarOperator><Identifier><ColumnReference Column="Expr1012"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(100)"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Const ConstValue="(0)"></Const></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Else></IF></ScalarOperator></Predicate></Assert></RelOp></Concat></RelOp></Spool></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| WITH r(n, label) AS (SELECT 1, CAST(name AS nvarchar(100)) FROM products UNION ALL SELECT n + 1, label FROM r WHERE n < 3) SELECT n, label FROM r
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| products
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id FROM orders UNION SELECT id FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id FROM orders UNION SELECT id FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id FROM orders UNION SELECT id FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0179308" StatementEstRows="1.5" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x93873B3448D3C47C" QueryPlanHash="0xF73F73820F7DB4A3" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Sort" LogicalOp="Distinct Sort" EstimateRows="1.5" EstimateIO="0.0112613" EstimateCPU="0.000103131" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0179308" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1006"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="1"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Column="Union1006"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="1" PhysicalOp="Concatenation" LogicalOp="Concatenation" EstimateRows="2" EstimateIO="0" EstimateCPU="2e-07" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0065664" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1006"></ColumnReference></OutputList><Concat><DefinedValues><DefinedValue><ColumnReference Column="Union1006"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Concat></RelOp></Sort></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id FROM orders UNION SELECT id FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id AS k, amount FROM orders UNION ALL SELECT order_id, price FROM order_items', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| k
: null
: null
: null
: null
: null
< row
: string
| amount
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id AS k, amount FROM orders UNION ALL SELECT order_id, price FROM order_items
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id AS k, amount FROM orders UNION ALL SELECT order_id, price FROM order_items" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0065664" StatementEstRows="2" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xCC1328EBDFE299E8" QueryPlanHash="0x488136717F76E3BE" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Concatenation" LogicalOp="Concatenation" EstimateRows="2" EstimateIO="0" EstimateCPU="2e-07" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0065664" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1006"></ColumnReference><ColumnReference Column="Union1007"></ColumnReference></OutputList><Concat><DefinedValues><DefinedValue><ColumnReference Column="Union1006"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Column="order_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Column="Union1007"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Column="price"></ColumnReference></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Column="order_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Column="price"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Column="order_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[order_items]" Column="price"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[order_items]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Concat></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id AS k, amount FROM orders UNION ALL SELECT order_id, price FROM order_items
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| order_items
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT name FROM customers INTERSECT SELECT name FROM hr.employees', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT name FROM customers INTERSECT SELECT name FROM hr.employees
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT name FROM customers INTERSECT SELECT name FROM hr.employees" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.017932" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x9FF0A5FCAE83E6D0" QueryPlanHash="0xAD6DF2BD32FBC5B5" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Semi Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="61" EstimatedTotalSubtreeCost="0.017932" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[customers].[name] = [shop].[hr].[employees].[name]"><Compare CompareOp="IS"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Sort" LogicalOp="Distinct Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100061" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="1"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT name FROM customers INTERSECT SELECT name FROM hr.employees
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| hr
: string
| employees
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT name FROM customers EXCEPT SELECT name FROM products', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT name FROM customers EXCEPT SELECT name FROM products
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT name FROM customers EXCEPT SELECT name FROM products" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.017932" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x69CFA293E380CE7B" QueryPlanHash="0x279690B260BF80F3" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Anti Semi Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="61" EstimatedTotalSubtreeCost="0.017932" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[customers].[name] = [shop].[dbo].[products].[name]"><Compare CompareOp="IS"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Sort" LogicalOp="Distinct Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100061" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="1"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[products]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT name FROM customers EXCEPT SELECT name FROM products
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| products
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT name FROM customers UNION SELECT name FROM hr.employees ORDER BY name', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT name FROM customers UNION SELECT name FROM hr.employees ORDER BY name
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT name FROM customers UNION SELECT name FROM hr.employees ORDER BY name" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0179308" StatementEstRows="1.5" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xBF32C6EE7939BCF1" QueryPlanHash="0xCA07CFFA8A0BBDCD" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Sort" LogicalOp="Distinct Sort" EstimateRows="1.5" EstimateIO="0.0112613" EstimateCPU="0.000103181" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0179308" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1006"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="1"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Column="Union1006"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="1" PhysicalOp="Concatenation" LogicalOp="Concatenation" EstimateRows="2" EstimateIO="0" EstimateCPU="2e-07" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0065664" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1006"></ColumnReference></OutputList><Concat><DefinedValues><DefinedValue><ColumnReference Column="Union1006"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Concat></RelOp></Sort></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT name FROM customers UNION SELECT name FROM hr.employees ORDER BY name
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| hr
: string
| employees
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, name FROM customers UNION SELECT id, name FROM products UNION SELECT id, name FROM hr.departments', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: null
: null
: null
: null
: null
< row
: string
| name
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, name FROM customers UNION SELECT id, name FROM products UNION SELECT id, name FROM hr.departments
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, name FROM customers UNION SELECT id, name FROM products UNION SELECT id, name FROM hr.departments" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0212183" StatementEstRows="2" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x5862787A531C99FD" QueryPlanHash="0xBCF54590391D416F" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Sort" LogicalOp="Distinct Sort" EstimateRows="2" EstimateIO="0.0112613" EstimateCPU="0.000107483" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0212183" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1009"></ColumnReference><ColumnReference Column="Union1010"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="1"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Column="Union1009"></ColumnReference></OrderByColumn><OrderByColumn Ascending="1"><ColumnReference Column="Union1010"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="1" PhysicalOp="Concatenation" LogicalOp="Concatenation" EstimateRows="3" EstimateIO="0" EstimateCPU="3e-07" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0098496" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1009"></ColumnReference><ColumnReference Column="Union1010"></ColumnReference></OutputList><Concat><DefinedValues><DefinedValue><ColumnReference Column="Union1009"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Column="Union1010"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Column="name"></ColumnReference></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[products]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="4" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[departments]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Concat></RelOp></Sort></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, name FROM customers UNION SELECT id, name FROM products UNION SELECT id, name FROM hr.departments
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| dbo
: string
| products
< row
: string
| shop
: string
| hr
: string
| departments
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM order_totals', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| customer_id
: string
| shop
: string
| dbo
: string
| orders
: string
| customer_id
: null
< row
: string
| gross
: null
: null
: null
: null
: null
< row
: string
| customer_name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM order_totals
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM order_totals" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243489" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x661B6BBADDE4DBB0" QueryPlanHash="0x86D93C034BA36965" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="48" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177796" AvgRowSize="78" EstimatedTotalSubtreeCost="0.0243489" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="[shop].[dbo].[orders].[amount] as [o].[amount]*(1.1)"><Arithmetic Operation="MULT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1.1)"></Const></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM order_totals
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| order_totals
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, gross FROM order_totals', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| gross
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, gross FROM order_totals
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, gross FROM order_totals" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243409" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x9DC5D977930EBA7C" QueryPlanHash="0xC66988F1B98696E7" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177716" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0243409" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="[shop].[dbo].[orders].[amount] as [o].[amount]*(1.1)"><Arithmetic Operation="MULT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1.1)"></Const></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, gross FROM order_totals
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| order_totals
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT customer_name FROM dbo.order_totals WHERE customer_id = 1', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| customer_name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT customer_name FROM dbo.order_totals WHERE customer_id = 1
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT customer_name FROM dbo.order_totals WHERE customer_id = 1" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657134" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xAEAB442C43FDF714" QueryPlanHash="0xA0366F4C2DC0F6E0" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160" ParameterizedText="(@1 tinyint)SELECT [customer_name] FROM [dbo].[order_totals] WHERE [customer_id]=@1"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="61" EstimatedTotalSubtreeCost="0.00657134" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[customers].[id] as [c].[id]=(1)"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=(1)"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></NestedLoops></RelOp><ParameterList><ColumnReference Column="@1" ParameterDataType="tinyint" ParameterCompiledValue="(1)"></ColumnReference></ParameterList></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT customer_name FROM dbo.order_totals WHERE customer_id = 1
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| order_totals
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM big_orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| customer_name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| gross
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM big_orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM big_orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243494" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x7182EE4B902CA736" QueryPlanHash="0x86D93C034BA36965" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="48" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177796" AvgRowSize="74" EstimatedTotalSubtreeCost="0.0243494" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="[shop].[dbo].[orders].[amount] as [o].[amount]*(1.1)"><Arithmetic Operation="MULT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1.1)"></Const></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[amount] as [o].[amount]*(1.1)&gt;(100.000)"><Compare CompareOp="GT"><ScalarOperator><Arithmetic Operation="MULT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1.1)"></Const></ScalarOperator></Arithmetic></ScalarOperator><ScalarOperator><Const ConstValue="(100.000)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></ComputeScalar></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM big_orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| big_orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT b.customer_name, o.status FROM big_orders b JOIN orders o ON o.id = b.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| customer_name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| status
: string
| shop
: string
| dbo
: string
| orders
: string
| status
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT b.customer_name, o.status FROM big_orders b JOIN orders o ON o.id = b.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT b.customer_name, o.status FROM big_orders b JOIN orders o ON o.id = b.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0454071" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x2374A23B9316C9A4" QueryPlanHash="0xD56671993E70E2F4" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="72" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="2048" SerialDesiredMemory="2112" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177769" AvgRowSize="68" EstimatedTotalSubtreeCost="0.0454071" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference></OutputList><MemoryFractions Input="0.5" Output="0.5"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[id] as [o].[id]=[shop].[dbo].[orders].[id] as [o].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="status"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" TableReferenceId="2" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177742" AvgRowSize="65" EstimatedTotalSubtreeCost="0.024344" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="0.5" Output="0.5"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="4" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" TableReferenceId="1" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[amount] as [o].[amount]*(1.1)&gt;(100.000)"><Compare CompareOp="GT"><ScalarOperator><Arithmetic Operation="MULT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1.1)"></Const></ScalarOperator></Arithmetic></ScalarOperator><ScalarOperator><Const ConstValue="(100.000)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp><RelOp NodeId="5" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT b.customer_name, o.status FROM big_orders b JOIN orders o ON o.id = b.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| big_orders
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM hr.staff', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| staff_id
: string
| shop
: string
| hr
: string
| employees
: string
| id
: null
< row
: string
| staff_name
: string
| shop
: string
| hr
: string
| employees
: string
| name
: null
< row
: string
| dept
: string
| shop
: string
| hr
: string
| departments
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM hr.staff
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM hr.staff" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657068" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x2CB698CEC2D65CD5" QueryPlanHash="0xB2699FAC2B99D17A" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Outer Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="117" EstimatedTotalSubtreeCost="0.00657068" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[hr].[departments].[id] as [d].[id]=[shop].[hr].[employees].[dept_id] as [e].[dept_id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM hr.staff
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| hr
: string
| staff
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT staff_name, dept FROM hr.staff WHERE staff_id = 1', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| staff_name
: string
| shop
: string
| hr
: string
| employees
: string
| name
: null
< row
: string
| dept
: string
| shop
: string
| hr
: string
| departments
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT staff_name, dept FROM hr.staff WHERE staff_id = 1
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT staff_name, dept FROM hr.staff WHERE staff_id = 1" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.00657116" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x168DF45541DA0301" QueryPlanHash="0xB2699FAC2B99D17A" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160" ParameterizedText="(@1 tinyint)SELECT [staff_name],[dept] FROM [hr].[staff] WHERE [staff_id]=@1"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Outer Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="113" EstimatedTotalSubtreeCost="0.00657116" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[hr].[departments].[id] as [d].[id]=[shop].[hr].[employees].[dept_id] as [e].[dept_id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[hr].[employees].[id] as [e].[id]=(1)"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp><ParameterList><ColumnReference Column="@1" ParameterDataType="tinyint" ParameterCompiledValue="(1)"></ColumnReference></ParameterList></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT staff_name, dept FROM hr.staff WHERE staff_id = 1
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| hr
: string
| staff
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM all_people', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| all_people
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM all_people
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM all_people" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0065664" StatementEstRows="2" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x3378DA36966ABB3E" QueryPlanHash="0x4DA4F8F49B0EA30E" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Concatenation" LogicalOp="Concatenation" EstimateRows="2" EstimateIO="0" EstimateCPU="2e-07" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0065664" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Union1004"></ColumnReference></OutputList><Concat><DefinedValues><DefinedValue><ColumnReference Column="Union1004"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Concat></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM all_people
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| all_people
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM v_products', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| products
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| products
: string
| name
: null
< row
: string
| category
: string
| shop
: string
| dbo
: string
| products
: string
| category
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM v_products
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM v_products" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x991C57B6D5581E06" QueryPlanHash="0xFA4818BE39808759" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="77" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="category"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[products]" Column="category"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[products]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM v_products
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| v_products
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM archived', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| archive
: string
| dbo
: string
| old_orders
: string
| id
: null
< row
: string
| amount
: string
| archive
: string
| dbo
: string
| old_orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM archived
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM archived" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x27A7FCF7FD93D04C" QueryPlanHash="0x3EECC6FA094179FC" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM archived
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| archived
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT customer_name FROM order_totals WHERE id IN (SELECT id FROM big_orders)', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| customer_name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT customer_name FROM order_totals WHERE id IN (SELECT id FROM big_orders)
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT customer_name FROM order_totals WHERE id IN (SELECT id FROM big_orders)" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0245089" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x3D68C685445748F9" QueryPlanHash="0xB5CAC1236852BC8C" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="1" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0245089" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><OuterReferences><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></OuterReferences><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" TableReferenceId="1" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="3" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-06" AvgRowSize="9" EstimatedTotalSubtreeCost="0.0212216" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList></OutputList><StreamAggregate><DefinedValues></DefinedValues><GroupBy><ColumnReference Column="Bmk1000"></ColumnReference></GroupBy><RelOp NodeId="4" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0212206" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Bmk1000"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="5" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="19" EstimatedTotalSubtreeCost="0.017933" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Bmk1000"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[id] as [o].[id]=[shop].[dbo].[orders].[id] as [o].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="6" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100019" AvgRowSize="19" EstimatedTotalSubtreeCost="0.0146449" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Bmk1000"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Column="Bmk1000"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="7" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="23" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Bmk1000"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Column="Bmk1000"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" TableReferenceId="1" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></Sort></RelOp><RelOp NodeId="9" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" TableReferenceId="2" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[orders].[amount] as [o].[amount]*(1.1)&gt;(100.000)"><Compare CompareOp="GT"><ScalarOperator><Arithmetic Operation="MULT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1.1)"></Const></ScalarOperator></Arithmetic></ScalarOperator><ScalarOperator><Const ConstValue="(100.000)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></NestedLoops></RelOp><RelOp NodeId="10" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" TableReferenceId="2" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></StreamAggregate></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT customer_name FROM order_totals WHERE id IN (SELECT id FROM big_orders)
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| big_orders
< row
: string
| shop
: string
| dbo
: string
| order_totals
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM [Vip Names]', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| VipName
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM [Vip Names]
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM [Vip Names]" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x808A948943E58435" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="68" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[customers].[region] as [c].[region]=&apos;vip&apos;"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="&apos;vip&apos;"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM [Vip Names]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| Vip Names
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT VipName FROM dbo.[Vip Names]', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| VipName
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT VipName FROM dbo.[Vip Names]
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT VipName FROM dbo.[Vip Names]" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x808A948943E58435" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="68" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[shop].[dbo].[customers].[region] as [c].[region]=&apos;vip&apos;"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="region"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="&apos;vip&apos;"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT VipName FROM dbo.[Vip Names]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| Vip Names
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT t.id, s.staff_name FROM order_totals t JOIN hr.staff s ON s.staff_id = t.customer_id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| staff_name
: string
| shop
: string
| hr
: string
| employees
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT t.id, s.staff_name FROM order_totals t JOIN hr.staff s ON s.staff_id = t.customer_id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT t.id, s.staff_name FROM order_totals t JOIN hr.staff s ON s.staff_id = t.customer_id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0486887" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x5AB97B7CCEC141A6" QueryPlanHash="0xC28306CFD07D7D96" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="72" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="2048" SerialDesiredMemory="2112" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Left Outer Join" EstimateRows="1" EstimateIO="0" EstimateCPU="4.18e-06" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0486887" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="[shop].[hr].[departments].[id] as [d].[id]=[shop].[hr].[employees].[dept_id] as [e].[dept_id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></Predicate><RelOp NodeId="1" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177716" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0454011" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></OutputList><MemoryFractions Input="0.5" Output="0.5"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[customers].[id] as [c].[id]=[shop].[hr].[employees].[id] as [e].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="3" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177742" AvgRowSize="73" EstimatedTotalSubtreeCost="0.0243434" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></OutputList><MemoryFractions Input="0.5" Output="0.5"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[hr].[employees].[id] as [e].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="4" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="5" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" Column="dept_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" Alias="[e]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></Hash></RelOp><RelOp NodeId="6" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[departments]" Alias="[d]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT t.id, s.staff_name FROM order_totals t JOIN hr.staff s ON s.staff_id = t.customer_id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| order_totals
< row
: string
| shop
: string
| hr
: string
| staff
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM archive.dbo.old_orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| archive
: string
| dbo
: string
| old_orders
: string
| id
: null
< row
: string
| customer_id
: string
| archive
: string
| dbo
: string
| old_orders
: string
| customer_id
: null
< row
: string
| amount
: string
| archive
: string
| dbo
: string
| old_orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM archive.dbo.old_orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM archive.dbo.old_orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x9CC706A495BFAA9E" QueryPlanHash="0x3EECC6FA094179FC" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="id"></ColumnReference><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[archive]" Schema="[dbo]" Table="[old_orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM archive.dbo.old_orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| archive
: string
| dbo
: string
| old_orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM archive.dbo.recent', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| archive
: string
| dbo
: string
| old_orders
: string
| id
: null
< row
: string
| amount
: string
| archive
: string
| dbo
: string
| old_orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM archive.dbo.recent
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM archive.dbo.recent" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xD0C0F7F9E8347435" QueryPlanHash="0x3EECC6FA094179FC" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[archive].[dbo].[old_orders].[id] as [o].[id]&gt;(100)"><Compare CompareOp="GT"><ScalarOperator><Identifier><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(100)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM archive.dbo.recent
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| archive
: string
| dbo
: string
| recent
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT c.name, a.amount FROM customers c JOIN archive.dbo.old_orders a ON a.customer_id = c.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
< row
: string
| amount
: string
| archive
: string
| dbo
: string
| old_orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT c.name, a.amount FROM customers c JOIN archive.dbo.old_orders a ON a.customer_id = c.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT c.name, a.amount FROM customers c JOIN archive.dbo.old_orders a ON a.customer_id = c.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243461" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xEF55DA84FBF8835D" QueryPlanHash="0x176F48CF8ABC2E33" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177769" AvgRowSize="70" EstimatedTotalSubtreeCost="0.0243461" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[a]" Column="amount"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[a]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[customers].[id] as [c].[id]=[archive].[dbo].[old_orders].[customer_id] as [a].[customer_id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[a]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[a]" Column="customer_id"></ColumnReference><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[a]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[a]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[a]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[a]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT c.name, a.amount FROM customers c JOIN archive.dbo.old_orders a ON a.customer_id = c.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| archive
: string
| dbo
: string
| old_orders
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT [Order ID], [Unit Price] * Qty AS total FROM [Order Details]', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| Order ID
: string
| shop
: string
| dbo
: string
| Order Details
: string
| Order ID
: null
< row
: string
| total
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT [Order ID], [Unit Price] * Qty AS total FROM [Order Details]
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT [Order ID], [Unit Price] * Qty AS total FROM [Order Details]" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032832" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x91A6EED228C236C2" QueryPlanHash="0x57E41CE59B9F45F9" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="19" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Order ID"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="[shop].[dbo].[Order Details].[Unit Price]*CONVERT_IMPLICIT(money,[shop].[dbo].[Order Details].[Qty],0)"><Arithmetic Operation="MULT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Unit Price"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Convert DataType="money" Style="0" Implicit="1"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Qty"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="23" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Order ID"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Unit Price"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Qty"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Order ID"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Unit Price"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Qty"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[Order Details]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT [Order ID], [Unit Price] * Qty AS total FROM [Order Details]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| Order Details
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT "Order ID", Qty FROM "Order Details"', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| Order ID
: string
| shop
: string
| dbo
: string
| Order Details
: string
| Order ID
: null
< row
: string
| Qty
: string
| shop
: string
| dbo
: string
| Order Details
: string
| Qty
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT "Order ID", Qty FROM "Order Details"
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT &quot;Order ID&quot;, Qty FROM &quot;Order Details&quot;" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xA5142F893E7F8518" QueryPlanHash="0x437CCDA24B3431E5" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Order ID"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Qty"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Order ID"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Qty"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[Order Details]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT "Order ID", Qty FROM "Order Details"
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| Order Details
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT d.[Order ID] FROM dbo.[Order Details] AS d', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| Order ID
: string
| shop
: string
| dbo
: string
| Order Details
: string
| Order ID
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT d.[Order ID] FROM dbo.[Order Details] AS d
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT d.[Order ID] FROM dbo.[Order Details] AS d" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x0EE4B07058A678E1" QueryPlanHash="0x437CCDA24B3431E5" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Alias="[d]" Column="Order ID"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Alias="[d]" Column="Order ID"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[Order Details]" Alias="[d]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT d.[Order ID] FROM dbo.[Order Details] AS d
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| Order Details
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM [Order Details]', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| Order ID
: string
| shop
: string
| dbo
: string
| Order Details
: string
| Order ID
: null
< row
: string
| Unit Price
: string
| shop
: string
| dbo
: string
| Order Details
: string
| Unit Price
: null
< row
: string
| Qty
: string
| shop
: string
| dbo
: string
| Order Details
: string
| Qty
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM [Order Details]
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM [Order Details]" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xE4B1F7329D54F81F" QueryPlanHash="0x437CCDA24B3431E5" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="23" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Order ID"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Unit Price"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Qty"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Order ID"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Unit Price"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[Order Details]" Column="Qty"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[Order Details]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM [Order Details]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| Order Details
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, camelname FROM mixedcase', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| MixedCase
: string
| Id
: null
< row
: string
| camelname
: string
| shop
: string
| dbo
: string
| MixedCase
: string
| CamelName
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, camelname FROM mixedcase
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, camelname FROM mixedcase" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xEDE26DA71112E5F5" QueryPlanHash="0xCEF3A912A60025CC" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="35" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="Id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="CamelName"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="Id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="CamelName"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[MixedCase]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, camelname FROM mixedcase
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| MixedCase
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT ID, NAME FROM CUSTOMERS', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| ID
: string
| shop
: string
| dbo
: string
| customers
: string
| id
: null
< row
: string
| NAME
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT ID, NAME FROM CUSTOMERS
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT ID, NAME FROM CUSTOMERS" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xF7B23518122D6AF7" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT ID, NAME FROM CUSTOMERS
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT [Id], [CamelName] FROM [dbo].[MixedCase]', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| Id
: string
| shop
: string
| dbo
: string
| MixedCase
: string
| Id
: null
< row
: string
| CamelName
: string
| shop
: string
| dbo
: string
| MixedCase
: string
| CamelName
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT [Id], [CamelName] FROM [dbo].[MixedCase]
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT [Id], [CamelName] FROM [dbo].[MixedCase]" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xEDE26DA71112E5F5" QueryPlanHash="0xCEF3A912A60025CC" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="35" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="Id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="CamelName"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="Id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="CamelName"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[MixedCase]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT [Id], [CamelName] FROM [dbo].[MixedCase]
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| MixedCase
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT c.NAME AS [Customer Name] FROM Customers AS C', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| Customer Name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT c.NAME AS [Customer Name] FROM Customers AS C
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT c.NAME AS [Customer Name] FROM Customers AS C" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x4E88BEC372FB0DD7" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[C]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[C]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[C]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT c.NAME AS [Customer Name] FROM Customers AS C
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM MIXEDCASE', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| Id
: string
| shop
: string
| dbo
: string
| MixedCase
: string
| Id
: null
< row
: string
| CamelName
: string
| shop
: string
| dbo
: string
| MixedCase
: string
| CamelName
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM MIXEDCASE
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM MIXEDCASE" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xEDE26DA71112E5F5" QueryPlanHash="0xCEF3A912A60025CC" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="35" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="Id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="CamelName"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="Id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[MixedCase]" Column="CamelName"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[MixedCase]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM MIXEDCASE
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| MixedCase
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT customers.id FROM orders AS customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT customers.id FROM orders AS customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT customers.id FROM orders AS customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xD919D631C3C6573F" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[customers]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[customers]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT customers.id FROM orders AS customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT orders.name FROM customers AS orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT orders.name FROM customers AS orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT orders.name FROM customers AS orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x5BB4DD9DD27951FB" QueryPlanHash="0xEB583F9F57ED400D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[orders]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[orders]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT orders.name FROM customers AS orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT customers.nickname FROM hr.customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| nickname
: string
| shop
: string
| hr
: string
| customers
: string
| nickname
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT customers.nickname FROM hr.customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT customers.nickname FROM hr.customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xDECD6BDEA4D3086F" QueryPlanHash="0x2C5655194DF7A24F" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="31" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Column="nickname"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Column="nickname"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT customers.nickname FROM hr.customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| hr
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT c.id, d.name FROM hr.customers c JOIN dbo.customers d ON d.id = c.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| hr
: string
| customers
: string
| id
: null
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT c.id, d.name FROM hr.customers c JOIN dbo.customers d ON d.id = c.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT c.id, d.name FROM hr.customers c JOIN dbo.customers d ON d.id = c.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243408" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xF8FA1A2ED3842EEE" QueryPlanHash="0x4895898DDAC01E49" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177716" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0243408" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[d]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[d]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[hr].[customers].[id] as [c].[id]=[shop].[dbo].[customers].[id] as [d].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[d]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[d]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[d]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[d]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[d]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[d]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT c.id, d.name FROM hr.customers c JOIN dbo.customers d ON d.id = c.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| hr
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT customers.name FROM hr.customers AS h JOIN customers ON customers.id = h.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT customers.name FROM hr.customers AS h JOIN customers ON customers.id = h.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT customers.name FROM hr.customers AS h JOIN customers ON customers.id = h.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243408" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x085F58ACD3E1CEC3" QueryPlanHash="0x4895898DDAC01E49" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177716" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0243408" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" Column="id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[hr].[customers].[id] as [h].[id]=[shop].[dbo].[customers].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[customers]" Alias="[h]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT customers.name FROM hr.customers AS h JOIN customers ON customers.id = h.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
< row
: string
| shop
: string
| hr
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, ROW_NUMBER() OVER (ORDER BY amount) AS rn FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| rn
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, ROW_NUMBER() OVER (ORDER BY amount) AS rn FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, ROW_NUMBER() OVER (ORDER BY amount) AS rn FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146445" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x0C0F0DB5C6C4714D" QueryPlanHash="0x900400B57957A996" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Sequence Project" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="8e-08" AvgRowSize="19" EstimatedTotalSubtreeCost="0.0146445" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><SequenceProject><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="row_number"><Sequence FunctionName="row_number"></Sequence></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Segment" LogicalOp="Segment" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="19" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Column="Segment1004"></ColumnReference></OutputList><Segment><GroupBy></GroupBy><SegmentColumn><ColumnReference Column="Segment1004"></ColumnReference></SegmentColumn><RelOp NodeId="2" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.00010002" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></Segment></RelOp></SequenceProject></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, ROW_NUMBER() OVER (ORDER BY amount) AS rn FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, SUM(amount) OVER (PARTITION BY customer_id) AS total FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| total
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, SUM(amount) OVER (PARTITION BY customer_id) AS total FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, SUM(amount) OVER (PARTITION BY customer_id) AS total FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146724" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xA1AECD507275D696" QueryPlanHash="0x4F969FD0A836F1A3" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="2.156e-05" AvgRowSize="28" EstimatedTotalSubtreeCost="0.0146724" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><NestedLoops Optimized="0"><RelOp NodeId="1" PhysicalOp="Table Spool" LogicalOp="Lazy Spool" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="24" EstimatedTotalSubtreeCost="0.01465" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><Spool><RelOp NodeId="2" PhysicalOp="Segment" LogicalOp="Segment" EstimateRows="1" EstimateIO="0" EstimateCPU="2.8e-06" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0146472" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Column="Segment1004"></ColumnReference></OutputList><Segment><GroupBy><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></GroupBy><SegmentColumn><ColumnReference Column="Segment1004"></ColumnReference></SegmentColumn><RelOp NodeId="3" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100024" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="4" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></Segment></RelOp></Spool></RelOp><RelOp NodeId="5" PhysicalOp="Nested Loops" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="2.8e-07" AvgRowSize="24" EstimatedTotalSubtreeCost="5.6e-07" Parallel="0" EstimateRebinds="1" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><NestedLoops Optimized="0"><Predicate><ScalarOperator ScalarString="(1)"><Const ConstValue="(1)"></Const></ScalarOperator></Predicate><RelOp NodeId="6" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="2.8e-08" AvgRowSize="24" EstimatedTotalSubtreeCost="3.08e-07" Parallel="0" EstimateRebinds="1" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CASE WHEN [Expr1005]=(0) THEN NULL ELSE [Expr1006] END"><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Column="Expr1005"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(0)"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Then><Else><ScalarOperator><Identifier><ColumnReference Column="Expr1006"></ColumnReference></Identifier></ScalarOperator></Else></IF></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="7" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="2.8e-07" AvgRowSize="24" EstimatedTotalSubtreeCost="2.8e-07" Parallel="0" EstimateRebinds="1" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1005"></ColumnReference><ColumnReference Column="Expr1006"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1005"></ColumnReference><ScalarOperator ScalarString="COUNT_BIG([shop].[dbo].[orders].[amount])"><Aggregate Distinct="0" AggType="COUNT_BIG"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1006"></ColumnReference><ScalarOperator ScalarString="SUM([shop].[dbo].[orders].[amount])"><Aggregate Distinct="0" AggType="SUM"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="8" PhysicalOp="Table Spool" LogicalOp="Lazy Spool" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="24" EstimatedTotalSubtreeCost="0" Parallel="0" EstimateRebinds="1" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><Spool PrimaryNodeId="1"></Spool></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp><RelOp NodeId="9" PhysicalOp="Table Spool" LogicalOp="Lazy Spool" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="24" EstimatedTotalSubtreeCost="0" Parallel="0" EstimateRebinds="1" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><Spool PrimaryNodeId="1"></Spool></RelOp></NestedLoops></RelOp></NestedLoops></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, SUM(amount) OVER (PARTITION BY customer_id) AS total FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, LAG(amount) OVER (ORDER BY created) AS prev FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| prev
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, LAG(amount) OVER (ORDER BY created) AS prev FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, LAG(amount) OVER (ORDER BY created) AS prev FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146477" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x1BCDC17D49077301" QueryPlanHash="0x5662E8D2BC651B49" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="6.28e-07" AvgRowSize="60" EstimatedTotalSubtreeCost="0.0146477" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="LAST_VALUE([shop].[dbo].[orders].[amount])"><Aggregate Distinct="0" AggType="LAST_VALUE"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ScalarOperator ScalarString="ANY([shop].[dbo].[orders].[id])"><Aggregate Distinct="0" AggType="ANY"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><GroupBy><ColumnReference Column="WindowCount1009"></ColumnReference></GroupBy><RelOp NodeId="1" PhysicalOp="Window Spool" LogicalOp="Window Spool" EstimateRows="2" EstimateIO="0" EstimateCPU="2.198e-06" AvgRowSize="60" EstimatedTotalSubtreeCost="0.0146471" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference><ColumnReference Column="RowNumber1004"></ColumnReference><ColumnReference Column="WindowCount1009"></ColumnReference><ColumnReference Column="Segment1008"></ColumnReference></OutputList><WindowSpool><RelOp NodeId="2" PhysicalOp="Segment" LogicalOp="Segment" EstimateRows="1" EstimateIO="0" EstimateCPU="6.28e-08" AvgRowSize="60" EstimatedTotalSubtreeCost="0.0146447" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference><ColumnReference Column="RowNumber1004"></ColumnReference><ColumnReference Column="TopRowNumber1005"></ColumnReference><ColumnReference Column="BottomRowNumber1006"></ColumnReference><ColumnReference Column="Segment1008"></ColumnReference></OutputList><Segment><GroupBy></GroupBy><SegmentColumn><ColumnReference Column="Segment1008"></ColumnReference></SegmentColumn><RelOp NodeId="3" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="52" EstimatedTotalSubtreeCost="0.0146446" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference><ColumnReference Column="RowNumber1004"></ColumnReference><ColumnReference Column="TopRowNumber1005"></ColumnReference><ColumnReference Column="BottomRowNumber1006"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="TopRowNumber1005"></ColumnReference><ScalarOperator ScalarString="[RowNumber1004]-(1)"><Arithmetic Operation="SUB"><ScalarOperator><Identifier><ColumnReference Column="RowNumber1004"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="BottomRowNumber1006"></ColumnReference><ScalarOperator ScalarString="[RowNumber1004]-(1)"><Arithmetic Operation="SUB"><ScalarOperator><Identifier><ColumnReference Column="RowNumber1004"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="4" PhysicalOp="Sequence Project" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="8e-08" AvgRowSize="36" EstimatedTotalSubtreeCost="0.0146445" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference><ColumnReference Column="RowNumber1004"></ColumnReference></OutputList><SequenceProject><DefinedValues><DefinedValue><ColumnReference Column="RowNumber1004"></ColumnReference><ScalarOperator ScalarString="row_number"><Sequence FunctionName="row_number"></Sequence></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="5" PhysicalOp="Segment" LogicalOp="Segment" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="36" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference><ColumnReference Column="Segment1007"></ColumnReference></OutputList><Segment><GroupBy></GroupBy><SegmentColumn><ColumnReference Column="Segment1007"></ColumnReference></SegmentColumn><RelOp NodeId="6" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100028" AvgRowSize="28" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="7" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="28" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></Segment></RelOp></SequenceProject></RelOp></ComputeScalar></RelOp></Segment></RelOp></WindowSpool></RelOp></StreamAggregate></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, LAG(amount) OVER (ORDER BY created) AS prev FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT name, RANK() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS r FROM hr.employees', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| hr
: string
| employees
: string
| name
: null
< row
: string
| r
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT name, RANK() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS r FROM hr.employees
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT name, RANK() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS r FROM hr.employees" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146445" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xB39F587311DD63F2" QueryPlanHash="0xD867F0FEEFD6432F" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="32" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Sequence Project" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="6e-08" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0146445" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><SequenceProject><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="rank"><Sequence FunctionName="rank"></Sequence></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Segment" LogicalOp="Segment" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0146445" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="salary"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference><ColumnReference Column="Segment1004"></ColumnReference><ColumnReference Column="Segment1005"></ColumnReference></OutputList><Segment><GroupBy><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="salary"></ColumnReference></GroupBy><SegmentColumn><ColumnReference Column="Segment1005"></ColumnReference></SegmentColumn><RelOp NodeId="2" PhysicalOp="Segment" LogicalOp="Segment" EstimateRows="1" EstimateIO="0" EstimateCPU="2e-08" AvgRowSize="69" EstimatedTotalSubtreeCost="0.0146445" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="salary"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference><ColumnReference Column="Segment1004"></ColumnReference></OutputList><Segment><GroupBy><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference></GroupBy><SegmentColumn><ColumnReference Column="Segment1004"></ColumnReference></SegmentColumn><RelOp NodeId="3" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100074" AvgRowSize="74" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="salary"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference></OrderByColumn><OrderByColumn Ascending="0"><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="salary"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="4" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="74" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="salary"></ColumnReference><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="salary"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[hr]" Table="[employees]" Column="dept_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[hr]" Table="[employees]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></Segment></RelOp></Segment></RelOp></SequenceProject></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT name, RANK() OVER (PARTITION BY dept_id ORDER BY salary DESC) AS r FROM hr.employees
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| hr
: string
| employees
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, FIRST_VALUE(status) OVER (PARTITION BY customer_id ORDER BY created) AS first_status FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| first_status
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, FIRST_VALUE(status) OVER (PARTITION BY customer_id ORDER BY created) AS first_status FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, FIRST_VALUE(status) OVER (PARTITION BY customer_id ORDER BY created) AS first_status FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146476" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xA0D2322B3376F314" QueryPlanHash="0x733FB11C3E95D30F" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0146476" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CASE WHEN IsFalseOrNull [Expr1004] THEN [Expr1005] ELSE NULL END"><IF><Condition><ScalarOperator><Logical Operation="IsFalseOrNull"><ScalarOperator><Identifier><ColumnReference Column="Expr1004"></ColumnReference></Identifier></ScalarOperator></Logical></ScalarOperator></Condition><Then><ScalarOperator><Identifier><ColumnReference Column="Expr1005"></ColumnReference></Identifier></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Else></IF></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="6.28e-07" AvgRowSize="48" EstimatedTotalSubtreeCost="0.0146475" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference><ColumnReference Column="Expr1005"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="FIRST_VALUE_IS_NULL([shop].[dbo].[orders].[status])"><Aggregate Distinct="0" AggType="FIRST_VALUE_IS_NULL"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1005"></ColumnReference><ScalarOperator ScalarString="FIRST_VALUE_IGNORE_NULLS([shop].[dbo].[orders].[status])"><Aggregate Distinct="0" AggType="FIRST_VALUE_IGNORE_NULLS"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ScalarOperator ScalarString="ANY([shop].[dbo].[orders].[id])"><Aggregate Distinct="0" AggType="ANY"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><GroupBy><ColumnReference Column="WindowCount1008"></ColumnReference></GroupBy><RelOp NodeId="2" PhysicalOp="Window Spool" LogicalOp="Window Spool" EstimateRows="2" EstimateIO="0" EstimateCPU="2.198e-06" AvgRowSize="48" EstimatedTotalSubtreeCost="0.0146469" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference><ColumnReference Column="WindowCount1008"></ColumnReference><ColumnReference Column="Segment1006"></ColumnReference><ColumnReference Column="Segment1007"></ColumnReference></OutputList><WindowSpool><RelOp NodeId="3" PhysicalOp="Segment" LogicalOp="Segment" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="48" EstimatedTotalSubtreeCost="0.0146445" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference><ColumnReference Column="Segment1006"></ColumnReference><ColumnReference Column="Segment1007"></ColumnReference></OutputList><Segment><GroupBy><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></GroupBy><SegmentColumn><ColumnReference Column="Segment1007"></ColumnReference></SegmentColumn><RelOp NodeId="4" PhysicalOp="Segment" LogicalOp="Segment" EstimateRows="1" EstimateIO="0" EstimateCPU="6.28e-08" AvgRowSize="40" EstimatedTotalSubtreeCost="0.0146445" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference><ColumnReference Column="Segment1006"></ColumnReference></OutputList><Segment><GroupBy><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></GroupBy><SegmentColumn><ColumnReference Column="Segment1006"></ColumnReference></SegmentColumn><RelOp NodeId="5" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100032" AvgRowSize="32" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OrderByColumn><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="6" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="32" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="created"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></Segment></RelOp></Segment></RelOp></WindowSpool></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, FIRST_VALUE(status) OVER (PARTITION BY customer_id ORDER BY created) AS first_status FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT COUNT(*) FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: null
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT COUNT(*) FROM orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT COUNT(*) FROM orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032842" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x99D3EC46E287D3C6" QueryPlanHash="0xBA2D65D212AD893D" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032842" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CONVERT_IMPLICIT(int,[Expr1004],0)"><Convert DataType="int" Style="0" Implicit="1"><ScalarOperator><Identifier><ColumnReference Column="Expr1004"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032842" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1004"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="Count(*)"><Aggregate AggType="countstar" Distinct="0"></Aggregate></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="9" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT COUNT(*) FROM orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT customer_id, SUM(amount) AS total FROM orders GROUP BY customer_id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| customer_id
: string
| shop
: string
| dbo
: string
| orders
: string
| customer_id
: null
< row
: string
| total
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT customer_id, SUM(amount) AS total FROM orders GROUP BY customer_id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT customer_id, SUM(amount) AS total FROM orders GROUP BY customer_id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146455" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x103CB9A2D037097E" QueryPlanHash="0x3E1102BA5941C929" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="28" EstimatedTotalSubtreeCost="0.0146455" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CASE WHEN [Expr1008]=(0) THEN NULL ELSE [Expr1009] END"><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Column="Expr1008"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(0)"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Then><Else><ScalarOperator><Identifier><ColumnReference Column="Expr1009"></ColumnReference></Identifier></ScalarOperator></Else></IF></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="28" EstimatedTotalSubtreeCost="0.0146455" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1008"></ColumnReference><ColumnReference Column="Expr1009"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1008"></ColumnReference><ScalarOperator ScalarString="COUNT_BIG([shop].[dbo].[orders].[amount])"><Aggregate Distinct="0" AggType="COUNT_BIG"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1009"></ColumnReference><ScalarOperator ScalarString="SUM([shop].[dbo].[orders].[amount])"><Aggregate Distinct="0" AggType="SUM"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><GroupBy><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></GroupBy><RelOp NodeId="2" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.00010002" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT customer_id, SUM(amount) AS total FROM orders GROUP BY customer_id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT customer_id, COUNT(*) FROM orders GROUP BY customer_id HAVING COUNT(*) > 1', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| customer_id
: string
| shop
: string
| dbo
: string
| orders
: string
| customer_id
: null
< row
: null
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT customer_id, COUNT(*) FROM orders GROUP BY customer_id HAVING COUNT(*) > 1
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT customer_id, COUNT(*) FROM orders GROUP BY customer_id HAVING COUNT(*) &gt; 1" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.014646" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xAAA5939E98D1006C" QueryPlanHash="0x1DDA93DE2126A0F1" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Filter" LogicalOp="Filter" EstimateRows="1" EstimateIO="0" EstimateCPU="4.8e-07" AvgRowSize="15" EstimatedTotalSubtreeCost="0.014646" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><Filter StartupExpression="0"><RelOp NodeId="1" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0146455" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CONVERT_IMPLICIT(int,[Expr1006],0)"><Convert DataType="int" Style="0" Implicit="1"><ScalarOperator><Identifier><ColumnReference Column="Expr1006"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="15" EstimatedTotalSubtreeCost="0.0146455" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1006"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1006"></ColumnReference><ScalarOperator ScalarString="Count(*)"><Aggregate AggType="countstar" Distinct="0"></Aggregate></ScalarOperator></DefinedValue></DefinedValues><GroupBy><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></GroupBy><RelOp NodeId="3" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100011" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="4" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp><Predicate><ScalarOperator ScalarString="[Expr1003]&gt;(1)"><Compare CompareOp="GT"><ScalarOperator><Identifier><ColumnReference Column="Expr1003"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></Filter></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT customer_id, COUNT(*) FROM orders GROUP BY customer_id HAVING COUNT(*) > 1
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT MAX(name) AS m FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| m
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT MAX(name) AS m FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT MAX(name) AS m FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032842" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xCBF80F0E9BCA0E12" QueryPlanHash="0xD424E6084E652DED" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032842" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="MAX([shop].[dbo].[customers].[name])"><Aggregate Distinct="0" AggType="MAX"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></StreamAggregate></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT MAX(name) AS m FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT region, COUNT(DISTINCT id) AS n FROM customers GROUP BY region', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| region
: string
| shop
: string
| dbo
: string
| customers
: string
| region
: null
< row
: string
| n
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT region, COUNT(DISTINCT id) AS n FROM customers GROUP BY region
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT region, COUNT(DISTINCT id) AS n FROM customers GROUP BY region" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146455" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x76296B44D2465DBC" QueryPlanHash="0x73209F6D9B4926BA" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0146455" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CONVERT_IMPLICIT(int,[Expr1006],0)"><Convert DataType="int" Style="0" Implicit="1"><ScalarOperator><Identifier><ColumnReference Column="Expr1006"></ColumnReference></Identifier></ScalarOperator></Convert></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0146455" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference><ColumnReference Column="Expr1006"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1006"></ColumnReference><ScalarOperator ScalarString="COUNT([shop].[dbo].[customers].[id])"><Aggregate Distinct="0" AggType="COUNT_BIG"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><GroupBy><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></GroupBy><RelOp NodeId="2" PhysicalOp="Sort" LogicalOp="Distinct Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.00010002" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="1"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OrderByColumn><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT region, COUNT(DISTINCT id) AS n FROM customers GROUP BY region
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT STRING_AGG(name, '','') AS names FROM customers', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| names
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT STRING_AGG(name, ',') AS names FROM customers
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT STRING_AGG(name, &apos;,&apos;) AS names FROM customers" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032842" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x2EB0290917F29155" QueryPlanHash="0xD424E6084E652DED" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="24" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="4011" EstimatedTotalSubtreeCost="0.0032842" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Column="Expr1003"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="STRING_AGG([shop].[dbo].[customers].[name],N&apos;,&apos;)"><Aggregate Distinct="0" AggType="STRING_AGG"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="N&apos;,&apos;"></Const></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></StreamAggregate></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT STRING_AGG(name, ',') AS names FROM customers
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT region, SUM(LEN(name)) AS l FROM customers GROUP BY ROLLUP (region)', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| region
: string
| shop
: string
| dbo
: string
| customers
: string
| region
: null
< row
: string
| l
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT region, SUM(LEN(name)) AS l FROM customers GROUP BY ROLLUP (region)
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT region, SUM(LEN(name)) AS l FROM customers GROUP BY ROLLUP (region)" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146473" StatementEstRows="2" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xBED22567BF5FF0C3" QueryPlanHash="0x74CDEB8BEDBE82F9" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="32" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="1" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="2" EstimateIO="0" EstimateCPU="2e-07" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0146473" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CASE WHEN [Expr1004]=(0) THEN NULL ELSE [Expr1005] END"><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Column="Expr1004"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(0)"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Then><Else><ScalarOperator><Identifier><ColumnReference Column="Expr1005"></ColumnReference></Identifier></ScalarOperator></Else></IF></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="2" EstimateIO="0" EstimateCPU="1.5e-06" AvgRowSize="28" EstimatedTotalSubtreeCost="0.0146471" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference><ColumnReference Column="Expr1005"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="SUM([Expr1006])"><Aggregate Distinct="0" AggType="SUM"><ScalarOperator><Identifier><ColumnReference Column="Expr1006"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1005"></ColumnReference><ScalarOperator ScalarString="SUM([Expr1007])"><Aggregate Distinct="0" AggType="SUM"><ScalarOperator><Identifier><ColumnReference Column="Expr1007"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><GroupBy><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></GroupBy><RollupInfo HighestLevel="1"><RollupLevel Level="0"></RollupLevel><RollupLevel Level="1"></RollupLevel></RollupInfo><RelOp NodeId="3" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="28" EstimatedTotalSubtreeCost="0.0146456" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference><ColumnReference Column="Expr1006"></ColumnReference><ColumnReference Column="Expr1007"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1006"></ColumnReference><ScalarOperator ScalarString="COUNT_BIG([Expr1009])"><Aggregate Distinct="0" AggType="COUNT_BIG"><ScalarOperator><Identifier><ColumnReference Column="Expr1009"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1007"></ColumnReference><ScalarOperator ScalarString="SUM([Expr1009])"><Aggregate Distinct="0" AggType="SUM"><ScalarOperator><Identifier><ColumnReference Column="Expr1009"></ColumnReference></Identifier></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><GroupBy><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></GroupBy><RelOp NodeId="4" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.00010002" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0146445" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference><ColumnReference Column="Expr1009"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="5" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference><ColumnReference Column="Expr1009"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1009"></ColumnReference><ScalarOperator ScalarString="len([shop].[dbo].[customers].[name])"><Intrinsic FunctionName="len"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></Identifier></ScalarOperator></Intrinsic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="6" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="68" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="name"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Column="region"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp></Sort></RelOp></StreamAggregate></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT region, SUM(LEN(name)) AS l FROM customers GROUP BY ROLLUP (region)
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM (SELECT customer_id, status, amount FROM orders) s PIVOT (SUM(amount) FOR status IN ([new], [paid])) p', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| customer_id
: string
| shop
: string
| dbo
: string
| orders
: string
| customer_id
: null
< row
: string
| new
: null
: null
: null
: null
: null
< row
: string
| paid
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM (SELECT customer_id, status, amount FROM orders) s PIVOT (SUM(amount) FOR status IN ([new], [paid])) p
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM (SELECT customer_id, status, amount FROM orders) s PIVOT (SUM(amount) FOR status IN ([new], [paid])) p" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0146455" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x6916C27563D24CF5" QueryPlanHash="0x3E1102BA5941C929" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="512" SerialDesiredMemory="544" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="0" AvgRowSize="45" EstimatedTotalSubtreeCost="0.0146455" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1003"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1003"></ColumnReference><ScalarOperator ScalarString="CASE WHEN [Expr1013]=(0) THEN NULL ELSE [Expr1014] END"><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Column="Expr1013"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(0)"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Then><Else><ScalarOperator><Identifier><ColumnReference Column="Expr1014"></ColumnReference></Identifier></ScalarOperator></Else></IF></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="CASE WHEN [Expr1015]=(0) THEN NULL ELSE [Expr1016] END"><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Column="Expr1015"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(0)"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Then><Else><ScalarOperator><Identifier><ColumnReference Column="Expr1016"></ColumnReference></Identifier></ScalarOperator></Else></IF></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="1" PhysicalOp="Stream Aggregate" LogicalOp="Aggregate" EstimateRows="1" EstimateIO="0" EstimateCPU="1.1e-06" AvgRowSize="45" EstimatedTotalSubtreeCost="0.0146455" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1013"></ColumnReference><ColumnReference Column="Expr1014"></ColumnReference><ColumnReference Column="Expr1015"></ColumnReference><ColumnReference Column="Expr1016"></ColumnReference></OutputList><StreamAggregate><DefinedValues><DefinedValue><ColumnReference Column="Expr1013"></ColumnReference><ScalarOperator ScalarString="COUNT_BIG(CASE WHEN [shop].[dbo].[orders].[status]=&apos;new&apos; THEN [shop].[dbo].[orders].[amount] ELSE NULL END)"><Aggregate Distinct="0" AggType="COUNT_BIG"><ScalarOperator><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="&apos;new&apos;"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Else></IF></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1014"></ColumnReference><ScalarOperator ScalarString="SUM(CASE WHEN [shop].[dbo].[orders].[status]=&apos;new&apos; THEN [shop].[dbo].[orders].[amount] ELSE NULL END)"><Aggregate Distinct="0" AggType="SUM"><ScalarOperator><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="&apos;new&apos;"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Else></IF></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1015"></ColumnReference><ScalarOperator ScalarString="COUNT_BIG(CASE WHEN [shop].[dbo].[orders].[status]=&apos;paid&apos; THEN [shop].[dbo].[orders].[amount] ELSE NULL END)"><Aggregate Distinct="0" AggType="COUNT_BIG"><ScalarOperator><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="&apos;paid&apos;"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Else></IF></ScalarOperator></Aggregate></ScalarOperator></DefinedValue><DefinedValue><ColumnReference Column="Expr1016"></ColumnReference><ScalarOperator ScalarString="SUM(CASE WHEN [shop].[dbo].[orders].[status]=&apos;paid&apos; THEN [shop].[dbo].[orders].[amount] ELSE NULL END)"><Aggregate Distinct="0" AggType="SUM"><ScalarOperator><IF><Condition><ScalarOperator><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="&apos;paid&apos;"></Const></ScalarOperator></Compare></ScalarOperator></Condition><Then><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></Identifier></ScalarOperator></Then><Else><ScalarOperator><Const ConstValue="NULL"></Const></ScalarOperator></Else></IF></ScalarOperator></Aggregate></ScalarOperator></DefinedValue></DefinedValues><GroupBy><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></GroupBy><RelOp NodeId="2" PhysicalOp="Sort" LogicalOp="Sort" EstimateRows="1" EstimateIO="0.0112613" EstimateCPU="0.000100029" AvgRowSize="29" EstimatedTotalSubtreeCost="0.0146444" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Sort Distinct="0"><OrderBy><OrderByColumn Ascending="1"><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></OrderByColumn></OrderBy><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="29" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="amount"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="status"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Sort></RelOp></StreamAggregate></RelOp></ComputeScalar></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM (SELECT customer_id, status, amount FROM orders) s PIVOT (SUM(amount) FOR status IN ([new], [paid])) p
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT value FROM STRING_SPLIT(''a,b'', '','')', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| value
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT value FROM STRING_SPLIT('a,b', ',')
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT value FROM STRING_SPLIT(&apos;a,b&apos;, &apos;,&apos;)" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="5.0157e-05" StatementEstRows="50" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x368E4D12F495BC51" QueryPlanHash="0xE2CC8B27CFDEF3F5" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table-valued function" LogicalOp="Table-valued function" EstimateRows="50" EstimateIO="0" EstimateCPU="5.0157e-05" AvgRowSize="12" EstimatedTotalSubtreeCost="5.0157e-05" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Table="[STRING_SPLIT]" Column="value"></ColumnReference></OutputList><TableValuedFunction><DefinedValues><DefinedValue><ColumnReference Table="[STRING_SPLIT]" Column="value"></ColumnReference></DefinedValue></DefinedValues><Object Table="[STRING_SPLIT]"></Object><ParameterList><ScalarOperator ScalarString="&apos;a,b&apos;"><Const ConstValue="&apos;a,b&apos;"></Const></ScalarOperator><ScalarOperator ScalarString="&apos;,&apos;"><Const ConstValue="&apos;,&apos;"></Const></ScalarOperator></ParameterList></TableValuedFunction></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT value FROM STRING_SPLIT('a,b', ',')
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id FROM orders FOR JSON PATH', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| JSON_F52E2B61-18A1-11d1-B105-00805F49916B
: null
: null
: null
: null
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id FROM orders FOR JSON PATH
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id FROM orders FOR JSON PATH" StatementId="1" StatementCompId="1" StatementType="JSON SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x3BB2392D98C6046B" QueryPlanHash="0x072A1A6A5ECBBA17" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Column="id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id FROM orders FOR JSON PATH
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT nope FROM orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: null
: null
: null
: null
: null
: int64 207
< row
: null
: null
: null
: null
: null
: int64 11501
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM missing_table', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: null
: null
: null
: null
: null
: int64 208
< row
: null
: null
: null
: null
: null
: int64 11529
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [shop]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id FROM orders o JOIN customers c ON c.id = o.customer_id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: null
: null
: null
: null
: null
: int64 209
< row
: null
: null
: null
: null
: null
: int64 11501
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [archive]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM old_orders', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| archive
: string
| dbo
: string
| old_orders
: string
| id
: null
< row
: string
| customer_id
: string
| archive
: string
| dbo
: string
| old_orders
: string
| customer_id
: null
< row
: string
| amount
: string
| archive
: string
| dbo
: string
| old_orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM old_orders
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM old_orders" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0x9CC706A495BFAA9E" QueryPlanHash="0x3EECC6FA094179FC" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="id"></ColumnReference><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="customer_id"></ColumnReference><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[archive]" Schema="[dbo]" Table="[old_orders]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM old_orders
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| archive
: string
| dbo
: string
| old_orders
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [archive]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT id, amount FROM recent', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| archive
: string
| dbo
: string
| old_orders
: string
| id
: null
< row
: string
| amount
: string
| archive
: string
| dbo
: string
| old_orders
: string
| amount
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT id, amount FROM recent
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT id, amount FROM recent" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0032831" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="TRIVIAL" QueryHash="0xD0C0F7F9E8347435" QueryPlanHash="0x3EECC6FA094179FC" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="16" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="0" SerialDesiredMemory="0" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.0032035" EstimateCPU="7.96e-05" AvgRowSize="20" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object><Predicate><ScalarOperator ScalarString="[archive].[dbo].[old_orders].[id] as [o].[id]&gt;(100)"><Compare CompareOp="GT"><ScalarOperator><Identifier><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(100)"></Const></ScalarOperator></Compare></ScalarOperator></Predicate></TableScan></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT id, amount FROM recent
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| archive
: string
| dbo
: string
| recent
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [archive]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT c.name FROM shop.dbo.customers c JOIN old_orders o ON o.customer_id = c.id', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT c.name FROM shop.dbo.customers c JOIN old_orders o ON o.customer_id = c.id
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT c.name FROM shop.dbo.customers c JOIN old_orders o ON o.customer_id = c.id" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243408" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0xDEECFC5C10E45DBA" QueryPlanHash="0x176F48CF8ABC2E33" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="40" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177716" AvgRowSize="61" EstimatedTotalSubtreeCost="0.0243408" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[customers].[id] as [c].[id]=[archive].[dbo].[old_orders].[customer_id] as [o].[customer_id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="11" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="customer_id"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue></DefinedValues><Object Database="[archive]" Schema="[dbo]" Table="[old_orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT c.name FROM shop.dbo.customers c JOIN old_orders o ON o.customer_id = c.id
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| archive
: string
| dbo
: string
| old_orders
< row
: string
| shop
: string
| dbo
: string
| customers
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
> exec
| USE [archive]
< ok 0 -
> query
| SELECT name, source_database, source_schema, source_table, source_column, error_number
| FROM sys.dm_exec_describe_first_result_set(N'SELECT * FROM shop.dbo.order_totals', NULL, 1)
| WHERE ISNULL(is_hidden, 0) = 0
| ORDER BY column_ordinal
< result
< column "name" "NVARCHAR"
< column "source_database" "NVARCHAR"
< column "source_schema" "NVARCHAR"
< column "source_table" "NVARCHAR"
< column "source_column" "NVARCHAR"
< column "error_number" "INT"
< row
: string
| id
: string
| shop
: string
| dbo
: string
| orders
: string
| id
: null
< row
: string
| customer_id
: string
| shop
: string
| dbo
: string
| orders
: string
| customer_id
: null
< row
: string
| gross
: null
: null
: null
: null
: null
< row
: string
| customer_name
: string
| shop
: string
| dbo
: string
| customers
: string
| name
: null
> exec
| SET SHOWPLAN_XML ON
< ok 0 -
> query
| SELECT * FROM shop.dbo.order_totals
< result
< column "Microsoft SQL Server 2005 XML Showplan" "NVARCHAR"
< row
: string
| <ShowPlanXML xmlns="http://schemas.microsoft.com/sqlserver/2004/07/showplan" Version="1.564" Build="16.0.4305.5"><BatchSequence><Batch><Statements><StmtSimple StatementText="SELECT * FROM shop.dbo.order_totals" StatementId="1" StatementCompId="1" StatementType="SELECT" RetrievedFromCache="false" StatementSubTreeCost="0.0243489" StatementEstRows="1" SecurityPolicyApplied="false" StatementOptmLevel="FULL" QueryHash="0x661B6BBADDE4DBB0" QueryPlanHash="0x86D93C034BA36965" StatementOptmEarlyAbortReason="GoodEnoughPlanFound" CardinalityEstimationModelVersion="160"><StatementSetOptions QUOTED_IDENTIFIER="true" ARITHABORT="false" CONCAT_NULL_YIELDS_NULL="true" ANSI_NULLS="true" ANSI_PADDING="true" ANSI_WARNINGS="true" NUMERIC_ROUNDABORT="false"></StatementSetOptions><QueryPlan CachedPlanSize="48" CompileTime="0" CompileCPU="0" CompileMemory="0"><MemoryGrantInfo SerialRequiredMemory="1024" SerialDesiredMemory="1056" GrantedMemory="0" MaxUsedMemory="0"></MemoryGrantInfo><OptimizerHardwareDependentProperties EstimatedAvailableMemoryGrant="0" EstimatedPagesCached="0" EstimatedAvailableDegreeOfParallelism="0" MaxCompileMemory="0"></OptimizerHardwareDependentProperties><RelOp NodeId="0" PhysicalOp="Hash Match" LogicalOp="Inner Join" EstimateRows="1" EstimateIO="0" EstimateCPU="0.0177796" AvgRowSize="78" EstimatedTotalSubtreeCost="0.0243489" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><MemoryFractions Input="1" Output="1"></MemoryFractions><Hash><DefinedValues></DefinedValues><HashKeysBuild><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></HashKeysBuild><HashKeysProbe><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></HashKeysProbe><ProbeResidual><ScalarOperator ScalarString="[shop].[dbo].[orders].[customer_id] as [o].[customer_id]=[shop].[dbo].[customers].[id] as [c].[id]"><Compare CompareOp="EQ"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></Identifier></ScalarOperator></Compare></ScalarOperator></ProbeResidual><RelOp NodeId="1" PhysicalOp="Compute Scalar" LogicalOp="Compute Scalar" EstimateRows="1" EstimateIO="0" EstimateCPU="1e-07" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032832" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Column="Expr1004"></ColumnReference></OutputList><ComputeScalar><DefinedValues><DefinedValue><ColumnReference Column="Expr1004"></ColumnReference><ScalarOperator ScalarString="[shop].[dbo].[orders].[amount] as [o].[amount]*(1.1)"><Arithmetic Operation="MULT"><ScalarOperator><Identifier><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></Identifier></ScalarOperator><ScalarOperator><Const ConstValue="(1.1)"></Const></ScalarOperator></Arithmetic></ScalarOperator></DefinedValue></DefinedValues><RelOp NodeId="2" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="24" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="customer_id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" Column="amount"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[orders]" Alias="[o]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></ComputeScalar></RelOp><RelOp NodeId="3" PhysicalOp="Table Scan" LogicalOp="Table Scan" EstimateRows="1" EstimatedRowsRead="1" EstimateIO="0.003125" EstimateCPU="0.0001581" AvgRowSize="65" EstimatedTotalSubtreeCost="0.0032831" TableCardinality="0" Parallel="0" EstimateRebinds="0" EstimateRewinds="0" EstimatedExecutionMode="Row"><OutputList><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></OutputList><TableScan Ordered="0" ForcedIndex="0" ForceScan="0" NoExpandHint="0" Storage="RowStore"><DefinedValues><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="id"></ColumnReference></DefinedValue><DefinedValue><ColumnReference Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" Column="name"></ColumnReference></DefinedValue></DefinedValues><Object Database="[shop]" Schema="[dbo]" Table="[customers]" Alias="[c]" IndexKind="Heap" Storage="RowStore"></Object></TableScan></RelOp></Hash></RelOp></QueryPlan></StmtSimple></Statements></Batch></BatchSequence></ShowPlanXML>
> exec
| SET SHOWPLAN_XML OFF
< ok 0 -
> exec
| CREATE PROCEDURE dbo.omni_span_probe AS
| SELECT * FROM shop.dbo.order_totals
< ok 0 -
> query
| SELECT DISTINCT ISNULL(referenced_database_name, DB_NAME()),
| 	OBJECT_SCHEMA_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME()))),
| 	OBJECT_NAME(referenced_id, DB_ID(ISNULL(referenced_database_name, DB_NAME())))
| FROM sys.dm_sql_referenced_entities(N'dbo.omni_span_probe', N'OBJECT')
| WHERE referenced_minor_name IS NULL AND referenced_id IS NOT NULL
< result
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< column "" "NVARCHAR"
< row
: string
| shop
: string
| dbo
: string
| order_totals
> exec
| DROP PROCEDURE dbo.omni_span_probe
< ok 0 -
> exec
| USE [omni_TestQuerySpanFidelity]
< ok 0 -
