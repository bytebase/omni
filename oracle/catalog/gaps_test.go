package catalog

import (
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestDefinitionGaps pins what a definition cannot restore, one documented gap
// per subtest; docs/migration/oracle/catalog.md lists them. A gap that closes
// fails here, so the case moves into the fidelity tests beside this file.
func TestDefinitionGaps(t *testing.T) {
	definition := func(t *testing.T, meta *metadatapb.DatabaseSchemaMetadata) string {
		t.Helper()
		got, err := DatabaseDefinition(meta)
		if err != nil {
			t.Fatalf("DatabaseDefinition: %v", err)
		}
		return got
	}
	load := func(t *testing.T, sql string) *Catalog {
		t.Helper()
		c, err := LoadSQL(sql)
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		return c
	}
	schemaWith := func(schema *metadatapb.SchemaMetadata) *metadatapb.DatabaseSchemaMetadata {
		return &metadatapb.DatabaseSchemaMetadata{
			Name:    "APP",
			Schemas: []*metadatapb.SchemaMetadata{schema},
		}
	}

	t.Run("tablespace and storage", func(t *testing.T) {
		// Neither is in the metadata, so a definition states the table without
		// them and the restored table lands in the default tablespace.
		def := definition(t, load(t, `CREATE TABLE t (id NUMBER) TABLESPACE USERS STORAGE (INITIAL 64K)`).Metadata())
		if strings.Contains(def, "TABLESPACE") {
			t.Fatalf("gap closed: render the tablespace\n%s", def)
		}
		if strings.Contains(def, "STORAGE") {
			t.Fatalf("gap closed: render the storage clause\n%s", def)
		}
	})

	t.Run("partitioning", func(t *testing.T) {
		// Sync reads no partitioning and the catalog applies none, so a
		// partitioned table restores as a single heap table.
		c := load(t, `CREATE TABLE p (a NUMBER) PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS THAN (10))`)
		def := definition(t, c.Metadata())
		if strings.Contains(def, "PARTITION") {
			t.Fatalf("gap closed: render partitioning\n%s", def)
		}
		tables := c.Metadata().Schemas[0].Tables
		if len(tables) == 1 && len(tables[0].Partitions) > 0 {
			t.Fatal("gap closed: record a table's partitions")
		}
	})

	t.Run("external table", func(t *testing.T) {
		// ORGANIZATION EXTERNAL and its location are not in the metadata, so an
		// external table restores as an empty heap table.
		def := definition(t, load(t, `CREATE TABLE e (a NUMBER) ORGANIZATION EXTERNAL (TYPE ORACLE_LOADER DEFAULT DIRECTORY d ACCESS PARAMETERS (FIELDS TERMINATED BY ',') LOCATION ('f.csv'))`).Metadata())
		for _, clause := range []string{"ORGANIZATION", "LOCATION", "DIRECTORY"} {
			if strings.Contains(def, clause) {
				t.Fatalf("gap closed: render %s\n%s", clause, def)
			}
		}
	})

	t.Run("editioning view", func(t *testing.T) {
		// EDITIONING is not a metadata kind. The parser does not accept the
		// keyword either, and ALL_VIEWS.TEXT carries only the query, so a
		// synced editioning view renders as an ordinary view.
		def := definition(t, schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "EV",
				Definition: "SELECT ID FROM T",
			}},
		}))
		if strings.Contains(def, "EDITIONING") {
			t.Fatalf("gap closed: render an editioning view\n%s", def)
		}
		if !strings.Contains(def, `CREATE VIEW "EV"`) {
			t.Fatalf("an editioning view did not render as a view:\n%s", def)
		}
	})

	t.Run("object and XMLType table", func(t *testing.T) {
		// Sync marks these tables SkipDump because CREATE TABLE ... OF needs
		// the type, which is not a synced object.
		for _, name := range []string{"OT", "XT"} {
			def := definition(t, schemaWith(&metadatapb.SchemaMetadata{
				Name:   "APP",
				Tables: []*metadatapb.TableMetadata{{Name: name, SkipDump: true}},
			}))
			if strings.Contains(def, name) {
				t.Fatalf("gap closed: render the object table %s\n%s", name, def)
			}
		}
		// DDL that names the type outright is refused: installing an empty
		// table would state an object the type defines.
		_, err := LoadSQL(`CREATE TYPE ty AS OBJECT (a NUMBER); CREATE TABLE ot OF ty`)
		if err == nil || !strings.Contains(err.Error(), "object and XMLType tables are not supported") {
			t.Fatalf("CREATE TABLE OF ty: %v, want the unsupported error", err)
		}
		_, err = LoadSQL(`CREATE TABLE xt OF XMLTYPE`)
		if err == nil || !strings.Contains(err.Error(), "object and XMLType tables are not supported") {
			t.Fatalf("CREATE TABLE OF XMLTYPE: %v, want the unsupported error", err)
		}
	})

	t.Run("view column aliases", func(t *testing.T) {
		// ALL_VIEWS.TEXT does not carry the alias list, and writing the synced
		// column names can produce fewer names than output columns for a view
		// with an invisible column (ORA-01730).
		def := definition(t, schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Views: []*metadatapb.ViewMetadata{{
				Name:       "V",
				Definition: "SELECT ID, NAME FROM T",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "RENAMED_ID"},
					{Name: "RENAMED_NAME"},
				},
			}},
		}))
		if strings.Contains(def, `CREATE VIEW "V" (`) {
			t.Fatalf("gap closed: render the view's column aliases\n%s", def)
		}
	})

	t.Run("collation", func(t *testing.T) {
		// USING_NLS_COMP is the column default, and writing it back fails on a
		// session with a different string-size setting.
		def := definition(t, schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Tables: []*metadatapb.TableMetadata{{
				Name: "T",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "V", Type: "VARCHAR2(20 BYTE)", Collation: "BINARY_CI", Nullable: true},
				},
			}},
		}))
		if strings.Contains(def, "COLLATE") {
			t.Fatalf("gap closed: render the column's collation\n%s", def)
		}
	})

	t.Run("package", func(t *testing.T) {
		// A package is a synced object and the catalog keeps it, but a
		// definition does not render it.
		def := definition(t, schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Packages: []*metadatapb.PackageMetadata{{
				Name:       "PKG",
				Definition: "CREATE OR REPLACE PACKAGE PKG AS PROCEDURE P; END PKG;",
			}},
		}))
		if strings.Contains(def, "PACKAGE") {
			t.Fatalf("gap closed: render a package\n%s", def)
		}
	})

	t.Run("synonym and database link", func(t *testing.T) {
		// A synonym is a synced object and a database link is a database one;
		// neither renders, so a name a synonym supplies is missing after a
		// restore and a link-dependent query fails.
		meta := schemaWith(&metadatapb.SchemaMetadata{
			Name: "APP",
			Synonyms: []*metadatapb.SynonymMetadata{{
				Name:         "S",
				TargetSchema: "APP",
				TargetName:   "T",
			}},
		})
		meta.PublicSynonyms = []*metadatapb.SynonymMetadata{{
			Name:         "PUB_S",
			TargetSchema: "APP",
			TargetName:   "T",
		}}
		meta.LinkedDatabases = []*metadatapb.LinkedDatabaseMetadata{{
			Name:     "REMOTE",
			Username: "U",
			Host:     "example.com",
		}}
		def := definition(t, meta)
		for _, object := range []string{"SYNONYM", "DATABASE LINK"} {
			if strings.Contains(def, object) {
				t.Fatalf("gap closed: render a %s\n%s", object, def)
			}
		}
	})

	t.Run("identity sequence options", func(t *testing.T) {
		// The metadata carries an identity column's generation, seed and
		// increment; its backing sequence's cache, cycle and bounds have no
		// field, so a definition drops them.
		def := definition(t, load(t, `CREATE TABLE i (id NUMBER GENERATED ALWAYS AS IDENTITY (START WITH 5 INCREMENT BY 2 CACHE 100 CYCLE MINVALUE 5))`).Metadata())
		if !strings.Contains(def, "START WITH 5 INCREMENT BY 2") {
			t.Fatalf("the identity seed and increment were lost:\n%s", def)
		}
		for _, option := range []string{"CACHE", "CYCLE", "MINVALUE"} {
			if strings.Contains(def, option) {
				t.Fatalf("gap closed: render the identity's %s\n%s", option, def)
			}
		}
	})

	t.Run("named NOT NULL provenance", func(t *testing.T) {
		// The parser does not carry the constraint name through
		// ALTER TABLE ... MODIFY (CONSTRAINT nn NOT NULL), so the replay keeps
		// the nullability and loses the name.
		c := load(t, `CREATE TABLE n (id NUMBER NOT NULL); ALTER TABLE n MODIFY (id CONSTRAINT nn NOT NULL)`)
		def := definition(t, c.Metadata())
		if strings.Contains(def, "NN") {
			t.Fatalf("gap closed: keep the named NOT NULL constraint\n%s", def)
		}
		if !strings.Contains(def, `"ID" NUMBER NOT NULL`) {
			t.Fatalf("the nullability was lost:\n%s", def)
		}
	})
}
