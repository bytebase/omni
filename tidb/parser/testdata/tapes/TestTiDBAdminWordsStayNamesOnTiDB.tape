-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBAdminWordsStayNamesOnTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SELECT * FROM recommend
< error
| 1146 42S02 Table 'test.recommend' doesn't exist
> exec ^
: string
| SELECT * FROM recommend
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM pause
< error
| 1146 42S02 Table 'test.pause' doesn't exist
> exec ^
: string
| SELECT * FROM pause
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM cancel
< error
| 1146 42S02 Table 'test.cancel' doesn't exist
> exec ^
: string
| SELECT * FROM cancel
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM br
< error
| 1146 42S02 Table 'test.br' doesn't exist
> exec ^
: string
| SELECT * FROM br
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM apply
< error
| 1146 42S02 Table 'test.apply' doesn't exist
> exec ^
: string
| SELECT * FROM apply
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM metadata
< error
| 1146 42S02 Table 'test.metadata' doesn't exist
> exec ^
: string
| SELECT * FROM metadata
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM tidb
< error
| 1146 42S02 Table 'test.tidb' doesn't exist
> exec ^
: string
| SELECT * FROM tidb
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM hosts
< error
| 1146 42S02 Table 'test.hosts' doesn't exist
> exec ^
: string
| SELECT * FROM hosts
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM config
< error
| 1146 42S02 Table 'test.config' doesn't exist
> exec ^
: string
| SELECT * FROM config
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM resume
< error
| 1146 42S02 Table 'test.resume' doesn't exist
> exec ^
: string
| SELECT * FROM resume
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM purge
< error
| 1146 42S02 Table 'test.purge' doesn't exist
> exec ^
: string
| SELECT * FROM purge
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM plugins
< error
| 1146 42S02 Table 'test.plugins' doesn't exist
> exec ^
: string
| SELECT * FROM plugins
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM d.recommend
< error
| 1146 42S02 Table 'd.recommend' doesn't exist
> exec ^
: string
| SELECT * FROM d.recommend
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t recommend WHERE recommend.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t recommend WHERE recommend.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t pause WHERE pause.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t pause WHERE pause.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t cancel WHERE cancel.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t cancel WHERE cancel.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t br WHERE br.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t br WHERE br.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t tidb WHERE tidb.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t tidb WHERE tidb.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t recommend SET recommend.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t recommend SET recommend.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT cancel.a FROM t cancel
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT cancel.a FROM t cancel
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS recommend
< ok 0 0
> exec ^
: string
| SELECT 1 AS tidb
< ok 0 0
> exec ^
: string
| SELECT recommend, pause, cancel, br, apply, metadata, tidb, hosts FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT recommend, pause, cancel, br, apply, metadata, tidb, hosts FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.recommend, t.tidb FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.recommend, t.tidb FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE recommend (pause int, cancel int, br int, apply int, metadata int, tidb int)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (recommend int, INDEX cancel (recommend))
< ok 0 0
> exec ^
: string
| CREATE INDEX pause ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX pause ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP TABLE cancel
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN recommend int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN recommend int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE br RENAME TO tidb
< ok 0 0
> exec ^
: string
| SET config = 1
< ok 0 0
> exec ^
: string
| SET config := 1
< ok 0 0
> exec ^
: string
| SET config.a = 1
< ok 0 0
> exec ^
: string
| SET GLOBAL config = 1
< ok 0 0
> exec ^
: string
| SET @@config = 1
< ok 0 0
> exec ^
: string
| KILL `tidb`(1)
< ok 0 0
> exec ^
: string
| SET CONFIG config a = 1
< ok 0 0
> exec ^
: string
| FLUSH TIDB PLUGINS tidb, plugins, status
< ok 0 0
> exec ^
: string
| RECOMMEND INDEX RUN WITH recommend = 1, run = 2, apply = 3
< ok 0 0
> exec ^
: string
| ALTER RANGE placement PLACEMENT POLICY policy
< ok 0 0
> exec ^
: string
| SELECT 1.5, .5, 1e3 FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1.5, .5, 1e3 FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a.b FROM a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| SELECT a.b FROM a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.* FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.* FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.*, 1 FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.*, 1 FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
