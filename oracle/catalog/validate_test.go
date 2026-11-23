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
