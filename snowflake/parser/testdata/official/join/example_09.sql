CREATE OR REPLACE TABLE d1 (
  id NUMBER,
  name VARCHAR);

INSERT INTO d1 (id, name) VALUES
  (1,'a'),
  (2,'b'),
  (4,'c');

CREATE OR REPLACE TABLE d2 (
  id NUMBER,
  value VARCHAR);

INSERT INTO d2 (id, value) VALUES
  (1,'xx'),
  (2,'yy'),
  (5,'zz');
