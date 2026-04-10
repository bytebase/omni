CREATE OR REPLACE TABLE merge_example_mult_target (
  id INTEGER,
  val INTEGER,
  status VARCHAR);

INSERT INTO merge_example_mult_target (id, val, status) VALUES
  (1, 10, 'Production'),
  (2, 20, 'Alpha'),
  (3, 30, 'Production');

CREATE OR REPLACE TABLE merge_example_mult_source (
  id INTEGER,
  marked VARCHAR,
  isnewstatus INTEGER,
  newval INTEGER,
  newstatus VARCHAR);

INSERT INTO merge_example_mult_source (id, marked, isnewstatus, newval, newstatus) VALUES
  (1, 'Y', 0, 10, 'Production'),
  (2, 'N', 1, 50, 'Beta'),
  (3, 'N', 0, 60, 'Deprecated'),
  (4, 'N', 0, 40, 'Production');
