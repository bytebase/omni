-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_2_11_CreateDropAlterDatabase default#1
> exec
| DROP DATABASE IF EXISTS db_create1
< ok 0 0
> exec
| CREATE DATABASE db_create1
< ok 0 0
> query
| SHOW CREATE DATABASE db_create1
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_create1
: bytes
| CREATE DATABASE `db_create1` /*!40100 DEFAULT CHARACTER SET utf8mb4 */
> exec
| DROP DATABASE IF EXISTS db_create1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_ine
< ok 0 0
> exec
| CREATE DATABASE db_ine
< ok 0 0
> exec
| CREATE DATABASE IF NOT EXISTS db_ine
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_ine
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_cs
< ok 0 0
> exec
| CREATE DATABASE db_cs CHARACTER SET latin1
< ok 0 0
> query
| SHOW CREATE DATABASE db_cs
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_cs
: bytes
| CREATE DATABASE `db_cs` /*!40100 DEFAULT CHARACTER SET latin1 */
> exec
| DROP DATABASE IF EXISTS db_cs
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_coll
< ok 0 0
> exec
| CREATE DATABASE db_coll COLLATE utf8mb4_unicode_ci
< ok 0 0
> query
| SHOW CREATE DATABASE db_coll
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_coll
: bytes
| CREATE DATABASE `db_coll` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */
> exec
| DROP DATABASE IF EXISTS db_coll
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_drop1
< ok 0 0
> exec
| CREATE DATABASE db_drop1
< ok 0 0
> exec
| DROP DATABASE db_drop1
< ok 0 0
> query
| SHOW CREATE DATABASE db_drop1
< error
| 1049 42000 Unknown database 'db_drop1'
> exec
| DROP DATABASE IF EXISTS db_drop_ine
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_drop_ine
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_alter_cs
< ok 0 0
> exec
| CREATE DATABASE db_alter_cs
< ok 0 0
> exec
| ALTER DATABASE db_alter_cs CHARACTER SET utf8mb4
< ok 0 0
> query
| SHOW CREATE DATABASE db_alter_cs
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_alter_cs
: bytes
| CREATE DATABASE `db_alter_cs` /*!40100 DEFAULT CHARACTER SET utf8mb4 */
> exec
| DROP DATABASE IF EXISTS db_alter_cs
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_alter_coll
< ok 0 0
> exec
| CREATE DATABASE db_alter_coll
< ok 0 0
> exec
| ALTER DATABASE db_alter_coll COLLATE utf8mb4_unicode_ci
< ok 0 0
> query
| SHOW CREATE DATABASE db_alter_coll
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_alter_coll
: bytes
| CREATE DATABASE `db_alter_coll` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */
> exec
| DROP DATABASE IF EXISTS db_alter_coll
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_nonexist_xyz
< ok 0 0
> exec
| DROP DATABASE db_nonexist_xyz
< error
| 1008 HY000 Can't drop database 'db_nonexist_xyz'; database doesn't exist
> exec
| ALTER DATABASE db_nonexist_xyz CHARACTER SET utf8mb4
< error
| 1049 42000 Unknown database 'db_nonexist_xyz'
