package catalog

import (
	"context"
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
	loadMeta := func(t *testing.T, meta *metadatapb.DatabaseSchemaMetadata) *Catalog {
		t.Helper()
		c := New()
		_, err := c.LoadMetadata(context.Background(), meta)
		require.NoError(t, err)
		return c
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
		require.Empty(t, GenerateMigration(c1, cAltered, nil).SQL(), "CREATE PACKAGE target should converge after CREATE OR REPLACE migration")

		// Synced spec-only package (no CREATE prefix) equals its dumped CREATE OR REPLACE definition
		syncedSpec := loadMeta(t, schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Packages: []*metadatapb.PackageMetadata{{
				Name:       "PKG",
				Definition: "PACKAGE PKG AS PROCEDURE P; END PKG;",
			}},
		}))
		dumpedSpec, err := DatabaseDefinition(syncedSpec.Metadata())
		require.NoError(t, err)
		loadedSpec, err := LoadSQL(dumpedSpec)
		require.NoError(t, err)
		require.Empty(t, GenerateMigration(syncedSpec, loadedSpec, nil).SQL())

		dropPlan := GenerateMigration(c1, empty, nil)
		require.Contains(t, dropPlan.SQL(), `DROP PACKAGE "PKG";`)
		_, err = c1.Exec(dropPlan.SQL(), nil)
		require.NoError(t, err)
		require.Empty(t, c1.Metadata().Schemas[0].Packages)

		// Removing package body while keeping specification emits DROP PACKAGE BODY and converges
		cSpecAndBody, err := LoadSQL(src)
		require.NoError(t, err)
		cSpecOnly, err := LoadSQL(`
CREATE OR REPLACE PACKAGE PKG AS
  PROCEDURE P;
END PKG;
/
`)
		require.NoError(t, err)
		dropBodyPlan := GenerateMigration(cSpecAndBody, cSpecOnly, nil)
		require.Contains(t, dropBodyPlan.SQL(), `DROP PACKAGE BODY "PKG";`)
		require.NotContains(t, dropBodyPlan.SQL(), `CREATE OR REPLACE PACKAGE PKG`)
		_, err = cSpecAndBody.Exec(dropBodyPlan.SQL(), nil)
		require.NoError(t, err)
		require.Empty(t, GenerateMigration(cSpecAndBody, cSpecOnly, nil).SQL())

		// Removing package body AND changing specification drops body and replaces spec
		cSpecAndBody2, err := LoadSQL(src)
		require.NoError(t, err)
		cSpecChanged, err := LoadSQL(`
CREATE OR REPLACE PACKAGE PKG AS
  PROCEDURE P;
  PROCEDURE Q;
END PKG;
/
`)
		require.NoError(t, err)
		dropBodyAndAlterSpec := GenerateMigration(cSpecAndBody2, cSpecChanged, nil)
		require.Contains(t, dropBodyAndAlterSpec.SQL(), `DROP PACKAGE BODY "PKG";`)
		require.Contains(t, dropBodyAndAlterSpec.SQL(), `CREATE OR REPLACE PACKAGE PKG`)
		_, err = cSpecAndBody2.Exec(dropBodyAndAlterSpec.SQL(), nil)
		require.NoError(t, err)
		require.Empty(t, GenerateMigration(cSpecAndBody2, cSpecChanged, nil).SQL())
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

		// Synced APP snapshot compared against its unqualified LoadSQL dump produces an empty diff
		syncedSyns := loadMeta(t, meta)
		require.Empty(t, GenerateMigration(syncedSyns, c, nil).SQL())

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

		// Quoted column alias named "INVISIBLE" must not be marked invisible
		cQuotedInv, err := LoadSQL(`CREATE VIEW V_QUOTED_INV ("INVISIBLE", X) AS SELECT ID, NAME FROM T;`)
		require.NoError(t, err)
		require.Len(t, cQuotedInv.Metadata().Schemas[0].Views[0].Columns, 2)
		require.False(t, cQuotedInv.Metadata().Schemas[0].Views[0].Columns[0].GetIsInvisible())
		quotedInvDef, err := DatabaseDefinition(cQuotedInv.Metadata())
		require.NoError(t, err)
		require.Equal(t, "CREATE VIEW \"V_QUOTED_INV\" (\"INVISIBLE\", \"X\") AS SELECT ID, NAME FROM T;\n\n", quotedInvDef)

		// Quoted column alias named "INVISIBLE" that ALSO has the INVISIBLE modifier
		cQuotedInvHidden, err := LoadSQL(`CREATE VIEW V_QUOTED_INV_H ("INVISIBLE" INVISIBLE, X) AS SELECT ID, NAME FROM T;`)
		require.NoError(t, err)
		require.True(t, cQuotedInvHidden.Metadata().Schemas[0].Views[0].Columns[0].GetIsInvisible())
		quotedInvHiddenDef, err := DatabaseDefinition(cQuotedInvHidden.Metadata())
		require.NoError(t, err)
		require.Equal(t, "CREATE VIEW \"V_QUOTED_INV_H\" (\"INVISIBLE\" INVISIBLE, \"X\") AS SELECT ID, NAME FROM T;\n\n", quotedInvHiddenDef)
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

		// Quoted column names matching sequence option keywords (e.g. "CYCLE", "ORDER") must not set those options
		cQuotedCol, err := LoadSQL(`CREATE TABLE i_kw ("CYCLE" NUMBER GENERATED ALWAYS AS IDENTITY)`)
		require.NoError(t, err)
		require.Empty(t, cQuotedCol.Metadata().Schemas[0].Tables[0].Columns[0].GetOptions())
		defQuotedCol, err := DatabaseDefinition(cQuotedCol.Metadata())
		require.NoError(t, err)
		require.Contains(t, defQuotedCol, `"CYCLE" NUMBER GENERATED ALWAYS AS IDENTITY`)
		require.NotContains(t, defQuotedCol, `(CYCLE)`)
		require.NotContains(t, defQuotedCol, ` CYCLE`)
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

		// Composite RANGE-LIST definition survives LoadSQL round-trip and initial migration
		cComp, err := LoadSQL(def)
		require.NoError(t, err)
		defComp2, err := DatabaseDefinition(cComp.Metadata())
		require.NoError(t, err)
		require.Equal(t, def, defComp2)
		initCompMig := GenerateMigration(New(), cComp, nil).SQL()
		require.Contains(t, initCompMig, `PARTITION BY RANGE ("Y") SUBPARTITION BY LIST ("R")`)
		require.Contains(t, initCompMig, `SUBPARTITION "P2025_US" VALUES ('US')`)
		require.Contains(t, initCompMig, `SUBPARTITION "P2025_OTHER" VALUES (DEFAULT)`)

		// Composite RANGE-LIST with SUBPARTITION TEMPLATE across multiple partitions generates unique names
		cTemplate, err := LoadSQL(`
CREATE TABLE comp_tpl (
  y NUMBER,
  r VARCHAR2(10)
) PARTITION BY RANGE (y) SUBPARTITION BY LIST (r) SUBPARTITION TEMPLATE (
  SUBPARTITION sp_us VALUES ('US'),
  SUBPARTITION sp_other VALUES (DEFAULT)
) (
  PARTITION p2025 VALUES LESS THAN (2026),
  PARTITION p2026 VALUES LESS THAN (2027)
)`)
		require.NoError(t, err)
		defTpl, err := DatabaseDefinition(cTemplate.Metadata())
		require.NoError(t, err)
		require.Contains(t, defTpl, `PARTITION BY RANGE ("Y") SUBPARTITION BY LIST ("R")`)
		require.Contains(t, defTpl, `SUBPARTITION "P2025_SP_US" VALUES ('US')`)
		require.Contains(t, defTpl, `SUBPARTITION "P2025_SP_OTHER" VALUES (DEFAULT)`)
		require.Contains(t, defTpl, `SUBPARTITION "P2026_SP_US" VALUES ('US')`)
		require.Contains(t, defTpl, `SUBPARTITION "P2026_SP_OTHER" VALUES (DEFAULT)`)

		// Quoted partition and subpartition key column containing a comma
		cComma, err := LoadSQL(`CREATE TABLE p_comma ("A,B" NUMBER) PARTITION BY RANGE ("A,B") (PARTITION p1 VALUES LESS THAN (10))`)
		require.NoError(t, err)
		defComma, err := DatabaseDefinition(cComma.Metadata())
		require.NoError(t, err)
		require.Contains(t, defComma, `PARTITION BY RANGE ("A,B")`)
		cComma2, err := LoadSQL(defComma)
		require.NoError(t, err)
		defComma2, err := DatabaseDefinition(cComma2.Metadata())
		require.NoError(t, err)
		require.Equal(t, defComma, defComma2)

		defSubComma, err := DatabaseDefinition(schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Tables: []*metadatapb.TableMetadata{{
				Name: "COMP_COMMA",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "A,B", Type: "NUMBER", Nullable: true},
					{Name: "C,D", Type: "VARCHAR2(10 BYTE)", Nullable: true},
				},
				Partitions: []*metadatapb.TablePartitionMetadata{{
					Name:       "P1",
					Type:       metadatapb.TablePartitionMetadata_RANGE,
					Expression: `"A,B"`,
					Value:      "10",
					Subpartitions: []*metadatapb.TablePartitionMetadata{
						{Name: "P1_S1", Type: metadatapb.TablePartitionMetadata_LIST, Expression: `"C,D"`, Value: "'X'"},
					},
				}},
			}},
		}))
		require.NoError(t, err)
		require.Contains(t, defSubComma, `PARTITION BY RANGE ("A,B") SUBPARTITION BY LIST ("C,D")`)

		// Unnamed partitions and subpartitions omit zero-length "" identifiers
		cUnnamed, err := LoadSQL(`CREATE TABLE p_unnamed (id NUMBER) PARTITION BY RANGE (id) (PARTITION VALUES LESS THAN (10))`)
		require.NoError(t, err)
		defUnnamed, err := DatabaseDefinition(cUnnamed.Metadata())
		require.NoError(t, err)
		require.Contains(t, defUnnamed, `PARTITION VALUES LESS THAN (10)`)
		require.NotContains(t, defUnnamed, `PARTITION ""`)
		cUnnamed2, err := LoadSQL(defUnnamed)
		require.NoError(t, err)
		defUnnamed2, err := DatabaseDefinition(cUnnamed2.Metadata())
		require.NoError(t, err)
		require.Equal(t, defUnnamed, defUnnamed2)
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

	t.Run("synced spec-only package preserves spec when Exec adds body", func(t *testing.T) {
		c := loadMeta(t, schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Packages: []*metadatapb.PackageMetadata{{
				Name:       "PKG",
				Definition: "PACKAGE PKG AS PROCEDURE P; END PKG;",
			}},
		}))
		_, err := c.Exec(`
CREATE OR REPLACE PACKAGE BODY PKG AS
  PROCEDURE P IS
  BEGIN
    NULL;
  END P;
END PKG;
/`, nil)
		require.NoError(t, err)

		def, err := DatabaseDefinition(c.Metadata())
		require.NoError(t, err)
		require.Contains(t, def, "CREATE OR REPLACE PACKAGE PKG AS PROCEDURE P; END PKG;")
		require.Contains(t, def, "CREATE OR REPLACE PACKAGE BODY PKG AS")

		reloaded, err := LoadSQL(def)
		require.NoError(t, err)
		require.Len(t, reloaded.Metadata().Schemas[0].Packages, 1)
		reloadedDef, err := DatabaseDefinition(reloaded.Metadata())
		require.NoError(t, err)
		require.Equal(t, def, reloadedDef)

		// Synced body-only package preserves body when Exec adds spec
		cBody := loadMeta(t, schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Packages: []*metadatapb.PackageMetadata{{
				Name:       "PKG",
				Definition: "PACKAGE BODY PKG AS PROCEDURE P IS BEGIN NULL; END P; END PKG;",
			}},
		}))
		_, err = cBody.Exec(`CREATE OR REPLACE PACKAGE PKG AS PROCEDURE P; END PKG;`, nil)
		require.NoError(t, err)
		bodyDef, err := DatabaseDefinition(cBody.Metadata())
		require.NoError(t, err)
		require.Contains(t, bodyDef, "CREATE OR REPLACE PACKAGE PKG AS PROCEDURE P; END PKG;")
		require.Contains(t, bodyDef, "CREATE OR REPLACE PACKAGE BODY PKG AS PROCEDURE P IS BEGIN NULL; END P; END PKG;")
	})

	t.Run("database-link synonyms preserve explicit target schema matching local schema", func(t *testing.T) {
		meta := schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Synonyms: []*metadatapb.SynonymMetadata{
				{Name: "S_SAME_REMOTE", TargetSchema: "APP", TargetName: "T", TargetDbLink: "REMOTE_LINK"},
				{Name: "S_UNQUAL_REMOTE", TargetSchema: "", TargetName: "T", TargetDbLink: "REMOTE_LINK"},
			},
		})
		def, err := DatabaseDefinition(meta)
		require.NoError(t, err)
		require.Contains(t, def, `CREATE OR REPLACE SYNONYM "S_SAME_REMOTE" FOR "APP"."T"@REMOTE_LINK;`)
		require.Contains(t, def, `CREATE OR REPLACE SYNONYM "S_UNQUAL_REMOTE" FOR "T"@REMOTE_LINK;`)

		// Round-trip through LoadSQL keeps qualified and unqualified remote targets distinct
		c, err := LoadSQL(`
CREATE SYNONYM APP.S_SAME_REMOTE FOR APP.T@REMOTE_LINK;
CREATE SYNONYM APP.S_UNQUAL_REMOTE FOR T@REMOTE_LINK;
`)
		require.NoError(t, err)
		synByName := make(map[string]*metadatapb.SynonymMetadata)
		for _, syn := range c.Metadata().Schemas[0].Synonyms {
			synByName[syn.GetName()] = syn
		}
		require.Equal(t, "APP", synByName["S_SAME_REMOTE"].GetTargetSchema())
		require.Equal(t, "", synByName["S_UNQUAL_REMOTE"].GetTargetSchema())

		// Migration preserves APP.T@REMOTE_LINK and detects change between APP.T@REMOTE_LINK and T@REMOTE_LINK
		empty := loadMeta(t, schemaWith(&metadatapb.SchemaMetadata{Name: "APP"}))
		createPlan := GenerateMigration(empty, c, nil)
		require.Contains(t, createPlan.SQL(), `CREATE OR REPLACE SYNONYM "APP"."S_SAME_REMOTE" FOR "APP"."T"@REMOTE_LINK;`)
		require.Contains(t, createPlan.SQL(), `CREATE OR REPLACE SYNONYM "APP"."S_UNQUAL_REMOTE" FOR "T"@REMOTE_LINK;`)

		cChanged, err := LoadSQL(`
CREATE SYNONYM APP.S_SAME_REMOTE FOR T@REMOTE_LINK;
CREATE SYNONYM APP.S_UNQUAL_REMOTE FOR APP.T@REMOTE_LINK;
`)
		require.NoError(t, err)
		alterPlan := GenerateMigration(c, cChanged, nil)
		require.Contains(t, alterPlan.SQL(), `CREATE OR REPLACE SYNONYM "APP"."S_SAME_REMOTE" FOR "T"@REMOTE_LINK;`)
		require.Contains(t, alterPlan.SQL(), `CREATE OR REPLACE SYNONYM "APP"."S_UNQUAL_REMOTE" FOR "APP"."T"@REMOTE_LINK;`)
	})

	t.Run("view column diff distinguishes absent LoadSQL columns from alias or visibility changes", func(t *testing.T) {
		syncedTable := &metadatapb.TableMetadata{
			Name: "T",
			Columns: []*metadatapb.ColumnMetadata{
				{Name: "ID", Type: "NUMBER", Position: 1, Nullable: true},
				{Name: "NAME", Type: "VARCHAR2(20 BYTE)", Position: 2, Nullable: true},
			},
		}
		syncedStarView := &metadatapb.ViewMetadata{
			Name:       "V",
			Definition: "SELECT * FROM T",
			Columns: []*metadatapb.ColumnMetadata{
				{Name: "ID", Type: "NUMBER", Position: 1, Nullable: true},
				{Name: "NAME", Type: "VARCHAR2(20 BYTE)", Position: 2, Nullable: true},
			},
		}
		// Oracle stores a view's * expanded, and so does LoadSQL when the
		// view's table is in the schema.
		syncedExpandedView := &metadatapb.ViewMetadata{
			Name:       "V",
			Definition: `SELECT "ID","NAME" FROM T`,
			Columns:    syncedStarView.Columns,
		}
		fromWithTable := loadMeta(t, schemaWith(&metadatapb.SchemaMetadata{
			Name:   "APP",
			Tables: []*metadatapb.TableMetadata{syncedTable},
			Views:  []*metadatapb.ViewMetadata{syncedExpandedView},
		}))
		toWithTable, err := LoadSQL(`
CREATE TABLE APP.T (ID NUMBER, NAME VARCHAR2(20 BYTE));
CREATE VIEW APP.V AS SELECT * FROM T;
`)
		require.NoError(t, err)
		require.Empty(t, GenerateMigration(fromWithTable, toWithTable, nil).SQL(), "equivalent SELECT * view should not drop/recreate")

		// Even without T in either schema, synced view columns with types vs LoadSQL SELECT * or expression should not spuriously recreate
		fromExprOnly := loadMeta(t, schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Views: []*metadatapb.ViewMetadata{
				syncedStarView,
				{
					Name:       "V_EXPR",
					Definition: "SELECT ID + 1 FROM T",
					Columns: []*metadatapb.ColumnMetadata{
						{Name: "ID+1", Type: "NUMBER", Position: 1, Nullable: true},
					},
				},
			},
		}))
		toExprOnly, err := LoadSQL(`
CREATE VIEW APP.V AS SELECT * FROM T;
CREATE FORCE VIEW APP.V_EXPR AS SELECT ID + 1 FROM T;
`)
		require.NoError(t, err)
		require.Empty(t, GenerateMigration(fromExprOnly, toExprOnly, nil).SQL())

		// Removing explicit aliases on SELECT * FROM T (when T is in schema) MUST recreate the view
		fromAliasedStar := loadMeta(t, schemaWith(&metadatapb.SchemaMetadata{
			Name:   "APP",
			Tables: []*metadatapb.TableMetadata{syncedTable},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "V",
				Definition: `SELECT "ID","NAME" FROM T`,
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "RENAMED_ID", Type: "NUMBER", Position: 1, Nullable: true},
					{Name: "RENAMED_NAME", Type: "VARCHAR2(20 BYTE)", Position: 2, Nullable: true},
				},
			}},
		}))
		require.Contains(t, GenerateMigration(fromAliasedStar, toWithTable, nil).SQL(), `DROP VIEW "APP"."V";`)

		// Removing INVISIBLE on SELECT * FROM T MUST recreate the view
		fromInvisibleStar := loadMeta(t, schemaWith(&metadatapb.SchemaMetadata{
			Name:   "APP",
			Tables: []*metadatapb.TableMetadata{syncedTable},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "V",
				Definition: `SELECT "ID","NAME" FROM T`,
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "ID", Type: "NUMBER", Position: 1, Nullable: true, IsInvisible: true},
					{Name: "NAME", Type: "VARCHAR2(20 BYTE)", Position: 2, Nullable: true},
				},
			}},
		}))
		require.Contains(t, GenerateMigration(fromInvisibleStar, toWithTable, nil).SQL(), `DROP VIEW "APP"."V";`)

		// Adding explicit aliases in LoadSQL on SELECT * FROM T MUST recreate the view
		toAliasedStar, err := LoadSQL(`
CREATE TABLE APP.T (ID NUMBER, NAME VARCHAR2(20 BYTE));
CREATE VIEW APP.V (A, B) AS SELECT * FROM T;
`)
		require.NoError(t, err)
		addAliasSQL := GenerateMigration(toWithTable, toAliasedStar, nil).SQL()
		require.Contains(t, addAliasSQL, `DROP VIEW "APP"."V";`)
		require.Contains(t, addAliasSQL, `CREATE VIEW "APP"."V" ("A", "B") AS SELECT "ID","NAME" FROM T;`)
	})

	t.Run("column collation change detaches and restores indexes, constraints, and incoming foreign keys", func(t *testing.T) {
		// Regular index on VARCHAR2 column changing BINARY -> BINARY_CI
		fromIdx, err := LoadSQL(`
CREATE TABLE APP.T (
  ID NUMBER NOT NULL,
  CODE VARCHAR2(20 BYTE) COLLATE BINARY NOT NULL,
  CONSTRAINT T_CK CHECK (CODE <> '')
);
CREATE INDEX APP.T_CODE_IDX ON APP.T (CODE);
`)
		require.NoError(t, err)
		toIdx, err := LoadSQL(`
CREATE TABLE APP.T (
  ID NUMBER NOT NULL,
  CODE VARCHAR2(20 BYTE) COLLATE BINARY_CI NOT NULL,
  CONSTRAINT T_CK CHECK (CODE <> '')
);
CREATE INDEX APP.T_CODE_IDX ON APP.T (CODE);
`)
		require.NoError(t, err)
		sqlIdx := GenerateMigration(fromIdx, toIdx, nil).SQL()
		dropIdxPos := strings.Index(sqlIdx, `DROP INDEX "APP"."T_CODE_IDX";`)
		dropCkPos := strings.Index(sqlIdx, `DROP CONSTRAINT "T_CK";`)
		modifyPos := strings.Index(sqlIdx, `MODIFY ("CODE" VARCHAR2(20 BYTE) COLLATE BINARY_CI);`)
		createIdxPos := strings.Index(sqlIdx, `CREATE INDEX "APP"."T_CODE_IDX"`)
		createCkPos := strings.Index(sqlIdx, `ADD CONSTRAINT "T_CK"`)
		require.True(t, dropIdxPos >= 0 && dropCkPos >= 0 && modifyPos > dropIdxPos && modifyPos > dropCkPos, "must drop index and check before MODIFY COLLATE:\n%s", sqlIdx)
		require.True(t, createIdxPos > modifyPos && createCkPos > modifyPos, "must restore index and check after MODIFY COLLATE:\n%s", sqlIdx)
		_, err = fromIdx.Exec(sqlIdx, nil)
		require.NoError(t, err)
		require.Empty(t, GenerateMigration(fromIdx, toIdx, nil).SQL())

		// Primary key / unique constraint referenced by another table's FK
		fromFK, err := LoadSQL(`
CREATE TABLE APP.PARENT (
  CODE VARCHAR2(20 BYTE) COLLATE BINARY NOT NULL,
  CONSTRAINT PARENT_PK PRIMARY KEY (CODE)
);
CREATE TABLE APP.CHILD (
  ID NUMBER NOT NULL,
  PARENT_CODE VARCHAR2(20 BYTE) COLLATE BINARY NOT NULL,
  CONSTRAINT CHILD_FK FOREIGN KEY (PARENT_CODE) REFERENCES APP.PARENT (CODE)
);
`)
		require.NoError(t, err)
		toFK, err := LoadSQL(`
CREATE TABLE APP.PARENT (
  CODE VARCHAR2(20 BYTE) COLLATE BINARY_CI NOT NULL,
  CONSTRAINT PARENT_PK PRIMARY KEY (CODE)
);
CREATE TABLE APP.CHILD (
  ID NUMBER NOT NULL,
  PARENT_CODE VARCHAR2(20 BYTE) COLLATE BINARY_CI NOT NULL,
  CONSTRAINT CHILD_FK FOREIGN KEY (PARENT_CODE) REFERENCES APP.PARENT (CODE)
);
`)
		require.NoError(t, err)
		sqlFK := GenerateMigration(fromFK, toFK, nil).SQL()
		dropChildFKPos := strings.Index(sqlFK, `ALTER TABLE "APP"."CHILD" DROP CONSTRAINT "CHILD_FK";`)
		dropParentPKPos := strings.Index(sqlFK, `ALTER TABLE "APP"."PARENT" DROP CONSTRAINT "PARENT_PK";`)
		modifyParentPos := strings.Index(sqlFK, `ALTER TABLE "APP"."PARENT" MODIFY ("CODE" VARCHAR2(20 BYTE) COLLATE BINARY_CI);`)
		addParentPKPos := strings.Index(sqlFK, `ALTER TABLE "APP"."PARENT" ADD CONSTRAINT "PARENT_PK" PRIMARY KEY (CODE);`)
		addChildFKPos := strings.LastIndex(sqlFK, `ALTER TABLE "APP"."CHILD" ADD CONSTRAINT "CHILD_FK"`)
		require.True(t, dropChildFKPos >= 0 && dropParentPKPos > dropChildFKPos && modifyParentPos > dropParentPKPos, "must drop incoming FK and PK before MODIFY COLLATE:\n%s", sqlFK)
		require.True(t, addParentPKPos > modifyParentPos && addChildFKPos > addParentPKPos, "must restore PK and incoming FK after MODIFY COLLATE:\n%s", sqlFK)
		_, err = fromFK.Exec(sqlFK, nil)
		require.NoError(t, err)
		require.Empty(t, GenerateMigration(fromFK, toFK, nil).SQL())

		// Type + collation change together with a function-based index must only drop and recreate the index once
		fromBoth, err := LoadSQL(`
CREATE TABLE APP.T2 (
  CODE VARCHAR2(20 BYTE) COLLATE BINARY
);
CREATE INDEX APP.T2_UPPER_IDX ON APP.T2 (UPPER(CODE));
`)
		require.NoError(t, err)
		toBoth, err := LoadSQL(`
CREATE TABLE APP.T2 (
  CODE VARCHAR2(40 BYTE) COLLATE BINARY_CI
);
CREATE INDEX APP.T2_UPPER_IDX ON APP.T2 (UPPER(CODE));
`)
		require.NoError(t, err)
		sqlBoth := GenerateMigration(fromBoth, toBoth, nil).SQL()
		require.Equal(t, 1, strings.Count(sqlBoth, `DROP INDEX "APP"."T2_UPPER_IDX";`), sqlBoth)
		require.Equal(t, 1, strings.Count(sqlBoth, `CREATE INDEX "APP"."T2_UPPER_IDX"`), sqlBoth)

		// Converting a collated VARCHAR2 column to a non-character type (NUMBER) must not emit COLLATE
		fromChar, err := LoadSQL(`CREATE TABLE APP.T3 (C VARCHAR2(20 BYTE) COLLATE BINARY_CI);`)
		require.NoError(t, err)
		toNum, err := LoadSQL(`CREATE TABLE APP.T3 (C NUMBER);`)
		require.NoError(t, err)
		sqlToNum := GenerateMigration(fromChar, toNum, nil).SQL()
		require.Contains(t, sqlToNum, `ALTER TABLE "APP"."T3" MODIFY "C" NUMBER;`)
		require.NotContains(t, sqlToNum, `COLLATE`)
		_, err = fromChar.Exec(sqlToNum, nil)
		require.NoError(t, err)
		require.Empty(t, GenerateMigration(fromChar, toNum, nil).SQL())
		require.Empty(t, Diff(fromChar, toNum).TableChanges)
	})

	t.Run("packages referenced by views are created before dependent views", func(t *testing.T) {
		c, err := LoadSQL(`
CREATE OR REPLACE PACKAGE PKG AS
  FUNCTION F RETURN NUMBER;
END PKG;
/
CREATE OR REPLACE PACKAGE BODY PKG AS
  FUNCTION F RETURN NUMBER IS
  BEGIN
    RETURN 1;
  END F;
END PKG;
/
CREATE VIEW V AS SELECT PKG.F() AS N FROM DUAL;
`)
		require.NoError(t, err)

		def, err := DatabaseDefinition(c.Metadata())
		require.NoError(t, err)
		pkgPos := strings.Index(def, "CREATE OR REPLACE PACKAGE PKG")
		viewPos := strings.Index(def, `CREATE VIEW "V"`)
		require.True(t, pkgPos >= 0 && viewPos > pkgPos, "package must precede dependent view in definition:\n%s", def)

		mig := GenerateMigration(New(), c, nil).SQL()
		migPkgPos := strings.Index(mig, "CREATE OR REPLACE PACKAGE PKG")
		migViewPos := strings.Index(mig, `CREATE VIEW "V"`)
		require.True(t, migPkgPos >= 0 && migViewPos > migPkgPos, "package must precede dependent view in migration:\n%s", mig)
	})

	t.Run("synonyms referenced by views are created before dependent views", func(t *testing.T) {
		c, err := LoadSQL(`
CREATE TABLE T (A NUMBER);
CREATE SYNONYM S FOR T;
CREATE VIEW V AS SELECT A FROM S;
`)
		require.NoError(t, err)

		def, err := DatabaseDefinition(c.Metadata())
		require.NoError(t, err)
		tblPos := strings.Index(def, `CREATE TABLE "T"`)
		synPos := strings.Index(def, `CREATE OR REPLACE SYNONYM "S" FOR "T";`)
		viewPos := strings.Index(def, `CREATE VIEW "V"`)
		require.True(t, tblPos >= 0 && synPos > tblPos && viewPos > synPos, "table -> synonym -> view order required in definition:\n%s", def)

		mig := GenerateMigration(New(), c, nil).SQL()
		migTblPos := strings.Index(mig, `CREATE TABLE "T"`)
		migSynPos := strings.Index(mig, `CREATE OR REPLACE SYNONYM "S" FOR "T";`)
		migViewPos := strings.Index(mig, `CREATE VIEW "V"`)
		require.True(t, migTblPos >= 0 && migSynPos > migTblPos && migViewPos > migSynPos, "table -> synonym -> view order required in migration:\n%s", mig)

		// Synced snapshot where ALL_DEPENDENCIES omitted the synonym still orders synonym before view
		synced := loadMeta(t, schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Tables: []*metadatapb.TableMetadata{{
				Name:    "T",
				Columns: []*metadatapb.ColumnMetadata{{Name: "A", Type: "NUMBER", Position: 1, Nullable: true}},
			}},
			Synonyms: []*metadatapb.SynonymMetadata{{
				Name:         "S",
				TargetSchema: "APP",
				TargetName:   "T",
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "V",
				Definition: "SELECT A FROM S",
			}},
		}))
		syncedDef, err := DatabaseDefinition(synced.Metadata())
		require.NoError(t, err)
		require.True(t, strings.Index(syncedDef, `CREATE OR REPLACE SYNONYM "S" FOR "T";`) < strings.Index(syncedDef, `CREATE VIEW "V"`), syncedDef)

		// Transitive synonym chain S1 -> S2 -> T must order T -> S2 -> S1 -> V in both definition and migration
		cChain, err := LoadSQL(`
CREATE TABLE T (A NUMBER);
CREATE SYNONYM S2 FOR T;
CREATE SYNONYM S1 FOR S2;
CREATE VIEW V AS SELECT A FROM S1;
`)
		require.NoError(t, err)
		chainDef, err := DatabaseDefinition(cChain.Metadata())
		require.NoError(t, err)
		cTblPos := strings.Index(chainDef, `CREATE TABLE "T"`)
		cS2Pos := strings.Index(chainDef, `CREATE OR REPLACE SYNONYM "S2" FOR "T";`)
		cS1Pos := strings.Index(chainDef, `CREATE OR REPLACE SYNONYM "S1" FOR "S2";`)
		cVPos := strings.Index(chainDef, `CREATE VIEW "V"`)
		require.True(t, cTblPos >= 0 && cS2Pos > cTblPos && cS1Pos > cS2Pos && cVPos > cS1Pos, "T -> S2 -> S1 -> V required in definition:\n%s", chainDef)

		chainMig := GenerateMigration(New(), cChain, nil).SQL()
		mTblPos := strings.Index(chainMig, `CREATE TABLE "T"`)
		mS2Pos := strings.Index(chainMig, `CREATE OR REPLACE SYNONYM "S2" FOR "T";`)
		mS1Pos := strings.Index(chainMig, `CREATE OR REPLACE SYNONYM "S1" FOR "S2";`)
		mVPos := strings.Index(chainMig, `CREATE VIEW "V"`)
		require.True(t, mTblPos >= 0 && mS2Pos > mTblPos && mS1Pos > mS2Pos && mVPos > mS1Pos, "T -> S2 -> S1 -> V required in migration:\n%s", chainMig)

		// Cyclic synonym chain must terminate without infinite recursion
		cCycle, err := LoadSQL(`
CREATE SYNONYM S1 FOR S2;
CREATE SYNONYM S2 FOR S1;
CREATE VIEW V AS SELECT A FROM S1;
`)
		require.NoError(t, err)
		cycleDef, err := DatabaseDefinition(cCycle.Metadata())
		require.NoError(t, err)
		require.Contains(t, cycleDef, `CREATE OR REPLACE SYNONYM "S1" FOR "S2";`)
		require.Contains(t, cycleDef, `CREATE OR REPLACE SYNONYM "S2" FOR "S1";`)
	})
}
