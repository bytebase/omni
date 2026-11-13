-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBReturningOracle default#1
> exec
| CREATE TABLE IF NOT EXISTS rt (id INT PRIMARY KEY, name VARCHAR(20))
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS rt2 (id INT, name VARCHAR(20))
< ok 0 0
> exec
| INSERT INTO rt (id, name) VALUES (1, 'a') RETURNING id
< ok 0 0
> exec
| INSERT INTO rt (id, name) VALUES (2, 'b') RETURNING id, name, id + 100 AS big
< ok 0 0
> exec
| INSERT INTO rt (id, name) VALUES (3, 'c') RETURNING *
< ok 0 0
> exec
| INSERT INTO rt SELECT 9, 'q' RETURNING id, name
< ok 0 0
> exec
| INSERT INTO rt (SELECT 1, 'x') RETURNING id
< error
| 1062 23000 Duplicate entry '1' for key 'PRIMARY'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| INSERT INTO rt (SELECT 1, 'x') RETURNING id
< ok 0 0
> exec
| INSERT INTO rt SET id = 1, name = 'x' RETURNING id
< error
| 1062 23000 Duplicate entry '1' for key 'PRIMARY'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| INSERT INTO rt SET id = 1, name = 'x' RETURNING id
< ok 0 0
> exec
| INSERT INTO rt (id, name) VALUES (4, 'd') ON DUPLICATE KEY UPDATE name = 'x' RETURNING id
< ok 0 0
> exec
| INSERT INTO rt (id) VALUES (1) RETURNING id, (SELECT COUNT(*) FROM rt2)
< error
| 1062 23000 Duplicate entry '1' for key 'PRIMARY'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| INSERT INTO rt (id) VALUES (1) RETURNING id, (SELECT COUNT(*) FROM rt2)
< ok 0 0
> exec
| INSERT INTO rt (id) VALUES (1) RETURNING COUNT(*)
< error
| 1111 HY000 Invalid use of group function
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| INSERT INTO rt (id) VALUES (1) RETURNING COUNT(*)
< error
| 1111 HY000 Invalid use of group function
> exec ^
: string
| INSERT INTO rt (id) VALUES (1) RETURNING COUNT(*)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| REPLACE INTO rt (id, name) VALUES (1, 'z') RETURNING id, name
< ok 0 0
> exec
| REPLACE INTO rt (id, name) VALUES (1, 'z') RETURNING *
< ok 0 0
> exec
| DELETE FROM rt WHERE id = 3 RETURNING id, name
< ok 0 0
> exec
| DELETE FROM rt WHERE id = 2 RETURNING *
< ok 0 0
> exec
| DELETE FROM rt RETURNING COUNT(*)
< error
| 1111 HY000 Invalid use of group function
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| DELETE FROM rt RETURNING COUNT(*)
< error
| 1111 HY000 Invalid use of group function
> exec ^
: string
| DELETE FROM rt RETURNING COUNT(*)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| UPDATE rt SET name = 'x' WHERE id = 1 RETURNING id
< error
| 1064 42000 …near 'RETURNING id' at line 1
> exec
| DELETE rt FROM rt JOIN rt2 ON rt.id = rt2.id RETURNING rt.id
< error
| 1064 42000 …near 'RETURNING rt.id' at line 1
> exec
| DELETE FROM rt USING rt JOIN rt2 ON rt.id = rt2.id RETURNING rt.id
< error
| 1064 42000 …near 'RETURNING rt.id' at line 1
> exec
| INSERT INTO rt (id) VALUES (1) RETURNING
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE TABLE rret (returning INT)
< error
| 1064 42000 …near 'returning INT)' at line 1
> exec
| SELECT returning FROM rt
< error
| 1064 42000 …near 'returning FROM rt' at line 1
> exec
| SELECT 1 RETURNING
< error
| 1064 42000 …near 'RETURNING' at line 1
> exec
| SELECT 1 AS RETURNING
< error
| 1064 42000 …near 'RETURNING' at line 1
> exec
| SELECT 1 RETURNING x
< error
| 1064 42000 …near 'RETURNING x' at line 1
> exec
| INSERT INTO rt VALUES (33, 'a') AS new
< error
| 1064 42000 …near 'AS new' at line 1
> exec
| INSERT INTO rt VALUES (34, 'b') AS new ON DUPLICATE KEY UPDATE id = new.id
< error
| 1064 42000 …near 'AS new ON DUPLICATE KEY UPDATE id = new.id' at line 1
> exec
| INSERT INTO rt VALUES (33, 'a') AS new RETURNING id
< error
| 1064 42000 …near 'AS new RETURNING id' at line 1
> exec
| INSERT INTO rt VALUES (34, 'b') AS new ON DUPLICATE KEY UPDATE id = new.id RETURNING id
< error
| 1064 42000 …near 'AS new ON DUPLICATE KEY UPDATE id = new.id RETURNING id' at line 1
