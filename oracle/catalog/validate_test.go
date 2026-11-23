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
}
