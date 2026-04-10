ALTER TABLE t1
  DROP ROW ACCESS POLICY rap_t1_version_1,
  ADD ROW ACCESS POLICY rap_t1_version_2 ON (empl_id);
