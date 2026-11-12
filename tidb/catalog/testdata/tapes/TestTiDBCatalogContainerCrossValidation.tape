-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBCatalogContainerCrossValidation default#1
> exec
| DROP DATABASE IF EXISTS testtidbcatalogcontainercrossvalidation
< ok 0 0
> exec
| CREATE DATABASE testtidbcatalogcontainercrossvalidation
< ok 0 0
> exec
| USE testtidbcatalogcontainercrossvalidation
< ok 0 0
> exec
| CREATE TABLE t_ar (id BIGINT AUTO_RANDOM(5) PRIMARY KEY)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ar
< ok 0 0
> exec
| CREATE TABLE t_srib (a INT) SHARD_ROW_ID_BITS = 4
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_srib
< ok 0 0
> exec
| CREATE TABLE t_ttl (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ttl
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_tf (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t_tf SET TIFLASH REPLICA 0
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_tf
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS t_ttl_rm (id INT, created_at DATETIME) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON'
< ok 0 0
> exec
| ALTER TABLE t_ttl_rm REMOVE TTL
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_ttl_rm
< ok 0 0
> exec
| CREATE TABLE t_cl (id INT, PRIMARY KEY (id) CLUSTERED)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_cl
< ok 0 0
> exec
| CREATE TABLE t_mix (id BIGINT AUTO_RANDOM(5) PRIMARY KEY CLUSTERED, created_at DATETIME) AUTO_ID_CACHE = 100 TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'ON' TTL_JOB_INTERVAL = '1h'
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_mix
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testtidbcatalogcontainercrossvalidation
< ok 0 0
