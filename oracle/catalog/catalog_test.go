package catalog

import (
	"encoding/json"
	"flag"
	"os"
	"slices"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

type definitionCase struct {
	Description string `yaml:"description"`
	Metadata    string `yaml:"metadata"`
	Expected    string `yaml:"expected"`
}

func TestDatabaseDefinition(t *testing.T) {
	content, err := os.ReadFile("testdata/get_database_definition.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var tests []definitionCase
	if err := yaml.Unmarshal(content, &tests); err != nil {
		t.Fatal(err)
	}
	for i, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			meta := &metadatapb.DatabaseSchemaMetadata{}
			if err := protojson.Unmarshal([]byte(tc.Metadata), meta); err != nil {
				t.Fatal(err)
			}
			got, err := DatabaseDefinition(meta)
			if err != nil {
				t.Fatal(err)
			}
			if *update {
				tests[i].Expected = got
			} else if got != tc.Expected {
				t.Errorf("definition mismatch\n--- got ---\n%s\n--- want ---\n%s", got, tc.Expected)
			}
			if _, err := metadataFromSQL(got); err != nil {
				t.Errorf("generated definition does not parse: %v", err)
			}
		})
	}
	if *update {
		content, err := yaml.Marshal(tests)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile("testdata/get_database_definition.yaml", content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

var update = flag.Bool("update", false, "rewrite the catalog fixtures from this run")

type metadataCase struct {
	Description string `yaml:"description"`
	Schema      string `yaml:"schema"`
	Metadata    string `yaml:"metadata"`
}

func TestDefinitionRoundTrip(t *testing.T) {
	content, err := os.ReadFile("testdata/get_database_definition.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var tests []definitionCase
	if err := yaml.Unmarshal(content, &tests); err != nil {
		t.Fatal(err)
	}
	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			meta, err := metadataFromSQL(tc.Expected)
			if err != nil {
				t.Fatal(err)
			}
			got, err := DatabaseDefinition(meta)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.Expected {
				t.Errorf("round trip mismatch\n--- got ---\n%s\n--- want ---\n%s", got, tc.Expected)
			}
		})
	}
}

func TestMetadataFromSQL(t *testing.T) {
	content, err := os.ReadFile("testdata/get_database_metadata.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var tests []metadataCase
	if err := yaml.Unmarshal(content, &tests); err != nil {
		t.Fatal(err)
	}
	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			meta, err := metadataFromSQL(tc.Schema)
			if err != nil {
				t.Fatal(err)
			}
			encoded, err := json.MarshalIndent(meta, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			got := string(encoded) + "\n"
			if *update {
				for i := range tests {
					if tests[i].Description == tc.Description {
						tests[i].Metadata = got
					}
				}
				return
			}
			if got != tc.Metadata {
				t.Errorf("metadata mismatch\n--- got ---\n%s\n--- want ---\n%s", got, tc.Metadata)
			}
		})
	}
	if *update {
		content, err := yaml.Marshal(tests)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile("testdata/get_database_metadata.yaml", content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

type migrationCase struct {
	Description string `yaml:"description"`
	OldMetadata string `yaml:"oldMetadata"`
	NewMetadata string `yaml:"newMetadata"`
	Expected    string `yaml:"expected"`
}

func TestGenerateMigration(t *testing.T) {
	content, err := os.ReadFile("testdata/generate_migration.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var tests []migrationCase
	if err := yaml.Unmarshal(content, &tests); err != nil {
		t.Fatal(err)
	}
	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			oldMeta := parseMetadataFixture(t, tc.OldMetadata)
			newMeta := parseMetadataFixture(t, tc.NewMetadata)
			migration, err := generateMigration(Diff(&Catalog{meta: oldMeta}, &Catalog{meta: newMeta}))
			if err != nil {
				t.Fatal(err)
			}
			if migration != tc.Expected {
				t.Errorf("migration mismatch\n--- got ---\n%s\n--- want ---\n%s", migration, tc.Expected)
			}
		})
	}
}

func TestMigrationReplay(t *testing.T) {
	content, err := os.ReadFile("testdata/generate_migration.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var tests []migrationCase
	if err := yaml.Unmarshal(content, &tests); err != nil {
		t.Fatal(err)
	}
	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			from := &Catalog{meta: parseMetadataFixture(t, tc.OldMetadata)}
			to := parseMetadataFixture(t, tc.NewMetadata)
			plan := GenerateMigration(from, &Catalog{meta: to}, nil)
			results, err := from.Exec(plan.SQL(), nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, result := range results {
				if result.Error != nil {
					t.Fatalf("%s: %v", result.SQL, result.Error)
				}
			}
			requireSameObjects(t, from.Metadata(), to)
		})
	}
}

