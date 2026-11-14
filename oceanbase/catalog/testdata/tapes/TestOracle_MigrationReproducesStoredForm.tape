-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_MigrationReproducesStoredForm default#1
> exec
| DROP DATABASE IF EXISTS storedform_t_column_first
< ok 0 0
> exec
| CREATE DATABASE storedform_t_column_first DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_column_first
< ok 0 0
> exec
| CREATE TABLE category (is_featured tinyint(1) NOT NULL DEFAULT '0', category_id int, name varchar(25))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_column_first' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| category
> query
| SHOW CREATE TABLE storedform_t_column_first.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `is_featured` tinyint(1) NOT NULL DEFAULT '0',
|   `category_id` int(11) DEFAULT NULL,
|   `name` varchar(25) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS storedform
< ok 0 0
> exec
| CREATE DATABASE storedform DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform
< ok 0 0
> exec
| CREATE TABLE category (category_id int, name varchar(25))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| category
> query
| SHOW CREATE TABLE storedform.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int(11) DEFAULT NULL,
|   `name` varchar(25) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `storedform`.`category` ADD COLUMN `is_featured` tinyint(1) NOT NULL DEFAULT '0' FIRST
< ok 0 0
> query
| SHOW CREATE TABLE storedform.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `is_featured` tinyint(1) NOT NULL DEFAULT '0',
|   `category_id` int(11) DEFAULT NULL,
|   `name` varchar(25) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE storedform_t_column_first.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `is_featured` tinyint(1) NOT NULL DEFAULT '0',
|   `category_id` int(11) DEFAULT NULL,
|   `name` varchar(25) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS storedform_t_generated_columns_in_target_order
< ok 0 0
> exec
| CREATE DATABASE storedform_t_generated_columns_in_target_order DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_generated_columns_in_target_order
< ok 0 0
> exec
| CREATE TABLE p (id int, price decimal(10,2), name varchar(100), price_with_tax decimal(10,2) GENERATED ALWAYS AS (price*1.1) STORED, name_upper varchar(100) GENERATED ALWAYS AS (upper(name)) VIRTUAL)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_generated_columns_in_target_order' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| p
> query
| SHOW CREATE TABLE storedform_t_generated_columns_in_target_order.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * 1.1)) STORED,
|   `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS storedform
< ok 1 0
> exec
| CREATE DATABASE storedform DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform
< ok 0 0
> exec
| CREATE TABLE p (id int, price decimal(10,2), name varchar(100))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| p
> query
| SHOW CREATE TABLE storedform.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `name` varchar(100) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `storedform`.`p` ADD COLUMN `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * 1.1)) STORED
< ok 0 0
> exec
| ALTER TABLE `storedform`.`p` ADD COLUMN `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL
< ok 0 0
> query
| SHOW CREATE TABLE storedform.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * 1.1)) STORED,
|   `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE storedform_t_generated_columns_in_target_order.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * 1.1)) STORED,
|   `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS storedform_t_columns_mid_table
< ok 0 0
> exec
| CREATE DATABASE storedform_t_columns_mid_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_columns_mid_table
< ok 0 0
> exec
| CREATE TABLE t (a int, b int, g int GENERATED ALWAYS AS (b+1) VIRTUAL, c int, d int)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_columns_mid_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_t_columns_mid_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`b` + 1)) VIRTUAL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS storedform
< ok 1 0
> exec
| CREATE DATABASE storedform DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform
< ok 0 0
> exec
| CREATE TABLE t (a int, d int)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `b` int DEFAULT NULL AFTER `a`
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `c` int DEFAULT NULL AFTER `b`
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS ((`b` + 1)) VIRTUAL AFTER `b`
< ok 0 0
> query
| SHOW CREATE TABLE storedform.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`b` + 1)) VIRTUAL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE storedform_t_columns_mid_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`b` + 1)) VIRTUAL,
|   `c` int(11) DEFAULT NULL,
|   `d` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS storedform_t_column_readded_in_place
< ok 0 0
> exec
| CREATE DATABASE storedform_t_column_readded_in_place DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_column_readded_in_place
< ok 0 0
> exec
| CREATE TABLE t (a int, g int GENERATED ALWAYS AS (z+1) VIRTUAL, z int)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_column_readded_in_place' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_t_column_readded_in_place.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   `z` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS storedform
< ok 1 0
> exec
| CREATE DATABASE storedform DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform
< ok 0 0
> exec
| CREATE TABLE t (a int, g int GENERATED ALWAYS AS (a+1) VIRTUAL)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE `storedform`.`t` DROP COLUMN `g`
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `z` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL AFTER `a`
< ok 0 0
> query
| SHOW CREATE TABLE storedform.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   `z` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE storedform_t_column_readded_in_place.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   `g` int(11) GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   `z` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS storedform_t_fk_no_action
< ok 0 0
> exec
| CREATE DATABASE storedform_t_fk_no_action DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_fk_no_action
< ok 0 0
> exec
| CREATE TABLE city (city_id int NOT NULL, city varchar(50) NOT NULL, PRIMARY KEY (city_id))
< ok 0 0
> exec
| CREATE TABLE address (address_id int NOT NULL, city_id int NOT NULL, PRIMARY KEY (address_id), KEY idx_fk_city_id (city_id), CONSTRAINT fk_address_city FOREIGN KEY (city_id) REFERENCES city (city_id) ON DELETE NO ACTION ON UPDATE CASCADE, CONSTRAINT fk_address_city_2 FOREIGN KEY (city_id) REFERENCES city (city_id) ON UPDATE NO ACTION)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_fk_no_action' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| address
< row
: bytes
| city
> query
| SHOW CREATE TABLE storedform_t_fk_no_action.address
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| address
: bytes
| CREATE TABLE `address` (
|   `address_id` int(11) NOT NULL,
|   `city_id` int(11) NOT NULL,
|   PRIMARY KEY (`address_id`),
|   CONSTRAINT `fk_address_city_2` FOREIGN KEY (`city_id`) REFERENCES `storedform_t_fk_no_action`.`city`(`city_id`) ON UPDATE NO ACTION ON DELETE RESTRICT ,
|   CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `storedform_t_fk_no_action`.`city`(`city_id`) ON UPDATE CASCADE ON DELETE NO ACTION ,
|   KEY `idx_fk_city_id` (`city_id`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE storedform_t_fk_no_action.city
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| city
: bytes
| CREATE TABLE `city` (
|   `city_id` int(11) NOT NULL,
|   `city` varchar(50) NOT NULL,
|   PRIMARY KEY (`city_id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS storedform
< ok 1 0
> exec
| CREATE DATABASE storedform DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform
< ok 0 0
> exec
| CREATE TABLE city (city_id int NOT NULL, city varchar(50) NOT NULL, PRIMARY KEY (city_id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "table_name" "CHAR"
< row
: bytes
| city
> query
| SHOW CREATE TABLE storedform.city
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| city
: bytes
| CREATE TABLE `city` (
|   `city_id` int(11) NOT NULL,
|   `city` varchar(50) NOT NULL,
|   PRIMARY KEY (`city_id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE TABLE `storedform`.`address` (
|   `address_id` int NOT NULL,
|   `city_id` int NOT NULL,
|   PRIMARY KEY (`address_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `storedform`.`address` ADD KEY `idx_fk_city_id` (`city_id`) LOCAL
< ok 0 0
> exec
| ALTER TABLE `storedform`.`address` ADD CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `storedform`.`city` (`city_id`) ON DELETE NO ACTION ON UPDATE CASCADE
< ok 0 0
> exec
| ALTER TABLE `storedform`.`address` ADD CONSTRAINT `fk_address_city_2` FOREIGN KEY (`city_id`) REFERENCES `storedform`.`city` (`city_id`) ON UPDATE NO ACTION
< ok 0 0
> query
| SHOW CREATE TABLE storedform.address
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| address
: bytes
| CREATE TABLE `address` (
|   `address_id` int(11) NOT NULL,
|   `city_id` int(11) NOT NULL,
|   PRIMARY KEY (`address_id`),
|   CONSTRAINT `fk_address_city_2` FOREIGN KEY (`city_id`) REFERENCES `storedform`.`city`(`city_id`) ON UPDATE NO ACTION ON DELETE RESTRICT ,
|   CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `storedform`.`city`(`city_id`) ON UPDATE CASCADE ON DELETE NO ACTION ,
|   KEY `idx_fk_city_id` (`city_id`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE storedform_t_fk_no_action.address
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| address
: bytes
| CREATE TABLE `address` (
|   `address_id` int(11) NOT NULL,
|   `city_id` int(11) NOT NULL,
|   PRIMARY KEY (`address_id`),
|   CONSTRAINT `fk_address_city_2` FOREIGN KEY (`city_id`) REFERENCES `storedform_t_fk_no_action`.`city`(`city_id`) ON UPDATE NO ACTION ON DELETE RESTRICT ,
|   CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `storedform_t_fk_no_action`.`city`(`city_id`) ON UPDATE CASCADE ON DELETE NO ACTION ,
|   KEY `idx_fk_city_id` (`city_id`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
