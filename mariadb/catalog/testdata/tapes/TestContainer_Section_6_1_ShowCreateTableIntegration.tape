-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_6_1_ShowCreateTableIntegration default#1
> exec
| DROP TABLE IF EXISTS t_all_types
< ok 0 0
> exec
| CREATE TABLE t_all_types (
| 			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
| 			tiny_col TINYINT,
| 			small_col SMALLINT,
| 			med_col MEDIUMINT,
| 			int_col INT,
| 			big_col BIGINT,
| 			float_col FLOAT,
| 			float_prec FLOAT(7,3),
| 			double_col DOUBLE,
| 			double_prec DOUBLE(15,5),
| 			decimal_col DECIMAL(10,2),
| 			decimal_bare DECIMAL,
| 			bool_col BOOLEAN,
| 			bit_col BIT(8),
| 			char_col CHAR(10),
| 			varchar_col VARCHAR(255),
| 			tinytext_col TINYTEXT,
| 			text_col TEXT,
| 			mediumtext_col MEDIUMTEXT,
| 			longtext_col LONGTEXT,
| 			enum_col ENUM('a','b','c'),
| 			set_col SET('x','y','z'),
| 			binary_col BINARY(16),
| 			varbinary_col VARBINARY(255),
| 			tinyblob_col TINYBLOB,
| 			blob_col BLOB,
| 			mediumblob_col MEDIUMBLOB,
| 			longblob_col LONGBLOB,
| 			date_col DATE,
| 			time_col TIME,
| 			time_frac TIME(3),
| 			datetime_col DATETIME,
| 			datetime_frac DATETIME(6),
| 			timestamp_col TIMESTAMP NULL,
| 			timestamp_frac TIMESTAMP(3) NULL,
| 			year_col YEAR,
| 			json_col JSON,
| 			geo_col GEOMETRY,
| 			point_col POINT,
| 			linestring_col LINESTRING,
| 			polygon_col POLYGON,
| 			multipoint_col MULTIPOINT,
| 			multiline_col MULTILINESTRING,
| 			multipoly_col MULTIPOLYGON,
| 			geocoll_col GEOMETRYCOLLECTION,
| 			PRIMARY KEY (id)
| 		)
< ok 0 0
> query
| SHOW CREATE TABLE t_all_types
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_all_types
: bytes
| CREATE TABLE `t_all_types` (
|   `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
|   `tiny_col` tinyint(4) DEFAULT NULL,
|   `small_col` smallint(6) DEFAULT NULL,
|   `med_col` mediumint(9) DEFAULT NULL,
|   `int_col` int(11) DEFAULT NULL,
|   `big_col` bigint(20) DEFAULT NULL,
|   `float_col` float DEFAULT NULL,
|   `float_prec` float(7,3) DEFAULT NULL,
|   `double_col` double DEFAULT NULL,
|   `double_prec` double(15,5) DEFAULT NULL,
|   `decimal_col` decimal(10,2) DEFAULT NULL,
|   `decimal_bare` decimal(10,0) DEFAULT NULL,
|   `bool_col` tinyint(1) DEFAULT NULL,
|   `bit_col` bit(8) DEFAULT NULL,
|   `char_col` char(10) DEFAULT NULL,
|   `varchar_col` varchar(255) DEFAULT NULL,
|   `tinytext_col` tinytext DEFAULT NULL,
|   `text_col` text DEFAULT NULL,
|   `mediumtext_col` mediumtext DEFAULT NULL,
|   `longtext_col` longtext DEFAULT NULL,
|   `enum_col` enum('a','b','c') DEFAULT NULL,
|   `set_col` set('x','y','z') DEFAULT NULL,
|   `binary_col` binary(16) DEFAULT NULL,
|   `varbinary_col` varbinary(255) DEFAULT NULL,
|   `tinyblob_col` tinyblob DEFAULT NULL,
|   `blob_col` blob DEFAULT NULL,
|   `mediumblob_col` mediumblob DEFAULT NULL,
|   `longblob_col` longblob DEFAULT NULL,
|   `date_col` date DEFAULT NULL,
|   `time_col` time DEFAULT NULL,
|   `time_frac` time(3) DEFAULT NULL,
|   `datetime_col` datetime DEFAULT NULL,
|   `datetime_frac` datetime(6) DEFAULT NULL,
|   `timestamp_col` timestamp NULL DEFAULT NULL,
|   `timestamp_frac` timestamp(3) NULL DEFAULT NULL,
|   `year_col` year(4) DEFAULT NULL,
|   `json_col` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`json_col`)),
|   `geo_col` geometry DEFAULT NULL,
|   `point_col` point DEFAULT NULL,
|   `linestring_col` linestring DEFAULT NULL,
|   `polygon_col` polygon DEFAULT NULL,
|   `multipoint_col` multipoint DEFAULT NULL,
|   `multiline_col` multilinestring DEFAULT NULL,
|   `multipoly_col` multipolygon DEFAULT NULL,
|   `geocoll_col` geometrycollection DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_all_defaults
< ok 0 0
> exec
| CREATE TABLE t_all_defaults (
| 			id INT NOT NULL AUTO_INCREMENT,
| 			int_def INT DEFAULT 0,
| 			int_null INT DEFAULT NULL,
| 			int_notnull INT NOT NULL,
| 			varchar_def VARCHAR(100) DEFAULT 'hello',
| 			varchar_empty VARCHAR(100) DEFAULT '',
| 			float_def FLOAT DEFAULT 3.14,
| 			decimal_def DECIMAL(10,2) DEFAULT 0.00,
| 			bool_true BOOLEAN DEFAULT TRUE,
| 			bool_false BOOLEAN DEFAULT FALSE,
| 			enum_def ENUM('a','b','c') DEFAULT 'a',
| 			set_def SET('x','y','z') DEFAULT 'x,y',
| 			bit_def BIT(8) DEFAULT b'00001111',
| 			blob_col BLOB,
| 			text_col TEXT,
| 			json_col JSON,
| 			ts_def TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
| 			dt_def DATETIME DEFAULT CURRENT_TIMESTAMP,
| 			ts3_def TIMESTAMP(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
| 			expr_int INT DEFAULT (FLOOR(RAND()*100)),
| 			expr_json JSON DEFAULT (JSON_ARRAY()),
| 			expr_varchar VARCHAR(100) DEFAULT (UUID()),
| 			dt_literal DATETIME DEFAULT '2024-01-01 00:00:00',
| 			date_literal DATE DEFAULT '2024-01-01',
| 			PRIMARY KEY (id)
| 		)
< ok 0 0
> query
| SHOW CREATE TABLE t_all_defaults
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_all_defaults
: bytes
| CREATE TABLE `t_all_defaults` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `int_def` int(11) DEFAULT 0,
|   `int_null` int(11) DEFAULT NULL,
|   `int_notnull` int(11) NOT NULL,
|   `varchar_def` varchar(100) DEFAULT 'hello',
|   `varchar_empty` varchar(100) DEFAULT '',
|   `float_def` float DEFAULT 3.14,
|   `decimal_def` decimal(10,2) DEFAULT 0.00,
|   `bool_true` tinyint(1) DEFAULT 1,
|   `bool_false` tinyint(1) DEFAULT 0,
|   `enum_def` enum('a','b','c') DEFAULT 'a',
|   `set_def` set('x','y','z') DEFAULT 'x,y',
|   `bit_def` bit(8) DEFAULT b'1111',
|   `blob_col` blob DEFAULT NULL,
|   `text_col` text DEFAULT NULL,
|   `json_col` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`json_col`)),
|   `ts_def` timestamp NULL DEFAULT current_timestamp(),
|   `dt_def` datetime DEFAULT current_timestamp(),
|   `ts3_def` timestamp(3) NULL DEFAULT current_timestamp(3),
|   `expr_int` int(11) DEFAULT floor(rand() * 100),
|   `expr_json` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT json_array() CHECK (json_valid(`expr_json`)),
|   `expr_varchar` varchar(100) DEFAULT uuid(),
|   `dt_literal` datetime DEFAULT '2024-01-01 00:00:00',
|   `date_literal` date DEFAULT '2024-01-01',
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_all_indexes
< ok 0 0
> exec
| CREATE TABLE t_all_indexes (
| 			id INT NOT NULL AUTO_INCREMENT,
| 			name VARCHAR(255) NOT NULL,
| 			email VARCHAR(255),
| 			score INT,
| 			rank_val INT,
| 			bio TEXT,
| 			description TEXT,
| 			geo_col POINT NOT NULL,
| 			PRIMARY KEY (id),
| 			KEY idx_name (name),
| 			KEY idx_email_prefix (email(10)),
| 			KEY idx_score_desc (score DESC),
| 			KEY idx_multi_mixed (score ASC, rank_val DESC),
| 			UNIQUE KEY uk_email (email),
| 			FULLTEXT KEY ft_bio (bio, description),
| 			SPATIAL KEY sp_geo (geo_col),
| 			KEY idx_comment (name) COMMENT 'name lookup index'
| 		) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE t_all_indexes
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_all_indexes
: bytes
| CREATE TABLE `t_all_indexes` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `name` varchar(255) NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   `score` int(11) DEFAULT NULL,
|   `rank_val` int(11) DEFAULT NULL,
|   `bio` text DEFAULT NULL,
|   `description` text DEFAULT NULL,
|   `geo_col` point NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_email` (`email`),
|   KEY `idx_name` (`name`),
|   KEY `idx_email_prefix` (`email`(10)),
|   KEY `idx_score_desc` (`score` DESC),
|   KEY `idx_multi_mixed` (`score`,`rank_val` DESC),
|   SPATIAL KEY `sp_geo` (`geo_col`),
|   KEY `idx_comment` (`name`) COMMENT 'name lookup index',
|   FULLTEXT KEY `ft_bio` (`bio`,`description`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_all_constraints
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent
< ok 0 0
> exec
| CREATE TABLE t_parent (
| 			id INT NOT NULL AUTO_INCREMENT,
| 			code VARCHAR(10) NOT NULL,
| 			PRIMARY KEY (id),
| 			UNIQUE KEY uk_code (code)
| 		)
< ok 0 0
> query
| SHOW CREATE TABLE t_parent
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_parent
: bytes
| CREATE TABLE `t_parent` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `code` varchar(10) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk_code` (`code`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_all_constraints
< ok 0 0
> exec
| CREATE TABLE t_all_constraints (
| 			id INT NOT NULL AUTO_INCREMENT,
| 			parent_id INT,
| 			parent_code VARCHAR(10),
| 			self_ref INT,
| 			val INT,
| 			score INT,
| 			PRIMARY KEY (id),
| 			CONSTRAINT fk_parent FOREIGN KEY (parent_id) REFERENCES t_parent(id) ON DELETE CASCADE,
| 			CONSTRAINT fk_code FOREIGN KEY (parent_code) REFERENCES t_parent(code) ON UPDATE SET NULL,
| 			CONSTRAINT fk_self FOREIGN KEY (self_ref) REFERENCES t_all_constraints(id) ON DELETE SET NULL ON UPDATE CASCADE,
| 			CONSTRAINT chk_val CHECK (val > 0),
| 			CONSTRAINT chk_score CHECK (score >= 0 AND score <= 100)
| 		)
< ok 0 0
> query
| SHOW CREATE TABLE t_all_constraints
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_all_constraints
: bytes
| CREATE TABLE `t_all_constraints` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `parent_id` int(11) DEFAULT NULL,
|   `parent_code` varchar(10) DEFAULT NULL,
|   `self_ref` int(11) DEFAULT NULL,
|   `val` int(11) DEFAULT NULL,
|   `score` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fk_parent` (`parent_id`),
|   KEY `fk_code` (`parent_code`),
|   KEY `fk_self` (`self_ref`),
|   CONSTRAINT `fk_code` FOREIGN KEY (`parent_code`) REFERENCES `t_parent` (`code`) ON UPDATE SET NULL,
|   CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `t_parent` (`id`) ON DELETE CASCADE,
|   CONSTRAINT `fk_self` FOREIGN KEY (`self_ref`) REFERENCES `t_all_constraints` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
|   CONSTRAINT `chk_val` CHECK (`val` > 0),
|   CONSTRAINT `chk_score` CHECK (`score` >= 0 and `score` <= 100)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_part_range
< ok 0 0
> exec
| CREATE TABLE t_part_range (
| 			id INT NOT NULL,
| 			created_date DATE NOT NULL,
| 			name VARCHAR(100),
| 			PRIMARY KEY (id, created_date)
| 		) ENGINE=InnoDB
| 		PARTITION BY RANGE (YEAR(created_date)) (
| 			PARTITION p2020 VALUES LESS THAN (2021),
| 			PARTITION p2021 VALUES LESS THAN (2022),
| 			PARTITION p2022 VALUES LESS THAN (2023),
| 			PARTITION pmax VALUES LESS THAN MAXVALUE
| 		)
< ok 0 0
> query
| SHOW CREATE TABLE t_part_range
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_range
: bytes
| CREATE TABLE `t_part_range` (
|   `id` int(11) NOT NULL,
|   `created_date` date NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`,`created_date`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY RANGE (year(`created_date`))
| (PARTITION `p2020` VALUES LESS THAN (2021) ENGINE = InnoDB,
|  PARTITION `p2021` VALUES LESS THAN (2022) ENGINE = InnoDB,
|  PARTITION `p2022` VALUES LESS THAN (2023) ENGINE = InnoDB,
|  PARTITION `pmax` VALUES LESS THAN MAXVALUE ENGINE = InnoDB)
> exec
| DROP TABLE IF EXISTS t_part_list
< ok 0 0
> exec
| CREATE TABLE t_part_list (
| 			id INT NOT NULL AUTO_INCREMENT,
| 			region VARCHAR(20) NOT NULL,
| 			data VARCHAR(100),
| 			PRIMARY KEY (id, region)
| 		) ENGINE=InnoDB
| 		PARTITION BY LIST COLUMNS(region) (
| 			PARTITION p_east VALUES IN ('east','northeast'),
| 			PARTITION p_west VALUES IN ('west','northwest'),
| 			PARTITION p_other VALUES IN ('central','south')
| 		)
< ok 0 0
> query
| SHOW CREATE TABLE t_part_list
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_list
: bytes
| CREATE TABLE `t_part_list` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `region` varchar(20) NOT NULL,
|   `data` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`,`region`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY LIST  COLUMNS(`region`)
| (PARTITION `p_east` VALUES IN ('east','northeast') ENGINE = InnoDB,
|  PARTITION `p_west` VALUES IN ('west','northwest') ENGINE = InnoDB,
|  PARTITION `p_other` VALUES IN ('central','south') ENGINE = InnoDB)
> exec
| DROP TABLE IF EXISTS t_part_hash
< ok 0 0
> exec
| CREATE TABLE t_part_hash (
| 			id INT NOT NULL AUTO_INCREMENT,
| 			val INT,
| 			PRIMARY KEY (id)
| 		) ENGINE=InnoDB
| 		PARTITION BY HASH(id)
| 		PARTITIONS 4
< ok 0 0
> query
| SHOW CREATE TABLE t_part_hash
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_part_hash
: bytes
| CREATE TABLE `t_part_hash` (
|   `id` int(11) NOT NULL AUTO_INCREMENT,
|   `val` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
|  PARTITION BY HASH (`id`)
| PARTITIONS 4
