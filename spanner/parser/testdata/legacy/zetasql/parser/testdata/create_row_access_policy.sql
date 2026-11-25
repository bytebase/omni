CREATE ROW ACCESS POLICY ON t1 GRANT TO ('foo@google.com') FILTER USING (c1 = 'foo');
CREATE ROW ACCESS POLICY p1 ON t1 GRANT TO ('mdbuser/bar') FILTER USING (c2 = 'foo');
CREATE ROW ACCESS POLICY ON t1 GRANT TO ('foo@google.com', 'mdbgroup/bar') FILTER USING (c1);
CREATE ROW ACCESS POLICY ON n1.t1 GRANT TO ('foo@google.com', 'mdbgroup/bar') FILTER USING (c1);
CREATE ROW ACCESS POLICY ON n1.t1 GRANT TO ('foo@google.com', 'mdbgroup/bar') FILTER USING (1);
CREATE ROW ACCESS POLICY ON t1 GRANT TO (@test_param_string) FILTER USING (true);
CREATE ROW ACCESS POLICY p1 ON t1 FILTER USING (true);
