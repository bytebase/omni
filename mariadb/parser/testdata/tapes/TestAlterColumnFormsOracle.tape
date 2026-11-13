-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestAlterColumnFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER TABLE t ALTER COLUMN b SET NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t ALTER b SET NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t ALTER b DROP NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec ^
: string
| alter table t alter column b set not null
< error
| 1064 42000 …near 'not null' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET NOT
< error
| 1064 42000 …near 'NOT' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP NOT
< error
| 1064 42000 …near 'NOT' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN t.b SET NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN t.b DROP NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 1, ALTER COLUMN b SET NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP NOT NULL, ALTER COLUMN b DROP DEFAULT
< error
| 1064 42000 …near 'NOT NULL, ALTER COLUMN b DROP DEFAULT' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ALTER b
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 1
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET DEFAULT NULL
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER b DROP DEFAULT
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER b DROP DEFAULT, ALTER b SET DEFAULT 2
< ok 0 0
> exec ^
: string
| ALTER TABLE t MODIFY b INT NOT NULL
< ok 0 0
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN b INT NULL
< ok 0 0
> exec ^
: string
| ALTER TABLE t CHANGE b b INT NOT NULL
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET VISIBLE
< error
| 1064 42000 …near 'VISIBLE' at line 1
> exec ^
: string
| ALTER TABLE t ALTER b SET VISIBLE
< error
| 1064 42000 …near 'VISIBLE' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET INVISIBLE
< error
| 1064 42000 …near 'INVISIBLE' at line 1
> exec ^
: string
| ALTER TABLE t ALTER b SET INVISIBLE
< error
| 1064 42000 …near 'INVISIBLE' at line 1
> exec ^
: string
| alter table t alter column b set invisible
< error
| 1064 42000 …near 'invisible' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 1, ALTER COLUMN b SET INVISIBLE
< error
| 1064 42000 …near 'INVISIBLE' at line 1
> exec ^
: string
| CREATE TABLE t2 (a INT, b INT INVISIBLE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t2 (a INT, b INT VISIBLE)
< error
| 1064 42000 …near 'VISIBLE)' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN c INT INVISIBLE
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN c INT VISIBLE
< error
| 1064 42000 …near 'VISIBLE' at line 1
> exec ^
: string
| ALTER TABLE t MODIFY b INT INVISIBLE
< ok 0 0
> exec ^
: string
| ALTER TABLE t MODIFY b INT VISIBLE
< error
| 1064 42000 …near 'VISIBLE' at line 1
> exec ^
: string
| ALTER TABLE t CHANGE b b INT INVISIBLE
< ok 0 0
> exec ^
: string
| ALTER TABLE t CHANGE b b INT VISIBLE
< error
| 1064 42000 …near 'VISIBLE' at line 1
