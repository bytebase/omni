SELECT *
  FROM monthly_sales
    UNPIVOT (sales FOR month IN (jan, feb, mar, apr)) unpvt
    JOIN LATERAL (SELECT unpvt.sales AS sales_value) jl;
