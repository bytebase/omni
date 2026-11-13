-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestNameTokenFormsOracle default#1
> exec
| PREPARE omni_name_form_stmt FROM ?
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
| 1146 42S02 Table 'omni_t_mysql_parser.@a' doesn't exist
> exec ^
: string
| SELECT * FROM `@a`
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| SELECT * FROM _utf8.t
< error
| 1064 42000 …near '_utf8.t' at line 1
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
| SELECT * FROM `_utf8`
< error
| 1146 42S02 Table 'omni_t_mysql_parser._utf8' doesn't exist
> exec ^
: string
| SELECT * FROM `_utf8`
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE `_utf8` (b int)
< ok 0 0
> exec ^
: string
| SELECT * FROM d._utf8
< error
| 1142 42000 SELECT command denied to user 'omni_test'@'172.17.0.1' for table '_utf8'
> exec ^
: string
| SELECT * FROM d._utf8
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a._utf8 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a._utf8 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| SELECT * FROM _filename
< error
| 1146 42S02 Table 'omni_t_mysql_parser._filename' doesn't exist
> exec ^
: string
| SELECT * FROM _filename
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM _utf8mb4_bin
< error
| 1146 42S02 Table 'omni_t_mysql_parser._utf8mb4_bin' doesn't exist
> exec ^
: string
| SELECT * FROM _utf8mb4_bin
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM _utf8_general_ci
< error
| 1146 42S02 Table 'omni_t_mysql_parser._utf8_general_ci' doesn't exist
> exec ^
: string
| SELECT * FROM _utf8_general_ci
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM _foo
< error
| 1146 42S02 Table 'omni_t_mysql_parser._foo' doesn't exist
> exec ^
: string
| SELECT * FROM _foo
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT _utf8mb4'x'
< ok 0 0
> exec ^
: string
| SELECT a.1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a.1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1b FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1b FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1e5 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1e5 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1e+5 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1e+5 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1e FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1e FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.0x10 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.0x10 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.0b1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.0b1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.123456789012345678901234567890 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.123456789012345678901234567890 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1.5 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1.5 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1.e1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1.e1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT d.t.1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT d.t.1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.`1` FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.`1` FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT select.1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT select.1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1a.1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT 1a.1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 12e.1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT 12e.1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a.$1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a.$1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a.1$ FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a.1$ FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1-1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1-1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.*, t.1 FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.*, t.1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1 'x' FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1 'x' FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1 AS x FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1 AS x FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT COUNT(t.1) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT COUNT(t.1) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t WHERE t.1 = 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE t.1 = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t ORDER BY t.1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t ORDER BY t.1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT t.1() FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT t.1() FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT d.1(1)
< error
| 1370 42000 execute command denied to user 'omni_test'@'%' for routine 'd.1'
> exec ^
: string
| SELECT d.1(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| UPDATE t SET t.1 = 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| UPDATE t SET t.1 = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.1 = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.1 = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT @@global.1
< error
| 1193 HY000 Unknown system variable '1'
> exec ^
: string
| SELECT @@autocommit.1
< error
| 1193 HY000 Unknown system variable 'autocommit.1'
> exec ^
: string
| SELECT @@session.autocommit.1
< error
| 1193 HY000 Unknown system variable 'autocommit.1'
> exec ^
: string
| CREATE TABLE a.1b (c int)
< error
| 1142 42000 CREATE command denied to user 'omni_test'@'172.17.0.1' for table '1b'
> exec ^
: string
| CREATE TABLE a.1b (c int)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE a.1 (c int)
< error
| 1142 42000 CREATE command denied to user 'omni_test'@'172.17.0.1' for table '1'
> exec ^
: string
| CREATE TABLE a.1 (c int)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE a.1e5 (c int)
< error
| 1142 42000 CREATE command denied to user 'omni_test'@'172.17.0.1' for table '1e5'
> exec ^
: string
| CREATE TABLE a.1e5 (c int)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE a. 1b (c int)
< error
| 1142 42000 CREATE command denied to user 'omni_test'@'172.17.0.1' for table '1b'
> exec ^
: string
| CREATE TABLE a. 1b (c int)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM 1db.1table
< error
| 1142 42000 SELECT command denied to user 'omni_test'@'172.17.0.1' for table '1table'
> exec ^
: string
| SELECT * FROM 1db.1table
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM d.1
< error
| 1142 42000 SELECT command denied to user 'omni_test'@'172.17.0.1' for table '1'
> exec ^
: string
| SELECT * FROM d.1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM d.1 x
< error
| 1142 42000 SELECT command denied to user 'omni_test'@'172.17.0.1' for table '1'
> exec ^
: string
| SELECT * FROM d.1 x
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO d.1 VALUES (1)
< error
| 1142 42000 INSERT command denied to user 'omni_test'@'172.17.0.1' for table '1'
> exec ^
: string
| INSERT INTO d.1 VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1.5
< ok 0 0
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
| SELECT 1e5.1 FROM t
< error
| 1064 42000 …near '.1 FROM t' at line 1
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
| 1064 42000 …near '.1' at line 1
> exec ^
: string
| CREATE TABLE a.1.5 (c int)
< error
| 1064 42000 …near '.5 (c int)' at line 1
> exec ^
: string
| CREATE TABLE a .1b (c int)
< error
| 1064 42000 …near '.1b (c int)' at line 1
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW .1 = 1
< error
| 1064 42000 …near '.1 = 1' at line 1
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
| 1064 42000 …near '' at line 1
> exec ^
: string
| FLUSH LOCAL
< error
| 1064 42000 …near '' at line 1
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
| FLUSH ERROR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| FLUSH ENGINE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| FLUSH BINARY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| FLUSH RELAY
< error
| 1064 42000 …near '' at line 1
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
| 1064 42000 …near 'QUERY HOSTS' at line 1
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
| FLUSH QUERY CACHE
< error
| 1064 42000 …near 'QUERY CACHE' at line 1
> exec ^
: string
| FLUSH DES_KEY_FILE
< error
| 1064 42000 …near 'DES_KEY_FILE' at line 1
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
| FLUSH RELAY LOG foo
< error
| 1064 42000 …near 'LOG foo' at line 1
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
| FLUSH BINARY LOGS, ENGINE LOGS, ERROR LOGS, GENERAL LOGS, SLOW LOGS, RELAY LOGS, HOSTS, PRIVILEGES, LOGS, STATUS, USER_RESOURCES, OPTIMIZER_COSTS
< ok 0 0
> exec ^
: string
| FLUSH HOSTS,
< error
| 1064 42000 …near '' at line 1
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
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL "c"
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL ''
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL 'c', RELAY LOGS
< ok 0 0
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL c
< error
| 1064 42000 …near 'c' at line 1
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL _utf8mb4'c'
< error
| 1064 42000 …near '_utf8mb4'c'' at line 1
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL N'c'
< error
| 1064 42000 …near 'N'c'' at line 1
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL 'a' 'b'
< error
| 1064 42000 …near ''b'' at line 1
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL 'a
| b'
< error
| 1525 HY000 Incorrect argument contains not-allowed LF value: 'a
| b'
> exec ^
: string
| FLUSH RELAY LOGS FOR CHANNEL 'a
| b'
| ; SELECT FROM
< error
| 1525 HY000 Incorrect argument contains not-allowed LF value: 'a
| b'
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
| FLUSH TABLES t, u
< ok 0 0
> exec ^
: string
| FLUSH TABLES d.t
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
| FLUSH TABLES t FOR EXPORT
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
| FLUSH TABLES t,
< error
| 1064 42000 …near '' at line 1
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
| FLUSH TABLES FOR EXPORT
< error
| 1064 42000 No tables used near 'FOR EXPORT' at line 1
> exec ^
: string
| FLUSH TABLES WITH READ
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| FLUSH TABLES WITH
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| FLUSH TABLES t FOR
< error
| 1064 42000 …near '' at line 1
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
> exec ^
: string
| FLUSH TABLES @a
< error
| 1064 42000 …near '@a' at line 1
> exec ^
: string
| FLUSH TABLES _utf8
< error
| 1064 42000 …near '_utf8' at line 1
