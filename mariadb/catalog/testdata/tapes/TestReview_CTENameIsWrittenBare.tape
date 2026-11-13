-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_CTENameIsWrittenBare default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE VIEW cv0 AS WITH cte AS (SELECT 1 AS x) SELECT x FROM cte
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv0'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with cte as (select 1 AS `x`)select `cte`.`x` AS `x` from `cte`
> exec
| SELECT * FROM cv0
< ok 0 0
> exec
| SHOW CREATE VIEW cv0
< ok 0 0
> exec
| CREATE VIEW cv1 AS WITH MyCte AS (SELECT 1 AS x) SELECT x FROM MyCte
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv1'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with MyCte as (select 1 AS `x`)select `MyCte`.`x` AS `x` from `MyCte`
> exec
| SELECT * FROM cv1
< ok 0 0
> exec
| SHOW CREATE VIEW cv1
< ok 0 0
> exec
| CREATE VIEW cv2 AS WITH `MyCte` AS (SELECT 1 AS x) SELECT x FROM `MyCte`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv2'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with MyCte as (select 1 AS `x`)select `MyCte`.`x` AS `x` from `MyCte`
> exec
| SELECT * FROM cv2
< ok 0 0
> exec
| SHOW CREATE VIEW cv2
< ok 0 0
> exec
| CREATE VIEW cv3 AS WITH `comment` AS (SELECT 1 AS x) SELECT x FROM `comment`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv3'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with comment as (select 1 AS `x`)select `comment`.`x` AS `x` from `comment`
> exec
| SELECT * FROM cv3
< ok 0 0
> exec
| SHOW CREATE VIEW cv3
< ok 0 0
> exec
| CREATE VIEW cv4 AS WITH `my cte` AS (SELECT 1 AS x) SELECT x FROM `my cte`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv4'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with my cte as (select 1 AS `x`)select `my cte`.`x` AS `x` from `my cte`
> exec
| SELECT * FROM cv4
< error
| 1064 42000 …near 'cte as (select 1 AS `x`)select `my cte`.`x` AS `x` from `my cte`' at line 1
> exec
| SHOW CREATE VIEW cv4
< error
| 1064 42000 …near 'cte as (select 1 AS `x`)select `my cte`.`x` AS `x` from `my cte`' at line 1
> exec
| CREATE VIEW cv5 AS WITH `select` AS (SELECT 1 AS x) SELECT x FROM `select`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv5'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with select as (select 1 AS `x`)select `select`.`x` AS `x` from `select`
> exec
| SELECT * FROM cv5
< error
| 1064 42000 …near 'select as (select 1 AS `x`)select `select`.`x` AS `x` from `select`' at line 1
> exec
| SHOW CREATE VIEW cv5
< error
| 1064 42000 …near 'select as (select 1 AS `x`)select `select`.`x` AS `x` from `select`' at line 1
> exec
| CREATE VIEW cv6 AS WITH `1cte` AS (SELECT 1 AS x) SELECT x FROM `1cte`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv6'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with 1cte as (select 1 AS `x`)select `1cte`.`x` AS `x` from `1cte`
> exec
| SELECT * FROM cv6
< ok 0 0
> exec
| SHOW CREATE VIEW cv6
< ok 0 0
> exec
| CREATE VIEW cv7 AS WITH `a.b` AS (SELECT 1 AS x) SELECT x FROM `a.b`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv7'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with a.b as (select 1 AS `x`)select `a.b`.`x` AS `x` from `a.b`
> exec
| SELECT * FROM cv7
< error
| 1064 42000 …near '.b as (select 1 AS `x`)select `a.b`.`x` AS `x` from `a.b`' at line 1
> exec
| SHOW CREATE VIEW cv7
< error
| 1064 42000 …near '.b as (select 1 AS `x`)select `a.b`.`x` AS `x` from `a.b`' at line 1
> exec
| CREATE VIEW cv8 AS WITH `a-b` AS (SELECT 1 AS x) SELECT x FROM `a-b`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv8'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with a-b as (select 1 AS `x`)select `a-b`.`x` AS `x` from `a-b`
> exec
| SELECT * FROM cv8
< error
| 1064 42000 …near '-b as (select 1 AS `x`)select `a-b`.`x` AS `x` from `a-b`' at line 1
> exec
| SHOW CREATE VIEW cv8
< error
| 1064 42000 …near '-b as (select 1 AS `x`)select `a-b`.`x` AS `x` from `a-b`' at line 1
> exec
| CREATE VIEW cv9 AS WITH `a``b` AS (SELECT 1 AS x) SELECT x FROM `a``b`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv9'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with a`b as (select 1 AS `x`)select `a``b`.`x` AS `x` from `a``b`
> exec
| SELECT * FROM cv9
< error
| 1064 42000 …near '`b as (select 1 AS `x`)select `a``b`.`x` AS `x` from `a``b`' at line 1
> exec
| SHOW CREATE VIEW cv9
< error
| 1064 42000 …near '`b as (select 1 AS `x`)select `a``b`.`x` AS `x` from `a``b`' at line 1
> exec
| CREATE VIEW cv10 AS WITH `my cte`(`x`) AS (SELECT 1) SELECT x FROM `my cte`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv10'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with my cte(`x`) as (select 1 AS `x`)select `my cte`.`x` AS `x` from `my cte`
> exec
| SELECT * FROM cv10
< error
| 1064 42000 …near 'cte(`x`) as (select 1 AS `x`)select `my cte`.`x` AS `x` from `my cte`' at line 1
> exec
| SHOW CREATE VIEW cv10
< error
| 1064 42000 …near 'cte(`x`) as (select 1 AS `x`)select `my cte`.`x` AS `x` from `my cte`' at line 1
> exec
| CREATE VIEW cv11 AS WITH `my cte` AS (SELECT 1 AS x), plain AS (SELECT 2 AS y) SELECT x FROM `my cte`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'cv11'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with my cte as (select 1 AS `x`), plain as (select 2 AS `y`)select `my cte`.`x` AS `x` from `my cte`
> exec
| SELECT * FROM cv11
< error
| 1064 42000 …near 'cte as (select 1 AS `x`), plain as (select 2 AS `y`)select `my cte`.`x` AS `x...' at line 1
> exec
| SHOW CREATE VIEW cv11
< error
| 1064 42000 …near 'cte as (select 1 AS `x`), plain as (select 2 AS `y`)select `my cte`.`x` AS `x...' at line 1
