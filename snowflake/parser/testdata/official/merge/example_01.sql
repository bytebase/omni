CREATE OR REPLACE TABLE merge_example_target (id INTEGER, description VARCHAR);

INSERT INTO merge_example_target (id, description) VALUES
  (10, 'To be updated (this is the old value)');

CREATE OR REPLACE TABLE merge_example_source (id INTEGER, description VARCHAR);

INSERT INTO merge_example_source (id, description) VALUES
  (10, 'To be updated (this is the new value)');
