-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestBindingKeywordsStayIdentifiersOnTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE TABLE binding (binding INT)
< ok 0 0
> exec ^
: string
| SELECT binding FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT binding FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.binding FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.binding FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS binding
< ok 0 0
> exec ^
: string
| SELECT * FROM t binding
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t binding
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS binding
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t AS binding
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM binding.t
< error
| 1146 42S02 Table 'binding.t' doesn't exist
> exec ^
: string
| SELECT * FROM binding.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t binding WHERE binding.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t binding WHERE binding.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t binding SET binding.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t binding SET binding.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO binding (binding) VALUES (1)
< error
| 1146 42S02 Table 'test.binding' doesn't exist
> exec ^
: string
| INSERT INTO binding (binding) VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX binding ON t (binding)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX binding ON t (binding)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX binding ON t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX binding ON t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN binding INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN binding INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO binding
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO binding
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE DATABASE binding
< ok 0 0
> exec ^
: string
| SAVEPOINT binding
< ok 0 0
> exec ^
: string
| RELEASE SAVEPOINT binding
< ok 0 0
> exec ^
: string
| SET binding = 1
< ok 0 0
> exec ^
: string
| SET @binding = 1
< ok 0 0
> exec ^
: string
| CREATE VIEW binding AS SELECT 1
< ok 0 0
> exec ^
: string
| SELECT * FROM t ORDER BY binding
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t ORDER BY binding
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t GROUP BY binding HAVING binding > 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t GROUP BY binding HAVING binding > 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE binding
< error
| 1146 42S02 Table 'test.binding' doesn't exist
> exec ^
: string
| ANALYZE TABLE binding
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRUNCATE TABLE binding
< error
| 1146 42S02 Table 'test.binding' doesn't exist
> exec ^
: string
| TRUNCATE TABLE binding
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE digest (digest INT)
< ok 0 0
> exec ^
: string
| SELECT digest FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT digest FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.digest FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.digest FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS digest
< ok 0 0
> exec ^
: string
| SELECT * FROM t digest
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t digest
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS digest
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t AS digest
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM digest.t
< error
| 1146 42S02 Table 'digest.t' doesn't exist
> exec ^
: string
| SELECT * FROM digest.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t digest WHERE digest.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t digest WHERE digest.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t digest SET digest.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t digest SET digest.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO digest (digest) VALUES (1)
< error
| 1146 42S02 Table 'test.digest' doesn't exist
> exec ^
: string
| INSERT INTO digest (digest) VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX digest ON t (digest)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX digest ON t (digest)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX digest ON t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX digest ON t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN digest INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN digest INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO digest
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO digest
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE DATABASE digest
< ok 0 0
> exec ^
: string
| SAVEPOINT digest
< ok 0 0
> exec ^
: string
| RELEASE SAVEPOINT digest
< ok 0 0
> exec ^
: string
| SET digest = 1
< ok 0 0
> exec ^
: string
| SET @digest = 1
< ok 0 0
> exec ^
: string
| CREATE VIEW digest AS SELECT 1
< ok 0 0
> exec ^
: string
| SELECT * FROM t ORDER BY digest
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t ORDER BY digest
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t GROUP BY digest HAVING digest > 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t GROUP BY digest HAVING digest > 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE digest
< error
| 1146 42S02 Table 'test.digest' doesn't exist
> exec ^
: string
| ANALYZE TABLE digest
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRUNCATE TABLE digest
< error
| 1146 42S02 Table 'test.digest' doesn't exist
> exec ^
: string
| TRUNCATE TABLE digest
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE plan (plan INT)
< ok 0 0
> exec ^
: string
| SELECT plan FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT plan FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.plan FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.plan FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS plan
< ok 0 0
> exec ^
: string
| SELECT * FROM t plan
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t plan
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS plan
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t AS plan
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM plan.t
< error
| 1146 42S02 Table 'plan.t' doesn't exist
> exec ^
: string
| SELECT * FROM plan.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t plan WHERE plan.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t plan WHERE plan.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t plan SET plan.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t plan SET plan.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO plan (plan) VALUES (1)
< error
| 1146 42S02 Table 'test.plan' doesn't exist
> exec ^
: string
| INSERT INTO plan (plan) VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX plan ON t (plan)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX plan ON t (plan)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX plan ON t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX plan ON t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN plan INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN plan INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO plan
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO plan
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE DATABASE plan
< ok 0 0
> exec ^
: string
| SAVEPOINT plan
< ok 0 0
> exec ^
: string
| RELEASE SAVEPOINT plan
< ok 0 0
> exec ^
: string
| SET plan = 1
< ok 0 0
> exec ^
: string
| SET @plan = 1
< ok 0 0
> exec ^
: string
| CREATE VIEW plan AS SELECT 1
< ok 0 0
> exec ^
: string
| SELECT * FROM t ORDER BY plan
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t ORDER BY plan
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t GROUP BY plan HAVING plan > 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t GROUP BY plan HAVING plan > 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE plan
< error
| 1146 42S02 Table 'test.plan' doesn't exist
> exec ^
: string
| ANALYZE TABLE plan
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRUNCATE TABLE plan
< error
| 1146 42S02 Table 'test.plan' doesn't exist
> exec ^
: string
| TRUNCATE TABLE plan
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE enabled (enabled INT)
< ok 0 0
> exec ^
: string
| SELECT enabled FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT enabled FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.enabled FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.enabled FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS enabled
< ok 0 0
> exec ^
: string
| SELECT * FROM t enabled
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t enabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS enabled
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t AS enabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM enabled.t
< error
| 1146 42S02 Table 'enabled.t' doesn't exist
> exec ^
: string
| SELECT * FROM enabled.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t enabled WHERE enabled.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t enabled WHERE enabled.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t enabled SET enabled.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t enabled SET enabled.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO enabled (enabled) VALUES (1)
< error
| 1146 42S02 Table 'test.enabled' doesn't exist
> exec ^
: string
| INSERT INTO enabled (enabled) VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX enabled ON t (enabled)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX enabled ON t (enabled)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX enabled ON t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX enabled ON t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN enabled INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN enabled INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO enabled
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO enabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE DATABASE enabled
< ok 0 0
> exec ^
: string
| SAVEPOINT enabled
< ok 0 0
> exec ^
: string
| RELEASE SAVEPOINT enabled
< ok 0 0
> exec ^
: string
| SET enabled = 1
< ok 0 0
> exec ^
: string
| SET @enabled = 1
< ok 0 0
> exec ^
: string
| CREATE VIEW enabled AS SELECT 1
< ok 0 0
> exec ^
: string
| SELECT * FROM t ORDER BY enabled
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t ORDER BY enabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t GROUP BY enabled HAVING enabled > 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t GROUP BY enabled HAVING enabled > 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE enabled
< error
| 1146 42S02 Table 'test.enabled' doesn't exist
> exec ^
: string
| ANALYZE TABLE enabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRUNCATE TABLE enabled
< error
| 1146 42S02 Table 'test.enabled' doesn't exist
> exec ^
: string
| TRUNCATE TABLE enabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE disabled (disabled INT)
< ok 0 0
> exec ^
: string
| SELECT disabled FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT disabled FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.disabled FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.disabled FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS disabled
< ok 0 0
> exec ^
: string
| SELECT * FROM t disabled
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t disabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS disabled
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t AS disabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM disabled.t
< error
| 1146 42S02 Table 'disabled.t' doesn't exist
> exec ^
: string
| SELECT * FROM disabled.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t disabled WHERE disabled.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t disabled WHERE disabled.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t disabled SET disabled.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t disabled SET disabled.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO disabled (disabled) VALUES (1)
< error
| 1146 42S02 Table 'test.disabled' doesn't exist
> exec ^
: string
| INSERT INTO disabled (disabled) VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX disabled ON t (disabled)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX disabled ON t (disabled)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP INDEX disabled ON t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DROP INDEX disabled ON t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD COLUMN disabled INT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN disabled INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO disabled
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t RENAME COLUMN a TO disabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE DATABASE disabled
< ok 0 0
> exec ^
: string
| SAVEPOINT disabled
< ok 0 0
> exec ^
: string
| RELEASE SAVEPOINT disabled
< ok 0 0
> exec ^
: string
| SET disabled = 1
< ok 0 0
> exec ^
: string
| SET @disabled = 1
< ok 0 0
> exec ^
: string
| CREATE VIEW disabled AS SELECT 1
< ok 0 0
> exec ^
: string
| SELECT * FROM t ORDER BY disabled
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t ORDER BY disabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t GROUP BY disabled HAVING disabled > 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t GROUP BY disabled HAVING disabled > 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE disabled
< error
| 1146 42S02 Table 'test.disabled' doesn't exist
> exec ^
: string
| ANALYZE TABLE disabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRUNCATE TABLE disabled
< error
| 1146 42S02 Table 'test.disabled' doesn't exist
> exec ^
: string
| TRUNCATE TABLE disabled
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
