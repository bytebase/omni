SELECT *
  FROM (
        SELECT *
          FROM (
                 SELECT *
                   FROM my_table
                   ORDER BY col1
               )
          LIMIT 6
       )
  LIMIT 100;
