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
		// EDITIONING is not a metadata kind, and ALL_VIEWS.TEXT carries only
		// the query, so a synced editioning view renders as an ordinary view.
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

	t.Run("database link", func(t *testing.T) {
		// Oracle's ALL_DB_LINKS withholds link passwords (and rejects CONNECT TO
		// without IDENTIFIED BY), and LinkedDatabaseMetadata carries neither
		// the password nor whether the link is PUBLIC or private.
		meta := schemaWith(&metadatapb.SchemaMetadata{Name: "APP"})
		meta.LinkedDatabases = []*metadatapb.LinkedDatabaseMetadata{{
			Name:     "REMOTE",
			Username: "U",
			Host:     "example.com",
		}}
		def := definition(t, meta)
		if strings.Contains(def, "DATABASE LINK") {
			t.Fatalf("gap closed: render a DATABASE LINK\n%s", def)
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
