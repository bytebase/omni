package catalog

import (
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

// TestDefinitionFidelity covers schema features that moved out of
// TestDefinitionGaps: packages, synonyms and public synonyms, view column
// aliases (including invisible columns), identity sequence options, explicit
// column collations, table partitioning, and materialized view refresh policies.
func TestDefinitionFidelity(t *testing.T) {
	schemaWith := func(schema *metadatapb.SchemaMetadata) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Name:    "APP",
			Schemas: []*metadatapb.SchemaMetadata{schema},
		}
	}

	t.Run("package spec and body definition, round trip, and migration", func(t *testing.T) {
		// Synced spec-only package (no CREATE prefix)
		def, err := DatabaseDefinition(schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Packages: []*metadatapb.PackageMetadata{{
				Name:       "PKG",
				Definition: "PACKAGE PKG AS PROCEDURE P; END PKG;",
			}},
		}))
		require.NoError(t, err)
		require.Equal(t, "CREATE OR REPLACE PACKAGE PKG AS PROCEDURE P; END PKG;\n\n", def)

		// Spec + body loaded from SQL round-trips cleanly
		src := `
CREATE OR REPLACE PACKAGE PKG AS
  PROCEDURE P;
END PKG;
/
CREATE OR REPLACE PACKAGE BODY PKG AS
  PROCEDURE P IS
  BEGIN
    NULL;
  END P;
END PKG;
/
`
		c1, err := LoadSQL(src)
		require.NoError(t, err)
		def1, err := DatabaseDefinition(c1.Metadata())
		require.NoError(t, err)
		require.Contains(t, def1, "CREATE OR REPLACE PACKAGE PKG")
		require.Contains(t, def1, "CREATE OR REPLACE PACKAGE BODY PKG")

		c2, err := LoadSQL(def1)
		require.NoError(t, err)
		def2, err := DatabaseDefinition(c2.Metadata())
		require.NoError(t, err)
		require.Equal(t, def1, def2)

		// Migration: create, alter (upgrading CREATE to CREATE OR REPLACE), drop
		empty := New()
		createPlan := GenerateMigration(empty, c1, nil)
		require.Contains(t, createPlan.SQL(), "CREATE OR REPLACE PACKAGE PKG")
		require.Contains(t, createPlan.SQL(), "CREATE OR REPLACE PACKAGE BODY PKG")

		cAltered, err := LoadSQL(`
CREATE PACKAGE PKG AS
  PROCEDURE P;
  PROCEDURE Q;
END PKG;
/
CREATE PACKAGE BODY PKG AS
  PROCEDURE P IS BEGIN NULL; END P;
  PROCEDURE Q IS BEGIN NULL; END Q;
END PKG;
/
`)
		require.NoError(t, err)
		alterPlan := GenerateMigration(c1, cAltered, nil)
		require.Contains(t, alterPlan.SQL(), "CREATE OR REPLACE PACKAGE PKG")
		require.Contains(t, alterPlan.SQL(), "CREATE OR REPLACE PACKAGE BODY PKG")
		_, err = c1.Exec(alterPlan.SQL(), nil)
		require.NoError(t, err)
		require.Len(t, c1.Metadata().Schemas[0].Packages, 1)
		require.Equal(t, 1, strings.Count(c1.Metadata().Schemas[0].Packages[0].Definition, "PACKAGE BODY"))

		dropPlan := GenerateMigration(c1, empty, nil)
		require.Contains(t, dropPlan.SQL(), `DROP PACKAGE "PKG";`)
		_, err = c1.Exec(dropPlan.SQL(), nil)
		require.NoError(t, err)
		require.Empty(t, c1.Metadata().Schemas[0].Packages)
	})

	t.Run("synonym and public synonym definition, filtering, round trip, and migration", func(t *testing.T) {
		meta := schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Synonyms: []*metadatapb.SynonymMetadata{
				{Name: "S_LOCAL", TargetSchema: "APP", TargetName: "T"},
				{Name: "S_REMOTE", TargetSchema: "OTHER", TargetName: "RT", TargetDbLink: "REMOTE_LINK"},
			},
		})
		meta.PublicSynonyms = []*metadatapb.SynonymMetadata{
			{Name: "DUAL", TargetSchema: "SYS", TargetName: "DUAL"},
			{Name: "OTHER_PUB", TargetSchema: "OTHER", TargetName: "OT"},
			{Name: "PUB_S", TargetSchema: "APP", TargetName: "T"},
		}
		def, err := DatabaseDefinition(meta)
		require.NoError(t, err)
		require.Contains(t, def, `CREATE OR REPLACE SYNONYM "S_LOCAL" FOR "T";`)
		require.Contains(t, def, `CREATE OR REPLACE SYNONYM "S_REMOTE" FOR "OTHER"."RT"@REMOTE_LINK;`)
		require.Contains(t, def, `CREATE OR REPLACE PUBLIC SYNONYM "PUB_S" FOR "T";`)
		require.NotContains(t, def, `"DUAL"`)
		require.NotContains(t, def, `"OTHER_PUB"`)

		// Round-trip via LoadSQL
		c, err := LoadSQL(def)
		require.NoError(t, err)
		def2, err := DatabaseDefinition(c.Metadata())
		require.NoError(t, err)
		require.Equal(t, def, def2)

		// Migration: create and drop
		empty := New()
		createPlan := GenerateMigration(empty, c, nil)
		require.Contains(t, createPlan.SQL(), `CREATE OR REPLACE SYNONYM "S_LOCAL" FOR "T";`)
		require.Contains(t, createPlan.SQL(), `CREATE OR REPLACE PUBLIC SYNONYM "PUB_S" FOR "T";`)

		dropPlan := GenerateMigration(c, empty, nil)
		require.Contains(t, dropPlan.SQL(), `DROP SYNONYM "S_LOCAL";`)
		require.Contains(t, dropPlan.SQL(), `DROP PUBLIC SYNONYM "PUB_S";`)
		_, err = c.Exec(dropPlan.SQL(), nil)
		require.NoError(t, err)
		require.Empty(t, c.Metadata().Schemas[0].Synonyms)
		require.Empty(t, c.Metadata().PublicSynonyms)
	})

	t.Run("view column aliases and invisible columns", func(t *testing.T) {
		// Renamed view columns emit an explicit alias list
		def, err := DatabaseDefinition(schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "V",
				Definition: "SELECT ID, NAME FROM T",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "RENAMED_ID", Position: 1, Nullable: true},
					{Name: "RENAMED_NAME", Position: 2, Nullable: true},
				},
			}},
		}))
		require.NoError(t, err)
		require.Equal(t, "CREATE VIEW \"V\" (\"RENAMED_ID\", \"RENAMED_NAME\") AS SELECT ID, NAME FROM T;\n\n", def)

		// Matching column names with no invisible columns omit redundant aliases
		plainDef, err := DatabaseDefinition(schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "PLAIN_V",
				Definition: "SELECT ID, NAME FROM T",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "ID", Position: 1, Nullable: true},
					{Name: "NAME", Position: 2, Nullable: true},
				},
			}},
		}))
		require.NoError(t, err)
		require.Equal(t, "CREATE VIEW \"PLAIN_V\" AS SELECT ID, NAME FROM T;\n\n", plainDef)

		// View with invisible column renders INVISIBLE and round-trips through LoadSQL
		c, err := LoadSQL(`CREATE VIEW V_INV (A INVISIBLE, B) AS SELECT ID, NAME FROM T;`)
		require.NoError(t, err)
		invDef, err := DatabaseDefinition(c.Metadata())
		require.NoError(t, err)
		require.Equal(t, "CREATE VIEW \"V_INV\" (\"A\" INVISIBLE, \"B\") AS SELECT ID, NAME FROM T;\n\n", invDef)

		c2, err := LoadSQL(invDef)
		require.NoError(t, err)
		invDef2, err := DatabaseDefinition(c2.Metadata())
		require.NoError(t, err)
		require.Equal(t, invDef, invDef2)
	})

	t.Run("identity sequence options definition, round trip, and migration", func(t *testing.T) {
		c, err := LoadSQL(`CREATE TABLE i (id NUMBER GENERATED ALWAYS AS IDENTITY (START WITH 5 INCREMENT BY 2 MAXVALUE 100000 MINVALUE 5 CACHE 100 CYCLE ORDER))`)
		require.NoError(t, err)
		def, err := DatabaseDefinition(c.Metadata())
		require.NoError(t, err)
		require.Contains(t, def, `GENERATED ALWAYS AS IDENTITY (START WITH 5 INCREMENT BY 2 MAXVALUE 100000 MINVALUE 5 CACHE 100 CYCLE ORDER)`)

		c2, err := LoadSQL(def)
		require.NoError(t, err)
		def2, err := DatabaseDefinition(c2.Metadata())
		require.NoError(t, err)
		require.Equal(t, def, def2)

		// NOCACHE round-trips as NOCACHE
		cNoCache, err := LoadSQL(`CREATE TABLE i2 (id NUMBER GENERATED BY DEFAULT ON NULL AS IDENTITY (START WITH 10 INCREMENT BY 1 NOCACHE))`)
		require.NoError(t, err)
		defNoCache, err := DatabaseDefinition(cNoCache.Metadata())
		require.NoError(t, err)
		require.Contains(t, defNoCache, `GENERATED BY DEFAULT ON NULL AS IDENTITY (START WITH 10 INCREMENT BY 1 NOCACHE)`)

		// Altering identity sequence options generates MODIFY
		cChanged, err := LoadSQL(`CREATE TABLE i (id NUMBER GENERATED ALWAYS AS IDENTITY (START WITH 5 INCREMENT BY 2 MAXVALUE 500000 MINVALUE 5 NOCACHE NOCYCLE ORDER))`)
		require.NoError(t, err)
		plan := GenerateMigration(c, cChanged, nil)
		require.Contains(t, plan.SQL(), `MAXVALUE 500000`)
		require.Contains(t, plan.SQL(), `NOCACHE`)
		require.Contains(t, plan.SQL(), `NOCYCLE`)
	})

	t.Run("explicit column collation emitted and USING_NLS_COMP omitted", func(t *testing.T) {
		def, err := DatabaseDefinition(schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Tables: []*metadatapb.TableMetadata{{
				Name: "T",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "EXPLICIT_COL", Type: "VARCHAR2(20 BYTE)", Collation: "BINARY_CI", Nullable: true},
					{Name: "DEFAULT_COL", Type: "VARCHAR2(20 BYTE)", Collation: "USING_NLS_COMP", Nullable: true},
				},
			}},
		}))
		require.NoError(t, err)
		require.Contains(t, def, `"EXPLICIT_COL" VARCHAR2(20 BYTE) COLLATE BINARY_CI`)
		require.NotContains(t, def, "USING_NLS_COMP")

		c, err := LoadSQL(def)
		require.NoError(t, err)
		def2, err := DatabaseDefinition(c.Metadata())
		require.NoError(t, err)
		require.Equal(t, def, def2)
	})

	t.Run("table partitioning RANGE, LIST, HASH, INTERVAL, and composite", func(t *testing.T) {
		cases := []struct {
			name string
			sql  string
			want []string
		}{
			{
				name: "range",
				sql:  `CREATE TABLE p (a NUMBER) PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS THAN (10), PARTITION p2 VALUES LESS THAN (MAXVALUE))`,
				want: []string{`PARTITION BY RANGE ("A")`, `PARTITION "P1" VALUES LESS THAN (10)`, `PARTITION "P2" VALUES LESS THAN (MAXVALUE)`},
			},
			{
				name: "range interval",
				sql:  `CREATE TABLE pi (d DATE) PARTITION BY RANGE (d) INTERVAL (NUMTOYMINTERVAL(1, 'MONTH')) (PARTITION p0 VALUES LESS THAN (DATE '2025-01-01'))`,
				want: []string{`PARTITION BY RANGE ("D") INTERVAL (NUMTOYMINTERVAL(1, 'MONTH'))`, `PARTITION "P0" VALUES LESS THAN (DATE '2025-01-01')`},
			},
			{
				name: "list",
				sql:  `CREATE TABLE pl (region VARCHAR2(10)) PARTITION BY LIST (region) (PARTITION p_us VALUES ('US', 'CA'), PARTITION p_other VALUES (DEFAULT))`,
				want: []string{`PARTITION BY LIST ("REGION")`, `PARTITION "P_US" VALUES ('US', 'CA')`, `PARTITION "P_OTHER" VALUES (DEFAULT)`},
			},
			{
				name: "hash",
				sql:  `CREATE TABLE ph (id NUMBER) PARTITION BY HASH (id) (PARTITION p1, PARTITION p2)`,
				want: []string{`PARTITION BY HASH ("ID")`, `PARTITION "P1"`, `PARTITION "P2"`},
			},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				c, err := LoadSQL(tc.sql)
				require.NoError(t, err)
				require.NotEmpty(t, c.Metadata().Schemas[0].Tables[0].Partitions)
				def, err := DatabaseDefinition(c.Metadata())
				require.NoError(t, err)
				for _, w := range tc.want {
					require.Contains(t, def, w)
				}
				c2, err := LoadSQL(def)
				require.NoError(t, err)
				def2, err := DatabaseDefinition(c2.Metadata())
				require.NoError(t, err)
				require.Equal(t, def, def2)
			})
		}

		// Composite RANGE-LIST subpartitions rendered from synced metadata
		def, err := DatabaseDefinition(schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Tables: []*metadatapb.TableMetadata{{
				Name: "COMP",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "Y", Type: "NUMBER", Nullable: true},
					{Name: "R", Type: "VARCHAR2(10 BYTE)", Nullable: true},
				},
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name:       "P2025",
					Type:       metadatapb.TablePartitionMetadata_RANGE,
					Expression: "Y",
					Value:      "2026",
					Subpartitions: []*metadatapb.TablePartitionMetadata{
						{Name: "P2025_US", Type: metadatapb.TablePartitionMetadata_LIST, Expression: "R", Value: "'US'"},
						{Name: "P2025_OTHER", Type: metadatapb.TablePartitionMetadata_LIST, Expression: "R", Value: "DEFAULT"},
					},
				}},
			}},
		}))
		require.NoError(t, err)
		require.True(t, strings.Contains(def, `PARTITION BY RANGE ("Y") SUBPARTITION BY LIST ("R")`), def)
		require.True(t, strings.Contains(def, `SUBPARTITION "P2025_US" VALUES ('US')`), def)
	})

	t.Run("materialized view refresh policy clauses", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			sql  string
			want string
		}{
			{
				name: "refresh complete on demand",
				sql:  `CREATE MATERIALIZED VIEW mv REFRESH COMPLETE ON DEMAND AS SELECT 1 AS c FROM DUAL`,
				want: `CREATE MATERIALIZED VIEW "MV" REFRESH COMPLETE ON DEMAND AS SELECT 1 AS c FROM DUAL`,
			},
			{
				name: "build deferred never refresh",
				sql:  `CREATE MATERIALIZED VIEW mv BUILD DEFERRED NEVER REFRESH AS SELECT 1 AS c FROM DUAL`,
				want: `CREATE MATERIALIZED VIEW "MV" BUILD DEFERRED NEVER REFRESH AS SELECT 1 AS c FROM DUAL`,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				c, err := LoadSQL(tc.sql)
				require.NoError(t, err)
				def, err := DatabaseDefinition(c.Metadata())
				require.NoError(t, err)
				require.Contains(t, def, tc.want)
			})
		}
	})
}
