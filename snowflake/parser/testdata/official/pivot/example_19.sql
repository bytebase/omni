WITH
  sales_without_amount AS
    (SELECT * EXCLUDE(amount) FROM quarterly_sales)
SELECT *
  FROM sales_without_amount
    PIVOT(AVG(discount_percent) FOR quarter IN (ANY ORDER BY quarter))
  ORDER BY empid;
