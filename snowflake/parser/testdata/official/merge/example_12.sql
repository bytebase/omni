CREATE OR REPLACE TABLE merge_example_target_orig (k NUMBER, v NUMBER);

INSERT INTO merge_example_target_orig VALUES (0, 10);

CREATE OR REPLACE TABLE merge_example_src (k NUMBER, v NUMBER);

INSERT INTO merge_example_src VALUES (0, 11), (0, 12), (0, 13);
