TRUNCATE table merge_example_target;

TRUNCATE table merge_example_source;

INSERT INTO merge_example_source (id, description) VALUES
  (50, 'This is a duplicate in the source and has no match in target'),
  (50, 'This is a duplicate in the source and has no match in target');
