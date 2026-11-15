-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestSyncForeignKeyReferencedDatabase default#1
> exec
|
| 		CREATE TABLE omni_fk_remote.parent (id INT PRIMARY KEY);
| 		CREATE TABLE omni_fk_local.parent (id INT PRIMARY KEY);
| 		CREATE TABLE omni_fk_local.child (
| 			id INT PRIMARY KEY,
| 			local_id INT,
| 			qualified_id INT,
| 			remote_id INT,
| 			CONSTRAINT fk_local FOREIGN KEY (local_id) REFERENCES parent (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
| 			CONSTRAINT fk_qualified FOREIGN KEY (qualified_id) REFERENCES omni_fk_local.parent (id) ON DELETE SET NULL ON UPDATE RESTRICT,
| 			CONSTRAINT fk_remote FOREIGN KEY (remote_id) REFERENCES omni_fk_remote.parent (id) ON DELETE CASCADE ON UPDATE RESTRICT
| 		);
| 	
< ok 0 0
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			REFERENCED_TABLE_NAME,
| 			DELETE_RULE,
| 			UPDATE_RULE,
| 			MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
| 		WHERE CONSTRAINT_SCHEMA = ?;
| 	
: string
| omni_fk_local
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
< row
: bytes
| child
: bytes
| fk_local
: bytes
| parent
: bytes
| RESTRICT
: bytes
| RESTRICT
: bytes
| NONE
< row
: bytes
| child
: bytes
| fk_qualified
: bytes
| parent
: bytes
| SET NULL
: bytes
| RESTRICT
: bytes
| NONE
< row
: bytes
| child
: bytes
| fk_remote
: bytes
| parent
: bytes
| CASCADE
: bytes
| RESTRICT
: bytes
| NONE
> query
|
| 		SELECT
| 			TABLE_SCHEMA,
| 			TABLE_NAME,
| 			CONSTRAINT_NAME,
| 			COLUMN_NAME,
| 			REFERENCED_TABLE_SCHEMA,
| 			REFERENCED_COLUMN_NAME
| 		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
| 		WHERE POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION;
| 	
: string
| omni_fk_local
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| omni_fk_local
: bytes
| child
: bytes
| fk_local
: bytes
| local_id
: bytes
| omni_fk_local
: bytes
| id
< row
: bytes
| omni_fk_local
: bytes
| child
: bytes
| fk_qualified
: bytes
| qualified_id
: bytes
| omni_fk_local
: bytes
| id
< row
: bytes
| omni_fk_local
: bytes
| child
: bytes
| fk_remote
: bytes
| remote_id
: bytes
| omni_fk_remote
: bytes
| id
