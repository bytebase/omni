-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_MigrationReproducesStoredForm default#1
> exec
| DROP DATABASE IF EXISTS storedform_f_column_first
< ok 0 0
> exec
| CREATE DATABASE storedform_f_column_first DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_column_first
< ok 0 0
> exec
| CREATE TABLE category (category_id int, name varchar(25))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_column_first' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| category
> query
| SHOW CREATE TABLE storedform_f_column_first.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int DEFAULT NULL,
|   `name` varchar(25) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| category
> query
| SHOW CREATE TABLE storedform_t_column_first.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `is_featured` tinyint(1) NOT NULL DEFAULT '0',
|   `category_id` int DEFAULT NULL,
|   `name` varchar(25) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_f_generated_columns_in_target_order
< ok 0 0
> exec
| CREATE DATABASE storedform_f_generated_columns_in_target_order DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_generated_columns_in_target_order
< ok 0 0
> exec
| CREATE TABLE p (id int, price decimal(10,2), name varchar(100))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_generated_columns_in_target_order' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| p
> query
| SHOW CREATE TABLE storedform_f_generated_columns_in_target_order.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `name` varchar(100) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| p
> query
| SHOW CREATE TABLE storedform_t_generated_columns_in_target_order.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * 1.1)) STORED,
|   `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_f_columns_mid_table
< ok 0 0
> exec
| CREATE DATABASE storedform_f_columns_mid_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_columns_mid_table
< ok 0 0
> exec
| CREATE TABLE t (a int, d int)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_columns_mid_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_f_columns_mid_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `d` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_t_columns_mid_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`b` + 1)) VIRTUAL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_f_column_readded_in_place
< ok 0 0
> exec
| CREATE DATABASE storedform_f_column_readded_in_place DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_column_readded_in_place
< ok 0 0
> exec
| CREATE TABLE t (id int, g int GENERATED ALWAYS AS (z+1) VIRTUAL, z smallint)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_column_readded_in_place' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_f_column_readded_in_place.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   `z` smallint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| CREATE TABLE t (id int, g int GENERATED ALWAYS AS (z+1) STORED, z bigint)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_column_readded_in_place' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_t_column_readded_in_place.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) STORED,
|   `z` bigint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_f_modify_reads_added_column
< ok 0 0
> exec
| CREATE DATABASE storedform_f_modify_reads_added_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_modify_reads_added_column
< ok 0 0
> exec
| CREATE TABLE t (a int, g int GENERATED ALWAYS AS (a+1) VIRTUAL)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_modify_reads_added_column' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_f_modify_reads_added_column.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`a` + 1)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_t_modify_reads_added_column
< ok 0 0
> exec
| CREATE DATABASE storedform_t_modify_reads_added_column DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_modify_reads_added_column
< ok 0 0
> exec
| CREATE TABLE t (a int, g int GENERATED ALWAYS AS (z+1) VIRTUAL, z int)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_modify_reads_added_column' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_t_modify_reads_added_column.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   `z` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_f_grouped_auto_increment_columns
< ok 0 0
> exec
| CREATE DATABASE storedform_f_grouped_auto_increment_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_grouped_auto_increment_columns
< ok 0 0
> exec
| CREATE TABLE t (a int)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_grouped_auto_increment_columns' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_f_grouped_auto_increment_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_t_grouped_auto_increment_columns
< ok 0 0
> exec
| CREATE DATABASE storedform_t_grouped_auto_increment_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_grouped_auto_increment_columns
< ok 0 0
> exec
| CREATE TABLE t (id int NOT NULL AUTO_INCREMENT, x int, a int, KEY k (id, x))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_grouped_auto_increment_columns' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| t
> query
| SHOW CREATE TABLE storedform_t_grouped_auto_increment_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `x` int DEFAULT NULL,
|   `a` int DEFAULT NULL,
|   KEY `k` (`id`,`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_f_fk_backing_index_new_table
< ok 0 0
> exec
| CREATE DATABASE storedform_f_fk_backing_index_new_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_fk_backing_index_new_table
< ok 0 0
> exec
| CREATE TABLE departments (dept_no char(4) NOT NULL, dept_name varchar(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE TABLE employees (emp_no int NOT NULL, first_name varchar(14) NOT NULL, PRIMARY KEY (emp_no))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_fk_backing_index_new_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| departments
< row
: bytes
| employees
> query
| SHOW CREATE TABLE storedform_f_fk_backing_index_new_table.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_f_fk_backing_index_new_table.employees
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| employees
: bytes
| CREATE TABLE `employees` (
|   `emp_no` int NOT NULL,
|   `first_name` varchar(14) NOT NULL,
|   PRIMARY KEY (`emp_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_t_fk_backing_index_new_table
< ok 0 0
> exec
| CREATE DATABASE storedform_t_fk_backing_index_new_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_fk_backing_index_new_table
< ok 0 0
> exec
| CREATE TABLE departments (dept_no char(4) NOT NULL, dept_name varchar(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE TABLE employees (emp_no int NOT NULL, first_name varchar(14) NOT NULL, PRIMARY KEY (emp_no))
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no int NOT NULL, dept_no char(4) NOT NULL, from_date date NOT NULL, PRIMARY KEY (emp_no, dept_no), KEY dept_no (dept_no), CONSTRAINT dept_emp_ibfk_1 FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE, CONSTRAINT dept_emp_ibfk_2 FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_fk_backing_index_new_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| departments
< row
: bytes
| dept_emp
< row
: bytes
| employees
> query
| SHOW CREATE TABLE storedform_t_fk_backing_index_new_table.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_backing_index_new_table.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`),
|   KEY `dept_no` (`dept_no`),
|   CONSTRAINT `dept_emp_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE,
|   CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_backing_index_new_table.employees
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| employees
: bytes
| CREATE TABLE `employees` (
|   `emp_no` int NOT NULL,
|   `first_name` varchar(14) NOT NULL,
|   PRIMARY KEY (`emp_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_f_fk_unnamed_new_table
< ok 0 0
> exec
| CREATE DATABASE storedform_f_fk_unnamed_new_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_fk_unnamed_new_table
< ok 0 0
> exec
| CREATE TABLE departments (dept_no char(4) NOT NULL, dept_name varchar(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_fk_unnamed_new_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| departments
> query
| SHOW CREATE TABLE storedform_f_fk_unnamed_new_table.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_t_fk_unnamed_new_table
< ok 0 0
> exec
| CREATE DATABASE storedform_t_fk_unnamed_new_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_fk_unnamed_new_table
< ok 0 0
> exec
| CREATE TABLE departments (dept_no char(4) NOT NULL, dept_name varchar(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE TABLE c (id int PRIMARY KEY, dept_no char(4), FOREIGN KEY (dept_no) REFERENCES departments (dept_no))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_fk_unnamed_new_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| c
< row
: bytes
| departments
> query
| SHOW CREATE TABLE storedform_t_fk_unnamed_new_table.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `dept_no` char(4) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `dept_no` (`dept_no`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_unnamed_new_table.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_f_fk_backing_index_existing_table
< ok 0 0
> exec
| CREATE DATABASE storedform_f_fk_backing_index_existing_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_fk_backing_index_existing_table
< ok 0 0
> exec
| CREATE TABLE departments (dept_no char(4) NOT NULL, dept_name varchar(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE TABLE employees (emp_no int NOT NULL, first_name varchar(14) NOT NULL, PRIMARY KEY (emp_no))
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no int NOT NULL, dept_no char(4) NOT NULL, from_date date NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_fk_backing_index_existing_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| departments
< row
: bytes
| dept_emp
< row
: bytes
| employees
> query
| SHOW CREATE TABLE storedform_f_fk_backing_index_existing_table.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_f_fk_backing_index_existing_table.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_f_fk_backing_index_existing_table.employees
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| employees
: bytes
| CREATE TABLE `employees` (
|   `emp_no` int NOT NULL,
|   `first_name` varchar(14) NOT NULL,
|   PRIMARY KEY (`emp_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_t_fk_backing_index_existing_table
< ok 0 0
> exec
| CREATE DATABASE storedform_t_fk_backing_index_existing_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_fk_backing_index_existing_table
< ok 0 0
> exec
| CREATE TABLE departments (dept_no char(4) NOT NULL, dept_name varchar(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE TABLE employees (emp_no int NOT NULL, first_name varchar(14) NOT NULL, PRIMARY KEY (emp_no))
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no int NOT NULL, dept_no char(4) NOT NULL, from_date date NOT NULL, PRIMARY KEY (emp_no, dept_no), KEY dept_no (dept_no), CONSTRAINT dept_emp_ibfk_1 FOREIGN KEY (emp_no) REFERENCES employees (emp_no) ON DELETE CASCADE, CONSTRAINT dept_emp_ibfk_2 FOREIGN KEY (dept_no) REFERENCES departments (dept_no) ON DELETE CASCADE)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_fk_backing_index_existing_table' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| departments
< row
: bytes
| dept_emp
< row
: bytes
| employees
> query
| SHOW CREATE TABLE storedform_t_fk_backing_index_existing_table.departments
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| departments
: bytes
| CREATE TABLE `departments` (
|   `dept_no` char(4) NOT NULL,
|   `dept_name` varchar(40) NOT NULL,
|   PRIMARY KEY (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_backing_index_existing_table.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`),
|   KEY `dept_no` (`dept_no`),
|   CONSTRAINT `dept_emp_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE,
|   CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_backing_index_existing_table.employees
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| employees
: bytes
| CREATE TABLE `employees` (
|   `emp_no` int NOT NULL,
|   `first_name` varchar(14) NOT NULL,
|   PRIMARY KEY (`emp_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_f_fk_restrict
< ok 0 0
> exec
| CREATE DATABASE storedform_f_fk_restrict DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_f_fk_restrict
< ok 0 0
> exec
| CREATE TABLE city (city_id smallint unsigned NOT NULL AUTO_INCREMENT, city varchar(50) NOT NULL, PRIMARY KEY (city_id))
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_f_fk_restrict' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| city
> query
| SHOW CREATE TABLE storedform_f_fk_restrict.city
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| city
: bytes
| CREATE TABLE `city` (
|   `city_id` smallint unsigned NOT NULL AUTO_INCREMENT,
|   `city` varchar(50) NOT NULL,
|   PRIMARY KEY (`city_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS storedform_t_fk_restrict
< ok 0 0
> exec
| CREATE DATABASE storedform_t_fk_restrict DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform_t_fk_restrict
< ok 0 0
> exec
| CREATE TABLE city (city_id smallint unsigned NOT NULL AUTO_INCREMENT, city varchar(50) NOT NULL, PRIMARY KEY (city_id))
< ok 0 0
> exec
| CREATE TABLE address (address_id smallint unsigned NOT NULL AUTO_INCREMENT, city_id smallint unsigned NOT NULL, PRIMARY KEY (address_id), KEY idx_fk_city_id (city_id), CONSTRAINT fk_address_city FOREIGN KEY (city_id) REFERENCES city (city_id) ON DELETE RESTRICT ON UPDATE CASCADE, CONSTRAINT fk_address_city_2 FOREIGN KEY (city_id) REFERENCES city (city_id) ON UPDATE RESTRICT)
< ok 0 0
> query
| SELECT table_name FROM information_schema.tables WHERE table_schema = 'storedform_t_fk_restrict' AND table_type='BASE TABLE' ORDER BY table_name
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| address
< row
: bytes
| city
> query
| SHOW CREATE TABLE storedform_t_fk_restrict.address
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| address
: bytes
| CREATE TABLE `address` (
|   `address_id` smallint unsigned NOT NULL AUTO_INCREMENT,
|   `city_id` smallint unsigned NOT NULL,
|   PRIMARY KEY (`address_id`),
|   KEY `idx_fk_city_id` (`city_id`),
|   CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT `fk_address_city_2` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_restrict.city
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| city
: bytes
| CREATE TABLE `city` (
|   `city_id` smallint unsigned NOT NULL AUTO_INCREMENT,
|   `city` varchar(50) NOT NULL,
|   PRIMARY KEY (`city_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/column-first default#1
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
> exec
| ALTER TABLE `storedform`.`category` ADD COLUMN `is_featured` tinyint(1) NOT NULL DEFAULT '0' FIRST
< ok 0 0
> query
| SHOW CREATE TABLE storedform.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `is_featured` tinyint(1) NOT NULL DEFAULT '0',
|   `category_id` int DEFAULT NULL,
|   `name` varchar(25) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_column_first.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `is_featured` tinyint(1) NOT NULL DEFAULT '0',
|   `category_id` int DEFAULT NULL,
|   `name` varchar(25) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/column-readded-in-place default#1
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
| CREATE TABLE t (id int, g int GENERATED ALWAYS AS (z+1) VIRTUAL, z smallint)
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` DROP COLUMN `g`
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` MODIFY COLUMN `z` bigint DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS (((`z` + 1))) STORED AFTER `id`
< ok 0 0
> query
| SHOW CREATE TABLE storedform.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) STORED,
|   `z` bigint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_column_readded_in_place.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) STORED,
|   `z` bigint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/columns-mid-table default#1
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
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `b` int DEFAULT NULL AFTER `a`
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `c` int DEFAULT NULL AFTER `b`
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `g` int GENERATED ALWAYS AS (((`b` + 1))) VIRTUAL AFTER `b`
< ok 0 0
> query
| SHOW CREATE TABLE storedform.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`b` + 1)) VIRTUAL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_columns_mid_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`b` + 1)) VIRTUAL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/fk-backing-index-existing-table default#1
> exec
| DROP DATABASE IF EXISTS storedform
< ok 2 0
> exec
| CREATE DATABASE storedform DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform
< ok 0 0
> exec
| CREATE TABLE departments (dept_no char(4) NOT NULL, dept_name varchar(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE TABLE employees (emp_no int NOT NULL, first_name varchar(14) NOT NULL, PRIMARY KEY (emp_no))
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no int NOT NULL, dept_no char(4) NOT NULL, from_date date NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> exec
| ALTER TABLE `storedform`.`dept_emp` ADD KEY `dept_no` (`dept_no`)
< ok 0 0
> exec
| ALTER TABLE `storedform`.`dept_emp` ADD CONSTRAINT `dept_emp_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE
< ok 0 0
> exec
| ALTER TABLE `storedform`.`dept_emp` ADD CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
< ok 0 0
> query
| SHOW CREATE TABLE storedform.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`),
|   KEY `dept_no` (`dept_no`),
|   CONSTRAINT `dept_emp_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE,
|   CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_backing_index_existing_table.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`),
|   KEY `dept_no` (`dept_no`),
|   CONSTRAINT `dept_emp_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE,
|   CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/fk-backing-index-new-table default#1
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
| CREATE TABLE departments (dept_no char(4) NOT NULL, dept_name varchar(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE TABLE employees (emp_no int NOT NULL, first_name varchar(14) NOT NULL, PRIMARY KEY (emp_no))
< ok 0 0
> exec
| CREATE TABLE `storedform`.`dept_emp` (
|   `emp_no` int NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `storedform`.`dept_emp` ADD KEY `dept_no` (`dept_no`)
< ok 0 0
> exec
| ALTER TABLE `storedform`.`dept_emp` ADD CONSTRAINT `dept_emp_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE
< ok 0 0
> exec
| ALTER TABLE `storedform`.`dept_emp` ADD CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
< ok 0 0
> query
| SHOW CREATE TABLE storedform.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`),
|   KEY `dept_no` (`dept_no`),
|   CONSTRAINT `dept_emp_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE,
|   CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_backing_index_new_table.dept_emp
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dept_emp
: bytes
| CREATE TABLE `dept_emp` (
|   `emp_no` int NOT NULL,
|   `dept_no` char(4) NOT NULL,
|   `from_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`),
|   KEY `dept_no` (`dept_no`),
|   CONSTRAINT `dept_emp_ibfk_1` FOREIGN KEY (`emp_no`) REFERENCES `employees` (`emp_no`) ON DELETE CASCADE,
|   CONSTRAINT `dept_emp_ibfk_2` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/fk-restrict default#1
> exec
| DROP DATABASE IF EXISTS storedform
< ok 3 0
> exec
| CREATE DATABASE storedform DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform
< ok 0 0
> exec
| CREATE TABLE city (city_id smallint unsigned NOT NULL AUTO_INCREMENT, city varchar(50) NOT NULL, PRIMARY KEY (city_id))
< ok 0 0
> exec
| CREATE TABLE `storedform`.`address` (
|   `address_id` smallint unsigned NOT NULL AUTO_INCREMENT,
|   `city_id` smallint unsigned NOT NULL,
|   PRIMARY KEY (`address_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `storedform`.`address` ADD KEY `idx_fk_city_id` (`city_id`)
< ok 0 0
> exec
| ALTER TABLE `storedform`.`address` ADD CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON DELETE RESTRICT ON UPDATE CASCADE
< ok 0 0
> exec
| ALTER TABLE `storedform`.`address` ADD CONSTRAINT `fk_address_city_2` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON UPDATE RESTRICT
< ok 0 0
> query
| SHOW CREATE TABLE storedform.address
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| address
: bytes
| CREATE TABLE `address` (
|   `address_id` smallint unsigned NOT NULL AUTO_INCREMENT,
|   `city_id` smallint unsigned NOT NULL,
|   PRIMARY KEY (`address_id`),
|   KEY `idx_fk_city_id` (`city_id`),
|   CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT `fk_address_city_2` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_restrict.address
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| address
: bytes
| CREATE TABLE `address` (
|   `address_id` smallint unsigned NOT NULL AUTO_INCREMENT,
|   `city_id` smallint unsigned NOT NULL,
|   PRIMARY KEY (`address_id`),
|   KEY `idx_fk_city_id` (`city_id`),
|   CONSTRAINT `fk_address_city` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON DELETE RESTRICT ON UPDATE CASCADE,
|   CONSTRAINT `fk_address_city_2` FOREIGN KEY (`city_id`) REFERENCES `city` (`city_id`) ON UPDATE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/fk-unnamed-new-table default#1
> exec
| DROP DATABASE IF EXISTS storedform
< ok 3 0
> exec
| CREATE DATABASE storedform DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE storedform
< ok 0 0
> exec
| CREATE TABLE departments (dept_no char(4) NOT NULL, dept_name varchar(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE TABLE `storedform`.`c` (
|   `id` int NOT NULL,
|   `dept_no` char(4) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| ALTER TABLE `storedform`.`c` ADD KEY `dept_no` (`dept_no`)
< ok 0 0
> exec
| ALTER TABLE `storedform`.`c` ADD CONSTRAINT `c_ibfk_1` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`)
< ok 0 0
> query
| SHOW CREATE TABLE storedform.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `dept_no` char(4) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `dept_no` (`dept_no`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_fk_unnamed_new_table.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int NOT NULL,
|   `dept_no` char(4) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `dept_no` (`dept_no`),
|   CONSTRAINT `c_ibfk_1` FOREIGN KEY (`dept_no`) REFERENCES `departments` (`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/generated-columns-in-target-order default#1
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
> exec
| ALTER TABLE `storedform`.`p` ADD COLUMN `price_with_tax` decimal(10,2) GENERATED ALWAYS AS (((`price` * 1.1))) STORED
< ok 0 0
> exec
| ALTER TABLE `storedform`.`p` ADD COLUMN `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL
< ok 0 0
> query
| SHOW CREATE TABLE storedform.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * 1.1)) STORED,
|   `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_generated_columns_in_target_order.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int DEFAULT NULL,
|   `price` decimal(10,2) DEFAULT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   `price_with_tax` decimal(10,2) GENERATED ALWAYS AS ((`price` * 1.1)) STORED,
|   `name_upper` varchar(100) GENERATED ALWAYS AS (upper(`name`)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/grouped-auto-increment-columns default#1
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
| CREATE TABLE t (a int)
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `id` int NOT NULL AUTO_INCREMENT FIRST, ADD COLUMN `x` int DEFAULT NULL AFTER `id`, ADD KEY `k` (`id`,`x`)
< ok 0 0
> query
| SHOW CREATE TABLE storedform.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `x` int DEFAULT NULL,
|   `a` int DEFAULT NULL,
|   KEY `k` (`id`,`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_grouped_auto_increment_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `x` int DEFAULT NULL,
|   `a` int DEFAULT NULL,
|   KEY `k` (`id`,`x`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_MigrationReproducesStoredForm/modify-reads-added-column default#1
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
> exec
| ALTER TABLE `storedform`.`t` ADD COLUMN `z` int DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE `storedform`.`t` MODIFY COLUMN `g` int GENERATED ALWAYS AS (((`z` + 1))) VIRTUAL
< ok 0 0
> query
| SHOW CREATE TABLE storedform.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   `z` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE storedform_t_modify_reads_added_column.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `g` int GENERATED ALWAYS AS ((`z` + 1)) VIRTUAL,
|   `z` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
