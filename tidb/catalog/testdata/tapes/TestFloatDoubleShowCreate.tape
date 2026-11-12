-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestFloatDoubleShowCreate default#1
> exec
| DROP DATABASE IF EXISTS zz_float
< ok 0 0
> exec
| CREATE DATABASE zz_float
< ok 0 0
> exec
| USE zz_float
< ok 0 0
> exec
| CREATE TABLE fd (a DOUBLE(10,0), b REAL(10,0), c FLOAT8(10,0), d DOUBLE PRECISION(10,0), e DOUBLE(10,0) UNSIGNED, f DOUBLE(10,2), g DOUBLE, h FLOAT8, i FLOAT(10,0), j FLOAT4(10,0), k FLOAT(10,2), l FLOAT(10), m FLOAT4(5), n FLOAT(24), o FLOAT(25), p FLOAT(30), q FLOAT(0), r FLOAT, w FLOAT(30,0), x FLOAT(25,0), y FLOAT(53), z FLOAT(30,2), s FLOAT(10,0) ZEROFILL, t DOUBLE(10,2) ZEROFILL, u INT ZEROFILL, v DECIMAL(5,0))
< ok 0 0
> query
| SHOW CREATE TABLE fd
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| fd
: bytes
| CREATE TABLE `fd` (
|   `a` double DEFAULT NULL,
|   `b` double DEFAULT NULL,
|   `c` double DEFAULT NULL,
|   `d` double DEFAULT NULL,
|   `e` double unsigned DEFAULT NULL,
|   `f` double(10,2) DEFAULT NULL,
|   `g` double DEFAULT NULL,
|   `h` double DEFAULT NULL,
|   `i` float DEFAULT NULL,
|   `j` float DEFAULT NULL,
|   `k` float(10,2) DEFAULT NULL,
|   `l` float DEFAULT NULL,
|   `m` float DEFAULT NULL,
|   `n` float DEFAULT NULL,
|   `o` double DEFAULT NULL,
|   `p` double DEFAULT NULL,
|   `q` float DEFAULT NULL,
|   `r` float DEFAULT NULL,
|   `w` float DEFAULT NULL,
|   `x` float DEFAULT NULL,
|   `y` double DEFAULT NULL,
|   `z` float(30,2) DEFAULT NULL,
|   `s` float unsigned zerofill DEFAULT NULL,
|   `t` double(10,2) unsigned zerofill DEFAULT NULL,
|   `u` int(10) unsigned zerofill DEFAULT NULL,
|   `v` decimal(5,0) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS zz_float
< ok 0 0
