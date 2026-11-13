-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDeparseCorrelatedSubqueryMariaDB default#1
> query
| SELECT DATABASE()
< result
< column "DATABASE()" "VARCHAR"
< row
: bytes
| omni_testdeparsecorrelatedsubquerymariadb
> exec
| CREATE TABLE employees (id INT PRIMARY KEY, name VARCHAR(100), salary DECIMAL(10,2), department_id INT)
< ok 0 0
> exec
| CREATE TABLE departments (id INT PRIMARY KEY, name VARCHAR(100))
< ok 0 0
> exec
| INSERT INTO departments VALUES (1, 'eng'), (2, 'ops'), (3, 'empty')
< ok 3 0
> exec
| INSERT INTO employees VALUES (1, 'ann', 100, 1), (2, 'bob', 200, 1), (3, 'cy', 150, 2), (4, 'eng', 50, NULL)
< ok 4 0
> query
| SELECT e.id, e.salary IN (SELECT d.id FROM departments d WHERE d.id = e.department_id) AS x FROM employees e ORDER BY 1
< result
< column "id" "INT"
< column "x" "INT"
< row
: int64 1
: int64 0
< row
: int64 2
: int64 0
< row
: int64 3
: int64 0
< row
: int64 4
: int64 0
> query
| select `e`.`id`, `e`.`salary` in (select `d`.`id` from `omni_testdeparsecorrelatedsubquerymariadb`.`departments` as `d` where (`d`.`id` = `e`.`department_id`)) as `x` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `e` order by `e`.`id`
< result
< column "id" "INT"
< column "x" "INT"
< row
: int64 1
: int64 0
< row
: int64 2
: int64 0
< row
: int64 3
: int64 0
< row
: int64 4
: int64 0
> query
| SELECT e.id FROM employees e WHERE EXISTS (SELECT 1 FROM departments d WHERE d.id = e.department_id) ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
> query
| select `e`.`id` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `e` where exists (select 1 as `1` from `omni_testdeparsecorrelatedsubquerymariadb`.`departments` as `d` where (`d`.`id` = `e`.`department_id`)) order by `e`.`id`
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
> query
| SELECT e.id, (SELECT d.name FROM departments d WHERE d.id = e.department_id) AS dept FROM employees e ORDER BY 1
< result
< column "id" "INT"
< column "dept" "VARCHAR"
< row
: int64 1
: bytes
| eng
< row
: int64 2
: bytes
| eng
< row
: int64 3
: bytes
| ops
< row
: int64 4
: null
> query
| select `e`.`id`, (select `d`.`name` from `omni_testdeparsecorrelatedsubquerymariadb`.`departments` as `d` where (`d`.`id` = `e`.`department_id`)) as `dept` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `e` order by `e`.`id`
< result
< column "id" "INT"
< column "dept" "VARCHAR"
< row
: int64 1
: bytes
| eng
< row
: int64 2
: bytes
| eng
< row
: int64 3
: bytes
| ops
< row
: int64 4
: null
> query
| SELECT e.id FROM employees e WHERE e.salary > ANY (SELECT o.salary FROM employees o WHERE o.department_id = e.department_id) ORDER BY 1
< result
< column "id" "INT"
< row
: int64 2
> query
| select `e`.`id` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `e` where (`e`.`salary` > any (select `o`.`salary` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `o` where (`o`.`department_id` = `e`.`department_id`))) order by `e`.`id`
< result
< column "id" "INT"
< row
: int64 2
> query
| SELECT e.id FROM employees e WHERE e.salary >= ALL (SELECT o.salary FROM employees o WHERE o.department_id = e.department_id) ORDER BY 1
< result
< column "id" "INT"
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| select `e`.`id` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `e` where (`e`.`salary` >= all (select `o`.`salary` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `o` where (`o`.`department_id` = `e`.`department_id`))) order by `e`.`id`
< result
< column "id" "INT"
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| SELECT e.id FROM employees e WHERE EXISTS (SELECT 1 FROM departments d WHERE d.id = e.department_id AND EXISTS (SELECT 1 FROM employees o WHERE o.department_id = d.id AND o.salary > e.salary)) ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
> query
| select `e`.`id` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `e` where exists (select 1 as `1` from `omni_testdeparsecorrelatedsubquerymariadb`.`departments` as `d` where ((`d`.`id` = `e`.`department_id`) and exists (select 1 as `1` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `o` where ((`o`.`department_id` = `d`.`id`) and (`o`.`salary` > `e`.`salary`))))) order by `e`.`id`
< result
< column "id" "INT"
< row
: int64 1
> query
| SELECT (SELECT e.name FROM departments d WHERE d.id = e.department_id) AS n FROM employees e ORDER BY 1
< result
< column "n" "VARCHAR"
< row
: null
< row
: bytes
| ann
< row
: bytes
| bob
< row
: bytes
| cy
> query
| select (select `e`.`name` from `omni_testdeparsecorrelatedsubquerymariadb`.`departments` as `d` where (`d`.`id` = `e`.`department_id`)) as `n` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `e` order by (select `e`.`name` from `omni_testdeparsecorrelatedsubquerymariadb`.`departments` as `d` where (`d`.`id` = `e`.`department_id`))
< result
< column "n" "VARCHAR"
< row
: null
< row
: bytes
| ann
< row
: bytes
| bob
< row
: bytes
| cy
> query
| SELECT e.id FROM employees e WHERE e.department_id IN (SELECT d.id FROM departments d WHERE d.name = e.name UNION ALL SELECT o.department_id FROM employees o WHERE o.id = e.id) ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
> query
| select `e`.`id` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `e` where `e`.`department_id` in (select `d`.`id` from `omni_testdeparsecorrelatedsubquerymariadb`.`departments` as `d` where (`d`.`name` = `e`.`name`) union all select `o`.`department_id` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `o` where (`o`.`id` = `e`.`id`)) order by `e`.`id`
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
> query
| SELECT e.id FROM employees e WHERE EXISTS (WITH c AS (SELECT d.id FROM departments d) SELECT 1 FROM c WHERE c.id = e.department_id) ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
> query
| select `e`.`id` from `omni_testdeparsecorrelatedsubquerymariadb`.`employees` as `e` where exists (with `c`(`id`) as (select `d`.`id` from `omni_testdeparsecorrelatedsubquerymariadb`.`departments` as `d`) select 1 as `1` from `c` where (`c`.`id` = `e`.`department_id`)) order by `e`.`id`
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
