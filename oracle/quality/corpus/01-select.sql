-- @name: simple select
-- @valid: true
-- @source: Oracle SQL Reference
SELECT 1 FROM dual

-- @name: select with alias
-- @valid: true
-- @source: Oracle SQL Reference
SELECT 1 AS one, 'hello' AS greeting FROM dual

-- @name: select from table
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id, first_name, last_name FROM employees

-- @name: select with where
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM employees WHERE department_id = 10

-- @name: select with and/or
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM employees WHERE salary > 5000 AND department_id IN (10, 20, 30)

-- @name: select with join
-- @valid: true
-- @source: Oracle SQL Reference
SELECT e.employee_id, d.department_name
FROM employees e
JOIN departments d ON e.department_id = d.department_id

-- @name: select with left outer join
-- @valid: true
-- @source: Oracle SQL Reference
SELECT e.employee_id, d.department_name
FROM employees e
LEFT OUTER JOIN departments d ON e.department_id = d.department_id

-- @name: select with oracle join syntax (+)
-- @valid: true
-- @source: Oracle SQL Reference
SELECT e.employee_id, d.department_name
FROM employees e, departments d
WHERE e.department_id = d.department_id(+)

-- @name: select with group by
-- @valid: true
-- @source: Oracle SQL Reference
SELECT department_id, COUNT(*) AS cnt
FROM employees
GROUP BY department_id

-- @name: select with group by having
-- @valid: true
-- @source: Oracle SQL Reference
SELECT department_id, AVG(salary) AS avg_sal
FROM employees
GROUP BY department_id
HAVING AVG(salary) > 5000

-- @name: select with order by
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id, salary FROM employees ORDER BY salary DESC, employee_id ASC

-- @name: select with order by nulls first
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id, commission_pct FROM employees ORDER BY commission_pct NULLS FIRST

-- @name: select distinct
-- @valid: true
-- @source: Oracle SQL Reference
SELECT DISTINCT department_id FROM employees

-- @name: select with subquery in where
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM employees WHERE department_id IN (
  SELECT department_id FROM departments WHERE location_id = 1700
)

-- @name: select with correlated subquery
-- @valid: true
-- @source: Oracle SQL Reference
SELECT e.employee_id, e.salary FROM employees e
WHERE e.salary > (SELECT AVG(salary) FROM employees WHERE department_id = e.department_id)

-- @name: select with exists
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM departments d
WHERE EXISTS (SELECT 1 FROM employees e WHERE e.department_id = d.department_id)

-- @name: select with case expression
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id,
  CASE WHEN salary > 10000 THEN 'HIGH'
       WHEN salary > 5000 THEN 'MID'
       ELSE 'LOW'
  END AS salary_band
FROM employees

-- @name: select with inline view
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM (
  SELECT employee_id, salary, ROWNUM rn FROM employees ORDER BY salary DESC
) WHERE rn <= 10

-- @name: select with CTE
-- @valid: true
-- @source: Oracle SQL Reference
WITH dept_avg AS (
  SELECT department_id, AVG(salary) AS avg_sal FROM employees GROUP BY department_id
)
SELECT e.employee_id, e.salary, da.avg_sal
FROM employees e
JOIN dept_avg da ON e.department_id = da.department_id

-- @name: select with recursive CTE
-- @valid: true
-- @source: Oracle SQL Reference
WITH RECURSIVE org_chart (employee_id, manager_id, lvl) AS (
  SELECT employee_id, manager_id, 1 FROM employees WHERE manager_id IS NULL
  UNION ALL
  SELECT e.employee_id, e.manager_id, oc.lvl + 1
  FROM employees e JOIN org_chart oc ON e.manager_id = oc.employee_id
)
SELECT * FROM org_chart

-- @name: select connect by
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id, manager_id, LEVEL
FROM employees
START WITH manager_id IS NULL
CONNECT BY PRIOR employee_id = manager_id

-- @name: select connect by with sys_connect_by_path
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id, SYS_CONNECT_BY_PATH(last_name, '/') AS path
FROM employees
START WITH manager_id IS NULL
CONNECT BY PRIOR employee_id = manager_id

-- @name: select with analytic function
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id, salary,
  ROW_NUMBER() OVER (PARTITION BY department_id ORDER BY salary DESC) AS rn
