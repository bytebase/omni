SELECT *
  FROM quarterly_sales
    PIVOT(SUM(amount)
      FOR quarter IN (
        '2023_Q1',
        '2023_Q2')
      DEFAULT ON NULL (0))
  ORDER BY empid;
