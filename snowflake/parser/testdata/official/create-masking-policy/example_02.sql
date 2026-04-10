create or replace masking policy json_location_mask as (val variant) returns variant ->
  CASE
    WHEN current_role() IN ('ANALYST') THEN val
    else full_location_masking(val)
  END;
