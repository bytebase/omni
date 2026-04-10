WITH
  src AS
  (
    SELECT *
      FROM quarterly_sales
        PIVOT(SUM(amount) FOR quarter IN (ANY ORDER BY quarter))
  )
SELECT em.managerid, src.*
  FROM emp_manager em
  JOIN src ON em.empid = src.empid
  ORDER BY empid;
