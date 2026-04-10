ALTER VIEW v1
  DROP ROW ACCESS POLICY rap_v1_version_1,
  ADD ROW ACCESS POLICY rap_v1_version_2 ON (empl_id);