// requireSameObjects checks that the catalog holds every object the target
// snapshot names. Column order and positions are the server's to assign on an
// ADD, so they are not compared.
func requireSameObjects(t *testing.T, got, want *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	gotSchemas := schemaByName(got)
	for name, wantSchema := range schemaByName(want) {
		gotSchema := gotSchemas[name]
		if gotSchema == nil {
			t.Fatalf("schema %q is missing", name)
		}
		gotTables, wantTables := tableByName(gotSchema), tableByName(wantSchema)
		for tableName, wantTable := range wantTables {
			if wantTable.GetSkipDump() {
				continue
			}
			gotTable := gotTables[tableName]
			if gotTable == nil {
				t.Fatalf("table %q is missing", tableName)
			}
			gotColumns := map[string]bool{}
			for _, column := range gotTable.GetColumns() {
				gotColumns[column.GetName()] = true
			}
			for _, column := range wantTable.GetColumns() {
				if !gotColumns[column.GetName()] {
					t.Fatalf("column %q of table %q is missing", column.GetName(), tableName)
				}
			}
			requireNames(t, "index of "+tableName, indexNames(gotTable.GetIndexes()), indexNames(wantTable.GetIndexes()))
			requireNames(t, "foreign key of "+tableName, foreignKeyNames(gotTable.GetForeignKeys()), foreignKeyNames(wantTable.GetForeignKeys()))
			requireNames(t, "check of "+tableName, checkNames(gotTable), checkNames(wantTable))
		}
		requireNames(t, "view", keys(viewByName(gotSchema)), keys(viewByName(wantSchema)))
		requireNames(t, "materialized view", keys(materializedViewByName(gotSchema)), keys(materializedViewByName(wantSchema)))
		requireNames(t, "function", keys(functionByName(gotSchema)), keys(functionByName(wantSchema)))
		requireNames(t, "procedure", keys(procedureByName(gotSchema)), keys(procedureByName(wantSchema)))
		requireNames(t, "sequence", keys(sequenceByName(gotSchema)), keys(sequenceByName(wantSchema)))
	}
}

func requireNames(t *testing.T, kind string, got, want []string) {
	t.Helper()
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Fatalf("%s names = %v, want %v", kind, got, want)
	}
}

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for key := range m {
		out = append(out, key)
	}
	return out
}

func indexNames(indexes []*metadatapb.IndexMetadata) []string {
	out := make([]string, 0, len(indexes))
	for _, index := range indexes {
		out = append(out, index.GetName())
	}
	return out
}

func foreignKeyNames(foreignKeys []*metadatapb.ForeignKeyMetadata) []string {
	out := make([]string, 0, len(foreignKeys))
	for _, fk := range foreignKeys {
		out = append(out, fk.GetName())
	}
	return out
}

// checkNames leaves out named NOT NULL constraints: they are the column's
// NOT NULL in the server's dictionary, and Oracle does not carry the name
// through ALTER TABLE MODIFY.
func checkNames(table *metadatapb.TableMetadata) []string {
	out := make([]string, 0, len(table.GetCheckConstraints()))
	for _, check := range table.GetCheckConstraints() {
		if strings.HasSuffix(strings.TrimSpace(check.GetExpression()), "IS NOT NULL") {
			continue
		}
		out = append(out, check.GetName())
	}
	return out
}

func parseMetadataFixture(t *testing.T, encoded string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	meta := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal([]byte(encoded), meta); err != nil {
		t.Fatal(err)
	}
	return meta
}

