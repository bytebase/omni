-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSDLRealWorldOriginalVsSynced default#1
> begin
< ok
> exec
| --  Sample employee database 
| --  See changelog table for details
| --  Copyright (C) 2007,2008, MySQL AB
| --  
| --  Original data created by Fusheng Wang and Carlo Zaniolo
| --  http://www.cs.aau.dk/TimeCenter/software.htm
| --  http://www.cs.aau.dk/TimeCenter/Data/employeeTemporalDataSet.zip
| -- 
| --  Current schema by Giuseppe Maxia 
| --  Data conversion from XML to relational by Patrick Crews
| -- 
| -- This work is licensed under the 
| -- Creative Commons Attribution-Share Alike 3.0 Unported License. 
| -- To view a copy of this license, visit 
| -- http://creativecommons.org/licenses/by-sa/3.0/ or send a letter to 
| -- Creative Commons, 171 Second Street, Suite 300, San Francisco, 
| -- California, 94105, USA.
| -- 
| --  DISCLAIMER
| --  To the best of our knowledge, this data is fabricated, and
| --  it does not correspond to real people. 
| --  Any similarity to existing people is purely coincidental.
| -- 
|
|
|
|
|
| /*!50503 set default_storage_engine = InnoDB */;
< ok 0 0
> exec
|
| /*!50503 select CONCAT('storage engine: ', @@default_storage_engine) as INFO */;
< ok 0 0
> exec
|
|
| CREATE TABLE employees (
|     emp_no      INT             NOT NULL,
|     birth_date  DATE            NOT NULL,
|     first_name  VARCHAR(14)     NOT NULL,
|     last_name   VARCHAR(16)     NOT NULL,
|     gender      ENUM ('M','F')  NOT NULL,    
|     hire_date   DATE            NOT NULL,
|     PRIMARY KEY (emp_no)
| );
< ok 0 0
> exec
|
|
| CREATE TABLE departments (
|     dept_no     CHAR(4)         NOT NULL,
|     dept_name   VARCHAR(40)     NOT NULL,
|     PRIMARY KEY (dept_no),
|     UNIQUE  KEY (dept_name)
| );
< ok 0 0
> exec
|
|
| CREATE TABLE dept_manager (
|    emp_no       INT             NOT NULL,
|    dept_no      CHAR(4)         NOT NULL,
|    from_date    DATE            NOT NULL,
|    to_date      DATE            NOT NULL,
|    FOREIGN KEY (emp_no)  REFERENCES employees (emp_no)    ON DELETE CASCADE,
|    FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE,
|    PRIMARY KEY (emp_no,dept_no)
| );
< ok 0 0
> exec
|  
|
| CREATE TABLE dept_emp (
|     emp_no      INT             NOT NULL,
|     dept_no     CHAR(4)         NOT NULL,
|     from_date   DATE            NOT NULL,
|     to_date     DATE            NOT NULL,
|     FOREIGN KEY (emp_no)  REFERENCES employees   (emp_no)  ON DELETE CASCADE,
|     FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE,
|     PRIMARY KEY (emp_no,dept_no)
| );
< ok 0 0
> exec
|
|
| CREATE TABLE titles (
|     emp_no      INT             NOT NULL,
|     title       VARCHAR(50)     NOT NULL,
|     from_date   DATE            NOT NULL,
|     to_date     DATE,
|     FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE,
|     PRIMARY KEY (emp_no,title, from_date)
| ) 
| ;
< ok 0 0
> exec
|  
|
| CREATE TABLE salaries (
|     emp_no      INT             NOT NULL,
|     salary      INT             NOT NULL,
|     from_date   DATE            NOT NULL,
|     to_date     DATE            NOT NULL,
|     FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE,
|     PRIMARY KEY (emp_no, from_date)
| ) 
| ;
< ok 0 0
> exec
|  
|
| CREATE OR REPLACE VIEW dept_emp_latest_date AS
|     SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date
|     FROM dept_emp
|     GROUP BY emp_no;
< ok 0 0
> exec
|
|
| # shows only the current department for each employee
| CREATE OR REPLACE VIEW current_dept_emp AS
|     SELECT l.emp_no, dept_no, l.from_date, l.to_date
|     FROM dept_emp d
|         INNER JOIN dept_emp_latest_date l
|         ON d.emp_no=l.emp_no AND d.from_date=l.from_date AND l.to_date = d.to_date;
< ok 0 0
> commit
< ok
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 8.0.46
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| departments
: bytes
| dept_name
: bytes
| dept_name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| departments
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_emp
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| dept_no
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dept_manager
: bytes
| PRIMARY
: bytes
| dept_no
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| employees
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| salaries
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| emp_no
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| title
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| titles
: bytes
| PRIMARY
: bytes
| from_date
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			IFNULL(COLUMN_NAME, ''),
| 			ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE,
| 			COLUMN_TYPE,
| 			IFNULL(CHARACTER_SET_NAME, ''),
| 			IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT),
| 			convert(GENERATION_EXPRESSION using BINARY),
| 			EXTRA,
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
< row
: bytes
| current_dept_emp
: bytes
| emp_no
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| current_dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_no
: int64 1
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| departments
: bytes
| dept_name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| emp_no
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| emp_no
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| from_date
: int64 2
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_emp_latest_date
: bytes
| to_date
: int64 3
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| emp_no
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| dept_no
: int64 2
: null
: bytes
| NO
: bytes
| char(4)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dept_manager
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| emp_no
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| birth_date
: int64 2
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| first_name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(14)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| last_name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| gender
: int64 5
: null
: bytes
| NO
: bytes
| enum('M','F')
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| employees
: bytes
| hire_date
: int64 6
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| emp_no
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| salary
: int64 2
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| salaries
: bytes
| to_date
: int64 4
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| emp_no
: int64 1
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(50)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| from_date
: int64 3
: null
: bytes
| NO
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| titles
: bytes
| to_date
: int64 4
: null
: bytes
| YES
: bytes
| date
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE,
| 				tc.ENFORCED
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| current_dept_emp
: bytes
| select `sdl_rw_employees_1`.`l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`sdl_rw_employees_1`.`l`.`from_date` AS `from_date`,`sdl_rw_employees_1`.`l`.`to_date` AS `to_date` from (`sdl_rw_employees_1`.`dept_emp` `d` join `sdl_rw_employees_1`.`dept_emp_latest_date` `l` on(((`d`.`emp_no` = `sdl_rw_employees_1`.`l`.`emp_no`) and (`d`.`from_date` = `sdl_rw_employees_1`.`l`.`from_date`) and (`sdl_rw_employees_1`.`l`.`to_date` = `d`.`to_date`))))
< row
: bytes
| dept_emp_latest_date
: bytes
| select `sdl_rw_employees_1`.`dept_emp`.`emp_no` AS `emp_no`,max(`sdl_rw_employees_1`.`dept_emp`.`from_date`) AS `from_date`,max(`sdl_rw_employees_1`.`dept_emp`.`to_date`) AS `to_date` from `sdl_rw_employees_1`.`dept_emp` group by `sdl_rw_employees_1`.`dept_emp`.`emp_no`
> query
| SHOW CREATE VIEW `sdl_rw_employees_1`.`current_dept_emp`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| current_dept_emp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from (`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`) and (`l`.`to_date` = `d`.`to_date`))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sdl_rw_employees_1`.`dept_emp_latest_date`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `dept_emp_latest_date` AS select `dept_emp`.`emp_no` AS `emp_no`,max(`dept_emp`.`from_date`) AS `from_date`,max(`dept_emp`.`to_date`) AS `to_date` from `dept_emp` group by `dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
|
| 	SELECT
| 		TRIGGER_NAME,
| 		EVENT_OBJECT_TABLE,
| 		EVENT_MANIPULATION,
| 		ACTION_TIMING,
| 		ACTION_STATEMENT,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.TRIGGERS
| 	WHERE EVENT_OBJECT_SCHEMA = ?
| 	ORDER BY EVENT_OBJECT_TABLE ASC, EVENT_MANIPULATION ASC, ACTION_TIMING ASC, ACTION_ORDER ASC;
| 	
: string
| sdl_rw_employees_1
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
> query
|
| 	SELECT
| 		EVENT_NAME,
| 		TIME_ZONE,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.EVENTS
| 	WHERE EVENT_SCHEMA = ?
| 	ORDER BY EVENT_NAME ASC;
| 	
: string
| sdl_rw_employees_1
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| departments
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| departments
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| titles
: bytes
| titles_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| salaries
: bytes
| salaries_ibfk_1
: bytes
| employees
: bytes
| CASCADE
: bytes
| NO ACTION
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
| sdl_rw_employees_1
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_emp
: bytes
| dept_emp_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_1
: bytes
| dept_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| dept_manager
: bytes
| dept_manager_ibfk_2
: bytes
| dept_no
: bytes
| sdl_rw_employees_1
: bytes
| dept_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| salaries
: bytes
| salaries_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
< row
: bytes
| sdl_rw_employees_1
: bytes
| titles
: bytes
| titles_ibfk_1
: bytes
| emp_no
: bytes
| sdl_rw_employees_1
: bytes
| emp_no
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			PARTITION_NAME,
| 			SUBPARTITION_NAME,
| 			PARTITION_METHOD,
| 			SUBPARTITION_METHOD,
| 			PARTITION_EXPRESSION,
| 			SUBPARTITION_EXPRESSION,
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
> query
|
| 		SELECT
| 			ROUTINE_NAME,
| 			ROUTINE_TYPE,
| 			SQL_MODE,
| 			CHARACTER_SET_CLIENT,
| 			COLLATION_CONNECTION,
| 			DATABASE_COLLATION,
| 			IFNULL(ROUTINE_COMMENT, '')
| 		FROM
| 			INFORMATION_SCHEMA.ROUTINES
| 		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE IN ('FUNCTION', 'PROCEDURE')
| 		ORDER BY ROUTINE_TYPE, ROUTINE_NAME;
| 	
: string
| sdl_rw_employees_1
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
> query
|
| 		SELECT
| 			TABLES.TABLE_NAME,
| 			TABLES.TABLE_TYPE,
| 			IFNULL(TABLES.ENGINE, ''),
| 			IFNULL(TABLES.TABLE_COLLATION, ''),
| 			IFNULL(TABLES.TABLE_ROWS, 0),
| 			IFNULL(TABLES.DATA_LENGTH, 0),
| 			IFNULL(TABLES.INDEX_LENGTH, 0),
| 			IFNULL(TABLES.DATA_FREE, 0),
| 			IFNULL(TABLES.CREATE_OPTIONS, ''),
| 			QUOTE(IFNULL(TABLES.TABLE_COMMENT, '')),
| 			IFNULL(CCSA.CHARACTER_SET_NAME, '')
| 		FROM information_schema.TABLES TABLES
| 		LEFT JOIN information_schema.COLLATION_CHARACTER_SET_APPLICABILITY CCSA
| 		ON TABLES.TABLE_COLLATION = CCSA.COLLATION_NAME
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_rw_employees_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< row
: bytes
| current_dept_emp
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| departments
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| dept_emp
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| dept_emp_latest_date
: bytes
| VIEW
: bytes
|
: bytes
|
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| 'VIEW'
: bytes
|
< row
: bytes
| dept_manager
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| employees
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| salaries
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| titles
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
|
: bytes
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sdl_rw_employees_1
< result
< column "NAME" "VARCHAR" "string"
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| sdl_rw_employees_1
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLRealWorldOriginalVsSynced default#2
> begin
< ok
> exec
| -- Roundcube Webmail initial database structure
|
|
| SET FOREIGN_KEY_CHECKS=0;
< ok 0 0
> exec
|
|
| -- Table structure for table `session`
|
| CREATE TABLE `session` (
|  `sess_id` varchar(128) NOT NULL,
|  `expires_at` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `ip` varchar(40) NOT NULL,
|  `vars` mediumtext NOT NULL,
|  PRIMARY KEY(`sess_id`),
|  INDEX `expires_at_index` (`expires_at`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `users`
|
| CREATE TABLE `users` (
|  `user_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `username` varchar(128) BINARY NOT NULL,
|  `mail_host` varchar(128) NOT NULL,
|  `created` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `last_login` datetime,
|  `failed_login` datetime,
|  `failed_login_counter` int(10) UNSIGNED,
|  `language` varchar(16),
|  `preferences` longtext,
|  PRIMARY KEY(`user_id`),
|  UNIQUE `username` (`username`, `mail_host`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache`
|
| CREATE TABLE `cache` (
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `cache_key` varchar(128) BINARY NOT NULL,
|  `expires` datetime,
|  `data` longtext NOT NULL,
|  PRIMARY KEY (`user_id`, `cache_key`),
|  CONSTRAINT `user_id_fk_cache` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `expires_index` (`expires`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache_shared`
|
| CREATE TABLE `cache_shared` (
|  `cache_key` varchar(255) BINARY NOT NULL,
|  `expires` datetime,
|  `data` longtext NOT NULL,
|  PRIMARY KEY (`cache_key`),
|  INDEX `expires_index` (`expires`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache_index`
|
| CREATE TABLE `cache_index` (
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `mailbox` varchar(255) BINARY NOT NULL,
|  `expires` datetime,
|  `valid` tinyint(1) NOT NULL DEFAULT '0',
|  `data` longtext NOT NULL,
|  CONSTRAINT `user_id_fk_cache_index` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `expires_index` (`expires`),
|  PRIMARY KEY (`user_id`, `mailbox`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache_thread`
|
| CREATE TABLE `cache_thread` (
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `mailbox` varchar(255) BINARY NOT NULL,
|  `expires` datetime,
|  `data` longtext NOT NULL,
|  CONSTRAINT `user_id_fk_cache_thread` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `expires_index` (`expires`),
|  PRIMARY KEY (`user_id`, `mailbox`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `cache_messages`
|
| CREATE TABLE `cache_messages` (
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `mailbox` varchar(255) BINARY NOT NULL,
|  `uid` int(11) UNSIGNED NOT NULL DEFAULT '0',
|  `expires` datetime,
|  `data` longtext NOT NULL,
|  `flags` int(11) NOT NULL DEFAULT '0',
|  CONSTRAINT `user_id_fk_cache_messages` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `expires_index` (`expires`),
|  PRIMARY KEY (`user_id`, `mailbox`, `uid`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `collected_addresses`
|
| CREATE TABLE `collected_addresses` (
|  `address_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `name` varchar(255) NOT NULL DEFAULT '',
|  `email` varchar(255) NOT NULL,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `type` int(10) UNSIGNED NOT NULL,
|  PRIMARY KEY(`address_id`),
|  CONSTRAINT `user_id_fk_collected_addresses` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  UNIQUE INDEX `user_email_collected_addresses_index` (`user_id`, `type`, `email`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `contacts`
|
| CREATE TABLE `contacts` (
|  `contact_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `del` tinyint(1) NOT NULL DEFAULT '0',
|  `name` varchar(128) NOT NULL DEFAULT '',
|  `email` text NOT NULL,
|  `firstname` varchar(128) NOT NULL DEFAULT '',
|  `surname` varchar(128) NOT NULL DEFAULT '',
|  `vcard` longtext,
|  `words` text,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  PRIMARY KEY(`contact_id`),
|  CONSTRAINT `user_id_fk_contacts` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `user_contacts_index` (`user_id`,`del`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `contactgroups`
|
| CREATE TABLE `contactgroups` (
|   `contactgroup_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|   `user_id` int(10) UNSIGNED NOT NULL,
|   `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|   `del` tinyint(1) NOT NULL DEFAULT '0',
|   `name` varchar(128) NOT NULL DEFAULT '',
|   PRIMARY KEY(`contactgroup_id`),
|   CONSTRAINT `user_id_fk_contactgroups` FOREIGN KEY (`user_id`)
|     REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|   INDEX `contactgroups_user_index` (`user_id`,`del`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `contactgroupmembers`
|
| CREATE TABLE `contactgroupmembers` (
|   `contactgroup_id` int(10) UNSIGNED NOT NULL,
|   `contact_id` int(10) UNSIGNED NOT NULL,
|   `created` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|   PRIMARY KEY (`contactgroup_id`, `contact_id`),
|   CONSTRAINT `contactgroup_id_fk_contactgroups` FOREIGN KEY (`contactgroup_id`)
|     REFERENCES `contactgroups`(`contactgroup_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|   CONSTRAINT `contact_id_fk_contacts` FOREIGN KEY (`contact_id`)
|     REFERENCES `contacts`(`contact_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|   INDEX `contactgroupmembers_contact_index` (`contact_id`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `identities`
|
| CREATE TABLE `identities` (
|  `identity_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `del` tinyint(1) NOT NULL DEFAULT '0',
|  `standard` tinyint(1) NOT NULL DEFAULT '0',
|  `name` varchar(128) NOT NULL,
|  `organization` varchar(128) NOT NULL DEFAULT '',
|  `email` varchar(128) NOT NULL,
|  `reply-to` varchar(128) NOT NULL DEFAULT '',
|  `bcc` varchar(128) NOT NULL DEFAULT '',
|  `signature` longtext,
|  `html_signature` tinyint(1) NOT NULL DEFAULT '0',
|  PRIMARY KEY(`identity_id`),
|  CONSTRAINT `user_id_fk_identities` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `user_identities_index` (`user_id`, `del`),
|  INDEX `email_identities_index` (`email`, `del`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `responses`
|
| CREATE TABLE `responses` (
|  `response_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `name` varchar(255) NOT NULL,
|  `data` longtext NOT NULL,
|  `is_html` tinyint(1) NOT NULL DEFAULT '0',
|  `changed` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  `del` tinyint(1) NOT NULL DEFAULT '0',
|  PRIMARY KEY (`response_id`),
|  CONSTRAINT `user_id_fk_responses` FOREIGN KEY (`user_id`)
|    REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  INDEX `user_responses_index` (`user_id`, `del`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `dictionary`
|
| CREATE TABLE `dictionary` (
|   `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY, -- redundant, for compat. with Galera Cluster
|   `user_id` int(10) UNSIGNED, -- NULL here is for "shared dictionaries"
|   `language` varchar(16) NOT NULL,
|   `data` longtext NOT NULL,
|   CONSTRAINT `user_id_fk_dictionary` FOREIGN KEY (`user_id`)
|     REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|   UNIQUE `uniqueness` (`user_id`, `language`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
|
| -- Table structure for table `searches`
|
| CREATE TABLE `searches` (
|  `search_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `type` int(3) NOT NULL DEFAULT '0',
|  `name` varchar(128) NOT NULL,
|  `data` text,
|  PRIMARY KEY(`search_id`),
|  CONSTRAINT `user_id_fk_searches` FOREIGN KEY (`user_id`)
|    REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  UNIQUE `uniqueness` (`user_id`, `type`, `name`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
| -- Table structure for table `filestore`
|
| CREATE TABLE `filestore` (
|  `file_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
|  `user_id` int(10) UNSIGNED NOT NULL,
|  `context` varchar(32) NOT NULL,
|  `filename` varchar(128) NOT NULL,
|  `mtime` int(10) NOT NULL,
|  `data` longtext NOT NULL,
|  PRIMARY KEY (`file_id`),
|  CONSTRAINT `user_id_fk_filestore` FOREIGN KEY (`user_id`)
|    REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
|  UNIQUE `uniqueness` (`user_id`, `context`, `filename`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
| -- Table structure for table `uploads`
|
| CREATE TABLE `uploads` (
|  `upload_id` varchar(64) NOT NULL,
|  `session_id` varchar(128) NOT NULL,
|  `group` varchar(128) NOT NULL,
|  `metadata` mediumtext NOT NULL,
|  `created` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
|  PRIMARY KEY (`upload_id`),
|  INDEX `uploads_session_group_index` (`session_id`, `group`, `created`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
| -- Table structure for table `system`
|
| CREATE TABLE `system` (
|  `name` varchar(64) NOT NULL,
|  `value` mediumtext,
|  PRIMARY KEY(`name`)
| ) ROW_FORMAT=DYNAMIC ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
< ok 0 0
> exec
|
|
| SET FOREIGN_KEY_CHECKS=1;
< ok 0 0
> commit
< ok
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR" "string"
< row
: bytes
| 8.0.46
> query
|
| 			SELECT
| 				TABLE_NAME,
| 				INDEX_NAME,
| 				COLUMN_NAME,
| 				COLLATION,
| 				IFNULL(SUB_PART, -1),
| 				EXPRESSION,
| 				SEQ_IN_INDEX,
| 				INDEX_TYPE,
| 				NON_UNIQUE,
| 				CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END,
| 				INDEX_COMMENT
| 			FROM information_schema.STATISTICS
| 			WHERE TABLE_SCHEMA = ?
| 			ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| sdl_rw_roundcube_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION" "VARCHAR"
< column "IFNULL(SUB_PART, -1)" "BIGINT" "int64"
< column "EXPRESSION" "TEXT"
< column "SEQ_IN_INDEX" "UNSIGNED INT" "uint32"
< column "INDEX_TYPE" "VARCHAR" "string"
< column "NON_UNIQUE" "INT" "int32"
< column "CASE IS_VISIBLE WHEN 'YES' THEN 1 ELSE 0 END" "INT" "int32"
< column "INDEX_COMMENT" "VARCHAR" "string"
< row
: bytes
| cache
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache
: bytes
| PRIMARY
: bytes
| cache_key
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_index
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_index
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_index
: bytes
| PRIMARY
: bytes
| mailbox
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| PRIMARY
: bytes
| mailbox
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_messages
: bytes
| PRIMARY
: bytes
| uid
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_shared
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_shared
: bytes
| PRIMARY
: bytes
| cache_key
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_thread
: bytes
| expires_index
: bytes
| expires
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| cache_thread
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| cache_thread
: bytes
| PRIMARY
: bytes
| mailbox
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| PRIMARY
: bytes
| address_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| user_email_collected_addresses_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| user_email_collected_addresses_index
: bytes
| type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| collected_addresses
: bytes
| user_email_collected_addresses_index
: bytes
| email
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contactgroupmembers
: bytes
| contactgroupmembers_contact_index
: bytes
| contact_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contactgroupmembers
: bytes
| PRIMARY
: bytes
| contactgroup_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contactgroupmembers
: bytes
| PRIMARY
: bytes
| contact_id
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contactgroups
: bytes
| contactgroups_user_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contactgroups
: bytes
| contactgroups_user_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contactgroups
: bytes
| PRIMARY
: bytes
| contactgroup_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contacts
: bytes
| PRIMARY
: bytes
| contact_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| contacts
: bytes
| user_contacts_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| contacts
: bytes
| user_contacts_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| dictionary
: bytes
| PRIMARY
: bytes
| id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dictionary
: bytes
| uniqueness
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| dictionary
: bytes
| uniqueness
: bytes
| language
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| PRIMARY
: bytes
| file_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| uniqueness
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| uniqueness
: bytes
| context
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| filestore
: bytes
| uniqueness
: bytes
| filename
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| email_identities_index
: bytes
| email
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| email_identities_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| PRIMARY
: bytes
| identity_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| user_identities_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| identities
: bytes
| user_identities_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| responses
: bytes
| PRIMARY
: bytes
| response_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| responses
: bytes
| user_responses_index
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| responses
: bytes
| user_responses_index
: bytes
| del
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| PRIMARY
: bytes
| search_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| uniqueness
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| uniqueness
: bytes
| type
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| searches
: bytes
| uniqueness
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| session
: bytes
| expires_at_index
: bytes
| expires_at
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| session
: bytes
| PRIMARY
: bytes
| sess_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| system
: bytes
| PRIMARY
: bytes
| name
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| PRIMARY
: bytes
| upload_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| uploads_session_group_index
: bytes
| session_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| uploads_session_group_index
: bytes
| group
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| uploads
: bytes
| uploads_session_group_index
: bytes
| created
: bytes
| A
: int64 -1
: null
: int64 3
: bytes
| BTREE
: int64 1
: int64 1
: bytes
|
< row
: bytes
| users
: bytes
| PRIMARY
: bytes
| user_id
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| users
: bytes
| username
: bytes
| username
: bytes
| A
: int64 -1
: null
: int64 1
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
< row
: bytes
| users
: bytes
| username
: bytes
| mail_host
: bytes
| A
: int64 -1
: null
: int64 2
: bytes
| BTREE
: int64 0
: int64 1
: bytes
|
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			IFNULL(COLUMN_NAME, ''),
| 			ORDINAL_POSITION,
| 			CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END,
| 			IS_NULLABLE,
| 			COLUMN_TYPE,
| 			IFNULL(CHARACTER_SET_NAME, ''),
| 			IFNULL(COLLATION_NAME, ''),
| 			QUOTE(COLUMN_COMMENT),
| 			convert(GENERATION_EXPRESSION using BINARY),
| 			EXTRA,
| 			SRS_ID
| 		FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| sdl_rw_roundcube_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "IFNULL(COLUMN_NAME, '')" "VARCHAR" "string"
< column "ORDINAL_POSITION" "UNSIGNED INT" "uint32"
< column "CASE WHEN COLUMN_DEFAULT is NULL THEN NULL ELSE QUOTE(COLUMN_DEFAULT) END" "TEXT"
< column "IS_NULLABLE" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IFNULL(CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< column "IFNULL(COLLATION_NAME, '')" "VARCHAR" "string"
< column "QUOTE(COLUMN_COMMENT)" "TEXT"
< column "convert(GENERATION_EXPRESSION using BINARY)" "BLOB" "[]uint8"
< column "EXTRA" "VARCHAR"
< column "SRS_ID" "UNSIGNED INT"
< row
: bytes
| cache
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache
: bytes
| cache_key
: int64 2
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache
: bytes
| expires
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| mailbox
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| expires
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| valid
: int64 4
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_index
: bytes
| data
: int64 5
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| mailbox
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| uid
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| expires
: int64 4
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| data
: int64 5
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_messages
: bytes
| flags
: int64 6
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_shared
: bytes
| cache_key
: int64 1
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_shared
: bytes
| expires
: int64 2
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_shared
: bytes
| data
: int64 3
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| mailbox
: int64 2
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| expires
: int64 3
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| cache_thread
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| address_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| collected_addresses
: bytes
| changed
: int64 2
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| name
: int64 3
: bytes
| ''
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| email
: int64 4
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| user_id
: int64 5
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| collected_addresses
: bytes
| type
: int64 6
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroupmembers
: bytes
| contactgroup_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroupmembers
: bytes
| contact_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroupmembers
: bytes
| created
: int64 3
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| contactgroup_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| contactgroups
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| changed
: int64 3
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| del
: int64 4
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contactgroups
: bytes
| name
: int64 5
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| contact_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| contacts
: bytes
| changed
: int64 2
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| del
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| name
: int64 4
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| email
: int64 5
: null
: bytes
| NO
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| firstname
: int64 6
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| surname
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| vcard
: int64 8
: null
: bytes
| YES
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| words
: int64 9
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| contacts
: bytes
| user_id
: int64 10
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dictionary
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| dictionary
: bytes
| user_id
: int64 2
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dictionary
: bytes
| language
: int64 3
: null
: bytes
| NO
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| dictionary
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| file_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| filestore
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| context
: int64 3
: null
: bytes
| NO
: bytes
| varchar(32)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| filename
: int64 4
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| mtime
: int64 5
: null
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| filestore
: bytes
| data
: int64 6
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| identity_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| identities
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| changed
: int64 3
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| del
: int64 4
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| standard
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| name
: int64 6
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| organization
: int64 7
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| email
: int64 8
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| reply-to
: int64 9
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| bcc
: int64 10
: bytes
| ''
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| signature
: int64 11
: null
: bytes
| YES
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| identities
: bytes
| html_signature
: int64 12
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| response_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| responses
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| data
: int64 4
: null
: bytes
| NO
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| is_html
: int64 5
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| changed
: int64 6
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| responses
: bytes
| del
: int64 7
: bytes
| '0'
: bytes
| NO
: bytes
| tinyint(1)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searches
: bytes
| search_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| searches
: bytes
| user_id
: int64 2
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searches
: bytes
| type
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| int
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searches
: bytes
| name
: int64 4
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| searches
: bytes
| data
: int64 5
: null
: bytes
| YES
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| sess_id
: int64 1
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| expires_at
: int64 2
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| ip
: int64 3
: null
: bytes
| NO
: bytes
| varchar(40)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| session
: bytes
| vars
: int64 4
: null
: bytes
| NO
: bytes
| mediumtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| system
: bytes
| name
: int64 1
: null
: bytes
| NO
: bytes
| varchar(64)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| system
: bytes
| value
: int64 2
: null
: bytes
| YES
: bytes
| mediumtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| upload_id
: int64 1
: null
: bytes
| NO
: bytes
| varchar(64)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| session_id
: int64 2
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| group
: int64 3
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| metadata
: int64 4
: null
: bytes
| NO
: bytes
| mediumtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| uploads
: bytes
| created
: int64 5
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| user_id
: int64 1
: null
: bytes
| NO
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| auto_increment
: null
< row
: bytes
| users
: bytes
| username
: int64 2
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| mail_host
: int64 3
: null
: bytes
| NO
: bytes
| varchar(128)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| created
: int64 4
: bytes
| '1000-01-01 00:00:00'
: bytes
| NO
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| last_login
: int64 5
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| failed_login
: int64 6
: null
: bytes
| YES
: bytes
| datetime
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| failed_login_counter
: int64 7
: null
: bytes
| YES
: bytes
| int unsigned
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| language
: int64 8
: null
: bytes
| YES
: bytes
| varchar(16)
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
< row
: bytes
| users
: bytes
| preferences
: int64 9
: null
: bytes
| YES
: bytes
| longtext
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
: bytes
| ''
: bytes
|
: bytes
|
: null
> query
|
| 			SELECT
| 				tc.TABLE_NAME,
| 				cc.CONSTRAINT_NAME,
| 				cc.CHECK_CLAUSE,
| 				tc.ENFORCED
| 			FROM information_schema.CHECK_CONSTRAINTS cc
| 				JOIN information_schema.TABLE_CONSTRAINTS tc
| 					ON cc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
| 					AND cc.CONSTRAINT_SCHEMA = tc.CONSTRAINT_SCHEMA
| 			WHERE tc.CONSTRAINT_TYPE = 'CHECK' AND tc.TABLE_SCHEMA = ?
| 			ORDER BY tc.TABLE_NAME, cc.CONSTRAINT_NAME
| 		
: string
| sdl_rw_roundcube_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_rw_roundcube_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
> query
|
| 	SELECT
| 		TRIGGER_NAME,
| 		EVENT_OBJECT_TABLE,
| 		EVENT_MANIPULATION,
| 		ACTION_TIMING,
| 		ACTION_STATEMENT,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.TRIGGERS
| 	WHERE EVENT_OBJECT_SCHEMA = ?
| 	ORDER BY EVENT_OBJECT_TABLE ASC, EVENT_MANIPULATION ASC, ACTION_TIMING ASC, ACTION_ORDER ASC;
| 	
: string
| sdl_rw_roundcube_2
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
> query
|
| 	SELECT
| 		EVENT_NAME,
| 		TIME_ZONE,
| 		SQL_MODE,
| 		CHARACTER_SET_CLIENT,
| 		COLLATION_CONNECTION
| 	FROM INFORMATION_SCHEMA.EVENTS
| 	WHERE EVENT_SCHEMA = ?
| 	ORDER BY EVENT_NAME ASC;
| 	
: string
| sdl_rw_roundcube_2
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
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
| sdl_rw_roundcube_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| cache
: bytes
| user_id_fk_cache
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache_index
: bytes
| user_id_fk_cache_index
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache_thread
: bytes
| user_id_fk_cache_thread
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| cache_messages
: bytes
| user_id_fk_cache_messages
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| collected_addresses
: bytes
| user_id_fk_collected_addresses
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contacts
: bytes
| user_id_fk_contacts
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contactgroups
: bytes
| user_id_fk_contactgroups
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contactgroupmembers
: bytes
| contactgroup_id_fk_contactgroups
: bytes
| contactgroups
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| contactgroupmembers
: bytes
| contact_id_fk_contacts
: bytes
| contacts
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| identities
: bytes
| user_id_fk_identities
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| responses
: bytes
| user_id_fk_responses
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| dictionary
: bytes
| user_id_fk_dictionary
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| searches
: bytes
| user_id_fk_searches
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| filestore
: bytes
| user_id_fk_filestore
: bytes
| users
: bytes
| CASCADE
: bytes
| CASCADE
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
| sdl_rw_roundcube_2
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| cache
: bytes
| user_id_fk_cache
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| cache_index
: bytes
| user_id_fk_cache_index
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| cache_messages
: bytes
| user_id_fk_cache_messages
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| cache_thread
: bytes
| user_id_fk_cache_thread
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| collected_addresses
: bytes
| user_id_fk_collected_addresses
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| contactgroupmembers
: bytes
| contact_id_fk_contacts
: bytes
| contact_id
: bytes
| sdl_rw_roundcube_2
: bytes
| contact_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| contactgroupmembers
: bytes
| contactgroup_id_fk_contactgroups
: bytes
| contactgroup_id
: bytes
| sdl_rw_roundcube_2
: bytes
| contactgroup_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| contactgroups
: bytes
| user_id_fk_contactgroups
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| contacts
: bytes
| user_id_fk_contacts
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| dictionary
: bytes
| user_id_fk_dictionary
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| filestore
: bytes
| user_id_fk_filestore
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| identities
: bytes
| user_id_fk_identities
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| responses
: bytes
| user_id_fk_responses
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
< row
: bytes
| sdl_rw_roundcube_2
: bytes
| searches
: bytes
| user_id_fk_searches
: bytes
| user_id
: bytes
| sdl_rw_roundcube_2
: bytes
| user_id
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			PARTITION_NAME,
| 			SUBPARTITION_NAME,
| 			PARTITION_METHOD,
| 			SUBPARTITION_METHOD,
| 			PARTITION_EXPRESSION,
| 			SUBPARTITION_EXPRESSION,
| 			PARTITION_DESCRIPTION,
| 			PARTITION_COMMENT
| 		FROM INFORMATION_SCHEMA.PARTITIONS
| 		WHERE TABLE_SCHEMA = ? AND PARTITION_NAME IS NOT NULL
| 		ORDER BY TABLE_NAME ASC, PARTITION_ORDINAL_POSITION ASC, SUBPARTITION_ORDINAL_POSITION ASC;
| 	
: string
| sdl_rw_roundcube_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "PARTITION_NAME" "VARCHAR"
< column "SUBPARTITION_NAME" "VARCHAR"
< column "PARTITION_METHOD" "VARCHAR"
< column "SUBPARTITION_METHOD" "VARCHAR"
< column "PARTITION_EXPRESSION" "VARCHAR"
< column "SUBPARTITION_EXPRESSION" "VARCHAR"
< column "PARTITION_DESCRIPTION" "TEXT"
< column "PARTITION_COMMENT" "TEXT" "string"
> query
|
| 		SELECT
| 			ROUTINE_NAME,
| 			ROUTINE_TYPE,
| 			SQL_MODE,
| 			CHARACTER_SET_CLIENT,
| 			COLLATION_CONNECTION,
| 			DATABASE_COLLATION,
| 			IFNULL(ROUTINE_COMMENT, '')
| 		FROM
| 			INFORMATION_SCHEMA.ROUTINES
| 		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE IN ('FUNCTION', 'PROCEDURE')
| 		ORDER BY ROUTINE_TYPE, ROUTINE_NAME;
| 	
: string
| sdl_rw_roundcube_2
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
> query
|
| 		SELECT
| 			TABLES.TABLE_NAME,
| 			TABLES.TABLE_TYPE,
| 			IFNULL(TABLES.ENGINE, ''),
| 			IFNULL(TABLES.TABLE_COLLATION, ''),
| 			IFNULL(TABLES.TABLE_ROWS, 0),
| 			IFNULL(TABLES.DATA_LENGTH, 0),
| 			IFNULL(TABLES.INDEX_LENGTH, 0),
| 			IFNULL(TABLES.DATA_FREE, 0),
| 			IFNULL(TABLES.CREATE_OPTIONS, ''),
| 			QUOTE(IFNULL(TABLES.TABLE_COMMENT, '')),
| 			IFNULL(CCSA.CHARACTER_SET_NAME, '')
| 		FROM information_schema.TABLES TABLES
| 		LEFT JOIN information_schema.COLLATION_CHARACTER_SET_APPLICABILITY CCSA
| 		ON TABLES.TABLE_COLLATION = CCSA.COLLATION_NAME
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_rw_roundcube_2
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "TABLE_TYPE" "ENUM" "string"
< column "IFNULL(TABLES.ENGINE, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_COLLATION, '')" "VARCHAR" "string"
< column "IFNULL(TABLES.TABLE_ROWS, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.INDEX_LENGTH, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.DATA_FREE, 0)" "DECIMAL" "string" precision 21 scale 0
< column "IFNULL(TABLES.CREATE_OPTIONS, '')" "VARCHAR" "string"
< column "QUOTE(IFNULL(TABLES.TABLE_COMMENT, ''))" "TEXT"
< column "IFNULL(CCSA.CHARACTER_SET_NAME, '')" "VARCHAR" "string"
< row
: bytes
| cache
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_index
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_messages
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_shared
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| cache_thread
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| collected_addresses
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| contactgroupmembers
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| contactgroups
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| contacts
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| dictionary
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| filestore
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| identities
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 32768
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| responses
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| searches
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| session
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| system
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 0
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| uploads
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| users
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_unicode_ci
: bytes
| 0
: bytes
| 16384
: bytes
| 16384
: bytes
| 0
: bytes
| row_format=DYNAMIC
: bytes
| ''
: bytes
| utf8mb4
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sdl_rw_roundcube_2
< result
< column "NAME" "VARCHAR" "string"
> query
|
| 		SELECT
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE SCHEMA_NAME = ?
: string
| sdl_rw_roundcube_2
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
