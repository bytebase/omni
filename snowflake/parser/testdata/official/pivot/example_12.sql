SELECT *
  FROM quarterly_sales
    PIVOT(SUM(amount) FOR quarter IN (
      '2023_Q1' AS q1,
      '2023_Q2' AS q2,
      '2023_Q3' AS q3,
      '2023_Q4' AS q4))
  ORDER BY empid;
