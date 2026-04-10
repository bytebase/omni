SELECT 'Average sale amount' AS aggregate, *
  FROM quarterly_sales
    PIVOT(AVG(amount) FOR quarter IN (ANY ORDER BY quarter))
UNION
SELECT 'Highest value sale' AS aggregate, *
  FROM quarterly_sales
    PIVOT(MAX(amount) FOR quarter IN (ANY ORDER BY quarter))
UNION
SELECT 'Lowest value sale' AS aggregate, *
  FROM quarterly_sales
    PIVOT(MIN(amount) FOR quarter IN (ANY ORDER BY quarter))
UNION
SELECT 'Number of sales' AS aggregate, *
  FROM quarterly_sales
    PIVOT(COUNT(amount) FOR quarter IN (ANY ORDER BY quarter))
UNION
SELECT 'Total amount' AS aggregate, *
  FROM quarterly_sales
    PIVOT(SUM(amount) FOR quarter IN (ANY ORDER BY quarter))
ORDER BY aggregate, empid;
