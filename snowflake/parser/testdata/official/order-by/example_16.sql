SELECT *
  FROM (
    SELECT branch_name
      FROM branch_offices
      ORDER BY monthly_sales DESC
      LIMIT 3
  );
