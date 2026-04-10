SELECT *
  FROM monthly_sales
    UNPIVOT INCLUDE NULLS (sales FOR month IN (jan, feb, mar, apr))
  ORDER BY empid;
