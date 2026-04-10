SELECT SUM($1) AS q1_sales_total,
       SUM($2) AS q2_sales_total,
       SUM($3) AS q3_sales_total,
       SUM($4) AS q4_sales_total,
       MAX($5) AS q1_maximum_discount,
       MAX($6) AS q2_maximum_discount,
       MAX($7) AS q3_maximum_discount,
       MAX($8) AS q4_maximum_discount
  FROM
    (SELECT amount,
            quarter AS quarter_amount,
            quarter AS quarter_discount,
            discount_percent
      FROM quarterly_sales)
  PIVOT (
    SUM(amount)
    FOR quarter_amount IN (
      '2023_Q1',
      '2023_Q2',
      '2023_Q3',
      '2023_Q4'))
  PIVOT (
    MAX(discount_percent)
    FOR quarter_discount IN (
      '2023_Q1',
      '2023_Q2',
      '2023_Q3',
      '2023_Q4'));
