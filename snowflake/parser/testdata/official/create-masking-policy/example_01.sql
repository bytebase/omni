CREATE OR REPLACE MASKING POLICY email_mask AS (val string) returns string ->
  CASE
    WHEN current_role() IN ('ANALYST') THEN VAL
    ELSE '********'
  END;
