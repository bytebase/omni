WITH
  sales_without_discount AS
    (SELECT * EXCLUDE(discount_percent) FROM quarterly_sales)
SELECT *
  FROM sales_without_discount
    PIVOT(SUM(amount) FOR quarter IN (ANY ORDER BY quarter))
  ORDER BY empid;
