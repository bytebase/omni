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
	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			meta := &metadatapb.DatabaseSchemaMetadata{}
			if err := protojson.Unmarshal([]byte(tc.Metadata), meta); err != nil {
				t.Fatal(err)
			}
			got, err := DatabaseDefinition(meta)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.Expected {
				t.Errorf("definition mismatch\n--- got ---\n%s\n--- want ---\n%s", got, tc.Expected)
			}
			if _, err := metadataFromSQL(got); err != nil {
				t.Errorf("generated definition does not parse: %v", err)
			}
		})
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
