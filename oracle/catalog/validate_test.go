package catalog

import (
	"context"
	"errors"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

func TestExecOracleErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		setup    string
		stmt     string
		wantCode string
	}{
		{
			name:     "alter missing table",
			setup:    `CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY)`,
			stmt:     `ALTER TABLE a ADD (c NUMBER)`,
			wantCode: "00942",
		},
		{
			name:     "create duplicate table",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `CREATE TABLE a (c NUMBER)`,
			wantCode: "00955",
		},
		{
			name:     "create table duplicate column",
			stmt:     `CREATE TABLE a (c NUMBER, c NUMBER)`,
			wantCode: "00957",
		},
		{
			name:     "create table missing datatype",
			stmt:     `CREATE TABLE a (c)`,
			wantCode: "02263",
		},
		{
			name:     "create table all virtual columns",
			stmt:     `CREATE TABLE a (c AS (1))`,
			wantCode: "54037",
		},
		{
			name:     "create table as select unaliased expression",
			stmt:     `CREATE TABLE a AS SELECT 1 FROM dual`,
			wantCode: "00998",
		},
		{
			name:     "create table missing domain",
			stmt:     `CREATE TABLE a (c DOMAIN missing_dom)`,
			wantCode: "11504",
		},
		{
			name:     "alter table add existing column",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `ALTER TABLE a ADD (b NUMBER)`,
			wantCode: "01430",
		},
		{
			name:     "alter table drop missing column",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `ALTER TABLE a DROP COLUMN c`,
			wantCode: "00904",
		},
		{
			name:     "alter table modify real column to virtual",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `ALTER TABLE a MODIFY (b AS (1))`,
			wantCode: "54026",
		},
		{
			name:     "alter table modify non-identity column as identity",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `ALTER TABLE a MODIFY (b GENERATED ALWAYS AS IDENTITY)`,
			wantCode: "30673",
		},
		{
			name:     "alter table modify nullable column to null",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `ALTER TABLE a MODIFY (b NULL)`,
			wantCode: "01451",
		},
		{
			name:     "alter table second primary key",
			setup:    `CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER)`,
			stmt:     `ALTER TABLE a ADD PRIMARY KEY (c)`,
			wantCode: "02260",
		},
		{
			name:     "alter table duplicate constraint name",
			setup:    `CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER)`,
			stmt:     `ALTER TABLE a ADD CONSTRAINT a_pk UNIQUE (c)`,
			wantCode: "02264",
		},
		{
			name:     "alter table foreign key without matching parent key",
			setup:    `CREATE TABLE a (b NUMBER CONSTRAINT a_pk PRIMARY KEY, c NUMBER)`,
			stmt:     `ALTER TABLE a ADD FOREIGN KEY (c) REFERENCES a (c)`,
			wantCode: "02270",
		},
		{
			name:     "alter table drop nonexistent constraint",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `ALTER TABLE a DROP CONSTRAINT no_such`,
			wantCode: "02443",
		},
		{
			name:     "alter table drop nonexistent primary key",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `ALTER TABLE a DROP PRIMARY KEY`,
			wantCode: "02441",
		},
		{
			name:     "alter table drop nonexistent unique key",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `ALTER TABLE a DROP UNIQUE (b)`,
			wantCode: "02442",
		},
		{
			name:     "alter table drop referenced primary key without cascade",
			setup:    `CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY); CREATE TABLE c (p_b NUMBER CONSTRAINT c_fk REFERENCES p (b))`,
			stmt:     `ALTER TABLE p DROP PRIMARY KEY`,
			wantCode: "02273",
		},
		{
			name:     "alter table drop referenced primary key online without cascade",
			setup:    `CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY); CREATE TABLE c (p_b NUMBER CONSTRAINT c_fk REFERENCES p (b))`,
			stmt:     `ALTER TABLE p DROP PRIMARY KEY ONLINE`,
			wantCode: "02297",
		},
		{
			name:     "alter table rename nonexistent constraint",
			setup:    `CREATE TABLE a (b NUMBER)`,
			stmt:     `ALTER TABLE a RENAME CONSTRAINT no_such TO other_c`,
			wantCode: "23292",
		},
		{
			name:     "alter table encrypt fk column",
			setup:    `CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY); CREATE TABLE c (p_b NUMBER CONSTRAINT c_fk REFERENCES p (b))`,
			stmt:     `ALTER TABLE c MODIFY (p_b ENCRYPT)`,
			wantCode: "28335",
		},
		{
			name:     "drop table referenced by fk without cascade",
			setup:    `CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY); CREATE TABLE c (p_b NUMBER CONSTRAINT c_fk REFERENCES p (b))`,
			stmt:     `DROP TABLE p`,
			wantCode: "02449",
		},
		{
			name:     "drop missing index",
			stmt:     `DROP INDEX missing_idx`,
			wantCode: "01418",
		},
		{
			name:     "drop non-domain index with force",
			setup:    `CREATE TABLE t (b NUMBER); CREATE INDEX idx_b ON t (b)`,
			stmt:     `DROP INDEX idx_b FORCE`,
			wantCode: "29862",
		},
		{
			name:     "drop missing sequence",
			stmt:     `DROP SEQUENCE missing_seq`,
			wantCode: "02289",
		},
		{
			name:     "drop missing private synonym",
			stmt:     `DROP SYNONYM missing_syn`,
			wantCode: "01434",
		},
		{
			name:     "drop missing public synonym",
			stmt:     `DROP PUBLIC SYNONYM missing_syn`,
			wantCode: "01432",
		},
		{
			name:     "create index on view",
			setup:    `CREATE TABLE t (b NUMBER); CREATE VIEW v AS SELECT b FROM t`,
			stmt:     `CREATE INDEX idx_v ON v (b)`,
			wantCode: "01702",
		},
		{
			name:     "create index on sequence",
			setup:    `CREATE SEQUENCE s`,
			stmt:     `CREATE INDEX idx_s ON s (b)`,
			wantCode: "02201",
		},
		{
			name:     "create view unaliased expression",
			stmt:     `CREATE VIEW v AS SELECT 1 FROM dual`,
			wantCode: "00998",
		},
		{
			name:     "create view all invisible columns",
			stmt:     `CREATE VIEW v (a INVISIBLE) AS SELECT 1 FROM dual`,
			wantCode: "54039",
		},
		{
			name:     "alter view read only on non-editioning view",
			setup:    `CREATE TABLE t (b NUMBER); CREATE VIEW v AS SELECT b FROM t`,
			stmt:     `ALTER VIEW v READ ONLY`,
			wantCode: "42310",
		},
		{
			name:     "create sequence cache 1",
			stmt:     `CREATE SEQUENCE s CACHE 1`,
			wantCode: "04010",
		},
		{
			name:     "create sequence cycle without maxvalue",
			stmt:     `CREATE SEQUENCE s CYCLE`,
			wantCode: "04015",
		},
		{
			name:     "create sequence maxvalue <= minvalue",
			stmt:     `CREATE SEQUENCE s MAXVALUE 1`,
			wantCode: "04004",
		},
		{
			name:     "create self-referencing private synonym",
			stmt:     `CREATE SYNONYM s FOR s`,
			wantCode: "01471",
		},
		{
			name:     "truncate view",
			setup:    `CREATE TABLE t (b NUMBER); CREATE VIEW v AS SELECT b FROM t`,
			stmt:     `TRUNCATE TABLE v`,
			wantCode: "01702",
		},
		{
			name:     "rename missing object",
			stmt:     `RENAME missing_obj TO other_obj`,
			wantCode: "04043",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var c *Catalog
			var err error
			if tc.setup != "" {
				c, err = LoadSQL(tc.setup)
				require.NoError(t, err)
			} else {
				c = New()
			}
			results, err := c.Exec(tc.stmt, nil)
			require.NoError(t, err)
			require.Len(t, results, 1)
			var oraErr *Error
			require.True(t, errors.As(results[0].Error, &oraErr), "expected *catalog.Error, got %v", results[0].Error)
			require.Equal(t, tc.wantCode, oraErr.Code)
		})
	}
}

