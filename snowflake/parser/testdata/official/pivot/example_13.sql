SELECT *
  FROM quarterly_sales
    PIVOT(SUM(amount) FOR quarter IN (ANY ORDER BY quarter)
      DEFAULT ON NULL (0))
  ORDER BY empid;
