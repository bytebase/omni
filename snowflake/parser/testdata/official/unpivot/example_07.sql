CREATE OR REPLACE TEMPORARY TABLE unpivot_result AS
  SELECT *
    FROM monthly_sales
      UNPIVOT (sales FOR month IN (jan, feb, mar, apr));

SELECT *
  FROM unpivot_result
    JOIN LATERAL (SELECT unpivot_result.sales AS sales_value) jl;
