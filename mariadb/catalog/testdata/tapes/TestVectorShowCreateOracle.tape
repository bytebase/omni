-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestVectorShowCreateOracle default#1
> exec
| CREATE TABLE m (
|   id int(11) NOT NULL,
|   v vector(3) NOT NULL,
|   w VECTOR(3,2) DEFAULT NULL,
|   x VECTOR(1.5) DEFAULT NULL,
|   a int(11) DEFAULT NULL,
|   b text DEFAULT NULL,
|   PRIMARY KEY (id),
|   UNIQUE KEY uk (a),
|   KEY ka (a),
|   FULLTEXT KEY ft (b),
|   VECTOR KEY k (v) M=6 DISTANCE=cosine COMMENT 'c'
| )
< ok 0 0
> query
| SHOW CREATE TABLE m
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| m
: bytes
| CREATE TABLE `m` (
|   `id` int(11) NOT NULL,
|   `v` vector(3) NOT NULL,
|   `w` vector(3) DEFAULT NULL,
|   `x` vector(1) DEFAULT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` text DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`a`),
|   KEY `ka` (`a`),
|   FULLTEXT KEY `ft` (`b`),
|   VECTOR KEY `k` (`v`) COMMENT 'c' `M`=6 `DISTANCE`=cosine
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
