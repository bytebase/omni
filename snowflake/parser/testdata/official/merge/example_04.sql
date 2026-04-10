MERGE INTO merge_example_target
  USING merge_example_source
  ON merge_example_target.id = merge_example_source.id
  WHEN MATCHED THEN
    UPDATE SET merge_example_target.description = merge_example_source.description;
