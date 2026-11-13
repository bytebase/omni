-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_6_3_InformationSchemaConsistency default#1
> exec
| DROP TABLE IF EXISTS t_is_test
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_is_parent
< ok 0 0
> exec
| CREATE TABLE t_is_parent (id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_is_test (
| 		id INT NOT NULL AUTO_INCREMENT,
| 		name VARCHAR(100) NOT NULL DEFAULT '',
| 		price DECIMAL(10,2) DEFAULT '0.00',
| 		status ENUM('active','inactive') DEFAULT 'active',
| 		parent_id INT,
| 		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
| 		PRIMARY KEY (id),
| 		UNIQUE KEY uk_name (name),
| 		KEY idx_status (status),
| 		KEY idx_parent (parent_id),
| 		CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES t_is_parent (id) ON DELETE SET NULL,
| 		CONSTRAINT chk_price CHECK (price >= 0)
| 	) ENGINE=InnoDB
< ok 0 0
> query
|
| 		SELECT COLUMN_NAME, ORDINAL_POSITION, DATA_TYPE, COLUMN_TYPE,
| 		       IS_NULLABLE, COLUMN_DEFAULT, COLUMN_KEY, EXTRA,
| 		       CHARACTER_SET_NAME, COLLATION_NAME,
| 		       CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE
| 		FROM INFORMATION_SCHEMA.COLUMNS
| 		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
| 		ORDER BY ORDINAL_POSITION
: string
| catalog_shared
: string
| t_is_test
< result
< column "COLUMN_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "UNSIGNED BIGINT"
< column "DATA_TYPE" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_DEFAULT" "TEXT"
< column "COLUMN_KEY" "VARCHAR"
< column "EXTRA" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< column "CHARACTER_MAXIMUM_LENGTH" "UNSIGNED BIGINT"
< column "NUMERIC_PRECISION" "UNSIGNED BIGINT"
< column "NUMERIC_SCALE" "UNSIGNED BIGINT"
< row
: bytes
| id
: uint64 1
: bytes
| int
: bytes
| int(11)
: bytes
| NO
: null
: bytes
| PRI
: bytes
| auto_increment
: null
: null
: null
: uint64 10
: uint64 0
< row
: bytes
| name
: uint64 2
: bytes
| varchar
: bytes
| varchar(100)
: bytes
| NO
: bytes
| ''
: bytes
| UNI
: bytes
|
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: uint64 100
: null
: null
< row
: bytes
| price
: uint64 3
: bytes
| decimal
: bytes
| decimal(10,2)
: bytes
| YES
: bytes
| 0.00
: bytes
|
: bytes
|
: null
: null
: null
: uint64 10
: uint64 2
< row
: bytes
| status
: uint64 4
: bytes
| enum
: bytes
| enum('active','inactive')
: bytes
| YES
: bytes
| 'active'
: bytes
| MUL
: bytes
|
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: uint64 8
: null
: null
< row
: bytes
| parent_id
: uint64 5
: bytes
| int
: bytes
| int(11)
: bytes
| YES
: bytes
| NULL
: bytes
| MUL
: bytes
|
: null
: null
: null
: uint64 10
: uint64 0
< row
: bytes
| created_at
: uint64 6
: bytes
| timestamp
: bytes
| timestamp
: bytes
| YES
: bytes
| current_timestamp()
: bytes
|
: bytes
|
: null
: null
: null
: null
: null
> query
|
| 		SELECT INDEX_NAME, SEQ_IN_INDEX, COLUMN_NAME, COLLATION,
| 		       NON_UNIQUE, INDEX_TYPE, NULLABLE
| 		FROM INFORMATION_SCHEMA.STATISTICS
| 		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
| 		ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| catalog_shared
: string
| t_is_test
< result
< column "INDEX_NAME" "VARCHAR"
< column "SEQ_IN_INDEX" "UNSIGNED INT"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_TYPE" "VARCHAR"
< column "NULLABLE" "VARCHAR"
< row
: bytes
| idx_parent
: int64 1
: bytes
| parent_id
: bytes
| A
: int64 1
: bytes
| BTREE
: bytes
| YES
< row
: bytes
| idx_status
: int64 1
: bytes
| status
: bytes
| A
: int64 1
: bytes
| BTREE
: bytes
| YES
< row
: bytes
| PRIMARY
: int64 1
: bytes
| id
: bytes
| A
: int64 0
: bytes
| BTREE
: bytes
|
< row
: bytes
| uk_name
: int64 1
: bytes
| name
: bytes
| A
: int64 0
: bytes
| BTREE
: bytes
|
> query
|
| 		SELECT CONSTRAINT_NAME, CONSTRAINT_TYPE
| 		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
| 		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
| 		ORDER BY CONSTRAINT_NAME
: string
| catalog_shared
: string
| t_is_test
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< row
: bytes
| chk_price
: bytes
| CHECK
< row
: bytes
| fk_parent
: bytes
| FOREIGN KEY
< row
: bytes
| PRIMARY
: bytes
| PRIMARY KEY
< row
: bytes
| uk_name
: bytes
| UNIQUE
> query
|
| 			SELECT CONSTRAINT_NAME, COLUMN_NAME, ORDINAL_POSITION,
| 			       REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
| 			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
| 			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't_is_test'
| 			ORDER BY CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "ORDINAL_POSITION" "BIGINT"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| fk_parent
: bytes
| parent_id
: int64 1
: bytes
| t_is_parent
: bytes
| id
< row
: bytes
| PRIMARY
: bytes
| id
: int64 1
: null
: null
< row
: bytes
| uk_name
: bytes
| name
: int64 1
: null
: null
> query
|
| 			SELECT CONSTRAINT_NAME, UNIQUE_CONSTRAINT_NAME,
| 			       MATCH_OPTION, UPDATE_RULE, DELETE_RULE,
| 			       REFERENCED_TABLE_NAME
| 			FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
| 			WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 't_is_test'
| 			ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "UNIQUE_CONSTRAINT_NAME" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| NONE
: bytes
| RESTRICT
: bytes
| SET NULL
: bytes
| t_is_parent
> query
|
| 			SELECT CONSTRAINT_NAME, CHECK_CLAUSE
| 			FROM INFORMATION_SCHEMA.CHECK_CONSTRAINTS
| 			WHERE CONSTRAINT_SCHEMA = DATABASE()
| 			ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| chk_price
: bytes
| `price` >= 0
