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

// TestCollationMigrationRunsOnOracle executes generated column-collation
// migrations (regular index, check constraint, PK/UK with incoming FK, and
// function-based index with type+collation change) against the pinned Oracle
// server and verifies the resulting column collations, indexes, and constraints.
func TestCollationMigrationRunsOnOracle(t *testing.T) {
	db := testcontainer.OpenExtended(t, "system")
	defer db.Close()
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	const schema = "OMNI_COLLATE"
	execScript(t, conn, fmt.Sprintf(`CREATE USER "%s" IDENTIFIED BY "%s" QUOTA UNLIMITED ON USERS`, schema, testcontainer.Password))
	execScript(t, conn, fmt.Sprintf(`GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE TO "%s"`, schema))
	execScript(t, conn, fmt.Sprintf(`ALTER SESSION SET CURRENT_SCHEMA = "%s"`, schema))

	fromSQL := `
CREATE TABLE OMNI_COLLATE.T (
  ID NUMBER NOT NULL,
  CODE VARCHAR2(20 BYTE) COLLATE BINARY NOT NULL,
  CONSTRAINT T_CK CHECK (CODE <> 'x')
);
CREATE INDEX OMNI_COLLATE.T_CODE_IDX ON OMNI_COLLATE.T (CODE);

CREATE TABLE OMNI_COLLATE.PARENT (
  CODE VARCHAR2(20 BYTE) COLLATE BINARY NOT NULL,
  CONSTRAINT PARENT_PK PRIMARY KEY (CODE)
);
CREATE TABLE OMNI_COLLATE.CHILD (
  ID NUMBER NOT NULL,
  PARENT_CODE VARCHAR2(20 BYTE) COLLATE BINARY NOT NULL,
  CONSTRAINT CHILD_FK FOREIGN KEY (PARENT_CODE) REFERENCES OMNI_COLLATE.PARENT (CODE)
);

CREATE TABLE OMNI_COLLATE.T2 (
  CODE VARCHAR2(20 BYTE) COLLATE BINARY
);
CREATE INDEX OMNI_COLLATE.T2_UPPER_IDX ON OMNI_COLLATE.T2 (UPPER(CODE));
`
	toSQL := `
CREATE TABLE OMNI_COLLATE.T (
  ID NUMBER NOT NULL,
  CODE VARCHAR2(20 BYTE) COLLATE BINARY_CI NOT NULL,
  CONSTRAINT T_CK CHECK (CODE <> 'x')
);
CREATE INDEX OMNI_COLLATE.T_CODE_IDX ON OMNI_COLLATE.T (CODE);

CREATE TABLE OMNI_COLLATE.PARENT (
  CODE VARCHAR2(20 BYTE) COLLATE BINARY_CI NOT NULL,
  CONSTRAINT PARENT_PK PRIMARY KEY (CODE)
);
CREATE TABLE OMNI_COLLATE.CHILD (
  ID NUMBER NOT NULL,
  PARENT_CODE VARCHAR2(20 BYTE) COLLATE BINARY_CI NOT NULL,
  CONSTRAINT CHILD_FK FOREIGN KEY (PARENT_CODE) REFERENCES OMNI_COLLATE.PARENT (CODE)
);

CREATE TABLE OMNI_COLLATE.T2 (
  CODE VARCHAR2(40 BYTE) COLLATE BINARY_CI
);
CREATE INDEX OMNI_COLLATE.T2_UPPER_IDX ON OMNI_COLLATE.T2 (UPPER(CODE));
`
	fromCat, err := LoadSQL(fromSQL)
	if err != nil {
		t.Fatal(err)
	}
	toCat, err := LoadSQL(toSQL)
	if err != nil {
		t.Fatal(err)
	}
	migration := GenerateMigration(fromCat, toCat, nil).SQL()
	if migration == "" {
		t.Fatal("expected non-empty migration SQL")
	}

	execScript(t, conn, fromSQL)
	execScript(t, conn, migration)

	rows, err := conn.QueryContext(ctx, `
		SELECT table_name, column_name, data_length, collation
		FROM all_tab_cols
		WHERE owner = 'OMNI_COLLATE'
		  AND column_name IN ('CODE', 'PARENT_CODE')
		ORDER BY table_name, column_name`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var gotCols []string
	for rows.Next() {
		var tbl, col, coll string
		var length int
		if err := rows.Scan(&tbl, &col, &length, &coll); err != nil {
			t.Fatal(err)
		}
		gotCols = append(gotCols, fmt.Sprintf("%s.%s:%d:%s", tbl, col, length, coll))
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	wantCols := []string{
		"CHILD.PARENT_CODE:20:BINARY_CI",
		"PARENT.CODE:20:BINARY_CI",
		"T.CODE:20:BINARY_CI",
		"T2.CODE:40:BINARY_CI",
	}
	if !slices.Equal(gotCols, wantCols) {
		t.Errorf("migrated columns = %v, want %v", gotCols, wantCols)
	}

	gotObjects := readObjects(t, conn, schema)
	wantObjects := expectedObjects(toCat.Metadata())
	if !slices.Equal(gotObjects, wantObjects) {
		t.Errorf("objects = %v, want %v", gotObjects, wantObjects)
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