func TestExecSchemaStateFidelityFixes(t *testing.T) {
	t.Run("DROP PRIMARY KEY KEEP INDEX preserves unique index", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, c NUMBER)`)
		require.NoError(t, err)
		results, err := c.Exec(`ALTER TABLE t DROP PRIMARY KEY KEEP INDEX`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		tbl := tableByName(c.Metadata().Schemas[0])["T"]
		require.True(t, tbl.Columns[0].Nullable)
		require.Len(t, tbl.Indexes, 1)
		require.Equal(t, "T_PK", tbl.Indexes[0].Name)
		require.False(t, tbl.Indexes[0].Primary)
		require.False(t, tbl.Indexes[0].IsConstraint)
		require.True(t, tbl.Indexes[0].Unique)
	})

	t.Run("ALTER TABLE MODIFY applies inline UNIQUE, CHECK, and REFERENCES", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY); CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, x NUMBER)`)
		require.NoError(t, err)
		for _, stmt := range []string{
			`ALTER TABLE t MODIFY (x UNIQUE)`,
			`ALTER TABLE t MODIFY (x CHECK (x > 0))`,
			`ALTER TABLE t MODIFY (x REFERENCES p ON DELETE CASCADE)`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			require.Nil(t, results[0].Error, stmt)
		}
		tbl := tableByName(c.Metadata().Schemas[0])["T"]
		require.Len(t, tbl.Indexes, 2)
		require.Len(t, tbl.CheckConstraints, 1)
		require.Len(t, tbl.ForeignKeys, 1)
		require.Equal(t, "CASCADE", tbl.ForeignKeys[0].OnDelete)
	})

	t.Run("COMMENT ON TABLE and CREATE INDEX resolve private synonyms", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE t (b NUMBER, c NUMBER); CREATE SYNONYM s FOR t`)
		require.NoError(t, err)
		results, err := c.Exec(`COMMENT ON TABLE s IS 'via synonym'; CREATE INDEX s_idx ON s (c)`, nil)
		require.NoError(t, err)
		for _, r := range results {
			require.Nil(t, r.Error)
		}
		tbl := tableByName(c.Metadata().Schemas[0])["T"]
		require.Equal(t, "via synonym", tbl.Comment)
		require.Len(t, tbl.Indexes, 1)
		require.Equal(t, "S_IDX", tbl.Indexes[0].Name)
	})

	t.Run("RENAME renames private synonyms", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE t (b NUMBER); CREATE SYNONYM s1 FOR t`)
		require.NoError(t, err)
		results, err := c.Exec(`RENAME s1 TO s2`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		syns := synonymByName(c.Metadata().Schemas[0].Synonyms)
		require.Nil(t, syns["S1"])
		require.NotNil(t, syns["S2"])
		require.Equal(t, "T", syns["S2"].TargetName)
	})
}

func TestExecForeignKeyAgainstConstraintMetadata(t *testing.T) {
	newCatalog := func(t *testing.T, parentIndexes []*metadatapb.IndexMetadata) *Catalog {
		t.Helper()
		c := New()
		_, err := c.LoadMetadata(context.Background(), &metadatapb.DatabaseSchemaMetadata{
			Name: "APP",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "APP",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "PARENT",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "ID", Type: "NUMBER", Position: 1, Nullable: false},
							{Name: "CODE", Type: "NUMBER", Position: 2, Nullable: true},
							{Name: "EXTRA", Type: "NUMBER", Position: 3, Nullable: true},
						},
						Indexes: parentIndexes,
					},
					{
						Name: "CHILD",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "P_ID", Type: "NUMBER", Position: 1, Nullable: true},
							{Name: "P_CODE", Type: "NUMBER", Position: 2, Nullable: true},
							{Name: "P_EXTRA", Type: "NUMBER", Position: 3, Nullable: true},
						},
					},
				},
			}},
		})
		require.NoError(t, err)
		return c
	}

	t.Run("PK reusing nonunique index accepts matching constraint columns and rejects index superset", func(t *testing.T) {
		c := newCatalog(t, []*metadatapb.IndexMetadata{{
			Name:         "PK_REUSED_IDX",
			Type:         "NORMAL",
			Primary:      true,
			Unique:       false,
			Expressions:  []string{`"ID"`, `"EXTRA"`},
			Visible:      true,
			IsConstraint: true,
			Constraints: []*metadatapb.IndexConstraintMetadata{
				{Name: "PK_REUSED", Columns: []string{`"ID"`}, Primary: true},
			},
		}})

		results, err := c.Exec(`ALTER TABLE child ADD CONSTRAINT fk_id FOREIGN KEY (p_id) REFERENCES parent (id)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		results, err = c.Exec(`ALTER TABLE child ADD CONSTRAINT fk_both FOREIGN KEY (p_id, p_extra) REFERENCES parent (id, extra)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "02270", oraErr.Code)
	})

	t.Run("UK reusing nonunique index accepts matching constraint columns", func(t *testing.T) {
		c := newCatalog(t, []*metadatapb.IndexMetadata{{
			Name:         "UK_REUSED_IDX",
			Type:         "NORMAL",
			Primary:      false,
			Unique:       false,
			Expressions:  []string{`"CODE"`, `"EXTRA"`},
			Visible:      true,
			IsConstraint: true,
			Constraints: []*metadatapb.IndexConstraintMetadata{
				{Name: "UK_REUSED", Columns: []string{`"CODE"`}, Primary: false},
			},
		}})

		results, err := c.Exec(`ALTER TABLE child ADD CONSTRAINT fk_code FOREIGN KEY (p_code) REFERENCES parent (code)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
	})

	t.Run("shared index enforcing PK and UK accepts FKs to either constraint", func(t *testing.T) {
		c := newCatalog(t, []*metadatapb.IndexMetadata{{
			Name:         "SHARED_IDX",
			Type:         "NORMAL",
			Primary:      true,
			Unique:       false,
			Expressions:  []string{`"ID"`, `"CODE"`},
			Visible:      true,
			IsConstraint: true,
			Constraints: []*metadatapb.IndexConstraintMetadata{
				{Name: "SHARED_PK", Columns: []string{`"ID"`}, Primary: true},
				{Name: "SHARED_UK", Columns: []string{`"CODE"`}, Primary: false},
			},
		}})

		results, err := c.Exec(`
			ALTER TABLE child ADD CONSTRAINT fk_id FOREIGN KEY (p_id) REFERENCES parent (id);
			ALTER TABLE child ADD CONSTRAINT fk_code FOREIGN KEY (p_code) REFERENCES parent (code);
		`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)
	})

	t.Run("plain non-constraint unique index is rejected for FK reference", func(t *testing.T) {
		c := newCatalog(t, []*metadatapb.IndexMetadata{{
			Name:         "PLAIN_UQ_IDX",
			Type:         "NORMAL",
			Primary:      false,
			Unique:       true,
			Expressions:  []string{`"CODE"`},
			Visible:      true,
			IsConstraint: false,
		}})

		results, err := c.Exec(`ALTER TABLE child ADD CONSTRAINT fk_code FOREIGN KEY (p_code) REFERENCES parent (code)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "02270", oraErr.Code)
	})
}

