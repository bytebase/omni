CREATE RECURSIVE VIEW employee_hierarchy_02 (title, employee_ID, manager_ID, "MGR_EMP_ID (SHOULD BE SAME)", "MGR TITLE") AS (
      SELECT title, employee_ID, manager_ID, NULL AS "MGR_EMP_ID (SHOULD BE SAME)", 'President' AS "MGR TITLE"
        FROM employees
        WHERE title = 'President'
      UNION ALL
      SELECT employees.title, employees.employee_ID, employees.manager_ID,
             employee_hierarchy_02.employee_id AS "MGR_EMP_ID (SHOULD BE SAME)",
             employee_hierarchy_02.title AS "MGR TITLE"
        FROM employees INNER JOIN employee_hierarchy_02
        WHERE employee_hierarchy_02.employee_ID = employees.manager_ID
);
