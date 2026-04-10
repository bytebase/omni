SELECT SUM(salary), ANY_VALUE(employment_state)
  FROM employees
  GROUP BY employment_state;
