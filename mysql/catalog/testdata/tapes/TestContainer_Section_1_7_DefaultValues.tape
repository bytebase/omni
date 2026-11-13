-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_1_7_DefaultValues default#1
> exec
| DROP TABLE IF EXISTS t_def_int0
< ok 0 0
> exec
| CREATE TABLE t_def_int0 (a INT DEFAULT 0)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_int0
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_int0
: bytes
| CREATE TABLE `t_def_int0` (
|   `a` int DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_intn
< ok 0 0
> exec
| CREATE TABLE t_def_intn (a INT DEFAULT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_intn
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_intn
: bytes
| CREATE TABLE `t_def_intn` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_intnn
< ok 0 0
> exec
| CREATE TABLE t_def_intnn (a INT NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_intnn
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_intnn
: bytes
| CREATE TABLE `t_def_intnn` (
|   `a` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_vch
< ok 0 0
> exec
| CREATE TABLE t_def_vch (a VARCHAR(50) DEFAULT 'hello')
< ok 0 0
> query
| SHOW CREATE TABLE t_def_vch
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_vch
: bytes
| CREATE TABLE `t_def_vch` (
|   `a` varchar(50) DEFAULT 'hello'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_vce
< ok 0 0
> exec
| CREATE TABLE t_def_vce (a VARCHAR(50) DEFAULT '')
< ok 0 0
> query
| SHOW CREATE TABLE t_def_vce
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_vce
: bytes
| CREATE TABLE `t_def_vce` (
|   `a` varchar(50) DEFAULT ''
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_flt
< ok 0 0
> exec
| CREATE TABLE t_def_flt (a FLOAT DEFAULT 3.14)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_flt
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_flt
: bytes
| CREATE TABLE `t_def_flt` (
|   `a` float DEFAULT '3.14'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_dec
< ok 0 0
> exec
| CREATE TABLE t_def_dec (a DECIMAL(10,2) DEFAULT 0.00)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_dec
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_dec
: bytes
| CREATE TABLE `t_def_dec` (
|   `a` decimal(10,2) DEFAULT '0.00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_bt
< ok 0 0
> exec
| CREATE TABLE t_def_bt (a BOOLEAN DEFAULT TRUE)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_bt
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_bt
: bytes
| CREATE TABLE `t_def_bt` (
|   `a` tinyint(1) DEFAULT '1'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_bf
< ok 0 0
> exec
| CREATE TABLE t_def_bf (a BOOLEAN DEFAULT FALSE)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_bf
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_bf
: bytes
| CREATE TABLE `t_def_bf` (
|   `a` tinyint(1) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_enum
< ok 0 0
> exec
| CREATE TABLE t_def_enum (a ENUM('a','b','c') DEFAULT 'a')
< ok 0 0
> query
| SHOW CREATE TABLE t_def_enum
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_enum
: bytes
| CREATE TABLE `t_def_enum` (
|   `a` enum('a','b','c') DEFAULT 'a'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_set
< ok 0 0
> exec
| CREATE TABLE t_def_set (a SET('x','y','z') DEFAULT 'x,y')
< ok 0 0
> query
| SHOW CREATE TABLE t_def_set
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_set
: bytes
| CREATE TABLE `t_def_set` (
|   `a` set('x','y','z') DEFAULT 'x,y'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_bit
< ok 0 0
> exec
| CREATE TABLE t_def_bit (a BIT(8) DEFAULT b'00001111')
< ok 0 0
> query
| SHOW CREATE TABLE t_def_bit
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_bit
: bytes
| CREATE TABLE `t_def_bit` (
|   `a` bit(8) DEFAULT b'1111'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_blob
< ok 0 0
> exec
| CREATE TABLE t_def_blob (a BLOB)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_blob
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_blob
: bytes
| CREATE TABLE `t_def_blob` (
|   `a` blob
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_text
< ok 0 0
> exec
| CREATE TABLE t_def_text (a TEXT)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_text
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_text
: bytes
| CREATE TABLE `t_def_text` (
|   `a` text
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_json
< ok 0 0
> exec
| CREATE TABLE t_def_json (a JSON)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_json
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_json
: bytes
| CREATE TABLE `t_def_json` (
|   `a` json DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_tsct
< ok 0 0
> exec
| CREATE TABLE t_def_tsct (a TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_tsct
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_tsct
: bytes
| CREATE TABLE `t_def_tsct` (
|   `a` timestamp NULL DEFAULT CURRENT_TIMESTAMP
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_dtct
< ok 0 0
> exec
| CREATE TABLE t_def_dtct (a DATETIME DEFAULT CURRENT_TIMESTAMP)
< ok 0 0
> query
| SHOW CREATE TABLE t_def_dtct
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_dtct
: bytes
| CREATE TABLE `t_def_dtct` (
|   `a` datetime DEFAULT CURRENT_TIMESTAMP
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_ts3
< ok 0 0
> exec
| CREATE TABLE t_def_ts3 (a TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3))
< ok 0 0
> query
| SHOW CREATE TABLE t_def_ts3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_ts3
: bytes
| CREATE TABLE `t_def_ts3` (
|   `a` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_expr
< ok 0 0
> exec
| CREATE TABLE t_def_expr (a INT DEFAULT (FLOOR(RAND()*100)))
< ok 0 0
> query
| SHOW CREATE TABLE t_def_expr
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_expr
: bytes
| CREATE TABLE `t_def_expr` (
|   `a` int DEFAULT (floor((rand() * 100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_exjson
< ok 0 0
> exec
| CREATE TABLE t_def_exjson (a JSON DEFAULT (JSON_ARRAY()))
< ok 0 0
> query
| SHOW CREATE TABLE t_def_exjson
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_exjson
: bytes
| CREATE TABLE `t_def_exjson` (
|   `a` json DEFAULT (json_array())
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_exvc
< ok 0 0
> exec
| CREATE TABLE t_def_exvc (a VARCHAR(36) DEFAULT (UUID()))
< ok 0 0
> query
| SHOW CREATE TABLE t_def_exvc
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_exvc
: bytes
| CREATE TABLE `t_def_exvc` (
|   `a` varchar(36) DEFAULT (uuid())
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_dtlit
< ok 0 0
> exec
| CREATE TABLE t_def_dtlit (a DATETIME DEFAULT '2024-01-01 00:00:00')
< ok 0 0
> query
| SHOW CREATE TABLE t_def_dtlit
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_dtlit
: bytes
| CREATE TABLE `t_def_dtlit` (
|   `a` datetime DEFAULT '2024-01-01 00:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_def_dtlit2
< ok 0 0
> exec
| CREATE TABLE t_def_dtlit2 (a DATE DEFAULT '2024-01-01')
< ok 0 0
> query
| SHOW CREATE TABLE t_def_dtlit2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_def_dtlit2
: bytes
| CREATE TABLE `t_def_dtlit2` (
|   `a` date DEFAULT '2024-01-01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
