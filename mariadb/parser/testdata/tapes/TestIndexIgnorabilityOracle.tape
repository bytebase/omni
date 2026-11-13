-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestIndexIgnorabilityOracle default#1
> exec
| CREATE TABLE t (a INT, KEY i (a) IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) NOT IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i (a) NOT IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (a INT, INDEX (a) IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, INDEX (a) IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (a INT, UNIQUE KEY u (a) COMMENT 'x' IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, UNIQUE KEY u (a) COMMENT 'x' IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) IGNORED COMMENT 'x')
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i (a) IGNORED COMMENT 'x')
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i USING HASH (a) IGNORED KEY_BLOCK_SIZE=8)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i USING HASH (a) IGNORED KEY_BLOCK_SIZE=8)
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) USING BTREE IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i (a) USING BTREE IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (a INT, b TEXT, FULLTEXT KEY f (b) IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, b TEXT, FULLTEXT KEY f (b) IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (g GEOMETRY NOT NULL, SPATIAL KEY s (g) IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (g GEOMETRY NOT NULL, SPATIAL KEY s (g) IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT c UNIQUE (a) IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, CONSTRAINT c UNIQUE (a) IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) IGNORED NOT IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i (a) IGNORED NOT IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) NOT IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) NOT IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) IGNORED)
< ok 0 0
> exec
| CREATE TABLE t (ignored INT, KEY ignored (ignored) IGNORED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (ignored INT, KEY ignored (ignored) IGNORED)
< ok 0 0
> exec
| CREATE INDEX i ON t (a) IGNORED
< error
| 1061 42000 Duplicate key name 'i'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE INDEX i ON t (a) IGNORED
< ok 0 0
> exec
| CREATE INDEX i ON t (a) NOT IGNORED ALGORITHM=INPLACE
< error
| 1061 42000 Duplicate key name 'i'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE INDEX i ON t (a) NOT IGNORED ALGORITHM=INPLACE
< ok 0 0
> exec
| CREATE UNIQUE INDEX i ON t (a) COMMENT 'u' IGNORED LOCK=NONE
< error
| 1061 42000 Duplicate key name 'i'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE UNIQUE INDEX i ON t (a) COMMENT 'u' IGNORED LOCK=NONE
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX i (a) IGNORED
< error
| 1061 42000 Duplicate key name 'i'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER TABLE t ADD INDEX i (a) IGNORED
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT u UNIQUE (a) NOT IGNORED
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (a) IGNORED
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER TABLE t ADD PRIMARY KEY (a) IGNORED
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) VISIBLE)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i (a) VISIBLE)
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) IGNORED VISIBLE)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i (a) IGNORED VISIBLE)
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) VISIBLE)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) VISIBLE)
< ok 0 0
> exec
| CREATE INDEX i ON t (a) VISIBLE
< error
| 1061 42000 Duplicate key name 'i'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE INDEX i ON t (a) VISIBLE
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX i (a) VISIBLE
< error
| 1061 42000 Duplicate key name 'i'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER TABLE t ADD INDEX i (a) VISIBLE
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX i IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX i NOT IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER KEY i IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX IF EXISTS i NOT IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER KEY IF EXISTS i IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX `PRIMARY` IGNORED
< error
| 1176 42000 Key 'PRIMARY' doesn't exist in table 't'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER TABLE t ALTER INDEX `PRIMARY` IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX ignored NOT IGNORED
< error
| 1176 42000 Key 'ignored' doesn't exist in table 't'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER TABLE t ALTER INDEX ignored NOT IGNORED
< ok 0 0
> exec
| ALTER TABLE t ALTER INDEX i IGNORED, ADD INDEX j (a) IGNORED
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) /*!80000 INVISIBLE */)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i (a) /*!80000 INVISIBLE */)
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) /*!50700 INVISIBLE */)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i (a) /*!50700 INVISIBLE */)
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) /*!99999 INVISIBLE */)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT, KEY i (a) /*!99999 INVISIBLE */)
< ok 0 0
> exec
| CREATE INDEX i ON t (a) /*!80000 INVISIBLE */
< error
| 1061 42000 Duplicate key name 'i'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE INDEX i ON t (a) /*!80000 INVISIBLE */
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX i (a) /*!80000 INVISIBLE */
< error
| 1061 42000 Duplicate key name 'i'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER TABLE t ADD INDEX i (a) /*!80000 INVISIBLE */
< ok 0 0
> exec
| CREATE TABLE t (a INT INVISIBLE, b INT)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (a INT INVISIBLE, b INT)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN e INT INVISIBLE
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY i (a) INVISIBLE)
< error
| 1064 42000 …near 'INVISIBLE)' at line 1
> exec
| CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a) INVISIBLE)
< error
| 1064 42000 …near 'INVISIBLE)' at line 1
> exec
| CREATE INDEX i ON t (a) INVISIBLE
< error
| 1064 42000 …near 'INVISIBLE' at line 1
> exec
| ALTER TABLE t ADD INDEX i (a) INVISIBLE
< error
| 1064 42000 …near 'INVISIBLE' at line 1
> exec
| CREATE TABLE t (a INT, KEY i (a) /*!INVISIBLE */)
< error
| 1064 42000 …near 'INVISIBLE */)' at line 1
> exec
| CREATE TABLE t (a INT, KEY i (a) /*!50699 INVISIBLE */)
< error
| 1064 42000 …near 'INVISIBLE */)' at line 1
> exec
| CREATE TABLE t (a INT, KEY i (a) /*!100000 INVISIBLE */)
< error
| 1064 42000 …near 'INVISIBLE */)' at line 1
> exec
| ALTER TABLE t ALTER INDEX i INVISIBLE
< error
| 1064 42000 …near 'INVISIBLE' at line 1
> exec
| ALTER TABLE t ALTER INDEX i VISIBLE
< error
| 1064 42000 …near 'VISIBLE' at line 1
> exec
| ALTER TABLE t ALTER INDEX i
< error
| 1064 42000 …near '' at line 1
> exec
| ALTER TABLE t ALTER INDEX i NOT
< error
| 1064 42000 …near '' at line 1
> exec
| ALTER TABLE t ALTER INDEX i IGNORED IGNORED
< error
| 1064 42000 …near 'IGNORED' at line 1
> exec
| ALTER TABLE t ALTER INDEX i `IGNORED`
< error
| 1064 42000 …near '`IGNORED`' at line 1
> exec
| ALTER TABLE t ALTER INDEX i /*!80000 INVISIBLE */
< error
| 1064 42000 …near '' at line 1
> exec
| ALTER TABLE t ALTER INDEX primary IGNORED
< error
| 1064 42000 …near 'primary IGNORED' at line 1
> exec
| ALTER TABLE t ALTER COLUMN a IGNORED
< error
| 1064 42000 …near 'IGNORED' at line 1
> exec
| CREATE TABLE t (a INT, KEY i (a) `IGNORED`)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a INT, KEY i (a) NOT)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (a INT, KEY i (a) NOT VISIBLE)
< error
| 1064 42000 …near 'VISIBLE)' at line 1
> exec
| CREATE TABLE t (a INT, KEY i (a) IGNORE)
< error
| 1064 42000 …near 'IGNORE)' at line 1
> exec
| CREATE TABLE t (a INT, KEY i (a) COMMENT '/*!80000' INVISIBLE)
< error
| 1064 42000 …near 'INVISIBLE)' at line 1
> exec
| CREATE TABLE t (a INT UNIQUE IGNORED)
< error
| 1064 42000 …near 'IGNORED)' at line 1
> exec
| CREATE TABLE t (a INT PRIMARY KEY IGNORED)
< error
| 1064 42000 …near 'IGNORED)' at line 1