FROM employees

-- @name: select with analytic rank
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id, department_id, salary,
  RANK() OVER (ORDER BY salary DESC) AS salary_rank,
  DENSE_RANK() OVER (ORDER BY salary DESC) AS dense_rank
FROM employees

-- @name: select with window frame
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id, salary,
  SUM(salary) OVER (ORDER BY employee_id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) AS moving_sum
FROM employees

-- @name: select with pivot
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM (
  SELECT department_id, job_id, salary FROM employees
)
PIVOT (
  SUM(salary) FOR job_id IN ('IT_PROG' AS it, 'SA_REP' AS sales, 'FI_ACCOUNT' AS finance)
)

-- @name: select with unpivot
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM (
  SELECT employee_id, first_name, last_name FROM employees
)
UNPIVOT (
  name FOR name_type IN (first_name AS 'FIRST', last_name AS 'LAST')
)

-- @name: select for update
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM employees WHERE department_id = 10 FOR UPDATE

-- @name: select for update nowait
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM employees WHERE employee_id = 100 FOR UPDATE NOWAIT

-- @name: select with fetch first
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM employees ORDER BY salary DESC FETCH FIRST 5 ROWS ONLY

-- @name: select with offset fetch
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM employees ORDER BY employee_id OFFSET 10 ROWS FETCH NEXT 20 ROWS ONLY

-- @name: select with union
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id, first_name FROM employees WHERE department_id = 10
UNION
SELECT employee_id, first_name FROM employees WHERE department_id = 20

-- @name: select with union all
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id FROM employees WHERE salary > 10000
UNION ALL
SELECT employee_id FROM employees WHERE department_id = 10

-- @name: select with intersect
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id FROM employees WHERE salary > 10000
INTERSECT
SELECT employee_id FROM employees WHERE department_id = 10

-- @name: select with minus
-- @valid: true
-- @source: Oracle SQL Reference
SELECT employee_id FROM employees WHERE department_id = 10
MINUS
SELECT employee_id FROM employees WHERE salary < 5000

-- @name: select with hints
-- @valid: true
-- @source: Oracle SQL Reference
SELECT /*+ FULL(e) PARALLEL(e, 4) */ * FROM employees e

-- @name: select with flashback query
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM employees AS OF TIMESTAMP SYSTIMESTAMP - INTERVAL '1' HOUR

-- @name: select with sample
-- @valid: true
-- @source: Oracle SQL Reference
SELECT * FROM employees SAMPLE (10)

-- @name: select with lateral inline view
-- @valid: true
-- @source: Oracle SQL Reference
SELECT d.department_name, e.max_sal
FROM departments d,
LATERAL (SELECT MAX(salary) AS max_sal FROM employees WHERE department_id = d.department_id) e

-- @name: select with model clause
-- @valid: true
-- @source: Oracle SQL Reference
SELECT country, product, year, sales
FROM sales_view
MODEL
  PARTITION BY (country)
  DIMENSION BY (product, year)
  MEASURES (sales)
  RULES (
    sales['Bounce', 2024] = sales['Bounce', 2023] * 1.1
  )

-- @name: select with json_table
-- @valid: true
-- @source: Oracle SQL Reference
SELECT jt.*
FROM j_purchaseorder,
JSON_TABLE(po_document, '$.LineItems[*]'
  COLUMNS (
    itemno NUMBER PATH '$.ItemNumber',
    description VARCHAR2(100) PATH '$.Part.Description'
  )
) jt

-- @name: select with listagg
-- @valid: true
-- @source: Oracle SQL Reference
SELECT department_id,
  LISTAGG(last_name, ', ') WITHIN GROUP (ORDER BY last_name) AS employees
FROM employees
GROUP BY department_id

-- @name: select with xmlagg
-- @valid: true
-- @source: Oracle SQL Reference
SELECT department_id,
  XMLAGG(XMLELEMENT("name", last_name) ORDER BY last_name) AS xml_names
FROM employees
GROUP BY department_id

-- @name: missing FROM clause
-- @valid: false
SELECT * WHERE 1=1

-- @name: double WHERE
-- @valid: false
SELECT * FROM dual WHERE 1=1 WHERE 2=2

-- @name: unclosed parenthesis
-- @valid: false
SELECT (1 + 2 FROM dual
