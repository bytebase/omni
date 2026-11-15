-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
< column "DELETE_RULE" "ENUM"
< column "UPDATE_RULE" "ENUM"
< column "MATCH_OPTION" "ENUM"
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
