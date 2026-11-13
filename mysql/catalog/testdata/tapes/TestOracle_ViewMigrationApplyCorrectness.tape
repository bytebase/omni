-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_ViewMigrationApplyCorrectness default#1

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-employees-current-dept-emp default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_employees_current_dept_emp
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_employees_current_dept_emp DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_employees_current_dept_emp
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no INT NOT NULL, dept_no CHAR(4) NOT NULL, from_date DATE NOT NULL, to_date DATE NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> exec
| CREATE TABLE departments (dept_no CHAR(4) NOT NULL, dept_name VARCHAR(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE VIEW dept_emp_latest_date AS SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date FROM dept_emp GROUP BY emp_no
< ok 0 0
> exec
| CREATE VIEW current_dept_emp AS SELECT l.emp_no, d.dept_no, l.from_date, l.to_date FROM dept_emp d JOIN dept_emp_latest_date l ON d.emp_no = l.emp_no AND d.from_date = l.from_date JOIN departments dp ON d.dept_no = dp.dept_no
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_employees_current_dept_emp.dept_emp
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
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_to_create_employees_current_dept_emp.departments
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
| SHOW CREATE VIEW vgen_to_create_employees_current_dept_emp.dept_emp_latest_date
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| SHOW CREATE VIEW vgen_to_create_employees_current_dept_emp.current_dept_emp
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| current_dept_emp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from ((`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`)))) join `departments` `dp` on((`d`.`dept_no` = `dp`.`dept_no`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-employees-current-dept-emp default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_employees_current_dept_emp
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_employees_current_dept_emp DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_employees_current_dept_emp
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no INT NOT NULL, dept_no CHAR(4) NOT NULL, from_date DATE NOT NULL, to_date DATE NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> exec
| CREATE TABLE departments (dept_no CHAR(4) NOT NULL, dept_name VARCHAR(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE VIEW dept_emp_latest_date AS SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date FROM dept_emp GROUP BY emp_no
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_employees_current_dept_emp.dept_emp
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
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_from_create_employees_current_dept_emp.departments
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
| SHOW CREATE VIEW vgen_from_create_employees_current_dept_emp.dept_emp_latest_date
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| dept_emp_latest_date
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `dept_emp_latest_date` AS select `dept_emp`.`emp_no` AS `emp_no`,max(`dept_emp`.`from_date`) AS `from_date`,max(`dept_emp`.`to_date`) AS `to_date` from `dept_emp` group by `dept_emp`.`emp_no`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-employees-current-dept-emp default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 6 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE dept_emp (emp_no INT NOT NULL, dept_no CHAR(4) NOT NULL, from_date DATE NOT NULL, to_date DATE NOT NULL, PRIMARY KEY (emp_no, dept_no))
< ok 0 0
> exec
| CREATE TABLE departments (dept_no CHAR(4) NOT NULL, dept_name VARCHAR(40) NOT NULL, PRIMARY KEY (dept_no))
< ok 0 0
> exec
| CREATE VIEW dept_emp_latest_date AS SELECT emp_no, MAX(from_date) AS from_date, MAX(to_date) AS to_date FROM dept_emp GROUP BY emp_no
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from ((`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`)))) join `departments` `dp` on((`d`.`dept_no` = `dp`.`dept_no`)))
< ok 0 0
> query
| SHOW CREATE TABLE vt.dept_emp
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
|   `to_date` date NOT NULL,
|   PRIMARY KEY (`emp_no`,`dept_no`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vt.departments
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
| SHOW CREATE VIEW vt.dept_emp_latest_date
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
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
| SHOW CREATE VIEW vt.current_dept_emp
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| current_dept_emp
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `current_dept_emp` AS select `l`.`emp_no` AS `emp_no`,`d`.`dept_no` AS `dept_no`,`l`.`from_date` AS `from_date`,`l`.`to_date` AS `to_date` from ((`dept_emp` `d` join `dept_emp_latest_date` `l` on(((`d`.`emp_no` = `l`.`emp_no`) and (`d`.`from_date` = `l`.`from_date`)))) join `departments` `dp` on((`d`.`dept_no` = `dp`.`dept_no`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-explicit-columns default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_explicit_columns
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_explicit_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_explicit_columns
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v (p, q) AS SELECT a, b FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_explicit_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_explicit_columns.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`p`,`q`) AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-explicit-columns default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_explicit_columns
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_explicit_columns DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_explicit_columns
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_explicit_columns.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-explicit-columns default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 1 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b VARCHAR(20))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` (`p`,`q`) AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` (`p`,`q`) AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-multi-join default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_multi_join
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_multi_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_multi_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, u.y, w.n FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_multi_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_to_create_multi_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_to_create_multi_join.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_multi_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`u`.`y` AS `y`,`w`.`n` AS `n` from ((`t` join `u` on((`t`.`c` = `u`.`x`))) join `w` on((`u`.`y` = `w`.`m`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-multi-join default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_multi_join
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_multi_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_multi_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_multi_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_from_create_multi_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_from_create_multi_join.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-multi-join default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 3 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `t`.`a` AS `a`,`u`.`y` AS `y`,`w`.`n` AS `n` from ((`t` join `u` on((`t`.`c` = `u`.`x`))) join `w` on((`u`.`y` = `w`.`m`)))
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vt.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vt.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`u`.`y` AS `y`,`w`.`n` AS `n` from ((`t` join `u` on((`t`.`c` = `u`.`x`))) join `w` on((`u`.`y` = `w`.`m`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-paren-subquery-compare default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_paren_subquery_compare
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_paren_subquery_compare DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_paren_subquery_compare
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT if(((SELECT count(0) FROM t) = 0),'NO','YES') AS x
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_paren_subquery_compare.t
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
> query
| SHOW CREATE VIEW vgen_to_create_paren_subquery_compare.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select if(((select count(0) from `t`) = 0),'NO','YES') AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-paren-subquery-compare default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_paren_subquery_compare
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_paren_subquery_compare DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_paren_subquery_compare
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_paren_subquery_compare.t
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

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-paren-subquery-compare default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 4 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select if(((select count(0) from `t`) = 0),'NO','YES') AS `x`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
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
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select if(((select count(0) from `t`) = 0),'NO','YES') AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-sakila-film-list default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_sakila_film_list
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_sakila_film_list DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_sakila_film_list
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT film.film_id AS FID, category.name AS category, actor.name AS actor FROM category LEFT JOIN film_category ON category.category_id = film_category.category_id LEFT JOIN film ON film_category.film_id = film.film_id JOIN film_actor ON film.film_id = film_actor.film_id JOIN actor ON film_actor.actor_id = actor.actor_id GROUP BY film.film_id
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int DEFAULT NULL,
|   `category_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int DEFAULT NULL,
|   `film_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_to_create_sakila_film_list.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_sakila_film_list.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `film`.`film_id` AS `FID`,`category`.`name` AS `category`,`actor`.`name` AS `actor` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`))) group by `film`.`film_id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-sakila-film-list default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_sakila_film_list
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_sakila_film_list DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_sakila_film_list
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int DEFAULT NULL,
|   `category_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int DEFAULT NULL,
|   `film_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_from_create_sakila_film_list.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-sakila-film-list default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 4 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE category (category_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_category (film_id INT, category_id INT)
< ok 0 0
> exec
| CREATE TABLE film (film_id INT, title VARCHAR(50))
< ok 0 0
> exec
| CREATE TABLE film_actor (actor_id INT, film_id INT)
< ok 0 0
> exec
| CREATE TABLE actor (actor_id INT, name VARCHAR(50))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `film`.`film_id` AS `FID`,`category`.`name` AS `category`,`actor`.`name` AS `actor` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`))) group by `film`.`film_id`
< ok 0 0
> query
| SHOW CREATE TABLE vt.category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| category
: bytes
| CREATE TABLE `category` (
|   `category_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vt.film_category
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_category
: bytes
| CREATE TABLE `film_category` (
|   `film_id` int DEFAULT NULL,
|   `category_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vt.film
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film
: bytes
| CREATE TABLE `film` (
|   `film_id` int DEFAULT NULL,
|   `title` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vt.film_actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| film_actor
: bytes
| CREATE TABLE `film_actor` (
|   `actor_id` int DEFAULT NULL,
|   `film_id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vt.actor
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| actor
: bytes
| CREATE TABLE `actor` (
|   `actor_id` int DEFAULT NULL,
|   `name` varchar(50) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `film`.`film_id` AS `FID`,`category`.`name` AS `category`,`actor`.`name` AS `actor` from ((((`category` left join `film_category` on((`category`.`category_id` = `film_category`.`category_id`))) left join `film` on((`film_category`.`film_id` = `film`.`film_id`))) join `film_actor` on((`film`.`film_id` = `film_actor`.`film_id`))) join `actor` on((`film_actor`.`actor_id` = `actor`.`actor_id`))) group by `film`.`film_id`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-simple default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_simple
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_simple DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_simple
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_simple.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_simple.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-simple default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_simple
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_simple DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_simple
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_simple.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-simple default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 0 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `t`.`a` AS `a` from `t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-tsdiff-view default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_tsdiff_view
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_tsdiff_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_tsdiff_view
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT id, TIMESTAMPDIFF(SECOND, created, NOW()) AS age_secs FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_tsdiff_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_tsdiff_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`id` AS `id`,timestampdiff(SECOND,`t`.`created`,now()) AS `age_secs` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-tsdiff-view default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_tsdiff_view
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_tsdiff_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_tsdiff_view
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_tsdiff_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-tsdiff-view default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 3 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `t`.`id` AS `id`,timestampdiff(SECOND,`t`.`created`,now()) AS `age_secs` from `t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int NOT NULL,
|   `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`id` AS `id`,timestampdiff(SECOND,`t`.`created`,now()) AS `age_secs` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-alias-name-clash default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_view_alias_name_clash
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_view_alias_name_clash DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_view_alias_name_clash
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| CREATE VIEW a AS SELECT id AS b FROM t
< ok 0 0
> exec
| CREATE VIEW b AS SELECT b FROM a
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_view_alias_name_clash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_alias_name_clash.a
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a` AS select `t`.`id` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_alias_name_clash.b
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b` AS select `a`.`b` AS `b` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-alias-name-clash default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_view_alias_name_clash
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_view_alias_name_clash DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_view_alias_name_clash
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_view_alias_name_clash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-alias-name-clash default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 4 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`a` AS select `t`.`id` AS `b` from `t`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`b` AS select `a`.`b` AS `b` from `a`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.a
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `a` AS select `t`.`id` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.b
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `b` AS select `a`.`b` AS `b` from `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-and-table default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_view_and_table
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_view_and_table DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_view_and_table
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_view_and_table.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_and_table.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-and-table default#2
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE `vt`.`t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `t`.`a` AS `a` from `t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-chain default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_view_chain
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_view_chain DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_view_chain
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v2 AS SELECT a FROM v1
< ok 0 0
> exec
| CREATE VIEW v3 AS SELECT a FROM v2
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_view_chain.t
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
> query
| SHOW CREATE VIEW vgen_to_create_view_chain.v1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_chain.v2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2` AS select `v1`.`a` AS `a` from `v1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_chain.v3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v3` AS select `v2`.`a` AS `a` from `v2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-chain default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_view_chain
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_view_chain DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_view_chain
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_view_chain.t
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

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-chain default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 3 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v1` AS select `t`.`a` AS `a` from `t`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v2` AS select `v1`.`a` AS `a` from `v1`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v3` AS select `v2`.`a` AS `a` from `v2`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
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
> query
| SHOW CREATE VIEW vt.v1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.v2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2` AS select `v1`.`a` AS `a` from `v1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.v3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v3` AS select `v2`.`a` AS `a` from `v2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-on-view default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_view_on_view
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT a, b FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM base
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_view_on_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_on_view.base
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `base` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_on_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `base`.`a` AS `a` from `base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-on-view default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_view_on_view
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_view_on_view.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-on-view default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`base` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `base`.`a` AS `a` from `base`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.base
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `base` AS select `t`.`a` AS `a`,`t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `base`.`a` AS `a` from `base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-on-view-table-primary default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_view_on_view_table_primary
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_view_on_view_table_primary DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_view_on_view_table_primary
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT id FROM t
< ok 0 0
> exec
| CREATE VIEW tv AS TABLE base
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_view_on_view_table_primary.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_on_view_table_primary.base
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `base` AS select `t`.`id` AS `id` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_to_create_view_on_view_table_primary.tv
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| tv
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `tv` AS select `base`.`id` AS `id` from `base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-on-view-table-primary default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_view_on_view_table_primary
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_view_on_view_table_primary DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_view_on_view_table_primary
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_view_on_view_table_primary.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-view-on-view-table-primary default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`base` AS select `t`.`id` AS `id` from `t`
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`tv` AS select `base`.`id` AS `id` from `base`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.base
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `base` AS select `t`.`id` AS `id` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vt.tv
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| tv
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `tv` AS select `base`.`id` AS `id` from `base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-with-options default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_create_with_options
< ok 0 0
> exec
| CREATE DATABASE vgen_to_create_with_options DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_create_with_options
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW v AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_create_with_options.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_create_with_options.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=MERGE DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0) WITH CASCADED CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-with-options default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_create_with_options
< ok 0 0
> exec
| CREATE DATABASE vgen_from_create_with_options DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_create_with_options
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_create_with_options.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/create-with-options default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=MERGE SQL SECURITY INVOKER VIEW `vt`.`v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0) WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=MERGE DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0) WITH CASCADED CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/drop-view default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_drop_view
< ok 0 0
> exec
| CREATE DATABASE vgen_to_drop_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_drop_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_drop_view.t
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

== TestOracle_ViewMigrationApplyCorrectness/8.0/drop-view default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_drop_view
< ok 0 0
> exec
| CREATE DATABASE vgen_from_drop_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_drop_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_drop_view.t
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
> query
| SHOW CREATE VIEW vgen_from_drop_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/drop-view default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| DROP VIEW IF EXISTS `vt`.`v`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
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
> query
| SHOW CREATE VIEW vt.v
< error
| 1146 42S02 Table 'vt.v' doesn't exist

== TestOracle_ViewMigrationApplyCorrectness/8.0/drop-view-on-view default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_drop_view_on_view
< ok 0 0
> exec
| CREATE DATABASE vgen_to_drop_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_drop_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_drop_view_on_view.t
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

== TestOracle_ViewMigrationApplyCorrectness/8.0/drop-view-on-view default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_drop_view_on_view
< ok 0 0
> exec
| CREATE DATABASE vgen_from_drop_view_on_view DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_drop_view_on_view
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM base
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_drop_view_on_view.t
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
> query
| SHOW CREATE VIEW vgen_from_drop_view_on_view.base
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| base
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `base` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SHOW CREATE VIEW vgen_from_drop_view_on_view.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `base`.`a` AS `a` from `base`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/drop-view-on-view default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 1 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW base AS SELECT a FROM t
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM base
< ok 0 0
> exec
| DROP VIEW IF EXISTS `vt`.`base`
< ok 0 0
> exec
| DROP VIEW IF EXISTS `vt`.`v`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
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
> query
| SHOW CREATE VIEW vt.base
< error
| 1146 42S02 Table 'vt.base' doesn't exist
> query
| SHOW CREATE VIEW vt.v
< error
| 1146 42S02 Table 'vt.v' doesn't exist

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-add-check-option default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_add_check_option
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_add_check_option DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_add_check_option
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0 WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_add_check_option.t
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
> query
| SHOW CREATE VIEW vgen_to_replace_add_check_option.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0) WITH CASCADED CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-add-check-option default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_add_check_option
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_add_check_option DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_add_check_option
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_add_check_option.t
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
> query
| SHOW CREATE VIEW vgen_from_replace_add_check_option.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-add-check-option default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t WHERE a > 0
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0) WITH CASCADED CHECK OPTION
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
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
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t` where (`t`.`a` > 0) WITH CASCADED CHECK OPTION
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-algorithm default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_algorithm
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_algorithm DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_algorithm
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_algorithm.t
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
> query
| SHOW CREATE VIEW vgen_to_replace_algorithm.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=MERGE DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-algorithm default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_algorithm
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_algorithm DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_algorithm
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_algorithm.t
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
> query
| SHOW CREATE VIEW vgen_from_replace_algorithm.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-algorithm default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE ALGORITHM=UNDEFINED VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=MERGE SQL SECURITY DEFINER VIEW `vt`.`v` AS select `t`.`a` AS `a` from `t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
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
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=MERGE DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-body default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_body
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_body
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_body.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_replace_body.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-body default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_body
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_body DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_body
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_body.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_from_replace_body.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-body default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `t`.`b` AS `b` from `t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-into-multi-join default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_into_multi_join
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_into_multi_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_into_multi_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.a, w.n FROM t JOIN u ON t.c = u.x JOIN w ON u.y = w.m
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_into_multi_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_to_replace_into_multi_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_to_replace_into_multi_join.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_replace_into_multi_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`w`.`n` AS `n` from ((`t` join `u` on((`t`.`c` = `u`.`x`))) join `w` on((`u`.`y` = `w`.`m`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-into-multi-join default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_into_multi_join
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_into_multi_join DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_into_multi_join
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_into_multi_join.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_from_replace_into_multi_join.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vgen_from_replace_into_multi_join.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_from_replace_into_multi_join.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-into-multi-join default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 4 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT, c INT)
< ok 0 0
> exec
| CREATE TABLE u (x INT, y INT)
< ok 0 0
> exec
| CREATE TABLE w (m INT, n INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select `t`.`a` AS `a`,`w`.`n` AS `n` from ((`t` join `u` on((`t`.`c` = `u`.`x`))) join `w` on((`u`.`y` = `w`.`m`)))
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vt.u
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| u
: bytes
| CREATE TABLE `u` (
|   `x` int DEFAULT NULL,
|   `y` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE vt.w
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| w
: bytes
| CREATE TABLE `w` (
|   `m` int DEFAULT NULL,
|   `n` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a`,`w`.`n` AS `n` from ((`t` join `u` on((`t`.`c` = `u`.`x`))) join `w` on((`u`.`y` = `w`.`m`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-into-temporal-medley default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_into_temporal_medley
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_into_temporal_medley DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_into_temporal_medley
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT TIMESTAMPDIFF(DAY, a, NOW()) AS c1, TIMESTAMPADD(HOUR, 1, a) AS c2, EXTRACT(DAY_HOUR FROM a) AS c3, GET_FORMAT(DATE, 'USA') AS c4, a + INTERVAL '1:1' DAY_HOUR AS c5, TRIM('x' FROM s) AS c6 FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_into_temporal_medley.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_to_replace_into_temporal_medley.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select timestampdiff(DAY,`t`.`a`,now()) AS `c1`,(`t`.`a` + interval 1 hour) AS `c2`,extract(day_hour from `t`.`a`) AS `c3`,get_format(DATE, 'USA') AS `c4`,(`t`.`a` + interval '1:1' day_hour) AS `c5`,trim('x' from `t`.`s`) AS `c6` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-into-temporal-medley default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_into_temporal_medley
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_into_temporal_medley DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_into_temporal_medley
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_into_temporal_medley.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vgen_from_replace_into_temporal_medley.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-into-temporal-medley default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a DATETIME, s VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `vt`.`v` AS select timestampdiff(DAY,`t`.`a`,now()) AS `c1`,(`t`.`a` + interval 1 hour) AS `c2`,extract(day_hour from `t`.`a`) AS `c3`,get_format(DATE, 'USA') AS `c4`,(`t`.`a` + interval '1:1' day_hour) AS `c5`,trim('x' from `t`.`s`) AS `c6` from `t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime DEFAULT NULL,
|   `s` varchar(20) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select timestampdiff(DAY,`t`.`a`,now()) AS `c1`,(`t`.`a` + interval 1 hour) AS `c2`,extract(day_hour from `t`.`a`) AS `c3`,get_format(DATE, 'USA') AS `c4`,(`t`.`a` + interval '1:1' day_hour) AS `c5`,trim('x' from `t`.`s`) AS `c6` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-sql-security default#1
> exec
| DROP DATABASE IF EXISTS vgen_to_replace_sql_security
< ok 0 0
> exec
| CREATE DATABASE vgen_to_replace_sql_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_to_replace_sql_security
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_to_replace_sql_security.t
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
> query
| SHOW CREATE VIEW vgen_to_replace_sql_security.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-sql-security default#2
> exec
| DROP DATABASE IF EXISTS vgen_from_replace_sql_security
< ok 0 0
> exec
| CREATE DATABASE vgen_from_replace_sql_security DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vgen_from_replace_sql_security
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW v AS SELECT a FROM t
< ok 0 0
> query
| SHOW CREATE TABLE vgen_from_replace_sql_security.t
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
> query
| SHOW CREATE VIEW vgen_from_replace_sql_security.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestOracle_ViewMigrationApplyCorrectness/8.0/replace-sql-security default#3
> exec
| DROP DATABASE IF EXISTS vt
< ok 2 0
> exec
| CREATE DATABASE vt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE vt
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW v AS SELECT a FROM t
< ok 0 0
> exec
| CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY INVOKER VIEW `vt`.`v` AS select `t`.`a` AS `a` from `t`
< ok 0 0
> query
| SHOW CREATE TABLE vt.t
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
> query
| SHOW CREATE VIEW vt.v
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY INVOKER VIEW `v` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
