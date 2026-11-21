-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSDLStressLargeSchemaIdempotence default#1
> begin
< ok
> exec
|
| CREATE TABLE department (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	name VARCHAR(100) NOT NULL,
| 	manager_id INT NULL,
| 	parent_id INT NULL,
| 	budget DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY uk_dept_name (name),
| 	CONSTRAINT fk_dept_parent FOREIGN KEY (parent_id) REFERENCES department (id) ON DELETE SET NULL,
| 	CONSTRAINT chk_budget CHECK (budget >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE employee (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	dept_id INT NOT NULL,
| 	first_name VARCHAR(50) NOT NULL,
| 	last_name VARCHAR(50) NOT NULL,
| 	full_name VARCHAR(101) AS (CONCAT(first_name, ' ', last_name)) VIRTUAL,
| 	email VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
| 	active BOOLEAN NOT NULL DEFAULT TRUE,
| 	salary DECIMAL(10,2) NOT NULL DEFAULT 0.00,
| 	hired_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
| 	rank_kind ENUM('junior','mid','senior','staff') NOT NULL DEFAULT 'junior',
| 	tags SET('remote','contractor','lead') NOT NULL DEFAULT '',
| 	notes TEXT,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY uk_emp_email (email),
| 	KEY idx_emp_name (last_name, first_name),
| 	KEY idx_emp_notes_prefix (notes(20)),
| 	FULLTEXT KEY ft_emp_notes (notes),
| 	CONSTRAINT chk_salary CHECK (salary >= 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE project (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	dept_id INT NOT NULL,
| 	code VARCHAR(20) NOT NULL,
| 	title VARCHAR(200) NOT NULL,
| 	budget DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	spent DECIMAL(14,2) NOT NULL DEFAULT 0.00,
| 	remaining DECIMAL(14,2) AS (budget - spent) STORED,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	PRIMARY KEY (id),
| 	UNIQUE KEY uk_proj_code (code),
| 	KEY idx_proj_dept (dept_id),
| 	CONSTRAINT fk_proj_dept FOREIGN KEY (dept_id) REFERENCES department (id) ON DELETE CASCADE ON UPDATE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE assignment (
| 	employee_id INT NOT NULL,
| 	project_id INT NOT NULL,
| 	role VARCHAR(50) NOT NULL DEFAULT 'member',
| 	allocation DECIMAL(5,2) NOT NULL DEFAULT 100.00,
| 	assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	PRIMARY KEY (employee_id, project_id),
| 	KEY idx_assign_project (project_id),
| 	CONSTRAINT fk_assign_emp FOREIGN KEY (employee_id) REFERENCES employee (id) ON DELETE CASCADE,
| 	CONSTRAINT fk_assign_proj FOREIGN KEY (project_id) REFERENCES project (id) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE event_log (
| 	id BIGINT NOT NULL AUTO_INCREMENT,
| 	occurred_on DATE NOT NULL,
| 	kind VARCHAR(40) NOT NULL,
| 	payload JSON,
| 	PRIMARY KEY (id, occurred_on)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
| PARTITION BY RANGE (YEAR(occurred_on)) (
| 	PARTITION p2023 VALUES LESS THAN (2024),
| 	PARTITION p2024 VALUES LESS THAN (2025),
| 	PARTITION p_future VALUES LESS THAN MAXVALUE
| );
< ok 0 0
> exec
|
|
| CREATE TABLE metric (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	bucket INT NOT NULL,
| 	name VARCHAR(60) NOT NULL,
| 	value DOUBLE NOT NULL DEFAULT 0,
| 	PRIMARY KEY (id, bucket)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
| PARTITION BY HASH (bucket) PARTITIONS 4;
< ok 0 0
> exec
|
|
| CREATE TABLE audit_trail (
| 	id BIGINT NOT NULL AUTO_INCREMENT,
| 	emp_id INT NOT NULL,
| 	old_salary DECIMAL(10,2),
| 	new_salary DECIMAL(10,2),
| 	changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	PRIMARY KEY (id)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE category (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	parent_id INT NULL,
| 	name VARCHAR(80) NOT NULL,
| 	PRIMARY KEY (id),
| 	KEY idx_cat_parent (parent_id),
| 	CONSTRAINT fk_cat_parent FOREIGN KEY (parent_id) REFERENCES category (id) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE doc (
| 	id INT NOT NULL AUTO_INCREMENT,
| 	category_id INT NOT NULL,
| 	owner_id INT NOT NULL,
| 	body MEDIUMTEXT,
| 	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
| 	PRIMARY KEY (id),
| 	KEY idx_doc_cat (category_id),
| 	KEY idx_doc_owner (owner_id),
| 	CONSTRAINT fk_doc_cat FOREIGN KEY (category_id) REFERENCES category (id),
| 	CONSTRAINT fk_doc_owner FOREIGN KEY (owner_id) REFERENCES employee (id) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE TABLE tag_map (
| 	doc_id INT NOT NULL,
| 	tag VARCHAR(40) NOT NULL,
| 	weight INT NOT NULL DEFAULT 1,
| 	PRIMARY KEY (doc_id, tag),
| 	CONSTRAINT fk_tagmap_doc FOREIGN KEY (doc_id) REFERENCES doc (id) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
< ok 0 0
> exec
|
|
| CREATE VIEW v_active_employees AS
| SELECT e.id, e.full_name, e.email, e.salary, d.name AS dept_name
| FROM employee e JOIN department d ON e.dept_id = d.id
| WHERE e.active = TRUE;
< ok 0 0
> exec
|
|
| CREATE VIEW v_dept_payroll AS
| SELECT dept_name, COUNT(*) AS headcount, SUM(salary) AS payroll
| FROM v_active_employees
| GROUP BY dept_name;
< ok 0 0
> exec
|
|
| CREATE FUNCTION emp_annual_salary(monthly DECIMAL(10,2)) RETURNS DECIMAL(12,2) DETERMINISTIC
| RETURN monthly * 12;
< ok 0 0
> exec
|
|
| CREATE FUNCTION dept_headcount(d_id INT) RETURNS INT READS SQL DATA
| RETURN (SELECT COUNT(*) FROM employee WHERE dept_id = d_id);
< ok 0 0
> exec
|
|
| CREATE PROCEDURE give_raise(IN emp INT, IN pct DECIMAL(5,2))
| BEGIN
| 	UPDATE employee SET salary = salary * (1 + pct / 100) WHERE id = emp;
| END;
< ok 0 0
> exec
|
|
| CREATE PROCEDURE close_project(IN proj INT)
| BEGIN
| 	DELETE FROM assignment WHERE project_id = proj;
| 	UPDATE project SET spent = budget WHERE id = proj;
| END;
< ok 0 0
> exec
|
|
| CREATE TRIGGER trg_emp_before_ins BEFORE INSERT ON employee FOR EACH ROW
| SET NEW.email = LOWER(NEW.email);
< ok 0 0
> exec
|
|
| CREATE TRIGGER trg_emp_after_upd AFTER UPDATE ON employee FOR EACH ROW
| BEGIN
| 	IF OLD.salary <> NEW.salary THEN
| 		INSERT INTO audit_trail (emp_id, old_salary, new_salary) VALUES (NEW.id, OLD.salary, NEW.salary);
| 	END IF;
| END;
< ok 0 0
> exec
|
|
| CREATE TRIGGER trg_proj_before_ins BEFORE INSERT ON project FOR EACH ROW
| SET NEW.created_at = NOW();
< ok 0 0
> exec
|
|
| CREATE EVENT IF NOT EXISTS ev_daily_metric ON SCHEDULE EVERY 1 DAY
| DO INSERT INTO metric (bucket, name, value) VALUES (0, 'daily', 1);
< ok 0 0
> exec
|
|
| CREATE VIEW v_project_load AS
| WITH per_proj AS (
| 	SELECT project_id, COUNT(*) AS members, SUM(allocation) AS total_alloc
| 	FROM assignment GROUP BY project_id
| )
| SELECT p.code, p.title, pp.members, pp.total_alloc
| FROM project p JOIN per_proj pp ON p.id = pp.project_id;
< ok 0 0
> exec
|
|
| ALTER TABLE department ADD CONSTRAINT fk_dept_manager FOREIGN KEY (manager_id) REFERENCES employee (id) ON DELETE SET NULL;
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
| sdl_stress_big_1
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
| assignment
: bytes
| idx_assign_project
: bytes
| project_id
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
| assignment
: bytes
| PRIMARY
: bytes
| employee_id
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
| assignment
: bytes
| PRIMARY
: bytes
| project_id
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
| audit_trail
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
| category
: bytes
| idx_cat_parent
: bytes
| parent_id
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
| category
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
| department
: bytes
| fk_dept_manager
: bytes
| manager_id
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
| department
: bytes
| fk_dept_parent
: bytes
| parent_id
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
| department
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
| department
: bytes
| uk_dept_name
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
| doc
: bytes
| idx_doc_cat
: bytes
| category_id
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
| doc
: bytes
| idx_doc_owner
: bytes
| owner_id
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
| doc
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
| employee
: bytes
| ft_emp_notes
: bytes
| notes
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| employee
: bytes
| idx_emp_name
: bytes
| last_name
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
| employee
: bytes
| idx_emp_name
: bytes
| first_name
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
| employee
: bytes
| idx_emp_notes_prefix
: bytes
| notes
: bytes
| A
: int64 20
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
| employee
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
| employee
: bytes
| uk_emp_email
: bytes
| email
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
| event_log
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
| event_log
: bytes
| PRIMARY
: bytes
| occurred_on
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
| metric
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
| metric
: bytes
| PRIMARY
: bytes
| bucket
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
| project
: bytes
| idx_proj_dept
: bytes
| dept_id
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
| project
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
| project
: bytes
| uk_proj_code
: bytes
| code
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
| tag_map
: bytes
| PRIMARY
: bytes
| doc_id
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
| tag_map
: bytes
| PRIMARY
: bytes
| tag
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
| sdl_stress_big_1
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
| assignment
: bytes
| employee_id
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
| assignment
: bytes
| project_id
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
| assignment
: bytes
| role
: int64 3
: bytes
| 'member'
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
| assignment
: bytes
| allocation
: int64 4
: bytes
| '100.00'
: bytes
| NO
: bytes
| decimal(5,2)
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
| assignment
: bytes
| assigned_at
: int64 5
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| audit_trail
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| bigint
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
| audit_trail
: bytes
| emp_id
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
| audit_trail
: bytes
| old_salary
: int64 3
: null
: bytes
| YES
: bytes
| decimal(10,2)
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
| audit_trail
: bytes
| new_salary
: int64 4
: null
: bytes
| YES
: bytes
| decimal(10,2)
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
| audit_trail
: bytes
| changed_at
: int64 5
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| category
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| category
: bytes
| parent_id
: int64 2
: null
: bytes
| YES
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
| category
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(80)
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
| department
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| department
: bytes
| name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(100)
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
| department
: bytes
| manager_id
: int64 3
: null
: bytes
| YES
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
| department
: bytes
| parent_id
: int64 4
: null
: bytes
| YES
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
| department
: bytes
| budget
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
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
| doc
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| doc
: bytes
| category_id
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
| doc
: bytes
| owner_id
: int64 3
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
| doc
: bytes
| body
: int64 4
: null
: bytes
| YES
: bytes
| mediumtext
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
| doc
: bytes
| created_at
: int64 5
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| employee
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| employee
: bytes
| dept_id
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
| employee
: bytes
| first_name
: int64 3
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
| employee
: bytes
| last_name
: int64 4
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
| employee
: bytes
| full_name
: int64 5
: null
: bytes
| YES
: bytes
| varchar(101)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| concat(`first_name`,_utf8mb4\' \',`last_name`)
: bytes
| VIRTUAL GENERATED
: null
< row
: bytes
| employee
: bytes
| email
: int64 6
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
| employee
: bytes
| active
: int64 7
: bytes
| '1'
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
| employee
: bytes
| salary
: int64 8
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(10,2)
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
| employee
: bytes
| hired_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| employee
: bytes
| updated_at
: int64 10
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| employee
: bytes
| rank_kind
: int64 11
: bytes
| 'junior'
: bytes
| NO
: bytes
| enum('junior','mid','senior','staff')
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
| employee
: bytes
| tags
: int64 12
: bytes
| ''
: bytes
| NO
: bytes
| set('remote','contractor','lead')
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
| employee
: bytes
| notes
: int64 13
: null
: bytes
| YES
: bytes
| text
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
| event_log
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| bigint
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
| event_log
: bytes
| occurred_on
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
| event_log
: bytes
| kind
: int64 3
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
| event_log
: bytes
| payload
: int64 4
: null
: bytes
| YES
: bytes
| json
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
| metric
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| metric
: bytes
| bucket
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
| metric
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(60)
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
| metric
: bytes
| value
: int64 4
: bytes
| '0'
: bytes
| NO
: bytes
| double
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
| project
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| project
: bytes
| dept_id
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
| project
: bytes
| code
: int64 3
: null
: bytes
| NO
: bytes
| varchar(20)
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
| project
: bytes
| title
: int64 4
: null
: bytes
| NO
: bytes
| varchar(200)
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
| project
: bytes
| budget
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
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
| project
: bytes
| spent
: int64 6
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
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
| project
: bytes
| remaining
: int64 7
: null
: bytes
| YES
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| (`budget` - `spent`)
: bytes
| STORED GENERATED
: null
< row
: bytes
| project
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| tag_map
: bytes
| doc_id
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
| tag_map
: bytes
| tag
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
| tag_map
: bytes
| weight
: int64 3
: bytes
| '1'
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
| v_active_employees
: bytes
| id
: int64 1
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
| v_active_employees
: bytes
| full_name
: int64 2
: null
: bytes
| YES
: bytes
| varchar(101)
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
| v_active_employees
: bytes
| email
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
| v_active_employees
: bytes
| salary
: int64 4
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(10,2)
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
| v_active_employees
: bytes
| dept_name
: int64 5
: null
: bytes
| NO
: bytes
| varchar(100)
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
| v_dept_payroll
: bytes
| dept_name
: int64 1
: null
: bytes
| NO
: bytes
| varchar(100)
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
| v_dept_payroll
: bytes
| headcount
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| bigint
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
| v_dept_payroll
: bytes
| payroll
: int64 3
: null
: bytes
| YES
: bytes
| decimal(32,2)
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
| v_project_load
: bytes
| code
: int64 1
: null
: bytes
| NO
: bytes
| varchar(20)
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
| v_project_load
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(200)
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
| v_project_load
: bytes
| members
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| bigint
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
| v_project_load
: bytes
| total_alloc
: int64 4
: null
: bytes
| YES
: bytes
| decimal(27,2)
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
| sdl_stress_big_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
< row
: bytes
| department
: bytes
| chk_budget
: bytes
| (`budget` >= 0)
: bytes
| YES
< row
: bytes
| employee
: bytes
| chk_salary
: bytes
| (`salary` >= 0)
: bytes
| YES
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_stress_big_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_active_employees
: bytes
| select `e`.`id` AS `id`,`e`.`full_name` AS `full_name`,`e`.`email` AS `email`,`e`.`salary` AS `salary`,`d`.`name` AS `dept_name` from (`sdl_stress_big_1`.`employee` `e` join `sdl_stress_big_1`.`department` `d` on((`e`.`dept_id` = `d`.`id`))) where (`e`.`active` = true)
< row
: bytes
| v_dept_payroll
: bytes
| select `sdl_stress_big_1`.`v_active_employees`.`dept_name` AS `dept_name`,count(0) AS `headcount`,sum(`sdl_stress_big_1`.`v_active_employees`.`salary`) AS `payroll` from `sdl_stress_big_1`.`v_active_employees` group by `sdl_stress_big_1`.`v_active_employees`.`dept_name`
< row
: bytes
| v_project_load
: bytes
| with `per_proj` as (select `sdl_stress_big_1`.`assignment`.`project_id` AS `project_id`,count(0) AS `members`,sum(`sdl_stress_big_1`.`assignment`.`allocation`) AS `total_alloc` from `sdl_stress_big_1`.`assignment` group by `sdl_stress_big_1`.`assignment`.`project_id`) select `p`.`code` AS `code`,`p`.`title` AS `title`,`pp`.`members` AS `members`,`pp`.`total_alloc` AS `total_alloc` from (`sdl_stress_big_1`.`project` `p` join `per_proj` `pp` on((`p`.`id` = `pp`.`project_id`)))
> query
| SHOW CREATE VIEW `sdl_stress_big_1`.`v_active_employees`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_active_employees
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_active_employees` AS select `e`.`id` AS `id`,`e`.`full_name` AS `full_name`,`e`.`email` AS `email`,`e`.`salary` AS `salary`,`d`.`name` AS `dept_name` from (`employee` `e` join `department` `d` on((`e`.`dept_id` = `d`.`id`))) where (`e`.`active` = true)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sdl_stress_big_1`.`v_dept_payroll`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_dept_payroll
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_dept_payroll` AS select `v_active_employees`.`dept_name` AS `dept_name`,count(0) AS `headcount`,sum(`v_active_employees`.`salary`) AS `payroll` from `v_active_employees` group by `v_active_employees`.`dept_name`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sdl_stress_big_1`.`v_project_load`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_project_load
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_project_load` AS with `per_proj` as (select `assignment`.`project_id` AS `project_id`,count(0) AS `members`,sum(`assignment`.`allocation`) AS `total_alloc` from `assignment` group by `assignment`.`project_id`) select `p`.`code` AS `code`,`p`.`title` AS `title`,`pp`.`members` AS `members`,`pp`.`total_alloc` AS `total_alloc` from (`project` `p` join `per_proj` `pp` on((`p`.`id` = `pp`.`project_id`)))
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
| sdl_stress_big_1
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| trg_emp_before_ins
: bytes
| employee
: bytes
| INSERT
: bytes
| BEFORE
: bytes
| SET NEW.email = LOWER(NEW.email)
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
< row
: bytes
| trg_emp_after_upd
: bytes
| employee
: bytes
| UPDATE
: bytes
| AFTER
: bytes
| BEGIN
| 	IF OLD.salary <> NEW.salary THEN
| 		INSERT INTO audit_trail (emp_id, old_salary, new_salary) VALUES (NEW.id, OLD.salary, NEW.salary);
| 	END IF;
| END
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
< row
: bytes
| trg_proj_before_ins
: bytes
| project
: bytes
| INSERT
: bytes
| BEFORE
: bytes
| SET NEW.created_at = NOW()
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sdl_stress_big_1
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_metric
: bytes
| SYSTEM
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `sdl_stress_big_1`.`ev_daily_metric`
< result
< column "Event" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "time_zone" "VARCHAR" "string"
< column "Create Event" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| ev_daily_metric
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_metric` ON SCHEDULE EVERY 1 DAY STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO INSERT INTO metric (bucket, name, value) VALUES (0, 'daily', 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
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
| sdl_stress_big_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| assignment
: bytes
| fk_assign_emp
: bytes
| employee
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| assignment
: bytes
| fk_assign_proj
: bytes
| project
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| category
: bytes
| fk_cat_parent
: bytes
| category
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| department
: bytes
| fk_dept_manager
: bytes
| employee
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| department
: bytes
| fk_dept_parent
: bytes
| department
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| doc
: bytes
| fk_doc_cat
: bytes
| category
: bytes
| NO ACTION
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| doc
: bytes
| fk_doc_owner
: bytes
| employee
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| project
: bytes
| fk_proj_dept
: bytes
| department
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| tag_map
: bytes
| fk_tagmap_doc
: bytes
| doc
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
| sdl_stress_big_1
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_stress_big_1
: bytes
| assignment
: bytes
| fk_assign_emp
: bytes
| employee_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| assignment
: bytes
| fk_assign_proj
: bytes
| project_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| category
: bytes
| fk_cat_parent
: bytes
| parent_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| department
: bytes
| fk_dept_manager
: bytes
| manager_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| department
: bytes
| fk_dept_parent
: bytes
| parent_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| doc
: bytes
| fk_doc_cat
: bytes
| category_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| doc
: bytes
| fk_doc_owner
: bytes
| owner_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| project
: bytes
| fk_proj_dept
: bytes
| dept_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| tag_map
: bytes
| fk_tagmap_doc
: bytes
| doc_id
: bytes
| sdl_stress_big_1
: bytes
| id
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
| sdl_stress_big_1
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
< row
: bytes
| event_log
: bytes
| p2023
: null
: bytes
| RANGE
: null
: bytes
| year(`occurred_on`)
: null
: bytes
| 2024
: bytes
|
< row
: bytes
| event_log
: bytes
| p2024
: null
: bytes
| RANGE
: null
: bytes
| year(`occurred_on`)
: null
: bytes
| 2025
: bytes
|
< row
: bytes
| event_log
: bytes
| p_future
: null
: bytes
| RANGE
: null
: bytes
| year(`occurred_on`)
: null
: bytes
| MAXVALUE
: bytes
|
< row
: bytes
| metric
: bytes
| p0
: null
: bytes
| HASH
: null
: bytes
| `bucket`
: null
: null
: bytes
|
< row
: bytes
| metric
: bytes
| p1
: null
: bytes
| HASH
: null
: bytes
| `bucket`
: null
: null
: bytes
|
< row
: bytes
| metric
: bytes
| p2
: null
: bytes
| HASH
: null
: bytes
| `bucket`
: null
: null
: bytes
|
< row
: bytes
| metric
: bytes
| p3
: null
: bytes
| HASH
: null
: bytes
| `bucket`
: null
: null
: bytes
|
> query
| SHOW CREATE TABLE `sdl_stress_big_1`.`event_log`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| event_log
: bytes
| CREATE TABLE `event_log` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `occurred_on` date NOT NULL,
|   `kind` varchar(40) NOT NULL,
|   `payload` json DEFAULT NULL,
|   PRIMARY KEY (`id`,`occurred_on`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`occurred_on`))
| (PARTITION p2023 VALUES LESS THAN (2024) ENGINE = InnoDB,
|  PARTITION p2024 VALUES LESS THAN (2025) ENGINE = InnoDB,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> query
| SHOW CREATE TABLE `sdl_stress_big_1`.`metric`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metric
: bytes
| CREATE TABLE `metric` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `bucket` int NOT NULL,
|   `name` varchar(60) NOT NULL,
|   `value` double NOT NULL DEFAULT '0',
|   PRIMARY KEY (`id`,`bucket`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`bucket`)
| PARTITIONS 4 */
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
| sdl_stress_big_1
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| dept_headcount
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| emp_annual_salary
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| close_project
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| give_raise
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
> query
| SHOW CREATE FUNCTION `sdl_stress_big_1`.`dept_headcount`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| dept_headcount
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `dept_headcount`(d_id INT) RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM employee WHERE dept_id = d_id)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION `sdl_stress_big_1`.`emp_annual_salary`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| emp_annual_salary
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `emp_annual_salary`(monthly DECIMAL(10,2)) RETURNS decimal(12,2)
|     DETERMINISTIC
| RETURN monthly * 12
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sdl_stress_big_1`.`close_project`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| close_project
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `close_project`(IN proj INT)
| BEGIN
| 	DELETE FROM assignment WHERE project_id = proj;
| 	UPDATE project SET spent = budget WHERE id = proj;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sdl_stress_big_1`.`give_raise`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| give_raise
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `give_raise`(IN emp INT, IN pct DECIMAL(5,2))
| BEGIN
| 	UPDATE employee SET salary = salary * (1 + pct / 100) WHERE id = emp;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
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
| sdl_stress_big_1
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
| assignment
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
| audit_trail
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
| category
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
| department
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
| 49152
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
| doc
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
| 32768
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
| employee
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
| 81920
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
| event_log
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 49152
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| metric
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 65536
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| project
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
| 32768
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
| tag_map
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
| v_active_employees
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
| v_dept_payroll
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
| v_project_load
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
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sdl_stress_big_1
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
| sdl_stress_big_1
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci

== TestSDLStressLargeSchemaIdempotence default#2
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
| sdl_stress_big_1
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
| assignment
: bytes
| idx_assign_project
: bytes
| project_id
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
| assignment
: bytes
| PRIMARY
: bytes
| employee_id
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
| assignment
: bytes
| PRIMARY
: bytes
| project_id
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
| audit_trail
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
| category
: bytes
| idx_cat_parent
: bytes
| parent_id
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
| category
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
| department
: bytes
| fk_dept_manager
: bytes
| manager_id
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
| department
: bytes
| fk_dept_parent
: bytes
| parent_id
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
| department
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
| department
: bytes
| uk_dept_name
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
| doc
: bytes
| idx_doc_cat
: bytes
| category_id
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
| doc
: bytes
| idx_doc_owner
: bytes
| owner_id
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
| doc
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
| employee
: bytes
| ft_emp_notes
: bytes
| notes
: null
: int64 -1
: null
: int64 1
: bytes
| FULLTEXT
: int64 1
: int64 1
: bytes
|
< row
: bytes
| employee
: bytes
| idx_emp_name
: bytes
| last_name
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
| employee
: bytes
| idx_emp_name
: bytes
| first_name
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
| employee
: bytes
| idx_emp_notes_prefix
: bytes
| notes
: bytes
| A
: int64 20
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
| employee
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
| employee
: bytes
| uk_emp_email
: bytes
| email
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
| event_log
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
| event_log
: bytes
| PRIMARY
: bytes
| occurred_on
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
| metric
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
| metric
: bytes
| PRIMARY
: bytes
| bucket
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
| project
: bytes
| idx_proj_dept
: bytes
| dept_id
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
| project
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
| project
: bytes
| uk_proj_code
: bytes
| code
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
| tag_map
: bytes
| PRIMARY
: bytes
| doc_id
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
| tag_map
: bytes
| PRIMARY
: bytes
| tag
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
| sdl_stress_big_1
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
| assignment
: bytes
| employee_id
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
| assignment
: bytes
| project_id
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
| assignment
: bytes
| role
: int64 3
: bytes
| 'member'
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
| assignment
: bytes
| allocation
: int64 4
: bytes
| '100.00'
: bytes
| NO
: bytes
| decimal(5,2)
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
| assignment
: bytes
| assigned_at
: int64 5
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| audit_trail
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| bigint
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
| audit_trail
: bytes
| emp_id
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
| audit_trail
: bytes
| old_salary
: int64 3
: null
: bytes
| YES
: bytes
| decimal(10,2)
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
| audit_trail
: bytes
| new_salary
: int64 4
: null
: bytes
| YES
: bytes
| decimal(10,2)
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
| audit_trail
: bytes
| changed_at
: int64 5
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| category
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| category
: bytes
| parent_id
: int64 2
: null
: bytes
| YES
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
| category
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(80)
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
| department
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| department
: bytes
| name
: int64 2
: null
: bytes
| NO
: bytes
| varchar(100)
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
| department
: bytes
| manager_id
: int64 3
: null
: bytes
| YES
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
| department
: bytes
| parent_id
: int64 4
: null
: bytes
| YES
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
| department
: bytes
| budget
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
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
| doc
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| doc
: bytes
| category_id
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
| doc
: bytes
| owner_id
: int64 3
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
| doc
: bytes
| body
: int64 4
: null
: bytes
| YES
: bytes
| mediumtext
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
| doc
: bytes
| created_at
: int64 5
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| employee
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| employee
: bytes
| dept_id
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
| employee
: bytes
| first_name
: int64 3
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
| employee
: bytes
| last_name
: int64 4
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
| employee
: bytes
| full_name
: int64 5
: null
: bytes
| YES
: bytes
| varchar(101)
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
: bytes
| ''
: bytes
| concat(`first_name`,_utf8mb4\' \',`last_name`)
: bytes
| VIRTUAL GENERATED
: null
< row
: bytes
| employee
: bytes
| email
: int64 6
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
| employee
: bytes
| active
: int64 7
: bytes
| '1'
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
| employee
: bytes
| salary
: int64 8
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(10,2)
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
| employee
: bytes
| hired_at
: int64 9
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| employee
: bytes
| updated_at
: int64 10
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
: null
< row
: bytes
| employee
: bytes
| rank_kind
: int64 11
: bytes
| 'junior'
: bytes
| NO
: bytes
| enum('junior','mid','senior','staff')
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
| employee
: bytes
| tags
: int64 12
: bytes
| ''
: bytes
| NO
: bytes
| set('remote','contractor','lead')
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
| employee
: bytes
| notes
: int64 13
: null
: bytes
| YES
: bytes
| text
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
| event_log
: bytes
| id
: int64 1
: null
: bytes
| NO
: bytes
| bigint
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
| event_log
: bytes
| occurred_on
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
| event_log
: bytes
| kind
: int64 3
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
| event_log
: bytes
| payload
: int64 4
: null
: bytes
| YES
: bytes
| json
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
| metric
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| metric
: bytes
| bucket
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
| metric
: bytes
| name
: int64 3
: null
: bytes
| NO
: bytes
| varchar(60)
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
| metric
: bytes
| value
: int64 4
: bytes
| '0'
: bytes
| NO
: bytes
| double
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
| project
: bytes
| id
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
| auto_increment
: null
< row
: bytes
| project
: bytes
| dept_id
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
| project
: bytes
| code
: int64 3
: null
: bytes
| NO
: bytes
| varchar(20)
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
| project
: bytes
| title
: int64 4
: null
: bytes
| NO
: bytes
| varchar(200)
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
| project
: bytes
| budget
: int64 5
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
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
| project
: bytes
| spent
: int64 6
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(14,2)
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
| project
: bytes
| remaining
: int64 7
: null
: bytes
| YES
: bytes
| decimal(14,2)
: bytes
|
: bytes
|
: bytes
| ''
: bytes
| (`budget` - `spent`)
: bytes
| STORED GENERATED
: null
< row
: bytes
| project
: bytes
| created_at
: int64 8
: bytes
| 'CURRENT_TIMESTAMP'
: bytes
| NO
: bytes
| timestamp
: bytes
|
: bytes
|
: bytes
| ''
: bytes
|
: bytes
| DEFAULT_GENERATED
: null
< row
: bytes
| tag_map
: bytes
| doc_id
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
| tag_map
: bytes
| tag
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
| tag_map
: bytes
| weight
: int64 3
: bytes
| '1'
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
| v_active_employees
: bytes
| id
: int64 1
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
| v_active_employees
: bytes
| full_name
: int64 2
: null
: bytes
| YES
: bytes
| varchar(101)
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
| v_active_employees
: bytes
| email
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
| v_active_employees
: bytes
| salary
: int64 4
: bytes
| '0.00'
: bytes
| NO
: bytes
| decimal(10,2)
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
| v_active_employees
: bytes
| dept_name
: int64 5
: null
: bytes
| NO
: bytes
| varchar(100)
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
| v_dept_payroll
: bytes
| dept_name
: int64 1
: null
: bytes
| NO
: bytes
| varchar(100)
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
| v_dept_payroll
: bytes
| headcount
: int64 2
: bytes
| '0'
: bytes
| NO
: bytes
| bigint
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
| v_dept_payroll
: bytes
| payroll
: int64 3
: null
: bytes
| YES
: bytes
| decimal(32,2)
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
| v_project_load
: bytes
| code
: int64 1
: null
: bytes
| NO
: bytes
| varchar(20)
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
| v_project_load
: bytes
| title
: int64 2
: null
: bytes
| NO
: bytes
| varchar(200)
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
| v_project_load
: bytes
| members
: int64 3
: bytes
| '0'
: bytes
| NO
: bytes
| bigint
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
| v_project_load
: bytes
| total_alloc
: int64 4
: null
: bytes
| YES
: bytes
| decimal(27,2)
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
| sdl_stress_big_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR" "string"
< column "CHECK_CLAUSE" "TEXT" "string"
< column "ENFORCED" "VARCHAR" "string"
< row
: bytes
| department
: bytes
| chk_budget
: bytes
| (`budget` >= 0)
: bytes
| YES
< row
: bytes
| employee
: bytes
| chk_salary
: bytes
| (`salary` >= 0)
: bytes
| YES
> query
|
| 		SELECT
| 			TABLE_NAME,
| 			VIEW_DEFINITION
| 		FROM information_schema.VIEWS
| 		WHERE TABLE_SCHEMA = ?
| 		ORDER BY TABLE_NAME
: string
| sdl_stress_big_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| v_active_employees
: bytes
| select `e`.`id` AS `id`,`e`.`full_name` AS `full_name`,`e`.`email` AS `email`,`e`.`salary` AS `salary`,`d`.`name` AS `dept_name` from (`sdl_stress_big_1`.`employee` `e` join `sdl_stress_big_1`.`department` `d` on((`e`.`dept_id` = `d`.`id`))) where (`e`.`active` = true)
< row
: bytes
| v_dept_payroll
: bytes
| select `sdl_stress_big_1`.`v_active_employees`.`dept_name` AS `dept_name`,count(0) AS `headcount`,sum(`sdl_stress_big_1`.`v_active_employees`.`salary`) AS `payroll` from `sdl_stress_big_1`.`v_active_employees` group by `sdl_stress_big_1`.`v_active_employees`.`dept_name`
< row
: bytes
| v_project_load
: bytes
| with `per_proj` as (select `sdl_stress_big_1`.`assignment`.`project_id` AS `project_id`,count(0) AS `members`,sum(`sdl_stress_big_1`.`assignment`.`allocation`) AS `total_alloc` from `sdl_stress_big_1`.`assignment` group by `sdl_stress_big_1`.`assignment`.`project_id`) select `p`.`code` AS `code`,`p`.`title` AS `title`,`pp`.`members` AS `members`,`pp`.`total_alloc` AS `total_alloc` from (`sdl_stress_big_1`.`project` `p` join `per_proj` `pp` on((`p`.`id` = `pp`.`project_id`)))
> query
| SHOW CREATE VIEW `sdl_stress_big_1`.`v_active_employees`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_active_employees
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_active_employees` AS select `e`.`id` AS `id`,`e`.`full_name` AS `full_name`,`e`.`email` AS `email`,`e`.`salary` AS `salary`,`d`.`name` AS `dept_name` from (`employee` `e` join `department` `d` on((`e`.`dept_id` = `d`.`id`))) where (`e`.`active` = true)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sdl_stress_big_1`.`v_dept_payroll`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_dept_payroll
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_dept_payroll` AS select `v_active_employees`.`dept_name` AS `dept_name`,count(0) AS `headcount`,sum(`v_active_employees`.`salary`) AS `payroll` from `v_active_employees` group by `v_active_employees`.`dept_name`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW `sdl_stress_big_1`.`v_project_load`
< result
< column "View" "VARCHAR" "string"
< column "Create View" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< row
: bytes
| v_project_load
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_project_load` AS with `per_proj` as (select `assignment`.`project_id` AS `project_id`,count(0) AS `members`,sum(`assignment`.`allocation`) AS `total_alloc` from `assignment` group by `assignment`.`project_id`) select `p`.`code` AS `code`,`p`.`title` AS `title`,`pp`.`members` AS `members`,`pp`.`total_alloc` AS `total_alloc` from (`project` `p` join `per_proj` `pp` on((`p`.`id` = `pp`.`project_id`)))
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
| sdl_stress_big_1
< result
< column "TRIGGER_NAME" "VARCHAR" "string"
< column "EVENT_OBJECT_TABLE" "VARCHAR" "string"
< column "EVENT_MANIPULATION" "ENUM" "string"
< column "ACTION_TIMING" "ENUM" "string"
< column "ACTION_STATEMENT" "TEXT" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| trg_emp_before_ins
: bytes
| employee
: bytes
| INSERT
: bytes
| BEFORE
: bytes
| SET NEW.email = LOWER(NEW.email)
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
< row
: bytes
| trg_emp_after_upd
: bytes
| employee
: bytes
| UPDATE
: bytes
| AFTER
: bytes
| BEGIN
| 	IF OLD.salary <> NEW.salary THEN
| 		INSERT INTO audit_trail (emp_id, old_salary, new_salary) VALUES (NEW.id, OLD.salary, NEW.salary);
| 	END IF;
| END
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
< row
: bytes
| trg_proj_before_ins
: bytes
| project
: bytes
| INSERT
: bytes
| BEFORE
: bytes
| SET NEW.created_at = NOW()
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| sdl_stress_big_1
< result
< column "EVENT_NAME" "VARCHAR" "string"
< column "TIME_ZONE" "VARCHAR" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< row
: bytes
| ev_daily_metric
: bytes
| SYSTEM
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE EVENT `sdl_stress_big_1`.`ev_daily_metric`
< result
< column "Event" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "time_zone" "VARCHAR" "string"
< column "Create Event" "VARCHAR" "string"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| ev_daily_metric
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
: bytes
| CREATE DEFINER=`root`@`%` EVENT `ev_daily_metric` ON SCHEDULE EVERY 1 DAY STARTS '2026-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE DO INSERT INTO metric (bucket, name, value) VALUES (0, 'daily', 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
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
| sdl_stress_big_1
< result
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR" "string"
< column "DELETE_RULE" "ENUM" "string"
< column "UPDATE_RULE" "ENUM" "string"
< column "MATCH_OPTION" "ENUM" "string"
< row
: bytes
| assignment
: bytes
| fk_assign_emp
: bytes
| employee
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| assignment
: bytes
| fk_assign_proj
: bytes
| project
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| category
: bytes
| fk_cat_parent
: bytes
| category
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| department
: bytes
| fk_dept_manager
: bytes
| employee
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| department
: bytes
| fk_dept_parent
: bytes
| department
: bytes
| SET NULL
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| doc
: bytes
| fk_doc_cat
: bytes
| category
: bytes
| NO ACTION
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| doc
: bytes
| fk_doc_owner
: bytes
| employee
: bytes
| CASCADE
: bytes
| NO ACTION
: bytes
| NONE
< row
: bytes
| project
: bytes
| fk_proj_dept
: bytes
| department
: bytes
| CASCADE
: bytes
| CASCADE
: bytes
| NONE
< row
: bytes
| tag_map
: bytes
| fk_tagmap_doc
: bytes
| doc
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
| sdl_stress_big_1
< result
< column "TABLE_SCHEMA" "VARCHAR" "string"
< column "TABLE_NAME" "VARCHAR" "string"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| sdl_stress_big_1
: bytes
| assignment
: bytes
| fk_assign_emp
: bytes
| employee_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| assignment
: bytes
| fk_assign_proj
: bytes
| project_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| category
: bytes
| fk_cat_parent
: bytes
| parent_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| department
: bytes
| fk_dept_manager
: bytes
| manager_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| department
: bytes
| fk_dept_parent
: bytes
| parent_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| doc
: bytes
| fk_doc_cat
: bytes
| category_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| doc
: bytes
| fk_doc_owner
: bytes
| owner_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| project
: bytes
| fk_proj_dept
: bytes
| dept_id
: bytes
| sdl_stress_big_1
: bytes
| id
< row
: bytes
| sdl_stress_big_1
: bytes
| tag_map
: bytes
| fk_tagmap_doc
: bytes
| doc_id
: bytes
| sdl_stress_big_1
: bytes
| id
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
| sdl_stress_big_1
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
< row
: bytes
| event_log
: bytes
| p2023
: null
: bytes
| RANGE
: null
: bytes
| year(`occurred_on`)
: null
: bytes
| 2024
: bytes
|
< row
: bytes
| event_log
: bytes
| p2024
: null
: bytes
| RANGE
: null
: bytes
| year(`occurred_on`)
: null
: bytes
| 2025
: bytes
|
< row
: bytes
| event_log
: bytes
| p_future
: null
: bytes
| RANGE
: null
: bytes
| year(`occurred_on`)
: null
: bytes
| MAXVALUE
: bytes
|
< row
: bytes
| metric
: bytes
| p0
: null
: bytes
| HASH
: null
: bytes
| `bucket`
: null
: null
: bytes
|
< row
: bytes
| metric
: bytes
| p1
: null
: bytes
| HASH
: null
: bytes
| `bucket`
: null
: null
: bytes
|
< row
: bytes
| metric
: bytes
| p2
: null
: bytes
| HASH
: null
: bytes
| `bucket`
: null
: null
: bytes
|
< row
: bytes
| metric
: bytes
| p3
: null
: bytes
| HASH
: null
: bytes
| `bucket`
: null
: null
: bytes
|
> query
| SHOW CREATE TABLE `sdl_stress_big_1`.`event_log`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| event_log
: bytes
| CREATE TABLE `event_log` (
|   `id` bigint NOT NULL AUTO_INCREMENT,
|   `occurred_on` date NOT NULL,
|   `kind` varchar(40) NOT NULL,
|   `payload` json DEFAULT NULL,
|   PRIMARY KEY (`id`,`occurred_on`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (year(`occurred_on`))
| (PARTITION p2023 VALUES LESS THAN (2024) ENGINE = InnoDB,
|  PARTITION p2024 VALUES LESS THAN (2025) ENGINE = InnoDB,
|  PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */
> query
| SHOW CREATE TABLE `sdl_stress_big_1`.`metric`
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| metric
: bytes
| CREATE TABLE `metric` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `bucket` int NOT NULL,
|   `name` varchar(60) NOT NULL,
|   `value` double NOT NULL DEFAULT '0',
|   PRIMARY KEY (`id`,`bucket`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY HASH (`bucket`)
| PARTITIONS 4 */
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
| sdl_stress_big_1
< result
< column "ROUTINE_NAME" "VARCHAR" "string"
< column "ROUTINE_TYPE" "ENUM" "string"
< column "SQL_MODE" "SET" "string"
< column "CHARACTER_SET_CLIENT" "VARCHAR" "string"
< column "COLLATION_CONNECTION" "VARCHAR" "string"
< column "DATABASE_COLLATION" "VARCHAR" "string"
< column "IFNULL(ROUTINE_COMMENT, '')" "TEXT" "string"
< row
: bytes
| dept_headcount
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| emp_annual_salary
: bytes
| FUNCTION
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| close_project
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
< row
: bytes
| give_raise
: bytes
| PROCEDURE
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
: bytes
|
> query
| SHOW CREATE FUNCTION `sdl_stress_big_1`.`dept_headcount`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| dept_headcount
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `dept_headcount`(d_id INT) RETURNS int
|     READS SQL DATA
| RETURN (SELECT COUNT(*) FROM employee WHERE dept_id = d_id)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE FUNCTION `sdl_stress_big_1`.`emp_annual_salary`
< result
< column "Function" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Function" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| emp_annual_salary
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` FUNCTION `emp_annual_salary`(monthly DECIMAL(10,2)) RETURNS decimal(12,2)
|     DETERMINISTIC
| RETURN monthly * 12
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sdl_stress_big_1`.`close_project`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| close_project
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `close_project`(IN proj INT)
| BEGIN
| 	DELETE FROM assignment WHERE project_id = proj;
| 	UPDATE project SET spent = budget WHERE id = proj;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
> query
| SHOW CREATE PROCEDURE `sdl_stress_big_1`.`give_raise`
< result
< column "Procedure" "VARCHAR" "string"
< column "sql_mode" "VARCHAR" "string"
< column "Create Procedure" "VARCHAR"
< column "character_set_client" "VARCHAR" "string"
< column "collation_connection" "VARCHAR" "string"
< column "Database Collation" "VARCHAR" "string"
< row
: bytes
| give_raise
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` PROCEDURE `give_raise`(IN emp INT, IN pct DECIMAL(5,2))
| BEGIN
| 	UPDATE employee SET salary = salary * (1 + pct / 100) WHERE id = emp;
| END
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
: bytes
| utf8mb4_0900_ai_ci
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
| sdl_stress_big_1
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
| assignment
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
| audit_trail
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
| category
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
| department
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
| 49152
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
| doc
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
| 32768
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
| employee
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
| 81920
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
| event_log
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 49152
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| metric
: bytes
| BASE TABLE
: bytes
| InnoDB
: bytes
| utf8mb4_0900_ai_ci
: bytes
| 0
: bytes
| 65536
: bytes
| 0
: bytes
| 0
: bytes
| partitioned
: bytes
| ''
: bytes
| utf8mb4
< row
: bytes
| project
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
| 32768
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
| tag_map
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
| v_active_employees
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
| v_dept_payroll
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
| v_project_load
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
> query
|
| 		SELECT t.NAME
| 		FROM information_schema.INNODB_TABLES t
| 		LEFT JOIN information_schema.INNODB_TABLESPACES ts ON t.SPACE = ts.SPACE
| 		WHERE t.NAME LIKE CONCAT(?, '/%') AND (ts.SPACE_TYPE IS NULL OR ts.SPACE_TYPE <> 'Single')
: string
| sdl_stress_big_1
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
| sdl_stress_big_1
< result
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_ai_ci
