SELECT *
  FROM quarterly_sales
    PIVOT(SUM(amount) FOR quarter IN (
      SELECT DISTINCT quarter
        FROM ad_campaign_types_by_quarter
        WHERE television = TRUE
        ORDER BY quarter))
  ORDER BY empid;
