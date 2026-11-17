-- The schema TestDump dumps and restores. CockroachDB's SHOW CREATE ALL TABLES
-- states tables, indexes, views, sequences and foreign keys; it does not state
-- functions or user-defined types, so those stay out of this corpus.

CREATE SCHEMA app;
COMMENT ON SCHEMA app IS 'the application';

CREATE TABLE app.parent (
    id INT8 PRIMARY KEY,
    label STRING NOT NULL
);
COMMENT ON TABLE app.parent IS 'a parent';
COMMENT ON COLUMN app.parent.label IS 'the label';

CREATE TABLE app.child (
    id INT8 PRIMARY KEY,
    parent_id INT8 NOT NULL REFERENCES app.parent(id) ON DELETE CASCADE ON UPDATE RESTRICT,
    note STRING DEFAULT 'note'
);
CREATE INDEX child_note_idx ON app.child (note);

CREATE VIEW app.child_view AS
    SELECT c.id, c.note, p.label
    FROM app.child AS c JOIN app.parent AS p ON p.id = c.parent_id;

CREATE SEQUENCE app.counter AS INT8 INCREMENT BY 2 START WITH 10;

INSERT INTO app.parent VALUES (1, 'one');
INSERT INTO app.child VALUES (1, 1, 'note');
ANALYZE app.parent;
ANALYZE app.child;
