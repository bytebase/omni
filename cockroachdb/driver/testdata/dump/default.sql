CREATE SCHEMA app;
COMMENT ON SCHEMA app IS 'the application';

CREATE TABLE app.parent (
	id INT8 NOT NULL,
	label STRING NOT NULL,
	CONSTRAINT parent_pkey PRIMARY KEY (id ASC)
);
COMMENT ON TABLE app.parent IS 'a parent';
COMMENT ON COLUMN app.parent.label IS 'the label';

CREATE TABLE app.child (
	id INT8 NOT NULL,
	parent_id INT8 NOT NULL,
	note STRING NULL DEFAULT 'note':::STRING,
	CONSTRAINT child_pkey PRIMARY KEY (id ASC),
	INDEX child_note_idx (note ASC)
);

CREATE VIEW app.child_view (
	id,
	note,
	label
) AS SELECT
		c.id, c.note, p.label
	FROM
		omni_testdump.app.child AS c JOIN omni_testdump.app.parent AS p ON p.id = c.parent_id;

CREATE SEQUENCE app.counter AS INT8 MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 2 START 10;

ALTER TABLE app.child ADD CONSTRAINT child_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES app.parent(id) ON DELETE CASCADE ON UPDATE RESTRICT;

-- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES

ALTER TABLE app.child VALIDATE CONSTRAINT child_parent_id_fkey;

