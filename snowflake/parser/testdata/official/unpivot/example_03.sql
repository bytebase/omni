SELECT *
  FROM monthly_sales
    UNPIVOT (sales FOR month IN (
      jan AS january,
      feb AS february,
      mar AS march,
      apr AS april)
    )
  ORDER BY empid;
