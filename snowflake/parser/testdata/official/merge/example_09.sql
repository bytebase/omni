MERGE INTO merge_example_target_all
  USING merge_example_source_all
  ON merge_example_target_all.id = merge_example_source_all.id
  WHEN MATCHED THEN
    UPDATE ALL BY NAME
  WHEN NOT MATCHED THEN
    INSERT ALL BY NAME;
