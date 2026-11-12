-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestNameTokenFormsOceanBaseOracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_name_token_forms
< ok 1 0
> exec
| USE omni_name_token_forms
< ok 0 0
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| CREATE TABLE @a (b int)
< error
| 1064 42000 …near '@a (b int)' at line 1
> exec ^
: string
| SELECT * FROM @a
< error
| 1064 42000 …near '@a' at line 1
> exec ^
: string
| SELECT * FROM t.@a
< error
| 1064 42000 …near '@a' at line 1
> exec ^
: string
| SELECT * FROM @a.t
< error
| 1064 42000 …near '@a.t' at line 1
> exec ^
: string
| SELECT 1 AS @a
< error
| 1064 42000 …near '@a' at line 1
> exec ^
: string
| SELECT * FROM `@a`
< error
| 1146 42S02 Table 'omni_name_token_forms.@a' doesn't exist
> exec ^
: string
| CREATE TABLE `@a` (b int)
< ok 0 0
> exec ^
: string
| SELECT @a
< ok 0 0
> exec ^
: string
| SELECT @a.1
< ok 0 0
> exec ^
: string
| SET @a = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE _utf8 (b int)
< error
| 1064 42000 …near '_utf8 (b int)' at line 1
> exec ^
: string
| SELECT * FROM _utf8
< error
| 1064 42000 …near '_utf8' at line 1
> exec ^
: string
| SELECT * FROM _UTF8MB4
< error
| 1064 42000 …near '_UTF8MB4' at line 1
> exec ^
: string
| SELECT * FROM _binary
< error
| 1064 42000 …near '_binary' at line 1
> exec ^
: string
| SELECT * FROM _gb18030
< error
| 1064 42000 …near '_gb18030' at line 1
> exec ^
: string
| SELECT * FROM _hkscs
< error
| 1064 42000 …near '_hkscs' at line 1
> exec ^
: string
| SELECT _utf8 FROM t
< error
| 1064 42000 …near 'FROM t' at line 1
> exec ^
: string
| SELECT 1 AS _utf8
< error
| 1064 42000 …near '_utf8' at line 1
> exec ^
: string
| CREATE TABLE t (_latin1 int)
< error
| 1064 42000 …near '_latin1 int)' at line 1
> exec ^
: string
| SELECT * FROM _utf8. t
< error
| 1064 42000 …near '_utf8. t' at line 1
> exec ^
: string
| SELECT * FROM _utf8 .t
< error
| 1064 42000 …near '_utf8 .t' at line 1
> exec ^
: string
| SELECT * FROM d. _utf8
< error
| 1064 42000 …near '_utf8' at line 1
> exec ^
: string
| SELECT d. _utf8 FROM t
< error
| 1064 42000 …near '_utf8 FROM t' at line 1
> exec ^
: string
| SELECT * FROM `_utf8`
< error
| 1146 42S02 Table 'omni_name_token_forms._utf8' doesn't exist
> exec ^
: string
| CREATE TABLE `_utf8` (b int)
< ok 0 0
> exec ^
: string
| SELECT * FROM d._utf8
< error
| 1146 42S02 Table 'd._utf8' doesn't exist
> exec ^
: string
| SELECT * FROM _utf8.t
< error
| 1146 42S02 Table '_utf8.t' doesn't exist
> exec ^
: string
| SELECT * FROM _utf8.1
< error
| 1146 42S02 Table '_utf8.1' doesn't exist
> exec ^
: string
| SELECT _utf8.t.c FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT * FROM t._utf8 x
< error
| 1146 42S02 Table 't._utf8' doesn't exist
> exec ^
: string
| SELECT a._utf8 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT * FROM _ucs2
< error
| 1146 42S02 Table 'omni_name_token_forms._ucs2' doesn't exist
> exec ^
: string
| SELECT * FROM _cp1250
< error
| 1146 42S02 Table 'omni_name_token_forms._cp1250' doesn't exist
> exec ^
: string
| SELECT * FROM _filename
< error
| 1146 42S02 Table 'omni_name_token_forms._filename' doesn't exist
> exec ^
: string
| SELECT * FROM _utf8mb4_bin
< error
| 1146 42S02 Table 'omni_name_token_forms._utf8mb4_bin' doesn't exist
> exec ^
: string
| SELECT * FROM _utf8_general_ci
< error
| 1146 42S02 Table 'omni_name_token_forms._utf8_general_ci' doesn't exist
> exec ^
: string
| SELECT * FROM _foo
< error
| 1146 42S02 Table 'omni_name_token_forms._foo' doesn't exist
> exec ^
: string
| SELECT _utf8mb4'x'
< ok 0 0
> exec ^
: string
| SELECT a.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1b FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1e5 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1e+5 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1e FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.0x10 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.0b1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.123456789012345678901234567890 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1.5 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1.e1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT d.t.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.`1` FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT select.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1a.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 12e.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT a.$1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT a.1$ FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1-1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.*, t.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1 'x' FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1 AS x FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT COUNT(t.1) FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE t.1 = 1
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t ORDER BY t.1
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT t.1() FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT d.1(1)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| UPDATE t SET t.1 = 1
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.1 = 1
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| CREATE TABLE a.1b (c int)
< error
| 1049 42000 Unknown database 'a'
> exec ^
: string
| CREATE TABLE a.1 (c int)
< error
| 1049 42000 Unknown database 'a'
> exec ^
: string
| CREATE TABLE a.1e5 (c int)
< error
| 1049 42000 Unknown database 'a'
> exec ^
: string
| CREATE TABLE a. 1b (c int)
< error
| 1049 42000 Unknown database 'a'
> exec ^
: string
| SELECT * FROM 1db.1table
< error
| 1146 42S02 Table '1db.1table' doesn't exist
> exec ^
: string
| SELECT * FROM d.1
< error
| 1146 42S02 Table 'd.1' doesn't exist
> exec ^
: string
| SELECT * FROM d.1 x
< error
| 1146 42S02 Table 'd.1' doesn't exist
> exec ^
: string
| INSERT INTO d.1 VALUES (1)
< error
| 1146 42S02 Table 'd.1' doesn't exist
> exec ^
: string
| SELECT 1.5
< ok 0 0
> exec ^
: string
| SELECT 1.a FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1e5.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1e5.a FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1e5.1.2 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 0x10.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 0b1.a FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 0x.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1e.1 FROM t
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a=1e5.1
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a=0x1F.b
< error
| 1146 42S02 Table 'omni_name_token_forms.t' doesn't exist
> exec ^
: string
| SELECT 1e5 .a FROM t
< error
| 1064 42000 …near '.a FROM t' at line 1
> exec ^
: string
| SELECT 1e5.`a` FROM t
< error
| 1064 42000 …near '.`a` FROM t' at line 1
> exec ^
: string
| SELECT 1e+5.1 FROM t
< error
| 1064 42000 …near '.1 FROM t' at line 1
> exec ^
: string
| SELECT 12.5.a FROM t
< error
| 1064 42000 …near '.a FROM t' at line 1
> exec ^
: string
| SELECT `t`.1 FROM t
< error
| 1064 42000 …near '.1 FROM t' at line 1
> exec ^
: string
| SELECT `t`.1b FROM t
< error
| 1064 42000 …near '.1b FROM t' at line 1
> exec ^
: string
| SELECT t .1 FROM t
< error
| 1064 42000 …near '.1 FROM t' at line 1
> exec ^
: string
| SELECT t. 1 FROM t
< error
| 1064 42000 …near '1 FROM t' at line 1
> exec ^
: string
| SELECT t . 1 FROM t
< error
| 1064 42000 …near '1 FROM t' at line 1
> exec ^
: string
| SELECT t.1. FROM t
< error
| 1064 42000 …near 'FROM t' at line 1
> exec ^
: string
| SELECT t.1.2.3 FROM t
< error
| 1064 42000 …near '.3 FROM t' at line 1
> exec ^
: string
| SELECT 1 .5
< error
| 1064 42000 …near '.5' at line 1
> exec ^
: string
| SELECT 'a'.1
< error
| 1064 42000 …near '.1' at line 1
> exec ^
: string
| SELECT x'0a'.1
< error
| 1064 42000 …near '.1' at line 1
> exec ^
: string
| SELECT (t).1 FROM t
< error
| 1064 42000 …near '.1 FROM t' at line 1
> exec ^
: string
| SELECT @`a`.1
< error
| 1064 42000 …near '@`a`.1' at line 1
> exec ^
: string
| SELECT @@global.1
< error
| 1064 42000 …near '.1' at line 1
> exec ^
: string
| SELECT @@autocommit.1
< error
| 1064 42000 …near '.1' at line 1
> exec ^
: string
| SELECT @@global.autocommit
< ok 0 0
> exec ^
: string
| CREATE TABLE a .1b (c int)
< error
| 1064 42000 …near '.1b (c int)' at line 1
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW .1 = 1
< error
| 1064 42000 …near 'SET NEW .1 = 1' at line 1
> exec ^
: string
| CREATE TABLE `a`.1b (c int)
< error
| 1064 42000 …near '.1b (c int)' at line 1
> exec ^
: string
| FLUSH tidb
< error
| 1064 42000 …near 'tidb' at line 1
> exec ^
: string
| FLUSH foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| FLUSH
< error
| 1064 42000 …near 'FLUSH' at line 1
> exec ^
: string
| FLUSH LOCAL
< error
| 1064 42000 …near 'LOCAL' at line 1
> exec ^
: string
| FLUSH `hosts`
< error
| 1064 42000 …near '`hosts`' at line 1
> exec ^
: string
| FLUSH 'hosts'
< error
| 1064 42000 …near ''hosts'' at line 1
> exec ^
: string
| FLUSH COLUMNS
< error
| 1064 42000 …near 'COLUMNS' at line 1
> exec ^
: string
| FLUSH RESOURCES
< error
| 1064 42000 …near 'RESOURCES' at line 1
> exec ^
: string
| FLUSH QUERY_CACHE
< error
| 1064 42000 …near 'QUERY_CACHE' at line 1
> exec ^
: string
| FLUSH ERROR
< error
| 1064 42000 …near 'ERROR' at line 1
> exec ^
: string
| FLUSH ENGINE
< error
| 1064 42000 …near 'ENGINE' at line 1
> exec ^
: string
| FLUSH BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec ^
: string
| FLUSH RELAY
< error
| 1064 42000 …near 'RELAY' at line 1
> exec ^
: string
| FLUSH LOGS RELAY
< error
| 1064 42000 …near 'RELAY' at line 1
> exec ^
: string
| FLUSH QUERY
< error
| 1064 42000 …near 'QUERY' at line 1
> exec ^
: string
| FLUSH QUERY HOSTS
< error
| 1064 42000 …near 'HOSTS' at line 1
> exec ^
: string
| FLUSH RELAY foo bar
< error
| 1064 42000 …near 'foo bar' at line 1
> exec ^
: string
| FLUSH RELAY HOSTS foo
< error
| 1064 42000 …near 'HOSTS foo' at line 1
> exec ^
: string
| FLUSH RELAY LOGS foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| FLUSH SSL
< error
| 1064 42000 …near 'SSL' at line 1
> exec ^
: string
| FLUSH THREADS
< error
| 1064 42000 …near 'THREADS' at line 1
> exec ^
: string
| FLUSH SLAVE
< error
| 1064 42000 …near 'SLAVE' at line 1
> exec ^
: string
| FLUSH REPLICA
< error
| 1064 42000 …near 'REPLICA' at line 1
> exec ^
: string
| FLUSH MASTER
< error
| 1064 42000 …near 'MASTER' at line 1
> exec ^
: string
| FLUSH SESSION STATUS
< error
| 1064 42000 …near 'SESSION STATUS' at line 1
> exec ^
: string
| FLUSH GLOBAL STATUS
< error
| 1064 42000 …near 'GLOBAL STATUS' at line 1
> exec ^
: string
| FLUSH user_statistics
< error
| 1064 42000 …near 'user_statistics' at line 1
> exec ^
: string
| FLUSH BINARY LOGS DELETE_DOMAIN_ID = (1)
< error
| 1064 42000 …near 'DELETE_DOMAIN_ID = (1)' at line 1
> exec ^
: string
| FLUSH RELAY LOGS 'c'
< error
| 1064 42000 …near ''c'' at line 1
> exec ^
: string
| FLUSH TABLES WITH READ LOCK AND DISABLE CHECKPOINT
< error
| 1064 42000 …near 'AND DISABLE CHECKPOINT' at line 1
> exec ^
: string
| FLUSH HOSTS
< ok 0 0
> exec ^
: string
| FLUSH LOGS
< ok 0 0
> exec ^
: string
| FLUSH PRIVILEGES
< ok 0 0
> exec ^
: string
| FLUSH PRIVILEGES;
< ok 0 0
> exec ^
: string
| FLUSH OPTIMIZER_COSTS
< ok 0 0
> exec ^
: string
| FLUSH STATUS
< ok 0 0
> exec ^
: string
| FLUSH USER_RESOURCES
< ok 0 0
> exec ^
: string
| FLUSH DES_KEY_FILE
< ok 0 0
> exec ^
: string
| FLUSH QUERY CACHE
< ok 0 0
> exec ^
: string
| FLUSH LOCAL QUERY CACHE
< ok 0 0
> exec ^
: string
| FLUSH ERROR LOGS
< ok 0 0
> exec ^
: string
| FLUSH ENGINE LOGS
< ok 0 0
> exec ^
: string
| FLUSH GENERAL LOGS
< ok 0 0
> exec ^
: string
| FLUSH SLOW LOGS
< ok 0 0
> exec ^
: string
| FLUSH BINARY LOGS
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOGS
< ok 0 0
> exec ^
: string
| FLUSH LOCAL HOSTS
< ok 0 0
> exec ^
: string
| FLUSH LOCAL STATUS
< ok 0 0
> exec ^
: string
| FLUSH NO_WRITE_TO_BINLOG PRIVILEGES
< ok 0 0
> exec ^
: string
| FLUSH hosts, status
< ok 0 0
> exec ^
: string
| FLUSH LOGS, LOGS
< ok 0 0
> exec ^
: string
| FLUSH STATUS, STATUS
< ok 0 0
> exec ^
: string
| FLUSH LOGS, RELAY LOGS
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOGS, RELAY LOGS
< ok 0 0
> exec ^
: string
| FLUSH BINARY LOGS, ENGINE LOGS, ERROR LOGS, GENERAL LOGS, SLOW LOGS, RELAY LOGS, HOSTS, PRIVILEGES, LOGS, STATUS, USER_RESOURCES, OPTIMIZER_COSTS
< ok 0 0
> exec ^
: string
| FLUSH HOSTS,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| FLUSH , HOSTS
< error
| 1064 42000 …near ', HOSTS' at line 1
> exec ^
: string
| FLUSH HOSTS HOSTS
< error
| 1064 42000 …near 'HOSTS' at line 1
> exec ^
: string
| FLUSH status, tables
< error
| 1064 42000 …near 'tables' at line 1
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL 'c'
< error
| 1064 42000 …near 'FOR CHANNEL 'c'' at line 1
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL c
< error
| 1064 42000 …near 'FOR CHANNEL c' at line 1
> exec ^
: string
| FLUSH RELAY LOG foo
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOG `foo`
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOG 1a
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOG _cp1250
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOG foo, HOSTS
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOG
< error
| 1064 42000 …near 'LOG' at line 1
> exec ^
: string
| FLUSH RELAY LOG 'foo'
< error
| 1064 42000 …near ''foo'' at line 1
> exec ^
: string
| FLUSH RELAY LOG a.b
< error
| 1064 42000 …near 'a.b' at line 1
> exec ^
: string
| FLUSH RELAY LOG @a
< error
| 1064 42000 …near '@a' at line 1
> exec ^
: string
| FLUSH RELAY LOG _utf8
< error
| 1064 42000 …near '_utf8' at line 1
> exec ^
: string
| FLUSH RELAY LOG 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| FLUSH RELAY LOG foo bar
< error
| 1064 42000 …near 'bar' at line 1
> exec ^
: string
| FLUSH RELAY LOG select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| FLUSH RELAY LOG status
< error
| 1064 42000 …near 'status' at line 1
> exec ^
: string
| FLUSH RELAY LOG hosts
< error
| 1064 42000 …near 'hosts' at line 1
> exec ^
: string
| FLUSH RELAY LOG relay
< error
| 1064 42000 …near 'relay' at line 1
> exec ^
: string
| FLUSH RELAY LOG logs
< error
| 1064 42000 …near 'logs' at line 1
> exec ^
: string
| FLUSH RELAY LOG tables
< error
| 1064 42000 …near 'tables' at line 1
> exec ^
: string
| FLUSH RELAY LOG account
< error
| 1064 42000 …near 'account' at line 1
> exec ^
: string
| FLUSH RELAY LOG names
< error
| 1064 42000 …near 'names' at line 1
> exec ^
: string
| FLUSH TABLES
< ok 0 0
> exec ^
: string
| FLUSH TABLE
< ok 0 0
> exec ^
: string
| FLUSH TABLES t
< ok 0 0
> exec ^
: string
| FLUSH TABLE t
< ok 0 0
> exec ^
: string
| FLUSH TABLES t, u
< ok 0 0
> exec ^
: string
| FLUSH TABLES test.t
< ok 0 0
> exec ^
: string
| FLUSH TABLES `t`
< ok 0 0
> exec ^
: string
| FLUSH TABLES WITH READ LOCK
< ok 0 0
> exec ^
: string
| FLUSH TABLES t WITH READ LOCK
< ok 0 0
> exec ^
: string
| FLUSH TABLES t, u WITH READ LOCK
< ok 0 0
> exec ^
: string
| FLUSH TABLES t FOR EXPORT
< ok 0 0
> exec ^
: string
| FLUSH TABLE t FOR EXPORT
< ok 0 0
> exec ^
: string
| FLUSH TABLES test.t FOR EXPORT
< ok 0 0
> exec ^
: string
| FLUSH TABLES t, u FOR EXPORT
< ok 0 0
> exec ^
: string
| FLUSH LOCAL TABLES
< ok 0 0
> exec ^
: string
| FLUSH LOCAL TABLES t WITH READ LOCK
< ok 0 0
> exec ^
: string
| FLUSH NO_WRITE_TO_BINLOG TABLES t WITH READ LOCK
< ok 0 0
> exec ^
: string
| FLUSH TABLES t,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| FLUSH TABLES 't'
< error
| 1064 42000 …near ''t'' at line 1
> exec ^
: string
| FLUSH TABLES t.*
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| FLUSH TABLES t AS x
< error
| 1064 42000 …near 'AS x' at line 1
> exec ^
: string
| FLUSH TABLES t PARTITION (p0)
< error
| 1064 42000 …near 'PARTITION (p0)' at line 1
> exec ^
: string
| FLUSH TABLES FOR EXPORT
< error
| 1064 42000 …near 'FOR EXPORT' at line 1
> exec ^
: string
| FLUSH TABLES WITH READ
< error
| 1064 42000 …near 'READ' at line 1
> exec ^
: string
| FLUSH TABLES WITH
< error
| 1064 42000 …near 'WITH' at line 1
> exec ^
: string
| FLUSH TABLES t FOR
< error
| 1064 42000 …near 'FOR' at line 1
> exec ^
: string
| FLUSH TABLES t HOSTS
< error
| 1064 42000 …near 'HOSTS' at line 1
> exec ^
: string
| FLUSH TABLES, HOSTS
< error
| 1064 42000 …near ', HOSTS' at line 1
> exec ^
: string
| FLUSH TABLES t FOR EXPORT, u
< error
| 1064 42000 …near ', u' at line 1
> exec ^
: string
| FLUSH TABLES t WITH READ LOCK FOR EXPORT
< error
| 1064 42000 …near 'FOR EXPORT' at line 1