func TestExecReviewRegressionFixes(t *testing.T) {
	t.Run("drop CHECK constraint sharing name with standalone index on referenced key succeeds", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE p (id NUMBER CONSTRAINT p_pk PRIMARY KEY, val NUMBER CONSTRAINT p_idx CHECK (val > 0));
			CREATE INDEX p_idx ON p (id);
			CREATE TABLE ch (p_id NUMBER CONSTRAINT ch_fk REFERENCES p (id));
		`)
		require.NoError(t, err)

		results, err := c.Exec(`ALTER TABLE p DROP CONSTRAINT p_idx`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		tbl := tableByName(c.Metadata().Schemas[0])["P"]
		require.Empty(t, tbl.CheckConstraints)
		require.Len(t, tbl.Indexes, 2) // p_pk and standalone p_idx remain
	})

	t.Run("drop FK constraint sharing name with standalone index on referenced key succeeds", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE p (id NUMBER CONSTRAINT p_pk PRIMARY KEY, parent_id NUMBER);
			CREATE INDEX same_name ON p (id);
			ALTER TABLE p ADD CONSTRAINT same_name FOREIGN KEY (parent_id) REFERENCES p (id);
		`)
		require.NoError(t, err)

		results, err := c.Exec(`ALTER TABLE p DROP CONSTRAINT same_name`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Empty(t, tableByName(c.Metadata().Schemas[0])["P"].ForeignKeys)
	})

	t.Run("ALTER VIEW retains ADD CONSTRAINT and CREATE VIEW constraints for MODIFY and DROP", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE t (id NUMBER, code NUMBER);
			CREATE VIEW v AS SELECT id, code FROM t;
		`)
		require.NoError(t, err)

		for _, stmt := range []string{
			`ALTER VIEW v ADD CONSTRAINT v_pk PRIMARY KEY (id) DISABLE NOVALIDATE`,
			`ALTER VIEW v MODIFY CONSTRAINT v_pk RELY`,
			`ALTER VIEW v MODIFY PRIMARY KEY NORELY`,
			`ALTER VIEW v DROP CONSTRAINT v_pk`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			require.Nil(t, results[0].Error, stmt)
		}

		// Dropping again after removal returns ORA-02443.
		results, err := c.Exec(`ALTER VIEW v DROP CONSTRAINT v_pk`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "02443", oraErr.Code)

		// Constraints declared directly in CREATE VIEW (out-of-line and inline) are also retained.
		c2, err := LoadSQL(`
			CREATE TABLE t (id NUMBER, code NUMBER);
			CREATE VIEW v2 (id CONSTRAINT v2_pk PRIMARY KEY DISABLE NOVALIDATE, code, CONSTRAINT v2_uk UNIQUE (code) DISABLE NOVALIDATE) AS SELECT id, code FROM t;
		`)
		require.NoError(t, err)
		for _, stmt := range []string{
			`ALTER VIEW v2 MODIFY CONSTRAINT v2_pk RELY`,
			`ALTER VIEW v2 DROP PRIMARY KEY`,
			`ALTER VIEW v2 MODIFY UNIQUE (code) RELY`,
			`ALTER VIEW v2 DROP UNIQUE (code)`,
		} {
			results, err := c2.Exec(stmt, nil)
			require.NoError(t, err)
			require.Nil(t, results[0].Error, stmt)
		}
	})

	t.Run("synthetic constraint names do not block explicit user constraint names", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE t (id NUMBER PRIMARY KEY, x NUMBER);
			CREATE TABLE u (id NUMBER);
		`)
		require.NoError(t, err)

		// Cross-table: PK_T was synthetic on T, so explicit PK_T on U must succeed.
		results, err := c.Exec(`ALTER TABLE u ADD CONSTRAINT pk_t PRIMARY KEY (id)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		// Second explicit PK_T must still fail with ORA-02264.
		results, err = c.Exec(`ALTER TABLE t ADD CONSTRAINT pk_t UNIQUE (x)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "02264", oraErr.Code)

		// Same-table: unnamed PK on A does not block explicit PK_A UNIQUE on A.
		c2, err := LoadSQL(`CREATE TABLE a (id NUMBER PRIMARY KEY, x NUMBER)`)
		require.NoError(t, err)
		results, err = c2.Exec(`ALTER TABLE a ADD CONSTRAINT pk_a UNIQUE (x)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		// Dropping explicit PK_A removes the unique constraint, leaving the primary key intact.
		results, err = c2.Exec(`ALTER TABLE a DROP CONSTRAINT pk_a`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		tblA := tableByName(c2.Metadata().Schemas[0])["A"]
		require.Len(t, tblA.Indexes, 1)
		require.True(t, tblA.Indexes[0].Primary)
	})

	t.Run("TRUNCATE TABLE accepts materialized views from CREATE and LoadMetadata", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE t (id NUMBER);
			CREATE MATERIALIZED VIEW mv AS SELECT id FROM t;
		`)
		require.NoError(t, err)
		results, err := c.Exec(`TRUNCATE TABLE mv`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		c2 := New()
		_, err = c2.LoadMetadata(context.Background(), &metadatapb.DatabaseSchemaMetadata{
			Name: "APP",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "APP",
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{Name: "MV_SYNCED", Definition: "SELECT 1 FROM DUAL"},
				},
			}},
		})
		require.NoError(t, err)
		results, err = c2.Exec(`TRUNCATE TABLE mv_synced`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
	})

	t.Run("ALTER TABLE accepts materialized-view containers from CREATE and LoadMetadata", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE t (id NUMBER);
			CREATE MATERIALIZED VIEW mv AS SELECT id FROM t;
		`)
		require.NoError(t, err)
		results, err := c.Exec(`ALTER TABLE mv LOGGING; COMMENT ON TABLE mv IS 'container comment'`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)

		c2 := New()
		_, err = c2.LoadMetadata(context.Background(), &metadatapb.DatabaseSchemaMetadata{
			Name: "APP",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "APP",
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{Name: "MV_SYNCED", Definition: "SELECT 1 FROM DUAL"},
				},
			}},
		})
		require.NoError(t, err)
		results, err = c2.Exec(`ALTER TABLE mv_synced LOGGING`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
	})

	t.Run("editioning views allow ALTER VIEW READ ONLY and READ WRITE", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE t (id NUMBER);
			CREATE EDITIONING VIEW ev AS SELECT id FROM t;
		`)
		require.NoError(t, err)
		for _, stmt := range []string{
			`ALTER VIEW ev READ ONLY`,
			`ALTER VIEW ev READ WRITE`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			require.Nil(t, results[0].Error, stmt)
		}

		// Replacing with a non-editioning view restores the ORA-42310 check.
		results, err := c.Exec(`CREATE OR REPLACE VIEW ev AS SELECT id FROM t; ALTER VIEW ev READ ONLY`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		var oraErr *Error
		require.True(t, errors.As(results[1].Error, &oraErr))
		require.Equal(t, "42310", oraErr.Code)
	})

	t.Run("RENAME COLUMN updates IndexMetadata.Constraints for FK validation", func(t *testing.T) {
		c := New()
		_, err := c.LoadMetadata(context.Background(), &metadatapb.DatabaseSchemaMetadata{
			Name: "APP",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "APP",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "PARENT",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "ID", Type: "NUMBER", Position: 1, Nullable: false},
						},
						Indexes: []*metadatapb.IndexMetadata{{
							Name:         "P_IDX",
							Type:         "NORMAL",
							Primary:      true,
							Unique:       false,
							Expressions:  []string{`"ID"`},
							Visible:      true,
							IsConstraint: true,
							Constraints: []*metadatapb.IndexConstraintMetadata{
								{Name: "P_PK", Columns: []string{`"ID"`}, Primary: true},
							},
						}},
					},
					{
						Name: "CHILD",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "P_ID", Type: "NUMBER", Position: 1, Nullable: true},
						},
					},
				},
			}},
		})
		require.NoError(t, err)

		results, err := c.Exec(`
			ALTER TABLE parent RENAME COLUMN id TO new_id;
			ALTER TABLE child ADD CONSTRAINT ch_fk FOREIGN KEY (p_id) REFERENCES parent (new_id);
		`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)
	})

	t.Run("encrypted columns can be decrypted and then reject a second DECRYPT", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE t (id NUMBER ENCRYPT)`)
		require.NoError(t, err)

		results, err := c.Exec(`ALTER TABLE t MODIFY (id DECRYPT)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		results, err = c.Exec(`ALTER TABLE t MODIFY (id DECRYPT)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "28333", oraErr.Code)
	})

	t.Run("sequence validation preserves out-of-int64 bounds", func(t *testing.T) {
		c, err := LoadSQL(`CREATE SEQUENCE s MINVALUE -10000000000000000000 MAXVALUE 0 START WITH -1`)
		require.NoError(t, err)
		require.NotNil(t, sequenceByName(c.Metadata().Schemas[0])["S"])

		results, err := c.Exec(`CREATE SEQUENCE s_bad MINVALUE 100000000000000000000 MAXVALUE 10000000000000000000`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "04004", oraErr.Code)
	})

	t.Run("DROP INDEX FORCE accepts domain indexes from LoadMetadata and rejects non-domain indexes", func(t *testing.T) {
		c := New()
		_, err := c.LoadMetadata(context.Background(), &metadatapb.DatabaseSchemaMetadata{
			Name: "APP",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "APP",
				Tables: []*metadatapb.TableMetadata{{
					Name:    "T",
					Columns: []*metadatapb.ColumnMetadata{{Name: "DOC", Type: "VARCHAR2(100)", Position: 1, Nullable: true}},
					Indexes: []*metadatapb.IndexMetadata{
						{Name: "DOM_IDX", Type: "DOMAIN", Expressions: []string{`"DOC"`}, Visible: true},
						{Name: "NORM_IDX", Type: "NORMAL", Expressions: []string{`"DOC"`}, Visible: true},
					},
				}},
			}},
		})
		require.NoError(t, err)

		results, err := c.Exec(`DROP INDEX dom_idx FORCE`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		results, err = c.Exec(`DROP INDEX norm_idx FORCE`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "29862", oraErr.Code)
	})

	t.Run("synthetic key-index names do not block explicit CREATE INDEX", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE t (id NUMBER PRIMARY KEY, x NUMBER)`)
		require.NoError(t, err)

		results, err := c.Exec(`CREATE INDEX pk_t ON t (x)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		// Second explicit index with the same name must still fail with ORA-00955.
		results, err = c.Exec(`CREATE INDEX pk_t ON t (x)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "00955", oraErr.Code)
	})

	t.Run("ALTER TABLE ADD REF column allows WITH ROWID and SCOPE IS while non-REF rejects", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE r (id NUMBER)`)
		require.NoError(t, err)

		results, err := c.Exec(`
			ALTER TABLE r ADD (c1 REF obj_t WITH ROWID);
			ALTER TABLE r ADD (c2 REF obj_t SCOPE IS other_tbl);
		`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)

		results, err = c.Exec(`ALTER TABLE r ADD (bad NUMBER WITH ROWID)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "22893", oraErr.Code)
	})

	t.Run("DROP PRIMARY KEY KEEP INDEX removes constraint metadata from retained synced index", func(t *testing.T) {
		c := New()
		_, err := c.LoadMetadata(context.Background(), &metadatapb.DatabaseSchemaMetadata{
			Name: "APP",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "APP",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "PARENT",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "ID", Type: "NUMBER", Position: 1, Nullable: false},
							{Name: "CODE", Type: "NUMBER", Position: 2, Nullable: true},
						},
						Indexes: []*metadatapb.IndexMetadata{{
							Name:         "SHARED_IDX",
							Type:         "NORMAL",
							Primary:      true,
							Unique:       false,
							Expressions:  []string{`"ID"`, `"CODE"`},
							Visible:      true,
							IsConstraint: true,
							Constraints: []*metadatapb.IndexConstraintMetadata{
								{Name: "SHARED_PK", Columns: []string{`"ID"`}, Primary: true},
								{Name: "SHARED_UK", Columns: []string{`"CODE"`}, Primary: false},
							},
						}},
					},
					{
						Name: "CHILD",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "P_ID", Type: "NUMBER", Position: 1, Nullable: true},
							{Name: "P_CODE", Type: "NUMBER", Position: 2, Nullable: true},
						},
					},
				},
			}},
		})
		require.NoError(t, err)

		// Dropping PK keeps the shared index and its remaining UK constraint.
		results, err := c.Exec(`ALTER TABLE parent DROP PRIMARY KEY KEEP INDEX`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		// FK to dropped PK column now fails with ORA-02270, while FK to remaining UK succeeds.
		results, err = c.Exec(`ALTER TABLE child ADD FOREIGN KEY (p_id) REFERENCES parent (id)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "02270", oraErr.Code)

		results, err = c.Exec(`ALTER TABLE child ADD FOREIGN KEY (p_code) REFERENCES parent (code)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		// Adding a new primary key now succeeds without ORA-02260.
		results, err = c.Exec(`ALTER TABLE parent ADD CONSTRAINT new_pk PRIMARY KEY (id)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
	})

	t.Run("ALTER TABLE ADD validates self-referencing FK after installing all added columns", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE t (x NUMBER)`)
		require.NoError(t, err)

		results, err := c.Exec(`ALTER TABLE t ADD (parent_id NUMBER REFERENCES t (id), id NUMBER PRIMARY KEY)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		tbl := tableByName(c.Metadata().Schemas[0])["T"]
		require.Len(t, tbl.ForeignKeys, 1)
		require.Len(t, tbl.Indexes, 1)
	})

	t.Run("ALTER TABLE MODIFY inline UNIQUE reuses existing standalone unique index", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE t (id NUMBER);
			CREATE UNIQUE INDEX t_idx ON t (id);
			CREATE TABLE ch (t_id NUMBER);
		`)
		require.NoError(t, err)

		results, err := c.Exec(`ALTER TABLE t MODIFY (id CONSTRAINT t_uk UNIQUE)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		// Foreign key can now reference the unique constraint on t(id).
		results, err = c.Exec(`ALTER TABLE ch ADD CONSTRAINT ch_fk FOREIGN KEY (t_id) REFERENCES t (id)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		// Dropping t_uk with CASCADE succeeds (constraint was retained on t_idx).
		results, err = c.Exec(`ALTER TABLE t DROP CONSTRAINT t_uk CASCADE`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		// The pre-existing unique index t_idx is preserved after dropping t_uk and can be dropped separately.
		results, err = c.Exec(`DROP INDEX t_idx`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
	})

	t.Run("dropping constraint backed by reused nonunique index preserves physical index for subsequent DROP INDEX", func(t *testing.T) {
		c := New()
		_, err := c.LoadMetadata(context.Background(), &metadatapb.DatabaseSchemaMetadata{
			Name: "APP",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "APP",
				Tables: []*metadatapb.TableMetadata{{
					Name:    "T",
					Columns: []*metadatapb.ColumnMetadata{{Name: "ID", Type: "NUMBER", Position: 1, Nullable: false}},
					Indexes: []*metadatapb.IndexMetadata{{
						Name:         "T_ID_IDX",
						Type:         "NORMAL",
						Primary:      true,
						Unique:       false,
						Expressions:  []string{`"ID"`},
						Visible:      true,
						IsConstraint: true,
						Constraints: []*metadatapb.IndexConstraintMetadata{
							{Name: "T_PK", Columns: []string{`"ID"`}, Primary: true},
						},
					}},
				}},
			}},
		})
		require.NoError(t, err)

		// Mirrors writeDropConstraintIndex migration replay: DROP CONSTRAINT then DROP INDEX.
		results, err := c.Exec(`ALTER TABLE t DROP CONSTRAINT t_pk; DROP INDEX t_id_idx`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)
		require.Empty(t, tableByName(c.Metadata().Schemas[0])["T"].Indexes)
	})

	t.Run("ALTER TABLE MOVE follows private and public synonyms while ADD does not", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE t (id NUMBER);
			CREATE SYNONYM s FOR t;
			CREATE PUBLIC SYNONYM ps FOR t;
		`)
		require.NoError(t, err)

		results, err := c.Exec(`ALTER TABLE s MOVE; ALTER TABLE ps MOVE`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)

		results, err = c.Exec(`ALTER TABLE s ADD (x NUMBER)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "00942", oraErr.Code)
	})

	t.Run("CREATE DOMAIN allows domain-backed columns and DROP DOMAIN revokes them", func(t *testing.T) {
		c, err := LoadSQL(`CREATE DOMAIN d AS NUMBER`)
		require.NoError(t, err)

		results, err := c.Exec(`CREATE TABLE t (id DOMAIN d)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		results, err = c.Exec(`DROP DOMAIN d; CREATE TABLE t2 (id DOMAIN d)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		var oraErr *Error
		require.True(t, errors.As(results[1].Error, &oraErr))
		require.Equal(t, "11504", oraErr.Code)
	})

	t.Run("COMMENT ON and CREATE INDEX resolve unqualified public synonyms with direct-object precedence", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE app.t (id NUMBER);
			CREATE TABLE app.other (id NUMBER);
			CREATE PUBLIC SYNONYM ps FOR t;
			CREATE PUBLIC SYNONYM other FOR t;
		`)
		require.NoError(t, err)

		// Unqualified public synonym resolves to local table T.
		results, err := c.Exec(`
			COMMENT ON TABLE ps IS 'via public synonym';
			COMMENT ON COLUMN ps.id IS 'col via public synonym';
			CREATE INDEX ps_idx ON ps (id);
		`, nil)
		require.NoError(t, err)
		for _, r := range results {
			require.Nil(t, r.Error)
		}
		tblT := tableByName(c.Metadata().Schemas[0])["T"]
		require.Equal(t, "via public synonym", tblT.Comment)
		require.Equal(t, "col via public synonym", tblT.Columns[0].Comment)
		require.Len(t, tblT.Indexes, 1)

		// Direct local object OTHER takes precedence over public synonym OTHER -> T.
		results, err = c.Exec(`COMMENT ON TABLE other IS 'direct other'`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Equal(t, "direct other", tableByName(c.Metadata().Schemas[0])["OTHER"].Comment)
		require.Equal(t, "via public synonym", tableByName(c.Metadata().Schemas[0])["T"].Comment)

		// Schema-qualified reference does NOT resolve public synonyms.
		results, err = c.Exec(`COMMENT ON TABLE app.ps IS 'should fail'`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "00942", oraErr.Code)
	})

	t.Run("CREATE VIEW accepts pseudocolumns and parenthesized column projections", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE t (id NUMBER)`)
		require.NoError(t, err)

		results, err := c.Exec(`
			CREATE VIEW v_paren AS SELECT (id) FROM t;
			CREATE VIEW v_pseudo AS SELECT SYSDATE FROM dual;
		`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)
	})

	t.Run("ALTER TABLE ADD CONSTRAINT resolves dotted object-attribute paths against root column", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE p (id NUMBER CONSTRAINT p_pk PRIMARY KEY);
			CREATE TABLE t (o obj_t);
		`)
		require.NoError(t, err)

		results, err := c.Exec(`ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (o.x) REFERENCES p (id)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		results, err = c.Exec(`ALTER TABLE t ADD CONSTRAINT fk_bad FOREIGN KEY (missing.x) REFERENCES p (id)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "00904", oraErr.Code)
	})

	t.Run("CREATE TABLE IF NOT EXISTS no-ops on existing table before unsupported-form checks", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE t (id NUMBER)`)
		require.NoError(t, err)

		results, err := c.Exec(`
			CREATE TABLE IF NOT EXISTS t OF obj_t;
			CREATE TABLE IF NOT EXISTS t AS SELECT id FROM t;
			CREATE GLOBAL TEMPORARY TABLE IF NOT EXISTS t (id NUMBER);
		`, nil)
		require.NoError(t, err)
		for _, r := range results {
			require.Nil(t, r.Error)
		}
	})

	t.Run("CREATE INDEX resolves object-attribute paths against root column", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE tt (o ty, o2 ty)`)
		require.NoError(t, err)

		results, err := c.Exec(`CREATE INDEX tt_ix ON tt (o.x, o2.x)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Len(t, tableByName(c.Metadata().Schemas[0])["TT"].Indexes, 1)

		results, err = c.Exec(`CREATE INDEX tt_bad_ix ON tt (missing.x)`, nil)
		require.NoError(t, err)
		var oraErr *Error
		require.True(t, errors.As(results[0].Error, &oraErr))
		require.Equal(t, "00904", oraErr.Code)
	})

	t.Run("ALTER TABLE preserves foreign keys referencing external schemas", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE app.c (p_id NUMBER)`)
		require.NoError(t, err)

		results, err := c.Exec(`ALTER TABLE app.c ADD CONSTRAINT c_fk FOREIGN KEY (p_id) REFERENCES other.p (id)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		fks := tableByName(c.Metadata().Schemas[0])["C"].ForeignKeys
		require.Len(t, fks, 1)
		require.Equal(t, "OTHER", fks[0].ReferencedSchema)
		require.Equal(t, "P", fks[0].ReferencedTable)
	})
}

