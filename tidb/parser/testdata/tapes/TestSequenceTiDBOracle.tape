-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSequenceTiDBOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_seq
< ok 0 0
> exec
| CREATE DATABASE omni_seq
< ok 0 0
> exec
| USE omni_seq
< ok 0 0
> exec
| CREATE SEQUENCE seq_min
< ok 0 0
> exec
| CREATE SEQUENCE IF NOT EXISTS seq_min
< ok 0 0
> exec
| CREATE SEQUENCE seq_min
< error
| 1050 42S01 Table 'omni_seq.seq_min' already exists
> exec
| CREATE SEQUENCE omni_seq.seq_qualified
< ok 0 0
> exec
| CREATE SEQUENCE `seq quoted`
< ok 0 0
> exec
| CREATE SEQUENCE seq_by INCREMENT BY 2 START WITH 10 MINVALUE 1 MAXVALUE 1000 CACHE 20 CYCLE
< ok 0 0
> exec
| CREATE SEQUENCE seq_eq INCREMENT = 2 START = 10 MINVALUE = 1 MAXVALUE = 1000 CACHE = 20 NOCYCLE
< ok 0 0
> exec
| CREATE SEQUENCE seq_bare INCREMENT 2 START 10 MINVALUE 1 MAXVALUE 1000 CACHE 20 NO CYCLE
< ok 0 0
> exec
| CREATE SEQUENCE seq_no1 NOMINVALUE NOMAXVALUE NOCACHE NOCYCLE
< ok 0 0
> exec
| CREATE SEQUENCE seq_no2 NO MINVALUE NO MAXVALUE NO CACHE NO CYCLE
< ok 0 0
> exec
| CREATE SEQUENCE seq_neg INCREMENT BY -1 START WITH -1 MINVALUE -100 MAXVALUE -1
< ok 0 0
> exec
| CREATE SEQUENCE seq_plus INCREMENT BY +1 START WITH +5
< ok 0 0
> exec
| CREATE SEQUENCE seq_repeat START WITH 1 START WITH 2
< ok 0 0
> exec
| CREATE SEQUENCE seq_comment COMMENT 'ids'
< ok 0 0
> exec
| CREATE SEQUENCE seq_tblopts START 1 COMMENT = 'ids', ENGINE = InnoDB
< ok 0 0
> exec
| CREATE SEQUENCE seq_engine ENGINE InnoDB COMMENT 'a' COMMENT 'b'
< ok 0 0
> exec
| CREATE SEQUENCE seq_ai AUTO_INCREMENT = 5
< error
| 8227 HY000 Unsupported sequence table-option 
> exec
| CREATE SEQUENCE seq_cs DEFAULT CHARSET = utf8mb4
< error
| 8227 HY000 Unsupported sequence table-option utf8mb4
> exec
| CREATE SEQUENCE seq_zero INCREMENT BY 0
< error
| 4136 HY000 Sequence 'omni_seq.seq_zero' values are conflicting
> exec
| CREATE SEQUENCE seq_floor MINVALUE -9223372036854775808
< error
| 4136 HY000 Sequence 'omni_seq.seq_floor' values are conflicting
> exec
| CREATE OR REPLACE SEQUENCE seq_x
< error
| 1064 42000 …line 1 column 26 near "SEQUENCE seq_x" 
> exec
| CREATE TEMPORARY SEQUENCE seq_x
< error
| 1064 42000 …line 1 column 25 near "SEQUENCE seq_x" 
> exec
| CREATE SEQUENCE IF EXISTS seq_x
< error
| 1064 42000 …line 1 column 25 near "EXISTS seq_x" 
> exec
| CREATE SEQUENCE seq_x AS BIGINT
< error
| 1064 42000 …line 1 column 24 near "AS BIGINT" 
> exec
| CREATE SEQUENCE seq_x RESTART
< error
| 1064 42000 …line 1 column 29 near "RESTART" 
> exec
| CREATE SEQUENCE seq_x CACHE
< error
| 1064 42000 …line 1 column 27 near "" 
> exec
| CREATE SEQUENCE seq_x MAXVALUE
< error
| 1064 42000 …line 1 column 30 near "" 
> exec
| CREATE SEQUENCE seq_x NO
< error
| 1064 42000 …line 1 column 24 near "" 
> exec
| CREATE SEQUENCE seq_x NO START
< error
| 1064 42000 …line 1 column 30 near "START" 
> exec
| CREATE SEQUENCE seq_x START WITH 1.5
< error
| 1064 42000 …line 1 column 36 near "1.5" 
> exec
| CREATE SEQUENCE seq_x START WITH (1)
< error
| 1064 42000 …line 1 column 34 near "(1)" 
> exec
| CREATE SEQUENCE seq_x START WITH 1+1
< error
| 1064 42000 …line 1 column 35 near "+1" 
> exec
| CREATE SEQUENCE seq_x START WITH '1'
< error
| 1064 42000 …line 1 column 36 near "'1'" 
> exec
| CREATE SEQUENCE seq_x START WITH 0x10
< error
| 1064 42000 …line 1 column 37 near "0x10" 
> exec
| CREATE SEQUENCE seq_x START WITH 9223372036854775808
< error
| 1064 42000 …line 1 column 52 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec
| CREATE SEQUENCE seq_x MINVALUE -9223372036854775809
< error
| 1064 42000 …line 1 column 51 near ""the Signed Value should be at the range of [-9223372036854775808, 9223372036854775807]. 
> exec
| CREATE SEQUENCE seq_x INCREMENT BY = 1
< error
| 1064 42000 …line 1 column 36 near "= 1" 
> exec
| CREATE SEQUENCE seq_x START WITH = 1
< error
| 1064 42000 …line 1 column 34 near "= 1" 
> exec
| CREATE SEQUENCE seq_x COMMENT 'x' START 1
< error
| 1064 42000 …line 1 column 39 near "START 1" 
> exec
| CREATE SEQUENCE seq_x COMMENT 'x',
< error
| 1064 42000 …line 1 column 34 near "" 
> exec
| CREATE SEQUENCE seq_x , COMMENT 'x'
< error
| 1064 42000 …line 1 column 23 near ", COMMENT 'x'" 
> exec
| CREATE SEQUENCE seq_x START TRANSACTION
< error
| 1064 42000 …line 1 column 39 near "TRANSACTION" 
> exec
| CREATE SEQUENCE seq_x PARTITION BY HASH(id)
< error
| 1064 42000 …line 1 column 31 near "PARTITION BY HASH(id)" 
> exec
| ALTER SEQUENCE seq_min RESTART
< ok 0 0
> exec
| ALTER SEQUENCE IF EXISTS seq_min RESTART WITH 5
< ok 0 0
> exec
| ALTER SEQUENCE IF EXISTS seq_missing RESTART
< error
| 1146 42S02 Table 'omni_seq.seq_missing' doesn't exist
> exec
| ALTER SEQUENCE seq_min RESTART = 5
< ok 0 0
> exec
| ALTER SEQUENCE seq_min RESTART 5
< ok 0 0
> exec
| ALTER SEQUENCE seq_neg RESTART -5
< ok 0 0
> exec
| ALTER SEQUENCE seq_min RESTART RESTART WITH 3
< ok 0 0
> exec
| ALTER SEQUENCE seq_min INCREMENT BY 3 CACHE 100 NOCYCLE
< ok 0 0
> exec
| ALTER SEQUENCE seq_min START WITH 3 RESTART
< ok 0 0
> exec
| ALTER SEQUENCE seq_min NO MINVALUE NO MAXVALUE NO CACHE NO CYCLE
< ok 0 0
> exec
| ALTER SEQUENCE seq_min
< error
| 1064 42000 …line 1 column 22 near "" 
> exec
| ALTER SEQUENCE IF NOT EXISTS seq_min RESTART
< error
| 1064 42000 …line 1 column 21 near "NOT EXISTS seq_min RESTART" 
> exec
| ALTER SEQUENCE seq_min RESTART WITH
< error
| 1064 42000 …line 1 column 35 near "" 
> exec
| ALTER SEQUENCE seq_min RESTART WITH = 5
< error
| 1064 42000 …line 1 column 37 near "= 5" 
> exec
| ALTER SEQUENCE seq_min RESTART, CACHE 5
< error
| 1064 42000 …line 1 column 31 near ", CACHE 5" 
> exec
| ALTER SEQUENCE seq_min COMMENT 'x'
< error
| 1064 42000 …line 1 column 30 near "COMMENT 'x'" 
> exec
| ALTER SEQUENCE seq_min ENGINE = InnoDB
< error
| 1064 42000 …line 1 column 29 near "ENGINE = InnoDB" 
> exec
| ALTER SEQUENCE seq_min AS INT
< error
| 1064 42000 …line 1 column 25 near "AS INT" 
> exec
| DROP SEQUENCE seq_by
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS seq_by
< ok 0 0
> exec
| DROP SEQUENCE seq_eq, seq_bare, omni_seq.seq_no1
< ok 0 0
> exec
| DROP SEQUENCE IF EXISTS seq_no2, seq_missing
< ok 0 0
> exec
| DROP SEQUENCE seq_missing
< error
| 4139 HY000 Unknown SEQUENCE: 'omni_seq.seq_missing'
> exec
| DROP SEQUENCE
< error
| 1064 42000 …line 1 column 13 near "" 
> exec
| DROP SEQUENCE seq_min,
< error
| 1064 42000 …line 1 column 22 near "" 
> exec
| DROP TEMPORARY SEQUENCE seq_min
< error
| 1064 42000 …line 1 column 23 near "SEQUENCE seq_min" 
> exec
| DROP SEQUENCE seq_min CASCADE
< error
| 1064 42000 …line 1 column 29 near "CASCADE" 
> exec
| DROP SEQUENCE seq_min RESTRICT
< error
| 1064 42000 …line 1 column 30 near "RESTRICT" 
> exec
| DROP SEQUENCE IF NOT EXISTS seq_min
< error
| 1064 42000 …line 1 column 20 near "NOT EXISTS seq_min" 
> exec
| CREATE SEQUENCE sequence
< ok 0 0
> exec
| CREATE SEQUENCE nocache
< ok 0 0
> exec
| ALTER SEQUENCE nocache RESTART
< ok 0 0
> exec
| DROP SEQUENCE sequence, nocache
< ok 0 0
> exec
| CREATE TABLE cycle (increment INT, minvalue INT, nominvalue INT, nomaxvalue INT, nocache INT, nocycle INT, sequence INT)
< ok 0 0
> exec
| SELECT increment, minvalue, nominvalue, nomaxvalue, nocache, nocycle, sequence FROM cycle
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_seq
< ok 0 0
