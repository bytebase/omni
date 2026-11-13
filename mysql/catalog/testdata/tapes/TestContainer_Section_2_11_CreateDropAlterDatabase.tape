-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_2_11_CreateDropAlterDatabase default#1
> exec
| DROP DATABASE IF EXISTS db_create1
< ok 0 0
> exec
| CREATE DATABASE db_create1
< ok 1 0
> query
| SHOW CREATE DATABASE db_create1
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_create1
: bytes
| CREATE DATABASE `db_create1` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */
> exec
| DROP DATABASE IF EXISTS db_create1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_ine
< ok 0 0
> exec
| CREATE DATABASE db_ine
< ok 1 0
> exec
| CREATE DATABASE IF NOT EXISTS db_ine
< ok 1 0
> exec
| DROP DATABASE IF EXISTS db_ine
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_cs
< ok 0 0
> exec
| CREATE DATABASE db_cs CHARACTER SET latin1
< ok 1 0
> query
| SHOW CREATE DATABASE db_cs
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_cs
: bytes
| CREATE DATABASE `db_cs` /*!40100 DEFAULT CHARACTER SET latin1 */ /*!80016 DEFAULT ENCRYPTION='N' */
> exec
| DROP DATABASE IF EXISTS db_cs
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_coll
< ok 0 0
> exec
| CREATE DATABASE db_coll COLLATE utf8mb4_unicode_ci
< ok 1 0
> query
| SHOW CREATE DATABASE db_coll
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_coll
: bytes
| CREATE DATABASE `db_coll` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */
> exec
| DROP DATABASE IF EXISTS db_coll
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_drop1
< ok 0 0
> exec
| CREATE DATABASE db_drop1
< ok 1 0
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
> exec ^
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_alter_cs
< ok 0 0
> exec
| CREATE DATABASE db_alter_cs
< ok 1 0
> exec
| ALTER DATABASE db_alter_cs CHARACTER SET utf8mb4
< ok 1 0
> query
| SHOW CREATE DATABASE db_alter_cs
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_alter_cs
: bytes
| CREATE DATABASE `db_alter_cs` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */
> exec
| DROP DATABASE IF EXISTS db_alter_cs
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_alter_coll
< ok 0 0
> exec
| CREATE DATABASE db_alter_coll
< ok 1 0
> exec
| ALTER DATABASE db_alter_coll COLLATE utf8mb4_unicode_ci
< ok 1 0
> query
| SHOW CREATE DATABASE db_alter_coll
< result
< column "Database" "VARCHAR"
< column "Create Database" "VARCHAR"
< row
: bytes
| db_alter_coll
: bytes
| CREATE DATABASE `db_alter_coll` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */
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
| 3503 42Y07 Database 'db_nonexist_xyz' doesn't exist