func TestServerConfigAndOracle23SimulationFidelity(t *testing.T) {
	requireCode := func(t *testing.T, err error, want string) {
		t.Helper()
		require.Error(t, err)
		var oraErr *Error
		require.True(t, errors.As(err, &oraErr), "expected *catalog.Error, got %T: %v", err, err)
		require.Equal(t, want, oraErr.Code)
	}

	t.Run("MAX_STRING_SIZE=STANDARD refuses column COLLATE and view DEFAULT COLLATION with ORA-43929", func(t *testing.T) {
		c := New()
		c.SetServerConfig(StandardOracle23ServerConfig())
		for _, stmt := range []string{
			`ALTER TABLE missing ADD (a NUMBER COLLATE binary_ci)`,
			`ALTER TABLE missing ADD (a COLLATE binary_ci AS (1))`,
			`ALTER TABLE missing MODIFY (a COLLATE binary_ci)`,
			`ALTER TABLE missing MODIFY (a COLLATE binary_ci AS (1))`,
			`CREATE TABLE t (a VARCHAR2(10) COLLATE binary_ci)`,
			`CREATE TABLE t (a VARCHAR2(10)) DEFAULT COLLATION binary_ci`,
			`CREATE VIEW v DEFAULT COLLATION binary_ci AS SELECT 1 FROM dual`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			require.Len(t, results, 1)
			requireCode(t, results[0].Error, "43929")
		}

		// Switching to EXTENDED allows collation clauses again.
		cfg := c.ServerConfig()
		cfg.MaxStringSize = "EXTENDED"
		c.SetServerConfig(cfg)
		results, err := c.Exec(`CREATE TABLE t (a VARCHAR2(10) COLLATE binary_ci); CREATE VIEW v DEFAULT COLLATION binary_ci AS SELECT a FROM t`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)
	})

	t.Run("sequence SHARD and NOSHARD fail with ORA-02511 when DisableShardDDL is enabled, while schema dumps stay permissive", func(t *testing.T) {
		// Schema dumps (GetDatabaseMetadata / default LoadSQL) stay permissive for sharded DDL.
		meta, err := GetDatabaseMetadata(`CREATE SEQUENCE s_dump SHARD; ALTER SEQUENCE s_dump NOSHARD;`)
		require.NoError(t, err)
		require.Len(t, meta.Schemas[0].Sequences, 1)

		c := New()
		c.SetServerConfig(StandardOracle23ServerConfig())
		for _, stmt := range []string{
			`CREATE SEQUENCE s SHARD`,
			`ALTER SEQUENCE missing SHARD`,
			`ALTER SEQUENCE missing NOSHARD`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			requireCode(t, results[0].Error, "02511")
		}

		// ALTER SESSION ENABLE SHARD DDL unlocks sequence shard clauses.
		results, err := c.Exec(`ALTER SESSION ENABLE SHARD DDL; CREATE SEQUENCE s SHARD; ALTER SEQUENCE s NOSHARD`, nil)
		require.NoError(t, err)
		for _, r := range results {
			require.Nil(t, r.Error)
		}

		results, err = c.Exec(`ALTER SESSION DISABLE SHARD DDL; ALTER SEQUENCE s SHARD`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		requireCode(t, results[1].Error, "02511")
	})

	t.Run("SHARING = METADATA and DATA fail with ORA-65021 outside Application Root, while schema dumps stay permissive", func(t *testing.T) {
		// Schema dumps (GetDatabaseMetadata / default LoadSQL) stay permissive for application-common DDL.
		meta, err := GetDatabaseMetadata(`
			CREATE SEQUENCE s_app SHARING = METADATA;
			CREATE SYNONYM sy_app SHARING = METADATA FOR s_app;
		`)
		require.NoError(t, err)
		require.Len(t, meta.Schemas[0].Sequences, 1)
		require.Len(t, meta.Schemas[0].Synonyms, 1)

		c := New()
		c.SetServerConfig(StandardOracle23ServerConfig())
		for _, stmt := range []string{
			`CREATE SEQUENCE s SHARING = METADATA`,
			`CREATE SEQUENCE s SHARING = DATA`,
			`CREATE SYNONYM s SHARING = METADATA FOR s`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			requireCode(t, results[0].Error, "65021")
		}

		results, err := c.Exec(`CREATE SEQUENCE s SHARING = NONE; CREATE SYNONYM sy SHARING = NONE FOR s`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)

		// Inside an application root, SHARING = METADATA is accepted.
		cfg := c.ServerConfig()
		cfg.NonApplicationRoot = false
		c.SetServerConfig(cfg)
		results, err = c.Exec(`CREATE SEQUENCE s2 SHARING = METADATA; CREATE SYNONYM sy2 SHARING = METADATA FOR s2`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)
	})

	t.Run("TDE keystore closed fails column ENCRYPT with ORA-28365 after table, column, and FK checks", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE p (id NUMBER CONSTRAINT p_pk PRIMARY KEY);
			CREATE TABLE t (id NUMBER, p_id NUMBER CONSTRAINT t_fk REFERENCES p (id));
		`)
		require.NoError(t, err)
		c.SetServerConfig(StandardOracle23ServerConfig())

		// Missing table -> ORA-00942 before wallet check.
		results, err := c.Exec(`ALTER TABLE missing ADD (x NUMBER ENCRYPT)`, nil)
		require.NoError(t, err)
		requireCode(t, results[0].Error, "00942")

		// Duplicate column -> ORA-01430 before wallet check.
		results, err = c.Exec(`ALTER TABLE t ADD (id NUMBER ENCRYPT)`, nil)
		require.NoError(t, err)
		requireCode(t, results[0].Error, "01430")

		// FK column -> ORA-28335 before wallet check.
		results, err = c.Exec(`ALTER TABLE t MODIFY (p_id ENCRYPT)`, nil)
		require.NoError(t, err)
		requireCode(t, results[0].Error, "28335")

		// Non-FK column ADD and MODIFY ENCRYPT -> ORA-28365.
		for _, stmt := range []string{
			`ALTER TABLE t ADD (secret NUMBER ENCRYPT)`,
			`ALTER TABLE t MODIFY (id ENCRYPT)`,
			`ALTER TABLE t MODIFY (id ENCRYPT SALT)`,
			`ALTER TABLE t MODIFY (id ENCRYPT NO SALT)`,
			`ALTER TABLE t MODIFY (id ENCRYPT IDENTIFIED BY secret_pw)`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			requireCode(t, results[0].Error, "28365")
		}
	})

	t.Run("ValidateExternalSchemaRefs rejects cross-schema REFERENCES with ORA-00942", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE app.t (id NUMBER)`)
		require.NoError(t, err)
		c.SetServerConfig(StandardOracle23ServerConfig())

		for _, stmt := range []string{
			`CREATE TABLE child1 (id REFERENCES other.p)`,
			`CREATE TABLE child2 (id NUMBER REFERENCES other.p)`,
			`CREATE TABLE child3 (id NUMBER, FOREIGN KEY (id) REFERENCES other.p)`,
			`ALTER TABLE t ADD (p_id NUMBER REFERENCES other.p)`,
			`ALTER TABLE t ADD FOREIGN KEY (id) REFERENCES other.p`,
			`ALTER TABLE t MODIFY (id REFERENCES other.p)`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			requireCode(t, results[0].Error, "00942")
		}
	})

	t.Run("untyped virtual column infers NUMBER, BINARY_FLOAT, and BINARY_DOUBLE with numeric precedence while keeping functions untyped", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE t (c VARCHAR2(10));
			ALTER TABLE t ADD (v_lit AS (1));
			ALTER TABLE t ADD (v_neg AS (-1));
			ALTER TABLE t ADD (v_arith AS (1 + 2 * 3));
			ALTER TABLE t ADD (v_float AS (1f));
			ALTER TABLE t ADD (v_double AS (1d));
			ALTER TABLE t ADD (v_mix_f AS (1 + 2.5F));
			ALTER TABLE t ADD (v_mix_d AS (1f + 2.5D));
			ALTER TABLE t ADD (v_fn AS (UPPER(c)));
		`)
		require.NoError(t, err)
		tbl := tableByName(c.Metadata().Schemas[0])["T"]
		require.Equal(t, "NUMBER", tableFindColumn(tbl, "V_LIT").Type)
		require.Equal(t, "NUMBER", tableFindColumn(tbl, "V_NEG").Type)
		require.Equal(t, "NUMBER", tableFindColumn(tbl, "V_ARITH").Type)
		require.Equal(t, "BINARY_FLOAT", tableFindColumn(tbl, "V_FLOAT").Type)
		require.Equal(t, "BINARY_DOUBLE", tableFindColumn(tbl, "V_DOUBLE").Type)
		require.Equal(t, "BINARY_FLOAT", tableFindColumn(tbl, "V_MIX_F").Type)
		require.Equal(t, "BINARY_DOUBLE", tableFindColumn(tbl, "V_MIX_D").Type)
		require.Equal(t, "", tableFindColumn(tbl, "V_FN").Type)

		// Modifying a virtual column's expression refreshes or clears its inferred type.
		results, err := c.Exec(`
			ALTER TABLE t MODIFY (v_lit AS (UPPER(c)));
			ALTER TABLE t MODIFY (v_fn AS (1.5f));
		`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)
		tbl = tableByName(c.Metadata().Schemas[0])["T"]
		require.Equal(t, "", tableFindColumn(tbl, "V_LIT").Type)
		require.Equal(t, "BINARY_FLOAT", tableFindColumn(tbl, "V_FN").Type)
		def, err := DatabaseDefinition(c.Metadata())
		require.NoError(t, err)
		require.NotContains(t, def, `"V_LIT" NUMBER`)
	})

	t.Run("virtual column edition clauses validate against catalog editions when ValidateEditions is enabled, while schema dumps stay permissive", func(t *testing.T) {
		// Schema dumps (GetDatabaseMetadata / default LoadSQL) stay permissive about pre-existing database editions.
		meta, err := GetDatabaseMetadata(`CREATE TABLE t (id NUMBER, v AS (1) EVALUATE USING EDITION app_v2)`)
		require.NoError(t, err)
		require.Len(t, meta.Schemas[0].Tables, 1)

		c, err := LoadSQL(`CREATE TABLE t (id NUMBER)`)
		require.NoError(t, err)
		c.SetServerConfig(StandardOracle23ServerConfig())

		for _, stmt := range []string{
			`ALTER TABLE t ADD (v1 AS (1) EVALUATE USING EDITION missing_ed)`,
			`ALTER TABLE t ADD (v2 AS (1) UNUSABLE BEFORE EDITION missing_ed)`,
			`ALTER TABLE t ADD (v3 AS (1) UNUSABLE BEGINNING WITH EDITION missing_ed)`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			requireCode(t, results[0].Error, "38802")
		}

		// Built-in ORA$BASE, CURRENT EDITION, NULL EDITION, and created editions succeed.
		results, err := c.Exec(`
			CREATE EDITION v2;
			ALTER TABLE t ADD (v_base AS (1) EVALUATE USING EDITION ora$base);
			ALTER TABLE t ADD (v_cur AS (1) EVALUATE USING CURRENT EDITION);
			ALTER TABLE t ADD (v_null AS (1) EVALUATE USING NULL EDITION);
			ALTER TABLE t ADD (v_v2 AS (1) UNUSABLE BEFORE EDITION v2);
			DROP EDITION v2;
		`, nil)
		require.NoError(t, err)
		for _, r := range results {
			require.Nil(t, r.Error)
		}

		// Dropping v2 revokes it.
		results, err = c.Exec(`ALTER TABLE t ADD (v_gone AS (1) EVALUATE USING EDITION v2)`, nil)
		require.NoError(t, err)
		requireCode(t, results[0].Error, "38802")
	})

	t.Run("PERIOD FOR synthesizes temporal validity check constraints and renders executable PERIOD FOR DDL", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE p1 (PERIOD FOR valid_time);
			CREATE TABLE p2 (id NUMBER, s_col TIMESTAMP, e_col TIMESTAMP);
			ALTER TABLE p2 ADD (PERIOD FOR vt);
			ALTER TABLE p2 ADD (PERIOD FOR vt2 (s_col, e_col));
		`)
		require.NoError(t, err)

		t1 := tableByName(c.Metadata().Schemas[0])["P1"]
		require.Len(t, t1.CheckConstraints, 1)
		require.Equal(t, `"VALID_TIME_START" < "VALID_TIME_END"`, t1.CheckConstraints[0].Expression)

		t2 := tableByName(c.Metadata().Schemas[0])["P2"]
		require.Len(t, t2.CheckConstraints, 2)
		require.Equal(t, `"VT_START" < "VT_END"`, t2.CheckConstraints[0].Expression)
		require.Equal(t, `("S_COL" < "E_COL") and ("VT2" > 0)`, t2.CheckConstraints[1].Expression)

		// DatabaseDefinition must emit PERIOD FOR clauses rather than standalone CHECKs
		// referencing hidden generated columns, and must round-trip cleanly.
		def, err := DatabaseDefinition(c.Metadata())
		require.NoError(t, err)
		require.Contains(t, def, `PERIOD FOR "VALID_TIME"`)
		require.Contains(t, def, `PERIOD FOR "VT"`)
		require.Contains(t, def, `PERIOD FOR "VT2" ("S_COL", "E_COL")`)
		require.NotContains(t, def, `VALID_TIME_START`)
		require.NotContains(t, def, `VT_START`)

		roundTrip, err := LoadSQL(def)
		require.NoError(t, err)
		require.True(t, Diff(c, roundTrip).IsEmpty())

		// GenerateMigration adding and dropping a PERIOD FOR emits ALTER TABLE ADD/DROP (PERIOD FOR ...).
		base, err := LoadSQL(`CREATE TABLE p2 (id NUMBER, s_col TIMESTAMP, e_col TIMESTAMP);`)
		require.NoError(t, err)
		withPeriod, err := LoadSQL(`CREATE TABLE p2 (id NUMBER, s_col TIMESTAMP, e_col TIMESTAMP, PERIOD FOR vt);`)
		require.NoError(t, err)
		addSQL := GenerateMigration(base, withPeriod, nil).SQL()
		require.Contains(t, addSQL, `ALTER TABLE "P2" ADD (PERIOD FOR "VT");`)
		require.NotContains(t, addSQL, `VT_START`)
		dropSQL := GenerateMigration(withPeriod, base, nil).SQL()
		require.Contains(t, dropSQL, `ALTER TABLE "P2" DROP (PERIOD FOR "VT");`)
		// Replaying dropSQL removes the period cleanly.
		results, err := withPeriod.Exec(dropSQL, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.True(t, Diff(withPeriod, base).IsEmpty())

		// Synced temporal tables carry Oracle-generated names like Ad6b29a; diffing
		// against DatabaseDefinition (which re-creates CHK_P2_VT) must converge with no diff.
		synced := &Catalog{meta: &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{{
				Tables: []*metadatapb.TableMetadata{{
					Name: "P2",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "ID", Type: "NUMBER", Nullable: true, Position: 1},
						{Name: "S_COL", Type: "TIMESTAMP(6)", Nullable: true, Position: 2},
						{Name: "E_COL", Type: "TIMESTAMP(6)", Nullable: true, Position: 3},
					},
					CheckConstraints: []*metadatapb.CheckConstraintMetadata{{
						Name:       "VTd6b29a",
						Expression: `"VT_START" < "VT_END"`,
					}},
				}},
			}},
		}}
		syncedDef, err := DatabaseDefinition(synced.Metadata())
		require.NoError(t, err)
		reloaded, err := LoadSQL(syncedDef)
		require.NoError(t, err)
		require.True(t, Diff(synced, reloaded).IsEmpty())
		require.True(t, Diff(reloaded, synced).IsEmpty())
	})

	t.Run("SET UNUSED executes when ApplySetUnused is enabled and enforces column/FK/last-column checks", func(t *testing.T) {
		c, err := LoadSQL(`
			CREATE TABLE p (b NUMBER CONSTRAINT p_pk PRIMARY KEY);
			CREATE TABLE t (b NUMBER CONSTRAINT t_pk PRIMARY KEY, a NUMBER CONSTRAINT t_fk REFERENCES p (b), d NUMBER);
			CREATE TABLE ch (id NUMBER, t_b NUMBER CONSTRAINT ch_fk REFERENCES t (b));
		`)
		require.NoError(t, err)
		c.SetServerConfig(StandardOracle23ServerConfig())

		// Referencing FK column can be set unused, removing column and FK.
		results, err := c.Exec(`ALTER TABLE t SET UNUSED COLUMN a ONLINE`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		tblT := tableByName(c.Metadata().Schemas[0])["T"]
		require.Nil(t, tableFindColumn(tblT, "A"))
		require.Empty(t, tblT.ForeignKeys)

		// Parent key referenced by incoming FK without CASCADE -> ORA-12992.
		results, err = c.Exec(`ALTER TABLE t SET UNUSED COLUMN b`, nil)
		require.NoError(t, err)
		requireCode(t, results[0].Error, "12992")

		// Setting all remaining columns unused -> ORA-12983.
		results, err = c.Exec(`ALTER TABLE t SET UNUSED (b, d) CASCADE CONSTRAINTS`, nil)
		require.NoError(t, err)
		requireCode(t, results[0].Error, "12983")

		// Duplicate column in SET UNUSED list -> ORA-00957.
		results, err = c.Exec(`ALTER TABLE t SET UNUSED (d, d)`, nil)
		require.NoError(t, err)
		requireCode(t, results[0].Error, "00957")

		// Parent key with CASCADE CONSTRAINTS succeeds and drops incoming FK.
		results, err = c.Exec(`ALTER TABLE t SET UNUSED COLUMN b CASCADE CONSTRAINTS`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Empty(t, tableByName(c.Metadata().Schemas[0])["CH"].ForeignKeys)

		// Multi-column CHECK, UNIQUE, and FK constraints require CASCADE CONSTRAINTS (ORA-12991)
		// when another referenced column remains in the table.
		mc, err := LoadSQL(`
			CREATE TABLE mp (x NUMBER, y NUMBER, CONSTRAINT mp_pk PRIMARY KEY (x, y));
			CREATE TABLE mt (
				id NUMBER,
				a NUMBER,
				b NUMBER,
				c NUMBER,
				CONSTRAINT mt_ck CHECK (a < b),
				CONSTRAINT mt_uk UNIQUE (b, c),
				CONSTRAINT mt_fk FOREIGN KEY (a, c) REFERENCES mp (x, y)
			);
		`)
		require.NoError(t, err)
		mc.SetServerConfig(StandardOracle23ServerConfig())

		for _, stmt := range []string{
			`ALTER TABLE mt SET UNUSED COLUMN a`,
			`ALTER TABLE mt SET UNUSED COLUMN b`,
			`ALTER TABLE mt SET UNUSED COLUMN c`,
		} {
			results, err := mc.Exec(stmt, nil)
			require.NoError(t, err)
			requireCode(t, results[0].Error, "12991")
		}
		// Table remains unchanged after rejected SET UNUSED statements.
		tblMT := tableByName(mc.Metadata().Schemas[0])["MT"]
		require.Len(t, tblMT.Columns, 4)
		require.Len(t, tblMT.CheckConstraints, 1)

		// Setting all columns of a multi-column check unused together succeeds when no
		// remaining column is referenced, or with CASCADE CONSTRAINTS.
		results, err = mc.Exec(`ALTER TABLE mt SET UNUSED COLUMN a CASCADE CONSTRAINTS`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		tblMT = tableByName(mc.Metadata().Schemas[0])["MT"]
		require.Nil(t, tableFindColumn(tblMT, "A"))
		require.Empty(t, tblMT.CheckConstraints)
		require.Empty(t, tblMT.ForeignKeys)

		// Dropping both B and C together (all columns of mt_uk) succeeds without CASCADE.
		results, err = mc.Exec(`ALTER TABLE mt SET UNUSED (b, c)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		tblMT = tableByName(mc.Metadata().Schemas[0])["MT"]
		require.Len(t, tblMT.Columns, 1)
		require.Empty(t, tblMT.Indexes)

		// Setting a column unused when referenced by a surviving virtual column -> ORA-54031
		// (even with CASCADE CONSTRAINTS), and leaving only virtual columns -> ORA-54037.
		vc, err := LoadSQL(`CREATE TABLE vt (a NUMBER, b NUMBER, v AS (a + 1))`)
		require.NoError(t, err)
		vc.SetServerConfig(StandardOracle23ServerConfig())

		for _, stmt := range []string{
			`ALTER TABLE vt SET UNUSED COLUMN a`,
			`ALTER TABLE vt SET UNUSED COLUMN a CASCADE CONSTRAINTS`,
		} {
			results, err := vc.Exec(stmt, nil)
			require.NoError(t, err)
			requireCode(t, results[0].Error, "54031")
		}
		require.Len(t, tableByName(vc.Metadata().Schemas[0])["VT"].Columns, 3)

		results, err = vc.Exec(`ALTER TABLE vt SET UNUSED (a, b)`, nil)
		require.NoError(t, err)
		requireCode(t, results[0].Error, "54037")

		// Setting both A and its dependent virtual column V unused together succeeds.
		results, err = vc.Exec(`ALTER TABLE vt SET UNUSED (a, v)`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		tblVT := tableByName(vc.Metadata().Schemas[0])["VT"]
		require.Len(t, tblVT.Columns, 1)
		require.Equal(t, "B", tblVT.Columns[0].Name)
	})

	t.Run("ALTER SEQUENCE enforces ORA-04007 and ORA-04009 against initial sequence current value", func(t *testing.T) {
		c, err := LoadSQL(`CREATE SEQUENCE s START WITH 1 INCREMENT BY 1`)
		require.NoError(t, err)

		for _, stmt := range []string{
			`ALTER SEQUENCE s MINVALUE 1`,
			`ALTER SEQUENCE s NOMINVALUE`,
		} {
			results, err := c.Exec(stmt, nil)
			require.NoError(t, err)
			requireCode(t, results[0].Error, "04007")
		}

		// MINVALUE 0 is <= initial current value (0), so it succeeds.
		results, err := c.Exec(`ALTER SEQUENCE s MINVALUE 0`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)

		// RESTART START WITH 10 moves current value to 9, allowing MINVALUE 5 and rejecting MAXVALUE 8 (ORA-04009).
		results, err = c.Exec(`ALTER SEQUENCE s RESTART START WITH 10; ALTER SEQUENCE s MINVALUE 5`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		require.Nil(t, results[1].Error)

		results, err = c.Exec(`ALTER SEQUENCE s MAXVALUE 8`, nil)
		require.NoError(t, err)
		requireCode(t, results[0].Error, "04009")

		// Descending sequence NOMAXVALUE (-1 < initial currval 0) -> ORA-04009.
		results, err = c.Exec(`CREATE SEQUENCE sd START WITH -1 INCREMENT BY -1; ALTER SEQUENCE sd NOMAXVALUE`, nil)
		require.NoError(t, err)
		require.Nil(t, results[0].Error)
		requireCode(t, results[1].Error, "04009")
	})
}
