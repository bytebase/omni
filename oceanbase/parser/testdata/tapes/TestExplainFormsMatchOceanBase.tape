-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestExplainFormsMatchOceanBase default#1
> exec
| EXPLAIN FORMAT = TRADITIONAL SELECT 1
< ok 0 0
> exec
| EXPLAIN FORMAT = JSON SELECT 1
< ok 0 0
> exec
| EXPLAIN FORMAT = traditional SELECT 1
< ok 0 0
> exec
| EXPLAIN FORMAT = json SELECT 1
< ok 0 0
> exec
| EXPLAIN FORMAT = 'traditional' SELECT 1
< error
| 1064 42000 …near ''traditional' SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = 'TRADITIONAL' SELECT 1
< error
| 1064 42000 …near ''TRADITIONAL' SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = "traditional" SELECT 1
< error
| 1064 42000 …near '"traditional" SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = `TRADITIONAL` SELECT 1
< error
| 1064 42000 …near '`TRADITIONAL` SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = 'json' SELECT 1
< error
| 1064 42000 …near ''json' SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = '' SELECT 1
< error
| 1064 42000 …near ''' SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = TREE SELECT 1
< error
| 1064 42000 …near 'TREE SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = ROW SELECT 1
< error
| 1064 42000 …near 'ROW SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = DOT SELECT 1
< error
| 1064 42000 …near 'DOT SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = BASIC SELECT 1
< error
| 1064 42000 …near 'BASIC SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = EXTENDED SELECT 1
< error
| 1064 42000 …near 'EXTENDED SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = EXTENDED_NOADDR SELECT 1
< error
| 1064 42000 …near 'EXTENDED_NOADDR SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = PLANREGRESS SELECT 1
< error
| 1064 42000 …near 'PLANREGRESS SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = nosuchformat SELECT 1
< error
| 1064 42000 …near 'nosuchformat SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = binary SELECT 1
< error
| 1064 42000 …near 'binary SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = select SELECT 1
< error
| 1064 42000 …near 'select SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = 1 SELECT 1
< error
| 1064 42000 …near '1 SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec
| EXPLAIN FORMAT SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec
| EXPLAIN FORMAT 'json' SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = JSON FORMAT = TRADITIONAL SELECT 1
< error
| 1064 42000 …near 'FORMAT = TRADITIONAL SELECT 1' at line 1
> exec
| EXPLAIN EXTENDED SELECT 1
< ok 0 0
> exec
| EXPLAIN PARTITIONS SELECT 1
< ok 0 0
> exec
| EXPLAIN EXTENDED SELECT * FROM t1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN PARTITIONS UPDATE t1 SET a = 1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN EXTENDED FORMAT = TRADITIONAL SELECT 1
< error
| 1064 42000 …near '= TRADITIONAL SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = TRADITIONAL EXTENDED SELECT 1
< error
| 1064 42000 …near 'EXTENDED SELECT 1' at line 1
> exec
| EXPLAIN ANALYZE SELECT 1
< error
| 1064 42000 …near 'ANALYZE SELECT 1' at line 1
> exec
| EXPLAIN ANALYZE FORMAT = TRADITIONAL SELECT 1
< error
| 1064 42000 …near 'ANALYZE FORMAT = TRADITIONAL SELECT 1' at line 1
> exec
| EXPLAIN ANALYZE FORMAT = 'traditional' SELECT 1
< error
| 1064 42000 …near 'ANALYZE FORMAT = 'traditional' SELECT 1' at line 1
> exec
| EXPLAIN ANALYZE UPDATE t1 SET a = 1
< error
| 1064 42000 …near 'ANALYZE UPDATE t1 SET a = 1' at line 1
> exec
| EXPLAIN ANALYZE t1
< error
| 1064 42000 …near 'ANALYZE t1' at line 1
> exec
| DESC ANALYZE SELECT 1
< error
| 1064 42000 …near 'ANALYZE SELECT 1' at line 1
> exec
| EXPLAIN FOR CONNECTION 1
< error
| 1064 42000 …near 'FOR CONNECTION 1' at line 1
> exec
| EXPLAIN FOR connection 42
< error
| 1064 42000 …near 'FOR connection 42' at line 1
> exec
| EXPLAIN FORMAT = TRADITIONAL FOR CONNECTION 1
< error
| 1064 42000 …near 'FOR CONNECTION 1' at line 1
> exec
| EXPLAIN FORMAT = JSON FOR CONNECTION 1
< error
| 1064 42000 …near 'FOR CONNECTION 1' at line 1
> exec
| EXPLAIN EXTENDED FOR CONNECTION 1
< error
| 1064 42000 …near 'FOR CONNECTION 1' at line 1
> exec
| EXPLAIN FOR CONNECTION
< error
| 1064 42000 …near 'FOR CONNECTION' at line 1
> exec
| EXPLAIN FOR 1
< error
| 1064 42000 …near 'FOR 1' at line 1
> exec
| EXPLAIN FOR CONNECTION 'x'
< error
| 1064 42000 …near 'FOR CONNECTION 'x'' at line 1
> exec
| EXPLAIN ANALYZE FOR CONNECTION 1
< error
| 1064 42000 …near 'ANALYZE FOR CONNECTION 1' at line 1
> exec
| EXPLAIN SELECT 1
< ok 0 0
> exec
| EXPLAIN INSERT INTO t1 VALUES (1)
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN REPLACE INTO t1 VALUES (1)
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN UPDATE t1 SET a = 1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN DELETE FROM t1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN VALUES ROW(1)
< ok 0 0
> exec
| EXPLAIN VALUES ROW(1), ROW(2) ORDER BY 1
< ok 0 0
> exec
| EXPLAIN FORMAT = TRADITIONAL UPDATE t1 SET a = 1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN FORMAT = TRADITIONAL DELETE FROM t1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN FORMAT = TRADITIONAL INSERT INTO t1 VALUES (1)
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN FORMAT = TRADITIONAL REPLACE INTO t1 VALUES (1)
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN (SELECT 1)
< ok 0 0
> exec
| EXPLAIN WITH c AS (SELECT 1) SELECT * FROM c
< ok 0 0
> exec
| EXPLAIN FORMAT = TRADITIONAL (SELECT 1)
< ok 0 0
> exec
| EXPLAIN FORMAT = TRADITIONAL WITH c AS (SELECT 1) SELECT * FROM c
< ok 0 0
> exec
| EXPLAIN TABLE t1
< error
| 1064 42000 …near 'TABLE t1' at line 1
> exec
| EXPLAIN FORMAT = TRADITIONAL TABLE t1
< error
| 1064 42000 …near 'TABLE t1' at line 1
> exec
| EXPLAIN VALUES (1)
< error
| 1064 42000 …near '(1)' at line 1
> exec
| EXPLAIN ALTER TABLE t1 ADD INDEX (a)
< error
| 1064 42000 …near 'ALTER TABLE t1 ADD INDEX (a)' at line 1
> exec
| EXPLAIN ALTER TABLE t1 ADD a varchar(255)
< error
| 1064 42000 …near 'ALTER TABLE t1 ADD a varchar(255)' at line 1
> exec
| EXPLAIN FORMAT = TRADITIONAL ALTER TABLE t1 ADD INDEX (a)
< error
| 1064 42000 …near 'ALTER TABLE t1 ADD INDEX (a)' at line 1
> exec
| EXPLAIN CREATE TABLE t2 (a int)
< error
| 1064 42000 …near 'CREATE TABLE t2 (a int)' at line 1
> exec
| EXPLAIN DROP TABLE t1
< error
| 1064 42000 …near 'DROP TABLE t1' at line 1
> exec
| EXPLAIN SET @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec
| EXPLAIN FORMAT = TRADITIONAL SET @a = 1
< error
| 1064 42000 …near 'SET @a = 1' at line 1
> exec
| EXPLAIN FORMAT = TRADITIONAL SHOW TABLES
< error
| 1064 42000 …near 'SHOW TABLES' at line 1
> exec
| EXPLAIN t1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN t1 a
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN t1 `a`
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN t1 'abc'
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN t1 "abc"
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN t1 '%a%'
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN d1.t1 a
< error
| 1049 42000 Unknown database
> exec
| DESC t1 a
< error
| 1046 3D000 No database selected
> exec
| DESC t1 'abc'
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN BEGIN
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN t1 a.b
< error
| 1064 42000 …near 'a.b' at line 1
> exec
| EXPLAIN t1 123
< error
| 1064 42000 …near '123' at line 1
> exec
| EXPLAIN t1 select
< error
| 1064 42000 …near 'select' at line 1
> exec
| EXPLAIN t1 a b
< error
| 1064 42000 …near 'b' at line 1
> exec
| EXPLAIN FORMAT = TRADITIONAL t1
< error
| 1064 42000 …near 't1' at line 1
> exec
| EXPLAIN FORMAT = TRADITIONAL t1 a
< error
| 1064 42000 …near 't1 a' at line 1
> exec
| EXPLAIN EXTENDED t1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN PARTITIONS t1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN EXTENDED
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN 'sqldigest'
< error
| 1064 42000 …near ''sqldigest'' at line 1
> exec
| EXPLAIN BASIC SELECT 1
< ok 0 0
> exec
| EXPLAIN OUTLINE SELECT 1
< ok 0 0
> exec
| EXPLAIN EXTENDED_NOADDR SELECT 1
< ok 0 0
> exec
| EXPLAIN PLANREGRESS SELECT 1
< ok 0 0
> exec
| EXPLAIN PRETTY SELECT 1
< ok 0 0
> exec
| EXPLAIN PRETTY_COLOR SELECT 1
< ok 0 0
> exec
| EXPLAIN BASIC PRETTY SELECT 1
< ok 0 0
> exec
| EXPLAIN PARTITIONS PRETTY_COLOR SELECT 1
< ok 0 0
> exec
| EXPLAIN EXTENDED PRETTY DELETE FROM t1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN FORMAT = OBJECT_NAME_DISPLAY SELECT 1
< ok 0 0
> exec
| EXPLAIN PRETTY BASIC SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec
| EXPLAIN BASIC OUTLINE SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec
| EXPLAIN PRETTY PRETTY_COLOR SELECT 1
< error
| 1064 42000 …near 'SELECT 1' at line 1
> exec
| EXPLAIN BASIC FORMAT = JSON SELECT 1
< error
| 1064 42000 …near '= JSON SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = JSON PRETTY SELECT 1
< error
| 1064 42000 …near 'PRETTY SELECT 1' at line 1
> exec
| EXPLAIN basic
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN basic a
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN pretty
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN pretty basic
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN format
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN format a
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN basic pretty
< error
| 1064 42000 …near 'pretty' at line 1
> exec
| EXPLAIN SET STATEMENT_ID = 'a' SELECT 1
< ok 0 0
> exec
| EXPLAIN SET STATEMENT_ID = 1 SELECT 1
< ok 0 0
> exec
| EXPLAIN INTO p SELECT 1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN INTO p SET STATEMENT_ID = 'a' SELECT 1
< error
| 1046 3D000 No database selected
> exec
| EXPLAIN SET STATEMENT_ID = a SELECT 1
< error
| 1064 42000 …near 'a SELECT 1' at line 1
> exec
| EXPLAIN SET STATEMENT_ID = 'a' t1
< error
| 1064 42000 …near 't1' at line 1
> exec
| EXPLAIN SET STATEMENT_ID = 'a' INTO p SELECT 1
< error
| 1064 42000 …near 'INTO p SELECT 1' at line 1
> exec
| EXPLAIN INTO p
< error
| 1064 42000 …near 'p' at line 1
> exec
| EXPLAIN INTO p t1
< error
| 1064 42000 …near 't1' at line 1
> exec
| EXPLAIN INTO d.p SELECT 1
< error
| 1064 42000 …near 'd.p SELECT 1' at line 1
> exec
| EXPLAIN INTO p BASIC SELECT 1
< error
| 1064 42000 …near 'BASIC SELECT 1' at line 1
> exec
| EXPLAIN BASIC INTO p SELECT 1
< error
| 1064 42000 …near 'INTO p SELECT 1' at line 1
> exec
| EXPLAIN FORMAT = JSON INTO p SELECT 1
< error
| 1064 42000 …near 'INTO p SELECT 1' at line 1
> exec
| DESC SELECT 1
< ok 0 0
> exec
| DESC FORMAT = TRADITIONAL SELECT 1
< ok 0 0
> exec
| DESCRIBE BASIC SELECT 1
< ok 0 0
> exec
| DESC INTO p SELECT 1
< error
| 1046 3D000 No database selected
