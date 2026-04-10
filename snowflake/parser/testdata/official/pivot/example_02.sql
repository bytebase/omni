SELECT *
  FROM quarterly_sales
    PIVOT(SUM(amount) FOR quarter IN (ANY ORDER BY quarter))
  ORDER BY empid;
