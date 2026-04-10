create masking policy mask_geo_point as (val geography) returns geography ->
  case
    when current_role() IN ('ANALYST') then val
    else to_geography('POINT(-122.35 37.55)')
  end;
