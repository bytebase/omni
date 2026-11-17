package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
	oracleparser "github.com/bytebase/omni/oracle/parser"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

func execScript(t *testing.T, conn *sql.Conn, script string) {
	t.Helper()
	ctx := context.Background()
	for _, statement := range oracleparser.Split(script) {
		if statement.Empty() {
			continue
		}
		if _, err := conn.ExecContext(ctx, statement.Text); err != nil {
			t.Fatalf("%s:\n%v", statement.Text, err)
		}
	}
}

// TestDefinitionRunsOnOracle executes each definition on the pinned server and
// reads back the objects the snapshot names. The tape replays it in CI; record
// with OMNI_RECORD=1 after changing a fixture.
func TestDefinitionRunsOnOracle(t *testing.T) {
	content, err := os.ReadFile("testdata/get_database_definition.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var tests []definitionCase
	if err := yaml.Unmarshal(content, &tests); err != nil {
		t.Fatal(err)
	}

	db := testcontainer.Open(t, "system")
	defer db.Close()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for i, tc := range tests {
		schema := fmt.Sprintf("OMNI_CAT_%d", i)
		execScript(t, conn, fmt.Sprintf(`CREATE USER "%s" IDENTIFIED BY "%s" QUOTA UNLIMITED ON USERS`, schema, testcontainer.Password))
		execScript(t, conn, fmt.Sprintf(`GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE PROCEDURE, CREATE TRIGGER, CREATE MATERIALIZED VIEW, CREATE TYPE TO "%s"`, schema))
		execScript(t, conn, fmt.Sprintf(`ALTER SESSION SET CURRENT_SCHEMA = "%s"`, schema))
		execScript(t, conn, tc.Expected)

		meta := &metadatapb.DatabaseSchemaMetadata{}
		if err := protojson.Unmarshal([]byte(tc.Metadata), meta); err != nil {
			t.Fatal(err)
		}
		want := expectedObjects(meta)
		got := readObjects(t, conn, schema)
		if !slices.Equal(got, want) {
			t.Errorf("%s: objects = %v, want %v", tc.Description, got, want)
		}
	}
}

// TestQuotedNamesSurviveOracle executes a definition with quoted,
// case-sensitive names and reads the names back as written.
func TestQuotedNamesSurviveOracle(t *testing.T) {
	db := testcontainer.Open(t, "system")
	defer db.Close()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	execScript(t, conn, fmt.Sprintf(`CREATE USER "OmniQuoted" IDENTIFIED BY "%s" QUOTA UNLIMITED ON USERS`, testcontainer.Password))
	execScript(t, conn, `GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE TO "OmniQuoted"`)
	execScript(t, conn, `ALTER SESSION SET CURRENT_SCHEMA = "OmniQuoted"`)

	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "OmniQuoted",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{
				Name: "Mixed Case",
				Columns: []*metadatapb.ColumnMetadata{
					{Name: "CamelCase", Type: "NUMBER", Position: 1},
					{Name: "lower", Type: "VARCHAR2(20 BYTE)", Position: 2},
				},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "Upper View",
				Definition: `SELECT "CamelCase" FROM "Mixed Case"`,
			}},
		}},
	}
	definition, err := DatabaseDefinition(meta)
	if err != nil {
		t.Fatal(err)
	}
	execScript(t, conn, definition)

	got := readObjects(t, conn, "OmniQuoted")
	want := expectedObjects(meta)
	if !slices.Equal(got, want) {
		t.Errorf("objects = %v, want %v", got, want)
	}
}

func readObjects(t *testing.T, conn *sql.Conn, schema string) []string {
	t.Helper()
	query := fmt.Sprintf(`SELECT object_type || '/' || object_name FROM all_objects
		WHERE owner = '%s' AND object_type IN ('TABLE', 'VIEW', 'MATERIALIZED VIEW', 'SEQUENCE', 'FUNCTION', 'PROCEDURE', 'INDEX', 'TRIGGER')
		ORDER BY object_type, object_name`, schema)
	rows, err := conn.QueryContext(context.Background(), query)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		if strings.HasPrefix(name, "SEQUENCE/ISEQ$$_") {
			continue
		}
		result = append(result, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	slices.Sort(result)
	return result
}

func expectedObjects(meta *metadatapb.DatabaseSchemaMetadata) []string {
	var result []string
	add := func(kind, name string) {
		if name != "" {
			result = append(result, strings.ToUpper(kind)+"/"+name)
		}
	}
	for _, schema := range meta.GetSchemas() {
		for _, table := range schema.GetTables() {
			if table.GetSkipDump() {
				continue
			}
			add("TABLE", table.GetName())
			for _, index := range table.GetIndexes() {
				add("INDEX", index.GetName())
			}
			for _, trigger := range table.GetTriggers() {
				add("TRIGGER", trigger.GetName())
			}
		}
		for _, view := range schema.GetViews() {
			add("VIEW", view.GetName())
			for _, trigger := range view.GetTriggers() {
				add("TRIGGER", trigger.GetName())
			}
		}
		for _, view := range schema.GetMaterializedViews() {
			add("MATERIALIZED VIEW", view.GetName())
			for _, index := range view.GetIndexes() {
				add("INDEX", index.GetName())
			}
			for _, trigger := range view.GetTriggers() {
				add("TRIGGER", trigger.GetName())
			}
		}
		for _, sequence := range schema.GetSequences() {
			if !strings.HasPrefix(sequence.GetName(), "ISEQ$$_") {
				add("SEQUENCE", sequence.GetName())
			}
		}
		for _, function := range schema.GetFunctions() {
			add("FUNCTION", function.GetName())
		}
		for _, procedure := range schema.GetProcedures() {
			add("PROCEDURE", procedure.GetName())
		}
	}
	slices.Sort(result)
	return slices.Compact(result)
}
