MERGE INTO merge_example_mult_target
  USING merge_example_mult_source
  ON merge_example_mult_target.id = merge_example_mult_source.id
  WHEN MATCHED AND merge_example_mult_source.marked = 'Y'
    THEN DELETE
  WHEN MATCHED AND merge_example_mult_source.isnewstatus = 1
    THEN UPDATE SET val = merge_example_mult_source.newval, status = merge_example_mult_source.newstatus
  WHEN MATCHED
    THEN UPDATE SET val = merge_example_mult_source.newval
  WHEN NOT MATCHED
    THEN INSERT (id, val, status) VALUES (
      merge_example_mult_source.id,
      merge_example_mult_source.newval,
      merge_example_mult_source.newstatus);
