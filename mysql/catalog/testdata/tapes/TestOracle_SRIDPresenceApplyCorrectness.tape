-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_SRIDPresenceApplyCorrectness default#1

== TestOracle_SRIDPresenceApplyCorrectness/8.0/add-column-with-srid-zero default#1
> exec
| DROP DATABASE IF EXISTS pa_add_column_with_srid_zero_a5e965a1_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_with_srid_zero_a5e965a1_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_with_srid_zero_a5e965a1_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_with_srid_zero_a5e965a1_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_column_with_srid_zero_a5e965a1_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_with_srid_zero_a5e965a1_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_with_srid_zero_a5e965a1_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_with_srid_zero_a5e965a1_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_column_with_srid_zero_a5e965a1
< ok 0 0
> exec
| CREATE DATABASE pa_add_column_with_srid_zero_a5e965a1 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_column_with_srid_zero_a5e965a1
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE `pa_add_column_with_srid_zero_a5e965a1`.`t` ADD COLUMN `g` geometry NOT NULL /*!80003 SRID 0 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_column_with_srid_zero_a5e965a1.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_SRIDPresenceApplyCorrectness/8.0/add-srid-none-to-4326 default#1
> exec
| DROP DATABASE IF EXISTS pa_add_srid_none_to_4326_e34e7a6d_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_4326_e34e7a6d_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_4326_e34e7a6d_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_4326_e34e7a6d_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_srid_none_to_4326_e34e7a6d_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_4326_e34e7a6d_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_4326_e34e7a6d_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_4326_e34e7a6d_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_srid_none_to_4326_e34e7a6d
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_4326_e34e7a6d DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_4326_e34e7a6d
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> exec
| ALTER TABLE `pa_add_srid_none_to_4326_e34e7a6d`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 4326 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_4326_e34e7a6d.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_SRIDPresenceApplyCorrectness/8.0/add-srid-none-to-zero default#1
> exec
| DROP DATABASE IF EXISTS pa_add_srid_none_to_zero_cb2aa32e_to
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_zero_cb2aa32e_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_zero_cb2aa32e_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_zero_cb2aa32e_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_srid_none_to_zero_cb2aa32e_from
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_zero_cb2aa32e_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_zero_cb2aa32e_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_zero_cb2aa32e_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_add_srid_none_to_zero_cb2aa32e
< ok 0 0
> exec
| CREATE DATABASE pa_add_srid_none_to_zero_cb2aa32e DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_add_srid_none_to_zero_cb2aa32e
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> exec
| ALTER TABLE `pa_add_srid_none_to_zero_cb2aa32e`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 0 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_add_srid_none_to_zero_cb2aa32e.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_SRIDPresenceApplyCorrectness/8.0/change-srid-4326-to-zero default#1
> exec
| DROP DATABASE IF EXISTS pa_change_srid_4326_to_zero_878ba292_to
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_4326_to_zero_878ba292_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_4326_to_zero_878ba292_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_4326_to_zero_878ba292_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_change_srid_4326_to_zero_878ba292_from
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_4326_to_zero_878ba292_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_4326_to_zero_878ba292_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_4326_to_zero_878ba292_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_change_srid_4326_to_zero_878ba292
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_4326_to_zero_878ba292 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_4326_to_zero_878ba292
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> exec
| ALTER TABLE `pa_change_srid_4326_to_zero_878ba292`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 0 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_4326_to_zero_878ba292.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_SRIDPresenceApplyCorrectness/8.0/change-srid-zero-to-4326 default#1
> exec
| DROP DATABASE IF EXISTS pa_change_srid_zero_to_4326_a46d49d2_to
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_zero_to_4326_a46d49d2_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_zero_to_4326_a46d49d2_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_zero_to_4326_a46d49d2_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_change_srid_zero_to_4326_a46d49d2_from
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_zero_to_4326_a46d49d2_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_zero_to_4326_a46d49d2_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_zero_to_4326_a46d49d2_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_change_srid_zero_to_4326_a46d49d2
< ok 0 0
> exec
| CREATE DATABASE pa_change_srid_zero_to_4326_a46d49d2 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_change_srid_zero_to_4326_a46d49d2
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> exec
| ALTER TABLE `pa_change_srid_zero_to_4326_a46d49d2`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 4326 */
< ok 0 0
> query
| SHOW CREATE TABLE pa_change_srid_zero_to_4326_a46d49d2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_SRIDPresenceApplyCorrectness/8.0/create-srid-zero-spatial-index default#1
> exec
| DROP DATABASE IF EXISTS pa_create_srid_zero_spatial_index_62e8ed76_to
< ok 0 0
> exec
| CREATE DATABASE pa_create_srid_zero_spatial_index_62e8ed76_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_srid_zero_spatial_index_62e8ed76_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, SPATIAL INDEX sp_g (g))
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_srid_zero_spatial_index_62e8ed76_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_create_srid_zero_spatial_index_62e8ed76
< ok 0 0
> exec
| CREATE DATABASE pa_create_srid_zero_spatial_index_62e8ed76 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_create_srid_zero_spatial_index_62e8ed76
< ok 0 0
> exec
| CREATE TABLE `pa_create_srid_zero_spatial_index_62e8ed76`.`t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `pa_create_srid_zero_spatial_index_62e8ed76`.`t` ADD SPATIAL KEY `sp_g` (`g`)
< ok 0 0
> query
| SHOW CREATE TABLE pa_create_srid_zero_spatial_index_62e8ed76.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `sp_g` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_SRIDPresenceApplyCorrectness/8.0/remove-srid-4326-to-none default#1
> exec
| DROP DATABASE IF EXISTS pa_remove_srid_4326_to_none_442a3192_to
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_4326_to_none_442a3192_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_4326_to_none_442a3192_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_4326_to_none_442a3192_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_remove_srid_4326_to_none_442a3192_from
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_4326_to_none_442a3192_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_4326_to_none_442a3192_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_4326_to_none_442a3192_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 4326 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_remove_srid_4326_to_none_442a3192
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_4326_to_none_442a3192 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_4326_to_none_442a3192
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326)
< ok 0 0
> exec
| ALTER TABLE `pa_remove_srid_4326_to_none_442a3192`.`t` MODIFY COLUMN `g` geometry NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_4326_to_none_442a3192.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_SRIDPresenceApplyCorrectness/8.0/remove-srid-zero-to-none default#1
> exec
| DROP DATABASE IF EXISTS pa_remove_srid_zero_to_none_bfbe1bd5_to
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_zero_to_none_bfbe1bd5_to DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_zero_to_none_bfbe1bd5_to
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_zero_to_none_bfbe1bd5_to.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_remove_srid_zero_to_none_bfbe1bd5_from
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_zero_to_none_bfbe1bd5_from DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_zero_to_none_bfbe1bd5_from
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_zero_to_none_bfbe1bd5_from.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL /*!80003 SRID 0 */,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS pa_remove_srid_zero_to_none_bfbe1bd5
< ok 0 0
> exec
| CREATE DATABASE pa_remove_srid_zero_to_none_bfbe1bd5 DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE pa_remove_srid_zero_to_none_bfbe1bd5
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> exec
| ALTER TABLE `pa_remove_srid_zero_to_none_bfbe1bd5`.`t` MODIFY COLUMN `g` geometry NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE pa_remove_srid_zero_to_none_bfbe1bd5.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `g` geometry NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
