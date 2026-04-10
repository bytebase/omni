SELECT *
  FROM quarterly_sales
    PIVOT(SUM(amount) AS total FOR quarter IN (ANY ORDER BY quarter))
  ORDER BY empid;
