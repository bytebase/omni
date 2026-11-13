-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestContainer_Section_6_1_ShowCreateTableIntegration default#1
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
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'expr_json' can't have a default value
