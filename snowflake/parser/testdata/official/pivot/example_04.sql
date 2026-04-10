CREATE OR REPLACE TABLE ad_campaign_types_by_quarter(
  quarter VARCHAR,
  television BOOLEAN,
  radio BOOLEAN,
  print BOOLEAN)
  AS SELECT * FROM VALUES
    ('2023_Q1', TRUE, FALSE, FALSE),
    ('2023_Q2', FALSE, TRUE, TRUE),
    ('2023_Q3', FALSE, TRUE, FALSE),
    ('2023_Q4', TRUE, FALSE, TRUE);
