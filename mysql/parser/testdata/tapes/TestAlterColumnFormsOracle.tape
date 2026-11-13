-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestAlterColumnFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
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
| 1064 42000 …near '.b SET NOT NULL' at line 1
> exec ^
: string
| ALTER TABLE t ALTER COLUMN t.b DROP NOT NULL
< error
| 1064 42000 …near '.b DROP NOT NULL' at line 1
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
| ALTER TABLE t ALTER COLUMN b SET INVISIBLE
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET VISIBLE
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
