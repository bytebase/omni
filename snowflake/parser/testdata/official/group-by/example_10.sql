SELECT SUM(salary), ANY_VALUE(employment_state) AS state
  FROM employees
  GROUP BY state;