func TestGetDatabaseMetadata(t *testing.T) {
	t.Run("empty and comment-only scripts fail with no parse results", func(t *testing.T) {
		for _, sql := range []string{"", "   \n\t", "-- comment only\n"} {
			if _, err := GetDatabaseMetadata(sql); err == nil {
				t.Fatalf("GetDatabaseMetadata(%q) = nil error, want error", sql)
			}
		}
	})

	t.Run("invalid syntax returns parse error", func(t *testing.T) {
		if _, err := GetDatabaseMetadata("CREATE TABLE ("); err == nil {
			t.Fatal("GetDatabaseMetadata(invalid) = nil error, want error")
		}
	})

	t.Run("slash-terminated PL/SQL units and schema authorization", func(t *testing.T) {
		meta, err := GetDatabaseMetadata(`
CREATE SCHEMA AUTHORIZATION HR
CREATE TABLE AUDIT_LOG (
    LOG_ID NUMBER PRIMARY KEY,
    STATUS VARCHAR2(20) DEFAULT ON NULL 'PENDING'
)
CREATE VIEW LOG_VIEW AS
SELECT LOG_ID, STATUS FROM AUDIT_LOG;

CREATE OR REPLACE TRIGGER audit_log_trigger
    BEFORE INSERT ON AUDIT_LOG
    FOR EACH ROW
BEGIN
    NULL;
END;
/

CREATE OR REPLACE PACKAGE financial_utils AS
    FUNCTION format_currency(p_amount NUMBER) RETURN VARCHAR2;
END financial_utils;
/

CREATE OR REPLACE PACKAGE BODY financial_utils AS
    FUNCTION format_currency(p_amount NUMBER) RETURN VARCHAR2 IS
    BEGIN
        RETURN TO_CHAR(p_amount, 'FM9999990.00');
    END format_currency;
END financial_utils;
/
`)
		if err != nil {
			t.Fatal(err)
		}
		if len(meta.Schemas) != 1 || meta.Schemas[0].Name != "HR" {
			t.Fatalf("schemas = %+v, want single schema HR", meta.Schemas)
		}
		schema := meta.Schemas[0]
		if len(schema.Tables) != 1 || schema.Tables[0].Name != "AUDIT_LOG" {
			t.Fatalf("tables = %+v, want AUDIT_LOG", schema.Tables)
		}
		if len(schema.Tables[0].Triggers) != 1 || schema.Tables[0].Triggers[0].Name != "AUDIT_LOG_TRIGGER" {
			t.Fatalf("triggers = %+v, want AUDIT_LOG_TRIGGER", schema.Tables[0].Triggers)
		}
		if len(schema.Views) != 1 || schema.Views[0].Name != "LOG_VIEW" {
			t.Fatalf("views = %+v, want LOG_VIEW", schema.Views)
		}
		if len(schema.Packages) != 1 || schema.Packages[0].Name != "FINANCIAL_UTILS" {
			t.Fatalf("packages = %+v, want FINANCIAL_UTILS", schema.Packages)
		}
	})

	t.Run("skips SQL*Plus commands and normalizes MV query and USING INDEX", func(t *testing.T) {
		meta, err := GetDatabaseMetadata(`
SET DEFINE OFF;
PROMPT Creating schema objects
REM Leading SQL*Plus comment
SPOOL out.log
WHENEVER SQLERROR EXIT SQL.SQLCODE
SET SERVEROUTPUT ON

CREATE TABLE "HR"."T" (
    "ID" NUMBER,
    "CODE" VARCHAR2(20)
);
CREATE UNIQUE INDEX "HR"."T_ID_UX" ON "HR"."T" ("ID");
CREATE INDEX "HR"."T_CODE_IX" ON "HR"."T" ("CODE");
ALTER TABLE "HR"."T" ADD CONSTRAINT "PK_T" PRIMARY KEY ("ID") USING INDEX "HR"."T_ID_UX";
ALTER TABLE "HR"."T" ADD CONSTRAINT "UK_T_CODE" UNIQUE ("CODE") USING INDEX "T_CODE_IX";

CREATE MATERIALIZED VIEW "HR"."MV_T"
    BUILD IMMEDIATE
    REFRESH COMPLETE ON DEMAND
    AS SELECT "ID", "CODE" FROM "HR"."T";
`)
		if err != nil {
			t.Fatal(err)
		}
		if len(meta.Schemas) != 1 || meta.Schemas[0].Name != "HR" {
			t.Fatalf("schemas = %+v, want HR", meta.Schemas)
		}
		schema := meta.Schemas[0]
		if len(schema.MaterializedViews) != 1 {
			t.Fatalf("mvs = %+v, want 1", schema.MaterializedViews)
		}
		mv := schema.MaterializedViews[0]
		if mv.Definition != "SELECT \"ID\", \"CODE\" FROM \"HR\".\"T\"\n" {
			t.Fatalf("mv.Definition = %q, want query-only SELECT with trailing newline", mv.Definition)
		}
		if len(schema.Tables) != 1 {
			t.Fatalf("tables = %+v, want 1", schema.Tables)
		}
		table := schema.Tables[0]
		if table.Columns[0].Nullable {
			t.Fatalf("ID column Nullable = true, want false after PK USING INDEX")
		}
		if len(table.Indexes) != 2 {
			t.Fatalf("indexes = %+v, want 2 reused indexes", table.Indexes)
		}
		idIdx := table.Indexes[0]
		if idIdx.Name != "T_ID_UX" || !idIdx.Primary || !idIdx.Unique || !idIdx.IsConstraint {
			t.Fatalf("T_ID_UX = %+v, want primary unique constraint index", idIdx)
		}
		if len(idIdx.Constraints) != 1 || idIdx.Constraints[0].Name != "PK_T" || !idIdx.Constraints[0].Primary {
			t.Fatalf("T_ID_UX.Constraints = %+v, want PK_T", idIdx.Constraints)
		}
		codeIdx := table.Indexes[1]
		if codeIdx.Name != "T_CODE_IX" || codeIdx.Primary || !codeIdx.IsConstraint {
			t.Fatalf("T_CODE_IX = %+v, want non-primary constraint index", codeIdx)
		}
		if len(codeIdx.Constraints) != 1 || codeIdx.Constraints[0].Name != "UK_T_CODE" || codeIdx.Constraints[0].Primary {
			t.Fatalf("T_CODE_IX.Constraints = %+v, want UK_T_CODE", codeIdx.Constraints)
		}
		// LoadCatalog/LoadSQL still preserves the BUILD/REFRESH prefix on MVs.
		cat, err := LoadSQL(`CREATE MATERIALIZED VIEW "HR"."MV_T" BUILD IMMEDIATE REFRESH COMPLETE ON DEMAND AS SELECT 1 FROM DUAL;`)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(cat.Metadata().Schemas[0].MaterializedViews[0].Definition, "BUILD IMMEDIATE") {
			t.Fatalf("LoadSQL mv.Definition = %q, want BUILD IMMEDIATE prefix preserved", cat.Metadata().Schemas[0].MaterializedViews[0].Definition)
		}
	})
}
