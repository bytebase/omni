package driver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestDumpRestores syncs a corpus, dumps it, empties the schema, applies the
// dump, and requires the resynced database to equal the source.
func TestDumpRestores(t *testing.T) {
	corpus, err := os.ReadFile(filepath.Join("testdata", "dump", "corpus.sql"))
	if err != nil {
		t.Fatal(err)
	}
	d := openTestDriver(t)
	source, restored, _ := dumpAndRestore(t, d, string(corpus))
	if diff := restoreDiff(source, restored); diff != "" {
		t.Errorf("restored database differs from source (-source +restored):\n%s", diff)
	}
}

func dumpAndRestore(t *testing.T, d *Driver, setup string) (*metadatapb.DatabaseSchemaMetadata, *metadatapb.DatabaseSchemaMetadata, string) {
	t.Helper()
	ctx := context.Background()
	execCorpus(t, d.db, setup)
	source, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var dump strings.Builder
	if err := d.Dump(ctx, &dump, source); err != nil {
		t.Fatal(err)
	}
	execCorpus(t, d.db, dropScript(source))
	execCorpus(t, d.db, dump.String())
	restored, err := d.SyncDatabase(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return source, restored, dump.String()
}

func restoreDiff(source, restored *metadatapb.DatabaseSchemaMetadata) string {
	return cmp.Diff(source, restored, protocmp.Transform(),
		protocmp.IgnoreFields(&metadatapb.TableMetadata{}, "data_size", "index_size", "data_free", "row_count"),
	)
}

// dropScript empties every schema in snapshot so a dump can create it again.
func dropScript(snapshot *metadatapb.DatabaseSchemaMetadata) string {
	var lines []string
	schemas := slices.Clone(snapshot.Schemas)
	slices.SortFunc(schemas, func(a, b *metadatapb.SchemaMetadata) int { return strings.Compare(a.Name, b.Name) })
	for _, schema := range schemas {
		views := slices.Clone(schema.Views)
		slices.SortFunc(views, func(a, b *metadatapb.ViewMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, view := range views {
			lines = append(lines, fmt.Sprintf("DROP VIEW IF EXISTS [%s].[%s];", schema.Name, view.Name))
		}
		procedures := slices.Clone(schema.Procedures)
		slices.SortFunc(procedures, func(a, b *metadatapb.ProcedureMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, procedure := range procedures {
			lines = append(lines, fmt.Sprintf("DROP PROCEDURE IF EXISTS [%s].[%s];", schema.Name, procedure.Name))
		}
		functions := slices.Clone(schema.Functions)
		slices.SortFunc(functions, func(a, b *metadatapb.FunctionMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, function := range functions {
			lines = append(lines, fmt.Sprintf("DROP FUNCTION IF EXISTS [%s].[%s];", schema.Name, function.Name))
		}
		synonyms := slices.Clone(schema.Synonyms)
		slices.SortFunc(synonyms, func(a, b *metadatapb.SynonymMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, synonym := range synonyms {
			lines = append(lines, fmt.Sprintf("DROP SYNONYM IF EXISTS [%s].[%s];", schema.Name, synonym.Name))
		}
		sequences := slices.Clone(schema.Sequences)
		slices.SortFunc(sequences, func(a, b *metadatapb.SequenceMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, sequence := range sequences {
			lines = append(lines, fmt.Sprintf("DROP SEQUENCE IF EXISTS [%s].[%s];", schema.Name, sequence.Name))
		}
	}
	for _, schema := range schemas {
		for _, table := range schema.Tables {
			fks := slices.Clone(table.ForeignKeys)
			slices.SortFunc(fks, func(a, b *metadatapb.ForeignKeyMetadata) int { return strings.Compare(a.Name, b.Name) })
			for _, fk := range fks {
				lines = append(lines, fmt.Sprintf("ALTER TABLE [%s].[%s] DROP CONSTRAINT IF EXISTS [%s];", schema.Name, table.Name, fk.Name))
			}
		}
	}
	for _, schema := range schemas {
		tables := slices.Clone(schema.Tables)
		slices.SortFunc(tables, func(a, b *metadatapb.TableMetadata) int { return strings.Compare(a.Name, b.Name) })
		for _, table := range tables {
			lines = append(lines, fmt.Sprintf("DROP TABLE IF EXISTS [%s].[%s];", schema.Name, table.Name))
		}
	}
	for _, schema := range schemas {
		if !strings.EqualFold(schema.Name, "dbo") {
			lines = append(lines, fmt.Sprintf("DROP SCHEMA IF EXISTS [%s];", schema.Name))
		}
	}
	return strings.Join(lines, "\n")
}
