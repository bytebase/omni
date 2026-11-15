-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSyncForeignKeyReferencedDatabase default#1
> exec
| CREATE TABLE omni_fk_remote.parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE omni_fk_local.parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE omni_fk_local.child (
| 			id INT PRIMARY KEY,
| 			local_id INT,
| 			qualified_id INT,
| 			remote_id INT,
| 			CONSTRAINT fk_local FOREIGN KEY (local_id) REFERENCES omni_fk_local.parent (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
| 			CONSTRAINT fk_remote FOREIGN KEY (remote_id) REFERENCES omni_fk_remote.parent (id) ON DELETE CASCADE ON UPDATE RESTRICT
| 		)
< ok 0 0
> query
|
| 		SELECT
| 			kcu.TABLE_SCHEMA,
| 			fks.TABLE_NAME,
| 			fks.CONSTRAINT_NAME,
| 			kcu.COLUMN_NAME,
| 			kcu.REFERENCED_TABLE_SCHEMA,
| 			fks.REFERENCED_TABLE_NAME,
| 			kcu.REFERENCED_COLUMN_NAME,
| 			fks.DELETE_RULE,
| 			fks.UPDATE_RULE,
| 			fks.MATCH_OPTION
| 		FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS fks
| 			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
| 			ON fks.CONSTRAINT_SCHEMA = kcu.TABLE_SCHEMA
| 				AND fks.TABLE_NAME = kcu.TABLE_NAME
| 				AND fks.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
| 		WHERE kcu.POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL AND fks.CONSTRAINT_SCHEMA = ?
| 		ORDER BY fks.TABLE_NAME, fks.CONSTRAINT_NAME, kcu.ORDINAL_POSITION;
| 	
: string
| omni_fk_local
< result
< column "TABLE_SCHEMA" "VARCHAR"
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "MATCH_OPTION" "VARCHAR"
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
| parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| RESTRICT
: bytes
| NONE
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
| parent
: bytes
| id
: bytes
| CASCADE
: bytes
| RESTRICT
: bytes
| NONE
