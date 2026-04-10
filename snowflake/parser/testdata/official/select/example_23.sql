SELECT t1.product_id AS prod_id, t2.prod_id
  FROM table1 AS t1 JOIN table2 AS t2
    ON t1.product_id=t2.prod_id
  GROUP BY prod_id, t2.prod_id;
