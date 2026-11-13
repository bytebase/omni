-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestCommentValueTiDBOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_comment_value
< ok 0 0
> exec
| CREATE DATABASE omni_comment_value
< ok 0 0
> exec
| USE omni_comment_value
< ok 0 0
> exec
| CREATE TABLE ta (a INT)
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT)'
< error
| 1064 42000 …line 1 column 30 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT, b INT)'
< error
| 1064 42000 …line 1 column 30 near ", b INT)" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT NOT NULL)'
< error
| 1064 42000 …line 1 column 33 near "NOT NULL)" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT DEFAULT 1)'
< error
| 1064 42000 …line 1 column 37 near "DEFAULT 1)" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT = ''x'')'
< error
| 1064 42000 …line 1 column 31 near "= 'x')" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT AS (1) COMMENT)'
< error
| 1064 42000 …line 1 column 37 near ")" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD COLUMN b INT COMMENT'
< error
| 1064 42000 …line 1 column 39 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta MODIFY a INT COMMENT'
< error
| 1064 42000 …line 1 column 35 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta CHANGE a b INT COMMENT'
< error
| 1064 42000 …line 1 column 37 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE INDEX i ON ta (a) COMMENT'
< error
| 1064 42000 …line 1 column 32 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) COMMENT)'
< error
| 1064 42000 …line 1 column 41 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) COMMENT, b INT)'
< error
| 1064 42000 …line 1 column 41 near ", b INT)" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, UNIQUE KEY i (a) COMMENT)'
< error
| 1064 42000 …line 1 column 48 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, PRIMARY KEY (a) COMMENT)'
< error
| 1064 42000 …line 1 column 47 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, CONSTRAINT c UNIQUE (a) COMMENT)'
< error
| 1064 42000 …line 1 column 55 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a TEXT, FULLTEXT KEY i (a) COMMENT)'
< error
| 1064 42000 …line 1 column 51 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) KEY_BLOCK_SIZE =)'
< error
| 1064 42000 …line 1 column 50 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) USING)'
< error
| 1064 42000 …line 1 column 39 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a TEXT, FULLTEXT KEY i (a) WITH PARSER)'
< error
| 1064 42000 …line 1 column 55 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) ENGINE_ATTRIBUTE)'
< error
| 1064 42000 …line 1 column 49 near "ENGINE_ATTRIBUTE)" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) SECONDARY_ENGINE_ATTRIBUTE =)'
< error
| 1064 42000 …line 1 column 59 near "SECONDARY_ENGINE_ATTRIBUTE =)" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD INDEX i (a) COMMENT'
< error
| 1064 42000 …line 1 column 38 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD UNIQUE (a) COMMENT'
< error
| 1064 42000 …line 1 column 37 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD PRIMARY KEY (a) COMMENT'
< error
| 1064 42000 …line 1 column 42 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD INDEX i (a) COMMENT, ADD COLUMN b INT'
< error
| 1064 42000 …line 1 column 39 near ", ADD COLUMN b INT" 
> exec
| PREPARE parse_probe FROM 'CREATE USER u COMMENT'
< error
| 1064 42000 …line 1 column 21 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER USER u COMMENT'
< error
| 1064 42000 …line 1 column 20 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE USER u ATTRIBUTE'
< error
| 1064 42000 …line 1 column 23 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER USER u ATTRIBUTE'
< error
| 1064 42000 …line 1 column 22 near "" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT COMMENT ''x'', b INT)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT NOT NULL COMMENT "x" DEFAULT 1)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT AS (1) COMMENT ''x'')'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta MODIFY a INT COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD COLUMN b INT COMMENT ''x'' FIRST'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE INDEX i ON ta (a) COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, KEY i (a) COMMENT ''x'' USING BTREE KEY_BLOCK_SIZE = 8)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, PRIMARY KEY (a) COMMENT ''x'', UNIQUE (a) COMMENT ''y'')'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a TEXT, FULLTEXT KEY i (a) WITH PARSER ngram)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE ta ADD INDEX i (a) COMMENT ''x'', ADD COLUMN b INT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE USER u COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER USER u COMMENT ''x'''
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE USER u ATTRIBUTE ''{"a": 1}'''
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_comment_value
< ok 0 0
